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

### SQLite

T-1801 ran on 2026-10-16 between 06:13Z and 06:35Z, on the same host, with a load average of 0.6 to 1.8. Its program, outputs, and commands are in `bin/experiments/S-0427/sqlite/`, with build and run steps in the folder's `README.txt`.

**The corpus.** The main checkout kept changing, so T-1801 copied it once with `cp -a`, which keeps modification times, to `/tmp/s0427-corpus`, at 06:13:39Z. It copied `system-flow.yaml`, `wip/kanban`, `wip/archive`, `wip/threads`, `wip/messages`, the three activity documents under `wip/agents`, `design/issues`, and `.flai-cache/conflicts.jsonl`. Every load, query, and check below reads that copy. Its items number 2,308, against 2,284 at T-1800's baseline. The copy is not a git checkout, so `storygit.Committed` fails on it. That leaves `claims.drift` and its `git log` out on both sides of every comparison.

Every query and every flai figure it is checked against uses a fixed now, `2026-10-16T06:04:02Z`, the reference's `generated_at`, and the reference's window from `2026-09-16T06:04:02Z`. `metrics.Compute` takes `Options.Now`, so the program sets it rather than read the clock.

**The driver.** `modernc.org/sqlite` v1.60.1, which bundles SQLite 3.53.4. It is SQLite's C translated to Go, so it needs no cgo. flai builds with `CGO_ENABLED=0` and cross-compiles in GoReleaser, and a cgo driver such as `mattn/go-sqlite3` would end both. `ncruces/go-sqlite3`, SQLite compiled to Wasm under wazero, also needs no cgo; T-1801 did not try it. Downloading the driver and `modernc.org/libc` took 11 minutes on this host's connection.

**The program.** It is in package `main` under `flai/internal/s0427sqlite`, mapped there by `go build -overlay`, so it parses with flai's own readers: `workitem.ReadItem`, `threads.Read`, and `messages.Read`. A second overlay file adds `metrics.S0427Parts` to package `metrics`, to time Compute's parts one at a time. `-modfile` points the build at a copy of `flai/go.mod` with the driver added. No file is written into the module, and `flai/go.mod` gains no dependency. In that copy only, `go get` also raised `golang.org/x/sync` from v0.20.0 to v0.23.0.

```bash
cd flai && GOMODCACHE=$S/gomodcache GOPROXY=file:///home/alex/go/pkg/mod/cache/download GOFLAGS=-modcacherw \
  go build -modfile=$S/sqlite/go.mod -overlay=$S/sqlite/overlay.json -o $S/sqlite/s0427sqlite ./internal/s0427sqlite
./s0427sqlite load /tmp/s0427-corpus a.db
./s0427sqlite query /tmp/s0427-corpus a.db 2026-10-16T06:04:02Z answers.json 0 tables
```

**The load.** It lists the six item folders under `wip/kanban` and `wip/archive/kanban`, then `TH-*.md` under `wip/threads` and `MS-*.md` under `wip/messages`, as flai lists them. That is 2,914 files: 2,308 items, 479 threads, and 127 conversations. It leaves out the archive's narratives and the READMEs, which no chart reads. In sorted path order and in one transaction, it parses each file and writes its rows. It then fills the derived tables and runs `ANALYZE`.

```sql
CREATE TABLE files (path TEXT PRIMARY KEY, kind TEXT, mtime_ns INTEGER, size INTEGER, ino INTEGER, trusted INTEGER) WITHOUT ROWID;
CREATE TABLE items (path TEXT PRIMARY KEY, id TEXT, type TEXT, nature TEXT, title TEXT, status TEXT, parent TEXT,
  archived INTEGER, created INTEGER, estimate TEXT, doc TEXT) WITHOUT ROWID;    -- doc: the parsed item as JSON
CREATE TABLE transitions (path TEXT, seq INTEGER, to_state TEXT, at INTEGER, by TEXT, PRIMARY KEY (path, seq)) WITHOUT ROWID;
CREATE TABLE blocks (path TEXT, seq INTEGER, from_at INTEGER, until_at INTEGER, reason TEXT, PRIMARY KEY (path, seq)) WITHOUT ROWID;
CREATE TABLE touches (path TEXT, seq INTEGER, touch TEXT, PRIMARY KEY (path, seq)) WITHOUT ROWID;
CREATE TABLE usage (path TEXT PRIMARY KEY, source TEXT, seconds INTEGER, estimated INTEGER, priced_by TEXT,
  empty_wakes INTEGER, n_models INTEGER, n_strategic INTEGER) WITHOUT ROWID;
CREATE TABLE usage_models (path TEXT, seq INTEGER, model TEXT, input INTEGER, output INTEGER, cache_read INTEGER,
  cache_write INTEGER, cost REAL, PRIMARY KEY (path, seq)) WITHOUT ROWID;
CREATE TABLE threads (path TEXT PRIMARY KEY, id TEXT, status TEXT, anchor TEXT, created INTEGER, updated INTEGER, doc TEXT, body TEXT) WITHOUT ROWID;
CREATE TABLE conversations (path TEXT PRIMARY KEY, id TEXT, from_story TEXT, to_story TEXT, status TEXT,
  created INTEGER, updated INTEGER, doc TEXT, body TEXT) WITHOUT ROWID;
-- derived, refreshed for a path whenever its file is loaded
CREATE TABLE item_facts (path TEXT PRIMARY KEY, id TEXT, type TEXT, nature TEXT, status TEXT, parent TEXT,
  created INTEGER, committed INTEGER, started INTEGER, first_done INTEGER, completed INTEGER) WITHOUT ROWID;
CREATE TABLE segments (path TEXT, seq INTEGER, state TEXT, from_at INTEGER, to_at INTEGER, PRIMARY KEY (path, seq)) WITHOUT ROWID;
```

Times are Unix seconds. Every table is keyed by the file's path, so reloading a file deletes its rows by path and writes them again. `item_facts` holds each item's first entry into ready, in progress, and done, and when it closed: its last transition while it is done or cancelled. `segments` holds each state an item was in, from when to when; `seq` -1 is backlog from creation, and `to_at` is NULL for the state it is in now.

```sql
INSERT INTO segments
  SELECT i.path, -1, 'backlog', i.created, (SELECT t.at FROM transitions t WHERE t.path = i.path ORDER BY t.seq LIMIT 1)
  FROM items i WHERE i.path = :path
  UNION ALL
  SELECT t.path, t.seq, t.to_state, t.at, LEAD(t.at) OVER (PARTITION BY t.path ORDER BY t.seq)
  FROM transitions t WHERE t.path = :path;
```

**Keeping it current.** A rescan works as `flai/internal/workitem/store.go` does. It lists the same folders, stats each file, and reloads a file only when its inode, modification time, or size has changed. It also reloads a file it read within 2 s of that file's modification time, as `racyWindow` does. A file that is gone has its rows deleted. Each one-shot query and each warm run starts with a rescan, so every answer comes from a current model. The rescan figures come from a second copy, `/tmp/s0427-corpus-rescan`, because the test appends a line to the body of one story before each changed run (`rescan-test.py`).

**Two loads from empty give the same rows.** A sorted dump of every table, `sqlite_stat1` included, hashed to `96612a61a2a4fc36baabc4f338b6409c020bf79b7ed2ab3b0c81559d954e3e06` for each of three loads. The three database files were byte for byte the same too, at SHA-256 `7433ca3bd69164593a8fc6bdae5bd6cc89fc8b16a878e5343b813e9c3b271bf3`.

**The queries.** Eight reads, each one prepared statement, are written twice. The first form computes everything at query time over views. The second form reads the derived tables. Both are in `queries.go`. The key ones in the second form:

```sql
-- time in state per story; an open story's state runs to now
SELECT f.id, s.state, SUM(COALESCE(s.to_at, :now) - s.from_at)
FROM item_facts f JOIN segments s ON s.path = f.path
WHERE f.type = :type AND (s.to_at IS NOT NULL OR f.status NOT IN ('done', 'cancelled'))
GROUP BY f.path, s.state;

-- each story's state at the end of each day: the CFD, the burn-up overall and per epic, claims.days[].in_progress
WITH RECURSIVE
it AS (SELECT path, parent, created FROM item_facts WHERE type = :type),
bounds AS (SELECT unixepoch(date(MAX(MIN(created), :start), 'unixepoch')) AS first, unixepoch(date(:now, 'unixepoch')) AS last FROM it),
days(d) AS (SELECT first FROM bounds UNION ALL SELECT d + 86400 FROM days, bounds WHERE d + 86400 <= last),
st AS (SELECT days.d, it.path, s.state FROM days JOIN it ON it.created <= days.d + 86399
       JOIN segments s ON s.path = it.path AND s.from_at <= days.d + 86399 AND (s.to_at IS NULL OR s.to_at > days.d + 86399))
SELECT date(d, 'unixepoch'), state, count(*) FROM st GROUP BY d, state;   -- the burn-up groups st by parent the same way

-- cycle, lead, and queue time distributions over the window, nearest rank as metrics.percentile
SELECT metric, n, MAX(CASE WHEN rn = MAX((50 * n + 99) / 100, 1) THEN x END) AS p50,
       MAX(CASE WHEN rn = MAX((85 * n + 99) / 100, 1) THEN x END) AS p85, MAX(x), AVG(x)
FROM (SELECT metric, x, ROW_NUMBER() OVER (PARTITION BY metric ORDER BY x) AS rn, COUNT(*) OVER (PARTITION BY metric) AS n FROM v)
GROUP BY metric;

-- throughput per ISO week, by nature, empty weeks included
WITH RECURSIVE weeks(m) AS (SELECT unixepoch(date(:start, 'unixepoch', '-6 days', 'weekday 1'))
                            UNION ALL SELECT m + 604800 FROM weeks WHERE m + 604800 <= :now)
SELECT strftime('%G-W%V', w.m, 'unixepoch'), f.nature, count(f.id) FROM weeks w LEFT JOIN item_facts f
  ON f.type = :type AND f.first_done BETWEEN :start AND :now AND f.first_done >= w.m AND f.first_done < w.m + 604800
GROUP BY w.m, f.nature;
```

Reading a state from `segments` gives `metrics.stateAt`'s answer only when each item's transitions are in time order. In this corpus no transition is earlier than the one before it, and the view form, which cuts at the first transition after the day's end as `stateAt` does, gave the same answers. The other reads: items per column now; per story, lead, cycle, queue, and blocked time; throughput per day; the summary's counts, WIP, flow efficiency, and share of lead time per state. The last is what agents spent per model over time: per type and day, from the first day an item with usage was done to today, the items, tokens, cost, and seconds, overall and per model, with the running means.

**The hold replay has no SQL form here.** `metrics.replayHolds` replays `Holds.Of` for every ready story at each moment a story or task was created or moved. In this corpus that is up to 5,951 moments, over as many as 2,283 stories and tasks. The `after` hold could be written in SQL: the state, at that moment, of the stories a story names. The overlap hold is a program, not a query. It reads each story's claim through the manifest's sub-project names and tags. It replaces a story touch with its tasks' narrower touches, which depends on which tasks are open at that moment. It compares paths by prefix and leaves out overlaps inside the manifest's shared patterns and inside conversations' shares. The replay also uses today's touches for every past moment, so a stored result goes stale for the whole history whenever any story's or task's touches change. T-1801 left the replay, and the held figures in `claims` and `items[].held_seconds`, in Go.

**Checks.** `compare.py` checks each SQL answer against a report, value by value. Integers must be equal. Floats may differ by a relative 1e-9, because SQL sums in a different order from Go. The first report is `metrics.Compute` over the same snapshot at the same now (`goref`). The second is T-1800's `reference-flai-stats.json`.

| Read | Report fields checked | Values | Differ from flai over the snapshot | Differ from the reference |
|------|-----------------------|--------|------------------------------------|---------------------------|
| `columns` | `items[].status`, counted | 4 | 0 | 2 |
| `cycle_time` | `items[]`: lead, cycle, queue, blocked seconds | 1,740 | 0 | 3 |
| `time_in_state` | `items[].time_in_state_seconds` | 1,587 | 0 | 46 |
| `summary` | `summary`: completed, cancelled, WIP, three distributions, flow efficiency, state shares | 23 | 0 | 1 |
| `throughput_week` | `throughput[]` with `by_nature` | 32 | 0 | 0 |
| `throughput_day` | no field of its own; its sum against `throughput[]` | 1 | 0 | 0 |
| `daily_states` | `cfd[]`, `burnup` (26 groups), `claims.days[].in_progress` | 1,151 | 0 | 4 |
| `spend_models` | `usage.spend.{epic,story,task}`: buckets, models per bucket, the window's models | 783 | 0 | 11 |
| Total | | 5,321 | 0 | 67 |

The SQL answers, in both forms, match flai's over the snapshot on all 5,321 values. Of the 67 that differ from the reference:

- 45 differ by 0.468341669 s. Each is a time measured to now: an open item's time in state or an open block. The reference's now was 06:04:02.468341669, and `generated_at` drops the fraction; the snapshot's now is the whole second.
- 11 are S-0413, accepted at 06:07:38Z, after the reference: done 379 against 378, in progress 2 against 3, and its lead, cycle, and time in state.
- 11 are task spend on 2026-10-16. The acceptance rewrote the usage of S-0413's tasks at 06:06:54Z, so T-1750 and T-1755 carry 907 agent seconds against the reference's 835.

**Reads replaced, compute kept.** The `passthru` command rebuilds what `metrics.Compute` takes. Items, threads, and conversations come from the `doc` and `body` columns. The activities, issues, board, manifest, and conflict log come from their files, as `statsread.Read` reads them. It times that against `statsread.Read` in the same process, then runs `metrics.Compute` on both. The two reports were byte for byte the same: 648,492 bytes, SHA-256 `5a83742a…`. On the snapshot, `statsread.Read`'s `git log` fails at once, so neither side includes the 370 ms it takes in the main checkout.

| Inputs of `metrics.Compute` | Cold: first read in the process | Warm: median of three | Measured with |
|-----------------------------|---------------------------------|-----------------------|---------------|
| `statsread.Read` from the files | 550 ms | 71.7 ms | `passthru`, `time.Now` around each read |
| From SQLite, with the rest from files | 110 ms, 84 ms of it the items' rows and their JSON | 98.8 ms, 78.7 ms of it the items | the same run |

**Timings.** Medians of three, each with its range. Wall time, CPU, and max RSS of a fresh process come from `measure.py`. In-process times come from `time.Now` around each step. "flai" in the first table is `metrics.S0427Parts`: each part of `Compute`, warm, from items already parsed, over the snapshot (`goref`, three runs). Over the same snapshot, `metrics.Compute` took 3,513 ms, median of three.

| Read | flai: Compute's part, warm | SQLite over views, warm | SQLite over derived tables, warm | SQLite over derived tables, one-shot, first run |
|------|----------------------------|-------------------------|----------------------------------|-------------------------------------------------|
| `columns` | within `items` | 0.32 ms | 0.29 ms | 0.28 ms |
| `cycle_time` | `items`, `Derive` for every story: 1.6 ms | 14.6 ms | 1.26 ms | 1.46 ms |
| `time_in_state` | within `items` | 132 ms | 4.4 ms | 5.3 ms |
| `summary` | 0.21 ms | 145 ms | 10.0 ms | 10.0 ms |
| `throughput_week` | 0.16 ms | 39.1 ms | 2.1 ms | 2.1 ms |
| `throughput_day` | none | 2.1 ms | 0.56 ms | 0.55 ms |
| `daily_states` | `cfd` 6.3 ms, `burnup` 9.2 ms, claims in progress 5.4 ms | 134 ms | 43.6 ms | 44.2 ms |
| `spend_models` | `spendOverTime` 1.8 ms | 70.0 ms | 20.6 ms | 21.6 ms |
| All eight | 24.5 ms, summed (derived) | 540 ms (538 to 543) | 82.8 ms (81.9 to 83.5) | 85.7 ms (85.2 to 86.4) |
| Hold replay | `replayHolds` 3,404 ms | no SQL form | no SQL form | no SQL form |

| What | Baseline | SQLite | Success-measure row |
|------|----------|--------|---------------------|
| Full load from empty, fresh process | `stats.get` cold, 4715 ms | 1088 ms wall (1083 to 1117); in process 1080 ms, of which inserts 920 ms and derived tables 65 ms; user 1036 ms, system 272 ms; max RSS 26,264 KB | A full load: met |
| Rescan, nothing changed | none | 19.6 ms (19.2 to 19.9) in process, nearly all of it reading the `files` table and stat-ing 2,914 files; 16.5 to 17.2 ms in the warm processes | Keeping current |
| Rescan, one story changed | none | 27.5 ms (25.3 to 28.3), of which 7.8 ms reloads the file and its derived rows | Keeping current, 50 ms per changed file: met |
| One-shot: open, rescan, the eight reads, write the answers | `flai stats --json`, 4701 ms | Derived tables: 119 ms wall (118 to 120), user 97 ms, system 28 ms, max RSS 24,456 KB. Views: 577 ms wall (573 to 578), max RSS 28,924 KB | Agent's one-shot: met for these eight reads only |
| Warm: the eight reads with prepared statements | `stats.get` warm, 4022 ms | Derived tables: 82.8 ms, user 97 ms and system 7 ms a run, max RSS 27,396 KB. Views: 540 ms, user 351 ms and system 219 ms a run, max RSS 31,316 KB | Operator's wait: met for these eight reads only |
| One-shot of the whole report, inputs from SQLite | `flai stats --json`, 4701 ms | about 4.0 s, derived: 110 ms of inputs, 3,513 ms of `Compute`, and the 370 ms `git log` | Agent's one-shot: not met |
| Database file | corpus 8.95 MB | 23,531,520 bytes. Without the pass-through `doc` and `body` columns, a vacuumed copy is 8,609,792 bytes | Storage: T-1805 |

The eight SQL reads answer in 83 ms warm over the derived tables, with every value matching flai's. But flai's Go computes the same parts from parsed items in about 25 ms. The hold replay, at 3.4 s, is 97 percent of `Compute` over the snapshot, and no read model takes it over. Inputs read from SQLite save about 440 ms on a cold call and cost about 27 ms on a warm one.

## Results

Later tasks of S-0427 write the results, against the success measure above.

## Recommendation

Provisional, from the baseline alone: adapt. A read model that only replaces the file reads saves at most about 80 to 670 ms of a 4 s answer, so it is worth having only if its queries also take over the computation. The SQLite figures point the same way: the chart reads in SQL take 83 ms where flai's Go takes about 25 ms over parsed items, and the 3.4 s hold replay has no SQL form, so the replay, not the reads, is what to change first. Later tasks of S-0427 settle it once the SQLite and Postgres comparisons are in.
