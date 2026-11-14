---
id: T-3345
type: task
nature: remediation
title: Remove the four branches and worktrees S-0777's test run left in the main repository, once the operator confirms
status: in-progress
parent: S-0781
owner: alex
created: 2026-11-13T22:53:18Z
updated: 2026-11-13T23:00:11Z
transitions:
  - to: ready
    at: 2026-11-13T23:00:10Z
    by: system-flow
  - to: in-progress
    at: 2026-11-13T23:00:11Z
    by: system-flow
stream: S-0781
tags: [git, check]
usage:
  source: log
  seconds: 2671
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 67
      output: 20386
      cache_read: 3397812
      cache_write: 151639
      cost: 2.0484
  strategic:
    - kind: orchestrator
      seconds: 8
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 1
          output: 12
          cache_read: 48333
          cache_write: 613
          cost: 0.0145
---
# T-3345 Remove the four branches and worktrees S-0777's test run left in the main repository, once the operator confirms

## Work

`story/S-0001`, `story/S-0002`, `story/S-0003`, and `story/S-0005` each point at 5c3b47d0c. That is S-0777's commit "test: [S-0777] flai/cmd's kluster fixtures resolve the project root", 2 commits past main's 289bbc15b, and S-0777 is accepted. Their worktrees, under S-0777's removed `bin/` folder, are prunable. They hold no work of S-0001 to S-0005, which are archived. While they exist, `flai check` in the main checkout warns `story.unaccepted` four times, every close-out notes it outside its story, and `make flai-test` fails on this host (I-0370).

Removing a branch is hard to reverse, so ask the operator first on a thread on S-0781, naming the four branches and their commit, unless they confirmed it on the plan's thread already. Once they confirm, run these in the main checkout and write their output in this task's `## Notes`:

```sh
git worktree prune
git branch -D story/S-0001 story/S-0002 story/S-0003 story/S-0005
```

Then run `scripts/flai.sh check` in the main checkout and check that it reports no `story.unaccepted`.

It waits for nothing in this story.

## Done when

- The operator confirmed the removal on a thread.
- `git branch --list 'story/S-000*'` lists nothing, and `git worktree list` lists no worktree for those branches.
- `flai check` in the main checkout reports no `story.unaccepted`.

## Notes
