package reporter_test

import (
	"bytes"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yuanyu90221/avengine/internal/reporter"
	"github.com/yuanyu90221/avengine/internal/scanner"
	"github.com/yuanyu90221/avengine/internal/sigdb"
)

func cleanReport() *scanner.ScanReport {
	now := time.Now()
	return &scanner.ScanReport{
		TotalFiles: 5,
		StartedAt:  now,
		FinishedAt: now.Add(time.Millisecond * 100),
	}
}

func infectedReport() *scanner.ScanReport {
	r := cleanReport()
	r.Detections = []scanner.Detection{
		{
			Path:   "/tmp/bad.bin",
			SHA256: "abcdef1234567890",
			MatchResult: sigdb.MatchResult{
				Name:     "FakeRansom",
				Category: "ransomware",
				Severity: "high",
			},
		},
	}
	return r
}

func TestText_Clean(t *testing.T) {
	rep, err := reporter.New("text")
	require.NoError(t, err) // 前置條件：無效 format 應在此失敗

	var buf bytes.Buffer
	err = rep.Write(&buf, cleanReport())
	require.NoError(t, err)

	assert.Contains(t, buf.String(), "未發現威脅")
}

func TestText_Infected(t *testing.T) {
	rep, err := reporter.New("text")
	require.NoError(t, err)

	var buf bytes.Buffer
	err = rep.Write(&buf, infectedReport())
	require.NoError(t, err)

	out := buf.String()
	assert.Contains(t, out, "FakeRansom", "威脅名稱應出現在文字報告中")
	assert.Contains(t, out, "abcdef12", "SHA256 前綴應出現在文字報告中")
}

func TestJSON_Infected(t *testing.T) {
	rep, err := reporter.New("json")
	require.NoError(t, err)

	var buf bytes.Buffer
	err = rep.Write(&buf, infectedReport())
	require.NoError(t, err)

	var v map[string]any
	err = json.Unmarshal(buf.Bytes(), &v)
	require.NoError(t, err, "輸出應為合法 JSON")

	assert.Contains(t, v, "detections", "JSON 應包含 detections 鍵")
	assert.Contains(t, v, "totalFiles", "JSON 應包含 totalFiles 鍵")
}

func TestNew_InvalidFormat(t *testing.T) {
	_, err := reporter.New("xml")
	assert.Error(t, err, "未知格式應回傳 error")
}
