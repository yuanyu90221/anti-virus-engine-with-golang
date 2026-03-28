package sigdb_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yuanyu90221/avengine/internal/sigdb"
)

func writeYAML(t *testing.T, dir, name, content string) {
	t.Helper()
	err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0644)
	require.NoError(t, err) // 前置條件：寫入失敗則測試無法繼續
}

func TestYAMLLoader_MultiFile(t *testing.T) {
	dir := t.TempDir()
	writeYAML(t, dir, "ransomware.yaml", `
version: "1.0"
category: ransomware
updated: "2024-01-01"
signatures:
  - sha256: aaa111
    name: FakeRansom
    severity: high
    added: "2024-01-01"
`)
	writeYAML(t, dir, "trojans.yaml", `
version: "1.0"
category: trojan
updated: "2024-01-01"
signatures:
  - sha256: bbb222
    name: FakeTrojan
    severity: medium
    added: "2024-01-01"
`)
	loader := sigdb.YAMLLoader{Dir: dir}
	sigs, err := loader.Load()
	require.NoError(t, err)
	assert.Len(t, sigs, 2, "應合併兩個 YAML 檔案共 2 筆特徵")
}

func TestYAMLLoader_MissingDir(t *testing.T) {
	loader := sigdb.YAMLLoader{Dir: "/nonexistent/dir"}
	_, err := loader.Load()
	assert.Error(t, err, "目錄不存在應回傳 error")
}

func TestYAMLLoader_BadYAML(t *testing.T) {
	dir := t.TempDir()
	writeYAML(t, dir, "bad.yaml", `{{{not valid yaml`)
	loader := sigdb.YAMLLoader{Dir: dir}
	_, err := loader.Load()
	assert.Error(t, err, "格式錯誤的 YAML 應回傳 error")
}
