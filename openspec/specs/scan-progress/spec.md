### Requirement: Scanner exposes progress callback
The `Options` struct SHALL include an `OnProgress func(path string, count int64)` field. The scanner SHALL invoke this callback after each file is processed, passing the file path and the cumulative count of files processed so far. If `OnProgress` is nil, the scanner SHALL behave identically to the current implementation.

#### Scenario: Callback invoked per file
- **WHEN** `Options.OnProgress` is set and a scan runs over N files
- **THEN** the callback is called exactly N times, once per processed file, with monotonically increasing count starting at 1

#### Scenario: Nil callback does not panic
- **WHEN** `Options.OnProgress` is nil
- **THEN** the scanner completes without error and produces the same report as before

### Requirement: Scanner provides file count utility
The `scanner` package SHALL expose a `CountFiles(opts Options) (int64, error)` function that walks `opts.Dir` and returns the number of files that would be processed by `Scan` (applying the same FollowLinks and MaxFileSizeB filters). It SHALL return an error only if `opts.Dir` does not exist or cannot be accessed.

#### Scenario: Count matches scannable files
- **WHEN** `CountFiles` is called on a directory with N eligible files
- **THEN** it returns N (matching what `Scan` would report as TotalFiles, excluding errors and skipped files)

#### Scenario: Empty directory returns zero
- **WHEN** `CountFiles` is called on an empty directory
- **THEN** it returns 0 with no error

### Requirement: Progress display on stderr in text mode
In text output mode, the CLI SHALL display a single-line progress indicator on stderr during scanning. The indicator SHALL show the current file count, the total file count, the percentage complete, and the file path being processed, truncated to fit within 80 characters. Each update SHALL overwrite the previous line using a carriage return (`\r`). After scanning completes, the progress line SHALL be cleared before the report is written to stdout. If the total file count cannot be determined (zero or error), the indicator SHALL fall back to the `[N] path` format without percentage.

#### Scenario: Progress updates show percentage
- **WHEN** scanning is in progress in text mode and stderr is a TTY and total > 0
- **THEN** each file update writes `\r[<count>/<total>] (<pct>%) <truncated-path>` to stderr

#### Scenario: Fallback when total is zero
- **WHEN** scanning is in progress and CountFiles returns 0 or an error
- **THEN** the indicator falls back to `\r[<count>] <truncated-path>` without percentage

#### Scenario: Progress cleared before report output
- **WHEN** scanning completes in text mode
- **THEN** stderr receives a clear-line sequence before the report is written to stdout

### Requirement: No progress display in JSON mode or non-TTY
The CLI SHALL NOT write progress output to stderr when the output format is `json`, or when stderr is not a TTY.

#### Scenario: JSON mode suppresses progress
- **WHEN** `--output json` is specified
- **THEN** no progress output appears on stderr

#### Scenario: Non-TTY suppresses progress
- **WHEN** stderr is not a TTY (e.g. piped to a file)
- **THEN** no progress output appears on stderr
