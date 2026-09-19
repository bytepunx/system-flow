---
id: I-0006
title: The agent shell is zsh and bash idioms fail silently
class: efficiency
status: open
count: 5
cost: 2m
first_reported: 2026-09-15T18:00:57Z
last_reported: 2026-09-19T07:56:04Z
updated: 2026-09-19T07:56:04Z
---

# I-0006 The agent shell is zsh and bash idioms fail silently

## Description
The environment shell is zsh. `${var^}` is unsupported and unquoted `$VAR` in `for` loops does not word-split, so a loop over task IDs ran once with all IDs as one argument and the rest of the command chain was skipped.

## Instances

### 2026-09-15T18:00:57Z
S-0001 bootstrap: `${a^}` broke a heredoc loop; three docs index files were rewritten.

### 2026-09-15T18:21:56Z
S-0010: `for t in $TASKS` passed one argument; the close-out chain silently stopped and was redone.

### 2026-09-18T18:08:38Z
S-0049 refinement: flags held in an unquoted variable were passed to flai as one argument; six task-new calls failed and were rerun with explicit flags

### 2026-09-19T02:47:12Z
S-0058: a git push with tags held in a shell variable was passed as one refspec and failed; rerun with the tags written out

### 2026-09-19T07:56:04Z
S-0052: a docker run command line held in a shell variable was passed as one word by zsh in the research lab; wrapped it in a script instead

## Remediation
Scripts use `#!/usr/bin/env sh` with POSIX constructs only; loops in ad hoc commands use literal lists. Close when scripts cover the routine flows.
