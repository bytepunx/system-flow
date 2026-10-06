---
id: T-1007
type: task
nature: improvement
title: delegation.md and the template's hooks say how to wait for a sub-agent, in a template release
status: done
parent: S-0285
owner: alex
created: 2026-10-06T10:38:23Z
updated: 2026-10-06T10:54:18Z
transitions:
  - to: ready
    at: 2026-10-06T10:51:05Z
    by: agent-S-0285
  - to: in-progress
    at: 2026-10-06T10:51:06Z
    by: agent-S-0285
  - to: done
    at: 2026-10-06T10:54:18Z
    by: agent-S-0285
stream: S-0285
tags: []
touches: [design/conventions/delegation.md, template/root/design/conventions/delegation.md, template/root/.claude/settings.json, template/CHANGELOG.md, template/template.yaml]
after: [T-1004, T-1006]
usage:
  source: log
  seconds: 192
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 16
      output: 6109
      cache_read: 828812
      cache_write: 24318
      cost: 0.4402
---
# T-1007 delegation.md and the template's hooks say how to wait for a sub-agent, in a template release

## Work

`delegation.md`, in `template/root/design/conventions/` first and then above the marker in `design/conventions/`, says how a story's agent waits for a sub-agent, in the way the start prompt (T-1006) says it, and that `wait_for_events` is for a thread awaiting the designer. The template's `.claude/settings.json` runs `flai guard` on `SubagentStart` and `SubagentStop` as well, passing on nothing from them. The template is released: `template/template.yaml`'s version bumped and a `template/CHANGELOG.md` entry. Waits for T-1004, whose guard the hooks call, and T-1006, whose wording the convention repeats.

## Done when

- [ ] both copies of `delegation.md` say how to wait and when `wait_for_events` is right, matching the prompt
- [ ] the template's `.claude/settings.json` has the two hooks
- [ ] the template's version and changelog record the release

## Notes
