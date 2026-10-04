---
id: T-0819
type: task
nature: feature
title: The design and the user guide describe the three commands and the planning model
status: done
parent: S-0210
owner: alex
created: 2026-10-04T19:46:09Z
updated: 2026-10-04T20:15:58Z
transitions:
  - to: ready
    at: 2026-10-04T20:07:07Z
    by: agent-S-0210
  - to: in-progress
    at: 2026-10-04T20:07:08Z
    by: agent-S-0210
  - to: done
    at: 2026-10-04T20:15:58Z
    by: agent-S-0210
stream: S-0210
tags: [flai]
touches: [design/system/strategic-agents.md, design/system/flai-cli.md, design/system/project-manifest.md, docs/users/flai.md, docs/users/flai-reference.md, docs/operators/settings.md]
after: [T-0817]
usage:
  source: log
  seconds: 530
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 104
      output: 39889
      cache_read: 5187089
      cache_write: 148405
      cost: 2.7116
---
# T-0819 The design and the user guide describe the three commands and the planning model

## Work

Describe the commands and the model where readers look: `design/system/strategic-agents.md` (how the planner enriches a story: touches from co-change, the forecast model, the cost of delay formula, the `### Planning` notes), `design/system/flai-cli.md` (the three commands), `design/system/project-manifest.md` (`planning.default_duration`), `docs/users/flai.md` (how to use them), and the generated `docs/users/flai-reference.md` and settings flag table (`make flai-reference`).

Waits for T-0817, the last command it describes. It waited for T-0818 too until T-0818 came down to the operator pasting `.claude/agents/planner.md` (TH-0109), which the docs do not depend on.

## Done when

- [ ] Each document describes the commands as built, with the forecast model and the cost of delay formula stated once in `strategic-agents.md` and linked from the others
- [ ] `make flai-reference` leaves no diff, and the docs tests pass

## Notes
