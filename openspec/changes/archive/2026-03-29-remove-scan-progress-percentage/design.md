## Context

The CLI currently runs a `CountFiles` pre-scan before the main `Scan` call so it can display `[count/total] (pct%)` progress. The percentage was added as a UX nicety but the total-count walk adds latency before scanning begins and the percentage value is not essential for the user to track progress.

The existing fallback path (`[N] path` format when `CountFiles` returns 0) is already implemented and tested; this change promotes it to the sole format.

## Goals / Non-Goals

**Goals:**
- Remove the `CountFiles` pre-scan from the CLI to eliminate the startup delay
- Simplify the progress closure in `cmd/avengine/main.go` to only track `count` and `path`
- Keep the `CountFiles` function in the `scanner` package (it has its own spec and may be used elsewhere)

**Non-Goals:**
- Changing the `OnProgress` callback signature — it already passes `(path string, count int64)`
- Modifying JSON mode or non-TTY suppression behaviour
- Updating the `CountFiles` unit tests

## Decisions

**Decision: Promote the fallback format as the only format**

The `[N] path` fallback already exists and is correct. Rather than adding a new format string, we simply remove the `if total > 0` branch and the `CountFiles` call, leaving the fallback as the sole code path. This minimises the diff and eliminates dead code.

Alternative considered: keep `CountFiles` but run it concurrently with `Scan`. Rejected — unnecessary complexity for a display-only feature.

## Risks / Trade-offs

- Users who relied on percentage to estimate remaining time will lose that indicator. → Acceptable; file count still conveys progress.
- `CountFiles` stays in the codebase as an unused CLI call. → Its spec requires it to exist; it may be useful for future features or tests.
