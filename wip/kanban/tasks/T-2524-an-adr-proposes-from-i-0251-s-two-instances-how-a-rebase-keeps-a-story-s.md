---
id: T-2524
type: task
nature: research
title: "An ADR proposes from I-0251's two instances how a rebase keeps a story's changelog lines under main's top Unreleased, refining ADR-0228 §4"
status: done
parent: S-0570
owner: alex
created: 2026-11-07T03:47:31Z
updated: 2026-11-07T04:18:23Z
transitions:
  - to: ready
    at: 2026-11-07T04:16:29Z
    by: system-flow
  - to: in-progress
    at: 2026-11-07T04:16:42Z
    by: system-flow
  - to: done
    at: 2026-11-07T04:18:23Z
    by: agent-S-0570
stream: S-0570
tags: [adr, git, release]
touches: [design/adrs/README.md]
usage:
  source: log
  seconds: 101
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 8
      output: 35
      cache_read: 1061753
      cache_write: 4441
      cost: 0.3811
---
# T-2524 An ADR proposes from I-0251's two instances how a rebase keeps a story's changelog lines under main's top Unreleased, refining ADR-0228 §4

## Work

First layer: it waits for no task. Every later task builds what this ADR decides.

Write an ADR with `flai adr new` that refines ADR-0228 §4 and ADR-0155, from the two instances in I-0251:

1. A release renamed the top `## Unreleased` while the branch's line sat under it. The union merge kept both sides, so the rebased file had two `## Unreleased` headings, and four released lines came back.
2. A release renamed the section again. The branch's line was replayed into `## 1.1.18`, and union dropped main's second `## 1.1.3` heading, which `checkUndone` reported as `undoes_main`.

The planner's proposal for the ADR to weigh, recommended first:

- **A. Put the lines back after the rebase.** After a clean rebase in `storygit.Rebase`, and after the merge queue batch's rebase, take each `CHANGELOG.md` the branch changes. Write it again as the base branch has it, with the lines the branch added against its old base (`git diff <old base> <old head> -- <path>`) at the end of the base's top `## Unreleased`. Create that section when the base has none. Leave out a line the base already holds. Commit the result on the branch as its own commit, as the issue fold does (ADR-0126). Run it before `checkUndone`, and leave the rewritten changelogs out of what `checkUndone` reports. `merge=union` stays, so no rebase stops on the file.
- **B. A custom merge driver** (`merge=flai-changelog`). This needs `git config` in every clone and in the main checkout's `.git/config`. The trial merge's `git merge-tree` reads it only there.
- **C. Changelog fragments,** one file per story, that `flai release` gathers. This changes git.md's rule and the template.

Name the choice, the step's order relative to `checkUndone`, `fold`, and `fix`, what `flai stream sync` prints, and what the batch does.

## Done when

- An ADR in `design/adrs/`, `status: proposed` or `accepted` as the story's agent judges, with `refines: [ADR-0228, ADR-0155]`, records the choice and the alternatives with reasons.
- `design/adrs/README.md` lists it. `flai check --strict` reports no `adr.index`.

## Notes

- Drafted by the planner. The ADR's own file is under `design/adrs`, which the story keeps as a folder touch, since its name is not known until `flai adr new` writes it.
