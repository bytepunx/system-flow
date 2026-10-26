---
id: T-2124
type: task
nature: improvement
title: The user docs, the analyzer's convention, and its template copy read the curator's runs beside the digest's repeats
status: done
parent: S-0504
owner: alex
created: 2026-10-26T04:09:59Z
updated: 2026-10-26T04:25:41Z
transitions:
  - to: ready
    at: 2026-10-26T04:14:09Z
    by: agent-S-0504
  - to: in-progress
    at: 2026-10-26T04:14:10Z
    by: agent-S-0504
  - to: done
    at: 2026-10-26T04:25:40Z
    by: agent-S-0504
stream: S-0504
tags: [docs, digest, template]
touches: [docs/users/flai.md, design/conventions/strategic-agents.md, template/root/design/conventions/strategic-agents.md, template/CHANGELOG.md, template/template.yaml]
after: [T-2121]
usage:
  source: log
  seconds: 690
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 19
      output: 6771
      cache_read: 980386
      cache_write: 42135
      cost: 0.593
---
# T-2124 The user docs, the analyzer's convention, and its template copy read the curator's runs beside the digest's repeats

## Work

Document what T-2121 defines, in the same shape it names:

- `docs/users/flai.md`: the `flai stats` sample output's `digest:` lines gain the `curator runs` figure, the `asked again –` reading, and the `curator` line; the section "Running the curator" says that until the curator has run on a project, `flai stats` reports no `asked_again` and `curator  never ran`, so a repeat count of 0 is a measurement only once it has, and that `flai serve enable curate` is what starts the measuring. Link the ADR.
- `design/conventions/strategic-agents.md` § As the analyzer, the bottlenecks bullet: read `digest.weeks[].curator_runs` and `digest.curator` beside `asked_again`, and say in the report when the repeat count was not measured because no curator run had reached the week, rather than reading its absence as 0. Keep the baseline edit above the `<!-- system-flow:end-of-baseline -->` marker, since the rule is the baseline's, and bump `updated`.
- `template/root/design/conventions/strategic-agents.md`: the same edit, so that the template and this project's copy stay the same above the marker, and a `template/CHANGELOG.md` entry naming it.

This waits on T-2121 because the figures' names and the printed text are the ADR's. It can run beside T-2123, which it shares no path with.

## Done when

- `docs/users/flai.md` shows the new `digest` lines in the stats sample and says, under "Running the curator", when the repeat count is a measurement.
- `design/conventions/strategic-agents.md` § As the analyzer tells the analyzer to read `curator_runs` and `digest.curator` and to say when the repeats were not measured, and `template/root/design/conventions/strategic-agents.md` carries the same text above the marker.
- `template/CHANGELOG.md` has an entry for the convention change.
- The markdown lint passes on the four files, and `flai check --strict` is clean for the story's scope.

## Notes

Drafted by planner-S-0504. `docs/users/flai-reference.md` and `docs/operators/settings.md` are left alone: no flag, command, or setting changes.
