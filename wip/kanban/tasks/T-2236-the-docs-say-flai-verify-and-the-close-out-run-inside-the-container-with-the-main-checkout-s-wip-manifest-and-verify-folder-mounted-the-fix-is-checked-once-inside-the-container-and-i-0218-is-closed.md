---
id: T-2236
type: task
nature: remediation
title: The docs say flai verify and the close-out run inside the container with the main checkout's wip, manifest, and verify folder mounted, the fix is checked once inside the container, and I-0218 is closed
status: backlog
parent: S-0520
owner: alex
created: 2026-10-28T13:58:16Z
updated: 2026-10-28T13:58:16Z
transitions: []
stream: S-0520
tags: [cli, docker, docs]
touches: [docs/users/flai.md, design/system/flai-cli.md, design/system/agent-adapters.md, design/issues/I-0218-scripts-close-out-sh-stops-at-its-first-step-inside-the-agent-container-since-scripts-flai-sh-builds-flai-with-go-which-the-image-lacks.md, design/issues/summary.md]
after: [T-2233, T-2234, T-2235]
---
# T-2236 The docs say flai verify and the close-out run inside the container with the main checkout's wip, manifest, and verify folder mounted, the fix is checked once inside the container, and I-0218 is closed

## Work

Waits for T-2233, T-2234, and T-2235, since it documents what they built and checks it as a whole. In `docs/users/flai.md` § Verify a story before review, after S-0521's paragraph on a worktree whose main checkout is absent, say what the container mounts so that `flai verify` runs inside it (the main checkout's `.git`, `wip/`, and `system-flow.yaml` read-only, its `.flai-cache/verify/` read-write), that the report it writes is the one `flai verify --last` reads on the host, that the markdown lint configuration is the worktree's when the main checkout's is not there, that the `queue: true` tiers run inside the container since its config carries no merge-queue host action, and that the close-out's commit step stops inside the container when `--record-issues` recorded something, the host committing through `task_done`. In `design/system/flai-cli.md`, add the same to `flai verify`'s row, in its words for builders, with the files. In `design/system/agent-adapters.md` § A container harness of flai's own, change the Mounts row from the worktree alone to the list S-0517 and this story give, citing both ADRs, and add one sentence to the paragraph on what the spike settled: the close-out's stop (I-0218) is remediated by S-0518's `scripts/flai.sh` fallback, these mounts, and S-0510's layer. Then check the fix once by hand: build the images from this branch (`make agent-images`), start a story or use this one's worktree, and run `scripts/agent-container.sh S-nnnn -- scripts/close-out.sh S-nnnn` with the worktree clean and committed: the first step runs `flai verify` through `scripts/flai.sh` (the image's flai by S-0518's fallback, since the image has no `go`), passes the rebase, sync, narrative, and check steps, and reaches the tiers, stopping, if it stops, at a tier for want of a tool the layer (S-0510) carries, never at exit 127 and never at S-0521's refusal; record the run's last line in the narrative. Then `flai issue close I-0218 --reason` naming the three parts: S-0518's `scripts/flai.sh` fallback for the `go` build, this story's mounts for the main checkout `flai verify` needs, and S-0510's layer for the tiers' tools. `flai issue close` rewrites `design/issues/summary.md`; commit both.

## Done when

The three documents say what is above and the `markdown` tier passes on them; the narrative records the close-out's last line from inside `flai-agent-claude-code:local` built from this branch, past `flai verify`'s own steps; I-0218 is closed with the reason, `design/issues/summary.md` no longer lists it, and `flai check --strict` scoped to the story is clean.

## Notes

The by-hand check needs S-0518 on main, which `after` on the story guarantees. If the layer of S-0510 is on main by then, run the close-out with it (`flai-agent-system-flow-claude-code:local`) too and record whether every tier passed: that is S-0510's criterion 5 seen from here, and worth a line on its thread.
