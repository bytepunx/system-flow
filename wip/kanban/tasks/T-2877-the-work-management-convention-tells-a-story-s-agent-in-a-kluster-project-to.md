---
id: T-2877
type: task
nature: feature
title: The work-management convention tells a story's agent in a kluster project to deploy, check, and tick before review
status: done
parent: S-0662
owner: alex
created: 2026-11-10T08:26:49Z
updated: 2026-11-13T13:56:09Z
transitions:
  - to: ready
    at: 2026-11-13T13:44:02Z
    by: agent-S-0662
  - to: in-progress
    at: 2026-11-13T13:44:02Z
    by: agent-S-0662
  - to: done
    at: 2026-11-13T13:56:09Z
    by: agent-S-0662
stream: S-0662
tags: [template, conventions]
touches: [design/conventions/work-management.md, template/root/design/conventions/work-management.md, design/system/workflow.md, template/CHANGELOG.md]
usage:
  source: log
  seconds: 727
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 33
      output: 157
      cache_read: 1606013
      cache_write: 64326
      cost: 0.6244
---
# T-2877 The work-management convention tells a story's agent in a kluster project to deploy, check, and tick before review

## Work

In `design/conventions/work-management.md`, its template copy `template/root/design/conventions/work-management.md`, and `design/system/workflow.md` (Definition of done), say that in a project with a kluster and a `deploy`, a story's agent deploys its story at its branch's head and runs the checks against its namespace before review, with the commands S-0660 and S-0661 added. It ticks each criterion that needs the running system once verified there, and leaves unticked only what it could not verify, with the reason in the story's Notes.

Keep the baseline and the project additions apart as the marker sets them. Add a line under `## Unreleased` in `template/CHANGELOG.md`.

It reads no code of this story, so it waits for none.

## Done when

- The two convention copies say the same above the marker, and `design/system/workflow.md` says it in its definition of done.
- `template/CHANGELOG.md` has the line.
- The markdown lint and `flai check --strict` pass on the four files.

## Notes

- Criterion 6 of S-0662.
