---
id: T-1269
type: task
nature: experiment
title: Write S-0241's results document under design/experiments with a recommendation
status: backlog
parent: S-0241
owner: alex
created: 2026-10-07T23:13:00Z
updated: 2026-10-07T23:13:00Z
transitions: []
stream: S-0241
tags: [template]
touches: [design/experiments/S-0241-measure-planned-stories-again-with-forked-task-sub-agents-once-a-headless-session-offers-forks.md]
after: [T-1268]
---
# T-1269 Write S-0241's results document under design/experiments with a recommendation

## Work

Copy `design/experiments/template.md` to `design/experiments/S-0241-measure-planned-stories-again-with-forked-task-sub-agents-once-a-headless-session-offers-forks.md` and write every section. The hypothesis: a forked task sub-agent starts from its parent's prompt cache, so it reads less than a fresh `general-purpose` one, and a planned story costs less. The success measure: the forked replays' sub-agent cache reads fall by at least a third against the general-purpose replays of the same stories, with each run's cost no higher and no more review defects. Link the figures in `design/system/agent-context.md` rather than copying the table. Recommend adopt, adapt, or drop forked task sub-agents, and name the story that would carry an adoption into `delegation.md` and the `claude-code` prompt.

Waits for T-1268: the results and the recommendation come from its runs.

## Done when

- The results document exists with every section of the template written, `story: S-0241` in its front matter, and a recommendation of adopt, adapt, or drop with its reason.
- `flai check --strict` and the markdown lint pass on it.

## Notes

The success measure's threshold is the planner's assumption, stated on the plan's thread; the story's agent may change it there before the replays run.
