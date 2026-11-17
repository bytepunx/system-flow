---
id: T-3574
type: task
nature: improvement
title: "flai manifest set writes and unsets agent.image.layer and agent.image.args.<NAME>, checked by the manifest's image validation"
status: backlog
parent: S-0839
owner: alex
created: 2026-11-17T12:29:41Z
updated: 2026-11-17T12:29:41Z
transitions: []
stream: S-0839
tags: [flai]
touches: [flai/internal/manifest/settings.go, flai/internal/manifest/settings_test.go]
---
# T-3574 flai manifest set writes and unsets agent.image.layer and agent.image.args.<NAME>, checked by the manifest's image validation

## Work

Add two settings to `flai/internal/manifest/settings.go`, resolved by pattern in `SettingFor` as `roleScheduleSetting` resolves `roles.<name>.schedule`:

- `agent.image.layer`, a text setting, the Dockerfile's path.
- `agent.image.args.<NAME>`, a text setting whose `get` reads `m.Agent.Image.Args[NAME]`.

`settingKeys` lists both patterns, so a refusal names them. `--unset` removes either key. Leave the checks of the value to `settingsProblems`, which already runs `m.Agent.Problems()` and so `AgentImage.problems`: it refuses `BASE`, a name that is not an `ARG` name, a value over lines, args with no layer, and a layer that is not a clean path in the repository. Write a value YAML would read as a number, such as `1.20`, quoted, so it reads back as the same string. Leave `flai agent set` and `flai agent clear` as they are: they keep `agent.image`.

In `flai/internal/manifest/settings_test.go`, cover `EditSettings`:

- It writes an arg into a manifest with no `image:` block, and into one with args already, keeping the other args, the order, and the comments.
- It writes the layer.
- It unsets one arg, and unsetting the layer while args remain is refused.
- It refuses `BASE`, a lower-case or digit-first name, and an arg with no layer.
- It writes `1.20` quoted, so the value reads back as the string `1.20`.

This task waits for none: it is the first layer.

## Done when

- `flai manifest set agent.image.args.GO_VERSION=1.26.0` and `agent.image.layer=<path>` write through `EditSettings`, and `--unset` removes them, each refusal naming the field and the validation's reason.
- `flai test flai/internal/manifest/settings.go flai/internal/manifest/settings_test.go` passes.

## Notes
