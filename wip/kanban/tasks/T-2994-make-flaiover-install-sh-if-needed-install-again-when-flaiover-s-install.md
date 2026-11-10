---
id: T-2994
type: task
nature: remediation
title: Make flaiover-install.sh --if-needed install again when flaiover's install configuration is newer than the last install
status: backlog
parent: S-0706
owner: alex
created: 2026-11-10T10:30:38Z
updated: 2026-11-10T10:30:38Z
transitions: []
stream: S-0706
tags: [scripts, flaiover]
touches: [scripts/flaiover-install.sh, scripts/README.md]
---
# T-2994 Make flaiover-install.sh --if-needed install again when flaiover's install configuration is newer than the last install

## Work

It waits for no task: it is the first layer, beside T-2993, and changes other files.

`scripts/flaiover-install.sh --if-needed` skips the install while `flaiover/pnpm-lock.yaml` is not newer than `node_modules/.modules.yaml`. A change to the install configuration therefore never reaches a `node_modules` already installed for one libc: not T-2993's `supportedArchitectures`, nor anything else in `package.json`, `pnpm-workspace.yaml`, or `.npmrc`. The main checkout and every open worktree would keep their single-libc install. flai accept's `install` would skip it too.

Make the check also install when any of `flaiover/package.json`, `flaiover/pnpm-workspace.yaml`, or `flaiover/.npmrc` is newer than the stamp. Keep the stamp and the `touch` after the install as they are. Say in `scripts/README.md` what `--if-needed` now compares.

## Done when

- `scripts/flaiover-install.sh --if-needed` installs when the lockfile or any of the three configuration files is newer than `node_modules/.modules.yaml`, and skips otherwise.
- `scripts/README.md` describes the rule.
- `flai test scripts/flaiover-install.sh` passes (ShellCheck included).

## Notes
