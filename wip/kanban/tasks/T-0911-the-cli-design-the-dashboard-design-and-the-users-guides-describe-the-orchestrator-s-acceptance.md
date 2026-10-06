---
id: T-0911
type: task
nature: feature
title: The CLI design, the dashboard design, and the users' guides describe the orchestrator's acceptance
status: done
parent: S-0221
owner: alex
created: 2026-10-05T04:48:12Z
updated: 2026-10-06T11:39:32Z
transitions:
  - to: ready
    at: 2026-10-06T11:35:19Z
    by: agent-S-0221
  - to: in-progress
    at: 2026-10-06T11:35:19Z
    by: agent-S-0221
  - to: done
    at: 2026-10-06T11:39:32Z
    by: agent-S-0221
stream: S-0221
tags: [docs]
touches: [design/system/flai-cli.md, design/system/flaiover-dashboard.md, docs/users/flai.md, docs/users/flai-reference.md, docs/users/flaiover.md]
after: [T-0907]
usage:
  source: log
  seconds: 252
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 62
      output: 345
      cache_read: 2688329
      cache_write: 119679
      cost: 1.2081
---
# T-0911 The CLI design, the dashboard design, and the users' guides describe the orchestrator's acceptance

## Work

Describe what T-0907 and the T-0898 ADR made true:

- `design/system/flai-cli.md`, at `flai accept`: `--by orchestrator`, the verified commit and evidence flags, the four extra blockers, the permission check, and the evidence in the story's Notes. Link the ADR.
- `docs/users/flai.md` § Accept and release: the same, for a reader turning `accept_reviews` on.
- `docs/users/flai-reference.md`: the new flags on `flai accept`.
- `design/system/flaiover-dashboard.md` and `docs/users/flaiover.md` § Reviewing a story: the review and story pages name who accepted, and link an orchestrator's evidence, as T-0909 shows them. Read T-0909's diff before you write this part.

This task waits for T-0907, whose flags and output it describes. It runs beside T-0909, which touches none of these files.

## Done when

- Each document above describes the orchestrator's acceptance, its flags, and where its evidence is shown.
- `scripts/lint-md.sh` and `flai check --strict` pass.

## Notes
