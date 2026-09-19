---
id: T-0190
type: task
nature: feature
title: "Verify: all tiers, and a review, a send back, and an acceptance through a real container"
status: done
parent: S-0041
owner: alex
created: 2026-09-19T04:08:29Z
updated: 2026-09-19T04:23:16Z
transitions:
  - to: ready
    at: 2026-09-19T04:16:23Z
    by: system-flow
  - to: in-progress
    at: 2026-09-19T04:16:23Z
    by: system-flow
  - to: done
    at: 2026-09-19T04:23:16Z
    by: system-flow
stream: S-0041
tags: []
---

# T-0190 Verify: all tiers, and a review, a send back, and an acceptance through a real container

## Work
Run `make flai-test`, `make flaiover-test`, and a production build. Build a local image from the branch and, against a scratch git project on its own container name and port: open the review page's data for a story in review on a branch (diff, plan), send it back with a reason and see it in progress with the reason in its notes, move it to review again, accept it through the streaming endpoint and see the progress lines, the result, the merge, and `by` equal to the owner; make an acceptance fail (an unchecked criterion) and see flai's message verbatim with the story still in review. Look at the review page in a browser at desktop and phone widths. Pull the released image before stopping anything when the operator's dashboard is next restarted; do not touch it for this verification. Tick the story's criteria for what was observed.

## Done when
- All tiers pass, with results in the narrative
- The container run is recorded with its outcome
- Every criterion on S-0041 is checked, or unchecked with the reason in the story notes

## Notes
