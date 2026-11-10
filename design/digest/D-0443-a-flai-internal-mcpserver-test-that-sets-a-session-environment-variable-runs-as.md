---
id: D-0443
title: A flai/internal/mcpserver test that sets a session environment variable runs as a parallel subtest of a serial parent through together, since the server reads the session from the environment
kind: decision
sources: [S-0548, T-2361, I-0234, S-0684]
topics: [testing, go, mcp, code]
paths: [flai/internal/mcpserver/server_test.go, flai/internal/mcpserver]
date: 2026-10-31
status: active
absorbed_by: ""
asked_again: []
cited: [S-0684, T-2904, T-2906, TH-0996]
---

# D-0443 A flai/internal/mcpserver test that sets a session environment variable runs as a parallel subtest of a serial parent through together, since the server reads the session from the environment

**No longer current: S-0684 took this rule away on 2026-11-10 (D-0571).** The server now takes the session and its test overrides through `mcpserver.Options`, `together` and the serial parents are gone, and `process_test.go` fails a test that sets the process. Follow D-0571, not the rule below.

S-0548's T-2361 made `flai/internal/mcpserver`'s tests run in parallel: 9m14s to 136 s under `-race` at load 40 to 57. The server then read the session's story, role, session, reasons, and started-by from the environment when it was made (`FLAI_STORY`, `FLAI_ROLE`, `FLAI_SESSION`, `FLAI_REASONS`, `FLAI_STARTED_BY`), and a test that set a process variable could not run beside the package's other tests.

What the story's narrative records: "`TestMain` clears the session variables, `together` runs tests that need one variable set as parallel subtests of a serial parent (54 tests renamed `testX`), and each story's prime is a parallel subtest; 41 top-level tests stay serial for a variable of their own" (wip/archive/agents/S-0548.md, Decisions).

Until S-0684, a new `mcpserver` test that needed one of those variables joined the serial parent that already set it, as a `testX` function passed to `together`. I-0234 recorded the remedy, the session in `Options`, which S-0684 built.

Sources: S-0548 (T-2361), I-0234, S-0684.
