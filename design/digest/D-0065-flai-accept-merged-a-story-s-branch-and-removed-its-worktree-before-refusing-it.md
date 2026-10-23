---
id: D-0065
title: flai accept merged a story's branch and removed its worktree before refusing it for unticked criteria or open tasks
kind: cause
sources: [I-0021]
topics: [flai, cli, go, flaiover, dashboard, sveltekit, code, all]
paths: []
date: 2026-10-01
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0065 flai accept merged a story's branch and removed its worktree before refusing it for unticked criteria or open tasks

I-0021 was closed on 2026-10-01.

The cause:

```text
flai accept merged a story's branch and removed its worktree before refusing it for unticked criteria or open tasks
```

The fix:

```text
Fixed by S-0041 (a9c9cd8): the done rules run on a copy before anything merges; TestAcceptChecksTheRulesBeforeItMergesAnything pins it.
```

Source: I-0021.
