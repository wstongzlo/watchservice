package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

const (
	PORT     = ":10088"
	LOG_FILE = "server.log"
)

func logToFile(message string) {
	file, err := os.OpenFile(LOG_FILE, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening log file:", err)
		return
	}
	defer file.Close()
	logger := fmt.Sprintf("%s - %s\n", time.Now().Format("2006-01-02 15:04:05"), message)
	file.WriteString(logger)
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	for {
		data, err := reader.ReadString(']') // 读取完整的数据包
		if err != nil {
			logToFile(fmt.Sprintf("Connection closed or error: %v", err))
			return
		}

		logToFile("Received: " + data)
		fmt.Println("Received:", data)

		// 解析数据，提取设备 ID 和长度
		if strings.HasPrefix(data, "[") && strings.Contains(data, "*") {
			parts := strings.Split(data[1:len(data)-1], "*") // 去掉首尾的 [] 并分割
			if len(parts) >= 3 {
				deviceID := parts[1]
				length := "0002" // 固定响应长度
				response := fmt.Sprintf("[DW*%s*%s*KA]", deviceID, length)

				logToFile("Sending Response: " + response)
				fmt.Println("Sending Response:", response)
				conn.Write([]byte(response))
			}
		}
	}
}

func main() {
	listener, err := net.Listen("tcp", PORT)
	if err != nil {
		logToFile(fmt.Sprintf("Error starting server: %v", err))
		fmt.Println("Error starting server:", err)
		return
	}
	defer listener.Close()

	logToFile("TCP Server listening on port" + PORT)
	fmt.Println("TCP Server listening on port", PORT)

	for {
		conn, err := listener.Accept()
		if err != nil {
			logToFile(fmt.Sprintf("Error accepting connection: %v", err))
			fmt.Println("Error accepting connection:", err)
			continue
		}
		logToFile("New connection from " + conn.RemoteAddr().String())
		fmt.Println("New connection from", conn.RemoteAddr())
		go handleConnection(conn) // 启动 goroutine 处理连接
	}
}
