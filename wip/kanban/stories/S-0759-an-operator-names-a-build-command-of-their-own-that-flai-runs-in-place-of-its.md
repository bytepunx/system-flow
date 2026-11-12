---
id: S-0759
type: story
nature: feature
title: An operator names a build command of their own that flai runs in place of its image build, or turns the automatic image job off and manages the images outside flai
status: backlog
parent: E-0036
owner: alex
created: 2026-11-12T16:04:53Z
updated: 2026-11-12T21:10:09Z
transitions: []
tags: [cli]
topics: [docker, security]
touches: [flai/internal/config/config.go, flai/internal/config/config_test.go, flai/cmd/serve_harness.go, flai/cmd/serve_settings.go, flai/cmd/serve_settings_test.go, flai/internal/guard/hostconfig.go, flai/internal/guard/hostconfig_test.go, flai/internal/agentimage/ensure.go, flai/internal/agentimage/ensure_test.go, flai/internal/agentimage/command.go, flai/internal/agentimage/command_test.go, flai/internal/serve/imagejob.go, flai/internal/serve/imagejob_test.go, flai/internal/serve/container.go, flai/internal/serve/container_test.go, docs/operators/settings.md, docs/users/flai.md, docs/users/flai-reference.md, docs/operators/runbooks/agent-container.md, design/system/flai-cli.md]
after: [S-0751]
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
    - kind: orchestrator
      seconds: 12
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 1
          output: 10
          cache_read: 93147
          cache_write: 17092
          cost: 0.0314
draft: true
forecast:
  duration: 42m
  delivery: 2026-11-13T11:04:00Z
  basis: "Its own forecast of 42m; 33rd in the pull order with an in-progress limit of 3, behind S-0659, S-0241, S-0660, S-0661, S-0662, S-0663, S-0730, S-0733, S-0734, S-0735, S-0736, S-0737, S-0738, S-0739, S-0740, S-0741, S-0742, S-0743, S-0744, S-0745, S-0746, S-0747, S-0748, S-0749, S-0750, S-0751, S-0752, S-0753, S-0754, S-0755, S-0756, S-0757 and S-0758."
  by: flai
  at: 2026-11-12T21:10:09Z
---
# S-0759 An operator names a build command of their own that flai runs in place of its image build, or turns the automatic image job off and manages the images outside flai

## Goal

What exists today stays: `agent.image.layer` for a project's tools, and `flai serve agent harness container --image <inner>=<image>` for an image of the operator's own. Two things are missing.

1. An operator whose images come from their own machinery, such as a registry of their own, a CI build, or this repository's `scripts/agent-image.sh`, cannot have flai run it when an image is missing.
2. An operator who manages images outside flai cannot stop `flai serve` building them.

After this story the host's configuration names a build command per project, which `flai image ensure` and the image job run in place of flai's own pull and build. A setting turns the automatic job off, and the start then refuses a missing image as it does today, naming `flai image ensure`. The command is named where S-0733's ADR puts it: by this plan's assumption, the host's configuration. A command there runs with the operator's rights, so no agent may set it.

## Acceptance criteria

- [ ] `flai serve agent harness container --build <command>` names the build command for a project, kept in the host's configuration and shown by `flai serve settings` and `flai serve agent show`. Clearing it brings back flai's own build
- [ ] With a build command named, `flai image ensure` and the image job run it from the project's root with the image's name, its inner harness, and whether it is the layer in named variables. They treat its exit as the outcome, with its output streamed and logged as flai's own build's is, and the step limit applied to it. `flai image status` then reads the images the command made
- [ ] `flai serve agent harness container --no-auto-image` turns the automatic image job off for the project. A start then refuses a missing or stale image, naming `flai image ensure` and the setting, and `flai image ensure` by hand still runs
- [ ] `flai guard` refuses every agent `flai serve` started either setting, as ADR-0220 refuses every command that writes the host's configuration, shown by a test
- [ ] `docs/operators/settings.md` has a row for each setting and each flag, `docs/users/flai.md` and the generated `docs/users/flai-reference.md` describe them, `docs/operators/runbooks/agent-container.md` gains a section on substituting one's own machinery, and `design/system/flai-cli.md` marks the substitution built

## Tasks
- T-3190 The host's configuration keeps, per project, the container harness's build command and whether flai serve brings images back by itself
- T-3191 flai serve agent harness container takes --build and --no-auto-image, flai serve settings shows them, and flai guard refuses them to served agents
- T-3192 flai image ensure runs the operator's build command in place of its own pull and build, with the image named in variables and the step limit applied
- T-3194 The image job runs the operator's build command when the host names one, and a start refuses a missing image when the automatic job is off
- T-3198 The settings index, the guides, the runbook, and flai-cli.md describe the operator's build command and the switch that turns the automatic image job off

## Notes

Proposed slice 5 of E-0036. The flag names are this plan's assumption, and S-0733's ADR may name them otherwise. If the ADR puts the command in a protected path of the project instead, the first criterion changes to read it from there, and `flai/internal/protected` lists that path.

### Planning

- Touches:
  - `flai/internal/config/config.go`, `config_test.go`: layout (`HarnessHost`, lines 444 to 473; `checkHarnesses`, lines 564 to 590).
  - `flai/cmd/serve_harness.go`, `serve_settings.go`, `serve_settings_test.go`: layout. The container's flags (lines 344 to 347) and the settings view live there; `serve_harness.go` has no test file of its own.
  - `flai/internal/guard/hostconfig.go`, `hostconfig_test.go`: layout. ADR-0220's refusal of host-configuration writes, which the fourth criterion tests.
  - `flai/internal/agentimage/command.go`, `command_test.go`, `ensure.go`, `ensure_test.go`: the command's runner, and S-0740's `Ensure` calling it.
  - `flai/internal/serve/imagejob.go`, `imagejob_test.go`, `container.go`, `container_test.go`: S-0751's job and start, for the command and the switch.
  - `docs/operators/settings.md`, `docs/users/flai.md`, `docs/users/flai-reference.md`, `docs/operators/runbooks/agent-container.md`, `design/system/flai-cli.md`: named in the last criterion.
  - Left out from `flai touches suggest`: `docs/operators/index.md` (21%) and the dashboard's documents. The dashboard does not set these settings.
- `after`:
  - S-0733 is not named, though this story reads where the ADR puts the command, and the variables' names: it waits for S-0733 through S-0751, S-0740, and S-0735, which names it.
  - S-0751: the switch and the command change that story's image job and start.
- Tasks, in four layers: T-3190 (the configuration) first; T-3191 (flags, view, guard) and T-3192 (the command's runner) share no file and run together; T-3194 wires the job and the start, after T-3190 and T-3192; T-3198 describes it, after T-3191 and T-3194.
- Forecast 42m stands as `flai forecast` gives it (109 s per unit over 118 done large feature stories, size 23). Its size is mostly documents and tests; the code is two fields, two flags, and one runner. Delivery is as flai gives it.
- Cost of delay: not set yet. E-0036 has no inputs; TH-1094 asks the operator for them.
