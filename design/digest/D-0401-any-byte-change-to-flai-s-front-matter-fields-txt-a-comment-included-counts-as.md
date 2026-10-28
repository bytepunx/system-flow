---
id: D-0401
title: Any byte change to flai's front-matter-fields.txt, a comment included, counts as a changed field list and raises flai.minimum on the next publish
kind: decision
sources: [S-0321]
topics: [release, flai, cli, code]
paths: [flai/internal/workitem/front-matter-fields.txt, flai/internal/release/release.go]
date: 2026-10-28
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0401 Any byte change to flai's front-matter-fields.txt, a comment included, counts as a changed field list and raises flai.minimum on the next publish

`design/system/work-hierarchy.md` says a release whose field list changed raises `flai.minimum`. S-0321 found that the comparison is of the file's bytes: `RaiseMinimum` in `flai/internal/release/release.go` reads a release as changing the file "when its fields file differs from the previous release's". S-0321 left the file's comment as it was, because "changing its bytes would itself make the next publish raise the minimum". Edit `front-matter-fields.txt` only when a field really changes.

Source: S-0321.
