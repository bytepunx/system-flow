---
id: S-0598
type: story
nature: improvement
title: A checked write runs only the rules an item write can affect, so creating or editing an item takes about 2 s, not 7 to 18 s
status: backlog
parent: E-0033
owner: alex
created: 2026-11-07T22:39:02Z
updated: 2026-11-07T23:48:54Z
transitions: []
tags: [cli]
touches: [flai/internal/itemnew/itemnew.go, flai/internal/itemedit/itemedit.go, flai/internal/check/check.go, flai/internal/check/itemwrite.go, flai/internal/check/itemwrite_test.go, flai/internal/itemnew/itemnew_test.go, flai/internal/itemedit/rules_test.go, design/system/server-performance.md, design/system/flai-cli.md, design/system/strategic-agents.md, docs/users/flai.md, design/adrs/0264-a-checked-write-of-a-work-item-runs-only-the-rules-the-write-can-affect.md, design/adrs/README.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
finalized:
  by: alex
  at: 2026-11-07T23:48:54Z
---
# S-0598 A checked write runs only the rules an item write can affect, so creating or editing an item takes about 2 s, not 7 to 18 s

## Goal

Cause 8 in [server-performance.md](../../../design/system/server-performance.md#where-the-time-goes): a checked creation or edit runs the whole of `flai check` twice, before and after. Each run lists all 3,233 items cold, lints all 5,068 markdown files under `wip`, starts `node` for the component examples, and scans the digest for secrets. `story new`, `epic new`, `task new` (and MCP `item_new`), `edit`, finalize, `criteria tick`, and `issue story` take 7 to 18 s wall, 8 to 12.7 s CPU, about 97% of it in the two checks. The serve log shows `item.new` at p50 7.9 s and up to 32 s.

## Acceptance criteria

- [ ] A checked creation or edit runs only the rules a work item write can affect: the written items and their parent, the after, overlap, and touches graph, the board, the story's narrative, and the markdown lint and conflict markers of the written files; and it computes the "before" result for those rules alone.
- [ ] It still refuses exactly the findings ADR-0197 and ADR-0213 refuse, which tests show for a creation and an edit that bring a finding of each kept rule.
- [ ] `flai story new`, `flai edit`, and MCP `item_new` of a task, timed in a scratch copy of this repository, each take under 3 s, and server-performance.md records the before and after.

## Tasks
- T-2620 check runs an item write's rule set: the written items and their parent, the graph rules, the board, the story's narrative, and the written files' lint and conflict markers
- T-2622 itemnew.Create checks a creation with the item-write rule set before and after, in place of check.Run
- T-2625 itemedit.Apply checks an edit with the item-write rule set before and after, in place of check.Run
- T-2635 Time story new, edit, and MCP item_new of a task in a scratch copy, before and after, and record them in server-performance.md
- T-2640 Write the ADR that a checked write runs only the rules an item write can affect, and say so in flai-cli.md, strategic-agents.md, and the user guide

## Notes

- Proposed fix: a rule set for item writes in `check`, used by `itemnew.Create` and `itemedit.Apply`, in place of `check.Run`. The whole check stays in `flai check`, the close-out, and CI.
- If the scope of what a checked write refuses changes, that refines ADR-0197 and needs an ADR.
