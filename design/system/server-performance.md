---
title: Server-side response times
updated: 2026-10-28
status: active
topics: [server-side, back-end, cli, dashboard]
---

# Server-side response times

The finding of S-0152, for E-0012. The board and other pages pause for long moments while they load. This document says how the time a request takes is now measured beneath its transport, what was measured on this repository, where the time goes, and what would take it away. Each cause has a story under E-0012.

> Since S-0195 ([ADR-0067](../adrs/0067-accepted-work-reaches-the-remote-only-when-it-is-published-and-agents-publish.md)) the dashboard has no unpushed notice: the method `push.pending` and the route `/api/unpushed` are gone, and a board load no longer asks for them. The measurements below keep them as they were measured.

## How it is measured

`flai serve` and `flai mcp` log one `request answered` event per request, timed inside flai from when the request was read to when its answer is ready to write, with the phases the time went to ([flai-cli.md](flai-cli.md#internal-structure), [Request timing](../../docs/users/flai.md#request-timing-s-0152)). The dashboard logs its own duration for each HTTP request. The difference is the transport: the channel's WebSocket and the dashboard's own handling. `flai hostapi --timing` times one method with no transport at all, and `FLAI_PPROF_ADDR` offers Go's profiles from `flai serve`.

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
