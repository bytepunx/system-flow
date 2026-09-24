---
id: T-0402
type: task
nature: feature
title: Document the config file and link the reference
status: in-progress
parent: S-0017
owner: alex
created: 2026-09-24T08:12:56Z
updated: 2026-09-24T08:17:42Z
transitions:
  - to: ready
    at: 2026-09-24T08:13:04Z
    by: system-flow
  - to: in-progress
    at: 2026-09-24T08:17:42Z
    by: system-flow
stream: S-0017
tags: [cli]
touches: [docs/users, design/system]
---
# T-0402 Document the config file and link the reference

## Work
Rewrite the Configuration section of `docs/users/flai.md` so it documents every key in `~/.flai/config.json`, including the keys `flai serve` manages (`host_actions`, `agent`, `checks`, `import_roots`), `worktrees`, and the retired push key fields, plus `FLAI_CONFIG` and `FLAI_CACHE_DIR`. Replace the stale "Implemented so far" line with a link to the generated reference, link it from `docs/users/index.md`, and describe the generator in `design/system/flai-cli.md`.

## Done when
Every field of `config.Config` is documented in `docs/users/flai.md`, the reference is linked from both user pages, the design says how the reference is produced, and `make lint-md` and `flai check --strict` pass.

## Notes
