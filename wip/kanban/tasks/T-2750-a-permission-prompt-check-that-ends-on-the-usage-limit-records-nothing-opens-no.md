---
id: T-2750
type: task
nature: remediation
title: A permission_prompt check that ends on the usage limit records nothing, opens no thread, and is checked again after the reset
status: done
parent: S-0629
owner: alex
created: 2026-11-09T17:09:05Z
updated: 2026-11-09T17:42:10Z
transitions:
  - to: ready
    at: 2026-11-09T17:27:23Z
    by: agent-S-0629
  - to: in-progress
    at: 2026-11-09T17:27:23Z
    by: agent-S-0629
  - to: done
    at: 2026-11-09T17:42:10Z
    by: agent-S-0629
stream: S-0629
tags: [flai, serve, claude-code]
touches: [flai/internal/serve/claudecheck.go, flai/internal/serve/claudecheck_test.go, flai/internal/usage/log.go, flai/internal/usage/log_test.go]
after: [T-2749]
usage:
  source: log
  seconds: 1
  estimated: true
  priced_by: estimate
  models: []
---
# T-2750 A permission_prompt check that ends on the usage limit records nothing, opens no thread, and is checked again after the reset

## Work

Build the rule T-2749's ADR decides. This task waits for T-2749 because the code follows that decision, and the decision changes ADR-0106 § 2.

- In `flai/internal/serve/claudecheck.go`, read the usage limit from the check's `claude -p` output. That output is `--output-format stream-json`. Use the reader `flai/internal/usage` already has for ADR-0214 (`ended`, `rateLimit`), not a match on Claude Code's text. `usage.Read` takes only file paths, so add a variant in `flai/internal/usage/log.go` that reads from bytes or an `io.Reader`.
- `run` and `runInImage` pass the limit back beside `why`. On a limit, `check` and `checkImage` call neither `conclude` nor `tell`. They log `claude code check refused on the usage limit` with the version, the provider, and the reset time.
- Change the in-process gate in `consider` (`c.through`) and its image counterpart so a version refused on the limit is considered again on the first look after the reset, or an hour after the refusal when no reset was given. Today the gate stops it until flai serve restarts.

Reproduce I-0296 in `flai/internal/serve/claudecheck_test.go` first. Give the stand-in claude the fixture `flai/internal/serve/testdata/usage-limit-result.jsonl` (or one like it), and show that today it records `failed` and opens a thread. Then make the test pass:

- No record is written.
- No thread is opened.
- A later look past the reset checks the version again and records its real outcome.

Keep `TestARecordedClaudeCodeIsNotCheckedAgain` and `TestAFailedClaudeCodeCheckOpensOneThreadQuotingItsError` passing for a failure that isn't a limit. Cover the host check and the image check.

## Done when

- A test in `flai/internal/serve/claudecheck_test.go` runs a check that ends on the usage limit and passes. It writes no `claude-checks.json` entry and opens no thread, and a later look checks the version again.
- The image check does the same, with a test.
- A test in `flai/internal/usage/log_test.go` covers the new reader.
- `make flai-test` passes.

## Notes
