---
id: TH-1142
title: "The publish step exited 1: another git process held the main checkout's index lock"
anchor:
  path: wip/kanban/stories/S-0764-publish-what-is-accepted-and-not-yet-released.md
  item: S-0764
status: answered
participants: [agent-S-0764, orchestrator]
created: 2026-11-13T08:58:00Z
updated: 2026-11-13T09:00:03Z
---

# TH-1142 The publish step exited 1: another git process held the main checkout's index lock

On wip/kanban/stories/S-0764-publish-what-is-accepted-and-not-yet-released.md.

## Entries

### 2026-11-13T08:58:00Z agent-S-0764
Step `publish` (T-3275), `flai release --pending --require-pending --json` in the main checkout, exited 1. The formula's stop does not list that exit.

```text
git add -A: exit status 128
fatal: Unable to create '/home/alex/git/bytepunx/system-flow/.git/index.lock': File exists.
```

What I found:

- Seconds later the lock no longer existed. Another git process held it for a moment: the orchestrator and the curator were both running then.
- HEAD is still `631c7c498`, no tag points at it, and nothing was committed or pushed.
- The fetch and evaluate steps passed. The policy is met by S-0660 (1250 USD/week, against 1000).

Recommendation: answer "run it again". I will then `flai unblock T-3275`, run the step again, and close the story out to review. If you would rather publish by hand, say so: I will then stop the run with no release.

### 2026-11-13T09:00:03Z orchestrator
Run it again.

- The failure was transient: `.git/index.lock` no longer exists, HEAD is still `631c7c498`, and no tag points at it, so nothing was half-done.
- The policy is still met: `release_evaluate` gives threshold, S-0660 at 1250 USD/week against 1000, count 1 against 1.
- ADR-0239 §3 covers this case: unblock T-3275 with `flai unblock T-3275`, run the `publish` step again as written, then go on and close the story out to review.
- If it fails again on the lock, block the task again and say so on this thread. Do not retry in a loop.

Source: design/adrs/0239-a-formula-s-instance-is-worked-by-its-story-s-agent-step-by-step-a-role-step.md § 3. A failed step
