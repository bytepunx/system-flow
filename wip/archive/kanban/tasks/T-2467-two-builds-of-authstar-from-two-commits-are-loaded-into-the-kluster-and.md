---
id: T-2467
type: task
nature: experiment
title: Two builds of authstar from two commits are loaded into the kluster and deployed into two story namespaces with their own vhost, database, and keys, and the golden path is run against each
status: done
parent: S-0575
owner: alex
created: 2026-11-06T21:41:30Z
updated: 2026-11-06T23:07:44Z
transitions:
  - to: ready
    at: 2026-11-06T22:17:36Z
    by: agent-S-0575
  - to: in-progress
    at: 2026-11-06T22:17:39Z
    by: agent-S-0575
  - to: done
    at: 2026-11-06T23:07:44Z
    by: agent-S-0575
stream: S-0575
tags: [kluster, experiment]
touches: [design/experiments/S-0575-two-builds-of-authstar-run-side-by-side-in-one-kluster-each-in-a-namespace.md]
after: [T-2466]
usage:
  source: log
  seconds: 3005
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 194
      output: 58173
      cache_read: 26845749
      cache_write: 189404
      cost: 8.0486
---
# T-2467 Two builds of authstar from two commits are loaded into the kluster and deployed into two story namespaces with their own vhost, database, and keys, and the golden path is run against each

## Work

It waits for T-2466, whose cluster it deploys into and whose results document it writes in.

Pick two commits of `~/git/bytepunx/authstar` and check each out into a scratch worktree outside the repository (`git worktree add` into the scratch folder, removed afterwards), so authstar's own checkout is untouched. Build each service's image with its `Dockerfile`, tagged with the namespace (such as `s-0901` and `s-0902`). Load the images into the cluster as `integration/scripts/deploy-local.sh` does (`docker save` piped to `ctr images import` in the k3d server node), and record whether `k3d image import` would serve as well.

For each namespace, make what authstar needs apart from the shared infrastructure: the namespace, its RabbitMQ vhost and user, its Postgres database, and its signet keys or tenant. Use copies of `integration/scripts/lib.sh`, `lib-secrets.sh`, `provision.sh`, and `manifests/*.yaml` in the scratch folder, with `AUTHSTAR_NAMESPACE` and the vhost changed. Deploy each build. Run `integration/scripts/scenario-golden-path.sh` against each namespace.

Record under `## What was done` each per-namespace setting and how it was made, every changed line of the copied scripts, and the time of each image load and deploy. Record each golden-path outcome, or the failing step and why, under `## Results`.

## Done when

- Two namespaces run two builds side by side, and the document lists their settings, the script changes, and the times.
- The golden path's outcome in each namespace is recorded.
- authstar's own checkout has no change, and the scratch worktrees are removed.

## Notes
