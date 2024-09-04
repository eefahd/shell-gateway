package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/shlex"
)

type Config struct {
	Port             string   `json:"port"`
	APIToken         string   `json:"api_token"`
	ScriptsDirectory string   `json:"scripts_directory"`
	AllowedCommands  []string `json:"allowed_commands"`
}

var config Config

type Request struct {
	Arguments string `json:"arguments"`
}

const configFilePathInHome = ".config/shell-gateway/config.json"

func loadConfig() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("Error getting user home directory")
	}
	configFilePath := filepath.Join(homeDir, configFilePathInHome)

	file, err := os.ReadFile(configFilePath)
	if err != nil {
		log.Fatalf("Cannot open config file: %v", err)
	}

	if err := json.Unmarshal(file, &config); err != nil {
		log.Fatalf("Cannot get configuration from file: %v", err)
	}

	// Set default values if not present
	if config.Port == "" {
		config.Port = "9090"
	}
}

func validateToken(r *http.Request) bool {
	token := r.Header.Get("Authorization")
	return token == "Bearer "+config.APIToken
}

func authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !validateToken(r) {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func getScriptFilePath(scriptName string) (string, error) {
	files, err := os.ReadDir(config.ScriptsDirectory)
	if err != nil {
		return "", err
	}

	for _, file := range files {
		if !file.IsDir() && strings.TrimSuffix(file.Name(), filepath.Ext(file.Name())) == scriptName {
			return filepath.Join(config.ScriptsDirectory, file.Name()), nil
		}
	}
	return "", fmt.Errorf("script %s not found", scriptName)
}

func isValidCommand(command string) bool {
	for _, allowed := range config.AllowedCommands {
		if command == allowed {
			return true
		}
	}
	return false
}

func executeScript(scriptPath string, arguments string) (string, error) {
	args, err := shlex.Split(arguments)
	if err != nil {
		return "", err
	}
	args = append([]string{scriptPath}, args...)
	cmd := exec.Command("sh", args...)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

func executeCommand(command string, arguments string) (string, error) {
	args, err := shlex.Split(arguments)
	if err != nil {
		return "", err
	}
	cmd := exec.Command(command, args...)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

func handler(w http.ResponseWriter, r *http.Request) {
	currentTime := time.Now().Format(time.RFC3339)

	name := strings.TrimPrefix(r.URL.Path, "/")
	if name == "" {
		http.Error(w, "No command or script specified", http.StatusBadRequest)
		return
	}

	var req Request
	if r.Body != nil && r.ContentLength != 0 {
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
	}

	// Check if the name is a valid command
	if isValidCommand(name) {
		log.Printf("[%s] [command] %s %s %s\n", currentTime, r.Method, r.RequestURI, r.RemoteAddr)
		output, err := executeCommand(name, req.Arguments)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write([]byte(output))
		return
	}

	log.Printf("[%s] [script] %s %s %s\n", currentTime, r.Method, r.RequestURI, r.RemoteAddr)
	// If not a command, check if it's a script
	scriptPath, err := getScriptFilePath(name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	output, err := executeScript(scriptPath, req.Arguments)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte(output))
}

func main() {
	loadConfig()

	http.Handle("/", authenticate(http.HandlerFunc(handler)))

	port := config.Port
	fmt.Printf("Server listening on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
