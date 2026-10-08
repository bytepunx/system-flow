---
id: T-1342
type: task
nature: remediation
title: Close I-0086 saying the installer and self-upgrade retry a dropped GitHub connection
status: backlog
parent: S-0291
owner: alex
created: 2026-10-08T08:00:17Z
updated: 2026-10-08T08:00:28Z
transitions: []
stream: S-0291
tags: [issues]
touches: [design/issues/I-0086-the-close-out-s-install-smoke-test-failed-once-and-passed-when-run-alone-with-no-cause-in-its-output.md, design/issues/summary.md]
after: [T-1338, T-1339, T-1340]
---
# T-1342 Close I-0086 saying the installer and self-upgrade retry a dropped GitHub connection

## Work

Run `flai issue close I-0086 --reason "<reason>"` in the story's worktree, since `flai issue` writes under the checkout it runs in. The reason names the cause: GitHub dropping the large release listing, with curl (92) or (56). It names the fix: `install.sh` and `flai self-upgrade` retry a call the network drops, `install.sh` resolves the latest release from a small page, and `scripts/install-test.sh` prints `install.sh`'s output when it fails. It names the tests that reproduce the drop, `flai/cmd/installsh_test.go` and the new test in `flai/internal/selfupgrade/selfupgrade_test.go`. The command regenerates `design/issues/summary.md`.

Waits for T-1338, T-1339, and T-1340: the reason names what they built and the tests they added.

## Done when

- I-0086 is closed with that reason, and `design/issues/summary.md` no longer lists it.
- `flai check --strict` is clean on both files.

## Notes
