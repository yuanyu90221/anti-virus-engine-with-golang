// Package sigdb 管理病毒特徵資料庫的載入與查詢。
//
// 核心設計：以「依賴注入 + 介面」取代硬編碼的資料來源，
// Loader 介面讓 DB 的建立與資料來源完全解耦——
// 生產環境使用 YAMLLoader 從磁碟讀取，測試時可注入假資料，
// 未來也能擴充為 HTTP、資料庫等來源，無需修改 DB 本身。
package sigdb

import "strings"

// Loader 是特徵資料來源的抽象介面。
// 任何實作 Load() 方法的型別都可作為 NewDB 的資料來源。
type Loader interface {
	Load() ([]Signature, error)
}

// Signature 代表一筆病毒特徵記錄，對應 YAML 檔案中的單一項目。
// Category 欄位不在 YAML 的 signatures 列表內，
// 而是從檔案層級的 category 欄位注入（由 YAMLLoader 負責填入）。
type Signature struct {
	SHA256   string `yaml:"sha256"`   // 64 字元小寫十六進位，為主要比對鍵
	Name     string `yaml:"name"`     // 人類可讀的威脅名稱
	Severity string `yaml:"severity"` // 嚴重程度：low / medium / high / critical
	Added    string `yaml:"added"`    // 加入日期（YYYY-MM-DD）
	Category string // 從 YAML 檔案層級的 category 欄位注入
}

// MatchResult 是 DB.Lookup 命中時回傳的輕量資料，
// 僅包含報告所需的欄位，不暴露完整的 Signature 結構。
type MatchResult struct {
	Name     string
	Category string
	Severity string
}

// DB 是特徵資料庫的記憶體索引，以小寫 SHA256 十六進位字串為鍵。
// 一旦建立後為唯讀，並行讀取安全（無鎖）。
type DB struct {
	index map[string]MatchResult // key：小寫 SHA256，value：比對結果
}

// NewDB 接受任意 Loader，載入所有特徵後建立記憶體索引。
// 若兩筆特徵有相同的 SHA256，後載入的會覆蓋先前的（最後寫入勝出）。
func NewDB(loader Loader) (*DB, error) {
	sigs, err := loader.Load()
	if err != nil {
		return nil, err
	}
	// 預分配 map 容量，避免插入時多次 rehash
	idx := make(map[string]MatchResult, len(sigs))
	for _, s := range sigs {
		// 統一轉小寫，使 hash 比對不區分大小寫
		key := strings.ToLower(s.SHA256)
		idx[key] = MatchResult{
			Name:     s.Name,
			Category: s.Category,
			Severity: s.Severity,
		}
	}
	return &DB{index: idx}, nil
}

// Lookup 以 sha256hex 查詢資料庫。
// 命中時回傳 (MatchResult, true)；未命中回傳 (零值, false)。
// 輸入不區分大小寫，內部統一轉小寫後比對。
func (db *DB) Lookup(sha256hex string) (MatchResult, bool) {
	r, ok := db.index[strings.ToLower(sha256hex)]
	return r, ok
}
