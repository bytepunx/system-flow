---
id: S-0805
type: story
nature: feature
title: A tier in the manifest belongs to a suite, flai test --suite runs one, and the template's tiers are labelled static, behavior, integration, and smoke
status: in-progress
parent: E-0034
owner: alex
created: 2026-11-16T23:41:56Z
updated: 2026-11-17T02:48:57Z
transitions:
  - to: ready
    at: 2026-11-16T23:51:58Z
    by: orchestrator
  - to: in-progress
    at: 2026-11-16T23:52:21Z
    by: agent-S-0805
blocked:
  - from: 2026-11-17T00:54:07Z
    until: 2026-11-17T01:18:29Z
    reason: awaiting the operator on TH-1258
  - from: 2026-11-17T01:18:50Z
    until: 2026-11-17T01:18:50Z
    reason: awaiting the operator on TH-1258
  - from: 2026-11-17T01:19:16Z
    until: 2026-11-17T02:47:35Z
    reason: awaiting the operator on TH-1258
  - from: 2026-11-17T02:48:57Z
    reason: awaiting the operator on TH-1258
tags: [flai, template]
topics: [verify, workflow]
touches: [design/adrs/0312-a-tier-belongs-to-a-suite-and-the-workflow-file-names-the-suites-a-transition.md, design/adrs/README.md, flai/internal/manifest/manifest_tests.go, flai/internal/manifest/manifest_test.go, flai/internal/manifest/settings.go, flai/internal/manifest/settings_test.go, flai/internal/check/tests_test.go, flai/internal/verify/verify.go, flai/internal/verify/manifest.go, flai/internal/verify/manifest_test.go, flai/internal/verify/select.go, flai/internal/verify/select_test.go, flai/cmd/test.go, flai/cmd/test_test.go, flai/internal/mcpserver/test.go, flai/internal/mcpserver/test_test.go, template/root/system-flow.yaml.tmpl, template/CHANGELOG.md, system-flow.yaml, design/system/project-manifest.md, design/system/flai-cli.md, docs/users/flai.md, docs/users/flai-reference.md, docs/operators/settings.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 3815
  priced_by: harness
  turns:
    - day: 2026-11-16
      work: 24
    - day: 2026-11-17
      ceremony: 3
      hand_edits: 1
      work: 44
  models:
    - model: claude-opus-5-5
      input: 366
      output: 80061
      cache_read: 24712600
      cache_write: 1178507
      cost: 14.482
  strategic:
    - kind: planner
      seconds: 49
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 81
          output: 4316
          cache_read: 8648876
          cache_write: 128431
          cost: 2.7132
    - kind: orchestrator
      seconds: 28
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 8
          output: 57
          cache_read: 413428
          cache_write: 47500
          cost: 0.1411
cost_of_delay:
  value: 275.4
  by: planner-S-0805
  at: 2026-11-16T23:49:25Z
forecast:
  duration: 2h30m
  delivery: 2026-11-17T04:28:00Z
  basis: "Two comparable tier stories took 1h45m (S-0547) and 2h38m (S-0548); this one adds a manifest key, a check, a flag, an MCP argument, two manifests, and four documents, so 2h30m, delivered after the work flai's forecast puts ahead of it in the pull order, which ends about 01:58Z."
  by: planner-S-0805
  at: 2026-11-16T23:49:25Z
finalized:
  by: orchestrator
  at: 2026-11-16T23:51:53Z
---
# S-0805 A tier in the manifest belongs to a suite, flai test --suite runs one, and the template's tiers are labelled static, behavior, integration, and smoke

## Goal

[ADR-0312](../../../design/adrs/0312-a-tier-belongs-to-a-suite-and-the-workflow-file-names-the-suites-a-transition.md) §1. A tier in the manifest's `tests:` may give `suite: <name>`; a tier that gives none is the suite of its own name. This is the label every later story of the ADR reads. Accept the ADR, or record what the operator changes in it, as this story's first task.

## Acceptance criteria

- [ ] ADR-0312 is accepted, or superseded by the ADR this story writes from the operator's changes.
- [x] A tier's `suite` is read from the manifest, kebab-case, defaulting to the tier's name, and `flai check` refuses one that is not kebab-case.
- [x] `flai test --suite <name>` runs the tiers of a suite that the paths select, and with `--all` every tier of it but a `manual` one.
- [x] The template's manifest labels its tiers `static`, `behavior`, `integration`, and `smoke`, and this repository's manifest labels its own without changing which tiers any command runs.
- [x] The manifest design, the reference, and the user guide describe `suite`.

## Tasks
- T-3444 ADR-0312 is accepted, or superseded by an ADR recording the operator's changes
- T-3446 A tier's suite is read from the manifest, defaults to its name, and flai check refuses one that is not kebab-case
- T-3447 flai test --suite and the MCP tool test run the tiers of one suite that the paths select, or under --all every one but a manual tier
- T-3448 The template's manifest labels its tiers static, behavior, integration, and smoke, and this repository's manifest labels its own, with no change to what runs
- T-3449 The manifest design, the CLI design, and the user guide describe a tier's suite and flai test --suite

## Notes

### Planning

Touches. The story declared none, and `flai touches suggest` needs a starting path, so every touch is predicted, all files, no folder:

| Touch | From |
|-------|------|
| `design/adrs/0312-…-transition.md` | design: the Goal links it, and criterion 1 accepts it |
| `design/adrs/README.md` | layout: `flai adr accept` updates the ADR's index row |
| `flai/internal/manifest/manifest_tests.go`, `manifest_test.go` | layout: `TestTier` and `TestErrors`, where `flai check`'s `manifest.tests` comes from |
| `flai/internal/manifest/settings.go`, `settings_test.go` | layout: `flai manifest set tests=` parses strictly and renders each tier's keys |
| `flai/internal/check/tests_test.go` | layout: the check's test of tier errors |
| `flai/internal/verify/verify.go`, `manifest.go`, `manifest_test.go`, `select.go`, `select_test.go` | layout: `verify.Tier`, `FromManifest`, and the selection `flai test` runs; co-change with the tier stories S-0547 and S-0548 |
| `flai/cmd/test.go`, `test_test.go` | layout: the `flai test` command and its flag |
| `flai/internal/mcpserver/test.go`, `test_test.go` | layout: the MCP tool `test`, which mirrors the command's arguments |
| `template/root/system-flow.yaml.tmpl`, `template/CHANGELOG.md` | design: criterion 4; the template's changes go in its changelog |
| `system-flow.yaml` | design: criterion 4 |
| `design/system/project-manifest.md`, `docs/users/flai.md`, `docs/users/flai-reference.md` | design: criterion 5 |
| `design/system/flai-cli.md` | layout: its Commands table describes `flai test`'s flags |
| `docs/operators/settings.md` | layout: `TestSettingsIndexHasEveryKey` needs a `tests[].suite` row, and `make flai-reference` writes its flag index |

After. None. No open story writes something this one reads. S-0806 to S-0815 read this story's `suite` label; the planner of those stories names S-0805 in their `after`.

Forecast. `flai forecast` gave 10m with no touches, then 54m with 24 (size 29 at 110 s a unit). Adjusted to 2h30m: the two tier stories nearest this one took 1h45m (S-0547) and 2h38m (S-0548), and this one also waits on the operator for the ADR. The delivery, 2026-11-17T04:28Z, is flai's estimate of when the work ahead of it in the pull order ends, about 01:58Z, plus the 2h30m.

Cost of delay. `flai cod` gives 275.40 USD a week. That is S-0805's share of E-0034's 1,500 USD a week, the operator's input of 10h lost a cycle, split by forecast: 2h30m of the 13h37m forecast over the epic's 20 open stories that have no inputs of their own. It stands as given.
