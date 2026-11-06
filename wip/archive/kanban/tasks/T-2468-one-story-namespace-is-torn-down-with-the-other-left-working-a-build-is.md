---
id: T-2468
type: task
nature: experiment
title: One story namespace is torn down with the other left working, a build is promoted to the default namespace, and the results and recommendation name what kluster lacks
status: done
parent: S-0575
owner: alex
created: 2026-11-06T21:41:42Z
updated: 2026-11-06T23:16:36Z
transitions:
  - to: ready
    at: 2026-11-06T23:08:05Z
    by: agent-S-0575
  - to: in-progress
    at: 2026-11-06T23:08:06Z
    by: agent-S-0575
  - to: done
    at: 2026-11-06T23:16:36Z
    by: agent-S-0575
stream: S-0575
tags: [kluster, experiment]
touches: [design/experiments/S-0575-two-builds-of-authstar-run-side-by-side-in-one-kluster-each-in-a-namespace.md, design/experiments/README.md]
after: [T-2467]
usage:
  source: log
  seconds: 510
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 78
      output: 23350
      cache_read: 10775618
      cache_write: 76025
      cost: 3.2306
---
# T-2468 One story namespace is torn down with the other left working, a build is promoted to the default namespace, and the results and recommendation name what kluster lacks

## Work

It waits for T-2467, whose two namespaces it tears down and promotes from.

Delete one story namespace, and its vhost, database, and keys. Run the golden path again against the other, and check that the shared services still answer. Deploy one build into the default `authstar` namespace, as an acceptance would promote it, and run the golden path there. Record the time of the teardown and the promotion. Then take the cluster down with `kluster down`.

Finish the document. Under `## Results`, set every measurement against the success measure. Under `## Recommendation`, write adopt, adapt, or drop, and name S-0576 as the story whose ADR carries it on. List what kluster lacks for this, each with the workaround used, such as an image import, a deploy into a namespace, `status --json`, and a host mapping for demos. List what each project must supply flai for a namespace of its own: its build, its deploy, its teardown, and its per-namespace settings. Bump `updated`.

## Done when

- The teardown and promotion outcomes and times are recorded, and the cluster is down.
- `## Results` and `## Recommendation` are written, with kluster's gaps and what a project must supply each listed.
- `flai check --strict` and the markdown lint pass on the results document and the experiments index.

## Notes
