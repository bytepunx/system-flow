---
id: T-0148
type: task
nature: improvement
title: "Conventions: the pulling agent writes the tasks; template baseline first, then this repository"
status: done
parent: S-0049
owner: alex
created: 2026-09-18T18:08:30Z
updated: 2026-09-18T18:14:31Z
transitions:
  - to: ready
    at: 2026-09-18T18:13:18Z
    by: alex
  - to: in-progress
    at: 2026-09-18T18:13:18Z
    by: alex
  - to: done
    at: 2026-09-18T18:14:31Z
    by: alex
stream: S-0049
tags: []
touches: [template, design/conventions]
---

# T-0148 Conventions: the pulling agent writes the tasks; template baseline first, then this repository

## Work
Edit the baseline in `template/root/design/conventions/work-management.md` and `session-start.md` first: definition of ready without tasks; the pull sequence (move to in-progress, open the narrative, read goal, criteria, and notes, write tasks with `## Work` and `## Done when`, work them in order); what the agent does when it cannot write tasks (no invented scope, `flai block --reason`, a thread on the story saying what is missing, move to other work); and "do not write tasks for backlog stories" reworded so tasks are written by whoever starts the story when it is started. Check the template's `CLAUDE.md` for the same wording. Copy the baseline above the marker into `design/conventions/` here, bump `updated`, and record the template change in `template/CHANGELOG.md` and `template/template.yaml` the way earlier stories did.

## Done when
- Template baseline and this repository's copies are identical above the marker
- No convention file says ready needs a task
- The template changelog has the entry

## Notes
The third done-when item was wrong as written: `flai accept` writes `template/CHANGELOG.md` and bumps `template/template.yaml` at acceptance (`flai/internal/release`), so no entry is written by hand. The story declares `template` in `touches` so the release gives it a patch. `CLAUDE.md` and `template/root/CLAUDE.md.tmpl` gained the clause "write its tasks if it has none" in the pull step.
