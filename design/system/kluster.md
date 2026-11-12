---
title: Running a project's kluster
updated: 2026-11-12
status: active
topics: [kluster, deployment, manifest]
---

# Running a project's kluster

How flai runs a project's local Kubernetes cluster, made by `bytepunx/kluster`, so that each story in review is deployed into a namespace of its own over the project's shared infrastructure, checked there, and promoted to the default namespace on acceptance. [ADR-0254](../adrs/0254-flai-runs-one-kluster-per-project-from-the-config-its-manifest-s-kluster-block.md) decides it and gives the reasons. [S-0575's results](../experiments/S-0575-two-builds-of-authstar-run-side-by-side-in-one-kluster-each-in-a-namespace.md) show it working by hand on authstar, with the times of each step. [ADR-0259](../adrs/0259-enabling-kluster-on-a-project-with-no-kluster-block-writes-one-naming-only-its.md) refines it (S-0586): turning kluster on gives a project with no block one that names only its config, `deploy` is optional, and the dashboard's Klusters page holds the whole setup. [ADR-0263](../adrs/0263-a-kluster-block-may-list-components-each-built-as-its-own-step-with-its-own-log.md) refines both (S-0595): a block may list components, each built as its own step with its own log, a job's record lists every step before the first runs, `kluster.log` reads a step's log from an offset, and the Deployments page holds the deploy setup, the builds, and the live logs in an in-page rail.

[ADR-0292](../adrs/0292-a-story-s-agent-asks-flai-serve-to-deploy-check-and-tear-down-its-own-story-s.md) refines ADR-0254 for a story in progress (S-0659): its agent asks flai serve to deploy, check, and tear down its own story's namespace before review, the namespace lives from its first deploy until acceptance or cancellation, `kluster.agent_namespaces` bounds how many stories in progress hold one, and the namespace's record at the story's head is the evidence the move to review, the close-out, and an acceptance read. None of it is built yet. [Deploys a story's agent asks for](#deploys-a-storys-agent-asks-for-adr-0292) describes it, and each part not built names the story that builds it: S-0660, S-0661, S-0662, or S-0663.

system-flow itself runs no cluster. This is what flai does for a project that names one in its manifest, behind the `kluster` host action.

## The project's cluster

- A project has at most one cluster. It is made from the kluster config file `kluster.config` names, as that file stands on the main branch. A story that changes the file changes nothing in the running cluster.
- The cluster's name is the config's `name`. Its provider is `k3d`; `flai check` reports any other.
- The shared infrastructure is what the config's profile installs. kluster v1.3.1 has no profile that installs nothing: a config with no `profile` gets `spire`, which installs cert-manager, SPIRE, and traefik-tls ([ADR-0259](../adrs/0259-enabling-kluster-on-a-project-with-no-kluster-block-writes-one-naming-only-its.md)). For authstar's `integration/kluster.yaml`, profile `authstar`, it is SPIRE, signet, RabbitMQ, Postgres, ClickHouse, and Dex. flai installs nothing in the cluster and changes no shared service's own settings.
- flai brings the cluster up with `kluster --config <config> up` and takes it down with `kluster down`, in the main checkout. A bring-up took 11m31s in S-0575.
- flai writes the cluster's kubeconfig to `.flai-cache/kluster/kubeconfig` with `kluster kubeconfig --name <name> --output`, and passes it to every command as `KUBECONFIG`. It never merges it into `~/.kube/config`. The context is `k3d-<name>`.

What flai keeps, all under `.flai-cache/kluster/`:

| Path | What it holds |
|------|---------------|
| `kubeconfig` | The cluster's kubeconfig |
| `deployments/<namespace>.json` | The namespace's record: its commit, the outcome and times of each step, and their logs |
| `logs/<namespace>/<step>.log` | Each step's output in the namespace's last job, and `images`, or `images-<name>` per component, the lists the builds wrote |
| `job.lock`, `job.json` | The project's job lock, and the job that holds it |
| `state/<namespace>/` | The namespace's state folder, `FLAI_KLUSTER_STATE` |
| `main/` | A worktree of main's head, where a promotion builds and deploys |

## The `kluster:` block

The block in `system-flow.yaml` is optional. Without it, flai runs no cluster for the project. Turning the `kluster` host action on gives a project with no block one, as [the `kluster` host action](#the-kluster-host-action) says ([ADR-0259](../adrs/0259-enabling-kluster-on-a-project-with-no-kluster-block-writes-one-naming-only-its.md)):

```yaml
kluster:
  config: kluster.yaml
```

with `kluster.yaml` at the repository root, written when absent and kept when there:

```yaml
name: <the manifest's name as a DNS label>-local
provider: k3d
profile: spire
```

A full block, such as authstar's:

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
| `components` | no | none: `build` builds, or nothing is loaded | In place of `build`: a list of components, each a `name`, a DNS label unique in the list, and its own `build`, an argument list, run one at a time as the steps `build-<name>` |
| `deploy` | no | none: the cluster comes up and goes down, and flai deploys nothing | Makes the namespace's settings in the shared services, applies the project's manifests into the namespace, and waits for their rollouts |
| `teardown` | no | none: flai deletes only the namespace | Removes what `deploy` made outside the namespace |
| `timeout` | no | `30m` | The most one job may take, a Go duration |
| `agent_namespaces` | no | `2` | The most stories in progress that hold a namespace at once, from deploys their agents asked for; `0` turns deploys before review off. Stories in review do not count. Decided by [ADR-0292](../adrs/0292-a-story-s-agent-asks-flai-serve-to-deploy-check-and-tear-down-its-own-story-s.md) §4, not built yet (S-0660) |

A project that builds several images names one build per component instead of `build` ([ADR-0263](../adrs/0263-a-kluster-block-may-list-components-each-built-as-its-own-step-with-its-own-log.md)):

```yaml
kluster:
  config: kluster.yaml
  components:
    - name: api
      build: [scripts/build.sh, api]
    - name: web
      build: [scripts/build.sh, web]
  deploy: [scripts/kluster-deploy.sh]
```

A block names `build` or `components`, not both.

Each command is an argument list, run as it stands, never through a shell, as the manifest's `tests` are. It runs in the checkout being deployed, with the variables below in its environment. A command that exits other than 0 fails its job. `deploy` is idempotent: a second deploy into the same namespace, at a new commit, updates it in place.

- A block with no `deploy` is a project whose deploys are not set up yet ([ADR-0259](../adrs/0259-enabling-kluster-on-a-project-with-no-kluster-block-writes-one-naming-only-its.md), replacing ADR-0254's `deploy` required). `flai serve` deploys and promotes nothing for it, and `flai kluster deploy` and `promote` refuse it.
- Each key is a setting `flai manifest set` writes, `kluster.config` to `kluster.timeout`, the commands as JSON argument lists such as `["scripts/kluster-deploy.sh","--wait"]`, `[]` to unset one, and `kluster.components` as a JSON list such as `[{"name":"api","build":["scripts/build.sh","api"]}]`, `[]` to unset it, written into the block as a YAML list. The setup form both Ops pages show writes them through `settings.manifest`, under the `settings` host action, since the commands run on the host as the operator.
- The config file `kluster.config` names is read and saved from the Klusters page through `file.get` and `file.save` ([ADR-0147](../adrs/0147-the-dashboard-reads-and-saves-a-text-file-that-is-not-markdown-through-file-get.md)), whatever `settings` says: it names a profile and addons, not a command.

## The namespaces

| Namespace | Who deploys into it | Name | Prefix |
|-----------|---------------------|------|--------|
| A story's | The story, while in review; and, decided by ADR-0292 and not built yet (S-0660), while in progress when its agent asks | Its ID in lower case, with the zero padding: `s-0576` | `s0576_`: the namespace without hyphens, then an underscore |
| The default | Main, on each acceptance | `kluster.namespace`, such as `authstar` | Empty, so main keeps the names the profile made |

- `flai check` reports a default namespace that is not a DNS label, or that has a story namespace's form, `s-` and digits.
- Today a story's namespace lives while the story is in review, and while a send-back keeps it. Decided by ADR-0292 and not built yet (S-0660), it lives from its first deploy, in progress or in review, until the story is accepted, cancelled, or moved back to ready or backlog, and `kluster.agent_namespaces` bounds how many stories in progress hold one.
- The prefix is for the names a namespace makes in shared services, such as databases: `s0576_tower`.
- What a story owns in the shared services is the project's to make and remove: its RabbitMQ vhost and user, its databases, its signet bundles and policies, its keys and tokens, and its tenant where the project has tenants, as authstar has. `deploy` makes them from the namespace and the prefix, and `teardown` removes them. flai does not know them.
- A project's Kubernetes manifests and scripts take the namespace as a parameter, with the default namespace as its default, and never name it as a constant. authstar's manifests still name `authstar`, which S-0575 renamed with `sed`; that is authstar's to change.

## The steps and their variables

A deploy is the first three steps, and a teardown the fourth:

1. **Build.** `build` builds the project's images in the checkout and writes each reference, one a line, to `FLAI_KLUSTER_IMAGES`. With `components`, each component's build is a step of its own, `build-<name>`, in the list's order, with its own log and its own images file, and `FLAI_KLUSTER_COMPONENT` set to its name. A failed build stops the deploy, and the steps after it are recorded `skipped`.
2. **Load.** flai loads each listed image into every node of the cluster: `docker save -o` into a temporary file, which is then fed to `docker exec -i <node> ctr -n k8s.io images import -`, for each container labelled `k3d.cluster=<name>` with `k3d.role` `server` or `agent`. Once kluster has an image import, flai calls that instead.
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
| `FLAI_KLUSTER_IMAGES` | A build only: the file it lists its images in, `logs/<namespace>/images` for `build` and `logs/<namespace>/images-<name>` for a component | The same |
| `FLAI_KLUSTER_COMPONENT` | A component's build only: the component's name | The same |
| `FLAI_KLUSTER_STATE` | `.flai-cache/kluster/state/s-0576` | `.flai-cache/kluster/state/authstar` |
| `KUBECONFIG` | flai's kubeconfig for the cluster | The same |
| `FLAI_KLUSTER_CONTEXT` | `k3d-` and the cluster's name | The same |

- `FLAI_KLUSTER_STATE` is where `deploy` keeps what the checks need, such as a namespace's operator tokens. It is git-ignored, holds only values that work against the local cluster, and is deleted after the teardown.
- flai gives no shared service's address. The project's scripts know their profile's.
- A story's agent will get the same variables for its own namespace from `flai kluster env` and `flai kluster exec`, but for `FLAI_KLUSTER_IMAGES`, and with `KUBECONFIG` a kubeconfig held to the namespace in place of flai's ([ADR-0292](../adrs/0292-a-story-s-agent-asks-flai-serve-to-deploy-check-and-tear-down-its-own-story-s.md) §2, not built yet: S-0660).

## What happens when

While the `kluster` host action is on for a project with a `kluster:` block that names a `deploy`, `flai serve` runs these jobs, with no model. With no `deploy` it runs only the teardowns. The table is [ADR-0292](../adrs/0292-a-story-s-agent-asks-flai-serve-to-deploy-check-and-tear-down-its-own-story-s.md) §3's, which replaces ADR-0254 §5's; its last column says what is built and, for each row that is not, what flai does until the story named builds it:

| When | Job | Where it runs | Built |
|------|-----|---------------|-------|
| The agent asks for a deploy, its story in progress with no namespace, under `kluster.agent_namespaces` | Deploy: the namespace is made | The story's worktree, at its head | S-0660; until then `flai guard` refuses agents every deploy |
| The agent asks for a deploy, its story's namespace there | Deploy again, in place | The story's worktree, at its head | S-0660 |
| The agent asks for a checks run | Checks run against the namespace as deployed | The story's worktree | S-0661 |
| The agent asks for a teardown | Teardown | The story's worktree | S-0660 |
| A story in progress gains commits, and its agent asks nothing | None: the namespace stays at the commit it was deployed at | Nowhere | Yes |
| A story enters review, its namespace's deploy passed at its head, and its checks too while `checks` is on | None: those stand as the review's | Nowhere | S-0660; until then a story enters review undeployed, sent back ones aside |
| A story enters review otherwise: no namespace, another commit, or a failure | Deploy, then checks, or the checks run alone when only it has not passed at the head | The story's worktree, at its head | Yes, the deploy and its checks; the checks run alone, S-0661 |
| A story in review gets a new head | Deploy again, in place | The story's worktree, at its new head | Yes |
| A story leaves review for in-progress, sent back by a failure or by anyone | None: the namespace stays | Nowhere | Sent back by a failure, yes; by anyone else, S-0660, and until then a teardown |
| A story moves back to ready or backlog | Teardown | The story's worktree, or the main checkout when it is gone | Yes |
| A story is cancelled, from any column | Teardown | The story's worktree, or the main checkout when it is gone | Yes |
| A story is accepted | Teardown of its namespace, then a promotion | The teardown as above; the promotion in `.flai-cache/kluster/main` | Yes |
| A job finds the cluster absent | Bring the cluster up first | The main checkout | Yes, for a review deploy; for a request, S-0660 |

- **At the head**, decided by ADR-0292 §3 and not built yet (S-0662): a record is at the head when its commit is the head, or when every commit from the record's commit to the head changes only paths under the manifest's `design`, `docs`, and `wip` folders. It will hold for the review row above, for the acceptance gate, and for what the move to review and the close-out report. Today the record's commit must be the head.
- **A send-back needs no case of its own** once a namespace stays while its story is in progress (S-0660). Until then the pass keeps a sent-back namespace only while its story stands where the send-back put it, as [flai serve (S-0579)](#flai-serve-s-0579) says.

- A promotion builds, loads, and deploys main's head into the default namespace. There `deploy` replaces what the profile deployed. Promotions queued behind one another run once, at main's newest head.
- Jobs run one at a time per project, in the order their causes came.
- A story never deployed is not torn down.
- `flai serve` never takes the cluster down.
- The close-out and `flai verify` do not use the cluster: they run before review, when the story has no namespace. Once its agent can deploy it before review, the close-out gains a `kluster` step that reads the namespace's record and runs nothing ([deploys a story's agent asks for](#deploys-a-storys-agent-asks-for-adr-0292), S-0662).

## Deploys a story's agent asks for (ADR-0292)

[ADR-0292](../adrs/0292-a-story-s-agent-asks-flai-serve-to-deploy-check-and-tear-down-its-own-story-s.md) decides it and gives the reasons. None of this section is built yet; each part names the story that builds it. Until then `flai guard` refuses every agent the `flai kluster` commands that write, and a story's namespace is made only in review.

### The requests (S-0660, S-0661)

A story's agent asks flai serve for a job on its own story's namespace, as it asks for the close-out ([ADR-0229](../adrs/0229-the-close-out-runs-as-a-flai-serve-job-the-story-s-agent-asks-for-with-the.md)):

| Request | MCP tool | On the host | Built by |
|---------|----------|-------------|----------|
| Deploy at the branch's head: build, load, deploy, then the checks while `checks` is on | `kluster_deploy` | `flai kluster deploy S-nnnn --job` | S-0660 |
| Run the checks again against the namespace as deployed | `kluster_checks` | `flai checks run S-nnnn --job` | S-0661 |
| Tear the namespace down | `kluster_teardown` | `flai kluster teardown S-nnnn --job` | S-0660 |

- Each writes `S-nnnn.<kind>.request.json` under `.flai-cache/kluster/requests/` and waits at most ten seconds. flai serve looks for requests every second, beside the close-out jobs, queues each in the project's one job queue, and answers the job's id and the number of jobs ahead. A request it does not take in ten seconds is withdrawn, and the agent is told flai serve is not running; it cannot run the job itself.
- A request is for the session's own story, `FLAI_STORY` or the story in `agent-S-nnnn`, and refused for any other. The story is in progress: one in review is answered with the review's job, or with its record when that passed at the head, and one in any other column is refused.
- A request is refused while the worktree has uncommitted changes to tracked files, naming them. The job records the head it built.
- It is refused while the `kluster` host action is off or the manifest has no `kluster:` block, naming `flai serve enable kluster` as the operator's. With no `deploy`, `kluster_deploy` and `kluster_checks` are refused, saying deploys are not set up yet. `kluster_checks` is refused while the `checks` host action is off, and before the namespace has a record that passed.
- A request while a job of the story waits or runs is answered with that job. A deploy request at a newer head takes the place of the one waiting. A request always runs: the pass's rule that a head is tried once does not hold it back.
- A deploy builds every component, as a review deploy does.

### The limit (S-0660)

- `kluster.agent_namespaces`, default 2, is the most stories in progress that hold a namespace at once. 0 turns deploys before review off. `flai manifest set kluster.agent_namespaces` writes it, the setup form shows it under `settings`, and `flai check` reports a value that is not a whole number of 0 or more as `manifest.kluster`.
- Stories in review do not count, and their deploys are never refused for the limit.
- The limit makes no namespace and tears none down: a story sent back from review keeps its namespace whatever the count. A request for a story whose namespace is there is never over it.
- A deploy request for a story with no namespace, over the limit, is refused at once, naming the limit and each story in progress that holds a namespace, with its last deploy's commit and time. It is not queued.

### What the agent runs against its namespace (S-0660, S-0661)

- The project's checks, through flai serve: the checks run after a passing deploy, or the one `kluster_checks` asks for, is the review's, with the namespace's variables and flai's whole kubeconfig, in the queue (S-0661).
- Its own commands: `flai kluster env S-nnnn` prints the namespace's variables, as `export` lines or `--json`, and `flai kluster exec S-nnnn -- <command>` runs a command with them in the worktree, for the agent's own story while its namespace has a record, a failed one too (S-0660).
- In what `env` and `exec` give, `KUBECONFIG` is `.flai-cache/kluster/agents/s-nnnn.kubeconfig`: each deploy makes it after the deploy step passes, from a ServiceAccount `flai-agent` in the namespace bound by a RoleBinding to the `edit` ClusterRole, with a token from `kubectl create token` valid for 24 hours. The API server refuses it anything outside the namespace. The teardown deletes it (S-0660).
- `flai kluster log S-nnnn --step <step>` prints a step's log as `kluster.log` reads it; with `flai kluster deployments` and `flai kluster status` it is open to every agent for its own story (S-0660).
- `flai guard` refuses every agent and sub-agent `flai kluster up`, `down`, and `promote`; `deploy` and `teardown` without `--job`; requests, `env`, `exec`, and `log` for another story's namespace or the default one; the `kluster` program but `status`; and a shell line naming `.flai-cache/kluster/kubeconfig` (S-0660, S-0661).

### Waiting and failures (S-0660)

- When a job an agent asked for ends, flai serve appends a `kluster` event to `.flai-cache/events.jsonl`, with the story as its ID and fields `job`, `kind` (`deploy`, `checks`, or `teardown`), `commit`, `outcome` (`passed`, `failed`, `refused`, or `interrupted`), `failed_step`, `exit_code`, and `log`. The inbox reports it to the story's agent and to readers with no story of their own. The agent holds `wait_for_events` for it, which answers it no `end` while a job of its story waits or runs.
- A job waits for those ahead of it, each bounded by `kluster.timeout`, then runs within it. A deploy and its checks run are one job. A job still waiting or running when flai serve stops ends `interrupted` when it starts again.
- A failure before review is the agent's: flai serve opens no thread, moves nothing, and starts no agent. The namespace stays as the failure left it. The failure reaches the operator only when the agent says on a thread on its story that the cause is outside the story, or at review, where [a failed deploy holds the story](#a-failed-deploy-holds-the-story).

### What is recorded and read (S-0660, S-0662)

- The namespace's record holds a deploy the agent asked for as it holds a review's, and gains `asked_by`, the agent or `flai serve`, and `checks_log`. A checks run the agent asks for writes `checks_commit`, `checks_outcome`, and `checks_log`. What the agent runs with `exec` is not recorded (S-0660, S-0661).
- The move to review reads the record and says what was deployed and checked at the head, or why not, and refuses nothing (S-0662).
- The close-out gains a step, `kluster`, that reads the record and runs nothing. It passes when the action is off, the block has no `deploy`, or `kluster.agent_namespaces` is 0, reports the record against the head otherwise, and fails only on a deploy or checks run that failed at the head (S-0662).
- An acceptance reads the record through the `kluster` blocker, with at the head as above. A deploy and checks run the agent asked for that passed at the head lift it without a review deploy (S-0662). The review page shows the record (S-0662), and the Deployments page shows the namespaces of stories in progress (S-0663).

## A failed deploy holds the story

This is ADR-0254 §5a's, built by S-0694. S-0579 built the thread alone: a failure moved nothing, and no acceptance waited for the story's deploy ([I-0268](../issues/I-0268-adr-0254-5a-s-send-back-of-a-story-whose-kluster-deploy-or-checks-failed-and.md)). [flai serve (S-0579)](#flai-serve-s-0579) says how it is built.

- A story's failed deploy, or a checks run against its namespace that does not pass, is recorded. `flai serve` opens one thread on the story, naming the job, the step, the commit, and the log, and sends the story back to in-progress, starting its agent. The agent fixes the change, or says on the thread that the cause is outside the story.
- A story that is no longer in review at the failed head, or whose action went off meanwhile, is not sent back. The thread says why.
- A story sent back keeps its namespace, as the failure left it, while it stands in progress where the send-back put it. Back in review, at a new head or the same one, it is deployed again in place. Cancelled, accepted, or moved on by anyone else, it is torn down. Once S-0660 keeps every namespace while its story is in progress, only cancellation, acceptance, or a move back to ready or backlog tears it down, and a deploy and checks run the agent asks for that pass at the head resolve the send-back's threads ([ADR-0292](../adrs/0292-a-story-s-agent-asks-flai-serve-to-deploy-check-and-tear-down-its-own-story-s.md) §3).
- This holds for failures in review. A deploy or checks run the agent asked for that fails before review sends nothing back and opens no thread: it is the agent's to fix, as [waiting and failures](#waiting-and-failures-s-0660) says (S-0660).
- While the action is on, every acceptance of a story is refused with the blocker `kluster`, unless the story's namespace was deployed at its head and, while `checks` is on, its checks run against it passed there. This holds for the operator's acceptance, the dashboard's, the orchestrator's, and the merge queue's landing. The blocker names the namespace's last job, its commit, and its outcome.
- The merge queue lists a story the gate refuses as ineligible, with the reason `kluster`, and does not batch it. The orchestrator's enqueue is not refused for it: the queue waits for the gate, so a story may be enqueued before its deploy has run.
- `flai accept --without-kluster`, and `flai move S-nnnn done --without-kluster`, land the story anyway, by hand and outside the merge queue. It is the operator's alone, refused with `--by orchestrator` and `--batch`, and recorded under `### Accepted without the kluster gate` in the story's Notes.
- A failed promotion has no story to send back. `flai serve` records it and opens one thread to the operator on `system-flow.yaml`, naming the commit, the step, and the log. The next acceptance promotes again.
- Until a project's baseline passes in the default namespace, each of its stories is sent back at its first deploy or checks run.

## The `kluster` host action

| State | What it allows |
|-------|----------------|
| Off, the default | `flai serve` starts no job and leaves every namespace as it is. Turning the action off tears nothing down. |
| On | `flai serve` runs the jobs above, and the dashboard may bring the cluster up or down and deploy, tear down, or promote through flai |

- The operator turns it on from the Klusters page's setup, its toggle on the dashboard's Features page, or `flai serve enable kluster`. All three run `flai serve enable kluster`. For a project with no `kluster:` block it writes the block and `kluster.yaml` shown in [the `kluster:` block](#the-kluster-block), and commits those two files alone, with the dashboard's trailer, unless the project sets `dashboard.autocommit: false`. It says what it wrote and the commit. `--all-projects` writes no block ([ADR-0259](../adrs/0259-enabling-kluster-on-a-project-with-no-kluster-block-writes-one-naming-only-its.md)). The dashboard's toggle needs the `settings` host action, which is turned on only in a shell ([ADR-0039](../adrs/0039-a-settings-host-action-turned-on-only-in-a-shell-lets-the-dashboard-change-the.md)).
- While it is on, the Deployments page starts and stops each service's port-forward, which `flai serve` runs bound to `0.0.0.0`, as [the Deployments page](#the-deployments-page) says.
- `flai kluster up`, `down`, `status`, `deploy`, `teardown`, and `promote` work from the operator's shell whatever the action says.
- `flai guard` refuses every agent and sub-agent the `flai kluster` commands that write. `flai kluster status` is a read any agent may run. ADR-0292 changes this, not built yet (S-0660, S-0661): a story's agent asks for deploys, checks runs, and teardowns of its own namespace through the requests above, the host action still gating each, and runs `env`, `exec`, and `log` for it, while the guard keeps refusing everything else that writes, as [what the agent runs against its namespace](#what-the-agent-runs-against-its-namespace-s-0660-s-0661) lists.

## What the review checks see

- With `kluster` and `checks` both on, `flai serve` starts a story's checks run once its deploy has passed at the story's head, not when the story enters review.
- A checks run for a story whose namespace is deployed at its head gets the variables above for that namespace, beside `{story}` and `{root}` in its arguments. A check such as authstar's golden path reaches the namespace through `KUBECONFIG` and `FLAI_KLUSTER_NAMESPACE`.
- A checks run for a story with no namespace deployed at its head gets none of them.
- Checks runs against the cluster run one at a time per project, after the job queue.
- With `checks` off, a story is deployed and nothing is checked against it.
- Before review, decided by ADR-0292 and not built yet (S-0661): the checks run after a deploy the agent asked for, or one it asks for with `kluster_checks`, is this same run, with the same variables, recorded in the namespace's record. A story that enters review with its deploy and checks passed at its head is not checked again there.

## The Ops pages in the dashboard

E-0032 gives the dashboard's rail an Ops group with two pages, each behind the `kluster` host action and acting only through flai:

| Page | What it shows | What it does |
|------|---------------|--------------|
| Klusters | The setup's steps and what each lacks; then the project's cluster: its name, profile, provider, state, and age, the up or down flai serve runs or last ran, `flai check`'s findings, and the nodes and pods | Turns the action on, edits the block and the config, makes a story for the scripts, brings the cluster up or takes it down |
| Deployments | Main's namespace and each story's, with its commit, outcome, each build and step with its live log, and how to reach it for a demo | Deploys, tears down, or promotes, starts and stops port-forwards, and edits the block's deploy setup |

## Commands (S-0577, S-0578)

S-0577 built `flai kluster up`, `down`, and `status`, the check of the `kluster:` block, the `kluster` host action, and the guard's rule. S-0578 built `deploy`, `teardown`, `promote`, and `deployments`, with the image load, the job lock, the records, and the logs.

The commands run `kluster` in the project's main checkout, with the block's config, whatever the host action says (`flai/internal/kluster`, `flai/cmd/kluster.go`). Each passes `--config <config> --provider <provider>`, the provider as the config names it.

| Command | What it runs | What it prints | Exit |
|---------|--------------|----------------|------|
| `flai kluster up` | `kluster up`, its output streamed; then `kluster kubeconfig --name <name> --output .flai-cache/kluster/kubeconfig` | kluster's output, then the cluster's name, the kubeconfig's path, and the context `k3d-<name>` | 1 when either fails, naming the step; a failed `up` writes no kubeconfig |
| `flai kluster down [--yes]` | `kluster down --name <name> --yes`, so kluster never prompts | kluster's output, then `<name> is down` | 1 when kluster fails, when the operator answers no, or without a terminal and without `--yes` |
| `flai kluster status [--json]` | `kluster status`, read by its `NAME`, `RUNNING`, and `AGE` headers, the row of the config's `name` | The name and running, stopped, or not created, then the profile, provider, context, and age; `--json` prints one object with `name`, `profile`, `provider`, `context`, `exists`, `running`, `age`, and `age_seconds` | 0 whether the cluster runs, is stopped, or is absent; 1 when kluster fails or prints a table without `NAME` and `RUNNING` |
| `flai kluster deploy S-nnnn [--json]` | `git rev-parse HEAD` in the story's worktree, then build, load, and deploy into the story's namespace from that worktree | A line per step, `<namespace>: <step>, log <path>`, then the record: outcome, commit, tag, images, and each step with its outcome and log; `--json` prints the record and sends the step lines to standard error | 1 when a step fails, naming the step, its exit code, and its log, after printing the record; 1 on a refusal |
| `flai kluster teardown S-nnnn [--yes] [--json]` | The block's `teardown`, when it names one, then `kubectl delete namespace`, in the story's worktree, or the main checkout when it is gone | `<story>'s namespace <namespace> is torn down`; `--json` prints the record removed | 0, saying nothing was deployed, for a story with no record; 1 when a step fails, or when the operator answers no, or without a terminal and without `--yes` |
| `flai kluster promote [--json]` | Moves the worktree `.flai-cache/kluster/main` to the main branch's head, then build, load, and deploy into the default namespace from it | As `deploy`, for the default namespace | As `deploy` |
| `flai kluster deployments [--json]` | Nothing: it reads the records | A table, `NAMESPACE`, `STORY` (`(main)` for the default namespace), `COMMIT` (12 digits), `OUTCOME` (or the running job's kind and `running`), and `ENDED`, or `no deployments`; `--json` prints an array of records, `[]` for none | 0; 1 when a record cannot be read |

- `up` and `down` stop on Ctrl-C or SIGTERM. Neither has a timeout of its own; `kluster.timeout` is for `flai serve`'s jobs.
- `down` asks on a terminal unless `--yes`, and refuses a non-terminal without it, since it deletes the cluster and every namespace in it.
- A column a later kluster adds to its table is passed over. Once kluster has a JSON status, flai reads that instead ([ADR-0254 § Proposals for kluster's repository](../adrs/0254-flai-runs-one-kluster-per-project-from-the-config-its-manifest-s-kluster-block.md#proposals-for-klusters-repository)).
- A project with no block is refused, saying to turn the `kluster` host action on, which writes the block (S-0586), and a host with no `kluster` on `PATH` with kluster's `install.sh` and `go install github.com/bytepunx/kluster/kluster@latest`. The tests fake `kluster` on `PATH`; none needs a cluster or Docker.

What else S-0577 built:

- **The check.** `flai check` reports each fault of the block as the error `manifest.kluster` on its line, one sentence each beginning with the key: an unknown key, a value of the wrong kind, a `config` missing, absent, unreadable, naming no cluster, or naming a provider other than `k3d`, a default namespace that is not a DNS label or has a story namespace's form, a command with no program, and a `timeout` that is not a duration above zero. Since S-0595 it also reports `build` and `components` both set, and a component that is not a block, has an unknown key, has no `name`, a name that is not a DNS label or that another component has, or has no `build` or no program ([ADR-0263](../adrs/0263-a-kluster-block-may-list-components-each-built-as-its-own-step-with-its-own-log.md)). `manifest.Load` never refuses the block ([project-manifest.md](project-manifest.md)). Since S-0586 a missing `deploy` is no fault ([ADR-0259](../adrs/0259-enabling-kluster-on-a-project-with-no-kluster-block-writes-one-naming-only-its.md)).
- **The host action.** `kluster` is in `hostapi.Actions`, seen by the dashboard, so the Features page lists it and turns it on and off with no change to flaiover. `flai serve enable kluster`, which that toggle runs, was refused for a project with no block until S-0586, which has it write the block instead ([the `kluster` host action](#the-kluster-host-action)). Since S-0579, `flai serve` runs its jobs while it is on ([flai serve (S-0579)](#flai-serve-s-0579)).
- **The guard.** `flai guard` refuses every session `flai serve` started and every sub-agent `flai kluster up`, `down`, `deploy`, `teardown`, and `promote`, named in one list, with the reason and what to do instead; `status` stays open ([flai-cli.md](flai-cli.md#commands)).

What S-0578 built, in `flai/internal/kluster/deploy.go`, for the commands above and for `flai serve`'s jobs (S-0579):

- **The refusals.** A deploy refuses a story that is not in progress or in review, a story with no worktree (open it with `flai stream open`), a block with no `deploy`, a cluster with no kubeconfig at `.flai-cache/kluster/kubeconfig` (bring it up with `flai kluster up`), and a commit that is not a whole hash. Teardown and promote refuse the same where they apply. Every job refuses while another holds the lock.
- **The lock.** One job runs at a time per project: a deploy, a teardown, or a promotion. A job takes `.flai-cache/kluster/job.lock`, waiting at most 200 ms, and names itself in `job.json`: its kind, story, namespace, start, and process. A job that finds the lock held is refused with `kluster.ErrBusy`, naming the job that holds it. The lock is released when the process ends, so a crashed job leaves none held; a `job.json` without a held lock names no running job.
- **The record.** Each namespace has one record, `.flai-cache/kluster/deployments/<namespace>.json`, written whole each time (`atomicfile`). A job writes it as `running` when it starts, with every step it will run listed `pending` (S-0595, [ADR-0263](../adrs/0263-a-kluster-block-may-list-components-each-built-as-its-own-step-with-its-own-log.md)): a deploy or a promotion its builds, `load`, and `deploy`, a teardown its `teardown` after the steps the record already holds. Each step is written `running`, with its start and its log, when it starts, and with its outcome when it ends, in place of its pending entry. At the job's end the steps still pending, those after a failure, are written `skipped`. So the record says at any moment which steps passed, which runs, which failed, and which wait, and a job cut short leaves what it did. A deploy again at a new commit replaces it.

  | Field | Meaning |
  |-------|---------|
  | `story` | The story deployed; empty for the default namespace |
  | `namespace` | The namespace |
  | `commit` | The checkout's head, the whole hash |
  | `tag` | `FLAI_KLUSTER_TAG` |
  | `checkout` | The folder the commands ran in |
  | `images` | The images the build listed |
  | `started_at`, `ended_at` | When the job started and ended; `ended_at` is absent while it runs |
  | `outcome` | `running`, then `passed` or `failed` |
  | `failed_step`, `exit_code`, `error` | Which step failed, its exit code (-1 when its program gave none, as when it was killed at `kluster.timeout`), and the error; absent when it passed |
  | `checks_commit`, `checks_outcome` | The commit and the outcome, `passed`, `failed`, `timed-out`, or `cancelled`, of the last checks run against the namespace since its deploy; absent until a run ends, and cleared by a deploy (S-0694) |
  | `steps` | Each step: `name` (`build`, `build-<name>`, `load`, `deploy`, or `teardown`), `started_at`, absent while it is pending, `ended_at`, `outcome` (`pending`, `running`, `passed`, `failed`, or `skipped`), `exit_code`, and `log` |

- **The list.** `kluster.Deployments` reads every record: the default namespace's first, then the stories' by their number, then any other by name. Each carries `default`, and `job` while a running job acts on its namespace. A record that still says `running` while no job runs for its namespace is shown `interrupted`.
- **The logs.** Each step's output goes to `.flai-cache/kluster/logs/<namespace>/<step>.log`, named in the record. A deploy clears the namespace's logs first, so they are its last job's. The build lists its images in `logs/<namespace>/images`, and a component's build in `logs/<namespace>/images-<name>`, the file `FLAI_KLUSTER_IMAGES` names; the load takes the images of every build.
- **The steps.** Each command runs as an argument list, never through a shell, in the checkout, with the namespace's variables over flai's own environment; a program named with a relative path, such as `scripts/kluster-deploy.sh`, is found from the checkout. A block with no `build` and no `components` records build and load as `skipped`. The whole job is bounded by `kluster.timeout`.
- **The load as built.** flai lists the nodes with `docker ps --filter label=k3d.cluster=<name>`, keeping the containers whose `k3d.role` label is `server` or `agent`. For each image it runs `docker save -o` into a temporary file, then gives that file to `docker exec -i <node> ctr -n k8s.io images import -` on each node. A build that lists no image loads nothing.
- **Teardown.** It reads the story's record, runs the block's `teardown`, when it names one, with the namespace's variables at the record's commit, then `kubectl --context k3d-<name> delete namespace <namespace> --wait --ignore-not-found`, both to the namespace's `teardown.log`. When both pass it removes the record, the state folder, and the logs. When either fails it keeps the record, marked `failed` at `teardown`. A story with no record is `kluster.ErrNotDeployed`: nothing runs.
- **Promote.** Under the lock, it resolves the main branch's head with `git rev-parse --verify <branch>^{commit}` in the main checkout. It moves the worktree `.flai-cache/kluster/main` there with `git checkout --detach --force`, or, when there is none, makes it with `git worktree prune` and `git worktree add --detach`. Then it deploys from that worktree into the default namespace, as a story's deploy does, with an empty `FLAI_STORY` and prefix and the tag `main-<commit>`. A failure of the git steps writes no record.

## flai serve (S-0579)

`flai/internal/serve/kluster.go` is the kluster pass. `serve.Run` starts one per served project, beside the forge pass. It does nothing, and runs no kluster command, unless the project's manifest has a `kluster:` block and the `kluster` host action is on for it. Both are asked at every look and again before every job. A block that names no `deploy` queues only teardowns: no deploy, no promotion, and no bring-up, and a deploy or promotion queued before `deploy` was removed does not run. The pass reads the block at every look, so a story in review is deployed at the first look after a `deploy` is set; an acceptance made meanwhile promotes nothing (S-0586, [ADR-0259](../adrs/0259-enabling-kluster-on-a-project-with-no-kluster-block-writes-one-naming-only-its.md)).

**The looks.** The pass looks at once and then every minute. Each look queues jobs:

| It finds | It queues |
|----------|-----------|
| A story in review whose branch's head, `refs/heads/story/S-nnnn` in the main checkout, is neither its namespace record's commit nor the head the pass last tried, or that is back in review after a send-back, at any head | A deploy at that head. A deploy waiting for the story at an older head takes the new one in its place |
| A story namespace with a record, whose story is in any column but review, archived ones included | A teardown, unless one failed while the story stood in that column, or a send-back keeps the namespace (S-0694) |
| A story moved to done since the last look, while the action was on | One promotion, after the teardowns. A promotion already waiting is not queued again |

- The pass finds acceptances with the curate scheduler's `news` and `accepted`, from its own start. So an acceptance made before `flai serve` started, or while the action was off, promotes nothing.
- A namespace whose story flai cannot find is left as it is.
- ADR-0292 changes the looks, not built yet (S-0660): a namespace whose story is in progress is not torn down, the pass queues no deploy for a story that enters review with its record passed at its head, and the agent's requests join the queue as jobs of their own ([deploys a story's agent asks for](#deploys-a-storys-agent-asks-for-adr-0292)).

**The worker.** One worker per project runs the queue in order, one job at a time. A look never waits for a job. Each job runs through `KlusterRunner`: `SystemKluster` calls `Cluster.Deploy`, `Teardown`, `Promote`, `Status`, and `Up` from `flai/internal/kluster`, as the commands do, and the tests give their own (`Options.Kluster`).

- **A deploy** first checks that the story is still in review and has a worktree. When the cluster has no kubeconfig, or `kluster status` does not list it running, the pass brings it up, as `flai kluster up` does. The output goes to `.flai-cache/kluster/up.log`, bounded by `kluster.timeout`. A bring-up the dashboard runs counts as busy. The pass then deploys from the worktree at the queued head.
- **A teardown** runs from the story's worktree, or the main checkout once the worktree is gone. A story with no record is not torn down.
- **A promotion** promotes the main branch's head as it is when the job runs, so promotions queued together deploy main's newest head once.
- **The lock.** A job refused with `kluster.ErrBusy`, because `flai kluster` in a shell or the dashboard holds the lock, is no failure. The next look queues it again.

**Not retrying a failure.** The pass keeps what it last tried for each story namespace in `.flai-cache/kluster/serve/<namespace>.json`: the head, the outcome, and, after a failed teardown, the story's column then.

- A head is tried once, passed or failed, even when a refusal wrote no record, until the story returns to review after a send-back.
- Since S-0694 the file also keeps the send-back, `sent_back`: the head that failed, the thread, the move's reason, and its time. The next deploy writes the file afresh, without it.
- A failed teardown is tried again only once the story moves to another column. A story back in review clears it.
- A failed promotion waits for the next acceptance.
- The file is removed after a teardown that passes, and for a namespace with no record once its story has left review, unless a send-back keeps it.

**The checks.** A checks run against the namespace runs in the same job as the deploy, before the next job. It runs only when the deploy passed, `Options.Checks` (`checksConfig`, wired by `flai serve`) has the `checks` host action on, and checks are named. It runs through `Dir.RunChecksWith`, with the namespace's variables. `KlusterEnv` decides those variables for every checks run, the pass's, `flai checks run`'s, and the dashboard's alike. It answers `Cluster.Env` for the story's namespace only when all of these hold:

- the action is on;
- the block is set;
- the namespace's record passed;
- the record's commit is the story branch's head.

Otherwise it answers nothing, and the run is as before. A run with the variables writes `=== against the kluster namespace <ns> ===` at the top of its log, and `flai checks run` says `S-nnnn: checks run against its kluster namespace s-nnnn` on standard error.

Since S-0694 a run that ends keeps its outcome, at the head it ran at, in the namespace's record as `checks_commit` and `checks_outcome` (`kluster.RecordChecks`), which an acceptance reads. Every write of a record takes `.flai-cache/kluster/deployments.lock`, so the outcome never overwrites a deploy's write. An outcome is not kept for a namespace deployed again at another commit while the run ran.

**Failures.** Each kind has one thread by `flai serve`, found again by its title while it is open. A later failure adds an entry to it, and a later success resolves it:

| Failure | Thread title | On |
|---------|--------------|----|
| A deploy, or its bring-up | `flai serve could not deploy S-nnnn into its kluster namespace` | The story |
| The checks run against the namespace | `flai serve's checks against S-nnnn's kluster namespace failed` | The story |
| A teardown | `flai serve could not tear down S-nnnn's kluster namespace` | The story, or `system-flow.yaml` once the story is archived, since a thread on an archived item is a `threads.archived` finding |
| A promotion | `flai serve could not promote main into the kluster's default namespace` | `system-flow.yaml` |

Each entry names the job, the step (or "before its first step" for a refusal), its exit code, the commit, the namespace, the log, and the error.

**The send-back (S-0694).** A failed deploy, its bring-up included, or a checks run that does not pass, sends the story back, as a failing required check does ([ADR-0167](../adrs/0167-flai-serve-reads-each-pull-request-s-required-checks-and-reviews-into-kept.md)). Once the thread is open, the pass reads the project and the story again, and moves the story from review to in-progress as `flai`, with a reason naming the job, the step, and the commit, such as `flai serve's kluster deploy of S-nnnn failed at the deploy step, at commit <commit>`. It moves nothing when the action went off meanwhile, the story is no longer in review, or its branch is no longer at the failed head; the thread's entry says which. The thread's entry says too that the namespace is kept as the failure left it, that the agent fixes the change or says the cause is outside the story, and that cancelling or accepting the story tears the namespace down. A failed teardown or promotion sends nothing back: it has no story in review.

- **The namespace kept.** The send-back is kept in what the pass tried. While the story stands where it put it, in progress, its newest move a move to in-progress by `flai` at or after the send-back, the pass queues no teardown. Once it moves on, or someone else moves it back to in-progress, the namespace is torn down as for any story out of review, unless it is in review again.
- **The agent.** With the `agent` host action on, `flai serve` starts the story's agent once per send-back, while the thread is open, its last run having started before the send-back and ended, as for a story the forge sent back (`startKlusterSentBack` in `serve/agents_restart.go`). It resumes the last run's session where it has one. `harness.SentBack.Kluster` tells the agent that `flai serve` deploys the story again in place once it is back in review, resolves the thread once that deploy and its checks pass, and that a cause outside the story is said on the thread.
- **Back in review.** The story is deployed again in place, at a new head or the same one. A deploy that passes resolves the deploy's thread, and a checks run that passes the checks'. A story cancelled or accepted is torn down, and its failure threads are resolved.

**The acceptance gate (S-0694).** `kluster.AcceptGate(root, story, head, checksOn)` passes only when the namespace's record says its deploy passed at `head`, the story branch's head (`kluster.StoryHead`), and, with `checksOn`, that a checks run against it passed there. Otherwise it names the last job, its commit, and its outcome: never deployed, at another commit, failed at a step, still running, interrupted, or deployed with no checks run ended there yet. While the `kluster` host action is on and the manifest has a `kluster:` block, `preview.AcceptWith` turns its reason into the blocker `kluster`, for every acceptance of a story with a branch that is not a resumed one: `flai accept` and `flai move S-nnnn done`, the dashboard's `accept.preview` and `accept.run`, `--by orchestrator`, among its `orchestrator_blockers`, and the queue's landing, `--batch`. `--dry-run --json` gives `kluster_blocker` and `kluster_why`. `--without-kluster` lifts it for the operator by hand, outside the merge queue, and writes `### Accepted without the kluster gate`, with `- Why:`, `- By:`, and `- At:`, under the story's Notes; it is refused with `--by orchestrator` and with `--batch`, and changes nothing while the action is off or the gate passes ([flai-cli.md](flai-cli.md#commands)). ADR-0292 widens `head` to a record at the head as [what happens when](#what-happens-when) defines it, so a close-out's commit of issue records does not undo the evidence, not built yet (S-0662).

**The merge queue (S-0694).** While the `kluster` host action is on and the manifest has a `kluster:` block, `mergequeue.ComputeWith` lists a story in review that the gate refuses as ineligible, with the reason kind `kluster` naming what failed, and the queue does not batch it. `flai serve`'s queue pass and `flai queue` read the `kluster` and `checks` host actions for it. The orchestrator's enqueue, `queue_add`, is not refused for the gate: the queue waits for it, so the orchestrator may enqueue a story whose deploy has not run yet.

**The journal.** Each job writes two entries to the host journal, one when it starts and one when it ends `done` or `failed`. Each entry has the action `kluster` and the method `serve.kluster`, and names the story, the namespace, the commit, and the failed step. A job refused for the lock ends `failed`, "not run, and queued again at the next look". Each entry has an slog line beside it with `component` `serve`. A send-back adds a `done` entry naming the story, the thread, and the reason.

## The dashboard (S-0580, S-0581)

The rail's Ops group sits between Design and Documentation. Its first page is Klusters, at `/ops/klusters` (S-0580), and its second Deployments, at `/ops/deployments` (S-0581). The page's layout and its routes are in [flaiover-dashboard.md](flaiover-dashboard.md#views).

### The host methods

`flai/internal/hostapi/kluster.go` gives the dashboard these methods. None runs a flai command: each reads the manifest and opens the cluster as `flai kluster` does. The block's keys are written through `settings.manifest`, the config file through `file.get` and `file.save`, and the action turned on through `settings.action` ([ADR-0259](../adrs/0259-enabling-kluster-on-a-project-with-no-kluster-block-writes-one-naming-only-its.md)).

| Method | Kind | What it does |
|--------|------|--------------|
| `kluster.status` | Read, whatever the host action says | The host action, the block, the cluster as `kluster status` reports it, and the up or down flai serve runs or last ran |
| `kluster.up` | Write, under the `kluster` host action | Starts kluster up in the background and answers the job at once |
| `kluster.down` | Write, under the `kluster` host action | Starts kluster down the same way |
| `kluster.deployments` | Read, whatever the host action says | The namespaces' deployment records, main's first, each with its story, head, steps' logs, and how to reach it for a demo; the stories in review with none; and the deploy, teardown, or promotion that runs (S-0581) |
| `kluster.deploy` | Write, under the `kluster` host action, with `{"story": "S-nnnn"}` | Starts what `flai kluster deploy` runs, from the story's worktree at its head, and answers the job at once |
| `kluster.teardown` | Write, under the `kluster` host action, with `{"story": "S-nnnn"}` | Starts what `flai kluster teardown` runs the same way |
| `kluster.promote` | Write, under the `kluster` host action | Starts what `flai kluster promote` runs, main's head into the default namespace, the same way |
| `kluster.workloads` | Read, whatever the host action says | The cluster's nodes and each namespace's pods (S-0586), below |
| `kluster.log` | Read, whatever the host action says, with `{"namespace", "step", "offset"}` | At most 64 KB of a step's log from an offset, for the live feed (S-0595), below |
| `kluster.forward.start`, `kluster.forward.stop` | Write, under the `kluster` host action | Start or stop a service port's port-forward, which flai serve runs bound to `0.0.0.0` (S-0586), below |

`kluster.status` answers these fields. Those marked optional are left out when empty.

| Field | What it holds |
|-------|---------------|
| `enabled` | Whether the `kluster` host action is on for the project |
| `enable` | The command that turns it on, `flai serve enable kluster` |
| `block` | Whether the manifest has a `kluster:` block |
| `no_block`, optional | flai's sentence, when there is no block, saying to turn the `kluster` host action on, which writes one |
| `cluster`, optional | `flai kluster status --json`'s object: `name`, `profile`, `provider`, `context`, `exists`, `running`, `age`, `age_seconds`. Absent with no block or a config flai cannot read. When `kluster status` fails, it holds what the config names, with `exists` false |
| `status_error`, optional | Why the cluster could not be read, such as `kluster` missing from the host's `PATH` |
| `running`, optional | The job that runs now |
| `last`, optional | The last job to end since flai serve started |
| `settings` | The block's keys as `flai manifest set` writes them, `kluster.config` to `kluster.timeout`, each with its `kind` (`text`, `command`, `components`, or `duration`), `default`, `meaning`, `value`, and `set`; listed with no block too, since setting `kluster.config` adds one (S-0586) |
| `settings_on` | Whether the `settings` host action is on, which the block's writes need |
| `settings_enable` | What the operator runs on the host to turn `settings` on, `flai serve enable settings` |
| `config`, optional | The config file the block names |
| `config_file`, optional | What that file names, as kluster reads it: `exists`, `name`, `provider`, `profile`, `default_profile` (`spire`), `profiles` (`spire`, `signet`, `authstar`), `addons`, and `error` when it cannot be read |
| `findings` | `flai check`'s sentences on the block |
| `deployable` | Whether the block names a `deploy` |
| `tools` | `kluster`, `docker`, and `kubectl`, each with `found`, its `path`, what flai `need`s it for, the `install` sentence for a missing one (none for Docker), and `docs`, a link to its install documentation |

A job carries `action` (`up`, `down`, `deploy`, `teardown`, or `promote`), `cluster`, `started_at`, and `log`, the file its output goes to: `.flai-cache/kluster/<action>.log`. A deploy, a teardown, or a promotion also carries `namespace`, `story` when it is a story's, and `logs`, the folder its steps' own logs go to, `.flai-cache/kluster/logs/<namespace>`. Once it ends it carries `ended_at`, `outcome` (`done` or `failed`), and `error` when it failed. `kluster.status` and `kluster.deployments` add `log_tail`, the last 4 KB of the log cut to whole lines, since the dashboard cannot read files under `.flai-cache`. The writes answer the job as it started, with no `log_tail`.

`kluster.deployments` answers these fields. Those marked optional are left out when empty.

| Field | What it holds |
|-------|---------------|
| `enabled`, `enable`, `block`, `no_block` | As `kluster.status` answers them |
| `status_error`, optional | Why the block's config could not be read |
| `cluster`, optional | The cluster's name |
| `default_namespace`, optional | The namespace main is promoted into |
| `deployments` | One row per namespace with a record, main's first, then the stories' by number, then any other by name |
| `review` | The stories in review with no record, each `story`, `title`, `head`, and `head_error` when the head could not be read |
| `running`, optional | The deploy, teardown, or promotion that runs now. When another process holds kluster's lock, such as `flai kluster deploy` in a shell, it is built from the lock's job, with no log |
| `last`, optional | The last deploy, teardown, or promotion this flai serve ran |

A row of `deployments` holds the record's fields, `story`, `namespace`, `commit`, `tag`, `checkout`, `images`, `started_at`, `ended_at`, `outcome`, `failed_step`, `exit_code`, and `error`, with the project's root left out of their paths, and these:

| Field | What it holds |
|-------|---------------|
| `default` | Whether it is the default namespace |
| `job`, optional | The job holding kluster's lock against the namespace: `kind`, `story`, `namespace`, `started_at`, `pid`. A record that says `running` with no such job reads `interrupted` |
| `steps` | Each step's `name`, `started_at`, `ended_at`, `outcome`, `exit_code`, `log`, from the root, and `log_tail`, its last 4 KB |
| `title`, `status`, optional | The story's, for a story's namespace |
| `head`, optional | The story branch's head now, or main's for the default namespace, so a deploy at an older commit shows; `head_error` says why it could not be read |
| `access` | How to reach the namespace for a demo, below |

`access` follows [ADR-0254](../adrs/0254-flai-runs-one-kluster-per-project-from-the-config-its-manifest-s-kluster-block.md): a demo reaches a namespace through a port-forward, until a project gives each namespace a host of its own. flai reads the services with one `kubectl get services --all-namespaces --output json` per list, with its kubeconfig and the cluster's context, within 10 seconds, and only when the kubeconfig exists. Its commands name the kubeconfig from the root, so they are run in the project's main checkout.

| Field | What it holds |
|-------|---------------|
| `kubeconfig` | flai's kubeconfig of the cluster, `.flai-cache/kluster/kubeconfig` |
| `context` | The cluster's kube context, `k3d-` and its name |
| `namespace` | The namespace |
| `list` | The `kubectl` command that lists the namespace's services |
| `services` | Each service's `name`, `host` (`<name>.<namespace>.svc.cluster.local`), `ports` (`name`, `port`, `protocol`), and `port_forward`, the `kubectl port-forward service/<name>` command that forwards each port to the same port on the host |
| `services_error`, optional | Why the services could not be read: no kubeconfig, `kubectl` missing from the host's `PATH`, or kubectl's error |
| `forwards` | The port-forwards flai serve runs, or ran, into the namespace (S-0586), each `namespace`, `service`, `port`, `local_port`, `address` (`0.0.0.0`), `started_at`, `running`, and, once it ended on its own, `ended_at` and the end of its `output` |

`kluster.workloads` answers `block`, `no_block`, `cluster`, `context`, `nodes` (each `name`, `ready`, `roles`, `version`), `namespaces` (each `namespace` and its `pods`, each `name`, `phase`, `ready` and `containers`, and `restarts`), and `workloads_error`, why they could not be read: no kubeconfig yet, `kubectl` missing, or kubectl's error. It runs `kubectl get nodes` and `kubectl get pods --all-namespaces`, each `--output json`, with flai's kubeconfig and the cluster's context, within 10 seconds, only when the kubeconfig exists.

`kluster.forward.start` runs `kubectl port-forward --address 0.0.0.0` for one service port of a namespace, on a local port flai picks, and `kluster.forward.stop` stops it ([ADR-0259](../adrs/0259-enabling-kluster-on-a-project-with-no-kluster-block-writes-one-naming-only-its.md)). Both take `{namespace, service, port}` and answer the forward. A start checks that the cluster lists that service port, needs flai's kubeconfig, and answers once kubectl says it is forwarding, or kubectl's output when it exits first or says nothing within 15 seconds; a second start of a running forward answers it and starts nothing. A stop of a port with no forward is refused. flai serve keeps the forwards in memory, on its kluster jobs; they stop when it stops. A teardown through `kluster.teardown` stops its namespace's forwards first. One run by flai serve's pass does not: their kubectl exits when the namespace goes, and the forward is kept as ended, with its output. A running forward is reachable by anything that reaches the host's network, with no sign-in of its own.

`kluster.log` reads `.flai-cache/kluster/logs/<namespace>/<step>.log` ([ADR-0263](../adrs/0263-a-kluster-block-may-list-components-each-built-as-its-own-step-with-its-own-log.md)). It needs no `kluster:` block, since it only reads files. `namespace` is a DNS label, `step` is `build`, `build-<name>`, `load`, `deploy`, or `teardown`, and `offset` is 0 or more; anything else is refused as invalid params. It answers:

| Field | What it holds |
|-------|---------------|
| `namespace`, `step` | As asked |
| `log` | The log's path from the root |
| `exists` | Whether the log is there; without it the text is empty |
| `offset` | Where `text` starts: the offset asked, or 0 after a reset |
| `next` | The offset to ask from next, `offset` and the length of `text` in bytes |
| `size` | The log's size when it was read |
| `text` | At most 64 KB of the log from `offset`, ending on a whole character |
| `reset`, optional | True when the log is shorter than the offset asked, or gone while the offset asked was above 0, as when a new job cleared it; `text` then starts at 0 |
| `outcome`, optional | The step's outcome in the namespace's record |

### The jobs

- flai serve runs the jobs in its own process, one per project at a time, each bounded by `kluster.timeout`. An up, a down, a deploy, a teardown, and a promotion share that one slot.
- A deploy, a teardown, or a promotion also runs under kluster's own lock, `.flai-cache/kluster/job.lock`, which the commands and flai serve's own kluster jobs take too, so none of them runs beside another.
- It keeps the running job, the last up or down, and the last deploy, teardown, or promotion in memory, apart, so the Klusters page's last job stays the last up or down. A restart of flai serve forgets them; the log files stay.
- Stopping flai serve cuts a running job short, and it ends `failed`.
- The start of each job goes in the host journal. Its end is an info log line, `kluster job ended`.

### Refusals

The writes refuse in this order, each with flai's own sentence:

| Case | Answer |
|------|--------|
| The host action is off | `Disabled`, naming `flai serve enable kluster` |
| A flai that runs no kluster jobs | A rule refusal saying flai serve runs them, or to run `flai kluster up` or `down` in a shell |
| No `kluster:` block | A rule refusal, the sentence `flai kluster up` gives, which says to turn the `kluster` host action on |
| A config flai cannot read | A rule refusal, the kluster package's error |
| An up or down already runs | A rule refusal naming the job, the cluster, and when it started |

A deploy, a teardown, and a promotion refuse in this order:

| Case | Answer |
|------|--------|
| The host action is off | `Disabled`, naming `flai serve enable kluster` |
| A story that is not an ID, for a deploy or a teardown | Invalid params, before the manifest is read |
| A flai that runs no kluster jobs | A rule refusal saying flai serve runs them, or to run `flai kluster deploy S-nnnn`, `teardown`, or `promote` in a shell |
| No `kluster:` block | A rule refusal, the sentence `flai kluster up` gives |
| A config flai cannot read | A rule refusal, the kluster package's error |
| No `deploy` command, for a deploy or a promotion, or no kubeconfig | A rule refusal, as `flai kluster deploy` gives it |
| A job holds kluster's lock | A rule refusal naming the job, the namespace, and when it started |
| A deploy of a story not in progress or in review, or with no worktree | A rule refusal, as `flai kluster deploy` gives it |
| A teardown of a story with no deployment recorded | A rule refusal |
| A promotion with no default namespace | A rule refusal naming `kluster.namespace` |
| A job already runs in flai serve | A rule refusal naming the job, the cluster, and when it started |

A host without `kluster` on its `PATH` is not refused: the job starts and ends `failed` with kluster's install sentence, which `status_error` shows beforehand.

### The Klusters page

The page holds everything setting up and running the project's cluster needs, with no command or file edit outside the dashboard but the two host steps below ([ADR-0259](../adrs/0259-enabling-kluster-on-a-project-with-no-kluster-block-writes-one-naming-only-its.md)). Its **Set up** panel lists five steps, each saying Done or To do. It is open until every step is done, then folded away:

1. **The host's tools.** `kluster`, `docker`, and `kubectl`, each found with its path or missing, with what flai needs it for, the sentence that installs it, and a link to its documentation. flai installs none of them.
2. **Turn kluster on.** A button that turns the `kluster` host action on through `settings.action`. For a project with no block it says what flai wrote and the commit.
3. **The cluster.** The config file's name, provider, profile, and addons; a choice of profile, saved by rewriting the file's `profile:` line; the whole file, edited in place; and the `kluster.config` path. The file is saved through `file.save`, whatever `settings` says.
4. **Deploying stories.** The setup form, `DeploySettings`, which the Deployments page's Setup pane shows too, so the two pages cannot drift (S-0595, [ADR-0263](../adrs/0263-a-kluster-block-may-list-components-each-built-as-its-own-step-with-its-own-log.md)). It has fields for `namespace`, `deploy`, `teardown`, and `timeout`, and a choice between one build command, `build`, and one build per component, `components`, edited as rows of a name and its command, with add, remove, and move up and down. Commands are one argument a line. A save goes through `settings.manifest`, sets the chosen build key and unsets the other in the same call, and shows flai's refusal beside the field it names. The choice is hidden when flai lists no `kluster.components`, as an older flai does. With no `deploy`, the form says the cluster can run but no story is deployed yet. **Write these with an agent** makes a backlog story, under an open epic the operator picks, to write the project's build, deploy, and teardown scripts and name them in the block, and links it; with components, the story names them. Below the form a link opens the Deployments page's Setup pane.
5. **Bring it up.** **Up**.

Below it the page shows the cluster's name, profile, provider, state (Running, Stopped, or Not created), and age, with **Up**, offered only while the cluster is not running, **Down**, which asks for confirmation first, **Turn kluster off**, a link to Deployments, `flai check`'s findings on the block, and the nodes and pods from `kluster.workloads`, read when the page opens and after each job. While a job runs the page says so, shows the end of its log, and asks for the status again every 3 seconds until it ends. Then it shows how the last job ended. flai's refusals and errors are shown word for word.

The two steps left on the host are turning `settings` on and installing the tools. With `settings` off, one notice, **One step on the host**, says what the action allows and shows `settings_enable`. The steps that write the action or the block are shown but cannot be saved; the config file still can. With the action off, the managing view points to step 2.

### The Deployments page

Since S-0595 the page groups what it shows in an in-page rail, the `PageSidebar` the item page uses, with one pane per entry, kept in the URL's hash ([ADR-0263](../adrs/0263-a-kluster-block-may-list-components-each-built-as-its-own-step-with-its-own-log.md)). It reads `kluster.deployments`, and `kluster.status` for the settings and `deployable`.

Above the rail a row of badges says whether the `kluster` host action is on, the cluster's name, whether deploys are set up, and whether a job runs. Below them come flai's errors word for word, the notice that the host action is off, linking the Klusters page's setup, and the job that runs, with **Follow it**, which opens Logs at its namespace, or the job a shell holds kluster's lock for, or the last job and its outcome. With no `kluster:` block the page shows only a notice linking the Klusters page's setup, with no rail, and a flai too old shows the notice the Klusters page shows.

| Pane | Hash | What it shows and does |
|------|------|------------------------|
| Overview | `#overview`, the default | Main's namespace first, then each story's, then the stories in review not yet deployed. Each row shows the story's ID and title, linking its item page, or `main`; the namespace; the commit, and the head when it differs; when it was deployed; the outcome and the step it failed at; and a badge per step with its outcome, which opens Build for a build step and Logs for any other. The entry counts the rows |
| Build | `#build` | A choice of namespace, then one panel per build step of its record, `build` or `build-<name>` per component, in order, each with its status and its live feed; a line counting them, such as "3 builds: 1 passed, 1 running, 1 waiting"; then the load step, with the images it loaded. A deploy that built nothing says so and links Setup. The entry counts the build steps |
| Logs | `#logs` | The running or last job's own log path and the end of its log, then one panel per step of the chosen namespace, each with its live feed and its log's path |
| Access | `#access` | Each namespace's `access`: each service's host and ports, with **Start** and **Stop** for each port's forward and a link to it once it runs, and `services_error`. The pane says once, while a forward runs, that it is reachable from the host's network, with no sign-in of its own ([ADR-0259](../adrs/0259-enabling-kluster-on-a-project-with-no-kluster-block-writes-one-naming-only-its.md)). The entry counts the running forwards |
| Setup | `#setup` | The setup form the Klusters page's step 4 shows, `DeploySettings`, which writes the new status back into the page after a save, and a link to the Klusters page for the cluster and its config file |

The chosen namespace is the one the reader picked, else the running job's, else main's, else the first. Overview's actions:

| Row | Its actions |
|-----|-------------|
| Main's namespace | **Promote** |
| A story's namespace | **Teardown**, after confirmation; **Deploy** when the story is in review and is deployed at a commit other than its head, or its deploy failed or was interrupted |
| A story in review not deployed | **Deploy** |
| Main not deployed | **Promote**, beside the sentence saying so |

With no `deploy` set, read from `kluster.status`'s `deployable`, Overview says deploys are not set up yet, and **Deploy** and **Promote** give way to **Set up deploys**, which opens Setup and asks flai nothing; **Teardown** stays. While a job runs the actions wait, and the page asks for the deployments again every 3 seconds until it ends. flai's refusals and errors are shown word for word.

**The live feed.** Each step's panel, `LogFeed`, shows the step's name, a badge of its outcome, its duration or when it started, and its exit code. It starts folded. Opened, it reads the step's log through `kluster.log`, from 0 and then from each answer's `next`, again at once while `next` is below `size`, about every second while the step runs, and once more when it ends. A `reset` clears the text. It keeps the newest 256 KB, cut at a whole line, saying earlier lines are in the log on the host, and stays scrolled to the end unless the reader scrolls up. After a failed read it tries again every 5 seconds while the step runs. Folding it or leaving the page stops the reads.
