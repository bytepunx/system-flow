---
id: T-1081
type: task
nature: remediation
title: flai upgrade without --ref applies the template's newest release, asks which version on a conflict with system-flow.yaml, and never reverts to the lock
status: cancelled
parent: S-0301
owner: alex
created: 2026-10-06T22:52:42Z
updated: 2026-10-06T22:53:26Z
transitions:
  - to: cancelled
    at: 2026-10-06T22:53:26Z
    by: agent-S-0301
stream: S-0301
tags: [cli]
touches: [flai/cmd/upgrade.go, flai/cmd/upgrade_test.go]
after: [T-1057, T-1078]
---
# T-1081 flai upgrade without --ref applies the template's newest release, asks which version on a conflict with system-flow.yaml, and never reverts to the lock

## Work

Criteria 3 and 4. With no `--ref`, `flai upgrade` on a git template takes the template's newest release, unless `system-flow.yaml` asks for a different one, and never falls back to what the lock recorded:

- The source is the manifest's `template.repo` and `template.ref`, as T-1057 leaves them. When `Releases` (T-1078) finds version tags and the ref follows releases (`FollowsReleases`), the target is the latest tag. A ref that names another branch or a commit is used as given, fetched again when it is a branch. With no version tags, the manifest's ref is used, fetched again when it is a branch. A local template directory is unchanged.
- Intent in `system-flow.yaml`: when the lock exists and the manifest's `template.ref` or `template.version` differs from the lock's, the operator edited it. A changed ref names that ref; a changed version names the tag `v<version>`, refused with a clear error naming the tags when there is no such tag.
- When that intent differs from the latest tag, it is a conflict. In a terminal, ask with `a.prompts().Select` which version to apply: the one `system-flow.yaml` asks for, or the latest. Without a terminal, or with `--yes`, refuse with exit 1, naming both and `--ref`, and change nothing.
- The version the project is at is the lock's `template.version` when there is a lock, else the manifest's. That is what "already at template" compares with and what `upgrade.Compute` is given as the version upgraded from. Today an edited `template.version` in the manifest makes the upgrade report "already at" and apply nothing.
- What was chosen is written to `template.ref` and `template.version` in `system-flow.yaml` and to the lock (`writeTemplateFields` with the ref changed), the same as for `--ref` (T-1057). Say which ref and version were chosen and why, in one line of output, and add them to `--json`.
- `--dry-run` resolves the same way and prints the choice. It never prompts: on a conflict it names both versions.

It waits for T-1057, which touches the same lines of `upgrade.go`, and for T-1078, whose `Releases` it calls.

## Done when

- [ ] Tests in `flai/cmd/upgrade_test.go` use a git template with tags, built in a temp directory as the existing git tests are, and cover each of these:
  - Manifest `ref: main` at 1.0.0, with a lock at the same version, upgrades to the newest tag and records it in the manifest and the lock.
  - The reported case: lock `main` at 1.0.0, the manifest's version hand-set to the latest. It upgrades to the latest. It does not report "already at" and does not revert.
  - A manifest version hand-set to an older tag than the latest is a conflict. Without a terminal it exits 1 naming both. With a scripted prompt it applies the chosen one.
  - A branch ref that does not follow releases is fetched again, so a commit added to it after the first clone is applied.
- [ ] Every existing upgrade test passes.
- [ ] `go test -race ./cmd/ -run Upgrade` passes.

## Notes
- 2026-10-06T22:53:26Z: moved to cancelled: duplicate of the planner's draft for S-0301, written before I saw it; the planner's tasks are kept and edited
