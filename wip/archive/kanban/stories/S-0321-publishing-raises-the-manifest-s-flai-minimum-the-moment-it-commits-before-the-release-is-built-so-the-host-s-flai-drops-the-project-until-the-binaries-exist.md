---
id: S-0321
type: story
nature: remediation
title: Publishing raises the manifest's flai minimum the moment it commits, before the release is built, so the host's flai drops the project until the binaries exist
status: done
owner: alex
created: 2026-10-07T18:59:52Z
updated: 2026-10-08T09:00:11Z
transitions:
  - to: ready
    at: 2026-10-08T05:02:36Z
    by: orchestrator
  - to: in-progress
    at: 2026-10-08T07:58:14Z
    by: agent-S-0321
  - to: review
    at: 2026-10-08T09:00:01Z
    by: alex
  - to: done
    at: 2026-10-08T09:00:11Z
    by: alex
tags: [flai, release, serve]
topics: [release]
touches: [flai/internal/release/release.go, flai/internal/release/release_test.go, flai/cmd/release.go, flai/internal/serve/serve.go, flai/internal/serve/serve_test.go, flai/internal/channel/channel.go, design/system/flai-cli.md, design/system/project-manifest.md, docs/users/flai.md, docs/operators/index.md, design/issues/I-0107-publishing-raises-the-manifest-s-flai-minimum-the-moment-it-commits-before-the-release-is-built-so-the-host-s-flai-drops-the-project-until-the-binaries-exist.md, design/issues/summary.md, design/system/work-hierarchy.md, flai/internal/manifest/manifest.go, design/issues/I-0079-testroundtriprepositoryitems-reads-the-live-main-checkout-and-fails-a-close-out-when-another-agent-edits-a-story-mid-run.md, design/issues/I-0113-a-failed-integration-tier-in-the-close-out-shows-only-the-last-lines-of-go-test-so-the-failing-test-is-not-named.md, design/issues/I-0086-the-close-out-s-install-smoke-test-failed-once-and-passed-when-run-alone-with-no-cause-in-its-output.md, design/issues/I-0123-flai-check-finds-board-wip-limit-outside-the-story-at-close-out.md, design/issues/I-0119-the-close-out-s-last-check-that-the-branch-contains-main-fails-when-flai-commits-wip-on-main-during-its-run.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 3752
  turns:
    - day: 2026-10-08
      ceremony: 5
      test_runs: 3
      hand_edits: 4
      work: 72
  models:
    - model: claude-opus-5-5
      input: 328
      output: 107914
      cache_read: 21645514
      cache_write: 506062
      cost: 9.7847
  strategic:
    - kind: orchestrator
      seconds: 442
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 101
          output: 1454
          cache_read: 24920862
          cache_write: 56919
          cost: 6.1541
cost_of_delay:
  inputs:
    time_lost_per_cycle: 1h
    by: flai
    at: 2026-10-07T18:59:52Z
  value: 150
  by: planner-S-0321
  at: 2026-10-08T00:24:08Z
forecast:
  duration: 45m
  delivery: 2026-10-08T08:31:00Z
  basis: "Its own forecast of 45m; 1st in the pull order with an in-progress limit of 3, behind S-0232, S-0326 and S-0339."
  by: flai
  at: 2026-10-08T07:40:40Z
finalized:
  by: orchestrator
  at: 2026-10-08T00:24:59Z
---
# S-0321 Publishing raises the manifest's flai minimum the moment it commits, before the release is built, so the host's flai drops the project until the binaries exist

## Goal

This story remediates [I-0107](../../../design/issues/I-0107-publishing-raises-the-manifest-s-flai-minimum-the-moment-it-commits-before-the-release-is-built-so-the-host-s-flai-drops-the-project-until-the-binaries-exist.md), "Publishing raises the manifest's flai minimum the moment it commits, before the release is built, so the host's flai drops the project until the binaries exist". The issue recommends this solution:

Directions to weigh: have `flai release` raise `flai.minimum` only after the release's binaries are published, in a later commit, or have the publish leave the bump for the first run of the upgraded flai; or have `flai serve` keep serving a project whose manifest newly demands a minimum above its own version, with the warning it already logs for a flai older than the project, and refuse only the writes that touch fields it does not know. Either way the publish should not cancel its own `publish.run`: the journal recorded it as failed with `flai exited with -1` although the tags had been pushed. A test that publishes a release raising the minimum against a running `flai serve` of the older version, and expects the project to stay served, would pin it. The release build's own failure that day (a duplicated `## 1.0.67` heading in `template/CHANGELOG.md`, which `TestRepositoryLintsClean` rejects, since `flai release` prepended a second section with the heading the stories had already written) is a separate defect and made the outage last longer.

## Acceptance criteria
- [x] The cause I-0107 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0107 is closed with `flai issue close I-0107 --reason` saying what fixed it

## Tasks
- T-1309 flai release raises flai.minimum no higher than the flai that publishes
- T-1310 flai serve lets a project's requests in flight finish when its manifest stops loading
- T-1311 Document when a publish raises flai.minimum and what the host does when a project stops loading
- T-1312 Close I-0107 with what fixed it

## Notes

Cost of delay inputs set by flai from I-0107. time_lost_per_cycle 1h: 30m per occurrence × 2 occurrences ÷ 1 cycle of 168h (first reported 2026-10-07T07:39:30Z, 0.5 days before this story; under one cycle counts as one).

The separate defect I-0107 names, the duplicated `## 1.0.67` heading in `template/CHANGELOG.md`, was already fixed on `main` by 7772bd34 (`TestBumpMergesIntoAChangelogSectionTheVersionAlreadyHas`), so no issue records it and this story leaves it out.

### Planning

Direction taken from the issue's options: the planner assumed it, and the story's agent may change it. The publish raises `flai.minimum` only to a flai release no newer than the flai running the publish (T-1309). A raise it cannot make yet waits for the first publish from the upgraded flai, so a publish can never drop the host that serves it. The host also lets a dropped project's requests in flight finish (T-1310), so a minimum raised by other means, such as a pull, no longer kills a `publish.run`. Serving a project below its minimum was not chosen: the CLI and `flai mcp` refuse it all the same.

Touches, all files, no folder:

| Touch | Source |
|-------|--------|
| `flai/internal/release/release.go`, `flai/internal/release/release_test.go` | layout: `RaiseMinimum` and `TestRaiseMinimum` |
| `flai/cmd/release.go` | layout: `computeApplyAndTagPending` calls `RaiseMinimum` and prints the warning |
| `flai/internal/serve/serve.go`, `flai/internal/serve/serve_test.go` | layout: `reconcile` drops the project with `running.halt` |
| `flai/internal/channel/channel.go` | layout: the client whose context cancels requests in flight |
| `design/system/flai-cli.md` | design: § Versions: the host's flai and the tree describes the raise |
| `design/system/project-manifest.md`, `docs/users/flai.md`, `docs/operators/index.md` | design and layout: each describes when `flai.minimum` rises |
| `design/issues/I-0107-….md`, `design/issues/summary.md` | goal: criterion 2, `flai issue close` |

`flai touches suggest` listed only hub documents that change with most commits, such as `design/system/flaiover-dashboard.md` at 15%, and none of them bears on this story, so none was taken. An ADR, if the story's agent records the new rule in one, is a new file no task can name yet; its touch is left to the agent's widening rather than kept as a folder.

Forecast: `flai forecast` gave 17m (73 s per unit of size over 17 done large-band remediation stories, times size 14). Raised to 45m: two independent code changes, the serve one with no existing test that drives a manifest which stops loading mid-request, plus four documents and the close-out's full run. The delivery is flai's 2026-10-08T07:21Z moved by the 28m added.

Cost of delay: `flai cod` gives 150 USD a week from the inputs flai set from I-0107 (1h lost per 168h cycle at 150 USD an hour). It stands, worked from those inputs as they are.
