# memos-cli

一个用于访问 [Memos](https://www.usememos.com/) HTTP API 的 Go 命令行工具。

它提供常见的备忘录查询、创建、更新、删除、评论、标签处理和用户查询能力，适合在终端中直接操作 Memos，也适合作为脚本或自动化流程中的 CLI 工具使用。

## 功能概览

- 检查当前 CLI 配置是否完整
- 列出、查看、搜索、过滤备忘录
- 创建、更新、删除备忘录
- 为备忘录追加标签并移除标签
- 为指定备忘录创建评论
- 通过管理员 API 获取用户列表
- 支持 `--json` 输出，便于脚本消费
- 支持分页查询 memo 列表

## 安装

需要 Go `1.25` 或更高版本。

```bash
go install github.com/rogeecn/memos-cli@latest
```

安装后可直接使用：

```bash
memos --help
```

如果你是在当前仓库内本地运行，也可以使用：

```bash
go run . --help
```

## 配置

CLI 从以下位置按优先级读取配置：

1. 当前 shell 环境变量
2. 当前工作目录下的 `.env` 文件
3. 未配置

也就是说：如果环境变量和 `.env` 同时存在，环境变量优先。

### 必要变量

```env
MEMOS_URL=http://localhost:5230
MEMOS_API_KEY=your-memos-api-key
```

### 可选变量

```env
MEMOS_ADMIN_API_KEY=your-memos-admin-api-key
DEFAULT_TAG=cli
```

变量说明：

- `MEMOS_URL`：Memos 实例地址，例如 `http://localhost:5230`
- `MEMOS_API_KEY`：普通 API Key，用于大多数 memo 操作
- `MEMOS_ADMIN_API_KEY`：管理员 API Key，仅 `user list` 等管理员接口需要
- `DEFAULT_TAG`：创建 memo 时自动追加的默认标签

可参考仓库中的示例文件： [`.env.example`](.env.example)

### 检查配置

建议在首次使用前先检查配置状态：

```bash
memos config check
```

输出会显示各变量是 `configured` 还是 `missing`，不会打印密钥原文。

## 快速开始

```bash
memos config check
memos memo list
memos memo get <memo-id>
```

创建一条 memo：

```bash
memos memo create "开始使用 memos-cli"
```

以 JSON 格式输出列表：

```bash
memos --json memo list
```

## 命令总览

### 配置

```bash
memos config check
```

### Memo 读取与查询

```bash
memos memo list
memos memo list --page-size 20
memos memo list --page-token <next-page-token>
memos memo get <memo-id>
memos search "keyword"
memos filter --expr "visibility == 'PRIVATE'"
```

### Memo 写入与维护

```bash
memos memo create <content>
memos memo create <content> --visibility PUBLIC
memos memo create <content> --tag release --tag cli
memos memo update <memo-id> --content "updated content"
memos memo update <memo-id> --visibility PUBLIC
memos memo delete <memo-id> --yes
```

### 评论

```bash
memos comment create <memo-id> <content>
```

### 标签

```bash
memos tag remove <memo-id> <tag>
```

### 用户

```bash
memos user list
```

`memos user list` 依赖 `MEMOS_ADMIN_API_KEY`。

## 输出说明

### 默认输出

默认情况下，CLI 以适合终端阅读的文本格式输出。

例如：

```bash
memos memo list
memos memo get <memo-id>
```

### JSON 输出

加上全局参数 `--json` 后，CLI 会输出原始结构化 JSON：

```bash
memos --json memo list
memos --json memo get <memo-id>
memos --json user list
```

当你需要把结果交给脚本、`jq` 或其他自动化工具处理时，推荐使用 `--json`。

### 分页

`memo list` 支持分页：

```bash
memos memo list --page-size 20
memos memo list --page-size 20 --page-token <next-page-token>
```

- 文本模式下，如果还有下一页，会输出 `Next page token: ...`
- JSON 模式下，可读取响应中的 `nextPageToken`

## 使用示例

### 列出 memo

```bash
memos memo list
```

### 获取单条 memo

```bash
memos memo get abc123
```

### 搜索包含关键词的 memo

```bash
memos search "项目复盘"
```

### 使用 CEL 表达式过滤 memo

```bash
memos filter --expr "createTime > timestamp('2026-01-01T00:00:00Z') && visibility == 'PRIVATE'"
```

### 创建带标签的 memo

```bash
memos memo create "发布检查已完成" --tag release --tag weekly
```

如果配置了 `DEFAULT_TAG`，CLI 会在显式 `--tag` 之外自动追加默认标签。

### 更新 memo 内容

```bash
memos memo update abc123 --content "已更新的内容"
```

### 删除 memo

```bash
memos memo delete abc123 --yes
```

删除操作必须显式传入 `--yes`。

### 创建评论

```bash
memos comment create abc123 "这条内容已经复核完成"
```

### 移除标签

```bash
memos tag remove abc123 weekly
```

### 查询用户列表

```bash
memos user list
```

## ID 与行为说明

- 大多数命令优先使用纯 memo ID，例如 `abc123`
- `comment create` 和 `tag remove` 会使用仓库内已实现的 ID 规范化逻辑
- `memo update` 至少需要传一个更新项：`--content` 或 `--visibility`
- `user list` 需要管理员 API Key

## 开发

仓库入口：

- `main.go`
- `internal/cli/`
- `internal/memos/`
- `internal/output/`
- `internal/config/`

本地查看帮助：

```bash
go run . --help
go run . memo --help
```

运行测试：

```bash
go test ./...
```

## 许可证

MIT
