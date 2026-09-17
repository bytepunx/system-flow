---
id: I-0009
title: flai default cache_dir wrote a template clone under the operator's home
class: defect
status: open
count: 2
cost: 3m
first_reported: 2026-09-17T00:10:55Z
last_reported: 2026-09-17T00:18:41Z
updated: 2026-09-17T00:18:41Z
---

# I-0009 flai default cache_dir wrote a template clone under the operator's home

## Description
flai default cache_dir wrote a template clone under the operator's home

## Instances

### 2026-09-17T00:10:55Z
S-0003: rendering from the published repo cloned into ~/.flai/cache because .flai-cache/config.json carried the default cache_dir. Fixed by setting cache_dir in the repo config and in scripts/env.sh on first run; the stray clone was removed.

### 2026-09-17T00:18:41Z
S-0006: a test resolved a relative template path against the test's working directory, treated it as a remote, and created ~/.flai/cache. Removed. Config now honours FLAI_CACHE_DIR; test helpers and scripts/env.sh set it so neither tests nor scripts can reach the home directory.

## Remediation
