---
id: T-2307
type: task
nature: research
title: The orchestrator in this repository runs under the new trigger for a day, and the narrative records its runs, their reasons, and the context tokens per action
status: backlog
parent: S-0536
owner: alex
created: 2026-10-28T21:46:49Z
updated: 2026-10-28T21:46:49Z
transitions: []
stream: S-0536
tags: [orchestrator, measurement]
after: [T-2306]
---
# T-2307 The orchestrator in this repository runs under the new trigger for a day, and the narrative records its runs, their reasons, and the context tokens per action

## Work

Criterion 6. ADR-0209 is accepted on main already (`docs: ADR-0209 accepted`), so that half of the criterion needs only checking. The other half is a day of the orchestrator in this repository under the new trigger:

- The host's flai serve runs the installed flai, not the story branch's. Before the day starts, ask the operator on a thread how this story's build reaches the host; the plan's thread records the planner's recommendation. Do not start a second flai serve.
- Over one day, read the `orchestrate` journal entries: each run, its reasons, how it ended (`reasons done`, the timeout, or a failure), and the activities it logged.
- Take the orchestrator's cost and activities over the day from `flai stats --json` (`strategic`, kind `orchestrator`), and the context tokens of each run from its session log under the serve folder's agents logs, and divide by the actions it logged. Set the figure beside ADR-0209's 7M per action on 2026-10-28.
- Record the runs, the reasons, and the figure in the narrative's `## Decisions` with `flai stream log`, and say which part comes from `flai stats` and which from the logs.

Waits for T-2306: the day measures the story as built and documented.

## Done when

- The narrative records the day's runs, each with its reasons and how it ended, and the context tokens per action with its sources.
- ADR-0209's status reads `accepted`.
- Criterion 6 is ticked, or, if the day could not be run before review, left unticked with the reason in the story's notes.

## Notes
