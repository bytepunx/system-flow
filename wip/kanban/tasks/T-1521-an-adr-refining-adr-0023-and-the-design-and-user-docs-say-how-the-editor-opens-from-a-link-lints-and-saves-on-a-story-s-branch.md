---
id: T-1521
type: task
nature: feature
title: An ADR refining ADR-0023 and the design and user docs say how the editor opens from a link, lints, and saves on a story's branch
status: done
parent: S-0378
owner: alex
created: 2026-10-13T16:57:37Z
updated: 2026-10-13T18:36:13Z
transitions:
  - to: ready
    at: 2026-10-13T18:32:21Z
    by: agent-S-0378
  - to: in-progress
    at: 2026-10-13T18:32:21Z
    by: agent-S-0378
  - to: done
    at: 2026-10-13T18:36:13Z
    by: agent-S-0378
stream: S-0378
tags: [dashboard, cli, docs]
touches: [design/adrs, design/adrs/README.md, design/system/flaiover-dashboard.md, design/system/flai-cli.md, docs/users/flai.md, docs/users/flaiover.md, docs/users/flai-reference.md]
after: [T-1520]
usage:
  source: log
  seconds: 232
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 63
      output: 24057
      cache_read: 3666301
      cache_write: 98275
      cost: 1.7721
---
# T-1521 An ADR refining ADR-0023 and the design and user docs say how the editor opens from a link, lints, and saves on a story's branch

## Work

- Write an ADR with `flai adr new` that refines ADR-0023. A save may name a story and then commits on that story's branch in its worktree. A wip path is never saved there. `doc.lint` answers flai's subset of the rules while the operator types, and the markdown tier stays the authority at acceptance.
- `design/system/flaiover-dashboard.md`: add `doc.lint`, `/api/docs/lint`, and the `story` parameter under § API (S-0011) and § Writes (S-0075, ADR-0029). Under § Workbench (E-0006), say how the editor opens over a page at a line, with a story, from a link, a report, a thread, or a refusal.
- `design/system/flai-cli.md`: the `flai doc` row in § Commands gains `lint` and `--story`. Say that `flai accept --json` prints the ADR-0140 refusal's findings as data.
- `docs/users/flai.md` § Edit a document through flai, and § Accept and release for the refusal. `docs/users/flaiover.md` § Editing documents and § Reviewing a story, for fixing a refused acceptance from the review page.
- Regenerate `docs/users/flai-reference.md` with `make flai-reference`.
- Waits for T-1520, the last of the behaviour it describes.

## Done when

- The ADR is accepted and listed in `design/adrs/README.md`, and ADR-0023 is linked from it as refined.
- The sections above describe the behaviour as built. `flai check --strict` and the markdown lint pass on them.

## Notes

- `design/adrs` is a folder touch because the ADR's file name is not known until `flai adr new` numbers it.
