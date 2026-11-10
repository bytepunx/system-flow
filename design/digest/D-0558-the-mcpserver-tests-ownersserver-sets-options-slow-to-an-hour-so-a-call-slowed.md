---
id: D-0558
title: The mcpserver tests' ownersServer sets Options.Slow to an hour, so a call slowed by host load logs no slow request event that a test reading the log would trip on (I-0329)
kind: cause
sources: [I-0329, S-0652, T-2843]
topics: [testing, flai, go, code]
paths: [flai/internal/mcpserver/server_test.go, design/issues/I-0329-teststartdevserversayswhatavitethatdoesnotanswer200answered-and.md]
date: 2026-11-10
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0558 The mcpserver tests' ownersServer sets Options.Slow to an hour, so a call slowed by host load logs no slow request event that a test reading the log would trip on (I-0329)

I-0329's second instance: `TestAnOperatorsReplyBecomesADigestEntry` in `flai/internal/mcpserver/server_test.go` failed the integration tier under host load. A `thread_reply` slowed past the server's slow-request threshold logged a "request answered" event, which the test did not expect.

S-0652 fixed it in the test (T-2843). From its Decisions (`wip/archive/agents/S-0652.md`):

```text
T-2843: `ownersServer` sets `Options.Slow` to an hour, so that no call logs a slow "request answered" event. The new `TestTheOwnersServerLogsNoSlowCall` sets `FLAI_SLOW_REQUEST=1ns`, which is not parallel because of `t.Setenv`, and reproduces the second instance on every run.
```

So a test server whose log a test reads sets `Options.Slow` far above any load's delay, and a test that wants to prove it forces the slow path with `FLAI_SLOW_REQUEST=1ns` instead of waiting for load. I-0329's first instance, the dev-server 404 test, is in D-0551.

Sources: I-0329; S-0652 (`wip/archive/agents/S-0652.md`, Decisions); T-2843.
