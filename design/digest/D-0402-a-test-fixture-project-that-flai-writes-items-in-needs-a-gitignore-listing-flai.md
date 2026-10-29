---
id: D-0402
title: A test fixture project that flai writes items in needs a .gitignore listing .flai-cache/, since the items lock lives there and the write check refuses a project without it
kind: cause
sources: [S-0479, S-0484, I-0191, S-0531, I-0221, S-0525]
topics: [testing, code, flai, flaiover, performance, cli]
paths: [flaiover/src/lib/server/writes.test.ts, flai/cmd, flai/internal/storygit/committed.go, design/system/server-performance.md, flai/internal/itemnew/itemnew.go, flai/internal/check/scope.go]
date: 2026-10-28
status: active
absorbed_by: ""
asked_again: []
cited: [TH-0696]
---

# D-0402 A test fixture project that flai writes items in needs a .gitignore listing .flai-cache/, since the items lock lives there and the write check refuses a project without it

Every work item write takes one items lock under `.flai-cache/` (ADR-0184, S-0479), and `flai check` reports a project whose `.gitignore` does not ignore it (`layout.gitignore`). S-0479 gave three `flai/cmd` test fixtures a `.gitignore` for `.flai-cache/`. S-0484 then found two `writes.test.ts` tests failing on main for the same reason: their temporary projects had none. That is I-0191.

S-0525 fixed I-0191 by [ADR-0213](../adrs/0213-a-checked-creation-of-a-work-item-refuses-only-a-new-finding-inside-the-new.md): a checked creation scopes its second check run to the new item's story, as ADR-0197 scopes a checked edit, so the first item written in a project with no `.gitignore` is kept with `layout.gitignore` among its warnings instead of refusing itself. The ADR keeps the lock where ADR-0184 puts it, and flai writes no `.gitignore` of its own: `flai check` on such a project still reports `layout.gitignore`.

So a new fixture that writes items through flai, or runs `flai check` after one, still writes that `.gitignore`. This is the fixture project's own ignore file, not a fixture of this repository left out of git.

The same holds for flai's own code that keeps a file under `.flai-cache`. S-0531's `storygit.Committed` first made the folder to keep `stories-commits.json`, and five `flai/cmd` check tests failed on `layout.gitignore` in `metrics/testdata/good`, a fixture with no `.gitignore`. S-0531 decided: "`Committed` writes the kept file only into a `.flai-cache` folder that already exists", so a project without one keeps nothing and reads `git log` each time (`design/system/server-performance.md`). The close-out's Go tier ran only the changed packages and passed over the failure, recorded as I-0221.

Sources: S-0479, S-0484, I-0191, S-0525 (`wip/archive/agents/S-0525.md` § Decisions), S-0531 (`wip/archive/agents/S-0531.md` § Decisions), I-0221.
