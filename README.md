# golang-concurrency-demo

30 個可獨立執行的 Go 併發範例，涵蓋 goroutine、channel、同步機制、取消與逾時、worker pool、錯誤傳遞，以及 race、deadlock、goroutine leak 的重現與修正。每個目錄只有一個 `main.go`，避免多個 `main` 衝突。

## 執行範例

需要 Go 1.22 以上、Python 3.9 以上及 C 編譯器（race detector 使用 cgo）。macOS 使用 Xcode Command Line Tools。依賴已放入 `vendor/`，不需要額外下載套件。

在專案根目錄執行：

```bash
go run ./examples/basic/04_atomic
# 最終計數: 1000

python3 scripts/run.py basic/04_atomic --race
python3 scripts/run.py supplement/pprof --race
```

`run.py` 將建置暫存檔放入 `.cache/`，執行後刪除 binary。在 macOS 上，腳本針對 pprof 範例使用系統 linker 並重新套用 ad-hoc signature，解決本機 Go 1.22 在新版 macOS 的 `LC_UUID` 與簽章問題。

## 驗證全部範例

```bash
python3 scripts/verify.py --repeat 3
```

驗證器檢查 `scripts/examples.json` 與 Go 程式的一對一覆蓋、`gofmt`、`go vet`，再逐一建置及執行 30 個範例。每次執行有 10 秒上限；只有符合預期輸出與退出碼才顯示 `PASS`。`--repeat 3` 會重用同一個 binary 執行三次。

執行驗證器的回歸測試：

```bash
python3 -m unittest discover -s scripts -p 'test_*.py'
```

正常與 leak 範例均以 `-race` 執行。race 範例必須偵測到 `WARNING: DATA RACE` 並退出 66；deadlock 範例不加 `-race`，必須由 Go runtime 回報死鎖並退出 2。這兩種失敗是預期結果，CI 會驗證它們。

## 重現錯誤與修正

```bash
# 預期退出 66；不可用計數器是否剛好等於 1000 判斷 race。
python3 scripts/run.py basic/08_race --race

# 預期退出 2；同步屏障確保兩個 goroutine 互相等待第二把鎖。
python3 scripts/run.py basic/09_deadlock

# 預期退出 0，輸出實際阻塞的 goroutine stack。
python3 scripts/run.py advanced/06_leak_send --race
python3 scripts/run.py ultimate/07_leak_receive --race

# 對照正確版本，等待背景工作真正退出。
python3 scripts/run.py supplement/leak_send_fixed --race
python3 scripts/run.py supplement/leak_receive_fixed --race
python3 scripts/run.py practical/03_lock_order --race
```

leak 範例只在程序存活期間保留阻塞的 goroutine；`main` 返回後整個程序結束，不會留下背景程序。pprof 範例透過 `httptest` 在記憶體內呼叫真正的 HTTP handler，不會綁定 port。

## 範例索引

`success` 預期退出 0；`race` 與 `deadlock` 預期依上述方式失敗；`leak` 以 stack 證明阻塞後退出 0。Goroutine 的輸出順序不固定，驗證器檢查值與數量，不比較排程順序。`sync.Pool` 允許隨時丟棄物件，因此不要求 `New` 只執行一次。

| Go 範例 | 預期 |
| --- | --- |
| [basic/01_waitgroup](examples/basic/01_waitgroup/main.go) | success |
| [basic/02_channel](examples/basic/02_channel/main.go) | success |
| [basic/03_rwmutex](examples/basic/03_rwmutex/main.go) | success |
| [basic/04_atomic](examples/basic/04_atomic/main.go) | success |
| [basic/05_context_cancel](examples/basic/05_context_cancel/main.go) | success |
| [basic/06_worker_pool](examples/basic/06_worker_pool/main.go) | success |
| [basic/07_semaphore](examples/basic/07_semaphore/main.go) | success |
| [basic/08_race](examples/basic/08_race/main.go) | race |
| [basic/09_deadlock](examples/basic/09_deadlock/main.go) | deadlock |
| [advanced/01_select](examples/advanced/01_select/main.go) | success |
| [advanced/02_once](examples/advanced/02_once/main.go) | success |
| [advanced/03_cond](examples/advanced/03_cond/main.go) | success |
| [advanced/04_sync_map](examples/advanced/04_sync_map/main.go) | success |
| [advanced/05_errgroup](examples/advanced/05_errgroup/main.go) | success |
| [advanced/06_leak_send](examples/advanced/06_leak_send/main.go) | leak |
| [ultimate/01_select_timeout](examples/ultimate/01_select_timeout/main.go) | success |
| [ultimate/02_once](examples/ultimate/02_once/main.go) | success |
| [ultimate/03_cond](examples/ultimate/03_cond/main.go) | success |
| [ultimate/04_sync_map](examples/ultimate/04_sync_map/main.go) | success |
| [ultimate/05_errgroup](examples/ultimate/05_errgroup/main.go) | success |
| [ultimate/06_pool](examples/ultimate/06_pool/main.go) | success |
| [ultimate/07_leak_receive](examples/ultimate/07_leak_receive/main.go) | leak |
| [practical/01_worker_pool_timeout](examples/practical/01_worker_pool_timeout/main.go) | success |
| [practical/02_select_no_deadlock](examples/practical/02_select_no_deadlock/main.go) | success |
| [practical/03_lock_order](examples/practical/03_lock_order/main.go) | success |
| [practical/04_errgroup](examples/practical/04_errgroup/main.go) | success |
| [supplement/mutex](examples/supplement/mutex/main.go) | success |
| [supplement/leak_send_fixed](examples/supplement/leak_send_fixed/main.go) | success |
| [supplement/leak_receive_fixed](examples/supplement/leak_receive_fixed/main.go) | success |
| [supplement/pprof](examples/supplement/pprof/main.go) | success |

## 發佈至公開 GitHub repository

需要可寫入 `.git`、可連線 GitHub 的環境，以及已登入 `aqutw` 的 GitHub CLI。

```bash
gh auth login
bash scripts/publish.sh
```

腳本先驗證全部範例，再提交檔案、建立公開的 `aqutw/golang-concurrency-demo` 並推送目前分支。若同名公開 repository 已存在，腳本會使用它；若 origin 指向其他 repository 或同名 repository 是 private，腳本會停止。腳本不會 force push，最後會核對遠端 commit 與公開狀態。

GitHub Actions 設定於 [verify.yml](.github/workflows/verify.yml)。第三方授權見 [THIRD_PARTY.md](THIRD_PARTY.md)。
