---
id: T-2108
type: task
nature: improvement
title: An ADR settles, from I-0194's instances, how task done checks a story whose tree builds the flai that knows its format
status: done
parent: S-0495
owner: alex
created: 2026-10-25T20:32:15Z
updated: 2026-10-25T20:37:33Z
transitions:
  - to: ready
    at: 2026-10-25T20:36:38Z
    by: agent-S-0495
  - to: in-progress
    at: 2026-10-25T20:36:43Z
    by: agent-S-0495
  - to: done
    at: 2026-10-25T20:37:33Z
    by: agent-S-0495
stream: S-0495
tags: [cli, mcp]
touches: [design/adrs/README.md, design/adrs/0201-when-the-manifest-names-flai-command-task-done-s-check-step-runs-that-command-s.md]
usage:
  source: log
  seconds: 50
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 20
      output: 7267
      cache_read: 1060307
      cache_write: 56832
      cost: 0.7139
---
# T-2108 An ADR settles, from I-0194's instances, how task done checks a story whose tree builds the flai that knows its format

## Work

I-0194 records five closes in three stories (S-0476, S-0429, S-0432) where the MCP `task_done` committed, synced, and moved the task, then stopped at its check step. The check ran in the MCP server's flai (1.67.13, then 1.68.1), which did not know a format the worktree carried: the `answers` key of an examples file (ADR-0187, the story's own change), the `doc.updated` exemption of digest entries (S-0428, accepted and unpublished), and the `curate` role (S-0431, accepted and unpublished). The tree's `scripts/flai.sh check --strict --story` passed each time. The check step is `run.check` in `flai/internal/taskdone/taskdone.go`, which calls `check.RunScoped` in-process, so the running binary's rules check the worktree. ADR-0160's close through the installed flai does not reach it: no release had the rules either.

Write the ADR with `flai adr new`, proposing the remedy from the instances and recording the alternatives:

- Recommended: a manifest key beside `flai.minimum`, such as `flai.command`, naming a command relative to the repository root that runs the flai the tree builds (`scripts/flai.sh` here). When it is set, task done's check step runs `<command> check --strict --story <S-nnnn> --json` in the story's worktree and reads its findings in place of the in-process check, and the result says which flai checked. Only a project that builds flai sets it; a project made from the template never meets the cause.
- Run the check with the installed flai at the server's path, extending ADR-0160: the instances rule it out, since no release knew the format.
- Say in the stop that the finding may come from a flai older than the branch's format: the stop and the fallback to the tree's CLI remain.
- Detect flai's own Go module in the tree: ties flai to this repository with no setting.

Settle in the ADR what the step does when the command cannot build or run (recommended: the step stops as a failed check does, naming the command's output, so a branch that does not build closes no task), and whether the change raises `flai/internal/buildinfo/task-done-rules.txt`, since ADR-0180 says a change to what task done checks raises it.

Waits for nothing: it is the first layer. The manifest task and the check-step task wait for it, since it fixes the key's name and the step's rules.

## Done when

The ADR is accepted with `flai adr accept`, carries the topics `cli` and `mcp`, has its row in `design/adrs/README.md`, names the key, the command task done runs, the rule for a command that fails, and the rules-number decision, and `flai check --strict --story S-0495` reports nothing on it.

## Notes

The ADR's file is created by `flai adr new`, so its path is not in the touches; `flai task done` adds it when the task closes (ADR-0144).
