---
id: T-1304
type: task
nature: improvement
title: scripts/flaiover-lint.sh runs prettier --check and eslint on the flaiover files it is given
status: done
parent: S-0319
owner: alex
created: 2026-10-08T00:14:59Z
updated: 2026-10-08T05:53:15Z
transitions:
  - to: ready
    at: 2026-10-08T05:52:39Z
    by: agent-S-0319
  - to: in-progress
    at: 2026-10-08T05:52:39Z
    by: agent-S-0319
  - to: done
    at: 2026-10-08T05:53:15Z
    by: agent-S-0319
stream: S-0319
tags: [flaiover, testing, scripts]
touches: [scripts/flaiover-lint.sh]
usage:
  source: log
  seconds: 36
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 6
      output: 20
      cache_read: 386117
      cache_write: 6796
      cost: 0.1751
---
# T-1304 scripts/flaiover-lint.sh runs prettier --check and eslint on the flaiover files it is given

## Work

Write `scripts/flaiover-lint.sh`, POSIX `sh` under `set -eu`, modelled on `scripts/flaiover-unit.sh`. It takes the files a `flai test` tier selected, relative to `flaiover/` (the tier's `dir`). It sources `scripts/env.sh`. In a story worktree it first runs `scripts/flaiover-install.sh --if-needed`; in the main checkout it exits 0 while `flaiover/node_modules` is missing, as `flaiover-unit.sh` does. It then runs, in `flaiover/`:

- `pnpm exec prettier --check --ignore-unknown <files>`, so that a file prettier does not format, or one `.prettierignore` lists, is passed over rather than failed;
- `pnpm exec eslint --no-warn-ignored <the .js, .ts, and .svelte files among them>`, skipped when there are none.

It exits non-zero when either finds a fault, with the tool's own output, which names each file. It waits for no other task: it is the first layer.

## Done when

- [ ] `scripts/flaiover-lint.sh src/routes/+page.svelte`, run from `flaiover/`, exits 0 on a clean file
- [ ] The same run exits non-zero and names the file when a prettier fault is planted in it; the fault is reverted after
- [ ] A file eslint or prettier ignores, such as `package.json` or a file under `static/`, does not make it fail

## Notes

Drafted by the planner for S-0319. svelte-check is left out on purpose: it checks the whole project and takes no file list, so it stays in the `all_only` `flaiover` tier (`scripts/flaiover-test.sh`).
