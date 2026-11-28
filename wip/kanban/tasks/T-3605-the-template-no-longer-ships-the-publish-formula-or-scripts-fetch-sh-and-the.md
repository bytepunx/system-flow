---
id: T-3605
type: task
nature: improvement
title: The template no longer ships the publish formula or scripts/fetch.sh, and the formulas' and scripts' READMEs no longer list them
status: backlog
parent: S-0837
owner: alex
created: 2026-11-17T12:38:06Z
updated: 2026-11-27T23:57:23Z
transitions: []
stream: S-0837
tags: [template, flai]
touches: [template/root/design/workflows/formulas/publish.md, template/root/design/workflows/formulas/README.md, design/workflows/formulas/README.md, scripts/fetch.sh, template/root/scripts/fetch.sh, scripts/README.md, template/root/scripts/README.md, flai/cmd/publish_judgement_test.go, flai/cmd/publish_formula_suites_test.go]
after: [T-3604]
---
# T-3605 The template no longer ships the publish formula or scripts/fetch.sh, and the formulas' and scripts' READMEs no longer list them

## Work

This task waits for T-3604, which teaches `flai upgrade` to remove the retired formula. The template's copy goes once an upgrade over this template still knows its path.

- Delete `template/root/design/workflows/formulas/publish.md`. This repository's copy is gone already: the stopgap commit 4052fc603 removed it.
- Remove the publish formula's row and link from both formulas READMEs. This repository's README links a file that no longer exists.
- `scripts/fetch.sh` (`git fetch --tags`) is named only by the publish formula and the scripts READMEs. The Makefile and flai do not use it. Delete it and its template copy, and remove their rows from both scripts READMEs.
- Delete `flai/cmd/publish_judgement_test.go`. It reads the template's formula and stubs `scripts/fetch.sh`.
- Rewrite any test or fixture that still reads the template's formula, such as a template render test, to expect it absent. Check that `flai check` passes and that `scripts/template-test.sh` renders a project with no publish formula.

## Done when

- Neither the template nor this repository holds `design/workflows/formulas/publish.md` or `scripts/fetch.sh`.
- No README links them.
- `flai test` passes on the changed paths, the template tier included.

## Notes

- Drafted by planner-S-0837.
- `template/CHANGELOG.md` gets its entry in T-3606, which writes the template's whole change in one place.
