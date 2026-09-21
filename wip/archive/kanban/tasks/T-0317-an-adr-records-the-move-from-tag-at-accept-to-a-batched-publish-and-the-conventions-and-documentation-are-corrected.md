---
id: T-0317
type: task
nature: feature
title: An ADR records the move from tag-at-accept to a batched publish, and the conventions and documentation are corrected
status: done
parent: S-0087
owner: alex
created: 2026-09-21T03:32:55Z
updated: 2026-09-21T04:19:13Z
transitions:
  - to: ready
    at: 2026-09-21T04:09:18Z
    by: system-flow
  - to: in-progress
    at: 2026-09-21T04:09:19Z
    by: system-flow
  - to: done
    at: 2026-09-21T04:19:13Z
    by: system-flow
stream: S-0087
tags: []
---
# T-0317 An ADR records the move from tag-at-accept to a batched publish, and the conventions and documentation are corrected

## Work
An ADR records that tagging and pushing move from accept-time to a deliberate publish, batched per component by the highest delivery type pending, refining ADR-0019 (which established tag-at-accept as the baseline) and cross-referencing ADR-0025 (research/experiment are unchanged) and ADR-0031/S-0078 (the push host action gates publish exactly as it gated the old automatic push). This story's whole point is that "tag on acceptance" (`design/conventions/git.md`'s baseline "Rules" list) stops being true, so the convention text changes with it: `template/root/design/conventions/git.md` first, then `design/conventions/git.md` identically, in this same story, per the layout's own rule for how conventions change. Its stale note that `flai release` "computes the bump, tags, and pushes once it exists" is corrected too: it exists, and after this story it pushes. `docs/operators` and any user-facing documentation describing acceptance and releases are updated to describe merge-then-publish and the done column's two states.

## Done when
- The ADR is accepted, refines and cross-references what is listed above
- The baseline "Tag on acceptance" rule and the stale `flai release` note are corrected in the template first, then identically in this repository's own copy
- flai check --strict and the markdown lint pass

## Notes
Written once T1 through T4 exist, so the ADR describes what was actually built rather than the plan.
