#!/bin/bash
# scripts/release.sh - 发布脚本

set -euo pipefail

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 检查版本标签
if [ -z "${TAG:-}" ]; then
    log_error "请设置 TAG 环境变量指定发布版本"
    exit 1
fi

log_info "准备发布版本 ${TAG}..."

# 运行测试
log_info "运行测试..."
go test ./... || exit 1

# 构建所有平台
log_info "构建所有平台版本..."
make build-all || exit 1

# 创建发布包
log_info "创建发布包..."
mkdir -p release/${TAG}

# 复制构建产物
cp bin/geelato-* release/${TAG}/

# 创建校验和文件
cd release/${TAG}
sha256sum geelato-* > SHA256SUMS
cd ../..

log_info "发布准备完成"
log_info "发布文件位于: release/${TAG}/"
ls -lh release/${TAG}/
