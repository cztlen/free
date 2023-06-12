package server

import (
	"context"
	"fmt"
	"net"
	"time"
)

const timeOut = 300 * time.Second
const messageTimeOut = 200 * time.Second

func Server() {
	// 创建一个监听器，在指定的地址和端口上监听TCP连接请求
	listener, err := net.Listen("tcp", "127.0.0.1:8000")
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		return
	}
	defer listener.Close()

	// 输出服务器已经开始监听
	fmt.Println("Server is listening on 127.0.0.1:8000")

	ctx, cancel := context.WithTimeout(context.Background(), timeOut)
	defer cancel()

	// 接收新的连接请求，并将连接请求放入通道中
	connections := make(chan net.Conn)
	go func() {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("context timeOut 1")
				return
			default:
				conn, err := listener.Accept()
				if err != nil {
					fmt.Println("Error accepting connection:", err.Error())
					continue
				}
				connections <- conn
			}
		}
	}()

	// 处理连接请求（每个连接都启动一个goroutine）
	for {
		select {
		case <-ctx.Done():
			fmt.Println("context timeOut 2")
			return
		case conn := <-connections:
			go handleConnection(ctx, conn)
		}
	}
}

// 处理连接请求
//由半双工改为全双工通信，区别在于客户端一次连接后，多次进行通讯
//todo 后续考虑读写分离
func handleConnection(ctx context.Context, conn net.Conn) {
	defer conn.Close()
	var (
		n   int
		err error
		buf = make([]byte, 1024)
	)
	//select {
	//case <-ctx.Done():
	//	fmt.Println("context timeOut 3")
	//	return
	//default:
	//	n, err = conn.Read(buf)
	//	if err != nil {
	//		fmt.Printf("read message err:%s", err.Error())
	//		return
	//	}
	//}
	////将读取的消息打印
	//fmt.Printf("received message:%s", string(buf[:n]))

	//创建存储客户端消息的通道
	messages := make(chan []byte)
	//设置超时，若超过30秒，则退出
	ctx, cancel := context.WithTimeout(context.Background(), messageTimeOut)
	defer cancel()
	//将读取到的消息放入通道中
	go func() {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("context timeOut 4")
				close(messages)
				return
			default:
				n, err = conn.Read(buf)
				if err != nil {
					fmt.Printf("read message err1:%s", err.Error())
					close(messages)
					return
				}
				data := make([]byte, n)
				copy(data, buf[:n])
				messages <- data //放入通道中，阻塞
			}
		}

	}()
	//循环读取消息通道中消息
	for {
		select {
		case <-ctx.Done():
			fmt.Println("context timeOut 5")
			return
		case message, ok := <-messages:
			if !ok {
				return
			}
			fmt.Println(string(message))
			//回复客户端
			select {
			case <-ctx.Done():
				fmt.Println("context timeOut 6")
				return
			default:
				_, err = conn.Write([]byte("received message."))
				if err != nil {
					fmt.Printf("error sending response:%s", err.Error())
					return
				}
			}
		}
	}

}
