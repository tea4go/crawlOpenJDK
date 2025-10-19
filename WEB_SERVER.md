# OpenJDK 下载中心 - Web 服务器使用指南

## 🚀 快速启动

### 启动服务器

```bash
go run server.go
```

或编译后运行：

```bash
# 编译
go build -o openjdk-server server.go

# 运行
./openjdk-server          # Linux/macOS
openjdk-server.exe        # Windows
```

### 访问地址

服务器启动后，可通过以下地址访问：

- **本地访问**: http://localhost:8080
- **局域网访问**: http://你的IP地址:8080

## 📋 服务器功能

### ✅ 主要功能

1. **静态文件服务**
   - 提供 index.html 网页访问
   - 提供 jdkindex.json 数据访问
   - 自动设置正确的 Content-Type

2. **请求日志**
   - 记录所有访问请求
   - 显示请求方法、路径、来源 IP
   - 显示响应状态码

3. **CORS 支持**
   - 允许跨域请求
   - 支持前端开发调试

4. **文件检查**
   - 启动时自动检查必需文件
   - 缺少文件时给出明确提示

### 📁 支持的文件类型

服务器自动识别并设置正确的 MIME 类型：

| 文件类型 | MIME 类型 |
|---------|-----------|
| .html   | text/html; charset=utf-8 |
| .css    | text/css; charset=utf-8 |
| .js     | application/javascript; charset=utf-8 |
| .json   | application/json; charset=utf-8 |
| .png/.jpg/.gif | image/* |
| .svg    | image/svg+xml |
| .txt    | text/plain; charset=utf-8 |

## 🔧 配置选项

### 自定义端口

**方法一：环境变量**

```bash
# Linux/macOS
export PORT=3000
go run server.go

# Windows (PowerShell)
$env:PORT="3000"
go run server.go

# Windows (CMD)
set PORT=3000
go run server.go
```

**方法二：修改代码**

编辑 `server.go` 中的 `defaultPort` 常量：

```go
const (
    defaultPort = "3000"  // 修改为你想要的端口
)
```

## 📊 日志示例

### 启动日志

```
╔═══════════════════════════════════════════════════════╗
║                                                       ║
║       ☕ OpenJDK 下载中心 - Web 服务器               ║
║                                                       ║
╚═══════════════════════════════════════════════════════╝

🚀 服务器启动成功！

📍 访问地址:
   - 本地访问: http://localhost:8080
   - 网络访问: http://0.0.0.0:8080

📁 服务文件:
   - index.html  (主页面)
   - jdkindex.json (数据文件)

⌨️  按 Ctrl+C 停止服务器
```

### 访问日志

```
2025/10/19 12:29:10 收到请求: GET / 来自 [::1]:37109
2025/10/19 12:29:10 200 - GET /
2025/10/19 12:29:29 收到请求: GET /jdkindex.json 来自 [::1]:33384
2025/10/19 12:29:29 200 - GET /jdkindex.json
```

## 🛠️ 开发模式

### 热重载开发

使用 `air` 工具实现热重载：

```bash
# 安装 air
go install github.com/cosmtrek/air@latest

# 创建 .air.toml 配置文件
cat > .air.toml << EOF
root = "."
tmp_dir = "tmp"

[build]
cmd = "go build -o ./tmp/server server.go"
bin = "tmp/server"
include_ext = ["go", "html", "json"]
exclude_dir = ["tmp"]

[color]
main = "magenta"
watcher = "cyan"
build = "yellow"
runner = "green"
EOF

# 启动热重载
air
```

### 调试模式

添加详细日志输出：

```go
// 在 server.go 中添加调试日志
log.SetFlags(log.LstdFlags | log.Lshortfile)
```

## 🔒 安全建议

### 生产环境部署

1. **使用反向代理**
   ```nginx
   # Nginx 配置示例
   server {
       listen 80;
       server_name your-domain.com;

       location / {
           proxy_pass http://localhost:8080;
           proxy_set_header Host $host;
           proxy_set_header X-Real-IP $remote_addr;
       }
   }
   ```

2. **添加 HTTPS**
   - 使用 Let's Encrypt 获取免费 SSL 证书
   - 配置 Nginx/Apache 处理 HTTPS

3. **限制访问**
   ```go
   // 添加 IP 白名单
   allowedIPs := []string{"127.0.0.1", "192.168.1.0/24"}
   ```

4. **设置防火墙**
   ```bash
   # 只允许本地访问
   ufw allow from 127.0.0.1 to any port 8080
   ```

## 🐳 Docker 部署

### Dockerfile

```dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY . .
RUN go build -o server server.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/server .
COPY index.html .
COPY jdkindex.json .

EXPOSE 8080
CMD ["./server"]
```

### 构建和运行

```bash
# 构建镜像
docker build -t openjdk-server .

# 运行容器
docker run -d -p 8080:8080 --name openjdk-web openjdk-server

# 查看日志
docker logs -f openjdk-web
```

## 📦 系统服务

### Linux (systemd)

创建服务文件 `/etc/systemd/system/openjdk-server.service`：

```ini
[Unit]
Description=OpenJDK Download Center Web Server
After=network.target

[Service]
Type=simple
User=www-data
WorkingDirectory=/opt/openjdk-server
ExecStart=/opt/openjdk-server/server
Restart=on-failure

[Install]
WantedBy=multi-user.target
```

启动服务：

```bash
sudo systemctl daemon-reload
sudo systemctl start openjdk-server
sudo systemctl enable openjdk-server
sudo systemctl status openjdk-server
```

### Windows (NSSM)

```powershell
# 下载 NSSM
# https://nssm.cc/download

# 安装服务
nssm install OpenJDKServer "C:\path\to\openjdk-server.exe"
nssm set OpenJDKServer AppDirectory "C:\path\to\project"
nssm start OpenJDKServer
```

## ❓ 常见问题

### Q: 端口已被占用

```
listen tcp :8080: bind: address already in use
```

**解决方案**：
1. 更改端口号（使用 PORT 环境变量）
2. 关闭占用端口的进程

```bash
# 查找占用端口的进程
netstat -ano | findstr :8080    # Windows
lsof -i :8080                   # Linux/macOS

# 结束进程
taskkill /PID <PID> /F          # Windows
kill -9 <PID>                   # Linux/macOS
```

### Q: 文件不存在错误

```
缺少必需文件: index.html
```

**解决方案**：
确保在项目目录下运行服务器，包含以下文件：
- index.html
- jdkindex.json

### Q: 无法从外网访问

**解决方案**：
1. 检查防火墙设置
2. 确保路由器端口转发配置正确
3. 使用 `0.0.0.0` 监听所有接口（默认已配置）

## 🎯 性能优化

### 启用 Gzip 压缩

```go
import "github.com/NYTimes/gziphandler"

func main() {
    handler := gziphandler.GzipHandler(http.HandlerFunc(serveFile))
    http.Handle("/", corsMiddleware(logRequest(handler)))
}
```

### 添加缓存

```go
// 已在代码中添加
w.Header().Set("Cache-Control", "public, max-age=3600")
```

### 使用连接池

```go
server := &http.Server{
    Addr:         addr,
    Handler:      nil,
    ReadTimeout:  10 * time.Second,
    WriteTimeout: 10 * time.Second,
    IdleTimeout:  60 * time.Second,
}
```

## 📚 API 端点

| 端点 | 方法 | 描述 |
|-----|------|------|
| `/` | GET | 返回 index.html 主页 |
| `/index.html` | GET | 返回 index.html |
| `/jdkindex.json` | GET | 返回 JSON 数据 |
| 任意静态文件 | GET | 返回对应文件 |

## 🔍 监控和日志

### 查看实时日志

```bash
# 运行时查看
go run server.go 2>&1 | tee server.log

# 只看错误日志
go run server.go 2>&1 | grep "ERROR"
```

### 日志轮转

使用 logrotate 管理日志文件（Linux）：

```
/var/log/openjdk-server.log {
    daily
    rotate 7
    compress
    delaycompress
    missingok
    notifempty
}
```

## 📞 技术支持

遇到问题？请检查：
1. Go 版本是否 >= 1.16
2. 必需文件是否存在
3. 端口是否被占用
4. 防火墙设置

---

**版本**: 1.0
**最后更新**: 2025-10-19
