## Why

目前防毒引擎已整合 YARA 偵測引擎，但缺乏實際的 YARA 規則檔案，導致引擎無法有效偵測常見病毒與惡意軟體。為了使引擎具備開箱即用的偵測能力，需要在 `rules/` 目錄下加入涵蓋常見病毒家族、惡意行為特徵的 YARA 規則文件。

## What Changes

- 新增 `rules/` 目錄，存放 YARA 格式的規則檔案
- 加入涵蓋以下常見威脅類別的規則：
  - Ransomware（勒索軟體）：WannaCry、Ryuk、LockBit 等特徵
  - Trojans（木馬）：常見後門與遠端存取工具特徵
  - Malware dropper / downloader 行為特徵
  - Webshells（網頁後門）
  - Coinminer（挖礦程式）
- 規則依威脅類別分檔存放（每個類別一個 `.yar` 檔案）
- 每條規則包含 `meta`（描述、嚴重性、來源參考）、`strings`、`condition` 區段

## Capabilities

### New Capabilities
- `yara-virus-rules`: 提供一組涵蓋常見病毒家族與惡意行為的 YARA 規則檔案，供偵測引擎載入使用

### Modified Capabilities
- `yara-engine`: 更新規則載入路徑以支援從 `rules/` 目錄批次載入規則檔

## Impact

- 新增 `rules/` 目錄於專案根目錄（或 `testdata/yara-e2e/` 依需求）
- 現有 YARA 引擎（`internal/yara/`）需能掃描並載入目錄下所有 `.yar` 規則
- 不影響現有 API 或 CLI 介面
- 規則文件為純文字，不引入新的 Go 依賴
