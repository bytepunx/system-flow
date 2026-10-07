---
id: T-1282
type: task
nature: improvement
title: Propose what flai task done commits, in an ADR refining ADR-0107
status: backlog
parent: S-0312
owner: alex
created: 2026-10-07T23:34:31Z
updated: 2026-10-07T23:34:31Z
transitions: []
stream: S-0312
tags: [flai]
touches: [design/adrs/README.md, design/system/flai-cli.md, design/system/workflow.md]
---
# T-1282 Propose what flai task done commits, in an ADR refining ADR-0107

## Work

Propose the fix from I-0104's four instances (S-0212, S-0215, S-0298, S-0328) before building it, as the story's goal asks. Today `taskdone.commit` in `flai/internal/taskdone/taskdone.go` runs `git add -A` and `git commit -m`, so closing one task of a layer commits every other task's uncommitted files under its message and widens its touches with their paths, and `-m` is required even when nothing is left to commit.

The planner's recommended shape, to confirm or replace with a stated reason:

- The commit stages the changed paths the closing task's `touches` cover, and every changed path no other open task of the story (not `done` or `cancelled`) covers, so a file the task added without declaring it still goes with it. `storygit.CommitPaths` in `flai/internal/storygit/commit.go` already stages and commits a path list.
- A changed path covered only by another open task's touches stays uncommitted, and the result lists it (for example `left`), so the agent sees what waits for the other task's close.
- A changed path both the closing task and another open task cover goes with the closing task, as declared.
- `-m` is needed only when there is something to commit; with nothing to commit the close goes on without it.

Write a new ADR with `flai adr new` that refines ADR-0107's commit step, add it to `design/adrs/README.md`, and update the `flai task done` row of `design/system/flai-cli.md` § Commands and the ADR-0107 paragraph of `design/system/workflow.md` § Branches and collisions (ADR-0019) to say what is committed. The new ADR's file is added to the touches when this task closes.

It waits for nothing: it is the first layer, and the code tasks wait for it.

## Done when

- A new ADR says which paths `flai task done` commits, what it leaves, what it reports, and when `-m` is needed, and links ADR-0107 and I-0104.
- `design/adrs/README.md`, `design/system/flai-cli.md`, and `design/system/workflow.md` describe the same behaviour.
- `flai check --strict` is clean.

## Notes

Drafted by the planner for S-0312.
