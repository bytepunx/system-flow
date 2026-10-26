---
id: T-2147
type: task
nature: improvement
title: The design, the strategic-agents convention in both copies, the user docs, the command reference, and the changelog describe confirm_after, the operator-alone mark, and flai serve's confirmation by timeout
status: backlog
parent: S-0502
owner: alex
created: 2026-10-26T04:46:19Z
updated: 2026-10-26T04:46:19Z
transitions: []
stream: S-0502
tags: [flai, template]
touches: [design/system/strategic-agents.md, design/system/flai-cli.md, design/conventions/strategic-agents.md, template/root/design/conventions/strategic-agents.md, docs/users/flai.md, docs/users/flaiover.md, docs/users/flai-reference.md, template/CHANGELOG.md]
after: [T-2144, T-2145, T-2146]
---
# T-2147 The design, the strategic-agents convention in both copies, the user docs, the command reference, and the changelog describe confirm_after, the operator-alone mark, and flai serve's confirmation by timeout

## Work

What the earlier tasks built, written where each audience reads it, in the same story as the behaviour (documentation.md):

- `design/system/strategic-agents.md`: § Answering threads (S-0220) gains a paragraph on the mark and on flai serve's confirmation by timeout, with the conditions T-2141's ADR fixes and the log entry; the permissions table in § Its permissions and the guard gains the `confirm_after` row, saying that it allows the orchestrator nothing and that flai serve acts on it; § The orchestrator › What it does with each permission names it. Link the ADR.
- `design/system/flai-cli.md`: `flai thread reply` gains `--operator-only`, in the thread commands.
- `design/conventions/strategic-agents.md`, § As the orchestrator, and `template/root/design/conventions/strategic-agents.md`, the same words: under `answer_threads`, mark a recommendation on money or a host change for the operator alone with `operator_only`; the rule "never confirm a recommendation: the operator does" gains "or flai serve, after `confirm_after`". The two copies stay identical below the baseline marker.
- `docs/users/flai.md`: the orchestrator's permissions section (`confirm_after` and what it never confirms), `flai thread reply --operator-only`, and the line on an agent waiting on a recommendation.
- `docs/users/flaiover.md`: the `confirm_after` row in the permissions table, with its risk, as the other rows have.
- `docs/users/flai-reference.md`: regenerated with `make flai-reference` for the new flag.
- `template/CHANGELOG.md`: an entry for the convention and the agent definition the template ships.

## Done when

- Each document above names `confirm_after`, and the convention's two copies are identical below the baseline marker.
- `flai check --strict` scoped to the story is clean, the markdown lint passes, and the reference's staleness test passes.

## Notes

Layer 4, waits for T-2144 (the metrics field), T-2145 (the prompt's words), and T-2146 (the look's conditions), so that it describes what was built, not what was planned. The closing task waits for it. S-0498's T-2132 rewrites § Answering threads and the convention too; this story waits for S-0498.
