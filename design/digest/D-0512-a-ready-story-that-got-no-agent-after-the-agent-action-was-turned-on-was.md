---
id: D-0512
title: "A ready story that got no agent after the agent action was turned on was stalled by its default agent naming no harness with no host command, not by when it entered ready; the reason is in serve.log and the story's waiting, never the journal"
kind: cause
sources: [S-0589, I-0300]
topics: [flai, cli, go, code]
paths: [flai/internal/serve/agents.go, flai/internal/serve/agents_test.go, docs/operators/index.md]
date: 2026-11-07
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0512 A ready story that got no agent after the agent action was turned on was stalled by its default agent naming no harness with no host command, not by when it entered ready; the reason is in serve.log and the story's waiting, never the journal

I-0300 guessed that a story already in ready when the operator turns the `agent` host action on gets no agent until someone presses **Start agent**. S-0589 found that guess does not hold. `(*launcher).look` in `flai/internal/serve/agents.go` reads the `agent` action at every look and starts a ready story that has had no run since it entered ready, whenever it entered ([ADR-0041](../adrs/0041-a-story-in-ready-gets-its-agent-whether-it-entered-ready-before-or-after-flai.md)).

The stall on the imported project came from its default agent naming no harness, with no command set on the host. `serve.log` read `agent not started`: `S-0001 names no harness, and no command is set on the host`, and S-0001 started twelve seconds after a default agent with harness claude-code was set. The operator saw nothing in `flai serve journal`, since the journal records only starts and failures. The reason is in `serve.log` and the story's `waiting` (`agent.status`, the card).

S-0589 changed no launcher code. It pinned the behaviour in a test in `flai/internal/serve/agents_test.go`, corrected `docs/operators/index.md` § "Starting an agent when a story becomes ready", and closed I-0300 naming S-0587 ("Default harness and agent should work out of the box"), which owns the missing-harness fix.

- S-0589: "So the stall was the imported project's default agent naming no harness, which S-0587 owns. The operator saw no journal entry because the journal records only starts and failures; the reason was in `serve.log` and the story's `waiting`."

Source: the Decisions of wip/archive/agents/S-0589.md; I-0300.
