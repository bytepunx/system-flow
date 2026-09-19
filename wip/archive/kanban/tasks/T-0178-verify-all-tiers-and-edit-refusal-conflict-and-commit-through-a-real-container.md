---
id: T-0178
type: task
nature: feature
title: "Verify: all tiers, and edit, refusal, conflict, and commit through a real container"
status: done
parent: S-0040
owner: alex
created: 2026-09-19T02:11:30Z
updated: 2026-09-19T02:32:09Z
transitions:
  - to: ready
    at: 2026-09-19T02:27:42Z
    by: alex
  - to: in-progress
    at: 2026-09-19T02:27:42Z
    by: alex
  - to: done
    at: 2026-09-19T02:32:09Z
    by: alex
stream: S-0040
tags: []
---

# T-0178 Verify: all tiers, and edit, refusal, conflict, and commit through a real container

## Work
Run `make flai-test` and `make flaiover-test` and a production build. Build a local image from the branch and, against a scratch git project on its own container name and port: edit a design document and see the commit on main with the scratch identity as author and the trailer; save a change that breaks the check and see it refused with the file restored; change the file on the host between load and save and see the conflict with its diff; edit a story body and see the front matter untouched; set `dashboard.autocommit: false` and see the change left uncommitted. Look at the editor in a browser at desktop and phone widths in both themes. Stop the container and remove the image; the operator's dashboard is not touched. Tick the story's criteria for what was observed.

## Done when
- All tiers pass, with results in the narrative
- The container run is recorded with its outcome
- Every criterion on S-0040 is checked, or unchecked with the reason in the story notes

## Notes
