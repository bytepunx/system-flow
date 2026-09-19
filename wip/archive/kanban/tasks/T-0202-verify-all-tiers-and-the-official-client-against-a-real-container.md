---
id: T-0202
type: task
nature: feature
title: "Verify: all tiers, and the official client against a real container"
status: done
parent: S-0043
owner: alex
created: 2026-09-19T05:17:08Z
updated: 2026-09-19T05:26:51Z
transitions:
  - to: ready
    at: 2026-09-19T05:24:48Z
    by: system-flow
  - to: in-progress
    at: 2026-09-19T05:24:48Z
    by: system-flow
  - to: done
    at: 2026-09-19T05:26:51Z
    by: system-flow
stream: S-0043
tags: []
---

# T-0202 Verify: all tiers, and the official client against a real container

## Work
Run `make flai-test`, `make flaiover-test`, and a production build. Build a local image from the branch and, against a scratch project on its own container name and port: run the Go SDK client test against `/mcp` with the scratch token; check 401 without it; check the identity headers and body member on a few `/api` responses, an array one included. Stop the container and remove the image. The operator's dashboard is not touched. Tick the story's criteria for what was observed.

## Done when
- All tiers pass, with results in the narrative
- The official client ran against the container and its outcome is recorded
- Every criterion on S-0043 is checked, or unchecked with the reason in the story notes

## Notes
