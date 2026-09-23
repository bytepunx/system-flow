---
id: T-0370
type: task
nature: feature
title: "Tried end to end with a real claude-code agent; ADR and docs"
status: done
parent: S-0104
owner: alex
created: 2026-09-23T17:25:02Z
updated: 2026-09-23T17:54:55Z
transitions:
  - to: ready
    at: 2026-09-23T17:39:14Z
    by: system-flow
  - to: in-progress
    at: 2026-09-23T17:39:14Z
    by: system-flow
  - to: done
    at: 2026-09-23T17:54:55Z
    by: system-flow
stream: S-0104
tags: []
---

# T-0370 Tried end to end with a real claude-code agent; ADR and docs

## Work
On a throwaway project, with a scratch `flai serve` from this tree and a dashboard on its own port: enable the agent action, set a default agent of claude-code with the Haiku model, and move a tiny story to ready. Watch a real `claude -p` pull it, open a thread, wait for the answer given from the dashboard, and move the story to review, with the dot going green, yellow, and green. Then an ADR and the docs.

## Done when
- The trial is recorded in the narrative.
- ADR written; flai-cli.md, flaiover-dashboard.md, and the user docs updated.

## Notes
