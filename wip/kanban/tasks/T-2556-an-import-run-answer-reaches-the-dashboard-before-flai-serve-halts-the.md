---
id: T-2556
type: task
nature: remediation
title: An import.run answer reaches the dashboard before flai serve halts the candidate's connection, so TestARepositoryIsOfferedImportedAndThenServed gets no EOF
status: in-progress
parent: S-0560
owner: alex
created: 2026-11-07T08:19:54Z
updated: 2026-11-07T08:35:27Z
transitions:
  - to: ready
    at: 2026-11-07T08:35:26Z
    by: system-flow
  - to: in-progress
    at: 2026-11-07T08:35:27Z
    by: system-flow
stream: S-0560
tags: [flai, tests]
touches: [flai/internal/serve/candidates.go, flai/internal/serve/candidates_test.go, flai/internal/hostapi/import.go]
---
# T-2556 An import.run answer reaches the dashboard before flai serve halts the candidate's connection, so TestARepositoryIsOfferedImportedAndThenServed gets no EOF

## Work

I-0246's instance of S-0555: `TestARepositoryIsOfferedImportedAndThenServed` failed with EOF on its `import.run` call in `flai/internal/serve/candidates_test.go`. This happened once and has not been reproduced, so the cause below is a reading of the code, to be tested first.

`hostapi.ImportMethods` calls `imported(p.Root)` inside the method, before it returns. `offers.imported` registers the project and sets `f.rescan`. The next `offers.reconcile` then finds the candidate unwanted and calls `r.halt()`, which cancels the client at once. A look that lands before the answer is written loses the answer, and the dashboard reads EOF on an import that succeeded. The test looks every 20 ms, which makes this likely under load. `running.drain` already answers requests in flight for a project, and `reconcile` does not use it for a candidate.

Reproduce it with `GOMAXPROCS=2 go test -race -count=50 -run TestARepositoryIsOfferedImportedAndThenServed ./internal/serve` beside a CPU load, run from `flai/`, and record the result in the narrative. If that is the cause, make the answer leave first, by draining the candidate or by registering once the answer is written, whichever changes less. Add a test that holds the look until the answer would be lost, if a seam allows one. If 50 runs find nothing, do not change the product on the reading alone: say how many runs found nothing, and record the suspected window with `flai issue new` for I-0246's closing reason to name. This task waits for no other: it shares no path with the mcpserver tasks.

## Done when

- [ ] The narrative records the reproduction runs, the load, and what they found
- [ ] Either `import.run`'s answer is written before the candidate is halted, with the test passing under the command above, or an issue records the unreproduced window
- [ ] `flai test flai/internal/serve flai/internal/hostapi` passes

## Notes
