## Why

The current scan progress indicator shows `[count/total] (pct%)`, which requires a pre-scan `CountFiles` pass and adds visual noise. Showing only the file count and current file path is simpler, faster to start, and provides the same actionable feedback to the user.

## What Changes

- Remove the `CountFiles` pre-scan call from CLI startup
- Remove `total`, `pct`, and the `[count/total] (pct%)` display format
- Simplify progress output to `[count] path` format only (the existing fallback format becomes the sole format)
- Remove the now-dead fallback branch in progress rendering

## Capabilities

### New Capabilities
<!-- none -->

### Modified Capabilities
- `scan-progress`: Remove percentage and total-count display from the progress indicator; the `[N] path` format becomes the only format. The `CountFiles` utility in the scanner package may remain (it has its own spec requirement) but the CLI will no longer call it during a scan.

## Impact

- `cmd/avengine/main.go`: Remove `CountFiles` call, remove `total`/`pct` variables, simplify `onProgress` closure
- `openspec/specs/scan-progress/spec.md`: Delta spec to drop percentage-related requirements and scenarios
