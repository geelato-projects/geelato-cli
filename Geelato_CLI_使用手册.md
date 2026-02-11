# Geelato CLI 使用操作说明手册

## 一、产品概述

Geelato CLI 是 Geelato 低代码平台的命令行工具，为开发者提供从项目初始化到云端部署的完整开发工作流支持。通过 Geelato CLI，开发者可以高效地管理数据模型定义、创建 API 接口、设计业务流程、处理云端同步等核心开发任务。本手册基于 Geelato CLI 的现有实现版本编写，全面介绍各命令的使用方法、参数配置以及最佳实践指南。

Geelato CLI 的设计理念是让低代码开发变得简单而强大。传统的低代码平台往往在易用性和灵活性之间存在权衡，而 Geelato CLI 通过声明式的配置文件和脚本化的业务逻辑实现了二者的平衡。开发者既可以使用预定义的模型规范快速搭建应用结构，也可以通过自定义的 API 脚本实现复杂的业务需求。所有这些能力都可以通过统一的命令行界面访问，无需复杂的图形界面操作，非常适合现代 DevOps 开发流程。

## 二、安装与环境配置

### 2.1 系统要求与前置条件

在安装 Geelato CLI 之前，请确保您的开发环境满足以下基本要求。操作系统方面，Geelato CLI 支持 Windows、macOS 和 Linux 三大主流平台，其中 Windows 系统建议使用 Windows 10 及以上版本以获得最佳兼容性。运行时环境方面，Geelato CLI 采用 Go 语言开发，需要 Go 1.21 或更高版本运行时支持。此外，为了使用云端同步功能，您还需要配置网络访问权限，确保 CLI 工具能够与 Geelato 云端平台正常通信。

版本管理工具 Git 是推荐安装的辅助工具，虽然并非必需，但在团队协作和代码版本管理场景中非常有用。Git 不仅可以用于管理您自己的项目代码，还可以方便地回滚修改、创建分支进行实验性开发。如果您的项目需要与远程代码仓库集成，Git 更是不可或缺的基础工具。建议安装 Git 2.0 或更高版本以获得完整的功能支持。

### 2.2 安装步骤详解

Geelato CLI 的安装过程非常简单，只需几个步骤即可完成。首先，从官方仓库下载适用于您操作系统的安装包。Windows 用户可以下载 `.exe` 安装程序或压缩包，macOS 用户可以使用 Homebrew 包管理器或下载 `.pkg` 安装包，Linux 用户则可以通过包管理器或直接下载二进制文件来完成安装。

```powershell
# Windows 安装示例
# 方式一：使用压缩包
Invoke-WebRequest -Uri "https://geelato.com/cli/geelato-windows-amd64.zip" -OutFile "geelato.zip"
Expand-Archive -Path "geelato.zip" -DestinationPath "C:\Tools\geelato"
$env:PATH += ";C:\Tools\geelato"

# 方式二：使用 scoop（推荐）
scoop install geelato

# 方式三：使用 winget
winget install Geelato.GeelatoCLI
```

```bash
# macOS 安装示例
# 方式一：使用 Homebrew（推荐）
brew install geelato/geelato-cli/geelato

# 方式二：直接下载
curl -L "https://geelato.com/cli/geelato-darwin-amd64.tar.gz" -o "geelato.tar.gz"
tar -xzf geelato.tar.gz
sudo mv geelato /usr/local/bin/geelato

# 方式三：使用 Cask
brew install --cask geelato
```

```bash
# Linux 安装示例
# 方式一：使用 APT（Debian/Ubuntu）
curl -L "https://geelato.com/cli/geelato-linux-amd64.deb" -o "geelato.deb"
sudo dpkg -i geelato.deb

# 方式二：使用 RPM（RHEL/CentOS）
curl -L "https://geelato.com/cli/geelato-linux-amd64.rpm" -o "geelato.rpm"
sudo rpm -ivh geelato.rpm

# 方式三：直接下载二进制文件
curl -L "https://geelato.com/cli/geelato-linux-amd64.tar.gz" -o "geelato.tar.gz"
tar -xzf geelato.tar.gz
sudo mv geelato /usr/local/bin/geelato
```

安装完成后，打开新的终端窗口并执行版本验证命令确认安装成功。成功的安装会显示当前 CLI 工具的版本号、构建信息以及可用的命令列表。如果遇到命令未找到的错误，请检查系统的 PATH 环境变量配置是否包含 Geelato CLI 的安装目录。

```bash
# 验证安装
geelato version

# 查看可用命令
geelato --help

# 查看特定命令的帮助
geelato init --help
geelato model --help
```

### 2.3 全局配置初始化

首次使用 Geelato CLI 时，建议先完成全局配置的初始化。全局配置存储在用户主目录下的 `.geelato` 文件夹中，包含 API 服务器地址、认证信息、默认编辑器等基础设置。合理的全局配置可以大幅减少日常使用时的重复输入，提升开发效率。

```bash
# 初始化全局配置
geelato config init

# 设置云端平台地址
geelato config set api.url "https://api.geelato.com"

# 设置默认编辑器
geelato config set editor.command "code"
geelato config set editor.args "--wait"

# 查看当前配置
geelato config list

# 查看特定配置项
geelato config get api.url
```

全局配置文件采用 YAML 格式存储，您可以手动编辑配置文件以完成高级配置。配置文件的默认路径为 `~/.geelato/config.yaml`（Linux/macOS）或 `%USERPROFILE%\.geelato\config.yaml`（Windows）。以下是典型配置文件的结构示例：

```yaml
# ~/.geelato/config.yaml
api:
  url: "https://api.geelato.com"
  timeout: 30
  retryCount: 3
  retryDelay: 1

sync:
  autoPush: false
  autoPull: false
  branch: "main"
  message: "Auto sync from Geelato CLI"

logging:
  level: "info"
  format: "text"
  file: "~/.geelato/logs/geelato.log"

editor:
  command: "code"
  args: "--wait"

git:
  username: "your-name"
  email: "your-email@example.com"
  remote: "origin"

cache:
  dir: "~/.geelato/cache"
  ttl: 3600
  maxSize: 104857600
```

## 三、快速开始

### 3.1 创建第一个应用

让我们从一个完整的示例开始，快速体验 Geelato CLI 的核心功能。在这个示例中，我们将创建一个包含用户管理功能的完整应用，包括数据模型定义、API 接口开发以及云端同步配置。整个过程大约需要 15 分钟，完成后您将掌握 Geelato CLI 的基本使用方法。

首先，选择一个合适的目录作为您的工作空间，然后执行初始化命令创建新应用。Geelato CLI 会引导您完成应用的基本信息设置，包括应用名称、描述、模板选择等。如果您已有明确的项目规划，也可以通过命令行参数一次性完成所有配置。

```bash
# 切换到工作目录
cd ~/projects

# 交互式创建应用
geelato init my-first-app

# 或使用命令行参数一次性完成配置
geelato init my-first-app \
  --desc "我的第一个 Geelato 应用" \
  --template "default" \
  --skip-git

# 进入应用目录
cd my-first-app
```

应用创建完成后，您会看到标准化的项目目录结构。Geelato CLI 会自动生成应用配置文件、模型定义目录、API 脚本目录等核心组件。目录结构遵循统一的规范，便于团队成员理解和维护。

```
my-first-app/
├── geelato.json              # 应用配置文件
├── .geelato/
│   └── sync-state.json       # 同步状态记录
├── api/                      # API 脚本目录
│   └── .gitkeep             # 目录占位符
├── meta/                     # 模型定义目录
│   ├── .gitkeep
│   └── User/                # 用户实体
│       ├── User.define.json  # 表定义
│       ├── User.columns.json # 字段定义
│       └── User.fk.json     # 外键定义
├── page/                     # 页面定义目录
│   └── .gitkeep
├── workflow/                 # 工作流定义目录
│   └── .gitkeep
└── doc/                      # 文档目录（可选）
```

### 3.2 验证应用结构

应用创建成功后，可以使用 `validate` 命令验证当前目录的应用结构是否完整有效。这个命令会检查必要的配置文件和目录结构是否存在，确保应用可以正常进行开发和同步操作。

```bash
# 验证当前目录的应用结构
geelato validate

# 输出示例：
# [INFO] Validating application...
# [INFO] Validation complete!
# [INFO] Found 2 models, 3 APIs, 1 workflows
# [SUCCESS] Application structure is valid!
```

如果应用结构不完整，验证命令会列出缺失的文件和目录，并给出修复建议。

### 3.3 查看应用基本信息

应用创建成功后，可以使用 `config list` 命令查看当前应用的基本信息和配置状态。这个命令会读取配置文件中的元数据，并以友好的格式展示出来。通过这个命令，您可以快速确认当前应用的身份标识、版本信息以及创建时间等关键属性。

```bash
# 查看当前应用信息
geelato config list

# 输出示例：
# 应用信息
# =========
#
# 名称:      my-first-app
# 应用ID:   app_my_first_app_001
# 描述:     我的第一个 Geelato 应用
# 版本:     1.0.0
# 创建时间: 2024-01-15T10:00:00Z
#
# 路径:     /Users/yourname/projects/my-first-app
```

## 四、克隆应用命令详解

### 4.1 clone 命令概述

`clone` 命令用于从 Geelato 服务器克隆已存在的应用到本地。该命令会连接服务器下载应用数据，包括模型定义、API 接口、工作流和页面配置，并渲染成标准的本地项目结构。克隆操作是团队协作和项目迁移的重要工具，支持通过 URL 快速定位和下载目标应用。

URL 格式要求：
```
http://{host}:{port}/{tenant}/{app-code}
```

其中：
- `{host}`：服务器地址，如 `localhost` 或 `api.geelato.com`
- `{port}`：端口号，如 `8080`，HTTPS 可省略端口
- `{tenant}`：租户标识，用于隔离不同租户的数据
- `{app-code}`：应用代码，标识具体应用

### 4.2 命令语法与参数

```bash
geelato clone <url> [选项]

# 或使用完整参数形式
geelato clone <url> --output <目录> --version <版本> --skip-extract
```

**参数说明：**

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| url | 字符串 | 是 | 应用仓库地址，格式为 `http://{host}:{port}/{tenant}/{app-code}` |

**选项说明：**

| 选项 | 简写 | 类型 | 默认值 | 说明 |
|------|------|------|--------|------|
| --output | -o | 字符串 | （应用代码） | 指定输出目录，不指定则使用应用代码作为目录名 |
| --version | 无 | 字符串 | latest | 指定要克隆的应用版本 |
| --skip-extract | 无 | 布尔值 | false | 跳过解压步骤，直接处理响应数据 |

### 4.3 克隆应用

```bash
# 基本用法 - 克隆应用到当前目录（自动以应用代码为目录名）
geelato clone http://localhost:8080/default/myapp

# 指定输出目录
geelato clone http://localhost:8080/default/myapp -o ./projects

# 指定租户
geelato clone http://localhost:8080/mytenant/myapp

# 指定应用版本
geelato clone http://localhost:8080/default/myapp --version latest

# 跳过解压步骤（高级用法）
geelato clone http://localhost:8080/default/myapp --skip-extract
```

**克隆过程详解：**

克隆命令的执行过程包含以下六个步骤，每个步骤都有明确的日志输出，方便跟踪进度和排查问题：

1. **URL 解析**：解析输入的 URL，提取 tenant 和 appCode，同时构建 API 服务器地址
2. **建立连接**：向服务器的 `/api/cli/app/clone` 端点发送 POST 请求，携带 appCode、tenant 和 version 参数
3. **下载数据**：接收服务器返回的应用数据，包含实体定义、页面配置、API 脚本和工作流信息
4. **渲染模型文件**：将服务器返回的元数据渲染为本地文件，包括 define.json、columns.json、check.json、fk.json 和视图 SQL 文件
5. **渲染页面文件**：处理页面数据，生成 source.json、release.json、preview.json 等页面资源文件
6. **保存到本地**：在指定目录创建完整的项目结构，包括必要的占位文件

**成功输出示例：**

```bash
$ geelato clone http://localhost:8080/default/myapp
[INFO] Parsing URL: tenant=default, appCode=myapp, apiURL=http://localhost:8080
[INFO] Cloning app 'myapp' to 'myapp'...
[INFO] Requesting: http://localhost:8080/api/cli/app/clone
[INFO] Response code: 200, data length: 1024
[INFO] Found nested format: v1
[INFO] Parsed: entities=3, pages=5, apis=8, workflows=2
[INFO] Processing entities: 3
[INFO] Processing pages: 5
[INFO] Processing APIs: 8
[INFO] Processing workflows: 2
[SUCCESS] Clone completed successfully!
[INFO] App cloned to: myapp
```

### 4.4 克隆后的目录结构

克隆成功后，会在指定目录创建完整的项目结构，目录组织遵循 Geelato 平台的标准规范：

```
myapp/
├── geelato.json              # 应用配置文件（包含 repo 信息）
├── api/                      # API 脚本目录
│   ├── .gitkeep             # 目录占位符（无 API 时）
│   └── userApi/             # API 脚本目录（以 API 代码命名）
│       ├── userApi.define.json
│       ├── userApi.source.js
│       └── userApi.release.js
├── meta/                     # 模型定义目录
│   └── User/                # 用户实体（以实体名称命名）
│       ├── User.define.json  # 表定义
│       ├── User.columns.json # 字段定义
│       ├── User.check.json  # 约束定义
│       ├── User.fk.json     # 外键定义
│       └── User.default.view.sql  # 默认视图 SQL
├── page/                     # 页面定义目录
│   ├── .gitkeep             # 目录占位符（无页面时）
│   └── userListPage/        # 页面目录（以页面代码命名）
│       ├── userListPage.define.json
│       ├── userListPage.source.json
│       ├── userListPage.release.json
│       └── userListPage.preview.json
└── workflow/                 # 工作流定义目录
    └── userApproval/         # 工作流目录
        ├── userApproval.define.json
        └── userApproval.json
```

### 4.5 应用配置文件说明

克隆成功后，会在应用根目录生成 `geelato.json` 配置文件，记录应用的核心元数据和仓库来源信息。

**repo.url 字段的作用：**

1. 标识应用的克隆来源，便于追溯代码的历史来源
2. 供 `push`、`pull`、`sync` 等命令解析服务器地址和租户信息
3. 可以通过 `geelato config repo` 命令更新仓库地址，实现应用迁移

### 4.6 错误处理

克隆过程中可能遇到以下常见错误，表格提供了错误原因分析和解决方法：

| 错误信息 | 原因分析 | 解决方法 |
|----------|----------|----------|
| invalid URL format | URL 格式不正确 | 检查 URL 是否包含协议头（http://）和必要的路径部分 |
| tenant and app code cannot be empty | URL 路径中缺少租户或应用代码 | 确保 URL 格式为 `/tenant/app-code` |
| clone failed with status 404 | 服务器端找不到指定应用 | 检查 appCode 和 tenant 是否正确，确认应用是否存在 |
| clone failed with status 500 | 服务器内部错误 | 检查服务器日志，联系管理员排查问题 |
| failed to request | 网络连接失败 | 检查网络连通性，确认服务器地址和端口是否正确 |
| failed to render and save | 文件写入失败 | 检查输出目录权限，确保有写入权限 |
| context deadline exceeded | 请求超时 | 增加超时时间或检查服务器响应速度 |

## 五、仓库配置命令详解

### 5.1 config repo 命令概述

`config repo` 命令用于管理应用的仓库地址（repo）。仓库地址记录了应用对应的 Geelato 服务器地址，用于后续的 `push`、`pull`、`sync` 等同步操作。

```bash
# 查看当前仓库地址
geelato config repo

# 设置仓库地址
geelato config repo http://localhost:8080/tenant/app-code
```

### 5.2 使用场景

**场景一：init 后设置**
```bash
# 初始化新应用
geelato init my-new-app
cd my-new-app

# 设置仓库地址（如果是克隆自某个应用）
geelato config repo http://localhost:8080/default/myapp
```

**场景二：更新仓库地址**
```bash
# 应用迁移到新的服务器
geelato config repo http://new-server:8080/tenant/new-app-code
```

## 五、模型管理命令详解

### 4.1 模型定义概述

模型是 Geelato 低代码平台的核心概念之一，用于描述应用的数据结构和业务实体。每个模型对应数据库中的一张表，包含字段定义、约束条件、外键关系、视图定义等多个维度的配置。Geelato CLI 提供了完整的模型管理命令，支持从创建到验证的全生命周期操作。

模型定义文件采用 JSON 格式存储在 `meta/{EntityName}/` 目录下，遵循统一的命名规范。例如，用户实体对应的模型文件存储在 `meta/User/` 目录中，包含 `User.define.json`（表定义）、`User.columns.json`（字段定义）、`User.fk.json`（外键定义）等多个文件。这种拆分设计使得模型的不同方面可以独立管理，便于团队协作和版本控制。

### 4.2 创建新模型

使用 `model create` 命令可以创建新的数据模型。创建过程支持交互式模式和命令行参数模式，您可以根据使用场景选择合适的方式。交互式模式会逐步引导您完成模型配置，适合初次接触 Geelato 平台的开发者；命令行参数模式则适合熟练用户进行批量操作或脚本化处理。

```bash
# 交互式创建模型
geelato model create Department

# 命令行参数模式创建
geelato model create Department \
  --table "platform_department" \
  --desc "部门信息"

# 创建模型同时添加字段
geelato model create Product \
  --table "platform_product" \
  --fields "name:string" \
  --fields "price:decimal" \
  --fields "stock:integer"
```

创建模型时，系统会自动在 `meta/` 目录下生成对应的实体文件夹。默认情况下，每个实体包含四个基础定义文件，分别描述表结构、字段列表、外键关系和约束条件。您可以根据实际需求修改这些文件，添加或删除模型属性。

```bash
# 模型创建后的目录结构
my-first-app/
└── meta/
    ├── User/
    │   ├── User.define.json     # 表定义
    │   ├── User.columns.json    # 字段定义
    │   ├── User.fk.json        # 外键定义
    │   └── User.check.json     # 约束定义
    └── Department/
        └── ...                 # 类似结构
```

### 4.3 表定义文件详解

表定义文件（`*.define.json`）描述模型的核心元数据，包括应用标识、表标识、业务名称等信息。这个文件是每个实体的入口定义，其他配置文件通过表标识与表定义建立关联。

```json
{
  "meta": {
    "version": "1.0.0",
    "appId": "app_my_first_app_001",
    "tenantId": "tenant_default",
    "createdAt": "2024-01-15T10:00:00Z",
    "updatedAt": "2024-01-15T10:00:00Z"
  },
  "table": {
    "id": "tbl_user_001",
    "appId": "app_my_first_app_001",
    "title": "用户",
    "entityName": "User",
    "tableName": "platform_user",
    "tableSchema": "platform",
    "tableType": "entity",
    "tableComment": "系统用户表",
    "enableStatus": 1,
    "linked": 1,
    "synced": false,
    "sourceType": "creation",
    "description": "存储系统用户信息",
    "acrossApp": false,
    "acrossWorkflow": false,
    "cacheType": "none"
  }
}
```

表定义中的关键字段说明如下：`id` 字段是表的唯一标识符，采用前缀加数字的命名格式，便于在不同应用间区分和引用；`entityName` 字段是实体的业务名称，用于代码生成和界面展示；`tableName` 字段对应数据库中的实际表名，建议使用小写字母和下划线的命名风格；`tableType` 字段标识表的类型，可选值包括 `entity`（实体表）、`dict`（字典表）、`link`（关联表）等。

### 4.4 字段定义文件详解

字段定义文件（`*.columns.json`）详细描述表中各字段的属性，包括字段名、数据类型、约束条件等。一个完善的字段定义是保证数据完整性和业务正确性的基础。Geelato 平台支持丰富的数据类型，可以满足各种业务场景的需求。

```json
{
  "meta": {
    "version": "1.0.0",
    "tableId": "tbl_user_001",
    "createdAt": "2024-01-15T10:00:00Z"
  },
  "columns": [
    {
      "id": "col_id_001",
      "columnName": "id",
      "columnComment": "主键ID",
      "dataType": "bigint",
      "length": 20,
      "isPrimaryKey": true,
      "isNullable": false,
      "defaultValue": null,
      "autoIncrement": true
    },
    {
      "id": "col_id_002",
      "columnName": "name",
      "columnComment": "用户姓名",
      "dataType": "varchar",
      "length": 50,
      "isPrimaryKey": false,
      "isNullable": false,
      "defaultValue": null
    },
    {
      "id": "col_id_003",
      "columnName": "login_name",
      "columnComment": "登录名",
      "dataType": "varchar",
      "length": 32,
      "isPrimaryKey": false,
      "isNullable": false,
      "defaultValue": null
    },
    {
      "id": "col_id_004",
      "columnName": "email",
      "columnComment": "邮箱地址",
      "dataType": "varchar",
      "length": 100,
      "isPrimaryKey": false,
      "isNullable": true,
      "defaultValue": null
    },
    {
      "id": "col_id_005",
      "columnName": "status",
      "columnComment": "用户状态",
      "dataType": "tinyint",
      "length": 1,
      "isPrimaryKey": false,
      "isNullable": false,
      "defaultValue": "1"
    },
    {
      "id": "col_id_006",
      "columnName": "created_at",
      "columnComment": "创建时间",
      "dataType": "datetime",
      "length": null,
      "isPrimaryKey": false,
      "isNullable": false,
      "defaultValue": "CURRENT_TIMESTAMP"
    }
  ]
}
```

### 4.5 外键定义文件详解

外键定义文件（`*.fk.json`）描述表间的关联关系，支持一对多、多对一、自引用等多种关联类型。外键定义是构建完整数据模型的关键组成部分，能够保证数据的引用完整性并简化业务逻辑的实现。

```json
{
  "meta": {
    "version": "1.0.0",
    "tableId": "tbl_user_001",
    "createdAt": "2024-01-15T10:00:00Z"
  },
  "foreignKeys": [
    {
      "id": "fk_user_dept",
      "tableId": "tbl_user_001",
      "tableName": "platform_user",
      "columnName": "dept_id",
      "foreignTable": "platform_department",
      "foreignColumn": "id",
      "onDelete": "SET_NULL",
      "onUpdate": "CASCADE",
      "foreignAppId": "app_user_mgmt_001"
    },
    {
      "id": "fk_dept_parent",
      "tableId": "tbl_dept_001",
      "tableName": "platform_department",
      "columnName": "parent_id",
      "foreignTable": "platform_department",
      "foreignColumn": "id",
      "onDelete": "SET_NULL",
      "onUpdate": "CASCADE"
    }
  ]
}
```

外键定义中的 `onDelete` 和 `onUpdate` 字段指定了级联操作策略。`CASCADE` 表示当父表记录被删除或更新时，自动对子表的相关记录进行相同操作；`SET_NULL` 表示将子表的外键字段设置为 NULL；`RESTRICT` 表示禁止删除或更新父表记录；`NO_ACTION` 表示不执行任何级联操作，与 RESTRICT 类似但实现时机不同。选择合适的级联策略需要根据业务需求仔细考虑，避免产生意外的数据变更。

## 五、API 管理命令详解

### 5.1 API 脚本概述

API 脚本是 Geelato 平台实现后端业务逻辑的核心机制。每个 API 脚本对应一个接口端点，可以处理 HTTP 请求、执行数据库操作、调用外部服务等任务。Geelato 的 API 脚本采用 JavaScript 语法编写，运行在平台提供的沙箱环境中，可以便捷地访问请求参数、数据库连接、缓存服务等核心功能。

API 脚本文件存储在 `api/` 目录下，采用 `{apiName}.api.{js|py|go}` 的命名格式。例如，用户列表接口对应的脚本文件名为 `getList.api.js`，用户详情接口对应的脚本文件名为 `getDetail.api.js`。CLI 工具会自动扫描 `api/` 目录下的所有脚本文件并注册为可用的 API 端点。

### 5.2 创建 API 脚本

使用 `api create` 命令可以创建新的 API 脚本模板。创建时需要指定 API 名称，CLI 工具会自动生成包含基本框架的脚本文件。您可以在此基础上编写具体的业务逻辑实现。

```bash
# 创建用户列表 API
geelato api create getUserList

# 创建用户详情 API
geelato api create getUserDetail

# 创建用户保存 API
geelato api create saveUser

# 使用不同语言
geelato api create getUserList --type python
geelato api create getUserList --type go
```

### 5.3 API 脚本结构详解

每个 API 脚本由三部分组成：元数据注释、参数定义和业务逻辑。元数据注释采用 `@` 符号标记，用于描述 API 的路由、方法、分组等信息；参数定义部分声明接口的输入参数及其类型；业务逻辑部分使用 JavaScript 编写具体的处理代码。

```javascript
/**
 * @api
 * @name getList
 * @path /api/user/getList
 * @method POST
 * @description 获取用户列表
 * @group user
 * @version 1.0.0
 * @appId app_user_mgmt_001
 * @tenantId tenant_default
 * @tableId tbl_user_001
 * @tableName User
 */

// @param
// name: pageNum
// type: Integer
// required: true
// default: 1
// description: 页码

// @param
// name: pageSize
// type: Integer
// required: true
// default: 10
// description: 每页条数

// @param
// name: keyword
// type: String
// required: false
// description: 搜索关键词

// @return
// type: PageResult
// description: 分页用户列表

(function() {
    // 获取请求参数
    var pageNum = parseInt($params.pageNum || 1);
    var pageSize = parseInt($params.pageSize || 10);
    var keyword = $params.keyword;

    // 构建查询条件
    var conditions = [];
    var params = [];

    if (keyword && keyword.trim() !== '') {
        conditions.push("(name LIKE ? OR login_name LIKE ?)");
        params.push('%' + keyword + '%');
        params.push('%' + keyword + '%');
    }

    var whereClause = conditions.length > 0 
        ? " WHERE " + conditions.join(" AND ")
        : "";

    // 查询总数
    var countSql = "SELECT COUNT(*) as total FROM platform_user" + whereClause;
    var countResult = $db.query(countSql, params);
    var total = countResult[0].total;

    // 查询列表
    var offset = (pageNum - 1) * pageSize;
    var listSql = "SELECT id, name, login_name, email, status, created_at " +
        "FROM platform_user" + whereClause +
        " ORDER BY created_at DESC LIMIT ? OFFSET ?";

    var queryParams = params.concat([pageSize, offset]);
    var list = $db.query(listSql, queryParams);

    // 返回结果
    return {
        code: 200,
        message: "success",
        data: {
            list: list,
            total: total,
            pageNum: pageNum,
            pageSize: pageSize,
            pages: Math.ceil(total / pageSize)
        }
    };
})();
```

### 5.4 API 脚本内置对象

API 脚本运行时提供了多个内置对象，便于快速实现业务逻辑。`$params` 对象包含客户端传入的所有请求参数，支持 GET 参数和 POST body 参数的合并访问；`$db` 对象提供数据库操作能力，支持查询、插入、更新、删除等操作；`$session` 对象用于访问和操作会话数据；`$cache` 对象提供缓存服务访问；`$http` 对象支持发起 HTTP 请求。

```javascript
// 参数访问示例
var userId = $params.id;
var name = $params.name;

// 数据库查询示例
var users = $db.query("SELECT * FROM platform_user WHERE status = ?", [1]);

// 数据库写入示例
var result = $db.execute(
    "INSERT INTO platform_user (name, login_name) VALUES (?, ?)",
    ["张三", "zhangsan"]
);

// 缓存操作示例
var cachedData = $cache.get("user:" + userId);
if (!cachedData) {
    cachedData = $db.query("SELECT * FROM platform_user WHERE id = ?", [userId]);
    $cache.set("user:" + userId, cachedData, 3600);
}

// HTTP 请求示例
var response = $http.get("https://api.example.com/users");
var externalData = response.data;
```

## 六、工作流管理命令详解

### 6.1 工作流定义概述

工作流是 Geelato 平台实现业务流程自动化的核心机制，支持 BPMN 2.0 标准规范。工作流定义文件采用 JSON 格式存储，描述流程的开始事件、结束事件、任务节点、网关决策、顺序流连接等元素。通过 Geelato CLI 的工作流管理命令，您可以创建、验证、部署符合业务需求的工作流定义。

工作流脚本存储在 `workflow/` 目录下，每个工作流对应一个独立的 JSON 文件。工作流定义包含元数据（名称、版本、描述）、流程结构（事件、任务、连接）两部分核心内容。CLI 工具提供了完整的工作流生命周期管理命令，覆盖从创建到部署的全过程。

### 6.2 创建工作流

使用 `workflow create` 命令可以创建新的工作流定义。创建时需要指定工作流名称，CLI 工具会自动生成包含开始事件、结束事件和处理任务的基础流程模板。

```bash
# 创建审批工作流
geelato workflow create approval

# 创建请假审批工作流
geelato workflow create leave-approval \
  --desc "员工请假审批流程"

# 查看所有工作流
geelato workflow list

# 验证工作流定义
geelato workflow validate
geelato workflow validate leave-approval
```

### 6.3 工作流结构详解

工作流定义文件采用标准化的 JSON 结构，包含元数据、事件定义、任务定义和顺序流四个主要部分。元数据描述工作流的基本属性；事件定义开始点和结束点的属性；任务定义流程中各处理节点的逻辑；顺序流定义节点间的执行路径。

```json
{
  "meta": {
    "name": "leave-approval",
    "description": "员工请假审批流程",
    "version": "1.0.0",
    "createdAt": "2024-01-15T10:00:00Z",
    "updatedAt": "2024-01-15T10:00:00Z"
  },
  "startEvents": [
    {
      "id": "start_1",
      "name": "开始",
      "type": "none"
    }
  ],
  "endEvents": [
    {
      "id": "end_approved",
      "name": "审批通过",
      "type": "none"
    },
    {
      "id": "end_rejected",
      "name": "审批驳回",
      "type": "none"
    }
  ],
  "tasks": [
    {
      "id": "task_submit",
      "name": "提交请假申请",
      "type": "userTask",
      "assignee": "${initiator}",
      "formKey": "leaveForm"
    },
    {
      "id": "task_approve",
      "name": "部门经理审批",
      "type": "userTask",
      "assignee": "manager",
      "candidateUsers": ["manager001", "manager002"],
      "formKey": "approveForm"
    },
    {
      "id": "task_notify",
      "name": "发送通知",
      "type": "serviceTask",
      "expression": "${notificationService.send(result)}"
    }
  ],
  "sequenceFlows": [
    {
      "id": "flow_1",
      "sourceRef": "start_1",
      "targetRef": "task_submit",
      "name": "提交"
    },
    {
      "id": "flow_2",
      "sourceRef": "task_submit",
      "targetRef": "task_approve",
      "name": "审批"
    },
    {
      "id": "flow_approved",
      "sourceRef": "task_approve",
      "targetRef": "end_approved",
      "name": "通过",
      "conditionExpression": "${outcome == 'approved'}"
    },
    {
      "id": "flow_rejected",
      "sourceRef": "task_approve",
      "targetRef": "end_rejected",
      "name": "驳回",
      "conditionExpression": "${outcome == 'rejected'}"
    }
  ]
}
```

### 6.4 工作流任务类型

Geelato 平台支持多种类型的任务节点，以适应不同的业务场景。用户任务（`userTask`）需要人工参与处理，适用于审批、确认等需要人工决策的环节；服务任务（`serviceTask`）自动执行预定义的逻辑，适用于数据处理、通知发送等自动化操作；脚本任务（`scriptTask`）执行预定义的脚本代码，提供更灵活的扩展能力；调用任务（`callActivity`）可以引用外部流程或子流程，支持流程的模块化设计。

```json
{
  "tasks": [
    {
      "id": "task_user",
      "name": "用户审批",
      "type": "userTask",
      "assignee": "${assignee}",
      "candidateGroups": ["managers"],
      "dueDate": "P3D",
      "priority": 50
    },
    {
      "id": "task_service",
      "name": "发送邮件",
      "type": "serviceTask",
      "implementation": "${emailService.send}",
      "topic": "emailNotification"
    },
    {
      "id": "task_script",
      "name": "计算折扣",
      "type": "scriptTask",
      "scriptFormat": "javascript",
      "script": "var discount = orderAmount * 0.1; execution.setVariable('discount', discount);"
    },
    {
      "id": "task_call",
      "name": "调用审批流程",
      "type": "callActivity",
      "calledElement": "approval-process",
      "independent": false
    }
  ]
}
```

## 七、云端同步命令详解

### 7.1 同步机制概述

云端同步是 Geelato CLI 的核心功能之一，用于将本地开发的应用同步到 Geelato 云端平台。同步功能支持完整的变更管理流程，包括变更检测、冲突检测、版本管理和变更推送。通过同步机制，团队成员可以协作开发，确保本地和云端的数据一致性。

Geelato CLI 提供两套同步命令：
- **直接同步命令**：直接对当前目录进行操作，无需指定应用
  - `geelato push` - 推送到云端
  - `geelato pull` - 从云端拉取
  - `geelato diff` - 查看差异
  - `geelato watch` - 监听变化自动同步

- **同步管理命令**（在 sync 子命令下）：
  - `geelato sync push` - 推送到云端
  - `geelato sync pull` - 从云端拉取
  - `geelato sync status` - 查看同步状态
  - `geelato sync resolve` - 解决冲突

两种方式功能相同，直接命令更简洁，子命令更规范。

同步状态记录在 `.geelato/sync-state.json` 文件中，包含当前版本号、最后同步时间和已同步文件哈希等信息。CLI 工具通过比对本地文件和已记录状态来识别新增、修改或删除的变更。

### 7.2 推送变更到云端

使用 `push` 命令可以将本地变更推送到云端平台。推送前，系统会自动检测所有变更并显示变更列表；确认后，系统会打包变更内容并上传到云端。

```bash
# 推送变更（交互式确认）
geelato push

# 推送并指定提交消息
geelato push "feat: 添加用户头像上传功能"

# 推送所有变更（不需确认）
geelato push --all

# 演练模式（显示变更但不实际推送）
geelato push --dry-run

# 强制推送（覆盖远程冲突）
geelato push --force
```

推送过程中，系统会执行以下操作：扫描 `meta/`、`api/`、`page/`、`workflow/` 目录下的变更文件；计算每个文件的哈希值并与已记录状态比对；生成变更报告包含新增、修改、删除的文件列表；打包变更内容并上传到云端平台；更新本地同步状态记录。

```bash
# 推送输出示例
开始推送变更: feat: 添加用户头像上传功能

发现 3 个变更
变更列表:
---------
  [新增] meta/Avatar/Avatar.define.json (新增文件)
  [新增] meta/Avatar/Avatar.columns.json (新增文件)
  [修改] api/user/saveUser.api.js (修改文件)

确认推送 3 个变更？(y/n): y
正在推送: meta/Avatar/Avatar.define.json
正在推送: meta/Avatar/Avatar.columns.json
正在推送: api/user/saveUser.api.js

推送完成
版本: 20240115143025
变更数量: 3
```

### 7.3 从云端拉取更新

使用 `pull` 命令可以从云端平台拉取最新版本到本地。拉取前，系统会显示云端版本与本地版本的差异；确认后，系统会下载并解压更新包覆盖本地文件。

```bash
# 拉取更新（交互式确认）
geelato pull

# 强制拉取（覆盖本地修改）
geelato pull --force

# 演练模式（显示变更但不实际拉取）
geelato pull --dry-run
```

### 7.4 查看差异

使用 `diff` 命令可以查看本地与云端之间的差异，包括新增、修改、删除的文件列表。

```bash
# 查看差异
geelato diff

# JSON 格式输出
geelato diff --json

# 输出示例：
# 差异对比结果
# ============
#
# 云端新增:
#   meta/Product/Product.define.json
#   meta/Product/Product.columns.json
#
# 本地修改:
#   api/user/saveUser.api.js
#
# 已删除（云端无此文件）:
#   api/old/delete-me.api.js
```

### 7.5 监听文件变化

使用 `watch` 命令可以监听本地文件变化，自动同步到云端。适用于开发过程中持续同步的场景。

```bash
# 启动监听模式
geelato watch

# 输出示例：
# [INFO] 启动文件监听...
# [INFO] 监听目录: /path/to/my-app
# [INFO] 模式: 自动同步
# [CREATED] meta/Product/Product.define.json
# [SYNCED] 已同步到云端
# [MODIFIED] api/user/saveUser.api.js
# [SYNCED] 已同步到云端
```

### 7.3 查看同步状态

使用 `sync status` 命令可以查看当前的同步状态，包括本地版本、云端版本、待推送变更数、待拉取变更数以及冲突情况。

```bash
# 查看同步状态
geelato sync status

# JSON 格式输出
geelato sync status --json

# 详细输出
geelato sync status --verbose

# 输出示例：
同步状态
========

本地版本:  20240115140000
云端版本:  20240115143025

已是最新版本
最后同步: 2024-01-15 14:00:00
```

### 7.4 解决同步冲突

当本地和云端同时修改同一文件时，系统会检测到冲突并阻止自动合并。冲突情况下，需要使用 `sync resolve` 命令手动解决。CLI 工具提供了多种冲突解决策略供选择。

```bash
# 查看冲突列表
geelato sync status

# 交互式解决冲突
geelato sync resolve

# 使用策略解决所有冲突
geelato sync resolve --all --strategy ours      # 保留本地版本
geelato sync resolve --all --strategy theirs     # 保留云端版本

# 解决特定文件冲突
geelato sync resolve api/user/saveUser.api.js --strategy manual
```

冲突解决策略说明：`ours` 策略保留本地版本，放弃云端修改；`theirs` 策略保留云端版本，放弃本地修改；`manual` 策略打开编辑器进行手动合并。手动合并时，CLI 工具会打开配置的编辑器，您需要直接编辑冲突文件，处理好冲突标记后保存退出。

## 八、MCP 平台能力管理详解

### 8.1 MCP 能力概述

MCP（Model-Context-Protocol，模型-上下文-协议）是 Geelato 平台的扩展能力机制。通过 MCP，平台可以提供预定义的数据模型、业务逻辑和集成能力，帮助开发者快速构建应用。CLI 工具提供了 MCP 能力管理命令，支持能力浏览、安装、同步等操作。

MCP 能力分为多种类型：认证能力（`auth`）提供用户认证、授权、令牌管理等安全相关功能；数据库能力（`database`）提供高级数据库查询和事务管理功能；消息队列能力（`mq`）提供消息发布和订阅功能；存储能力（`storage`）提供文件存储和对象存储集成；AI 能力（`ai`）提供人工智能和机器学习服务接口。

### 8.2 浏览可用能力

使用 `mcp list` 命令可以浏览本地和云端平台可用的 MCP 能力。列表显示能力的详细信息，包括名称、版本、分类和描述。

```bash
# 列出所有可用能力
geelato mcp list

# 按分类筛选
geelato mcp list --category database
geelato mcp list --category auth
geelato mcp list --category ai

# JSON 格式输出
geelato mcp list --json

# 输出示例：
可用 MCP 能力
===============

- 数据库能力 (database)
  版本: 1.0.0
  描述: 提供数据库查询、更新、事务等能力

- 认证能力 (auth)
  版本: 1.0.0
  描述: 提供用户认证、授权、令牌管理等能力
```

### 8.3 搜索和安装能力

使用 `mcp search` 命令可以搜索特定的能力，使用 `mcp sync` 命令可以将能力同步到本地或推送到云端。

```bash
# 搜索能力
geelato mcp search database
geelato mcp search auth

# 查看能力详情
geelato mcp info database

# 同步能力到本地（从云端拉取）
geelato mcp sync --direction pull

# 推送能力到云端（发布能力）
geelato mcp sync --direction push
```

## 九、应用配置文件详解

### 9.1 geelato.json 结构

每个 Geelato 应用根目录下的 `geelato.json` 是应用的核心配置文件，定义了应用的身份标识、依赖关系、资源路径等关键信息。这个文件采用 JSON 格式存储，在应用初始化时自动生成，后续可根据需要手动编辑。

```json
{
  "name": "user-management",
  "version": "1.0.0",
  "description": "用户管理应用",
  "appId": "app_user_mgmt_001",
  "tenantId": "tenant_default",
  "appCode": "user",
  "dependencies": [
    {
      "app": "base-system",
      "version": "^1.0.0"
    }
  ],
  "tables": [
    "meta/tables/*.table.json"
  ],
  "columns": [
    "meta/columns/*.columns.json"
  ],
  "views": [
    "meta/views/*.view.sql"
  ],
  "foreignKeys": [
    "meta/foreign-keys/*.fk.json"
  ],
  "checks": [
    "meta/checks/*.check.json"
  ],
  "apis": [
    "api/**/*.api.js"
  ],
  "pages": [
    "page/**/page.json"
  ],
  "workflows": [
    "workflow/**/*.workflow.json"
  ],
  "scripts": {
    "pre-push": "echo 'Running pre-push validation...'",
    "post-pull": "echo 'Sync completed!'"
  }
}
```

配置项详细说明如下：`name` 是应用的显示名称，会在界面和报告中展示；`appId` 是应用的唯一标识符，由系统自动生成；`version` 是应用的版本号，遵循语义化版本规范；`dependencies` 声明了应用依赖的其他应用或模块；`tables`、`columns`、`views` 等路径模式定义了各类资源文件的扫描路径；`scripts` 节点定义了钩子脚本，用于在特定操作前后执行自定义逻辑。

### 9.2 应用模板结构

Geelato CLI 支持从模板创建应用，不同模板包含不同的预置内容。默认模板（`default`）包含基础目录结构和配置文件；空白模板（`blank`）只创建最小化的应用骨架；示例模板（`sample`）包含完整示例代码和参考文档。

```bash
# 使用默认模板创建
geelato init my-app --template default

# 使用空白模板创建
geelato init my-app --template blank

# 查看可用的应用模板
geelato init --help
```

## 十、最佳实践指南

### 10.1 项目结构规范

良好的项目结构是保证团队协作效率的基础。建议遵循以下规范组织您的 Geelato 应用代码：按业务域划分实体目录，将相关的模型、API 和页面组织在同一个子目录下；API 脚本按功能模块组织在 `api/{module}/` 子目录中；页面组件按页面组织在 `page/{pageName}/` 目录下；工作流文件按流程类别组织在 `workflow/{category}/` 目录下。

```
推荐的项目结构示例：
my-app/
├── api/
│   ├── user/
│   │   ├── getList.api.js
│   │   ├── getDetail.api.js
│   │   └── saveUser.api.js
│   ├── order/
│   │   ├── createOrder.api.js
│   │   └── getOrderList.api.js
│   └── common/
│       └── upload.api.js
├── meta/
│   ├── User/
│   ├── Order/
│   └── Product/
├── page/
│   ├── user-list/
│   ├── user-detail/
│   └── order-form/
├── workflow/
│   ├── approval/
│   └── notification/
└── doc/
```

### 10.2 版本控制建议

强烈建议将 Geelato 应用纳入 Git 版本控制系统管理。由于模型定义和 API 脚本都是纯文本文件，Git 的差异比较和版本回滚功能可以很好地支持这类内容的变更追踪。

```bash
# 初始化 Git 仓库
git init
git add .
git commit -m "feat: 初始化用户管理应用"

# 创建特性分支
git checkout -b feature/user-avatar

# 开发完成后合并
git checkout main
git merge feature/user-avatar

# 查看变更历史
git log --oneline
```

### 10.3 命名规范

一致的命名规范可以提升代码的可读性和可维护性。遵循以下建议进行命名：实体名称使用 PascalCase（如 `UserProfile`、`OrderItem`）；表名和字段名使用 snake_case（如 `user_profile`、`order_item`）；API 名称使用 camelCase（如 `getUserList`、`saveOrder`）；文件命名与对应实体或 API 保持一致。

## 十一、常见问题与解决方案

### 11.1 安装问题排查

在安装 Geelato CLI 过程中可能遇到的问题及其解决方案。如果遇到命令未找到的错误，请检查 PATH 环境变量是否包含 CLI 安装目录；如果遇到权限错误，Linux/macOS 用户可能需要使用 `sudo` 提升权限，Windows 用户请以管理员身份运行安装程序；如果下载速度慢，可以尝试使用镜像站点或代理服务器。

```bash
# 检查安装是否成功
which geelato    # Linux/macOS
where geelato    # Windows

# 手动添加 PATH（临时）
export PATH=$PATH:/usr/local/bin  # Linux/macOS
$env:PATH += ";C:\Tools\geelato"  # Windows PowerShell
```

### 11.2 使用问题排查

日常使用中可能遇到的常见问题及其解决方案。API 脚本执行报错时，检查 `$params` 参数访问方式是否正确，确认 SQL 语句语法无误；同步推送失败时，检查网络连接和 API 认证配置；模型验证不通过时，核对 JSON 格式是否正确，关键字段是否完整。

```bash
# 检查配置文件
geelato config list

# 验证应用结构
geelato validate

# 查看详细错误日志
geelato --verbose push
```

## 十二、命令参考速查

### 12.1 完整命令列表

以下是 Geelato CLI 所有命令的快速参考，方便日常使用时快速查阅。

| 命令 | 功能 | 常用参数 |
|------|------|----------|
| **应用管理** |
| `geelato init [name]` | 初始化新应用 | `--desc`、`--template`、`--skip-git` |
| `geelato validate` | 验证当前应用结构 | 无 |
| `geelato config list` | 查看当前应用信息 | 无 |
| **模型管理** |
| `geelato model create [name]` | 创建新模型 | `--table`、`--desc`、`--fields` |
| **API管理** |
| `geelato api create [name]` | 创建新 API | `--type` |
| **工作流管理** |
| `geelato workflow create [name]` | 创建新工作流 | `--desc` |
| `geelato workflow list` | 列出所有工作流 | 无 |
| `geelato workflow validate` | 验证工作流 | `--strict` |
| **云端同步** |
| `geelato push` | 推送变更到云端 | `--message`、`--all`、`--dry-run` |
| `geelato pull` | 从云端拉取更新 | `--force`、`--dry-run` |
| `geelato diff` | 查看本地与云端差异 | `--json` |
| `geelato watch` | 监听文件变化自动同步 | 无 |
| `geelato sync status` | 查看同步状态 | `--json`、`--verbose` |
| `geelato sync resolve` | 解决同步冲突 | `--strategy`、`--all` |
| **MCP能力管理** |
| `geelato mcp list` | 列出可用能力 | `--category`、`--json` |
| `geelato mcp search [keyword]` | 搜索能力 | 无 |
| **配置管理** |
| `geelato config get [key]` | 查看配置项 | 无 |
| `geelato config set [key] [value]` | 设置配置项 | 无 |

### 12.2 全局参数说明

Geelato CLI 支持以下全局参数，可以在所有命令中使用。

| 参数 | 简写 | 说明 |
|------|------|------|
| `--help` | `-h` | 显示命令帮助信息 |
| `--version` | `-v` | 显示 CLI 版本 |
| `--verbose` | `-v` | 输出详细日志 |
| `--json` | 无 | 使用 JSON 格式输出 |
| `--config` | `-c` | 指定配置文件路径 |

```bash
# 显示帮助
geelato --help
geelato push --help
geelato pull --help

# 详细输出模式
geelato --verbose push

# 指定配置文件
geelato --config /path/to/config.yaml push

# JSON 格式输出
geelato diff --json
```

### 12.3 快速入门命令序列

以下是典型的开发工作流命令序列：

```bash
# 1. 初始化新应用
geelato init my-app
cd my-app

# 2. 配置云端地址
geelato config set api.url "https://api.geelato.com"

# 3. 创建模型
geelato model create User

# 4. 创建 API
geelato api create getUserList

# 5. 验证应用结构
geelato validate

# 6. 推送到云端
geelato push "feat: 初始化用户模块"

# 7. 开发过程中监听变化
geelato watch
```

## 附录：示例应用参考

Geelato CLI 项目提供了完整的示例应用 `user-management-app`，位于 `geelato_new_app/example/` 目录下。该示例演示了用户管理功能的标准实现方式，包括用户实体定义、部门实体定义、用户相关 API 脚本、工作流定义等内容，是学习和参考的绝佳资源。

```bash
# 示例应用目录结构
geelato_new_app/example/user-management-app/
├── api/
│   └── user/
│       ├── getList.api.js       # 用户列表
│       ├── getDetail.api.js     # 用户详情
│       └── saveUser.api.js      # 保存用户
├── meta/
│   ├── User/
│   │   ├── User.define.json     # 表定义
│   │   ├── User.columns.json     # 字段定义
│   │   ├── User.fk.json         # 外键
│   │   └── User.check.json      # 约束
│   └── Department/
│       └── ...                   # 类似结构
├── page/
│   └── user-list/
│       └── page.json            # 页面配置
└── workflow/
    └── leave-approval/
        └── leave-approval.workflow.json  # 工作流
```

通过学习和参考示例应用，您可以快速掌握 Geelato 平台的开发规范和最佳实践，从而更高效地构建自己的业务应用。
