---
id: T-2500
type: task
nature: feature
title: A guard refusal of the planner inside its container is shown from the run's log, or from a call made by hand in the container with the planner's role
status: backlog
parent: S-0567
owner: alex
created: 2026-11-07T03:43:37Z
updated: 2026-11-07T03:44:14Z
transitions: []
stream: S-0567
tags: [docker, cli]
after: [T-2499]
---
# T-2500 A guard refusal of the planner inside its container is shown from the run's log, or from a call made by hand in the container with the planner's role

## Work

Waits for T-2499: it reads the log of the run T-2499 watched.

Search the run's session log under `.flai-cache/agents/plan` for a `PreToolUse` refusal by `flai guard`: an Edit or Write, a git write, a move past `backlog`, or another call the planner's role file refuses on the host (ADR-0247 § Decision 4). Record the call, the refusal's text, and the log's path.

When the log holds none, make the call by hand: run the image the run used with ADR-0247's mounts and `FLAI_ROLE=plan`, `FLAI_STARTED_BY=flai-serve`, and the run's `FLAI_AGENT`, and feed the checkout's guard hook a `PreToolUse` call the planner is refused on the host, such as an Edit of a file under `design/`. It must exit 2 with the same refusal the host gives. Run the same call through the host's `flai guard` and compare the two. Write the commands as a fenced block in the Notes.

Record the result in S-0567's Notes under `### Live run`, with `flai edit S-0567`. A refusal that differs from the host's, or a call that passes inside but is refused on the host, is a limitation for T-2501 to file.

## Done when

- S-0567's Notes show at least one refusal of the planner inside the container, from the run's log or from the hand call, with the host's refusal of the same call beside it.

## Notes
