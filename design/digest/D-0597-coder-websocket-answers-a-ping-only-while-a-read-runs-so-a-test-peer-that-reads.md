---
id: D-0597
title: "coder/websocket answers a ping only while a read runs, so a test peer that reads its connection only on demand lets the client's ping time out under load and close it: channeltest.Dashboard now reads each connection for its whole life (I-0267)"
kind: cause
sources: [I-0267, S-0703, T-3044, T-3042]
topics: [testing, go, flai]
paths: [flai/internal/channel/channeltest/dashboard.go, flai/internal/serve/candidates_test.go, flai/internal/hostapi/components_test.go, flai/internal/serve/serve_test.go]
date: 2026-11-11
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0597 coder/websocket answers a ping only while a read runs, so a test peer that reads its connection only on demand lets the client's ping time out under load and close it: channeltest.Dashboard now reads each connection for its whole life (I-0267)

I-0267: `TestARepositoryIsOfferedImportedAndThenServed` failed S-0580's close-out at integration with `failed to get reader: failed to read frame header: EOF` on a host at load 56, on a branch that does not touch `flai/internal/serve`.

S-0703 found the cause in the test dashboard, not in flai serve. From I-0267's close reason:

```text
channeltest.Dashboard read a proven connection only inside Ask, Next, and Last, and coder/websocket answers a ping only while a read runs, so the client's 4 s ping went unanswered whenever a loaded host left a connection unread past its tick and timeout; the client then closed the connection, and the next Ask read 'failed to read frame header: EOF'. The dashboard now reads each connection for its whole life and hands its messages to Ask, Next, and Last (T-3044).
```

So a test peer built on `coder/websocket` that should stay connected must keep a read running for the connection's whole life, not only when it wants a message. The reproduction is a seam, not load: `TestACandidateLeftQuietPastItsPingsAnswersImportRun` pings every 300 ms and leaves the candidate quiet for four pings before `import.run`; it failed 3 of 3 runs against the old dashboard and passes with the fix (T-3042), as D-0550 describes for load remediations.

T-3044 also dropped the comments in `flai/internal/hostapi/components_test.go` and `flai/internal/serve/serve_test.go` saying the test dashboard answers no ping, which are no longer true; their `PingEvery: time.Hour` settings stay, as harmless. The other path to the same EOF, flai serve halting a candidate under its `import.run` answer, was closed earlier by S-0560.

Sources: I-0267 (Remediation); S-0703 (`wip/archive/agents/S-0703.md`, Decisions); T-3044; T-3042.
