---
id: E-0032
type: epic
nature: feature
title: Integrate kluster for demo and testing
status: in-progress
owner: alex
created: 2026-11-06T21:29:53Z
updated: 2026-11-06T21:57:39Z
transitions:
  - to: ready
    at: 2026-11-06T21:30:04Z
    by: system-flow
  - to: in-progress
    at: 2026-11-06T21:57:39Z
    by: agent-S-0575
tags: []
usage:
  source: sum
  seconds: 5513
  priced_by: harness
  turns:
    - day: 2026-11-06
      ceremony: 3
      hand_edits: 3
      work: 227
    - day: 2026-11-07
      ceremony: 1
      hand_edits: 2
      work: 17
  models:
    - model: claude-haiku-5-5
      input: 10
      output: 2794
      cache_read: 85448
      cache_write: 33283
      cost: 0.0064
    - model: claude-opus-5-5
      input: 538
      output: 169898
      cache_read: 70180782
      cache_write: 688899
      cost: 22.9007
  strategic:
    - kind: planner
      seconds: 1128
      estimated: true
      models:
        - model: claude-haiku-5-5
          input: 32
          output: 17102
          cache_read: 765943
          cache_write: 83369
          cost: 0.0266
        - model: claude-opus-5-5
          input: 180
          output: 85118
          cache_read: 14099397
          cache_write: 301635
          cost: 6.8042
    - kind: orchestrator
      seconds: 608
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 191
          output: 1233
          cache_read: 7049870
          cache_write: 902141
          cost: 1.906
        - model: claude-sonnet-5-5
          input: 20
          output: 163
          cache_read: 238044
          cache_write: 84467
          cost: 0.229
cost_of_delay:
  inputs:
    revenue_per_week: 1250
    by: alex
    at: 2026-11-06T21:36:25Z
  value: 1250
  by: planner-E-0032
  at: 2026-11-06T21:50:29Z
---
# E-0032 Integrate kluster for demo and testing

## Outcome

For projects like signet (<https://github.com/bytepunx/signet>, `~/git/bytepunx/signet`) authstar (<https://github.com/bytepunx/authstar>, `~/git/bytepunx/authstar`), a local kubernetes managed through kluster is necessary for the integration and/or smoke tiers of their testing in order to ensure the system work as intended after a given changeset.

Since we control all 3 of these projects, integrating kluster support into flai and flaiover under an Ops submenu in the rail should come with sub menus that lead to pages for managing klusters for both testing and demoing the software without needing to wait for a full build/deploy loop.

Lets build out the stories necessary for an initial pass at this. The Features page should include a toggle (default off) for controlling the configuration, provisioning, and management of a kluster instance once a story has reached the review step. This would allow for stories to deploy a locally built docker image containing changes or operations automation to get tested against a kluster instance.

We may need to issue one or more research and/or experiment stories that can help determine the best way to handle multiple stories (i.e. a story's serial could provide the namespace for a given build out for testing and then once accepted, it could be merged into main then re-deployed to the default namespace). This approach would allow us to use a single set of shared infrastructure (i.e. databases, helm charts, etc.) and then use story-specific configuration (name-spaces, tenants re: authstar, virtual hosts re: rabbitmq, etc.).

## Stories
- S-0575 Two builds of authstar run side by side in one kluster, each in a namespace named for a story over the shared signet, RabbitMQ, and Postgres, and the results say what a namespace per story needs
- S-0576 An ADR settles how flai runs a project's kluster for a story in review: one shared cluster, a namespace per story, the story's own settings, promotion on acceptance, and the manifest block a project fills in
- S-0577 flai kluster up, down, and status run a project's kluster from the manifest's kluster block, flai check validates the block, and a kluster host action, off by default, appears on the Features page
- S-0578 flai kluster deploy builds a story's image, loads it into the project's kluster, and runs the project's deploy command in the story's namespace, and teardown and promote remove it and deploy main to the default namespace
- S-0579 With the kluster host action on, flai serve deploys each story entering review into its namespace, runs the review checks against it, tears it down when the story leaves review, and promotes main on acceptance
- S-0580 The rail gains an Ops group whose Klusters page shows the project's kluster and brings it up or takes it down through flai, behind the kluster host action
- S-0581 The Ops group's Deployments page lists main's namespace and each story's, with its commit, outcome, log, and how to reach it for a demo, and deploys, tears down, or promotes through flai

## Notes
- 2026-11-06T21:57:39Z: moved to in-progress: follows S-0575, which moved to in-progress
