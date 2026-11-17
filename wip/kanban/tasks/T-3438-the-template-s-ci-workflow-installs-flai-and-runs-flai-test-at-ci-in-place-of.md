---
id: T-3438
type: task
nature: feature
title: The template's CI workflow installs flai and runs flai test --at ci in place of commands of its own
status: backlog
parent: S-0812
owner: alex
created: 2026-11-16T23:47:26Z
updated: 2026-11-16T23:47:26Z
transitions: []
stream: S-0812
tags: [template]
touches: [template/root/.github/workflows/system-flow-check.yml, template/CHANGELOG.md, design/system/template.md, flai/cmd/new_test.go, scripts/template-test.sh]
after: [T-3437]
---
# T-3438 The template's CI workflow installs flai and runs flai test --at ci in place of commands of its own

## Work

ADR-0312 §6, criterion 2. In `template/root/.github/workflows/system-flow-check.yml`, install flai first, then run `flai test --at ci`. Its `markdown` tier then runs the markdown lint, so the separate `DavidAnson/markdownlint-cli2-action` step goes. Keep `flai check --strict`, which is the standard's conformance check and no test tier. Keep the install the file uses now, unless a signed release install through `install.sh` fits the runner better; give the reason in the narrative.

A project made from the template then runs its suites in CI with no command of its own. Its workflow file gives no `verify_at`, so the baseline's `ci: all` runs every tier but the manual ones.

Prove it in `scripts/template-test.sh`, which renders the template and checks the result: run `flai test --at ci` in the rendered project and require it to pass. Extend `flai/cmd/new_test.go` so it asserts that the rendered CI file calls `flai test --at ci`. Add a line under `## Unreleased` in `template/CHANGELOG.md`, and say in `design/system/template.md` what the CI file runs.

It waits for T-3437, because the command it calls must exist for the template test to run it.

## Done when

- The rendered template's `system-flow-check.yml` installs flai, runs `flai test --at ci` and `flai check --strict`, and holds no test or lint command of its own.
- `scripts/template-test.sh` runs `flai test --at ci` in the rendered project and passes.
- `flai/cmd/new_test.go` asserts the call.
- `template/CHANGELOG.md` and `design/system/template.md` say what changed.
- `flai test` on the changed paths passes.

## Notes
