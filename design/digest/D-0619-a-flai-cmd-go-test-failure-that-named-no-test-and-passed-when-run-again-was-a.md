---
id: D-0619
title: "A flai/cmd go-test failure that named no test and passed when run again was a test binary that crashed with its output lost; its cause is still unknown, so a recurrence, which go-test-json now names, is filed as a new issue (I-0365)"
kind: cause
sources: [I-0365, S-0782, T-3344, TH-1186]
topics: [testing, flai, go, code]
paths: [flai/internal/verify/gotest.go, flai/internal/verify/parse_test.go, flai/internal/verify/testdata/go-test-crash.json, design/system/project-manifest.md]
date: 2026-11-14
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0619 A flai/cmd go-test failure that named no test and passed when run again was a test binary that crashed with its output lost; its cause is still unknown, so a recurrence, which go-test-json now names, is filed as a new issue (I-0365)

I-0365 was closed on 2026-11-13 by S-0782. Its cause: "flai test reports a go-test failure in flai/cmd that names no test, and the same paths pass when run again".

The finding was empty because the package's test binary died with no test failing and no timeout, as on a goroutine panic, a fatal error, or `os.Exit`. `go test -json` then fails only the package. S-0782 did not find the flake itself. Its narrative (`wip/archive/agents/S-0782.md`, Decisions) says: "the instance kept no output, so the flake itself cannot be named; a readable finding lets a recurrence be recorded and fixed as its own story". The planner and the orchestrator read criterion 1 that way on TH-1186. So a story that meets this finding again reads the test and crash line it now names, and files a new issue. Reopening I-0365 is not the way.

How the `go-test-json` format names a crash is in `design/system/project-manifest.md`, under the tier's `format`, and is not repeated here. The design leaves out two facts from T-3344's recording on go1.26.0 (`flai/internal/verify/testdata/go-test-crash.json`):

- `go test -json` prints no `exit status n` line, and the parser sees only go's own exit. So a crash's message gives its `panic:` or `fatal error:` line alone.
- An `os.Exit` with no crash line reads "the test binary exited before the test ended".

Sources: I-0365 and its `## Remediation`; S-0782's narrative, § Decisions; T-3344; TH-1186.
