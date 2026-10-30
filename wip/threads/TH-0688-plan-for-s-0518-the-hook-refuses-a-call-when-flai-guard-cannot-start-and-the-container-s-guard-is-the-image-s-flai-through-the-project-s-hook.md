---
id: TH-0688
title: "Plan for S-0518: the hook refuses a call when flai guard cannot start, and the container's guard is the image's flai through the project's hook"
anchor:
  path: wip/kanban/stories/S-0518-flai-guard-does-not-run-as-the-hook-inside-the-agent-container-and-the-hook-lets-every-call-through-when-the-guard-cannot-start.md
  item: S-0518
status: resolved
participants: [planner-S-0518, orchestrator, flai serve]
created: 2026-10-28T13:35:01Z
updated: 2026-10-30T13:37:28Z
---

# TH-0688 Plan for S-0518: the hook refuses a call when flai guard cannot start, and the container's guard is the image's flai through the project's hook

On wip/kanban/stories/S-0518-flai-guard-does-not-run-as-the-hook-inside-the-agent-container-and-the-hook-lets-every-call-through-when-the-guard-cannot-start.md.

## Entries

### 2026-10-28T13:35:01Z planner-S-0518
The plan, by planner-S-0518 on 2026-10-28. S-0518 stays a draft for you to finalize.

Cause. The guard exits only 0 or 2; the hook's shell line in both `.claude/settings.json` files, `[ $? -eq 2 ] || exit 0`, allows every other exit, a missing binary included. Inside `flai-agent-claude-code:local` this repository's hook runs `scripts/flai.sh guard`, whose `go build` exits 127 for want of `go`, so S-0513's run passed every call unguarded. A project made from the template is guarded in the image today, since its hook runs the `flai` the image carries.

Remediation, in two parts. The hook refuses the call, exit 2 with `flai guard did not run (exit N): <output>`, when the guard command exits 126 or 127, and keeps allowing on 0 and on any other exit as ADR-0060 decided. Inside the container the guard is the image's flai through the project's own hook: `scripts/flai.sh` runs the installed flai, with a warning, when `go` is absent and `bin/flai` is not built.

Tasks and layers:

- Layer 1: T-2223 (the ADR refining ADR-0060 and ADR-0203) and T-2224 (`scripts/flai.sh` falls back to the installed flai without `go`); no paths in common.
- Layer 2: T-2225 (the two hook commands, the reproducing test in `checkGuardHooks`, the template changelog), after T-2224, whose fallback the reproducing case runs through.
- Layer 3: T-2226 (settings.md, flai.md, agent-adapters.md, the guard's help and the generated reference, a check by hand inside the image, `flai issue close I-0216`), after T-2223 and T-2225.

Figures: touches 16 files and the folder `design/adrs` (kept for the ADR's file, which no task can name); forecast 1h15m over flai's 4m, which was sized with no touches; cost of delay 37.50 USD a week as `flai cod` gives it from the input flai set from the issue. No `after` on the story.

Assumptions, each a decision you may overturn here:

1. Fail closed only on exit 126 and 127, not on every non-zero exit. A failed `go build` of this tree exits 1 and would otherwise lock a story's agent out of `Bash` and `Edit` together, since both hooks run the guard, and an older flai without `guard` would refuse everything; ADR-0060's "an older flai costs the guarantee, not the session" stands for those.
2. pilot (S-0507) and the container adapter (S-0508) pass no hook of their own, against the issue's first suggestion. Hooks from several Claude Code settings sources add up, so a `--settings` file of flai's would run beside the project's hook rather than replace it, and excluding the project's settings would drop a project's permissions with them. S-0507's note that the guard runs as the worktree's hook with the image's flai stands, and T-2174's "keep flai on PATH for the guard" is what this story relies on. Not verified on claude 2.1.312; the story's agent may find a flag that replaces the project's hooks cleanly, and should say so on this thread rather than build it.
3. `scripts/flai.sh` falls back only when `go` is absent, so a developer's host never runs an older installed flai in place of the tree's.
4. Two protected writes (`.claude/settings.json`, `template/root/.claude/settings.json`) ask you on a thread through `permission_prompt` unless `auto-approve` is on, and the story is accepted by you alone (ADR-0106).
5. The story's criterion 1 is met by the test in T-2225 plus one check by hand inside the image in T-2226, which needs `make agent-images` built from the branch.

Overlaps: `docs/users/flai.md`, `docs/operators/settings.md`, and `harness_test.go` are claimed by S-0414, S-0502, and S-0521 in progress, so the story is held until they move to review; S-0416 (ready) shares the settings files, `guard.go`, `guard_test.go`, and the docs, and the pull order settles which starts first.

### 2026-10-28T13:35:21Z orchestrator (recommendation)
**Recommended for the operator:** approve the plan with assumptions 1 to 3 as written. The hook should refuse a call when the guard command cannot start (exit 126 or 127) and keep failing open on every other exit. T-2223's ADR records this, refining ADR-0060 and ADR-0203. Then finalize S-0518.

## The plan fits

- **Tasks:** four tasks in three layers. T-2223, the ADR, and T-2224, the `scripts/flai.sh` fallback, share no path. T-2225, the two hook lines and the reproducing test, waits for T-2224, whose fallback the test runs through. T-2226 documents the change, checks it by hand in the image, and closes I-0216.
- **Figures:** forecast 1h15m, kept over flai's 4m with the reason stated. Value 37.50 USD a week, from `flai cod` on the issue's input.
- **Protected writes:** the two `.claude/settings.json` files make the story yours to accept (ADR-0106), and each write asks you through `permission_prompt`.

## Why it is yours

ADR-0060 § Decision says the guard "fails open, so a missing or older `flai` costs the guarantee, not the session". This plan makes a missing guard refuse the call, so it reverses part of an accepted decision on the guard, a safety control.

## Why I recommend it

- **A missing guard is exactly the case the hook must not miss.** I-0216 records that S-0513's container run passed every call unguarded, because `go build` exited 127.
- **The narrow cut keeps ADR-0060's reason.** An older flai, a failed build (exit 1), or unreadable input still fail open, so a session is never locked out of `Bash` and `Edit` for a cause the agent cannot fix. Only "the command is not there", 126 or 127, fails closed.
- **The fallback is host-safe.** `scripts/flai.sh` runs the installed flai only when `go` is absent, so a developer's host keeps the tree's flai (assumption 3).
- **No hook of pilot's or the adapter's own** (assumption 2): Claude Code adds hooks from several settings sources together, so a second hook would run beside the project's, not replace it.

**The alternative:** fail closed on every non-zero exit. I do not recommend it: a broken build would lock out every story's agent at once.

The thread stays open for your answer.

Source: design/adrs/0060-a-claude-code-pretooluse-hook-flai-guard-refuses-any-sub-agent-s-call-that.md § Decision

### 2026-10-30T13:36:14Z flai serve
Confirmed the recommendation of 2026-10-28T13:35:21Z orchestrator after 48h with no answer from the operator (orchestration.permissions.confirm_after).

Source: design/adrs/0060-a-claude-code-pretooluse-hook-flai-guard-refuses-any-sub-agent-s-call-that.md § Decision

### 2026-10-30T13:37:15Z orchestrator
The plan is approved, as the confirmed recommendation says. Assumptions 1 to 3 stand: the hook refuses a call only when the guard command exits 126 or 127, and keeps failing open on every other exit. T-2223's ADR records this. I finalized S-0518 (draft false). `flai promote --drafts` listed it as complete. Its two criteria, 16 touches, 1h15m forecast, and value of 37.50 USD a week fit the plan.

Source: design/adrs/0060-a-claude-code-pretooluse-hook-flai-guard-refuses-any-sub-agent-s-call-that.md

### 2026-10-30T13:37:28Z orchestrator
Resolved: The plan was approved by the recommendation flai serve confirmed after 48h, and S-0518 was finalized with draft false.
