---
id: T-1340
type: task
nature: remediation
title: scripts/install-test.sh prints install.sh's output when the default-path install fails
status: backlog
parent: S-0291
owner: alex
created: 2026-10-08T08:00:03Z
updated: 2026-10-08T08:00:03Z
transitions: []
stream: S-0291
tags: [smoke, install]
touches: [scripts/install-test.sh]
---
# T-1340 scripts/install-test.sh prints install.sh's output when the default-path install fails

## Work

In the default-path case, `scripts/install-test.sh` sends `install.sh`'s output to `.flai-cache/install-test-home.log`. Under `set -e`, a failing `install.sh` ends the script before any `cat "$OUT"`, so the tier prints `exit status 1` with no cause. S-0318, S-0319, and S-0336 each had to find the curl error in that log by hand.

- Run that `install.sh` so that its exit code is caught, such as `if ! env -i ... sh "$ROOT/install.sh" >"$OUT" 2>&1; then ...`. On failure, print a line naming the step, then the log to stderr, then exit 1.
- Leave the other steps as they are. Their output already reaches the tier.
- Keep the script POSIX `sh` under `set -eu`, as `tooling.md` says.

Waits for nothing: no other task changes this script.

## Done when

- With `install.sh` made to fail, for example with `FLAI_API` set to an address that refuses connections in a scratch run, the default-path step prints `install.sh`'s own error, such as "could not list releases", and the script exits 1.
- `flai test scripts/install-test.sh` passes.

## Notes
