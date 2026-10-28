---
id: D-0402
title: A test fixture project that flai writes items in needs a .gitignore listing .flai-cache/, since the items lock lives there and the write check refuses a project without it
kind: cause
sources: [S-0479, S-0484, I-0191]
topics: [testing, code, flai, flaiover]
paths: [flaiover/src/lib/server/writes.test.ts, flai/cmd]
date: 2026-10-28
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0402 A test fixture project that flai writes items in needs a .gitignore listing .flai-cache/, since the items lock lives there and the write check refuses a project without it

Every work item write takes one items lock under `.flai-cache/` (ADR-0184, S-0479), and the write's check refuses a project whose `.gitignore` does not ignore it (`layout.gitignore`). S-0479 gave three `flai/cmd` test fixtures a `.gitignore` for `.flai-cache/`. S-0484 then found two `writes.test.ts` tests failing on main for the same reason: their temporary projects had none. That is I-0191, left for the operator to choose a remedy.

A new fixture that writes items through flai writes that `.gitignore` too. This is the fixture project's own ignore file, not a fixture of this repository left out of git.

Sources: S-0479, S-0484, I-0191.
