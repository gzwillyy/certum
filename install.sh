#!/bin/bash

# 输出日志的函数
log() {
  echo -e "[INFO] $1"
}

error() {
  echo -e "[ERROR] $1" >&2
  exit 1
}

# GitHub 项目信息
GITHUB_REPO="gzwillyy/certum"
API_URL="https://api.github.com/repos/$GITHUB_REPO/releases/latest"

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

log "Fetching the latest release from $API_URL..."
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

TARGET_FILE="certum_validation"
if [[ "$OS" == "windows" ]]; then
  TARGET_FILE+=".exe"
fi

log "Downloading from $DOWNLOAD_URL..."
curl -L -o "$TARGET_FILE" "$DOWNLOAD_URL"
chmod +x "$TARGET_FILE"

# 提示用户输入验证文件内容
read -p "请输入验证文件内容：" VALIDATION_CONTENT

# 检查输入内容是否为空
if [[ -z "$VALIDATION_CONTENT" ]]; then
  error "验证文件内容不能为空，请重新运行脚本。"
fi

# 运行程序
log "Running $TARGET_FILE with user-provided content..."
./$TARGET_FILE -content "$VALIDATION_CONTENT" -port 80

log "Cleaning up temporary files..."
rm -f "$TARGET_FILE"
log "Cleanup complete. Exiting."
