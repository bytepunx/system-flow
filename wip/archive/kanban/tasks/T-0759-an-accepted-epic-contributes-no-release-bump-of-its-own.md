---
id: T-0759
type: task
nature: improvement
title: An accepted epic contributes no release bump of its own
status: done
parent: S-0200
owner: alex
created: 2026-10-03T17:53:22Z
updated: 2026-10-03T17:56:54Z
transitions:
  - to: ready
    at: 2026-10-03T17:53:26Z
    by: agent-S-0200
  - to: in-progress
    at: 2026-10-03T17:53:26Z
    by: agent-S-0200
  - to: done
    at: 2026-10-03T17:56:54Z
    by: agent-S-0200
stream: S-0200
tags: []
touches: [flai/internal/release, flai/cmd/release.go, template/root/design/conventions/git.md, design/conventions/git.md, template/template.yaml, template/CHANGELOG.md, design/tech/ci.md, docs/users/flai.md, docs/users/flai-reference.md, design/adrs]
after: [T-0752]
usage:
  source: log
  seconds: 208
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 89
      output: 30426
      cache_read: 6720412
      cache_write: 150979
      cost: 2.9662
---
# T-0759 An accepted epic contributes no release bump of its own

## Work

The designer answered TH-0082 with option 2: an epic accepted with its last story, or by hand, contributes no bump of its own; its stories carry theirs. `release.LevelFor` returns none for an epic. `git.md`'s release rule changes in the template first, then here above the marker, with a template patch release and its changelog entry. An ADR records it, refining ADR-0076. `design/tech/ci.md`, `flai release`'s help, and the user guide stop saying an epic is major.

It waits for T-0752, whose documents it corrects.

## Done when

- [x] `release.LevelFor` gives an epic no bump, and the release tests say so
- [x] git.md (template and here), the ADR, ci.md, the help, the reference, and the user guide agree

## Notes
