---
id: T-2477
type: task
nature: feature
title: "The kluster package deploys a story: it builds in the story's worktree, loads the images, runs the deploy command in the story's namespace, logs each step, and records the deployment under .flai-cache"
status: backlog
parent: S-0578
owner: alex
created: 2026-11-06T21:43:18Z
updated: 2026-11-06T21:43:18Z
transitions: []
stream: S-0578
tags: [kluster, flai]
touches: [flai/internal/kluster/deploy.go, flai/internal/kluster/deploy_test.go]
---
# T-2477 The kluster package deploys a story: it builds in the story's worktree, loads the images, runs the deploy command in the story's namespace, logs each step, and records the deployment under .flai-cache

## Work

Read S-0576's ADR and `design/system/kluster.md` for the namespace's name, the variables, and how images are loaded.

Add `flai/internal/kluster/deploy.go`. A deploy, given the project, the block, a story, and its worktree:

1. refuses a story not in progress or in review, and a second deploy while one runs for the project, by a lock file under `.flai-cache/kluster/`;
2. runs the block's build command in the worktree, with the variables the ADR names (the story, the namespace, the image tag from the story and its head commit);
3. loads each image the build names into the cluster, as the ADR decides;
4. runs the deploy command with the namespace, the story, the image tag, and the kube context.

Each step's output goes to its own log file under `.flai-cache/kluster/logs/`. A failing step stops the deploy with the step's name and exit code. The outcome is a record, written to `.flai-cache/kluster/deployments.json` with the atomic file helper: story, namespace, commit, images, start and end times, outcome, failing step, and the log paths. The record is JSON with stable field names, since the dashboard reads it.

Run commands through `flai/internal/execx`. Test in `deploy_test.go` with fake `docker`, `kluster`, build, and deploy scripts: a deploy that passes, each step failing, a story in the wrong state, a deploy already running, and a redeploy at a new commit replacing the story's record.

This is the first layer: it waits for no task.

## Done when

- A deploy runs the three steps in order, logs each, and records the outcome, with every case above tested.
- No test needs a real cluster or Docker.
- `flai test` on the two paths passes.

## Notes
