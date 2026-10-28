---
id: I-0198
title: The component props cache (S-0478) can serve a stale skeleton when only @sveltejs/kit is added or removed, and throws instead of missing on malformed stamp entries
class: defect
status: open
count: 1
first_reported: 2026-10-23T07:04:46Z
last_reported: 2026-10-23T07:04:46Z
updated: 2026-10-28T12:54:18Z
---

# I-0198 The component props cache (S-0478) can serve a stale skeleton when only @sveltejs/kit is added or removed, and throws instead of missing on malformed stamp entries

## Description
The component props cache (S-0478) can serve a stale skeleton when only @sveltejs/kit is added or removed, and throws instead of missing on malformed stamp entries

## Instances

### 2026-10-23T07:04:46Z
Found by the orchestrator's verifier reviewing S-0478 at 886afe03 (accepted through the merge queue). In flai/internal/components/props.mjs, isSvelteKit() reads package.json through require, outside the ts.sys stamps, and adds the $lib path mapping when @sveltejs/kit is a dependency; neither package.json nor that result is in the cache key, so adding or removing @sveltejs/kit alone could leave a stale skeleton. And in readCache the loop over cache.stamps sits outside the try, so a cache file that is valid JSON with malformed stamp entries throws and fails the run rather than missing. Both are rare; no test covers a corrupt cache file.

## Remediation

Story S-0528 remediates this issue, created from it at 2026-10-28T12:54:18Z.
