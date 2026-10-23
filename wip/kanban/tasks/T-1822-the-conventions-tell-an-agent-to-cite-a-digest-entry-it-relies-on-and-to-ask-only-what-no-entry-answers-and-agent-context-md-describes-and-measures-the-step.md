---
id: T-1822
type: task
nature: feature
title: The conventions tell an agent to cite a digest entry it relies on and to ask only what no entry answers, and agent-context.md describes and measures the step
status: done
parent: S-0430
owner: alex
created: 2026-10-16T06:05:07Z
updated: 2026-10-23T19:49:07Z
transitions:
  - to: ready
    at: 2026-10-23T19:43:33Z
    by: agent-S-0430
  - to: in-progress
    at: 2026-10-23T19:43:34Z
    by: agent-S-0430
  - to: done
    at: 2026-10-23T19:49:07Z
    by: agent-S-0430
stream: S-0430
tags: [cli, template]
touches: [design/conventions/work-management.md, template/root/design/conventions/work-management.md, design/conventions/strategic-agents.md, template/root/design/conventions/strategic-agents.md, design/system/agent-narrative.md, design/system/agent-context.md, design/system/flai-cli.md, docs/users/flai.md, template/CHANGELOG.md, template/template.yaml]
after: [T-1820]
usage:
  source: log
  seconds: 333
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 63
      output: 25447
      cache_read: 3982383
      cache_write: 135122
      cost: 2.1131
---
# T-1822 The conventions tell an agent to cite a digest entry it relies on and to ask only what no entry answers, and agent-context.md describes and measures the step

## Work

- `work-management.md` (template and here): where it says to ask the designer on a thread, add that an agent first looks for an entry in its pack or with `doc_search`, cites `D-nnnn` in its narrative or thread when it relies on one, and asks only what no entry answers.
- `strategic-agents.md` (template and here), "As the orchestrator": cite the entry as a recommendation's `source` when one answers the question.
- `design/system/agent-narrative.md`: what an agent is told about the digest and how its citations keep an entry from compaction.
- `design/system/agent-context.md`: a section on the digest step as built, and the pack sizes of three archived stories with and without it, measured with `scripts/flai.sh prime --story <id> --json` after the seed, none over its budget by the step alone.
- `docs/users/flai.md`, `template/CHANGELOG.md`, and the version in `template/template.yaml`.

Waits for T-1820, so the measured packs hold both steps. Runs beside the prime and prompt task.

## Done when

The documents say what the code does, the measurement is recorded, and `flai check --strict` and the markdown lint report nothing on them.

## Notes

The epic names a convention `agent-narrative.md`, which does not exist. The narrative's rules live in `work-management.md` and `design/system/agent-narrative.md`, so this task changes those.
