---
id: T-2996
type: task
nature: improvement
title: Measure S-0425's and S-0586's packs again, record the parts in agent-context.md, and close I-0261
status: backlog
parent: S-0699
owner: alex
created: 2026-11-10T10:30:50Z
updated: 2026-11-11T09:04:21Z
transitions: []
stream: S-0699
tags: [flai, docs]
touches: [design/system/agent-context.md, design/issues/I-0261-a-story-s-prime-pack-came-in-22-parts-most-of-each-part-the-same-topic-sources.md, design/issues/summary.md]
after: [T-2992, T-3111]
---
# T-2996 Measure S-0425's and S-0586's packs again, record the parts in agent-context.md, and close I-0261

## Work

This task waits for every other task of S-0699, because it measures what they changed. The planner sets its `after` to each of them.

Build flai from the story's branch with `scripts/flai.sh`. Run `flai prime --story S-0425 --json` and `flai prime --story S-0586 --json`, and the same with main's flai. For each pack, record its bytes, its parts, and its topics, conventions, named, briefs, and catalog bytes. The planner measured these on main on 2026-11-10:

| Story | Pack | Topics | Conventions | Named | Briefs |
|-------|------|--------|-------------|-------|--------|
| S-0425 | 840 KB | 12 KB | 106 KB | 362 KB | 342 KB |
| S-0586 | 552 KB | 13 KB | 106 KB | 74 KB | 342 KB |

Under `## Fitting the pack to a budget` › `### Delivered in parts` in `design/system/agent-context.md`, record:

- The before and after table.
- That later parts never repeated the topic sources: since T-0987 each part after the first carries an empty `topics` list. I-0261's first instance mistook part 1's list for one that every part repeats.
- What each task of S-0699 changed.

Bump the document's `updated`. Then run `flai issue close I-0261 --reason "<what fixed it>"`. Name the tasks and the bytes they saved in the reason. Commit `design/issues/summary.md` as flai writes it.

## Done when

- `design/system/agent-context.md` holds the before and after measurement of both packs, from runs you observed.
- I-0261 is closed with a reason naming what fixed it, and `design/issues/summary.md` no longer lists it.
- `flai check --strict` and the markdown lint pass on the changed files.

## Notes
