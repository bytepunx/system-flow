---
id: T-3581
type: task
nature: improvement
title: The settings index, the manifest design, the CLI design, the user guide, and the template's manifest say flai manifest set writes the layer and its args
status: backlog
parent: S-0839
owner: alex
created: 2026-11-17T12:30:01Z
updated: 2026-11-17T12:30:01Z
transitions: []
stream: S-0839
tags: [flai, template, docs]
touches: [docs/operators/settings.md, design/system/project-manifest.md, design/system/flai-cli.md, docs/users/flai.md, template/root/system-flow.yaml.tmpl, template/CHANGELOG.md]
after: [T-3574]
---
# T-3581 The settings index, the manifest design, the CLI design, the user guide, and the template's manifest say flai manifest set writes the layer and its args

## Work

Say what T-3574 built, in each document that says these keys are written by hand:

- `docs/operators/settings.md`: the rows `agent.image.layer` and `agent.image.args.<name>` say "Changed with `flai manifest set agent.image.layer=<path>`" and "`agent.image.args.<name>=<value>`", and "removed with `--unset`", in place of "Changed by hand".
- `design/system/project-manifest.md`:
  - § `image.args` names the command.
  - The paragraph on this repository drops "written by hand since no flai command writes `args` (I-0364)".
  - The key list under `flai manifest set` ("Not among them") adds `agent.image.layer` and `agent.image.args.<NAME>` to the keys outside the strategic blocks that it writes.
- `design/system/flai-cli.md`: the `flai manifest set` row names the two keys.
- `docs/users/flai.md`: where it documents `flai manifest set` and the agent image's layer, give the command with an example.
- `template/root/system-flow.yaml.tmpl`: the commented `layer:` and `args:` lines name `flai manifest set agent.image.layer=<path>` and `agent.image.args.<NAME>=<value>`.
- `template/CHANGELOG.md`: one line under `## Unreleased` for the comment.

This task waits for T-3574: the documents describe the keys, values, and refusals that T-3574 settles. It shares no path with T-3577, so the two run together.

## Done when

- No document under `docs/` or `design/system/` says `agent.image.layer` or `agent.image.args` is changed by hand only, or that no flai command writes them.
- `flai test docs/operators/settings.md design/system/project-manifest.md design/system/flai-cli.md docs/users/flai.md template/root/system-flow.yaml.tmpl template/CHANGELOG.md` passes, the settings index and the markdown lint among its tiers.

## Notes
