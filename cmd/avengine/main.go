// Package main 是 avengine 防毒引擎的命令列入口點。
//
// 使用方式：
//
//	avengine scan --dir <目標目錄> [旗標...]
//
// 整體流程：
//  1. 解析 CLI 旗標（使用標準庫 flag.FlagSet，支援 --help 自動產生）
//  2. sigdb.YAMLLoader 從特徵目錄載入所有 .yaml 並建立記憶體索引
//  3. scanner.Scan 遞迴走訪目標目錄，對每個檔案計算 SHA256 並查詢索引
//  4. reporter.Write 將結果以文字表格或 JSON 輸出至 stdout
//  5. 依掃描結果以對應結束碼退出（0=乾淨 / 1=偵測到威脅 / 2=執行錯誤）
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/yuanyu90221/avengine/internal/config"
	"github.com/yuanyu90221/avengine/internal/reporter"
	"github.com/yuanyu90221/avengine/internal/scanner"
	"github.com/yuanyu90221/avengine/internal/sigdb"
	"github.com/yuanyu90221/avengine/internal/yara"
)

// isTerminal 回報 f 是否連接到互動式終端機（TTY）。
// 使用 os.ModeCharDevice 旗標判斷，相容 Linux / macOS，無需外部套件。
func isTerminal(f *os.File) bool {
	fi, err := f.Stat()
	if err != nil {
		return false
	}
	return fi.Mode()&os.ModeCharDevice != 0
}

func main() {
	// 只支援 "scan" 子命令；未來若需擴充（如 "update-sigs"），可在此新增 switch
	if len(os.Args) < 2 || os.Args[1] != "scan" {
		fmt.Fprintln(os.Stderr, "usage: avengine scan [flags]")
		os.Exit(reporter.ExitError)
	}

	// 使用獨立的 FlagSet 而非 flag.Parse()，
	// 讓 "scan" 子命令有自己的旗標命名空間，避免與未來其他子命令衝突
	fs := flag.NewFlagSet("scan", flag.ContinueOnError)
	fs.SetOutput(os.Stderr) // 錯誤訊息輸出至 stderr，不污染 stdout 的報告輸出

	dir := fs.String("dir", "", "target directory to scan (required)")
	sigs := fs.String("sigs", "./signatures", "signatures YAML directory")
	output := fs.String("output", "text", "output format: text|json")
	followLinks := fs.Bool("follow-links", false, "follow symbolic links")
	maxSizeMB := fs.Int("max-size", 0, "skip files larger than N MB (0 = no limit)")
	_ = fs.Bool("verbose", false, "show all scanned files") // 預留旗標，尚未實作詳細模式
	yaraRules := fs.String("yara-rules", "", "path to YARA rules file or directory (optional)")
	yaraTimeout := fs.Int("yara-timeout", config.DefaultYARATimeout, "per-file YARA timeout in seconds")

	if err := fs.Parse(os.Args[2:]); err != nil {
		os.Exit(reporter.ExitError)
	}

	// --dir 為必填旗標，缺少時輸出使用說明至 stderr
	if *dir == "" {
		fmt.Fprintln(os.Stderr, "error: --dir is required")
		fs.Usage()
		os.Exit(reporter.ExitError)
	}

	// 建立 Reporter（在實際 IO 前先驗證格式，避免掃描完才發現格式錯誤）
	rep, err := reporter.New(*output)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(reporter.ExitError)
	}

	// 載入特徵資料庫：YAMLLoader 讀取目錄中所有 .yaml 並建立 hash → MatchResult 索引
	loader := sigdb.YAMLLoader{Dir: *sigs}
	db, err := sigdb.NewDB(loader)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading signatures: %v\n", err)
		os.Exit(reporter.ExitError)
	}

	// 將 --max-size（MB 整數）轉換為位元組，傳入 scanner
	// MaxFileSizeB = 0 表示不限制（scanner 內部以 > 0 判斷是否套用限制）
	var maxBytes int64
	if *maxSizeMB > 0 {
		maxBytes = int64(*maxSizeMB) * 1024 * 1024
	}

	// 僅在 text 模式且 stderr 為 TTY 時顯示進度（避免污染 json pipeline）
	showProgress := *output == "text" && isTerminal(os.Stderr)
	var onProgress scanner.ProgressFunc
	if showProgress {
		onProgress = func(path string, count int64) {
			display := path
			if len(display) > 70 {
				display = "..." + display[len(display)-67:]
			}
			fmt.Fprintf(os.Stderr, "\r[%d] %s", count, display)
		}
	}

	// 初始化額外偵測引擎（目前支援 YARA）
	var extraEngines []scanner.DetectionEngine
	if *yaraRules != "" {
		eng, engErr := yara.New(*yaraRules)
		if engErr != nil {
			// binary 找不到或規則路徑無效 → 警告後以 hash-only 繼續
			fmt.Fprintf(os.Stderr, "warning: YARA engine unavailable: %v\n", engErr)
		} else {
			extraEngines = append(extraEngines, eng)
		}
	}

	// 執行掃描：遞迴走訪 --dir，逐檔計算 SHA256 並查詢 db
	report, err := scanner.Scan(context.Background(), db, scanner.Options{
		Dir:          *dir,
		FollowLinks:  *followLinks,
		MaxFileSizeB: maxBytes,
		OnProgress:   onProgress,
		ExtraEngines: extraEngines,
		FileTimeout:  time.Duration(*yaraTimeout) * time.Second,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "scan error: %v\n", err)
		os.Exit(reporter.ExitError)
	}

	// 清除進度列，讓報告從乾淨的行開始輸出
	if showProgress {
		fmt.Fprint(os.Stderr, "\r\033[K")
	}

	// 輸出報告至 stdout（與錯誤訊息分離，便於 shell 管線處理）
	if err := rep.Write(os.Stdout, report); err != nil {
		fmt.Fprintf(os.Stderr, "output error: %v\n", err)
		os.Exit(reporter.ExitError)
	}

	// 結束碼語意：讓 CI/CD 腳本可直接以 $? 判斷掃描結果
	if len(report.Detections) > 0 {
		os.Exit(reporter.ExitDetected) // 1：偵測到威脅
	}
	os.Exit(reporter.ExitClean) // 0：乾淨
}
