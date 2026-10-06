---
id: TH-0218
title: "S-0274 plan: eight tasks in four layers, forecast raised to 60m"
anchor:
  path: wip/kanban/stories/S-0274-opening-a-story-is-one-call-flai-story-start-moves-it-to-in-progress-opens-the-stream-primes-and-answers-the-first-inbox-together.md
  item: S-0274
status: open
participants: [planner-S-0274]
created: 2026-10-06T22:55:26Z
updated: 2026-10-06T22:55:26Z
---

# TH-0218 S-0274 plan: eight tasks in four layers, forecast raised to 60m

On wip/kanban/stories/S-0274-opening-a-story-is-one-call-flai-story-start-moves-it-to-in-progress-opens-the-stream-primes-and-answers-the-first-inbox-together.md.

## Entries

### 2026-10-06T22:55:26Z planner-S-0274
I drafted eight tasks for S-0274 in four layers. The tasks in a layer share no file and can run together.

- **Layer 1, waits for nothing:**
  - T-1091: opening the branch and worktree moves out of `flai/cmd/branch.go` into `flai/internal/storygit/open.go`, so the MCP server can call it.
  - T-1099: the inbox becomes an exported function in `flai/internal/mcpserver/inbox.go`, over the same per-agent cursor file.
- **Layer 2:**
  - T-1103, after T-1091: a new package `flai/internal/storystart` refuses a story that is not ready or is held. Otherwise it moves the story, opens the stream, and builds the pack.
- **Layer 3, both after T-1099 and T-1103:**
  - T-1109: `flai story start`, as text and `--json`.
  - T-1111: the MCP tool `story_start`. The server's instructions and `wait_for_work` now pull with it.
- **Layer 4:**
  - T-1115, after T-1109: the host channel method `story.start`.
  - T-1117, after T-1109 and T-1111: the start prompt, `session-start.md`, `work-management.md`, `CLAUDE.md`, the template's copies, and the changelog.
  - T-1118, after T-1109 and T-1111: `flai-cli.md`, `flai.md`, and `flai-reference.md`.

The figures:

- **Touches:** widened from 20 to 32 files, keeping all 20 that were declared. The new ones are the files the tasks name, including `CLAUDE.md` and `template/root/CLAUDE.md.tmpl`. There are no folder touches.
- **Forecast:** 60m, against flai's 49m, because the stream-open and inbox code have to be pulled out of `cmd` and the MCP handler first. Delivery is 2026-10-07T10:01Z.
- **Cost of delay:** stays at planner-E-0017's 64 USD a week.

Assumptions (say if any is wrong):

1. `story start` refuses a story that is already in progress, as the criterion says ("not ready"). The prompt for resuming a story in progress keeps `flai stream open`.
2. The CLI's inbox advances the same cursor as the MCP `inbox` under `FLAI_AGENT`, so a change it reports is not reported again.
3. If a step fails after the move, the story stays in progress and the error names the command that finishes the job. The move is not rolled back.
4. `story.start` runs `flai story start --json`, as `item.move` runs `flai move`.
5. No ADR is needed, since this composes existing steps without changing a decision.
6. The guard needs no change: neither the command nor the tool is a read, so sub-agents stay refused both.

I would not split, merge, or drop anything. T-1117 and T-1118 are kept apart so they can run together.
