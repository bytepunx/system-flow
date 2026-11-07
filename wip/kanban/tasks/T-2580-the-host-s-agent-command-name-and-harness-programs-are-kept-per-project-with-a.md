---
id: T-2580
type: task
nature: remediation
title: The host's agent command, name, and harness programs are kept per project, with a host-wide fallback, so setting them for one project changes no other
status: done
parent: S-0587
owner: alex
created: 2026-11-07T10:12:11Z
updated: 2026-11-07T10:52:40Z
transitions:
  - to: ready
    at: 2026-11-07T10:32:00Z
    by: system-flow
  - to: in-progress
    at: 2026-11-07T10:32:03Z
    by: system-flow
  - to: done
    at: 2026-11-07T10:52:40Z
    by: agent-S-0587
stream: S-0587
tags: [cli, config]
touches: [flai/internal/config/config.go, flai/internal/config/config_test.go, flai/cmd/serve_actions.go, flai/cmd/serve_actions_test.go, flai/internal/hostapi/settings.go, flai/internal/hostapi/settings_test.go, docs/operators/settings.md, docs/users/flai-reference.md]
after: [T-2577]
usage:
  source: log
  seconds: 1237
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 99
      output: 29261
      cache_read: 19091341
      cache_write: 165105
      cost: 5.4102
---
# T-2580 The host's agent command, name, and harness programs are kept per project, with a host-wide fallback, so setting them for one project changes no other

## Work

It waits for T-2577: ADR-0261 § 3 decides how per-project host settings are kept and how "every project" is asked for.

- In `flai/internal/config/config.go`, add `agent.projects` to `AgentStart`: a map from a project's root to that project's own `command`, `name`, and `harnesses`. Give a lookup, such as `AgentStart.ForProject(root)`, that answers the project's settings: a command or name set for the project wins, a harness entry set for the project replaces the host-wide entry of that harness whole, and the host-wide `agent.command`, `agent.name`, and `agent.harnesses` are the fallback. Validate each project's harness entries as `checkHarnesses` does. A configuration without `projects` loads and answers as before.
- In `flai/cmd/serve_actions.go`:
  - `agentConfig(root)` reads `ForProject(root)`.
  - `flai serve agent set`, `agent clear`, and `agent harness` write the entry of the project in the working directory, or of `--project <root>`, and the host-wide fallback only with `--all-projects`. Outside a project with neither flag, they refuse naming both.
  - `flai serve agent show` (and what `host.Settings` answers) says, for the project, each value and whether it is the project's or the host-wide fallback.
- In `flai/internal/hostapi/settings.go`, the dashboard's `settings.agent` and `settings.harness` write the project's entry, and the host-wide fallback only when the request says `everywhere: true`; `settings.get` answers both.
- Providers and `auto_restarts` stay host-wide.
- Add `agent.projects` and the flags to `docs/operators/settings.md`, and regenerate `docs/users/flai-reference.md` with `make flai-reference`.

## Done when

- A test sets the command and a harness program for project A and finds project B's unchanged; another sets them with `--all-projects` and finds both changed.
- An old host-wide config still loads and applies to every project.
- `flai test` on the changed paths passes, and so does the settings-index test.

## Notes
