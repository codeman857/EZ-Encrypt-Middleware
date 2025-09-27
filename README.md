# API Proxy Server

一个基于Go语言开发的API代理服务器，支持路径加密解密、CORS跨域处理、支付回调免验证等功能。

部署教程 🤌 [查看](https://github.com/codeman857/EZ-Encrypt-Middleware/wiki/aapanel-%E9%83%A8%E7%BD%B2%E6%95%99%E7%A8%8B)

## 功能特性

1. **路径加密解密**：对请求路径进行AES解密后再转发到后端API
2. **CORS跨域支持**：支持通配符和特定域名的跨域配置
3. **支付回调免验证**：特定路径不进行加密解密直接转发
4. **请求超时控制**：可配置的请求超时时间
5. **日志记录**：可开关的请求日志记录功能
6. **环境配置**：通过.env文件进行配置管理

## 项目结构

```
.
├── main.go                 # 主程序入口
├── go.mod                  # Go模块定义
├── go.sum                  # Go模块校验和
├── .env                    # 环境配置文件
├── README.md               # 项目说明文档
├── config/
│   ├── config.go           # 配置管理
│   └── config_test.go      # 配置测试
├── proxy/
│   └── proxy.go            # 代理处理逻辑
└── utils/
    └── encryption.go       # 加密解密工具
```

## 配置说明

### 环境变量配置 (.env文件)

```bash
# 1. 基础服务器设置
PORT=3000                                  # 服务器监听端口
BACKEND_API_URL=https://skhsn6q4pnv95.ezdemo.xyz # 后端真实 API 根地址不带 /api/v1（无尾斜杠）
PATH_PREFIX=/ez/ez                               # 路径前缀，为空则处理所有路径，否则只处理匹配前缀的路径

# 2. CORS / 安全设置
CORS_ORIGIN=*                              # 允许的 CORS 源；* 表示全部
ALLOWED_ORIGINS=*                          # 请求来源白名单，逗号分隔或 * 通配
REQUEST_TIMEOUT=30000                      # 请求超时(ms)
ENABLE_LOGGING=false                       # 是否输出请求日志
DEBUG_MODE=false                           # 是否输出调试日志

# 3. 支付回调免验证路径
# 多条用英文逗号分隔，须写完整路径（含前缀）
# 例如: /api/v1/guest/payment/notify/EPay/12345, /api/v1/guest/payment/notify/Alipay/ABC123
ALLOWED_PAYMENT_NOTIFY_PATHS=

# 5. AES 加解密配置
# 中间件加密KEY必须是16位的16进制字符串，必须和前端key保持一致
AES_KEY=4c6f8e5f9467dc71
```

### 配置项详解

1. **PORT**：服务器监听端口，默认3000
2. **BACKEND_API_URL**：后端API的基础URL，不包含/api/v1部分
3. **PATH_PREFIX**：路径前缀，为空则处理所有路径，否则只处理匹配前缀的路径
4. **CORS_ORIGIN**：
   - `*`：允许所有来源
   - 具体域名：只允许指定域名
4. **ALLOWED_ORIGINS**：
   - `*`：允许所有来源
   - 多个域名用逗号分隔：`https://a.com, https://b.com`
5. **REQUEST_TIMEOUT**：请求超时时间，单位毫秒
6. **ENABLE_LOGGING**：是否启用请求日志记录
7. **DEBUG_MODE**：是否启用调试模式
8. **ALLOWED_PAYMENT_NOTIFY_PATHS**：支付回调免验证路径，多个路径用逗号分隔
9. **AES_KEY**：AES解密密钥，必须是16位的16进制字符串

## 运行方式

### 1. 使用预编译版本（推荐）

从 [Releases](https://github.com/codeman857/EZ-Encrypt-Middleware/releases) 页面下载对应平台的预编译版本：

- **Linux ARM64**: `EZ-Encrypt-Middleware-linux-arm64.tar.gz`
- **Linux AMD64**: `EZ-Encrypt-Middleware-linux-amd64.tar.gz`
- **Windows AMD64**: `EZ-Encrypt-Middleware-windows-amd64.zip`
- **macOS AMD64**: `EZ-Encrypt-Middleware-darwin-amd64.tar.gz`
- **macOS ARM64**: `EZ-Encrypt-Middleware-darwin-arm64.tar.gz`

下载后解压，压缩包中已包含配置好的 `.env` 文件，可直接运行：

```bash
# Linux/macOS
./proxy-server

# Windows
proxy-server.exe
```

> **注意**: 如果下载的版本是通过手动触发构建的，`.env` 文件已包含你设置的配置参数。如果是通过标签触发的构建，会使用 GitHub Secrets 中的默认配置。

### 2. 编译运行

```bash
# 编译
GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -o proxy-server // arm64

GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o proxy-server // amd64

# 运行
./proxy-server
```

### 3. 直接运行

```bash
go run main.go
```

## GitHub Actions 自动构建

本项目配置了 GitHub Actions 自动构建和发布：

### 自动构建触发条件

1. **推送标签**: 当推送以 `v` 开头的标签时（如 `v1.0.0`），会自动构建并发布到 GitHub Releases
2. **手动触发**: 在 GitHub Actions 页面可以手动触发构建，支持自定义配置参数
3. **CI 测试**: 每次推送到 `main` 或 `develop` 分支时会运行测试

### 构建的架构

- Linux ARM64
- Linux AMD64  
- Windows AMD64
- macOS AMD64
- macOS ARM64

### 配置方式

#### 方式1: 手动触发时设置参数

在 GitHub Actions 页面手动触发构建时，可以设置以下参数：

- **PORT**: 服务器监听端口（默认: 3000）
- **BACKEND_API_URL**: 后端API地址（必填）
- **PATH_PREFIX**: 路径前缀（默认: /ez/ez）
- **CORS_ORIGIN**: CORS源（默认: *）
- **ALLOWED_ORIGINS**: 允许的来源（默认: *）
- **REQUEST_TIMEOUT**: 请求超时毫秒数（默认: 30000）
- **ENABLE_LOGGING**: 是否启用日志（默认: false）
- **DEBUG_MODE**: 是否启用调试模式（默认: false）
- **ALLOWED_PAYMENT_NOTIFY_PATHS**: 支付回调免验证路径（默认: 空）
- **AES_KEY**: AES加密密钥，16位16进制（必填）

#### 方式2: 使用GitHub Secrets设置默认值

在仓库的 Settings → Secrets and variables → Actions 中添加以下secrets：

```
BACKEND_API_URL=https://your-api.com
PATH_PREFIX=/your/prefix
CORS_ORIGIN=https://yourdomain.com
ALLOWED_ORIGINS=https://yourdomain.com,https://anotherdomain.com
REQUEST_TIMEOUT=60000
ENABLE_LOGGING=true
DEBUG_MODE=false
ALLOWED_PAYMENT_NOTIFY_PATHS=/api/v1/payment/notify
AES_KEY=your16hexkey123456
```

### 使用自动构建

1. **通过标签发布**：
   ```bash
   git tag v1.0.0
   git push origin v1.0.0
   ```

2. **手动触发构建**：
   - 访问 [Actions](https://github.com/codeman857/EZ-Encrypt-Middleware/actions) 页面
   - 选择 "Build and Release" 工作流
   - 点击 "Run workflow"
   - 填写必要的参数（BACKEND_API_URL 和 AES_KEY 为必填）

3. **查看构建结果**：
   - 构建完成后，在 [Releases](https://github.com/codeman857/EZ-Encrypt-Middleware/releases) 页面下载预编译版本
   - 下载的压缩包中已包含配置好的 `.env` 文件，可直接运行

## 使用说明

### 1. 加密路径请求

对于需要加密的请求：
- 路径：`/{base64_encrypted_path}`
- 请求头：`X-IV`（加密使用的IV）

服务器会：
1. 从路径中提取Base64编码的加密路径
2. 使用`X-IV`头中的IV和配置的`AES_KEY`进行AES解密
3. 将解密后的路径拼接到`BACKEND_API_URL`后转发请求

### 2. 支付回调请求

对于配置在`ALLOWED_PAYMENT_NOTIFY_PATHS`中的路径：
- 请求会直接转发到后端API，不进行加密解密处理
- 不进行CORS验证

### 3. CORS处理

根据`CORS_ORIGIN`和`ALLOWED_ORIGINS`配置自动处理跨域请求。

## 依赖库

- [Gin](https://github.com/gin-gonic/gin)：Web框架
- [gin-contrib/cors](https://github.com/gin-contrib/cors)：CORS中间件
- [joho/godotenv](https://github.com/joho/godotenv)：环境变量加载
- [deatil/go-cryptobin](https://github.com/deatil/go-cryptobin)：加密解密库

## API接口

服务器会处理所有进入的请求：

1. **支付回调路径**：直接转发，不加密解密
2. **其他路径**：进行Base64解码和AES解密后转发

## 日志输出

启用日志后会输出：
- 服务器启动信息
- 请求处理信息
- 错误信息
- 配置信息

## 部署建议

1. 生产环境建议将`DEBUG_MODE`设置为`false`
2. 根据实际需求配置`REQUEST_TIMEOUT`
3. 合理配置CORS策略，避免安全风险
4. 定期更新`AES_KEY`提高安全性
