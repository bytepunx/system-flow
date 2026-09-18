---
id: T-0158
type: task
nature: remediation
title: "Verify: all tiers, and acceptance of a story with a worktree through a real container"
status: done
parent: S-0050
owner: alex
created: 2026-09-18T19:47:24Z
updated: 2026-09-18T20:18:53Z
transitions:
  - to: ready
    at: 2026-09-18T19:57:54Z
    by: alex
  - to: in-progress
    at: 2026-09-18T19:57:54Z
    by: alex
  - to: done
    at: 2026-09-18T20:18:53Z
    by: alex
stream: S-0050
tags: []
---

# T-0158 Verify: all tiers, and acceptance of a story with a worktree through a real container

## Work
Run `make flai-test` and `make flaiover-test`. Then prove the fix against a real container without touching the operator's running dashboard: render a scratch project with git, open a story with a worktree, commit on the branch, start a second dashboard container from a locally built image on another port and name, and accept the story through the dashboard's API; check the merge, the move to done, and that `git worktree list` inside the container shows nothing prunable. Stop and remove the scratch container. Tick the story's criteria only for what was observed; the opt-in's real-git case cannot run on this host's git 2.47.3 and is verified inside the container's git or left unchecked with the reason.

## Done when
- All tiers and lint pass, with results in the narrative
- The container acceptance ran and its outcome is recorded
- Every criterion on S-0050 is checked, or unchecked with the reason in the story notes

## Notes
