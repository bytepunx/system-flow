---
id: T-2203
type: task
nature: remediation
title: The design and the docs say how a curator run on the usage limit is recorded, retried, and raised, and I-0212 is closed
status: backlog
parent: S-0514
owner: alex
created: 2026-10-28T12:54:40Z
updated: 2026-10-28T12:54:40Z
transitions: []
stream: S-0514
tags: [docs, design, curator]
touches: [design/system/strategic-agents.md, docs/users/flai.md, docs/users/flai-reference.md, docs/operators/index.md, design/issues/I-0212-a-curator-run-that-ends-on-the-harness-s-usage-limit-is-recorded-failed-logs-no-activity-and-is-not-started-again-until-the-next-acceptance-or-schedule.md, design/issues/summary.md]
after: [T-2202]
---
# T-2203 The design and the docs say how a curator run on the usage limit is recorded, retried, and raised, and I-0212 is closed

## Work

Waits for T-2202, the last of the behaviour, so that the documents describe what was built.

In `design/system/strategic-agents.md`, under "The curator": in "The curator's run and how it ended", add the limit to the outcome table (`failed`, with its why naming the limit and the reset time, whether or not the exit was seen), the run's new fields (the reset time, the retry mark, `limit_thread`), that such a run logs no activity, and the thread of T-2202; in "The curator after acceptances and on a schedule", add a row for the retry's trigger and say that triggers are held until the reset and that the state survives a `flai serve` restart; link the ADR of T-2199. Bump `updated`.

In `docs/users/flai.md` "Running the curator", `docs/users/flai-reference.md` `curate`, and `docs/operators/index.md` "The curate host action", say what happens when the harness's usage limit ends a run: recorded `failed` on the limit with the reset time, started again once at the first look after the reset with the acceptances meanwhile, and, when that ends on the limit too or no reset time is read, a thread on `system-flow.yaml` by `flai serve`; replace or qualify the sentences that say a start is not tried again. Keep each to a few sentences.

Then close the issue with `flai issue close I-0212 --reason "..."` run in the story's worktree, the reason saying what fixed it (the usage reader's mark, the retry at the reset, and the thread), which writes the issue's file and `design/issues/summary.md`. Run `flai test` on the paths for the markdown lint, and `flai check --strict`.

## Done when

- `design/system/strategic-agents.md` describes the limit outcome, the retry, the held triggers, and the thread, with the ADR linked.
- The three docs sections say the same for users and operators, and no longer say a run on the limit is not tried again.
- I-0212 is closed with a reason, and `design/issues/summary.md` no longer lists it as open.
- `flai check --strict` and the markdown lint pass.

## Notes

Layer 5 of S-0514's plan, the last. No dashboard page shows the curator's runs yet, so `docs/users/flaiover.md` and flaiover do not change; `agent.status` carries the new fields of the run as data.
