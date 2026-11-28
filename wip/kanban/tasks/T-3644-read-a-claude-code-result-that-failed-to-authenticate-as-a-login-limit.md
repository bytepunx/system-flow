---
id: T-3644
type: task
nature: remediation
title: Read a claude-code result that failed to authenticate as a login limit
status: backlog
parent: S-0855
owner: alex
created: 2026-11-27T23:59:39Z
updated: 2026-11-27T23:59:39Z
transitions: []
stream: S-0855
tags: [flai]
touches: [flai/internal/usage/event.go, flai/internal/usage/claudecode.go, flai/internal/usage/claudecode_test.go, flai/internal/serve/testdata/auth-expired-result.jsonl]
after: [T-3643]
---
# T-3644 Read a claude-code result that failed to authenticate as a login limit

## Work

Waits for T-3643, since the ADR decides which result fields mark a failed login and what kind the limit carries.

- Give `usage.Limit` in `flai/internal/usage/event.go` the kind the ADR names, so that a login limit is told apart from a usage limit.
- In `flai/internal/usage/claudecode.go` (`claudeCodeReader.ended`), return such a limit for a result the ADR names as a failed login. Update the reader's contract comment at the top of the file.
- Add `flai/internal/serve/testdata/auth-expired-result.jsonl`. Copy it from a real failed run's log, with its session and any account details replaced.
- Add cases to `claudecode_test.go`: the failed login is read as a login limit, a 429 is still a usage limit, and a result with another error is no limit.

## Done when

- `flai test` on the usage package passes, with the new cases.
- The fixture holds no token, no email, and no account ID.

## Notes
