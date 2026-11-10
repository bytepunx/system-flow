---
title: Server-side response times
updated: 2026-11-10
status: active
topics: [server-side, back-end, cli, dashboard]
---

# Server-side response times

The finding of S-0152, for E-0012. The board and other pages pause for long moments while they load. This document says how the time a request takes is now measured beneath its transport, what was measured on this repository, where the time goes, and what would take it away. Each cause has a story under E-0012. [A second pass at 596 stories (S-0596)](#a-second-pass-at-596-stories-s-0596) measures it again and names causes 8 to 21.

> Since S-0195 ([ADR-0067](../adrs/0067-accepted-work-reaches-the-remote-only-when-it-is-published-and-agents-publish.md)) the dashboard has no unpushed notice: the method `push.pending` and the route `/api/unpushed` are gone, and a board load no longer asks for them. The measurements below keep them as they were measured.

## How it is measured

`flai serve` and `flai mcp` log one `request answered` event per request, timed inside flai from when the request was read to when its answer is ready to write, with the phases the time went to. Since S-0611, over the channel, the event is logged after the answer is written, and `send_ms` gives the write's time apart, the wait for the connection's write lock included ([flai-cli.md](flai-cli.md#internal-structure), [Request timing](../../docs/users/flai.md#request-timing-s-0152)). The dashboard logs its own duration for each HTTP request. The difference is the transport: the channel's WebSocket and the dashboard's own handling. `flai hostapi --timing` times one method with no transport at all, and `FLAI_PPROF_ADDR` offers Go's profiles from `flai serve`.

## What was measured

On 2026-09-29, on the operator's host (WSL2, 24 cores, load 3 to 4 with two agents working), against this repository: 1078 Markdown files under `design`, `docs`, and `wip`, and 895 under `wip`, most of them archived work items. Each method ran three times in-process with `flai hostapi --timing`; the median is shown. MCP tools ran three times over stdio against `flai mcp`. The dashboard's durations are from the `flaiover` container's request log over seven minutes of normal use.

### Dashboard methods

| Method | Asked by | Median | Answer | Where the time went |
|--------|----------|--------|--------|---------------------|
| `board.get` | board, `/api/projects` every 30 s and on change | 190 ms | 6 KB | `repo.list` 114, `release.pending` 75 (25 git processes: 13 `diff-tree`, 6 `log`, 6 `tag`) |
| `inbox.designer` | inbox badge, `/api/projects` | 163 ms | 1 KB | `check.run` 153 |
| `item.get` | story page | 118 ms | 8 KB | `repo.list` 119 |
| `items.list` (archive, bodies) | items page | 122 ms | 1.28 MB | `repo.list` 115, `encode` 6 |
| `docs.tree` | documents page | 135 ms | 772 KB | `docs.walk` 127, `encode` 8 |
| `publish.preview` | board, publish banner | 235 to 276 ms | 52 B | `exec.flai.release` all of it |
| `push.pending` | unpushed notice, review (both removed in S-0195) | 157 ms | 66 B | `exec.flai.push` all of it |
| `stats.get` | charts | 282 ms | 82 KB | `exec.flai.stats` all of it |
| `stream.diff` | story page | 261 ms | 82 KB | `exec.flai.stream.diff` all of it |
| `item.show` | story page, editor | 167 ms | 2 KB | `exec.flai.edit` all of it |
| `search.query` | search | 256 ms cold | 13 KB | `search.index` 297 cold; `flai serve` keeps the index until a file changes |
| `threads.list` | story page, threads | 7 ms | 97 KB | `threads.read` 5 |
| `activity.get` | activity | 5 ms | 1 KB | `repo.list` without the archive 5 |
| `agent.status` | board and activity every 15 s | 6 ms | 2 KB | `agent.state` 6 |
| `project.info`, `settings.get`, `doc.get`, `adrs.list` | every page | under 6 ms | | |

### MCP tools

| Tool | Median | Where the time went |
|------|--------|---------------------|
| `inbox` | 308 ms | `changes.read` 107 and `repo.list` 106: every item, archive included, listed twice; `release.pending` 75 |
| `board` | 194 ms | `repo.list` 111, `release.pending` 76, `pending.detect` 5 |
| `wait_for_work` | 190 ms per look | as `board`, again at every change to a work item, thread, or narrative while it waits |
| `prime` | 148 ms | the context pack; 98 KB |
| `item_get` | 111 ms | every item listed to find the children |
| `doc_search` | 68 ms | |
| `who_touches`, `doc_get` | under 4 ms | |

### The dashboard's view

| Route | Requests in 7 min | Median | p90 | flai's own time for it |
|-------|-------------------|--------|-----|-----------------------|
| `/api/stats` | 10 | 298 ms | 304 ms | `stats.get` 282 ms |
| `/api/projects` | 42 | 244 ms | 268 ms | `board.get`, `inbox.designer`, `agent.status` together |
| `/api/publish` | 56 | 234 ms | 249 ms | `publish.preview` 235 ms |
| `/api/board` | 18 | 192 ms | 209 ms | `board.get` 190 ms |
| `/api/unpushed` (removed in S-0195) | 58 | 161 ms | 172 ms | `push.pending` 157 ms |
| `/api/inbox` | 13 | 159 ms | 179 ms | `inbox.designer` 163 ms |
| `/api/items/[id]` | 1 | 114 ms | | `item.get` 118 ms |
| `/api/agent`, `/api/host-agent` | 287 | 2 to 10 ms | | |

## What it says

The transport is not where the time goes. Each dashboard route takes within a few milliseconds of what flai takes to answer it in-process, so the WebSocket channel and the dashboard's handling add almost nothing. The time is flai's work, and seven causes account for it.

| # | Cause | Cost per request | Who pays |
|---|-------|------------------|----------|
| 1 | Every request re-reads and re-parses every work item, the archive included (over 800 files; 18 are open), from disk. flai serve and flai mcp keep nothing between requests. | 105 to 120 ms | `board.get`, `item.get`, `items.list`; MCP `inbox` twice, `board`, `wait_for_work`, `item_get` |
| 2 | Pending releases are worked out by starting 25 git processes, one `diff-tree`, `log`, or `tag` at a time, on every board. | 75 ms | `board.get`; MCP `inbox`, `board`, `wait_for_work` |
| 3 | The inbox runs the whole of `flai check` to find overlapping touches, one rule of many. | 153 ms | `inbox.designer`, so every `/api/projects` and inbox badge |
| 4 | Read methods that the dashboard polls answer by starting a flai process, which reads the repository again from nothing. | 157 to 282 ms | `publish.preview`, `push.pending`, `stats.get`, `stream.diff`, `item.show`, and the other `read` methods of the write table |
| 5 | A flai process takes 145 ms to start on this host before it does anything: `atotto/clipboard`, which `huh` brings in through `bubbles/textarea`, looks for clipboard programs on `PATH` when the package loads, and this `PATH` has 54 entries, 17 of them Windows folders under `/mnt`. With a short `PATH` it starts in 8 ms. | 145 ms per process | every cause 4 request, every write, and every `flai` an agent runs |
| 6 | The dashboard forgets every answer it holds when any file of the project changes, so a narrative log line or a thread reply makes the next look at every page ask flai for everything again; with agents working, a file changes every few seconds. The board also asks for its board, `/api/publish`, and the agents again at each `change` event, one per file and whatever the file, with nothing gathering changes that arrive together. `/api/publish` and `/api/unpushed` were asked about every 7.5 s. | the sum of what a page asks | every page while agents work |
| 7 | Two answers are large: `items.list` with the archive and bodies is 1.28 MB, and `docs.tree` with every file's front matter is 772 KB, walked and encoded on each ask. | 122 to 135 ms and the bytes | items page, documents page |

Cause 2 is gone since S-0157: `release.Pending` reads git through a history kept per repository root ([flai-cli.md](flai-cli.md#internal-structure)), and while HEAD and the tags are unchanged it starts one git process, `show-ref`. On 2026-09-29, on the same host and repository, `release.pending` in a kept `board.get`, called through the host API's method table as `flai serve` answers it, took 8.2 to 9.5 ms with that one process, and 7.9 to 10.2 ms in a kept MCP `board`, against 242 to 259 ms and 69 processes (45 `diff-tree`, 12 `log`, 12 `tag`) in a `board.get` before the change the same hour, since more items are in range than when this was first measured. After an acceptance it took 10.8 ms with two processes, and after a publish 8.3 ms with one, measured in a clone. The first look in a process, and every one-shot `flai board`, reads the whole history with four processes in about 60 ms. Most of what remains beyond `show-ref` is parsing work items in `repo.Get`, which cause 1's story addresses.

The pauses the operator sees come from these adding up. A board load asked `board.get` (190 ms), `publish.preview` (235 ms), and `push.pending` (157 ms; since S-0195 the board no longer asks for it), and `/api/projects` asks `board.get`, `inbox.designer`, and `agent.status` again (about 360 ms). Each of them reads every work item from disk, and the three that start flai also pay the 145 ms start. They run at once, on a host where agents keep asking the same MCP tools, so they compete for the disk and the CPU.

Cause 3 is gone since S-0158: `inbox.designer` runs the `wip.overlap` rule alone, over the items it has already listed, and `check.run` is no longer one of its phases. On 2026-09-29, on the same host and repository, it took 11 to 14 ms, median 12.7 ms, against 170 to 195 ms, median 178 ms, before the change the same hour; its phases were `threads.read` 5.6, `repo.list` 4.7, and `check.overlap` under 0.1 ms.

What is not a cause: the file watchers. `flai serve` walks the three folders every 300 ms and each waiting MCP call every 250 ms; one walk of 1078 files takes 5.6 ms.

## After the stories

Each story remeasures what its cause cost, on the same host and repository, and records it here.

### Cause 1: work items kept between requests (S-0156)

`flai serve` and `flai mcp` keep the items they have parsed and read again only the files that changed ([flai-cli.md](flai-cli.md#internal-structure)). Measured on 2026-09-29 with 1078 Markdown files, load about 1: the dashboard's methods called five times in one process through the table `flai serve` answers with, timed by the same recorder, and the MCP tools five times each over stdio against `flai mcp` built from the story's branch. The first call of a process is cold; the others are warm.

| Request | `repo.list` cold | `repo.list` warm | Warm total |
|---------|------------------|------------------|------------|
| `board.get` | 122 ms | 4 to 5 ms | 232 ms, of which `release.pending` 227 (cause 2) |
| `item.get` | | 3.4 to 4.5 ms | 4 ms |
| `items.list` (archive, bodies) | | 3.7 to 4.1 ms | 8 ms, of which `encode` 4 |
| MCP `inbox` | 117 ms | 4 to 5 ms, once | 241 ms, of which `release.pending` 226 |
| MCP `board` | | 4.1 to 4.3 ms | 236 ms, of which `release.pending` 226 |
| MCP `wait_for_work` | | 4 to 5 ms per look | `release.pending` 223 to 248 per look |
| MCP `item_get` | | 3.4 to 4.8 ms | 4 ms |

`inbox` lists the items once for the board and the changes, where it listed them twice. What is left of a board load is cause 2: this measurement found 69 git processes where S-0152 found 25, as more items were accepted since the last release.

### Cause 4: the dashboard's reads answered in flai serve's process (S-0159)

`publish.preview`, `push.pending` (removed in S-0195), `stats.get`, `stream.diff`, `item.show`, `item.move.preview`, `accept.preview`, and since S-0217 `order.by`, `promote.candidates`, and `release.evaluate` start no flai: `flai serve` answers them from the code the commands print with, over the items and git history it keeps (causes 1 and 2), with the answer the command gave ([flai-cli.md](flai-cli.md#commands)). Their `request answered` events have no `exec.flai` phase; each names its work, such as `release.pending` or `stats.compute`. S-0217's three came after the measurement below and are not in its table. Measured on 2026-09-29, on the same host and repository, load about 2: each method called six times in one process through the table `flai serve` answers with, the first call cold, beside the command run three times by a flai built from `main` the same hour, which already starts in about 6 ms (cause 5) and keeps nothing between runs.

| Method | A flai process (S-0152) | A flai process, same hour | In-process, cold | In-process, warm | Warm, where the time went |
|--------|-------------------------|---------------------------|------------------|------------------|---------------------------|
| `publish.preview` | 235 to 276 ms | 81 ms | 75 ms | 7.2 to 8.3 ms | `release.pending` 6.9, one `git show-ref` |
| `push.pending` | 157 ms | 24 to 26 ms | 20 ms | 17 to 19 ms | seven git processes: `tag` 6.8, `rev-list` 5.1, `rev-parse` 3.5, `log` 2.6 |
| `stats.get` | 282 ms | 139 to 142 ms | 125 ms | 9.1 to 10.6 ms | `stats.compute` 8.2 |
| `stream.diff`, 24 files | 261 ms | 83 to 89 ms | 74 ms | 73 to 75 ms | one `git diff` per file, 64 ms |
| `item.show` | 167 ms | 9 to 12 ms | 0.6 ms | 0.3 to 1.3 ms | |
| `item.move.preview`, an epic | | 16 ms | 5.8 ms | 4.9 to 5.8 ms | three `git rev-parse` |
| `accept.preview`, refused | | 9 ms | 2.1 ms | 1.9 to 2.6 ms | one `git rev-parse` |

`publish.preview` and `push.pending`, which the board asked for at every change, answered in under 20 ms once warm; since S-0195 the board asks for `publish.preview` alone. What is left is git's: `push.pending` asked git where the branch stands against its remote at every call, and `stream.diff` reads each file's patch with a process of its own, which grows with the story.

### Cause 5: a flai process starts without searching PATH (S-0160)

flai asks its questions through a prompt package of its own and no longer depends on `huh`, so nothing searches `PATH` when flai starts ([ADR-0052](../adrs/0052-flai-asks-its-questions-a-line-at-a-time-with-its-own-prompt-package-not-with.md), [go-libraries.md](../tech/go-libraries.md#prompts-s-0160)). On 2026-09-29, on the same host with its own `PATH` of 54 entries, 17 under `/mnt`, `flai version` took 5.5 to 8.3 ms over 30 runs, median 5.8 ms, against 160 to 178 ms, median 163 ms, for `main` built the same hour. `GODEBUG=inittrace=1 flai version` shows no package init over 0.6 ms, where `atotto/clipboard` took 171 ms.

### Cause 6: the dashboard forgets and asks only for what a change affects (S-0161)

A change reaches the dashboard's server with the path's kind, read against the manifest's layout: a work item (`board.md` and the archive included), a narrative, a thread, an ADR, another document, the manifest, or anything else. The server forgets only the answers flai reads from that kind ([flaiover-dashboard.md](flaiover-dashboard.md#where-the-data-comes-from-s-0073)), so a narrative log line keeps `board.get`, the items, the statistics, and the ADRs, and forgets the inbox, the activity, and the documentation tree. Pages ask again only for the kinds they show, once for the changes that arrive within 500 ms of each other, and at least every 2 s during a burst that does not settle. The board asks for `board.get` and `/api/publish` (and `/api/unpushed` until S-0195 removed it, [ADR-0067](../adrs/0067-accepted-work-reaches-the-remote-only-when-it-is-published-and-agents-publish.md)) again only when a work item or the manifest changes, and for the agents when a work item or a thread changes or flai serve says an agent started or ended. Behaviour tests pin which paths forget which answers (`repo-forget.test.ts`) and what the board asks for each kind (`routes/board/board.svelte.test.ts`). The request rate was not measured again in the running dashboard: that needs an image built from this story in place of the operator's container.

### Cause 7: the items and documents pages ask for what they show (S-0162)

The dashboard asks `items.list` without bodies, which it never showed from the list: a body is asked for with `item.get` when its item is opened. It asks for the archive only where a page shows it, and the overview, which shows only how many items are archived, asks `items.count` instead. `docs.tree` answers each file's title and not its front matter, which the explorer reads with `doc.get` when a file is opened, and flai keeps each title by file while the file has the same identity, modification time, and size, as it keeps items (cause 1), so a warm answer lists the folders and reads no file ([flai-cli.md](flai-cli.md#commands), [flaiover-dashboard.md](flaiover-dashboard.md#api-s-0011)). `/_ready` asks for the manifest alone. Measured on 2026-09-29, on the same host and repository, now 1128 Markdown files under the three folders, load 0.3 to 3.5: each method called six times in one process through the table `flai serve` answers with, the first call cold, the size as encoded.

| Request | Before (S-0152) | Cold | Warm | Answer |
|---------|-----------------|------|------|--------|
| `items.list`, archive and bodies, as the overview asked | 122 ms, 1.28 MB | | 7.1 to 11.6 ms | 1.36 MB |
| `items.list`, active, no bodies, as the overview asks now | | 2.0 ms | 0.3 to 0.4 ms | 8.4 KB |
| `items.count` | | 123 ms | 3.9 to 5.3 ms | 28 B |
| `items.list`, epics with the archive, as the charts ask | | | 3.7 to 4.8 ms | 8.4 KB |
| `docs.tree` | 135 ms, 772 KB | 129 to 131 ms | 6.1 to 8.2 ms, of which the walk 5.5 to 6.7 | 320 KB |
| `project.info`, all `/_ready` asks now | | | 0.2 ms | 782 B |

The overview's first answers from flai are 8.4 KB and 28 B, where they were 1.36 MB. `docs.tree` stays over 200 KB: with 1128 files, the names and paths alone are 229 KB. What is left of a warm answer is the walk that tells flai which files changed, one `lstat` per file. Neither was measured in the running dashboard, for the reason given under cause 6.

## A regression in the chart data: the hold replay (S-0471)

On 2026-09-29, `stats.get` answered in 125 ms cold and 9.1 to 10.6 ms warm (cause 4). On 2026-10-16, S-0427 measured the same call, with the parameters the charts page sends (`{"since":"30d","type":"story","bucket":"day"}`), at 4,715 ms cold and 4,022 ms warm, and `flai stats --json` at 4,701 ms ([S-0427's results](../experiments/S-0427-sqlite-and-postgres-read-models.md#what-was-done)). A CPU profile of the warm calls put 67 percent of the samples under `metrics.replayHolds` ([I-0155](../issues/I-0155-the-chart-data-takes-about-4-s-warm-against-10-ms-measured-on-2026-09-29-almost-all-of-it-in-metrics-replayholds.md)). The charts page asks for the statistics again on every item change, and every agent's `flai stats` waits on the same work.

### Cause

`metrics.replayHolds` came with S-0205 on 2026-10-03. It gives each story the seconds it was held, and `claims` its held days and weeks. At every distinct event time in the history, it built a `workitem.Holds` over every item created so far and judged every ready story again. Its cost grew with the square of the history.

### Fix

`replayHolds` judges the holds again only at an event time that can change one. That is a story's creation, or its move into or out of ready, in progress, or done. It is also the creation of a task of a story in ready or in progress, or that task's move into or out of done or cancelled. It judges only the stories in ready or in progress, the stories their `after` names, and their tasks. A `workitem.HoldsMemo` remembers claims, shared-path coverage, and how one claim holds another, between judgements. `PathsOverlap` no longer allocates. Tests compare the replay with a reference that judges at every event time, on the hold rules and on a generated history of this repository's shape, and find the same held seconds (`flai/internal/metrics/claims_test.go`, `flai/internal/metrics/claims_bench_test.go`).

On a generated history at this repository's scale, `BenchmarkReplayHolds` in `flai/internal/metrics/claims_bench_test.go` gave these figures for one replay. `BenchmarkComputeClaims`, the whole report over the same history, takes 134 to 147 ms after the fix.

| One replay | Time | Memory | Allocations |
|------------|------|--------|-------------|
| Before | 2,726 ms | 1.07 GB | 14.3 million |
| After | 71 ms | 26.8 MB | 275 thousand |

### Measured on this repository

On 2026-10-21, between 12:41 and 12:47 UTC, on the operator's host (WSL2, 24 cores, Go 1.26.0), against this repository: 4160 Markdown files under `design`, `docs`, and `wip`, 3713 of them under `wip`, and 483 stories in the report. Two flai binaries were built the same hour: `main` at b465e0a28 and `story/S-0471` at db7a2c03a. Both ran from the main checkout, with `FLAI_CONFIG` and `FLAI_CACHE_DIR` as `scripts/env.sh` sets them.

- Cold is a fresh process: `flai hostapi --timing stats.get '{"since":"30d","type":"story","bucket":"day"}'`, three runs each, the binaries alternating. The figure is the median wall time; the phases are from the `request answered` event.
- Warm is a throw-away program that called the host API's method table six times in one process, with the same recorder as `--timing`, as S-0427 did. It was built with `go build -overlay`, which wrote nothing into the module. The figure is the median of calls 2 to 6; CPU is from `getrusage`.
- `flai stats --json` is a fresh process, three runs each, wall time.

The host was busy: the one-minute load average was 11 to 21, with other flai processes using one to one and a half cores each.

| Request | `main` | `main`, where the time went | S-0471 | S-0471, where the time went |
|---------|--------|-----------------------------|--------|-----------------------------|
| `stats.get` cold | 7,552 ms (6,888 to 8,357) | `stats.compute` 7,485, of which `exec.git.log` 448 | 1,569 ms (1,505 to 1,764) | `stats.compute` 1,510, of which `exec.git.log` 432; `encode` 9 |
| `stats.get` warm, calls 2 to 6 | 5,695 ms (5,454 to 5,896), 7.8 s of CPU a call | `stats.compute` 5,672, of which `exec.git.log` 452 | 894 ms (884 to 945), 0.7 s of CPU a call | `stats.compute` 874, of which `exec.git.log` 453 |
| `flai stats --json` | 7,730 ms (6,521 to 8,112) | | 1,569 ms (1,528 to 1,799) | |

A first round, a few minutes earlier at a load of 10 to 30, gave the same picture: `stats.get` cold 12,033 ms against 2,680 ms, warm 10,514 ms against 1,079 ms, and `flai stats --json` 10,513 ms against 2,717 ms. In CPU profiles of that round's warm calls, `metrics.replayHolds` held 73 percent of the samples for `main` and 22 percent for S-0471.

The two binaries give the same report. Their `flai stats --json`, run seconds apart, differ only in fields measured from now: `generated_at`, `window_start`, ages, times in state, cost of delay incurred and the times of its order, waits clipped to the window's start, and the blocked seconds of a story still blocked. `claims`, with its days and weeks, `coordination`, and every story's held seconds are identical in each of the four pairs compared. The six pairs of `stats.get` answers agree the same way. One of them also differs in turns and spend, because an agent's session went on between its two runs: the `flai stats --json` run between them already has the later figures.

A warm answer is still about 90 times the 9.1 to 10.6 ms of 2026-09-29. Half of it is `exec.git.log`: at every call, `statsread.Read` reads the stories' commits for `claims.drift` with one `git log` (`storygit.Committed`), which S-0427 measured at about 370 ms on 2026-10-16. S-0500 takes that log out of a warm call ([below](#the-stories-commits-kept-between-calls-s-0500)). By the S-0471 profile, a warm call spends about 330 ms of CPU in `metrics.Compute`, about 160 ms of it in `replayHolds`, and about 170 ms in `statsread.Read`, two fifths of that reading the threads. These figures were taken under the load above.

## The stories' commits kept between calls (S-0500)

After S-0471, half of a warm `stats.get` was one `git log` ([I-0184](../issues/I-0184-every-stats-get-and-flai-stats-runs-one-git-log-for-claims-drift-about-450-ms-of-a-900-ms-warm-answer.md)). S-0471 measured it at 448 to 453 ms of an 894 ms warm answer.

### Cause

`statsread.Read` calls `storygit.Committed` for `claims.drift` at every `stats.get` and every `flai stats`. `Committed` ran one `git log --name-only` over the main branch and every `story/` branch at every call, and parsed its output. In `flai serve`'s long-running process, the charts page asks for the statistics again on every item change, so the same log ran again and again while no branch had moved.

### Fix

`Committed` keeps its answer per repository root for the life of the process (`flai/internal/storygit/committed.go`). The key is the main branch's name, the wip folder, and the output of `git show-ref --heads`, which names every local branch's tip. A call runs `git rev-parse --abbrev-ref HEAD` and `git show-ref --heads`. While the key is unchanged, it answers what it kept and runs no log. When any local branch's tip moves, the next call reads the whole log again: a commit on the main branch or a story branch, a story branch made or deleted, and the rebase `flai stream sync` makes. Each call gets its own deep copy of the answer, and one lock per process keeps two calls from reading the log at once.

With S-0500 alone, a fresh process had nothing kept. A one-shot `flai stats`, `flai hostapi --timing`, and the first `stats.get` of a `flai serve` process still ran the log once, and so did the first call after a tip moved. S-0531 keeps the answer between processes ([below](#the-stories-commits-kept-between-processes-s-0531)).

Tests count the `git log` processes a real runner starts: a second call with nothing moved starts none, a commit on main or a story branch or a new story branch makes the next call read again, a caller who changes its answer does not change the next one, and each root keeps its own (`flai/internal/storygit/committed_test.go`, and `TestReadStartsNoLogWhileNothingChanged` in `flai/internal/statsread/statsread_test.go`).

### Measured on this repository

On 2026-10-26, between 04:37 and 04:39 UTC, on the operator's host (WSL2, 24 cores, Go 1.26.0), against this repository: 4878 Markdown files under `design`, `docs`, and `wip`, 4002 of them under `wip`, and 504 stories in the report. Two flai binaries were built the same hour: `main` at b458a6f90 and `story/S-0500` at db1a2f400. The story branch was `main` and one commit, so the `main` binary was built from a copy of the story's `flai` folder with `committed.go` as `main` has it; the two trees differed in nothing else but tests. Both ran from the main checkout, with `FLAI_CONFIG` and `FLAI_CACHE_DIR` as `scripts/env.sh` sets them. The method was S-0471's:

- Cold is a fresh process: `flai hostapi --timing stats.get '{"since":"30d","type":"story","bucket":"day"}'`, three runs each, the binaries alternating. The figure is the median wall time; the phases are the medians from the `request answered` event.
- Warm is a throw-away program that called the host API's method table six times in one process, with the same recorder as `--timing`, built with `go build -overlay`, which wrote nothing into the module. It ran twice for each binary, alternating. The figure is the median of calls 2 to 6 of both runs, ten calls; CPU is the process's own from `getrusage`, which leaves out git's.
- `flai stats --json` is a fresh process, three runs each, wall time.

The host was very busy: the one-minute load average was 21 to 23 during the cold runs, 24 to 27 during the warm runs, and 31 to 39 during the `flai stats --json` runs.

| Request | `main` | `main`, where the time went | S-0500 | S-0500, where the time went |
|---------|--------|-----------------------------|--------|-----------------------------|
| `stats.get` cold | 2,374 ms (2,170 to 2,416) | `stats.compute` 2,298, of which `exec.git.log` 651; `encode` 10 | 2,494 ms (2,329 to 2,555) | `stats.compute` 2,403, of which `exec.git.log` 655 and `exec.git.show-ref` 3; `encode` 12 |
| `stats.get` warm, calls 2 to 6 | 1,317 ms (1,205 to 1,385), 0.95 s of CPU a call | `stats.compute` 1,286, of which `exec.git.log` 648 | 648 ms (592 to 725), 0.76 s of CPU a call | `stats.compute` 603, of which `exec.git.show-ref` 3 and `exec.git.rev-parse` 3; no `exec.git.log` in any of the ten calls |
| `flai stats --json` | 3,558 ms (2,827 to 3,858) | | 3,458 ms (3,409 to 4,330) | |

A warm answer takes half the time it did. The log, about 650 ms, and the parsing of its output are gone from every call after the first while no branch moves; what is left of git is two processes of about 3 ms each. The first call of each warm run read the log, as a fresh process does: 655 and 659 ms of git log in a call of about 2,360 ms.

A cold `stats.get` and a `flai stats --json` are the same within the noise, as they should be: each is a fresh process, which runs the log once. In a cold `stats.get` with S-0500, `exec.git.log` was 655 of 2,444 ms, about a quarter. `git show-ref --heads` adds about 3 ms to it.

The two binaries give the same report. Their `flai stats --json`, run seconds apart, have identical `claims`: every one of the 467 entries of `claims.drift`, and its days and weeks. They differ only in fields measured from now: `generated_at`, `window_start`, ages, times in backlog and in progress, cost of delay incurred and the times and horizon of its order, the seconds waited on threads, and the blocked seconds of a story still blocked. One pair was compared.

`main`'s warm answer, 1,317 ms, is longer than the 894 ms S-0471 measured for its binary on 2026-10-21, under a heavier load and a larger history (4878 Markdown files against 4160, 504 stories against 483). This round took no profile. What is left of a warm answer is `stats.compute`'s own work, about 600 ms, which S-0471's profile put in `metrics.Compute` and `statsread.Read`.

## The stories' commits kept between processes (S-0531)

After S-0500, a fresh process still ran the log once ([I-0209](../issues/I-0209-a-one-shot-flai-stats-and-the-first-stats-get-of-a-flai-serve-process-still-run-one-git-log-for-claims-drift-about-a-quarter-of-a-cold-answer.md)). S-0500 measured it at 655 of 2,444 ms in a cold `stats.get`, about a quarter.

### Cause

`storygit.Committed` kept its answer in memory only. Every one-shot `flai stats`, every `flai hostapi --timing`, and the first `stats.get` of each `flai serve` process started with nothing kept, so each ran one `git log` over the main branch and every `story/` branch, even when no branch had moved since the last process read it.

### Fix

`Committed` also keeps its answer in `.flai-cache/stories-commits.json` in the main checkout (`flai/internal/storygit/committed.go`). The file holds the repository root, the key, and the answer. The key is S-0500's: the main branch's name, the wip folder, and the output of `git show-ref --heads`. A process with nothing in memory reads the file. When the root and the key match, it answers from the file and runs no log. A file that is missing, cannot be parsed, or names another root or another key is a miss: the process reads the log, answers what it read, and writes the file again. It writes a file beside it and renames it over the old one, so a reader finds the old answer or the new one whole. A write that fails never fails the call; it costs the next process a read of the log. It never makes the cache folder: a project without a `.flai-cache` keeps nothing and reads the log in each process, since `flai check` reports a `.flai-cache` that no `.gitignore` ignores (`layout.gitignore`).

What still reads the log is the first call after any local branch's tip moved, in whichever process makes it. That call writes the file again, so the processes after it read none.

Tests run the log in one process and count the `git log` processes the next one starts: none while nothing moved, a new read after a branch moved, a read over a damaged file or one for another root or key, and an answer when the cache folder cannot be written (`flai/internal/storygit/committed_test.go`). `TestReadStartsNoLogInANewProcessWhileNothingChanged` in `flai/internal/statsread/statsread_test.go` starts the test binary again as a new process and finds that its `statsread.Read` starts no log.

### Measured on this repository

On 2026-10-28, between 18:33 and 18:34 UTC, on the operator's host (WSL2, 24 cores, Go 1.26.0), against this repository: 5161 Markdown files under `design`, `docs`, and `wip`, 4233 of them under `wip`, and 531 stories in the report. Two flai binaries were built at 18:32 UTC: `main` at c26ac99b0 and `story/S-0531` at 9e9eb49af. The story branch was `main` and one commit, which changed `committed.go` and two test files; the `main` binary was built from a copy of the story's `flai` folder with those three files as `main` has them. Both ran from the main checkout, with `FLAI_CONFIG` and `FLAI_CACHE_DIR` as `scripts/env.sh` sets them. The method was S-0500's, for the fresh processes only:

- Cold is a fresh process: `flai hostapi --timing stats.get '{"since":"30d","type":"story","bucket":"day"}'`, three runs each, the binaries alternating. The figure is the median wall time; the phases are the medians from the `request answered` event.
- `flai stats --json` is a fresh process, three runs each, alternating, wall time.
- Before the rounds, the kept file did not exist. One S-0531 `stats.get` ran first and wrote it. It is shown on its own: a miss, the path the first call after a tip moved takes.
- The local branch tips were hashed before and after each run. No tip moved during the 13 runs, and the kept file was written once, by the first run.

The host was busy: the one-minute load average was 18 to 22.

| Request | `main` | `main`, where the time went | S-0531 | S-0531, where the time went |
|---------|--------|-----------------------------|--------|-----------------------------|
| `stats.get` cold, no kept file | | | 2,278 ms, one run | `stats.compute` 2,206, of which `exec.git.log` 553, `exec.git.rev-parse` 2.8 and `exec.git.show-ref` 2.6; `encode` 10.2 |
| `stats.get` cold | 2,154 ms (2,148 to 2,514) | `stats.compute` 2,079, of which `exec.git.log` 552, `exec.git.rev-parse` 3.0 and `exec.git.show-ref` 2.9; `encode` 11.5 | 1,751 ms (1,666 to 1,823) | `stats.compute` 1,672, of which `exec.git.rev-parse` 3.5 and `exec.git.show-ref` 3.1; `encode` 12.6; no `exec.git.log` in any of the three runs |
| `flai stats --json` | 2,248 ms (2,242 to 2,277) | | 1,591 ms (1,553 to 1,628) | |

A fresh process no longer reads the log while no branch moves. A cold `stats.get` takes 403 ms less, about a fifth, and a `flai stats --json` 657 ms less, about three tenths. `main`'s log took 552 ms of each cold call. What is left of git is the two processes S-0500 left in a warm call, about 3 ms each. The first call with no kept file costs what `main`'s does: it reads the log, 553 ms, and writes the file, 490 KB.

`main`'s median `stats.compute` less its median log was 1,527 ms; S-0531's median `stats.compute` was 1,672 ms. The difference includes reading and parsing the kept file, and is within the spread of S-0531's three runs, 1,590 to 1,753 ms. This round took no profile to split it.

The two binaries give the same report. In each of the six pairs compared, three of `stats.get` and three of `flai stats --json`, run about two seconds apart, `claims.drift`, all 491 entries, and `claims.days` are identical. `claims.weeks` differs only in the current week's held seconds, `after` and `overlap`, by 1 to 14 seconds, which grow with now while a story is held: `main` ran first in each pair, and S-0531's figures are the larger.

## Parses and lint findings kept between processes (S-0600)

This section covers cause 10 of [the second pass](#a-second-pass-at-596-stories-s-0596).

### Cause

The item cache S-0156 added lived only in a process's memory, so every one-shot `flai`, every `flai move`, and every `flai check` parsed all the items one file after another (`flai/internal/workitem/store.go` `list`). Every `flai check` linted every markdown file under `wip` again, one after another (`flai/internal/check/check.go` `markdown`).

### Fix

- `store.list` stats each file of a folder under the folder's lock, takes the kept item of each file whose key still matches, and parses the rest on min(GOMAXPROCS, 16) workers, keeping them in file name order.
- `mdlint.LintFolder` walks the folder and lints its files on the same bound of workers, in walk order (`flai/internal/mdlint/folder.go`).
- A folder's first list in a process reads its kept items file, `.flai-cache/items-<folder>.json`, such as `items-wip-archive-kanban-tasks.json` (`flai/internal/workitem/kept.go`). A file whose size and modification time match its kept entry takes the kept item without a parse. A restored item equals `ReadItem`'s, its unknown fields and the form it was read as included.
- `flai check` keeps each file's findings in `.flai-cache/lint-findings-<hash of the folder>.json` (`flai/internal/mdlint/kept.go`). An entry is keyed by the file's path, size, and modification time, and the file as a whole by the configuration's hash: its bytes, its rules, a rule-set version, and the format version. A file whose key matches takes its kept findings unread.
- Only a read made at least two seconds after the file last changed is kept, as S-0156 trusts one in memory.
- Both files are written again only when what they hold changed. Each is written beside the old one and renamed over it, and only into a `.flai-cache` that exists, as S-0531's kept commits are.
- A file that is missing, damaged, of another version, or another folder's or configuration's is a miss, never an error or a finding. It costs that process the parses or the lint, and is written again.

### Measured on this repository

On 2026-11-08, between 01:04 and 01:14 UTC, on the operator's host (WSL2, 24 cores, Go 1.26.0), in a scratch copy of the main checkout under `/tmp`. The copy held its working tree with the uncommitted `wip`: 3,343 work items and 5,209 markdown files under `wip`. It had an empty `.flai-cache`, which its `.gitignore` lists. Two binaries were built just before: `main` at 9fa3fec41, and `story/S-0600` at bf246d07d. Each ran with a scratch `FLAI_CONFIG`.

- Each run is a fresh process, timed by the shell's `time`. The three cases ran in turn in each round:
  - `main`;
  - S-0600 with the kept items and lint files deleted first, the first process after a cache is lost;
  - S-0600 with the kept files a previous process wrote, the usual case.
- `flai move` moved backlog task T-1269 to `ready` and back, both moves timed, five rounds: ten moves a case. `flai check --strict` ran five rounds; `flai show T-1269` ran five.
- The figure is the median, of wall time and of CPU time (user and system).
- The list alone was timed in process by a small program built against each tree with `go build -overlay`, calling `Repo.List(true)` from a fresh process, ten rounds.

The host was very busy: the one-minute load average was 84 to 115 during the moves, 51 to 83 during the checks, 101 to 118 during the shows, and 34 to 66 during the lists. Wall times swing with it, so the CPU times are given beside them.

| Fresh process | `main` | S-0600, no kept files | S-0600, kept files |
|---------------|--------|-----------------------|--------------------|
| `flai move` | 4.54 s wall, 2.86 s CPU | 2.87 s wall, 3.26 s CPU | 0.99 s wall, 0.89 s CPU |
| `flai check --strict` | 6.81 s wall, 7.14 s CPU | 4.51 s wall, 8.26 s CPU | 3.30 s wall, 3.81 s CPU |
| `flai show` | 5.26 s wall, 2.88 s CPU | 4.04 s wall, 3.48 s CPU | 1.31 s wall, 0.91 s CPU |
| `Repo.List(true)` in process | 4.54 s | 2.70 s | 0.94 s (0.31 to 1.37) |

Criterion 3 of S-0600 holds. With the kept files, a fresh `flai move` takes 3.55 s less wall time and 1.97 s less CPU than `main`. A fresh `flai check` takes 3.51 s less wall time and 3.33 s less CPU. With no kept files, the first process after a cache is lost takes 1.67 s and 2.30 s less wall time, but 0.40 s and 1.12 s more CPU. That CPU goes to the workers and to writing the kept files whole.

The two binaries give the same `flai check --strict` output, byte for byte, and both exit 1 on the same findings.

What is left:

- A fresh list with the kept files took 0.31 s at the lowest load, a load average of 35, and 0.94 s at the median. That is not yet the 0.3 s the story's title names. Most of it is decoding the kept items files: 10.5 MB of JSON, 6.0 MB of it the archived tasks and 3.8 MB the archived stories. The kept lint findings file is 0.9 MB.
- About 2.8 s of a check's CPU is its other rules, cause 11's (S-0601).
- `flai move` still lists twice (`flai/internal/workitem/transition.go`). In one process the second list stats the files and parses only what the move wrote.

## A second pass at 596 stories (S-0596)

The operator found flai slow again with 596 stories: the board, creating and editing items, the charts (especially over windows longer than 7 days), the close-out, and acceptance. S-0596 measured where the time goes now. Causes 1 to 7 above stay fixed. The time now goes to work that grows with the number of files, and to steps that run one after another.

### How it was measured

- The operator's serve log, `~/.flai/serve/serve.log`, from 2026-11-04 to 2026-11-07: its `request answered` events. It logs at `INFO` only the requests at or over `FLAI_SLOW_REQUEST`, so its percentiles are of the slow requests.
- In-process Go programs, built with `go build -overlay` under `/tmp`, that call what each method calls. They ran on the main checkout and wrote nothing. `flai guard` refuses `flai hostapi` to a sub-agent, so the programs stood in for `hostapi --timing`.
- The writes were run in a scratch copy of the main checkout under `/tmp`, with a profiling build of flai and shims that log every process it starts.
- The close-out's figures come from the serve journal and the job logs of the week before: 112 jobs, 83 passed.
- The repository then held 3,233 work items, 3,179 of them archived; 5,068 markdown files under `wip`; 884 threads, 351 conversations, 521 digest entries, and 295 issues. flai serve served 13 projects.
- The host was loaded throughout: load average 28 to 78 on 24 cores, with other agents running Go tests. Wall times are inflated about 2 to 3 times; CPU times are given beside them where they were taken.

### What the serve log says

| Method | Slow calls | p50 | p90 | Max | Where the time went (mean) |
|--------|-----------:|----:|----:|----:|----------------------------|
| `agent.status` | 84,417 | 796 ms | 2,473 ms | 17.6 s | `agent.state` 1,212 ms; 390 KB answers. 104,000 s over 4 days, about 30% of serve's wall time |
| `board.get` | 16,641 | 800 ms | 1,822 ms | 8.7 s | `release.pending` 438, `merge.queue` 279, `repo.list` 154, `board.view` 104 |
| `settings.get` | 4,104 | 825 ms | 2,072 ms | 7.0 s | no phases recorded |
| `inbox.designer` | 4,054 | 782 ms | 1,700 ms | 5.4 s | `threads.read` 603 |
| `threads.list` | 748 | 718 ms | 1,287 ms | 4.6 s | `threads.read` 622, `threads.view` 159 |
| `item.get` | 54 | 612 ms | 860 ms | 1.3 s | `repo.list` 578 |
| `stats.get` | 121 | 1,178 ms | 2,415 ms | 6.5 s | `stats.compute` 1,428, `exec.git.log` 168 |
| `item.move` | 35 | 1.3 s | 3.6 s | 6.3 s | `exec.flai.move` 1,808 |
| `item.new` | 8 | 7.9 s | 32 s | 32 s | `exec.flai.story.new` 10,951 |
| `item.edit` | 5 | 8.5 s | 12.4 s | 12.4 s | `exec.flai.edit` 8,598 |
| `item.finalize` | 10 | 6.2 s | 18.8 s | 18.8 s | `exec.flai.edit` 10,046 |
| `issue.story` | 5 | 10.4 s | 19 s | 19 s | `exec.flai.issue.story` 11,562 |
| `accept.run` | 20 | 16.6 s | 40 s | 42 s | `exec.flai.accept` 17,029 |

On the same day an MCP `item_new` of a task took about 20 s.

### Where the time goes

Each cause is numbered after causes 1 to 7. Each "fix" below keeps the files under `wip/` and `design/` as they are, with the same schema. What it keeps is in memory, or in a cache under `.flai-cache` that can be rebuilt from those files.

| # | Cause | Measured | Proposed fix | Estimated saving |
|---|-------|----------|--------------|------------------|
| 8 | A checked write runs the whole of `flai check` twice, before and after (`flai/internal/itemnew/itemnew.go` `Create`, `flai/internal/itemedit/itemedit.go` `Apply`). Each run lists all 3,233 items cold, lints all 5,068 markdown files under `wip`, starts `node` for the component examples, and scans the digest for secrets. An edit also lists the items again and rewrites `wip/agents/index.md`. | `story new`, `epic new`, `task new`, `edit`, finalize, `criteria tick`, `issue story`: 7 to 18 s wall, 8 to 12.7 s CPU, about 97% of it in the two checks | Run only the rules an item write can affect: the written items and their parent, the after, overlap, and touches graph, the board, the story's narrative, and the lint and conflict markers of the written files. Compute the "before" result for those rules only. Done in S-0598: [measured](#a-checked-write-runs-only-the-rules-it-can-affect-s-0598). | 7 to 13 s down to about 1.5 to 2 s a write; measured 12 s down to 2.5 s of CPU a write, and with cause 10 fixed, 6 to 10 s down to about 1 s wall |
| 9 | flai serve's replanner writes each forecast it moves through `itemedit.Apply`, one after another (`flai/internal/serve/replan.go`, `replan`), so each moved forecast pays cause 8's two checks. | 5 to 6.4 s a forecast. Over the week, 98 replans moved 1,200 forecasts, about 1.7 to 2 h of serve CPU. An acceptance's replan commit lands p50 80 s and p90 127 s after the acceptance. | Write every moved forecast in one pass, with one scoped check and one index write. | About 80 s down to under 5 s a replan |
| 10 | Every flai process starts cold. The item cache S-0156 added lives in flai serve's memory, so each `flai` a write starts, each `flai check`, and each `flai move` parses all 3,233 items, one file after another (`flai/internal/workitem/store.go` `list`). `flai move` lists the archive only so that `Follow` sees an epic's archived stories (`flai/internal/workitem/transition.go`). | Item parse 1.1 to 1.5 s cold: 1.28 s with one worker, 0.65 s with 8, 0.49 s with 16. `move` 1 to 2.5 s, 80% of it the list. The lint of `wip`: 0.98 s, 0.25 s on 8 workers. | Parse on a bounded pool of workers. Keep parsed items and lint results under `.flai-cache`, keyed by path, size, and modification time, as S-0531 keeps the stories' commits. Remedied by S-0600: see [Parses and lint findings kept between processes](#parses-and-lint-findings-kept-between-processes-s-0600). | About 1 s a `move`; 1.5 to 2 s a check run |
| 11 | `check.Run` runs its rules one after another (`flai/internal/check/check.go` `Run`), and `componentExamples` starts `node` on every run (`flai/internal/components/props.go`). | Whole check: 3.3 s cold and 1.8 s warm without `node`; 7.4 to 8 s cold and 2.6 s warm with it. `node` alone: 0.6 to 4.2 s. | Run the rules at once; they only read, so the findings need one lock. Keep the props `node` read, keyed by the component files' stat. | Whole check to about 1.2 to 1.5 s; helps `flai verify` and `flai task done` too |
| 12 | `agent.status` reads the host's whole `agents.json`, 663 KB for 14 projects, under one host-wide `agentsMu` (`flai/internal/serve/agents.go` `AgentStates`). It then reads all 884 threads, and all 351 conversations once for each live run, through `awaitingOther`, and again for `held` (`Repo.Holds`). It answers 437 finished runs and 288 planner runs. The dashboard asks for it every 30 s for all 13 projects at once, and the board, the activity page, and a story's agent panel every 15 s. | 1.4 to 2.4 s wall, about 450 ms CPU, 392 KB per call. 8 calls at once: 435 ms each; 32 at once: 1,185 ms. A project with a 1 KB answer still takes p50 780 ms. | Keep `agents.json` parsed, checked by size and modification time. Read the conversations once per call. | 15 to 150 ms and the lock queue; with cause 13, agent.state's CPU from about 450 to under 60 ms |
| 13 | Threads, conversations, digest entries, issues, and activity documents are read and parsed again on every call (`flai/internal/threads/threads.go` `List`, `flai/internal/messages`, `flai/internal/digest`, `flai/internal/issues`). Readers: `agent.state`, `inbox.designer`, `threads.list`, `messages.list`, `board.get`'s `merge.queue`, and `stats.get`. | Threads 80 to 110 ms, conversations 35 to 50 ms, digest 60 to 85 ms, issues 22 to 30 ms per read; YAML parsing is about a fifth of the CPU. | Keep them parsed between requests as S-0156 keeps items: stat each file, and parse again only the files that changed. | About 90 ms down to about 5 ms per thread read; about 180 to 200 ms a `stats.get` |
| 14 | `board.get` holds `release`'s `kept` mutex, one for every repository, across a git process (`flai/internal/release/history.go`), so the 13 projects' boards queue on it. It runs `git show-ref` on every call. `merge.queue` reads all threads even with no story in review (`flai/internal/mergequeue/eligible.go`). It reads each story's undone paths one after another (`flai/internal/hostapi/hostapi.go`). | `release.pending` 23 to 88 ms, 250 ms with 8 at once; `merge.queue` 185 to 630 ms, 94% of it threads; `repo.Holds` 82 to 195 ms | One lock per repository root. Skip `show-ref` while the main-refs watcher saw no change. Read no threads when no story is in review. Read the undone paths at once. | 100 to 500 ms a board under load |
| 15 | `settings.get` works out, for each of the 13 served projects, which projects are below the import folders (`flai/cmd/serve_actions.go` `servedView`, `flai/internal/serve/below.go` `FindBelow`): 13 walks and 169 manifest loads a request. It takes no context, so it records no phases. | 13 × 45 to 85 ms, which matches its p50 of 825 ms | Find them once per request, or take what the running flai serve already found. | About 750 ms a call |
| 16 | `item.get` lists every item, the archive included, to answer one (`flai/internal/hostapi/hostapi.go`), and each folder's lock in the item store is held across its whole stat loop (`flai/internal/workitem/store.go`). The file watcher stats 6,302 files every 300 ms for this project alone, about 3,800 of them under `wip/archive` (`flai/internal/watch/watch.go`). | `item.get` 27 to 41 ms warm, 109 ms with 8 at once, 423 ms with 32. A watcher walk 104 to 217 ms under load, a third to two thirds of a core. | Keep one snapshot of the item list, refreshed when the watcher sees a change. Stat `wip/archive` by its folders rather than its files. | Up to about 500 ms an `item.get` under load; most of a core |
| 17 | `stats.get` runs its 11 reads one after another (`flai/internal/statsread/statsread.go` `Read`), then its sections one after another (`flai/internal/metrics/metrics.go` `Compute`), and two identical requests in a burst each compute in full. The dashboard asks again on every item change. | Warm 0.63 to 1.1 s, cold 1.8 to 4 s; about 200 MB allocated a warm call. Of the 120 slow calls in the serve log, 73 came within 3 s of the one before. | Run the reads at once, then the sections at once; `claims` waits only for `replayHolds`. Answer identical requests in flight once. | Reads 230 to 90 ms; compute about 450 to 280 ms at 30 days |
| 18 | `replayHolds` replays all history on every call, whatever the window (`flai/internal/metrics/claims.go`). `costOfDelay` and `claims` loop day by day over every item, parse every transition's time again each day (`stateAt`), and start at the window's start even before the first item existed. Thread entries are parsed twice, and `costOrder` builds its forecaster three times. | `replayHolds` 175 to 310 ms; `costOfDelay` 2 ms at 1 day up to 170 to 217 ms at 365 days; `time.Parse` about 9% of the CPU | Keep `replayHolds` and `touchDrift` until the items or the commits change. Parse each item's transitions once and sweep the days once; answer the days before the first item without a loop. Parse entries once and build the forecaster once. Done in S-0608: [measured](#cause-18-fixed-the-hold-replay-kept-and-the-days-swept-once-s-0608). | About 200 to 250 ms warm; at 365 days, about 290 ms down to about 30 ms |
| 19 | `flai accept` runs `scripts/lint-md.sh` five times, one after another, each through `npx --yes markdownlint-cli2` (`flai/cmd/accept_markdown.go`), and checks again after each `--fix`. `flai task done` runs it through `npx` twice. `accept` starts 58 git processes, `git status --porcelain` six times, and lists all items for `EpicWalk`. | Each lint call 1.0 to 1.2 s, of which markdownlint is 0.27 to 0.36 s. `accept` 8.7 to 16.2 s in the scratch copy; serve log p50 16.6 s. | Run the installed `markdownlint-cli2` directly, once for the uncommitted and branch sets, and drop the check after `--fix`. | About 4 s an acceptance, about 1.5 s a task done |
| 20 | The close-out runs its tiers one after another (`flai/internal/verify/run.go`). The `flaiover` tier, which is node, waits for `integration`, which is Go, and `smoke` waits for both. `smoke` runs the template test, the check, and the markdown lint again after the close-out's own tiers ran them. | Passed jobs p50 247 s, p90 510 s, max 1,408 s. `integration` mean 238 s; `flaiover` mean 176 s on 12 of 60 jobs; `smoke` mean 44 s. | Run independent tiers at once, and cancel the rest when one fails. Skip the steps of `smoke` the close-out already ran. Built by S-0610, which measured the saving on this host: see [Cause 20: the close-out runs marked tiers at once, and smoke skips what the run passed](#cause-20-the-close-out-runs-marked-tiers-at-once-and-smoke-skips-what-the-run-passed-s-0610). | About 150 s on a story that runs `flaiover`; about 40 s on every story that runs `smoke` |
| 21 | Two answers are large: `threads.list` with `all: true` is 3.28 MB, since the dashboard asks for resolved threads too (`flaiover/src/lib/server/repo.ts`), and `agent.status` carries every finished run. Each holds its connection's one write lock while it is written (`flai/internal/channel/channel.go`), and that time is not in `duration_ms`, which is logged before the send. | 3.28 MB and 392 KB a call | Ask for open threads, and the resolved ones when the page shows them. Send finished runs only to the pages that show them. Log the send time with the request. Done in S-0611: [measured](#cause-21-fixed-the-dashboard-asks-for-open-threads-and-open-runs-and-the-send-is-timed-s-0611). | About 3.2 MB down to under 100 KB a thread list |

### Serial I/O that could run in parallel

The channel does not answer one request at a time. Each connection has one read loop, but each request runs in its own goroutine (`flai/internal/channel/channel.go`). What serializes the work is elsewhere:

- **Shared locks.** `agentsMu` serializes every project's `agent.status` (cause 12). `release`'s `kept` mutex serializes every project's `board.get` and `publish.preview`, git included (cause 14). The item store's folder locks are held across a whole stat loop (cause 16). Each connection's write lock serializes large answers (cause 21).
- **One file at a time.** Items, threads, conversations, digest entries, and issues are each read and parsed in a loop over their files (causes 10 and 13), as is the markdown lint of `wip` (cause 10).
- **One step at a time.** `check.Run`'s rules (cause 11), `stats.get`'s reads and sections (cause 17), the replanner's forecasts (cause 9), `accept`'s lint calls and its last steps, `releaseDue`, `tellAccepted`, and `installAccepted` (cause 19), and the close-out's tiers (cause 20).

On this loaded host, keeping a parse beats running it in parallel. Parallel parsing pays most on a cold start: a one-shot `flai` command, or the first request after flai serve starts.

### What would need a change to how flai keeps state

No cause above needs one. Two fixes the exploration weighed would need one, and no story proposes them:

- One agent state file per project instead of the host-wide `~/.flai/serve/agents.json`. Cause 12's fix removes most of the cost without it.
- A database or an index as the canonical store for work items, threads, and conversations. The kept parses and the rebuildable caches under `.flai-cache` take most of what it would.

Two choices are the operator's and need no story. Turning the `merge-queue` host action on leaves `integration` and `smoke` to the merge queue's batch (ADR-0174), so a close-out passes in seconds rather than minutes. Letting a chart ask for only the sections it shows would change the `stats.get` answer, which `metrics.md` makes a contract, so it needs an ADR first.

### Stories

Each cause from 8 to 21 has a draft story under [E-0033](../../wip/kanban/epics/E-0033-flai-stays-responsive-at-600-stories-checked-writes-kept-parses-and-concurrent.md), with the cause's numbers and its proposed fix. They are drafts for the operator to finalize. Causes 8 and 9 cost the operator most, on every dashboard write and every acceptance; cause 12 costs flai serve most.

| # | Story |
|---|-------|
| 8 | [S-0598](../../wip/kanban/stories/S-0598-a-checked-write-runs-only-the-rules-an-item-write-can-affect-so-creating-or.md) A checked write runs only the rules an item write can affect, so creating or editing an item takes about 2 s, not 7 to 18 s |
| 9 | [S-0599](../../wip/kanban/stories/S-0599-flai-serve-s-replanner-writes-every-forecast-it-moves-in-one-pass-with-one.md) flai serve's replanner writes every forecast it moves in one pass with one check, so an acceptance's replan lands in seconds, not 80 s |
| 10 | [S-0600](../../wip/kanban/stories/S-0600-a-flai-process-parses-work-items-on-a-pool-of-workers-and-keeps-parses-and-lint.md) A flai process parses work items on a pool of workers and keeps parses and lint results under .flai-cache, so a cold list takes under 0.3 s, not 1.1 to 1.5 s |
| 11 | [S-0601](../../wip/kanban/stories/S-0601-flai-check-runs-its-rules-at-once-and-keeps-the-component-props-node-read-so-a.md) flai check runs its rules at once and keeps the component props node read, so a whole check takes about 1.5 s, not 3.3 to 8 s |
| 12 | [S-0603](../../wip/kanban/stories/S-0603-agent-status-keeps-agents-json-parsed-and-reads-the-conversations-once-a-call.md) agent.status keeps agents.json parsed and reads the conversations once a call, so it stops queueing on agentsMu and takes under 200 ms, not 0.8 to 2.4 s |
| 13 | [S-0602](../../wip/kanban/stories/S-0602-flai-serve-keeps-threads-conversations-digest-entries-issues-and-activity.md) flai serve keeps threads, conversations, digest entries, issues, and activity documents parsed between requests, so a thread read takes about 5 ms, not 80 to 110 ms |
| 14 | [S-0604](../../wip/kanban/stories/S-0604-board-get-takes-one-release-lock-per-repository-runs-show-ref-only-when-a-ref.md) board.get takes one release lock per repository, runs show-ref only when a ref moved, and reads no threads with nothing in review, so 13 boards at once stop queueing |
| 15 | [S-0605](../../wip/kanban/stories/S-0605-settings-get-finds-the-projects-below-the-import-folders-once-a-request-not.md) settings.get finds the projects below the import folders once a request, not once a served project, so it answers in under 100 ms, not 825 ms |
| 16 | [S-0606](../../wip/kanban/stories/S-0606-item-get-answers-from-a-kept-item-list-and-the-file-watcher-stops-statting.md) item.get answers from a kept item list and the file watcher stops statting every archived file every 300 ms |
| 17 | [S-0607](../../wip/kanban/stories/S-0607-stats-get-runs-its-reads-at-once-and-its-sections-at-once-and-answers-identical.md) stats.get runs its reads at once and its sections at once, and answers identical requests in flight once |
| 18 | [S-0608](../../wip/kanban/stories/S-0608-the-chart-data-keeps-the-hold-replay-until-items-or-commits-change-and-sweeps.md) The chart data keeps the hold replay until items or commits change and sweeps the days once, so a 365-day window costs what a 30-day one does |
| 19 | [S-0609](../../wip/kanban/stories/S-0609-flai-accept-and-flai-task-done-lint-markdown-once-per-set-with-the-installed.md) flai accept and flai task done lint markdown once per set with the installed markdownlint, not five times through npx, saving about 4 s an acceptance |
| 20 | [S-0610](../../wip/kanban/stories/S-0610-the-close-out-runs-independent-test-tiers-at-once-and-smoke-skips-what-the.md) The close-out runs independent test tiers at once and smoke skips what the close-out already ran, saving about 40 to 150 s a close-out |
| 21 | [S-0611](../../wip/kanban/stories/S-0611-the-dashboard-asks-threads-list-for-open-threads-and-agent-status-for-open-runs.md) The dashboard asks threads.list for open threads and agent.status for open runs, and the request timing counts the send, so a thread list is under 100 KB, not 3.3 MB |

### A checked write runs only the rules it can affect (S-0598)

This section covers cause 8.

#### Cause

A checked creation or edit ran the whole of `flai check` twice, before and after the write. The cost was in the rules no item write can affect: the lint of every file under `wip`, `node` for the component examples, the digest, the documentation, and the threads.

#### Fix

`check.RunItemWrite` replaces `check.Run` in `itemnew.Create` and `itemedit.Apply`. It runs only the rules an item write can affect:

- the item rules on the written items, their parents, and their children;
- the graph rules and the board;
- the written stories' narratives;
- the markdown lint and conflict markers of the written files.

Both runs, before and after, cover the same items and files. What a write refuses is unchanged. [ADR-0264](../adrs/0264-a-checked-write-of-a-work-item-runs-only-the-rules-the-write-can-affect-before.md) records it.

#### Measured on this repository

Setup:

- A scratch clone of the main checkout under `/tmp`, holding 3,278 work items, with the dashboard's `node_modules` linked in so that the old check ran `node` as on the host.
- `flai` built from `main` (before) and from `story/S-0598` (after).
- Five runs of each write, alternating story, edit, and task:
  - `flai story new --epic E-0033 --body-stdin`;
  - `flai edit <that story> --tag cli`;
  - the MCP tool `item_new` of a task under it, through `flai mcp` over stdio, timing the call alone.

The host was heavily loaded: 24 cores, load average 46 to 131 during the runs, with other agents running Go tests. Wall times swung with the load. CPU time, user, is the steadier figure. The figures below are medians, with the range in brackets.

| Write | Before, wall | Before, CPU | After, wall | After, CPU |
|-------|-------------:|------------:|------------:|-----------:|
| `story new` | 18.7 s (9.8 to 47.3) | 12.6 s (10.0 to 21.6) | 2.7 s (2.3 to 4.0) | 2.5 s (2.4 to 3.0) |
| `edit` | 17.5 s (11.5 to 37.9) | 12.0 s (11.1 to 13.2) | 2.4 s (2.2 to 4.3) | 2.4 s (2.3 to 3.2) |
| MCP `item_new` of a task, the call | 18.5 s (8.9 to 21.0) | 12.6 s (10.1 to 13.2), the process | 2.5 s (2.1 to 5.0) | 2.9 s (2.5 to 3.4), the process |

A second set of five after-runs, while the load climbed from 62 to 128, took 4 to 15 s wall each. They still used only 2.9 to 3.6 s of CPU, because the process got 30 to 40% of a core.

A CPU profile of one `story new` after the change, 3.3 s of samples:

| Where | CPU |
|-------|----:|
| `Repo.List`, the first, cold parse of all 3,278 items | 2.1 s |
| `check.ScopeToItem`, which reads the threads and conversations | 0.35 s |
| `touchesFolder` | 0.19 s |
| `unaccepted`, which reads `packed-refs` once per archived done story | 0.13 s |

The cold parse is cause 10. [S-0600](../../wip/kanban/stories/S-0600-a-flai-process-parses-work-items-on-a-pool-of-workers-and-keeps-parses-and-lint.md) parses on a pool of workers and keeps the parses under `.flai-cache`.

I measured again on 2026-11-09, once S-0600 was on `main`. The setup was the same, in a fresh clone of 3,349 items. "Before" is `main` with S-0600's kept parses, which still runs the whole check twice. "After" is `story/S-0598` rebased onto it. The host's load average was 17 to 51.

| Write | Before, wall | Before, CPU | After, wall | After, CPU |
|-------|-------------:|------------:|------------:|-----------:|
| `story new` | 9.6 s (6.1 to 12.1) | 6.4 s (5.8 to 16.5) | 1.0 s (0.8 to 1.3) | 1.0 s (1.0 to 1.0) |
| `edit` | 7.3 s (6.5 to 11.6) | 6.0 s (5.7 to 6.5) | 0.9 s (0.7 to 1.3) | 1.0 s (0.9 to 1.1) |
| MCP `item_new` of a task, the call | 7.9 s (7.2 to 10.9) | 6.3 s (6.3 to 6.9), the process | 1.0 s (0.7 to 1.1) | 1.3 s (1.1 to 1.3), the process |

Every after-run took under 1.5 s wall.

### Cause 9 fixed: the replanner's forecasts written in one pass (S-0599)

#### Cause

flai serve's replanner wrote each forecast it moved through `itemedit.Apply`, one after another (`flai/internal/serve/replan.go`, `replan`). `Apply` is the checked edit `flai edit` makes. It runs the whole `flai check` before and after each write and lists the items again to write `wip/agents/index.md`. So a replan paid cause 8's two checks once per forecast it moved.

#### Fix

The replanner writes every forecast it moves in one call to `itemedit.ApplyForecasts` (`flai/internal/itemedit/forecasts.go`):

- `Repo.SaveAll` saves the items under one hold of the items lock (`flai/internal/workitem/itemlock.go`). Each item is merged as `Repo.Save` merges it ([ADR-0184](../adrs/0184-every-work-item-write-takes-one-items-lock-and-merges-a-copy-read-before.md)). A part both writers changed refuses that item alone, and the rest are written.
- `wip/agents/index.md` is written once, from one list.
- `check.Items` runs once over the written items (`flai/internal/check/items.go`): the per-item rules and the markdown lint of their files. It walks no folder, starts no `node`, and scans no digest. A forecast is a value flai owns, so only an error on an item's own file refuses it, as in `Apply`. That item is put back with `Repo.Restore`.
- Each written item gets one edit notice, as before.

The replanner passes the items its replay read. Each save therefore merges against what the replay read. A forecast another writer changed since the replay is refused and warned of as `forecast not replanned`. A change to another part of the item is kept. The files are still committed in one commit, `chore: replan forecasts after <triggers>`.

#### Measured on this repository

On 2026-11-08, between 00:48 and 01:10 UTC, on the operator's host (WSL2, 24 cores, Go 1.26). The scratch copy was the main checkout copied to `/tmp` without `.git`, `.flai-cache`, or `node_modules`, so the check's git rules found no repository, in both runs alike.

- One program, built from `story/S-0599` with `go build -overlay`, picked 20 open stories in backlog or ready. Only 10 had a forecast, so 10 without one were added. Their edit set a duration of 1h as well.
- It wrote a new delivery and basis on each. "Before" made one `itemedit.Apply` per story, one after another, as the old replanner did. `Apply` on this branch is `main`'s. "After" made one `itemedit.ApplyForecasts` call with the list the program had read.
- It timed the write alone, wall clock and the process's CPU. The commit, which is the same both ways, was not timed.
- `wip` and `.flai-cache` were restored from a copy before each run, so each run moved the same 20 forecasts. Each way ran three times, and every run wrote all 20.

The host was very busy: the one-minute load average was 54 to 136 during the "before" runs and about 83 during the "after" runs.

| Write of 20 forecasts | Runs, wall (CPU) | Median wall | Median CPU |
|-----------------------|------------------|-------------|------------|
| Before: `itemedit.Apply` per forecast | 197 s (157 s), 273 s (163 s), 336 s (176 s) | 273 s | 163 s |
| After: one `itemedit.ApplyForecasts` | 0.08 s (0.09 s), 0.10 s (0.10 s), 0.09 s (0.10 s) | 0.09 s | 0.10 s |

Before, each forecast cost about 14 s of wall time and 8 s of CPU on this host. That is above S-0596's 5 to 6.4 s, which was measured under a lighter load. After, the whole write of 20 takes under a tenth of a second. A replan's time is now its replay and its commit.

### Cause 11: flai check's rules at once and the props node read kept (S-0601)

Two changes:

- `check.Run` runs `workItems` first and alone, since it fills the index of items by ID that the other rules read, and works out the topic vocabulary before any rule. It then runs the other 28 rules at once on at most `GOMAXPROCS` workers. Each rule runs on a copy of the checker with a result of its own, so a rule that marks its last finding, such as `advise`, marks its own. The results are merged in the order `Run` ran the rules before, and sorted stably as before, so the findings, their order, and the counts do not change. `TestRunAtOnceMatchesTheRulesOneAfterAnother` (`flai/internal/check/concurrent_test.go`) checks that over a fixture with a finding of every rule.
- The props `node` read is kept, so a check whose components did not change starts no `node` ([component-library.md › The kept read](component-library.md#the-kept-read-s-0601)).

Measured on 2026-11-08 between 01:55 and 02:00 UTC on the operator's host (WSL2, 24 cores, Go 1.26.0), in the main checkout, with `node` 24.21 on `PATH` and flaiover's `node_modules` installed. That checkout held 3,349 work items. Two binaries were built just before, `main` at 3087d9cf4, which has S-0600's kept parses and lint findings, and `story/S-0601` at 1582a0dd3. Each ran `flai check --strict` as a fresh process, timed by the shell's `time`, five times each, taking turns, after a first run of each that wrote the kept files. Both gave the same output and exited 1 on the same 19 warnings. The host was very busy: the one-minute load average was 84 to 103.

| Warm `flai check --strict` | Wall | CPU (user and system) |
|----------------------------|-----:|----------------------:|
| `main` | 7.9 to 11.7 s | 4.9 to 5.3 s |
| S-0601 | 3.3 to 4.9 s | 3.8 to 3.9 s |

Before S-0600 was accepted, at a load of 61 to 96, `main` at e4e6dff2a took 7.7 to 18.7 s wall and 7.5 to 8.4 s CPU, and S-0601 on it 3.5 to 9.6 s wall and 6.0 to 7.1 s CPU. A first run with no kept props file took 9.3 s against 6.7 to 7.1 s warm.

S-0601 halves the wall time and takes about 1.2 s of CPU off a check. It does not reach the story's 2 s on this host. A small program built with `go build -overlay` timed each step in one process at a load of 103 to 114:

| Step | Time |
|------|-----:|
| `Repo.List(true)`, with S-0600's kept items | 0.66 to 1.10 s |
| `workItems`, alone | about 0.45 to 0.8 s |
| `digest`, the slowest of the rules at once | 0.76 to 1.77 s |
| `conflictMarkers` | 0.45 to 1.10 s |
| `documentation` | 0.28 to 0.85 s |
| `componentExamples`, from the kept read | 0.19 to 0.36 s |
| `markdown`, from S-0600's kept findings | 0.13 to 0.29 s |
| `check.Run` in all, the list included | 2.7 to 3.3 s |

What is left is that critical path: the list, `workItems` alone, and the digest rule, which reads and scans every digest entry. It is recorded as [I-0312](../issues/I-0312-a-whole-flai-check-takes-3-3-to-4-9-s-on-a-loaded-host-after-s-0601-not-under-2.md). S-0677 took `workItems` and the digest rule off that path ([below](#cause-11-again-the-per-item-rules-on-the-pool-and-the-digest-and-conflict-marker-results-kept-s-0677)).

### Cause 11 again: the per-item rules on the pool, and the digest and conflict-marker results kept (S-0677)

#### Cause

After S-0601 three steps still ran on every check, whatever had changed. `workItems` ran first and alone, because it fills the index of items by ID. It also ran each item's rules, and each item scanned every item for its children. The digest rule read and validated every digest entry. The conflict-marker rule read every markdown file under the design, docs, and wip folders.

#### Fix

- `check.Run` runs only `indexItems` first and alone. It fills `byID` and a children index, each parent's children in item order. The per-item rules then run on the pool beside the other rules, in chunks of 128 items. Each chunk has a Result of its own, merged in item order where `workItems` stood, so the findings, their order, and the counts do not change. `TestRunAtOnceMatchesTheRulesOneAfterAnother` checks that at chunk sizes of 1, 2, and 128.
- The digest rule keeps what each entry's file alone decides in `.flai-cache/digest-problems-<hash>.json`, keyed by the entry's size and modification time and a hash of the rules (`flai/internal/digest/kept.go`). Whether a source ID or an `absorbed_by` path names something follows other files, so it is asked on every check.
- The conflict-marker rule keeps each file's marker lines in `.flai-cache/conflict-markers-<hash>.json`, keyed by the file's size and modification time (`flai/internal/conflictmark/kept.go`). A warm run stats each file and reads none that is unchanged.

Both files follow the rules of the kept items and lint findings ([flai-cli.md](flai-cli.md)): written atomically and only into a `.flai-cache` that exists, nothing kept for a file changed within two seconds, and a damaged or foreign file a miss.

#### Measured on this repository

Measured on 2026-11-10 at about 20:08 UTC on the operator's host (WSL2, 24 cores), in the main checkout, which held 3,811 work items. Two binaries were built to `/tmp`: the story's branch at 010fa0f3d, and `main` at 9c2b2392c, the story's base, built from the same tree with the story's files under `flai/` taken from `main`. Nothing else under `flai/` differs between them. Each ran `flai check --strict` as a fresh process, timed by the shell's `time`, five times each, taking turns, after a first run of each that wrote the kept files. Both gave the same output and exited 1 on the same one warning. The one-minute load average was 47 to 56, against 84 to 103 when Cause 11 was measured.

| Warm `flai check --strict` | Wall | CPU (user and system) |
|----------------------------|-----:|----------------------:|
| `main` | 0.97 to 1.15 s | 2.9 to 3.1 s |
| S-0677 | 0.57 to 0.71 s | 2.2 to 2.4 s |

At this load `main` itself is under 2 s, so the fall measured against it is the figure that holds: about 0.4 s of wall time and 0.7 s of CPU off a check. A throwaway test in the `check` package timed each step in one process, six rounds, at a load of 37 to 42. The list's first round in each process is the one a fresh `flai check` pays; the list's code is the same in both.

| Step | `main` | S-0677 |
|------|-------:|-------:|
| `Repo.List(true)`, first in the process | 206 ms | 427 ms |
| The step alone: `workItems`, then `indexItems` | 122 to 167 ms | 2 to 7 ms |
| The per-item rules, run serially to time them | in `workItems` | 52 to 91 ms |
| `digest` | 412 to 513 ms | 30 to 78 ms |
| `conflictMarkers` | 116 to 145 ms | 74 to 134 ms |
| `check.Run` in all, after the first round | 600 to 719 ms | 258 to 544 ms |

What is left on the path is the list, whose decode of the kept items is S-0676's to fix ([I-0310](../issues/I-0310-a-fresh-flai-list-decodes-10-5-mb-of-kept-items-json-so-it-takes-0-3-to-1-4-s.md)), and the conflict-marker rule's walk and stat of every markdown file. At this load the conflict-marker rule reads from the page cache, so keeping its results saves less here than on the loaded host I-0312 measured, where it took 0.45 to 1.10 s.

### Cause 12 fixed: agents.json kept parsed and the conversations read once (S-0603)

#### Cause

Each `agent.status` read and parsed the host's whole `agents.json`, 701 KB on the day measured, under the host-wide `agentsMu`. The 13 projects' calls queued on that lock. `serve.Activity` then read every conversation once for each live run, through `awaitingOther`, and again for `held`, through `Repo.Holds`. The host's `Agent` took no context, so the serve log showed one phase, `agent.state`, and nothing inside it.

#### Fix

- `Dir.AgentStates` (`flai/internal/serve/agents_state.go`) keeps what it parsed with the file's stat. It parses again only when the file, its size, or its modification time changed. It parses outside `agentsMu`, which it holds only to read or replace what is kept. Every write flai makes renames a new file into place, so the clock's granularity cannot hide one.
- What `AgentStates` returns is shared, and its callers only read it. `updateAgent` works on a deep copy, keeps what it wrote with the stat after the write, and serializes writers on a lock of its own, `agentsWriteMu`, with the file lock for other processes. A reader never waits on a write.
- `serve.Activity` reads the conversations once a call, with `messages.ReadAll`. Each live run is judged from that list with `messages.AwaitingOtherIn`. The holds are built from its shares and asks with `Repo.HoldsFrom` and `workitem.SharesAndAsks`. The answer is the same.
- `Host.Agent` takes the request's context. Inside `agent.state` the serve log now shows `agent.states`, reading the agent states; `agent.roles`, the roles' states; and `agent.activity`, `serve.Activity`.

#### Measured on this repository

On 2026-11-08 between about 02:08 and 02:11 UTC on the operator's host (WSL2, 24 cores, Go 1.26.0). An in-process Go program, built under `/tmp` against each tree, did per call what the host's `Agent` does, without the configuration: `AgentStates`, `RoleStates`, and `Activity`. It read the host's own `~/.flai/serve/agents.json` and the 13 projects it names that exist. Before is `main` at 3087d9cf4; after is `story/S-0603` with T-2621, T-2629, and T-2623. The two ran in turn, twice each.

- Sequential: one warm call per project, after a warm-up round.
- At once: 13 calls, one per project, released together, five rounds a run: 65 calls.

The host was very busy throughout: the one-minute load average was 95 to 101, on 24 cores, with other agents running Go tests. Wall times swing with it.

| Case | Before | After |
|------|--------|-------|
| This repository, sequential | 1,494 to 2,391 ms | 966 to 1,201 ms |
| The other 12, sequential, p50 | 54 to 90 ms | 64 to 73 ms |
| 13 at once, a call, p50 | 369 to 463 ms | 221 to 283 ms |
| 13 at once, a call, p90 | 675 to 925 ms | 395 to 451 ms |
| 13 at once, a call, max | 2,620 to 3,288 ms | 1,542 to 1,593 ms |
| 13 at once, a round's wall time | 1.8 to 3.3 s | 1.0 to 1.6 s |
| 13 at once, a round's CPU | 1.44 to 1.62 s | 0.74 to 1.14 s |
| 13 at once, calls under 200 ms | 11 and 11 of 65 | 26 and 14 of 65 |
| 13 at once, this repository, p50 | 2,092 to 2,952 ms | 1,097 to 1,107 ms |

| Phase, 13 at once | Before p50 (max) | After p50 (max) |
|-------------------|------------------|-----------------|
| `agent.states` | 315 to 327 ms (892) | 4 to 159 ms (289) |
| `agent.roles` | 32 to 51 ms (182) | 99 to 111 ms (367) |
| `agent.activity` | 17 ms (2,482) | 21 to 29 ms (1,262) |

The queue on `agentsMu` is gone. Reading the agent states took a third of a second at the median with 13 at once; it is now a stat, and a sequential call takes under 1 ms. What it still takes with 13 at once, and the rise in `agent.roles`, whose code did not change, is the process waiting for a core on a host at a load of 100. With five live runs, this repository's `Activity` reads the conversations once instead of six times, and its call takes about half the time.

The story's target, every call under 200 ms with 13 at once, is not met on this host. At head, 14 and 26 of 65 calls answered under 200 ms. This repository's call took about 1.1 s in every round, and the other projects' p50 was about 200 to 280 ms, waiting for cores. Alone, this repository's `Activity` took 540 to 600 ms, of which reading its 909 threads took 180 to 210 ms, its 381 conversations 70 to 180 ms, and building the holds 85 to 150 ms; the rest is the trial merges of the holds, one `git merge-tree` for each overlap with a story in progress, and `elsewhere`. The threads and the conversations are cause 13's, which S-0602 keeps parsed between requests.

The program wrote nothing of its own. Like every `agent.status` in flai serve, its trial merges wrote tree objects into each project's object store, with no ref and no worktree, and its item lists may have written the rebuildable kept items files under `.flai-cache` that S-0600 added.

### Cause 13 fixed: threads, conversations, digest entries, issues, and activity documents kept parsed (S-0602)

Every call read and parsed each of these files again: 884 threads, 351 conversations, 521 digest entries, 295 issues, and the strategic agents' activity documents. Now a long-running flai keeps each kind parsed, as S-0156 keeps items.

- `workitem.KeptFiles` (`flai/internal/workitem/keptfiles.go`) keeps the parse of files of one kind, by absolute path. `List` reads a folder and stats each file it matches. `Read` stats one file. A file is parsed again only when it is new, is not the same file (`os.SameFile`), or its size or modification time changed. A file read within `racyWindow` of its change is parsed again on every look until it is older. A file gone from the folder is dropped. One lock per folder lets readers at once share one parse. Each caller gets its own copy through the kind's clone, and a parse error is not kept.
- Each kind keeps one store, process-wide: `threads.List` and `Get`, `messages.List` and `Get`, `digest.List` and `Get`, `issues.List` and `Get`, and `Repo.Activity` and `Repo.Activities`. `Repo.Conversations`, which `Repo.Holds` calls for `agent.state` and `board.get`, keeps its own store of each conversation's shares and asks. `workitem` cannot import `messages`, so the two stores are separate. Unknown fields are warned of once per parse, not on every list. A write replaces the file, or rewrites it in place as `issues.Save` does, and the next read sees it by its stat.
- Each kind's tests show that a changed, an added, and a removed file is seen on the next read, that a second read with nothing changed parses no file, that readers at once get equal answers, and that changing a returned value changes nothing kept.

Measured on 2026-11-08 between 03:25 and 03:30 UTC on the operator's host (WSL2, 24 cores, Go 1.26.0), on the main checkout. Before is `main` at 17097cea1; after is `story/S-0602` at cdb0441a5. Both binaries were built under `/tmp` with `go build -overlay`, which changes no file in the repository. A one-shot `flai hostapi --timing` is always cold, so the methods were timed by a build of each tree whose `hostapi --timing` calls the method six times in one process. The first call is cold and the other five warm. The two binaries took turns for five rounds, 25 warm calls each per method. The load average over one minute was 37 to 49 at the rounds' starts.

| Method | Phase | Before, warm | After, warm | Before, cold | After, cold |
|--------|-------|-------------:|------------:|-------------:|------------:|
| `threads.list` | `threads.read` | 262 ms | 14 ms | 345 ms | 385 ms |
| `threads.list` | whole call | 271 ms | 28 ms | 367 ms | 402 ms |
| `inbox.designer` | `threads.read` | 233 ms | 13 ms | 217 ms | 333 ms |
| `inbox.designer` | whole call | 354 ms | 131 ms | 385 ms | 508 ms |
| `stats.get` | `stats.compute` | 942 ms | 807 ms | 1,388 ms | 1,531 ms |

The figures are medians. A cold call parses every file, as before. `stats.get` records no phase of its own for its reads. So each read, and `statsread.Read` as `stats.get` calls it, was timed apart. A program built from each tree under `/tmp` opened the main checkout and called each read six times in one process. It recorded the wall time and the process's CPU time of each call. The two binaries took turns for eight rounds; the load average was 15 to 19.

| Read | Before, warm | After, warm | Before, warm CPU | After, warm CPU |
|------|-------------:|------------:|-----------------:|----------------:|
| `threads.List` | 99.5 ms | 7.5 ms | 133.0 ms | 7.6 ms |
| `messages.List` | 56.2 ms | 2.4 ms | 80.8 ms | 2.4 ms |
| `digest.List` | 95.4 ms | 6.2 ms | 140.2 ms | 6.3 ms |
| `issues.List` | 33.8 ms | 1.9 ms | 47.5 ms | 1.9 ms |
| `Repo.Activities` | 7.3 ms | 0.3 ms | 11.5 ms | 0.3 ms |
| `statsread.Read` | 95.8 ms | 49.5 ms | 388.8 ms | 128.9 ms |

The five reads together took 292 ms warm before and 18 ms after: a warm thread read is about 7 ms, not 80 to 110 ms. `threads.list` and `inbox.designer` spend about 220 to 250 ms less in `threads.read`. `stats.get`'s reads spend 260 ms less CPU. Their wall time falls only 46 ms, from 96 to 50 ms, because S-0607 runs the twelve reads at once, so the wall time is that of the slowest. What is left of `statsread.Read` is the reads this story does not keep: the stories' commits, the board, the conflict log, and the merge queue's events.

### Cause 14 fixed: board.get's release lock per root, the watcher's word, and no threads with nothing in review (S-0604)

- The kept history in `flai/internal/release/history.go` is one `keptRoot` per repository root, each with its own lock. The map's lock is held only to find or add a root's entry, never across a git process. `PendingBatch` holds only its root's lock, so one project's `git show-ref` and `git log` never hold another project's board. `TestPendingOnOneRootAnswersWhileAnotherRootsGitRuns` holds one root inside its `show-ref` and shows a second root answers meanwhile.
- flai serve's main-refs watcher, `mainRefs` in `flai/internal/serve/mergequeue.go`, now watches `.git/refs/tags` beside `HEAD`, `packed-refs`, and `refs/heads`. Each served project's merge queue pass calls `release.Watch` as it starts, and `release.RefsMoved` on every change the watcher reports and once a tick after its first look. `PendingBatchWatched`, which `PendingIDs` calls for `board.get`, the inbox, and the orchestrator's figures, reuses the kept history while a watcher covers the root and no move came since it was read. A generation counter, read before the `show-ref`, keeps a move reported during a read from being lost. The publish that tags and pushes, `release.EvaluateRepo`, and `publish.preview` call `PendingBatch`, which still reads git. Processes other than flai serve never call `Watch`, so they read git as before.
- `mergequeue.ComputeAt` (`flai/internal/mergequeue/eligible.go`) finds the stories in review first. With none, it answers the queue's saved state with no entries and no window, and reads no thread and runs no `git status`.
- `board.get`'s `undone.read` phase (`flai/internal/hostapi/hostapi.go`, `boardUndone`) reads the ready, in-progress, and review stories' undone paths at once, on at most 8 goroutines, each writing only its own card. A story in review is still read again by the merge queue; sharing the read would change `mergequeue`'s API.

Measured on 2026-11-08 between 02:34 and 02:36 UTC on the operator's host (WSL2, 24 cores, Go 1.26.0). An in-process Go program, built under `/tmp` against each tree, called `board.get` through `hostapi.MethodsFor` with a request timer, for the 13 projects in `~/.flai/serve/projects.json`. It ran a watcher over `mainRefs` for each root, as flai serve does; only the story's build had `release.Watch` and `release.RefsMoved` to call. Each run made a warm-up call and one sequential call per project, then five rounds of 13 calls released together: 65 calls a run. Before is `main` at 3087d9cf4; after is `story/S-0604` at 481db7030, with T-2659, T-2666, T-2670, and T-2677. The two ran in turn, twice each, so 130 calls a build. The one-minute load average was 80 to 104 on 24 cores, with other agents running Go tests.

| 13 at once, a call | Before | After |
|--------------------|--------|-------|
| All 13 projects, p50 | 41 ms | 14 ms |
| All 13 projects, p90 | 127 ms | 217 ms |
| All 13 projects, max | 529 ms | 761 ms |
| The other 12, p50 (max) | 36 ms (190) | 10 ms (293) |
| This repository, p50 (range) | 415 ms (274 to 529) | 442 ms (342 to 761) |
| Calls under 300 ms | 121 of 130 | 120 of 130 |
| Calls that ran `git show-ref` | 70 of 130 | 0 of 130 |
| Calls that ran `git status` | 130 of 130 | 40 of 130 |

| Phase, 13 at once, p50 (max) | Before | After |
|------------------------------|--------|-------|
| `release.pending`, all | 3.9 ms (96) | 0.0 ms (202) |
| `merge.queue`, all | 21.8 ms (309) | 0.2 ms (339) |
| `repo.list`, this repository | 56 ms (105) | 93 ms (269) |
| `release.pending`, this repository | 34 ms (91) | 26 ms (76) |
| `board.view`, this repository | 71 ms (101) | 69 ms (124) |
| `merge.queue`, this repository | 216 ms (309) | 224 ms (339) |
| `undone.read`, all | under 1 ms | under 1 ms (154) |

The board's own git is gone. No call ran `show-ref`, and `git status` ran only for the four projects with a story in review. With nothing in review, a project's `merge.queue` is now a read of the queue's state, and the other 12 projects' p50 fell from 36 to 10 ms. `release.pending` still takes 20 to 76 ms for this repository without starting git: that is working out its plans from every item, the archive included. The after run's maxima on other projects, up to 202 ms in `release.pending` with no git process, came in one round, together, while the host stalled.

The story's target, every call under 300 ms with 13 at once, is not met for this repository. The other 12 answered under 300 ms in every call of both runs. This repository answered in 342 to 761 ms after, as before. Four of its stories were in review, so its `merge.queue` read all its threads, about 200 ms (cause 13, S-0602), and its `repo.list` and `board.view`, the item list and the holds, took 60 to 270 ms (cause 16, S-0606). The serve log's p50 of 800 ms was taken inside flai serve, beside the other requests and the shared lock this story removed. One process making only board calls does not have that contention, so `main` measured 41 ms at p50 here. On TH-0912 the operator moved the target of every board under 300 ms to [S-0636](../../wip/kanban/stories/S-0636-13-board-get-calls-at-once-answer-under-300-ms-on-every-served-project-measured.md), which measures again once S-0602 and S-0606 have landed: [13 boards at once, measured again after S-0602 and S-0606 (S-0636)](#13-boards-at-once-measured-again-after-s-0602-and-s-0606-s-0636).

The program wrote nothing of its own. Like every `board.get`, its item lists may have written the rebuildable kept items files under `.flai-cache` that S-0600 added.

### 13 boards at once, measured again after S-0602 and S-0606 (S-0636)

S-0636 measures S-0604's target again: 13 `board.get` calls at once, one per served project, each under 300 ms.

#### How it was measured

On 2026-11-09 between 23:04:36 and 23:05:33 UTC on the operator's host (WSL2, 24 cores, Go 1.26.0), as S-0604 measured cause 14.

- Two trees were copied to `/tmp` and built there. Before is `main` at 8b29cf861, the commit before S-0602's acceptance: it has none of S-0602, S-0604, or S-0606. After is `main` at f667efc2d, with all three, and with every story accepted between them.
- One in-process program, built against each tree, called `board.get` through `hostapi.MethodsFor` with a request timer, for the 13 projects in `~/.flai/serve/projects.json`.
- Each build ran per project what its flai serve runs:
  - Before: the file watcher over the manifest and its three folders, and the merge queue's watcher over `mainRefs`.
  - After: the file watcher with its slow tier, telling the item store of each change, with the item snapshot on (S-0606); and the watcher over `mainRefs`, with `release.Watch`, `release.RefsMoved` on each change, and the word a tick after its first look (S-0604).
- Each run made a warm-up call and one sequential call per project, then five rounds of 13 calls released together: 65 calls. The builds ran in turn, before, after, before, after: 130 calls a build.
- The one-minute load average was 123 at the start and 96 at the end, on 24 cores, with other agents running Go tests.
- This repository had 4 stories in review as the first run started and 5 from its end on.

#### What it found

| 13 at once, a call | Before | After |
|--------------------|--------|-------|
| All 13 projects, p50 | 71 ms | 13 ms |
| All 13 projects, p90 | 283 ms | 136 ms |
| All 13 projects, max | 657 ms | 404 ms |
| The other 12, p50 (max) | 66 ms (359) | 9 ms (216) |
| This repository, p50 (range) | 405 ms (244 to 657) | 297 ms (169 to 404) |
| Calls under 300 ms | 118 of 130 | 125 of 130 |
| This repository's calls under 300 ms | 3 of 10 | 5 of 10 |
| Calls that ran `git show-ref` | 70 of 130 | 0 of 130 |
| Calls that ran `git status` | 130 of 130 | 30 of 130 |

| Phase, 13 at once, p50 (max) | Before | After |
|------------------------------|--------|-------|
| `release.pending`, all | 10.4 ms (296) | 0.0 ms (39) |
| `merge.queue`, all | 33.6 ms (337) | 0.0 ms (241) |
| `repo.list`, all | 0.4 ms (219) | 0.0 ms (62) |
| `board.view`, all | 1.6 ms (83) | 1.3 ms (90) |
| `undone.read`, all | under 0.1 ms (33) | under 0.1 ms (20) |
| `repo.open`, this repository | 13.8 ms (31) | 22.9 ms (33) |
| `repo.list`, this repository | 103 ms (219) | 32 ms (62) |
| `release.pending`, this repository | 41 ms (240) | 20 ms (28) |
| `board.view`, this repository | 37 ms (83) | 21 ms (90) |
| `merge.queue`, this repository | 156 ms (337) | 180 ms (241) |
| `git status` in `merge.queue`, this repository | 88 ms (244) | 95 ms (148) |
| five `git rev-parse` together, this repository | 31 ms (55) | 41 ms (81) |

The other 12 projects answered under 300 ms in every call of both after runs, against 115 of 120 before. Their p50 fell from 66 to 9 ms.

This repository did not answer under 300 ms in every call. Its after calls took 169 to 404 ms, and 5 of 10 took 300 ms or more: 373, 318, 304, 350, and 404 ms. S-0602 and S-0606 took out what S-0604 left to them. `repo.list` fell from 103 to 32 ms at the median, and `release.pending` from 41 to 20 ms. The thread read no longer shows in `merge.queue`.

What is left is the merge queue's git, which runs only while stories are in review (`mergequeue.ComputeAt`, `flai/internal/mergequeue/eligible.go`):

- `mainDirt` runs `git status --porcelain` on the main checkout: 52 to 148 ms a call after.
- `git rev-parse` ran five times a call, one after another, for the heads `verified` and `mainDirt` read: 21 to 81 ms together.

With five stories in review, `merge.queue` took 91 to 241 ms of each call, most of it that git. Its median did not fall: the thread read it lost was a smaller share than the git, whose time swings with the host's load. The rest of a call, `repo.list`, `board.view`, `repo.open`, `release.pending`, and the wait for a core, took 63 to 194 ms. The load average was 96 to 123 on 24 cores throughout, so every figure here is inflated, as in the earlier sections.

The program wrote nothing of its own. No file under this repository's `.flai-cache` changed while it ran, worktrees aside.

### Cause 15 fixed: settings.get finds each project's import folder from its path (S-0605)

The planner's first idea was to take the projects below the import folders from the running flai serve's status. That status cannot say which import folder a registered project is below, since `serve.Place` leaves registered projects out of `ImportProjects` and `Unserved`. So the walk is replaced, not cached.

- `serve.ImportRootOf` (`flai/internal/serve/below.go`) answers which import folder a root is below from the root's own path. It follows each import folder's walk down that path alone, by the walk's rules: no hidden or skipped folder, no symbolic link, at most `candidateDepth` deep, no repository above the root. The root must be a repository with a `system-flow.yaml`. That is a few stats a project, with no walk and no manifest load. `TestImportRootOfAgreesWithTheWalk` checks it against the walk on every folder of a tree that has each of those cases, with overlapping import folders.
- `servedView` (`flai/cmd/serve_view.go`) reads the configuration once a call. It no longer asks for the candidates, which the settings page does not show, so `serve.FindCandidates` no longer walks the import folders for it. While flai serve runs, a `settings.get` walks no folder. When it does not run, `belowImportRoots`' one `FindBelow` is the only walk. `TestServedViewSaysWhichImportFolderEachProjectIsBelow` (`flai/cmd/serve_project_test.go`) checks the answer for three served projects under two import folders, with flai serve running and stopped.
- `settings.get` passes its context to `Host.Settings`, which records three phases: `settings.config`, the host configuration and what it says of the project; `settings.served`, the projects flai serve serves; and `settings.project`, the project's own manifest and MCP server.

Measured on 2026-11-08 between 00:50 and 00:55 UTC on the operator's host (WSL2, 24 cores, Go 1.26.0), with the host's own configuration: flai serve running, 13 served projects, and 3 import folders. Both binaries ran `flai hostapi --timing settings.get` from the main checkout. Before is `main` at 9fa3fec41; after is `story/S-0605` at 2e6f35ff4. The two binaries took turns for eight rounds, and then the story's binary ran 20 more times. Both answered the same `host.projects`, byte for byte, and the same keys. The host was very busy: the load average was 58 to 63 over one minute and about 128 over fifteen.

| Binary | Runs | p50 | Range | Where the time went (typical run) |
|--------|-----:|----:|-------|-----------------------------------|
| `main` | 8 | 336 ms | 268 to 601 ms | no phases recorded |
| S-0605 | 28 | 91 ms | 62 to 211 ms | `settings.served` 33, `settings.project` 26, `settings.config` 17 |

The serve log's p50 of 825 ms was taken inside flai serve, beside the other requests it was answering. A one-shot process does not have that contention, so `main` measured 336 ms here. The 13 walks are gone. In the 28 runs, 18 answered under 100 ms; the slowest ran during the busiest moments. What is left is in `settings.served`: `readServeStatus` opens every registered project to describe it. `settings.project` is the project's own `workitem.Open`. Neither grows with the import folders.

### Cause 16: the item snapshot and the watcher's slow tier (S-0606)

#### Cause

`item.get` listed every item, the archive included, and each list statted every item file, holding each folder's lock across the loop. flai serve's watcher statted every file under `design`, `docs`, and `wip` every 300 ms, about 3,800 of them under `wip/archive`, which changes only when an item is archived or edited by hand.

#### Fix

- While flai serve serves a project, the item store keeps a snapshot of each item folder (`flai/internal/workitem/store.go`). A list or a get of a fresh folder reads and stats nothing. The watcher's changes, the process's own writes, and the flai subprocesses serve runs mark a folder stale, and so does a look once a minute; a stale folder refreshes through the stat-and-compare path of S-0156. [flai-cli.md](flai-cli.md) says what marks it.
- The watcher has a slow tier (`flai/internal/watch/watch.go`). Its folders are statted on every poll, which shows a file added, renamed in, or removed. Its files are statted every 10 s, which shows an edit in place.
- flai serve puts six folders on the slow tier (`slowPaths` in `flai/internal/serve/serve.go`): `wip/archive`, `wip/threads`, `wip/messages`, `design/digest`, `design/issues`, and `design/adrs`. The archive was there from the start. The other five joined after the first measurement below, as the operator chose on TH-0908.
  - flai writes every one of these folders by rename. Threads, conversations, and digest entries already were. Since T-2699, issues, the issues summary, ADRs, and the ADR index are too (`flai/internal/atomicfile`).
  - A rename moves the folder's modification time, so the fast look sees flai's writes at once.
  - Only an edit made in place by hand in those folders waits for the slow look, up to 10 s.
- `item.get` copies only what it answers (T-2697, `flai/internal/workitem/view.go`). It takes the item with `Repo.Get`, its children with `Repo.ChildrenOf`, which copies only them, and the stories' usage for the cost per agent hour with `Repo.Visit`, which reads the kept items in place. It no longer copies and sorts the whole list.

#### Measured on this repository

On 2026-11-08, between 01:21 and 01:35 UTC, on the operator's host (WSL2, 24 cores, Go 1.26.0), against the main checkout, read only: 6,478 files under the watched paths, 3,783 of them under `wip/archive`, and 3,347 work items, 3,202 of them archived. Two in-process programs were built at 01:20 UTC from copies of the `flai` folder, one with `main` (e4e6dff2a) and one with `story/S-0606` (52f9dce39), each with a hook exporting the watcher's walk. They ran one after the other, twice each.

- The walk is `Watcher` over the manifest, `design`, `docs`, and `wip`, 20 walks a run after one first walk, with `GOMAXPROCS=1`. CPU is the process's user and system time across each walk. S-0606's fast look is a walk between slow looks; its slow look is one where the archive's files are statted.
- `item.get` is the host API's method table, `hostapi.MethodsFor`, in one process: 13 calls at once, for `S-0241`, `S-0612`, `S-0615`, `S-0408`, `T-1267`, `T-1695`, `E-0033`, `S-0001`, `S-0100`, `S-0301`, `S-0506`, `T-0800`, and `E-0003`, three warm rounds and then 20 rounds, 260 calls a run. None is a story in ready, in progress, or review, so no call reads the undone paths through git (cause 14). S-0606 turned the snapshot on before the first round.

The host was very busy: the one-minute load average was 110 to 172 on 24 cores throughout. Wall times are inflated several times over; the CPU times are the steadier figure.

| What | `main` | S-0606 |
|------|--------|--------|
| Files statted a walk | 6,478 | fast look: 2,695 and 6 archive folders; slow look: 6,478 |
| Walk, wall, median (p90) | 169 ms (380), 159 ms (228) | fast look 99 ms (403), 111 ms (198) |
| Walk, CPU, median (p90) | 74 ms (83), 77 ms (85) | fast look 44 ms (51), 45 ms (50); slow look 80 ms (84), 84 ms (93) |
| `item.get`, 13 at once, wall per call, median (p90) | 1,112 ms (2,062), 1,378 ms (3,249), 811 ms (1,719) | 171 ms (280), 346 ms (726), 218 ms (343) |
| `item.get`, CPU of a whole run of 299 calls | 18.1 s user, 8.9 s system | 11.1 s user, 1.3 s system |

A walk between slow looks costs about 45 ms of CPU, against about 75 ms. `item.get` with 13 at once is about five times faster per call, and spends a seventh of the system time.

In this first run neither met S-0606's targets, a walk under 30 ms and `item.get` with 13 at once under 50 ms, on this host under this load:

- The fast look still stats the 2,695 files outside the archive: 908 threads, 526 digest entries, 370 conversations, 297 issues, and 264 ADRs make up most of them. A third program put `wip/threads`, `wip/messages`, `design/digest`, `design/issues`, and `design/adrs` on the slow tier too: a fast look then cost 22 to 23 ms of CPU, p90 27 to 30 ms. That change was not made then: those folders change often, and an edit in place there would reach the dashboard up to 10 s late. TH-0908 later chose it, with flai's own writes there made atomic; it is measured below.
- With the snapshot on, a call reads no file, but it still copies and sorts the whole list. A CPU profile of S-0606's `item.get` over 23 rounds found `Repo.List` in 67% of the samples: the deep copy of every item 46%, its sort 20%. `workitem.Open`, which loads the manifest again for every call, took 12%. That is about 27 ms of CPU a call, so 13 calls at once queue on a loaded host. Copying only the items an answer returns, and keeping the sorted list with the snapshot, would take most of it.

#### Measured again after T-2697

On 2026-11-08, between 02:09 and 02:11 UTC, on the same host, against the main checkout, read only: 6,498 files under the watched paths, 3,791 of them under `wip/archive`, and 3,349 work items, 3,209 of them archived. The two programs were built at 02:08 UTC, one from `main` at 3087d9cf4, which now holds S-0600's worker pool and kept items file, and one from `story/S-0606` at d0d2cff8c, with T-2697. Both copies read the kept items file but were changed not to write it, so that nothing was written in the main checkout. The calls, the IDs, and the rounds are those above. The runs were interleaved, two of each. The one-minute load average was 86 to 101 on 24 cores. CPU of a run is the process's user and system time, its start and first list included.

| What | `main` | S-0606 |
|------|--------|--------|
| Walk, wall, median (p90) | 56 ms (66), 65 ms (93) | fast look 44 ms (95), 34 ms (66); slow look 72 ms (81), 74 ms (104) |
| Walk, CPU, median (p90) | 57 ms (64), 60 ms (66) | fast look 37 ms (43), 34 ms (55); slow look 67 ms (69), 66 ms (72) |
| `item.get`, 13 at once, wall per call, median (p90) | 456 ms (787), 415 ms (724) | 15 ms (32), 13 ms (25) |
| `item.get`, 13 at once, wall per round of 13, median (p90) | 819 ms (884), 738 ms (1,022) | 40 ms (58), 29 ms (41) |
| `item.get`, CPU of a whole run of 299 calls | 10.9 s user, 7.4 s system; 10.6 s user, 6.9 s system | 2.3 s user, 0.3 s system; 2.0 s user, 0.3 s system |

`item.get` with 13 at once now takes 13 to 15 ms a call at the median and 25 to 32 ms at p90, under the 50 ms target, about 30 times faster than `main` with S-0600's pool. A walk between slow looks takes 34 to 44 ms wall and 34 to 37 ms of CPU at the median, against 56 to 65 ms and 57 to 60 ms on `main`. At that point it still missed the 30 ms target, for the reason above: the fast look statted the 2,700 files outside the archive.

#### Measured again after T-2700

The walk was measured again on 2026-11-09, between 08:08 and 08:09 UTC, on the same host, against the main checkout, read only. The watched paths held 6,560 files, 6,269 of them below the six slow folders.

- Two programs were built at 08:08 UTC from copies of the `flai` folder: one from `main` at 21b0c7729, one from `story/S-0606` at 83c527f70, which has T-2699 and T-2700. Each had a hook exporting the watcher's walk.
- The branch's watcher took serve's paths, the manifest with `design`, `docs`, and `wip`, and serve's six slow folders, with an hour's slow interval, so that every fast look fell between slow looks. Its slow looks were forced one at a time.
- Each run was 20 walks after one first walk, with `GOMAXPROCS=1`. The runs were interleaved, two each.
- The one-minute load average was 49 to 50 on 24 cores.

| What | `main` | S-0606 |
|------|--------|--------|
| Files statted a walk | 6,560 | fast look: 291 files and 11 slow folders; slow look: all 6,560 files and the 11 folders |
| Walk, wall, median (p90) | 63 ms (94), 67 ms (142) | fast look 13 ms (16), 16 ms (22); slow look 58 ms (63), 70 ms (100) |
| Walk, CPU, median (p90) | 57 ms (64), 60 ms (65) | fast look 13 ms (16), 16 ms (20); slow look 58 ms (61), 65 ms (73) |

A look between slow looks now takes 13 to 16 ms of CPU and wall at the median, and 16 to 22 ms at p90. That is under the 30 ms target, about a quarter of `main`'s 57 to 60 ms. A slow look, every 10 s, still costs what a walk on `main` does.

### Cause 17 fixed: stats.get runs its reads and its sections at once and shares a computation in flight (S-0607)

`stats.get` ran its reads one after another, then its sections one after another, and two identical requests in a burst each computed in full. All three now overlap, and the answer is the same.

- `statsread.Read` (`flai/internal/statsread/statsread.go`) runs its twelve reads at once, one goroutine each: the items, the activity documents, the roles, the threads, the issues, the board, the stories' commits, the conversations, the conflict log, the merge queue's events, the digest, and, under pull-request acceptance, the forge events. Each goroutine writes only its own variables, and the options are filled once all have ended. When several reads fail, the error is still that of the first in the old order, with its wording; each warning is still logged with its detail, though the warnings may now come in another order. `TestReadReturnsTheFirstReadsErrorWhenTwoFail` checks the error.
- `metrics.Compute` (`flai/internal/metrics/metrics.go`) runs its sections in three stages. The sections that need only the inputs start at once and run beside everything else. The per-item loop's inputs, `threadWaits`, `forgeWaits`, `mergeQueue`, and `CostPerAgentHour`, are built at once beside `replayHolds`, and `claims` starts as soon as `replayHolds`, which it reads, is done. After the loop, `summarise`, the groups, `forecasts`, `waiting`, and `strategicDays` run at once, with `aging` and `readyWait` after `summarise`, whose figures they read. Each goroutine writes only its own field of the report, and no section changes what another reads.
- `TestReportOverTheGoodFixtureIsTheGolden` (`flai/internal/statsread/report_test.go`) pins the report over `flai/internal/metrics/testdata/good` for 1, 7, 30, and 365 days, by nature, and by the hour, as the serial code wrote it in `testdata/good-report.json`. It passes under `-race` with the reads and the sections at once.
- `stats.get` (`flai/internal/hostapi/reads.go`) shares one computation between the requests that ask for the same report while it runs. The key is the project's root and the parameters after their defaults: the window, the type, `by`, and the bucket. The work runs under `context.WithoutCancel`, so the end of the request that started it does not end it for the others. Each request stops waiting when its own context ends. The key is dropped when the computation ends, so no answer is kept past it. A request that joins records the time it waited as its `stats.compute` phase. `TestStatsGetAnswersIdenticalRequestsInFlightOnce` checks three cases. Two requests with the same parameters, started together, compute once and get the same answer. Two with different parameters compute twice. A request after the first ended computes again. `TestAStatsGetThatJoinedStopsWaitingOnItsOwnContext` checks that a request whose context ends stops waiting, while the first still gets its answer.

Measured on 2026-11-08 between 02:36 and 02:42 UTC on the operator's host (WSL2, 24 cores, Go 1.26.0), on the main checkout. One in-process program, built once from each tree, opened the main checkout and ran `statsread.Read` and then `metrics.Compute` as `stats.get` does, with the real runner, timing the two apart. Before is `main` at 3087d9cf4; after is `story/S-0607` at 16e98f73b. Each run made one call per window that is left out, then 5 timed calls per window. The two binaries took turns for five rounds, 25 timed calls each per window. With `Now` fixed, both gave the same report, byte for byte, for every window. The host was very busy: the load average over one minute was 48 to 79 at the rounds' starts, and 83 to 94 over fifteen.

The table gives the median of the five rounds' medians, in milliseconds, and the range of the rounds' total medians.

| Window | Before: read | Before: compute | Before: total | Before: range | After: read | After: compute | After: total | After: range |
|--------|-------------:|----------------:|--------------:|---------------|------------:|---------------:|-------------:|--------------|
| 1 day | 697 | 1,010 | 1,714 | 897 to 1,839 | 229 | 491 | 703 | 409 to 1,211 |
| 7 days | 719 | 996 | 1,776 | 1,011 to 2,191 | 263 | 478 | 785 | 466 to 1,292 |
| 30 days | 686 | 1,138 | 1,861 | 1,373 to 2,526 | 197 | 476 | 654 | 391 to 1,122 |
| 365 days | 550 | 1,358 | 1,855 | 1,259 to 2,285 | 151 | 499 | 716 | 600 to 1,489 |

A warm 30-day call took about 1.2 s less, and the story's binary was faster in every round. Its saving per round, base against story at the same window, was 929 to 1,741 ms at 30 days, and 162 to 1,741 ms across all windows. That is more than S-0596 estimated, reads from 230 to 90 ms and compute from 450 to 280 ms. Both figures are inflated on this host: S-0596 found wall times two to three times as long under such load, and work spread over goroutines takes more of a loaded machine's cores than work on one does, so the saving on a quiet host was not measured here. What is left is the per-item loop, which waits for `replayHolds`, and `replayHolds` itself, which S-0608 proposes to keep between calls. The in-flight sharing is not in these figures: it saves the whole of the second request of a burst, which the serve log counted as 73 of its 120 slow `stats.get` calls.

### Cause 18 fixed: the hold replay kept and the days swept once (S-0608)

`stats.get` replayed the hold rules over all history on every call, whatever the window. The sections laid out by the day parsed every transition's time again on each day, and looped over every item even on the days before the first one. A thread's entries were parsed twice, and `costOrder` built its forecaster three times. Now the replay and the touches drift are kept between calls, and the rest is done once a call. The answer is the same, byte for byte.

- The item store counts a generation of each item folder (`flai/internal/workitem/store.go`), and `Repo.Generation` sums it over every folder `List` reads. It moves when the store is told of a change, when this process saves an item, and when a list or get finds an item file added, changed, or removed. It stays the same otherwise.
- `replayHolds` (`flai/internal/metrics/claims.go`) is now a replay into a timeline and a slice of it. The replay runs over every story and all history, and records what was held, by whom and under which code, from each judgement to the next. It does not depend on the window or on now. Each call lays the timeline out over its own window's days and weeks, keeps the stories it reports, and closes the holds still open at its own now. `Compute` keeps the timeline per repository root while the item generation and the manifest's projects and shared paths stay the same. A kept timeline answers a now no earlier than the last event it applied and before any event it left out; any other now replays it again. The touches drift is kept for every story the same way, with the commits' key too, which `storygit.CommittedKeyed` now answers. The replay reads no commit, so a moved branch tip works out the drift again but does not replay the holds. `statsread.Read` gives the root and the two keys, and gives no generation when it moved during its list. With no keys, as in `flai stats`, nothing is kept and the same code runs on a fresh replay.
- `Compute` parses each item's creation and transitions once into a timeline that every section reads (`flai/internal/metrics/metrics.go`, `newTimelines`). Burn-up, the cfd, `claims`, and `costOfDelay` move a cursor over each item's transitions as the days pass (`dayStates`), instead of calling `stateAt` for every item on every day. A day before every item is answered as a zero row, without reading the items. `costOfDelay` adds each item's waiting spans in the same order as before, so the sums round the same. `coordination` still calls `stateAt`, once for each share, not once a day.
- A thread keeps its parse of its entries while its body stays the same, and the copies the kept threads answer share it (`flai/internal/threads/threads.go`). So `waiting` and `digest` read one parse, and so does the next call.
- `planning.Puller` (`flai/internal/planning/forecast.go`) builds the forecaster once and plays each sequence out on fresh lanes. `costOrder` plays its three orders on one.

`TestReportAtEveryWindowIsTheGolden` (`flai/internal/metrics/windows_test.go`) pins the report over `testdata/good` at 1, 7, 30, 90, and 365 days, with commits, projects, and holds of every kind given, as the code wrote it before the change. Several tests in `claims_test.go` cover the kept replay:

- `TestTheHoldReplayIsKeptUntilItsKeysMove` and `TestTheTouchesDriftIsKeptUntilItsKeysMove` show that a second call with the same keys does not replay, and that a moved key does.
- `TestAKeptReplayAtALaterNowIsAFreshReplay` compares a kept timeline sliced an hour, three days, and eight days later with the reference replay, at each of the five windows.
- `TestComputeKeyedGivesTheReportUnkeyed` shows that a keyed call answers what an unkeyed one does.

Measured on 2026-11-09 between 10:15 and 10:17 UTC on the operator's host (WSL2, 24 cores, Go 1.26.0), on the main checkout. One in-process program, built once from each tree, opened the main checkout and turned on the item snapshot as `flai serve` does. Each call then ran `statsread.Read` and `metrics.Compute` as `stats.get` does, with the real runner, and timed the two apart. Before is the branch's base at c61142c7f. After is `story/S-0608` at 71282642c. For each window, a run made one call that is left out, then five timed calls. The two binaries took turns for five rounds. With `Now` fixed, both gave the same report, byte for byte, for every window. The load average over one minute was 22 to 42 at the rounds' starts, and 73 to 76 over fifteen.

The table gives the median of the five rounds' medians, in milliseconds, and the range of the rounds' total medians.

| Window | Before: read | Before: compute | Before: total | Before: range | After: read | After: compute | After: total | After: range |
|--------|-------------:|----------------:|--------------:|---------------|------------:|---------------:|-------------:|--------------|
| 1 day | 51 | 248 | 303 | 258 to 309 | 44 | 30 | 77 | 69 to 107 |
| 7 days | 46 | 241 | 286 | 249 to 320 | 46 | 32 | 82 | 75 to 102 |
| 30 days | 46 | 260 | 306 | 273 to 381 | 51 | 32 | 87 | 76 to 90 |
| 90 days | 51 | 268 | 319 | 283 to 417 | 49 | 35 | 79 | 76 to 88 |
| 365 days | 49 | 364 | 407 | 351 to 450 | 47 | 33 | 79 | 75 to 97 |

A warm 365-day call took 0.91 times a warm 30-day one, against 1.33 times before. Round by round the ratio was 0.86 to 1.24; it passed 1.2 only in the first round. Compute fell from 248 to 364 ms to 30 to 35 ms, at every window. That is more than S-0596 estimated, 200 to 250 ms saved warm, and at 365 days about 290 ms down to about 30 ms. What is left of a warm call is mostly the reads. These are warm calls served from the keep. The first call after an item changes replays the holds once, as every call did before. The host was less busy than in S-0607's measurement, and both binaries ran with the snapshot on, so these figures are not comparable with that table's.

### Cause 19 fixed: flai accept lints each set once with the installed markdownlint (S-0609)

- `scripts/lint-md.sh` runs the `markdownlint-cli2` on `PATH` when it is the pinned version, read from its `package.json` without starting node, and falls back to `npx` otherwise. `scripts/install-tools.sh` installs the pinned version into `.flai-cache/markdownlint` and links it into `bin/`. `flai task done` gains from this alone: it still runs the fix twice.
- `verify.FixRemaining` runs each tier's fix command once and answers what its exit says remains. `flai accept` refuses from that, at preflight over the main checkout's uncommitted files and after the rebase over the branch's files, and runs the check after the fix only when the fix cannot say ([ADR-0265](../adrs/0265-flai-accept-takes-a-tier-s-fix-command-s-exit-as-what-the-fix-left-exit-0.md)). The fixes before the acceptance's commits run only over files new or changed since the preflight fix, by content. So an acceptance runs the markdown tier three times, not five: once for each set, and once over the files it wrote itself, such as the archive.
- The preview reads the main checkout's `git status --porcelain` once and answers every question it asks from it. The preflight lint and the fix before the shared-set commit take the preview's lists. An acceptance reads the main checkout's status twice, not three times: before anything changes, and after the merge and the archive. Resuming an acceptance whose commit failed reads it twice, not five times.

Measured on 2026-11-08 between 03:05 and 03:21 UTC on the operator's host (WSL2, 24 cores, Go 1.26.0). Before is `main` at 3087d9cf4; after is `story/S-0609` at 183ef4f4b. Each was built into `/tmp` and run in its own `git clone --no-local` of the main checkout, so each ran its own tree's `scripts/lint-md.sh`, with a scratch `FLAI_CONFIG`. In each clone a story under E-0033 was started, its task closed, and the story moved to review. Its branch then held a markdown line the fix fixes, `main` moved one commit so the acceptance rebased, and the main checkout held an uncommitted design file with a line the fix fixes. Every run started from a copy of that state. Shims on `PATH` logged each lint call with its duration and each `git status --porcelain` with its checkout. The installed `markdownlint-cli2` 0.20.0 was a scratch install under `/tmp`; `scripts/install-tools.sh` was not run. The three configurations took turns for seven rounds. The host was very busy: the load average rose from 13 to 124 over one minute and was 54 at the end.

| Build | Lint through | Runs | `flai accept` p50 | Range | Lint calls | Time in lint, p50 | Main checkout status reads |
|-------|--------------|-----:|------------------:|-------|-----------:|------------------:|---------------------------:|
| `main` | `npx` | 7 | 16.5 s | 9.8 to 31.7 s | 5 | 11.4 s | 3 |
| S-0609 | installed | 7 | 8.1 s | 4.0 to 12.1 s | 3 | 2.5 s | 2 |
| S-0609 | `npx` | 7 | 16.1 s | 5.5 to 21.2 s | 3 | 8.0 s | 2 |

- Within each round, the acceptance with the installed markdownlint took 2.9 to 22.5 s less than `main`, 7.5 s at the median; six rounds of seven saved more than 4 s; the seventh saved 2.9 s, 11.1 s against 8.1 s.
- In the first round, at a load of 33 as it started, an installed lint call took 0.35 to 0.39 s and an `npx` call 1.1 to 1.5 s: `main` spent 6.7 s in lint, S-0609 1.1 s, and the acceptance fell from 9.8 to 4.0 s.
- Without markdownlint installed, S-0609 still saves two `npx` calls, but under this load the rounds varied more than that: 5.5 s against 9.8 s in the first round, about the same at the p50.
- `flai task done`, timed once per build in each clone's setup without markdownlint installed, ran the fix twice in both, through `npx`, and took 10 to 54 s, varying with the load. Its gain is the installed lint's alone, 0.7 to 1.1 s a call in the first round's figures.
- Unfixable markdown, a heading that skips a level, was refused alike before and after: in the main checkout's uncommitted files before anything changed, and in the branch's files before the merge. Each refusal named the same file, line, and error, and said where to correct it. The two differ only in the tool line the plain tier echoes: `installed` after, and npm's warning before.

### Cause 20: the close-out runs marked tiers at once, and smoke skips what the run passed (S-0610)

- A test tier may carry `parallel: true` ([ADR-0266](../adrs/0266-a-test-tier-marked-parallel-starts-once-the-unmarked-tiers-before-it-have-ended.md)). `verify.Run` (`flai/internal/verify/run.go`) starts a marked tier once every selected tier before it that is not marked has ended, beside the marked ones; an unmarked tier waits for every tier before it. When a tier fails, no tier after it starts, a running tier after it is stopped, and every tier after the first failure in list order is listed `not-reached`. So the verdict, the tier it names, and the findings are a serial run's. With no tier marked the run is serial, as before. `verify.Tiers`, the merge queue's batch, still runs one tier at a time.
- Each tier runs with `FLAI_VERIFY_PASSED`: the tiers it waited for that passed, as `<tier>` or `<tier>:all`, and, in a close-out, `step:check` and the reused tiers. It names only what the tier waited for, so it does not depend on timing.
- `scripts/smoke.sh` skips the template test when the `template` tier passed. It skips the repository check when the close-out's check step passed and the branch changes nothing under `flai/`, since that step runs the running flai's check and `check.sh` runs the tree's. It skips the markdown lint when the `markdown` tier passed whole, or for the paths with no change to the markdownlint configuration or `lint-md.sh`.
- `scripts/flai-build.sh` builds `bin/flai` into a temporary file and renames it, so tiers at once never run a half-written binary.
- This repository marks `flaiover`, `integration`, and `smoke`. `template` and the cheap per-path tiers stay unmarked, so `smoke` starts after `template` and `markdown` and reads their verdicts.

Measured on 2026-11-09 between 09:07 and 10:31 UTC on the operator's host (WSL2, 24 cores, Go 1.26.0), in S-0610's worktree at 2d0bff408 with the `parallel` marks and the `flai-build.sh` change uncommitted. Each run was `flai test --all --json --max 20` with `CLOSE_OUT_STORY=S-0610` exported. That scopes `check.sh` and `lint-md.sh` to the story and skips the integration tier's whole-monorepo tests, as in a close-out. Before is the installed flai 1.86.1, which ignores `parallel` and sets no `FLAI_VERIFY_PASSED`: the serial run, with smoke running every step, as every close-out ran before. After is `bin/flai` built from the same commit, built before the timed run. The Go and node caches were warm for every run. The two took turns for three rounds. The load is the one-minute average at each run's start and end.

| Round | Run | Start | Load | Wall | `flaiover` | `integration` | `smoke` | Cheap tiers |
|------:|-----|-------|------|-----:|-----------:|--------------:|--------:|------------:|
| 1 | before | 09:07 | 52 to 132 | 911 s | 422.0 s | 334.5 s | 137.3 s | 16.7 s |
| 1 | after | 09:34 | 98 to 108 | 547 s | 511.5 s | 317.7 s | 57.2 s | 35.3 s |
| 2 | before | 09:58 | 56 to 70 | 532 s | 190.5 s | 265.1 s | 62.9 s | 13.6 s |
| 2 | after | 10:07 | 73 to 48 | 524 s | 497.1 s | 414.0 s | 36.5 s | 25.9 s |
| 3 | before | 10:16 | 39 to 43 | 429 s | 165.6 s | 208.3 s | 38.8 s | 15.9 s |
| 3 | after | 10:23 | 38 to 52 | 443 s | 427.6 s | 325.8 s | 57.2 s | 15.2 s |

The cheap tiers are `gofmt`, `vet`, `golangci-lint`, `shellcheck`, `markdown`, and `template`, run one after another in both.

- The after run saved 364 s in round 1, 8 s in round 2, and lost 14 s in round 3. Only round 1 met the 150 s the story asked for, and it ran at the highest load, where the serial run's `flaiover` and `smoke` were slowed most.
- The two expensive tiers each fill the host's cores: `go test -race` runs packages in parallel and vitest runs its own workers. Alone at loads of 39 to 70, `flaiover` took 166 to 191 s and `integration` 208 to 265 s; beside each other they took 428 to 497 s and 326 to 414 s. They share the CPU rather than use idle CPU, so their overlap ends about when the serial run would. On a host with idle cores the overlap would save up to the shorter tier; that was not measured here.
- In the after runs `smoke` got `gofmt:all vet:all golangci-lint:all shellcheck:all markdown:all template:all`. It skipped the template test and the markdown lint and ran the repository check, since `--all` has no check step, and the install test. Run alone in round 3 with every step, it took 38.8 s; the lint it skips takes about 9 s, as the `markdown` tier's whole run did, and the template test well under a second. In a close-out the `markdown` tier runs for the changed files, so the lint is skipped only when the branch leaves the markdownlint configuration and `lint-md.sh` alone.
- Three further runs failed on one vitest test, `components.svelte.test.ts` › the props panel and the event log › "lists the calls of the component's function props, and clears them on another pick", which timed out at 5,000 ms. One was an after run at a load of 125, and two were before runs, serial, at 91 to 102; an untimed after run at 08:31, at a load of 86 to 119, failed on it too. It is a page test that posts about 60 messages to the page and touches no flai, so it fails under load whether or not the tiers run at once. A run that failed was run again; none is in the table.
- `--all` runs every cheap tier for the whole checkout, and a close-out selects fewer. So the difference between before and after, not the totals, is what carries over to a close-out that selects `flaiover`, `integration`, and `smoke`.

### Cause 21 fixed: the dashboard asks for open threads and open runs, and the send is timed (S-0611)

#### Fix

- The dashboard asks `threads.list` for the unresolved threads alone (`flaiover/src/lib/server/repo.ts`). It asks for every thread only when a page sends `all=1`: the Threads page with its show-resolved box ticked, and an item's page for that item's threads. The two lists are cached apart, and a thread or item change forgets both.
- `agent.status` takes a parameter, `finished` ([flaiover-dashboard.md § API](flaiover-dashboard.md#api-s-0011)). Without it the answer keeps its keys but holds only these runs:
  - the live and waiting runs;
  - the newest run of each open story and each open planned item;
  - the newest orchestrator, analyzer, and curator runs, and each role's run under way.

  It leaves out the past runs and the runs of finished work. `/api/host-agent`, which the board, the activity page, and a story's agent panel ask every 15 s, asks without it. The planner, orchestrator, analyzer, role, and item-plan routes ask with `finished: true`, and so does a done or cancelled story's agent panel.
- Over the channel, the `request answered` event is logged after the answer is written, with the write's time as `send_ms`: the encoding, the wait for the connection's one write lock, and the write ([Request timing](../../docs/users/flai.md#request-timing-s-0152)). The slow threshold compares `duration_ms` and `send_ms` together.

#### Measured on this repository

The measurement ran on 2026-11-09 between 17:25 and 17:26 UTC, on the operator's host (WSL2, 24 cores, Go 1.26.0).

- Two binaries were built into `/tmp`: one from `main` at 8f934d797 and one from `story/S-0611` at d9f5c6653.
- Each ran `flai hostapi --timing` in the main checkout with the host's own configuration, so `agent.status` read the host's `~/.flai/serve` state.
- "Before" asks as the dashboard asked: `threads.list` with `all: true`, and `agent.status` with nothing. "After" asks as the board now does: `threads.list` with nothing, and `agent.status` with nothing.
- They took turns for three rounds. Each call was a fresh process, so the times are cold.
- The repository held 940 threads, 15 of them unresolved.
- The one-minute load average was 58 to 62.

| Call | Before, bytes | After, bytes | Before, `duration_ms` | After, `duration_ms` |
|------|--------------:|-------------:|----------------------:|---------------------:|
| `threads.list`, the board's | 3,491,017 | 52,753 | 402 to 749 | 280 to 411 |
| `agent.status`, `/api/host-agent`'s | 420,952 | 14,163 | 613 to 759 | 566 to 735 |
| `agent.status` with `finished: true`, the past-runs pages' | — | 420,952 | — | 588 to 727 |

- The board's thread list is 1.5% of its former size and its `agent.status` 3.4%. Both are under the story's targets of 100 KB and 50 KB.
- With `finished: true` the answer is the same size as before. Only the pages that list past runs ask for it, and none of them polls every 15 s.
- The filter that leaves out finished work, `agent.unfinished`, took 3.7 to 9.2 ms a call. That is about 1% of `agent.state`, which is still mostly `agent.activity` (cause 13's reads and the holds' trial merges).
- These one-shot times include the cold reads. What the smaller answers save in flai serve is the encoding, and the time each answer holds the connection's write lock while the next request's answer waits.
- `send_ms` itself was not measured. flai serve on this host runs the installed flai 1.86.4, which does not log it, and it will be logged once a flai release with S-0611 runs there. The serve log's `send_ms` will then show what is left of the wait on the write lock.

## Stories

Each cause has a backlog story under [E-0012](../../wip/kanban/epics/E-0012-performance-analysis-and-improvements.md), with the cause's numbers and a proposed solution. Causes 1 to 4 are the largest share of a board load; 5 multiplies 4 and every write on hosts with a long `PATH`. S-0471 remedies a regression found later, recorded as I-0155, and S-0500 the `git log` that was left in every warm chart answer after it, recorded as I-0184, and S-0531 the one left in every cold answer, recorded as I-0209.

| # | Story |
|---|-------|
| 1 | [S-0156](../../wip/kanban/stories/S-0156-flai-serve-and-flai-mcp-keep-parsed-work-items-between-requests-and-read-again-only-the-files-that-changed.md) flai serve and flai mcp keep parsed work items between requests and read again only the files that changed |
| 2 | [S-0157](../../wip/kanban/stories/S-0157-pending-releases-are-worked-out-without-starting-a-git-process-per-accepted-item.md) Pending releases are worked out without starting a git process per accepted item |
| 3 | [S-0158](../../wip/kanban/stories/S-0158-the-designer-s-inbox-finds-overlapping-touches-without-running-the-whole-check.md) The designer's inbox finds overlapping touches without running the whole check |
| 4 | [S-0159](../../wip/kanban/stories/S-0159-the-dashboard-s-reads-are-answered-in-flai-serve-s-process-not-by-starting-flai.md) The dashboard's reads are answered in flai serve's process, not by starting flai |
| 5 | [S-0160](../../wip/kanban/stories/S-0160-a-flai-process-starts-in-milliseconds-whatever-the-host-s-path.md) A flai process starts in milliseconds whatever the host's PATH |
| 6 | [S-0161](../../wip/kanban/stories/S-0161-the-dashboard-forgets-only-the-answers-a-changed-file-affects-and-gathers-changes-that-arrive-together.md) The dashboard forgets only the answers a changed file affects, and gathers changes that arrive together |
| 7 | [S-0162](../../wip/kanban/stories/S-0162-the-items-and-documents-pages-ask-for-what-they-show-not-the-whole-archive-and-every-file-s-front-matter.md) The items and documents pages ask for what they show, not the whole archive and every file's front matter |
| I-0155 | [S-0471](../../wip/kanban/stories/S-0471-the-chart-data-takes-about-4-s-warm-against-10-ms-measured-on-2026-09-29-almost-all-of-it-in-metrics-replayholds.md) The chart data takes about 4 s warm, against 10 ms measured on 2026-09-29, almost all of it in metrics.replayHolds |
| I-0184 | [S-0500](../../wip/kanban/stories/S-0500-every-stats-get-and-flai-stats-runs-one-git-log-for-claims-drift-about-450-ms-of-a-900-ms-warm-answer.md) Every stats.get and flai stats runs one git log for claims drift, about 450 ms of a 900 ms warm answer |
| I-0209 | [S-0531](../../wip/kanban/stories/S-0531-a-one-shot-flai-stats-and-the-first-stats-get-of-a-flai-serve-process-still-run-one-git-log-for-claims-drift-about-a-quarter-of-a-cold-answer.md) A one-shot flai stats, and the first stats.get of a flai serve process, still run one git log for claims drift, about a quarter of a cold answer |
