---
id: T-0912
type: task
nature: feature
title: The orchestrator's prompt and definition tell it to recommend, answer, or escalate each thread awaiting the operator, citing its source
status: backlog
parent: S-0220
owner: alex
created: 2026-10-05T04:48:21Z
updated: 2026-10-05T04:48:21Z
transitions: []
stream: S-0220
tags: [flai]
touches: [flai/internal/harness/harness.go, flai/internal/harness/harness_test.go, ".claude/agents/orchestrator.md", template/root/.claude/agents/orchestrator.md, design/system/strategic-agents.md, design/conventions/strategic-agents.md, template/root/design/conventions/strategic-agents.md, template/CHANGELOG.md]
after: [T-0902]
---
# T-0912 The orchestrator's prompt and definition tell it to recommend, answer, or escalate each thread awaiting the operator, citing its source

## Work

S-0218 writes the orchestrator's prompt in `flai/internal/harness/harness.go` and its definition in `.claude/agents/orchestrator.md`, with a copy in the template. Add what it does with threads, by `answer_threads`, which the prompt names:

- Each turn, take the threads awaiting the operator from `inbox`, leaving out those it opened and those with a recommendation already pending.
- With `recommend`: reply to each with `thread_reply`, `recommendation` set, and a `source`, the ADR, design section, or convention its answer rests on, read with `doc_get` first.
- With `autonomous`: answer with a `source` when one settles the question. Post a recommendation instead, escalating to the operator, when no source settles it, or when the question asks for the operator's judgement: a decision not yet recorded, a change of scope, or money, such as a cost of delay input, an estimate, or spend.
- Never resolve a thread it did not open, and never answer one it opened.

Add the same rules to `## As the orchestrator` in `design/conventions/strategic-agents.md` and its template copy, with a `template/CHANGELOG.md` entry. Describe answering threads under `## The orchestrator and the analyzer` in `design/system/strategic-agents.md`: the two permissions, the escalation rule, the guard's part (T-0910), and the decision log entry (T-0902).

It waits for T-0902, whose `recommendation` and `source` parameters the prompt names. It runs with T-0908 and T-0910, whose paths it does not share.

## Done when

- a harness test finds the thread rules in the orchestrator's prompt for `recommend` and for `autonomous`, and none when `answer_threads` is unset
- the definition and its template copy match, and the convention and its template copy carry the same rules
- `strategic-agents.md` describes answering threads
- `go test ./internal/harness/` passes and `flai check --strict` is clean on the convention

## Notes
