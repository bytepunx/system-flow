---
id: TH-0214
title: "S-0273 plan: six tasks in four layers, and three assumptions to confirm"
anchor:
  path: wip/kanban/stories/S-0273-flai-test-runs-the-project-s-test-and-lint-tiers-for-a-path-or-a-package-and-answers-pass-or-the-first-failures-as-findings.md
  item: S-0273
status: resolved
participants: [planner-S-0273, alex]
created: 2026-10-06T22:54:29Z
updated: 2026-10-06T23:30:10Z
---

# TH-0214 S-0273 plan: six tasks in four layers, and three assumptions to confirm

On wip/kanban/stories/S-0273-flai-test-runs-the-project-s-test-and-lint-tiers-for-a-path-or-a-package-and-answers-pass-or-the-first-failures-as-findings.md.

## Entries

### 2026-10-06T22:54:29Z planner-S-0273
The plan for S-0273. I drafted six tasks and updated the story's touches, topics, forecast, and `### Planning` notes. I have not changed the goal or the criteria.

The tasks and their layers:

| Layer | Task | What | Waits for |
|------:|------|------|-----------|
| 1 | T-1067 | The manifest's `tests` key: name, command, the `paths` globs that select a tier, `format`, and `all_only`. Validated, writable through `flai manifest set`, and shipped in the template manifest | none |
| 1 | T-1068 | The `flai/internal/verify` package: picks tiers by path, runs them cheapest first, stops at the first failure, and turns go test, vitest, golangci-lint, gofmt, and plain output into findings, five by default | none |
| 2 | T-1069 | `flai test [path\|package]... [--all] [--max n] [--json]`, and this repository's tiers in `system-flow.yaml` | T-1067, T-1068 |
| 2 | T-1070 | The MCP tool `test` (`story`, `paths`, `all`, `max`), answering what `flai test --json` answers | T-1067, T-1068 |
| 3 | T-1073 | The host method `test.run` beside `checks.run`, which builds `flai test --json` | T-1069 |
| 4 | T-1079 | `delegation.md`, `work-management.md`, `tooling.md` and their template copies, the harness prompt, `devex.md`, `flai-cli.md`, and the user guide; the `make` targets checked | T-1067, T-1069, T-1070, T-1073 |

The tasks in each layer share no path, so they can run together.

The figures:

- **Forecast:** 75m, against flai's 40m. The work is a new manifest key, five output parsers, and three surfaces (CLI, MCP, host channel). Delivery 2026-10-07T08:47Z.
- **Cost of delay:** 213 USD a week stands, against flai's 141.36. It is E-0017's 900 shared by the hand turns this story removes (1,449 of about 6,100), not by forecast duration.

The assumptions. Each is my recommendation; reply if you want otherwise.

1. **Where the tiers live.** Criterion 1 says "the matching tiers from `scripts/`". I read that as the tiers today's scripts run. They become a `tests` list in the manifest, one command per tier, so flai can select them by path and parse their output. A whole script gives only pass or fail and a log tail.
   - This repository's tiers mirror `scripts/flai-test.sh` and `scripts/test.sh`, and each script gets a comment saying so. Both scripts and the `make` targets work as before.
   - A project with no `tests` key gets one `plain` tier running `scripts/test.sh`. The template manifest declares its three scripts as tiers.
   - This follows assumption 1 of the epic plan (TH-0176). If you want the tiers to stay as the scripts, say so: T-1067 then shrinks, and findings for those tiers are only plain.
2. **Who may call `test.run`.** It is allowed only while the existing `checks` host action is on, since it runs the project's commands on the host as `checks.run` does. No new action.
3. **Where it runs.** The MCP tool and `test.run` take a story and run in its worktree, and otherwise in the main checkout. With no paths given, they test what the story's branch changed against the main branch.

Other notes:

- `flai/internal/verify` stays a folder touch. It is a new package whose files and fixtures no task can name yet, and S-0270, which waits for this story, extends it.
- `flai/internal/manifest` is narrowed to the four files T-1067 changes.
- Nothing to split, merge, or drop.

### 2026-10-06T23:30:10Z alex
Resolved.
