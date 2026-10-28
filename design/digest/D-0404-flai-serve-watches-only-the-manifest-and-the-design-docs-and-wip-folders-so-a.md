---
id: D-0404
title: flai serve watches only the manifest and the design, docs, and wip folders, so a component file added under flaiover/src reaches the Components page only at the next ask
kind: decision
sources: [S-0439]
topics: [flai, serve, flaiover, dashboard]
paths: [flai/internal/serve/serve.go, design/system/component-library.md]
date: 2026-10-28
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0404 flai serve watches only the manifest and the design, docs, and wip folders, so a component file added under flaiover/src reaches the Components page only at the next ask

`watchedPaths` in `flai/internal/serve/serve.go` names "the manifest's three folders and the manifest itself". A file added under a sub-project's source, such as a new component under `flaiover/src`, sends no change, so the Components page lists it only when the page asks again. S-0439 found the page's comments claimed otherwise, and corrected them (T-1869).

Source: S-0439.
