---
id: T-3099
type: task
nature: improvement
title: Say in the design and the user docs that the section merge settles a paragraph both sides change in different sentences
status: done
parent: S-0719
owner: alex
created: 2026-11-11T08:54:49Z
updated: 2026-11-11T14:35:38Z
transitions:
  - to: ready
    at: 2026-11-11T14:24:27Z
    by: agent-S-0719
  - to: in-progress
    at: 2026-11-11T14:24:27Z
    by: agent-S-0719
  - to: done
    at: 2026-11-11T14:35:38Z
    by: agent-S-0719
stream: S-0719
tags: [docs, design]
touches: [design/system/workflow.md, design/system/agent-coordination.md, design/system/flai-cli.md, docs/users/flai.md]
after: [T-3094]
usage:
  source: log
  seconds: 651
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 33
      output: 223
      cache_read: 1136580
      cache_write: 49098
      cost: 0.439
---
# T-3099 Say in the design and the user docs that the section merge settles a paragraph both sides change in different sentences

## Work

Where the documents describe the section merge and its same-lines stop, say that a paragraph both sides change is merged by its sentences, both kept when they touch different sentences, and that the stop stays when both change one sentence, a table row, or a fenced line; link the new ADR:

- `design/system/workflow.md`: the S-0401 paragraph on ADR-0175 and the S-0637 paragraph on ADR-0276; add a sentence or a short paragraph for this story beside them.
- `design/system/agent-coordination.md`: the ADR-0175 paragraph.
- `design/system/flai-cli.md`: the `flai accept` and `flai stream sync` rows where they name the same-lines stop.
- `docs/users/flai.md`: the Shared documents bullet under sync.

Change only those sentences; bump each file's `updated`.

Waits for T-3094: the documents link and restate that ADR's decision.

## Done when

- The four documents say the same thing as the ADR, each links it, and none still says that any change two sides make to one line stops the merge.
- `flai check --strict` and the markdown lint pass on the four files.

## Notes
