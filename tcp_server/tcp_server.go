package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"flag"
	"fmt"
	"log"
	"math/big"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

const (
	certFile = "server.crt"
	keyFile  = "server.key"
)

var (
	wg       sync.WaitGroup
	shutdown = make(chan struct{}) // 用于通知服务器关闭
)

// 自动生成自签名证书
func generateSelfSignedCert() error {
	log.Println("Generating self-signed certificate...")
	priv, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	if err != nil {
		return fmt.Errorf("failed to generate private key: %v", err)
	}

	template := x509.Certificate{
		SerialNumber: big.NewInt(time.Now().UnixNano()),
		Subject: pkix.Name{
			Organization: []string{"Self-Signed"},
			CommonName:   "localhost",
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(365 * 24 * time.Hour),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		return fmt.Errorf("failed to create certificate: %v", err)
	}

	certOut, err := os.Create(certFile)
	if err != nil {
		return fmt.Errorf("failed to create cert file: %v", err)
	}
	defer certOut.Close()
	if err := pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: certDER}); err != nil {
		return fmt.Errorf("failed to write certificate: %v", err)
	}

	keyOut, err := os.Create(keyFile)
	if err != nil {
		return fmt.Errorf("failed to create key file: %v", err)
	}
	defer keyOut.Close()
	privBytes, err := x509.MarshalECPrivateKey(priv)
	if err != nil {
		return fmt.Errorf("failed to marshal private key: %v", err)
	}
	if err := pem.Encode(keyOut, &pem.Block{Type: "EC PRIVATE KEY", Bytes: privBytes}); err != nil {
		return fmt.Errorf("failed to write private key: %v", err)
	}

	log.Println("Self-signed certificate generated successfully.")
	return nil
}

// 检查文件是否存在
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// 启动明文 TCP 服务
func startTCPServer(port int) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("Error starting TCP server on port %d: %v", port, err)
	}
	defer listener.Close()
	log.Printf("TCP server started on port %d", port)

	for {
		select {
		case <-shutdown:
			log.Println("TCP server shutting down...")
			return
		default:
		}

		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting TCP connection: %v", err)
			continue
		}

		wg.Add(1)
		go func(conn net.Conn) {
			defer wg.Done()
			handleConnection(conn, "TCP")
		}(conn)
	}
}

// 启动加密的 TLS 服务
func startTLSServer(port int) {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		log.Fatalf("Error loading TLS certificate: %v", err)
	}

	tlsConfig := &tls.Config{Certificates: []tls.Certificate{cert}}
	listener, err := tls.Listen("tcp", fmt.Sprintf(":%d", port), tlsConfig)
	if err != nil {
		log.Fatalf("Error starting TLS server on port %d: %v", port, err)
	}
	defer listener.Close()
	log.Printf("TLS server started on port %d", port)

	for {
		select {
		case <-shutdown:
			log.Println("TLS server shutting down...")
			return
		default:
		}

		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting TLS connection: %v", err)
			continue
		}

		wg.Add(1)
		go func(conn net.Conn) {
			defer wg.Done()
			handleConnection(conn, "TLS")
		}(conn)
	}
}

// 处理客户端连接
func handleConnection(conn net.Conn, protocol string) {
	defer conn.Close()

	clientAddr := conn.RemoteAddr().String()
	log.Printf("[%s] New connection from %s", protocol, clientAddr)

	responseBody := `
<!DOCTYPE html>
<html>
<head>
    <title>Test Page</title>
</head>
<body>
    <h1>Welcome to the %s server</h1>
    <p>This is a simple HTML test page served by the %s server.</p>
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
		fmt.Sprintf(responseBody, protocol, protocol),
	)

	_, err := conn.Write([]byte(response))
	if err != nil {
		log.Printf("[%s] Error sending response to client %s: %v", protocol, clientAddr, err)
	}
	log.Printf("[%s] HTML page sent to %s", protocol, clientAddr)
}

func main() {
	// 定义命令行参数
	tcpPort := flag.Int("tcp-port", 25125, "Port for the TCP server")
	tlsPort := flag.Int("tls-port", 25126, "Port for the TLS server")
	flag.Parse()

	// 检查证书是否存在，否则自动生成
	if !fileExists(certFile) || !fileExists(keyFile) {
		if err := generateSelfSignedCert(); err != nil {
			log.Fatalf("Failed to generate self-signed certificate: %v", err)
		}
	}

	// 捕获退出信号
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	// 启动 TCP 和 TLS 服务
	go startTCPServer(*tcpPort)
	go startTLSServer(*tlsPort)

	// 等待退出信号
	<-signalChan
	close(shutdown)
	wg.Wait()
	log.Println("Server gracefully shut down.")
}

// GOOS=linux GOARCH=amd64 go build -o tcp_server tcp_server.go
// chmod +x ./tcp_server
// nohup ./tcp_server --tcp-port 25125 --tls-port 25126 > server.log 2>&1 &
// ps aux | grep tcp_server
// kill 12345
