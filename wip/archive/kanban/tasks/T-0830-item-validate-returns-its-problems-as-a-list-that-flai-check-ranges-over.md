---
id: T-0830
type: task
nature: feature
title: Item.Validate returns its problems as a list that flai check ranges over
status: done
parent: S-0248
owner: alex
created: 2026-10-04T23:14:45Z
updated: 2026-10-04T23:16:58Z
transitions:
  - to: ready
    at: 2026-10-04T23:14:58Z
    by: agent-S-0248
  - to: in-progress
    at: 2026-10-04T23:14:58Z
    by: agent-S-0248
  - to: done
    at: 2026-10-04T23:16:58Z
    by: agent-S-0248
stream: S-0248
tags: []
touches: [flai/internal/workitem, flai/internal/manifest, flai/internal/check, design/system/flai-cli.md]
usage:
  source: log
  seconds: 120
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 37
      output: 7919
      cache_read: 1877012
      cache_write: 50873
      cost: 0.9409
---

# T-0830 Item.Validate returns its problems as a list that flai check ranges over

## Work

`Item.Validate` returns a `workitem.Problems` error, one phrase per problem, and `check.oneItem` ranges over it instead of splitting the joined message on "; ". The agent block's problems join the list one by one (`manifest.Agent.Problems`). It waits for nothing: it is the story's only task.

## Done when

- A front-matter message that holds "; ", such as an agent role that sets nothing, is one `item.front-matter` finding on its key's line, with a test that fails without the change.
- Two agent problems are still two findings.
- `design/system/flai-cli.md` says check reports one finding per problem.

## Notes
