---
id: S-0839
type: story
nature: improvement
title: No flai command writes agent.image.args, so a project's layer build arguments can only be written by hand
status: backlog
owner: alex
created: 2026-11-17T08:32:50Z
updated: 2026-11-17T12:35:58Z
transitions: []
tags: [flai, template]
touches: [flai/internal/manifest/settings.go, flai/internal/manifest/settings_test.go, flai/cmd/manifest.go, flai/cmd/manifest_test.go, docs/users/flai-reference.md, docs/operators/settings.md, docs/users/flai.md, design/system/project-manifest.md, design/system/flai-cli.md, template/root/system-flow.yaml.tmpl, template/CHANGELOG.md, design/issues/I-0364-no-flai-command-writes-agent-image-args-so-a-project-s-layer-build-arguments.md, design/issues/summary.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: sum
  seconds: 0
  models: []
  strategic:
    - kind: planner
      seconds: 4
      estimated: true
      models:
        - model: claude-haiku-5-5
          input: 2
          output: 7
          cache_read: 184573
          cache_write: 2317
          cost: 0.0128
        - model: claude-opus-5-5
          input: 6
          output: 1727
          cache_read: 502256
          cache_write: 14423
          cost: 0.2521
    - kind: orchestrator
      seconds: 12
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 1
          output: 18
          cache_read: 124457
          cache_write: 4104
          cost: 0.0402
cost_of_delay:
  inputs:
    time_lost_per_cycle: 10m
    by: flai
    at: 2026-11-17T08:32:50Z
  value: 25
  by: planner-S-0839
  at: 2026-11-17T12:29:23Z
forecast:
  duration: 1h
  delivery: 2026-11-17T21:48:00Z
  basis: "Its own forecast of 1h; 19th in the pull order with an in-progress limit of 3, behind S-0806, S-0834, S-0811, S-0838, S-0829, S-0840, S-0809, S-0842, S-0812, S-0808, S-0807, S-0810, S-0241, S-0813, S-0814, S-0815, S-0822, S-0824, S-0835 and S-0837."
  by: flai
  at: 2026-11-17T12:35:58Z
finalized:
  by: orchestrator
  at: 2026-11-17T12:35:34Z
---
# S-0839 No flai command writes agent.image.args, so a project's layer build arguments can only be written by hand

## Goal

This story remediates [I-0364](../../../design/issues/I-0364-no-flai-command-writes-agent-image-args-so-a-project-s-layer-build-arguments.md), "No flai command writes agent.image.args, so a project's layer build arguments can only be written by hand". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0364 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0364 is closed with `flai issue close I-0364 --reason` saying what fixed it

## Tasks
- T-3574 flai manifest set writes and unsets agent.image.layer and agent.image.args.<NAME>, checked by the manifest's image validation
- T-3577 flai manifest set's help names the image keys, and a command test reproduces I-0364
- T-3581 The settings index, the manifest design, the CLI design, the user guide, and the template's manifest say flai manifest set writes the layer and its args
- T-3582 Close I-0364 with the reason that flai manifest set now writes agent.image.layer and agent.image.args

## Notes

Cost of delay inputs set by flai from I-0364. time_lost_per_cycle 10m: 10m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-13T20:00:47Z, 3.5 days before this story; under one cycle counts as one).

### Planning

Proposed remedy, from I-0364's one instance (S-0760, TH-1172): `flai manifest set` writes `agent.image.args.<NAME>=<value>` and `agent.image.layer=<path>`, and `--unset` removes either. They are two settings resolved by pattern in `SettingFor`, as `roles.<name>.schedule` is. They are checked by the validation the manifest already has (`AgentImage.problems` in `flai/internal/manifest/agent.go`, which `settingsProblems` runs through `m.Agent.Problems()`): no `BASE`, an `ARG` name, one line, and no args without a layer. `agent.image.layer` is in scope because the instance names it as the same exception. Without it, a project made from the template, which names no layer, still has to hand-write a line before any arg can be set. Out of scope: the dashboard's settings page (`flai/internal/hostapi/settings.go`) and `projects[].install` (D-0593). The `tests` tiers that the instance cites through D-0449 are written by `flai manifest set` already (`testsLines`).

Touches, and where each came from:

| Path | Source | Why |
|------|--------|-----|
| `flai/internal/manifest/settings.go` | design, co-change | `SettingFor`, `settingKeys`, the two settings |
| `flai/internal/manifest/settings_test.go` | co-change (73%) | write, unset, and refusals by `EditSettings` |
| `flai/cmd/manifest.go` | layout | the command's help names the new keys |
| `flai/cmd/manifest_test.go` | layout | the test that reproduces I-0364 through the command |
| `docs/users/flai-reference.md` | co-change (24%) | generated from the help by `make flai-reference` |
| `docs/operators/settings.md` | co-change (79%) | the two rows say "Changed by hand" today |
| `docs/users/flai.md` | co-change, documentation.md | `flai manifest set`'s user documentation |
| `design/system/project-manifest.md` | co-change (24%) | `image.args` (I-0364 named there) and the key list of `flai manifest set` |
| `design/system/flai-cli.md` | design | the `flai manifest set` row |
| `template/root/system-flow.yaml.tmpl` | co-change (12%) | the commented `layer` and `args` lines name the command |
| `template/CHANGELOG.md` | layout | the template change's entry |
| `design/issues/I-0364-…md`, `design/issues/summary.md` | criterion 2 | `flai issue close` writes them |

`flai/internal/manifest/manifest.go` and `manifest_test.go` co-change at 82–85%. They are left out: the image's struct and validation are in `agent.go`, which needs no change. No folder touch.

Tags `flai` and `template`: the story delivers the command to flai and a changed manifest comment to the template.

No `after`: the story reads nothing that another open story writes.

Forecast: 1h, delivery 2026-11-17T21:10:00Z. Before the touches were written, flai gave 2m from size 2. With the 13 touches it gives 28m: 108 s per unit of size over 186 large improvement stories, times size 15, 17th in the pull order. That figure is raised to 1h. Each of the four tasks closes with a commit and a sync, one regenerates the reference, and the close-out alone took 9m43s on S-0828 today. The delivery is flai's 20:37 plus the 32m added.

Cost of delay: 25 USD a week, as `flai cod` worked it out from the inputs flai set from I-0364 (10m lost per 168h cycle at the project's 150 USD an hour). It stands: the issue's one instance cost what its input says, and nothing else depends on the command.
