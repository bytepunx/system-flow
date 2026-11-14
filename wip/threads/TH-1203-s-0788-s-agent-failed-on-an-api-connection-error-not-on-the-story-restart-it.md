---
id: TH-1203
title: "S-0788's agent failed on an API connection error, not on the story: restart it"
anchor:
  path: wip/kanban/stories/S-0788-stories-in-review-that-cannot-be-merged-need-an-action.md
  item: S-0788
status: resolved
participants: [orchestrator, agent-S-0788]
created: 2026-11-14T19:55:27Z
updated: 2026-11-14T23:13:51Z
---

# TH-1203 S-0788's agent failed on an API connection error, not on the story: restart it

On wip/kanban/stories/S-0788-stories-in-review-that-cannot-be-merged-need-an-action.md.

## Entries

### 2026-11-14T19:55:27Z orchestrator
**Recommended:** restart S-0788's agent with `flai serve agent restart S-0788`, or Retry on the story's page. Nothing in the story needs to change.

**Why (TH-1202):** all three runs ended with the same error from the harness. Their logs are in `~/.flai/serve/agents/sf-S-0788-*.log`:

```text
API Error: Unable to connect to API: SSL certificate hostname mismatch
```

- The two restarts each lasted about 2 seconds. Both failed on their first API call.
- The failure came from the host's network at 18:58Z, not from the story's work.
- That network fault is gone: my own session reaches the API now.

**Where the story stands:** its close-out passed at a2a1652 at 18:57:41Z (job S-0788-20261114T184712Z). T-3364 and T-3365 are done. Its agent most likely only needs to move it to review.

Restarting an agent is not among my permissions, so it is yours.

### 2026-11-14T23:13:51Z agent-S-0788
Resolved: The agent was restarted as recommended, and the API connection error did not happen again. Nothing in the story needed to change; S-0788 is moving to review.
