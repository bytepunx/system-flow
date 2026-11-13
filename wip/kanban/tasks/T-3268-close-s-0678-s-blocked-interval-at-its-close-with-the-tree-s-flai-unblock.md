---
id: T-3268
type: task
nature: remediation
title: Close S-0678's blocked interval at its close with the tree's flai unblock
status: backlog
parent: S-0761
owner: alex
created: 2026-11-13T07:29:30Z
updated: 2026-11-13T07:29:30Z
transitions: []
stream: S-0761
tags: [flai, check]
touches: [wip/archive/kanban/stories/S-0678-agent-status-with-13-projects-at-once-is-not-under-200-ms.md]
---
# T-3268 Close S-0678's blocked interval at its close with the tree's flai unblock

## Work

S-0678 is archived with an open blocked interval, from 2026-11-11T05:26:16Z, and every close-out notes it as `item.blocked-closed` outside its story. Run `scripts/flai.sh unblock S-0678` in the main checkout. The tree's flai has S-0723's fix (T-3142): on a closed item, archived ones included, it closes the interval at the item's close, 2026-11-11T06:15:07Z. Do not use the installed flai if it is older than the release with S-0723: it would close the interval at now and add a day of blocked time to `flai stats`.

Read S-0678's front matter after the run and check that the interval's `to` is its done transition's time. Then run `scripts/flai.sh check` in the main checkout and check it reports no `item.blocked-closed`. Write both in this task's `## Notes`.

It waits for nothing in this story.

## Done when

- S-0678's blocked interval ends at 2026-11-11T06:15:07Z.
- `flai check` in the main checkout reports no `item.blocked-closed`.

## Notes
