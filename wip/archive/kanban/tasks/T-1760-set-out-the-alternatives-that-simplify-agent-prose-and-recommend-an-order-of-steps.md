---
id: T-1760
type: task
nature: research
title: Set out the alternatives that simplify agent prose and recommend an order of steps
status: done
parent: S-0413
owner: alex
created: 2026-10-16T05:46:29Z
updated: 2026-10-16T06:06:15Z
transitions:
  - to: ready
    at: 2026-10-16T05:46:56Z
    by: agent-S-0413
  - to: in-progress
    at: 2026-10-16T06:04:42Z
    by: agent-S-0413
  - to: done
    at: 2026-10-16T06:06:15Z
    by: agent-S-0413
stream: S-0413
tags: []
touches: [design/system/agent-prose.md]
after: [T-1758]
usage:
  source: log
  seconds: 93
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 78
      output: 32403
      cache_read: 3503337
      cache_write: 123001
      cost: 2.1092
---
# T-1760 Set out the alternatives that simplify agent prose and recommend an order of steps

## Work

Set out the alternatives that reach the same end, simpler and shorter prose across flai, with or without a linter and a local model: tighter prompts and conventions, length budgets flai enforces on what it writes, structured fields in place of prose, a small hosted model as the rewriter, summaries the dashboard shows above long text, fewer or shorter documents an agent must read, and others the research finds. Compare them with the linter and model path on cost, effort, risk, and effect, and end the document with a recommended order of steps and the follow-on stories they would become.

Waits for T-1758, since the comparison and the recommendation weigh the model path.

## Done when

- `design/system/agent-prose.md` has a `## Alternatives` section and a `## Recommendation` that ends the document.
- Acceptance criterion 3 is met.

## Notes
