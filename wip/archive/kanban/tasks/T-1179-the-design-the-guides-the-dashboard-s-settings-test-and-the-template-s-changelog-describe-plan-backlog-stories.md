---
id: T-1179
type: task
nature: feature
title: The design, the guides, the dashboard's settings test, and the template's changelog describe plan_backlog_stories
status: done
parent: S-0328
owner: alex
created: 2026-10-07T19:40:42Z
updated: 2026-10-07T20:17:37Z
transitions:
  - to: ready
    at: 2026-10-07T20:02:02Z
    by: agent-S-0328
  - to: in-progress
    at: 2026-10-07T20:02:03Z
    by: agent-S-0328
  - to: done
    at: 2026-10-07T20:17:37Z
    by: agent-S-0328
stream: S-0328
tags: []
touches: [design/system/strategic-agents.md, design/system/project-manifest.md, design/system/flai-cli.md, docs/operators/index.md, docs/users/flaiover.md, docs/users/flai.md, flaiover/src/routes/workflow/orchestrator/orchestrator.svelte.test.ts, template/CHANGELOG.md]
after: [T-1176]
usage:
  source: log
  seconds: 934
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 89
      output: 38626
      cache_read: 6552308
      cache_write: 164709
      cost: 3.1131
---
# T-1179 The design, the guides, the dashboard's settings test, and the template's changelog describe plan_backlog_stories

## Work

Describe the permission as the ADR T-1174 records it, beside `plan_backlog_epics` everywhere that one is described:

- `design/system/strategic-agents.md`: the planner's "Starting it", "What it does with each permission", "Answering threads", and "Its permissions and the guard", linking the ADR.
- `design/system/project-manifest.md` and `design/system/flai-cli.md` (`flai plan --candidates`).
- `docs/operators/index.md`, `docs/users/flaiover.md`, and `docs/users/flai.md`.
- `orchestrator.svelte.test.ts`: the key in `ORCHESTRATION` and in the mocked settings.
- `template/CHANGELOG.md`: an entry for the convention change T-1178 makes.

It waits on T-1176, whose names and texts it quotes. It runs in layer 2 beside T-1177 and T-1178, sharing no path with them.

## Done when

- [ ] `flai test` passes on the paths changed.
- [ ] Every place that describes `plan_backlog_epics` describes `plan_backlog_stories` too.

## Notes
