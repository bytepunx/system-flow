---
id: D-0058
title: Files needed by tests or checks were missing in a fresh clone
kind: cause
sources: [I-0003]
topics: []
paths: []
date: 2026-10-01
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0058 Files needed by tests or checks were missing in a fresh clone

I-0003 was closed on 2026-10-01.

The cause:

```text
The working tree passed but a clone did not: a test fixture under `bin/` was excluded by the root `.gitignore`, and `wip/kanban/tasks` was empty after archiving so git dropped the folder and `flai check` failed on the clone.
```

The fix:

```text
Fixed: CI runs on a clean checkout (.github/workflows/flai.yml, system-flow-check.yml) and has been green; no instance since 2026-09-15.
```

Source: I-0003.
