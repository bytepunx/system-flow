---
id: D-0472
title: A test fixture that edits a file in place uses plain sed through a temporary file, never sed -i, since BSD sed on the flai workflow's macOS job reads the next word as a backup suffix
kind: decision
sources: [S-0572]
topics: [testing, ci, go, code]
paths: [flai/internal/mergequeue/preflight_test.go, ".github/workflows/flai.yml"]
date: 2026-11-03
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0472 A test fixture that edits a file in place uses plain sed through a temporary file, never sed -i, since BSD sed on the flai workflow's macOS job reads the next word as a backup suffix

S-0572 fixed the `flai` workflow's macOS job and the failed `flai/v1.81.0` release. `TestThePreflightConfirmsALintFindingInAScratchWorktree/the_fix_clears_it` in `flai/internal/mergequeue/preflight_test.go` fixed a file with `sed -i s/FIXABLE/fixed/`. BSD sed reads the word after `-i` as a backup suffix, so on macOS the fix never ran and the finding stayed. This host runs GNU sed, so the test passed here (S-0572's story, Goal).

The fix, from S-0572's Decisions (`wip/archive/agents/S-0572.md`):

```text
T-2446: the fixture rewrites each file through a temporary one with plain `sed`, not `sed -i.bak` or `perl -pi`, so it depends on no flag that differs between GNU and BSD and on no perl.
```

A later test or script that rewrites a file does the same: `sed '...' file > tmp && mv tmp file`, with no `-i` and no perl.

Source: S-0572.
