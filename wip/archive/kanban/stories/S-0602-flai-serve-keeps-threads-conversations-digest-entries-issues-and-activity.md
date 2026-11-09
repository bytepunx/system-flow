---
id: S-0602
type: story
nature: improvement
title: flai serve keeps threads, conversations, digest entries, issues, and activity documents parsed between requests, so a thread read takes about 5 ms, not 80 to 110 ms
status: done
parent: E-0033
owner: alex
created: 2026-11-07T22:40:05Z
updated: 2026-11-09T07:17:27Z
transitions:
  - to: ready
    at: 2026-11-08T00:11:21Z
    by: orchestrator
  - to: in-progress
    at: 2026-11-08T02:58:53Z
    by: agent-S-0602
  - to: review
    at: 2026-11-09T07:15:48Z
    by: agent-S-0602
  - to: done
    at: 2026-11-09T07:17:27Z
    by: orchestrator
blocked:
  - from: 2026-11-08T03:45:26Z
    until: 2026-11-09T07:15:00Z
    reason: awaiting the operator on TH-0913
tags: [cli]
touches: [flai/internal/threads/threads.go, flai/internal/messages/messages.go, flai/internal/digest/digest.go, flai/internal/issues/issues.go, flai/internal/threads/threads_test.go, flai/internal/messages/messages_test.go, flai/internal/workitem/share.go, flai/internal/workitem/share_test.go, flai/internal/digest/digest_test.go, flai/internal/issues/issues_test.go, flai/internal/workitem/activity.go, flai/internal/workitem/activity_test.go, design/system/server-performance.md, design/system/flai-cli.md, flai/internal/workitem/keptfiles.go, flai/internal/workitem/keptfiles_test.go]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 2843
  priced_by: harness
  turns:
    - day: 2026-11-08
      ceremony: 2
      hand_edits: 2
      work: 86
    - day: 2026-11-09
      ceremony: 1
      hand_edits: 1
      work: 6
  models:
    - model: claude-opus-5-5
      input: 340
      output: 99480
      cache_read: 49409787
      cache_write: 1583134
      cost: 22.3116
  strategic:
    - kind: orchestrator
      seconds: 323
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 51
          output: 300
          cache_read: 2012942
          cache_write: 251463
          cost: 0.5812
        - model: claude-sonnet-5-5
          input: 6
          output: 43
          cache_read: 41353
          cache_write: 36496
          cost: 0.0572
cost_of_delay:
  inputs:
    penalty_per_week: 50
    by: orchestrator
    at: 2026-11-07T23:55:48Z
  value: 50
  by: planner-E-0033
  at: 2026-11-08T00:00:43Z
forecast:
  duration: 1h30m
  delivery: 2026-11-08T04:57:00Z
  basis: "Its own forecast of 1h30m; 1st in the pull order with an in-progress limit of 5, behind S-0603, S-0604, S-0606 and S-0609."
  by: flai
  at: 2026-11-08T02:57:43Z
finalized:
  by: alex
  at: 2026-11-07T23:49:33Z
---
# S-0602 flai serve keeps threads, conversations, digest entries, issues, and activity documents parsed between requests, so a thread read takes about 5 ms, not 80 to 110 ms

## Goal

Cause 13 in [server-performance.md](../../../design/system/server-performance.md#where-the-time-goes): threads (884 files), conversations (351), digest entries (521), issues (295), and activity documents are read and parsed again on every call. A thread read costs 80 to 110 ms, conversations 35 to 50 ms, the digest 60 to 85 ms, issues 22 to 30 ms, and YAML parsing is about a fifth of serve's CPU. `agent.state`, `inbox.designer`, `threads.list`, `messages.list`, `board.get`'s `merge.queue`, and `stats.get` all pay it; `threads.read` averages about 600 ms in the serve log.

## Acceptance criteria

- [x] flai serve and flai mcp keep each of the five parsed between requests, as S-0156 keeps items: each read stats every file and parses again only those whose size or modification time changed, or that are new, and drops those that are gone.
- [x] Tests show a changed, added, and removed file of each kind is seen on the next read, and that two readers at once get the same answer.
- [x] A warm `threads.list`, `inbox.designer`, and `stats.get` on this repository, timed with `flai hostapi --timing` in one process or through flai serve, each spend at least 80 ms less in their reads, and server-performance.md records the before and after.

## Tasks
- T-2660 workitem keeps parsed files of any kind by path, stat, and the racy window, for the five readers to share
- T-2667 threads.List and Get read threads through the kept store
- T-2671 messages.List and workitem's Conversations read conversations through the kept store
- T-2673 digest.List and Get read digest entries through the kept store
- T-2676 issues.List and Get read issues through the kept store
- T-2680 Repo.Activities and Activity read activity documents through the kept store
- T-2684 Measure warm threads.list, inbox.designer, and stats.get before and after, and record them in the design

## Notes

- Proposed fix: one kept-files store shared by the five packages, on the model of `flai/internal/workitem/store.go`.
- The stories for causes 12 and 14 gain from this one but do not wait for it.
- Criterion 3 is met by the operator's confirmation on TH-0913 (2026-11-09) of the orchestrator's recommendation. `stats.get` was judged on its reads' CPU, 389 ms down to 129 ms warm, and on the kept reads taken one at a time, 292 ms down to 18 ms. Its reads' wall time fell only 46 ms, 96 ms to 50 ms, because S-0607 runs them at once and the slowest is a read this story does not keep. `design/system/server-performance.md` § Cause 13 fixed records both.

### Planning

Touches, by where each came from:

| Touch | From |
|-------|------|
| `flai/internal/threads/threads.go`, `messages/messages.go`, `digest/digest.go`, `issues/issues.go` | declared |
| `flai/internal/threads/threads_test.go`, `messages/messages_test.go`, `issues/issues_test.go` | co-change: `flai touches suggest` finds them in 21 to 27% of the commits that changed the declared files |
| `flai/internal/digest/digest_test.go` | layout: the test file beside `digest.go` |
| `flai/internal/workitem/kept.go`, `kept_test.go` | design: the Notes' one shared store. It goes in `workitem` beside `store.go`, which all five packages import already, and reuses its `racyWindow` |
| `flai/internal/workitem/activity.go`, `activity_test.go` | layout: `Repo.Activities` and `ReadActivity` parse the activity documents |
| `flai/internal/workitem/share.go`, `share_test.go` | layout: `Repo.Conversations` parses every conversation again for `Repo.Holds`, which `agent.state` and `board.get` call |
| `design/system/server-performance.md` | criterion 3 |
| `design/system/flai-cli.md` | design: its Internal structure describes the S-0156 store and gains the kept files |

No folder touch. `flai touches suggest` also ranks `docs/users/flai-reference.md`, `docs/users/flai.md`, `docs/operators/settings.md`, and files under `flai/cmd`. They are left out because no command, flag, or setting changes. `flai/internal/workitem/store.go` is left out too, because S-0600 and S-0606 change it. `kept.go` reuses its `racyWindow` and changes nothing in it.

No `after`. S-0596, which wrote cause 13, is done. S-0603 and S-0604 read these files too, but neither writes anything this story reads.

Forecast: 1h30m, against flai's 13m. flai's figure is 103 s a unit of size times size 7, from 3 criteria and 4 touches. S-0156 kept one kind in 703 s. This story builds a shared store, moves five readers onto it, tests each kind, and measures on the host. The delivery moves with the duration, from flai's 00:10 to 01:30.

Cost of delay: no value yet. Neither S-0602 nor E-0033 has inputs, so `flai cod` refuses. TH-0888 asks the operator for `penalty_per_week: 50` on this story. The value is worked out once that input is set.

### Accepted by the orchestrator

- Verified: 8b29cf861eef8a98a2a42170c1b20a907b4b14f7
- At: 2026-11-09T07:17:27Z

Verdict: pass
- 1: flai/internal/workitem/keptfiles.go, flai/internal/threads/threads.go, flai/internal/messages/messages.go, flai/internal/digest/digest.go, flai/internal/issues/issues.go, flai/internal/workitem/activity.go, flai/internal/workitem/share.go
- 2: flai/internal/workitem/keptfiles_test.go, flai/internal/threads/threads_test.go, flai/internal/messages/messages_test.go, flai/internal/digest/digest_test.go, flai/internal/issues/issues_test.go, flai/internal/workitem/activity_test.go, flai/internal/workitem/share_test.go
- 3: design/system/server-performance.md, design/system/flai-cli.md
