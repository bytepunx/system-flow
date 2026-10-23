---
id: D-0091
title: Close-out chain was not gated on the narrative rewrite's exit code, so a commit went out with an empty narrative
kind: cause
sources: [I-0012]
topics: [conventions, template, all]
paths: [design/conventions, template, scripts, design/system/conventions.md, design/system/repository-layout.md, Makefile, design/issues/I-0006-shell-is-zsh.md, design/issues/I-0012-close-out-chain-was-not-gated-on-the-narrative-rewrite-s-exit-code-so-a-commit-went-out-with-an-empty-narrative.md, design/issues/summary.md]
date: 2026-10-01
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0091 Close-out chain was not gated on the narrative rewrite's exit code, so a commit went out with an empty narrative

I-0012 was closed on 2026-10-01.

The cause:

```text
Close-out chain was not gated on the narrative rewrite's exit code, so a commit went out with an empty narrative
```

The fix:

```text
S-0187: the baseline tooling.md says each step of a chain that must succeed is joined with && or runs under set -e, never ; or a pipe that hides its exit code, and what a commit records is checked before the commit; scripts/close-out.sh runs lint, tests, flai check --strict, the narrative check, and the commit under set -eu and stops at the first failure, and work-management.md points to it
```

Source: I-0012.
