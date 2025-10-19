package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/html"
)

// TOpenJDK 表示一个 OpenJDK 下载条目
type TOpenJDK struct {
	Version      string `json:"version"`       // 版本号
	Filename     string `json:"filename"`      // 文件名
	URL          string `json:"url"`           // 下载链接
	Size         string `json:"size"`          // 文件大小
	LastModified string `json:"last_modified"` // 最后修改时间
	GOOS         string `json:"goos"`          // 操作系统类型
	GOARCH       string `json:"goarch"`        // 系统架构
}

// String 返回 TOpenJDK 的字符串表示
// 返回:
//   - string: 格式化的 OpenJDK 信息字符串
func (jdk TOpenJDK) String() string {
	return fmt.Sprintf("OpenJDK %-10s | 系统: %-10s | 架构: %-10s | 大小: %-10s | 文件: %s",
		jdk.Version, jdk.GOOS, jdk.GOARCH, jdk.Size, jdk.Filename)
}

// 清华大学开源软件镜像站
type TWebTuna struct {
	BaseURL string // 入口地址
}

// 兰州大学开源软件镜像站
type TWebLzu struct {
	BaseURL string // 入口地址
}

// 华为软件镜像站
type TWebHuawei struct {
	BaseURL string // 入口地址
}

// injdk 软件镜像站
type TWebInjdk struct {
	BaseURL string // 入口地址
}

// TWebFileInfo 表示从 HTML 表格中解析出的文件信息
type TWebFileInfo struct {
	Name         string // 文件名
	LastModified string // 最后修改时间
	Size         string // 文件大小
}

// mapOSToGOOS 将操作系统目录名映射为 GOOS 标准名称
// 参数:
//   - osDir: 操作系统目录名（如 "windows/", "linux/"）
//
// 返回:
//   - string: GOOS 标准名称
func mapOSToGOOS(osDir string) string {
	osDir = strings.TrimSuffix(osDir, "/")
	osMap := map[string]string{
		"windows":      "windows",
		"linux":        "linux",
		"mac":          "darwin",
		"macos":        "darwin",
		"osx":          "darwin",
		"alpine-linux": "linux",
		"aix":          "aix",
		"solaris":      "solaris",
	}
	if goos, ok := osMap[osDir]; ok {
		return goos
	}
	return osDir
}

// mapArchToGOARCH 将架构目录名映射为 GOARCH 标准名称
// 参数:
//   - archDir: 架构目录名（如 "x64/", "aarch64/"）
//
// 返回:
//   - string: GOARCH 标准名称
func mapArchToGOARCH(archDir string) string {
	archDir = strings.TrimSuffix(archDir, "/")
	archMap := map[string]string{
		"x64":     "amd64",
		"x32":     "386",
		"aarch64": "arm64",
		"arm":     "arm",
		"ppc64":   "ppc64",
		"ppc64le": "ppc64le",
		"s390x":   "s390x",
		"riscv64": "riscv64",
		"sparcv9": "sparc64",
	}
	if goarch, ok := archMap[archDir]; ok {
		return goarch
	}
	return archDir
}

// extractVersion 从文件名中提取版本号
// 参数:
//   - filename: 文件名
//
// 返回:
//   - string: 版本号
func extractVersion(filename string) string {
	// 匹配版本号模式，例如 "8u462b08", "25", "21.0.1", "17.0.2+8" 等
	patterns := []string{
		`(\d+u\d+b\d+)`,        // 8u462b08
		`(\d+\.\d+\.\d+\+\d+)`, // 17.0.2+8
		`(\d+\.\d+\.\d+)`,      // 21.0.1
		`(\d+)`,                // 25
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		if matches := re.FindStringSubmatch(filename); len(matches) > 1 {
			return matches[1]
		}
	}

	return "unknown"
}

// parseTime 解析时间字符串并转换为标准格式
// 参数:
//   - timeStr: 时间字符串（如 "22 Jul 2025 15:07:40 +0000"）
//
// 返回:
//   - string: 格式化后的时间字符串（"2006-01-02 15:04:05"）
func parseTime(timeStr string) string {
	// 尝试解析多种时间格式
	formats := []string{
		"02 Jan 2006 15:04:05 -0700", // 22 Jul 2025 15:07:40 +0000
		"2006-Jan-02 15:04",
		"2006-01-02 15:04:05", // 已经是目标格式
		"2006-01-02 15:04",    // 已经是目标格式
		time.RFC3339,          // 2006-01-02T15:04:05Z07:00
	}
	for _, format := range formats {
		if t, err := time.Parse(format, timeStr); err == nil {
			// 转换为目标格式
			return t.Format("2006-01-02 15:04")
		}
	}

	// 如果无法解析，返回原始字符串
	return timeStr
}

// fetchHTML 从指定 URL 获取 HTML 内容
// 参数:
//   - url: 要获取的网页地址
//
// 返回:
//   - string: HTML 内容
//   - error: 错误信息
func fetchHTML(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("获取网页失败: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取网页内容失败: %v", err)
	}

	return string(body), nil
}

// parseLinks 解析 HTML 中的所有 <a> 标签链接
// 参数:
//   - htmlContent: HTML 内容字符串
//   - pattern: 正则表达式模式，用于过滤链接
//
// 返回:
//   - []string: 匹配的链接列表
//   - error: 错误信息
func parseLinks(htmlContent string, pattern string) ([]string, error) {
	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		return nil, fmt.Errorf("解析 HTML 失败: %v", err)
	}

	var links []string
	var re *regexp.Regexp
	if pattern != "" {
		re = regexp.MustCompile(pattern)
	}

	// 递归遍历 HTML 节点
	var traverse func(*html.Node)
	traverse = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					href := attr.Val
					// 如果有正则模式，进行匹配
					if re == nil || re.MatchString(href) {
						links = append(links, href)
					}
					break
				}
			}
		}
		// 遍历子节点
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			traverse(c)
		}
	}

	traverse(doc)
	return links, nil
}

// parseFileTable 解析 HTML 表格中的文件信息（包括文件名、大小、修改时间）
// 参数:
//   - htmlContent: HTML 内容字符串
//   - pattern: 正则表达式模式，用于过滤文件名
//
// 返回:
//   - []TWebFileInfo: 文件信息列表
//   - error: 错误信息
func parseFileTable(htmlContent string, pattern string) ([]TWebFileInfo, error) {
	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		return nil, fmt.Errorf("解析 HTML 失败: %v", err)
	}

	var files []TWebFileInfo
	var re *regexp.Regexp
	if pattern != "" {
		re = regexp.MustCompile(pattern)
	}

	// 递归遍历找到表格行
	var traverse func(*html.Node)
	traverse = func(n *html.Node) {
		// 找到 <tr> 标签
		if n.Type == html.ElementNode && n.Data == "tr" {
			var fileName, fileSize, fileDate string
			colIndex := 0

			// 遍历表格列 <td>
			for td := n.FirstChild; td != nil; td = td.NextSibling {
				if td.Type == html.ElementNode && td.Data == "td" {
					// 获取列的文本内容
					text := getTextContent(td)

					// 第一列通常是文件名（在 <a> 标签中）
					if colIndex == 0 {
						for a := td.FirstChild; a != nil; a = a.NextSibling {
							if a.Type == html.ElementNode && a.Data == "a" {
								for _, attr := range a.Attr {
									if attr.Key == "href" {
										fileName = attr.Val
										break
									}
								}
							}
						}
					} else if colIndex == 1 {
						// 第二列是文件大小
						fileSize = strings.TrimSpace(text)
					} else if colIndex == 2 {
						// 第三列是修改时间
						fileDate = strings.TrimSpace(text)
					}
					colIndex++
				}
			}

			// 如果找到文件名且匹配模式，添加到列表
			if fileName != "" && (re == nil || re.MatchString(fileName)) {
				files = append(files, TWebFileInfo{
					Name:         fileName,
					LastModified: fileDate,
					Size:         fileSize,
				})
			}
		}

		// 遍历子节点
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			traverse(c)
		}
	}

	traverse(doc)
	return files, nil
}

// getTextContent 获取节点的文本内容
// 参数:
//   - n: HTML 节点
//
// 返回:
//   - string: 文本内容
func getTextContent(n *html.Node) string {
	if n.Type == html.TextNode {
		return n.Data
	}
	var text string
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		text += getTextContent(c)
	}
	return text
}

// getVerDirs 获取所有版本目录（如 25.0.3/, 24/, 等）
// 返回:
//   - []string: 版本目录列表
//   - error: 错误信息
func getVerDirs(url string) ([]string, error) {
	fmt.Println("正在获取版本目录列表...")
	htmlContent, err := fetchHTML(url)
	if err != nil {
		return nil, err
	}

	// 匹配数字开头的目录，例如 "25/"
	links, err := parseLinks(htmlContent, `^\d+(?:\.\d+)*/$`)
	if err != nil {
		return nil, err
	}

	fmt.Printf("找到 %d 个版本目录\n", len(links))
	return links, nil
}

// getJDKDirectory 获取 jdk/ 目录
// 参数:
//   - versionURL: 版本目录的 URL
//
// 返回:
//   - string: jdk 目录路径
//   - error: 错误信息
func getJDKDirectory(versionURL string) (string, error) {
	htmlContent, err := fetchHTML(versionURL)
	if err != nil {
		return "", err
	}

	links, err := parseLinks(htmlContent, `^jdk/$`)
	if err != nil {
		return "", err
	}

	if len(links) == 0 {
		return "", fmt.Errorf("未找到 jdk 目录")
	}

	return links[0], nil
}

// getArchDirs 获取所有架构目录（如 x64/, aarch64/, 等）
// 参数:
//   - jdkURL: jdk 目录的 URL
//
// 返回:
//   - []string: 架构目录列表
//   - error: 错误信息
func getArchDirs(jdkURL string) ([]string, error) {
	htmlContent, err := fetchHTML(jdkURL)
	if err != nil {
		return nil, err
	}

	// 匹配以 / 结尾的目录
	links, err := parseLinks(htmlContent, `^[^/]+/$`)
	if err != nil {
		return nil, err
	}

	// 过滤掉父目录链接
	var archs []string
	for _, link := range links {
		if link != "../" {
			archs = append(archs, link)
		}
	}

	return archs, nil
}

func formatFileSize(filesize string) string {
	// 将字符串转为 float64
	bytes, err := strconv.ParseFloat(filesize, 64)
	if err != nil {
		return "0"
	}
	units := []string{"B", "KiB", "MiB", "GiB", "TiB", "PiB", "EiB"}

	if bytes < 1 {
		return "0 B"
	}

	// 计算适合的单位
	base := 1024.0
	exp := int(math.Log(bytes) / math.Log(base))
	if exp >= len(units) {
		exp = len(units) - 1
	}

	// 计算值并格式化
	value := bytes / math.Pow(base, float64(exp))

	// 根据值的大小决定小数位数
	var format string
	if value < 10 {
		format = "%.2f %s"
	} else if value < 100 {
		format = "%.1f %s"
	} else {
		format = "%.0f %s"
	}

	return fmt.Sprintf(format, value, units[exp])
}

// getOSDirs 获取所有操作系统目录（如 windows/, linux/, mac/, 等）
// 参数:
//   - archURL: 架构目录的 URL
//
// 返回:
//   - []string: 操作系统目录列表
//   - error: 错误信息
func getOSDirs(archURL string) ([]string, error) {
	htmlContent, err := fetchHTML(archURL)
	if err != nil {
		return nil, err
	}

	links, err := parseLinks(htmlContent, `^[^/]+/$`)
	if err != nil {
		return nil, err
	}

	// 过滤掉父目录链接
	var osDirs []string
	for _, link := range links {
		if link != "../" {
			osDirs = append(osDirs, link)
		}
	}

	return osDirs, nil
}

// saveToJSON 将下载列表保存为 JSON 文件
// 参数:
//   - downloads: JDK 下载条目列表
//   - filename: 输出文件名
//
// 返回:
//   - error: 错误信息
func saveToJSON(downloads []TOpenJDK, filename string) error {
	// 按版本和文件名排序
	sort.Slice(downloads, func(i, j int) bool {
		if downloads[i].Version != downloads[j].Version {
			return downloads[i].Version > downloads[j].Version
		}
		return downloads[i].Filename < downloads[j].Filename
	})

	// 格式化 JSON，使用缩进
	data, err := json.MarshalIndent(downloads, "", "  ")
	if err != nil {
		return fmt.Errorf("JSON 序列化失败: %v", err)
	}

	// 写入文件
	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		return fmt.Errorf("写入文件失败: %v", err)
	}

	fmt.Printf("\n成功保存 %d 个下载地址到 %s\n", len(downloads), filename)
	return nil
}

// ParseURL 爬取所有 清华大学 JDK 下载地址
// 目录层次:
//   - Adoptium/25/jdk/x64/windows/OpenJDK25U-jdk_x64_windows_hotspot_25_36.zip
//
// 返回:
//   - []TOpenJDK: 所有 JDK 下载条目
//   - error: 错误信息
func (s *TWebTuna) ParseURL() ([]TOpenJDK, error) {
	var allDownloads []TOpenJDK

	// 1. 获取版本目录
	versions, err := getVerDirs(s.BaseURL)
	if err != nil {
		return nil, err
	}

	// 2. 遍历每个版本
	for _, version := range versions {
		versionURL := s.BaseURL + version
		// 去除版本号中的 '/' 字符用于显示
		versionDisplay := strings.TrimSuffix(version, "/")
		fmt.Println("==================================================")
		fmt.Printf("-= 处理 JDK %s 版本 =-\n", versionDisplay)
		fmt.Println("==================================================")

		// 3. 获取 jdk 目录
		jdkDir, err := getJDKDirectory(versionURL)
		if err != nil {
			fmt.Printf("  获取jdk目录失败: %v\n", err)
			continue
		}
		jdkURL := versionURL + jdkDir

		// 4. 获取架构目录
		archs, err := getArchDirs(jdkURL)
		if err != nil {
			fmt.Printf("  获取架构目录失败: %v\n", err)
			continue
		}

		// 5. 遍历每个架构（只处理 amd64 和 arm64）
		for _, arch := range archs {
			// 过滤架构：只处理 x64(amd64) 和 aarch64(arm64)
			archName := strings.TrimSuffix(arch, "/")
			if archName != "x64" && archName != "aarch64" {
				continue
			}

			archURL := jdkURL + arch

			// 6. 获取操作系统目录
			osDirs, err := getOSDirs(archURL)
			if err != nil {
				fmt.Printf("      获取操作系统目录失败: %v\n", err)
				continue
			}

			// 7. 遍历每个操作系统（只处理 windows, darwin, linux）
			for _, osDir := range osDirs {
				// 过滤操作系统：只处理 windows, mac(darwin), linux
				goos := mapOSToGOOS(osDir)
				if goos != "windows" && goos != "darwin" && goos != "linux" {
					continue
				}

				osURL := archURL + osDir

				// 8. 获取 JDK 文件
				downloads, err := s.GetJDKFiles(osURL, osDir, arch)
				if err != nil {
					fmt.Printf("        获取文件失败: %v\n", err)
					continue
				}

				// 打印找到的 OpenJDK 信息
				if len(downloads) > 0 {
					for _, jdk := range downloads {
						fmt.Printf("%s\n", jdk.String())
					}
				}
				allDownloads = append(allDownloads, downloads...)
			}
		}
	}

	return allDownloads, nil
}

// getJDKFiles 获取指定 URL 中的所有 JDK 文件下载地址及详细信息
// 参数:
//   - fileURL: 文件列表页面的 URL
//   - osDir: 操作系统目录名
//   - archDir: 架构目录名
//
// 返回:
//   - []TOpenJDK: JDK 下载条目列表
//   - error: 错误信息
func (s *TWebTuna) GetJDKFiles(fileURL string, osDir string, archDir string) ([]TOpenJDK, error) {
	htmlContent, err := fetchHTML(fileURL)
	if err != nil {
		return nil, err
	}

	// 匹配 .zip 和 .tar.gz 文件
	files, err := parseFileTable(htmlContent, `\.(zip|tar\.gz)$`)
	if err != nil {
		return nil, err
	}

	var downloads []TOpenJDK
	for _, file := range files {
		// 从文件名中提取版本号
		version := extractVersion(file.Name)

		// 格式化时间
		formattedTime := parseTime(file.LastModified)

		downloads = append(downloads, TOpenJDK{
			Version:      version,
			Filename:     file.Name,
			URL:          fileURL + file.Name,
			Size:         file.Size,
			LastModified: formattedTime,
			GOOS:         mapOSToGOOS(osDir),
			GOARCH:       mapArchToGOARCH(archDir),
		})
	}

	return downloads, nil
}

// ParseURL 爬取所有 兰州大学 JDK 下载地址
// 目录层次:
//   - /openjdk/11.0.2/openjdk-11.0.2_windows-x64_bin.zip
//
// 返回:
//   - []TOpenJDK: 所有 JDK 下载条目
//   - error: 错误信息
func (s *TWebLzu) ParseURL() ([]TOpenJDK, error) {
	var allDownloads []TOpenJDK

	// 1. 获取版本目录
	versions, err := getVerDirs(s.BaseURL)
	if err != nil {
		return nil, err
	}

	// 2. 遍历每个版本
	for _, version := range versions {
		if !strings.Contains(version, "20.") {
			continue
		}
		versionURL := s.BaseURL + version
		// 去除版本号中的 '/' 字符用于显示
		versionDisplay := strings.TrimSuffix(version, "/")
		fmt.Println("==================================================")
		fmt.Printf("-= 处理 JDK %s 版本 =-\n", versionDisplay)
		fmt.Println("==================================================")

		// 8. 获取 JDK 文件
		downloads, err := s.GetJDKFiles(versionURL)
		if err != nil {
			fmt.Printf("        获取文件失败: %v\n", err)
			continue
		}

		// 打印找到的 OpenJDK 信息
		if len(downloads) > 0 {
			for _, jdk := range downloads {
				fmt.Printf("%s\n", jdk.String())
			}
		}
		allDownloads = append(allDownloads, downloads...)

	}

	return allDownloads, nil
}

func (s *TWebLzu) GetJDKFiles(fileURL string) ([]TOpenJDK, error) {
	htmlContent, err := fetchHTML(fileURL)
	if err != nil {
		return nil, err
	}

	// 匹配 .zip 和 .tar.gz 文件
	files, err := parseFileTable(htmlContent, `\.(zip|tar\.gz)$`)
	if err != nil {
		return nil, err
	}

	var downloads []TOpenJDK
	for _, file := range files {
		// 从文件名中提取版本号
		version := extractVersion(file.Name)

		// 格式化时间
		formattedTime := parseTime(file.LastModified)

		goos, goarch, _, err := s.ParseWebFileName(file.Name)
		if err != nil {
			fmt.Println(err)
			continue
		}

		downloads = append(downloads, TOpenJDK{
			Version:      version,
			Filename:     file.Name,
			URL:          fileURL + file.Name,
			Size:         formatFileSize(file.Size),
			LastModified: formattedTime,
			GOOS:         goos,
			GOARCH:       goarch,
		})
	}

	return downloads, nil
}

// 输入:openjdk-10.0.1_windows-x64_bin.tar.gz
// 返回:Goos,Arch,Version
func (s *TWebLzu) ParseWebFileName(filename string) (string, string, string, error) {
	filename = strings.TrimSpace(filename)
	if filename == "" {
		return "", "", "", fmt.Errorf("文件名不能为空")
	}

	// 解析文件名获取版本、GOOS和GOARCH
	filenameParts := strings.Split(filename, "_")
	if len(filenameParts) < 3 {
		return "", "", "", fmt.Errorf("文件名格式错误(%s)", filename)
	}

	version := strings.TrimPrefix(filenameParts[0], "openjdk-")
	filenameParts1 := strings.Split(filenameParts[1], "-")
	var goos, goarch string
	if len(filenameParts1) > 1 {
		goos = filenameParts1[0]
		goarch = filenameParts1[1]
	}

	// 标准化GOOS值
	goos = strings.ReplaceAll(goos, "osx", "darwin")
	goos = strings.ReplaceAll(goos, "macos", "darwin")
	switch goos {
	case "linux", "darwin", "windows":
		// 已经是标准值
	default:
		return "", "", "", fmt.Errorf("未知操作系统(%s)", goos)
	}

	// 标准化GOARCH值
	goarch = strings.ReplaceAll(goarch, "aarch64", "arm64")
	switch goarch {
	case "x64", "amd64", "arm64", "aarch64":
		if goarch == "x64" {
			goarch = "amd64" // Go标准中使用amd64而不是x64
		}
	default:
		return "", "", "", fmt.Errorf("未知系统架构(%s)", goarch)
	}

	return goos, goarch, version, nil
}

func main() {
	fmt.Println("====================================")
	fmt.Println("OpenJDK 下载地址爬取工具")
	fmt.Println("====================================")

	//WebJDK1 := TWebTuna{}
	//WebJDK1.BaseURL = "https://mirrors.tuna.tsinghua.edu.cn/Adoptium/"

	WebJDK2 := TWebLzu{}
	WebJDK2.BaseURL = "https://mirror4.lzu.edu.cn/openjdk/"

	// 爬取所有 JDK 下载地址
	downloads, err := WebJDK2.ParseURL()
	if err != nil {
		fmt.Printf("错误: %v\n", err)
		os.Exit(1)
	}

	// 保存为 JSON
	err = saveToJSON(downloads, "jdkindex.json")
	if err != nil {
		fmt.Printf("错误: %v\n", err)
		os.Exit(1)
	}

}
