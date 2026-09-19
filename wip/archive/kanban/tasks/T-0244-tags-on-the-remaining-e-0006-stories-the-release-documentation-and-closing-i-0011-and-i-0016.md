---
id: T-0244
type: task
nature: feature
title: Tags on the remaining E-0006 stories, the release documentation, and closing I-0011 and I-0016
status: done
parent: S-0047
owner: alex
created: 2026-09-19T10:30:01Z
updated: 2026-09-19T10:35:14Z
transitions:
  - to: ready
    at: 2026-09-19T10:32:32Z
    by: system-flow
  - to: in-progress
    at: 2026-09-19T10:32:33Z
    by: system-flow
  - to: done
    at: 2026-09-19T10:35:14Z
    by: system-flow
stream: S-0047
tags: []
---

# T-0244 Tags on the remaining E-0006 stories, the release documentation, and closing I-0011 and I-0016

## Work
Give the remaining E-0006 stories, and any other open story, tags that say where they deliver. The release paragraph in `docs/users/flai.md`, the `flai release` row in `design/system/flai-cli.md`, and `design/conventions/git.md` if it describes the choice (template baseline first) say the new rule. Close I-0011 and I-0016 against this story on the branch. Tick the criteria.

## Done when
- Open stories have delivery tags
- The documents state the rule
- I-0011 and I-0016 are closed and `flai check --strict` is clean

## Notes
