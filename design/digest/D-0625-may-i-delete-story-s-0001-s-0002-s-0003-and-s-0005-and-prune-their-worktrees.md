---
id: D-0625
title: "May I delete story/S-0001, S-0002, S-0003, and S-0005 and prune their worktrees: Recommended: yes. agent-S-0781 may run, in the main checkout, `git worktree prune` and then `git branch -D…"
kind: answer
sources: [TH-1190]
topics: [cli, conventions, template, testing, flai, go, code, all]
paths: [flai/cmd/cmd_test.go, flai/internal/mcpserver/server_test.go, flai/cmd/story_start_test.go, design/adrs, design/adrs/README.md, flai/internal/check/scope.go, flai/internal/check/scope_test.go, flai/cmd/check.go, flai/cmd/check_test.go, docs/users/flai-reference.md, design/conventions/work-management.md, template/root/design/conventions/work-management.md, template/CHANGELOG.md, design/system/flai-cli.md, design/system/continuous-improvement.md, docs/users/flai.md, design/issues/I-0371-flai-check-finds-story-unaccepted-outside-the-story-at-close-out.md, design/issues/summary.md, flai/cmd/verify_test.go, flai/internal/verify/story_test.go]
date: 2026-11-14
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0625 May I delete story/S-0001, S-0002, S-0003, and S-0005 and prune their worktrees: Recommended: yes. agent-S-0781 may run, in the main checkout, `git worktree prune` and then `git branch -D…

The operator, alex, answered TH-1190 at 2026-11-14T17:43:45Z.

What was asked:

```text
**Question.** May I delete the four branches S-0777's test run leaked into the main repository and prune their worktrees? (T-3345, S-0781)
```

The answer:

```text
**Recommended: yes.** agent-S-0781 may run, in the main checkout, `git worktree prune` and then `git branch -D story/S-0001 story/S-0002 story/S-0003 story/S-0005`. It should delete nothing else.

**What I checked in the main checkout:**

- All four branches point at 5c3b47d0c, S-0777's pre-rebase commit "flai/cmd's kluster fixtures resolve the project root". Its parent is 9ac5c38fe.
- main has both commits under the same subjects: 7eee9cc4b and 2605d60d3. So deleting the branches loses no work.
- `git worktree list` marks all four worktrees `prunable`. They sit under S-0777's test temp folders in `.flai-cache/worktrees/S-0777/bin/…`.

**Why this needs your word:** `branch -D` cannot be undone, and these are branches the agent did not create. `safety.md` asks for confirmation before deleting or force-removing anything like that. Your confirmation is that confirmation, for this one action only.

**Alternative:** keep them. Then the main checkout's `flai check` keeps warning `story.unaccepted` (I-0371).

Source: design/conventions/safety.md § Rules
```

Source: TH-1190.
