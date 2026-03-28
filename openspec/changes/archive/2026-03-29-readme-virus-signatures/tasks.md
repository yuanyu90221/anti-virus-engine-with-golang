## 1. 建立 README.md 主體結構

- [x] 1.1 在儲存庫根目錄建立 `README.md`，包含以下章節：專案標題、概述、架構、快速開始、特徵格式、病毒 hash 參考表、結束碼、參考資源

## 2. 撰寫架構說明

- [x] 2.1 在「架構」章節加入 ASCII 流程圖，呈現 `hasher → sigdb → scanner → reporter` 的資料流
- [x] 2.2 簡述各模組職責（hasher：SHA256 計算；sigdb：特徵載入與查詢；scanner：目錄遞迴掃描；reporter：輸出報告）

## 3. 撰寫快速開始

- [x] 3.1 列出環境需求：Go 1.22+、Pants 2.23.0
- [x] 3.2 提供建置指令：`./pants package cmd/avengine:`
- [x] 3.3 提供使用範例：`avengine scan --dir ./testdata --sigs ./signatures`

## 4. 撰寫特徵 YAML 格式說明

- [x] 4.1 加入 YAML 格式範例程式碼區塊，包含 `version`、`category`、`updated`、`signatures[]`（含 `sha256`、`name`、`severity`、`added`）欄位

## 5. 加入病毒特徵 hash 參考表

- [x] 5.1 建立 Markdown 表格，欄位為：名稱、SHA256、分類、嚴重程度、來源
- [x] 5.2 填入 EICAR 測試檔案：`275a021bbfb6489e54d471899f7db9d1663fc695ec2fe2a2c4538aabf651fd0f`，來源 eicar.org
- [x] 5.3 填入 WannaCry：`24d004a104d4d54034dbcffc2a4b19a11f39008a575aa614ea04703480b1022c`，來源 CISA
- [x] 5.4 填入 NotPetya：`027cc450ef5f8c5f653329641ec1fed91f694e0d229928963b30f6b0d7d3a745`，來源 US-CERT
- [x] 5.5 填入 Mirai：`9a024b9ef95a1ed9e5acba3e2fe2427395e866c42b5bce04d35ca2cefd8d2e4d`，來源 Malware Traffic Analysis
- [x] 5.6 確認每個 SHA256 值均為 64 個小寫十六進位字元

## 6. 加入結束碼說明與參考資源

- [x] 6.1 加入結束碼對照表：`0` = 無威脅、`1` = 偵測到威脅、`2` = 執行錯誤
- [x] 6.2 在參考資源章節加入 YouTube 影片連結：`https://www.youtube.com/watch?v=s_M1vKp69hA`
