package main

import (
	"bufio"
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/quic-go/quic-go/http3"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

// WebSocket 升级器，用于将 HTTP 升级为 WebSocket
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// QUIC 处理函数
func handleQUIC(certFile, keyFile string) {
	// 定义一个简单的处理器
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Unable to read request body", http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		// 打印请求体
		log.Printf("\n---------\n%s\n---------\n", string(body))
	})

	server := &http3.Server{
		Addr:    ":4433", // 监听端口，HTTP/3 通常在 443 端口上
		Handler: http.DefaultServeMux,
	}

	// 使用 TLS，配置证书和私钥文件
	log.Println("Starting QUIC server on :4433")
	if err := server.ListenAndServeTLS(certFile, keyFile); err != nil {
		log.Fatalf("Failed to start QUIC server: %v", err)
	}
}

// 从命令行读取输入
func inputCmd() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("error reading input: %v", err)
	}
	// 去除输入的换行符
	return strings.TrimSpace(input), nil
}

// WebSocket 处理函数
func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	// 将 HTTP 连接升级为 WebSocket 连接
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		http.Error(w, "WebSocket upgrade failed", http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	// 记录连接建立的日志
	log.Println("WebSocket connection established from:", r.RemoteAddr)
	for {
		ms, err := inputCmd()
		if err != nil {
			log.Println("Error reading command:", err)
			break
		}

		// 回应客户端
		if err := conn.WriteMessage(websocket.TextMessage, []byte(ms)); err != nil {
			log.Println("WebSocket write error:", err)
			break
		}
	}
}

// 启动 WebSocket 服务器
func startWebSocketServer(certFile, keyFile string) {
	http.HandleFunc("/ws", handleWebSocket)

	server := &http.Server{
		Addr:    ":8888",
		Handler: nil,
	}

	go func() {
		log.Println("Starting WebSocket server on wss://localhost:8888/ws")
		if err := server.ListenAndServeTLS(certFile, keyFile); err != nil && err != http.ErrServerClosed {
			log.Fatalf("WebSocket server error: %v", err)
		}
	}()

	// 捕获中断信号以优雅关闭服务器
	waitForShutdown(server)
}

// 捕获 OS 信号以优雅关闭服务器
func waitForShutdown(server *http.Server) {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	<-sig
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	log.Println("Server gracefully stopped")
}

func main() {
	// 替换为实际的证书路径
	certFile := "server.crt"
	keyFile := "server.key"

	// 启动 WebSocket 服务器
	go startWebSocketServer(certFile, keyFile)

	// 启动 QUIC 服务器
	go handleQUIC(certFile, keyFile)

	// 阻塞主 Goroutine
	select {}
}
