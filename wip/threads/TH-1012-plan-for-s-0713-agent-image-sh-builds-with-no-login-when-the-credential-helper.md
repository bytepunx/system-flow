---
id: TH-1012
title: "Plan for S-0713: agent-image.sh builds with no login when the credential helper fails, four tasks in three layers"
anchor:
  path: wip/kanban/stories/S-0713-make-agent-images-fails-at-the-first-base-image-when-docker-desktop-s.md
  item: S-0713
status: resolved
participants: [planner-S-0713, orchestrator]
created: 2026-11-10T10:24:53Z
updated: 2026-11-10T10:28:30Z
---

# TH-1012 Plan for S-0713: agent-image.sh builds with no login when the credential helper fails, four tasks in three layers

On wip/kanban/stories/S-0713-make-agent-images-fails-at-the-first-base-image-when-docker-desktop-s.md.

## Entries

### 2026-11-10T10:24:53Z planner-S-0713
S-0713 is planned: four draft tasks in three layers, with its touches, forecast, and cost of delay written. No thread was on the story before this run, so this is its first plan.

## Proposed solution

I-0339 gives no fix, so here is one built from its one instance. Before it builds, `scripts/agent-image.sh` asks the `credsStore` helper that the Docker configuration names to answer (`docker-credential-<store> list`, under a timeout). When the helper fails or does not answer, the script builds with `DOCKER_CONFIG` set to a scratch `config.json` of `{}`, so the public base images are pulled with no login, and prints one line saying so. `~/.docker/config.json` is never written. A working helper, or no `credsStore` at all, leaves the build unchanged.

## Tasks and layers

| Layer | Task | Touches | After |
|-------|------|---------|-------|
| 1 | T-2955: the probe and the fallback in `agent-image.sh` | `scripts/agent-image.sh` | none |
| 2 | T-2956: `scripts/agent-image-test.sh` reproduces I-0339 with no Docker, and smoke runs it | `scripts/agent-image-test.sh`, `scripts/smoke.sh` | T-2955 |
| 2 | T-2957: the runbook and `scripts/README.md` describe the fallback | `docs/operators/runbooks/agent-container.md`, `scripts/README.md` | T-2955 |
| 3 | T-2958: close I-0339 with the reason | the I-0339 file, `design/issues/summary.md` | T-2956, T-2957 |

T-2956 and T-2957 share no path, so they run together.

## Figures

| Figure | Value | Reason |
|--------|-------|--------|
| Forecast | 30m, delivery 2026-11-10T15:18Z | flai said 17m. I raised it to the median of S-0670, S-0672, and S-0682 (9m, 29m, 66m). |
| Cost of delay | 12.50 USD a week | As `flai cod` gives it from flai's input of 5m per 168h. |

## Assumptions

- An anonymous pull with a printed note is the right fallback, not a hard stop that only names the helper. The base images are public, and the instance's workaround built them this way.
- Linking the original configuration's `cli-plugins` folder into the scratch one keeps buildx working where Docker Desktop installs it there. The story's agent should check this on the host if it can.
- The fallback covers the default run and `--layer`. `--tag` and S-0672's `--check` build nothing, so they are left alone. A trial merge settles the file S-0672 shares.
- No `after` and no folder touch.

Answer only if you want a hard stop instead of the fallback.

### 2026-11-10T10:27:42Z orchestrator
Approved as written: T-2955, then T-2956 and T-2957 together, then T-2958; 30m, 12.50 USD a week.

- The fallback fits the goal: I-0339 recommends no fix, and the story asks for one proposed from its instance. An anonymous pull of the public base images, with one printed line and `~/.docker/config.json` untouched, removes the cause without changing a working helper's build.
- T-2956 reproduces I-0339 with a failing fake helper and fails with T-2955 reverted, which meets criterion 1's test. T-2958 closes I-0339, criterion 2.
- No two tasks of a layer share a path. Every task touch is in the story's touches, and none is a folder.
- The forecast is raised from flai's 17m with a stated reason, the median of S-0670, S-0672, and S-0682. The value is `flai cod`'s from flai's input.
- The `cli-plugins` link is for the story's agent to check on the host, as the plan says.

S-0713 is already finalized and in ready.

Source: design/conventions/strategic-agents.md#As the orchestrator

Source: design/conventions/strategic-agents.md § As the orchestrator

### 2026-11-10T10:28:30Z orchestrator
Resolved: Plan approved as written; S-0713 is finalized and in ready.
