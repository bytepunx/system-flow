---
id: T-2527
type: task
nature: improvement
title: Say in the design, the user guide, and the work-management convention that a close-out leaves out narrative.missing
status: done
parent: S-0555
owner: alex
created: 2026-11-07T03:47:42Z
updated: 2026-11-07T04:14:28Z
transitions:
  - to: ready
    at: 2026-11-07T03:59:38Z
    by: system-flow
  - to: in-progress
    at: 2026-11-07T03:59:40Z
    by: system-flow
  - to: done
    at: 2026-11-07T04:14:28Z
    by: agent-S-0555
stream: S-0555
tags: [flai, check, docs]
touches: [design/system/flai-cli.md, design/system/continuous-improvement.md, docs/users/flai.md, design/conventions/work-management.md, template/root/design/conventions/work-management.md, template/CHANGELOG.md]
after: [T-2518]
usage:
  source: log
  seconds: 888
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 20
      output: 119
      cache_read: 2385882
      cache_write: 23940
      cost: 0.8613
---
# T-2527 Say in the design, the user guide, and the work-management convention that a close-out leaves out narrative.missing

## Work

Each place below lists what a close-out's check leaves out. In each, say that a `narrative.missing` outside the story is left out too, and cite T-2518's ADR:

- the `flai check` row of `design/system/flai-cli.md`
- the close-out item in `design/system/continuous-improvement.md`, on the findings a close-out meets outside the story
- the `--story` and `--record-issues` paragraphs of `docs/users/flai.md`
- the close-out rule in `design/conventions/work-management.md`, and the same rule in `template/root/design/conventions/work-management.md`, worded alike. Say there that a story in progress with no narrative is its own agent's or the operator's to fix, and that only `flai stream open` on that story clears it.
- an entry for S-0555 in `template/CHANGELOG.md`, as a patch, naming the ADR, as S-0366's entry does

Put it beside the `narrative.state` text from ADR-0125 where a list has one, so the two narrative findings read alike. The convention's baseline changes, so `template/` changes in the same story, as `CLAUDE.md` says.

It waits for T-2518 for the ADR's number. It runs beside T-2521, since they share no path.

## Done when

- Each place above names `narrative.missing` among what a close-out leaves out, with the ADR linked.
- The convention reads the same in `design/conventions/` and `template/root/design/conventions/`.
- `flai test` on the changed paths passes, the markdown lint included.

## Notes
