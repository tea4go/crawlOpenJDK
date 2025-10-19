# OpenJDK 下载地址爬取工具 - 软件需求规格说明书 (SRS)

## 项目概述

开发一个 Golang 程序，自动从 Adoptium 镜像站爬取 OpenJDK 的下载地址，解析文件详细信息并保存为 JSON 格式。

## 核心需求

### 1. 网页爬取与解析

#### 1.1 目标网站
- **基础 URL**: `https://mirrors.tuna.tsinghua.edu.cn/Adoptium/`
- **网页结构**: 多层级目录结构

#### 1.2 爬取流程（5层目录结构）

**第一层：版本目录**
- 分析首页，解析出所有版本目录链接
- 匹配模式：`<a href="25/" title="25">25/</a>`
- 正则表达式：`^\d+/$`（数字开头，以 `/` 结尾）

**第二层：JDK 目录**
- 进入每个版本目录
- 查找并进入 `jdk/` 目录
- 匹配模式：`<a href="jdk/" title="jdk">jdk/</a>`

**第三层：架构目录**
- 查找所有架构目录
- 示例：`<a href="x64/" title="x64">x64/</a>`, `<a href="aarch64/" title="aarch64">aarch64/</a>`
- **过滤条件**：只处理 `x64` 和 `aarch64` 两种架构

**第四层：操作系统目录**
- 查找所有操作系统目录
- 示例：`<a href="windows/" title="windows">windows/</a>`, `<a href="linux/" title="linux">linux/</a>`, `<a href="mac/" title="mac">mac/</a>`
- **过滤条件**：只处理 `windows`, `mac`(映射为 darwin), `linux` 三种系统
- 注意：排除 `alpine-linux`, `aix`, `solaris` 等其他系统

**第五层：文件列表**
- 解析 HTML 表格，提取文件信息
- 需要提取的信息：
  - 文件名
  - 文件大小
  - 最后修改时间
- **文件类型过滤**：只处理 `.tar.gz` 和 `.zip` 文件
- **排除的文件类型**：`.msi`, `.pkg`

### 2. 数据结构定义

#### 2.1 主要结构体

```go
type TOpenJDK struct {
    Version      string `json:"version"`       // 版本号
    Filename     string `json:"filename"`      // 文件名
    URL          string `json:"url"`           // 下载链接
    Size         string `json:"size"`          // 文件大小
    LastModified string `json:"last_modified"` // 最后修改时间
    GOOS         string `json:"goos"`          // 操作系统类型
    GOARCH       string `json:"goarch"`        // 系统架构
}
```

#### 2.2 辅助结构体

```go
type TWebFileInfo struct {
    Name         string // 文件名
    LastModified string // 最后修改时间
    Size         string // 文件大小
}
```

### 3. 数据提取与转换规则

#### 3.1 版本号提取
从文件名中智能提取版本号，支持多种格式：
- `8u462b08` - 匹配模式：`(\d+u\d+b\d+)`
- `17.0.2+8` - 匹配模式：`(\d+\.\d+\.\d+\+\d+)`
- `21.0.1` - 匹配模式：`(\d+\.\d+\.\d+)`
- `25` - 匹配模式：`(\d+)`

#### 3.2 操作系统映射（GOOS）
将目录名映射为 Go 标准的 GOOS 名称：
- `windows` → `windows`
- `linux` → `linux`
- `mac` → `darwin`
- `alpine-linux` → `linux`（但此类型已被过滤）
- `aix` → `aix`（已被过滤）
- `solaris` → `solaris`（已被过滤）

#### 3.3 架构映射（GOARCH）
将目录名映射为 Go 标准的 GOARCH 名称：
- `x64` → `amd64`
- `aarch64` → `arm64`
- `x32` → `386`（已被过滤）
- `arm` → `arm`（已被过滤）
- `ppc64` → `ppc64`（已被过滤）
- `ppc64le` → `ppc64le`（已被过滤）
- `s390x` → `s390x`（已被过滤）
- `riscv64` → `riscv64`（已被过滤）
- `sparcv9` → `sparc64`（已被过滤）

#### 3.4 时间格式转换
- **输入格式**：`22 Jul 2025 15:07:40 +0000` (HTTP 日期格式)
- **输出格式**：`2025-07-22 15:07:05` (Go 时间格式)
- **Go 时间格式模板**：`2006-01-02 15:04:05`

### 4. 过滤规则总结

#### 4.1 架构过滤
**保留**：
- `x64` (映射为 amd64)
- `aarch64` (映射为 arm64)

**排除**：
- `x32` (386)
- `arm`
- `ppc64`
- `ppc64le`
- `s390x`
- `riscv64`
- `sparcv9`

#### 4.2 操作系统过滤
**保留**：
- `windows`
- `mac` (映射为 darwin)
- `linux`

**排除**：
- `alpine-linux`
- `aix`
- `solaris`

#### 4.3 文件类型过滤
**保留**：
- `.tar.gz`
- `.zip`

**排除**：
- `.msi`
- `.pkg`

### 5. 输出要求

#### 5.1 JSON 文件（jdkindex.json）

**格式要求**：
- 使用 2 空格缩进
- 按版本和文件名排序
- 字段顺序固定

**示例**：
```json
[
  {
    "version": "8u462b08",
    "filename": "OpenJDK8U-jdk_x64_windows_hotspot_8u462b08.zip",
    "url": "https://mirrors.tuna.tsinghua.edu.cn/Adoptium/8/jdk/x64/windows/OpenJDK8U-jdk_x64_windows_hotspot_8u462b08.zip",
    "size": "102.0 MiB",
    "last_modified": "2025-07-28 18:11:22",
    "goos": "windows",
    "goarch": "amd64"
  }
]
```

#### 5.2 README.md 文档

**必须包含的章节**：
1. 项目标题和简介
2. 使用说明
3. 数据来源
4. 统计信息（总数量、抓取时间）
5. JSON 文件格式示例
6. 字段说明
7. 支持的平台列表（操作系统和架构）
8. 许可证信息

### 6. 程序功能要求

#### 6.1 TOpenJDK 结构体方法

**String() 方法**：
- 返回格式化的字符串表示
- 格式：`OpenJDK %-10s | 系统: %-10s | 架构: %-10s | 大小: %-10s | 文件: %s`
- 使用左对齐，宽度为 10 的格式化

示例输出：
```
OpenJDK 8u462b08   | 系统: windows    | 架构: amd64      | 大小: 102.0 MiB  | 文件: OpenJDK8U-jdk_x64_windows_hotspot_8u462b08.zip
```

#### 6.2 日志输出要求

**版本处理日志**：
```
==================================================
-= 处理 JDK 11 版本 =-
==================================================
```

**注意**：
- 版本号显示时去除尾部的 `/` 字符
- 使用 `strings.TrimSuffix(version, "/")` 处理

**文件发现日志**：
- 每找到一个文件，立即打印详细信息
- 使用 `jdk.String()` 方法输出
- 不再显示"找到 X 个文件"的汇总信息

**架构和系统日志**：
- 架构显示：去除尾部 `/`
- 系统显示：去除尾部 `/`

### 7. 代码质量要求

#### 7.1 注释要求

**所有函数必须包含**：
1. 功能说明（中文）
2. 参数说明（格式：`// 参数:`，每个参数单独一行，包含中文说明）
3. 返回值说明（格式：`// 返回:`，每个返回值单独一行，包含中文说明）

**示例**：
```go
// mapOSToGOOS 将操作系统目录名映射为 GOOS 标准名称
// 参数:
//   - osDir: 操作系统目录名（如 "windows/", "linux/"）
//
// 返回:
//   - string: GOOS 标准名称
func mapOSToGOOS(osDir string) string {
    // 实现代码
}
```

#### 7.2 第三方组件注释

对使用的第三方包和函数添加中文注释：
- `golang.org/x/net/html` - HTML 解析库
- `html.Parse()` - 解析 HTML 字符串为 DOM 树
- `time.Parse()` - 解析时间字符串
- `json.MarshalIndent()` - 格式化 JSON 序列化

### 8. 自动执行要求

程序必须**自动执行**以下操作，不需要用户确认：

1. `go mod init wget_jdk` - 初始化 Go 模块
2. `go mod tidy` - 下载并整理依赖
3. `go run main.go` - 运行程序
4. 生成所有输出文件（jdkindex.json, README.md）

### 9. 依赖包

**必需的依赖**：
```go
import (
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "os"
    "regexp"
    "sort"
    "strings"
    "time"

    "golang.org/x/net/html"
)
```

### 10. 核心函数列表

#### 10.1 映射函数
- `mapOSToGOOS(osDir string) string` - 操作系统映射
- `mapArchToGOARCH(archDir string) string` - 架构映射
- `extractVersion(filename string) string` - 版本号提取
- `parseTime(timeStr string) string` - 时间格式转换

#### 10.2 HTML 解析函数
- `fetchHTML(url string) (string, error)` - 获取网页内容
- `parseLinks(htmlContent string, pattern string) ([]string, error)` - 解析链接
- `parseFileTable(htmlContent string, pattern string) ([]TWebFileInfo, error)` - 解析文件表格
- `getTextContent(n *html.Node) string` - 获取节点文本

#### 10.3 爬取函数
- `getVersionDirectories() ([]string, error)` - 获取版本目录
- `getJDKDirectory(versionURL string) (string, error)` - 获取 jdk 目录
- `getArchitectureDirectories(jdkURL string) ([]string, error)` - 获取架构目录
- `getOSDirectories(archURL string) ([]string, error)` - 获取操作系统目录
- `getJDKFiles(fileURL string, osDir string, archDir string) ([]TOpenJDK, error)` - 获取文件列表

#### 10.4 主控函数
- `crawlAllJDKs() ([]TOpenJDK, error)` - 爬取所有 JDK 下载地址
- `saveToJSON(downloads []TOpenJDK, filename string) error` - 保存为 JSON
- `generateReadme(downloads []TOpenJDK) error` - 生成 README
- `main()` - 主函数

### 11. HTML 解析细节

#### 11.1 表格结构
HTML 表格包含三列：
1. **第一列（colIndex=0）**：文件名（在 `<a>` 标签的 href 属性中）
2. **第二列（colIndex=1）**：文件大小
3. **第三列（colIndex=2）**：最后修改时间

#### 11.2 解析逻辑
```go
// 遍历 <tr> 标签
for each <tr>:
    for each <td>:
        if colIndex == 0:
            从 <a> 标签的 href 属性获取文件名
        else if colIndex == 1:
            获取文件大小
        else if colIndex == 2:
            获取最后修改时间
```

### 12. 错误处理

#### 12.1 日志输出
- 跳过的版本：`fmt.Printf("  跳过版本 %s: %v\n", versionDisplay, err)`
- 获取失败：`fmt.Printf("      获取操作系统目录失败: %v\n", err)`
- 文件获取失败：`fmt.Printf("        获取文件失败: %v\n", err)`

#### 12.2 继续执行
- 遇到错误时使用 `continue` 继续处理下一个项目
- 不中断整个爬取流程

### 13. 排序规则

JSON 输出排序：
```go
sort.Slice(downloads, func(i, j int) bool {
    if downloads[i].Version != downloads[j].Version {
        return downloads[i].Version > downloads[j].Version
    }
    return downloads[i].Filename < downloads[j].Filename
})
```

优先按版本号降序，版本相同时按文件名升序。

### 14. 预期输出结果

基于当前需求，程序应该：
- 爬取 8 个主要版本（8, 11, 17, 18, 19, 20, 21, 25）
- 每个版本包含 2 种架构（amd64, arm64）
- 每种架构包含 3 种操作系统（windows, darwin, linux）
- 每个组合包含 1-2 个文件（.tar.gz 和/或 .zip）
- **预期总数**：约 50 个下载地址

## 15. 验收标准

程序成功运行后应该：
- ✅ 生成 `jdkindex.json` 文件，包含约 50 个条目
- ✅ 生成 `README.md` 文件
- ✅ JSON 中只包含 amd64/arm64 架构
- ✅ JSON 中只包含 windows/darwin/linux 系统
- ✅ JSON 中只包含 .tar.gz/.zip 文件
- ✅ 所有时间格式为 `2006-01-02 15:04:05`
- ✅ 所有函数有完整的中文注释
- ✅ 日志输出格式正确，版本号无尾部斜杠

## 16. 项目文件清单

最终项目应包含：
1. `main.go` - 主程序文件
2. `go.mod` - Go 模块定义
3. `go.sum` - 依赖校验文件
4. `jdkindex.json` - 输出的 JSON 数据（运行后生成）
5. `README.md` - 项目文档（运行后生成）
6. `SRS.md` - 本需求文档（可选）
