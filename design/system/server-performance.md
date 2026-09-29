---
title: Server-side response times
updated: 2026-09-29
status: active
topics: [server-side, back-end, cli, dashboard]
---

# Server-side response times

The finding of S-0152, for E-0012. The board and other pages pause for long moments while they load. This document says how the time a request takes is now measured beneath its transport, what was measured on this repository, where the time goes, and what would take it away. Each cause has a story under E-0012.

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
| `push.pending` | unpushed notice, review | 157 ms | 66 B | `exec.flai.push` all of it |
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
| `/api/unpushed` | 58 | 161 ms | 172 ms | `push.pending` 157 ms |
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

The pauses the operator sees come from these adding up. A board load asks `board.get` (190 ms), `publish.preview` (235 ms), and `push.pending` (157 ms), and `/api/projects` asks `board.get`, `inbox.designer`, and `agent.status` again (about 360 ms). Each of them reads every work item from disk, and the three that start flai also pay the 145 ms start. They run at once, on a host where agents keep asking the same MCP tools, so they compete for the disk and the CPU.

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

### Cause 5: a flai process starts without searching PATH (S-0160)

flai asks its questions through a prompt package of its own and no longer depends on `huh`, so nothing searches `PATH` when flai starts ([ADR-0052](../adrs/0052-flai-asks-its-questions-a-line-at-a-time-with-its-own-prompt-package-not-with.md), [go-libraries.md](../tech/go-libraries.md#prompts-s-0160)). On 2026-09-29, on the same host with its own `PATH` of 54 entries, 17 under `/mnt`, `flai version` took 5.5 to 8.3 ms over 30 runs, median 5.8 ms, against 160 to 178 ms, median 163 ms, for `main` built the same hour. `GODEBUG=inittrace=1 flai version` shows no package init over 0.6 ms, where `atotto/clipboard` took 171 ms.

### Cause 6: the dashboard forgets and asks only for what a change affects (S-0161)

A change reaches the dashboard's server with the path's kind, read against the manifest's layout: a work item (`board.md` and the archive included), a narrative, a thread, an ADR, another document, the manifest, or anything else. The server forgets only the answers flai reads from that kind ([flaiover-dashboard.md](flaiover-dashboard.md#where-the-data-comes-from-s-0073)), so a narrative log line keeps `board.get`, the items, the statistics, and the ADRs, and forgets the inbox, the activity, and the documentation tree. Pages ask again only for the kinds they show, once for the changes that arrive within 500 ms of each other, and at least every 2 s during a burst that does not settle. The board asks for `board.get`, `/api/publish`, and `/api/unpushed` again only when a work item or the manifest changes, and for the agents when a work item or a thread changes or flai serve says an agent started or ended. Behaviour tests pin which paths forget which answers (`repo-forget.test.ts`) and what the board asks for each kind (`routes/board/board.svelte.test.ts`). The request rate was not measured again in the running dashboard: that needs an image built from this story in place of the operator's container.

## Stories

Each cause has a backlog story under [E-0012](../../wip/kanban/epics/E-0012-performance-analysis-and-improvements.md), with the cause's numbers and a proposed solution. Causes 1 to 4 are the largest share of a board load; 5 multiplies 4 and every write on hosts with a long `PATH`.

| # | Story |
|---|-------|
| 1 | [S-0156](../../wip/kanban/stories/S-0156-flai-serve-and-flai-mcp-keep-parsed-work-items-between-requests-and-read-again-only-the-files-that-changed.md) flai serve and flai mcp keep parsed work items between requests and read again only the files that changed |
| 2 | [S-0157](../../wip/kanban/stories/S-0157-pending-releases-are-worked-out-without-starting-a-git-process-per-accepted-item.md) Pending releases are worked out without starting a git process per accepted item |
| 3 | [S-0158](../../wip/kanban/stories/S-0158-the-designer-s-inbox-finds-overlapping-touches-without-running-the-whole-check.md) The designer's inbox finds overlapping touches without running the whole check |
| 4 | [S-0159](../../wip/kanban/stories/S-0159-the-dashboard-s-reads-are-answered-in-flai-serve-s-process-not-by-starting-flai.md) The dashboard's reads are answered in flai serve's process, not by starting flai |
| 5 | [S-0160](../../wip/kanban/stories/S-0160-a-flai-process-starts-in-milliseconds-whatever-the-host-s-path.md) A flai process starts in milliseconds whatever the host's PATH |
| 6 | [S-0161](../../wip/kanban/stories/S-0161-the-dashboard-forgets-only-the-answers-a-changed-file-affects-and-gathers-changes-that-arrive-together.md) The dashboard forgets only the answers a changed file affects, and gathers changes that arrive together |
| 7 | [S-0162](../../wip/kanban/stories/S-0162-the-items-and-documents-pages-ask-for-what-they-show-not-the-whole-archive-and-every-file-s-front-matter.md) The items and documents pages ask for what they show, not the whole archive and every file's front matter |
