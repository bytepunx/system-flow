---
id: T-2480
type: task
nature: feature
title: The kluster design, the CLI design, and the user guide describe deploy, teardown, promote, deployments, the record, and the logs
status: done
parent: S-0578
owner: alex
created: 2026-11-06T21:43:58Z
updated: 2026-11-07T02:07:01Z
transitions:
  - to: ready
    at: 2026-11-07T01:58:08Z
    by: agent-S-0578
  - to: in-progress
    at: 2026-11-07T01:58:08Z
    by: agent-S-0578
  - to: done
    at: 2026-11-07T02:00:27Z
    by: agent-S-0578
stream: S-0578
tags: [kluster, docs]
touches: [design/system/kluster.md, design/system/flai-cli.md, docs/users/flai.md, docs/users/flai-reference.md, docs/operators/settings.md]
after: [T-2479]
usage:
  source: log
  seconds: 139
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 45
      output: 17493
      cache_read: 6607983
      cache_write: 75352
      cost: 2.2156
---
# T-2480 The kluster design, the CLI design, and the user guide describe deploy, teardown, promote, deployments, the record, and the logs

## Work

It waits for T-2479, whose commands it documents as built.

- `design/system/kluster.md`: fill in the deploy section: the steps, the variables, the lock, the record's fields under `.flai-cache/kluster/deployments.json`, the logs, teardown, and promote.
- `design/system/flai-cli.md`: the four commands.
- `docs/users/flai.md`: the four commands, with an example of a deploy and its failure. Regenerate `docs/users/flai-reference.md` with `make flai-reference`.

Bump `updated` on each.

## Done when

- The three documents say what the story built, and the reference is regenerated.
- `flai check --strict` and the markdown lint pass on the four paths.

## Notes
