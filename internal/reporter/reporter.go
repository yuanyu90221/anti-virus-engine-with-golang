// Package reporter 負責將 ScanReport 格式化輸出至任意 io.Writer。
//
// 支援兩種輸出格式：
//   - text：人類可讀的中文表格，適合終端機互動使用
//   - json：camelCase 鍵名的縮排 JSON，適合 CI/CD 流程的機器解析
//
// 格式選擇透過工廠函式 New(format) 在執行期決定，
// 呼叫端只依賴 Reporter 介面，不需知道底層實作。
package reporter

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/yuanyu90221/avengine/internal/scanner"
)

// 結束碼常數，遵循 UNIX 慣例：
//   - 0 表示成功且無異常
//   - 非零表示需要人工介入
const (
	ExitClean    = 0 // 掃描完成，未發現威脅
	ExitDetected = 1 // 偵測到至少一個威脅
	ExitError    = 2 // 執行階段錯誤（參數錯誤、IO 失敗等）
)

// Reporter 定義將 ScanReport 寫入 io.Writer 的統一介面。
// 使用介面而非具體型別，讓 main.go 可在不修改邏輯的情況下切換輸出格式。
type Reporter interface {
	Write(w io.Writer, report *scanner.ScanReport) error
}

// New 是 Reporter 的工廠函式，依 format 字串回傳對應實作。
// 未知格式時回傳 nil 與描述性錯誤，由呼叫端決定如何處理。
func New(format string) (Reporter, error) {
	switch format {
	case "text":
		return textReporter{}, nil
	case "json":
		return jsonReporter{}, nil
	default:
		return nil, fmt.Errorf("reporter: unknown format %q", format)
	}
}

// ---- 文字格式實作 ----

type textReporter struct{}

// Write 輸出中文摘要行與偵測結果表格。
// 表格欄位寬度以空白對齊，方便終端機閱讀。
func (textReporter) Write(w io.Writer, r *scanner.ScanReport) error {
	duration := r.FinishedAt.Sub(r.StartedAt)
	fmt.Fprintf(w, "掃描完成\n")
	fmt.Fprintf(w, "檔案總數: %d  威脅: %d  錯誤: %d  略過: %d  耗時: %s\n\n",
		r.TotalFiles, len(r.Detections), r.ErrorFiles, r.SkippedFiles, duration)

	if len(r.Detections) == 0 {
		fmt.Fprintln(w, "未發現威脅。")
		return nil
	}

	// 表格標題
	fmt.Fprintf(w, "%-60s  %-16s  %-30s  %-8s  %-12s\n",
		"路徑", "SHA256(前16)", "威脅名稱", "嚴重度", "分類")
	fmt.Fprintf(w, "%s\n", "------------------------------------------------------------------------")

	// SHA256 僅顯示前 16 字元（8 位元組），兼顧可讀性與識別度
	for _, d := range r.Detections {
		short := d.SHA256
		if len(short) > 16 {
			short = short[:16]
		}
		fmt.Fprintf(w, "%-60s  %-16s  %-30s  %-8s  %-12s\n",
			d.Path, short, d.Name, d.Severity, d.Category)
	}
	return nil
}

// ---- JSON 格式實作 ----

type jsonReporter struct{}

// jsonReport 是 ScanReport 的可序列化版本，使用 camelCase 鍵名。
// 獨立定義此結構體以控制 JSON 輸出格式，
// 避免 scanner.ScanReport 的欄位變動直接影響 API 契約。
type jsonReport struct {
	Detections   []jsonDetection `json:"detections"`
	TotalFiles   int             `json:"totalFiles"`
	ErrorFiles   int             `json:"errorFiles"`
	SkippedFiles int             `json:"skippedFiles"`
	StartedAt    string          `json:"startedAt"`  // RFC 3339 格式
	FinishedAt   string          `json:"finishedAt"` // RFC 3339 格式
}

// jsonDetection 是單筆偵測結果的 JSON 表示。
type jsonDetection struct {
	Path     string `json:"path"`
	SHA256   string `json:"sha256"`
	Name     string `json:"name"`
	Category string `json:"category"`
	Severity string `json:"severity"`
}

// Write 將 ScanReport 序列化為縮排 JSON 並寫入 w。
// 時間欄位使用 RFC 3339（含時區偏移）格式，便於跨系統解析。
func (jsonReporter) Write(w io.Writer, r *scanner.ScanReport) error {
	jd := make([]jsonDetection, len(r.Detections))
	for i, d := range r.Detections {
		jd[i] = jsonDetection{
			Path:     d.Path,
			SHA256:   d.SHA256,
			Name:     d.Name,
			Category: d.Category,
			Severity: d.Severity,
		}
	}
	out := jsonReport{
		Detections:   jd,
		TotalFiles:   r.TotalFiles,
		ErrorFiles:   r.ErrorFiles,
		SkippedFiles: r.SkippedFiles,
		StartedAt:    r.StartedAt.Format("2006-01-02T15:04:05Z07:00"),
		FinishedAt:   r.FinishedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ") // 縮排 2 格空白，提升 JSON 可讀性
	return enc.Encode(out)
}
