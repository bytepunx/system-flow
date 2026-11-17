---
id: T-3591
type: task
nature: improvement
title: The strategic-agents convention, the users guide, the reference, and the changelog say the orchestrator asks for rebase and resolve
status: backlog
parent: S-0845
owner: alex
created: 2026-11-17T12:34:27Z
updated: 2026-11-17T12:34:27Z
transitions: []
stream: S-0845
tags: [docs, orchestrator]
touches: [design/conventions/strategic-agents.md, template/root/design/conventions/strategic-agents.md, docs/users/flai.md, docs/users/flai-reference.md, template/CHANGELOG.md]
after: [T-3589, T-3590]
---
# T-3591 The strategic-agents convention, the users guide, the reference, and the changelog say the orchestrator asks for rebase and resolve

## Work

Waits for T-3589 and T-3590, so the words describe the tool, the guard, and the role file as built and tested. In `design/conventions/strategic-agents.md` § Accepting a story, step 1 or step 4, and its template copy, say that a close-out that did not pass after main moved is asked for again as rebase and resolve with the MCP tool `resolve`, not left on a thread to the operator, and that a story held by `undone_main` stays the operator's. In `docs/users/flai.md`, where it says `--resolve` is the operator's alone, say the orchestrator may ask too under `accept_reviews`, and list the MCP tool. Regenerate `docs/users/flai-reference.md` with `make flai-reference` if the flag's help changed. Add a line to `template/CHANGELOG.md` under `## Unreleased` for the convention and role file change.

## Done when

- The convention and its template copy, the users guide, and the changelog say it, linking T-3587's ADR where a link fits.
- `flai test` passes on the paths changed, the markdown lint and `TestReferenceIsCurrent` among its tiers.

## Notes
