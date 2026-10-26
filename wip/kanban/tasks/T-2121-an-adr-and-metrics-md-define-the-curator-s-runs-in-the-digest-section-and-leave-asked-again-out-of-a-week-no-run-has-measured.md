---
id: T-2121
type: task
nature: feature
title: An ADR and metrics.md define the curator's runs in the digest section and leave asked_again out of a week no run has measured
status: done
parent: S-0504
owner: alex
created: 2026-10-26T04:09:13Z
updated: 2026-10-26T04:13:56Z
transitions:
  - to: ready
    at: 2026-10-26T04:12:57Z
    by: agent-S-0504
  - to: in-progress
    at: 2026-10-26T04:12:58Z
    by: agent-S-0504
  - to: done
    at: 2026-10-26T04:13:56Z
    by: agent-S-0504
stream: S-0504
tags: [digest, metrics, adr]
touches: [design/adrs/README.md, design/system/metrics.md, design/adrs/0204-flai-stats-reports-the-curator-s-runs-in-the-digest-section-and-leaves-asked.md]
after: [T-2120]
usage:
  source: log
  seconds: 58
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 11
      output: 3847
      cache_read: 556940
      cache_write: 23936
      cost: 0.3368
---
# T-2121 An ADR and metrics.md define the curator's runs in the digest section and leave asked_again out of a week no run has measured

## Work

Write the ADR with `flai adr new`, refining ADR-0200, and accept it with `flai adr accept` once the design is written. It decides, in one decision sentence: `flai stats` reads the curator's activity document, `wip/agents/curator.md`, into the `digest` section, so that a repeat count of 0 can be told from a curator that never ran. State:

- `digest.weeks[].curator_runs`: the curator's activities whose time falls in the week, from the activity document's entries, 0 in a week with none, and 0 in every week while the document does not exist.
- `digest.weeks[].asked_again` is absent in a week in which no curator run has ended since the week's start: the newest activity's time is before the week's `start`, or there is no document. Present, as today, otherwise. A week whose runs ended before some of its threads were opened still reads low, as ADR-0200's consequence says.
- `digest.curator`: `runs`, the activities in the window, and `last_run`, the newest activity's time over the whole document, absent while it has none. Absent while the document does not exist.
- The text `flai stats` prints: `asked again –` in a week with no `asked_again`, and `curator runs n` on each week's line, and a `curator` line, `curator  runs 3 · last run 2026-10-25T20:11:03Z`, or `curator  never ran` when `digest.curator` is absent.

Then change `design/system/metrics.md` § Digest to match: the Source table gains the curator's activity document, the Value tables gain the three values with their precision, the print format gains the new text, and the paragraph after it says that `asked_again` is left out while no run could have counted it. Bump `updated`. Add the ADR to the ADR index as `flai adr new` does, and link it from the section as ADR-0200 is linked.

This waits on T-2120 because the ADR records the remedy the operator agreed to. It names what the remaining tasks build, so it comes before the code.

## Done when

- A new ADR under `design/adrs/`, accepted, refines ADR-0200 and states the three values and the absent `asked_again` as above, and `design/adrs/README.md` lists it.
- `design/system/metrics.md` § Digest defines `digest.weeks[].curator_runs`, `digest.curator.runs`, `digest.curator.last_run`, and when `asked_again` is absent, with the printed text, and its `updated` is bumped.
- `flai check --strict` is clean for the story's scope.

## Notes

Drafted by planner-S-0504. The ADR's file is not named in the touches because `flai adr new` allocates its number; the story's `design/adrs` folder touch covers it, as its Planning note records.
