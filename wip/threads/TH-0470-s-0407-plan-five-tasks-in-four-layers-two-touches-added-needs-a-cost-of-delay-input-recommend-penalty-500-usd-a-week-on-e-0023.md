---
id: TH-0470
title: "S-0407 plan: five tasks in four layers, two touches added; needs a cost of delay input (recommend penalty 500 USD a week on E-0023)"
anchor:
  path: wip/kanban/stories/S-0407-flai-guard-judges-an-opencode-agent-s-calls-through-a-plugin-on-tool-execute-before-and-the-opencode-adapter-states-the-guard-only-where-the-plugin-can-tell-a-sub-agent-s-call.md
  item: S-0407
status: open
participants: [planner-S-0407, orchestrator]
created: 2026-10-16T05:26:41Z
updated: 2026-10-16T05:28:16Z
---

# TH-0470 S-0407 plan: five tasks in four layers, two touches added; needs a cost of delay input (recommend penalty 500 USD a week on E-0023)

On wip/kanban/stories/S-0407-flai-guard-judges-an-opencode-agent-s-calls-through-a-plugin-on-tool-execute-before-and-the-opencode-adapter-states-the-guard-only-where-the-plugin-can-tell-a-sub-agent-s-call.md.

## Entries

### 2026-10-16T05:26:41Z planner-S-0407
The E-0023 planner drafted S-0407's five tasks. I revisited them against the story. They are kept as written, with two touches added. The cost of delay input is still missing.

## 1. Cost of delay input (blocks the value)

`flai cod S-0407` refuses: neither the story nor E-0023 has inputs.

- **Recommended:** `penalty_per_week: 500` on E-0023, as TH-0467 question 1 recommends: the figure E-0019 has. Every E-0023 story then shares it by duration, and I write S-0407's value on my next run.
- Alternatives: the same input on S-0407 alone, which would value this story at 500 USD a week by itself; or `time_lost_per_cycle: 1h` on E-0023 (150 USD a week).

## Tasks and layers

| Layer | Task | After | Touches |
|-------|------|-------|---------|
| 1 | T-1690: `opencode` reader in `flai/internal/guard`, `flai guard --harness opencode` | none | `guard/opencode.go`, `opencode_test.go`, `guard.go`, `subagents.go`, `subagents_test.go`, `flai/cmd/guard.go`, `flai/cmd/guard_test.go`, `docs/users/flai-reference.md` |
| 2 | T-1691: `flai-guard.js` plugin, in this repository and the template | T-1690 | `.opencode/plugins/flai-guard.js`, `template/root/.opencode/plugins/flai-guard.js`, `template/CHANGELOG.md` |
| 2 | T-1692: `Guard` capability on the opencode adapter, set from S-0406's finding | T-1690 | `harness/opencode.go`, `opencode_test.go`, `harness_test.go` |
| 3 | T-1693: integration test, sub-agent refused and the agent let through | T-1691, T-1692 | `flai/tests/integration/opencode_test.go` |
| 4 | T-1694: operator, user, and design docs | T-1693 | `docs/operators/index.md`, `docs/users/flai.md`, `design/system/flai-cli.md`, `design/system/agent-adapters.md` |

T-1691 and T-1692 share no path, so they can run together.

## What changed

- **Added `flai/cmd/guard.go` and `docs/users/flai-reference.md`** to T-1690 and to the story. The guard's help says claude-code is "the only reader there is", and the `--harness` flag lists `guard.Harnesses()`. The generated reference has to be regenerated with the help in the same task.
- **Forecast:** 45m kept. flai's figure is now 39m, up from 36m, because of the two added touches. The raise for the plugin and the real-binary check drops from a third to about a sixth. Delivery is 2026-10-16T13:40Z.

## Assumptions

1. S-0406 is accepted first. It creates `flai/internal/harness/opencode.go`, `flai/tests/integration/opencode_test.go`, and § Harnesses in `docs/operators/index.md`, and it records the sub-agent finding that decides criterion 3.
2. `flai serve` needs no change. `startAllowed` in `flai/internal/serve/agents.go` reads `Capabilities().Guard` for any harness, so `Guard: true` lifts the role refusal by itself. A role may still be refused for role models until S-0409.
3. A read-only sub-agent on OpenCode is recognised by its agent name (explorer, verifier), as on Claude Code. S-0409 generates those agents.
4. No ADR is needed. ADR-0130 already decides that the guard is registered as an OpenCode plugin, and ADR-0060 decides that it fails open.
5. `subagents.go` may not change, because a reader hands sub-agent starts and stops to the existing record as `Hook.SubAgent`. The touch is kept because the story declared it.

## Proposals (not made)

- **T-1692 `## Work`:** add a line that tests `startAllowed` with an opencode story that has roles, so criterion 3's "no longer refused for the guard" is shown in a test. No touch is needed in `flai/internal/serve`. I did not change the text because T-1692 is another planner's task.
- **T-1694 `after`:** it could wait for T-1691 and T-1692 instead of T-1693 and run beside the integration test, since the two share no path. That cuts the plan to three layers. I left it as it is, because its `## Work` says the docs describe the tested behaviour.

### 2026-10-16T05:26:48Z planner-S-0407
A correction to the T-1692 proposal. A test of `startAllowed` would go in `flai/internal/serve/agents_test.go`, which T-1692 does not touch. If you agree to the proposal, I add that file to T-1692's touches and to the story's. If not, T-1692's harness test, which asserts `Guard: true`, stands for criterion 3, together with the existing `startAllowed` tests.

### 2026-10-16T05:26:55Z orchestrator (recommendation)
**Recommendation:** answer TH-0467 question 1 with `penalty_per_week: 500` on E-0023, as both planners recommend. Every E-0023 story then shares it by duration, and this planner writes S-0407's value on its next run.

## Why I do not set an input on S-0407 myself

- Under `plan_backlog_stories` I may set inputs only on the story. The input recommended here belongs on the epic. Its alternative, 500 on S-0407 alone, would value this one story at the whole epic's figure.
- flai refuses inputs once a story or its epic has any. An input on S-0407 now could keep you from setting E-0023's. The epic-level input is money and yours, and it is already asked on TH-0467.

## The plan: approved

- **The tasks fit.** Five tasks in four layers. T-1691 and T-1692 share no path. The two added touches (`flai/cmd/guard.go` and the reference) are right, because the guard's help changes.
- **The forecast is flai's.** 45m kept against flai's 39m, with the reason stated.
- **Proposal 1: yes.** Add a `startAllowed` test with an opencode story that has roles to T-1692's `## Work`, so criterion 3 is shown in a test. This changes another planner's task, so it is approved here.
- **Proposal 2: leave it.** T-1694 waits for T-1693, so the docs describe tested behaviour.

The thread stays open for the cost of delay input only.

Source: design/adrs/0119-with-plan-backlog-stories-on-the-orchestrator-asks-the-planner-for-each-backlog.md § Decision

### 2026-10-16T05:28:16Z planner-S-0407
Done, as approved:

- T-1692's `## Work` now has the `startAllowed` test.
- `flai/internal/serve/agents_test.go` is in T-1692's touches and the story's.
- The forecast stays at 45m against flai's new 40m, delivered 2026-10-16T13:42Z.

T-1694 is unchanged. The cost of delay value waits for the operator's input on E-0023 (TH-0467). I write it once `flai cod S-0407` answers.
