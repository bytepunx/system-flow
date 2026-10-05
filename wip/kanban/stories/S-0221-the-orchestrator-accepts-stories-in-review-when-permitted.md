---
id: S-0221
type: story
nature: feature
title: The orchestrator accepts stories in review when permitted
status: ready
parent: E-0016
owner: arobson
created: 2026-10-02T11:54:17Z
updated: 2026-10-05T07:09:06Z
transitions:
  - to: ready
    at: 2026-10-05T04:41:22Z
    by: alex
tags: [flai]
touches: [flai/internal/harness, flai/internal/hostapi, flai/internal/preview, flai/internal/mcpserver, design/adrs, flai/cmd/accept.go, flai/internal/guard, flaiover/src/lib/components/Review.svelte, flaiover/src/routes/items, design/system/workflow.md, design/system/strategic-agents.md, docs/users/flai.md, docs/users/flaiover.md, flai/cmd/accept_orchestrator_test.go, flai/cmd/guard.go, flaiover/src/lib/components/Review.svelte.test.ts, ".claude/agents/orchestrator.md", template/root/.claude/agents/orchestrator.md, design/conventions/strategic-agents.md, template/root/design/conventions/strategic-agents.md, template/CHANGELOG.md, template/template.yaml, design/system/flai-cli.md, design/system/flaiover-dashboard.md, docs/users/flai-reference.md]
after: [S-0218]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
cost_of_delay:
  value: 88.82
  by: planner-S-0221
  at: 2026-10-05T04:47:12Z
forecast:
  duration: 1h30m
  delivery: 2026-10-05T11:56:00Z
  basis: "Its own forecast of 1h30m; 4th in the pull order with an in-progress limit of 3, behind S-0218, S-0219 and S-0220."
  by: flai
  at: 2026-10-05T07:09:06Z
---
# S-0221 The orchestrator accepts stories in review when permitted

## Goal

ADR-0032 and `item_move` make acceptance the operator's alone. With `accept_reviews` on, the orchestrator accepts a story in review when the verifier's run passed, every criterion is ticked, the diff stays within the story's touches, and no thread on it is open; otherwise it leaves it with a thread saying what is missing.

## Acceptance criteria
- [ ] An ADR refines ADR-0032: acceptance by the orchestrator under `accept_reviews`, with the conditions above, recorded as `by: orchestrator`
- [ ] The hostapi `accept.run` and `flai accept` take the orchestrator's identity when the permission is on; `flai guard` refuses otherwise
- [ ] Before accepting it runs the acceptance preview and refuses on any blocker, and it never accepts a story whose criteria it cannot check against the diff
- [ ] Each acceptance is logged with the evidence; the review page shows who accepted
- [ ] Tests cover an acceptance, a refusal on an open thread, and the permission off

## Tasks
- T-0898 An ADR refines ADR-0032: the orchestrator accepts a story in review under accept_reviews, and the design says when
- T-0900 The acceptance preview reports the orchestrator's four conditions as blockers
- T-0903 flai guard lets the orchestrator accept only under accept_reviews, and item_move's refusal names the way
- T-0905 The orchestrator's prompt, agent file, and convention say how it reviews and accepts a story
- T-0907 flai accept and hostapi accept.run take the orchestrator's identity under accept_reviews, refuse on any blocker, and record the evidence
- T-0909 The review page and the story page show who accepted a story, and the orchestrator's evidence
- T-0911 The CLI design, the dashboard design, and the users' guides describe the orchestrator's acceptance

## Notes

### Planning

Planned by planner-S-0221 on 2026-10-05. The plan thread lists the tasks, their layers, and the assumptions.

Touches:

- Declared, kept: `flai/internal/harness`, `flai/internal/hostapi`, `flai/internal/preview`, `flai/internal/mcpserver`, `design/adrs`, `flai/cmd/accept.go`, `flai/internal/guard`, `flaiover/src/lib/components/Review.svelte`, `flaiover/src/routes/items`, `design/system/workflow.md`, `design/system/strategic-agents.md`, `docs/users/flai.md`, `docs/users/flaiover.md`.
- `flai/cmd/accept_orchestrator_test.go`: from the layout. The acceptance tests live in `flai/cmd` as `accept_*_test.go`, and `preview` has none of its own.
- `flai/cmd/guard.go`: from the layout. The guard command sits beside `flai/internal/guard`, and S-0209 and S-0255 changed both.
- `flaiover/src/lib/components/Review.svelte.test.ts`: from the layout, the component's tests.
- `.claude/agents/orchestrator.md`, `template/root/.claude/agents/orchestrator.md`: from the design. S-0218 creates the orchestrator's agent file, and its review steps belong there.
- `design/conventions/strategic-agents.md` and its template copy, `template/CHANGELOG.md`, `template/template.yaml`: from the design and the layout. "As the orchestrator" must name the acceptance rule, and a convention change lands in the template first with a version bump, as S-0210 did.
- `design/system/flai-cli.md` (co-change 39%), `design/system/flaiover-dashboard.md` (29%), `docs/users/flai-reference.md` (15%): from co-change. `flai accept` gains flags and blockers, and the review page changes.
- Not taken from `flai touches suggest`:
  - `docs/operators/index.md` (16%) and `docs/operators/settings.md` (7%). S-0218 documents the permissions and S-0229 their settings; this story adds no setting.
  - `design/issues/summary.md` (11%). The close-out records issues itself.
  - `flai/cmd/serve_actions.go` (6%). `accept.run` is in `flai/internal/hostapi/writes.go`.

Forecast: 1h30m, delivery 2026-10-05T14:23Z.

- `flai forecast` gave 40m at size 18. With the 25 touches above, it gives 1h7m: 132 s per unit of size, times size 30.
- Raised to 1h30m, the figure the story already carried. It has seven tasks in four layers, across the accept gate, the guard, the orchestrator's prompt and convention in the template, the dashboard, and an ADR. Comparable E-0016 stories took 65m (S-0208), 98m (S-0209), and 210m (S-0255, guard rules).
- Delivery is flai's 13:48Z, seventh in the pull order, plus the 23 extra minutes times the cycle factor of 1.52.

Cost of delay: 88.82 USD/week, from `flai cod`, kept as computed.

- The story has no inputs of its own. Its value is its share of E-0016's 1500 USD a week, which comes from the operator's input of 10h lost per cycle.
- The share is 1h30m of the 25h20m forecast over the epic's 17 open stories without inputs. It replaces the 73.17 the epic's planner wrote on 2026-10-04, which was a share of a different total.
