---
id: D-0060
title: flai default cache_dir wrote a template clone under the operator's home
kind: cause
sources: [I-0009]
topics: [all]
paths: []
date: 2026-10-01
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0060 flai default cache_dir wrote a template clone under the operator's home

I-0009 was closed on 2026-10-01.

The cause:

```text
flai default cache_dir wrote a template clone under the operator's home
```

The fix:

```text
Fixed by S-0003 and S-0006: config honours FLAI_CACHE_DIR and scripts/env.sh sets the cache inside the repository; ~/.flai/cache does not exist.
```

Source: I-0009.
