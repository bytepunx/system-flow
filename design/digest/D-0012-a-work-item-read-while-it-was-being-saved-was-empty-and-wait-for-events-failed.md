---
id: D-0012
title: A work item read while it was being saved was empty, and wait_for_events failed with no front matter
kind: cause
sources: [I-0036]
topics: [flai, cli, go, code, all]
paths: [flai/internal]
date: 2026-09-23
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0012 A work item read while it was being saved was empty, and wait_for_events failed with no front matter

I-0036 was closed on 2026-09-23.

The cause:

```text
A work item read while it was being saved was empty, and wait_for_events failed with no front matter
```

The fix:

```text
fixed by S-0100: items, the board, narratives, and threads are replaced in one step (atomicfile.WriteFile)
```

Source: I-0036.
