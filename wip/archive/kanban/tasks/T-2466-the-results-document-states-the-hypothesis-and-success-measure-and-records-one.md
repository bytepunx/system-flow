---
id: T-2466
type: task
nature: experiment
title: The results document states the hypothesis and success measure, and records one authstar-profile kluster brought up with its shared services answering and the time it took
status: done
parent: S-0575
owner: alex
created: 2026-11-06T21:41:17Z
updated: 2026-11-06T22:16:32Z
transitions:
  - to: ready
    at: 2026-11-06T22:02:12Z
    by: agent-S-0575
  - to: in-progress
    at: 2026-11-06T22:02:13Z
    by: agent-S-0575
  - to: done
    at: 2026-11-06T22:16:32Z
    by: agent-S-0575
stream: S-0575
tags: [kluster, experiment]
touches: [design/experiments/S-0575-two-builds-of-authstar-run-side-by-side-in-one-kluster-each-in-a-namespace.md, design/experiments/README.md]
usage:
  source: log
  seconds: 859
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 57
      output: 17209
      cache_read: 7941511
      cache_write: 56030
      cost: 2.3809
---
# T-2466 The results document states the hypothesis and success measure, and records one authstar-profile kluster brought up with its shared services answering and the time it took

## Work

Copy `design/experiments/template.md` to `design/experiments/S-0575-two-builds-of-authstar-run-side-by-side-in-one-kluster-each-in-a-namespace.md` and index it in `design/experiments/README.md`.

Write `## Hypothesis`: one shared kluster can hold several builds of authstar, each in a namespace named for a story, over one set of shared signet, RabbitMQ, Postgres, and ClickHouse, each with its own vhost, database, and keys. Write `## Success measure`: both namespaces pass the golden path, deleting one leaves the other working, and a deploy into a namespace takes minutes, not the length of a full build and deploy loop.

Install `kluster` into a scratch folder if it is missing, as `safety.md` says, and record the version. Bring up a cluster from `~/git/bytepunx/authstar/integration/kluster.yaml` with `kluster --config <path> up`, then `kluster use`. Do not edit authstar's file: copy it to the scratch folder if a setting must change. Check that signet, RabbitMQ, Postgres, and ClickHouse answer in their namespaces with `kubectl`. Record the commands and the time from `kluster up` to all four answering under `## What was done`.

This is the first layer: it waits for no task.

## Done when

- The results document exists from the template, with its hypothesis and success measure written, and is indexed.
- A cluster runs with the four shared services answering, and the document records the commands, the kluster version, and the time.
- Nothing is changed in the authstar, signet, or kluster repositories.

## Notes
