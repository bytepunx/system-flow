---
id: T-1065
type: task
nature: remediation
title: flai new and flai import default to the template's newest version tag and record it as the project's template.ref
status: backlog
parent: S-0301
owner: alex
created: 2026-10-06T22:50:35Z
updated: 2026-10-06T22:54:07Z
transitions: []
stream: S-0301
tags: [cli, template]
touches: [flai/cmd/new.go, flai/cmd/new_test.go, flai/cmd/import.go, flai/cmd/import_test.go]
after: [T-1060, T-1062]
---
# T-1065 flai new and flai import default to the template's newest version tag and record it as the project's template.ref

## Work

Criterion 2: a new project, or one brought under flai with `flai import`, starts at the newest template version, not at whatever the cache's clone of `main` holds.

- In `resolveTemplate` in `flai/cmd/new.go`, used by `flai new`, `flai import` and `flai template show`, a git template with no `--ref` resolves to T-1060's `Latest` tag when the configured ref follows releases (`FollowsReleases`: empty, the default branch such as `main`, or a version tag). Another ref the operator set with `flai template use --ref`, or passed as `--ref`, is used as given, with `1.0.60` matched to `v1.0.60`. A repository with no version tags falls back to its default branch, fetched fresh.
- `flai new` and `flai import` write the tag they rendered to `template.ref` in `system-flow.yaml` and to the lock, so the project's first upgrade starts from a known version.
- The config's default `template.ref` stays `main`: every config written so far holds it, and as the default branch it follows releases. Changing the default would leave every existing install on the branch (agent-S-0301's review). `resolveTemplate` is shared with `flai upgrade`, which passes its ref explicitly and does its own resolution in T-1064.
- This task waits for T-1060, whose `Latest` it calls, and for T-1062, whose ADR it implements. It shares no path with T-1064 and runs alongside it. `flai upgrade` passes its ref explicitly and does not depend on this default.

## Done when

- [ ] A test in `flai/cmd/new_test.go` uses a template repository with tags v1.0.18 and v1.0.60. `flai new` with no `--ref` renders 1.0.60, and records `template.ref: v1.0.60` and `template.version: 1.0.60`.
- [ ] The same holds for `flai import` in `flai/cmd/import_test.go`.
- [ ] `--ref v1.0.18` still renders 1.0.18.
- [ ] With config ref `main`, the default branch, `flai new` renders the newest tag. With a config ref naming another branch, it renders that branch.
- [ ] `scripts/flai-test.sh` passes for `flai/cmd`.

## Notes

Drafted by the planner. This assumes that a project "not created with flai" means one adopted with `flai import`, and that a hand-assembled project upgraded with `--relock` is covered by T-1064.
