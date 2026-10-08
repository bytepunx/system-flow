---
id: TH-0385
title: S-0321 in review needs a sync and a passing verify at its head
anchor:
  path: wip/archive/kanban/stories/S-0321-publishing-raises-the-manifest-s-flai-minimum-the-moment-it-commits-before-the-release-is-built-so-the-host-s-flai-drops-the-project-until-the-binaries-exist.md
  item: S-0321
status: resolved
participants: [orchestrator]
created: 2026-10-08T09:00:18Z
updated: 2026-10-08T09:00:24Z
---

# TH-0385 S-0321 in review needs a sync and a passing verify at its head

On wip/archive/kanban/stories/S-0321-publishing-raises-the-manifest-s-flai-minimum-the-moment-it-commits-before-the-release-is-built-so-the-host-s-flai-drops-the-project-until-the-binaries-exist.md.

## Entries

### 2026-10-08T09:00:18Z orchestrator
Recommendation for agent-S-0321 (or the operator): run `flai stream sync S-0321`, then `flai verify S-0321`, and I accept it when verify passes at the branch head.

What is missing under ADR-0093:

- The stored verify passed at 3a1c0fc2 (08:44Z), but the branch head is now 12771d32.
- The branch is 10 commits behind main. Those commits include S-0340's acceptance (de293abf), whose code (`flai/internal/releaseserver`, `scripts/install-test.sh`, `flai/.goreleaser.yaml`) changes what smoke runs. So the sync needs a fresh run, not only the sync step.

I have not run the verifier against the criteria yet, because there is no passing run at the head to review.

### 2026-10-08T09:00:24Z orchestrator
Resolved: Moot: the operator accepted S-0321 (f4389078) just before this thread opened.
