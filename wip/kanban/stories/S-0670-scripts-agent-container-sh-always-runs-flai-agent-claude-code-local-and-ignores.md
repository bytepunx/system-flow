---
id: S-0670
type: story
nature: remediation
title: "scripts/agent-container.sh always runs flai-agent-claude-code:local and ignores the manifest's agent.image.layer, so a close-out run by hand in it has none of the project's test tools"
status: backlog
owner: alex
created: 2026-11-10T08:34:46Z
updated: 2026-11-10T08:49:28Z
transitions: []
tags: [scripts, docker]
topics: [docker, code]
touches: [scripts/agent-image.sh, scripts/agent-container.sh, scripts/agent-container-test.sh, scripts/smoke.sh, scripts/README.md, docs/operators/runbooks/agent-container.md, design/issues/I-0279-scripts-agent-container-sh-always-runs-flai-agent-claude-code-local-and-ignores.md, design/issues/summary.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
cost_of_delay:
  inputs:
    time_lost_per_cycle: 10m
    by: flai
    at: 2026-11-10T08:34:46Z
  value: 25
  by: planner-S-0670
  at: 2026-11-10T08:48:24Z
forecast:
  duration: 40m
  delivery: 2026-11-10T13:14:00Z
  basis: "Its own forecast of 40m; 12th in the pull order with an in-progress limit of 5, behind S-0685, S-0241, S-0659, S-0660, S-0661, S-0662, S-0663, S-0665, S-0666, S-0667, S-0668 and S-0669."
  by: flai
  at: 2026-11-10T08:49:28Z
---
# S-0670 scripts/agent-container.sh always runs flai-agent-claude-code:local and ignores the manifest's agent.image.layer, so a close-out run by hand in it has none of the project's test tools

## Goal

This story remediates [I-0279](../../../design/issues/I-0279-scripts-agent-container-sh-always-runs-flai-agent-claude-code-local-and-ignores.md), "scripts/agent-container.sh always runs flai-agent-claude-code:local and ignores the manifest's agent.image.layer, so a close-out run by hand in it has none of the project's test tools". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0279 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0279 is closed with `flai issue close I-0279 --reason` saying what fixed it

## Tasks
- T-2908 scripts/agent-container.sh runs the manifest's layer image, or FLAI_AGENT_IMAGE, and refuses an image not built
- T-2911 scripts/agent-container-test.sh reproduces I-0279 without Docker, and the smoke tier runs it
- T-2914 The agent-container runbook and scripts/README.md say the script runs the manifest's layer image
- T-2916 I-0279 is closed with flai issue close, naming the fix and its test

## Notes

Cost of delay inputs set by flai from I-0279. time_lost_per_cycle 10m: 10m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-07T06:58:47Z, 3.1 days before this story; under one cycle counts as one).

### Planning

Proposed remedy, from I-0279's one instance (S-0565's T-2452): `scripts/agent-container.sh` runs the image `FLAI_AGENT_IMAGE` names, else the layer's tag `flai-agent-<key>-claude-code:local` when `system-flow.yaml` names `agent.image.layer` (`flai-agent-sf-claude-code:local` here), else `flai-agent-claude-code:local`. That is the tag `scripts/agent-image.sh` builds and `manifest.LayerImage` names, so `flai serve`'s container adapter and the script agree. `scripts/agent-image.sh --tag` prints it, so the tag is read in one place. The script refuses an image not built before `docker run`, naming the build command. A Docker-free test reads the choice through a new `--image` flag.

Touches, each a file, none a folder:

| Touch | Source | Why |
|-------|--------|-----|
| `scripts/agent-container.sh` | design (the issue) | the script that hard-codes the image |
| `scripts/agent-image.sh` | layout | holds `manifest_layer` and the tag's naming, reused through `--tag` |
| `scripts/agent-container-test.sh` | layout | the new test criterion 1 asks for |
| `scripts/smoke.sh` | layout | the tier that runs the script checks, beside `ci-pins.sh` |
| `scripts/README.md` | co-change (4 of 5 commits) | its rows for both scripts name the image |
| `docs/operators/runbooks/agent-container.md` | co-change (3 of 5 commits) | says the script does not run the layer image |
| `design/issues/I-0279-…md` | design (criterion 2) | closed by `flai issue close` |
| `design/issues/summary.md` | design (criterion 2) | regenerated when the issue closes |

`flai touches suggest S-0670 scripts/agent-container.sh` listed the two co-change files; the story declared no touches before.

Figures:

- Forecast 40m, adjusted from flai's 22m. Two reasons. A change under `scripts/` selects the smoke tier at close-out, which runs the installer test and the repository check. Both scripts are on Claude Code's protected list (T-2157), so each edit waits on `permission_prompt` unless auto-approve is on. Delivery moves by the same 18m.
- Cost of delay value 25 USD a week, as `flai cod` works it out from the operator's input, 10m lost per 168h cycle at 150 USD an hour. It stands: one instance in three days gives no reason to change it.

No `after`: the story reads nothing another open story writes. S-0682 changes `scripts/agent-container.sh` too, and S-0672 changes the layer's Dockerfile. A trial merge judges those overlaps, so neither is named.
