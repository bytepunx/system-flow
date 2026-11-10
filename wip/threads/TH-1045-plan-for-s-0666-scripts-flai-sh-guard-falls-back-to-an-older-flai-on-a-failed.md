---
id: TH-1045
title: "Plan for S-0666: scripts/flai.sh guard falls back to an older flai on a failed build and exits 127 with none, five tasks in three layers"
anchor:
  path: wip/kanban/stories/S-0666-this-repository-s-guard-hooks-let-a-call-through-when-scripts-flai-sh-exits.md
  item: S-0666
status: open
participants: [planner-S-0666]
created: 2026-11-10T10:36:50Z
updated: 2026-11-10T10:36:50Z
---

# TH-1045 Plan for S-0666: scripts/flai.sh guard falls back to an older flai on a failed build and exits 127 with none, five tasks in three layers

On wip/kanban/stories/S-0666-this-repository-s-guard-hooks-let-a-call-through-when-scripts-flai-sh-exits.md.

## Entries

### 2026-11-10T10:36:50Z planner-S-0666
This is the first plan for S-0666. The start prompt said the operator had written on a thread called "asked", but no such thread exists, and S-0666 had no thread, touches, forecast, or tasks before this run. I planned it as new.

## Proposed fix

When `scripts/flai.sh guard` fails to build, it runs the existing `bin/flai`, older than the sources, or else the installed `flai`, each with a warning on standard error. With neither, it exits 127, which the hooks refuse.

- Other commands still fail on a failed build.
- The text of the hooks in `.claude/settings.json` stays as it is.
- This repository's `.opencode/plugins/flai-guard.js` also runs `scripts/flai.sh`, and it passes an exit 127 today. It gets the same 126 and 127 refusal.

## Tasks

| Layer | Task | Touches | After |
|-------|------|---------|-------|
| 1 | T-3061 ADR refining ADR-0235 | `design/adrs/README.md` (and the new ADR, under the story's `design/adrs` touch) | none |
| 2 | T-3066 `scripts/flai.sh` guard fallback, with cases in `TestTheGuardHooksRefuseACallWhenTheGuardCannotStart` using a fake `go` that exits 1 | `scripts/flai.sh`, `flai/cmd/guard_test.go` | T-3061 |
| 2 | T-3069 OpenCode plugin refuses on 126 and 127, tested in `TestFlaiGuardPluginUnderNode` | `.opencode/plugins/flai-guard.js`, `flai/tests/integration/opencode_test.go` | T-3061 |
| 3 | T-3070 docs | `scripts/README.md`, `docs/users/flai.md`, `design/system/agent-adapters.md`, `docs/operators/index.md`, `design/system/strategic-agents.md` | T-3066, T-3069 |
| 3 | T-3071 close I-0275 | the issue, `design/issues/summary.md` | T-3066, T-3069 |

## Figures

- Forecast: 45m. `flai forecast` gave 26m; I raised it for the three layers worked in order and the integration and smoke tiers that the `flai/` touches select.
- Cost of delay: none yet. TH-1039 asks for an input and recommends `time_lost_per_cycle: 5m`.

## Assumptions

1. A new ADR is needed. ADR-0235 says a failed build fails open, and that `scripts/flai.sh` never runs an older installed flai where `go` is present. This fix reverses both, for `guard` only.
2. The fallback covers `guard` only, so an agent working on flai still sees its build fail.
3. Running a stale `bin/flai` as the guard is acceptable. ADR-0235 already accepts that when `go` is absent.
4. The template's hooks and plugin need no change. They run `flai` from `PATH`, so a missing flai is already exit 127 or a spawn error there.
5. The story keeps one folder touch, `design/adrs`, because no task can name the new ADR's number yet.

## Proposals

- T-3069, the OpenCode plugin, goes slightly beyond I-0275's wording. Its claude-code agents never run that plugin. I kept it so both guards in this repository judge exit 127 the same way. It can be dropped, or split into its own story, if you would rather keep S-0666 to the Claude Code hooks.
