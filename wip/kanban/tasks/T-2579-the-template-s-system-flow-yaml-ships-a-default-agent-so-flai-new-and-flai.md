---
id: T-2579
type: task
nature: remediation
title: The template's system-flow.yaml ships a default agent, so flai new and flai import write one into every new project
status: done
parent: S-0587
owner: alex
created: 2026-11-07T10:12:05Z
updated: 2026-11-07T10:51:56Z
transitions:
  - to: ready
    at: 2026-11-07T10:31:56Z
    by: system-flow
  - to: in-progress
    at: 2026-11-07T10:31:58Z
    by: system-flow
  - to: done
    at: 2026-11-07T10:51:56Z
    by: agent-S-0587
stream: S-0587
tags: [template, cli]
touches: [template/root/system-flow.yaml.tmpl, template/CHANGELOG.md, flai/cmd/new_test.go, flai/cmd/import_test.go]
after: [T-2577]
usage:
  source: log
  seconds: 1198
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 26
      output: 7564
      cache_read: 4934943
      cache_write: 42678
      cost: 1.3985
---
# T-2579 The template's system-flow.yaml ships a default agent, so flai new and flai import write one into every new project

## Work

It waits for T-2577, whose ADR names the harness and model the template ships.

- In `template/root/system-flow.yaml.tmpl`, uncomment the `agent:` block with that harness and model. Keep `provider` and `image` as commented examples, since a new project has neither.
- Add a line under `## Unreleased` in `template/CHANGELOG.md`.
- Add tests to `flai/cmd/new_test.go` and `flai/cmd/import_test.go` showing that a project made by `flai new`, and one made by `flai import`, carry that default agent in `system-flow.yaml`, and that `flai story new` copies it into a new story.

## Done when

- A new project and an imported one have `agent.harness` set in their manifest.
- The two tests pass, and the template tier passes with `flai test` on the changed paths.

## Notes
