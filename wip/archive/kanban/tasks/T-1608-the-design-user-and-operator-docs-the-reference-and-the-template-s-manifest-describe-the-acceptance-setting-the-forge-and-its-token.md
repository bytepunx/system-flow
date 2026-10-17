---
id: T-1608
type: task
nature: feature
title: The design, user, and operator docs, the reference, and the template's manifest describe the acceptance setting, the forge, and its token
status: done
parent: S-0392
owner: alex
created: 2026-10-16T04:33:56Z
updated: 2026-10-17T07:19:36Z
transitions:
  - to: ready
    at: 2026-10-17T07:13:14Z
    by: agent-S-0392
  - to: in-progress
    at: 2026-10-17T07:13:14Z
    by: agent-S-0392
  - to: done
    at: 2026-10-17T07:19:35Z
    by: agent-S-0392
stream: S-0392
tags: [cli, template]
touches: [design/system/project-manifest.md, design/system/flai-cli.md, docs/users/flai.md, docs/users/flai-reference.md, docs/operators/settings.md, docs/operators/index.md, template/root/system-flow.yaml.tmpl, template/CHANGELOG.md, docs/contributors/index.md]
after: [T-1605, T-1606]
usage:
  source: log
  seconds: 381
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 200
      output: 89767
      cache_read: 11059153
      cache_write: 377314
      cost: 6.1631
---
# T-1608 The design, user, and operator docs, the reference, and the template's manifest describe the acceptance setting, the forge, and its token

## Work

- `design/system/project-manifest.md`: the `acceptance` block.
- `design/system/flai-cli.md`: `flai forge check`, `flai serve forge`, and the `forge` package under Internal structure.
- `docs/users/flai.md`: `flai forge check`. Regenerate `docs/users/flai-reference.md` with `make flai-reference`.
- `docs/operators/settings.md`: rows for the three manifest keys, the host's `forges` entry, and the token variable. The settings index test then passes.
- `docs/operators/index.md`: how to make a fine-grained token with only the permissions pull-request acceptance needs (contents, pull requests, checks, read and write), where to put it, and `flai forge check` to try it.
- `template/root/system-flow.yaml.tmpl`: a commented `acceptance` block showing the default. Add a `template/CHANGELOG.md` entry.
- Waits for T-1605 (the ADR the docs link) and T-1606 (the command's final flags, which the reference is generated from).

## Done when

Every document above describes the setting, the forge, and its token, the settings index and reference tests pass, and the markdown lint is clean.

## Notes
