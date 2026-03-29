// Package scanner 實作遞迴目錄掃描，對每個檔案計算 SHA256 並查詢特徵資料庫。
//
// 掃描流程：
//
//	filepath.WalkDir（深度優先遍歷）
//	  └─ 每個非目錄項目
//	       ├─ 符號連結？→ FollowLinks=false 時略過
//	       ├─ 超過大小上限？→ 略過
//	       ├─ hasher.HashFile → 計算 SHA256
//	       ├─ sigdb.DB.Lookup → 命中則加入 Detections（hash 引擎）
//	       └─ ExtraEngines[i].Inspect → 每個額外引擎的偵測結果
package scanner

import (
	"context"
	"io/fs"
	"path/filepath"
	"time"

	"github.com/yuanyu90221/avengine/internal/hasher"
	"github.com/yuanyu90221/avengine/internal/sigdb"
)

// ProgressFunc 是掃描進度的回呼型別。
// path 為當前處理的檔案路徑，count 為已處理的檔案累計數（從 1 開始）。
type ProgressFunc func(path string, count int64)

// EngineDetection 是單一偵測引擎回傳的原始命中結果。
type EngineDetection struct {
	Name     string
	Category string
	Severity string
}

// DetectionEngine 是可插拔偵測後端的擴充點。
// Scan() 對每個通過大小/符號連結過濾的檔案呼叫 Inspect。
type DetectionEngine interface {
	// Name 回傳穩定的識別名稱，用於 Detection.Engine（例如 "hash"、"yara"）。
	Name() string
	// Inspect 檢查 path 指定的檔案，回傳零或多個偵測結果。
	// ctx 攜帶每個檔案的截止時間。回傳非 nil error 時計入 ErrorFiles。
	Inspect(ctx context.Context, path string) ([]EngineDetection, error)
}

// Options 設定單次掃描的行為參數。
type Options struct {
	Dir          string           // 遞迴掃描的根目錄（必填）
	FollowLinks  bool             // true = 追蹤符號連結；false = 略過（預防迴圈）
	MaxFileSizeB int64            // 超過此大小（位元組）的檔案會被略過；0 = 不限制
	OnProgress   ProgressFunc     // 每處理完一個檔案後呼叫；nil = 不通知
	ExtraEngines []DetectionEngine // 額外偵測引擎；nil = 僅使用 hash 引擎
	FileTimeout  time.Duration    // 每個檔案的 context 截止時間；0 = 不限制
}

// Detection 記錄一個命中偵測引擎的檔案。
// 嵌入 sigdb.MatchResult 以直接存取 Name、Category、Severity。
type Detection struct {
	Path   string // 檔案的完整路徑
	SHA256 string // 64 字元小寫十六進位雜湊值
	Engine string // 偵測引擎名稱，例如 "hash"、"yara"
	sigdb.MatchResult
}

// ScanReport 是一次掃描的完整結果。
//
// 計數說明：
//   - TotalFiles：成功計算 hash 的檔案數（不含略過與錯誤）
//   - ErrorFiles：開啟或讀取失敗的檔案數
//   - SkippedFiles：因符號連結或大小限制而略過的檔案數
type ScanReport struct {
	Detections   []Detection // 命中特徵的檔案列表
	TotalFiles   int         // 成功掃描的檔案總數
	ErrorFiles   int         // 讀取失敗的檔案數
	SkippedFiles int         // 略過的檔案數
	StartedAt    time.Time   // 掃描開始時間
	FinishedAt   time.Time   // 掃描結束時間
}

// CountFiles 走訪 opts.Dir 並套用相同的過濾條件（FollowLinks、MaxFileSizeB），
// 回傳符合條件的檔案總數。
func CountFiles(opts Options) (int64, error) {
	var count int64
	err := filepath.WalkDir(opts.Dir, func(_ string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil // 略過無法進入的路徑，不中斷計數
		}
		if d.IsDir() {
			return nil
		}
		if d.Type()&fs.ModeSymlink != 0 && !opts.FollowLinks {
			return nil
		}
		if opts.MaxFileSizeB > 0 {
			info, err := d.Info()
			if err != nil || info.Size() > opts.MaxFileSizeB {
				return nil
			}
		}
		count++
		return nil
	})
	return count, err
}

// Scan 遞迴走訪 opts.Dir，對每個符合條件的檔案計算 SHA256 並查詢 db，
// 以及執行 opts.ExtraEngines 中的所有額外偵測引擎。
//
// 錯誤處理策略：
//   - WalkDir 回呼收到的進入錯誤（如無讀取權限）→ 計入 ErrorFiles，繼續掃描
//   - hasher.HashFile 失敗 → 計入 ErrorFiles，TotalFiles 不加，繼續掃描
//   - opts.Dir 本身不存在 → 回傳 error，立即中止
func Scan(ctx context.Context, db *sigdb.DB, opts Options) (*ScanReport, error) {
	report := &ScanReport{StartedAt: time.Now()}

	err := filepath.WalkDir(opts.Dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			// 無法進入此路徑（如權限不足），記錄後繼續走訪其他項目
			report.ErrorFiles++
			return nil
		}
		if d.IsDir() {
			return nil // 目錄本身不計算 hash，繼續遞迴進入
		}

		// 符號連結處理：預設略過以避免迴圈參照
		if d.Type()&fs.ModeSymlink != 0 {
			if !opts.FollowLinks {
				report.SkippedFiles++
				return nil
			}
		}

		// 大小限制：取得 FileInfo 並比較，避免對超大檔案計算 hash
		if opts.MaxFileSizeB > 0 {
			info, err := d.Info()
			if err != nil {
				report.ErrorFiles++
				return nil
			}
			if info.Size() > opts.MaxFileSizeB {
				report.SkippedFiles++
				return nil
			}
		}

		report.TotalFiles++
		hash, err := hasher.HashFile(path)
		if err != nil {
			// 檔案在 WalkDir 列舉後、Hash 前被刪除或失去權限
			report.ErrorFiles++
			report.TotalFiles-- // 未能成功掃描，不計入掃描數
			return nil
		}

		// 建立每個檔案的 context（攜帶截止時間）
		fileCtx := ctx
		var cancel context.CancelFunc
		if opts.FileTimeout > 0 {
			fileCtx, cancel = context.WithTimeout(ctx, opts.FileTimeout)
		} else {
			fileCtx, cancel = context.WithCancel(ctx)
		}
		defer cancel()

		// hash 引擎：以 SHA256 查詢記憶體索引，O(1) 時間複雜度
		if match, ok := db.Lookup(hash); ok {
			report.Detections = append(report.Detections, Detection{
				Path:        path,
				SHA256:      hash,
				Engine:      "hash",
				MatchResult: match,
			})
		}

		// 額外引擎（例如 YARA）：依序執行，每個引擎的錯誤不中斷後續引擎
		for _, eng := range opts.ExtraEngines {
			hits, engErr := eng.Inspect(fileCtx, path)
			if engErr != nil {
				report.ErrorFiles++
			}
			for _, h := range hits {
				report.Detections = append(report.Detections, Detection{
					Path:   path,
					SHA256: hash,
					Engine: eng.Name(),
					MatchResult: sigdb.MatchResult{
						Name:     h.Name,
						Category: h.Category,
						Severity: h.Severity,
					},
				})
			}
		}

		if opts.OnProgress != nil {
			opts.OnProgress(path, int64(report.TotalFiles))
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	report.FinishedAt = time.Now()
	return report, nil
}
