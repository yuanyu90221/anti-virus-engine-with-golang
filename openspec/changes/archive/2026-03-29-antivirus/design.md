## Context

本專案為全新建立的 Go 防毒引擎，使用 Pants Build 管理 monorepo。引擎透過 SHA256 雜湊比對，偵測目錄中的惡意檔案，並以 YAML 作為特徵資料庫格式。所有模組均為內部套件，透過單一 CLI 二進位檔對外暴露。

## Goals / Non-Goals

**Goals:**
- 實作可獨立運行的 CLI 掃描工具
- 以 `Loader` 介面解耦特徵資料庫格式與核心邏輯
- 支援文字與 JSON 兩種輸出格式，方便整合 CI/CD
- 使用 Pants Build 管理建構、測試與打包
- 僅依賴一個外部套件（`gopkg.in/yaml.v3`）

**Non-Goals:**
- 即時監控（inotify / FSEvents）
- 啟發式或行為分析（僅做 hash 比對）
- 自動更新特徵資料庫
- GUI 介面
- 跨平台打包或安裝程式

## Decisions

### D1：Loader 介面解耦格式
**決定**：`sigdb.DB` 接受 `Loader` 介面，`YAMLLoader` 為內建實作。
**原因**：未來若需支援 CSV、JSON 或遠端特徵庫，只需實作新 Loader，不需修改核心索引邏輯。
**替代方案**：直接在 `DB` 內硬編碼 YAML 解析 → 擴充性差，捨棄。

### D2：串流 SHA256（io.Copy）
**決定**：使用 `io.Copy` 將檔案內容串流至 `sha256.New()`，不一次讀入記憶體。
**原因**：掃描大型檔案時避免 OOM；與 `--max-size` 旗標配合可在超限時提前略過。
**替代方案**：`os.ReadFile` → 小檔案可行，大檔案有風險，捨棄。

### D3：`fs.WalkDir` 遞迴掃描
**決定**：使用標準函式庫 `fs.WalkDir`，預設不追蹤符號連結。
**原因**：避免循環連結造成無限遞迴；`--follow-links` 旗標需明確啟用。
**替代方案**：`filepath.Walk` → 已被 `WalkDir` 取代，效能較差，捨棄。

### D4：結束碼語意
**決定**：`0` = 乾淨、`1` = 偵測到威脅、`2` = 致命錯誤。
**原因**：符合 Unix 慣例與 CI 工具期望（非零 = 失敗）；`1` 與 `2` 分開讓呼叫端可區分「有威脅」與「工具本身錯誤」。

### D5：Pants Build 作為建構工具
**決定**：使用 `pants.toml` 設定 Go backend，以 `./pants tailor ::` 自動產生 `BUILD` 檔。
**原因**：monorepo 結構下 Pants 可精確追蹤依賴、增量測試與打包，比純 `go build` 更適合長期維護。

## Risks / Trade-offs

- **Pants 學習曲線** → 提供完整 `pants.toml` 範例與 `tailor` 指令說明降低門檻
- **YAML 解析效能**（數百筆特徵於啟動時一次載入）→ 啟動時全量載入至記憶體 map，掃描期間 O(1) 查詢，可接受
- **符號連結循環**（`--follow-links` 啟用時）→ 文件說明風險，使用者自行負責
- **SHA256 碰撞**（理論上）→ 對防毒場景碰撞機率可忽略，不作額外處理

## Migration Plan

全新專案，無既有程式碼需遷移。部署步驟：

1. `go mod init github.com/yuanyu/avengine && go get gopkg.in/yaml.v3`
2. 撰寫所有 Go 原始碼
3. `./pants tailor ::` 自動產生 BUILD 檔
4. `./pants test ::` 執行所有測試
5. `./pants package cmd/avengine:` 產生二進位檔

回滾：刪除整個目錄即可（無資料庫、無外部狀態）。

## Open Questions

- 特徵資料庫是否需要支援版本控制（e.g., 只載入 `version >= x` 的特徵）？→ 目前不需要，YAML 中的 `updated` 欄位僅供人工參考。
- `--max-size` 單位是否應為 MB（整數）或 bytes？→ 計畫規格定為 MB，旗標說明需明確標示。
