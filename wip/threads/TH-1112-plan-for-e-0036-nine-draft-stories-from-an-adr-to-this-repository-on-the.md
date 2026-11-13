---
id: TH-1112
title: "Plan for E-0036: nine draft stories from an ADR to this repository on the general path"
anchor:
  path: wip/kanban/epics/E-0036-any-project-using-flai-gets-its-agent-images-from-flai-which-supplies-rebuilds.md
  item: E-0036
status: resolved
participants: [planner-E-0036, orchestrator, alex]
created: 2026-11-12T16:12:05Z
updated: 2026-11-13T20:22:48Z
---

# TH-1112 Plan for E-0036: nine draft stories from an ADR to this repository on the general path

On wip/kanban/epics/E-0036-any-project-using-flai-gets-its-agent-images-from-flai-which-supplies-rebuilds.md.

## Entries

### 2026-11-12T16:12:05Z planner-E-0036
E-0036 now has nine draft stories in the backlog, each with tasks, file-level touches, a forecast, and a `### Planning` note. None has a cost of delay yet: TH-1094 asks you for the epic's inputs.

## Stories, in order

| # | Story | Forecast | `after`, and what it reads there |
|---|---|---|---|
| 1 | S-0733 ADR: where the images come from, the label, the job, the substitution | 25m | none |
| 2 | S-0734 Each flai release publishes the signed bedrock and claude-code images, and `pilot` for Linux | 45m | S-0733: registry, names, label, list format |
| 3 | S-0735 `flai image status`, and one package that decides every agent image | 33m | S-0733: label, staleness rule, command name |
| 4 | S-0740 `flai image ensure` and `build`: signed pull, host build as fallback | 1h | S-0734: the list and the Linux `pilot`. S-0735: the status and `flai/cmd/image.go` |
| 5 | S-0745 flai builds the project's layer, with build arguments from the manifest | 40m | S-0740: its runner, `ensure`, and `build` |
| 6 | S-0751 `flai serve` brings an image back as a host job that holds the start | 50m | S-0740: `Ensure`. S-0745: the layer's build |
| 7 | S-0757 Dashboard: Ops › Agent images, the live job, a failure's output, Retry | 50m | S-0751: the job, its record in `agents.json`, its log |
| 8 | S-0759 The operator's own build command, and a switch for the automatic job | 42m | S-0751: the job and the start it changes |
| 9 | S-0760 This repository on the general path, and the template's documents | 32m | S-0745: the manifest's argument shape. S-0759: the build command's setting and variables |

S-0740, S-0745, S-0751, and S-0759 read S-0733's ADR too. They wait for it through S-0735, so they do not name it, which keeps S-0733 off `wip.after-hub`. S-0734 and S-0735 can run side by side, and so can S-0757 and S-0759.

## Tasks and layers

- S-0733: T-3149 → T-3151 → T-3152 (the ADR, the design, then your confirmation and acceptance).
- S-0734: T-3159, T-3160, and T-3161 together → T-3162 (the release job) → T-3163 (the documents).
- S-0735: T-3164 → T-3165 and T-3166 together → T-3167.
- S-0740: T-3168 and T-3169 together → T-3170 (after T-3168) → T-3171 → T-3172.
- S-0745: T-3173 → T-3174 → T-3175 → T-3176.
- S-0751: T-3177 → T-3178 → T-3179 → T-3180.
- S-0757: T-3181 → T-3182, T-3183, and T-3184 together → T-3185 → T-3186.
- S-0759: T-3190 → T-3191 and T-3192 together → T-3194 → T-3198.
- S-0760: T-3207, T-3209, and T-3212 together → T-3214.

## Assumptions, each for S-0733's ADR to confirm or change

- **Source of the images:** pull by a signed digest first, build on the host as the fallback. T-3149 recommends this, and S-0734 holds only if publication is chosen.
- **Command name:** `flai image status`, `ensure`, and `build`. `flai agent` already sets the project's default agent.
- **Layer arguments:** a map under `agent.image` in the manifest, with the Dockerfile's `ARG` defaults for the rest.
- **Your own build command:** kept in the host's configuration (`--build`), with `--no-auto-image` to turn the job off. Agents cannot set either (ADR-0220).
- **Consent:** none asked for the first automatic build. It is reported as it runs, and the switch turns it off.
- **Testing:** a story's container has no Docker socket. Every build and pull is tested over a fake docker, and a real one is yours to run on the host at review.
- **Strategic agents:** the epic says they "never get the project's layer", but they do today. `startOf` and `containerImageHost` apply the layer to every role. S-0733's ADR settles whether they should, and S-0751 follows it.

## Proposals

- **Fold S-0734 into S-0740** if the ADR chooses a host build alone. Only its `pilot` criterion would stay.
- **Merge nothing else.** S-0735's refactor of the image decision is worth landing alone, because it touches five start paths.
- **Add or drop nothing.** S-0732 (the MCP servers listening on the engine's gateway) came from the same change of engine, but it is not image work, so I left it outside E-0036.

## What you need to do

1. Answer TH-1094 with the cost of delay inputs. I then write each story's value.
2. Finalize the drafts you accept.

### 2026-11-12T16:14:34Z orchestrator (recommendation)
**Recommended:** take the plan as drafted, with no story added, merged, or dropped now.

**Why it holds:**

- The nine stories run in one line from S-0733's ADR to S-0760, and each `after` names what the story reads there. S-0734 and S-0735 can run side by side, and so can S-0757 and S-0759.
- Each has tasks in layers, file-level touches, and a forecast. The only gap `flai promote --drafts` reports is the cost of delay value.
- **Fold S-0734 into S-0740:** decide that when S-0733's ADR is accepted, not now. If the ADR chooses a host build alone, S-0734 keeps only its `pilot` criterion, and the planner can merge it then.
- **Strategic agents and the layer:** leave it to S-0733's ADR, as the planner says. S-0751 follows it.

**Your two steps:**

1. **Cost of delay:** TH-1094 carries my recommendation, `revenue_per_week: 300` and `time_lost_per_cycle: 2h`. That figure is yours alone. Once you answer or confirm it, I set it under `set_cost_of_delay`, and the planner writes each story's value.
2. **Finalizing:** nothing for you to do. Under `finalize_drafts` I finalize each draft once `flai promote --drafts` lists it complete. TH-1102 to TH-1111 track each one.

**Know before they run:** S-0740 and S-0760 write paths the container protects (`flai/images`, `scripts/agent-image.sh`). Run their agents on the host, or make those changes yourself (TH-1106, TH-1111).

Source: design/conventions/strategic-agents.md § As the orchestrator

### 2026-11-13T20:22:48Z alex
Resolved: S-0760 was accepted
