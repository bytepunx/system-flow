---
id: TH-0403
title: S-0362 is ready to accept, but main has an untracked scripts/release-key-setup.sh
anchor:
  path: wip/kanban/stories/S-0362-flai-task-done-commits-the-changed-paths-of-another-open-task-of-the-same-layer.md
  item: S-0362
status: resolved
participants: [orchestrator, alex]
created: 2026-10-12T14:53:49Z
updated: 2026-10-12T15:10:25Z
---

# TH-0403 S-0362 is ready to accept, but main has an untracked scripts/release-key-setup.sh

On wip/kanban/stories/S-0362-flai-task-done-commits-the-changed-paths-of-another-open-task-of-the-same-layer.md.

## Entries

### 2026-10-12T14:53:49Z orchestrator
**Recommendation:** move `scripts/release-key-setup.sh` out of the main checkout, or commit it on the story it belongs to (S-0232's worktree looks right). Or stash it with `git stash -u`. Then reply here, and I accept S-0362 at once.

## Where S-0362 stands

- `flai verify` passed every step at the branch head, f6fa8938.
- The verifier matched both criteria to the diff:
  1. The tests in `flai/internal/taskdone/taskdone_test.go` and `flai/internal/mcpserver/task_test.go` reproduce I-0126's layer.
  2. I-0126 is closed, and `design/issues/summary.md` is updated.
- The one file outside its touches is I-0134, which the close-out recorded.

## What stops it

`flai accept S-0362 --by orchestrator --verified f6fa8938 --dry-run` says:

```text
uncommitted outside wip: scripts/release-key-setup.sh; the real run refuses until they are committed or stashed, or --yes includes them
```

The file is untracked in the main checkout, mode 755, written today at 08:46 local time. It is not mine, and I will not use `--yes`: that would commit your file into S-0362's acceptance.

S-0362 stays in review until then.

### 2026-10-12T15:10:25Z alex
Resolved: S-0362 was accepted
