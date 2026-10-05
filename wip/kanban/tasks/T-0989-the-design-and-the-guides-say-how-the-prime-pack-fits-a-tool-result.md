---
id: T-0989
type: task
nature: improvement
title: The design and the guides say how the prime pack fits a tool result
status: backlog
parent: S-0261
owner: alex
created: 2026-10-05T05:52:26Z
updated: 2026-10-05T05:52:26Z
transitions: []
stream: S-0261
tags: [flai, docs]
touches: [design/system/agent-context.md, design/system/flai-cli.md, docs/users/flai.md, docs/users/flai-reference.md, docs/operators/settings.md]
after: [T-0986]
---
# T-0989 The design and the guides say how the prime pack fits a tool result

## Work

Describe T-0986's ADR where the pack's budget is described now:

- `design/system/agent-context.md`, under "Fitting the pack to a budget";
- the `flai prime` entry of `design/system/flai-cli.md`;
- `flai prime` in `docs/users/flai.md` and `docs/users/flai-reference.md`;
- `prime.budget` in `docs/operators/settings.md`, if the default or a key changes.

Link the ADR from each. Say what an agent does when a pack is larger than one result: read the next part, or fetch what was catalogued with `doc_get`.

It waits for T-0986 because it describes what that ADR decides. It shares no path with T-0987, so the two run together.

## Done when

- Each of the five documents matches the ADR, and none still says that every brief is printed if the ADR changed that.
- `flai check --strict` and the markdown lint pass.

## Notes
