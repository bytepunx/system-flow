---
id: T-3047
type: task
nature: remediation
title: The candidate tests pass repeatedly under load, and I-0267 is closed with what fixed it
status: done
parent: S-0703
owner: alex
created: 2026-11-10T10:34:54Z
updated: 2026-11-11T23:19:39Z
transitions:
  - to: ready
    at: 2026-11-11T23:18:08Z
    by: agent-S-0703
  - to: in-progress
    at: 2026-11-11T23:18:08Z
    by: agent-S-0703
  - to: done
    at: 2026-11-11T23:19:39Z
    by: agent-S-0703
stream: S-0703
tags: [flai, tests, issues]
touches: [design/issues/I-0267-testarepositoryisofferedimportedandthenserved-fails-the-close-out-s-integration.md, design/issues/summary.md]
after: [T-3042, T-3044]
usage:
  source: log
  seconds: 89
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 4
      output: 32
      cache_read: 464715
      cache_write: 2652
      cost: 0.1735
---
# T-3047 The candidate tests pass repeatedly under load, and I-0267 is closed with what fixed it

## Work

Run the candidate tests of `flai/internal/serve` under load: `go test -race -count=20 -run 'TestARepositoryIsOffered|TestFindCandidates'` with `GOMAXPROCS=2`, while another package's tests load the host's CPU, then the whole `flai/internal/serve` package once. Record the runs and their results in the narrative's log. Then close I-0267 with `flai issue close I-0267 --reason`, naming S-0560's drain and this story's change to the test dashboard, which regenerates `design/issues/summary.md`.

Waits on T-3044: the runs test its fix, and the reason names it.

## Done when

- Every run passes, and I-0267 is closed with a reason that names what fixed it.

## Notes
