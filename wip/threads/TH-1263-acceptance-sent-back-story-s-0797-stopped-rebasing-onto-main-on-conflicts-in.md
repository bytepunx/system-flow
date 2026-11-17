---
id: TH-1263
title: "Acceptance sent back: story/S-0797 stopped rebasing onto main on conflicts in design/system/flai-cli.md, docs/users/flai.md"
anchor:
  path: wip/kanban/stories/S-0797-the-merge-queue-s-preflight-holds-the-queue-on-main-s-files-when-the-installed.md
  item: S-0797
status: resolved
participants: [flai, flai serve]
created: 2026-11-17T02:02:10Z
updated: 2026-11-17T02:25:29Z
---

# TH-1263 Acceptance sent back: story/S-0797 stopped rebasing onto main on conflicts in design/system/flai-cli.md, docs/users/flai.md

On wip/kanban/stories/S-0797-the-merge-queue-s-preflight-holds-the-queue-on-main-s-files-when-the-installed.md.

## Entries

### 2026-11-17T02:02:10Z flai
flai serve's rebase and resolve job S-0797-20261117T020209Z, which alex asked for, rebased `story/S-0797` onto `main`, and the rebase stopped on conflicts that the merge of the issue files, the ADR index, the issue summary, and the shared markdown files did not settle, in:

- `design/system/flai-cli.md`
- `docs/users/flai.md`

flai aborted the rebase, so `story/S-0797` and `main` are as they were, and moved S-0797 back to in-progress. This thread is for S-0797's agent, which settles the conflicts; the operator does not. In its worktree, run `flai stream sync S-0797`, resolve each path it stops on, `git add` it, and run `flai stream sync S-0797` again, which continues the rebase; then close out again, with the MCP tool `verify` or `flai verify S-0797 --job`, and move S-0797 to review again. flai serve resolves this thread once S-0797 is back in review. The job's whole output is in `/home/alex/git/bytepunx/system-flow/.flai-cache/verify/jobs/S-0797.log`.

### 2026-11-17T02:25:29Z flai serve
Resolved: S-0797 is back in review since 2026-11-17T02:25:28Z, after the send back; its next acceptance takes it up again
