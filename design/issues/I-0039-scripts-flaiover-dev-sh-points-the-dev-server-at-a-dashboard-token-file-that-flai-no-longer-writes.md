---
id: I-0039
title: scripts/flaiover-dev.sh points the dev server at a dashboard token file that flai no longer writes
class: defect
status: open
count: 2
cost: 4m
first_reported: 2026-09-24T07:50:21Z
last_reported: 2026-09-24T09:41:52Z
updated: 2026-09-24T09:41:52Z
---

# I-0039 scripts/flaiover-dev.sh points the dev server at a dashboard token file that flai no longer writes

## Description
scripts/flaiover-dev.sh points the dev server at a dashboard token file that flai no longer writes

## Instances

### 2026-09-24T07:50:21Z
S-0113: scripts/flaiover-dev.sh defaults FLAIOVER_TOKEN_FILE to <project>/.flai-cache/dashboard.token. flai dashboard token reported the token in ~/.flai/serve/dashboard.token. Every request to the dev server failed with ENOENT in auth.ts initAuth (a 500 page) until FLAIOVER_TOKEN_FILE was set to the file flai reports.

### 2026-09-24T09:41:52Z
S-0028: flaiover/compose.yaml has the same stale token path (${PROJECT}/.flai-cache/dashboard.token) and still mounts the repository and sets PROJECT_DIR, which the container has not used since ADR-0031; found while indexing flaiover's settings, not fixed here

## Remediation
