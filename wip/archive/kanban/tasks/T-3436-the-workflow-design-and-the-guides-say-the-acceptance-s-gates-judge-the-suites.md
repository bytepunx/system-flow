---
id: T-3436
type: task
nature: feature
title: The workflow design and the guides say the acceptance's gates judge the suites its transition's verify lists
status: done
parent: S-0808
owner: alex
created: 2026-11-16T23:46:50Z
updated: 2026-11-27T23:35:32Z
transitions:
  - to: ready
    at: 2026-11-27T23:30:53Z
    by: agent-S-0808
  - to: in-progress
    at: 2026-11-27T23:30:53Z
    by: agent-S-0808
  - to: done
    at: 2026-11-27T23:35:32Z
    by: agent-S-0808
stream: S-0808
tags: [flai]
touches: [design/system/workflow.md, design/system/flai-cli.md, docs/users/flai.md, docs/operators/index.md, docs/users/flai-reference.md, flai/cmd/accept.go, flai/internal/preview/accept.go, flai/internal/preview/accept_test.go]
after: [T-3434, T-3435]
usage:
  source: log
  seconds: 279
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 88
      output: 28030
      cache_read: 6302228
      cache_write: 246760
      cost: 3.2853
---
# T-3436 The workflow design and the guides say the acceptance's gates judge the suites its transition's verify lists

## Work

ADR-0312 §3. It waits for T-3434 and T-3435, since it describes the gates and the refusal words they write.

- `design/system/workflow.md`, at the paragraphs on `queue_tiers` and `verify_behind_main`: the acceptance's gates read the suites the `verify:` of the transition with the `accept` effect lists. A suite it does not list never holds a story in review or sends it back. Describe the new refusal for a suite the last verify did not run. On the baseline, `verify: all`, the gates are as before.
- `design/system/flai-cli.md`, at the `flai accept` row and the paragraph on ADR-0291: the same, and the new blocker code.
- `docs/users/flai.md`, at the acceptance's waits, and `docs/operators/index.md`, at what the merge queue changes in an acceptance: which suites hold a story, and how to clear the new refusal with `flai verify S-nnnn --job`.

## Done when

- The four documents name the suites the gates judge, link ADR-0312, and give the refusal's way out.
- The markdown lint and `flai check --strict` pass on them.

## Notes
