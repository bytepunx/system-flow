---
id: T-3191
type: task
nature: feature
title: flai serve agent image command and auto set them, flai serve settings and agent show show them, and flai guard refuses them to served agents
status: done
parent: S-0759
owner: alex
created: 2026-11-12T16:05:38Z
updated: 2026-11-13T16:55:02Z
transitions:
  - to: ready
    at: 2026-11-13T16:30:12Z
    by: system-flow
  - to: in-progress
    at: 2026-11-13T16:30:13Z
    by: system-flow
  - to: done
    at: 2026-11-13T16:55:02Z
    by: agent-S-0759
stream: S-0759
tags: [cli]
touches: [flai/cmd/serve_agent.go, flai/cmd/serve_agent_image.go, flai/cmd/serve_agent_image_test.go, flai/cmd/serve_settings.go, flai/cmd/serve_settings_test.go, flai/cmd/serve_view.go, flai/internal/guard/hostconfig.go, flai/internal/guard/hostconfig_test.go, flai/internal/guard/guard_tables.go, docs/users/flai-reference.md, docs/operators/settings.md]
after: [T-3190]
usage:
  source: log
  seconds: 1489
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 53
      output: 20120
      cache_read: 3992851
      cache_write: 144348
      cost: 2.0483
---
# T-3191 flai serve agent image command and auto set them, flai serve settings and agent show show them, and flai guard refuses them to served agents

## Work

ADR-0295 point 7 names the commands, in place of the planned `--build` and `--no-auto-image` on `flai serve agent harness container`:

- `flai serve agent image command [--project <root>] -- <program> [args...]` sets `agent.projects.<root>.image.command`; `--clear` removes it, bringing back flai's own pull and build; with neither it prints the command set, or that flai's own build runs.
- `flai serve agent image auto on|off [--project <root>]` sets `agent.projects.<root>.image.auto`; with no argument it prints it.
- `flai serve agent image` alone prints both for the project.
- The project is the one in the working directory, or `--project <root>`, as `flai serve agent set` finds it; outside a project with no `--project` they refuse, naming the flag.

Put them in a new `flai/cmd/serve_agent_image.go`, registered from `serve_agent.go`, with tests in `serve_agent_image_test.go`. Show both in `flai serve agent show` and in `flai serve settings` (`serve_settings.go`, test in `serve_settings_test.go`), as the other per-project agent settings are shown.

`flai guard` refuses served agents the commands that write the host's configuration (ADR-0220, `flai/internal/guard/hostconfig.go`). Refuse `image command` with a program or `--clear`, and `image auto on|off`; leave the forms that only print open. Test both sides in `hostconfig_test.go`, in every form a shell line takes as the existing cases do.

It waits for T-3190, whose fields it sets. It shares no file with T-3192.

## Done when

- `flai test` on the changed files passes.
- `flai serve agent image --help` names both settings.

## Notes
