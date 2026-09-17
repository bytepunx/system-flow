---
id: T-099
type: task
nature: feature
title: "Local image script and Makefile target; run it with flai dashboard end to end"
status: done
parent: S-015
owner: alex
created: 2026-09-17T05:29:33Z
updated: 2026-09-17T05:33:21Z
transitions:
  - to: ready
    at: 2026-09-17T05:33:21Z
    by: agent
  - to: in-progress
    at: 2026-09-17T05:33:21Z
    by: agent
  - to: done
    at: 2026-09-17T05:33:21Z
    by: agent
stream: S-015
tags: [dashboard, docker]
---

# T-099 Local image script and Makefile target; run it with flai dashboard end to end

## Work
scripts/flaiover-image.sh builds flaiover:local from the repo root; Makefile target flaiover-image; run against this repository with flai dashboard --image flaiover --tag local; verify manifest, a board read, a write through the bundled flai, and stop.

## Done when
End to end through flai dashboard with the local image.

## Notes
