---
id: T-1067
type: task
nature: improvement
title: "The manifest declares a project's test tiers: name, command, the paths that select each, and the format of its output"
status: done
parent: S-0273
owner: alex
created: 2026-10-06T22:51:37Z
updated: 2026-10-07T00:09:43Z
transitions:
  - to: ready
    at: 2026-10-06T23:52:58Z
    by: agent-S-0273
  - to: in-progress
    at: 2026-10-06T23:52:59Z
    by: agent-S-0273
  - to: done
    at: 2026-10-07T00:09:43Z
    by: agent-S-0273
stream: S-0273
tags: [manifest, cli]
touches: [flai/internal/manifest/manifest.go, flai/internal/manifest/manifest_test.go, flai/internal/manifest/settings.go, flai/internal/manifest/settings_test.go, flai/internal/check/check.go, flai/internal/check/tests_test.go, template/root/system-flow.yaml.tmpl, template/CHANGELOG.md, template/template.yaml, design/system/project-manifest.md, docs/operators/settings.md, docs/users/flai-reference.md]
usage:
  source: log
  seconds: 1004
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 136
      output: 60189
      cache_read: 8072223
      cache_write: 194285
      cost: 3.9412
---
# T-1067 The manifest declares a project's test tiers: name, command, the paths that select each, and the format of its output

## Work

Criterion 1 runs "the tiers the project defines". Today they exist only as scripts and `Makefile` targets, which flai cannot select from by path. This task makes them project data. The epic's plan thread TH-0176 assumed this as assumption 1, and alex resolved that thread.

- Add a `tests` key to the manifest beside `checks`: an ordered list, cheapest first. Each tier has:
  - `name`
  - `command`: an argv, with a placeholder for the packages or files selected
  - `paths`: globs that select the tier for a changed or named path
  - `format`: how its output becomes findings, such as `go-test-json`, `vitest-json`, `golangci-json`, `gofmt-list`, or `plain`
  - `all_only`: true for a tier that runs only under `--all`, such as integration and smoke
- Validate the key with the manifest's other checks, so `flai check` names a bad tier by field and reason.
- Add `tests` to what `flai manifest set` writes. Then a project's tiers are set through flai, not by hand, as `CLAUDE.md` says of `system-flow.yaml`.
- With no `tests` key, the default is one `plain` tier running `scripts/test.sh` when it exists. Projects made from the template before this change still get a tier that way.
- `template/root/system-flow.yaml.tmpl` ships the template's `scripts/test.sh`, `integration.sh`, and `smoke.sh` as `plain` tiers, the last two `all_only`. Add a line to `template/CHANGELOG.md`.
- Describe the key in `design/system/project-manifest.md` and in `docs/operators/settings.md` under Project manifest.
- This task waits for nothing. It runs alongside the runner task, which defines its own tier type: the two share no path, and the command task joins them.

## Done when

- [ ] A manifest with a `tests` list loads, round-trips, and validates, and a bad tier is refused with its field and reason (tests in `manifest_test.go`).
- [ ] `flai manifest set` writes and reads back a `tests` list (tests in `settings_test.go`).
- [ ] A manifest without `tests` yields the default `scripts/test.sh` tier.
- [ ] The template manifest, the design, the operator guide, and `template/CHANGELOG.md` describe the key.
- [ ] `go test ./internal/manifest/...`, `flai check --strict`, and the markdown lint pass.

## Notes

Drafted by the planner for S-0273.
