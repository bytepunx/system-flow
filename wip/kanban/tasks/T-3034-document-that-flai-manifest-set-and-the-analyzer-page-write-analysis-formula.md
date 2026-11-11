---
id: T-3034
type: task
nature: improvement
title: Document that flai manifest set and the Analyzer page write analysis.formula, and regenerate the command reference
status: done
parent: S-0692
owner: alex
created: 2026-11-10T10:33:45Z
updated: 2026-11-11T02:17:03Z
transitions:
  - to: ready
    at: 2026-11-11T02:09:44Z
    by: system-flow
  - to: in-progress
    at: 2026-11-11T02:09:45Z
    by: system-flow
  - to: done
    at: 2026-11-11T02:17:03Z
    by: agent-S-0692
stream: S-0692
tags: [docs, manifest]
touches: [docs/users/flai-reference.md, docs/users/flai.md, docs/operators/settings.md, docs/operators/index.md, design/system/strategic-agents.md, design/system/project-manifest.md, design/system/flai-cli.md, design/system/flaiover-dashboard.md, template/CHANGELOG.md]
after: [T-3031]
usage:
  source: log
  seconds: 438
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 47
      output: 9743
      cache_read: 5457200
      cache_write: 276605
      cost: 2.9103
---
# T-3034 Document that flai manifest set and the Analyzer page write analysis.formula, and regenerate the command reference

## Work

Waits for T-3031: the command reference is generated from `flai manifest set`'s help, which lists the catalog's keys, so it can only be regenerated once the key is in the catalog.

- Regenerate `docs/users/flai-reference.md` with `make flai-reference` (`scripts/flai-reference.sh`); `TestReferenceIsCurrent` fails while it is stale.
- Replace the sentences that say `analysis.formula` is changed by hand, or that `flai manifest set` and the Analyzer page do not write it, with how it is set now: `design/system/strategic-agents.md` § The analysis schedule, `docs/operators/settings.md` (the `analysis.formula` row), `docs/operators/index.md`, and `docs/users/flai.md` (the by-hand paragraph, and a row in the `manifest set` key table).
- Add the key to the settings table of `design/system/project-manifest.md`, to the key list of the `flai manifest set` row in `design/system/flai-cli.md`, and to the Analyzer page's settings panel in `design/system/flaiover-dashboard.md`, which lists `analysis.agent` and `analysis.schedule`. `StrategicSettings.svelte` filters the catalog by prefix, so the row should appear with no flaiover change: check it on the page or in its test, and if it does not, say so on the plan thread rather than widen this task.
- Add an entry to `template/CHANGELOG.md` under `## Unreleased`: `- S-0692 ... (patch, I-0292).`

## Done when

- No document under `design/` or `docs/` says `analysis.formula` is changed by hand or not written by `flai manifest set`.
- `docs/users/flai-reference.md` is current, and `flai test` passes on the changed paths, the `settings-index` tier among them.

## Notes
