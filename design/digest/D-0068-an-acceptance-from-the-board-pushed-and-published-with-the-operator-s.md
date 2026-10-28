---
id: D-0068
title: An acceptance from the board pushed and published with the operator's credentials, with nothing enabling it
kind: cause
sources: [I-0028]
topics: [flai, cli, go, flaiover, dashboard, sveltekit, code, all]
paths: [flai/cmd, flai/internal, flaiover/src, design/system, docs/operators]
date: 2026-10-01
status: absorbed
absorbed_by: "design/system/pushing-from-the-board.md#pushing-an-acceptance-made-from-the-board"
asked_again: []
cited: []
---

# D-0068 An acceptance from the board pushed and published with the operator's credentials, with nothing enabling it

I-0028 was closed on 2026-10-01.

The cause:

```text
An acceptance from the board pushed and published with the operator's credentials, with nothing enabling it
```

The fix:

```text
Fixed by S-0078 (414d1ae) and S-0087: accept.run only merges and archives; push.run and publish.run sit behind the push host action, off until the operator enables it.
```

Source: I-0028.
