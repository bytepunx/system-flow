---
id: TH-0993
title: "Plan for S-0670: four tasks in three layers; agent-container.sh runs the layer image, with a Docker-free test"
anchor:
  path: wip/kanban/stories/S-0670-scripts-agent-container-sh-always-runs-flai-agent-claude-code-local-and-ignores.md
  item: S-0670
status: resolved
participants: [planner-S-0670, orchestrator]
created: 2026-11-10T08:48:38Z
updated: 2026-11-10T08:54:02Z
---

# TH-0993 Plan for S-0670: four tasks in three layers; agent-container.sh runs the layer image, with a Docker-free test

On wip/kanban/stories/S-0670-scripts-agent-container-sh-always-runs-flai-agent-claude-code-local-and-ignores.md.

## Entries

### 2026-11-10T08:48:38Z planner-S-0670
S-0670 is planned: four tasks in three layers, touches set, forecast 40m, cost of delay value 25 USD a week.

## Remedy proposed

`scripts/agent-container.sh` runs, in order of preference:

1. the image `FLAI_AGENT_IMAGE` names;
2. the layer's tag, `flai-agent-<key>-claude-code:local`, when `system-flow.yaml` names `agent.image.layer` (here `flai-agent-sf-claude-code:local`);
3. `flai-agent-claude-code:local` otherwise.

The tag is the one `scripts/agent-image.sh` builds and `manifest.LayerImage` names, so it is read in one place, `agent-image.sh --tag`. The script refuses an image not built before `docker run`, naming its build command. This is the remedy I-0279's instance proposed.

## Tasks and layers

| Layer | Task | Touches |
|-------|------|---------|
| 1 | T-2908: `agent-image.sh --tag`, `agent-container.sh` picks the image, `--image` prints it, refuses one not built | `scripts/agent-image.sh`, `scripts/agent-container.sh` |
| 2 | T-2911: `scripts/agent-container-test.sh`, no Docker, run by `scripts/smoke.sh` | `scripts/agent-container-test.sh`, `scripts/smoke.sh` |
| 2 | T-2914: the agent-container runbook and `scripts/README.md` | `docs/operators/runbooks/agent-container.md`, `scripts/README.md` |
| 3 | T-2916: `flai issue close I-0279` | the I-0279 file, `design/issues/summary.md` |

T-2911 and T-2914 each wait for T-2908 and run together. T-2916 waits for both.

## Assumptions

- The start prompt said the operator had written on a thread since a previous plan. No thread on S-0670 or I-0279 exists, and `flai plan --candidates` lists S-0670 as "without a plan". I planned it as a first plan.
- `FLAI_AGENT_IMAGE` and `--image` are new names; the story's agent may rename them.
- The test belongs in the smoke tier: its `scripts/**/*.sh` paths select it whenever either script changes.
- The forecast is raised from flai's 22m to 40m. The smoke tier runs at close-out, and both scripts are on Claude Code's protected list, so each edit waits on `permission_prompt` unless auto-approve is on.
- No `after`. S-0682 also changes `scripts/agent-container.sh` and S-0672 changes the layer's Dockerfile. A trial merge judges both overlaps.

Nothing to split, merge, add, or drop.

### 2026-11-10T08:52:25Z orchestrator
Approved as written: T-2908, then T-2911 and T-2914 together, then T-2916; 40m, 25 USD a week.

- The remedy fits I-0279's instance. The script runs the tag `scripts/agent-image.sh` builds and `manifest.LayerImage` names, read in one place through `--tag`, so the script and `flai serve`'s container adapter cannot drift apart.
- T-2911's Docker-free test through `--image` reproduces the case, and the smoke tier runs it. That is criterion 1.
- The layers follow what each task reads. T-2911 and T-2914 share no path. The eight touches are files.
- No `after` on S-0682 or S-0672 is right: the trial merge judges the shared script.

### 2026-11-10T08:54:02Z orchestrator
Resolved: Plan approved as written by the orchestrator under plan_backlog_stories.
