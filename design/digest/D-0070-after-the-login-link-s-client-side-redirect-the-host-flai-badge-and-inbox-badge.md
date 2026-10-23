---
id: D-0070
title: After the login link's client-side redirect, the host flai badge and inbox badge never start polling
kind: cause
sources: [I-0032]
topics: [flai, cli, go, flaiover, dashboard, sveltekit, code, all]
paths: [flai/cmd, flai/internal, flaiover/src, design/system, design/adrs, docs/operators]
date: 2026-10-01
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0070 After the login link's client-side redirect, the host flai badge and inbox badge never start polling

I-0032 was closed on 2026-10-01.

The cause:

```text
After the login link's client-side redirect, the host flai badge and inbox badge never start polling
```

The fix:

```text
Fixed by S-0080 (dacde66): the layout derives signedIn reactively from page.url, so the badges start polling after the client-side redirect.
```

Source: I-0032.
