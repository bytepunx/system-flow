---
id: T-2997
type: task
nature: remediation
title: "Reproduce I-0259 in flaiover-install-test.sh: one install holds both libc variants, and a newer configuration installs again"
status: backlog
parent: S-0706
owner: alex
created: 2026-11-10T10:30:52Z
updated: 2026-11-11T08:54:25Z
transitions: []
stream: S-0706
tags: [scripts, flaiover]
touches: [scripts/flaiover-install-test.sh, scripts/README.md]
after: [T-2993, T-2994]
---
# T-2997 Reproduce I-0259 in flaiover-install-test.sh: one install holds both libc variants, and a newer configuration installs again

## Work

It waits for T-2993, whose `supportedArchitectures` setting the test copies, and for T-2994, whose `--if-needed` rule it tests and whose `scripts/README.md` change it builds on.

`scripts/flaiover-install-test.sh` already reproduces I-0177 offline with local-folder packages. Add I-0259's case the same way, with no network. Use two local optional packages, one with `"libc": ["glibc"]` and one with `"libc": ["musl"]`, depended on by a fixture project. First install without `supportedArchitectures` and show that only the host's libc variant lands: that is the issue. Then install with the setting T-2993 chose, in the file it chose, and assert that both land. Also assert that `flaiover-install.sh --if-needed` installs again once the fixture's configuration file is newer than `node_modules/.modules.yaml`, and skips while it is not.

If pnpm cannot be made to resolve libc-specific optional packages offline, keep the `--if-needed` case. Then also assert that flaiover's own configuration declares both libcs. Say in this task's Notes why the fixture case could not run.

Add I-0259's case to the test's row in `scripts/README.md`.

## Done when

- `scripts/flaiover-install-test.sh` fails on the code before T-2993 and T-2994 and passes after them, offline.
- `scripts/README.md`'s row for the test names I-0259's case.
- `flai test scripts/flaiover-install-test.sh` passes, and so does `scripts/flaiover-test.sh`, which runs the test.

## Notes
