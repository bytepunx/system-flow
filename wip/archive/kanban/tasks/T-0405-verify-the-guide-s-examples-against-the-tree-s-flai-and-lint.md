---
id: T-0405
type: task
nature: feature
title: Verify the guide's examples against the tree's flai and lint
status: done
parent: S-0018
owner: alex
created: 2026-09-24T08:12:56Z
updated: 2026-09-24T08:19:31Z
transitions:
  - to: ready
    at: 2026-09-24T08:13:18Z
    by: agent-S-0018
  - to: in-progress
    at: 2026-09-24T08:15:47Z
    by: agent-S-0018
  - to: done
    at: 2026-09-24T08:19:31Z
    by: agent-S-0018
stream: S-0018
tags: []
---

# T-0405 Verify the guide's examples against the tree's flai and lint

## Work
Run the guide's examples with `scripts/flai.sh` in a scratch directory: fork the template, add a variable and a `.tmpl` file, render with `--var` and `--layout`, check the result. Run `make smoke`, `make lint-md` in the worktree and the main checkout, and `flai check --strict`.

## Done when
Every example ran and did what the guide says; smoke, markdown lint, and `flai check --strict` pass.

## Notes
