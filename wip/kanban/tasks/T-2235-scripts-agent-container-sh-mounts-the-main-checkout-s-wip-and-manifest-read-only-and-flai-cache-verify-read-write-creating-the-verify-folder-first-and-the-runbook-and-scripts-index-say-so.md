---
id: T-2235
type: task
nature: remediation
title: scripts/agent-container.sh mounts the main checkout's wip and manifest read-only and .flai-cache/verify read-write, creating the verify folder first, and the runbook and scripts index say so
status: backlog
parent: S-0520
owner: alex
created: 2026-10-28T13:57:52Z
updated: 2026-10-28T13:57:52Z
transitions: []
stream: S-0520
tags: [docker]
touches: [scripts/agent-container.sh, scripts/README.md, docs/operators/runbooks/agent-container.md]
after: [T-2232]
---
# T-2235 scripts/agent-container.sh mounts the main checkout's wip and manifest read-only and .flai-cache/verify read-write, creating the verify folder first, and the runbook and scripts index say so

## Work

Waits for T-2232, whose ADR settles the mounts. In `scripts/agent-container.sh`, beside the two git mounts S-0517's T-2220 adds by default, mount `$CACHE_ROOT/wip` and `$CACHE_ROOT/system-flow.yaml` read-only and `$CACHE_ROOT/.flai-cache/verify` read-write, each at its host path, with `mkdir -p "$CACHE_ROOT/.flai-cache/verify"` before `docker run`, so the folder is the operator's and not root's. Keep the script POSIX `sh` under `set -eu`, every path quoted, and the mount list written out as the existing `-v` lines are, not split from a variable; keep the header's list of what is mounted true. In `docs/operators/runbooks/agent-container.md` § What it mounts, add the three rows with their modes and why each is there (the narrative and the board `flai verify` checks, the manifest by which flai knows the main checkout, and the report the host reads with `flai verify --last`), and in § Limitations replace the close-out's line with what holds after this story: the close-out runs to its tiers, the tiers need the project's layer (S-0510), and the commit step stops inside the container when `--record-issues` recorded something, since commits are the host's. Update the `agent-container.sh` row of `scripts/README.md` to name the mounts.

## Done when

`scripts/agent-container.sh S-nnnn -- flai verify S-nnnn --sync-only` on a started story passes inside `flai-agent-claude-code:local` built from this branch, with `<main>/.flai-cache/verify/S-nnnn.json` on the host owned by the operator; the `shellcheck` and `markdown` tiers pass on the three files; the runbook's mount table lists eight mounts with modes.

## Notes

The `--sync-only` run is the script-level check; the whole close-out is checked once in T-2236. The script is also S-0507's (T-2174), S-0517's (T-2220), and S-0521's (T-2209, done); `after` orders S-0517 and S-0521 before this story, and whichever of S-0507 and this story starts second syncs over the other.
