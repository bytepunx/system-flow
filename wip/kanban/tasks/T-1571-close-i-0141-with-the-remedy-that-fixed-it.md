---
id: T-1571
type: task
nature: remediation
title: Close I-0141 with the remedy that fixed it
status: backlog
parent: S-0387
owner: alex
created: 2026-10-14T20:55:02Z
updated: 2026-10-14T20:55:02Z
transitions: []
stream: S-0387
tags: [issues]
touches: [design/issues/I-0141-the-mcp-test-tool-of-the-installed-flai-1-40-1-refuses-main-s-manifest-so-a-story-s-agent-runs-flai-test-from-the-tree.md, design/issues/summary.md]
after: [T-1570]
---
# T-1571 Close I-0141 with the remedy that fixed it

## Work

Close I-0141 in the story's worktree with `flai issue close I-0141 --reason`, naming the ADR, `flai serve`'s periodic check with `auto-upgrade` and its thread (T-1568), and the version-aware tier refusal (T-1569). It waits for T-1570, the last of the work, so the reason names what was built and documented.

## Done when

- I-0141's status is closed with the reason, and `design/issues/summary.md` no longer lists it.
- `flai check --strict` passes.

## Notes
