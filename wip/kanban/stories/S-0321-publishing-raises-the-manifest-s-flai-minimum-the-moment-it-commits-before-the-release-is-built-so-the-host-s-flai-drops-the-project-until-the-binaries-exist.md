---
id: S-0321
type: story
nature: remediation
title: Publishing raises the manifest's flai minimum the moment it commits, before the release is built, so the host's flai drops the project until the binaries exist
status: backlog
owner: alex
created: 2026-10-07T18:59:52Z
updated: 2026-10-07T18:59:52Z
transitions: []
tags: []
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: sum
  seconds: 0
  models: []
  strategic:
    - kind: orchestrator
      seconds: 4
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 2
          output: 12
          cache_read: 64217
          cache_write: 5432
          cost: 0.0172
draft: true
cost_of_delay:
  inputs:
    time_lost_per_cycle: 1h
    by: flai
    at: 2026-10-07T18:59:52Z
---
# S-0321 Publishing raises the manifest's flai minimum the moment it commits, before the release is built, so the host's flai drops the project until the binaries exist

## Goal

This story remediates [I-0107](../../../design/issues/I-0107-publishing-raises-the-manifest-s-flai-minimum-the-moment-it-commits-before-the-release-is-built-so-the-host-s-flai-drops-the-project-until-the-binaries-exist.md), "Publishing raises the manifest's flai minimum the moment it commits, before the release is built, so the host's flai drops the project until the binaries exist". The issue recommends this solution:

Directions to weigh: have `flai release` raise `flai.minimum` only after the release's binaries are published, in a later commit, or have the publish leave the bump for the first run of the upgraded flai; or have `flai serve` keep serving a project whose manifest newly demands a minimum above its own version, with the warning it already logs for a flai older than the project, and refuse only the writes that touch fields it does not know. Either way the publish should not cancel its own `publish.run`: the journal recorded it as failed with `flai exited with -1` although the tags had been pushed. A test that publishes a release raising the minimum against a running `flai serve` of the older version, and expects the project to stay served, would pin it. The release build's own failure that day (a duplicated `## 1.0.67` heading in `template/CHANGELOG.md`, which `TestRepositoryLintsClean` rejects, since `flai release` prepended a second section with the heading the stories had already written) is a separate defect and made the outage last longer.

## Acceptance criteria
- [ ] The cause I-0107 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0107 is closed with `flai issue close I-0107 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0107. time_lost_per_cycle 1h: 30m per occurrence × 2 occurrences ÷ 1 cycle of 168h (first reported 2026-10-07T07:39:30Z, 0.5 days before this story; under one cycle counts as one).
