---
id: D-0003
title: Moving a story to done from the board marks it done without accepting it
kind: cause
sources: [I-0015]
topics: [flai, cli, go, flaiover, dashboard, sveltekit, code, all]
paths: [flai/cmd/move.go, flai/cmd/accept.go, flaiover/src/routes/board, flaiover/src/routes/items]
date: 2026-09-18
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0003 Moving a story to done from the board marks it done without accepting it

I-0015 was closed on 2026-09-18.

The cause:

```text
Moving a story to done from the board marks it done without accepting it
```

The fix:

```text
S-0046: moving a story to done runs acceptance; accept completes stranded stories; check flags them
```

Source: I-0015.
