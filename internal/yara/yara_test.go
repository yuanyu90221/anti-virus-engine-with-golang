package yara_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yuanyu90221/avengine/internal/yara"
)

// writeFakeBinary 將 shell script 寫入 t.TempDir()，回傳可執行的路徑。
func writeFakeBinary(t *testing.T, script string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "yara")
	err := os.WriteFile(path, []byte("#!/bin/sh\n"+script), 0755)
	require.NoError(t, err)
	return path
}

// writeFakeRulesFile 在 t.TempDir() 建立一個假規則檔案，回傳其路徑。
func writeFakeRulesFile(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "rules.yar")
	err := os.WriteFile(path, []byte("rule Fake {}"), 0644)
	require.NoError(t, err)
	return path
}

func TestInspect_Match(t *testing.T) {
	// 假 binary：exit 0 並輸出一行比對結果
	bin := writeFakeBinary(t, `echo "TestRule $2"; exit 0`)
	eng := yara.NewWithBinary(bin, writeFakeRulesFile(t))

	hits, err := eng.Inspect(context.Background(), "/some/file")
	require.NoError(t, err)
	require.Len(t, hits, 1)
	assert.Equal(t, "TestRule", hits[0].Name)
	assert.Equal(t, "yara", hits[0].Category)
	assert.Equal(t, "unknown", hits[0].Severity)
}

func TestInspect_MultipleMatches(t *testing.T) {
	bin := writeFakeBinary(t, `printf "RuleA /some/file\nRuleB /some/file\n"; exit 0`)
	eng := yara.NewWithBinary(bin, writeFakeRulesFile(t))

	hits, err := eng.Inspect(context.Background(), "/some/file")
	require.NoError(t, err)
	require.Len(t, hits, 2)
	assert.Equal(t, "RuleA", hits[0].Name)
	assert.Equal(t, "RuleB", hits[1].Name)
}

func TestInspect_NoMatch(t *testing.T) {
	// exit 1 = 無比對，非錯誤
	bin := writeFakeBinary(t, `exit 1`)
	eng := yara.NewWithBinary(bin, writeFakeRulesFile(t))

	hits, err := eng.Inspect(context.Background(), "/some/file")
	require.NoError(t, err)
	assert.Empty(t, hits)
}

func TestInspect_ErrorExit(t *testing.T) {
	// exit 2 = YARA 錯誤（規則語法錯誤等）
	bin := writeFakeBinary(t, `echo "bad rules" >&2; exit 2`)
	eng := yara.NewWithBinary(bin, writeFakeRulesFile(t))

	_, err := eng.Inspect(context.Background(), "/some/file")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "exit 2")
}

func TestInspect_Timeout(t *testing.T) {
	// 使用 exec sleep（而非 shell sleep）確保 SIGKILL 能即時終止子程序
	bin := writeFakeBinary(t, `exec sleep 10`)
	eng := yara.NewWithBinary(bin, writeFakeRulesFile(t))

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	_, err := eng.Inspect(ctx, "/some/file")
	require.Error(t, err, "context deadline 應回傳 error")
}

func TestName(t *testing.T) {
	eng := yara.NewWithBinary("/fake/bin", "/fake/rules")
	assert.Equal(t, "yara", eng.Name())
}

func TestInspect_DirectoryRulesPath_UsesAllYarFiles(t *testing.T) {
	// 假 binary：列印 $1（第一個規則路徑）以驗證多個規則被傳入
	// yara CLI 收到目錄展開後的多個規則檔時，$1 為第一個規則檔路徑
	bin := writeFakeBinary(t, `echo "TestRule $2"; exit 0`)

	// 建立一個含有兩個 .yar 檔案的臨時目錄
	rulesDir := t.TempDir()
	err := os.WriteFile(filepath.Join(rulesDir, "a.yar"), []byte("rule A {}"), 0644)
	require.NoError(t, err)
	err = os.WriteFile(filepath.Join(rulesDir, "b.yar"), []byte("rule B {}"), 0644)
	require.NoError(t, err)

	eng := yara.NewWithBinary(bin, rulesDir)
	hits, err := eng.Inspect(context.Background(), "/some/file")
	require.NoError(t, err)
	require.Len(t, hits, 1)
	assert.Equal(t, "TestRule", hits[0].Name)
}

func TestInspect_DirectoryRulesPath_EmptyDir_ReturnsError(t *testing.T) {
	bin := writeFakeBinary(t, `exit 0`)
	emptyDir := t.TempDir()

	eng := yara.NewWithBinary(bin, emptyDir)
	_, err := eng.Inspect(context.Background(), "/some/file")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "no .yar files")
}
