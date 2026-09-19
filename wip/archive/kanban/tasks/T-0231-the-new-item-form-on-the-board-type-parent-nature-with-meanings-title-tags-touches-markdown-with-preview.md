---
id: T-0231
type: task
nature: feature
title: "The new item form on the board: type, parent, nature with meanings, title, tags, touches, markdown with preview"
status: done
parent: S-0059
owner: alex
created: 2026-09-19T09:21:36Z
updated: 2026-09-19T09:30:18Z
transitions:
  - to: ready
    at: 2026-09-19T09:27:40Z
    by: system-flow
  - to: in-progress
    at: 2026-09-19T09:27:40Z
    by: system-flow
  - to: done
    at: 2026-09-19T09:30:18Z
    by: system-flow
stream: S-0059
tags: []
---

# T-0231 The new item form on the board: type, parent, nature with meanings, title, tags, touches, markdown with preview

## Work
A "new" action on the board, hidden on a read-only dashboard, opening `/new`: type (epic or story); for a story the parent from the open epics; the nature from the five natures, each with what it means; a required title; optional tags and touches; and the markdown, starting from the template body, with the explorer's preview beside it. No front matter is shown. A refusal shows the findings with the text kept. After creating, the operator lands on the item's page. Component tests for the form.

## Done when
- The form is tested as a component: required title, parent only for a story, the body from the template, a refusal keeps the text, success navigates
- `make flaiover-test` and `make flaiover-build` pass

## Notes
