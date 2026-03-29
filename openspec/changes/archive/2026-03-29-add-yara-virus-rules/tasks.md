## 1. 建立 rules 目錄結構

- [x] 1.1 在專案根目錄建立 `rules/` 目錄
- [x] 1.2 建立 `rules/ransomware.yar`，加入 WannaCry、Ryuk、LockBit 特徵規則（含完整 meta）
- [x] 1.3 建立 `rules/trojan.yar`，加入常見 RAT／後門特徵規則（含完整 meta）
- [x] 1.4 建立 `rules/dropper.yar`，加入常見 dropper／downloader 行為特徵規則（含完整 meta）
- [x] 1.5 建立 `rules/webshell.yar`，加入常見 PHP／ASP webshell 特徵規則（含完整 meta）
- [x] 1.6 建立 `rules/coinminer.yar`，加入挖礦程式特徵規則（含完整 meta）

## 2. 規則驗證

- [x] 2.1 對每個 `.yar` 檔案執行 `yara <file>.yar /dev/null` 確認語法正確

## 3. 更新 YARA 引擎支援目錄載入

- [x] 3.1 修改 `internal/yara/engine.go`（或對應實作），使 `Inspect()` 在 `rulesPath` 為目錄時，批次傳入目錄下所有 `.yar` 檔案給 `yara` CLI
- [x] 3.2 新增或更新對應單元測試，驗證目錄載入行為

## 4. 文件與整合

- [x] 4.1 在 `README.md` 或相關說明文件中補充如何使用 `--rules ./rules` 載入規則目錄
