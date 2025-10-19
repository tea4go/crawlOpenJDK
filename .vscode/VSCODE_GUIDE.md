# VSCode 配置使用指南

本项目已配置完整的 VSCode 开发环境，包括调试、任务和代码格式化等功能。

## 📁 配置文件说明

### 1. `launch.json` - 调试配置

包含 5 个调试配置：

#### 🔧 可用的调试配置

1. **运行爬虫程序**
   - 调试运行 `main.go`
   - 爬取 OpenJDK 下载地址并生成 JSON 文件

2. **运行 Web 服务器**
   - 调试运行 `server.go`
   - 默认端口：8080

3. **运行 Web 服务器 (自定义端口)**
   - 调试运行 `server.go`
   - 自定义端口：3000

4. **调试当前文件**
   - 调试当前打开的 Go 文件

5. **附加到进程**
   - 附加调试器到正在运行的 Go 进程

### 2. `tasks.json` - 任务配置

包含 12 个预定义任务：

#### 🚀 运行任务

- **运行爬虫程序**: `Ctrl+Shift+P` → `Tasks: Run Task` → `运行爬虫程序`
- **运行 Web 服务器**: `Ctrl+Shift+P` → `Tasks: Run Task` → `运行 Web 服务器`

#### 🔨 编译任务

- **编译爬虫程序**: 编译 `main.go` 到 `bin/wget_jdk.exe`
- **编译 Web 服务器**: 编译 `server.go` 到 `bin/server.exe`
- **编译所有程序**: 编译所有 Go 程序
- **清理构建文件**: 删除 `bin/` 和 `tmp/` 目录

#### ✅ 测试和检查任务

- **运行测试**: `go test -v ./...`
- **代码检查**: `go vet ./...`
- **格式化代码**: `go fmt ./...`
- **安装依赖**: `go mod tidy`

#### 🌐 Web 相关任务

- **在浏览器中打开**: 自动打开 http://localhost:8080
- **完整工作流**: 爬取数据 → 启动服务器 → 打开浏览器

### 3. `settings.json` - 项目设置

配置了以下内容：

- **Go 语言**: 启用语言服务器、自动格式化、自动导入整理
- **编辑器**: Tab 大小、保存时格式化
- **HTML/JSON/Markdown**: 各自的格式化设置
- **文件排除**: 排除 `bin/`, `tmp/` 等目录

### 4. `extensions.json` - 推荐扩展

推荐安装的 VSCode 扩展：

- `golang.go` - Go 语言支持
- `streetsidesoftware.code-spell-checker` - 拼写检查
- `eamodio.gitlens` - Git 增强
- `ritwickdey.liveserver` - 实时预览 HTML
- 等等...

## 🎯 快速开始

### 方法一：使用调试面板

1. 按 `F5` 或点击左侧调试图标
2. 在下拉菜单中选择配置：
   - `运行爬虫程序` - 爬取数据
   - `运行 Web 服务器` - 启动服务器
3. 点击绿色播放按钮开始调试

### 方法二：使用任务

1. 按 `Ctrl+Shift+P` 打开命令面板
2. 输入 `Tasks: Run Task`
3. 选择要执行的任务

### 方法三：使用快捷键

- `Ctrl+Shift+B` - 运行默认构建任务（编译爬虫程序）
- `Ctrl+Shift+T` - 运行默认测试任务
- `F5` - 启动调试

## 🔍 调试技巧

### 设置断点

1. 在代码行号左侧点击，设置红色断点
2. 按 `F5` 启动调试
3. 程序会在断点处暂停

### 调试控制

- `F5` - 继续执行
- `F10` - 单步跳过
- `F11` - 单步进入
- `Shift+F11` - 单步跳出
- `Ctrl+Shift+F5` - 重新启动
- `Shift+F5` - 停止调试

### 查看变量

- **悬停**: 鼠标悬停在变量上查看值
- **变量面板**: 左侧调试面板的"变量"标签
- **监视**: 添加表达式到"监视"面板
- **调试控制台**: 输入表达式查看值

## 📝 常用任务工作流

### 开发工作流

```bash
# 1. 编辑代码
# 2. 保存时自动格式化
# 3. 运行任务：代码检查
Ctrl+Shift+P → Tasks: Run Task → 代码检查

# 4. 运行调试
F5 → 选择"运行爬虫程序"
```

### 完整构建流程

```bash
# 1. 安装依赖
Ctrl+Shift+P → Tasks: Run Task → 安装依赖

# 2. 格式化代码
Ctrl+Shift+P → Tasks: Run Task → 格式化代码

# 3. 代码检查
Ctrl+Shift+P → Tasks: Run Task → 代码检查

# 4. 编译所有程序
Ctrl+Shift+P → Tasks: Run Task → 编译所有程序
```

### Web 开发流程

```bash
# 方法一：分步执行
1. F5 → 运行爬虫程序 (生成 jdkindex.json)
2. F5 → 运行 Web 服务器
3. Ctrl+Shift+P → Tasks: Run Task → 在浏览器中打开

# 方法二：一键执行
Ctrl+Shift+P → Tasks: Run Task → 完整工作流
```

## ⚙️ 自定义配置

### 修改端口

编辑 `.vscode/launch.json`：

```json
{
    "name": "运行 Web 服务器",
    "env": {
        "PORT": "3000"  // 修改为你想要的端口
    }
}
```

### 添加新任务

编辑 `.vscode/tasks.json`：

```json
{
    "label": "我的自定义任务",
    "type": "shell",
    "command": "你的命令",
    "problemMatcher": []
}
```

### 添加新调试配置

编辑 `.vscode/launch.json`：

```json
{
    "name": "我的调试配置",
    "type": "go",
    "request": "launch",
    "mode": "debug",
    "program": "${workspaceFolder}/your_file.go"
}
```

## 🐛 故障排除

### Go 扩展未安装

```bash
Ctrl+Shift+P → Go: Install/Update Tools
选择所有工具 → 点击 OK
```

### 代码格式化失败

```bash
# 检查 gofmt 是否安装
go version

# 手动格式化
go fmt ./...
```

### 调试器无法启动

```bash
# 安装 delve 调试器
go install github.com/go-delve/delve/cmd/dlv@latest

# 验证安装
dlv version
```

### 任务执行失败

1. 检查终端输出的错误信息
2. 确保在项目根目录下执行
3. 检查 Go 环境是否正确配置

## 🔗 快捷键速查

| 功能 | Windows/Linux | macOS |
|------|--------------|-------|
| 开始调试 | `F5` | `F5` |
| 运行（不调试） | `Ctrl+F5` | `Cmd+F5` |
| 停止调试 | `Shift+F5` | `Shift+F5` |
| 重启调试 | `Ctrl+Shift+F5` | `Cmd+Shift+F5` |
| 单步跳过 | `F10` | `F10` |
| 单步进入 | `F11` | `F11` |
| 单步跳出 | `Shift+F11` | `Shift+F11` |
| 切换断点 | `F9` | `F9` |
| 运行任务 | `Ctrl+Shift+P` | `Cmd+Shift+P` |
| 构建 | `Ctrl+Shift+B` | `Cmd+Shift+B` |

## 📚 更多资源

- [VSCode Go 扩展文档](https://github.com/golang/vscode-go/wiki)
- [VSCode 调试文档](https://code.visualstudio.com/docs/editor/debugging)
- [Go 语言官方文档](https://golang.org/doc/)

---

**提示**: 第一次打开项目时，VSCode 会提示安装推荐的扩展，建议全部安装以获得最佳开发体验。
