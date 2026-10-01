---
id: ADR-0064
title: "A story in ready or in progress with no agent run on this host is started here on the operator's word, and its agent is told who began it, where, and when"
status: accepted
date: 2026-10-01
supersedes: []
superseded_by: []
refines: [ADR-0043]
topics: [cli, dashboard, server-side, client-side]
---

# ADR-0064 A story in ready or in progress with no agent run on this host is started here on the operator's word, and its agent is told who began it, where, and when

## Context

A story's agent runs on the host whose `flai serve` started it, and the record of that run is the host's own `serve/agents.json`. Work items, narratives, and threads travel between clones through git. Story branches and worktrees do not: they stay where the agent made them until the story is accepted.

On 2026-10-01 an agent on another machine moved S-0173 to in-progress, opened its narrative, and wrote its tasks, and those reached this clone. On this host:

- the story page showed no agent block;
- `flai serve agent restart` refused it, because this host's serve had started no agent for it (ADR-0043);
- `flai serve agent start` refused it, because it was not in ready.

The designer answered the agent's question here, and nothing resumed, because the run that waited for the answer was on the other host. The only way to get an agent was to move the story back to ready, which writes a false transition into its history and metrics. ADR-0043 rejected that for the same reason.

## Decision

A story in ready or in progress that this host has no agent run for can be started here on the operator's word. Its agent is told who began the story, where, and when, and reconciles from what is committed.

1. **Restart starts it.** `flai serve agent restart`, `agent.restart`, and the dashboard's Start agent no longer refuse a story with no run here. A story in progress starts at once. One in ready starts, or is queued like any other retry when the in-progress limit is full or a claim holds it (ADR-0044). Restart still refuses while an agent this host started runs for the story or waits for an answer, and the rest of ADR-0043's refusals stand.
2. **The agent is told where the story was begun.** The run's prompt names the agent that last moved the story to in-progress and when. It also names the narrative's agent and the host that opened it, when the narrative says. The agent is asked to reconcile rather than start over: open the stream with `flai stream open`, which uses `story/S-nnnn` from this clone or the remote and creates it from the main branch otherwise; read the narrative's Current state, Next steps, and log, and the story's tasks; and go on from what is committed. The prompt names each of the story's threads with an entry by someone other than the agent since the story was begun, and asks the agent to read them.
3. **The narrative records its host.** `flai stream open` writes `host`, the host's name, into the narrative's front matter, and sets it again when it reopens a stream on another host. It is how this host tells a story begun elsewhere from one begun on this host outside `flai serve`.
4. **`flai stream open` reopens a stream.** For a story whose narrative exists and whose worktree does not, it leaves the narrative as it is and checks out the story branch: the local branch, else the remote's (`origin`), fetched, else a new one from the main branch. It refuses, as before, when the worktree exists too.
5. **The dashboard shows it.** `agent.status` gives each story in progress with no run here an entry: `waiting`, its `elsewhere` (agent, when, host), and a why that says where it was begun and that no agent is here. The story page offers Start agent, which restarts it, and the board card shows the line.
6. **The MCP server can start one.** The tools `agent_start` and `agent_restart` do what `flai serve agent start` and `restart` do, under the same `agent` host action, so that an operator's own agent can do it. `flai guard` refuses them to a sub-agent, as it refuses every tool that is not a read (ADR-0060).

## Consequences

- A host that is gone, or whose agent stopped, no longer strands a story that is already in progress. The operator starts it on another host without a false transition.
- An agent started here cannot see uncommitted work, or commits that were not pushed, on the host that began the story. It goes on from what reached the remote, and the narrative says what was done.
- Nothing stops two hosts from running agents for the same story at once. Starting one here is the operator's word, and the dashboard says where the story was begun before they give it.
- An agent waiting on a thread on one host is still not resumed when the thread is answered on another. Starting an agent here is the remedy when the first host will not resume it.
- A narrative opened before this has no `host`. Its story reads as begun on another host, which is right for any story this host's serve did not start, save one the operator works in a session of their own here.

## Alternatives considered

- **Move the story back to ready to get an agent.** ADR-0055 allows the move, but made for this it writes a false transition into the story's history and metrics. ADR-0043 rejected it for restarts.
- **Let `agent start` take a story in progress.** Start is for a story in ready, whatever the launcher's rules say. Restart already covers a story in progress whose agent is gone, and a missing run is one more way for it to be gone.
- **Push story branches as the agent works.** It would make the work itself portable, but it changes when flai pushes, which is the operator's (ADR-0048). The agent here fetches the branch when it was pushed, and starts from main otherwise.
- **Record the host in the story's front matter.** The story is the designer's document and its front matter is the workflow's. The narrative is the agent's, and the host is about the agent.
