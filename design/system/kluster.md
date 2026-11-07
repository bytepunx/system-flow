---
title: Running a project's kluster
updated: 2026-11-06
status: active
topics: [kluster, deployment, manifest]
---

# Running a project's kluster

How flai runs a project's local Kubernetes cluster, made by `bytepunx/kluster`, so that each story in review is deployed into a namespace of its own over the project's shared infrastructure, checked there, and promoted to the default namespace on acceptance. [ADR-0254](../adrs/0254-flai-runs-one-kluster-per-project-from-the-config-its-manifest-s-kluster-block.md) decides it and gives the reasons. [S-0575's results](../experiments/S-0575-two-builds-of-authstar-run-side-by-side-in-one-kluster-each-in-a-namespace.md) show it working by hand on authstar, with the times of each step.

system-flow itself runs no cluster. This is what flai does for a project that names one in its manifest, behind the `kluster` host action.

## The project's cluster

- A project has at most one cluster. It is made from the kluster config file `kluster.config` names, as that file stands on the main branch. A story that changes the file changes nothing in the running cluster.
- The cluster's name is the config's `name`. Its provider is `k3d`; `flai check` reports any other.
- The shared infrastructure is what the config's profile installs. For authstar's `integration/kluster.yaml`, profile `authstar`, it is SPIRE, signet, RabbitMQ, Postgres, ClickHouse, and Dex. flai installs nothing in the cluster and changes no shared service's own settings.
- flai brings the cluster up with `kluster --config <config> up` and takes it down with `kluster down`, in the main checkout. A bring-up took 11m31s in S-0575.
- flai writes the cluster's kubeconfig to `.flai-cache/kluster/kubeconfig` with `kluster kubeconfig --name <name> --output`, and passes it to every command as `KUBECONFIG`. It never merges it into `~/.kube/config`. The context is `k3d-<name>`.

What flai keeps, all under `.flai-cache/kluster/`:

| Path | What it holds |
|------|---------------|
| `kubeconfig` | The cluster's kubeconfig |
| one record per namespace | Its commit, the outcome and times of each step, and the job's log |
| `state/<namespace>/` | The namespace's state folder, `FLAI_KLUSTER_STATE` |
| `main/` | A worktree of main's head, where a promotion builds and deploys |

## The `kluster:` block

The block in `system-flow.yaml` is optional. Without it, flai runs no cluster for the project.

```yaml
kluster:
  config: integration/kluster.yaml
  namespace: authstar
  build: [scripts/kluster-build.sh]
  deploy: [scripts/kluster-deploy.sh]
  teardown: [scripts/kluster-teardown.sh]
  timeout: 30m
```

| Key | Required | Default | Meaning |
|-----|----------|---------|---------|
| `config` | yes | none | The kluster config file, from the repository root |
| `namespace` | no | the manifest's `name` as a DNS label | The default namespace, which main deploys into |
| `build` | no | none: nothing is loaded | Builds the images and lists their references in the file `FLAI_KLUSTER_IMAGES` names |
| `deploy` | yes | none | Makes the namespace's settings in the shared services, applies the project's manifests into the namespace, and waits for their rollouts |
| `teardown` | no | none: flai deletes only the namespace | Removes what `deploy` made outside the namespace |
| `timeout` | no | `30m` | The most one job may take, a Go duration |

Each command is an argument list, run as it stands, never through a shell, as the manifest's `tests` are. It runs in the checkout being deployed, with the variables below in its environment. A command that exits other than 0 fails its job. `deploy` is idempotent: a second deploy into the same namespace, at a new commit, updates it in place.

## The namespaces

| Namespace | Who deploys into it | Name | Prefix |
|-----------|---------------------|------|--------|
| A story's | The story, while in review | Its ID in lower case, with the zero padding: `s-0576` | `s0576_`: the namespace without hyphens, then an underscore |
| The default | Main, on each acceptance | `kluster.namespace`, such as `authstar` | Empty, so main keeps the names the profile made |

- `flai check` reports a default namespace that is not a DNS label, or that has a story namespace's form, `s-` and digits.
- The prefix is for the names a namespace makes in shared services, such as databases: `s0576_tower`.
- What a story owns in the shared services is the project's to make and remove: its RabbitMQ vhost and user, its databases, its signet bundles and policies, its keys and tokens, and its tenant where the project has tenants, as authstar has. `deploy` makes them from the namespace and the prefix, and `teardown` removes them. flai does not know them.
- A project's Kubernetes manifests and scripts take the namespace as a parameter, with the default namespace as its default, and never name it as a constant. authstar's manifests still name `authstar`, which S-0575 renamed with `sed`; that is authstar's to change.

## The steps and their variables

A deploy is the first three steps, and a teardown the fourth:

1. **Build.** `build` builds the project's images in the checkout and writes each reference, one a line, to `FLAI_KLUSTER_IMAGES`.
2. **Load.** flai loads each listed image into every node of the cluster: `docker save <image>` piped into `docker exec -i <node> ctr -n k8s.io images import -`, for each container labelled `k3d.cluster=<name>` with `k3d.role` `server` or `agent`. Once kluster has an image import, flai calls that instead.
3. **Deploy.** `deploy` makes the namespace's settings and applies the manifests, with the images tagged `FLAI_KLUSTER_TAG`.
4. **Teardown.** `teardown` removes the namespace's settings in the shared services; then flai deletes the namespace with `kubectl delete namespace --wait`. A namespace already gone is torn down.

Every command, and a checks run against a namespace, gets these variables. No placeholder is replaced in their arguments.

| Variable | A story's namespace | The default namespace |
|----------|---------------------|-----------------------|
| `FLAI_STORY` | `S-0576` | empty |
| `FLAI_KLUSTER_NAME` | The cluster's name, such as `authstar-integration` | The same |
| `FLAI_KLUSTER_NAMESPACE` | `s-0576` | `authstar` |
| `FLAI_KLUSTER_PREFIX` | `s0576_` | empty |
| `FLAI_KLUSTER_COMMIT` | The story branch's head, the whole hash | Main's head |
| `FLAI_KLUSTER_TAG` | `s-0576-` and the head's first 12 hex digits | `main-` and the head's first 12 hex digits |
| `FLAI_KLUSTER_IMAGES` | `build` only: the file it lists its images in | The same |
| `FLAI_KLUSTER_STATE` | `.flai-cache/kluster/state/s-0576` | `.flai-cache/kluster/state/authstar` |
| `KUBECONFIG` | flai's kubeconfig for the cluster | The same |
| `FLAI_KLUSTER_CONTEXT` | `k3d-` and the cluster's name | The same |

- `FLAI_KLUSTER_STATE` is where `deploy` keeps what the checks need, such as a namespace's operator tokens. It is git-ignored, holds only values that work against the local cluster, and is deleted after the teardown.
- flai gives no shared service's address. The project's scripts know their profile's.

## What happens when

While the `kluster` host action is on for a project with a `kluster:` block, `flai serve` runs these jobs, with no model:

| When | Job | Where it runs |
|------|-----|---------------|
| A story enters review | Deploy into its namespace | The story's worktree, at its head |
| A story in review gets a new head | Deploy again, in place | The story's worktree, at its new head |
| A story leaves review for in-progress or cancelled, not by acceptance or a failure | Teardown | The story's worktree, or the main checkout when it is gone |
| A story is sent back by a failed deploy or checks run | None: its namespace stays for diagnosis | Nowhere |
| A story sent back returns to review | Deploy again, in place | The story's worktree |
| A story is accepted | Teardown of its namespace, then a promotion | The teardown as above; the promotion in `.flai-cache/kluster/main` |
| A deploy finds the cluster absent | Bring the cluster up first | The main checkout |

- A promotion builds, loads, and deploys main's head into the default namespace. There `deploy` replaces what the profile deployed. Promotions queued behind one another run once, at main's newest head.
- Jobs run one at a time per project, in the order their causes came.
- A story never deployed is not torn down.
- `flai serve` never takes the cluster down.
- The close-out and `flai verify` do not use the cluster: they run before review, when the story has no namespace.

## A failed deploy holds the story

- A story's failed deploy, or a failed checks run against its namespace, is recorded. `flai serve` opens one thread to the operator on the story, naming the job, the step, the commit, and the log, and sends the story back to in-progress, starting its agent. The agent fixes the change, or says on the thread that the cause is outside the story.
- While the action is on, every acceptance of a story is refused with the blocker `kluster`, unless the story's namespace was deployed at its head and its checks run against it passed there, or `checks` is off. This holds for the operator's acceptance, the orchestrator's, and the merge queue's. The blocker names the last job and its outcome.
- `flai accept --without-kluster` lands the story anyway. It is the operator's alone, and it is recorded in the story's Notes.
- A failed promotion has no story to send back. `flai serve` records it and opens one thread to the operator on `system-flow.yaml`, naming the commit, the step, and the log. The next acceptance promotes again.
- Until a project's baseline passes in the default namespace, each of its stories is sent back at its first deploy or checks run.

## The `kluster` host action

| State | What it allows |
|-------|----------------|
| Off, the default | `flai serve` starts no job and leaves every namespace as it is. Turning the action off tears nothing down. |
| On | `flai serve` runs the jobs above, and the dashboard may bring the cluster up or down and deploy, tear down, or promote through flai |

- The operator turns it on with `flai serve enable kluster`, or its toggle on the dashboard's Features page. It is refused for a project with no `kluster:` block.
- `flai kluster up`, `down`, `status`, `deploy`, `teardown`, and `promote` work from the operator's shell whatever the action says.
- `flai guard` refuses every agent and sub-agent the `flai kluster` commands that write. `flai kluster status` is a read any agent may run.

## What the review checks see

- With `kluster` and `checks` both on, `flai serve` starts a story's checks run once its deploy has passed at the story's head, not when the story enters review.
- A checks run for a story whose namespace is deployed at its head gets the variables above for that namespace, beside `{story}` and `{root}` in its arguments. A check such as authstar's golden path reaches the namespace through `KUBECONFIG` and `FLAI_KLUSTER_NAMESPACE`.
- A checks run for a story with no namespace deployed at its head gets none of them.
- Checks runs against the cluster run one at a time per project, after the job queue.
- With `checks` off, a story is deployed and nothing is checked against it.

## The Ops pages in the dashboard

E-0032 gives the dashboard's rail an Ops group with two pages, each behind the `kluster` host action and acting only through flai:

| Page | What it shows | What it does |
|------|---------------|--------------|
| Klusters | The project's cluster: its name, its state, and its shared services | Brings the cluster up or takes it down |
| Deployments | Main's namespace and each story's, with its commit, outcome, log, and how to reach it for a demo | Deploys, tears down, or promotes |

## Commands (S-0577, S-0578)

To be written by S-0577 and S-0578. This section will hold the `flai kluster` commands as built: `up`, `down`, and `status` from S-0577, with what `flai check` reports of the `kluster:` block; and `deploy`, `teardown`, and `promote` from S-0578, with the image load, their flags, their JSON, and their exit codes.

## flai serve (S-0579)

To be written by S-0579. This section will hold how `flai serve` queues and runs the jobs, records their outcomes, opens the threads of a failure and sends the story back, starts the checks run after a deploy with the namespace's variables, and gives `flai accept` its `kluster` blocker.

## The dashboard (S-0580, S-0581)

To be written by S-0580 and S-0581. This section will hold the Ops group's Klusters and Deployments pages as built, the host methods they call, and what each shows while the host action is off.
