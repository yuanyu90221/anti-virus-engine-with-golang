// Package yara 提供透過 YARA CLI subprocess 執行規則比對的偵測引擎。
//
// YARA 為可選功能：僅在呼叫端提供 --yara-rules 旗標時啟用。
// 此套件不依賴 CGo 或任何外部 Go binding，保持 pure Go 可攜性。
package yara

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/yuanyu90221/avengine/internal/scanner"
)

// EngineName 是此引擎在 Detection.Engine 中使用的識別名稱。
const EngineName = "yara"

// Engine 透過 YARA CLI binary 執行規則比對。
// 所有欄位在建構後為唯讀，並發使用安全。
type Engine struct {
	rulesPath string // .yar 規則檔案或包含 .yar/.yara 檔案的目錄路徑
	yaraPath  string // yara binary 的絕對路徑
}

// New 解析 yara binary 路徑並驗證 rulesPath 可用。
// 若找不到 binary 則回傳 error，讓呼叫端可警告並降級為 hash-only 模式。
func New(rulesPath string) (*Engine, error) {
	yaraPath, err := exec.LookPath("yara")
	if err != nil {
		return nil, fmt.Errorf("yara: binary not found in PATH: %w", err)
	}
	return &Engine{rulesPath: rulesPath, yaraPath: yaraPath}, nil
}

// NewWithBinary 使用指定的 binary 路徑建立 Engine，用於測試注入假 binary。
func NewWithBinary(binaryPath, rulesPath string) *Engine {
	return &Engine{rulesPath: rulesPath, yaraPath: binaryPath}
}

// Name 實作 scanner.DetectionEngine 介面。
func (e *Engine) Name() string { return EngineName }

// resolveRuleFiles 回傳要傳給 yara CLI 的規則檔案路徑清單。
// 若 rulesPath 為目錄，則展開目錄下所有 .yar 檔案；否則直接回傳單一路徑。
func resolveRuleFiles(rulesPath string) ([]string, error) {
	info, err := os.Stat(rulesPath)
	if err != nil {
		return nil, fmt.Errorf("yara: rules path: %w", err)
	}
	if !info.IsDir() {
		return []string{rulesPath}, nil
	}
	entries, err := filepath.Glob(filepath.Join(rulesPath, "*.yar"))
	if err != nil {
		return nil, fmt.Errorf("yara: glob rules dir: %w", err)
	}
	if len(entries) == 0 {
		return nil, fmt.Errorf("yara: no .yar files found in directory %s", rulesPath)
	}
	return entries, nil
}

// Inspect 實作 scanner.DetectionEngine 介面。
// 執行：yara <rulesPath> <filePath>（rulesPath 為目錄時展開為多個規則檔案）
// 並解析 stdout 中格式為 "RuleName /abs/path" 的每一行。
//
// YARA CLI 結束碼語意：
//   - 0：規則編譯成功且有比對
//   - 1：無比對（非錯誤）
//   - 2+：YARA 錯誤（規則語法錯誤、檔案無法讀取等）
func (e *Engine) Inspect(ctx context.Context, filePath string) ([]scanner.EngineDetection, error) {
	ruleFiles, err := resolveRuleFiles(e.rulesPath)
	if err != nil {
		return nil, err
	}
	args := append(ruleFiles, filePath)
	cmd := exec.CommandContext(ctx, e.yaraPath, args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			switch exitErr.ExitCode() {
			case 1:
				// 無比對，非錯誤
				return nil, nil
			default:
				return nil, fmt.Errorf("yara: subprocess error (exit %d): %s",
					exitErr.ExitCode(), strings.TrimSpace(stderr.String()))
			}
		}
		// context deadline 或 signal
		return nil, fmt.Errorf("yara: exec: %w", err)
	}

	// exit 0 = 有比對，解析 stdout
	return parseOutput(stdout.String()), nil
}

// parseOutput 解析 YARA CLI stdout，每行格式為 "RuleName /path/to/file"。
// Category 預設為 "yara"，Severity 預設為 "unknown"（CLI 輸出不含 metadata）。
func parseOutput(output string) []scanner.EngineDetection {
	var results []scanner.EngineDetection
	sc := bufio.NewScanner(strings.NewReader(output))
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, " ", 2)
		if len(parts) != 2 {
			continue
		}
		results = append(results, scanner.EngineDetection{
			Name:     parts[0],
			Category: "yara",
			Severity: "unknown",
		})
	}
	return results
}
