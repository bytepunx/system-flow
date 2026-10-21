---
id: ADR-0185
title: "flai accept runs a sub-project's install command in the main checkout when the accepted change touched the paths the manifest's projects[].install names, and a failed install warns without failing the acceptance"
status: accepted
date: 2026-10-21
supersedes: []
superseded_by: []
---

# ADR-0185 flai accept runs a sub-project's install command in the main checkout when the accepted change touched the paths the manifest's projects[].install names, and a failed install warns without failing the acceptance

## Context

A story worktree has no `flaiover/node_modules` of its own. Since S-0482 its integration tier installs one when it is missing or older than the lock file (I-0178). The main checkout's install is never brought up to date. When an accepted story adds a dependency to `flaiover/package.json` and `flaiover/pnpm-lock.yaml`, the main checkout keeps the old `node_modules`. Every run there that needs the new package then fails: an unscoped `flai test`, the components host action's dev server ([ADR-0163](0163-flai-runs-a-sub-project-s-own-vite-dev-server-on-127-0-0-1-behind-the.md)), and, before S-0482, every story's close-out, which borrowed that install (I-0179). Someone had to notice and install by hand.

flai knows which paths an acceptance changed, since it tells the open stories about them (S-0132). It does not know how a sub-project installs its dependencies. That differs per kind and per project: `pnpm install`, `npm ci`, `go mod download`, or a script that uses the project's own store.

## Decision

A sub-project may declare an install command in the manifest. `flai accept` runs it in the main checkout after the acceptance commit when the accepted change touched one of the paths that declaration names.

- `projects[].install` has `command`, the argument list run with no shell in the sub-project's folder, and `paths`, globs relative to that folder written as `components` globs are. Both are required when `install` is set. `flai check` reports a bad one as `manifest.install`.
- After the commit step, `flai accept` compares the paths the acceptance brought into the main branch with each sub-project's `install.paths`. The comparison uses the paths it already reads for the overlap notices. It runs each matching command once, in sub-project order, in the main checkout. A resumed acceptance that only commits uses the story's committed paths, as its overlap notices do.
- A command that fails, or cannot start, is logged at warn with its exit and its last lines of output. The text output says the install failed and names the command to run by hand. `--json` lists each run under `installs` with its sub-project, command, and outcome. The acceptance still succeeds: the merge and the commit are behind it.
- `--dry-run` runs nothing. The command's output goes to standard error, so standard output stays flai's.
- This repository sets flaiover's `install` to `../scripts/flaiover-install.sh --if-needed` on `package.json` and `pnpm-lock.yaml`. That script uses the pnpm store `scripts/env.sh` sets (I-0177).

## Consequences

- An acceptance that changes flaiover's lock file leaves the main checkout's `node_modules` current, so nothing that runs there afterwards misses the new package.
- The merge queue's landing runs through `flai accept --batch` and installs the same way.
- An acceptance that changes a lock file takes as long as the install. An acceptance that changes none runs nothing.
- An install that fails leaves the main checkout as it was before this decision, and says so, rather than failing an acceptance that has already merged.
- A project made from the template sets `install` on its own sub-projects. Without it, nothing changes.

## Alternatives considered

- Installing in the test that needs the dependencies: it fixes one test, not the dashboard's dev server or a run in the main checkout. S-0482 already covers the worktree's side.
- An install command inferred from the sub-project's kind: the right command depends on the project's package manager and store, and this repository's needs its own script.
- Refusing the acceptance when the install fails: the merge is already committed, so a refusal would only hide a finished acceptance behind an error.
- A `post_accept` hook of any command: wider than the cause, and it would run on every acceptance rather than only on a changed lock file.
