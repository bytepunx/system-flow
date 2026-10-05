---
id: T-0910
type: task
nature: feature
title: "flai guard holds the orchestrator's thread calls to answer_threads: recommend only, an answer only with a source, never resolving another's thread or answering its own"
status: backlog
parent: S-0220
owner: alex
created: 2026-10-05T04:48:11Z
updated: 2026-10-05T04:48:11Z
transitions: []
stream: S-0220
tags: [flai]
touches: [flai/internal/guard/guard.go, flai/internal/guard/guard_test.go]
after: [T-0902]
---
# T-0910 flai guard holds the orchestrator's thread calls to answer_threads: recommend only, an answer only with a source, never resolving another's thread or answering its own

## Work

S-0218 gives `flai guard` (`flai/internal/guard/guard.go`) the orchestrate role and refuses a call outside `orchestration.permissions`. For the orchestrator's thread calls, refine it by `answer_threads`:

- Unset: `thread_reply` and `thread_resolve` are refused on threads it did not open, naming `answer_threads` as the permission that would allow them.
- `recommend`: `thread_reply` on another's thread is allowed only with `recommendation` set.
- `autonomous`: `thread_reply` on another's thread is allowed as a recommendation, or as an answer with a `source`. An answer without a source is refused, saying to post it as a recommendation. That is the escalation flai enforces; when to escalate a question that names the operator's judgement is the prompt's (T-0912).
- `thread_resolve` is refused on any thread the orchestrator did not open, whatever the permission.
- On a thread the orchestrator opened, `thread_reply` is allowed only as a follow-up, without `recommendation` or `source`: it never answers its own question.

The guard reads the thread's opener from `wip/threads` in the main checkout, through `threads.Get`. Each refusal is logged as S-0218's refusals are.

It waits for T-0902, which adds the `recommendation` and `source` parameters the guard reads. It runs with T-0908 and T-0912, whose paths it does not share.

## Done when

- a table test covers each permission value against a recommendation, a sourced answer, an unsourced answer, and a resolve, on a thread the orchestrator opened and on one it did not
- each refusal names the permission or the parameter that would allow the call
- `go test ./internal/guard/` passes

## Notes
