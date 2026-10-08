---
id: T-1374
type: task
nature: feature
title: The user and design documents describe the provider, and the template's manifest shows provider and providers
status: backlog
parent: S-0349
owner: alex
created: 2026-10-08T08:43:22Z
updated: 2026-10-08T08:43:22Z
transitions: []
stream: S-0349
tags: [cli, template]
touches: [docs/users/flai.md, design/system/project-manifest.md, design/system/work-hierarchy.md, design/system/flai-cli.md, template/root/system-flow.yaml.tmpl, template/CHANGELOG.md]
after: [T-1373]
---
# T-1374 The user and design documents describe the provider, and the template's manifest shows provider and providers

## Work

- `docs/users/flai.md`, "Starting an agent": the provider, where its entry lives, that `key_env` names a variable and never holds a key, and the commands of T-1372 and T-1373.
- `design/system/project-manifest.md`: `agent.provider` and `providers.<name>` with the keys a project may override and the one it may not.
- `design/system/work-hierarchy.md`: `agent.provider` in a story's front matter.
- `design/system/flai-cli.md`: `flai serve agent provider`, the `--provider` flags, and the resolver's refusals.
- `template/root/system-flow.yaml.tmpl`: `provider:` under the commented agent, and a commented `providers:` example with `api`, `base_url`, and `models`; an entry in `template/CHANGELOG.md`.

It waits for T-1373, so that it describes the commands as built.

## Done when

- The five documents and the template say what the code does, with no example that holds a key.
- `flai test docs/users/flai.md design/system/ template/` passes, the template test among it.

## Notes

Layer 4 of S-0349.
