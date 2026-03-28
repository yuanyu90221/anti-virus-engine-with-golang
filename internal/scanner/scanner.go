// Package scanner 實作遞迴目錄掃描，對每個檔案計算 SHA256 並查詢特徵資料庫。
//
// 掃描流程：
//
//	filepath.WalkDir（深度優先遍歷）
//	  └─ 每個非目錄項目
//	       ├─ 符號連結？→ FollowLinks=false 時略過
//	       ├─ 超過大小上限？→ 略過
//	       ├─ hasher.HashFile → 計算 SHA256
//	       └─ sigdb.DB.Lookup → 命中則加入 Detections
package scanner

import (
	"io/fs"
	"path/filepath"
	"time"

	"github.com/yuanyu90221/avengine/internal/hasher"
	"github.com/yuanyu90221/avengine/internal/sigdb"
)

// Options 設定單次掃描的行為參數。
type Options struct {
	Dir          string // 遞迴掃描的根目錄（必填）
	FollowLinks  bool   // true = 追蹤符號連結；false = 略過（預防迴圈）
	MaxFileSizeB int64  // 超過此大小（位元組）的檔案會被略過；0 = 不限制
}

// Detection 記錄一個命中特徵資料庫的檔案。
// 嵌入 sigdb.MatchResult 以直接存取 Name、Category、Severity。
type Detection struct {
	Path   string // 檔案的完整路徑
	SHA256 string // 64 字元小寫十六進位雜湊值
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

// Scan 遞迴走訪 opts.Dir，對每個符合條件的檔案計算 SHA256 並查詢 db。
//
// 錯誤處理策略：
//   - WalkDir 回呼收到的進入錯誤（如無讀取權限）→ 計入 ErrorFiles，繼續掃描
//   - hasher.HashFile 失敗 → 計入 ErrorFiles，TotalFiles 不加，繼續掃描
//   - opts.Dir 本身不存在 → 回傳 error，立即中止
func Scan(db *sigdb.DB, opts Options) (*ScanReport, error) {
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

		// 以 SHA256 查詢記憶體索引，O(1) 時間複雜度
		if match, ok := db.Lookup(hash); ok {
			report.Detections = append(report.Detections, Detection{
				Path:        path,
				SHA256:      hash,
				MatchResult: match,
			})
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	report.FinishedAt = time.Now()
	return report, nil
}
