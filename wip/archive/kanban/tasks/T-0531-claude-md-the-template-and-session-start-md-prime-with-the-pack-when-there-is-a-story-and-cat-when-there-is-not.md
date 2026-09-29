---
id: T-0531
type: task
nature: feature
title: CLAUDE.md, the template, and session-start.md prime with the pack when there is a story and --cat when there is not
status: done
parent: S-0148
owner: alex
created: 2026-09-29T05:24:53Z
updated: 2026-09-29T05:26:52Z
transitions:
  - to: ready
    at: 2026-09-29T05:25:01Z
    by: agent-S-0148
  - to: in-progress
    at: 2026-09-29T05:26:14Z
    by: agent-S-0148
  - to: done
    at: 2026-09-29T05:26:52Z
    by: agent-S-0148
stream: S-0148
tags: []
touches: [CLAUDE.md, template/root/CLAUDE.md.tmpl, template/root/design/conventions, design/conventions, template/template.yaml, template/CHANGELOG.md]
---
# T-0531 CLAUDE.md, the template, and session-start.md prime with the pack when there is a story and --cat when there is not

## Work

`CLAUDE.md`, `template/root/CLAUDE.md.tmpl`, the conventions `README.md`, and the baseline `session-start.md`, in the template and here above the marker, say: with a story, prime with `flai prime --story <id>` (or the MCP `prime`); the pack is a brief within a budget, and a briefed document's body is read with `doc_get` and a heading, or `flai doc show --heading`, before changing what it describes, and whole when its brief shows it bears on the story; without a story, `flai prime --cat`. Bump the template to 1.0.22 with a changelog entry.

## Done when

- [x] The four files say the same in the template and here; the baseline copies are identical above the marker.
- [x] `template/template.yaml` is 1.0.22 and `template/CHANGELOG.md` names S-0148.

## Notes
