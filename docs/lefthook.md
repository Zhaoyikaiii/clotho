# Lefthook Setup

Lefthook 已配置为自动在 git commit/push 时运行。

## 自动运行的钩子

| 钩子 | 时机 | 检查内容 |
|------|------|----------|
| `pre-commit` | git commit 前 | gofmt, goimports, golangci-lint, generate, secret-scan |
| `pre-push` | git push 前 | go test -race, vuln-scan, license-check |
| `post-merge` | git pull/merge 后 | go mod tidy |

## 手动运行

```bash
# 运行 pre-commit 钩子
lefthook run pre-commit

# 运行 pre-push 钩子
lefthook run pre-push

# 安装/更新钩子
make lefthook-install
```

## 跳过钩子

```bash
# 跳过本次 commit 的钩子
git commit --no-verify -m "Skip hooks"

# 跳过本次 push 的钩子
git push --no-verify
```
