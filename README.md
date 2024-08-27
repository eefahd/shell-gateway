# ShellGateway

ShellGateway is a minimal HTTP-based command execution service for local use. This tool allows you to execute command-line programs and shell scripts over HTTP requests. It's intended for internal use only, primarily to facilitate communication with local applications that only support frontend plugin development.

**Disclaimer**: This script is provided "as is" with no guarantees. Use it at your own risk.

## Design Philosophy
ShellGateway is designed to be a minimal, straightforward solution with minimal configuration requirements:
- **Simplified Script Management**: All scripts are expected to reside in the specified scripts directory.
- **Direct URL Mapping**: The request URL corresponds directly to the script or command name, eliminating the need for mapping dictionaries.
- **Simple Argument Passing**: All arguments are passed as a single string, maintaining simplicity.
- **No Unnecessary Complexity**: There are no plans to add advanced features or complications, preserving its lightweight nature.

## Installation

### Prerequisites

- Go (Golang) must be installed. You can download it from [here](https://go.dev/doc/install).
- Ensure your `GOPATH` and `GOROOT` environment variables are set correctly.

### Installation with `go install`

```sh
go install github.com/eefahd/shell-gateway@latest
```

## Environment Setup

Ensure that your Go environment variables are set:

- Add the Go bin directory to your `PATH`:

```sh
export PATH=$PATH:$(go env GOPATH)/bin
```

## Configuration

Create a `config.json` file in the `.config/shell-gateway/` directory under your home directory with the following content:

```json
{
  "port": "9090",
  "api_token": "YOUR_TOKEN",
  "scripts_directory": "YOUR_SCRIPTS_DIR_PATH",
  "allowed_commands": [
    "YOUR_COMMAND", // example: echo
    "YOUR_COMMAND2"
    //...
  ]
}
```

### Note
- **Allowed Commands**: A list of commands that are allowed to be executed. This setting is only needed to restrict command execution for more security, and does not apply to scripts.

## Running the Server

Start the server by running the binary:

```sh
shell-gateway
```

## Making Requests

To execute a command or script, send an HTTP POST request with the `Authorization` header and, if needed, a JSON body.

### Example: Executing a Command

```sh
curl -X POST "http://localhost:9090/echo" -H "Authorization: Bearer YOUR_TOKEN" -d '{"arguments": "Hello, World!"}'
```

### Example: Executing a Script

Assuming you have a script named `test.sh` in your scripts directory:

```sh
curl -X POST "http://localhost:9090/test" -H "Authorization: Bearer YOUR_TOKEN" -d '{"arguments": "arg1 arg2"}'
```

## License

This project is licensed under the MIT License. See the [LICENSE](https://opensource.org/license/mit) file for details.