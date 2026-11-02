---
id: T-2424
type: task
nature: remediation
title: The preflight's check part leaves out the files a landing leaves, and the real-git fixtures give the check a clean main
status: done
parent: S-0549
owner: alex
created: 2026-11-02T01:07:53Z
updated: 2026-11-02T01:08:22Z
transitions:
  - to: ready
    at: 2026-11-02T01:08:03Z
    by: agent-S-0549
  - to: in-progress
    at: 2026-11-02T01:08:04Z
    by: agent-S-0549
  - to: done
    at: 2026-11-02T01:08:22Z
    by: agent-S-0549
stream: S-0549
tags: []
touches: [flai/internal/mergequeue/preflight.go, flai/internal/mergequeue/preflight_test.go, flai/internal/mergequeue/batch_test.go, flai/internal/serve/mergequeue_test.go, design/system/flai-cli.md, docs/users/flai.md]
usage:
  source: log
  seconds: 18
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 5
      output: 2062
      cache_read: 410394
      cache_write: 8885
      cost: 0.1772
---
# T-2424 The preflight's check part leaves out the files a landing leaves, and the real-git fixtures give the check a clean main

## Work

The integration tier, which the close-out leaves to the merge queue, failed 42 real-git tests: their fixtures had no `design` folder and stories in review with no narrative, so `flai check`'s errors on main held the queue in the preflight. And `TestUncommittedFilesNoStoryChangesStopNoLanding` showed the check part holding on an uncommitted, malformed file outside `wip/`, which a landing leaves (ADR-0178) and the batch never sees. Give the fixtures a `design` folder and a narrative per story in review; leave out of the check part the findings on the files a landing leaves; say so in flai-cli.md and the user docs.

## Done when

The integration tier passes, and the check part holds nothing on an uncommitted file outside `wip/`, with a test. Story criterion 1.

## Notes

Found by running the integration tier before the close-out.
