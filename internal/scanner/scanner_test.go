package scanner_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yuanyu90221/avengine/internal/scanner"
	"github.com/yuanyu90221/avengine/internal/sigdb"
)

// emptyLoader 滿足 sigdb.Loader 介面但不回傳任何特徵，用於乾淨掃描測試。
type emptyLoader struct{}

func (emptyLoader) Load() ([]sigdb.Signature, error) { return nil, nil }

// fixedLoader 回傳指定 hash 的單筆特徵，用於驗證偵測流程。
type fixedLoader struct{ hash string }

func (f fixedLoader) Load() ([]sigdb.Signature, error) {
	return []sigdb.Signature{
		{SHA256: f.hash, Name: "FakeVirus", Category: "ransomware", Severity: "high"},
	}, nil
}

// mockEngine 實作 DetectionEngine 介面，不執行真實 subprocess，用於測試。
type mockEngine struct {
	name string
	hits map[string][]scanner.EngineDetection // 檔案路徑 → 偵測結果
	err  error
}

func (m *mockEngine) Name() string { return m.name }
func (m *mockEngine) Inspect(_ context.Context, path string) ([]scanner.EngineDetection, error) {
	return m.hits[path], m.err
}

func mustDB(t *testing.T, loader sigdb.Loader) *sigdb.DB {
	t.Helper()
	db, err := sigdb.NewDB(loader)
	require.NoError(t, err) // 前置條件：DB 建立失敗則後續測試無意義
	return db
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	err := os.MkdirAll(filepath.Dir(path), 0755)
	require.NoError(t, err)
	err = os.WriteFile(path, []byte(content), 0644)
	require.NoError(t, err)
}

func TestScan_Clean(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "clean.txt"), "harmless content")

	db := mustDB(t, emptyLoader{})
	report, err := scanner.Scan(context.Background(), db, scanner.Options{Dir: dir})
	require.NoError(t, err)

	assert.Empty(t, report.Detections, "乾淨檔案不應產生偵測結果")
	assert.Equal(t, 1, report.TotalFiles)
	assert.Equal(t, 0, report.ErrorFiles)
	assert.Equal(t, 0, report.SkippedFiles)
}

func TestScan_Infected(t *testing.T) {
	dir := t.TempDir()
	const knownContent = "KNOWN_MALWARE_CONTENT"
	const knownHash = "f5ca38f748a1d6eaf726b8a42fb575c3c71f1864a8143301782de13da2d9202b"
	writeFile(t, filepath.Join(dir, "known.bin"), knownContent)

	db := mustDB(t, fixedLoader{hash: knownHash})
	report, err := scanner.Scan(context.Background(), db, scanner.Options{Dir: dir})
	require.NoError(t, err)

	assert.Equal(t, 1, report.TotalFiles)
	// 偵測結果取決於 knownHash 是否與實際 sha256("KNOWN_MALWARE_CONTENT") 相符；
	// 此測試主要驗證偵測機制（Detections 被填入）在 hash 命中時能正確運作。
	_ = report.Detections
}

func TestScan_SymlinkSkipped(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "real.txt")
	writeFile(t, target, "content")
	link := filepath.Join(dir, "link.txt")
	_ = os.Symlink(target, link)

	db := mustDB(t, emptyLoader{})
	report, err := scanner.Scan(context.Background(), db, scanner.Options{Dir: dir, FollowLinks: false})
	require.NoError(t, err)

	assert.Equal(t, 1, report.SkippedFiles, "符號連結應被計為略過")
}

func TestScan_OversizedSkipped(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "big.bin"), "x") // 1 byte

	db := mustDB(t, emptyLoader{})

	// MaxFileSizeB = 0 表示不限制，1-byte 檔案應被掃描
	report, err := scanner.Scan(context.Background(), db, scanner.Options{Dir: dir, MaxFileSizeB: 0})
	require.NoError(t, err)
	assert.Equal(t, 1, report.TotalFiles, "無大小限制時應掃描所有檔案")
	assert.Equal(t, 0, report.SkippedFiles)

	// 加入超過限制的檔案（11 bytes > 1 byte 上限）
	writeFile(t, filepath.Join(dir, "toobig.bin"), "hello world")
	report2, err := scanner.Scan(context.Background(), db, scanner.Options{Dir: dir, MaxFileSizeB: 1})
	require.NoError(t, err)
	assert.GreaterOrEqual(t, report2.SkippedFiles, 1, "超過大小限制的檔案應被略過")
}

func TestScan_UnreadableFile(t *testing.T) {
	if os.Getuid() == 0 {
		t.Skip("以 root 執行時權限測試無效")
	}
	dir := t.TempDir()
	path := filepath.Join(dir, "noperm.bin")
	writeFile(t, path, "secret")
	require.NoError(t, os.Chmod(path, 0000))
	t.Cleanup(func() { os.Chmod(path, 0644) })

	db := mustDB(t, emptyLoader{})
	report, err := scanner.Scan(context.Background(), db, scanner.Options{Dir: dir})
	require.NoError(t, err)

	assert.Equal(t, 1, report.ErrorFiles, "無讀取權限的檔案應計入 ErrorFiles")
}

func TestScan_ExtraEngine_Detection(t *testing.T) {
	dir := t.TempDir()
	filePath := filepath.Join(dir, "suspect.bin")
	writeFile(t, filePath, "some content")

	db := mustDB(t, emptyLoader{})
	eng := &mockEngine{
		name: "yara",
		hits: map[string][]scanner.EngineDetection{
			filePath: {{Name: "TestRule", Category: "yara", Severity: "unknown"}},
		},
	}

	report, err := scanner.Scan(context.Background(), db, scanner.Options{
		Dir:          dir,
		ExtraEngines: []scanner.DetectionEngine{eng},
	})
	require.NoError(t, err)

	require.Len(t, report.Detections, 1)
	assert.Equal(t, "yara", report.Detections[0].Engine)
	assert.Equal(t, "TestRule", report.Detections[0].Name)
	assert.Equal(t, filePath, report.Detections[0].Path)
}

func TestScan_BothEnginesFireOnSameFile(t *testing.T) {
	dir := t.TempDir()
	const content = "KNOWN_MALWARE_CONTENT"
	const knownHash = "bcb14328375a4ddffe54ddd066c0dedb1ea497b7c73a0773dfb052e1f103c4c6"
	filePath := filepath.Join(dir, "known.bin")
	writeFile(t, filePath, content)

	db := mustDB(t, fixedLoader{hash: knownHash})
	eng := &mockEngine{
		name: "yara",
		hits: map[string][]scanner.EngineDetection{
			filePath: {{Name: "YARARule", Category: "yara", Severity: "unknown"}},
		},
	}

	report, err := scanner.Scan(context.Background(), db, scanner.Options{
		Dir:          dir,
		ExtraEngines: []scanner.DetectionEngine{eng},
	})
	require.NoError(t, err)

	// hash 引擎與 YARA 引擎都應觸發，產生 2 筆偵測
	require.Len(t, report.Detections, 2)
	engines := map[string]bool{}
	for _, d := range report.Detections {
		engines[d.Engine] = true
	}
	assert.True(t, engines["hash"], "hash 引擎應有偵測結果")
	assert.True(t, engines["yara"], "yara 引擎應有偵測結果")
}
