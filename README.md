# OpenJDK 下载地址索引

本文件由程序自动生成，包含从 Adoptium 镜像站抓取的 OpenJDK 下载地址。

## 使用说明

1. 运行程序: `go run main.go`
2. 程序会自动抓取所有 JDK 版本的下载地址
3. 结果保存在 `jdkindex.json` 文件中

## 数据来源

- 镜像站: https://mirrors.tuna.tsinghua.edu.cn/Adoptium/

## 统计信息

- 总计下载地址数量: 64
- 抓取时间: 自动生成

## JSON 文件格式

```json
[
  {
    "version": "8u462b08",
    "filename": "OpenJDK8U-jdk_x64_windows_hotspot_8u462b08.zip",
    "url": "https://mirrors.tuna.tsinghua.edu.cn/Adoptium/8/jdk/x64/windows/OpenJDK8U-jdk_x64_windows_hotspot_8u462b08.zip",
    "size": "105.2M",
    "last_modified": "2024-10-22 08:30:00",
    "goos": "windows",
    "goarch": "amd64"
  }
]
```

## 字段说明

- **version**: JDK 版本号（如 8u462b08、17、21 等）
- **filename**: 文件名
- **url**: 完整下载链接
- **size**: 文件大小
- **last_modified**: 最后修改时间
- **goos**: 操作系统类型（windows, linux, darwin, aix, solaris）
- **goarch**: 系统架构（amd64, 386, arm64, arm, ppc64, ppc64le, s390x, riscv64, sparc64）

## 支持的平台

### 操作系统
- Windows
- Linux
- macOS (darwin)
- AIX
- Solaris
- Alpine Linux

### 架构
- amd64 (x64)
- 386 (x32)
- arm64 (aarch64)
- arm
- ppc64
- ppc64le
- s390x
- riscv64
- sparc64

## 许可证

本程序遵循 MIT 许可证。OpenJDK 本身遵循 GPL v2 + Classpath Exception 许可证。
