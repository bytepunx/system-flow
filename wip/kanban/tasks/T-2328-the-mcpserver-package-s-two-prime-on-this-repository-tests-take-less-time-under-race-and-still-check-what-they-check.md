---
id: T-2328
type: task
nature: remediation
title: The mcpserver package's two prime-on-this-repository tests take less time under -race and still check what they check
status: backlog
parent: S-0540
owner: alex
created: 2026-10-29T07:37:19Z
updated: 2026-10-29T07:37:34Z
transitions: []
stream: S-0540
tags: [flai, testing]
touches: [flai/internal/mcpserver/prime_test.go]
---
# T-2328 The mcpserver package's two prime-on-this-repository tests take less time under -race and still check what they check

## Work

Without `-race`, `TestPrimeOnThisRepository` takes about 33s and `TestPrimeComesInPartsOnThisRepository` about 31s, half of the package's 130s. `TestPrimeOnThisRepository` primes every open story of the repository one after another, so its cost grows with the backlog. Both stay under `-race`: the integration tier runs the package only with `-race` too, so they would lose their only run otherwise.

- Measure first. Run `go test -json -race -count=1 -skip '^TestGateway' ./internal/mcpserver` in `flai` and take each test's `Elapsed` from the `pass` events. Note the five slowest tests, the package's total, and the host's load average.
- Make the two tests cheaper without narrowing what they assert. For `TestPrimeOnThisRepository`, prime each story in a parallel subtest named for its ID, against the one fixture `thisRepository` gives, if that fixture is safe to share; otherwise give each subtest its own. For `TestPrimeComesInPartsOnThisRepository`, run each case as a parallel subtest. If a test the measurement shows slower than these two is cheap to fix the same way, fix it too, and add its file to this task's touches first.
- Measure again with the same command on the same host, and record both measurements in the narrative's `## Decisions`.

It waits for no task: no other task of this story touches `prime_test.go`. S-0536, in progress, claims `flai/internal/mcpserver`, and S-0537 touches `prime_test.go`. Agree on the file in their conversation before you change it.

## Done when

- Both tests still fail on the cause each guards: a story whose pack the schema refuses, and a part over `ctxpack.PartLimit`.
- The two tests' combined `Elapsed` under `-race` is at most half the first measurement's, on the same host under similar load.
- Both measurements, the five slowest tests and the package's total each time, are in the narrative's `## Decisions`.
- `flai test flai/internal/mcpserver` passes.

## Notes
