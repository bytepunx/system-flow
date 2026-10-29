---
id: T-2315
type: task
nature: research
title: The orchestrator in this repository runs a day on reason-primed packs, and the narrative records the context tokens per action from flai stats beside S-0536's figure
status: cancelled
parent: S-0537
owner: alex
created: 2026-10-28T22:08:20Z
updated: 2026-10-29T16:27:59Z
transitions:
  - to: cancelled
    at: 2026-10-29T16:27:59Z
    by: agent-S-0537
stream: S-0537
tags: [orchestrator, measurement]
after: [T-2313, T-2314]
---
# T-2315 The orchestrator in this repository runs a day on reason-primed packs, and the narrative records the context tokens per action from flai stats beside S-0536's figure

## Work

Criterion 5, and criterion 2 measured on this repository. Once the story's change runs in the orchestrator here, through a flai built from the story's branch or the release that carries it as the operator chooses, let it run a day on reason-primed packs. Then:

- Measure `flai prime --role orchestrate --reason thread:TH-nnnn --json` against `flai prime --role orchestrate --json` on this repository, for a thread awaiting the operator, and record both sizes.
- Read the orchestrator's context tokens per action from `flai stats` over the day, and the figure S-0536's narrative recorded for the day before, and record both in the narrative's `## Decisions` with the window each covers.
- Record the reasons the runs primed for and any reason a run was refused, with its message.

Running a flai other than the installed release on the host is the operator's to decide: ask on a thread on the story first, recommending a build from the branch under `bin/`.

Waits for T-2313 and T-2314: the prompt must prime with the reasons, and the docs must say so, before the day's run is the story's behaviour.

## Done when

- The narrative records the two pack sizes, the one-thread pack under a quarter of the general pack, and the two per-action figures with their windows.
- Any reason the day refused is either fixed in the story or recorded as an issue.

## Notes
- 2026-10-29T16:27:59Z: moved to cancelled: TH-0719: the operator confirmed moving the day of runs on reason-primed packs to S-0536's follow-up story under E-0030
