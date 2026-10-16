---
id: T-1742
type: task
nature: research
title: Measure the prose flai's agents write and set the readability baseline
status: cancelled
parent: S-0413
owner: alex
created: 2026-10-16T05:45:23Z
updated: 2026-10-16T06:01:27Z
transitions:
  - to: cancelled
    at: 2026-10-16T06:01:27Z
    by: agent-S-0413
stream: S-0413
tags: [prose, research]
touches: [design/system/agent-prose.md]
after: [T-1750]
---
# T-1742 Measure the prose flai's agents write and set the readability baseline

## Work

Add to `design/system/agent-prose.md`, which T-1750 creates, a section that measures how readable the prose is, beside T-1750's volumes.

- Take a sample of what agents write and the operator reads: narratives under `wip/agents/` and `wip/archive/agents/`, thread entries under `wip/threads/`, story notes, `wip/agents/planner.md`, and final reports. Name the files sampled.
- Measure each kind: words per sentence and one readability score (Flesch reading ease or similar). Run the tools in a scratch folder, never installed on the host.
- Name the recurring kinds of verbosity with a short quoted example of each, such as restated context, hedging, long qualifier chains, or repeated IDs.
- State the target style by link to the rules T-1750 lists: `design/conventions/communication.md`, the Language section of `design/system/documentation-standard.md`, and flai's markdown lint of `wip/` (ADR-0061).

Waits for T-1750: it adds to the document T-1750 creates. The later tasks score linters and models against this sample and these rules.

## Done when

- `design/system/agent-prose.md` has a section giving the sample, its readability figures in a table, and the kinds of verbosity with examples.
- The scratch scripts and the sample list are named in the section, so the measurement can be run again.

## Notes
- 2026-10-16T06:01:27Z: moved to cancelled: Duplicates agent-S-0413's plan (T-1750, T-1755, T-1758, T-1760), as the planner recommends on TH-0472
