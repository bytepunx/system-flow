---
id: I-0038
title: A smoke test that runs the latest published flai went red once a release changed what it pinned, after that release's story was accepted
class: defect
status: open
count: 1
cost: 5m
first_reported: 2026-09-24T06:27:32Z
last_reported: 2026-09-24T06:27:32Z
updated: 2026-09-24T06:27:32Z
---

# I-0038 A smoke test that runs the latest published flai went red once a release changed what it pinned, after that release's story was accepted

## Description
A smoke test that runs the latest published flai went red once a release changed what it pinned, after that release's story was accepted

## Instances

### 2026-09-24T06:27:32Z
S-0112 (2026-09-24): scripts/install-test.sh installs the latest release under .flai-cache inside the project and expected self-upgrade --check to name that path. S-0111 made a binary inside a project upgrade into ~/.flai/bin; its own smoke run passed against 1.15.2, and the tier went red only once flai/v1.15.3 was published. Fixed in S-0112's T-0391 by checking a copy outside any project.

## Remediation
