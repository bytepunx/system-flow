---
id: ADR-0254
title: "flai runs one kluster per project from the config its manifest's kluster block names, deploys each story in review into a namespace named for its serial through the project's own build, deploy, and teardown commands, and on acceptance tears the story's namespace down and deploys main to the default namespace, behind the kluster host action, off by default"
status: proposed
date: 2026-11-06
supersedes: []
superseded_by: []
---

# ADR-0254 flai runs one kluster per project from the config its manifest's kluster block names, deploys each story in review into a namespace named for its serial through the project's own build, deploy, and teardown commands, and on acceptance tears the story's namespace down and deploys main to the default namespace, behind the kluster host action, off by default

## Context

E-0032 wants a story in review deployed to a local Kubernetes cluster, so that a project such as authstar or signet can run its integration or smoke checks against the change, and the operator can demo it, without a full build and deploy loop. The stories after this one, S-0577 to S-0581, build that into flai and the dashboard. They need one decision to build on.

S-0575's experiment ([results](../experiments/S-0575-two-builds-of-authstar-run-side-by-side-in-one-kluster-each-in-a-namespace.md)) tried the shape E-0032 guesses at, by hand, on authstar with kluster 1.3.1:

- One `kluster up` with authstar's `integration/kluster.yaml` (profile `authstar`, provider `k3d`) brought up SPIRE, signet, RabbitMQ, Postgres, ClickHouse, and Dex in 11m31s.
- Two builds of authstar ran at once in the namespaces `s-0901` and `s-0902`. Each had its own RabbitMQ vhost and user, Postgres and ClickHouse databases with the prefix `s0901_` or `s0902_`, and signet bundles and policies under a signet namespace equal to the Kubernetes namespace. SPIRE named each workload by its namespace with no change. Both passed authstar's golden path.
- A story namespace cost about 90 seconds to deploy and 70 to tear down, plus what its change rebuilt: 41 seconds for one changed service on a warm cache, 18 minutes cold.
- Deleting `s-0901` left `s-0902` and the shared services working. Build A was then promoted to the default namespace, `authstar`, in 87 seconds, and passed there.
- What made a namespace was authstar's own manifests and config with the namespace put in, and settings only authstar knows: which databases, which secrets are keys, which services grant which writes. kluster deploys authstar only into the constant namespace `authstar`, from Docker Hub's `latest`.
- kluster has no image import, no deploy into a namespace, no per-namespace setup or teardown of the shared services, and no machine-readable status. Its signet chart is unpinned, and its Dex does not serve authstar's tenant.
- authstar's checks reach a namespace through fixed local ports and leak their port-forwards, so two namespaces' checks cannot run at once.

The recommendation was to adapt: one cluster per project and a namespace per story, with the deploy the project's and the cluster kluster's.

## Decision

**flai runs one kluster per project from the config its manifest's `kluster:` block names, deploys each story in review into a namespace named for its serial through the project's own build, deploy, and teardown commands, and on acceptance tears the story's namespace down and deploys main to the default namespace, behind the `kluster` host action, off by default.**

### 1. One cluster per project

- A project has at most one kluster. Its config is the kluster config file `kluster.config` names, from the repository root, such as authstar's `integration/kluster.yaml`. flai reads the file as it stands on the main branch: a story that changes it changes nothing in the running cluster.
- The cluster's name is the config's `name`. Its provider must be `k3d`; `flai check` reports any other. kind, kluster's CI provider, is left out until a project needs it.
- The shared infrastructure is what the config's profile installs: for authstar, SPIRE, signet, RabbitMQ, Postgres, ClickHouse, and Dex. flai installs nothing in the cluster itself and never changes a shared service's own settings.
- flai brings the cluster up with `kluster --config <config> up` and takes it down with `kluster down`, run in the main checkout. `flai kluster up` and `down` do so on the operator's word, from a shell or the dashboard's Klusters page. While the `kluster` host action is on, `flai serve` brings the cluster up itself when a deploy finds it absent, and never takes it down.
- flai writes the cluster's kubeconfig with `kluster kubeconfig --name <name> --output .flai-cache/kluster/kubeconfig` and passes it to every command as `KUBECONFIG`. It never merges it into `~/.kube/config`. The context is `k3d-<name>`.
- What flai keeps of the cluster and its deployments is under `.flai-cache/kluster/`: the kubeconfig, one record per namespace with its commit, outcome, step times, and log, and a state folder per namespace (§4).

### 2. The namespaces

- A story deploys into the namespace of its ID in lower case, with the ID's zero padding: S-0576 deploys into `s-0576`. It is a DNS label, the ID is never reused, so a namespace is never a second story's, and the operator reads the story from it.
- Main deploys into the default namespace, `kluster.namespace`. It defaults to the manifest's `name` made a DNS label. For authstar it is `authstar`, the namespace its profile deploys into.
- `flai check` reports a default namespace that is not a DNS label or that has a story namespace's form, `s-` and digits.
- Each namespace has a prefix for the names it makes in shared services, such as databases: the namespace with its hyphens taken out and an underscore after, `s0576_`. The default namespace's prefix is empty, so main keeps the names the profile made.
- The settings a story owns in the shared services, its vhost and user, its databases, its signet bundles and policies, its keys and tokens, are made by the project's deploy command from the namespace and the prefix, and removed by its teardown command. flai does not know them.

### 3. The manifest's `kluster:` block

```yaml
kluster:                                     # optional (S-0577): the project's kluster; absent, flai runs none
  config: integration/kluster.yaml           # required: the kluster config file, from the repository root
  namespace: authstar                        # optional: the default namespace main deploys into; default the manifest's name as a DNS label
  build: [scripts/kluster-build.sh]          # optional: builds the images and lists them; default none, nothing is loaded
  deploy: [scripts/kluster-deploy.sh]        # required: deploys the build into FLAI_KLUSTER_NAMESPACE
  teardown: [scripts/kluster-teardown.sh]    # optional: removes what deploy made outside the namespace; default none
  timeout: 30m                               # optional: the most one job may take, a Go duration; default 30m
```

- Each command is an argument list, run as it stands, never through a shell, as the manifest's `tests` and the `checks` host action run theirs. It runs in the checkout being deployed (§5) with the variables of §4 in its environment. A command that exits other than 0 fails its job.
- `build` builds the project's images and writes the reference of each, one a line, to the file `FLAI_KLUSTER_IMAGES` names. flai loads each (§6). With no `build`, flai loads nothing: the deploy then pulls what it needs or builds it itself.
- `deploy` makes the namespace's settings in the shared services and applies the project's manifests into the namespace, with the images tagged `FLAI_KLUSTER_TAG`, and waits for their rollouts. It must be idempotent: a second deploy into the same namespace, at a new commit, updates it in place. In the default namespace it replaces what the profile deployed, as S-0575's promotion did by deleting the profile's Deployments and CronJobs first.
- `teardown` removes what `deploy` made outside the namespace: its vhost and user, its databases, its signet bundles and policies. flai then deletes the namespace itself, with `kubectl delete namespace --wait`. With no `teardown`, flai deletes only the namespace.

### 4. The variables

flai gives `build`, `deploy`, `teardown`, and a checks run against a namespace (§8) these variables, and no placeholders in their arguments:

| Variable | A story's namespace | The default namespace |
|----------|---------------------|-----------------------|
| `FLAI_STORY` | `S-0576` | empty |
| `FLAI_KLUSTER_NAME` | the cluster's name, such as `authstar-integration` | the same |
| `FLAI_KLUSTER_NAMESPACE` | `s-0576` | `authstar` |
| `FLAI_KLUSTER_PREFIX` | `s0576_` | empty |
| `FLAI_KLUSTER_COMMIT` | the story branch's head, the whole hash | main's head |
| `FLAI_KLUSTER_TAG` | `s-0576-` and the head's first 12 hex digits | `main-` and the head's first 12 hex digits |
| `FLAI_KLUSTER_IMAGES` | `build` only: the file it lists its images in | the same |
| `FLAI_KLUSTER_STATE` | a folder kept for the namespace, `.flai-cache/kluster/state/s-0576` | `.flai-cache/kluster/state/authstar` |
| `KUBECONFIG` | flai's kubeconfig for the cluster (§1) | the same |
| `FLAI_KLUSTER_CONTEXT` | `k3d-` and the cluster's name | the same |

- The tag carries the commit so that a deploy at a new commit never runs an image of the old one under `imagePullPolicy: IfNotPresent`.
- `FLAI_KLUSTER_STATE` is where the deploy keeps what its checks need, such as the operator tokens S-0575 generated per namespace. It is outside the repository and git-ignored; flai deletes it after the teardown. It holds only values that work against the local cluster, which `safety.md` allows.
- flai gives no shared service's address. Those are the profile's and the project's deploy knows them, as authstar's scripts do.

### 5. When flai deploys, tears down, and promotes

While the `kluster` host action is on for a project with a `kluster:` block, `flai serve` runs these jobs, with no model:

| When | Job | Where it runs |
|------|-----|---------------|
| A story enters review, or is in review at a new head | Deploy the story: `build`, load, `deploy`, into its namespace | The story's worktree, at its head |
| A story leaves review other than by acceptance: back to in-progress, or cancelled | Tear the story down: `teardown`, then delete its namespace | The story's worktree, or the main checkout when it is gone |
| A story is accepted | Tear the story down, then promote: `build`, load, `deploy`, into the default namespace | Teardown as above; the promotion in a worktree of main's head that flai keeps, `.flai-cache/kluster/main` |

- Jobs run one at a time per project, in the order their causes came. Builds share the host's Docker cache and loads share the nodes, so running them at once would gain little and race. Promotions that queue up behind one another run once, at main's newest head.
- A story never deployed is not torn down. A teardown of a namespace that is gone succeeds.
- Each job's outcome, step times, and log are kept per namespace (§1), for `flai kluster status` and the dashboard's Deployments page (S-0581).
- A job that fails is recorded and stops there. It opens no thread and moves no work item: the operator sees it on the story's review page and in the journal. A failed deploy is not a finding against the story, since S-0575's baseline failed for reasons that were no story's.
- The close-out and `flai verify` do not use the cluster. They run before review, when the story has no namespace. A project's integration and smoke tiers stay its own.

### 6. How an image reaches the cluster

kluster has no image import, and the host has no `k3d` command: kluster drives k3d as a library. So flai loads each image `build` lists itself, as S-0575 did and authstar's own `deploy-local.sh` does: `docker save <image>` piped into `docker exec -i <node> ctr -n k8s.io images import -`, for each node container whose label `k3d.cluster` names the cluster and whose `k3d.role` is `server` or `agent`. It takes 31 to 34 seconds for authstar's eight images.

Once kluster gains an import (Consequences), flai calls it instead, and keeps this path for an older kluster.

### 7. The `kluster` host action

- `kluster` is a host action like `checks` or `merge-queue`: off for every project until the operator turns it on, with `flai serve enable kluster` or its toggle on the dashboard's Features page. It is refused for a project with no `kluster:` block.
- On, it lets `flai serve` run the jobs of §5 and lets the dashboard bring the cluster up or down and deploy, tear down, or promote through flai (S-0580, S-0581).
- Off, `flai serve` starts no job and leaves every namespace as it is. Turning it off tears nothing down.
- `flai kluster up`, `down`, `status`, `deploy`, `teardown`, and `promote` work from the operator's shell whatever the action says, as `flai checks run` does. `flai guard` refuses every agent and sub-agent those that write, since they change shared infrastructure on the operator's machine (`safety.md`). `flai kluster status` is a read every agent may run.

### 8. What the review checks see

- While both `kluster` and `checks` are on, `flai serve` starts the story's checks run once its deploy has passed at the story's head, not when the story enters review. With `checks` off, it deploys and runs no check.
- That run, and any checks run started for a story whose namespace is deployed at its head, gets the variables of §4 for the story's namespace in its environment, beside `{story}` and `{root}` in its arguments. So a project's check, such as authstar's golden path, reaches the story's own namespace through `KUBECONFIG` and `FLAI_KLUSTER_NAMESPACE`.
- A checks run for a story whose namespace is not deployed at its head gets none of them, and runs as it does today.
- Checks runs against the cluster run one at a time per project, after the job queue of §5, because authstar's scenario holds fixed local ports. A project whose checks can share the cluster may be allowed to run them at once later.

## Consequences

- S-0577 adds the `kluster:` block, its checks, `flai kluster up`, `down`, and `status`, and the host action. S-0578 adds `deploy`, `teardown`, and `promote` with the image load. S-0579 adds the jobs of §5 and the checks' environment of §8 to `flai serve`. S-0580 and S-0581 add the dashboard's Ops group. `design/system/kluster.md` describes the whole.
- A project fills in a build, a deploy, and a teardown script of its own. authstar's would be S-0575's `build.sh`, `ns-up.sh`, and `ns-down.sh` made into its repository's scripts.
- Signet keeps a torn-down namespace's bundles and policies until kluster or signet can delete them. A story's serial is never reused, so they only accumulate.
- A browser demo of two authstar namespaces at once needs a tenant host per namespace in portcullis's config and Dex's client. That is authstar's to change. Until then a demo reaches a namespace through a port-forward.
- Each story in review holds a namespace and its pods. On a 24-core, 62 GB host, two authstar namespaces and the default one ran together with room to spare; the review WIP limit bounds how many there are.
- system-flow itself runs no cluster. Its tooling convention says so, and that flai drives a project's kluster behind the host action.

### Proposals for kluster's repository

These are changes this decision would ask of `bytepunx/kluster`, proposed for that repository and made there. They are not work items here.

1. **An image import**, such as `kluster load <image>... --name <cluster>`, that puts a local image into every node. flai then calls it in place of `docker save` into `ctr`.
2. **A JSON status**, `kluster status --json`: each cluster's name, provider, state, and context, and each addon's version and whether it answers. flai's `kluster status` and the Klusters page read it instead of a table of `NAME`, `RUNNING`, and `AGE`.
3. **A namespace setup and teardown for the shared services**, such as `kluster namespace add <ns> --prefix <prefix>` and `kluster namespace remove <ns>`, for the profile's services: a RabbitMQ vhost and user, Postgres and ClickHouse databases with the prefix, and signet bundles and policies under the namespace, deleted on removal. kluster-lib has `SeedSignetBundle` and `EnsureSignetPolicy` but no command offers them. A project's deploy and teardown scripts would then call these rather than `rabbitmqctl`, `psql`, `clickhouse-client`, and their own signet seeding.
4. **A pinned signet chart**, as every other addon is, so that a bring-up installs a signet the profile's services are known to work with.
5. **A Dex that serves authstar's tenant**: the redirect URI `https://acme.authstar.app/oauth2/dex/callback` on the static client, and a password hash that matches the password its comment names.
6. **Machine-readable progress from `kluster up`**, such as `--json` events per step, so that a bring-up of 11 minutes shows where it is on the dashboard.
7. **A deploy of the profile's own services into a named namespace from a named tag**, in place of the constant `authstar` namespace and Docker Hub's `latest`, for a project that wants kluster to deploy it rather than scripts of its own.

Outside kluster, S-0575 found faults for signet (a bundle with config and no secrets fails `SyncBundle`; the CLI cannot delete a bundle's secrets or config) and drift in authstar's integration harness (the signet port, tower's scope, steward's and web's signet config, the scenario's fixed database names, its leaked port-forwards, and the fixed tenant host). Those go to those repositories as issues, and a project's checks against a namespace mean little until its baseline passes.

## Alternatives considered

- **A cluster per story.** Isolates fully, but costs 11m31s and every shared service again for each story, against 90 seconds for a namespace. Rejected.
- **kluster deploys the story.** kluster's `authstar` profile deploys into a constant namespace from `latest`, and what makes a namespace is per-service knowledge that is the project's. Proposal 7 would let a project choose it later; the deploy command stays the project's either way.
- **flai makes the per-namespace settings itself.** It would have to know each profile's services, their admin tools, and each project's databases and secrets. Rejected: the deploy command makes them from the namespace and the prefix, and proposal 3 would give it kluster's help.
- **Wait for kluster's image import before deploying anything.** Holds E-0032 on another repository for a step S-0575 did in 31 seconds with `docker` alone. Rejected; flai uses the import once it exists.
- **`k3d image import`.** The host has no `k3d` command, since kluster drives k3d as a library, and installing one for flai is the operator's. `docker` is already flai's ([ADR-0010](0010-shell-out-to-git-and-docker.md)).
- **Name the namespace from the story's title or branch.** A title changes and is long; `story-s-0576` adds nothing. The lower-case ID is short, stable, never reused, and a DNS label.
- **A fixed image tag per namespace, as S-0575 used.** Under `IfNotPresent` a deploy at a new commit would keep the old image. Rejected for the namespace and the commit.
- **Placeholders such as `{namespace}` in the commands' arguments.** The checks host action's `{story}` and `{root}` stay, but a deploy script needs ten values and its checks the same ones; environment variables reach both alike, and reach the scripts they call.
- **Give the shared services' addresses as variables.** They are the profile's, not flai's, and change with it. The project's scripts know them.
- **Gate acceptance on a passing deploy.** S-0575's baseline failed for drift no story caused. The operator sees the outcome and decides; the orchestrator's acceptance may take it as evidence in a later decision.
- **Deploy on the move to review in the story's agent's turn.** A deploy takes minutes and changes shared infrastructure on the operator's machine. A job in `flai serve` with no model, behind a host action, as the close-out and the merge queue are, keeps it the operator's.
