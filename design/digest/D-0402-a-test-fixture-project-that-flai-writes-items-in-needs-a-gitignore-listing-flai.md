---
id: D-0402
title: A test fixture project that flai writes items in needs a .gitignore listing .flai-cache/, since the items lock lives there and the write check refuses a project without it
kind: cause
sources: [S-0479, S-0484, I-0191, S-0531, I-0221]
topics: [testing, code, flai, flaiover, performance]
paths: [flaiover/src/lib/server/writes.test.ts, flai/cmd, flai/internal/storygit/committed.go, design/system/server-performance.md]
date: 2026-10-28
status: active
absorbed_by: ""
asked_again: []
cited: [TH-0696]
---

# D-0402 A test fixture project that flai writes items in needs a .gitignore listing .flai-cache/, since the items lock lives there and the write check refuses a project without it

Every work item write takes one items lock under `.flai-cache/` (ADR-0184, S-0479), and the write's check refuses a project whose `.gitignore` does not ignore it (`layout.gitignore`). S-0479 gave three `flai/cmd` test fixtures a `.gitignore` for `.flai-cache/`. S-0484 then found two `writes.test.ts` tests failing on main for the same reason: their temporary projects had none. That is I-0191, left for the operator to choose a remedy.

A new fixture that writes items through flai writes that `.gitignore` too. This is the fixture project's own ignore file, not a fixture of this repository left out of git.

The same holds for flai's own code that keeps a file under `.flai-cache`. S-0531's `storygit.Committed` first made the folder to keep `stories-commits.json`, and five `flai/cmd` check tests failed on `layout.gitignore` in `metrics/testdata/good`, a fixture with no `.gitignore`. S-0531 decided: "`Committed` writes the kept file only into a `.flai-cache` folder that already exists", so a project without one keeps nothing and reads `git log` each time (`design/system/server-performance.md`). The narrative calls it "the same class of bug as S-0525". The close-out's Go tier ran only the changed packages and passed over the failure, recorded as I-0221.

Sources: S-0479, S-0484, I-0191, S-0531 (`wip/archive/agents/S-0531.md` § Decisions), I-0221.
