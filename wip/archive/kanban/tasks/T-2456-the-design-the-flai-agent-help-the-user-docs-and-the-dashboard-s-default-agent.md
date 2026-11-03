---
id: T-2456
type: task
nature: feature
title: The design, the flai agent help, the user docs, and the dashboard's default agent panel say a story with no agent of its own starts with the default
status: done
parent: S-0573
owner: alex
created: 2026-11-03T17:28:57Z
updated: 2026-11-03T17:36:45Z
transitions:
  - to: ready
    at: 2026-11-03T17:30:02Z
    by: system-flow
  - to: in-progress
    at: 2026-11-03T17:30:03Z
    by: system-flow
  - to: done
    at: 2026-11-03T17:36:45Z
    by: agent-S-0573
stream: S-0573
tags: []
touches: [design/adrs/0251-a-story-whose-agent-names-no-harness-is-started-with-the-project-s-default.md, design/system/work-hierarchy.md, design/system/project-manifest.md, design/system/flai-cli.md, flai/cmd/agent.go, docs/users/flai-reference.md, docs/users/flai.md, flaiover/src/lib/components/SettingsPanel.svelte]
usage:
  source: log
  seconds: 402
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 28
      output: 6801
      cache_read: 1511918
      cache_write: 46465
      cost: 0.7501
---
# T-2456 The design, the flai agent help, the user docs, and the dashboard's default agent panel say a story with no agent of its own starts with the default

## Work

Say what ADR-0251 decides wherever the docs say the default reaches new stories only: `design/system/work-hierarchy.md` (a story's `agent`), `design/system/project-manifest.md` (`agent`), `design/system/flai-cli.md` where it describes the default or the start's harness check, the `flai agent` help and its `default agent:` line in `flai/cmd/agent.go`, `docs/users/flai-reference.md` regenerated from the help (`make flai-reference`), `docs/users/flai.md` where it says so, the dashboard's Default agent heading in `SettingsPanel.svelte`, and a line under `## Unreleased` in `CHANGELOG.md`. It also carries the corrected ADR-0242 link in ADR-0251. Waits for nothing: first layer, no path in common with T-2454.

## Done when

- No document says a story made before the default was set cannot start or must be edited.
- Each changed design paragraph links ADR-0251.
- `flai test` on the changed paths passes.

## Notes
