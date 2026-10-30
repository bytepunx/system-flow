---
id: S-0548
type: story
nature: remediation
title: "A queue verification takes under five minutes: the mcpserver tests are split or parallelized, the timing-dependent serve tests wait for what they check, and every go test tier sets its timeout"
status: backlog
parent: E-0031
owner: alex
created: 2026-10-30T01:52:08Z
updated: 2026-10-30T01:52:08Z
transitions: []
tags: [flai]
topics: [testing]
touches: [flai/internal/mcpserver, flai/internal/serve, flai/internal/dashboard, system-flow.yaml, ".github/workflows/flai.yml", design/tech/go-libraries.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
---
# S-0548 A queue verification takes under five minutes: the mcpserver tests are split or parallelized, the timing-dependent serve tests wait for what they check, and every go test tier sets its timeout

## Goal

A merge-queue verification took 10.4 minutes at the median, 15.2 hours in the fortnight, and some of its failures were time alone: the mcpserver package runs 8 to 10 minutes under the race detector against Go's 10-minute default (I-0225, I-0227), two serve tests check before what they wait for has happened (I-0226), the dashboard notification test fails its temp folder's cleanup (I-0233), and the macOS CI run exceeded the default on 2026-10-25. This story brings a verification under five minutes and makes it fail only on defects: the mcpserver tests share their fixtures or run in parallel with `t.Parallel`, the serve and dashboard tests wait for the condition they assert, and every go test tier in the manifest and in CI names its timeout so a slow host fails with the package's name rather than a panic.

## Acceptance criteria
- [ ] `go test -race ./internal/mcpserver` finishes under three minutes on the host, and the package's slowest tests are named in the story's close-out record
- [ ] The tests I-0226 and I-0233 name pass ten times in a row under `-race -count=10` with `GOMAXPROCS=2`, and the issues are closed with the reason
- [ ] Every go test tier in system-flow.yaml and both `go test` lines in the CI workflow carry an explicit `-timeout`, and `flai check` warns of a go test tier without one
- [ ] The fortnight after landing shows verification p50 under five minutes in `flai stats`, recorded in the close-out, and I-0225 and I-0227 are closed

## Tasks

The planner drafts them.

## Notes

Evidence: `flai stats` merge-queue figures and the CI run 40921428484 fixed by hand on 2026-10-29; the race fixes from that day are in the main checkout uncommitted.
