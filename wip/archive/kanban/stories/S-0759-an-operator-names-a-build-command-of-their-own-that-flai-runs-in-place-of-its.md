---
id: S-0759
type: story
nature: feature
title: An operator names a build command of their own that flai runs in place of its image build, or turns the automatic image job off and manages the images outside flai
status: done
parent: E-0036
owner: alex
created: 2026-11-12T16:04:53Z
updated: 2026-11-13T19:45:12Z
transitions:
  - to: ready
    at: 2026-11-13T07:31:00Z
    by: alex
  - to: in-progress
    at: 2026-11-13T16:13:20Z
    by: agent-S-0759
  - to: review
    at: 2026-11-13T19:44:51Z
    by: agent-S-0759
  - to: done
    at: 2026-11-13T19:45:12Z
    by: alex
tags: [cli]
topics: [docker, security]
touches: [flai/internal/config/config.go, flai/internal/config/config_test.go, flai/cmd/serve_harness.go, flai/cmd/serve_settings.go, flai/cmd/serve_settings_test.go, flai/internal/guard/hostconfig.go, flai/internal/guard/hostconfig_test.go, flai/internal/agentimage/ensure.go, flai/internal/agentimage/ensure_test.go, flai/internal/agentimage/command.go, flai/internal/agentimage/command_test.go, flai/internal/serve/imagejob.go, flai/internal/serve/imagejob_test.go, flai/internal/serve/container.go, flai/internal/serve/container_test.go, docs/operators/settings.md, docs/users/flai.md, docs/users/flai-reference.md, docs/operators/runbooks/agent-container.md, design/system/flai-cli.md, flai/cmd/serve_agent.go, flai/cmd/serve_agent_image.go, flai/cmd/serve_agent_image_test.go, flai/cmd/serve_view.go, flai/internal/guard/guard_tables.go, flai/cmd/image.go, flai/cmd/image_test.go, flai/internal/serve/agents.go, flai/internal/serve/checks.go, flai/internal/check/image.go, flai/internal/check/image_test.go]
after: [S-0751]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 6077
  priced_by: harness
  turns:
    - day: 2026-11-13
      test_runs: 1
      work: 71
  models:
    - model: claude-opus-5-5
      input: 458
      output: 174954
      cache_read: 34719970
      cache_write: 1255186
      cost: 17.8108
  strategic:
    - kind: orchestrator
      seconds: 30
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 10
          output: 71
          cache_read: 342819
          cache_write: 46942
          cost: 0.1129
forecast:
  duration: 42m
  delivery: 2026-11-13T17:01:00Z
  basis: "Its own forecast of 42m; 2nd in the pull order with an in-progress limit of 3, behind S-0760."
  by: flai
  at: 2026-11-13T16:13:06Z
finalized:
  by: alex
  at: 2026-11-13T07:28:25Z
---
# S-0759 An operator names a build command of their own that flai runs in place of its image build, or turns the automatic image job off and manages the images outside flai

## Goal

What exists today stays: `agent.image.layer` for a project's tools, and `flai serve agent harness container --image <inner>=<image>` for an image of the operator's own. Two things are missing.

1. An operator whose images come from their own machinery, such as a registry of their own, a CI build, or this repository's `scripts/agent-image.sh`, cannot have flai run it when an image is missing.
2. An operator who manages images outside flai cannot stop `flai serve` building them.

After this story the host's configuration names a build command per project, which `flai image ensure` and the image job run in place of flai's own pull and build. A setting turns the automatic job off, and the start then refuses a missing image as it does today, naming `flai image ensure`. The command is named where S-0733's ADR, [ADR-0295](../../../design/adrs/0295-an-installed-flai-pulls-its-agent-images-by-the-digests-a-signed-list-on-its.md) point 7, puts it: the host's configuration. A command there runs with the operator's rights, so no agent may set it.

## Acceptance criteria

- [x] `flai serve agent image command [--project <root>] -- <program> [args...]` names the build command for a project, kept in the host's configuration as `agent.projects.<root>.image.command` and shown by `flai serve settings` and `flai serve agent show`. Clearing it brings back flai's own build
- [x] With a build command named, `flai image ensure` and the image job run it from the project's root with the `FLAI_IMAGE_*` variables ADR-0295 names: the version, the inner harness, the tags and recipes the images must have, and which are missing or stale. They treat its exit as the outcome, with its output streamed and logged as flai's own build's is, and the step limit applied to it. `flai image status` then reads the images the command made, and a missing image or other labels after it fail the run
- [x] `flai serve agent image auto off [--project <root>]` turns the automatic image job off for the project. A start then refuses a missing or stale image, naming `flai image ensure` and the setting, and `flai image ensure` by hand still runs
- [x] `flai guard` refuses every agent `flai serve` started either setting, as ADR-0220 refuses every command that writes the host's configuration, shown by a test
- [x] `docs/operators/settings.md` has a row for each setting and each flag, `docs/users/flai.md` and the generated `docs/users/flai-reference.md` describe them, `docs/operators/runbooks/agent-container.md` gains a section on substituting one's own machinery, and `design/system/flai-cli.md` marks the substitution built

## Tasks
- T-3190 The host's configuration keeps, per project, the container harness's build command and whether flai serve brings images back by itself
- T-3191 flai serve agent image command and auto set them, flai serve settings and agent show show them, and flai guard refuses them to served agents
- T-3192 flai image ensure runs the operator's build command in place of its own pull and build, with the image named in variables and the step limit applied
- T-3194 The image job runs the operator's build command when the host names one, and a start refuses a missing image when the automatic job is off
- T-3198 The settings index, the guides, the runbook, and flai-cli.md describe the operator's build command and the switch that turns the automatic image job off

## Notes

Proposed slice 5 of E-0036. The flag names were this plan's assumption, and S-0733's ADR could name them otherwise. It did: ADR-0295 point 7 names `flai serve agent image command [--project <root>] -- <program> <args>…`, stored as `agent.projects.<root>.image.command`, `flai serve agent image auto off [--project <root>]`, and the `FLAI_IMAGE_*` variables. On starting the story, its agent rewrote criteria 1 to 3 to the ADR's names (2026-11-13). The goal and what each criterion asks are unchanged. Planned as `flai serve agent harness container --build` and `--no-auto-image`.

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
