---
id: T-2472
type: task
nature: feature
title: flai's manifest reads the kluster block the ADR names, and flai check reports an invalid one with the key and the reason
status: done
parent: S-0577
owner: alex
created: 2026-11-06T21:42:28Z
updated: 2026-11-07T01:08:31Z
transitions:
  - to: ready
    at: 2026-11-07T00:55:40Z
    by: system-flow
  - to: in-progress
    at: 2026-11-07T00:55:41Z
    by: system-flow
  - to: done
    at: 2026-11-07T01:08:31Z
    by: agent-S-0577
stream: S-0577
tags: [kluster, manifest, check]
touches: [flai/internal/manifest/manifest.go, flai/internal/manifest/manifest_test.go, flai/internal/check/check.go, flai/internal/check/kluster_test.go]
usage:
  source: log
  seconds: 770
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 48
      output: 15112
      cache_read: 6750468
      cache_write: 75467
      cost: 2.1847
---
# T-2472 flai's manifest reads the kluster block the ADR names, and flai check reports an invalid one with the key and the reason

## Work

Read S-0576's ADR and `design/system/kluster.md` for the `kluster:` block's keys, which are required, and their defaults.

In `flai/internal/manifest/manifest.go`, add the block to the manifest's type, parsed from `system-flow.yaml`, with its defaults applied and a project with no block left as none. Test it in `manifest_test.go`: a full block, a block with only the required keys, no block, and an unknown key.

In `flai/internal/check/check.go`, add a finding of its own, named after the other manifest findings (such as `manifest.kluster`), for an invalid block: a required key missing, an unknown key, a config path that does not exist, a default namespace that is not a DNS label or that has a story namespace's form (`s-` and digits), a `timeout` that is not a Go duration above zero, a command list that is empty, and, read from the kluster config file itself, a config with no `name` or whose `provider` is other than `k3d` (absent means kluster's default, `k3d`). Each finding names the key and the reason. Test each case, and a valid block with no finding, in `flai/internal/check/kluster_test.go`.

This is the first layer: it waits for no task.

## Done when

- The manifest reads the block with its defaults, and a project with none parses as before.
- `flai check` reports each invalid case with its key and reason, and nothing for a valid block or none.
- `flai test` on the four paths passes.

## Notes
