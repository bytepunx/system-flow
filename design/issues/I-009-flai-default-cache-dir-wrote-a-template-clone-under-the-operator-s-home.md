---
id: I-009
title: flai default cache_dir wrote a template clone under the operator's home
class: defect
status: open
count: 1
cost: 2m
first_reported: 2026-09-17T00:10:55Z
last_reported: 2026-09-17T00:10:55Z
updated: 2026-09-17T00:10:55Z
---

# I-009 flai default cache_dir wrote a template clone under the operator's home

## Description
flai default cache_dir wrote a template clone under the operator's home

## Instances

### 2026-09-17T00:10:55Z
S-003: rendering from the published repo cloned into ~/.flai/cache because .flai-cache/config.json carried the default cache_dir. Fixed by setting cache_dir in the repo config and in scripts/env.sh on first run; the stray clone was removed.

## Remediation
