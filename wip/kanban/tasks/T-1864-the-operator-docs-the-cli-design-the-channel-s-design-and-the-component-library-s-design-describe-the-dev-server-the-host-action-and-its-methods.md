---
id: T-1864
type: task
nature: feature
title: The operator docs, the CLI design, the channel's design, and the component library's design describe the dev server, the host action, and its methods
status: done
parent: S-0438
owner: alex
created: 2026-10-16T21:34:18Z
updated: 2026-10-17T00:39:46Z
transitions:
  - to: ready
    at: 2026-10-17T00:32:50Z
    by: system-flow
  - to: in-progress
    at: 2026-10-17T00:32:51Z
    by: system-flow
  - to: done
    at: 2026-10-17T00:39:46Z
    by: system-flow
stream: S-0438
tags: [cli, dashboard]
touches: [docs/operators/index.md, docs/operators/settings.md, design/system/flai-cli.md, docs/users/flai.md, design/system/dashboard-host-channel.md, design/system/component-library.md]
after: [T-1862, T-1863]
usage:
  source: log
  seconds: 405
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 96
      output: 517
      cache_read: 5135875
      cache_write: 155524
      cost: 2.4599
---
# T-1864 The operator docs, the CLI design, the channel's design, and the component library's design describe the dev server, the host action, and its methods

## Work

- `docs/operators/index.md`: a section "The components host action", on the pattern of the checks section: what it runs on the host, how to enable it, the idle stop, and the loopback port.
- `docs/operators/settings.md`: rows for the action and `components.idle`.
- `design/system/flai-cli.md` and `docs/users/flai.md`: `flai components serve`.
- `design/system/dashboard-host-channel.md`: the four methods and the action, beside the other host actions.
- `design/system/component-library.md`: the dev server as built, with what T-1859 found about SvelteKit.

Waits for T-1862 and T-1863, so that the docs describe what was built and seen.

## Done when

The six files describe the dev server as it behaves, the settings test passes, and `flai test` on them passes.

## Notes

None.
