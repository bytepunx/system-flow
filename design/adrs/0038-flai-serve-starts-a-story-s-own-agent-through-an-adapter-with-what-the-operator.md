---
id: ADR-0038
title: "flai serve starts a story's own agent through an adapter, with what the operator allows"
status: accepted
date: 2026-09-23
supersedes: []
superseded_by: []
refines: [ADR-0029, ADR-0037]
---

# ADR-0038 flai serve starts a story's own agent through an adapter, with what the operator allows

## Context

S-0079 had `flai serve` start one command the operator wrote, for one ready story at a time, and only while nobody was attending the project ([ADR-0029](0029-the-dashboard-reaches-a-project-only-through-flai-on-the-host.md)). S-0103 gave each story an agent, a harness, a model, and options ([ADR-0037](0037-a-story-carries-its-agent-copied-from-the-project-s-default-when-it-is-made.md)). S-0104 asks for three things: a story moved to or created in ready is worked to review by the agent it names; questions go through the inbox and the answer is watched for; and the board shows each agent working, waiting, or failed.

## Decision

**flai serve starts each ready story's own agent, through an adapter for its harness, with the program and permissions the operator set on the host.** The operator asked for it in S-0104 (2026-09-23).

- **Adapters** (`internal/harness`) turn a story's agent into an argument list, run as it stands, never through a shell.
  - `claude-code` runs `claude -p` headless in the project's directory. It passes the story's model and a prompt to work that story alone to review by the project's conventions. Its only MCP server is flai's (`--mcp-config` with this flai's `mcp`, `--strict-mcp-config`), and its session is one flai names.
  - `command` is the operator's S-0079 command, now with `{model}` and `{harness}` too. A story that names no harness uses it, when it is set.
- **What runs, and what the agent may do, are the operator's alone.** Each harness has a program and arguments set on the host (`flai serve agent harness`). For `claude-code` the default is `--permission-mode acceptEdits --allowedTools Bash,mcp__flai`. A story names only the tunables its adapter checks (for `claude-code`: `effort`, `max_budget_usd`, `fallback_model`), because whoever can edit a story in the dashboard would otherwise choose what runs on the operator's host.
- **One agent per story, many per project, within the in-progress limit.** An agent started for a story still in ready counts against the limit. A story is started once each time it enters ready; a restart of flai serve still starts nothing that was already ready. Each agent works as `<name>-<story>`, so its cursor and `wait_for_work` are its own, and its own signs do not count as someone attending.
- **An agent's outcome is where it left its story.** Review or done means it worked. An open question of its own on the story means it asked. Anything else means it failed, with the reason: where it left the story, the block, the exit code.
- **An agent that asked is waiting, and is started again when it is answered.** The prompt tells it to hold `wait_for_events` until the thread is answered. If it ends anyway, flai starts it again, in the same session, as soon as someone else writes on the thread or resolves it, and tells it which question was answered. The designer answers where every question is answered: the inbox and the story's page.
- **The board shows it.** `agent.status` carries each story's activity: working, waiting (it asked, or its story is blocked), failed, or worked. A card has a dot, green, yellow, or red, and the story's page says which harness and model, since when, the question, or why it failed.

## Consequences

- Moving a card to ready now spends the operator's model budget without anyone at a keyboard, once the action is enabled. Enabling it remains a deliberate step on the host, and `max_budget_usd` lets a story cap its own spend.
- A headless agent's context ends with its process, except through the harness's session: `claude-code` resumes it. The operator's command gets `FLAI_ANSWERED` and must resume on its own.
- The trial found that the template's `wip/agents/README.md` counted as "someone attending", and that the operator's own MCP connectors loaded into a flai-started agent. Both are fixed here: only a story's narrative counts, and the agent's MCP servers are flai's alone.
- An agent that fails is not retried. Its story waits, red, until someone moves it to ready again.

## Alternatives considered

- **The story's config as flags passed straight to the harness.** It needs no adapter, but a story would then choose permissions and programs.
- **Keep one agent per project.** Simpler, but the criterion asks that every ready story be worked, until the WIP limit or the inbox stops it.
- **Have the agent hold `wait_for_events` for as long as it takes, and nothing else.** The trial's first run showed a model ending its turn instead. Resuming on the answer makes the outcome independent of the model, and a waiting session costs nothing.
- **An SDK-based runner inside flai.** It would tie flai to one vendor's API. The CLI harnesses already carry the tool loop, permissions, and sessions.
