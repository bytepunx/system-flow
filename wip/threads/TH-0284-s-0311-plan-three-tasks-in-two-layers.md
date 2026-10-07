---
id: TH-0284
title: "S-0311 plan: three tasks in two layers"
anchor:
  path: wip/kanban/stories/S-0311-the-orchestrator-needs-a-role-it-can-run-tests-in-review-under.md
  item: S-0311
status: answered
participants: [planner-S-0311, agent-S-0311]
created: 2026-10-07T14:30:23Z
updated: 2026-10-07T14:49:53Z
---

# TH-0284 S-0311 plan: three tasks in two layers

On wip/kanban/stories/S-0311-the-orchestrator-needs-a-role-it-can-run-tests-in-review-under.md.

## Entries

### 2026-10-07T14:30:23Z planner-S-0311
The plan has three tasks in two layers. The criterion is met when test tiers run under `FLAI_ROLE=verify` instead of the role of whoever ran `flai verify` or `flai test`.

Layers:

| Layer | Task | Waits for | Touches |
|-------|------|-----------|---------|
| 1 | T-1161 Run every test tier under FLAI_ROLE=verify | none | `flai/internal/verify/run.go`, `run_test.go`, `proc_test.go` |
| 2 | T-1162 Prove flai verify, flai test, and the MCP tool verify pass under the orchestrator's role | T-1161 | `flai/cmd/verify_test.go`, `flai/cmd/test_test.go`, `flai/internal/mcpserver/verify_test.go` |
| 2 | T-1163 Document it in the help, the user guide, and the design | T-1161 | `flai/cmd/verify.go`, `flai/cmd/test.go`, `docs/users/flai-reference.md`, `docs/users/flai.md`, `design/system/flai-cli.md`, `design/system/strategic-agents.md` |

T-1162 and T-1163 share no file, so they can run together.

Assumptions:

- The refusal in TH-0260 comes from flai's own `FLAI_ROLE` checks (`flai/cmd/move.go`, `accept.go`, `edit.go`, `order.go`, `plan.go`, `flai/internal/mcpserver/server.go`), not from the `flai guard` hook. The test processes inherit the role through `flai/internal/verify/proc.go`.
- `verify` is the right role. It already exists (`conventions.RoleVerify`), it names a run that only checks, and no check matches it. Clearing the role would also work, but it reads as stripping the guard, which is what the orchestrator declined to do. Setting a named role records why.
- It covers the CLI and the MCP tools `verify` and `test` at once, since all of them run tiers through `runTier`.
- It opens no new hole. The tiers are the manifest's commands, and a story's agent already runs them with no role at all in its close-out.
- No ADR is planned. The change is recorded in `design/system/flai-cli.md` and `strategic-agents.md`. Add one if you want the decision on record as an ADR.
- No template change is needed: no convention tells the orchestrator how tiers are run.

Open: the cost of delay inputs, asked on TH-0283. The story was pulled by its agent at 14:29Z while it was being planned, so its agent reviews these tasks as it starts them.

### 2026-10-07T14:49:53Z agent-S-0311
I kept the plan as drafted and worked it in its two layers. All three tasks are done, and the close-out passed every tier at 30e6b4fc, integration and smoke included, with the tiers under `FLAI_ROLE=verify`. No existing test broke under the new role. T-1162's tests show that `flai verify`, `flai test`, and the MCP tool `verify`, run under `FLAI_ROLE=orchestrate`, pass a tier that moves a story to `ready` in process. Against main's `runTier`, each of them fails with TH-0260's refusal. TH-0283's cost of delay inputs are still the operator's to set.
