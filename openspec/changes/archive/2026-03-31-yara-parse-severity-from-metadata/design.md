## Context

YARA 引擎（`internal/yara/yara.go`）以 subprocess 方式呼叫 `yara` CLI。目前呼叫格式為：

```
yara <rulesPath> <filePath>
```

輸出只有：
```
RuleName /path/to/file
```

因此 `parseOutput()` 無法取得規則的 metadata，`Severity` 只能寫死為 `"unknown"`。YARA CLI 提供 `--print-meta`（或 `-m`）旗標，加上後輸出格式變為：

```
RuleName [key="value",...] /path/to/file
```

例如：
```
FakeC2Beacon [description="Detects fake C2 beacon",severity="high",date="2026-03-29"] /tmp/sample
```

## Goals / Non-Goals

**Goals:**
- 從 YARA metadata 中解析 `severity` 值，取代寫死的 `"unknown"`
- 無 `severity` metadata 的規則 fallback 為 `"unknown"`（向下相容）
- 不改變 `EngineDetection` 的資料型別或介面

**Non-Goals:**
- 解析 severity 以外的其他 metadata 欄位（如 `description`、`date`）
- 支援使用 YARA C library binding（維持 subprocess 架構）
- 對既有規則檔進行任何修改

## Decisions

### 決策 1：使用 `--print-meta` 旗標而非側邊映射檔

**選擇**：加 `--print-meta` 旗標讓 CLI 直接輸出 metadata。

**理由**：資料來源為 YARA 規則本身的 metadata，不需維護額外的對照表，不存在同步問題。

**考慮過的替代方案**：
- 側邊 YAML 對照表（rule name → severity）：需要手動維護，容易與規則脫節
- 從規則名稱前綴解析（如 `HIGH_RuleName`）：侵入性高，需修改所有現有規則

---

### 決策 2：解析邏輯拆為獨立輔助函數

**選擇**：新增 `parseYaraLine()` 與 `parseSeverityFromMeta()` 兩個 unexported 輔助函數。

**理由**：`parseOutput()` 職責單純（迭代行），解析邏輯獨立出來便於單元測試與日後擴充（如解析其他 metadata 欄位）。

---

### 決策 3：新格式與舊格式相容

**選擇**：`parseYaraLine()` 支援有無 `[...]` block 兩種格式。

**理由**：若未來有其他測試環境或舊版 YARA 不支援 `--print-meta`，不會直接崩潰，而是 fallback 為 `"unknown"`。

## Risks / Trade-offs

- **YARA CLI 版本相容性** → `--print-meta` 為 YARA 長期支援的旗標，風險低；若遇到不支援的版本，輸出格式不含 `[...]`，`parseYaraLine()` fallback 處理，不影響功能
- **metadata 值含特殊字元**（如逗號、等號）→ 目前以簡單的 `strings.Split(",")` 解析，若 value 含逗號會解析錯誤；severity 的值（`low/medium/high/critical`）不含特殊字元，此風險不影響當前使用場景
