---
id: T-2699
type: task
nature: feature
title: flai writes issues and ADRs by write and rename, so a change moves its folder's modification time
status: done
parent: S-0606
owner: alex
created: 2026-11-09T07:35:38Z
updated: 2026-11-09T07:47:07Z
transitions:
  - to: ready
    at: 2026-11-09T07:36:27Z
    by: system-flow
  - to: in-progress
    at: 2026-11-09T07:36:29Z
    by: system-flow
  - to: done
    at: 2026-11-09T07:47:07Z
    by: agent-S-0606
stream: S-0606
tags: []
touches: [flai/internal/issues/issues.go, flai/internal/issues/issues_test.go, flai/internal/adr/adr.go, flai/internal/adr/adr_test.go, flai/internal/adrindex/adrindex.go, flai/internal/issues/generated.go]
usage:
  source: log
  seconds: 638
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 40
      output: 264
      cache_read: 10220458
      cache_write: 585077
      cost: 5.239
---

# T-2699 flai writes issues and ADRs by write and rename, so a change moves its folder's modification time

## Work

TH-0908's answer: flai's own issue and ADR writes become atomic, written to a temporary file and renamed over the target, as threads, conversations, and digest entries already are, through `flai/internal/atomicfile`. A rename moves the folder's modification time, which the watcher's slow tier sees at its next fast look. An in-place write does not move it. Change every `os.WriteFile` of an issue file in `flai/internal/issues` and of an ADR or the ADR index in `flai/internal/adr` to `atomicfile.WriteFile`, keeping each file's mode.

Waits for no open task.

## Done when

- [ ] Every issue and ADR write in those two packages goes through `atomicfile`, and their tests pass.
- [ ] A test in each package shows that a write replaces the file, giving it a new inode, rather than writing it in place.

## Notes
