---
id: T-1339
type: task
nature: remediation
title: flai self-upgrade retries a release listing or download the network drops
status: backlog
parent: S-0291
owner: alex
created: 2026-10-08T07:59:57Z
updated: 2026-10-08T07:59:57Z
transitions: []
stream: S-0291
tags: [self-upgrade, network]
touches: [flai/internal/selfupgrade/selfupgrade.go, flai/internal/selfupgrade/selfupgrade_test.go]
---
# T-1339 flai self-upgrade retries a release listing or download the network drops

## Work

The smoke tier's `scripts/install-test.sh` runs `flai self-upgrade --check` and `--dir` four times against GitHub after `install.sh`. `selfupgrade.List` pages `releases?per_page=100`, a larger body than `install.sh`'s, through `Options.getPage`, and neither it nor `getBytes` retries. The connection drop I-0086 records would fail those steps the same way. S-0283's instance stopped with self-upgrade reaching GitHub and no reason printed.

- Retry a request in `getPage` and `getBytes` up to three times with a short pause when the transport fails, or when reading the body fails, such as `unexpected EOF` or an HTTP/2 stream reset. Do not retry a non-2xx answer: a 401 or 404 keeps its "private repository?" hint and fails at once. Honour the context's cancellation between attempts.
- Keep the error of the last attempt, saying how many attempts were made.
- In `selfupgrade_test.go`, add a test whose `httptest` server cuts the first releases page off mid-body by hijacking the connection, and check that `Resolve` still answers the newest release. Add a test that a 404 is not retried. Make the pause a field or variable the test sets to zero, so the tests stay fast.

Waits for nothing: its paths are its own, and it shares no behaviour with `install.sh` beyond the fault.

## Done when

- The new test fails against the current `selfupgrade.go` and passes with the change.
- A 404 is asked for once, which the test checks by counting requests.
- `flai test flai/internal/selfupgrade` passes.

## Notes
