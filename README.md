# 萌娘百科 CLI 工具

一个简单的命令行工具，用于搜索和查看萌娘百科的内容。

## 功能

- 搜索萌娘百科页面
- 查看页面内容（按标题或页面ID）
- 将完整页面内容保存到本地文件

## 安装

首先确保您已安装 Go 1.18 或更高版本。

克隆或下载此项目：

```bash
git clone <repository-url>
cd moegirlwiki
```

## 使用方法

### 搜索页面

```bash
# 搜索关键词，返回最多10个结果
go run main.go search <关键词>

# 搜索关键词，返回指定数量的结果
go run main.go search <关键词> <数量>
```

### 查看页面

```bash
# 通过标题查看页面
go run main.go view <页面标题>

# 通过页面ID查看页面
go run main.go viewid <页面ID>
```

## 示例

```bash
# 搜索"萌娘"相关页面，返回3个结果
go run main.go search 萌娘 3

# 查看标题为"萌娘"的页面
go run main.go view 萌娘

# 查看页面ID为23528的页面
go run main.go viewid 23528
```

## 特性

- 搜索功能：支持关键词搜索并显示结果摘要
- 页面查看：支持通过标题或ID查看页面内容
- 本地保存：完整页面内容会自动保存为本地文本文件，方便查阅
- 命令行友好：简洁的命令行界面

## 依赖

- Go 1.18+
- 互联网连接以访问萌娘百科API

## 文件说明

- `main.go`: 主程序文件
- `go.mod`: Go模块定义文件
- `<页面标题>.txt`: 保存的页面内容文件（当使用view功能时自动生成）
