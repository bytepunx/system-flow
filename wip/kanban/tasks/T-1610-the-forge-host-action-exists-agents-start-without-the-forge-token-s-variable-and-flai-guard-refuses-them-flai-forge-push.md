---
id: T-1610
type: task
nature: feature
title: The forge host action exists, agents start without the forge token's variable, and flai guard refuses them flai forge push
status: backlog
parent: S-0393
owner: alex
created: 2026-10-16T04:35:26Z
updated: 2026-10-16T04:35:26Z
transitions: []
stream: S-0393
tags: [cli]
touches: [flai/internal/hostapi/writes.go, flai/internal/hostapi/writes_test.go, flai/internal/serve/start.go, flai/internal/serve/start_test.go, flai/internal/guard/guard.go, flai/internal/guard/guard_test.go]
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
