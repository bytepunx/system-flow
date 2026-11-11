---
id: D-0571
title: "A flai/internal/mcpserver test passes the session and its overrides through mcpserver.Options and calls t.Parallel; it never sets the process environment or a package variable, and process_test.go fails one that does (I-0234)"
kind: decision
sources: [S-0684, I-0234, T-2904, T-2905, T-2906, T-2910, T-2915, S-0665, T-3079]
topics: [testing, go, mcp, code]
paths: [flai/internal/mcpserver, flai/internal/mcpserver/server.go, flai/internal/mcpserver/server_test.go, flai/cmd/mcp.go, flai/cmd/mcp_http.go, flai/internal/mcpserver/stream_sections_test.go]
date: 2026-11-10
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0571 A flai/internal/mcpserver test passes the session and its overrides through mcpserver.Options and calls t.Parallel; it never sets the process environment or a package variable, and process_test.go fails one that does (I-0234)

S-0684 took the remedy I-0234 named, and D-0443 recorded as pending: the server no longer reads the session from the environment when it is made. I-0234's close reason (`design/issues/I-0234-flai-internal-mcpserver-s-tests-still-take-136-to-162-s-under-race-since-the.md`, Remediation):

```text
The server takes the session (Options.Story, Role, Session, Reasons, StartedBy) through mcpserver.Options, which flai mcp fills with SessionFromEnv. It also takes its test overrides that way: PermissionPoll, FlaiExecutable, SubprojectEnv, Releases, CloseOutWait, and Version. So no test of flai/internal/mcpserver sets the process or runs serially; together and the serial parents are gone, and process_test.go fails on any test that sets the process or skips t.Parallel.
```

So a new `mcpserver` test gives what it needs as `Options` fields and calls `t.Parallel()`. It does not set `FLAI_STORY`, `FLAI_ROLE`, `FLAI_SESSION`, `FLAI_REASONS`, `FLAI_STARTED_BY`, or a package variable, and there is no `together` or serial parent to join any more. A new process-wide setting the server needs becomes an `Options` field that `flai mcp` fills, as S-0684's T-2910 did for `CloseOutWait` and the version check. `Options.Role` is filled by `SessionFromEnv`, and the HTTP server overrides it only when `mcpHTTPRole` names the orchestrator: S-0684 and S-0687 agreed on that one field in MS-0454 (`wip/archive/agents/S-0684.md`, Decisions).

Measured with `go test -race -count=1` at load 53 to 156, alternating with main under the same load: 317.7, 200.0, and 155.5 s against main's 354.0, 352.7, and 414.6 s; 57.2 s against 177.6 s with `-short` (T-2915's Notes).

The rule bit at once on a test written against the old one. S-0668 and S-0684, accepted apart, left main red: `TestNoTestSetsTheProcess` refused `TestStreamStateWritesContextDecisionsAndQuestions`'s `t.Setenv`. S-0665's T-3079 fixed it the same way (`wip/archive/agents/S-0665.md`, Decisions):

```text
The test now calls `t.Parallel`, since `streamFixture` already gives the session through `Options`. Bumped I-0229, the issue for stories accepted apart leaving main red.
```

Sources: S-0684 (T-2904, T-2905, T-2906, T-2910, T-2915), I-0234, S-0665 (T-3079). Replaces the rule in D-0443.
