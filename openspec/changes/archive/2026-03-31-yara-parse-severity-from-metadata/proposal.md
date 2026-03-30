## Why

YARA 掃描結果的 `severity` 欄位永遠顯示 `"unknown"`，因為 CLI 呼叫時未加 `--print-meta` 旗標，導致輸出不含 metadata，無法取得規則定義的嚴重程度。修正此問題可讓回報結果更具可操作性，使用者能依嚴重程度排序或篩選偵測結果。

## What Changes

- **YARA CLI 呼叫加入 `--print-meta` 旗標**，使輸出包含規則的 metadata block
- **更新輸出解析邏輯**，從 `[severity="..."]` 格式中擷取 severity 值
- **找不到 severity metadata 時 fallback 為 `"unknown"`**（向下相容）
- **更新 `yara-engine` spec**，修正 `Severity` 的規格描述
- **更新相關測試**，反映新的輸出格式並驗證 severity 解析正確

## Capabilities

### New Capabilities

_（無新增 capability）_

### Modified Capabilities

- `yara-engine`：`Inspect()` 的輸出解析行為變更——原本 `Severity` 固定為 `"unknown"`，改為從 `--print-meta` 輸出的 metadata 中解析實際值；無 severity metadata 時 fallback 為 `"unknown"`

## Impact

- `internal/yara/yara.go`：CLI 呼叫加 flag、`parseOutput()` 解析邏輯更新、新增輔助函數
- `internal/yara/yara_test.go`：fake binary 輸出格式更新、severity 斷言更新、新增 fallback 測試案例
- `openspec/specs/yara-engine/spec.md`：更新 `Inspect()` requirement 的 severity 描述
