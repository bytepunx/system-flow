---
id: T-2875
type: task
nature: feature
title: kluster.Evidence reads a story namespace's last deploy and checks run and says whether both passed at the branch's head
status: backlog
parent: S-0662
owner: alex
created: 2026-11-10T08:26:45Z
updated: 2026-11-10T08:26:45Z
transitions: []
stream: S-0662
tags: [cli]
touches: [flai/internal/kluster/evidence.go, flai/internal/kluster/evidence_test.go]
---
# T-2875 kluster.Evidence reads a story namespace's last deploy and checks run and says whether both passed at the branch's head

## Work

Add `flai/internal/kluster/evidence.go`: one function that, for a story, reads its namespace's deploy record (`ReadRecord`) and the last checks run S-0661 records, and answers each with its commit, outcome, who asked, when it ended, and its log paths. It also answers whether the deploy passed at the story branch's head and whether the checks passed there, and the commands that deploy and check, as S-0660 and S-0661 named them.

It answers that nothing applies while the `kluster` host action is off or the block names no `deploy`. Read the checks record by the path and shape S-0661 writes, not through `flai/internal/serve`, which imports `kluster`. Follow S-0659's ADR on what the record holds.

Every later code task reads this, so it waits for none.

## Done when

- The function answers for: no block, the action off, no `deploy`, no record, a deploy at an older commit, a failed deploy, a passed deploy with no checks run, checks failed at the head, and both passed at the head.
- `evidence_test.go` covers each case with fake records, and needs no cluster or Docker.
- `flai test` passes on the two files.

## Notes

- `KlusterEnv` in `flai/internal/serve/kluster.go` judges a passed record at the branch head today; this is the same test, importable outside `serve`.
