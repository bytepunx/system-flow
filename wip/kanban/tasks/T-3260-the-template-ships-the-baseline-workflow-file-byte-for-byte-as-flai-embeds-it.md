---
id: T-3260
type: task
nature: feature
title: The template ships the baseline workflow file byte for byte as flai embeds it, and flai upgrade keeps a project's changes to it
status: backlog
parent: S-0755
owner: alex
created: 2026-11-12T16:10:15Z
updated: 2026-11-12T16:10:15Z
transitions: []
stream: S-0755
tags: [cli, template]
touches: [template/root/design/workflows/board.md, flai/internal/workflow/baseline_test.go, flai/internal/upgrade/upgrade.go, flai/internal/upgrade/upgrade_test.go, template/CHANGELOG.md, template/template.yaml]
after: [T-3253]
---
# T-3260 The template ships the baseline workflow file byte for byte as flai embeds it, and flai upgrade keeps a project's changes to it

## Work

Waits for T-3253: a project made from the template must have its file read by the load that task adds. Add the baseline file to the template, byte for byte the embedded copy, and extend `baseline_test.go` to hold the two equal, as `roles/baseline_test.go` holds the role files. Make `flai upgrade` keep a project's own changes to the file, as the ADR says, and add it to a project that lacks it only where the ADR says so. Record the change in `template/CHANGELOG.md` and bump `template/template.yaml`.

## Done when

`baseline_test.go` and `upgrade_test.go` pass, and `scripts/template-test.sh` renders a project whose check is clean.

## Notes
