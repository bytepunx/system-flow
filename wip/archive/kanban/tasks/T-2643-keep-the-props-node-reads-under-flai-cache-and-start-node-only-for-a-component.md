---
id: T-2643
type: task
nature: improvement
title: Keep the props node reads under .flai-cache and start node only for a component whose key changed
status: done
parent: S-0601
owner: alex
created: 2026-11-07T23:48:46Z
updated: 2026-11-08T01:36:34Z
transitions:
  - to: ready
    at: 2026-11-08T01:11:54Z
    by: agent-S-0601
  - to: in-progress
    at: 2026-11-08T01:12:00Z
    by: agent-S-0601
  - to: done
    at: 2026-11-08T01:36:34Z
    by: agent-S-0601
stream: S-0601
tags: [cli]
touches: [flai/internal/components/props.go, flai/internal/components/propscache.go, flai/internal/components/propscache_test.go, flai/internal/components/props.mjs]
usage:
  source: log
  seconds: 1474
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 35
      output: 8420
      cache_read: 5495079
      cache_write: 144047
      cost: 2.173
---
# T-2643 Keep the props node reads under .flai-cache and start node only for a component whose key changed

## Work

`ReadProps` in `flai/internal/components/props.go` starts `node` for each sub-project on every call, 0.6 to 4.2 s. S-0478's skeleton cache, kept by `props.mjs` under `PropsCacheDir`, saves typescript's work but not node's start.

Add `flai/internal/components/propscache.go`, a kept read of each component's answer, under `.flai-cache` in a file `props.mjs`'s own pruning of `PropsCacheDir` cannot remove. Key each entry by:

- the component file's path, size, and modification time;
- the sub-project's svelte version, read from its installed `svelte/package.json`;
- `scriptHash`, so a flai with another script reads again;
- the stamps of the files `props.mjs` read for the component's skeletons. Have the script name them in its answer, and stat them as the key's other files are. Without them, a change to an imported type would keep a stale skeleton.

In `ReadProps`, after the `node` lookup, answer each component whose key is unchanged from the kept read. Start `node` only for the rest, once per sub-project as now, and keep what it answers. A component that errors is kept with its error under the same key. A kept file that does not parse is a miss, not an error. Write it through a temporary file and a rename, so two flai processes cannot leave half a file.

Test it in `flai/internal/components/propscache_test.go` with a fake `Runner` that counts starts: a second call with nothing changed starts no `node`; a change to one component's file, to the svelte version, or to a stamped file starts it for that component's sub-project only, with only the changed component on its input.

This task waits for none: it shares no file with T-2633.

## Done when

- A second `ReadProps` with nothing changed starts no `node`, and one with one component changed asks `node` for that component alone.
- `propscache_test.go` passes, and so does `go test -race ./internal/components/...`.

## Notes
