---
id: T-3321
type: task
nature: improvement
title: The guard's design, the user guide, and the tooling convention say a work item file is never replaced whole by a Write
status: done
parent: S-0730
owner: alex
created: 2026-11-13T19:54:15Z
updated: 2026-11-13T20:21:46Z
transitions:
  - to: ready
    at: 2026-11-13T20:18:20Z
    by: system-flow
  - to: in-progress
    at: 2026-11-13T20:18:20Z
    by: system-flow
  - to: done
    at: 2026-11-13T20:21:46Z
    by: agent-S-0730
stream: S-0730
tags: [flai, template]
touches: [design/system/agent-adapters.md, docs/users/flai.md, design/conventions/tooling.md, template/root/design/conventions/tooling.md, template/CHANGELOG.md]
after: [T-3320]
usage:
  source: log
  seconds: 204
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 50
      output: 307
      cache_read: 1653956
      cache_write: 74168
      cost: 0.6486
---
# T-3321 The guard's design, the user guide, and the tooling convention say a work item file is never replaced whole by a Write

## Work

Document the rule T-3320 adds, in its final words. That is why this task waits on T-3320.

- `design/system/agent-adapters.md`, under `### The guard`: add the refusal to the list of what the guard refuses. Also say that the hook reads `content` from a `Write` and from OpenCode's `write`.
- `docs/users/flai.md`: in the section that lists what `flai guard` refuses an agent `flai serve` started, add the Write rule and what the agent does instead.
- `design/conventions/tooling.md` and its template copy `template/root/design/conventions/tooling.md`: change the baseline line "Hand-edit item bodies (goals, criteria, notes) freely." Agents may edit a body with `Edit`. A Write that replaces an item file whole must keep the item's front matter, and `flai guard` refuses one that does not. Change both copies in the same words.
- `template/CHANGELOG.md`: an entry for the convention change.
- If the story's agent judges the refusal a decision of its own, as ADR-0102 and ADR-0141 were for earlier guard rules, write an ADR with `flai adr new` and link it from the design section.

## Done when

- The four documents state the rule in the same terms as the refusal T-3320 gives.
- `design/conventions/tooling.md` and `template/root/design/conventions/tooling.md` match each other in the changed line.
- `flai test` on the changed paths passes the markdown lint.

## Notes
