# 项目名称和版本
NAME := geelato
VERSION := 0.1.0
COMMIT := unknown
DATE := 
LDFLAGS := -X main.version= -X main.commit= -X main.date=

# Go 命令配置
GO := go
GOCMD := 
GOBIN := /bin
GOFLAGS := -mod=readonly

# 安装目录
PREFIX ?= /usr/local
BINDIR ?= /bin

# 颜色输出配置
RED := \033[31m
GREEN := \033[32m
YELLOW := \033[33m
BLUE := \033[34m
NC := \033[0m

# 默认目标
.PHONY: all
all: help

# 帮助信息
.PHONY: help
help:
@echo Geelato
