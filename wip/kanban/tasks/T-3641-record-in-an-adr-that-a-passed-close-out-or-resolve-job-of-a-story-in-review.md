---
id: T-3641
type: task
nature: improvement
title: Record in an ADR that a passed close-out or resolve job of a story in review starts the orchestrator, and update the design and conventions
status: backlog
parent: S-0854
owner: alex
created: 2026-11-27T23:59:22Z
updated: 2026-11-27T23:59:22Z
transitions: []
stream: S-0854
tags: [flai, template]
touches: [design/adrs/README.md, design/adrs/0322-a-close-out-or-resolve-job-that-passed-on-a-story-in-review-starts-the.md, design/system/strategic-agents.md, design/workflows/roles/README.md, template/root/design/workflows/roles/README.md, design/conventions/strategic-agents.md, template/root/design/conventions/strategic-agents.md]
---
# T-3641 Record in an ADR that a passed close-out or resolve job of a story in review starts the orchestrator, and update the design and conventions

## Work

This task waits for nothing. It writes the decision that the other tasks build.

Write the ADR with `flai adr new`, refining ADR-0211 and ADR-0209. Its number is the next one free, and the touch `0322-…` is a guess that `flai task done` corrects. It decides:

- `verified` joins the trigger events of a role file.
- A new condition, `passed-in-review`, holds for a `verified` event when three things are true: its verdict is `passed`, its `recorded` field is not `true`, and the story is in review when flai serve looks. The close-out job and the rebase-and-resolve job both end in `announce` in `flai/internal/serve/verifyjob.go` with such an event.
- The orchestrator's role file, and flai serve's built-in triggers, gain `{event: verified, when: passed-in-review}`. A matched event is the reason `review:S-nnnn`, primed as a review's story and its last verify result.
- A verdict answered from the record is no trigger. The orchestrator's own `verify` call answers from the record, and on a story it cannot accept it would start itself again, in a loop. A job that ended `failed`, `stopped`, or `interrupted` is no trigger either: the story's agent or the operator acts on it.
- Coalescing and the dedupe key stay as they are, so one story's event and its move to review in one window give one reason.

Cite I-0413's instance as the context: S-0839 and S-0843 waited six days in review with passing verdicts.

Then update the design to match, linking the ADR:

- `design/system/strategic-agents.md`: the trigger table under § The orchestrator › Starting and stopping it, and the reason-to-form table under § Primed for its reasons (S-0537).
- `design/workflows/roles/README.md` and `template/root/design/workflows/roles/README.md`: the event list and the condition table.
- `design/conventions/strategic-agents.md` and `template/root/design/conventions/strategic-agents.md`, § As the orchestrator › Accepting a story, step 1: say that the job's passing verdict on a story still in review starts the later run.

## Done when

- The ADR is accepted, is listed in `design/adrs/README.md`, and names the event, the condition, the trigger entry, the reason, and what it leaves out.
- The design, the roles README, and the convention, each with its template copy, say the same and link the ADR.
- `flai check --strict` and the markdown lint pass on the changed files.

## Notes
