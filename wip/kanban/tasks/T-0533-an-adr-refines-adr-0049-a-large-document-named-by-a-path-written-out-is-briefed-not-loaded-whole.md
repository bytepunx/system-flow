---
id: T-0533
type: task
nature: improvement
title: "An ADR refines ADR-0049: a large document named by a path written out is briefed, not loaded whole"
status: done
parent: S-0149
owner: alex
created: 2026-09-29T05:25:25Z
updated: 2026-09-29T05:27:16Z
transitions:
  - to: ready
    at: 2026-09-29T05:25:47Z
    by: agent-S-0149
  - to: in-progress
    at: 2026-09-29T05:25:48Z
    by: agent-S-0149
  - to: done
    at: 2026-09-29T05:27:16Z
    by: agent-S-0149
stream: S-0149
tags: []
touches: [design/adrs, design/system/agent-context.md, design/system/flai-cli.md]
---
# T-0533 An ADR refines ADR-0049: a large document named by a path written out is briefed, not loaded whole

## Work

- Record with `flai adr new --refines ADR-0049 --status accepted` the decision from TH-0032: a document the story, its epic, or its tasks name only by its path written out in plain text is briefed when it is larger than an eighth of the pack's budget (10 KB at the default 80 KB), and loads whole at or under it. A markdown link, a `#fragment` link, and an ADR ID still load as ADR-0049 says. The brief says it was named, and tells the agent to read the document, or the sections it changes, before relying on it or changing it.
- Update the living design: the pack's rules and the measurements in `agent-context.md`, and the `flai prime` row in `flai-cli.md`, each linking the new ADR.

## Done when

- The ADR is accepted, refines ADR-0049, and `flai check --strict` is clean.
- `agent-context.md` and `flai-cli.md` describe the rule and link the ADR.

## Notes
