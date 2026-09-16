---
id: I-008
title: Path-filtered flai workflow did not run on the initial push to main
class: impression
status: closed
count: 1
first_reported: 2026-09-16T04:30:44Z
last_reported: 2026-09-16T04:30:44Z
updated: 2026-09-16T04:35:16Z
---

# I-008 Path-filtered flai workflow did not run on the initial push to main

## Description
`.github/workflows/flai.yml` filters on `flai/**`. The first push of the repository to GitHub triggered `system-flow check` but not `flai`, although every file under `flai/` was new. It did run for Dependabot pull requests minutes later. Suspected: path filters need a previous commit to diff against on `push`.

## Instances

### 2026-09-16T04:30:44Z
First push of main to github.com/bytepunx/system-flow. No `flai` run listed for the push event.

## Remediation
Closed 2026-09-16T04:35:16Z: the next push touching `flai/` triggered the workflow. First-push quirk of path filters; no change needed.
