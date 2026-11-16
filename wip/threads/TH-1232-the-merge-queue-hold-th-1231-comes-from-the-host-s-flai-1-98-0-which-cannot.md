---
id: TH-1232
title: "The merge queue hold (TH-1231) comes from the host's flai 1.98.0, which cannot read S-0748's role files: publish S-0748, then upgrade the host"
anchor:
  path: wip/archive/kanban/stories/S-0748-a-role-file-s-to-review-and-to-backlog-conditions-become-one-condition-that.md
  item: S-0748
status: open
participants: [orchestrator]
created: 2026-11-15T21:47:54Z
updated: 2026-11-15T21:47:54Z
---

# TH-1232 The merge queue hold (TH-1231) comes from the host's flai 1.98.0, which cannot read S-0748's role files: publish S-0748, then upgrade the host

On wip/archive/kanban/stories/S-0748-a-role-file-s-to-review-and-to-backlog-conditions-become-one-condition-that.md.

## Entries

### 2026-11-15T21:47:54Z orchestrator
**Recommended:** leave the role files as they are. Let S-0794 publish S-0748, then run `flai host upgrade`. The queue's next look then finds main clean and clears TH-1231 by itself.

**Why the queue is held:**

- S-0748 changed `design/workflows/roles/orchestrator.md` and `planner.md` to its new trigger form, `{event: moved, when: {to: review}}`, and the new `{call: …, when: own-thread}` tool entries.
- The flai on this host is 1.98.0 (`inbox` reports `flai_outdated`: 1.99.0 is in the project's history). It predates S-0748, so its `flai check` cannot parse those files. The queue's preflight runs that check. My own `prime` fails on the same two files.
- The files are correct for the flai S-0748 built. Reverting them would undo S-0748.

**What I did:** `release_evaluate` is met under `threshold` (1 accepted story, threshold 1). I started the publish formula: instance S-0794, now ready.

**What may need you:**

1. S-0794 cannot be pulled yet: in-progress is full (3 of 3: S-0750 blocked, S-0755, S-0790). If you want the queue moving sooner, run `flai release --pending` yourself, or install the tree's build, then `flai host upgrade`.
2. If S-0794's publish step fails because 1.98.0 cannot read the role files, its agent blocks the task and asks on its own thread. The way out is the same: install the flai built from main, then upgrade the host.
