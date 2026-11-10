---
id: T-2899
type: task
nature: improvement
title: The runbook, flai-cli.md, and docker.md say how to probe an agent image without acting on its story
status: done
parent: S-0682
owner: alex
created: 2026-11-10T08:46:24Z
updated: 2026-11-10T09:11:51Z
transitions:
  - to: ready
    at: 2026-11-10T09:08:42Z
    by: agent-S-0682
  - to: in-progress
    at: 2026-11-10T09:08:43Z
    by: agent-S-0682
  - to: done
    at: 2026-11-10T09:11:51Z
    by: agent-S-0682
stream: S-0682
tags: [docs, containers]
touches: [docs/operators/runbooks/agent-container.md, design/system/flai-cli.md, design/tech/docker.md]
after: [T-2897, T-2898]
usage:
  source: log
  seconds: 188
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 12
      output: 2306
      cache_read: 1558446
      cache_write: 23282
      cost: 0.5157
---
# T-2899 The runbook, flai-cli.md, and docker.md say how to probe an agent image without acting on its story

## Work

Waits for T-2897 and T-2898: it describes what pilot and the script do with a probe once both are built.

- In `docs/operators/runbooks/agent-container.md` § Run, replace the `prompt.md` lines, which no longer match the script (pilot builds the prompt since S-0507), with the three forms the script takes: `S-nnnn`, `--probe S-nnnn`, and `[--probe] S-nnnn -- cmd`. Say that a run without `--probe` is the story's agent working the live story, as `agent-S-nnnn`, and that `--probe` is the trial of an image: `probe-S-nnnn`, a new session, a prompt that reads and reports, and no `FLAI_STORY`. Give the capped trial's form, `FLAI_AGENT_MODEL=haiku FLAI_AGENT_MAX_BUDGET_USD=0.02 scripts/agent-container.sh --probe S-nnnn`.
- In `design/system/flai-cli.md` § pilot, the agent images' entrypoint, say what a start request's `probe` does and that `flai serve` never sets it.
- In `design/tech/docker.md`'s Agent images row, name `--probe` beside `scripts/agent-container.sh S-nnnn`.
- Bump each file's `updated`, and run `flai test` on the three files (markdown lint).

## Done when

- The runbook's Run section matches the script's usage, and all three documents name the probe.
- The markdown lint passes on the three files.

## Notes
