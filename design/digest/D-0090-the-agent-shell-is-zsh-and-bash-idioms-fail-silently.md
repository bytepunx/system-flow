---
id: D-0090
title: The agent shell is zsh and bash idioms fail silently
kind: cause
sources: [I-0006]
topics: [conventions, template, all]
paths: [design/conventions, template, scripts, design/system/conventions.md, design/system/repository-layout.md, Makefile, design/issues/I-0006-shell-is-zsh.md, design/issues/I-0012-close-out-chain-was-not-gated-on-the-narrative-rewrite-s-exit-code-so-a-commit-went-out-with-an-empty-narrative.md, design/issues/summary.md]
date: 2026-10-01
status: absorbed
absorbed_by: "design/conventions/tooling.md#rules"
asked_again: []
cited: []
---

# D-0090 The agent shell is zsh and bash idioms fail silently

I-0006 was closed on 2026-10-01.

The cause:

```text
The environment shell is zsh. `${var^}` is unsupported and unquoted `$VAR` in `for` loops does not word-split, so a loop over task IDs ran once with all IDs as one argument and the rest of the command chain was skipped.
```

The fix:

```text
S-0187: the baseline tooling.md says the host shell may be zsh, so quote every glob, URL, and variable, compare with [ a = b ], write a loop's list out, and put a sequence of more than a few commands in a POSIX sh script under scripts/ (set -eu); scripts/close-out.sh is the routine close-out flow as such a script
```

Source: I-0006.
