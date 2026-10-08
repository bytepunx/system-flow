---
id: T-1330
type: task
nature: improvement
title: Record what flai task done commits, in an ADR refining ADR-0107
status: done
parent: S-0322
owner: alex
created: 2026-10-08T04:35:05Z
updated: 2026-10-08T08:16:41Z
transitions:
  - to: ready
    at: 2026-10-08T08:14:14Z
    by: agent-S-0322
  - to: in-progress
    at: 2026-10-08T08:14:14Z
    by: agent-S-0322
  - to: done
    at: 2026-10-08T08:16:41Z
    by: agent-S-0322
stream: S-0322
tags: [flai]
touches: [design/adrs/README.md, design/system/flai-cli.md, design/system/workflow.md, design/adrs/0128-flai-task-done-commits-the-paths-the-closing-task-covers-and-those-no-other.md]
usage:
  source: log
  seconds: 147
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 24
      output: 7571
      cache_read: 1790199
      cache_write: 35027
      cost: 0.7318
---
# T-1330 Record what flai task done commits, in an ADR refining ADR-0107

## Work

Record the remedy for I-0104 and I-0108 before building it. Today `run.commit` in `flai/internal/taskdone/taskdone.go` runs `git add -A` and `git commit -m`. Closing one task of a layer therefore commits every sibling's uncommitted files under its message and widens its touches with their paths, and `-m` is required even when nothing is left to commit. I-0104 counts four instances (S-0212, S-0215, S-0298, S-0328); I-0108's is S-0274's layer 4, three tasks in one worktree.

The operator confirmed this remedy on TH-0337, when it was planned for S-0312, since cancelled as a duplicate of this story:

- The commit stages the changed paths the closing task's `touches` cover, and every changed path no other open task of the story (not `done` or `cancelled`) covers, so a file the task added without declaring it still goes with it. `storygit.CommitPaths` in `flai/internal/storygit/commit.go` already stages and commits a path list.
- A changed path that only another open task's touches cover stays uncommitted, and the result lists it, so the agent sees what waits for the other task's close.
- A changed path that both the closing task and another open task cover goes with the closing task, as declared.
- `-m` is needed only when there is something to commit; with nothing to commit the close goes on without it. The same holds for the MCP tool `task_done` and the host method `task.done`.

Write a new ADR with `flai adr new` that refines ADR-0107's commit step and links I-0104, I-0108, and TH-0337. Add it to `design/adrs/README.md`. Update the `flai task done` entry of `design/system/flai-cli.md` § Commands, and the `task.done` sentence of its `flai hostapi` entry, which says it runs `flai task done <id> -m <message>`. Update the ADR-0107 paragraph of `design/system/workflow.md` § Branches and collisions (ADR-0019) to say what is committed. The new ADR's file is added to the touches when this task closes.

It waits for nothing: it is the first layer, and the code tasks wait for it.

## Done when

- A new ADR says which paths `flai task done` commits, what it leaves, what it reports, and when `-m` is needed, and links ADR-0107, I-0104, and I-0108.
- `design/adrs/README.md`, `design/system/flai-cli.md`, and `design/system/workflow.md` describe the same behaviour.
- `flai check --strict` is clean.

## Notes

Drafted by the planner for S-0322, from S-0312's cancelled T-1282. The host method `task.done` was added on the planner's second run.
