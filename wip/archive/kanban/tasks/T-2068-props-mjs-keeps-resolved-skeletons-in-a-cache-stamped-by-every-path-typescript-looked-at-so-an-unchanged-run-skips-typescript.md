---
id: T-2068
type: task
nature: improvement
title: props.mjs keeps resolved skeletons in a cache stamped by every path typescript looked at, so an unchanged run skips typescript
status: done
parent: S-0478
owner: alex
created: 2026-10-23T06:52:13Z
updated: 2026-10-23T06:53:26Z
transitions:
  - to: ready
    at: 2026-10-23T06:52:50Z
    by: agent-S-0478
  - to: in-progress
    at: 2026-10-23T06:52:51Z
    by: agent-S-0478
  - to: done
    at: 2026-10-23T06:53:26Z
    by: agent-S-0478
stream: S-0478
tags: []
touches: [flai/internal/components/props.mjs, flai/internal/components/props.go, flai/internal/components/props_test.go, design/system/component-library.md, docs/operators/runbooks/backup.md]
usage:
  source: log
  seconds: 35
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 9
      output: 3723
      cache_read: 730864
      cache_write: 16351
      cost: 0.3515
---
# T-2068 props.mjs keeps resolved skeletons in a cache stamped by every path typescript looked at, so an unchanged run skips typescript

## Work

Waits for nothing: the first layer. props.mjs keys a cache file by the script's hash, the sub-project, where its typescript resolves, and each typed component's typing and json props, and stamps every path typescript's system looks at while it resolves the skeletons (a missing file by its folder's listing). A run whose key has a file and whose stamps all hold takes the skeletons from it without loading typescript. flai passes the cache folder, `.flai-cache/component-props/<sub-project path>/`, and the script's hash. The folder keeps the four files used last, so the components list and the examples check do not evict each other.

## Done when

- A second run over flaiover with nothing changed answers the same props as the first in about 0.45s, against 2.7s before (measured by hand on the main checkout's flaiover).
- Integration tests show a run with nothing changed answering from the cache, a change to an imported type's file resolving again, a module made after a run that missed it resolving, and two sets of components keeping a file each.
- `design/system/component-library.md` describes the cache, and the backup runbook lists it.

## Notes
