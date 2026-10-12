---
id: T-1457
type: task
nature: improvement
title: The manifest's markdown tier fixes with lint-md.sh --fix, and the docs and an ADR say what fix_command does
status: done
parent: S-0367
owner: alex
created: 2026-10-12T16:52:55Z
updated: 2026-10-12T17:07:15Z
transitions:
  - to: ready
    at: 2026-10-12T16:57:58Z
    by: agent-S-0367
  - to: in-progress
    at: 2026-10-12T16:57:59Z
    by: agent-S-0367
  - to: done
    at: 2026-10-12T17:07:15Z
    by: agent-S-0367
stream: S-0367
tags: []
touches: [system-flow.yaml, design/system/project-manifest.md, design/system/flai-cli.md, docs/users/flai.md, docs/operators/settings.md, template/root/system-flow.yaml.tmpl, template/CHANGELOG.md, design/adrs]
after: [T-1455]
usage:
  source: log
  seconds: 556
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 56
      output: 21746
      cache_read: 3436601
      cache_write: 95116
      cost: 1.7128
---
# T-1457 The manifest's markdown tier fixes with lint-md.sh --fix, and the docs and an ADR say what fix_command does

## Work

- `system-flow.yaml`'s `markdown` tier gets `fix_command: [scripts/lint-md.sh, --fix, "{files}"]`.
- `design/system/project-manifest.md`, `design/system/flai-cli.md` (`flai stream sync`, `flai task done`, `flai accept`), `docs/users/flai.md`, and `docs/operators/settings.md` (`tests[].fix_command` row) say what it does and when it runs.
- An ADR records the decision: a tier's `fix_command` runs before `flai task done` commits, after each clean sync's rebase, and before `flai accept` commits in the main checkout, and never stops any of them.
- The template's manifest comment names `fix_command`; `template/CHANGELOG.md` notes it.
- Waits for T-1455, whose key the settings test checks for a row.

## Done when

- [ ] The settings test finds a row for `tests[].fix_command`, and the markdown lint passes on the docs changed.

## Notes
