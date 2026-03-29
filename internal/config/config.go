// Package config 定義 CLI 掃描子命令的旗標結構與預設值。
// 此套件僅負責資料容器，不包含任何業務邏輯，讓 main.go 的旗標解析
// 與下游套件（scanner、sigdb）的呼叫保持解耦。
package config

// 預設值常數，集中定義以便日後維護。
const (
	DefaultSigsDir     = "./signatures" // 特徵 YAML 目錄的預設路徑
	DefaultOutput      = "text"         // 預設輸出格式：純文字表格
	DefaultYARATimeout = 10             // YARA 每個檔案的預設逾時（秒）
)

// CLIConfig 儲存 `avengine scan` 子命令解析後的所有旗標值。
// 欄位對應關係：
//
//	Dir         → --dir         掃描目標目錄（必填）
//	SigsDir     → --sigs        特徵資料庫目錄
//	Output      → --output      輸出格式（text | json）
//	FollowLinks → --follow-links 是否追蹤符號連結
//	MaxSizeMB   → --max-size    略過超過 N MB 的檔案（0 = 不限制）
//	Verbose     → --verbose     顯示所有掃描結果（含乾淨檔案）
type CLIConfig struct {
	Dir         string
	SigsDir     string
	Output      string
	FollowLinks bool
	MaxSizeMB   int
	Verbose     bool
	YARARules   string // --yara-rules：YARA 規則檔案或目錄路徑；空字串 = 停用 YARA
	YARATimeout int    // --yara-timeout：每個檔案的 YARA 逾時（秒）
}
