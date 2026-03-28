package sigdb

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// yamlFile 對應單一特徵 YAML 檔案的頂層結構。
// 檔案層級的 category 欄位會在載入後注入到每筆 Signature.Category，
// 這樣 Signature 本身不需要重複記錄所屬分類。
type yamlFile struct {
	Version    string      `yaml:"version"`
	Category   string      `yaml:"category"`   // 例如 "ransomware"、"trojan"
	Updated    string      `yaml:"updated"`    // 最後更新日期
	Signatures []Signature `yaml:"signatures"` // 特徵列表
}

// YAMLLoader 從指定目錄讀取所有 .yaml 特徵檔案，實作 Loader 介面。
//
// 目錄結構範例：
//
//	signatures/
//	  ransomware.yaml
//	  trojans.yaml
//
// 每個 .yaml 檔案對應一個惡意軟體分類，啟動時全部合併為單一索引。
type YAMLLoader struct {
	Dir string // 特徵 YAML 檔案所在目錄
}

// Load 讀取 Dir 下所有 .yaml 檔案，解析後合併回傳完整特徵列表。
//
// 處理流程：
//  1. os.ReadDir 列舉目錄（已按檔名排序，確保載入順序一致）
//  2. 過濾非 .yaml 檔案與子目錄
//  3. 逐一讀取並以 yaml.Unmarshal 反序列化
//  4. 將檔案層級的 category 注入每筆 Signature.Category
//  5. 合併至結果切片
//
// 任何單一檔案的讀取或解析錯誤都會立即中止並回傳錯誤，
// 確保引擎不會在特徵不完整的狀態下啟動。
func (y YAMLLoader) Load() ([]Signature, error) {
	entries, err := os.ReadDir(y.Dir)
	if err != nil {
		return nil, fmt.Errorf("sigdb: read dir %s: %w", y.Dir, err)
	}

	var all []Signature
	for _, e := range entries {
		// 跳過子目錄與非 YAML 檔案
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".yaml") {
			continue
		}
		path := filepath.Join(y.Dir, e.Name())
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("sigdb: read %s: %w", e.Name(), err)
		}
		var yf yamlFile
		if err := yaml.Unmarshal(data, &yf); err != nil {
			return nil, fmt.Errorf("sigdb: parse %s: %w", e.Name(), err)
		}
		// 將檔案層級的 category 注入每筆特徵，讓 MatchResult 可直接取用
		for i := range yf.Signatures {
			yf.Signatures[i].Category = yf.Category
		}
		all = append(all, yf.Signatures...)
	}
	return all, nil
}
