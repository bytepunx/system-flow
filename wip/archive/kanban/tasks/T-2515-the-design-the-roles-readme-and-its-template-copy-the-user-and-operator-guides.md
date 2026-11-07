---
id: T-2515
type: task
nature: remediation
title: The design, the roles README and its template copy, the user and operator guides, and the template changelog say a defined role runs once per ready story under each-ready-story
status: done
parent: S-0571
owner: alex
created: 2026-11-07T03:46:59Z
updated: 2026-11-07T04:06:29Z
transitions:
  - to: ready
    at: 2026-11-07T03:50:58Z
    by: agent-S-0571
  - to: in-progress
    at: 2026-11-07T04:05:33Z
    by: system-flow
  - to: done
    at: 2026-11-07T04:06:29Z
    by: agent-S-0571
stream: S-0571
tags: [docs, roles]
touches: [design/system/strategic-agents.md, design/workflows/roles/README.md, template/root/design/workflows/roles/README.md, docs/users/flai.md, docs/operators/settings.md, template/CHANGELOG.md]
after: [T-2511, T-2512]
usage:
  source: log
  seconds: 56
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 10
      output: 2103
      cache_read: 1210418
      cache_write: 22498
      cost: 0.4585
---
# T-2515 The design, the roles README and its template copy, the user and operator guides, and the template changelog say a defined role runs once per ready story under each-ready-story

## Work

Describe what T-2511 and T-2512 built, in the same story, as `documentation.md` asks:

- `design/system/strategic-agents.md` § On an event or a schedule (S-0418): the schedule row of the table gives the run's item as each ready story under `each-ready-story`. The first bullet under **Schedules** drops "does not yet start a defined role once per ready story". It says instead that each ready story gets one run, in pull order, with a held start kept per story. It also says `flai check` reports the condition on a role that does not take stories.
- `design/workflows/roles/README.md` and its identical template copy `template/root/design/workflows/roles/README.md`: the `roles.file` paragraph and the `each-ready-story` row say the condition applies to a defined role that takes stories as it does to the planner. Keep the two files identical.
- `docs/users/flai.md`, the **Schedules** bullet under § What starts a role your project defines: replace "A schedule starts a role that takes no item only" with what a schedule starts for a role that takes stories under `each-ready-story`.
- `docs/operators/settings.md`, the `roles.<name>.schedule` row: replace "a role that takes items is not started by one" with the same.
- `template/CHANGELOG.md`: one line under `## Unreleased`.

It waits for T-2511 and T-2512 because it describes the behaviour and the check message they settle.

## Done when

- None of the six files says a defined role that takes items is not started by a schedule under `each-ready-story`.
- The two roles READMEs are identical.
- `flai test` on the six paths passes, the markdown lint included.

## Notes
