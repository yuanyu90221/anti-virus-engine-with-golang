## Context

防毒引擎原本只有 SHA256 hash 比對一種偵測方式。hash 比對對已知樣本的偵測率為 100%，但對於修改過 1 個位元組的變種完全無效。YARA 是業界標準的規則語言，透過字串、正規表達式、PE 結構等多維度比對，可涵蓋 hash 無法觸及的威脅。

本設計採用 CLI subprocess 方式整合 YARA，而非 CGo binding，以保持 pure Go 可攜性並降低建置複雜度。

## Goals / Non-Goals

**Goals:**
- 定義 `DetectionEngine` 介面，讓掃描器支援多個可插拔的偵測後端
- 透過 YARA CLI subprocess 實現 YARA 規則比對
- `--yara-rules` 未提供時，行為與原有 hash-only 模式 100% 相同（向下相容）
- 偵測結果標示來源引擎（`Engine` 欄位），文字與 JSON 輸出均呈現

**Non-Goals:**
- CGo 或 Go YARA binding（不引入 C 依賴）
- 自行解析 YARA 規則格式
- YARA 版本管理或自動安裝（由使用者或 CI 預先安裝）

## Decisions

**決策 1：`DetectionEngine` 介面放在 `scanner` 套件**

介面由 `scanner.Scan()` 消費，放在同一套件可避免 `yara` 套件 import `scanner` 而形成循環依賴。`yara` 套件只 import `scanner`，單向依賴清晰。

備選方案：另立 `engine` 套件。拒絕原因：僅為放置介面而增加一層 package 無實質收益。

**決策 2：hash 引擎不包裝成 `DetectionEngine`**

hash 引擎永遠先跑（後續引擎也需要 SHA256 值），且其效能已是 O(1)，不需要透過介面呼叫。`ExtraEngines` 只放 hash 以外的引擎，避免過度設計。

**決策 3：YARA 以 subprocess 方式呼叫（逐檔）**

`exec.CommandContext(ctx, "yara", rulesPath, filePath)` 逐檔呼叫，配合 `FileTimeout` 傳遞截止時間。YARA CLI 的 exit 1 語意為「無比對」（非錯誤），需特別處理。

備選方案：`yara -r` 掃描整個目錄一次。拒絕原因：無法對每個檔案套用 context 截止時間，且錯誤隔離性差。

**決策 4：`Scan()` 加入 `context.Context` 第一參數**

這是唯一的 breaking change。加入 ctx 是 Go 慣例，讓 per-file timeout 可以沿呼叫鏈傳遞。舊有呼叫端只需補上 `context.Background()`。

## Risks / Trade-offs

- [效能] YARA subprocess 逐檔呼叫有 process spawn overhead。→ 僅在提供 `--yara-rules` 時才啟用，hash-only 路徑不受影響。大量小檔案掃描時可考慮未來改用批次模式。
- [YARA 安裝] `yara` binary 不在 Go 依賴中，需使用者自行安裝。→ 以 `exec.LookPath` 在啟動時檢查，找不到時 warn 並降級（不中止掃描）。
- [規則品質] YARA 規則誤報率取決於規則撰寫品質，與引擎無關。→ 記錄在使用文件中，非引擎責任。
- [Breaking change] `Scan()` 簽名加入 `context.Context`。→ 影響範圍僅限 `main.go` 與測試檔，已同步更新。

## Migration Plan

1. 現有呼叫 `scanner.Scan(db, opts)` 改為 `scanner.Scan(context.Background(), db, opts)`
2. 若需啟用 YARA：安裝 `yara` binary（見 `yara-cli-setup` spec），並加入 `--yara-rules` 旗標
3. 降級：移除 `--yara-rules` 即可回到 hash-only 模式，無資料遷移需求
