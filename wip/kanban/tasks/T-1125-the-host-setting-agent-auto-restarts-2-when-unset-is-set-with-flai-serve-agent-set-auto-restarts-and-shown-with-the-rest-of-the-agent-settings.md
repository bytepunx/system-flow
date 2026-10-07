---
id: T-1125
type: task
nature: remediation
title: The host setting agent.auto_restarts, 2 when unset, is set with flai serve agent set --auto-restarts and shown with the rest of the agent settings
status: done
parent: S-0294
owner: alex
created: 2026-10-06T23:07:46Z
updated: 2026-10-07T01:49:51Z
transitions:
  - to: ready
    at: 2026-10-07T01:43:22Z
    by: agent-S-0294
  - to: in-progress
    at: 2026-10-07T01:43:23Z
    by: agent-S-0294
  - to: done
    at: 2026-10-07T01:49:51Z
    by: agent-S-0294
stream: S-0294
tags: [cli, serve]
touches: [flai/internal/config/config.go, flai/internal/config/config_test.go, flai/cmd/serve_actions.go, flai/cmd/serve_actions_test.go, docs/operators/settings.md, docs/users/flai-reference.md]
usage:
  source: log
  seconds: 388
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 66
      output: 359
      cache_read: 3132129
      cache_write: 105950
      cost: 1.443
---
# T-1125 The host setting agent.auto_restarts, 2 when unset, is set with flai serve agent set --auto-restarts and shown with the rest of the agent settings

## Work

Add the limit of automatic restarts to this host's agent settings.

- **The field.** In `flai/internal/config/config.go`, add `AutoRestarts` to `AgentStart`, as `agent.auto_restarts`. Unset means 2, and 0 means off. A pointer or an explicit "set" marker keeps 0 apart from unset. Give it a method that answers the effective limit, and reject a negative value when the config is loaded or set.
- **The command.** In `flai/cmd/serve_actions.go`, `flai serve agent set` takes `--auto-restarts N`, and changes only that when given alone. The command that shows the agent settings prints the limit and whether it is the default.
- **The wiring.** Pass the limit into `serve.AgentConfig` where `serve_actions.go` builds it. Adding the field to `AgentConfig` itself is T-1126's.
- **The tests.** In `config_test.go`:
  - unset gives 2;
  - 0 gives off;
  - a negative value is refused;
  - the JSON round-trips.

  In `serve_actions_test.go`:
  - `--auto-restarts` alone keeps the command and the harnesses;
  - the setting is shown.

## Done when

- `flai serve agent set --auto-restarts 0` and `--auto-restarts 3` write `agent.auto_restarts` to the host config, and leave the other agent settings as they were.
- With the key unset, the effective limit is 2.
- The tests above pass with `scripts/flai-test.sh`.

## Notes

Drafted by the planner. The operator confirmed TH-0208's recommendation on 2026-10-06, which puts the setting on the host beside `agent.command`, not in the manifest.
