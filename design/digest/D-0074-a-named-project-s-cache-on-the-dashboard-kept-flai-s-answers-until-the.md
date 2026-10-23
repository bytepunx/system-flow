---
id: D-0074
title: A named project's cache on the dashboard kept flai's answers until the container restarted, because only the default project's Repo watched
kind: cause
sources: [I-0046]
topics: [flaiover, dashboard, sveltekit, code, all]
paths: [flaiover/src, design/issues]
date: 2026-10-01
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0074 A named project's cache on the dashboard kept flai's answers until the container restarted, because only the default project's Repo watched

I-0046 was closed on 2026-10-01.

The cause:

```text
flaiover keeps one cache of flai's answers per project, a `Repo` made by `repo()` for the project key a request names. A cache is only forgotten when its `Repo` listens for flai's `change`, `connected`, and `gone` events, and only some `Repo`s were ever told to listen: the default project's, at startup, and the one `/api/events` or `/api/stats` was asked for. A page that asked `/api/board?project=<key>` without either of those having reached that project saw what flai said the first time until the container restarted.
```

The fix:

```text
Fixed by S-0117: repo() watches every Repo it makes; project-events.test.ts pins it.
```

Source: I-0046.
