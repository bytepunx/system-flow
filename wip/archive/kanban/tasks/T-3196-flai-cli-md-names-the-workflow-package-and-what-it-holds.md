---
id: T-3196
type: task
nature: feature
title: flai-cli.md names the workflow package and what it holds
status: done
parent: S-0737
owner: alex
created: 2026-11-12T16:05:54Z
updated: 2026-11-14T23:23:35Z
transitions:
  - to: ready
    at: 2026-11-14T23:21:38Z
    by: system-flow
  - to: in-progress
    at: 2026-11-14T23:21:38Z
    by: system-flow
  - to: done
    at: 2026-11-14T23:23:35Z
    by: agent-S-0737
stream: S-0737
tags: [cli]
touches: [design/system/flai-cli.md]
after: [T-3188]
usage:
  source: log
  seconds: 117
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 9
      output: 4275
      cache_read: 1040753
      cache_write: 23948
      cost: 0.4698
---
# T-3196 flai-cli.md names the workflow package and what it holds

## Work

Waits for T-3188: it describes the package that task adds. Add `flai/internal/workflow` to `design/system/flai-cli.md` § Internal structure: what it parses, the baseline it embeds, the lookups it gives, and that `workitem` reads its states and transitions from it.

## Done when

§ Internal structure names the package, its baseline, and its readers, and links the ADR.

## Notes
