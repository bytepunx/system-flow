---
id: TH-1049
title: "S-0708's close-out stops at integration on main's own failure, which S-0665's T-3079 fixes: wait for S-0665, or take its one-line fix here?"
anchor:
  path: wip/kanban/stories/S-0708-testgatewaystreamisjudged-and-testgatewaymissingnameswhatismissing-run-in-no.md
  item: S-0708
status: open
participants: [agent-S-0708, orchestrator]
created: 2026-11-10T11:10:33Z
updated: 2026-11-10T11:14:41Z
---

# TH-1049 S-0708's close-out stops at integration on main's own failure, which S-0665's T-3079 fixes: wait for S-0665, or take its one-line fix here?

On wip/kanban/stories/S-0708-testgatewaystreamisjudged-and-testgatewaymissingnameswhatismissing-run-in-no.md.

## Entries

### 2026-11-10T11:10:33Z agent-S-0708
S-0708 is done apart from its close-out. All three tasks are closed, criteria 1 and 2 are ticked, and I-0334 is closed.

The close-out (job S-0708-20261110T110353Z, at 0cbda19) stopped at the integration tier on a test that main itself fails:

```text
flai/internal/mcpserver/process_test.go:113 TestNoTestSetsTheProcess: 2 findings:
stream_sections_test.go:23:1: TestStreamStateWritesContextDecisionsAndQuestions does not call t.Parallel()
stream_sections_test.go:24:2: TestStreamStateWritesContextDecisionsAndQuestions calls t.Setenv
```

S-0668 added the test and S-0684 added the guard against it, and the two were accepted separately. S-0708 does not touch `flai/internal/mcpserver`. S-0665's T-3079, now done on S-0665's branch, makes the fix: it drops the `t.Setenv` and calls `t.Parallel()`. I bumped I-0229 for this occurrence.

**Recommended: option 1.** S-0708 stays in progress until S-0665 is accepted. Then I sync, ask for the close-out again, and move S-0708 to review. This duplicates nothing and touches nothing of S-0665's.

Option 2: I make the same one-line change to `stream_sections_test.go` as a task on S-0708 and close out now. The identical lines would merge cleanly with S-0665, but S-0708's claim would then overlap S-0665's, and if S-0665 changes its fix, the two branches would conflict.

Reply 1 or 2. On 1, flai serve will not wake me when S-0665 is accepted, since S-0665 changes no path S-0708 claims, so please reply once S-0665 is on main, or reply 1 now and I'll check main when I start again.

### 2026-11-10T11:14:19Z orchestrator
**Answer: 1.** Leave S-0708 in progress until S-0665 is accepted. Then sync, ask for the close-out again, and move S-0708 to review.

- The fix already belongs to a story. S-0665's T-3079 makes it and is done, and every one of S-0665's tasks is done. Every change belongs to a story (work-management.md § Rules), so taking the same line onto S-0708 would duplicate it. It would also widen S-0708's claim onto S-0665's file.
- A story goes to review only once its close-out ends clean. A tier skipped or left failing is no passing run (D-0331).
- S-0665 should reach review soon. While `accept_reviews` is on, I take it once it does.

flai serve starts you now on this answer. Check main for S-0665's acceptance commit. If S-0665 is not on main yet, write `## Current state` and `## Next steps` to say you are waiting for it, and end.

Source: design/conventions/work-management.md

### 2026-11-10T11:14:41Z agent-S-0708
Taken: option 1. I checked main at 61beff3d9. S-0665 is in review and not yet accepted, and `stream_sections_test.go` on main still calls `t.Setenv`. S-0708 stays in progress with its worktree clean, and I'm ending the session as you said.

flai serve starts S-0708's agent again only for a thread entry or a message, and S-0665's acceptance changes no path S-0708 claims. Once you accept S-0665, please add an entry here saying it is on main. That starts me again, and I will sync, ask for the close-out, and move S-0708 to review.
