---
id: TH-0392
title: "S-0360 is a complete draft, but it adds scope to E-0019: finalize it or cancel it"
anchor:
  path: wip/kanban/stories/S-0360-flai-serve-checks-each-new-claude-version-through-the-provider-a-project-s-agents-name-so-a-host-that-reaches-claude-only-through-a-gateway-is-not-reported-failing.md
  item: S-0360
status: open
participants: [orchestrator]
created: 2026-10-09T14:17:47Z
updated: 2026-10-09T14:17:47Z
---

# TH-0392 S-0360 is a complete draft, but it adds scope to E-0019: finalize it or cancel it

On wip/kanban/stories/S-0360-flai-serve-checks-each-new-claude-version-through-the-provider-a-project-s-agents-name-so-a-host-that-reaches-claude-only-through-a-gateway-is-not-reported-failing.md.

## Entries

### 2026-10-09T14:17:47Z orchestrator
**Recommendation:** finalize S-0360. Your 500 USD a week penalty on E-0019 (TH-0375) rests on the subscription cap. The planner's reason for S-0360 is the same case: on a host that reaches Claude only through a gateway, the check of each new `claude` version calls Anthropic directly and would fail.

`flai promote --drafts` lists it as complete:

- Goal, criteria, and two tasks (T-1426, T-1427).
- Touches: `flai/internal/serve/claudecheck.go`, its test, and `design/system/flai-cli.md`.
- Forecast 20m, value 22.12 USD a week.
- After S-0351.

These are consistent with its criteria.

Why I did not finalize it: the planner proposed it on TH-0386 as an addition to E-0019, a change of scope. You finalized the other eleven E-0019 stories on 2026-10-08 and left this one as a draft. A change of scope is yours to decide, so I leave it as a draft.

If you do not want it, cancel it. It holds nothing meanwhile: it waits for S-0351 either way.
