# Micro ⚡️

**Micro** is a command-line interface (CLI) tool built with Cobra that simplifies the process of running microservices. It provides a set of commands to streamline development tasks such as running the service, managing metrics, generating protobuf files, and more.

## Features

* **Easy to use:** Simple and intuitive commands for managing your microservice.
* **Efficient:** Quickly start and stop your services with minimal overhead.
* **Integrated tools:** Built-in support for common microservice development tools.
* **Configurable:** Customize behavior using environment variables and configuration files.

## Installation

1. **Prerequisites:** Ensure you have Go installed on your system.
2. **Clone the repository:** `git clone https://github.com/your-username/micro.git`
3. **Build the CLI:** `go build cmd/main.go`
4. **Move the binary:** `mv main micro` (or your preferred location)
5. **(Optional) Add to PATH:** For easy access, add the binary's location to your system's PATH environment variable.

## Usage

Micro provides a set of commands to manage your microservice. Here are some of the key commands:

* **`task run`:** Runs the microservice.
* **`task metrics`:** Starts `expvarmon` to monitor service metrics.
* **`task proto`:** Generates protobuf files for gRPC communication.
* **`task swagger`:** Launches Swagger UI for API documentation.
* **`task evans`:** Starts `evans` for gRPC interaction and debugging.
* **`task test`:** Runs tests and generates a coverage report.

