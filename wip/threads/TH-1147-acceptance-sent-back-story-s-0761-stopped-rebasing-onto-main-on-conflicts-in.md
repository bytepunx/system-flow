---
id: TH-1147
title: "Acceptance sent back: story/S-0761 stopped rebasing onto main on conflicts in design/system/flai-cli.md"
anchor:
  path: wip/kanban/stories/S-0761-flai-check-finds-item-blocked-closed-outside-the-story-at-close-out.md
  item: S-0761
status: resolved
participants: [flai, agent-S-0761, orchestrator, flai serve]
created: 2026-11-13T09:19:42Z
updated: 2026-11-13T09:32:24Z
---

# TH-1147 Acceptance sent back: story/S-0761 stopped rebasing onto main on conflicts in design/system/flai-cli.md

On wip/kanban/stories/S-0761-flai-check-finds-item-blocked-closed-outside-the-story-at-close-out.md.

## Entries

### 2026-11-13T09:19:42Z flai
The acceptance of S-0761 by orchestrator rebased `story/S-0761` onto `main`, and the rebase stopped on conflicts that the merge of the issue files, the ADR index, the issue summary, and the shared markdown files did not settle, in:

- `design/system/flai-cli.md`

flai aborted the rebase, so `story/S-0761` and `main` are as they were, and moved S-0761 back to in-progress. This thread is for S-0761's agent, which settles the conflicts; the operator does not. In its worktree, run `flai stream sync S-0761`, resolve each path it stops on, `git add` it, and run `git rebase --continue`; then close out and move S-0761 to review again. flai serve resolves this thread once S-0761 is back in review.

### 2026-11-13T09:22:01Z agent-S-0761
Settled in S-0761's worktree. The rebase stopped on line 54 of `design/system/flai-cli.md`, the `flai check` row of § Commands, which S-0735 and S-0761 both changed. I kept main's line as S-0735 left it and added S-0761's two edits: the `item.blocked-closed` sentence, and the rule with ADR-0296 in the `--record-issues` list. `git diff --word-diff main` now shows only those additions. The sync then finished on its own.

One thing for the acceptance. The sync warns that the branch "removes 1 line main added" to `design/system/flai-cli.md`. That is a false positive. The check compares whole lines, and that row is one long line, so S-0761's change to it looks like removing main's version. Every word of main's line is still there. The acceptance will be held on it (ADR-0155) until you accept with `--keep-removal design/system/flai-cli.md --reason "the row keeps every word of main's; S-0761 only adds to it"`. Do not run `flai stream restore`: it would drop S-0761's own edit. I recorded this as I-0360.

Source: design/adrs/0155-a-rebased-story-branch-that-undoes-main-s-change-to-a-path-is-recorded-at-sync.md

### 2026-11-13T09:23:16Z orchestrator (recommendation)
**Recommended:** once S-0761 is back in review and its close-out passes, the operator accepts it by hand, keeping the removal:

```bash
flai accept S-0761 --keep-removal design/system/flai-cli.md --reason "the flai check row keeps every word of main's; S-0761 only adds to it"
```

**Why the operator.** ADR-0155 § Decision makes the keep the operator's decision. flai refuses `--keep-removal` with `--by orchestrator`, so my acceptance would stop on the `undoes_main` blocker. ADR-0225 lands a `--keep-removal` acceptance by hand, not through the queue.

**What I checked.** `.flai-cache/undone/S-0761.json` records one path, `design/system/flai-cli.md`, and one line, the `flai check` row of § Commands. I ran `git diff --word-diff=porcelain main story/S-0761` on that path at 09d428b. It removes only two words, and puts each back:

| Removed | Added in its place |
|---------|--------------------|
| `and` | `and an item.blocked-closed outside the story,` |
| `ADR-0256),` | `ADR-0256, ADR-0296),` |

It also adds one new sentence on `item.blocked-closed`. None of main's words are lost, so agent-S-0761 is right that this is a false positive (I-0360). `flai stream restore` would drop S-0761's own edit.

**Meanwhile.** I leave S-0761 alone until it is in review. Then I verify it and run the acceptance preview. If the only blocker is this path, I leave the acceptance to you.

Source: design/adrs/0155-a-rebased-story-branch-that-undoes-main-s-change-to-a-path-is-recorded-at-sync.md § Decision

### 2026-11-13T09:32:24Z flai serve
Resolved: S-0761 is back in review since 2026-11-13T09:32:24Z, after the send back; its next acceptance takes it up again
