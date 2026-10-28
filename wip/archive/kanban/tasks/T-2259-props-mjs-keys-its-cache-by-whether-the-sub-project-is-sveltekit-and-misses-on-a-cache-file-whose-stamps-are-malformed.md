---
id: T-2259
type: task
nature: remediation
title: props.mjs keys its cache by whether the sub-project is SvelteKit and misses on a cache file whose stamps are malformed
status: done
parent: S-0528
owner: alex
created: 2026-10-28T18:04:04Z
updated: 2026-10-28T18:12:35Z
transitions:
  - to: ready
    at: 2026-10-28T18:05:31Z
    by: agent-S-0528
  - to: in-progress
    at: 2026-10-28T18:05:32Z
    by: agent-S-0528
  - to: done
    at: 2026-10-28T18:12:35Z
    by: agent-S-0528
stream: S-0528
tags: [flai, components]
touches: [flai/internal/components/props.mjs, flai/internal/components/props_test.go, design/system/component-library.md]
usage:
  source: log
  seconds: 423
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 34
      output: 8309
      cache_read: 1567757
      cache_write: 83657
      cost: 1.0741
---
# T-2259 props.mjs keys its cache by whether the sub-project is SvelteKit and misses on a cache file whose stamps are malformed

## Work

Fix both causes I-0198 names in `flai/internal/components/props.mjs`, with a Go test for each in `flai/internal/components/props_test.go`. It waits for no task: it is the first layer.

1. **Stale skeleton on `@sveltejs/kit`.** `compilerOptions` calls `isSvelteKit()`, which reads `package.json` through `require`, outside the `ts.sys` stamps, and maps `$lib` when `@sveltejs/kit` is a dependency. Neither `package.json` nor that answer is in the key. Add the answer of `isSvelteKit()` to the array `cachedSkeletons` hashes into the key, so that adding or removing `@sveltejs/kit` alone gives a key of its own. This is the proposed fix. Stamping `package.json` in `stamps` would work too, but would miss on every edit to `package.json`, such as a version bump. If you choose that instead, say why in the narrative's `## Decisions`.
2. **Throw on malformed stamps.** In `readCache`, the loop `for (const [[kind, p], was] of cache.stamps)` sits outside any `try`. A cache file that is valid JSON but holds a stamp entry that is not `[[kind, path], stamp]`, such as `1`, `[]`, or `[[null]]`, throws and fails the run. Make such a file a miss: check each entry's shape, or put the loop in a `try` that returns `null`. A miss resolves the skeletons and writes the file again.
3. **Tests.**
   - With the typed fixture, a `tsconfig.json` that extends the missing `.svelte-kit/tsconfig.json`, and a `package.json` without `@sveltejs/kit`, a first `List` gives `story` its fallback. Then add `@sveltejs/kit` to `devDependencies` with `writeLater`. A second `List` gives `storySkeleton`, as `TestListMapsLibInASvelteKitSubProjectWhoseGeneratedTsconfigIsMissing` does. Fail the test against the current script first.
   - After a first `List`, rewrite the one cache file as valid JSON whose `stamps` holds malformed entries. A second `List` returns no error and gives the resolved skeletons.
4. Update `design/system/component-library.md` § Props in the same commit: the **Key** bullet names whether the sub-project is SvelteKit, and the **Miss** bullet says a cache file that cannot be read, or whose stamps are malformed, is a miss.

`props.go` needs no change: `scriptHash` is the script's own hash, so the changed script keys new files and old ones miss.

## Done when

- The cache key holds whether the sub-project is SvelteKit, and adding `@sveltejs/kit` to a cached sub-project's `package.json` gives the `$lib`-resolved skeleton on the next run.
- A cache file with malformed stamp entries is a miss, not an error.
- Both new tests fail against the script before the change and pass after it.
- `design/system/component-library.md` § Props says both.
- `flai test` on the three changed paths passes.

## Notes

Drafted by planner-S-0528. The tests need `node` and the fixture's `typescript`, as the other cache tests in `props_test.go` do.
