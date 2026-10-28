---
id: T-2284
type: task
nature: improvement
title: While its provider's budget state is set, no story's agent is started or restarted on any served project, and a start or restart asked for by hand says it is held and until when
status: backlog
parent: S-0534
owner: alex
created: 2026-10-28T21:26:06Z
updated: 2026-10-28T21:26:06Z
transitions: []
stream: S-0534
tags: [flai, serve, budget]
touches: [flai/internal/serve/agents.go, flai/internal/serve/agents_test.go, flai/internal/serve/start.go, flai/internal/serve/start_test.go, flai/internal/serve/restart.go, flai/internal/serve/restart_test.go]
after: [T-2283]
---
# T-2284 While its provider's budget state is set, no story's agent is started or restarted on any served project, and a start or restart asked for by hand says it is held and until when

## Work

Check the budget state (T-2282) before a story's agent is started. Resolve the story's provider first, from its `agent` and the host's providers, as `provide` does. Check before any run is recorded: a refusal inside `startOf` or `fail` is recorded as a failed run (`agents.go`). Check at every start path:

- `look` starting a ready story
- `restartEnded`
- `resume`, on an answered thread
- `startSentBack`, `startQueueSentBack`, and `startAcceptSentBack`
- `StartNow`/`startNow` in `restart.go`

While the provider is held, the story waits where it is. It stays in its column, takes no new run, and counts nothing against `auto_restarts`. Its next look starts it once the state clears. Once `until` has passed, the first start the launcher would make on that provider is the probe (`Probe`). Every other start stays held until the probe ends.

Every launcher serving a project reads the one `budget.json`, so a limit one project's run hit holds the others at their next look.

A start the operator asks for by hand is refused with nothing recorded. That covers `Start` in `start.go` and `Restart` in `restart.go`, and through them `flai serve agent start` and `restart`, the host API's, and the MCP tools'. The refusal says: `held: <provider> usage limit (<limit>) until <until>; flai serve starts it when the limit resets`.

This task waits for T-2283, which changes `agents.go` first. It runs beside the roles' hold, which touches none of these files.

## Done when

- [ ] With the state set, no ready, answered, or sent-back story starts on that provider, and a story on another provider starts as before.
- [ ] No ended story is restarted, and its `auto_restarts` count is unchanged.
- [ ] Once `until` has passed, exactly one story's agent starts as the probe. The rest start at the look after the probe clears the state.
- [ ] `flai serve agent start` and `restart` refuse with the provider and `until`, and record no run.
- [ ] Two projects under one serve folder both hold on the one state, tested.
- [ ] `flai test` passes on the changed paths.

## Notes

Planned by planner-S-0534.
