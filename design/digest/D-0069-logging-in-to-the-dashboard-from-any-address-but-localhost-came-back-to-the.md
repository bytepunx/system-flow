---
id: D-0069
title: Logging in to the dashboard from any address but localhost came back to the login page
kind: cause
sources: [I-0031]
topics: [flaiover, dashboard, sveltekit, code, all]
paths: [flaiover/src, flaiover/server.js, flaiover/Dockerfile, docs/operators]
date: 2026-10-01
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0069 Logging in to the dashboard from any address but localhost came back to the login page

I-0031 was closed on 2026-10-01.

The cause:

```text
Logging in to the dashboard from any address but localhost came back to the login page
```

The fix:

```text
Fixed by S-0083 (b6e263a): a login over the LAN sets a cookie without Secure.
```

Source: I-0031.
