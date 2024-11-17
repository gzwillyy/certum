package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
)

const (
	logFilePath = "server.log"
)

var (
	server    net.Listener
	logFile   *os.File
	connCount int
	shutdown  = make(chan struct{}) // 用于通知服务器关闭
	wg        sync.WaitGroup        // 等待所有 goroutine 完成
)

// 初始化日志系统
func initLogging() error {
	var err error
	logFile, err = os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("failed to open log file: %v", err)
	}
	multiWriter := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(multiWriter)
	return nil
}

// 处理客户端连接
func handleConnection(conn net.Conn) {
	defer conn.Close()

	clientAddr := conn.RemoteAddr().String()
	log.Printf("New connection from %s", clientAddr)

	// 使用 bufio 读取客户端数据
	reader := bufio.NewReader(conn)
	_, err := reader.ReadString('\n') // 简单读取客户端请求
	if err != nil {
		log.Printf("Error reading from client %s: %v", clientAddr, err)
		return
	}

	// 构造HTML响应
	responseBody := `
<!DOCTYPE html>
<html>
<head>
    <title>Test Page</title>
</head>
<body>
    <h1>Welcome to TCP Server</h1>
    <p>This is a simple HTML test page served by TCP server.</p>
</body>
</html>
`
	response := fmt.Sprintf(
		"HTTP/1.1 200 OK\r\n"+
			"Content-Type: text/html\r\n"+
			"Content-Length: %d\r\n"+
			"Connection: close\r\n\r\n"+
			"%s",
		len(responseBody),
		responseBody,
	)

	// 发送响应到客户端
	_, err = conn.Write([]byte(response))
	if err != nil {
		log.Printf("Error sending response to client %s: %v", clientAddr, err)
	}
	log.Printf("HTML page sent to %s", clientAddr)
}

// 清理资源
func cleanup() {
	log.Println("Server shutting down...")

	// 关闭服务器
	if server != nil {
		server.Close()
		log.Println("Server socket closed.")
	}

	// 等待所有 goroutine 完成
	wg.Wait()

	// 删除日志文件
	if logFile != nil {
		logFile.Close()
		os.Remove(logFilePath)
		log.Println("Log file removed.")
	}
	log.Println("Cleanup complete.")
}

func main() {
	// 初始化日志系统
	if err := initLogging(); err != nil {
		fmt.Printf("Failed to initialize logging: %v\n", err)
		os.Exit(1)
	}

	// 手动输入端口
	var port string
	fmt.Print("Enter the port number to start the server: ")
	_, err := fmt.Scanln(&port)
	if err != nil {
		log.Fatalf("Invalid input: %v", err)
	}

	// 验证端口号有效性
	portInt, err := strconv.Atoi(port)
	if err != nil || portInt <= 0 || portInt > 65535 {
		log.Fatalf("Invalid port number. Please enter a number between 1 and 65535.")
	}

	// 捕获退出信号
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	// 启动服务器
	log.Printf("Starting TCP server on port %s...\n", port)
	server, err = net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	log.Printf("Server started on port %s", port)

	// 在后台监听信号
	go func() {
		<-signalChan
		close(shutdown)
		cleanup()
		os.Exit(0)
	}()

	// 主循环，接受连接
	for {
		conn, err := server.Accept()
		select {
		case <-shutdown:
			log.Println("Server is shutting down, stopping Accept loop.")
			return
		default:
		}

		if err != nil {
			if opErr, ok := err.(*net.OpError); ok && opErr.Err.Error() == "use of closed network connection" {
				log.Println("Server socket closed, stopping Accept loop.")
				return
			}
			log.Printf("Error accepting connection: %v", err)
			continue
		}

		connCount++
		log.Printf("Active connections: %d", connCount)

		wg.Add(1) // 增加计数
		go func() {
			defer wg.Done() // 减少计数
			handleConnection(conn)
			connCount--
			log.Printf("Connection closed. Active connections: %d", connCount)
		}()
	}
}
