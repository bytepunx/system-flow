---
id: ADR-0042
title: "Someone attending holds a ready story back for the attended window, then flai serve starts it"
status: accepted
date: 2026-09-24
supersedes: []
superseded_by: [ADR-0043]
refines: [ADR-0038, ADR-0041]
---

# ADR-0042 Someone attending holds a ready story back for the attended window, then flai serve starts it

## Context

[ADR-0038](0038-flai-serve-starts-a-story-s-own-agent-through-an-adapter-with-what-the-operator.md) has `flai serve` start nothing while someone is attending the project: an agent connected over MCP, or one writing a story's narrative, will see the ready story in its inbox and pull it. [ADR-0041](0041-a-story-in-ready-gets-its-agent-whether-it-entered-ready-before-or-after-flai.md) names attendance as one of the reasons a ready story waits. Nothing bounded that wait. An agent at work on its own story keeps its sign fresh for as long as it works and pulls nothing else, so a ready story waited as long as any agent was active while the in-progress limit had room.

On 2026-09-24 the operator watched ready stories sit between agent runs and filed S-0114. The serve log showed the launcher looking every minute, as designed, and refusing each time because someone was attending. The sign was the MCP cursor `system-flow.json`: every agent flai serve had started connected under that one name, because the project's `.claude/settings.local.json` sets `FLAI_AGENT` in its `env` and Claude Code applies that to the MCP server it spawns, over the name flai serve set. The launcher could not tell its own agents' signs from anyone else's (I-0037), so each of its own agents ending held the next start back six minutes. The reason logged while the limit was full also said "attending", because attendance was checked first.

## Decision

**Someone attending holds a ready story back for the attended window, and no longer. Once attendance has been the only thing keeping a story from its agent for `attended_minutes` (6 by default), `flai serve` starts it, whoever is attending.** The operator asked for it in S-0114 (2026-09-24). This refines ADR-0038 and ADR-0041; it reverses neither.

- The launcher remembers when attendance became the only thing holding each story back: the story could start (its agent action on, a harness or command to start it, no run since it entered ready) and the in-progress limit had room. From then the story has the window to be pulled. The `waiting` reason names the time it has until.
- The window is `attended_minutes`, the same setting that says how recent a sign must be to count as attending. One knob: an idle agent holding `wait_for_work` pulls within seconds, an attending agent that has not pulled in six minutes is at work on something else.
- A full in-progress limit is the reason a story waits, whoever is attending. Attendance holds a story back only when the limit would let it start, and the hold begins only then.
- The hold is the launcher's memory alone, not in `serve/agents.json`. A restart of `flai serve` gives the attending agent the window again, which costs at most one window.
- `agent not started` is logged when a story's reason changes in kind, not when only the sign's age or the time in it changes.
- The agent's name reaches the MCP server as an argument: the `claude-code` adapter runs `flai mcp --agent <name>`, which wins over `FLAI_AGENT`. An environment a harness sets for the servers it starts can no longer merge every agent into one name, so the launcher's own agents' signs are excluded as ADR-0038 intended, and a thread one agent opens reads as awaiting the other.

## Consequences

- A ready story is never held back by attendance for longer than the window. An operator working in an interactive session sees the story in their inbox for six minutes, then an agent is started for it as if nobody were there. Enabling the `agent` action means that.
- Two agents can meet on a story that was held by the limit: when the limit frees, an idle attending agent and the launcher both act, the launcher only once the window has passed. The second `item_move` to in-progress is refused, and the flai-started agent's prompt tells it to pull the next ready story or hold `wait_for_work` instead.
- Shell commands the agent runs still see the harness's own `FLAI_AGENT`. An operator who set one in `.claude/settings.local.json` should remove it: flai serve sets the name itself.

## Alternatives considered

- **Another periodic look.** The launcher already looks every minute and on every change under `kanban/` and `threads/`. A look that refuses for the same reason changes nothing.
- **Count the window from when the story entered ready.** Simpler to state, but a story held by the limit for an hour would then be started the moment the limit frees, racing an idle agent whose `wait_for_work` answers in the same second. Counting from when attendance became the only hold gives that agent the window first.
- **A setting of its own for the hold.** The attended window already says how long a sign of an agent is trusted; a second number to explain and to set would say the same thing.
- **Fix only the agent's name.** It removes the six-minute gaps the operator saw, but an interactive session at work on its own story would still hold every ready story back for as long as it works.
