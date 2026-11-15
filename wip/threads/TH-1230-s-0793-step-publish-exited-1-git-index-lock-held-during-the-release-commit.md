---
id: TH-1230
title: "S-0793 step publish exited 1: .git/index.lock held during the release commit"
anchor:
  path: wip/kanban/stories/S-0793-publish-what-is-accepted-and-not-yet-released.md
  item: S-0793
status: answered
participants: [agent-S-0793, orchestrator]
created: 2026-11-15T21:31:50Z
updated: 2026-11-15T21:32:46Z
---

# TH-1230 S-0793 step publish exited 1: .git/index.lock held during the release commit

On wip/kanban/stories/S-0793-publish-what-is-accepted-and-not-yet-released.md.

## Entries

### 2026-11-15T21:31:50Z agent-S-0793
The step `publish` (T-3386) failed. `scripts/flai.sh release --pending --require-pending --json` exited 1 in the main checkout. The formula's stop lists only exit 2.

```text
git commit -q -m chore: publish flai 1.99.0: exit status 128
fatal: Unable to create '/home/alex/git/bytepunx/system-flow/.git/index.lock': File exists.
```

What I found afterwards (I only read; I changed nothing):

- `.git/index.lock` is gone, so another git process held it for a moment. It was probably another agent or flai serve writing in the main checkout.
- There is no `flai/v1.99.0` tag; the newest is `flai/v1.98.0`. HEAD is still `2fd4f33d4`, and nothing was pushed.
- The main checkout's `wip/` and `design/digest/` changes are now staged in the index, not committed.
- The fetch (T-3384) and the evaluate (T-3385) both exited 0. The policy is met, with S-0741 pending.

Recommendation: answer "retry". I will then `flai unblock T-3386` and run the publish step again, unchanged, from the main checkout. The failure looks transient. If it happens again, I will block and ask again rather than remove a lock.

### 2026-11-15T21:32:46Z orchestrator
**Answer: retry, once.** This is D-0603, the answer given on TH-1142 for the same failure in S-0764.

1. Unblock T-3386 with `flai unblock T-3386`.
2. Run the `publish` step again, as written, from the main checkout.
3. Record the versions and tags it prints in S-0793's Notes. Then close the story out and move it to review.

The conditions D-0603 sets are met:

- The lock is gone.
- No `flai/v1.99.0` tag exists, and HEAD is still `2fd4f33d4`.
- Nothing was committed or pushed.
- Fetch and evaluate both exited 0, and the threshold policy is met with S-0741 pending.

The `wip/` and `design/digest/` files that are staged in main's index come from the failed run's own staging. The retry stages and commits them again.

If it fails again on the lock, block T-3386 again and say so on this thread. Do not retry in a loop, and do not remove a lock. I-0359 records the defect, and S-0793 already touches it, so add this occurrence there.

Source: design/digest/D-0603-a-publish-step-that-exits-1-because-another-git-process-briefly-held-the-main.md
