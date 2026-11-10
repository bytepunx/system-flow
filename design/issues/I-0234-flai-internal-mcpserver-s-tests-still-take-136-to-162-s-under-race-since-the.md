---
id: I-0234
title: flai/internal/mcpserver's tests still take 136 to 162 s under -race, since the server reads the session's story, role, session, reasons, and started-by from the environment and 95 tests must run serially or in env groups before the rest
class: efficiency
status: closed
count: 1
cost: 1m
first_reported: 2026-10-31T09:50:04Z
last_reported: 2026-10-31T09:50:04Z
updated: 2026-11-10T10:34:20Z
---

# I-0234 flai/internal/mcpserver's tests still take 136 to 162 s under -race, since the server reads the session's story, role, session, reasons, and started-by from the environment and 95 tests must run serially or in env groups before the rest

## Description
flai/internal/mcpserver's tests still take 136 to 162 s under -race, since the server reads the session's story, role, session, reasons, and started-by from the environment and 95 tests must run serially or in env groups before the rest

## Instances

### 2026-10-31T09:50:04Z
Story: S-0548.
S-0548 made the package's tests parallel: 9m14s to 136 s (load 40 to 57), 141 s (14 to 31), 162 s (28 to 55), against criterion 4's two minutes. The critical path is about 70 s of serial tests and env groups (TestInAGitSession, TestAsTheStorysAgent, and the like, which set FLAI_STORY, FLAI_ROLE, FLAI_SESSION, FLAI_REASONS, FLAI_STARTED_BY, git identity, or permissionPoll, flaiExecutable, buildinfo.Version), then the two prime-on-this-repository tests, which wait for that phase and take 55 to 70 s of wall time and about 550 user-seconds under race in ctxpack.ForStory. The server reads those variables at New (flai/internal/mcpserver/server.go:203-209, shared.go:180, issues.go:196). Remedy: Options carry the session's story, role, session, reasons, and started-by, filled from the environment by flai mcp, and the tests pass them, so the serial phase goes and the package's time becomes the prime tests'; or prime the repository's stories faster under race. A search change in Selection.refs measured no gain under race and was reverted (T-2420, cancelled).

## Remediation

Story S-0684 remediates this issue, created from it at 2026-11-10T08:37:09Z.
Closed 2026-11-10T10:34:20Z: S-0684 fixed it. The server takes the session (Options.Story, Role, Session, Reasons, StartedBy) through mcpserver.Options, which flai mcp fills with SessionFromEnv. It also takes its test overrides that way: PermissionPoll, FlaiExecutable, SubprojectEnv, Releases, CloseOutWait, and Version. So no test of flai/internal/mcpserver sets the process or runs serially; together and the serial parents are gone, and process_test.go fails on any test that sets the process or skips t.Parallel. Measured with go test -race -count=1 on a 24-core host at load 53 to 156, alternating with main under the same load: 317.7, 200.0, and 155.5 s against main's 354.0, 352.7, and 414.6 s, and 57.2 s against 177.6 s with -short (T-2915's Notes).
