## 1. 新增 Pants 建置產物忽略規則

- [x] 1.1 在 `.gitignore` 中新增 `.pants.d/` 忽略規則（Pants 建置快取目錄）
- [x] 1.2 在 `.gitignore` 中新增 `dist/` 忽略規則（pants package 輸出目錄）

## 2. 新增執行時產出忽略規則

- [x] 2.1 在 `.gitignore` 中新增 `scan-report.json` 忽略規則（avengine 掃描報告輸出）

## 3. 新增 IDE 設定忽略規則

- [x] 3.1 在 `.gitignore` 中新增 `.idea/` 忽略規則（JetBrains IDE）
- [x] 3.2 在 `.gitignore` 中新增 `.vscode/` 忽略規則（VS Code）

## 4. 新增作業系統暫存檔忽略規則

- [x] 4.1 在 `.gitignore` 中新增 `.DS_Store` 忽略規則（macOS Finder）
- [x] 4.2 在 `.gitignore` 中新增 `Thumbs.db` 忽略規則（Windows）

## 5. 驗證

- [x] 5.1 確認 `dist/`、`.pants.d/`、`scan-report.json` 均不出現在 `git status` 的未追蹤列表中
