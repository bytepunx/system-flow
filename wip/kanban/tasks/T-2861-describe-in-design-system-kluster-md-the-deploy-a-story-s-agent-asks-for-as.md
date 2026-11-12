---
id: T-2861
type: task
nature: feature
title: Describe in design/system/kluster.md the deploy a story's agent asks for as decided, marking each part not yet built
status: backlog
parent: S-0659
owner: alex
created: 2026-11-10T08:21:00Z
updated: 2026-11-12T15:54:10Z
transitions: []
stream: S-0659
tags: [cli]
touches: [design/system/kluster.md]
after: [T-2860]
---
# T-2861 Describe in design/system/kluster.md the deploy a story's agent asks for as decided, marking each part not yet built

## Work

Edit `design/system/kluster.md` to say what T-2860's ADR decides, with a link to it, and bump `updated`. It waits for T-2860 because it links the ADR by the number and file name `flai adr new` gives it, and describes what that ADR decides.

- § The namespaces: a story's namespace lives from its first deploy, in progress or review, until acceptance or cancellation, and the limit on how many live at once.
- § The `kluster:` block: the key that sets that limit, and its default.
- § What happens when: the table the ADR puts in place of ADR-0254 §5's, with the send-back rows S-0694 added.
- § A failed deploy holds the story: how a failure of a deploy the agent asked for differs from a failure in review, which S-0694 built.
- § The `kluster` host action and § flai serve (S-0579): how the agent asks, how it waits, and what `flai guard` still refuses.
- § What the review checks see: the agent's checks runs and its own commands with the namespace's variables, and what is recorded before review.

Mark each part flai does not do yet with the story that builds it, S-0660, S-0661, S-0662, or S-0663, so a reader never takes the design for what is built. Leave every section that describes what is built today as it is, beyond saying what changes.

## Done when

- `design/system/kluster.md` describes the ADR's design as decided, links the ADR, and marks each part not yet built with the story that builds it.
- `flai check --strict` and the markdown lint pass on it.

## Notes
