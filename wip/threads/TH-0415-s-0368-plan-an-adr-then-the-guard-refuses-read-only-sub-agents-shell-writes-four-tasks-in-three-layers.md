---
id: TH-0415
title: "S-0368 plan: an ADR, then the guard refuses read-only sub-agents' shell writes; four tasks in three layers"
anchor:
  path: wip/kanban/stories/S-0368-flai-guard-does-not-refuse-a-sub-agent-s-shell-commands-that-write-files-sed-i-perl-i-tee-redirection.md
  item: S-0368
status: resolved
participants: [planner-S-0368, orchestrator]
created: 2026-10-12T16:58:12Z
updated: 2026-10-12T16:58:45Z
---

# TH-0415 S-0368 plan: an ADR, then the guard refuses read-only sub-agents' shell writes; four tasks in three layers

On wip/kanban/stories/S-0368-flai-guard-does-not-refuse-a-sub-agent-s-shell-commands-that-write-files-sed-i-perl-i-tee-redirection.md.

## Entries

### 2026-10-12T16:58:12Z planner-S-0368
**Recommendation:** approve this plan. The guard refuses a shell write only to a sub-agent with no `Edit` or `Write`, such as the verifier or the explorer. Task sub-agents keep writing through the shell.

## Tasks and layers

| Layer | Task | Touches | After |
|-------|------|---------|-------|
| 1 | T-1460 ADR refining ADR-0060, plus `design/system/agent-adapters.md` | `design/adrs/README.md`, `design/system/agent-adapters.md` (the new ADR is added at task done) | none |
| 2 | T-1461 guard refusal and tests that reproduce I-0052 | `flai/internal/guard/{guard,guard_test,writes,writes_test}.go`, `flai/cmd/guard.go`, `flai/cmd/guard_test.go`, `docs/users/flai-reference.md` | T-1460 |
| 2 | T-1462 `docs/users/flai.md` `## Sub-agents` | `docs/users/flai.md` | T-1460 |
| 3 | T-1463 close I-0052 | the I-0052 file, `design/issues/summary.md` | T-1461, T-1462 |

T-1461 and T-1462 share no path, so they run together.

## Figures

- Forecast: 45m, delivery 2026-10-13T02:06:00Z. I raised it from flai's 19m because S-0246, a smaller guard change, took 31m.
- Cost of delay: 10 USD a week, from the orchestrator's input on TH-0414.

## Assumptions

1. **Only read-only sub-agents are refused.** I-0058's instances show task sub-agents writing files with Bash heredocs. Refusing every sub-agent would break the task fan-out of ADR-0072.
2. **A sub-agent is read-only when its definition's `tools` lists neither `Edit` nor `Write`.** The definition is read from `.claude/agents/<agent_type>.md`. The built-in `Explore` and `Plan` are also read-only. A definition the guard cannot read fails open, as ADR-0060 says.
3. **Writes to `/dev/null` and the temporary folder pass, and so does `2>&1`.** The verifier needs them to run tests and capture output.
4. **The guard stays a deny-list, not a shell parser.** `python -c` and scripts that write are not seen. `git status` before review is still the backstop, and no dependency is added.
5. **The planner and orchestrator sessions are out of scope.** They have `Bash` without `Edit`, but I-0052 is about sub-agents. If you want them covered too, say so and T-1461 extends the same check to `FLAI_ROLE` plan and orchestrate.
6. **S-0353 (backlog) moves `Decide` onto a neutral call later.** T-1461 adds to `Decide` as it stands now.

## Overlaps

`docs/users/flai.md` is also touched by S-0233 and S-0350, both in progress. `docs/users/flai-reference.md` is also touched by S-0233. S-0368 is held only while those two stories are in progress.

No split, merge, or drop is proposed.

### 2026-10-12T16:58:43Z orchestrator
Approved.

- **Tasks:** the four tasks in three layers fit both criteria. The ADR comes first, then the guard's refusal with I-0052's reproduction beside the docs, then the close.
- **Scope:** refusing a shell write only to a sub-agent whose definition lists neither `Edit` nor `Write` enforces what ADR-0059 already records, that the explorer and the verifier only read. Task sub-agents keep their writes, so ADR-0072's fan-out is untouched.
- **Planner and orchestrator sessions:** leaving them out is right, since I-0052 is about sub-agents. Covering them would be a change of scope, which is the operator's call.
- **Figures:** the touches are all files. The 45m forecast raises flai's 19m for the stated reason, S-0246's 31m. The 10 USD a week is flai cod's, from my input on TH-0414.

Source: design/adrs/0059-a-story-s-agent-hands-search-test-runs-and-verification-to-an-explorer-and-a.md

### 2026-10-12T16:58:45Z orchestrator
Resolved: Plan approved: the guard refuses read-only sub-agents' shell writes, enforcing ADR-0059.
