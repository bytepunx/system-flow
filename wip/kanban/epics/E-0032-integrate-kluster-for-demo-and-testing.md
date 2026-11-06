---
id: E-0032
type: epic
nature: feature
title: Integrate kluster for demo and testing
status: ready
owner: alex
created: 2026-11-06T21:29:53Z
updated: 2026-11-06T21:36:25Z
transitions:
  - to: ready
    at: 2026-11-06T21:30:04Z
    by: system-flow
tags: []
cost_of_delay:
  inputs:
    revenue_per_week: 1250
    by: alex
    at: 2026-11-06T21:36:25Z
---
# E-0032 Integrate kluster for demo and testing

## Outcome

For projects like signet (<https://github.com/bytepunx/signet>, `~/git/bytepunx/signet`) authstar (<https://github.com/bytepunx/authstar>, `~/git/bytepunx/authstar`), a local kubernetes managed through kluster is necessary for the integration and/or smoke tiers of their testing in order to ensure the system work as intended after a given changeset.

Since we control all 3 of these projects, integrating kluster support into flai and flaiover under an Ops submenu in the rail should come with sub menus that lead to pages for managing klusters for both testing and demoing the software without needing to wait for a full build/deploy loop.

Lets build out the stories necessary for an initial pass at this. The Features page should include a toggle (default off) for controlling the configuration, provisioning, and management of a kluster instance once a story has reached the review step. This would allow for stories to deploy a locally built docker image containing changes or operations automation to get tested against a kluster instance.

We may need to issue one or more research and/or experiment stories that can help determine the best way to handle multiple stories (i.e. a story's serial could provide the namespace for a given build out for testing and then once accepted, it could be merged into main then re-deployed to the default namespace). This approach would allow us to use a single set of shared infrastructure (i.e. databases, helm charts, etc.) and then use story-specific configuration (name-spaces, tenants re: authstar, virtual hosts re: rabbitmq, etc.).

## Stories

## Notes
