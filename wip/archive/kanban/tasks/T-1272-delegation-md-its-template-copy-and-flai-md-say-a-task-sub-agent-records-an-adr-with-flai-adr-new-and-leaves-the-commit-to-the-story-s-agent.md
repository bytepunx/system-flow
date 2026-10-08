---
id: T-1272
type: task
nature: remediation
title: delegation.md, its template copy, and flai.md say a task sub-agent records an ADR with flai adr new and leaves the commit to the story's agent
status: done
parent: S-0287
owner: alex
created: 2026-10-07T23:19:23Z
updated: 2026-10-08T07:22:43Z
transitions:
  - to: ready
    at: 2026-10-08T07:22:19Z
    by: agent-S-0287
  - to: in-progress
    at: 2026-10-08T07:22:19Z
    by: agent-S-0287
  - to: done
    at: 2026-10-08T07:22:43Z
    by: agent-S-0287
stream: S-0287
tags: [conventions, docs, template]
touches: [design/conventions/delegation.md, template/root/design/conventions/delegation.md, docs/users/flai.md, template/CHANGELOG.md]
after: [T-1270, T-1271]
usage:
  source: log
  seconds: 24
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 6
      output: 1858
      cache_read: 371992
      cache_write: 14530
      cost: 0.2119
---
# T-1272 delegation.md, its template copy, and flai.md say a task sub-agent records an ADR with flai adr new and leaves the commit to the story's agent

## Work

In `design/conventions/delegation.md` § As a task sub-agent, and in the same words in `template/root/design/conventions/delegation.md`, say that a task whose work needs an ADR records it with `flai adr new` and sets its topics with `flai adr topics`, without `--commit` or `--autocommit`, and names the ADR in its final message for the story's agent to commit; never copy the template or guess the number. Edit the baseline rule above the marker in both, since it changes in the template in this story.

In `docs/users/flai.md`, in the paragraph that lists what a sub-agent may run under `flai guard`, add the ADR commands and say the commit is refused.

It waits for T-1270, so that the words match the rule as built, and for T-1271, so that both cite the new ADR's number.

## Done when

- Both copies of `delegation.md` carry the same new rule, linking the new ADR, with `updated` bumped.
- `docs/users/flai.md` names what a sub-agent may do with an ADR.
- The markdown lint passes on the three files with `flai test`.

## Notes
