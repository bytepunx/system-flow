---
id: D-0067
title: More than three tags in one push create no tag events, so release flai did not run
kind: cause
sources: [I-0026]
topics: [flai, cli, go, flaiover, dashboard, sveltekit, template, conventions, code, all]
paths: [flai/cmd, flai/internal, flaiover/src, design/system, design/adrs, template]
date: 2026-10-01
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0067 More than three tags in one push create no tag events, so release flai did not run

I-0026 was closed on 2026-10-01.

The cause:

```text
More than three tags in one push create no tag events, so release flai did not run
```

The fix:

```text
Fixed by S-0076 (c1ae725) and S-0087 (130213b): pushes are batched at three tags (pending.TagsPerPush).
```

Source: I-0026.
