---
id: T-2875
type: task
nature: feature
title: kluster.Evidence reads a story namespace's last deploy and checks run and says whether both passed at the branch's head
status: done
parent: S-0662
owner: alex
created: 2026-11-10T08:26:45Z
updated: 2026-11-13T13:55:42Z
transitions:
  - to: ready
    at: 2026-11-13T13:44:01Z
    by: agent-S-0662
  - to: in-progress
    at: 2026-11-13T13:44:02Z
    by: agent-S-0662
  - to: done
    at: 2026-11-13T13:55:42Z
    by: agent-S-0662
stream: S-0662
tags: [cli]
touches: [flai/internal/kluster/evidence.go, flai/internal/kluster/evidence_test.go]
usage:
  source: log
  seconds: 700
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 26
      output: 10089
      cache_read: 1938538
      cache_write: 83015
      cost: 1.0428
---
# T-2875 kluster.Evidence reads a story namespace's last deploy and checks run and says whether both passed at the branch's head

## Work

Add `flai/internal/kluster/evidence.go`: one function that, for a story, reads its namespace's deploy record (`ReadRecord`) and the last checks run S-0661 records, and answers each with its commit, outcome, who asked, when it ended, and its log paths. It also answers whether the deploy passed at the story branch's head and whether the checks passed there, and the commands that deploy and check, as S-0660 and S-0661 named them.

"At the head" is ADR-0292 §3's: a record is at the head when its commit is the head, or when every commit from the record's commit to the head changes only paths under the manifest's `design`, `docs`, and `wip` folders (`layout`), as git reads it in the main checkout through `execx.Runner`. Give that judgement as an exported function of its own, so that T-3306 can make the acceptance gate and the review pass use it, and report the commit that was judged and whether it is the exact head.

It answers that nothing applies while the `kluster` host action is off or the block names no `deploy`. Read the checks record by the path and shape S-0661 writes, not through `flai/internal/serve`, which imports `kluster`. Follow S-0659's ADR on what the record holds.

Every later code task reads this, so it waits for none.

## Done when

- The function answers for: no block, the action off, no `deploy`, no record, a deploy at an older commit, a failed deploy, a passed deploy with no checks run, checks failed at the head, both passed at the head, and both passed at a commit followed only by commits under `design`, `docs`, and `wip` (at the head) or by a commit outside them (not at the head).
- `evidence_test.go` covers each case with fake records, and needs no cluster or Docker.
- `flai test` passes on the two files.

## Notes

- `KlusterEnv` in `flai/internal/serve/kluster.go` judges a passed record at the branch head today; this is the same test, importable outside `serve`.
