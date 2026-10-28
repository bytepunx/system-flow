---
id: T-2240
type: task
nature: remediation
title: The main checkout's items pass item.chronology, and I-0183 is closed with the fix named
status: done
parent: S-0522
owner: alex
created: 2026-10-28T14:22:12Z
updated: 2026-10-28T14:40:02Z
transitions:
  - to: ready
    at: 2026-10-28T14:39:38Z
    by: agent-S-0522
  - to: in-progress
    at: 2026-10-28T14:39:39Z
    by: agent-S-0522
  - to: done
    at: 2026-10-28T14:40:02Z
    by: agent-S-0522
stream: S-0522
tags: [cli]
touches: [design/issues/I-0183-flai-check-finds-item-chronology-outside-the-story-at-close-out.md, design/issues/summary.md]
after: [T-2237, T-2238, T-2239]
usage:
  source: log
  seconds: 23
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 13
      output: 5434
      cache_read: 579634
      cache_write: 30872
      cost: 0.4254
---
# T-2240 The main checkout's items pass item.chronology, and I-0183 is closed with the fix named

## Work

With the three changes on the branch, run `flai check --strict` from the story's worktree, which reads the main checkout's items, and confirm no `item.chronology` stands on any item; one that does is a file a hand edit left, so name it in the narrative and leave it, since the story's branch does not change the main checkout's `wip/`. Run the tests of the three packages once more together. Then close the issue from the worktree, as `issue commands` are run there: `flai issue close I-0183 --reason` naming the ADR, that `Repo.Save` now keeps `updated` at or after the last transition and `created`, and that the replanner's forecast write, the merge queue's send-back, and the touches widening take their clock after reading the item; `flai issue close` rewrites the issue file and `design/issues/summary.md`. Tick the story's two criteria with `flai criteria tick S-0522 1` and `2` once each is verified.

This task waits for T-2237, T-2238, and T-2239: the issue is closed only once the ADR records the fix and the two code changes pass their tests.

## Done when

`flai check --strict` reports no `item.chronology`, I-0183 reads `status: closed` with a reason that names the ADR and the two code changes, `design/issues/summary.md` no longer lists it, both of S-0522's criteria are ticked, and `flai check --strict` and the markdown lint are clean on the files changed.

## Notes

`design/issues` is in the manifest's shared paths, so these touches hold no ready story (ADR-0096). The verifier's review of the diff against the criteria and the close-out follow this task, as `work-management.md` says.
