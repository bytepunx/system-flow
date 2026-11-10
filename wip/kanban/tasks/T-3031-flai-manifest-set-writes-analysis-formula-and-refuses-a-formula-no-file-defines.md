---
id: T-3031
type: task
nature: improvement
title: flai manifest set writes analysis.formula and refuses a formula no file defines or one with an input that has no default
status: backlog
parent: S-0692
owner: alex
created: 2026-11-10T10:33:30Z
updated: 2026-11-10T10:33:30Z
transitions: []
stream: S-0692
tags: [flai, manifest]
touches: [flai/internal/manifest/settings.go, flai/internal/manifest/settings_test.go, flai/internal/formulas/formulas.go, flai/internal/formulas/formulas_test.go, flai/internal/check/check_roles.go, flai/cmd/manifest.go, flai/cmd/manifest_test.go, flai/cmd/serve_actions_test.go]
---
# T-3031 flai manifest set writes analysis.formula and refuses a formula no file defines or one with an input that has no default

## Work

The analyzer ran once in thirty days because this project's `analysis` block sets no `analysis.schedule` and no `analysis.formula` (I-0292). The schedule can be set with `flai manifest set`, but the formula cannot: `design/system/strategic-agents.md` § The analysis schedule says it "is changed by hand", and `CLAUDE.md` says the manifest's keys beyond `name`, `description`, and `dashboard` are flai's. This task makes the formula settable through flai, so that the next task can set both without a hand edit.

- Add `analysis.formula` to the settings catalog in `flai/internal/manifest/settings.go`, after `analysis.schedule`, as `KindText`, with the kebab-case shape check it has already (`manifest_strategic.go`).
- Move the checks `flai check`'s `manifest.analysis` makes on the formula (`analysisFormula` in `flai/internal/check/check_roles.go`: a formula a file under `design/workflows/formulas/` defines, every input with a default) into a pure function in `flai/internal/formulas/formulas.go`, and call it from both the check and `flai/cmd/manifest.go` before `WriteSettings`. The `manifest` package cannot call it: `workitem` imports `manifest` and `formulas` imports `workitem`.
- `flai manifest set analysis.formula=<name>` then writes it, refuses with the check's own message a name no file defines or one whose formula has an input without a default, and writes nothing on a refusal. `--unset analysis.formula` removes it.
- Tests: the catalog and validation tests in `flai/internal/manifest/settings_test.go`; the shared function in `flai/internal/formulas/formulas_test.go`; in `flai/cmd/manifest_test.go`, set writes `analysis-cadence` in a project that defines it and refuses `nope` and a formula with a required input, leaving the manifest as it was. Update the expected analysis keys in `flai/cmd/serve_actions_test.go` (`TestSettingsShowAndSaveTheStrategicAgentsSettings`, which lists `analysis.schedule`).

Waits for no task: it is the code the others build on.

## Done when

- `flai manifest set analysis.formula=analysis-cadence` writes the key in a project that defines the formula, and refuses an undefined formula or one with an input that has no default, writing nothing.
- `flai check` reports `manifest.analysis` as before, through the shared function.
- The tests above pass with `flai test` on the changed paths.

## Notes
