---
id: T-3448
type: task
nature: feature
title: The template's manifest labels its tiers static, behavior, integration, and smoke, and this repository's manifest labels its own, with no change to what runs
status: done
parent: S-0805
owner: alex
created: 2026-11-16T23:48:56Z
updated: 2026-11-17T00:39:52Z
transitions:
  - to: ready
    at: 2026-11-17T00:13:45Z
    by: agent-S-0805
  - to: in-progress
    at: 2026-11-17T00:13:45Z
    by: agent-S-0805
  - to: done
    at: 2026-11-17T00:39:52Z
    by: agent-S-0805
stream: S-0805
tags: [template, flai]
touches: [template/root/system-flow.yaml.tmpl, template/CHANGELOG.md, system-flow.yaml]
after: [T-3446]
usage:
  source: log
  seconds: 1567
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 25
      output: 5449
      cache_read: 1682102
      cache_write: 80217
      cost: 0.9857
---
# T-3448 The template's manifest labels its tiers static, behavior, integration, and smoke, and this repository's manifest labels its own, with no change to what runs

## Work

In `template/root/system-flow.yaml.tmpl`, give `markdown` `suite: static` and `test` `suite: behavior`; `integration` and `smoke` are already the suites of their names, so label them explicitly for a reader, and say in the `tests:` comment what a suite is. Add a line under `## Unreleased` in `template/CHANGELOG.md`. In this repository's `system-flow.yaml`, label the eighteen tiers, changing no `name`, `paths`, `covers`, `queue`, `when`, `all_only`, or `manual`. The planner's proposal, for the story's agent to settle in `## Decisions`: `static` for `gofmt`, `vet`, `golangci-lint`, `shellcheck`, `markdown`, `flaiover-lint`, and `svelte-check`; `behavior` for `go-test`, `settings-index`, `vitest`, `flaiover-flai`, `design-routes`, `template`, and `flaiover`; `integration` for `integration` and the manual `race`; `smoke` for `smoke` and the manual `gateway`. Waits for T-3446, whose check accepts the key; it shares no path with T-3447 and runs beside it.

## Done when

- [ ] The template's four tiers carry `static`, `behavior`, `integration`, and `smoke`, and `scripts/template-test.sh` renders a project that passes `flai check --strict`.
- [ ] Every tier of `system-flow.yaml` carries a suite, and `flai test --all` and `flai test` on a sample of paths select the same tiers before and after, as the task's Notes record.
- [ ] `flai check --strict` is clean.

## Notes

What runs is unchanged. `flai test` has no list or dry-run mode, so the comparison was made on the parsed manifests. Each file was loaded as YAML beside main's copy (`git show main:<path>`), with the `suite` key dropped from every tier. This repository's whole manifest then equals main's, and the template's `tests` list equals main's. Selection reads a tier's suite only under `--suite` (T-3447), so `flai test --all`, `flai test <paths>`, `flai verify`, and the batch choose the same tiers as before. `scripts/template-test.sh` exited 0, and `flai test` on the three paths passed.
