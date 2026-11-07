---
id: S-0602
type: story
nature: improvement
title: flai serve keeps threads, conversations, digest entries, issues, and activity documents parsed between requests, so a thread read takes about 5 ms, not 80 to 110 ms
status: backlog
parent: E-0033
owner: alex
created: 2026-11-07T22:40:05Z
updated: 2026-11-07T23:49:37Z
transitions: []
tags: [cli]
touches: [flai/internal/threads/threads.go, flai/internal/messages/messages.go, flai/internal/digest/digest.go, flai/internal/issues/issues.go, flai/internal/workitem/kept.go, flai/internal/workitem/kept_test.go, flai/internal/threads/threads_test.go, flai/internal/messages/messages_test.go, flai/internal/workitem/share.go, flai/internal/workitem/share_test.go, flai/internal/digest/digest_test.go, flai/internal/issues/issues_test.go, flai/internal/workitem/activity.go, flai/internal/workitem/activity_test.go, design/system/server-performance.md, design/system/flai-cli.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
forecast:
  duration: 1h30m
  delivery: 2026-11-08T01:30:00Z
  basis: "flai's 13m (103 s a unit of size times 7) scaled up: S-0156 kept one kind in 703 s, and this story keeps five kinds on a shared store, with tests for each and a measurement on the host."
  by: planner-S-0602
  at: 2026-11-07T23:48:53Z
finalized:
  by: alex
  at: 2026-11-07T23:49:33Z
---
# S-0602 flai serve keeps threads, conversations, digest entries, issues, and activity documents parsed between requests, so a thread read takes about 5 ms, not 80 to 110 ms

## Goal

Cause 13 in [server-performance.md](../../../design/system/server-performance.md#where-the-time-goes): threads (884 files), conversations (351), digest entries (521), issues (295), and activity documents are read and parsed again on every call. A thread read costs 80 to 110 ms, conversations 35 to 50 ms, the digest 60 to 85 ms, issues 22 to 30 ms, and YAML parsing is about a fifth of serve's CPU. `agent.state`, `inbox.designer`, `threads.list`, `messages.list`, `board.get`'s `merge.queue`, and `stats.get` all pay it; `threads.read` averages about 600 ms in the serve log.

## Acceptance criteria

- [ ] flai serve and flai mcp keep each of the five parsed between requests, as S-0156 keeps items: each read stats every file and parses again only those whose size or modification time changed, or that are new, and drops those that are gone.
- [ ] Tests show a changed, added, and removed file of each kind is seen on the next read, and that two readers at once get the same answer.
- [ ] A warm `threads.list`, `inbox.designer`, and `stats.get` on this repository, timed with `flai hostapi --timing` in one process or through flai serve, each spend at least 80 ms less in their reads, and server-performance.md records the before and after.

## Tasks
- T-2660 workitem keeps parsed files of any kind by path, stat, and the racy window, for the five readers to share

## Notes

- Proposed fix: one kept-files store shared by the five packages, on the model of `flai/internal/workitem/store.go`.
- The stories for causes 12 and 14 gain from this one but do not wait for it.

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
