---
id: T-3176
type: task
nature: feature
title: The manifest design, the settings index, the guides, and the template's manifest say how a project names its layer and its arguments, and that flai builds it
status: backlog
parent: S-0745
owner: alex
created: 2026-11-12T16:00:56Z
updated: 2026-11-12T16:07:59Z
transitions: []
stream: S-0745
tags: [cli, template]
touches: [design/system/project-manifest.md, design/system/flai-cli.md, docs/operators/settings.md, docs/users/flai.md, docs/users/flai-reference.md, docs/operators/runbooks/agent-container.md, template/root/system-flow.yaml.tmpl]
after: [T-3175]
---
# T-3176 The manifest design, the settings index, the guides, and the template's manifest say how a project names its layer and its arguments, and that flai builds it

## Work

- `design/system/project-manifest.md`: the `agent.image` key in the YAML example (line 84) and the paragraph on `image.layer` (line 238) gain the arguments, and say that `flai image ensure` builds the layer. Bump `updated`.
- `design/system/flai-cli.md` § The container harness (S-0508), § The project's layer image: flai builds the layer, with the manifest's arguments. Mark it built in the section S-0733 wrote.
- `docs/operators/settings.md`: a row for the new manifest key. `TestSettingsIndex` fails on a manifest key with no row.
- `docs/users/flai.md` § In a container, and `docs/operators/runbooks/agent-container.md` § Build: `flai image ensure` and `flai image build --layer` build the layer. `scripts/agent-image.sh --layer` stays this repository's until S-0760.
- `docs/users/flai-reference.md`: `make flai-reference` for the `--layer` flag, never by hand.
- `template/root/system-flow.yaml.tmpl` (line 48): the commented `agent.image` block shows a layer with its arguments.

It waits for T-3175, whose behaviour it describes.

## Done when

- The seven files describe the key and the build as built.
- `TestSettingsIndex`, `TestReferenceIsCurrent`, `flai check --strict`, the markdown lint, and the template's tests pass.

## Notes
