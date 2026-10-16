---
id: T-1857
type: task
nature: feature
title: The CLI design and the user docs describe flai components list and flai components config
status: done
parent: S-0437
owner: alex
created: 2026-10-16T21:32:56Z
updated: 2026-10-16T22:29:45Z
transitions:
  - to: ready
    at: 2026-10-16T22:26:37Z
    by: agent-S-0437
  - to: in-progress
    at: 2026-10-16T22:26:38Z
    by: agent-S-0437
  - to: done
    at: 2026-10-16T22:29:45Z
    by: agent-S-0437
stream: S-0437
tags: [cli]
touches: [design/system/flai-cli.md, docs/users/flai.md]
after: [T-1856]
usage:
  source: log
  seconds: 187
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 48
      output: 19957
      cache_read: 2257896
      cache_write: 80861
      cost: 1.3202
---
# T-1857 The CLI design and the user docs describe flai components list and flai components config

## Work

Describe `flai components list` and `flai components config` in `design/system/flai-cli.md` and `docs/users/flai.md`: what the first lists, the prop kinds, what it needs on the host (`node` and the sub-project's installed `svelte`), and its errors; what the second writes and refuses.

Waits for T-1856, so that the docs describe the commands as built.

## Done when

Both files describe the commands as they behave, and `flai test` on them passes.

## Notes

None.
