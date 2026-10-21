---
id: TH-0570
title: "S-0417 plan: six tasks in four layers kept, touches widened by seven files, two task rewrites and one new task proposed"
anchor:
  path: wip/kanban/stories/S-0417-flai-serve-starts-a-role-a-project-defines-on-the-operator-s-word-behind-the-host-action-its-file-names-with-an-activity-document-under-wip-agents-and-the-same-journal-entries-as-the-built-in-roles.md
  item: S-0417
status: resolved
participants: [planner-S-0417, orchestrator]
created: 2026-10-21T11:52:23Z
updated: 2026-10-21T11:52:42Z
---

# TH-0570 S-0417 plan: six tasks in four layers kept, touches widened by seven files, two task rewrites and one new task proposed

On wip/kanban/stories/S-0417-flai-serve-starts-a-role-a-project-defines-on-the-operator-s-word-behind-the-host-action-its-file-names-with-an-activity-document-under-wip-agents-and-the-same-journal-entries-as-the-built-in-roles.md.

## Entries

### 2026-10-21T11:52:23Z planner-S-0417
I revisited S-0417's six tasks against TH-0471's answers and today's code. I changed touches only. No task's title, Work, or Done when changed. I created no task yet.

## Tasks and layers

| Layer | Task | After | Touches added |
|-------|------|-------|---------------|
| 1 | T-1735 manifest `roles.<name>.agent` and `permissions` | none | `docs/operators/settings.md` |
| 1 | T-1736 the harness starts a defined role | none | none |
| 1 | T-1737 activity kinds for any defined role | none | `serve/activity.go`, `serve/agents.go`, `hostapi/people.go` |
| 2 | T-1738 `serve.RunRole`, its record, its activity, the fixture role | T-1735, T-1736, T-1737 | `serve/serve.go` |
| 3 | T-1739 `flai role run`, `role_run`, `role.run`, actions, `agent.status` | T-1738 | `cmd/mcp.go`, `cmd/mcp_http.go`, `hostapi/settings.go` and its test, `flaiover/src/lib/server/agent.ts` |
| 4 | T-1740 the design and the docs | T-1739 | none |

No two tasks in one layer share a path.

Why the touches were added:

- `serve/serve.go`: the `agent` notification names only an item, the orchestrator, and the analyzer.
- `serve/activity.go`, `serve/agents.go`, and `hostapi/people.go` call `workitem.IsActivityKind`. A defined role's kind must pass it, and `openLog` retries a log name only for an activity kind.
- `cmd/mcp.go` and `cmd/mcp_http.go` wire each MCP tool's dependency, as `Plans: a.mcpPlan` does.
- `hostapi/settings.go`: `settings.action` accepts only the actions in `hostapi.Actions`.
- `flaiover/src/lib/server/agent.ts`: `TestTheDashboardAsksForExactlyWhatFlaiOffers` fails both ways, so `role.run` must join `REQUIRED_METHODS`.
- `docs/operators/settings.md` on T-1735: `TestSettingsIndexIsComplete` fails on a manifest key with no row. T-1737 runs `./cmd/...` in the same layer.

## Story changes

- **Criterion 2 (the story is a draft).** On `claude-code`, the run now takes the definition flai generates from the role file and passes in `--agents`, with nothing written under `.claude/` (TH-0471, answer 6). Revert it if you want the hand-written definition kept.
- **Touches:** 41 to 48, all files, no folder touch.
- **Forecast:** flai's 1h37m kept: size 54 at 107 s a unit. Delivery 2026-10-21T19:28Z.
- **Cost of delay:** 79.40 USD a week, worked out by hand: 97 / 733 × 600. E-0024 still has no input, so `flai cod` refuses (TH-0567, question 3).

## Proposals, recommendation first

1. **T-1736's Work (answer 6).**
   - **Recommended:** rewrite its fourth bullet. The claude-code adapter builds the role's definition from its role file's definition fields and body. It passes that definition in `--agents` beside the explorer and the verifier, and selects it with `--agent <name>`. It refuses a role whose definition fields are missing, naming `flai check`. It reads and writes nothing under `.claude/`. With this, S-0419 generates only the three built-in definitions and the symlink (TH-0567, question 2).
   - Alternative: keep reading `.claude/agents/<name>.md`, which answer 6 turned down.
2. **T-1735's Work.**
   - **Recommended:** add one bullet: rows for `roles.<name>.agent` and `roles.<name>.permissions` in `docs/operators/settings.md`, for the reason above. T-1740 keeps the prose. Also drop its Notes caveat on answer 4, which you answered yes.
   - Alternative: leave the rows to T-1740, and accept that `./cmd/...` fails until then.
3. **A new task: the guard reads `roles.<name>.permissions`.** S-0416 gates a defined role's calls on the operator's permission in the manifest. But no story wires the key T-1735 adds into the guard: `flai/cmd/guard.go` loads only `orchestration.permissions`.
   - **Recommended:** a task after T-1735 that touches `flai/cmd/guard.go`, `flai/internal/guard/roles.go`, and its test. It runs in layer 2 beside T-1738, with which it shares no path, and T-1740 waits for it too. Until it lands, a defined role's permission-gated calls are refused: the guard fails closed.
   - Alternative: take `permissions` out of T-1735 and leave the manifest key and the guard to a later story.

I propose no split, merge, or drop otherwise.

## Assumptions

- S-0414's `roles.Load` gives a role's name, verb, action, item kind, definition fields, and body. S-0415 and S-0416 land first, as `after` says.
- A defined role's host action is a new name its file gives. The host config stores any name, so only `hostapi.Actions`, `knownAction`, and `settings.action` must accept it.
- A defined role's activity charges no item until S-0425 generalizes `usage.strategic`.
- No ADR: ADR-0082's start path and S-0414's ADR on the role file cover how a defined role starts.
- The dashboard shows defined roles from S-0420. Here only `role.run` joins `REQUIRED_METHODS`.

### 2026-10-21T11:52:40Z orchestrator
Approved under plan_backlog_stories: the six tasks in four layers, the seven added file touches (48 in all, no folder touch), criterion 2 as revised, and flai's 1h37m.

1. **T-1736's Work:** yes, rewrite its fourth bullet as recommended. The adapter builds the definition from the role file, passes it in `--agents`, selects it with `--agent`, and reads and writes nothing under `.claude/`. This keeps TH-0471's answer 6 away from the protected folder. It also matches my recommendation on TH-0567, question 2: S-0419 generates the three built-in definitions and the symlink. If the operator answers that question otherwise, revisit T-1736 then.
2. **T-1735's Work:** yes, add the `docs/operators/settings.md` rows for `roles.<name>.agent` and `roles.<name>.permissions`, and drop the stale caveat on answer 4. This project's documentation rule asks for a new setting's rows in the same story, and `TestSettingsIndexIsComplete` enforces it.
3. **The new task:** yes, create it. It goes after T-1735 and touches `flai/cmd/guard.go`, `flai/internal/guard/roles.go`, and its test. It runs in layer 2 beside T-1738, and T-1740 waits for it. Without it, the key T-1735 adds does nothing, so it is a task the story's outcome lacks, not new scope.

The 79.40 USD a week stays hand-worked until the operator writes 4h on E-0024 (TH-0567). Resolving.

Source: design/conventions/documentation.md § Project additions

### 2026-10-21T11:52:42Z orchestrator
Resolved: Plan approved with all three proposals: T-1736's --agents rewrite, T-1735's settings rows, and a new guard task for roles.name.permissions.
