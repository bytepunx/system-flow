---
id: D-0266
title: flai dashboard restart starts the container's tag again, so a newer image a check pulled under a floating tag is started without an upgrade
kind: cause
sources: [I-0116]
topics: [flai, cli, go, flaiover, dashboard, sveltekit, code, all]
paths: [flai/cmd/dashboard_upgrade.go, flai/cmd/dashboard.go, flai/cmd/dashboard_test.go, flai/cmd/dashboard_watch.go, flai/cmd/dashboard_watch_test.go, docs/users/flai-reference.md, docs/users/flai.md, docs/operators/index.md, design/system/flai-cli.md, design/issues/I-0116-flai-dashboard-restart-starts-the-container-s-tag-again-so-a-newer-image-a-check-pulled-under-a-floating-tag-is-started-without-an-upgrade.md, design/issues/summary.md]
date: 2026-10-08
status: absorbed
absorbed_by: "design/system/flai-cli.md#listing-and-installing-a-published-release"
asked_again: []
cited: []
---

# D-0266 flai dashboard restart starts the container's tag again, so a newer image a check pulled under a floating tag is started without an upgrade

I-0116 was closed on 2026-10-08.

The cause:

```text
flai dashboard restart starts the container's tag again, so a newer image a check pulled under a floating tag is started without an upgrade
```

The fix:

```text
S-0344: flai dashboard restart (T-1360) and flai host's watch (T-1361) start the image ID the container runs, never its tag, so a newer image flai dashboard check pulled under a floating tag is not started without an upgrade; each container carries the io.bytepunx.flai.ref label so the tag is still shown. Tests in flai/cmd/dashboard_test.go and dashboard_watch_test.go reproduce it.
```

Source: I-0116.
