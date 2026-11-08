---
id: TH-0908
title: "S-0606 criterion 3: the walk misses 30 ms with only the archive on the slow tier"
anchor:
  path: wip/kanban/stories/S-0606-item-get-answers-from-a-kept-item-list-and-the-file-watcher-stops-statting.md
  item: S-0606
status: open
participants: [agent-S-0606, orchestrator]
created: 2026-11-08T01:40:45Z
updated: 2026-11-08T01:43:10Z
---

# TH-0908 S-0606 criterion 3: the walk misses 30 ms with only the archive on the slow tier

On wip/kanban/stories/S-0606-item-get-answers-from-a-kept-item-list-and-the-file-watcher-stops-statting.md.

## Entries

### 2026-11-08T01:40:45Z agent-S-0606
S-0606's watcher walk still misses criterion 3's 30 ms target. With `wip/archive` on the slow tier, a walk between slow looks costs about 45 ms of CPU, against about 75 ms on main. The measurement was taken at load 110 to 172 on 24 cores; the figures are in `design/system/server-performance.md` § Cause 16. The fast walk still stats 2,695 files outside the archive, most of them threads (908), digest entries (526), conversations (370), issues (297), and ADRs (264). With those five folders on the slow tier too, a walk measured 22 to 23 ms of CPU (p90 27 to 30 ms).

The cost of moving them is that an edit made in place in one of those folders reaches the dashboard up to 10 s late. A file added, removed, or replaced by a rename is still seen at the next poll. flai writes threads, conversations, and digest entries by rename, so the dashboard sees those at once. It writes issues (`os.WriteFile` in `internal/issues`) and ADRs (`internal/adr`) in place, and a hand edit is often in place too.

Which do you want?

1. Recommended: put `wip/threads`, `wip/messages`, `design/digest`, `design/issues`, and `design/adrs` on the slow tier, and make flai's own issue and ADR writes atomic (write and rename) so the dashboard sees them at once. Only an in-place hand edit in those folders waits up to 10 s. The walk then meets 30 ms.
2. Keep only the archive on the slow tier and change criterion 3 to record about 45 ms of CPU a walk, down from 75 ms.

Meanwhile I am making `item.get` copy only the items it answers, which is what keeps it above 50 ms. That needs no answer from you.

### 2026-11-08T01:43:10Z orchestrator (recommendation)
**Recommended:** option 1. Put `wip/threads`, `wip/messages`, `design/digest`, `design/issues`, and `design/adrs` on the slow tier. Make flai's own issue and ADR writes atomic (write and rename). Keep criterion 3 as written.

**Why option 1:**

- Criterion 3's goal is a walk under 30 ms. agent-S-0606's third program measured 22 to 23 ms of CPU a fast look with those five folders on the slow tier, p90 27 to 30 ms, on a host at load 110 to 172. Option 2 accepts about 45 ms, half again the target.
- The cost is small. Every flai write to those folders reaches the dashboard at the next poll: threads, conversations, and digest entries are written by rename already, and issues and ADRs will be once their writes are atomic. A `git` merge on main replaces files too, so an accepted story's ADRs and issues are seen at once. Only an in-place hand edit in those five folders waits up to 10 s.
- The extra work is two small write changes, in `flai/internal/issues` and `flai/internal/adr`. No story in progress or review touches either.

**What the story's agent should add, if you confirm:**

- Widen S-0606's touches with the two packages and their tests.
- Add a test that a flai issue write and an ADR write are each seen at the next fast look.
- Record the five folders and the reason under cause 16's Fix in `design/system/server-performance.md`, and the watcher's tiers in `design/system/flai-cli.md`.

**Why I recommend rather than answer:** both options change S-0606's scope. Option 1 goes beyond cause 16's proposed fix, which names only `wip/archive`. Option 2 rewrites a criterion. Either way, the call is yours.

**Alternative:** option 2. Keep only the archive on the slow tier, and reword criterion 3 to record about 45 ms of CPU a walk, down from about 75 ms.

Source rests on cause 16's row in the section below, and on the measurements under § Cause 16 on `story/S-0606`.

Source: design/system/server-performance.md § Where the time goes
