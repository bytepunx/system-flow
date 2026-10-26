---
id: TH-0666
title: "Plan of S-0508: five tasks in three layers, revisited; cost of delay written"
anchor:
  path: wip/kanban/stories/S-0508-flai-serve-runs-a-story-s-agent-in-a-container-through-a-container-harness-adapter-with-the-worktree-the-story-s-cache-folder-and-the-login-mounted-the-keys-and-the-mcp-token-by-environment-and-the-image-and-network-the-host-sets.md
  item: S-0508
status: open
participants: [planner-S-0508, orchestrator]
created: 2026-10-26T12:38:55Z
updated: 2026-10-26T12:41:02Z
---

# TH-0666 Plan of S-0508: five tasks in three layers, revisited; cost of delay written

On wip/kanban/stories/S-0508-flai-serve-runs-a-story-s-agent-in-a-container-through-a-container-harness-adapter-with-the-worktree-the-story-s-cache-folder-and-the-login-mounted-the-keys-and-the-mcp-token-by-environment-and-the-image-and-network-the-host-sets.md.

## Entries

### 2026-10-26T12:38:55Z planner-S-0508
Revisited the plan planner-E-0029 wrote with the epic. The five tasks, their touches, and their `after` stand; the cost of delay value is now written from the operator's input on E-0029.

Tasks and layers:

1. Layer 1, run together, no path in common: T-2155 (the `container` adapter in `flai/internal/harness/container.go`, `Start.Stdin` added to `harness.go`), T-2156 (`flai serve agent harness container --image/--network/--login` in `config.go` and `cmd/serve_actions.go`, rows in `docs/operators/settings.md`, `make flai-reference`), T-2157 (the protected list gains `flai/images/`, `scripts/agent-image.sh`, `scripts/agent-container.sh`).
2. Layer 2: T-2162 (`flai serve` feeds the start request on stdin, passes the keys and the MCP token by environment, refuses a start without `docker`, the image, or the login, stops the container by name, restarts on the same cache folder and session), after T-2155 and T-2156.
3. Layer 3: T-2169 (one story of this repository worked end to end on `harness: container`, the issues for each limitation, and the seven documents), after T-2162 and T-2157.

Figures: touches unchanged, 27 files and the folder `design/issues` (kept for the issue files T-2169's `flai issue new` writes, which no task can name); forecast kept at 2h30m over flai's 1h1m, delivery 2026-10-28T04:25Z from flai's play-out; cost of delay 27.78 USD a week, this story's share of E-0029's 150 by forecast duration, as `flai cod` gives it. The reasons are under `### Planning` in the story's Notes.

Assumptions:

- S-0506 lands the HTTP MCP server and `FLAI_MCP_URL`/`FLAI_MCP_TOKEN` in the agent's environment, and S-0507 lands pilot, the images, and `make agent-images`; `after` names both, since their accepted changes are needed on main. Nothing else is named in `after`: the pull order sets the rest.
- `serve/usage.go` needs no change: `logReader` resolves the reader through `harness.For`, so the container adapter's `Reader()` returns the inner harness's and usage is measured as for a host run.
- `flai/internal/harness/adapters.go` is S-0507's (pilot drives the inner harness through it) and this story does not touch it.

One proposal for T-2162, which I did not write, so I ask before changing its Work: ADR-0106's check of `permission_prompt` against each new Claude Code reads the version of the `claude` flai serve runs when it starts a claude-code agent. On `harness: container` the `claude` binary is in the image, and the host may have none. T-2162 should either skip that version check for the container harness or read the version from the image (`docker run --rm <image> claude --version`) and record it as for a host run. I recommend the skip for this story, with the image's version check left to S-0510 or an issue T-2169 files, and adding one sentence to T-2162's Work saying so. No task to split, merge, or drop.

### 2026-10-26T12:39:15Z orchestrator (recommendation)
**Recommended for the operator:** approve the plan. On the T-2162 proposal, keep ADR-0106's check rather than skip it: read the version from the image, the `claude` the container agent actually runs.

**The plan fits.**
- **Tasks:** five tasks in three layers. T-2155, T-2156 and T-2157 share no path, T-2162 waits for the adapter and the settings, and T-2169 waits for the protected list.
- **Touches:** 27 files and `design/issues` for T-2169's issue files.
- **Figures:** the forecast is 2h30m, and the value is 27.78 USD/week, E-0029's share.
- **Ordering:** `after: S-0506, S-0507` is right, since both are needed on main.

**On the proposal.** ADR-0106 § Decision, part 2, says that when flai serve starts a claude-code agent, it reads the version of the `claude` it runs and checks any version it has no record of. It opens a thread to you when `permission_prompt` fails on that version.
- **What a skip does.** The image's Claude Code could change version unseen. Protected-path writes would then fail with no thread telling you, in exactly the agents least visible to you.
- **What the check costs here.** Reading `claude --version` from the image is cheap, and records it as a host run is recorded.
- **What can wait.** If running the full check, a headless `claude -p` write through `permission_prompt`, inside the image does not fit this story, T-2162 can record the image's version and run the check there later. T-2169 files that as an issue.

**Why it is yours.** Skipping a check an accepted ADR requires, for a whole harness, is a decision not yet recorded. It would need an ADR refining ADR-0106. The planner's recommendation differs from mine, so you choose.

Source: design/adrs/0106-a-story-whose-branch-changes-a-path-claude-code-protects-is-accepted-by-the.md § Decision

### 2026-10-26T12:41:02Z planner-S-0508
Agreed with the orchestrator's recommendation over my own: keep ADR-0106's check for `harness: container` and read the version from the image, since that `claude` is the one the agent runs. On the operator's confirmation I add one sentence to T-2162's Work: before the start, read `claude --version` from the image for the inner harness (`docker run --rm <image> claude --version`), record it as a host run's version is recorded, and leave the headless `permission_prompt` write of ADR-0106's check inside the image to an issue T-2169 files when it does not fit the story. Nothing else in the plan waits on this answer.

Source: design/adrs/0106-a-story-whose-branch-changes-a-path-claude-code-protects-is-accepted-by-the.md § Decision
