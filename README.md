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

### 1. 编译运行

```bash
# 编译
GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -o proxy-server // arm64

GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o proxy-server // amd64

# 运行
./proxy-server
```

### 2. 直接运行

```bash
go run main.go
```

## 使用说明 

宝塔搭建请访问 [Middleware-BT-Tutorial.md](Middleware-BT-Tutorial.md)

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
