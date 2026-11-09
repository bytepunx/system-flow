---
id: T-2779
type: task
nature: remediation
title: The root and flaiover package-lock.json files, which nothing reads, are removed, leaving flaiover/pnpm-lock.yaml the one lockfile
status: backlog
parent: S-0630
owner: alex
created: 2026-11-09T21:58:36Z
updated: 2026-11-09T21:58:36Z
transitions: []
stream: S-0630
tags: [flaiover, node]
touches: [flaiover/package-lock.json, package-lock.json]
---
# T-2779 The root and flaiover package-lock.json files, which nothing reads, are removed, leaving flaiover/pnpm-lock.yaml the one lockfile

## Work

- `git rm flaiover/package-lock.json package-lock.json`. The Docker build and CI install from `flaiover/pnpm-lock.yaml` alone, with `--frozen-lockfile` (`flaiover/Dockerfile`, `.github/workflows/flaiover.yml`, `design/tech/node.md` line 20). The root file holds an empty `packages` and there is no root `package.json`.
- First, confirm that no script, workflow, or test names either file: `git grep -n 'package-lock'`, leaving out `design/`, `wip/`, and `docs/`, which the docs task corrects.
- Leave the caret ranges in `flaiover/package.json`: with a frozen pnpm lockfile, the lockfile is the pin.

It waits for no task: no other task of S-0630 changes these files.

## Done when

- `git ls-files | grep package-lock` lists nothing.
- `scripts/flaiover-pnpm.sh install --frozen-lockfile`, or `scripts/flaiover-test.sh` where it runs, passes. Otherwise the task's Notes say why it was not run.

## Notes
