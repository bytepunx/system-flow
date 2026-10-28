---
id: TH-0691
title: "Plan for S-0520: the close-out's flai verify runs inside the container with the main checkout's wip, manifest, and verify folder mounted, between S-0518's flai.sh fallback and S-0510's layer"
anchor:
  path: wip/kanban/stories/S-0520-scripts-close-out-sh-stops-at-its-first-step-inside-the-agent-container-since-scripts-flai-sh-builds-flai-with-go-which-the-image-lacks.md
  item: S-0520
status: open
participants: [planner-S-0520]
created: 2026-10-28T13:58:40Z
updated: 2026-10-28T13:58:40Z
---

# TH-0691 Plan for S-0520: the close-out's flai verify runs inside the container with the main checkout's wip, manifest, and verify folder mounted, between S-0518's flai.sh fallback and S-0510's layer

On wip/kanban/stories/S-0520-scripts-close-out-sh-stops-at-its-first-step-inside-the-agent-container-since-scripts-flai-sh-builds-flai-with-go-which-the-image-lacks.md.

## Entries

### 2026-10-28T13:58:40Z planner-S-0520
The plan, by planner-S-0520 on 2026-10-28. S-0520 stays a draft for you to finalize.

Cause. The close-out's first step is `scripts/flai.sh verify`, which runs `go build` when the tree's `bin/flai` is missing or stale; a fresh worktree has no `bin/` and `flai-agent-claude-code:local` has no `go`, so the step exits 127 before `flai verify` starts.

What other stories already carry. S-0518's T-2224 makes `scripts/flai.sh` run the installed flai without `go`, the cause the issue names. S-0510's layer carries the tiers' tools. S-0521 makes `flai verify` refuse, naming the MCP tools, when the main checkout is absent. Once those land, the close-out inside the container still stops at its first step, now with S-0521's refusal: `flai verify` reads the main checkout's `wip/` (the narrative's two sections, the board, the items, the threads), its `system-flow.yaml` (how `workitem.Open` knows the main checkout is there), its markdownlint configuration, and its `.git` (S-0517's mounts), and writes its report under `<main>/.flai-cache/verify/`, which `flai verify --last` and the orchestrator read on the host. S-0521's thread (TH-0684) assigned "flai verify inside the container with what it needs" to S-0510, but S-0510's plan has no task for it, so S-0520 takes it.

Remediation proposed. The container mounts, beside the worktree and S-0517's two git mounts, the main checkout's `wip/` and `system-flow.yaml` read-only and `.flai-cache/verify/` read-write, at their host paths; the markdown lint configuration falls back to the worktree's when the main checkout's cannot be read; a test with a main checkout stripped to `.git`, `wip/`, and `system-flow.yaml` guards the list; an ADR refining ADR-0203 records it. `scripts/close-out.sh` itself does not change.

Tasks and layers:

- Layer 1: T-2232, the ADR (`design/adrs/README.md`; the folder `design/adrs` kept on the story for the ADR's file).
- Layer 2, after T-2232, no path in common: T-2233 (`check.go`'s lint fallback and the stripped-root test in `flai/cmd/verify_test.go`), T-2234 (the three mounts in `storygit/mounts.go` beside S-0517's two, and the container adapter's `Start` and its argv test), T-2235 (`scripts/agent-container.sh`, the runbook's mount table, `scripts/README.md`).
- Layer 3: T-2236, after all three: `docs/users/flai.md`, `flai-cli.md`, `agent-adapters.md`, the close-out run once by hand inside the image built from the branch, and `flai issue close I-0218`.

Figures: touches 16 files and the folder `design/adrs`; forecast 1h15m over flai's 4m, which was sized with no touches; cost of delay 25.00 USD a week as `flai cod` gives it from the input flai set from the issue. `after`: S-0517 (the mounts go beside its helper and into S-0508's adapter), S-0518 (its fallback is what the by-hand check runs through without `go`), S-0521 (its script change and the refusal this turns into a pass). No `after` on S-0510: its criterion 5, a close-out passing inside the container, needs these mounts rather than the reverse; the pull order settles which lands first.

Assumptions, each a decision you may overturn here:

1. Keep S-0520 rather than drop it. The alternative is to cancel it as covered by S-0518 and S-0510 and let S-0510's agent discover the mounts under its criterion 5. I recommend keeping it: no story names what `flai verify` needs from the main checkout, and S-0510 is sized for the layer alone.
2. Selective mounts, not the main checkout read-only whole. The whole checkout would serve every read with one mount but put `.flai-cache` (the config, the MCP token, the LiteLLM key file, other stories' agent folders) and the operator's untracked files in the container, against ADR-0203; T-2232's ADR names it as the alternative.
3. The lint configuration falls back to the worktree's instead of being mounted. The file's name varies by project (`mdlint.ConfigFiles`), and the worktree's is the branch's tracked copy; on the host nothing changes, since `Root` is `MainRoot` there.
4. The commit step is outside this story. S-0517's plan keeps `.git/objects` read-only and commits the host's, so the close-out inside a container stops at `commit` whenever `--record-issues` recorded something or anything is uncommitted; the agent then commits through `task_done`. This is a limitation S-0517's ADR (T-2218) should weigh: either it keeps the rule and the docs say so (what T-2235 and T-2236 write), or it makes `.git/objects` and the story branch's ref writable, which would let the close-out commit inside the container. I have not changed S-0517.
5. The container's own config (S-0521's `FLAI_CONFIG` under the agent folder) carries no merge-queue host action, so `flai verify` inside it runs the `queue: true` tiers itself: longer, but never fewer checks.
6. Criterion 1 is met when the close-out's first step inside the image built from the branch runs `flai verify` past its own steps to the tiers; a tier stopping for want of a tool is S-0510's, and T-2236 records which, and runs with S-0510's layer too when it is on main by then.

Overlaps: `scripts/agent-container.sh`, the runbook, and `docs/users/flai.md` are claimed by S-0521 in progress, and `docs/users/flai.md` by S-0414 and S-0502; `after` orders S-0517, S-0518, and S-0521 before this story, and S-0507 (T-2174) shares the script, left to the pull order.
