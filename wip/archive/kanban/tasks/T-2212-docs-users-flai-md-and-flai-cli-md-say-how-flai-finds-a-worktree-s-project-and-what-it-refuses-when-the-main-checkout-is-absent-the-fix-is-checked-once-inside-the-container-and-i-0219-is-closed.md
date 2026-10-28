---
id: T-2212
type: task
nature: remediation
title: docs/users/flai.md and flai-cli.md say how flai finds a worktree's project and what it refuses when the main checkout is absent, the fix is checked once inside the container, and I-0219 is closed
status: done
parent: S-0521
owner: alex
created: 2026-10-28T13:03:41Z
updated: 2026-10-28T13:42:18Z
transitions:
  - to: ready
    at: 2026-10-28T13:38:57Z
    by: agent-S-0521
  - to: in-progress
    at: 2026-10-28T13:38:58Z
    by: agent-S-0521
  - to: done
    at: 2026-10-28T13:42:18Z
    by: agent-S-0521
stream: S-0521
tags: [cli, docker, docs]
touches: [docs/users/flai.md, design/system/flai-cli.md, design/issues/I-0219-flai-verify-inside-the-agent-container-takes-the-story-s-worktree-for-the-project-root-since-the-main-checkout-is-not-mounted-and-writes-a-config-under-the-container-s-home.md, design/issues/summary.md]
after: [T-2209, T-2211]
usage:
  source: log
  seconds: 200
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 44
      output: 11208
      cache_read: 2054304
      cache_write: 118018
      cost: 1.3629
---
# T-2212 docs/users/flai.md and flai-cli.md say how flai finds a worktree's project and what it refuses when the main checkout is absent, the fix is checked once inside the container, and I-0219 is closed

## Work

This task waits for T-2209 and T-2211, since it documents and checks what they built. `docs/users/flai.md` § Story branches gains the rule in two or three sentences: flai run in a story's worktree finds the project through the worktree's `.git` file, keeps `wip/` and `.flai-cache` in the main checkout, and, when the main checkout is not there, as in an agent container that mounts the worktree alone, refuses `verify`, `stream sync`, `task done`, and `check`, naming the main checkout and the MCP tools that run them on the host; § Configuration says `FLAI_CONFIG` and `FLAI_CACHE_DIR` are how a container's flai is kept out of the home. `design/system/flai-cli.md` records the same for builders where it describes the worktree and the configuration, with `updated` bumped. Then check the fix where I-0219 saw it, with the images built from this tree (`make agent-images`): `FLAI_AGENT_GIT_MOUNT=1 scripts/agent-container.sh S-nnnn -- flai verify S-nnnn` in this story's own worktree prints the refusal naming the main checkout and the MCP tool, and `scripts/agent-container.sh S-nnnn -- ls -A /home/agent` shows no `.flai`; record both outputs in the narrative. Last, close the issue with `flai issue close I-0219 --reason "<what fixed it>"` from the story's worktree, which writes the issue file and `design/issues/summary.md`, and tick the story's criteria you verified with `flai criteria tick`.

## Done when

- `docs/users/flai.md` and `design/system/flai-cli.md` describe the root rule for a worktree, the refusal when the main checkout is absent, and the two variables, and the markdown lint passes
- Inside `flai-agent-claude-code:local` built from this tree, `flai verify S-nnnn` in the mounted worktree refuses naming the main checkout and the MCP tool, and nothing is written under `/home/agent`, both recorded in the narrative
- I-0219 is closed with `flai issue close I-0219 --reason` naming T-2207, T-2211, and T-2209's change, and `design/issues/summary.md` no longer lists it

## Notes

If the image cannot be built on the host at the time, say so in the narrative and leave the by-hand check unticked rather than claim it; the Go tests of T-2211 reproduce the instance on their own.
