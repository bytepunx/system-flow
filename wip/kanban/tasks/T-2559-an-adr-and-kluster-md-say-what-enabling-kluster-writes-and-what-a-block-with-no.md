---
id: T-2559
type: task
nature: improvement
title: An ADR and kluster.md say what enabling kluster writes and what a block with no deploy means
status: backlog
parent: S-0586
owner: alex
created: 2026-11-07T08:29:22Z
updated: 2026-11-07T08:29:40Z
transitions: []
stream: S-0586
tags: [kluster, docs]
touches: [design/adrs/README.md, design/system/kluster.md, design/system/project-manifest.md]
usage:
  source: log
  seconds: 0
  estimated: true
  priced_by: estimate
  models:
    - model: claude-haiku-5-5
      input: 10
      output: 41
      cache_read: 106965
      cache_write: 36155
      cost: 0.0103
---
# T-2559 An ADR and kluster.md say what enabling kluster writes and what a block with no deploy means

## Work

Write the ADR that refines ADR-0254, with `flai adr new`, and change `design/system/kluster.md` and `design/system/project-manifest.md` to match. ADR-0254 is accepted, so it stays as it is. The ADR decides:

- Enabling the `kluster` host action on a project with no `kluster:` block writes one, and a kluster config file, and commits both, in place of the refusal `kluster.ErrNoBlock` gives today.
- The block it writes names only `config`. `deploy` becomes optional: a block with no `deploy` is a project whose cluster can come up and go down, and whose deploys are not set up yet. This replaces ADR-0254's `deploy` required.
- The config file is `kluster.yaml` at the repository root, since the layout allows no new top-level folder. It names the cluster, as the manifest's name made a DNS label with `-local` after it, and the provider `k3d`, and no profile: the empty cluster the story's goal asks for.
- Every key of the block, and the config file, is changed from the dashboard. Those writes stay under the `settings` host action, as the Features toggle does (ADR-0039), because a dashboard that sets these commands runs them on the host as the operator.
- What the operator's answers on this story's plan thread settle about the host's tools and the port-forwards. Read that thread first.

Check kluster's config against the real command before the ADR says a config with no profile is an empty cluster: run `kluster --config <a copy> up` where `kluster` is on `PATH`. The planner could read neither kluster's README nor a `kluster` binary. If `kluster` is not on `PATH` where you run, or it refuses a config with no profile, ask on a thread on the story, naming what you tried.

This task waits for nothing: the tasks after it build what it decides.

## Done when

- The ADR is accepted, names ADR-0254 as what it refines, and `flai check --strict` is clean.
- `design/system/kluster.md` § The `kluster:` block, § The `kluster` host action, and § Refusals, and `design/system/project-manifest.md`'s `kluster` entry, say what the ADR decides, and no longer say `deploy` is required or that enabling is refused for a project with no block.
- The narrative's `## Decisions` records what `kluster` did with a config that names no profile, or the thread that asks.

## Notes
