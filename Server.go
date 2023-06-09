package main

import (
	"fmt"
	"net"
)

func main() {
	// 创建一个监听器，在指定的地址和端口上监听TCP连接请求
	listener, err := net.Listen("tcp", "127.0.0.1:8000")
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		return
	}
	defer listener.Close()

	// 输出服务器已经开始监听
	fmt.Println("Server is listening on 127.0.0.1:8000")

	// 接收新的连接请求，并将连接请求放入通道中
	connections := make(chan net.Conn)
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				fmt.Println("Error accepting connection:", err.Error())
				continue
			}
			connections <- conn
		}
	}()

	// 处理连接请求（每个连接都启动一个goroutine）
	for {
		select {
		case conn := <-connections:
			go handleConnection(conn)
		}
	}
}

// 处理连接请求
func handleConnection(conn net.Conn) {
	defer conn.Close()

	// 创建一个 1024 字节的缓冲区，用于读取客户端发来的消息
	buf := make([]byte, 1024)

	// 读取客户端发来的消息
	_, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
		return
	}

	// 输出接收到的消息
	fmt.Printf("Received message: %s\n", string(buf))

	// 向客户端发送确认消息
	conn.Write([]byte("Message received."))

	// 等待客户端发送新的消息，并将消息放入通道中
	messages := make(chan []byte)
	go func() {
		for {
			n, err := conn.Read(buf)
			if err != nil {
				fmt.Println("Error reading:", err.Error())
				close(messages)
				return
			}
			data := make([]byte, n)
			copy(data, buf[:n])
			messages <- data
		}
	}()

	// 处理客户端发来的消息
	for {
		select {
		case message, ok := <-messages:
			if !ok {
				return
			}
			// 输出接收到的消息
			fmt.Printf("Received message: %s\n", string(message))

			// 回复客户端
			conn.Write([]byte("Message received."))
		}
	}
}
