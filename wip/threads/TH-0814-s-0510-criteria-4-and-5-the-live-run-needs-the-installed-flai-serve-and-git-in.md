---
id: TH-0814
title: "S-0510 criteria 4 and 5: the live run needs the installed flai serve and git in the container, so move it to S-0565?"
anchor:
  path: wip/kanban/stories/S-0510-a-project-s-manifest-names-an-image-layer-for-its-test-tiers-tools-flai-serve-runs-it-flai-check-warns-when-it-is-not-built-and-this-repository-s-layer-runs-the-close-out-inside-the-container.md
  item: S-0510
status: open
participants: [agent-S-0510, orchestrator]
created: 2026-11-02T18:18:11Z
updated: 2026-11-02T18:19:19Z
---

# TH-0814 S-0510 criteria 4 and 5: the live run needs the installed flai serve and git in the container, so move it to S-0565?

On wip/kanban/stories/S-0510-a-project-s-manifest-names-an-image-layer-for-its-test-tiers-tools-flai-serve-runs-it-flai-check-warns-when-it-is-not-built-and-this-repository-s-layer-runs-the-close-out-inside-the-container.md.

## Entries

### 2026-11-02T18:18:11Z agent-S-0510
**Recommended:** answer "as recommended". Criterion 5 moves to S-0565, and criterion 4's close-out part is proven there too. S-0510 proves criterion 4 with a container run by hand of the layered image on its own worktree.

## Why

- Criterion 5 needs `flai serve` to run the layer. Only the installed release runs here, and it gains the layer only once S-0510 is accepted and published. S-0508's criterion 6 hit the same wall, and TH-0812 moved it to S-0565 (D-0461).
- `scripts/close-out.sh` runs git: sync, the commit, and the branch check. Inside the container the worktree's `.git` file names the main checkout's `.git/worktrees`, which is not mounted. That is I-0215, S-0517's, still in progress. Until it lands, no close-out passes in a container, whatever the layer holds.
- S-0565 already is the live container run. It would gain the layer, `flai test`, and the close-out passing inside the container.

## What S-0510 does meanwhile

- Builds `flai-agent-sf-claude-code:local` from `flai/images/system-flow/Dockerfile`: Go from `flai/go.mod`, golangci-lint v2.5.0, Node 24, pnpm, and ShellCheck.
- Runs `flai test` inside it on this story's worktree by hand, with the main checkout's `.git` mounted read-only only so git works, and records each tier that passes or fails in the story's notes.
- Records each limitation it meets as an issue.

## The alternatives

- **Keep criterion 5 here.** S-0510 would sit in review, blocked, until S-0517 lands and S-0510 is published.
- **Drop it.** S-0565 would not cover the layer.

S-0510 goes on meanwhile with T-2164, T-2165, T-2171, and T-2172. Only T-2176's live run waits on your answer.

### 2026-11-02T18:19:19Z orchestrator (recommendation)
**Recommended:** say yes, as you answered for S-0508 on TH-0812 (D-0466) and for S-0423 on TH-0803 (D-0461).

- S-0510's agent removes criterion 5 from S-0510, and narrows criterion 4 to what a hand run of the layered image on its own worktree shows: each tier `flai test` runs inside it, recorded in the story's notes.
- It adds criterion 5's words and criterion 4's close-out half to S-0565, which is still a draft, and adds `S-0510` and `S-0517` to S-0565's `after`. S-0565's Notes already say its proof needs the layer (S-0510) and git in the container (S-0517).
- S-0510's Notes record the move.

**Why:** the same case as D-0466. The layer exists only on `story/S-0510`, and this host's `flai serve` is the installed release. The close-out also needs I-0215's fix, S-0517, still in progress. So S-0510 could not meet the definition of done with criterion 5 kept.

**Why I recommend rather than answer:** this takes a criterion out of a finalized story and widens another. That changes scope, so it is your call. Finalizing S-0565 stays yours.

**Alternative:** keep criterion 5. I don't recommend it: S-0510 would wait in review until S-0517 lands and S-0510 is published.

Source: design/digest/D-0466-s-0508-the-live-container-run-criterion-6-t-2169-s-first-half-cannot-run-before.md
