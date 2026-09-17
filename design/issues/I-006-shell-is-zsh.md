---
id: I-006
title: The agent shell is zsh and bash idioms fail silently
class: efficiency
status: open
count: 2
cost: 3m
first_reported: 2026-09-15T18:00:57Z
last_reported: 2026-09-15T18:21:56Z
updated: 2026-09-15T22:40:33Z
---

# I-006 The agent shell is zsh and bash idioms fail silently

## Description
The environment shell is zsh. `${var^}` is unsupported and unquoted `$VAR` in `for` loops does not word-split, so a loop over task IDs ran once with all IDs as one argument and the rest of the command chain was skipped.

## Instances

### 2026-09-15T18:00:57Z
S-0001 bootstrap: `${a^}` broke a heredoc loop; three docs index files were rewritten.

### 2026-09-15T18:21:56Z
S-0010: `for t in $TASKS` passed one argument; the close-out chain silently stopped and was redone.

## Remediation
Scripts use `#!/usr/bin/env sh` with POSIX constructs only; loops in ad hoc commands use literal lists. Close when scripts cover the routine flows.
