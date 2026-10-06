---
id: S-0301
type: story
nature: remediation
title: flai upgrade consistently pulls an old template no matter what
status: in-progress
owner: alex
created: 2026-10-06T22:46:13Z
updated: 2026-10-06T22:54:09Z
transitions:
  - to: ready
    at: 2026-10-06T22:46:14Z
    by: alex
  - to: in-progress
    at: 2026-10-06T22:50:14Z
    by: agent-S-0301
tags: [cli]
topics: [template]
touches: [flai/cmd/upgrade.go, flai/cmd/upgrade_test.go, flai/cmd/new.go, flai/cmd/new_test.go, flai/cmd/import.go, flai/cmd/import_test.go, flai/internal/template, design/adrs, design/system/template.md, design/system/project-manifest.md, design/system/flai-cli.md, docs/users/flai.md, docs/users/flai-reference.md, docs/contributors/template.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
forecast:
  duration: 2h30m
  delivery: 2026-10-07T01:30:00Z
  basis: "flai forecast's 12m sized the story at one touch and four criteria; the plan is six tasks in three layers, with new tag resolution, a cache refetch, an upgrade prompt, an ADR and docs, so it is sized like the 1h30m to 1h45m six-task stories plus a margin for git test fixtures."
  by: planner-S-0301
  at: 2026-10-06T22:51:34Z
---
# S-0301 flai upgrade consistently pulls an old template no matter what

## Goal

### flai upgrade

flai upgrade should default to pulling the latest available version tag from the template repository it is configured for.

If there is a ref passed, it should overwrite the lock file.
If there is no ref passed and the operator either set a different version in system-flow.yaml or there is a newer version, flai should prompt the user about their intent rather than silently reverting back to the template version in the lockfile.

### Description

In another project, the template defaulted to 1.0.18 (not 1.0.60 which was the latest at the time of this writing). When I change the version manually in system-flow.yaml, it reverts (possibly due to the version in the lock file). But even when providing a version reference to the tag 1.0.60, I still see it revert to 1.0.18.

## Acceptance criteria
- [ ] flai correctly updates system-flow files to reflect the new version supplied in the --ref argument before completing the upgrade
- [ ] flai defaults to the latest available version of the template when creating a new project or upgrading an existing project not created with flai
- [ ] flai assumes that `flai upgrade` should find and use the latest available tag, update the system-flow files, and then perform the upgrade
- [ ] if there is a conflict, flai's CLI must prompt the operator with the version they intend before proceeding rather than reverting to the version found in the lock file.

## Tasks
- T-1057 flai upgrade --ref wins over the manifest's template.ref and records the ref and version it applied
- T-1060 The template package lists a git template's version tags, names the newest, and fetches a cached branch ref again
- T-1062 An ADR and the template design say which template ref flai new, flai import, and flai upgrade take, and when upgrade asks
- T-1064 flai upgrade with no --ref takes the newest template tag and asks which version to apply when the manifest, the lock, and the newest tag disagree
- T-1065 flai new and flai import default to the template's newest version tag and record it as the project's template.ref
- T-1066 The flai docs and the CLI design say flai new, import, and upgrade take the newest template tag and when upgrade asks
- T-1078 flai/internal/template finds the template repository's version tags, its newest release, and whether a ref follows releases
- T-1081 flai upgrade without --ref applies the template's newest release, asks which version on a conflict with system-flow.yaml, and never reverts to the lock
- T-1084 flai new and flai import without --ref make the project at the template's newest release
- T-1086 An ADR, the design, and the guides say new, import, and upgrade follow the template's releases

## Notes

### Planning

Planned by planner-S-0301 on 2026-10-06. The plan thread is on S-0301.

**Diagnosis.** The code shows two causes.

- `flai/cmd/upgrade.go`, lines 54 to 56, replaces `--ref` with the manifest's `template.ref` whenever `--template` is not given, so `--ref 1.0.60` is thrown away.
- `template.Source.Ensure` clones each repo@ref once into the cache and never fetches it again. A project on `main`, the config's default, keeps rendering whatever `main` was on the first clone, here 1.0.18. Editing `template.version` in `system-flow.yaml` then reads as "older than the template", and the upgrade "reverts" to the cached version.

No code lists the template's tags today.

**Touches.**

- **Declared:** `flai/cmd`. It is kept, as the planner keeps every declared touch. It is a folder touch, and under ADR-0096 the tasks' file touches below it replace it in the story's claim. The plan thread proposes dropping it.
- **Layout:** read from the code.
  - `flai/cmd/upgrade.go` and `upgrade_test.go`: the bug and the prompt.
  - `flai/cmd/new.go` (`resolveTemplate`) and `new_test.go`.
  - `flai/cmd/import.go` and `import_test.go`.
  - `flai/internal/template/source.go`, plus the new `tags.go`, `tags_test.go` and `source_test.go`.
  - `flai/internal/config/config.go` and `config_test.go`: the default `template.ref` is `main`.
- **Design:**
  - `design/system/template.md`, "Upgrading a project", which today says that upgrade never prompts.
  - `design/system/project-manifest.md`: `template.ref`.
  - `design/system/flai-cli.md`: the rows for `flai new` and `flai upgrade`.
  - `docs/users/flai.md`, "Upgrade to a newer template".
  - `docs/users/flai-reference.md`, which `scripts/flai-reference.sh` regenerates.
  - `docs/contributors/template.md`, which warns that a branch ref goes stale.
  - A new ADR and `design/adrs/README.md`.
- **Co-change:** `flai touches suggest` ranks by co-change with `flai/cmd` as a whole. Its top hits, `docs/users/flai.md`, `design/system/flai-cli.md` and `docs/users/flai-reference.md`, are taken. The rest are other commands' files that this story does not reach.
- **Folder touches kept:**
  - `flai/cmd`, as declared.
  - `design/adrs`, because the ADR's number is not known until `flai adr new` runs. It is in `claims.shared`, so it holds no ready story.

**Forecast.** The duration is 2h30m and delivery is 2026-10-07T01:30Z, adjusted from flai's 12m.

- flai sized the story from one declared touch and four criteria, which gave size 5 in the small band.
- The plan is six tasks in three layers, with new code for tag listing, a cache refetch and an interactive prompt, plus an ADR and docs.
- Six-task stories here have been forecast at 1h30m to 1h45m. Git-backed test fixtures for tags add margin.
- Delivery counts from the pull at 22:50Z.

**Cost of delay.** The story has no inputs and no epic, so `flai cod` refuses. The inputs are the operator's, asked for on the plan's cost of delay thread, with `time_lost_per_cycle: 30m` recommended. The value is written once they are answered.
