package hasher_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yuanyu90221/avengine/internal/hasher"
)

func TestHashFile_Normal(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "hashtest")
	require.NoError(t, err) // 前置條件：建立暫存檔失敗則無法繼續

	_, err = f.WriteString("hello world")
	require.NoError(t, err)
	f.Close()

	got, err := hasher.HashFile(f.Name())
	require.NoError(t, err) // 前置條件：hash 失敗則後續比對無意義

	assert.Len(t, got, 64, "SHA256 應回傳 64 字元十六進位字串")
}

func TestHashFile_Empty(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "empty")
	require.NoError(t, err)
	f.Close()

	got, err := hasher.HashFile(f.Name())
	require.NoError(t, err)

	const emptySHA256 = "e3b0c44298fc1c149afbf4c8996fb924" +
		"27ae41e4649b934ca495991b7852b855"
	assert.Equal(t, emptySHA256, got, "空檔案的 SHA256 應符合已知值")
}

func TestHashFile_Missing(t *testing.T) {
	_, err := hasher.HashFile("/nonexistent/path/file.bin")
	assert.Error(t, err, "不存在的路徑應回傳 error")
}
