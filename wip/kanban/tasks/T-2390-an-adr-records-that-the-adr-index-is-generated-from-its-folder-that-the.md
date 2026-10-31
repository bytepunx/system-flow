---
id: T-2390
type: task
nature: improvement
title: An ADR records that the ADR index is generated from its folder, that the changelogs merge by union, and that the acceptance notice skips the shared paths
status: in-progress
parent: S-0550
owner: alex
created: 2026-10-31T07:43:04Z
updated: 2026-10-31T08:17:24Z
transitions:
  - to: ready
    at: 2026-10-31T08:17:23Z
    by: agent-S-0550
  - to: in-progress
    at: 2026-10-31T08:17:24Z
    by: agent-S-0550
stream: S-0550
tags: [cli]
touches: [design/adrs/README.md, design/adrs/0157-flai-stream-sync-and-flai-accept-merge-design-adrs-readme-md-s-index-row-by-row.md]
---
# T-2390 An ADR records that the ADR index is generated from its folder, that the changelogs merge by union, and that the acceptance notice skips the shared paths

## Work

Write one ADR with `flai adr new` for the three changes this story makes to what earlier ADRs decided. It refines ADR-0096 and ADR-0098, and supersedes ADR-0157 in part.

- **The ADR index is generated.** ADR-0157 keeps `design/adrs/README.md` hand-editable and rejects regenerating it, because the hand-written status cells carry notes such as "accepted, superseded in part by 0016" and "accepted, refines 0019, superseded by 0031". Decide how the generated status cell is derived from each ADR's front matter (`status`, `refines`, `supersedes`, `superseded_by`), including when a supersession reads "in part" (for example, the superseded ADR's `status` stays `accepted`). Decide what a rebase that stops on the index does: regenerate it, as ADR-0098 does for the issues summary, in place of ADR-0157's row-by-row merge.
- **The issues summary is checked against its folder.** `flai check` reports a summary that differs from what the issue files render, and `flai check --fix` rewrites it and the ADR index. No other file is fixed.
- **The changelogs merge by union.** `CHANGELOG.md` files take `merge=union`. Union keeps both sides' lines, so a story writes its lines under a `## Unreleased` section, and `flai release` names that section with the version. The digest has no index (`design/digest/README.md` says so), so no digest file takes the attribute: say so.
- **The acceptance notice skips the shared paths.** ADR-0096 §2 says the notice at acceptance still reports shared paths. This ADR reverses that for the acceptance notice alone. The trial merge at `flai stream sync` keeps reporting them, and so does `flai task done`'s message, which is about a task's commit while both stories are open. Give the reason: the trial merge already decides whether a shared edit conflicts.

Set ADR-0157's `superseded_by` through flai, never by hand. First layer: it waits for nothing. The ADR index task waits for it, since that task builds the status cell this ADR decides, and both write `design/adrs/README.md`.

## Done when

The ADR has its Context, Decision, Consequences, and Alternatives considered, cites the story's evidence, and has its row in `design/adrs/README.md`. ADR-0157 names it in `superseded_by`. `flai check --strict --story S-0550` reports nothing new. This task supplies the decisions behind story criteria 1 to 3.

## Notes

Drafted by the planner of S-0550 on 2026-10-31. The ADR's file name takes the next free number when it is written, so the task's touches name only the index and ADR-0157. `flai task done` widens them with the new file.
