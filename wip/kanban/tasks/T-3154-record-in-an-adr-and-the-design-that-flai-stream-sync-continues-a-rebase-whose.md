---
id: T-3154
type: task
nature: improvement
title: Record in an ADR and the design that flai stream sync continues a rebase whose stop the agent resolved
status: backlog
parent: S-0726
owner: alex
created: 2026-11-12T15:56:17Z
updated: 2026-11-12T15:56:17Z
transitions: []
stream: S-0726
tags: [flai]
touches: [design/adrs/0291-flai-stream-sync-continues-a-rebase-whose-stop-the-agent-resolved.md, design/adrs/README.md, design/system/flai-cli.md, design/system/workflow.md]
---
# T-3154 Record in an ADR and the design that flai stream sync continues a rebase whose stop the agent resolved

## Work

Write the decision the story's fix rests on. Use `flai adr new`, then accept it. The ADR says the following:

- When `flai stream sync` finds a rebase in progress in the story's worktree whose stop has no unmerged path and no conflict marker in a staged file, it continues that rebase itself. It runs the loop its own rebase runs, `ContinueOverResolved`, and so settles later stops on issue files (ADR-0126), `design/issues/summary.md` (ADR-0098), the ADR index (ADR-0228), and markdown merged by sections (ADR-0175, ADR-0276).
- It finishes with the steps after a clean rebase (ADR-0286, ADR-0257, ADR-0155).
- A rebase whose stop still has unmerged paths or markers is refused as now, with those paths named.
- Every message that tells an agent to run `git rebase --continue` after a sync's stop says instead to `git add` the resolved paths and run `flai stream sync` again.
- `git rebase --abort` stays the way back.

It refines ADR-0069: the agent still never rebases by hand. Say how the undone check (ADR-0155) learns the branch's old base and head from the rebase's state in the git directory (`rebase-merge/onto`, `orig-head`, or what the sync recorded when it started).

Update the stream sync bullet in `design/system/flai-cli.md` § Commands, and `design/system/workflow.md` § Branches and collisions (ADR-0019), to match.

The ADR's number and slug in this task's touches are predicted. `flai task done` narrows them to the file written.

This task waits for nothing. It is the first layer, because the code and the advice follow its decision.

## Done when

- An accepted ADR records the decision above.
- `design/adrs/README.md` lists it.
- `design/system/flai-cli.md` and `design/system/workflow.md` say what the sync now does with a rebase in progress.

## Notes
