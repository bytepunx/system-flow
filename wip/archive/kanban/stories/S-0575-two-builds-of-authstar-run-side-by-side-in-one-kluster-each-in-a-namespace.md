---
id: S-0575
type: story
nature: experiment
title: Two builds of authstar run side by side in one kluster, each in a namespace named for a story over the shared signet, RabbitMQ, and Postgres, and the results say what a namespace per story needs
status: done
parent: E-0032
owner: alex
created: 2026-11-06T21:39:06Z
updated: 2026-11-06T23:19:09Z
transitions:
  - to: ready
    at: 2026-11-06T21:56:36Z
    by: orchestrator
  - to: in-progress
    at: 2026-11-06T21:57:39Z
    by: agent-S-0575
  - to: review
    at: 2026-11-06T23:17:17Z
    by: agent-S-0575
  - to: done
    at: 2026-11-06T23:19:09Z
    by: orchestrator
tags: [kluster, experiment]
topics: [kluster, testing, deployment]
touches: [design/experiments/S-0575-two-builds-of-authstar-run-side-by-side-in-one-kluster-each-in-a-namespace.md, design/experiments/README.md, design/issues/I-0265-a-story-s-sub-agents-write-and-monitor-are-refused-every-path-outside-its.md, design/issues/summary.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 4803
  priced_by: harness
  turns:
    - day: 2026-11-06
      ceremony: 1
      hand_edits: 2
      work: 193
  models:
    - model: claude-haiku-5-5
      input: 10
      output: 2794
      cache_read: 85448
      cache_write: 33283
      cost: 0.0064
    - model: claude-opus-5-5
      input: 402
      output: 120506
      cache_read: 55611261
      cache_write: 392353
      cost: 16.6728
  strategic:
    - kind: orchestrator
      seconds: 234
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 63
          output: 402
          cache_read: 2619517
          cache_write: 246774
          cost: 0.6866
        - model: claude-sonnet-5-5
          input: 16
          output: 123
          cache_read: 225689
          cache_write: 52052
          cost: 0.1972
cost_of_delay:
  value: 490.2
  by: planner-E-0032
  at: 2026-11-06T21:49:56Z
forecast:
  duration: 3h
  delivery: 2026-11-07T10:05:00Z
  basis: "Its own forecast of 3h; 3rd in the pull order with an in-progress limit of 5, behind S-0545 and S-0520."
  by: flai
  at: 2026-11-06T21:56:47Z
finalized:
  by: orchestrator
  at: 2026-11-06T21:55:26Z
---
# S-0575 Two builds of authstar run side by side in one kluster, each in a namespace named for a story over the shared signet, RabbitMQ, and Postgres, and the results say what a namespace per story needs

## Goal

E-0032 wants a story in review deployed to a local kluster for its integration or smoke tier, without a full build and deploy loop. Its outcome guesses at one shared cluster, with one namespace per story and the story's own configuration (an authstar tenant, a RabbitMQ vhost) over shared infrastructure. Nothing in kluster, signet, or authstar does that today. kluster installs a fixed set of addons into fixed namespaces, and authstar's `integration/scripts/lib.sh` hard-codes the `authstar` namespace and the `/authstar` vhost.

This experiment tries the idea by hand on authstar before flai automates it. Two builds share one `authstar`-profile kluster, each in its own story namespace. The results say what each namespace needs, what it costs in time, and what kluster lacks. The ADR that S-0576 writes rests on them.

## Acceptance criteria

- [x] One kluster is brought up with authstar's `integration/kluster.yaml` (profile `authstar`). The results record the time from `kluster up` to a cluster whose shared signet, RabbitMQ, Postgres, and ClickHouse answer.
- [x] Two builds of authstar's services, from two different commits of `~/git/bytepunx/authstar`, are built locally, loaded into the cluster, and deployed into two namespaces named as a story's serial would name them (such as `s-0901` and `s-0902`). Each namespace gets its own RabbitMQ vhost, Postgres database, and signet keys or authstar tenant, as far as authstar allows. The results list each per-namespace setting and how it was made.
- [x] authstar's golden-path scenario (`integration/scripts/scenario-golden-path.sh`) is run against each namespace. The results say whether it passed, or which step failed and why.
- [x] Deleting one story namespace leaves the other namespace and the shared infrastructure working. One build is then deployed to the default `authstar` namespace, as an acceptance would promote it. The results record the time of each load, deploy, and teardown.
- [x] A results document under `design/experiments/`, made from `template.md` and named as this story's file is, records the hypothesis, the success measure, what was done (with the commands and scripts used, written out), the results, and a recommendation. It lists what kluster lacks for this, such as an image import, a deploy into a namespace, and machine-readable status. It is indexed in `design/experiments/README.md`.
- [x] Nothing is committed in the authstar, signet, or kluster repositories. The experiment reads them and runs their scripts with variables overridden or with copies kept outside them, and the results say which.

## Tasks

- T-2466 The results document states the hypothesis and success measure, and records one authstar-profile kluster brought up with its shared services answering and the time it took
- T-2467 Two builds of authstar from two commits are loaded into the kluster and deployed into two story namespaces with their own vhost, database, and keys, and the golden path is run against each
- T-2468 One story namespace is torn down with the other left working, a build is promoted to the default namespace, and the results and recommendation name what kluster lacks

## Notes

This story works outside this repository on purpose: it runs `kluster`, `docker`, and authstar's scripts on the host against a local cluster. It may change files only in its own worktree and in a scratch folder.

### Planning

- Touches:
  - `design/experiments/S-0575-two-builds-of-authstar-run-side-by-side-in-one-kluster-each-in-a-namespace.md` comes from the design: `design/experiments/template.md` names a results document after its story's file.
  - `design/experiments/README.md` comes from the layout: it is the folder's index.
  - `flai touches suggest` found nothing changed with them often enough (3 of 2956 commits).
- `after`: none. The story reads nothing another story writes. It reads the authstar, signet, and kluster repositories as they are.
- Tasks: three layers, one task each, since all three write the one results document. T-2466 brings the cluster up. T-2467 deploys into it. T-2468 tears down, promotes, and concludes.
- Forecast 3h, raised from `flai forecast`'s 23m. The median comes from experiment stories with no outside infrastructure. This one waits on wall-clock work: `kluster up` with the authstar profile (SPIRE, signet, CockroachDB, RabbitMQ, Dex, Postgres, ClickHouse), two builds of eight service images, and three golden-path runs. The delivery moves by the same 2h37m.
- Cost of delay 490.20 USD/week: `flai cod`'s share of E-0032's 1250 USD/week revenue input (alex, TH-0829), by forecast duration over the epic's seven open stories. It is the largest share only because this story's forecast is the longest. The rule is the project's, so the figure stands as flai gives it.

### Accepted by the orchestrator

- Verified: f7c8f7a4474361c3d9f4acbe2d33bdf94ceade08
- At: 2026-11-06T23:19:09Z

Verdict: pass. Every criterion is met by the results document; flai verify passed every step at f7c8f7a; authstar, signet, and kluster have no commits since the story started.
- 1: design/experiments/S-0575-two-builds-of-authstar-run-side-by-side-in-one-kluster-each-in-a-namespace.md
- 2: design/experiments/S-0575-two-builds-of-authstar-run-side-by-side-in-one-kluster-each-in-a-namespace.md
- 3: design/experiments/S-0575-two-builds-of-authstar-run-side-by-side-in-one-kluster-each-in-a-namespace.md
- 4: design/experiments/S-0575-two-builds-of-authstar-run-side-by-side-in-one-kluster-each-in-a-namespace.md
- 5: design/experiments/S-0575-two-builds-of-authstar-run-side-by-side-in-one-kluster-each-in-a-namespace.md, design/experiments/README.md
- 6: design/experiments/S-0575-two-builds-of-authstar-run-side-by-side-in-one-kluster-each-in-a-namespace.md
