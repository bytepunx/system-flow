---
id: S-0759
type: story
nature: feature
title: An operator names a build command of their own that flai runs in place of its image build, or turns the automatic image job off and manages the images outside flai
status: backlog
parent: E-0036
owner: alex
created: 2026-11-12T16:04:53Z
updated: 2026-11-12T16:04:53Z
transitions: []
tags: [cli]
topics: [docker, security]
touches: [flai/internal/config/config.go, flai/internal/config/config_test.go, flai/cmd/serve_harness.go, flai/cmd/serve_settings.go, flai/cmd/serve_settings_test.go, flai/internal/agentimage/ensure.go, flai/internal/agentimage/ensure_test.go, flai/internal/agentimage/command.go, flai/internal/agentimage/command_test.go, flai/internal/serve/imagejob.go, flai/internal/serve/imagejob_test.go, flai/internal/serve/container.go, flai/internal/serve/container_test.go, docs/operators/settings.md, docs/users/flai.md, docs/users/flai-reference.md, docs/operators/runbooks/agent-container.md, design/system/flai-cli.md]
after: [S-0733, S-0751]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
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

## Notes

Proposed slice 5 of E-0036. The flag names are this plan's assumption, and S-0733's ADR may name them otherwise. If the ADR puts the command in a protected path of the project instead, the first criterion changes to read it from there, and `flai/internal/protected` lists that path.
