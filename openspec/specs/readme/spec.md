## Requirements

### Requirement: 儲存庫根目錄存在 README.md
儲存庫根目錄必須存在 `README.md`，包含專案概述、架構圖、快速開始指引、特徵 YAML 格式範例、病毒 hash 參考表、結束碼說明，以及參考資源章節。

#### Scenario: README 包含必要章節
- **WHEN** 使用者在 GitHub 或本機開啟 `README.md`
- **THEN** 文件必須包含「架構」、「快速開始」、「病毒特徵 hash 參考表」、「結束碼」與「參考資源」等章節

### Requirement: 病毒特徵 hash 參考表列出公開已記錄的 SHA256 值
README 必須包含一張表格，列出至少四筆來自公開來源的病毒特徵 SHA256 hash，每筆須標示名稱、hash 值、分類，以及資料來源出處。

#### Scenario: 參考表包含 EICAR 測試檔案 hash
- **WHEN** 使用者查閱病毒特徵 hash 參考表
- **THEN** 表格必須包含 EICAR 測試檔案的 SHA256（`275a021bbfb6489e54d471899f7db9d1663fc695ec2fe2a2c4538aabf651fd0f`）並標示來源為 eicar.org

#### Scenario: 參考表包含 WannaCry hash
- **WHEN** 使用者查閱病毒特徵 hash 參考表
- **THEN** 表格必須包含 WannaCry 勒索軟體的 SHA256（`24d004a104d4d54034dbcffc2a4b19a11f39008a575aa614ea04703480b1022c`）並標示來源為 CISA

#### Scenario: 所有 hash 值格式正確
- **WHEN** 使用者讀取參考表中的任一 hash 值
- **THEN** 每個 hash 值必須為 64 個小寫十六進位字元的 SHA256 字串

### Requirement: 參考資源章節包含 YouTube 影片連結
README 的參考資源章節必須包含指向 `https://www.youtube.com/watch?v=s_M1vKp69hA` 的連結。

#### Scenario: YouTube 連結可被點擊
- **WHEN** 使用者在 GitHub README 中點擊參考資源章節的 YouTube 連結
- **THEN** 連結必須導向該教學影片的正確 URL
