---
id: T-0534
type: task
nature: improvement
title: flai prime --story briefs a large document named by a path written out, and tells the agent to read it before changing it
status: done
parent: S-0149
owner: alex
created: 2026-09-29T05:25:25Z
updated: 2026-09-29T05:30:27Z
transitions:
  - to: ready
    at: 2026-09-29T05:25:48Z
    by: agent-S-0149
  - to: in-progress
    at: 2026-09-29T05:27:16Z
    by: agent-S-0149
  - to: done
    at: 2026-09-29T05:30:27Z
    by: agent-S-0149
stream: S-0149
tags: []
touches: [flai/internal/context, flai/cmd/prime.go]
---
# T-0534 flai prime --story briefs a large document named by a path written out, and tells the agent to read it before changing it

## Work

- In `flai/internal/context`, tell a path written out apart from a link and an ADR ID when finding what a work item names.
- In the named step, brief a document named only by a path written out when it is larger than an eighth of the budget, with the reason `named in <ID>`; a link or an ID to the same document still loads it whole.
- Say in the briefs group heading and on such a brief that it was named and is briefed for its size, and that the agent reads it, or the sections it changes, with `doc_get` and a heading, or `flai doc show`, before relying on it or changing it.
- Behaviour tests for each case: small and large plain path, link, fragment link, ADR ID, a document both linked and written out, and the heading text.

## Done when

- `make test` and lint pass, and the tests fail without the change.

## Notes
