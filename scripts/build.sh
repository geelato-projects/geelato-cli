#!/bin/bash
# scripts/build.sh - 构建脚本

set -euo pipefail

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

# 日志函数
log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 获取版本信息
VERSION="${VERSION:-dev}"
COMMIT="$(git rev-parse --short HEAD 2>/dev/null || echo 'unknown')"
DATE="$(date -u +"%Y-%m-%dT%H:%M:%SZ")"

# 构建标签
LDFLAGS="-X main.version=${VERSION} -X main.commit=${COMMIT} -X main.date=${DATE}"

log_info "开始构建 Geelato CLI ${VERSION}..."

# 确保依赖已下载
go mod download

# 清理旧构建
rm -f bin/geelato

# 构建
log_info "编译主程序..."
go build -ldflags "${LDFLAGS}" -o bin/geelato main.go

if [ -f bin/geelato ]; then
    log_info "构建成功: bin/geelato"
    ls -lh bin/geelato
else
    log_error "构建失败"
    exit 1
fi
