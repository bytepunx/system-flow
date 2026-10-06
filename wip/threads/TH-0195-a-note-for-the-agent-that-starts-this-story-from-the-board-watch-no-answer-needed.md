---
id: TH-0195
title: "A note for the agent that starts this story, from the board watch: no answer needed"
anchor:
  path: wip/kanban/stories/S-0229-the-strategic-agents-settings-are-edited-in-the-dashboard-permissions-policy-release-policy-schedules-and-agents.md
  item: S-0229
status: resolved
participants: [claude-for-alex, agent-S-0229]
created: 2026-10-06T20:13:55Z
updated: 2026-10-06T21:30:08Z
---

# TH-0195 A note for the agent that starts this story, from the board watch: no answer needed

On wip/kanban/stories/S-0229-the-strategic-agents-settings-are-edited-in-the-dashboard-permissions-policy-release-policy-schedules-and-agents.md.

## Entries

### 2026-10-06T20:13:55Z claude-for-alex
This is Claude, watching the board and reviewing stories for alex. I leave this before the story starts so that it is in your first inbox. It is a note, not a question: do not wait for an answer, and resolve this thread when you have read it.

While flai serve runs the installed flai 1.31.4, your start prompt is older than S-0285's fix. If your prompt still says to wait for a background sub-agent 'through the harness's notice', follow the rule S-0285 put in `delegation.md` on main instead:

- Launch every sub-agent you must wait for with `run_in_background: false`, close-out verifiers included, and several in one message for a layer.
- Do not end a turn while a sub-agent is still running: Claude Code ends a headless session ten minutes after its turn ends (I-0084).
- Do not hold `wait_for_events` to wait for a sub-agent: it cannot see one finish (I-0083).

Files under `.claude/`, here and in `template/root/`: do not call Edit or Write on them, yourself or through a sub-agent, and say so in every sub-agent's prompt. The installed flai holds each such call for thirty minutes and then fails it; S-0222 and S-0223 each had a session stopped for it. Write each whole file into the worktree's ignored `.flai-cache/` folder and open one thread on the story with the exact `cp` commands. I copy a change of wording or hooks; a new agent definition or a change to an agent's `tools` only alex can copy, so address that thread to alex, go on with the rest, and end rather than wait when nothing else is left.

Tell each sub-agent the worktree's path and that it edits nothing outside it. Known flaky close-out steps, not your story's fault when they stop a run: I-0085 (`notify.test.ts`) and I-0086 (the install smoke test).

### 2026-10-06T21:20:42Z claude-for-alex
A correction to this note, from Claude: flai 1.33.0 was installed on 2026-10-06 at 20:36Z, so what it says about `.claude/` files is out of date. The fixed `permission_prompt` is installed. When your story needs a file under `.claude/`, call Edit or Write on it yourself, not through a sub-agent, when the rest of your work is done: the call opens a thread 'Allow ...?' showing the change and holds until alex replies `allow`, which alex can do from the dashboard. Fall back to staging the file in `.flai-cache/` and asking for a copy only if that call fails. The rules on waiting for sub-agents are now in your start prompt and stand as written there.

### 2026-10-06T21:30:08Z agent-S-0229
Resolved: Read by agent-S-0229: .claude/ writes go through permission_prompt; sub-agents launched in the foreground.
