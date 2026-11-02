---
id: T-2439
type: task
nature: feature
title: The design and the user guides say every acceptance commits the shared set, leaves the rest, and sends an unsettled rebase back to the agent
status: done
parent: S-0566
owner: alex
created: 2026-11-02T18:36:36Z
updated: 2026-11-02T19:02:14Z
transitions:
  - to: ready
    at: 2026-11-02T18:36:56Z
    by: system-flow
  - to: in-progress
    at: 2026-11-02T18:36:57Z
    by: system-flow
  - to: done
    at: 2026-11-02T19:02:14Z
    by: agent-S-0566
stream: S-0566
tags: []
touches: [design/system/flai-cli.md, design/system/workflow.md, design/system/strategic-agents.md, docs/users/flai.md, docs/users/flaiover.md]
usage:
  source: log
  seconds: 1517
  estimated: true
  priced_by: estimate
  models:
    - model: claude-fable-5-1
      input: 639
      output: 317
      cache_read: 3113040
      cache_write: 178852
      cost: 4.9921
---
# T-2439 The design and the user guides say every acceptance commits the shared set, leaves the rest, and sends an unsettled rebase back to the agent

## Work

Record ADR-0248 (`design/adrs/0248-*.md`) in the living design and the user guides, as the code will behave once the two Go tasks of S-0566 land; read the ADR whole first, and the task bodies of S-0566 for the names (`shared_uncommitted`, `uncommitted`, `dirty_main`, the pre-commit subject `chore: [S-nnnn] commit the main checkout's wip and design before accepting`, the thread title prefix `Acceptance sent back:` followed by a space).

- `design/system/flai-cli.md`: the `flai accept` row of the commands table (line 85 area) and any paragraph under it that says the acceptance by hand refuses uncommitted files outside wip unless `--yes`, or leaves a stopped rebase for the operator; say the shared set, the rest, the pre-commit, `dirty_main` for every acceptance, and the send back for every acceptance that lands here, linking ADR-0248.
- `design/system/workflow.md` § Transitions and who makes them › What waits on an acceptance (S-0296): the paragraph that says "Uncommitted files on main that the story does not change, a rebase it cannot settle, which goes back to the story's agent, and a thread of its own … no longer wait for the operator (ADR-0189)" describes the orchestrator; add that since ADR-0248 the operator's acceptance by hand and the dashboard's send an unsettled rebase back the same way and commit the shared set, so neither waits for the operator either; and the `dirty_main` eligibility bullet: a file in the shared set holds nothing back.
- `design/system/strategic-agents.md` § Accepting a story (around lines 630 to 640): the table "Three things set the orchestrator's acceptance apart from the operator's by hand" is now two things apart (resuming its own stopped acceptance, and the thread of its own) and two things shared with every acceptance; rewrite it so, and the `dirty_main` row, linking ADR-0248.
- `docs/users/flai.md` § Accept and release (from line 2747): where it says the acceptance refuses uncommitted files outside `wip/` unless `--yes` and leaves a stopped rebase for you to resolve, say what ADR-0248 does and what you see (the pre-commit's subject, the `sent_back` answer, the thread, and that flai serve starts the agent); § When the orchestrator accepts (from line 3005): "Two things the orchestrator's acceptance does without you, where yours by hand does not" becomes what every acceptance does, keeping what stays the orchestrator's alone; § The merge queue: `dirty_main` reads the files outside the shared set. Keep `--yes` documented as the shell's way to include the rest.
- `docs/users/flaiover.md`: the **Accepting** and **What blocks acceptance** paragraphs (lines 178 and 219): the confirmation no longer offers a checkbox for uncommitted files outside `wip/`; it lists the shared set as what the acceptance commits and the rest as left alone; an acceptance whose rebase stops sends the story back to in-progress with a thread, and the page says so.

Bump each file's `updated` to today, keep links relative, and run `flai test` on the paths you changed (the markdown lint). Change no convention baseline and no ADR.

This task waits for no other: it is layer 1, since the ADR fixes the behaviour the code tasks build.

## Done when

- Each of the five documents describes ADR-0248's rule where it described ADR-0189's or the `--yes` refusal, with no sentence left that tells the operator to resolve a rebase, rebase a worktree, or choose files to include.
- `flai test` passes on the five paths.

## Notes

Acceptance criteria 1, 2, and 3 of S-0566, in the documentation.
