## Context

防毒引擎已具備 YARA 偵測能力（`internal/yara/`），可透過 `yara` CLI subprocess 執行規則比對。然而目前專案中缺乏實際的 YARA 規則檔案，引擎無法偵測任何威脅。本設計說明如何建立 `rules/` 目錄結構並加入涵蓋常見威脅的規則集。

## Goals / Non-Goals

**Goals:**
- 建立 `rules/` 頂層目錄，依威脅類別分檔存放 `.yar` 規則
- 提供涵蓋 ransomware、trojan、dropper、webshell、coinminer 的規則
- 每條規則包含完整 meta 資訊（description、severity、reference）
- 更新 yara-engine 規則載入路徑以支援目錄批次載入

**Non-Goals:**
- 不自動更新或同步外部規則來源
- 不修改偵測引擎核心邏輯（`internal/yara/`）
- 不提供規則版本管理機制

## Decisions

### 目錄結構：頂層 `rules/` vs 嵌入 `testdata/`

選擇在專案根目錄建立 `rules/`，原因：
- 規則檔案是產品的一部分，不是測試資料
- 讓 CLI 預設可用 `--rules ./rules` 載入，直覺明確
- `testdata/yara-e2e/` 保留給 E2E 測試專用的最小規則集

### 規則分類方式：以威脅類別分檔

每個威脅類別一個 `.yar` 檔案（例如 `ransomware.yar`、`trojan.yar`），優點：
- 使用者可選擇性載入特定類別
- 維護時改動範圍明確
- 規則條數增加時不至於單一檔案過大

### 規則來源：手工整理的公開規則特徵

基於公開威脅情報（VT、AlienVault OTX、Malpedia）的特徵字串，不依賴外部工具或付費規則集。每條規則須通過 `yara` CLI 語法驗證才加入。

## Risks / Trade-offs

- **誤報（False Positive）** → 使用高特異性字串（magic bytes + 多重 and 條件），盡量避免泛用字串
- **規則過時** → meta 欄位標註 `date` 與 `reference`，方便未來審查更新
- **規則數量有限** → 初版以代表性規則為主，建立結構後可持續擴充
