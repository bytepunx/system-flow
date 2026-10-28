---
id: T-2265
type: task
nature: improvement
title: flai test runs svelte-check on the flaiover project when a changed TypeScript, JavaScript, or Svelte file under flaiover selects its new tier
status: done
parent: S-0530
owner: alex
created: 2026-10-28T18:20:37Z
updated: 2026-10-28T21:02:17Z
transitions:
  - to: ready
    at: 2026-10-28T20:50:02Z
    by: agent-S-0530
  - to: in-progress
    at: 2026-10-28T20:50:03Z
    by: agent-S-0530
  - to: done
    at: 2026-10-28T21:02:17Z
    by: agent-S-0530
stream: S-0530
tags: [flaiover, testing, scripts]
touches: [system-flow.yaml, scripts/flaiover-check.sh, scripts/flaiover-lint.sh, scripts/README.md, flai/internal/verify/select_test.go]
usage:
  source: log
  seconds: 734
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 42
      output: 8874
      cache_read: 1770448
      cache_write: 91460
      cost: 1.1835
---
# T-2265 flai test runs svelte-check on the flaiover project when a changed TypeScript, JavaScript, or Svelte file under flaiover selects its new tier

## Work

The cause in I-0206: svelte-check takes no file list, so it runs only in the `--all` `flaiover` tier, and a type error that `flaiover-lint` and `vitest` pass over first shows in the close-out. The fix is a tier that runs svelte-check on the whole project whenever a changed file can change its types. It does not wait for another task.

- Add `scripts/flaiover-check.sh`, POSIX `sh` under `set -eu`, modelled on `scripts/flaiover-lint.sh`. It sources `env.sh`, installs flaiover's dependencies with `flaiover-install.sh --if-needed` in a story worktree, and exits 0 in the main checkout while `flaiover/node_modules` is missing. Then it runs `pnpm check` in `flaiover/` (`svelte-kit sync && svelte-check --tsconfig ./tsconfig.json`). It ignores the file arguments, so that an error a change causes in another file, such as a caller of a changed prop, is caught too. The main branch passes svelte-check, since the close-out's `flaiover` tier gates it, so every error it reports is the story's. Try `--incremental` if repeat runs in one worktree are slow, and keep it only if it makes them faster.
- In `system-flow.yaml` `tests`, add a tier named `svelte-check` after `vitest`, with `dir: flaiover`, `command: ["../scripts/flaiover-check.sh"]`, `format: plain`, and paths `flaiover/src/**/*.ts`, `flaiover/src/**/*.js`, `flaiover/src/**/*.svelte`, `flaiover/*.ts`, `flaiover/*.js`, `flaiover/tsconfig.json`, `flaiover/package.json`, and `scripts/flaiover-check.sh`. Add `svelte-check` to the `flaiover` tier's `covers`, so `--all` and the close-out run svelte-check once.
- In `scripts/flaiover-lint.sh`, replace the comment saying svelte-check stays in `flaiover-test.sh` with one naming the new tier. Add a `scripts/README.md` row for `flaiover-check.sh`.
- In `flai/internal/verify/select_test.go`, add the `svelte-check` line to the `want` of `TestRepositoryManifestLintsChangedFlaioverFiles`. Add a test that this repository's manifest selects `svelte-check` for a changed `flaiover/src/lib/*.ts` file and selects no `svelte-check` tier for a changed `flaiover/README.md`. Name I-0206 in the comment above the new test. The new test fails without the tier, so it reproduces the issue's cause.

## Done when

- `flai test flaiover/src/lib/<some file>.ts` runs `flaiover-lint`, `vitest`, and `svelte-check`, and passes on the branch.
- With a type error put in that file by hand, such as an object literal whose spread overwrites a key (TS2783, the I-0206 instance), `flai test` on the file fails at the `svelte-check` tier. Revert the error afterwards; the wall-clock time of the tier is recorded in the narrative.
- `flai test --all` and the close-out run svelte-check once, through the `flaiover` tier.
- The new test fails without the tier, and the `go-test` and `shellcheck` tiers pass on the changed files.

## Notes
