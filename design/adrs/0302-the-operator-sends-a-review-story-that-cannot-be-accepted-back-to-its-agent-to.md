---
id: ADR-0302
title: "The operator sends a review story that cannot be accepted back to its agent to rebase and resolve what holds it, with flai accept --resolve and the host method accept.resolve"
status: accepted
date: 2026-11-14
supersedes: []
superseded_by: []
refines: [ADR-0248, ADR-0291]
topics: [flai, cli, dashboard]
---

# ADR-0302 The operator sends a review story that cannot be accepted back to its agent to rebase and resolve what holds it, with flai accept --resolve and the host method accept.resolve

## Context

A story in review often cannot be accepted because main moved on after its close-out: its last flai verify ran over an older main (verify_behind_main, ADR-0291), the tiers marked queue did not run at its head, its worktree holds uncommitted changes, and so on. `flai accept --dry-run` and the dashboard's review page list each blocker with what to do, but the operator could only read them: settling them is the story's agent's work, in its worktree, and nothing sent the story back to it. ADR-0248 sends a story back to its agent when an acceptance's own rebase stops on conflicts it cannot settle, through a thread flai serve starts the agent on; a blocker found before anything is merged sent nothing back.

## Decision

`flai accept S-nnnn --resolve`, and the host method `accept.resolve` with `{id}`, which runs `flai accept <id> --resolve --by=<owner> --json`, are the operator's request that a review story's agent rebase its branch onto the main branch and resolve what holds its acceptance. It accepts nothing and merges nothing. It runs the acceptance preview as `flai accept --dry-run` does, and when the story is in review and the preview lists at least one blocker it sends the story back through ADR-0248's send-back:

1. it moves the story from review to in-progress as `flai`, with a reason naming who asked, the branch, the main branch, and the first blocker;
2. it opens the send-back thread on the story, by `flai` with a title beginning `Acceptance sent back:`, or writes on the one still open, listing each blocker as the preview words it and telling the agent to sync the branch with `flai stream sync`, resolve each path the rebase stops on, take each blocker's recommended action, close out again, and move the story to review;
3. flai serve starts the story's agent on that thread and resolves it once the story is back in review, as it does for every acceptance's send-back, with no change of its own.

With `--json` it prints the send-back object, `{id, reason, sent_back, branch, base, conflicts, thread, blockers}`. It refuses, exiting 4 and printing `{"refused": {id, reason}}` with `--json`, a story that is not in review and a story whose preview lists no blocker, which the operator accepts instead. `--resolve` is the operator's: it is refused with `--dry-run`, `--yes`, `--batch`, the orchestrator's acceptance, under `FLAI_ROLE=orchestrate`, and with the flags of an acceptance (`--verified`, `--evidence`, `--keep-removal`, `--reason`, `--without-queue-tiers`, `--without-kluster`, `--trailer`). `accept.resolve` is gated by no host action: it changes work items and a thread in the main checkout and reaches no forge. The dashboard's Rebase and resolve button on a review story that cannot be accepted runs it.

## Consequences

A blocker the preview finds is something the operator can hand back with one action rather than relay by hand. The agent's prompt is the one flai serve already gives an acceptance's send-back, so the thread's title reads after "the acceptance of S-nnnn found that". A blocker the agent cannot settle, such as a git identity missing where the dashboard runs, sends the story back all the same; the agent says so on the thread. The send-back thread now carries two kinds of request, a rebase that stopped and blockers before a merge, told apart by their words only.

## Alternatives considered

- The dashboard sending the agent a message of its own: a second path to the agent beside ADR-0248's, which flai serve would need to learn to start and to resolve.
- Rebasing the branch onto main from the acceptance itself: the rebase alone does not settle a verify behind main, which needs a close-out over the new main, nor most other blockers.
- Gating it on a host action: it writes nothing a move and a thread from the dashboard do not already write.
