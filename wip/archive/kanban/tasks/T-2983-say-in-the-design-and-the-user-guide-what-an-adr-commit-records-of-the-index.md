---
id: T-2983
type: task
nature: remediation
title: Say in the design and the user guide what an ADR commit records of the index
status: done
parent: S-0710
owner: alex
created: 2026-11-10T10:29:30Z
updated: 2026-11-11T14:17:26Z
transitions:
  - to: ready
    at: 2026-11-11T14:16:50Z
    by: agent-S-0710
  - to: in-progress
    at: 2026-11-11T14:16:50Z
    by: agent-S-0710
  - to: done
    at: 2026-11-11T14:17:26Z
    by: agent-S-0710
stream: S-0710
tags: [cli, docs]
touches: [design/system/flai-cli.md, docs/users/flai.md]
after: [T-2980]
usage:
  source: log
  seconds: 36
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 6
      output: 753
      cache_read: 326174
      cache_write: 4634
      cost: 0.1174
---
# T-2983 Say in the design and the user guide what an ADR commit records of the index

## Work

Write down what T-2980 changed. In `design/system/flai-cli.md`, in the `flai adr` row, say that `--autocommit` and `--commit` record the index that HEAD's ADR files, with the files the commit records, render, and that the working tree keeps the index its own files render, so an ADR file left uncommitted never puts its row in a commit (I-0336). In `docs/users/flai.md`, under "Record a decision", say the same in a sentence for users.

This task waits for T-2980, since it describes the behaviour T-2980 builds, in the words T-2980's code settles. It runs beside T-2986, which changes no path of this task.

## Done when

- Both documents say what an ADR commit records of the index, and match T-2980's behaviour.
- The markdown lint passes on both files.

## Notes
