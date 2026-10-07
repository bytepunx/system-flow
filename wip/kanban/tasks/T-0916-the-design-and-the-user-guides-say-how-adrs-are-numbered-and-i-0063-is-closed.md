---
id: T-0916
type: task
nature: remediation
title: The design and the user guides say how ADRs are numbered, and I-0063 is closed
status: in-progress
parent: S-0245
owner: alex
created: 2026-10-05T05:44:20Z
updated: 2026-10-07T02:10:00Z
transitions:
  - to: ready
    at: 2026-10-07T02:10:00Z
    by: agent-S-0245
  - to: in-progress
    at: 2026-10-07T02:10:00Z
    by: agent-S-0245
stream: S-0245
tags: [docs, adr]
touches: [design/system/flai-cli.md, design/system/documentation-standard.md, docs/users/flai.md, docs/users/flaiover.md, design/issues/I-0063-flai-adr-new-numbers-from-the-story-s-worktree-only-so-parallel-story-branches-take-the-same-adr-number.md, design/issues/summary.md]
after: [T-0915]
---
# T-0916 The design and the user guides say how ADRs are numbered, and I-0063 is closed

## Work

Each of these places says the ADR number is "one more than the highest file present". Change each to say what T-0915 built: one past the highest ADR in the main checkout, every linked worktree, the main branch, and every story branch, read through `storygit.FolderNames`.

- `design/system/flai-cli.md`, in the `flai adr new` row. Write it as the `flai issue` row describes issue numbering, naming S-0245 and I-0063.
- `design/system/documentation-standard.md`, in the rule on how ADRs are numbered.
- `docs/users/flai.md`, in its `flai adr new` paragraph.
- `docs/users/flaiover.md`, where the "+ new ADR" form says how the number is given.

Then close I-0063 with `flai issue close I-0063 --reason`, naming S-0245 and what fixed it. Run it in the story's worktree, so that `design/issues/summary.md` is regenerated on the branch.

It waits for T-0915. The documents describe what that task built, and the issue closes only once its reproduction test passes.

## Done when

- The four documents say how the next ADR number is found, and the markdown lint passes on them.
- I-0063 is closed with a reason naming S-0245, and `design/issues/summary.md` no longer lists it.
- `flai check --strict` is clean.

## Notes
