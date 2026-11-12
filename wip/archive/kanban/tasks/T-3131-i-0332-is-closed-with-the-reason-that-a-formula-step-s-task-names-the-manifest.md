---
id: T-3131
type: task
nature: feature
title: I-0332 is closed with the reason that a formula step's task names the manifest's flai.command
status: done
parent: S-0681
owner: alex
created: 2026-11-12T00:43:02Z
updated: 2026-11-12T00:47:02Z
transitions:
  - to: ready
    at: 2026-11-12T00:46:57Z
    by: agent-S-0681
  - to: in-progress
    at: 2026-11-12T00:46:57Z
    by: agent-S-0681
  - to: done
    at: 2026-11-12T00:47:02Z
    by: agent-S-0681
stream: S-0681
tags: []
touches: [design/issues/I-0332-a-formula-step-s-bare-flai-command-resolves-to-the-installed-flai-on-path-which.md, design/issues/summary.md]
after: [T-3130]
usage:
  source: log
  seconds: 5
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 2
      output: 502
      cache_read: 195031
      cache_write: 6139
      cost: 0.0982
---
# T-3131 I-0332 is closed with the reason that a formula step's task names the manifest's flai.command

## Work

Close I-0332 with `flai issue close I-0332 --reason ... --commit` from the worktree, saying that `flai formula run` now writes a flai step's task with the manifest's `flai.command`.

Waits for T-3130: the reason names the fix that task makes, so it closes only once that fix is done and tested.

## Done when

I-0332 is closed with that reason, committed on story/S-0681.

## Notes
