---
id: TH-1015
title: "Plan for S-0708: rename the two no-network gateway tests out of the TestGateway prefix and guard the prefix with a test; three tasks in two layers"
anchor:
  path: wip/kanban/stories/S-0708-testgatewaystreamisjudged-and-testgatewaymissingnameswhatismissing-run-in-no.md
  item: S-0708
status: resolved
participants: [planner-S-0708, orchestrator]
created: 2026-11-10T10:28:27Z
updated: 2026-11-10T10:29:39Z
---

# TH-1015 Plan for S-0708: rename the two no-network gateway tests out of the TestGateway prefix and guard the prefix with a test; three tasks in two layers

On wip/kanban/stories/S-0708-testgatewaystreamisjudged-and-testgatewaymissingnameswhatismissing-run-in-no.md.

## Entries

### 2026-11-10T10:28:27Z planner-S-0708
**Recommended remedy:** keep the `-skip '^TestGateway'` in the go-test tier (`system-flow.yaml`) and in `scripts/integration.sh` as it is. Reserve the `TestGateway` prefix for the paid checks, then:

- Rename `TestGatewayStreamIsJudged` and `TestGatewayMissingNamesWhatIsMissing` in `flai/internal/serve/gateway_test.go` out of the prefix, as S-0685 did with `TestCreditRefusalIsFoundInAGatewayRun`.
- Add a guard test, named outside the prefix. It fails on any `TestGateway*` function that `scripts/gateway-smoke.sh` does not run by name. Run against the current names, it reproduces I-0334.

**Alternative not taken:** narrow the skip to `^TestGateway(OpenRouter|LiteLLM)$`, as the issue also suggests. That changes the manifest, `scripts/integration.sh`, and the pin in `flai/internal/verify/select_test.go`. Worse, a paid check added later under the prefix would run in every close-out and batch and spend money on the gateway's key. Say so here if you want it instead, and I will replan.

### Tasks and layers

| Layer | Task | Touches | After |
| --- | --- | --- | --- |
| 1 | T-2961: guard test first (fails on the old names), then rename the two tests | `flai/internal/serve/gateway_test.go` | none |
| 1 | T-2962: `design/conventions/code-quality.md` (below the marker) and `docs/contributors/index.md` state the prefix rule | those two files | none |
| 2 | T-2964: `flai issue close I-0334` naming the tests and the rule | the I-0334 file, `design/issues/summary.md` | T-2961, T-2962 |

T-2961 and T-2962 share no path, so they can run together.

### Figures

- Forecast: 45m, delivery 2026-11-10T15:36:00Z. `flai forecast` gave 16m. I raised it because a close-out for a story that touches `flai/` also runs integration and smoke, about 12 minutes.
- Cost of delay: 25 USD a week, from flai's input of 10m lost per cycle. It stands as computed.

### Assumptions

- The `TestGateway` prefix keeps meaning a paid check. D-0451's advice to pass `-skip '^TestGateway'` when running the tests by hand still holds.
- No `after` on another story, and no folder touch.
- The co-change paths `flai touches suggest` listed (5 percent or less each) are not changed by this remedy, so none were added.

### 2026-11-10T10:28:46Z orchestrator
Approved as written: T-2961 and T-2962 together, then T-2964; 45m, 25 USD a week.

- The remedy fits the goal: I-0334 recommends no fix, and the story asks for one proposed from its instance. Keeping the `^TestGateway` skip and reserving the prefix for paid checks means no later paid check runs in a close-out by accident.
- T-2961 writes the guard test first, and it fails on the current names, which reproduces I-0334 for criterion 1. T-2964 closes I-0334, criterion 2.
- T-2962 changes `code-quality.md` below its marker, a project rule, so no template change is owed.
- T-2961 and T-2962 share no path. Every task touch is in the story's touches, and none is a folder.
- The forecast is raised from flai's 16m with a stated reason, the close-out's integration and smoke tiers. The value is `flai cod`'s from flai's input.

S-0708 is finalized: `flai promote --drafts` listed it as complete.

Source: design/conventions/strategic-agents.md#As the orchestrator

Source: design/conventions/strategic-agents.md § As the orchestrator

### 2026-11-10T10:29:39Z orchestrator
Resolved: Plan approved as written; S-0708 is finalized.
