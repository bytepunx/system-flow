---
id: T-0445
type: task
nature: improvement
title: settings.get reports every project flai serve serves and every one below an import folder it does not, and settings.serve and settings.unserve register and unregister one
status: done
parent: S-0122
owner: alex
created: 2026-09-26T06:05:08Z
updated: 2026-09-26T06:08:41Z
transitions:
  - to: ready
    at: 2026-09-26T06:05:19Z
    by: agent-S-0119
  - to: in-progress
    at: 2026-09-26T06:05:20Z
    by: agent-S-0119
  - to: done
    at: 2026-09-26T06:08:41Z
    by: agent-S-0119
stream: S-0122
tags: []
touches: [flai/internal/hostapi, flai/cmd/serve_actions.go, flai/cmd/serve_project.go]
---
# T-0445 settings.get reports every project flai serve serves and every one below an import folder it does not, and settings.serve and settings.unserve register and unregister one

## Work

In the flai serve process, `hostSettings` (`flai/cmd/serve_actions.go`) adds `projects` to `settings.get`'s host object: every served project with key, name, root, where it is served from, state, since, last error, and reason, from `listServedProjects`, and the unserved projects below import folders with why. Two new write specs in `flai/internal/hostapi/settings.go`, gated by `ActionSettings` for the project the request goes through: `settings.serve {root}` runs `flai serve project add -- <root>`, and `settings.unserve {key}` runs `flai serve project remove -- <key>`. Both are journalled as the other settings are. `serve project add` and `remove` accept `--`. The contract list in `flaiover/src/lib/server/agent.ts` gains both names.

## Done when

- a hostapi test covers each method's command line, its refusal without the settings action, and its journal entry
- `settings.get` carries the projects in a test through the serve process's host
- `make test` and the hostapi contract test pass

## Notes
