# MCP Configuration Guide

## 1. Overview

MCP (Model Context Protocol) provides access to cloud platform real-time capabilities.

## 2. Configuration Steps

# 1. Configure API Key
echo "GEELATO_API_KEY=your_api_key" > .env

# 2. Initialize MCP
geelato mcp init

# 3. Sync capabilities
geelato mcp sync
