package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
)

const (
	defaultValidationPath = "/.well-known/pki-validation/certum.txt" // 默认验证文件访问路径
	defaultValidationDir  = ".well-known/pki-validation"             // 验证文件存储目录
	defaultFileName       = "certum.txt"                             // 验证文件名
)

var server *http.Server // 用于优雅关闭服务器

func main() {
	// 定义命令行参数
	port := flag.Int("port", 80, "HTTP 服务监听端口")
	validationContent := flag.String("content", "", "验证文件的内容")
	flag.Parse()

	// 如果未提供验证内容，通过交互式输入
	if *validationContent == "" {
		fmt.Println("请输入验证文件内容：")
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		*validationContent = input
	}

	// 创建验证文件
	err := createValidationFile(*validationContent)
	if err != nil {
		fmt.Printf("创建验证文件失败: %v\n", err)
		return
	}

	// 启动 HTTP 服务
	go startServer(*port)

	// 捕获 Ctrl+C 信号以安全结束服务并清理文件
	waitForExitSignal()
	cleanup()
	fmt.Println("服务已安全关闭并清理完成。")
}

// 创建验证文件
func createValidationFile(content string) error {
	// 确保目录存在
	dirPath := filepath.Join(".", defaultValidationDir)
	err := os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		return fmt.Errorf("无法创建目录: %v", err)
	}

	// 写入验证文件内容
	filePath := filepath.Join(dirPath, defaultFileName)
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("无法创建验证文件: %v", err)
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		return fmt.Errorf("无法写入验证文件: %v", err)
	}

	fmt.Printf("验证文件已创建：%s\n", filePath)
	return nil
}

// 提供验证文件的 HTTP 处理器
func serveValidationFile(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == defaultValidationPath {
		filePath := filepath.Join(".", defaultValidationDir, defaultFileName)
		file, err := os.Open(filePath)
		if err != nil {
			http.Error(w, "无法访问验证文件", http.StatusInternalServerError)
			return
		}
		defer file.Close()
		io.Copy(w, file)
	} else {
		http.NotFound(w, r)
	}
}

// 启动 HTTP 服务
func startServer(port int) {
	mux := http.NewServeMux()
	mux.HandleFunc(defaultValidationPath, serveValidationFile)
	address := fmt.Sprintf("0.0.0.0:%d", port)
	server = &http.Server{Addr: address, Handler: mux}

	fmt.Printf("验证服务已启动：http://%s%s\n", address, defaultValidationPath)
	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		fmt.Printf("服务器启动失败: %v\n", err)
	}
}

// 捕获退出信号
func waitForExitSignal() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan // 等待信号
}

// 清理文件和目录
func cleanup() {
	fmt.Println("正在清理临时文件...")
	dirPath := filepath.Join(".", defaultValidationDir)
	err := os.RemoveAll(dirPath)
	if err != nil {
		fmt.Printf("清理临时文件失败: %v\n", err)
	} else {
		fmt.Println("临时文件已清理。")
	}

	// 优雅关闭 HTTP 服务
	if server != nil {
		err := server.Close()
		if err != nil {
			fmt.Printf("关闭服务器失败: %v\n", err)
		} else {
			fmt.Println("服务器已安全关闭。")
		}
	}
}

// GOARCH=amd64 GOOS=linux go build -o certum_validation_linux certum_validation.go

// GOARCH=amd64 GOOS=windows go build -o certum_validation_windows.exe certum_validation.go

// GOARCH=amd64 GOOS=darwin go build -o certum_validation_darwin certum_validation.go
