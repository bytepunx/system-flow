---
id: T-2354
type: task
nature: improvement
title: "An ADR refining ADR-0207 § 4 settles a role file's start prompt as data and the built-in roles' triggers read from their files"
status: done
parent: S-0554
owner: alex
created: 2026-10-30T20:48:31Z
updated: 2026-11-02T20:39:07Z
transitions:
  - to: ready
    at: 2026-11-02T20:37:02Z
    by: agent-S-0554
  - to: in-progress
    at: 2026-11-02T20:37:03Z
    by: agent-S-0554
  - to: done
    at: 2026-11-02T20:39:07Z
    by: agent-S-0554
stream: S-0554
tags: [cli]
touches: [design/adrs/README.md, design/adrs/0250-a-role-file-s-start-prompt-is-a-text-template-in-its-front-matter-s-start-field.md]
usage:
  source: log
  seconds: 124
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 20
      output: 7571
      cache_read: 1616077
      cache_write: 32418
      cost: 0.6672
---
# T-2354 An ADR refining ADR-0207 § 4 settles a role file's start prompt as data and the built-in roles' triggers read from their files

## Work

Write an ADR refining ADR-0207 § 3 and § 4. Decide where the start prompt lives in a role file, apart from the body, which stays the definition's prompt: a front matter field or a fenced section, with the reason. Name the template syntax and the run fields it may use: the item and its kind, who asked (`AskedOn`), the reasons and their prime forms, the merge queue, the focus, the window since the last run, a defined role's inputs, the root, and the run's name. Say how a project's additions reach the start prompt, and how `flai upgrade` merges it in a baseline file. Say that the planner's, the analyzer's, and the curator's triggers are read from their files through the launchers S-0418 built, how an item-taking schedule (`each-ready-story`) and the three trigger conditions are evaluated, and that the orchestrator's reasons-and-end run (ADR-0209) stays. Add it to `design/adrs/README.md`.

First layer: it waits for nothing. Every other task waits for it.

## Done when

The ADR is accepted, `design/adrs/README.md` lists it, and `flai check --strict` reports nothing.

## Notes

Drafted by the planner of E-0024 on 2026-10-30. TH-0760 confirmed that the body stays the definition (ADR-0207 § 4) and that moving the start prompt is a new ADR, not a rewrite of the body. The ADR file's name is not known until it is written, so the story keeps `design/adrs` as a folder touch.
