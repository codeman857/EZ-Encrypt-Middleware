# EZ 加密中间件教程

本教程以 中间件 和 V2Board 管理员后端在同一服务器的 宝塔 上面搭建为例，不同服务器上搭建方法类似，不再赘述

## 安装 Go Version 1.24 环境

### 执行如下命令安装 Go 环境 (以 AMD 架构为例)

```bash
# 下载安装包 以 AMD 架构为例
wget https://go.dev/dl/go1.24.0.linux-amd64.tar.gz

# 解压到 /usr/local
tar -C /usr/local -xzf go1.24.0.linux-amd64.tar.gz

# 配置环境变量，编辑 ~/.bashrc 文件
nano ~/.bashrc

# 在 ~/.bashrc 文件末尾增加 如下一行内容
export PATH=$PATH:/usr/local/go/bin
# 保存并退出：在 nano 中，按 Ctrl+X，然后按 Y，再按 Enter

# 让配置生效
source ~/.bashrc
```

### 检查是否安装成功

```bash
# 执行检查 go 版本命令
go version

# 如果一切顺利 应该会看到类似以下内容的输出
go version go1.24.0 linux/amd64
```

## 下载并编译 EZ 加密中间件 (以 root 用户根目录为例)

### 1. 下载并编译

```bash
# 下载 R 佬的中间件源码
git clone https://github.com/codeman857/EZ-Encrypt-Middleware.git

# 切换到中间件目录
cd EZ-Encrypt-Middleware/

# 执行 go 语言编译命令 以 "ez-api" 命名编译好的二进制文件
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ez-api

# 执行 ls 命令查看是否生成了 ez-api 二进制文件
ls

# 赋予二进制文件 ez-api 执行权限
chmod 777 /root/EZ-Encrypt-Middleware/ez-api
```

### 2. 修改 /root/EZ-Encrypt-Middleware/ 下的 .env 文件配置

```bash
# 查看 VPS 端口占用情况
ss -tulpn

# 找一个未被占用的端口号，这里我以 3939 为例 (九宫格 "e" "z" 对应数字 "3" 和 "9")
nano .env
```

```bash
# 1. 基础服务器设置
PORT=3939                                  # 服务器监听端口
BACKEND_API_URL=https://api.xxx.com        # 后端真实 API 根地址 不带 /api/v1（无尾斜杠）

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

# 4. AES 加解密配置
# 中间件加密 KEY 必须是 16 位的 16 进制字符串
# 必须和 EZ 前端配置文件 /src/config/index.js 中的 API_MIDDLEWARE_KEY 保持一致
# 在线生成地址 https://www.bejson.com/math/hex_gen/
AES_KEY=4c6f8e5f9467dc71
```

### 3. 宝塔 进程守护管理器 添加 进程守护

这一步是将 ez-api 中间件服务运行在 VPS 的 3939 端口上

```bash
# Name
ez-api

# Run User
root

# Process directory
/root/EZ-Encrypt-Middleware/

# Start command
/root/EZ-Encrypt-Middleware/ez-api
```

<img width="1180" height="596" alt="添加进程守护" src="https://github.com/user-attachments/assets/2553b171-7d10-4688-8754-baf39bbc2617" />


### 4. 宝塔 添加用于和前端通信的反向代理 api 站点

添加站点 ( ez.api 域名需要解析到 VPS IP 上)<br />
PHP Version 选择 Static

<img width="701" height="628" alt="添加静态站点" src="https://github.com/user-attachments/assets/97255110-ae81-4ca3-b001-2ba95e2b3d69" />


申请 SSL 并设置反向代理

```bash
# Target URL
http://127.0.0.1:3939
```

<img width="859" height="738" alt="设置反向代理" src="https://github.com/user-attachments/assets/9bf2ca73-927c-4c67-a2a6-cdea56763123" />


### 5. EZ 主题前端 修改 `/src/config/index.js` 相应的 API 配置

```bash
# API_MIDDLEWARE_URL
https://ez.api

# API_MIDDLEWARE_KEY
4c6f8e5f9467dc71
```

<img width="1095" height="330" alt="修改前端配置文件" src="https://github.com/user-attachments/assets/2e7891a4-63d8-49b2-b927-56c2cfe4b485" />

