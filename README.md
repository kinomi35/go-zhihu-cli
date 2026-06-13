# go-zhihu-cli

一个基于 Cobra 的知乎命令行工具，参考了上级目录里的 Python `zhihu-cli` 项目。

当前支持：

- 使用浏览器 Cookie 登录。
- 查看知乎推荐流。
- 按回答 ID 阅读推荐回答。
- 将知乎 Web API 地址放在可编辑的 JSON 配置里。

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

## 命令

```bash
./zhihu status
./zhihu feed -l 10
./zhihu feed --json
./zhihu read <回答ID>
./zhihu read <回答ID> --comments 5
./zhihu logout
```

当推荐项是回答时，`feed` 命令会打印回答 ID：

```text
阅读: zhihu read <回答ID>
```

## 接口配置

默认的知乎 Web API 地址配置在：

```text
configs/endpoints.json
```

运行时也可以指定其他配置文件：

```bash
./zhihu --endpoints ./configs/endpoints.json feed
```

这样接口地址变化时，不需要改命令实现代码。
