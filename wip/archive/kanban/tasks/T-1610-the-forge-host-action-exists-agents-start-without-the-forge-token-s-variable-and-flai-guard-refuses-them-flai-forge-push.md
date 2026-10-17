---
id: T-1610
type: task
nature: feature
title: The forge host action exists, agents start without the forge token's variable, and flai guard refuses them flai forge push
status: done
parent: S-0393
owner: alex
created: 2026-10-16T04:35:26Z
updated: 2026-10-17T07:57:25Z
transitions:
  - to: ready
    at: 2026-10-17T07:44:47Z
    by: agent-S-0393
  - to: in-progress
    at: 2026-10-17T07:44:50Z
    by: agent-S-0393
  - to: done
    at: 2026-10-17T07:57:25Z
    by: agent-S-0393
stream: S-0393
tags: [cli]
touches: [flai/internal/hostapi/writes.go, flai/internal/hostapi/writes_test.go, flai/internal/serve/start.go, flai/internal/serve/start_test.go, flai/internal/guard/guard.go, flai/internal/guard/guard_test.go, flai/cmd/serve_actions.go, flai/cmd/serve_actions_test.go, flai/internal/guard/call.go, flai/internal/serve/agents.go, flai/internal/serve/agents_test.go]
usage:
  source: log
  seconds: 755
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 102
      output: 42850
      cache_read: 5516883
      cache_write: 162541
      cost: 2.8881
---
# T-1610 The forge host action exists, agents start without the forge token's variable, and flai guard refuses them flai forge push

## Work

- In `flai/internal/hostapi/writes.go`, add `ActionForge`, `forge`, to the `Actions` map. Its text says what it lets flai serve do: push story branches and open, update, and merge pull requests with the token the host names. It is off by default and enabled with `flai serve enable forge`.
- In `flai/internal/serve/start.go`, remove the variable each forge entry names, `GITHUB_TOKEN` by default, from the environment of every agent flai serve starts: story agents, sub-agents, and strategic agents. Do this where the harness's environment is built, if that is elsewhere.
- In `flai/internal/guard/guard.go`, refuse `flai forge push` and `flai forge merge` to every agent, naming the host action as the way.
- Test each part.
- First in its layer: it waits for no task.

## Done when

`flai serve actions` lists `forge`, an agent's environment lacks the token's variable, the guard refuses `flai forge push` to an agent, and the tests pass.

## Notes
