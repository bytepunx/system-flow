---
id: S-0171
type: story
nature: improvement
title: "Add a \"New {item}\" button to the top of the item pages"
status: done
parent: E-0013
owner: alex
created: 2026-09-30T00:22:14Z
updated: 2026-09-30T01:17:27Z
transitions:
  - to: ready
    at: 2026-09-30T00:22:24Z
    by: alex
  - to: in-progress
    at: 2026-09-30T01:11:45Z
    by: agent-S-0171
  - to: review
    at: 2026-09-30T01:17:10Z
    by: agent-S-0171
  - to: done
    at: 2026-09-30T01:17:27Z
    by: alex
tags: [dashboard]
topics: [client-side]
touches: [flaiover/src, design/system/flaiover-dashboard.md, docs/users/flaiover.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 344
  models:
    - model: claude-opus-5-5
      input: 104
      output: 23278
      cache_read: 4947977
      cache_write: 126964
      cost: 2.4713
---
# S-0171 Add a "New {item}" button to the top of the item pages

## Goal

At the top of the item page, display a "New {item}" button to allow the operator to make another of the same type of work item they're currently viewing.

This prevents the operator from having to go back to the board to create a new item.

## Acceptance criteria
- [x] The button is displayed at the top of the page
- [x] The item type named in the button matches the item type being viewed
- [x] The item type being created matches the type that was being viewed

## Tasks
- T-0606 The item page offers New story or New epic, starting the form with the same type
- T-0607 The design and the user guide say an item page offers a New link of its type

## Notes

- Stories and epics only: a task's page has no New link, because the dashboard makes only epics and stories and tasks are the agent's to write (`design/system/flaiover-dashboard.md`, `/items/:id`).
- A story's **New story** starts the form under the same epic (`?parent=`); the form drops an epic that has closed since. The link shows whatever the viewed item's state, when the dashboard can write.
- Verified by behavior tests in `flaiover/src/routes/items/[id]/item.svelte.test.ts` and `flaiover/src/lib/components/NewItemForm.svelte.test.ts`, with `scripts/flaiover-test.sh` passing. Not looked at in a running dashboard, and the Playwright tests were not run.
