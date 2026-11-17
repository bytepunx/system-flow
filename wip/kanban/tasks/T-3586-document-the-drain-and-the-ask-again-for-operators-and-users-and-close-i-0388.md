---
id: T-3586
type: task
nature: improvement
title: Document the drain and the ask again for operators and users, and close I-0388
status: backlog
parent: S-0846
owner: alex
created: 2026-11-17T12:33:47Z
updated: 2026-11-17T12:33:47Z
transitions: []
stream: S-0846
tags: [docs]
touches: [docs/operators/index.md, docs/users/flai.md, design/issues/I-0388-flai-serve-s-auto-upgrade-restart-interrupts-the-close-out-jobs-it-was-running.md, design/issues/summary.md]
after: [T-3585]
---
# T-3586 Document the drain and the ask again for operators and users, and close I-0388

## Work

It waits for T-3585, the last code task, so the documents and the issue's reason describe what was built, not what was planned.

- `docs/operators/index.md`: § Close-out jobs (S-0546), its paragraph When `flai serve` restarts, and § Keeping the host's flai current say that an auto-upgrade waits for running close-out and rebase and resolve jobs up to the ADR's bound, takes no new job meanwhile, and that serve asks again once a job a restart interrupted.
- `docs/users/flai.md` § The close-out job: the `interrupted` verdict is now followed by a job flai serve asks again, so the agent need not ask again itself.
- Close the issue: `flai issue close I-0388 --reason "<what fixed it, naming the ADR and S-0846>"`, which updates `design/issues/summary.md`.

## Done when

- Both documents say what the ADR decides, and link it.
- I-0388 is closed with a reason naming the fix, and `flai check --strict` passes.

## Notes
