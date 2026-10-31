---
id: T-2244
type: task
nature: improvement
title: The design, the convention and its template copy, the orchestrator's prompt, the permission's settings text, and the guides say the orchestrator sets cost of delay inputs the operator decided on an epic or an open story
status: done
parent: S-0523
owner: alex
created: 2026-10-28T14:32:46Z
updated: 2026-10-30T17:44:31Z
transitions:
  - to: ready
    at: 2026-10-30T17:33:02Z
    by: agent-S-0523
  - to: in-progress
    at: 2026-10-30T17:33:04Z
    by: agent-S-0523
  - to: done
    at: 2026-10-30T17:44:31Z
    by: agent-S-0523
stream: S-0523
tags: [docs, conventions, orchestrator, cost-of-delay]
touches: [design/system/strategic-agents.md, design/system/flai-cli.md, design/system/project-manifest.md, design/conventions/strategic-agents.md, template/root/design/conventions/strategic-agents.md, design/workflows/roles/orchestrator.md, template/root/design/workflows/roles/orchestrator.md, template/CHANGELOG.md, flai/internal/harness/harness.go, flai/internal/harness/harness_test.go, docs/users/flai.md, docs/users/flaiover.md, docs/operators/settings.md, docs/operators/index.md]
after: [T-2344]
usage:
  source: log
  seconds: 687
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 107
      output: 29149
      cache_read: 8382693
      cache_write: 365759
      cost: 4.4139
---
# T-2244 The design, the convention and its template copy, the orchestrator's prompt, the permission's settings text, and the guides say the orchestrator sets cost of delay inputs the operator decided on an epic or an open story

## Work

Fourth layer: waits for T-2344, whose code and settings text it describes. Rewritten by agent-S-0523 after TH-0764: the operator chose a permission of its own, `set_cost_of_delay` (ADR-0224), so wherever the planner's text below says "under `plan_backlog_stories`" for an epic or an open story, read "under `set_cost_of_delay`". `plan_backlog_stories` keeps a backlog story's first inputs. Every document gains the new permission's row beside `plan_backlog_stories`'s. The settings text and `docs/users/flai-reference.md` are done in T-2344. The template's changelog entry goes under a new `1.1.12` heading, since 1.1.11 is published.

Design, `design/system/strategic-agents.md`: in What it does with each permission (S-0219), the `plan_backlog_stories` row's "It writes" and "flai also refuses" cells, and the paragraph below the table that says flai lets the orchestrator add inputs only to a backlog story when neither it nor its epic has any; in A story planner's threads (S-0328), say that a figure the operator decides on any thread, an answer or a confirmed recommendation, is the orchestrator's to set on the epic or story it concerns; in Its permissions and the guard, the row that says the guard admits a story alone. Link the new ADR. `design/system/flai-cli.md`: the `flai edit` row's clause on `plan_backlog_stories` ("a story in the backlog, whose epic has none either, its first inputs") and the `flai mcp` row where it names `item_edit`'s rule. `design/system/project-manifest.md`: the permission's row and its risk.

Convention, `design/conventions/strategic-agents.md`, and the same text in `template/root/design/conventions/strategic-agents.md`: under As the orchestrator, the `plan_backlog_stories` bullet's cost of delay sentence, and the `autonomous` bullet's money exception; under Asking the operator, "the orchestrator's to give a backlog story that has none" becomes: the orchestrator's to set, under `plan_backlog_stories`, on an epic or an open story, for a figure the operator decided on a thread or the planner recommended on a story planner's thread for a backlog story without any; it never sets a value and never changes an input on its own judgement. Add the entry to `template/CHANGELOG.md` under a new patch version, as the entries for S-0498 and S-0504 are written: what changed in the convention and the role file, and that it needs the flai release that carries S-0523, since an older flai still refuses the orchestrator an epic and a started story.

Prompt, in three places that must say the same thing. The orchestrator's role file, `design/workflows/roles/orchestrator.md`, and its identical template copy `template/root/design/workflows/roles/orchestrator.md` (S-0414): rule 3 says "giving the story cost of delay inputs with `item_edit`", and rule 10 says what to do with a cost of delay input a story's planner asks for; say what it may now set and when: on an epic or an open story, a figure the operator decided on a thread, or the planner's recommendation for a backlog story without any. `flai/internal/harness/harness.go` still holds the same prompt text until S-0419 reads it from the role file: change the same two sentences there, and the expectations `flai/internal/harness/harness_test.go` pins at the same lines ("giving its story cost of delay inputs with item_edit", "save the cost of delay inputs a story's planner asks for").

Settings text: in `flai/internal/manifest/settings.go`, the permission's description and risk (the risk reads "sets cost of delay inputs, money decisions, without you"); say it sets them on an epic or an open story the operator decided on. Run `make flai-reference`, which regenerates `docs/users/flai-reference.md` from it, and copy the same words to the permission tables in `docs/operators/settings.md`, `docs/operators/index.md`, and `docs/users/flaiover.md`. In `docs/users/flai.md`, When the orchestrator plans stories, rewrite the paragraph "It sets cost of delay inputs alone ... only on a backlog story when neither the story nor its epic has any", and the `plan_backlog_stories` row of the permissions table.

## Done when

- Each document above says that under `set_cost_of_delay` (ADR-0224) the orchestrator adds or changes the cost of delay inputs of an epic or a story that is not done, cancelled, or archived, inputs only, for a figure the operator decided on a thread. Each also says that `plan_backlog_stories` keeps a backlog story's first inputs, as the planner recommends on its thread.
- `design/workflows/roles/orchestrator.md` and `template/root/design/workflows/roles/orchestrator.md` are identical, and `flai check --strict` passes on them.
- `design/system/strategic-agents.md` links the new ADR, and its `updated` date and the conventions' are bumped.
- `go test ./internal/harness/... ./internal/manifest/...` passes in `flai/`, and `docs/users/flai-reference.md` matches `make flai-reference`'s output.
- `template/CHANGELOG.md` has the entry, and `scripts/flai.sh test` on the changed paths is clean, the markdown lint among its tiers.

## Notes

Drafted by planner-S-0523; the role files were added after S-0414's acceptance put the orchestrator's prompt in `design/workflows/roles/orchestrator.md`. The agent definitions `.claude/agents/orchestrator.md` and `template/root/.claude/agents/orchestrator.md` say "giving the story cost of delay inputs" and are left out: a write under `.claude/` needs the operator's permission prompt, and the sentence stays true. Widen the touches and ask on the plan thread if the story's agent judges them worth the prompt.
