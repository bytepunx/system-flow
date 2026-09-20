---
id: T-0275
type: task
nature: feature
title: flai push and flai accept send tags three at a time, the branch last (I-0026)
status: done
parent: S-0076
owner: alex
created: 2026-09-20T12:16:18Z
updated: 2026-09-20T12:18:01Z
transitions:
  - to: ready
    at: 2026-09-20T12:16:25Z
    by: system-flow
  - to: in-progress
    at: 2026-09-20T12:16:25Z
    by: system-flow
  - to: done
    at: 2026-09-20T12:18:01Z
    by: system-flow
stream: S-0076
tags: []
---
# T-0275 flai push and flai accept send tags three at a time, the branch last (I-0026)

## Work
GitHub creates no tag events when one push carries more than three tags, so release flai never ran for 1.5.1 to 1.5.3. pending gains the batches a push is made of: tags three at a time, the branch with the last batch so an interrupted push is still detected as pending. flai push --pending and flai accept both use it.

## Done when
- A test with seven tags shows three pushes, the branch in the last
- The dry run and the messages are unchanged for three tags or fewer
- make flai-test passes

## Notes
Found while releasing S-0073 to S-0075; small, and the next stacked acceptances would hit it again.
