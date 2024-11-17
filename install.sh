#!/bin/bash

# 输出日志的函数
log() {
  echo -e "[INFO] $1"
}

error() {
  echo -e "[ERROR] $1" >&2
  exit 1
}

# 检测系统类型
OS=$(uname | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

# 映射架构
case "$ARCH" in
"x86_64")
  ARCH="amd64"
  ;;
"arm64" | "aarch64")
  ARCH="arm64"
  ;;
*)
  error "Unsupported architecture: $ARCH"
  ;;
esac

# GitHub 项目信息
GITHUB_REPO="gzwillyy/certum"
API_URL="https://api.github.com/repos/$GITHUB_REPO/releases/latest"

log "Fetching the latest release from $API_URL..."

# 获取最新版本号
LATEST_RELEASE=$(curl -s "$API_URL" | grep "tag_name" | cut -d '"' -f 4)
if [[ -z "$LATEST_RELEASE" ]]; then
  error "Failed to fetch the latest release information. Please check your network or the repository name."
fi

log "Latest version: $LATEST_RELEASE"

# 构造下载 URL
case "$OS" in
"linux")
  DOWNLOAD_URL="https://github.com/$GITHUB_REPO/releases/download/$LATEST_RELEASE/certum_validation_linux"
  ;;
"darwin")
  DOWNLOAD_URL="https://github.com/$GITHUB_REPO/releases/download/$LATEST_RELEASE/certum_validation_mac"
  ;;
"windows")
  DOWNLOAD_URL="https://github.com/$GITHUB_REPO/releases/download/$LATEST_RELEASE/certum_validation_windows.exe"
  ;;
*)
  error "Unsupported operating system: $OS"
  ;;
esac

# 设置目标文件名
TARGET_FILE="certum_validation"
if [[ "$OS" == "windows" ]]; then
  TARGET_FILE+=".exe"
fi

# 定义清理函数
cleanup() {
  log "Cleaning up temporary files..."
  rm -f "$TARGET_FILE"
  log "Cleanup complete. Exiting."
}

# 捕获退出信号
trap cleanup EXIT

# 检查是否已经下载
if [[ -f "$TARGET_FILE" ]]; then
  log "$TARGET_FILE already exists. Skipping download."
else
  # 下载文件
  log "Downloading from $DOWNLOAD_URL..."
  curl -L -o "$TARGET_FILE" "$DOWNLOAD_URL"
  
  # 检查下载是否成功
  if [[ $? -ne 0 ]]; then
    error "Failed to download $DOWNLOAD_URL. Please check your network or the repository name."
  fi

  log "Download complete: $TARGET_FILE"
  chmod +x "$TARGET_FILE"
fi

# 运行程序
log "Running $TARGET_FILE with default arguments..."
./$TARGET_FILE


