package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

const (
	defaultPort = "8080" // 默认端口号
)

// serveFile 提供静态文件服务
// 参数:
//   - w: HTTP 响应写入器
//   - r: HTTP 请求对象
func serveFile(w http.ResponseWriter, r *http.Request) {
	// 获取请求的文件路径
	filePath := r.URL.Path
	if filePath == "/" {
		filePath = "/index.html"
	}

	// 移除开头的斜杠
	filePath = filePath[1:]

	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.Error(w, "文件不存在", http.StatusNotFound)
		log.Printf("404 - 文件不存在: %s", filePath)
		return
	}

	// 设置正确的 Content-Type
	ext := filepath.Ext(filePath)
	contentType := getContentType(ext)
	w.Header().Set("Content-Type", contentType)

	// 设置缓存控制
	w.Header().Set("Cache-Control", "public, max-age=3600")

	// 提供文件服务
	http.ServeFile(w, r, filePath)
	log.Printf("200 - %s %s", r.Method, r.URL.Path)
}

// getContentType 根据文件扩展名返回 MIME 类型
// 参数:
//   - ext: 文件扩展名（如 ".html", ".json"）
//
// 返回:
//   - string: MIME 类型
func getContentType(ext string) string {
	contentTypes := map[string]string{
		".html": "text/html; charset=utf-8",
		".css":  "text/css; charset=utf-8",
		".js":   "application/javascript; charset=utf-8",
		".json": "application/json; charset=utf-8",
		".png":  "image/png",
		".jpg":  "image/jpeg",
		".jpeg": "image/jpeg",
		".gif":  "image/gif",
		".svg":  "image/svg+xml",
		".ico":  "image/x-icon",
		".txt":  "text/plain; charset=utf-8",
		".xml":  "application/xml; charset=utf-8",
	}

	if contentType, ok := contentTypes[ext]; ok {
		return contentType
	}
	return "application/octet-stream"
}

// logRequest 记录请求日志的中间件
// 参数:
//   - next: 下一个处理器
//
// 返回:
//   - http.HandlerFunc: 包装后的处理器
func logRequest(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("收到请求: %s %s 来自 %s", r.Method, r.URL.Path, r.RemoteAddr)
		next(w, r)
	}
}

// corsMiddleware 添加 CORS 头部的中间件
// 参数:
//   - next: 下一个处理器
//
// 返回:
//   - http.HandlerFunc: 包装后的处理器
func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 设置 CORS 头部
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// 处理 OPTIONS 预检请求
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}

// printBanner 打印启动横幅
// 参数:
//   - port: 服务器端口号
func printBanner(port string) {
	banner := `
╔═══════════════════════════════════════════════════════╗
║                                                       ║
║       ☕ OpenJDK 下载中心 - Web 服务器               ║
║                                                       ║
╚═══════════════════════════════════════════════════════╝

🚀 服务器启动成功！

📍 访问地址:
   - 本地访问: http://localhost:%s
   - 网络访问: http://0.0.0.0:%s

📁 服务文件:
   - index.html  (主页面)
   - jdkindex.json (数据文件)

⌨️  按 Ctrl+C 停止服务器

════════════════════════════════════════════════════════
`
	fmt.Printf(banner, port, port)
}

// checkRequiredFiles 检查必需的文件是否存在
// 返回:
//   - error: 如果文件缺失则返回错误
func checkRequiredFiles() error {
	requiredFiles := []string{"index.html", "jdkindex.json"}

	for _, file := range requiredFiles {
		if _, err := os.Stat(file); os.IsNotExist(err) {
			return fmt.Errorf("缺少必需文件: %s", file)
		}
	}

	log.Println("✅ 所有必需文件检查通过")
	return nil
}

func webserver() {
	// 获取端口号（可以从环境变量或使用默认值）
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	// 检查必需文件
	if err := checkRequiredFiles(); err != nil {
		log.Fatalf("❌ 启动失败: %v", err)
	}

	// 设置路由
	http.HandleFunc("/", corsMiddleware(logRequest(serveFile)))

	// 打印启动信息
	printBanner(port)

	// 启动服务器
	addr := ":" + port
	log.Printf("🌐 HTTP 服务器正在监听端口 %s...\n", port)

	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("❌ 服务器启动失败: %v", err)
	}
}
