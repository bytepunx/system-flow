---
id: T-1676
type: task
nature: feature
title: The operators' guide, the settings index, the template, and the CLI and dashboard docs describe turning the queue on, enqueueing, holding, and batch_max
status: done
parent: S-0399
owner: alex
created: 2026-10-16T05:03:31Z
updated: 2026-10-19T20:04:13Z
transitions:
  - to: ready
    at: 2026-10-19T19:55:56Z
    by: agent-S-0399
  - to: in-progress
    at: 2026-10-19T19:55:57Z
    by: agent-S-0399
  - to: done
    at: 2026-10-19T20:04:13Z
    by: agent-S-0399
stream: S-0399
tags: [cli, dashboard]
touches: [template/root/system-flow.yaml.tmpl, template/template.yaml, template/CHANGELOG.md, docs/operators/index.md, docs/operators/settings.md, design/system/flai-cli.md, design/system/flaiover-dashboard.md, docs/users/flai.md, docs/users/flai-reference.md, docs/users/flaiover.md, design/system/dashboard-host-channel.md, design/system/project-manifest.md]
after: [T-1668, T-1658]
usage:
  source: log
  seconds: 496
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 153
      output: 70804
      cache_read: 10692136
      cache_write: 392162
      cost: 5.6586
---
# T-1676 The operators' guide, the settings index, the template, and the CLI and dashboard docs describe turning the queue on, enqueueing, holding, and batch_max

## Work

- `docs/operators/index.md`: turning the `merge-queue` host action on, what it needs (`accept_reviews`), holding it, and what runs on the host.
- `docs/operators/settings.md`: rows for `merge_queue.batch_max` and the host action; the settings test fails without them.
- The template's `system-flow.yaml.tmpl` carries `merge_queue.batch_max` commented out. Bump `template.yaml`'s version and name the key in `template/CHANGELOG.md`.
- `design/system/flai-cli.md`, `docs/users/flai.md`: `flai queue add`, `remove`, `hold`, `resume`, and the batch. Regenerate `flai-reference.md`.
- `design/system/flaiover-dashboard.md` and `docs/users/flaiover.md`: the queue button, the hold switch, and the host writes.
- Waits for T-1668 and T-1658, whose behaviour it describes.

## Done when

The documents say the above, the settings test passes, and the markdown lint and `flai check --strict` pass.

## Notes
