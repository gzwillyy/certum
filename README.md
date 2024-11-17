# Certum 验证服务器

用于验证 Certum 证书的轻量级 HTTP 服务器。此服务器创建一个验证文件，并根据证书验证的需要通过特定的 URL 路径提供该文件。

---

## 功能

- 在 `.well-known/pki-validation/certum.txt` 处创建一个验证文件。
- 在 HTTP 上提供验证文件。
- 在收到 `Ctrl+C` (SIGINT) 时处理正常关机。
- 终止时自动清理临时文件。