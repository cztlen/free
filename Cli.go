package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"time"
)

func main() {
	// 建立 TCP 连接，连接到指定的服务器地址和端口
	conn, err := net.Dial("tcp", "127.0.0.1:8000")
	if err != nil {
		fmt.Println("Error connecting:", err.Error())
		os.Exit(1)
	}
	defer conn.Close()

	// 启动一个 goroutine 用于读取服务器响应并输出
	messages := make(chan string)
	go func() {
		for {
			//response, err := bufio.NewReader(conn).ReadString('\n')
			buf := make([]byte, 1024)
			n, err := conn.Read(buf)
			if err != nil {
				fmt.Println("Error reading response:", err.Error())
				close(messages)
				return
			}
			response := string(buf[:n])
			messages <- response
		}
	}()

	// 启动一个 goroutine 用于从标准输入中读取用户输入并发送到服务器
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Enter message: ")
		text, _ := reader.ReadString('\n')
		fmt.Fprintf(conn, text)
		select {
		case message, ok := <-messages:
			if !ok {
				fmt.Println("Server disconnected.")
				return
			}
			fmt.Println("Response from server:", message)
		case <-time.After(30 * time.Second): // 等待1秒钟，如果没有服务器响应则继续下一次用户输入
			fmt.Println("No response from server.")
		}
	}
}
