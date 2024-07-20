
# JRPC2 Client

This Go program implements a JSON-RPC 2.0 client that connects to a server, sends requests, and handles responses. It includes functionality for sending single requests, as well as performing a concurrency stress test.

## Features

- Connects to a JSON-RPC 2.0 server
- Sends individual string operations to the server
- Handles server notifications
- Performs concurrency stress testing

## Prerequisites

- Go 1.22.5 or later
- A running version https://github.com/karat-1/jsonrpcserver

## Installation

1. Clone the repository:

   ```sh
   git clone https://github.com/karat-1/jrpc2client.git
   cd jrpc2client
   ```

2. Install dependencies:

   ```sh
   go mod tidy
   ```

## Usage

1. Compile the program:

   ```sh
   go build -o jrpc2client
   ```

2. Run the program:

   ```sh
   ./jrpc2client -server localhost:8080
   ```

   Replace `localhost:8080` with the address of your JSON-RPC 2.0 server.

3. Interact with the client:

   - Enter a string to send to the server for counting characters.
   - Enter `-concurrency` to perform a concurrency stress test.
   - Enter `-exit` to exit the program.

## Code Overview

### Main Components

- **Client struct**: Manages the connection and context for each client.
- **newClient function**: Creates a new client and connects to the server.
- **countString function**: Sends a string to the server and receives the count of characters.
- **createRandomString function**: Generates a random string encoded in base64.
- **concurrencyStressTest function**: Performs a concurrency stress test by creating multiple clients and sending requests simultaneously.
- **clientSendMessage function**: Sends a message to the server and logs the response.

### Main Function

The `main` function initializes the client, parses command-line flags, and handles user input from the console.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgements

- [creachadair/jrpc2](https://github.com/creachadair/jrpc2) for the JSON-RPC 2.0 implementation
```
