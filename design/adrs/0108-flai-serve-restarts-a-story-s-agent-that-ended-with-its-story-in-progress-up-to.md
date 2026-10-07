---
id: ADR-0108
title: "flai serve restarts a story's agent that ended with its story in progress, up to agent.auto_restarts times, then asks the operator on a thread"
status: proposed
date: 2026-10-07
supersedes: []
superseded_by: []
refines: [ADR-0043]
---

# ADR-0108 flai serve restarts a story's agent that ended with its story in progress, up to agent.auto_restarts times, then asks the operator on a thread

## Context

[ADR-0043](0043-flai-serve-starts-a-ready-story-s-agent-whenever-the-in-progress-limit-has-room.md) gave the restart of a story's agent that dropped or failed to a person: `flai serve agent restart`, or Retry on the story's page. flai serve itself starts a story's agent again only on an answer to the question it ended on.

[I-0084](../issues/I-0084-claude-code-ends-a-headless-agent-ten-minutes-after-its-turn-ends-even-while-its-background-sub-agent-is-still-working-and-flai-serve-leaves-the-story-in-progress-with-no-agent.md) recorded two runs that `claude -p` ended ten minutes after their turn ended, with a background sub-agent still at work. Each left its story in progress with no agent, judged failed, until the operator noticed and restarted it. [ADR-0092](0092-a-story-s-agent-waits-for-a-sub-agent-by-launching-it-in-the-foreground-and.md) gave the agent a safe way to wait, which removes that cause. Any other cause that ends a run early, such as a crash, a lost connection, or an agent that ends its turn by mistake, still leaves the story waiting for a person. On TH-0208 the operator chose to have flai serve restart such an agent itself, a limited number of times.

## Decision

flai serve restarts a story's agent that ended with its story in progress, up to `agent.auto_restarts` times, then asks the operator on a thread.

1. **When.** At a look, the launcher restarts the agent of a story whose newest run ended, when all of these hold:
   - the story is in progress and not blocked;
   - the run started and ended failed: it did not end asking a question, which the restart on an answer covers, and the operator did not stop it ([ADR-0058](0058-the-agent-host-action-lets-the-dashboard-stop-a-story-s-agent-and-a-stop.md));
   - the `agent` host action is on for the project;
   - the story's count of automatic restarts is under the limit.
2. **How.** As `flai serve agent restart` does: a new session, the same agent name, the run recorded in `serve/agents.json`, and the new agent told how the last one ended and that flai serve started it again on its own. The in-progress limit, review, and claims do not hold it back: the story is in progress already.
3. **The limit.** The host setting `agent.auto_restarts`, set with `flai serve agent set --auto-restarts N`, is how many times; 2 when unset, and 0 turns the restart off. The count is per story and is kept on its runs. It starts again from 0 with any run flai serve did not start on its own: the operator's restart or start, and the agent started when the story enters ready again.
4. **At the limit.** When a run ends that would be restarted but for the limit, flai serve opens one thread on the story, by `flai serve`. It says how many times the agent ended, how the last run ended, and that `flai serve agent restart` or Retry starts it again. The story is then left as before ADR-0043's restart, for the operator.

## Consequences

- A story whose agent ends early gets another within a look, without a person, up to the limit.
- A story whose agent fails the same way every time costs the limit's runs before the operator hears of it, and then hears of it once, on a thread.
- A story in progress whose newest run had already failed when a flai serve with this change first looks is restarted then, as if it had just ended.
- A run that could not be started at all is not restarted: it never ran, and the same start would fail again.

## Alternatives considered

- **No automatic restart, only the thread.** Offered on TH-0208 and not chosen: the story would still wait for a person each time.
- **A manifest key per project.** The agent settings are the host's, beside `agent.command`, and so is the `agent` action that gates restarting; the operator chose the host on TH-0208.
- **Restart without a limit.** An agent that fails the same way every time would run, and cost, without end.
