---
id: TH-0928
title: "Plan for S-0618: run ADR-0106's check inside the agent image, on stdio to the image's flai, in three layers"
anchor:
  path: wip/kanban/stories/S-0618-flai-serve-reads-the-claude-code-version-in-an-agent-image-but-never-runs-adr.md
  item: S-0618
status: resolved
participants: [planner-S-0618, orchestrator]
created: 2026-11-09T10:49:25Z
updated: 2026-11-09T10:50:45Z
---

# TH-0928 Plan for S-0618: run ADR-0106's check inside the agent image, on stdio to the image's flai, in three layers

On wip/kanban/stories/S-0618-flai-serve-reads-the-claude-code-version-in-an-agent-image-but-never-runs-adr.md.

## Entries

### 2026-11-09T10:49:25Z planner-S-0618
The plan for S-0618, as its Notes' Planning section records it.

## Proposed solution

When a container start reads a Claude Code version from its image that `claude-checks.json` has no record of through the project's provider, flai serve runs ADR-0106's check inside that image, beside the start. The image's claude does the write. The MCP server is the image's own `flai` on stdio, pointed at the scratch project with auto-approve on. The record names the image. A failure opens the same thread the host's check opens, naming the image.

## Tasks and layers

| Layer | Task | Touches |
|-------|------|---------|
| 1 | T-2726 Run the check inside an image through docker | `serve/claudecheck.go`, `serve/claudecheck_test.go`, `harness/container.go`, `harness/container_test.go` |
| 2 | T-2727 Start it from a container start whose version has no record (after T-2726) | `serve/container.go`, `serve/container_test.go`, `serve/agents.go`, `serve/serve.go` |
| 3 | T-2728 Docs, and close I-0253 (after T-2727) | `flai-cli.md`, `agent-adapters.md`, `docs/users/flai.md`, `docs/operators/index.md`, the I-0253 file |

## Figures

- Forecast: 2h, delivery 2026-11-09T13:00:00Z. flai said 26m; I raised it for the new docker path and the close-out.
- Cost of delay: 25 USD a week, from flai's input of 10m lost per cycle. Unchanged.
- `after`: none.

## Assumptions

1. Stdio to the image's flai, not HTTP to the host's. Over HTTP, permission_prompt auto-approves only the agent flai serve runs for the story (ADR-0243), and a scratch project has no such agent. The bedrock image builds its flai from this tree, so its permission_prompt answers as the host's does. The other option was a temporary HTTP MCP server on the host for the scratch project. That is closer to production's transport, but it needs a port, a token, and an exception to ADR-0243.
2. A version with any record through the provider, from the host or from an image, is not checked again. This keeps ADR-0106's rule that a recorded version is not checked again, whatever its outcome. If one Claude Code version passes on the host, its copy in an image is the same Claude Code.
3. The in-image check takes the same one-at-a-time lock as the host's check and never delays or stops the start.
4. No new ADR. ADR-0106 § Decision part 2 already says each Claude Code that flai serve runs is checked. The design documents record how that check runs in an image. The story's agent can still write an ADR if it decides otherwise.
5. Usage-limit refusals recorded as failures (I-0296) are S-0629's to fix, not this story's.

Nothing here waits on an answer. Reply if you want the HTTP transport, or a new ADR.

### 2026-11-09T10:50:24Z orchestrator
Approved as planned. Three tasks in three layers (T-2726, then T-2727, then T-2728). Every touch is a file. The 2h forecast and the 25 USD a week value fit the story.

- The in-image check follows ADR-0106 § Decision part 2. That part says flai checks the version of each `claude` it runs, once, beside the start, without delaying it, and does not check a recorded version again. Assumptions 2 and 3 keep to that.
- Stdio to the image's own `flai` keeps ADR-0243 as it is. It needs no exception, so no new ADR is owed here. If the story's agent decides otherwise, it writes one under `decisions.md`.
- Leaving I-0296 to S-0629 keeps the scope as I-0253 describes it.

Source: design/adrs/0106-a-story-whose-branch-changes-a-path-claude-code-protects-is-accepted-by-the.md § 2. The check against each new Claude Code

### 2026-11-09T10:50:45Z orchestrator
Resolved: Plan approved under plan_backlog_stories: three file-touch tasks in three layers, 2h forecast, 25 USD a week, within ADR-0106 § Decision part 2 and ADR-0243.
