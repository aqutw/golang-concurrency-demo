# 第三方原始碼

`vendor/golang.org/x/sync/errgroup` 使用 `golang.org/x/sync v0.5.0` 的 errgroup 實作與 Go 版本相容 helper，保留 BSD 授權及專利聲明。此版本支援本專案使用的 Go 1.22。

來源：[Go Authors 的 errgroup v0.5.0](https://github.com/golang/sync/tree/v0.5.0/errgroup)。授權：[LICENSE](vendor/golang.org/x/sync/LICENSE)、[PATENTS](vendor/golang.org/x/sync/PATENTS)。本專案將依賴放入 `vendor/`，讓執行與驗證不需要下載套件。

在可連線的環境重新取得官方依賴：

```bash
go mod download
go mod vendor
python3 scripts/verify.py --repeat 3
```
