---
id: T-2565
type: task
nature: improvement
title: The docs say a project's kluster is set up from the Klusters page
status: done
parent: S-0586
owner: alex
created: 2026-11-07T08:32:12Z
updated: 2026-11-07T10:29:26Z
transitions:
  - to: ready
    at: 2026-11-07T10:18:20Z
    by: agent-S-0586
  - to: in-progress
    at: 2026-11-07T10:18:24Z
    by: agent-S-0586
  - to: done
    at: 2026-11-07T10:29:26Z
    by: agent-S-0586
stream: S-0586
tags: [kluster, docs]
touches: [design/system/kluster.md, design/system/flai-cli.md, design/system/flaiover-dashboard.md, design/tech/kluster.md, docs/users/flaiover.md, docs/users/flai.md, docs/operators/index.md, docs/operators/settings.md]
after: [T-2559, T-2561, T-2563, T-2564, T-2569]
usage:
  source: log
  seconds: 662
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 72
      output: 21471
      cache_read: 10778463
      cache_write: 137087
      cost: 3.3669
---
# T-2565 The docs say a project's kluster is set up from the Klusters page

## Work

Waits for T-2561, T-2563, and T-2564: it describes what they built. Any task the plan thread adds also goes in its `after`.

- `docs/users/flaiover.md`: § Klusters walks through setup. § Features no longer says `kluster` is listed only for a project with a block.
- `docs/operators/index.md` § The kluster host action: enabling writes the block and the config, and the sentence saying it is refused goes.
- `docs/operators/settings.md`: the `kluster.*` rows no longer say "Changed by hand", and `kluster.deploy` is no longer required.
- `docs/users/flai.md`, `design/system/flai-cli.md`: `flai serve enable kluster` and `flai manifest set kluster.*`.
- `design/system/kluster.md` § The dashboard: the new fields of `kluster.status`, `kluster.workloads`, and the Klusters page as built, and § flai serve: the pass with no `deploy`.
- `design/system/flaiover-dashboard.md` § Views: the page's setup.
- `design/tech/kluster.md`: what kluster's config takes, as T-2559 found it.

The story's criterion 3 is that nothing about a project's kluster needs a step outside the dashboard. Read each of these documents for such a step. Where one stays, such as the one host action turned on in a shell, say so plainly and name the plan thread's answer that keeps it.

## Done when

- `flai check --strict` and the markdown lint pass.
- The settings-index test passes.
- No document says the block is added by hand, or that enabling is refused for want of one.
- Every step left outside the dashboard is named in `docs/operators/index.md` with the reason it stays, and the story's review note lists them for the operator.

## Notes
