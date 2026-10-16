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

### Overhead

T-1802 ran on 2026-10-16 between 06:40Z and 06:58Z, on the same host, with a load average of 0.7 to 4.1. Its programs, scripts, outputs, and commands are in `bin/experiments/S-0427/overhead/`, with the steps appended to the folder's `README.txt`. `overhead.go` adds commands to T-1801's program through the same overlay; T-1801's files are unchanged. Every query uses T-1801's derived tables and its fixed now. `/tmp` is tmpfs, so the corpora are in RAM. The databases are on the worktree's ext4 disk unless a row says tmpfs.

**The corpora.** `make-corpus.py` copies the snapshot whole, then adds N-1 copies of every work item, thread, and conversation. Copy k adds k × 10000 to every four-digit `E-`, `S-`, `T-`, `TH-`, and `MS-` ID, in the file and in its name, so S-0241 becomes S-10241, S-20241, and so on. Parents, `after`, anchors, and a conversation's stories resolve within each copy. Issues, the board, the manifest, the activity documents, and the conflict log stay once. Dates inside the files are not shifted. Each copy repeats the snapshot's history on the same days, so at any moment there are N times as many open stories, all claiming the same paths. That is a heavier case for the hold replay than ten times the history spread over more time.

| Corpus | Files loaded | Items | Bytes of those files | Measured with |
|--------|--------------|-------|----------------------|---------------|
| 1x, `/tmp/s0427-corpus` | 2,914 | 2,308 | 7,432,477 | `corpus-bytes.py <corpus>...` |
| 2x, `/tmp/s0427-corpus-2x` | 5,828 | 4,616 | 14,906,975 | the same |
| 3x, `/tmp/s0427-corpus-3x` | 8,742 | 6,924 | 22,381,473 | the same |
| 10x, `/tmp/s0427-corpus-10x` | 29,140 | 23,080 | 74,702,959 | the same |

The SQL answers matched flai's report over the same corpus on every value at every scale: 5,321, 9,517, 13,713, and 43,085 values, none differing (`compare.py runs/answers-<x>x.json runs/baseline-<x>x.json`).

**Storage.** The database grows in step with the items, at about 10.2 MB per thousand items, or 3.7 MB without the pass-through columns. It is 3.2 times the bytes it loads, and 1.16 times without them. At 1x, the `items` table holds 11.9 MB, most of it the `doc` JSON; `threads` holds 2.5 MB and `touches` 2.4 MB. In WAL mode the load's one transaction writes a WAL as large as the database, so a load needs twice the database's space until the close checkpoints it. The close removes the WAL and the shared-memory file.

| Corpus | Database after a full load | Per 1,000 items | Vacuumed copy | Vacuumed, pass-through columns blanked | Per 1,000 items | WAL and `-shm` at the end of a load in WAL mode, before close |
|--------|----------------------------|-----------------|---------------|----------------------------------------|-----------------|------------------------------------------------------------------|
| 1x | 23,531,520 bytes | 10.20 MB | 22,769,664 | 8,609,792 | 3.73 MB | 23,751,832 and 65,536 |
| 2x | 47,058,944 | 10.19 MB | 45,465,600 | 17,113,088 | 3.71 MB | 47,417,112 and 98,304 |
| 3x | 70,574,080 | 10.19 MB | 68,173,824 | 25,636,864 | 3.70 MB | 71,070,032 and 163,840 |
| 10x | 234,643,456 | 10.17 MB | 227,119,104 | 85,295,104 | 3.70 MB | 236,100,752 and 458,752 |
| Measured with | `s0427overhead load <corpus> <db>` | derived | `dbsize.py <db> runs` | the same | derived | `s0427overhead walload <corpus> <db>` |

| Other storage | Figure | Measured with |
|---------------|--------|---------------|
| WAL after a rescan with nothing changed, 1x and 10x | 0 bytes, `-shm` 32,768 | `wal-rescan.sh 1 10` |
| WAL after a rescan with one story changed, 1x and 10x | 49,472 bytes, `-shm` 32,768 | the same |
| flai's binary with the driver linked in | 37,353,257 bytes against 31,323,518: 6,029,739 more, 19 percent | `go build` of `flai/` with and without an overlay that imports `modernc.org/sqlite` |

**RAM.** Max RSS is from `wait4`'s rusage through `measure.py`. VmRSS is from `/proc/self/status`, read by the program at each step. The baseline row is flai's own report over the same corpus, one process: `statsread.Read`, then `metrics.Compute`, then at 1x to 3x Compute's parts timed once more.

| What | 1x | 2x | 3x | 10x | Measured with |
|------|----|----|----|-----|---------------|
| Full load, max RSS | 25,956 KB | 30,196 KB | 33,760 KB | 58,620 KB | `measure.py runs/load-<x>x 3 -- s0427overhead load <corpus> <db>`, median |
| Cold open, rescan, and the eight reads, max RSS | 24,140 KB | 30,572 KB | 37,468 KB | 70,256 KB | `measure.py ... 3 -- s0427overhead query <corpus> <db> <now> <out> 0 tables`, median |
| A process holding the read model warm, after three warm rounds, max RSS | 26,468 KB | 34,068 KB | 39,856 KB | 81,700 KB | `measure.py ... 1 -- s0427overhead query ... 3 tables` |
| The same kind of process, about to idle, VmRSS | 24,636 KB | | | 71,004 KB | `s0427overhead idle <corpus> <db> <now> db 130` |
| A process holding flai's parsed items, VmRSS | 44,912 KB | 74,236 KB | 103,284 KB | 298,512 KB | `s0427overhead serverss <corpus> <db> <now> items 5` |
| The same, after it opens the read model and answers the eight reads five times | 58,336 KB, 13,424 more | 91,212 KB, 16,976 more | 122,584 KB, 19,300 more | 349,780 KB, 51,268 more | the same run |
| A process after `statsread.Read` and one `metrics.Compute`, VmRSS | 62,336 KB | | 146,696 KB | | `s0427overhead serverss ... compute 5` |
| The same, after the read model and the eight reads five times | 70,676 KB, 8,340 more | | 155,032 KB, 8,336 more | | the same run |
| Baseline: flai's report over the corpus, max RSS | 61,016 KB | 103,328 KB | 145,752 KB | 436,696 KB | `measure.py runs/baseline-<x>x 1 -- s0427overhead baseline <corpus> <now> <out> <0 or 1>` |

The baseline's `stats.get` in the main checkout peaked at 79,368 KB, and a warm process held 71,304 KB after six calls. A process that holds only the parsed items gains 13 MB at 1x and 51 MB at 10x from the read model. Most of the gain is Go heap for the answers' rows: at 10x the heap in use rose from 247 MB to 307 MB over the reads. A process that has just run `metrics.Compute` gains 8.3 MB at 1x and at 3x, because the reads reuse the heap Compute freed. Of that, 2.5 MB is file-backed: the driver's code pages (`RssFile` 8,008 KB to 10,556 KB). A long-running flai serve would be the second kind once it has answered one `stats.get`.

**CPU.** The warm figures are the means of 30 rescans or requests in one process (`s0427overhead cpu <corpus-copy> <db> <now> 30`), from `getrusage` around each, with the median wall time in brackets. For the changed rescans the program appends a line to a story and sets its modification time 5 s back, so the rescan reads it once and trusts it. The per-hour and per-thousand rows are derived: 3,600 or 1,000 times the measured CPU.

| What | 1x | 2x | 3x | 10x | Measured with |
|------|----|----|----|-----|---------------|
| Full load, user and system CPU | 1,013 and 266 ms | 1,975 and 401 ms | 2,780 and 487 ms | 8,893 and 1,325 ms | the program's `getrusage` at the end of each load, median of three |
| Rescan, nothing changed, warm | 21.7 ms (18.2 ms) | 41.3 ms (34.7 ms) | 64.4 ms (53.8 ms) | 215.4 ms (183.2 ms) | `cpu` |
| Rescan, one story changed, warm | 25.2 ms (23.7 ms) | 45.6 ms (40.7 ms) | 66.3 ms (58.1 ms) | 218.4 ms (187.9 ms) | `cpu` |
| A request: rescan and the eight reads, warm | 108.1 ms (103.1 ms) | 215.1 ms (205.3 ms) | 335.1 ms (318.7 ms) | 1,170.3 ms (1,122.8 ms) | `cpu` |
| Derived: a rescan every second, nothing changed, CPU an hour | 78 s, 2.2 percent of a core | 149 s | 232 s | 776 s, 21.5 percent of a core | 3,600 × the warm rescan |
| Derived: per thousand requests, the rescan alone, or the rescan and the eight reads | 21.7 s, or 108 s | 41 s, or 215 s | 64 s, or 335 s | 215 s, or 1,170 s | 1,000 × the warm figures |
| Idle for 130 s, holding the read model open warm | 2.6 ms | | | 3.1 ms | `idle.sh`: `s0427overhead idle ... db 130` |
| Idle for 130 s, holding flai's parsed items, for comparison | 10.5 ms | | | 95.9 ms | `idle.sh`: `s0427overhead idle ... items 130` |

For comparison, a warm `stats.get` costs 5,502 ms of CPU a call on the main checkout (5,208 ms user, 294 ms system). The charts page asks for the stats again on every item change (`follow(['item'], ...)` in `flaiover/src/routes/charts/[kind]/+page.svelte`), so per request is how a rescan before each answer would be paid. Idle, the read model costs nothing measurable: 2.6 ms in 130 s at 1x and 3.1 ms at 10x, less than a process that holds only the parsed items.

Nearly all of a rescan is stat-ing every file: 19 to 20 ms of 20 to 21 ms at 1x, and 172 to 189 ms of 188 to 211 ms at 10x, about 6.5 µs a file. Reloading the one changed file and its derived rows took 6.7 to 7.3 ms at 1x and 22 to 25 ms at 10x. A rescan keyed on inode also reloads every file after anything that rewrites them all. A first, mistaken run rescanned a WAL database loaded from the snapshot against a copy with new inodes. It reloaded all 2,914 files in 2.2 s at 1x and all 29,140 in 20.0 s at 10x, about twice a full load (`runs/walrescan-<x>x-1.out`).

**Growth.** flai's report over the 10x corpus took 199 s in one run, under the 10-minute bound, so no estimate from a smaller corpus was needed. The other rows are medians of three, except the warm reads, the median of three rounds in one process.

| What | 1x | 2x | 3x | 10x | 10x over 1x | Measured with |
|------|----|----|----|-----|-------------|---------------|
| Baseline: `statsread.Read`, cold | 572 ms | 1,084 ms | 1,560 ms | 4,894 ms | 8.6 | `s0427overhead baseline`, one run each |
| Baseline: `metrics.Compute` | 3,715 ms | 10,975 ms | 21,826 ms | 193,865 ms | 52 | the same |
| `replayHolds`, timed apart | 3,647 ms | 10,826 ms | 23,095 ms | not timed | | the same, with parts |
| Compute's other parts that the SQL covers, summed | 25.8 ms | 54.5 ms | 88.7 ms | not timed | | the same |
| Baseline: user CPU of read and Compute | 5,986 ms | 17,683 ms | 35,497 ms | 331,018 ms | 55 | the same |
| SQLite: full load, database on ext4 | 1,098 ms | 2,076 ms | 2,890 ms | 9,186 ms | 8.4 | `measure.py runs/load-<x>x 3 -- ...` |
| SQLite: full load, database on tmpfs | 979 ms | 1,946 ms | 2,713 ms | 8,763 ms | 9.0 | `measure.py runs/loadtmp-<x>x 3 -- ...` |
| SQLite: one-shot, open, rescan, eight reads | 119 ms | 226 ms | 334 ms | 1,153 ms | 9.7 | `measure.py runs/oneshot-<x>x 3 -- ...` |
| SQLite: the eight reads, warm | 82.0 ms | 169.1 ms | 257.8 ms | 917.5 ms | 11.2 | `measure.py runs/warm-<x>x 1 -- ...` |
| of which `daily_states` | 43.9 ms | 91.3 ms | 137.3 ms | 498.3 ms | 11.4 | the same |
| SQLite: rescan in a fresh process, nothing changed | 21.1 ms | 38.1 ms | 57.1 ms | 191.3 ms | 9.1 | `rescan-test.py` |
| SQLite: rescan in a fresh process, one story changed | 27.0 ms | 45.0 ms | 64.1 ms | 200.0 ms | 7.4 | the same |
| SQLite: database | 23.5 MB | 47.1 MB | 70.6 MB | 234.6 MB | 9.97 | the load |

The read model grows linearly: its load, rescans, reads, and size each grow 7 to 11 times from 1x to 10x. flai's file read grows the same way, 8.6 times. `metrics.Compute` grows 52 times, about as the 1.7th power of the corpus, because `replayHolds` does: it grew 3.0 times at 2x and 6.3 times at 3x, while Compute's other parts grew 2.1 and 3.4 times.

Against the thresholds of the success measure at 10x, for the eight reads only: the warm reads take 918 ms, over 400 ms. The one-shot takes 1,153 ms, under 1,200 ms. A rescan for one changed file takes 200 ms, over 50 ms, almost all of it the stat of every file. The database is 235 MB against twice the corpus, 149 MB; without the pass-through columns it is 85 MB. The RAM added is 51 MB to a process holding only the items, and 8.3 MB to one that has run Compute (at 1x and 3x). The full load, 9.2 s, is far under the 10x cold baseline of about 199 s.

**Where the database would live.** flai keeps its per-checkout state in `.flai-cache/` of the main checkout, which git ignores. `Repo.CacheDir` resolves to the main checkout even from a story worktree, and `Repo.WipDir` is always the main checkout's `wip/` (`flai/internal/workitem/repo.go`). So one database per main checkout, such as `.flai-cache/readmodel.db`, would serve `flai serve` and every story worktree, because they all read the same files. A story worktree would pay nothing to build it. A fresh clone pays one full load, 1.1 s at 1x and 9.2 s at 10x, the first time anything reads the stats. It then holds 23.5 MB on disk at 1x and 235 MB at 10x, or 8.6 MB and 85 MB without the pass-through columns. With the load in WAL mode it briefly needs twice that. With one database shared across processes, `flai serve` and agents' one-shot `flai stats --json` from worktrees would all rescan into it. That needs WAL mode and a busy timeout, and the rescans then take turns for SQLite's one write lock. T-1802 did not measure concurrent writers.

### Postgres

T-1803 ran on 2026-10-16 between 07:09Z and 07:23Z, on the same host, with a load average of 0.6 to 2.8. Its program, scripts, outputs, and commands are in `bin/experiments/S-0427/postgres/`, with the steps appended to the scratch folder's `README.txt`. Every query uses T-1801's corpora, its fixed now, and flai's reports over each corpus (`overhead/runs/baseline-<x>x.json`) as the reference.

**What a server could give that SQLite cannot here.**

| Claimed advantage | What it means for this read model |
|-------------------|-----------------------------------|
| One server shared by `flai serve`, each project's MCP server, the dashboard, and every story worktree | A SQLite file in the main checkout's `.flai-cache/` is already one model shared by every process on the host (see [Where the database would live](#overhead)). What a server adds is reach over the network: a process in another container or on another machine. The dashboard does not need that, because it reads through `flai serve`'s API, not the database. One server could also hold every project's database, but each project's data is its own, so sharing a server shares nothing. |
| Concurrent writers | Postgres locks rows, not the database, so writers of different rows do not wait for each other. Every writer here is a rescan of the same checkout, which writes the same rows. The test below measures what that gives. |
| `LISTEN`/`NOTIFY` to push changes to the dashboard instead of polling | The dashboard is already pushed its changes by `flai serve`, whose watcher polls three folders every 300 ms (`flai/internal/watch/watch.go`, after [ADR-0029](../adrs/0029-the-dashboard-reaches-a-project-only-through-flai-on-the-host.md), by which the dashboard reaches a project only through flai). The files stay the source of truth, and agents write files, not the database. So something must still notice a file change, and a notification can only follow a rescan; it removes no poll. It would let a process other than `flai serve`, such as a worktree's MCP server, learn of a rescan without a watcher of its own. |

What it costs:

- **A daemon and a container to keep running.** flai is one static binary built with `CGO_ENABLED=0`. With a server, every machine that runs flai also needs Docker or a Postgres install, and every one-shot `flai stats --json` must cope with the server being down, which means keeping the file reads as a fallback. The image `postgres:16-alpine` takes 420 MB on disk. The first start, initdb included, took 2.65 s to accept connections, and a restart took 1.8 s.
- **A port.** 55432 here, bound to 127.0.0.1. Each host needs one free port, which every flai process must be told.
- **A credential.** Every flai process needs the password in its connection string, from configuration or the environment, never the repository. Here it was a throw-away `local`.
- **A backup story.** The read model is derived: it reloads from the files in 1.3 s at 1x and 10 s at 10x, so it needs no backup. The volume still needs looking after. Postgres keeps WAL up to `max_wal_size`, 1 GB by default, and a major version upgrade needs the data dumped and reloaded or `pg_upgrade`. For a derived model, dropping and reloading would do.
- **RAM and idle CPU**, measured below.

**The setup.** The container is the stock image, started as `docker run -d --name s0427-pg -p 127.0.0.1:55432:5432 -e POSTGRES_PASSWORD=local -v s0427-pg-data:/var/lib/postgresql/data postgres:16-alpine`. It runs PostgreSQL 16.15 with the image's defaults: `shared_buffers` 128MB, `work_mem` 4MB, `fsync` and `synchronous_commit` on, `max_wal_size` 1GB. The volume is on docker's root, `/var/lib/docker`, on the host's disk. Each corpus has its own database, created with `LC_COLLATE 'C'` so text sorts by bytes, as in SQLite. The client connects over TCP to 127.0.0.1 with SCRAM-SHA-256 and `sslmode=disable`.

**The driver and the program.** `github.com/jackc/pgx/v5` v5.10.0, pure Go like the SQLite driver. It was already in the operator's module cache, so `go get` against that cache as a file proxy took 10 s, with no download; it added `pgpassfile`, `pgservicefile`, and `golang.org/x/text` v0.29.0. The program `s0427pg` is in package `main` under `flai/internal/s0427pg`, mapped there by `go build -overlay`, with `-modfile` pointing at `postgres/go.mod`, T-1801's copy plus pgx. As with T-1801, nothing is written into the module. It uses pgx's own interface, not `database/sql`. It parses each file with flai's readers into T-1801's rows. The full load writes them with `COPY`, one table at a time, in one transaction, then fills the derived tables with two `INSERT ... SELECT` and runs `ANALYZE`. A second load mode sends one batch of `INSERT`s per file instead. The rescan is T-1801's: it reads the `files` table, stats every file, and reloads what changed. For each changed file it queues the deletes, the inserts, and the derived rows in one pgx batch, in one transaction.

```bash
cd flai && GOMODCACHE=$S/gomodcache GOPROXY=file:///home/alex/go/pkg/mod/cache/download GOFLAGS=-modcacherw \
  go build -modfile=$S/postgres/go.mod -overlay=$S/postgres/overlay.json -o $S/postgres/s0427pg ./internal/s0427pg
./s0427pg load /tmp/s0427-corpus "postgres://postgres:local@127.0.0.1:55432/s0427_1x?sslmode=disable" copy
./s0427pg query /tmp/s0427-corpus "<same dsn>" 2026-10-16T06:04:02Z answers-1x.json 3
```

**The port.** The schema is T-1801's, with `INTEGER` as `BIGINT`, `REAL` as `DOUBLE PRECISION` (Postgres's `REAL` has four bytes), no `WITHOUT ROWID`, and the column `by` quoted, because it is a reserved word. The eight reads are T-1801's over the derived tables, in `postgres/queries.go`, with these changes:

| SQLite | Postgres |
|--------|----------|
| `unixepoch(date(x, 'unixepoch'))` | `(x / 86400) * 86400`, `x` a non-negative `BIGINT` |
| `date(x, 'unixepoch')` | `to_char(to_timestamp(x) AT TIME ZONE 'UTC', 'YYYY-MM-DD')` |
| `date(s, 'unixepoch', '-6 days', 'weekday 1')`, the Monday that starts the week | `date_trunc('week', to_timestamp(s) AT TIME ZONE 'UTC')` |
| `strftime('%G-W%V', x, 'unixepoch')` | `to_char(to_timestamp(x) AT TIME ZONE 'UTC', 'IYYY-"W"IW')` |
| `MAX(a, b)` with two arguments | `GREATEST(a, b)` |
| `SUM(state <> 'cancelled')` | `count(*) FILTER (WHERE state <> 'cancelled')` |
| `SUM` and `AVG` of integers | cast back with `::bigint` and `::float8`: Postgres returns `NUMERIC` for both |
| `x * 1.0 / y` | `x::float8 / y` |
| `GROUP BY metric`, selecting `n` | `GROUP BY metric, n`: Postgres requires every selected column grouped |
| `:now`, `:start`, `:type` | pgx named arguments `@now::bigint`, `@start::bigint`, `@type` |

**Checks.** The answers matched flai's report over the same corpus on every value at every scale: 5,321, 9,517, 13,713, and 43,085 values, none differing. That held for the one-shot and for the warm answers (`compare.py runs/answers-<x>x.json ../overhead/runs/baseline-<x>x.json`), and again after a restart of the container and after every round of the concurrency test. At each scale a load with `COPY` and a load with batched inserts dumped to the same SHA-256, `ebd18f6b…` at 1x (`s0427pg dump`).

**Timings.** Medians of three, each with its range, except where a row says one run. Wall time, CPU, and max RSS of a fresh process come from `measure.py`; in-process times come from `time.Now` around each step. The baseline and SQLite figures are those above.

| What | Baseline | SQLite, 1x and 10x | Postgres, 1x | Postgres, 10x | Measured with |
|------|----------|--------------------|--------------|---------------|---------------|
| Full load from empty, fresh process | `stats.get` cold, 4715 ms | 1,098 ms and 9,186 ms | 1,345 ms (1,326 to 1,371); in process: parse 561 ms, `COPY` 313 ms, derived tables and commit 142 ms, `ANALYZE` 197 ms, schema 103 ms, connect 15 ms | 10,081 ms (10,010 to 10,096); parse 5,031 ms, `COPY` 2,830 ms | `measure.py runs/load-<x>x 3 -- s0427pg load <corpus> <dsn> copy` |
| Full load with a batch of `INSERT`s per file, one run | | | 3,598 ms | 30,521 ms | `measure.py runs/loadbatch-<x>x 1 -- s0427pg load ... batch` |
| Rescan, nothing changed, fresh process | none | 21.1 ms and 191.3 ms | 22.4 ms (20.8 to 23.2), and a connect of about 8 ms before it | 197.3 ms (184.5 to 208.8), and about 7 ms to connect | `rescan-test.py ./s0427pg <corpus copy> <dsn> 3` |
| Rescan, one story changed, fresh process | none | 27.0 ms and 200.0 ms | 36.2 ms (34.5 to 37.4), of which 13.6 to 15.6 ms writes the file's rows in one transaction | 203.1 ms (199.3 to 203.1), of which 30.9 to 33.4 ms writes | the same |
| One-shot: connect, rescan, the eight reads, write the answers | `flai stats --json`, 4701 ms | 119 ms and 1,153 ms | 186.6 ms (183.2 to 188.9): connect 9.6 ms, rescan 22.4 ms, reads 146.4 ms; user 24.5 ms, system 18.3 ms; max RSS 16,996 KB | 965.2 ms (957.5 to 1,003.1); user 124.6 ms, system 134.2 ms; max RSS 52,500 KB | `measure.py runs/oneshot-<x>x 3 -- s0427pg query ... 0` |
| The same, first query after a container restart, one run | | | 197.9 ms | | `measure.py runs/oneshot-after-restart-1x 1 -- ...` |
| Warm: the eight reads, rounds 2 to 4 on one connection | `stats.get` warm, 4022 ms | 82.0 ms and 917.5 ms | 111.6 ms (109.9 to 112.1) | 710.5 ms (708.0 to 742.6) | `measure.py runs/warm-<x>x 1 -- s0427pg query ... 3` |
| Warm, a long-running process: the eight reads for 20 s, median round | `stats.get` warm, 4022 ms | as above; SQLite's rounds do not change | 63.8 ms (61.3 to 151.6), 309 rounds | 684.4 ms (678.6 to 759.9), 29 rounds | `s0427pg loop <dsn> <now> 20` |

The two warm rows differ because of how Postgres plans. pgx prepares each query once per connection. Postgres plans a prepared statement afresh for each of its first five runs, then may keep one generic plan, and here it did. The loop's first six rounds at 1x took 151.6, 116.1, 113.7, 114.1, 111.2, and 68.6 ms. A long-running `flai serve` would answer at the loop's figure. An agent's one-shot pays the planning every time.

| Read, warm | SQLite, 1x | Postgres rounds 2 to 4, 1x | Postgres long-running, 1x | SQLite, 10x | Postgres long-running, 10x |
|------------|------------|----------------------------|---------------------------|-------------|----------------------------|
| `columns` | 0.29 ms | 1.37 ms | 0.72 ms | 1.28 ms | 1.19 ms |
| `cycle_time` | 1.18 ms | 2.16 ms | 2.17 ms | 12.0 ms | 27.5 ms |
| `time_in_state` | 4.37 ms | 6.59 ms | 6.94 ms | 50.7 ms | 59.9 ms |
| `summary` | 9.82 ms | 6.83 ms | 5.94 ms | 102.4 ms | 64.4 ms |
| `throughput_week` | 2.01 ms | 1.47 ms | 1.36 ms | 22.3 ms | 8.46 ms |
| `throughput_day` | 0.55 ms | 1.15 ms | 1.05 ms | 5.89 ms | 7.44 ms |
| `daily_states` | 43.9 ms | 29.4 ms | 28.8 ms | 498.3 ms | 359.0 ms |
| `spend_models` | 19.7 ms | 62.1 ms | 16.5 ms | 224.9 ms | 155.7 ms |
| All eight | 82.0 ms | 111.6 ms | 63.8 ms | 917.5 ms | 684.4 ms |
| Measured with | `overhead/runs/warm-<x>x-1.out` | `postgres/runs/warm-1x-1.out` | `postgres/runs/loopq-1x-1.out` | `overhead/runs/warm-10x-1.out` | `postgres/runs/loopq-10x-1.out` |

**Growth.** Medians of three, as above.

| What | 1x | 2x | 3x | 10x | 10x over 1x |
|------|----|----|----|-----|-------------|
| Full load with `COPY` | 1,345 ms | 2,471 ms | 3,472 ms | 10,081 ms | 7.5 |
| One-shot | 186.6 ms | 270.5 ms | 334.5 ms | 965.2 ms | 5.2 |
| Warm, rounds 2 to 4 | 111.6 ms | 165.8 ms | 219.7 ms | 710.5 ms | 6.4 |
| Rescan, nothing changed | 22.4 ms | 38.3 ms | 56.3 ms | 197.3 ms | 8.8 |
| Rescan, one story changed | 36.2 ms | 54.0 ms | 71.9 ms | 203.1 ms | 5.6 |
| Database after the first load | 36,225,551 bytes | 60,744,207 | 87,777,807 | 264,094,223 | 7.3 |

**Storage.**

| What | 1x | 10x | Measured with |
|------|----|-----|---------------|
| SQLite database file, for comparison | 23,531,520 bytes | 234,643,456 bytes | T-1802 |
| `pg_database_size` after the first load | 36,225,551 bytes | 264,094,223 bytes | `s0427pg load`, which reads it at the end |
| of which an empty database, the system catalog | 7,537,167 bytes | 7,537,167 bytes | `pg_database_size` of a database just created |
| The read model's tables with their indexes and TOAST, after five loads | 27,860,992 bytes, 13,672,448 of it table rows | 256,032,768 bytes, 116,563,968 of it rows | `s0427pg size`: `pg_total_relation_size`, `pg_relation_size` |
| The `items` table, `doc` JSON included | 6,176,768 bytes, against SQLite's 11.9 MB: Postgres compresses long text | 57,335,808 bytes | the same |
| The volume, after every load of this task: 1x to 10x, a rescan copy, and the system databases | 1,820,669,768 bytes: databases 746,449,087, WAL 1,073,750,016 (64 segments) | | `du -sb` in the container; `docker system df -v` reported 1.82GB |

Each reload with `DROP SCHEMA` left the database about 0.4 MB larger, in the catalog. The WAL had grown to `max_wal_size` after about 25 loads in a few minutes. Postgres recycles it towards `min_wal_size`, 80 MB, as checkpoints pass. A volume holding one project's model at 1x would need the database, about 23 MB of system databases, and between 80 MB and 1 GB of WAL. That is derived from the figures above, not measured on its own.

**RAM.** `docker stats` reports the container's memory less its inactive page cache, so it moves with what the loads wrote. The server's own memory is its processes' resident sets, from `/proc/<pid>/status` in the container.

| What | Figure | Measured with |
|------|--------|---------------|
| Server, idle after a restart on the existing volume | 17.53 MiB; six processes with 12,516 KB of private memory | `docker stats --no-stream`; `RssAnon` per process |
| Server, idle after the loads and queries | 485 to 523 MiB. Of the cgroup's 2.2 GB, private memory 5.0 to 6.3 MB, shared memory 147 MB (`shared_buffers` filled), the rest page cache: 2.16 GB, 1.65 GB of it inactive | `docker stats`; `memory.stat` of the container's cgroup |
| Server, during the query loop | 533 to 538 MiB at 1x, 555 to 575 MiB at 10x, against 522 MiB idle just before | `docker stats` three times in each 32 s loop (`pg-resources.sh`) |
| Client, one-shot, max RSS | 16,996 KB at 1x, 52,500 KB at 10x; SQLite's, engine included, 24,140 KB and 70,256 KB | `measure.py` |
| Client, long-running loop, max RSS | 18,236 KB at 1x, 20,556 KB at 10x; SQLite's warm process 26,468 KB and 81,700 KB | the same |
| Client, full load, max RSS | 46,016 KB at 1x, 297,716 KB at 10x, because `COPY` here takes every row from memory; 21,428 KB and 53,668 KB with batched inserts | the same |

With Postgres the database engine runs in the container, not in flai's process. So flai's process is lighter than with SQLite, and the server is a second process that holds its own memory whether or not anything asks it.

**CPU.** The container's CPU is its cgroup's `usage_usec` from `cpu.stat`, read before and after each window. `docker stats` gives a percentage over about a second. The per-round rows are derived: the measured CPU divided by the rounds.

| What | Postgres | SQLite, for comparison | Measured with |
|------|----------|------------------------|---------------|
| Idle 130 s, one client connection held open | Container 270 ms, about 0.2 percent of a core (derived); `docker stats` 0.00 percent; client 0.6 ms | 2.6 ms in 130 s at 1x | `pg-resources.sh 130 32 1 10` |
| Query loop, 32 s at 1x | Container 29.8 s of CPU in 32.2 s, `docker stats` 91.8 to 92.3 percent; client 1.9 s; 495 rounds | | the same |
| Derived: CPU per round of the eight reads, 1x | 60.2 ms server and 3.9 ms client, 64.1 ms | 87.5 ms, the eight reads alone (`cpu`, T-1802) | the same |
| Query loop, 32 s at 10x | Container 32.2 s in 32.5 s, 97.4 to 98.8 percent; client 1.1 s; 47 rounds | | the same |
| Derived: CPU per round, 10x | 685 ms server and 23.2 ms client, 708 ms | 943 ms | the same |

The server answers one connection's queries on one core, as SQLite does in flai's process.

**Two rescans at once.** `concurrent.py` copies a corpus, loads a model from it, and in each of three rounds changes k work items. It sets their modification time 5 s back and starts two rescans as two processes at once. A third rescan then runs alone, and for comparison k other items change and one rescan runs alone. SQLite is T-1801's rescan unchanged, on a WAL database with `busy_timeout` 10000 ms and `BEGIN IMMEDIATE` (`postgres/sqlite-wal/concur.go`, mapped beside T-1801's files). Postgres is the same rescan as ported, and a second form that first takes `pg_advisory_xact_lock` and reads the `files` table inside that transaction. Times are each rescan's in-process total.

| Corpus, items changed | SQLite, WAL and busy timeout | Postgres, T-1801's rescan as ported | Postgres, advisory lock first |
|-----------------------|------------------------------|-------------------------------------|-------------------------------|
| 1x, 1 | Both succeeded in 3 of 3 rounds, and both reloaded the file: 26.7 to 38.6 ms, against 30.4 to 31.0 ms alone | One of the two failed in 3 of 3 rounds with `duplicate key value violates unique constraint "files_pkey"`; the other took 37.3 to 37.6 ms, against 35.1 to 36.0 ms alone | Both succeeded. The second waited 30.8 to 35.2 ms for the lock, then found nothing changed: 55.0 to 58.8 ms, against 36.5 to 36.9 ms alone |
| 1x, 100 | Both succeeded; the first took 105 to 152 ms, the second waited and reloaded the same files, 213 to 292 ms; alone 104 to 178 ms | One failed in 3 of 3, the same error; the other took 145 to 261 ms; alone 136 to 282 ms | Both succeeded; the second waited 150 to 253 ms, total 174 to 276 ms; alone 139 to 275 ms |
| 10x, 1 | Both succeeded: 201 to 224 ms, against 210 to 216 ms alone. Stat-ing 29,140 files dominates, and the write is short | One failed in 3 of 3; the other took 202 to 208 ms; alone 201 to 203 ms | Both succeeded; the second waited 198 to 205 ms, total 393 to 407 ms; alone 199 to 216 ms |
| Measured with | `python3 -I concurrent.py <corpus> <x> <k> 3`, `runs/concurrent-<x>x-k<k>.out` | the same | the same |

After every round the rescan alone found nothing left to change, and the answers matched flai's on every value (`runs/conc-answers-*.compare`). Two rescans of one checkout write the same rows, so they take turns on either database. On SQLite the busy timeout makes them take turns with no change to the code, and the second repeats the first's reload. On Postgres, row locks do not let them run together. T-1801's delete-then-insert fails outright when two writers race, so the rescan has to be written for it, here with an advisory lock that makes the second wait for the first.

**`LISTEN`/`NOTIFY`.** One connection listened on a channel. A second connection appended a line to a story, rescanned it with `pg_notify` in the same transaction, and committed. From sending the commit to the notification's arrival took 1.36 ms, median of 20 (1.25 to 1.86 ms). The commit itself took 1.31 ms (`s0427pg listen <corpus copy> <dsn> 20`).

**Against the success measure.** Postgres is worth its server only if it answers the same queries at least twice as fast as SQLite. Derived from the figures above: a long-running process answers the eight reads 1.29 times as fast at 1x (63.8 ms against 82.0 ms) and 1.34 times as fast at 10x (684 ms against 918 ms). An agent's one-shot is slower at 1x (187 ms against 119 ms) and 1.19 times as fast at 10x (965 ms against 1,153 ms). The full load and the rescans are about the same or slower. Neither database gives the hold replay a SQL form, and that replay is 97 percent of `metrics.Compute`. The threshold is not met.

## Results

Later tasks of S-0427 write the results, against the success measure above.

## Recommendation

Provisional, from the baseline alone: adapt. A read model that only replaces the file reads saves at most about 80 to 670 ms of a 4 s answer, so it is worth having only if its queries also take over the computation. The SQLite figures point the same way: the chart reads in SQL take 83 ms where flai's Go takes about 25 ms over parsed items, and the 3.4 s hold replay has no SQL form, so the replay, not the reads, is what to change first. Postgres does not change that: a long-running process answers the eight reads 1.3 times as fast as SQLite and an agent's one-shot at most 1.2 times, short of the twice the success measure asks, and two rescans still take turns, because they write the same rows. At ten times the corpus `metrics.Compute` grows 52 times, to 194 s, while the read model's load, reads, and rescans grow about ten times, and a read model would need flai's file watch rather than a stat of every file to stay within 50 ms a change. Later tasks of S-0427 settle it once the SQLite and Postgres comparisons are in.
