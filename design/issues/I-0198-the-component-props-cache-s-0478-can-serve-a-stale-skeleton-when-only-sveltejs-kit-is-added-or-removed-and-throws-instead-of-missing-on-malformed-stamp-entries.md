---
id: I-0198
title: The component props cache (S-0478) can serve a stale skeleton when only @sveltejs/kit is added or removed, and throws instead of missing on malformed stamp entries
class: defect
status: closed
count: 1
first_reported: 2026-10-23T07:04:46Z
last_reported: 2026-10-23T07:04:46Z
updated: 2026-10-28T18:13:52Z
---

# I-0198 The component props cache (S-0478) can serve a stale skeleton when only @sveltejs/kit is added or removed, and throws instead of missing on malformed stamp entries

## Description
The component props cache (S-0478) can serve a stale skeleton when only @sveltejs/kit is added or removed, and throws instead of missing on malformed stamp entries

## Instances

### 2026-10-23T07:04:46Z
Found by the orchestrator's verifier reviewing S-0478 at 886afe03 (accepted through the merge queue). In flai/internal/components/props.mjs, isSvelteKit() reads package.json through require, outside the ts.sys stamps, and adds the $lib path mapping when @sveltejs/kit is a dependency; neither package.json nor that result is in the cache key, so adding or removing @sveltejs/kit alone could leave a stale skeleton. And in readCache the loop over cache.stamps sits outside the try, so a cache file that is valid JSON with malformed stamp entries throws and fails the run rather than missing. Both are rare; no test covers a corrupt cache file.

## Remediation

Story S-0528 remediates this issue, created from it at 2026-10-28T12:54:18Z.
Closed 2026-10-28T18:13:52Z: Fixed by S-0528 (T-2259): props.mjs now puts whether the sub-project is SvelteKit, the answer of isSvelteKit() it reads from package.json outside the stamps, in the cache key, so adding or removing @sveltejs/kit alone gives a key of its own while a version bump does not; and readCache checks each stamp is [[kind, path], stamp] of strings, inside a try, so a cache file with malformed stamps, or a path stat refuses, is a miss rather than an error. TestListMapsLibOnceSvelteKitIsAddedToACachedSubProject and TestListTakesACacheFileWithMalformedStampsAsAMiss in flai/internal/components/props_test.go reproduce the two causes.
