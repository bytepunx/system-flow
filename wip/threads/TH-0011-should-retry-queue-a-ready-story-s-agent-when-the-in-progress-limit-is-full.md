---
id: TH-0011
title: Should Retry queue a ready story's agent when the in-progress limit is full?
anchor:
  path: wip/kanban/stories/S-0118-when-viewing-a-story-s-page-with-an-agent-pane-if-the-agent-has-failed-a-retry-button-should-appear.md
  item: S-0118
status: resolved
participants: [agent-S-0118, alex]
created: 2026-09-26T03:02:18Z
updated: 2026-09-26T03:14:52Z
---

# TH-0011 Should Retry queue a ready story's agent when the in-progress limit is full?

On wip/kanban/stories/S-0118-when-viewing-a-story-s-page-with-an-agent-pane-if-the-agent-has-failed-a-retry-button-should-appear.md.

## Entries

### 2026-09-26T03:02:18Z agent-S-0118
The story page already has a Restart agent button for a failed agent (S-0116). S-0118 changes it to Retry: it moves to the top right of the agent pane, and it hides once clicked. One criterion says Retry should "queue the agent as it would normally do". Today, restarting a failed agent on a story in ready is refused when the in-progress limit is full ("the in-progress limit leaves no room for S-…"). Nothing is queued.

Recommended: queue it. For a story in ready with the limit full, Retry marks the failed run as queued in serve/agents.json and does not refuse. flai serve then starts the agent itself as soon as the limit has room, the same way it starts any ready story. Until then the pane says the agent is queued and shows no Retry button. When there is room, and for a story in progress, Retry starts the agent at once, as Restart does today. `flai serve agent restart` works the same way from a shell.

Alternative: keep refusing when the limit is full, and only move and rename the button.

I am doing the pane changes now. The flai side waits for your answer.

### 2026-09-26T03:05:33Z alex
No restart button appears on the failed story I have on the board now (for the record). Queue the story as recommended but change the indicator to yellow so it's clear that the story agent is waiting for something.

### 2026-09-26T03:06:01Z agent-S-0118
Understood. I will queue it, and a queued agent shows a yellow dot with why "queued: … starts it when the in-progress limit has room". I am doing that now in T-0430.

About the missing button: this host runs versions older than S-0116. The dashboard container has no Restart agent button in its build, and the installed flai is 1.16.1. S-0116's restart shipped in flai/v1.17.0 and flaiover/v0.27.5. `flai host upgrade`, or `flai self-upgrade` followed by a new dashboard image, brings the button. S-0118's Retry needs the same, once it is released.

Also, for the record: S-0117's agent failed for a different reason. Its log (`~/.flai/serve/agents/sf-S-0117-20260926T024509Z.log`) says "Claude Code 2.1.250 does not support this model; version 2.1.280 or newer is required". Retry cannot fix that. Running `claude update` on the host does.

### 2026-09-26T03:09:28Z alex
I have run that command but because I had the older version I was not seeing a way to retry the story. Let's complete the changes requested as part of this story.

### 2026-09-26T03:13:37Z agent-S-0118
Understood. The pane's Retry is done (T-0429). Queueing is done too (T-0430): with the in-progress limit full, Retry on a ready story queues its agent and the dot turns yellow. Now I am updating the docs and the design (T-0431). Then S-0118 goes to review. Once it is accepted and published, `flai host upgrade` and a new dashboard image bring Retry to this host.

### 2026-09-26T03:14:52Z agent-S-0118
Resolved: Queued, with a yellow dot, as decided; ADR-0044, built in T-0430.
