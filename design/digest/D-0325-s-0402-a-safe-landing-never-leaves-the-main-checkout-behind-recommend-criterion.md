---
id: D-0325
title: "S-0402: a safe landing never leaves the main checkout behind; recommend criterion 4 report the stories a dirty file holds back instead: Recommendation: take option (1)"
kind: answer
sources: [TH-0521]
topics: [orchestration, cli, analysis, flai, go, flaiover, dashboard, sveltekit, code, all]
paths: [design/adrs, flai/internal/mergequeue/land.go, flai/internal/mergequeue/land_test.go, flai/internal/storygit/storygit.go, flai/internal/storygit/storygit_test.go, flai/internal/serve/mergequeue.go, flai/internal/serve/mergequeue_test.go, flai/cmd/queue.go, flai/cmd/queue_test.go, flaiover/src/lib/components/MergeQueue.svelte, flaiover/src/lib/components/MergeQueue.svelte.test.ts, design/system/workflow.md, design/system/flai-cli.md, design/system/flaiover-dashboard.md, docs/operators/index.md, docs/users/flai.md, design/issues/I-0157-the-host-s-installed-flai-is-again-behind-the-project-s-releases-after-i-0149-closed-so-the-mcp-task-done-refuses-every-close.md, design/issues/summary.md, design/issues/I-0168-a-story-asked-to-share-a-path-agrees-in-a-settled-reply-but-records-no-share-so-the-asking-story-s-task-stays-refused-on-the-overlap.md, flai/cmd/accept.go, flai/cmd/accept_batch_dirty.go, flai/cmd/accept_batch_dirty_test.go, flai/cmd/accept_markdown.go, flai/internal/mergequeue/eligible.go, flai/internal/mergequeue/progress.go, design/components/flaiover/src/lib/components/MergeQueue.examples.json, flai/cmd/board.go, flai/cmd/queue_held_test.go]
date: 2026-10-20
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0325 S-0402: a safe landing never leaves the main checkout behind; recommend criterion 4 report the stories a dirty file holds back instead: Recommendation: take option (1)

The operator, alex, answered TH-0521 at 2026-10-20T04:27:37Z.

What was asked:

```text
The queue lands each story with `flai accept --batch` in the main checkout. That landing is `git merge --ff-only`. Git's own check refuses only a dirty file the change would overwrite, and it updates every other file. So a landing either moves main and the checkout's files together, or it does not move main at all. Criterion 4 ("after a landing that left the checkout's files behind main") then never happens.
```

The answer:

```text
**Recommendation:** take option (1).
- **The ADR:** S-0402's ADR decides that a landing stays in the main checkout with git's own check, so it never leaves the checkout's files behind main.
- **Criterion 4 reads instead:** "`flai queue`, `flai board`, and the dashboard name each story a dirty file in the main checkout holds back, the files, and the command that frees them (commit them, or `git stash push -- <files>`), until they are clean. The held story stays enqueued and lands with the next batch."

**Why this needs you.** It rewrites one of your criteria, which is a change of scope. Confirm it, or answer (2).

**Why (1):**
- **It is what ADR-0176 decided.** Each story lands as "`flai accept <story> --by <whoever enqueued it> --batch <batch ID>`, a process of this flai's own binary in the main checkout. That is the acceptance of today." The agent says that is a `--ff-only` merge, so a landing either moves main and the checkout together or does not move main at all.
- **Option (2) works against it.** Moving the ref from the batch's worktree would bypass that acceptance. It would also leave the reversed landing staged under your checkout, so one `git commit -a` would quietly undo a landed story.
- **It still delivers the story's goal.** A dirty file no queued story claims never stops a landing it does not overlap, and you are told exactly which files hold which story back.

**Note:** criterion 4's "what their checkout is behind" then has nothing to report, because the checkout is never behind.

If you agree, confirm this, and the story's agent can reword criterion 4 to match. If you want the checkout left behind on purpose, answer (2).

Source: design/adrs/0176-the-merge-queue-lands-each-story-of-a-passing-batch-through-flai-accept-batch.md § Decision
```

Source: TH-0521.
