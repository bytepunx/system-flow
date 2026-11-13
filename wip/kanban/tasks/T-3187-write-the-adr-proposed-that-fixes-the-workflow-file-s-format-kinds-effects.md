---
id: T-3187
type: task
nature: feature
title: Write the ADR, proposed, that fixes the workflow file's format, kinds, effects, rules, timers, and migration
status: backlog
parent: S-0736
owner: alex
created: 2026-11-12T16:05:15Z
updated: 2026-11-12T16:05:15Z
transitions: []
stream: S-0736
tags: [cli]
touches: [design/adrs/README.md]
---
# T-3187 Write the ADR, proposed, that fixes the workflow file's format, kinds, effects, rules, timers, and migration

## Work

Write the ADR with `flai adr new`, status proposed, in ADR-0207's shape: Context, Decision in numbered sections, Consequences, and Alternatives considered. Decide, each with a table where it lists fields or values: the file's path under `design/workflows/` and its front matter; the state kinds and each of today's six states' kind, item types, WIP limit, and order; the transition fields (from, to, who may make it, reason, rules, effects) and today's transitions mapped to them, ADR-0055's moves back and a task's `in-progress` to `done` included; the rule vocabulary, refusing or warning, and the baseline's rules per transition; a timer's two forms and the baseline's lead, cycle, queue, and review time; how renaming or removing a state with items on the board works; and that a project with no file runs the baseline. Name what supersedes or refines ADR-0003 and ADR-0004. Put each choice the operator should make in a list at the end of Context, with your recommendation.

## Done when

The ADR exists, proposed, covering each of the story's first seven criteria, and the index lists it.

## Notes
