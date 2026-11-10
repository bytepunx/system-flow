---
id: T-2969
type: task
nature: remediation
title: Give repo.test.ts's whole-monorepo describe a 30 s timeout, as search.test.ts's monorepo describe has
status: backlog
parent: S-0695
owner: alex
created: 2026-11-10T10:28:26Z
updated: 2026-11-10T10:28:26Z
transitions: []
stream: S-0695
tags: [flaiover, testing]
touches: [flaiover/src/lib/server/repo.test.ts]
---
# T-2969 Give repo.test.ts's whole-monorepo describe a 30 s timeout, as search.test.ts's monorepo describe has

## Work

`flaiover/src/lib/server/repo.test.ts` › Repo on the monorepo › reads every item with a consistent status runs `flai hostapi items.list` with the archive on the whole monorepo under vitest's 5 s default. I-0269 records three failures, each on a loaded host, in stories that do not change what the test reads. The planner timed the same read on the host at 0.75 to 0.90 s with a load average of 130 on 24 cores. A 5 s limit leaves too little room when vitest's workers share that host.

- Give the `describe.skipIf(...)('Repo on the monorepo', ...)` block `{ timeout: 30000 }`, as `flaiover/src/lib/server/search.test.ts`'s `adrs on the monorepo` describe has.
- Before you change anything, check D-0524's other cause: other agents write items in the main checkout while the test reads them. Read how `flai hostapi items.list` reads an item that a save is writing at the same moment. If a half-written file can give a status that does not match its last transition, record it in the narrative's Decisions and on the plan thread, and leave it to its own issue. That is a different cause from the timeout, and this task does not fix it.
- No test can reproduce the timeout reliably, because it depends on the host's load. Say so in the story's Notes, which the first criterion allows ("where one fits").

Waits for no task: it is the first layer.

## Done when

- The monorepo describe in `flaiover/src/lib/server/repo.test.ts` names `{ timeout: 30000 }`.
- `flai test flaiover/src/lib/server/repo.test.ts` passes.
- The narrative's Decisions give the measured read time, the reason for 30 s, and the finding on concurrent writes.

## Notes
