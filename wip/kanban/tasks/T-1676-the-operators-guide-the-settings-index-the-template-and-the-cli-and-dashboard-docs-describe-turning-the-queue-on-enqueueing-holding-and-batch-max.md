---
id: T-1676
type: task
nature: feature
title: The operators' guide, the settings index, the template, and the CLI and dashboard docs describe turning the queue on, enqueueing, holding, and batch_max
status: backlog
parent: S-0399
owner: alex
created: 2026-10-16T05:03:31Z
updated: 2026-10-16T05:03:31Z
transitions: []
stream: S-0399
tags: [cli, dashboard]
touches: [template/root/system-flow.yaml.tmpl, template/template.yaml, template/CHANGELOG.md, docs/operators/index.md, docs/operators/settings.md, design/system/flai-cli.md, design/system/flaiover-dashboard.md, docs/users/flai.md, docs/users/flai-reference.md, docs/users/flaiover.md]
after: [T-1668, T-1658]
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
