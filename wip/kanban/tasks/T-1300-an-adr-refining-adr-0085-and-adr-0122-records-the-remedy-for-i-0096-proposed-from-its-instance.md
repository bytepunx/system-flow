---
id: T-1300
type: task
nature: improvement
title: An ADR refining ADR-0085 and ADR-0122 records the remedy for I-0096, proposed from its instance
status: done
parent: S-0318
owner: alex
created: 2026-10-08T00:10:28Z
updated: 2026-10-08T04:35:27Z
transitions:
  - to: ready
    at: 2026-10-08T04:34:44Z
    by: agent-S-0318
  - to: in-progress
    at: 2026-10-08T04:34:44Z
    by: agent-S-0318
  - to: done
    at: 2026-10-08T04:35:27Z
    by: agent-S-0318
stream: S-0318
tags: [flai]
touches: [design/adrs]
usage:
  source: log
  seconds: 43
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 12
      output: 113
      cache_read: 786343
      cache_write: 9476
      cost: 0.3549
---
# T-1300 An ADR refining ADR-0085 and ADR-0122 records the remedy for I-0096, proposed from its instance

## Work

Write an ADR with `flai adr new` that records the remedy for I-0096, proposed from its one instance. It refines ADR-0085 and ADR-0122. It waits for no task.

What the instance shows:

- S-0227's close-out, on 2026-10-06, found `markdown.MD038` in `wip/agents/S-0229.md`, the narrative of S-0229, then in progress. The finding was a code span with a trailing space, `` `rule: ` ``, in its `## Decisions`.
- `## Context`, `## Decisions`, and `## Open questions` are written by hand (`tooling.md`), so flai's wip lint guard (ADR-0061), which `flai stream log` and `flai stream state` run, never saw it. The text reached main in S-0227's acceptance commit, e645e6ee.
- S-0229's agent fixed it fifteen minutes later, in c9544681. Its own close-out checks its own narrative, since `ScopeToStory` (`flai/internal/check/scope.go`) counts a story's narrative as inside it.
- S-0227's agent could not fix it, because the narrative is S-0229's. Yet `ScopeToStory` marked it outside, and `recordOutside` (`flai/cmd/check.go`) recorded it in I-0096.

Propose, and decide in the ADR:

1. A check scoped to a story leaves out every `markdown.*` finding on the narrative of another story that is still open, that is, not done or cancelled. It takes back the counts, as ADR-0122 does for `item.archive`. That story's own close-out, which checks its own narrative, finds it.
2. A `markdown.*` finding on any other file outside the story is still a note recorded in an issue, as ADR-0085 says. Such a finding shows a gap in flai's own wip lint, as I-0056, I-0070, and I-0072 did, and it is worth recording.
3. Without `--story`, every finding is still reported, so the main checkout's check and CI still warn.

Name the alternatives you rejected:

- A flai command that writes the hand-edited sections through the lint guard. It changes `tooling.md` and is a story of its own.
- Leaving out every finding on another story's narrative. That also covers `narrative.state`, which S-0323 and S-0325 remediate; leave it to them.

## Done when

- The ADR exists under `design/adrs/`, status `proposed`, with `refines: [ADR-0085, ADR-0122]`, and it states the decision above with I-0096's instance as its context.
- `flai check --strict` passes, and the markdown lint passes on the ADR.

## Notes

Touches the `design/adrs` folder because `flai adr new` allocates the ADR's number when it runs, so no file can be named now.
