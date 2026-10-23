---
id: T-2070
type: task
nature: feature
title: flaiover-install.sh purges an install from another store instead of aborting, and flaiover-pnpm.sh adds a dependency with the project's pnpm and store
status: done
parent: S-0481
owner: alex
created: 2026-10-23T07:17:37Z
updated: 2026-10-23T07:19:56Z
transitions:
  - to: ready
    at: 2026-10-23T07:17:54Z
    by: agent-S-0481
  - to: in-progress
    at: 2026-10-23T07:17:55Z
    by: agent-S-0481
  - to: done
    at: 2026-10-23T07:19:20Z
    by: agent-S-0481
stream: S-0481
tags: []
touches: [scripts/flaiover-install.sh, scripts/flaiover-pnpm.sh, scripts/flaiover-install-test.sh, scripts/flaiover-test.sh, scripts/README.md, design/tech/node.md]
usage:
  source: log
  seconds: 85
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 19
      output: 5593
      cache_read: 1245954
      cache_write: 39042
      cost: 0.6735
---

# T-2070 flaiover-install.sh purges an install from another store instead of aborting, and flaiover-pnpm.sh adds a dependency with the project's pnpm and store

## Work

Waits for no task: it is the first layer.

- `scripts/flaiover-install.sh` runs `pnpm install` with `--config.confirm-modules-purge=false`, so that a `node_modules` linked from another store, such as a bare `pnpm add` made from the default store, is removed and installed again from the project's store instead of aborting with `ERR_PNPM_ABORTED_REMOVE_MODULES_DIR_NO_TTY`. `node_modules` is git-ignored and rebuilt from the lockfile, so the purge loses nothing.
- `scripts/flaiover-pnpm.sh` runs the pinned pnpm in `flaiover/` with `env.sh`'s store and caches, such as `scripts/flaiover-pnpm.sh add <package>`, so that a dependency is added without the default store or the host's other pnpm.
- `scripts/flaiover-install-test.sh` reproduces I-0177 in a scratch copy: an install from another store, then `flaiover-install.sh --if-needed`, which must pass and leave the install on the project's store. `scripts/flaiover-test.sh` runs it first.
- `scripts/README.md`, `design/tech/node.md`, and `design/system/devex.md` say so.

## Done when

- The test fails against `main`'s `flaiover-install.sh` and passes with this one.
- `flai test` passes on the changed paths.

## Notes
