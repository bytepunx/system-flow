---
id: T-1496
type: task
nature: improvement
title: Document the close-out's rounds and verify's reuse across a rebase in the design, the user docs, and the work-management convention
status: backlog
parent: S-0374
owner: alex
created: 2026-10-12T17:34:24Z
updated: 2026-10-12T17:34:24Z
transitions: []
stream: S-0374
tags: [cli, conventions, template]
touches: [docs/users/flai.md, design/system/workflow.md, design/system/flai-cli.md, design/conventions/work-management.md, template/root/design/conventions/work-management.md]
after: [T-1494, T-1495]
---
# T-1496 Document the close-out's rounds and verify's reuse across a rebase in the design, the user docs, and the work-management convention

## Work

Waits for T-1494 and T-1495, so it describes what they built.

- `docs/users/flai.md` § Verify a story before review: the resume across a clean rebase, which tiers it reuses and which main's paths send back, and the close-out's rounds where it covers the close-out's check.
- `design/system/flai-cli.md`, the verify package's description, and `design/system/workflow.md`, where the close-out and the definition of done are described: the same, linking the new ADR.
- `design/conventions/work-management.md` § Rules and its template copy `template/root/design/conventions/work-management.md`: replace "Any other commit the branch lacks, such as an acceptance, stops the close-out: run `flai stream sync` and run it again" with what the close-out now does on its own, and when it still stops. Both copies change in the same commit; this is a baseline rule, so the plan thread asks the operator first.

## Done when

- The four documents and the template copy say what T-1494 and T-1495 built, and link the ADR.
- The markdown lint and `flai check --strict` are clean.

## Notes

Layer 3.
