---
id: T-0637
type: task
nature: research
title: Find two flai serve runs made with a released flai that includes S-0175, each delegating
status: in-progress
parent: S-0188
owner: arobson
created: 2026-10-01T08:31:38Z
updated: 2026-10-01T08:31:52Z
transitions:
  - to: ready
    at: 2026-10-01T08:31:52Z
    by: agent-S-0188
  - to: in-progress
    at: 2026-10-01T08:31:52Z
    by: agent-S-0188
stream: S-0188
tags: []
touches: [design/system/agent-context.md]
---
# T-0637 Find two flai serve runs made with a released flai that includes S-0175, each delegating

## Work

Confirm the installed flai includes S-0175 (`flai version` against the release tag that carries it), then pick two story runs in `~/.flai/serve/agents/` that started after it was installed and whose logs show an `Agent` call to the explorer or the verifier.

## Done when

Two runs are named, with the flai version each ran under and the sub-agents each started.

## Notes
