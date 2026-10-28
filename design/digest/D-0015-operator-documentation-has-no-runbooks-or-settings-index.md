---
id: D-0015
title: Operator documentation has no runbooks or settings index
kind: cause
sources: [I-0005]
topics: [flai, cli, go, code, all]
paths: [docs/operators, flai/cmd, scripts, Makefile, design/system, design/conventions, design/issues]
date: 2026-09-24
status: absorbed
absorbed_by: "design/conventions/documentation.md#project-additions"
asked_again: []
cited: []
---

# D-0015 Operator documentation has no runbooks or settings index

I-0005 was closed on 2026-09-24.

The cause:

```text
The documentation convention requires operator docs to cover every setting with an index and to include runbooks for installation, updates, deletion, backups, restores, and migrations. `docs/operators/index.md` is a draft with the dashboard run command and a security note.
```

The fix:

```text
S-0028: docs/operators/settings.md indexes every setting and docs/operators/runbooks/ has install, update, backup, restore, migrate, and delete
```

Source: I-0005.
