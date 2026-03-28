// Package hasher 提供以串流方式計算檔案 SHA256 雜湊值的功能。
//
// 設計重點：使用 io.Copy 將檔案資料串流寫入 hash.Hash，
// 無論檔案多大都只佔用固定的記憶體（hash 內部緩衝區），
// 避免先將整個檔案讀進記憶體再計算的 OOM 風險。
package hasher

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

// HashFile 計算指定路徑檔案的 SHA256 雜湊值。
//
// 實作方式：
//  1. 開啟檔案取得 io.Reader
//  2. 建立 sha256.New()（實作 hash.Hash 與 io.Writer）
//  3. io.Copy 將資料分塊（預設 32 KB）從檔案寫入 hash
//  4. hash.Sum(nil) 取得原始位元組，hex.EncodeToString 轉為 64 字元小寫十六進位字串
//
// 回傳值：64 字元小寫十六進位字串；發生錯誤時回傳空字串與錯誤。
func HashFile(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("hasher: open %s: %w", path, err)
	}
	defer f.Close()

	h := sha256.New()
	// io.Copy 內部使用固定大小緩衝區分批傳輸，不會將整個檔案載入記憶體
	if _, err := io.Copy(h, f); err != nil {
		return "", fmt.Errorf("hasher: read %s: %w", path, err)
	}
	// h.Sum(nil) 將目前 hash 狀態附加到空 slice，取得原始 32 位元組摘要
	return hex.EncodeToString(h.Sum(nil)), nil
}
