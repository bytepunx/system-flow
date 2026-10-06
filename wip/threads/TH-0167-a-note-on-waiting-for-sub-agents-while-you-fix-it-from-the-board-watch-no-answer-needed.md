---
id: TH-0167
title: "A note on waiting for sub-agents while you fix it, from the board watch: no answer needed"
anchor:
  path: wip/kanban/stories/S-0285-a-story-s-agent-waits-for-its-sub-agents-with-wait-for-events-which-cannot-see-them-finish-so-each-wait-runs-to-its-timeout.md
  item: S-0285
status: resolved
participants: [claude-for-alex, alex]
created: 2026-10-06T10:32:15Z
updated: 2026-10-06T10:37:53Z
---

# TH-0167 A note on waiting for sub-agents while you fix it, from the board watch: no answer needed

On wip/kanban/stories/S-0285-a-story-s-agent-waits-for-its-sub-agents-with-wait-for-events-which-cannot-see-them-finish-so-each-wait-runs-to-its-timeout.md.

## Entries

### 2026-10-06T10:32:15Z claude-for-alex
This is Claude, watching the board for alex, and the author of this story's goal. It is a note, not a question: do not wait for an answer, and resolve this thread when you have read it.

You run on the installed flai, whose prompt still has the old sentence, so the two failures this story fixes can happen to you:

- Do not end a turn while a sub-agent is still running unless it will finish within a few minutes. Claude Code ends a headless session ten minutes after its turn ends (I-0084).
- Do not hold `wait_for_events` to wait for a sub-agent (I-0083).

Launch every sub-agent you must wait for with `run_in_background: false`, close-out verifiers included.

Evidence since the story was written, for its first criterion: S-0220's last session and S-0283's agent both launched their verifiers with `run_in_background: false` on Claude Code 2.1.290. In S-0283 a verifier ran ten minutes (10:11:00Z to 10:21:09Z) and the agent acted on its result about twenty seconds later, with the session intact. A layer of several launched in one message is not yet tried.

### 2026-10-06T10:37:53Z alex
Resolved.
