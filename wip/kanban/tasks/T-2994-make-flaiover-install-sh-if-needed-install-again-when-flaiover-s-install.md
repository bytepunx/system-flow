---
id: T-2994
type: task
nature: remediation
title: Make flaiover-install.sh --if-needed install again when flaiover's install configuration is newer than the last install
status: done
parent: S-0706
owner: alex
created: 2026-11-10T10:30:38Z
updated: 2026-11-11T09:20:35Z
transitions:
  - to: ready
    at: 2026-11-11T09:18:57Z
    by: agent-S-0706
  - to: in-progress
    at: 2026-11-11T09:18:57Z
    by: agent-S-0706
  - to: done
    at: 2026-11-11T09:20:35Z
    by: agent-S-0706
stream: S-0706
tags: [scripts, flaiover]
touches: [scripts/flaiover-install.sh, scripts/README.md, system-flow.yaml]
usage:
  source: log
  seconds: 98
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 10
      output: 2485
      cache_read: 1002349
      cache_write: 22879
      cost: 0.4332
---
# T-2994 Make flaiover-install.sh --if-needed install again when flaiover's install configuration is newer than the last install

## Work

It waits for no task: it is the first layer, beside T-2993, and changes other files.

`scripts/flaiover-install.sh --if-needed` skips the install while `flaiover/pnpm-lock.yaml` is not newer than `node_modules/.modules.yaml`. A change to the install configuration therefore never reaches a `node_modules` already installed for one libc: not T-2993's `supportedArchitectures`, nor anything else in `package.json`, `pnpm-workspace.yaml`, or `.npmrc`. The main checkout and every open worktree would keep their single-libc install.

Make the check also install when any of `flaiover/package.json`, `flaiover/pnpm-workspace.yaml`, or `flaiover/.npmrc` is newer than the stamp. Keep the stamp and the `touch` after the install as they are. Say in `scripts/README.md` what `--if-needed` now compares.

flai accept runs flaiover's `install` only when the accepted change touches its `paths` in `system-flow.yaml`, which name only `package.json` and `pnpm-lock.yaml`. Add `pnpm-workspace.yaml` and `.npmrc` there, so that an acceptance that changes only the configuration installs again in the main checkout. The `install` key is written by hand: `flai subproject` keeps it and does not set it.

## Done when

- `scripts/flaiover-install.sh --if-needed` installs when the lockfile or any of the three configuration files is newer than `node_modules/.modules.yaml`, and skips otherwise.
- flaiover's `install` `paths` in `system-flow.yaml` name `package.json`, `pnpm-lock.yaml`, `pnpm-workspace.yaml`, and `.npmrc`, and `flai check --strict` reports nothing on the manifest.
- `scripts/README.md` describes the rule.
- `flai test scripts/flaiover-install.sh` passes (ShellCheck included).

## Notes
