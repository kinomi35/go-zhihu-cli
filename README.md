# go-zhihu-cli

一个基于 Cobra 和 Bubble Tea 的知乎命令行工具，支持cli和tui操作。参考了 `pyzhihu-cli` 项目。

当前支持：

- 使用浏览器 Cookie 登录。
- 查看知乎推荐流。
- 按回答 ID 阅读推荐回答。
- 在 TUI 界面中浏览推荐流、查看详情和展开评论。
- 将知乎 Web API 地址放在可编辑的 JSON 配置里，并在构建时内置默认配置。

## 构建

```bash
go build -o zhihu ./cmd/zhihu
```

## 登录

```bash
./zhihu login --cookies "z_c0=...; _xsrf=...; d_c0=..."
```

从已登录知乎的浏览器请求里复制 `Cookie` 请求头。传入的字符串至少需要包含 `z_c0`、`_xsrf` 和 `d_c0`。

Cookie 会保存到：

```text
~/.go-zhihu-cli/cookies.json
```

检查登录状态：

```bash
./zhihu status
```

退出登录：

```bash
./zhihu logout
```

## 推荐流

```bash
./zhihu feed -l 10
./zhihu feed --json
```

当推荐项是回答时，`feed` 命令会打印回答 ID：

```text
阅读: zhihu read <回答ID>
```

## 阅读回答

```bash
./zhihu read <回答ID>
./zhihu read <回答ID> --comments 5
./zhihu read <回答ID> --json
```

`read` 会输出问题标题、作者、回答正文、赞同数和评论数。传入 `--comments` 后会额外拉取并显示指定数量的评论。

## TUI 界面

```bash
./zhihu tui
```

TUI 会启动全屏终端界面，自动加载推荐流。列表和详情使用同一个已登录知乎客户端。

常用按键：

- `↑/↓` 或 `k/j`：移动或滚动。
- `enter`：打开当前推荐项详情。
- `r`：刷新推荐列表。
- `/`：筛选列表。
- `c`：在详情页展开评论。
- `z` 或 `esc`：返回上一页。
- `q`：退出 TUI。
- `?`：显示更多列表帮助。

详情页内容较长时可以用 `↑/↓` 滚动查看。展开评论后，正文和评论会放在同一个滚动视图里。

## 接口配置

默认的知乎 Web API 地址配置在：

```text
configs/endpoints.json
```

默认配置会通过 `go:embed` 编进二进制。运行时也可以指定其他配置文件：

```bash
./zhihu --endpoints ./configs/endpoints.json feed
./zhihu --endpoints ./configs/endpoints.json tui
```

这样接口地址变化时，不需要改命令实现代码。
