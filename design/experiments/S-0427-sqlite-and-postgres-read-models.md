---
title: A SQLite or Postgres read model loaded from the files, for the chart data's broad reads
updated: 2026-10-16
status: active
story: S-0427
---

# A SQLite or Postgres read model loaded from the files, for the chart data's broad reads

The results of S-0427, an `experiment` story run on 2026-10-16. flai answers the charts by reading every work item, thread, conversation, and issue from their files and computing the metrics of [metrics.md](../system/metrics.md) over them. This experiment asks whether a database loaded from those files, and kept in step with them, answers the same reads faster by enough to pay for itself. The files stay the source of truth; the database would only be a read model.

## Hypothesis

A read model, a SQLite database loaded deterministically from the files and updated as they change, answers the chart data's broad reads faster than flai's file reads do. The margin is worth the storage, RAM, and CPU the database adds, both for the operator waiting on the charts page of a long-running `flai serve` and for an agent's one-shot `flai stats --json`. Postgres would add a running server; it is worth it only if it answers the same queries clearly faster than SQLite.

The baseline below qualifies this. Most of a chart answer's time today is computation, not reading files. A read model saves time only if its queries take over that computation, not just the reads.

## Success measure

Each figure is set against the baseline in [What was done](#what-was-done), on the same corpus, host, and 30-day window. An answer counts only if it matches the reference `flai stats --json` on every field except those that depend on the current time.

| Measure | Baseline | Threshold to confirm |
|---------|----------|----------------------|
| Operator's wait: the chart data from a warm read model in a long-running process | `stats.get` warm, 4022 ms | 400 ms or less, a tenth of the baseline |
| Keeping the read model current after one work item changes | none | 50 ms or less per changed file |
| Agent's one-shot: a fresh process that opens an existing read model, checks it is current, and answers | `flai stats --json`, 4701 ms | 1200 ms or less, a quarter of the baseline |
| A full load of the read model from the files, when there is none | `stats.get` cold, 4715 ms | no more than the cold baseline |
| Storage: the database file | corpus 8.95 MB | 18 MB or less, twice the corpus |
| RAM: max RSS added to the long-running process | 71 MB for six `stats.get` calls | 50 MB or less added |
| CPU while idle, once loaded | none | none: no polling beyond flai's existing file watch |
| Postgres over SQLite | the SQLite figures | worth its server only if it answers the same queries at least twice as fast |

## What was done

T-1800 measured the baseline on 2026-10-16 between 06:02Z and 06:06Z. The corpus is the main checkout, `/home/alex/git/bytepunx/system-flow`, read only. The flai is the one built from the story's worktree by `scripts/flai.sh`: `flai dev`, commit `4282b31c`, Go 1.26.0, git 2.47.3. The host is WSL2 with 24 cores; the load average was 1.1 to 3.0, with other agents working in the same checkout.

Every command ran from the main checkout with `FLAI_CONFIG` and `FLAI_CACHE_DIR` set to its `.flai-cache` as `scripts/env.sh` sets them. `/usr/bin/time` is not installed on this host, so a short Python script stood in for `time -v`. It starts each run as a fresh process and reads wall time, user and system CPU, and max RSS from the child's `wait4` rusage: `python3 -I measure.py <prefix> <runs> -- <command>`. The script, the outputs, and the exact commands are in the git-ignored scratch folder `bin/experiments/S-0427/` of the story's worktree, with a `README.txt`.

### The corpus

Counted at 06:02:39Z with `find <dir> -type f -name '*.md'` and the sum of `-printf '%s'`.

| Folder | Markdown files | Bytes |
|--------|----------------|-------|
| `wip/kanban` | 260 (5 epics, 45 stories, 209 tasks, the board) | 602,575 |
| `wip/archive` | 2,390 (20 epics, 390 stories, 1,598 tasks, 381 narratives, a README) | 6,517,808 |
| `wip/threads` | 479 | 1,464,948 |
| `wip/messages` | 128 | 366,931 |
| Total | 3,257 | 8,952,262 |

The checkout changed while it was measured: `statsread.Read` found 2,283 work items at 06:05Z and 2,284 a few seconds later. The chart data reads more than these folders: the strategic agents' activity documents, the issues under `design/issues`, the board, the manifest, the conflict log in the cache, and one `git log --name-only` over `main` and every story branch, for claims drift.

### The window

The stats cover the 30 days before the call, flai's default and the charts page's default (`flaiover/src/lib/chartwindow.svelte.ts`), for stories, laid out by day. The reference `flai stats --json` was generated at `2026-10-16T06:04:02Z` with `window_start` `2026-09-16T06:04:02Z`. Later tasks query that window. It counts 369 stories completed, 12 cancelled, and 3 in progress, and its report has 435 entries under `items`.

flai has no way to fix its clock, so two runs differ in the fields measured from now: `generated_at`, `window_start`, the ages and times in state under `items` and `aging`, the times under `cost_of_delay.order`, and `coordination.weeks[].shares.held_seconds`. Three runs of `flai stats --json` agreed on every other field.

### How each figure was taken

- Cold `stats.get` is a fresh process per call: `flai hostapi --timing stats.get '{"since":"30d","type":"story","bucket":"day"}'`, the parameters the charts page sends (`/api/stats?since=30d&type=story&bucket=day`). `--timing` logs one `request answered` event with its phases, as `flai serve` times a request ([server-performance.md](../system/server-performance.md#how-it-is-measured)).
- Warm `stats.get` cannot be measured with `flai hostapi`, which answers one call and exits. A throw-away program called the same method table six times in one process, the first call cold and the rest warm, as `flai serve` keeps parsed items between requests (S-0156). It used the same recorder as `--timing`, and read per-call CPU from `getrusage`. It was built with `go build -overlay`, which mapped it into `flai/internal/` without writing a file into the module. The same program then timed `statsread.Read` and `metrics.Compute` apart, three times warm, under a CPU profile.
- The agent's one-shot is `flai stats --json`, a fresh process per call.

`stats.get` marks no `repo.list` phase: it lists the items inside `statsread.Read`, within `stats.compute`. The phases it marks are `repo.open`, `stats.compute`, the git processes it starts (`exec.git.log`, `exec.git.rev-parse`), and `encode`.

### Baseline

Medians of three runs, except warm `stats.get`, the median of the five warm calls in one process. The answer of `stats.get` is about 1.08 MB of compact JSON; `flai stats --json` prints the same report indented, 1.61 MB.

| What | Wall | In-process | `stats.compute` | `exec.git.log` | `encode` | User CPU | System CPU | Max RSS | Measured with |
|------|------|------------|-----------------|----------------|----------|----------|------------|---------|---------------|
| `stats.get` cold | 4715 ms (4701 to 4801) | 4690 ms | 4671 ms | 376 ms | 6.2 ms | 6356 ms | 555 ms | 79,368 KB | `measure.py ... 3 -- flai hostapi --timing stats.get '{"since":"30d","type":"story","bucket":"day"}'` |
| `stats.get` warm, calls 2 to 6 | | 4022 ms (4011 to 4131) | 4006 ms | 371 ms | 6.8 ms | 5208 ms a call | 294 ms a call | 71,304 KB after six calls | `measure.py ... 1 -- warm <main-checkout> 6 '<same params>'` |
| `stats.get`, first call of that process | | 4558 ms | 4542 ms | 377 ms | 6.0 ms | 5818 ms | 398 ms | 65,832 KB | the same run |
| `flai stats --json` | 4701 ms (4654 to 4803) | | | | | 6307 ms | 611 ms | 79,784 KB | `measure.py ... 3 -- flai stats --json` |

`repo.open` took about 1 ms and `exec.git.rev-parse` about 2 ms in every call.

### Where a warm answer's time goes

| Half of `stats.compute`, warm | Median of three | Range | Measured with |
|-------------------------------|-----------------|-------|---------------|
| `statsread.Read`: items, activities, threads, issues, board, git log, conversations, conflicts | 450 ms | 439 to 459 ms | the warm program's split, `time.Now` around each call |
| `metrics.Compute` | 3612 ms | 3553 to 3712 ms | the same |

The git log is about 370 ms of the read, by the phases of the same process, which leaves about 80 ms for reading the files warm. A cold call took 4558 ms in the warm program and 4690 ms in a fresh process, against 4022 ms warm: the parsing that kept items save is about 540 to 670 ms.

The CPU profile of the three splits (`go tool pprof -top -cum warm/warm warm/compute.pprof`) puts 67 percent of the samples under `metrics.replayHolds`. It replays the hold rules at each creation and transition of a story or task, through `workitem.Holds.Claim`, `Holds.Overlaps`, and `PathsOverlap`. The profile's largest leaf costs are string concatenation and allocation, and the garbage collector takes about a fifth of all samples, which is why user CPU exceeds wall time.

On 2026-09-29, [server-performance.md](../system/server-performance.md#cause-4-the-dashboards-reads-answered-in-flai-serves-process-s-0159) measured `stats.get` on this repository at 125 ms cold and 9.1 to 10.6 ms warm. Today it takes 4.7 s cold and 4.0 s warm, almost all of it in computing.

## Results

Later tasks of S-0427 write the results, against the success measure above.

## Recommendation

Provisional, from the baseline alone: adapt. A read model that only replaces the file reads saves at most about 80 to 670 ms of a 4 s answer, so it is worth having only if its queries also take over the computation. Later tasks of S-0427 settle it once the SQLite and Postgres comparisons are in.
