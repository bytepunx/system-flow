---
id: T-2079
type: task
nature: improvement
title: An ADR refining ADR-0156 records that the patch id leaves out the generated paths, and flai-cli.md and workflow.md say so
status: done
parent: S-0487
owner: alex
created: 2026-10-23T18:05:03Z
updated: 2026-10-23T18:26:59Z
transitions:
  - to: ready
    at: 2026-10-23T18:26:26Z
    by: agent-S-0487
  - to: in-progress
    at: 2026-10-23T18:26:27Z
    by: agent-S-0487
  - to: done
    at: 2026-10-23T18:26:59Z
    by: agent-S-0487
stream: S-0487
tags: [cli, design]
touches: [design/adrs/README.md, design/system/flai-cli.md, design/system/workflow.md, design/adrs/0192-flai-verify-s-patch-id-leaves-out-the-files-a-sync-s-rebase-writes-again-so-a.md]
after: [T-2078]
usage:
  source: log
  seconds: 32
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 6
      output: 1448
      cache_read: 323224
      cache_write: 13424
      cost: 0.1804
---
# T-2079 An ADR refining ADR-0156 records that the patch id leaves out the generated paths, and flai-cli.md and workflow.md say so

## Work

Write the ADR with `flai adr new`, refining ADR-0156 and ADR-0098: the patch id `flai verify` records for reuse across a clean rebase is `git patch-id --stable` of the branch's diff from the merge base without the generated paths, the files flai regenerates at sync (`issues.Generated`, `design/issues/summary.md` alone today). Give it I-0196's instance as context, the reason the exclusion is safe (the generated content derives from the issue files, which stay in the diff, and from main's, which select tiers through main's new paths), and the alternatives: leaving the summary out of the branch's commits, which the dashboard, the review page, and the prime pack read (ADR-0098's alternatives); and comparing the change path by path, which ADR-0156 rejected. Accept it with `flai adr accept`, which writes its row in `design/adrs/README.md`.

Then bring the design in step, in the same commit: the `flai verify` row in `design/system/flai-cli.md` (its `patch` field and reuse across a clean rebase) and the sentence in `design/system/workflow.md` § Branches and collisions that says what `flai verify` reuses across the clean rebase, each naming the new ADR.

Waits for T-2078, so that the ADR records what was built, with the edge cases its tests found.

## Done when

The ADR is accepted and indexed, `design/system/flai-cli.md` and `design/system/workflow.md` name it where they describe the patch id, and `flai check --strict` and `flai test design/adrs design/system` pass.

## Notes

`work-management.md`'s close-out rule says only that the tiers are reused "while the branch's own change stays the same", which stays true: the convention is not changed.
