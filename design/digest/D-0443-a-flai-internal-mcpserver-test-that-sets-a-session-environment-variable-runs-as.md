---
id: D-0443
title: A flai/internal/mcpserver test that sets a session environment variable runs as a parallel subtest of a serial parent through together, since the server reads the session from the environment
kind: decision
sources: [S-0548, T-2361, I-0234]
topics: [testing, go, mcp, code]
paths: [flai/internal/mcpserver/server_test.go, flai/internal/mcpserver]
date: 2026-10-31
status: active
absorbed_by: ""
asked_again: []
cited: [S-0684, T-2904, T-2906, TH-0996]
---

# D-0443 A flai/internal/mcpserver test that sets a session environment variable runs as a parallel subtest of a serial parent through together, since the server reads the session from the environment

S-0548's T-2361 made `flai/internal/mcpserver`'s tests run in parallel: 9m14s to 136 s under `-race` at load 40 to 57. The server reads the session's story, role, session, reasons, and started-by from the environment when it is made (`FLAI_STORY`, `FLAI_ROLE`, `FLAI_SESSION`, `FLAI_REASONS`, `FLAI_STARTED_BY`), and a test that sets a process variable cannot run beside the package's other tests.

What the story's narrative records: "`TestMain` clears the session variables, `together` runs tests that need one variable set as parallel subtests of a serial parent (54 tests renamed `testX`), and each story's prime is a parallel subtest; 41 top-level tests stay serial for a variable of their own" (wip/archive/agents/S-0548.md, Decisions). `together` in `flai/internal/mcpserver/server_test.go` names each subtest for its function, `testX` running as `X`.

So a new `mcpserver` test that sets none of those variables, no git identity, and no package variable such as `permissionPoll`, `flaiExecutable`, or `buildinfo.Version` calls `t.Parallel()`. One that needs one of them joins the serial parent that already sets it, such as `TestInAGitSession` or `TestAsTheStorysAgent`, as a `testX` function passed to `together`. It does not start a new serial top-level test. The serial phase is still about 70 s of the package's time. I-0234 records the remedy: the session goes in `Options`, so the serial phase goes.

Sources: S-0548 (T-2361), I-0234.
