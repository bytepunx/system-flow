---
id: T-0167
type: task
nature: improvement
title: "Verify: all tiers and the choice through a real container"
status: done
parent: S-0051
owner: alex
created: 2026-09-18T21:00:24Z
updated: 2026-09-18T21:09:26Z
transitions:
  - to: ready
    at: 2026-09-18T21:06:32Z
    by: alex
  - to: in-progress
    at: 2026-09-18T21:06:32Z
    by: alex
  - to: done
    at: 2026-09-18T21:09:26Z
    by: alex
stream: S-0051
tags: []
---

# T-0167 Verify: all tiers and the choice through a real container

## Work
Run `make flai-test` and `make flaiover-test`. Build a local image from the branch and, against a scratch git project on its own container name and port, with a stray file outside `wip/`: the preview lists it; the move without the choice is refused and changes nothing; the move with the choice accepts and the file is in the acceptance commit. Then with a file ignored only by a global excludes file: the container's `git status` does not report it and the preview lists nothing. Stop the scratch container and remove the image. The operator's dashboard is not touched. Tick the story's criteria for what was observed.

## Done when
- All tiers pass, with results in the narrative
- The container run is recorded with its outcome
- Every criterion on S-0051 is checked, or unchecked with the reason in the story notes

## Notes
