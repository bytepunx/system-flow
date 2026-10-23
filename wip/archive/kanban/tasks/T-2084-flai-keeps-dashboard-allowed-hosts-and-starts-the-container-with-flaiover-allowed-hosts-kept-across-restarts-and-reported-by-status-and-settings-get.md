---
id: T-2084
type: task
nature: feature
title: flai keeps dashboard.allowed_hosts and starts the container with FLAIOVER_ALLOWED_HOSTS, kept across restarts and reported by status and settings.get
status: done
parent: S-0489
owner: alex
created: 2026-10-23T18:08:11Z
updated: 2026-10-23T18:43:43Z
transitions:
  - to: ready
    at: 2026-10-23T18:22:01Z
    by: agent-S-0489
  - to: in-progress
    at: 2026-10-23T18:22:02Z
    by: agent-S-0489
  - to: done
    at: 2026-10-23T18:43:43Z
    by: agent-S-0489
stream: S-0489
tags: [flai, dashboard]
touches: [flai/internal/config/config.go, flai/internal/config/config_test.go, flai/cmd/dashboard.go, flai/cmd/dashboard_upgrade.go, flai/cmd/dashboard_watch.go, flai/cmd/dashboard_test.go, flai/cmd/dashboard_watch_test.go, flai/cmd/serve_actions.go, flai/cmd/serve_actions_test.go, docs/operators/settings.md, docs/users/flai.md, design/system/flai-cli.md]
usage:
  source: log
  seconds: 1301
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 112
      output: 36094
      cache_read: 7520441
      cache_write: 210035
      cost: 3.4799
---
# T-2084 flai keeps dashboard.allowed_hosts and starts the container with FLAIOVER_ALLOWED_HOSTS, kept across restarts and reported by status and settings.get

## Work

Add `dashboard.allowed_hosts` to the host configuration (`config.Dashboard` in `flai/internal/config/config.go`, and `Keys()`): a comma-separated list of host names, each with an optional port, that the dashboard accepts in a browser's `Host` and `Origin` beside the address it is published on; empty by default. A string, so `flai config set dashboard.allowed_hosts "flow.example.com, board.example.net:8443"` works as `Set` is; trim each entry and lower-case it when flai passes it on. It is a fact about one machine's shared container, so it lives in `~/.flai/config.json` and not in the manifest, like `dashboard.push_key` did.

Pass it to the container: `dashboardSettings` in `flai/cmd/dashboard.go` reads it, and `containerArgs` in `flai/cmd/dashboard_upgrade.go` adds `--env FLAIOVER_ALLOWED_HOSTS=<list>` when it is not empty. The container record (`dashboardRecord` in `flai/cmd/dashboard_watch.go`) keeps `allowed_hosts`, so that the watch starts the container again the same way; `flai dashboard restart` and `flai dashboard upgrade` take the configuration's current value, as `keepAllowance` does for the allowance, so that a changed setting applies on the next restart. `flai dashboard status` and `--json` report the list the running container was started with (`allowed_hosts`), so the Host page can say when the configuration and the container differ.

Put the configured value in `settings.get`'s host block (`hostSettings` in `flai/cmd/serve_actions.go`) as `allowed_hosts`, so the dashboard shows it.

Tests: `config_test.go` for the key and `Set`; `dashboard_test.go` for the env argument, the record, restart and upgrade keeping the current value, and status; `dashboard_watch_test.go` for a restart from the record; `serve_actions_test.go` for the host block. Document the key in `docs/operators/settings.md` § flai configuration and `FLAIOVER_ALLOWED_HOSTS` in its flaiover container table (`TestSettingsIndexIsComplete` fails otherwise), and the key in the configuration table of `docs/users/flai.md`.

This task waits for no other. The dashboard's reading of `FLAIOVER_ALLOWED_HOSTS` is T-2083's; the two agree on the variable's name and shape here: a comma-separated list of lower-case host names with optional ports.

## Done when

- `flai config set dashboard.allowed_hosts "<list>"` is accepted, `flai config get` prints it, and a start passes it to the container as `FLAIOVER_ALLOWED_HOSTS`.
- A restart by the watch, `flai dashboard restart`, and `flai dashboard upgrade` start the container with the configuration's current list, and `flai dashboard status --json` reports the list the container runs with.
- `settings.get` carries the configured list in its host block.
- `docs/operators/settings.md` has a row for the key and for the variable, `docs/users/flai.md` names the key, and `flai test` on the changed paths passes, `TestSettingsIndexIsComplete` among them.

## Notes

Planned by the planner on 2026-10-23. The key is a string rather than a list so that `config.Set`, which parses by the current value's type, needs no new case; the story's agent may make it a list if it extends `Set`.
