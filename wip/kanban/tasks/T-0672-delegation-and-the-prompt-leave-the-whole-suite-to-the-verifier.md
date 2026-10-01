---
id: T-0672
type: task
nature: improvement
title: Delegation and the prompt leave the whole suite to the verifier
status: done
parent: S-0189
owner: arobson
created: 2026-10-01T10:44:01Z
updated: 2026-10-01T10:46:53Z
transitions:
  - to: ready
    at: 2026-10-01T10:44:11Z
    by: agent-S-0189
  - to: in-progress
    at: 2026-10-01T10:45:06Z
    by: agent-S-0189
  - to: review
    at: 2026-10-01T10:46:50Z
    by: agent-S-0189
  - to: done
    at: 2026-10-01T10:46:53Z
    by: agent-S-0189
stream: S-0189
tags: []
touches: [flai/internal/harness, design/conventions/delegation.md, template/root/design/conventions/delegation.md]
usage:
  source: log
  seconds: 104
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 57
      output: 12838
      cache_read: 2999961
      cache_write: 59103
      cost: 1.3298
---
# T-0672 Delegation and the prompt leave the whole suite to the verifier

## Work

Change `delegation.md` (template baseline first, then this repository's copy) and the `claude-code` prompt (`harness.delegation`) so the story's agent runs only the tests for what it changed. It leaves the whole suite, lint, and `flai check` to one verifier before review, plus one more after fixing what that verifier found, and fixes what a verifier finds itself. Update the prompt's test, the template changelog, and the user docs that quote the prompt.

## Done when

- `delegation.md` in both places and the prompt say it, and a harness test pins it.
- `make flai-test` passes.

## Notes
