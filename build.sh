#!/bin/bash

# 设置输出目录
OUTPUT_DIR="dist"

# 定义源文件和目标平台架构
SOURCE_FILE="certum_validation.go"

# 构建 Linux 二进制文件
echo "Building for Linux (amd64)..."
GOOS=linux GOARCH=amd64 go build -o "$OUTPUT_DIR/certum_validation_linux" "$SOURCE_FILE"

# 构建 Windows 二进制文件
echo "Building for Windows (amd64)..."
GOOS=windows GOARCH=amd64 go build -o "$OUTPUT_DIR/certum_validation_windows.exe" "$SOURCE_FILE"

# 构建 macOS 二进制文件
echo "Building for macOS (amd64)..."
GOOS=darwin GOARCH=amd64 go build -o "$OUTPUT_DIR/certum_validation_mac" "$SOURCE_FILE"

echo "Build complete. Binaries are in the $OUTPUT_DIR directory."
