---
id: T-2450
type: task
nature: experiment
title: A stop from the dashboard and a restart after an answered thread are each tried once on the first proof story's container, and what happened is in S-0565's Notes
status: backlog
parent: S-0565
owner: alex
created: 2026-11-03T01:24:29Z
updated: 2026-11-03T01:28:19Z
transitions: []
stream: S-0565
tags: [docker]
after: [T-2449]
---
# T-2450 A stop from the dashboard and a restart after an answered thread are each tried once on the first proof story's container, and what happened is in S-0565's Notes

## Work

Waits for T-2449: the first proof story must carry `harness: container`, and the host must be ready to start it.

Once `flai serve` starts the first proof story (the documentation or shell one) in its container, hold `wait_for_events` for its `agent` events and try two things on it, once each.

- The stop. Ask the operator on the plan's thread on S-0565 to stop that story's agent from the dashboard while it works. Record the `agent` event (`stopped`, and its `why`), whether `docker ps` still lists `flai-agent-<story>-<session>`, the card's status, and what the restart from the dashboard did: the same session resumed, or a new one.
- The restart after a thread. When the proof story's agent ends its session asking on a thread, the operator's answer should make `flai serve` start the container again in the same session, from the story's writable Claude Code configuration folder in its cache folder (ADR-0203, ADR-0205). Record the `ended` and `started` events, the session, and whether the agent went on from where it stopped. If its agent asks nothing in its run, ask the operator on S-0565's thread how to bring a question about; do not write on the proof story's threads yourself.

Write what happened in S-0565's `## Notes` under a heading of its own with `flai edit S-0565`, in the main checkout. Hand each limitation met to T-2451, which files it as an issue, by naming it in those Notes.

## Done when

- S-0565's Notes say, for the stop and for the restart after an answered thread, what was done, the events and container state seen, and how the run went on.
- Each limitation met is named there for T-2451 to file.

## Notes
