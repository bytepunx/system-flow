---
id: D-0063
title: Acceptance from the dashboard is refused for files only the host's global gitignore hides, and the preview does not show it
kind: cause
sources: [I-0019]
topics: [flai, cli, go, flaiover, dashboard, sveltekit, template, conventions, code, all]
paths: [flaiover, flai/cmd, template]
date: 2026-10-01
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0063 Acceptance from the dashboard is refused for files only the host's global gitignore hides, and the preview does not show it

I-0019 was closed on 2026-10-01.

The cause:

```text
Acceptance from the dashboard is refused for files only the host's global gitignore hides, and the preview does not show it
```

The fix:

```text
Fixed by S-0051 (2e85a74): the acceptance preview lists uncommitted paths outside wip, the designer can include them, and the container gets the host's global excludes.
```

Source: I-0019.
