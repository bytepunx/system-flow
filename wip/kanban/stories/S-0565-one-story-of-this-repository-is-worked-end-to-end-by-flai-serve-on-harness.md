---
id: S-0565
type: story
nature: experiment
title: One story of this repository is worked end to end by flai serve on harness container, from story_start to review, once S-0508's container harness is in the installed flai
status: backlog
parent: E-0029
owner: alex
created: 2026-11-02T18:08:56Z
updated: 2026-11-02T18:08:56Z
transitions: []
tags: []
topics: [harness, docker]
after: [S-0508]
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
      seconds: 11
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 3
          output: 25
          cache_read: 57078
          cache_write: 32925
          cost: 0.0212
draft: true
---
# S-0565 One story of this repository is worked end to end by flai serve on harness container, from story_start to review, once S-0508's container harness is in the installed flai

## Goal

Prove the container harness S-0508 built in the place it is meant for: `flai serve` on this host, the installed release, starts one story of this repository on `harness: container` and its agent works it from `story_start` to review inside the container, with its usage measured and its card's status shown as for a host run. Each limitation met becomes an issue.

## Acceptance criteria

- [ ] One story of this repository is worked end to end by `flai serve` with `harness: container` on its agent, from `story_start` to review, its usage measured and its card's status shown as for a host run, and each limitation met is an issue under `design/issues`
- [ ] A stop from the dashboard and a restart after a thread is answered are each tried once on that story's container, and what happened is recorded in this story's notes

## Tasks

## Notes

Moved out of S-0508 (its criterion 6 and the live-run half of T-2169) on TH-0812, as D-0461 records for S-0423. `after: S-0508`: the `container` adapter, its host settings (`flai serve agent harness container --image/--network/--login`), and the serve start path exist only on `story/S-0508`; they reach this host's `flai serve`, the installed release, only once S-0508 is accepted and published and the host upgraded.

The work, from T-2169: with the host settings in place (`flai serve agent harness container --login ~/.claude`) and the images built (`make agent-images`), set `harness: container` on one small story of this repository (`flai edit S-nnnn --harness container`), and let `flai serve` start it. The proof story should change documentation or shell only: the Go tiers and the close-out need the project's image layer (S-0510, S-0520), and git does not work in the container yet (S-0517). Record each limitation met with `flai issue new`, naming the story of E-0029 that carries it on; I-0253 and I-0254 are S-0508's already.
