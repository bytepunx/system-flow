---
id: ADR-0292
title: "A story's agent asks flai serve to deploy, check, and tear down its own story's kluster namespace while the story is in progress, through MCP tools and --job flags behind the kluster host action, and waits for a kluster event; the namespace lives from its first deploy until acceptance or cancellation, kluster.agent_namespaces bounds how many stories in progress hold one, and the namespace's record at the story's head is what the move to review, the close-out, and an acceptance read"
status: proposed
date: 2026-11-12
supersedes: []
superseded_by: []
refines: [ADR-0254]
topics: [kluster, deployment, manifest]
---

# ADR-0292 A story's agent asks flai serve to deploy, check, and tear down its own story's kluster namespace while the story is in progress, through MCP tools and --job flags behind the kluster host action, and waits for a kluster event; the namespace lives from its first deploy until acceptance or cancellation, kluster.agent_namespaces bounds how many stories in progress hold one, and the namespace's record at the story's head is what the move to review, the close-out, and an acceptance read

## Context

[ADR-0254](0254-flai-runs-one-kluster-per-project-from-the-config-its-manifest-s-kluster-block.md) deploys a story into its kluster namespace only once the story is in review, and tears the namespace down when the story goes back to in-progress. Its §7 has `flai guard` refuse every agent `flai kluster deploy`, `teardown`, and `promote`, with no other way to ask. So a story's agent has nothing running to verify against while it works. A criterion that needs the whole system running reaches review unticked, and the operator learns there whether it holds. E-0035 wants the agent to deploy and check its story from its worktree before review.

What stands today, which this decision starts from:

- **The send-back is built (S-0694).** A failed deploy or checks run in review opens a thread to the operator, sends the story back to in-progress, starts its agent, and keeps its namespace while the story stands where the send-back put it. Every acceptance is refused with the blocker `kluster` unless the namespace's record passed at the story's head, its checks too while `checks` is on (§5a, `kluster.AcceptGate`). The record keeps `checks_commit` and `checks_outcome`.
- **`deploy` is optional** ([ADR-0259](0259-enabling-kluster-on-a-project-with-no-kluster-block-writes-one-naming-only-its.md)). A block with no `deploy` brings the cluster up and deploys nothing.
- **Components** ([ADR-0263](0263-a-kluster-block-may-list-components-each-built-as-its-own-step-with-its-own-log.md)). A deploy may build several components, one step each, and its record lists every step before the first runs.
- **One job at a time per project.** Deploys, teardowns, and promotions share one queue and one lock (§5, `.flai-cache/kluster/job.lock`). A job is bounded by `kluster.timeout`, 30 minutes by default.
- **The close-out already leaves the agent's turn** ([ADR-0229](0229-the-close-out-runs-as-a-flai-serve-job-the-story-s-agent-asks-for-with-the.md)). The agent asks with the MCP tool `verify` or `flai verify --job`, gets a job id at once, and waits for a `verified` event with `wait_for_events`. That is the pattern a deploy the agent asks for follows.

The cluster is shared infrastructure on the operator's machine. That is why §7 keeps agents from it. What replaces the refusal must keep the operator in charge: flai serve runs every job, the `kluster` host action gates it, and an agent reaches only its own story's namespace.

## Decision

**A story's agent asks flai serve to deploy, check, and tear down its own story's kluster namespace while the story is in progress, through MCP tools and `--job` flags behind the `kluster` host action, and waits for a `kluster` event. The namespace lives from its first deploy until acceptance or cancellation. `kluster.agent_namespaces` bounds how many stories in progress hold one. The namespace's record at the story's head is what the move to review, the close-out, and an acceptance read.**

This refines ADR-0254. It amends §4 (the agent's variables), §5 (the table of jobs, replaced by the table in 3 below), §5a (failures before review, and what counts as at the head), §7 (what the guard refuses), and §8 (checks before review). §1 to §3 and §6 stand.

### 1. How the agent asks

| Request | MCP tool | On the host |
|---------|----------|-------------|
| Deploy the story at its branch's head: build, load, deploy, then its checks run while `checks` is on | `kluster_deploy` | `flai kluster deploy S-nnnn --job` |
| Run the checks again against the namespace as deployed | `kluster_checks` | `flai checks run S-nnnn --job` |
| Tear the namespace down | `kluster_teardown` | `flai kluster teardown S-nnnn --job` |

- **A request file, as ADR-0229's.** Each writes `S-nnnn.<kind>.request.json` under `.flai-cache/kluster/requests/` and waits at most ten seconds for flai serve to take it. flai serve looks for requests every second, beside the close-out jobs, and queues each as a job of §5's queue. It answers with the job's id and the number of jobs ahead of it. A request flai serve does not take in ten seconds is withdrawn, and the agent is told that flai serve is not running. The agent cannot run the job itself: `flai kluster deploy` without `--job` stays the operator's (2).
- **Its own story only.** The story is the session's own: `FLAI_STORY`, else the story in the agent's name, `agent-S-nnnn`. A request naming another story is refused, by the MCP tool and by `flai guard` alike.
- **While the story is in progress.** A story in review is deployed by flai serve at each new head already; a request for it is answered with that job, or with the record when it passed at the head. A story in any other column is refused.
- **At a committed head.** The job builds in the worktree, so a request is refused while the worktree has uncommitted changes to tracked files, naming them: commit first. The job records the head it built. Deploying uncommitted work is out of scope.
- **Gated by the `kluster` host action.** While it is off, or the manifest has no `kluster:` block, every request is refused, naming `flai serve enable kluster` as the operator's. A block with no `deploy` refuses `kluster_deploy` and `kluster_checks`, saying deploys are not set up yet ([ADR-0259](0259-enabling-kluster-on-a-project-with-no-kluster-block-writes-one-naming-only-its.md)). `kluster_checks` is refused while the `checks` host action is off, and before the namespace has a record that passed. The agent then says in the story's Notes which criteria it could not verify and why.
- **One job per story in the queue.** A request while a job of the story waits or runs is answered with that job's id and queues nothing. A deploy request at a newer head takes the place of the one waiting, as the review pass's deploy does. An agent's request always runs: the review pass's rule that a head is tried once does not hold it back.
- **Components.** A deploy runs every component's build, as a review deploy does. Building one component alone is left for a later story.

### 2. What the agent may run against its namespace

- **The project's checks, through flai serve.** The checks run of `kluster_checks`, or the one that follows a passing deploy, is the review's checks run (§8): the named checks, with §4's variables for the story's namespace and flai's whole kubeconfig, in the queue after the job before it. Its outcome and log go in the namespace's record and the story's checks status.
- **Its own commands, with the namespace's variables.** `flai kluster env S-nnnn` prints §4's variables for the story's namespace, as `export` lines or `--json`, and `flai kluster exec S-nnnn -- <command>` runs a command with them, in the worktree. The agent runs its project's own tests, `kubectl get`, `kubectl logs`, or a port-forward with them. Both are reads of flai's state, open to the story's agent and its sub-agents for their own story while its namespace has a record, a failed one included, for diagnosis.
- **A kubeconfig held to the namespace.** In what `env` and `exec` give, `KUBECONFIG` is not flai's kubeconfig for the cluster (§1) but one for the namespace alone: `.flai-cache/kluster/agents/s-nnnn.kubeconfig`. Each deploy makes it after the deploy step passes, from a ServiceAccount `flai-agent` in the namespace bound by a RoleBinding to the cluster's `edit` ClusterRole, with a token from `kubectl create token` valid for 24 hours. So the agent's `kubectl` and its project's scripts read and change what is in its namespace, and the API server refuses them anything in another namespace and anything cluster-wide. The teardown deletes the file. The other variables of §4 are as they are for the deploy, but `FLAI_KLUSTER_IMAGES`, which is the build's alone.
- **Reads of the record and the logs.** `flai kluster deployments`, `flai kluster status`, and `flai kluster log S-nnnn --step <step>`, which prints a step's log as the host method `kluster.log` reads it, are open to every agent for its own story's namespace. An agent in a container reaches the logs only through `flai kluster log`, since `.flai-cache/kluster` is not mounted.
- **What `flai guard` still refuses.** Every agent and sub-agent: `flai kluster up`, `down`, and `promote`; `deploy` and `teardown` without `--job`; `--job` requests, `env`, `exec`, and `log` for another story's namespace or the default one; the `kluster` program itself but `status`; and a shell line that names flai's cluster kubeconfig, `.flai-cache/kluster/kubeconfig`. The API server's refusal through the namespace's kubeconfig is the boundary; the guard's refusals keep an agent from reaching around it from its shell. A script the agent runs is not read by the guard, which is why the kubeconfig, not the guard, holds the namespace.

### 3. When a namespace is made, deployed again, and torn down

This table replaces ADR-0254 §5's, and the rows `kluster.md` § What happens when gained with S-0694's send-back.

| When | Job | Where it runs |
|------|-----|---------------|
| The agent asks for a deploy, its story in progress with no namespace, under the limit (4) | Deploy: the namespace is made | The story's worktree, at its head |
| The agent asks for a deploy, its story's namespace there | Deploy again, in place | The story's worktree, at its head |
| The agent asks for a checks run | Checks run against the namespace as deployed | The story's worktree |
| The agent asks for a teardown | Teardown | The story's worktree |
| A story in progress gains commits, and its agent asks nothing | None: the namespace stays at the commit it was deployed at | Nowhere |
| A story enters review, its namespace's deploy passed at its head, and its checks too while `checks` is on | None: those stand as the review's | Nowhere |
| A story enters review otherwise: no namespace, another commit, or a failure | Deploy, then checks, or the checks run alone when only it has not passed at the head | The story's worktree, at its head |
| A story in review gets a new head | Deploy again, in place | The story's worktree, at its new head |
| A story leaves review for in-progress, sent back by a failure or by anyone | None: the namespace stays | Nowhere |
| A story moves back to ready or backlog | Teardown | The story's worktree, or the main checkout when it is gone |
| A story is cancelled, from any column | Teardown | The story's worktree, or the main checkout when it is gone |
| A story is accepted | Teardown of its namespace, then a promotion | As ADR-0254 §5 has them |
| A job finds the cluster absent | Bring the cluster up first | The main checkout |

- **At the head.** A record is at the head when its commit is the head, or when every commit from the record's commit to the head changes only paths under the manifest's `design`, `docs`, and `wip` folders. So the issues a close-out records and commits do not undo the evidence. This holds for the row above, for the acceptance gate (§5a), and for what the move to review and the close-out report (6).
- **The send-back needs no case of its own.** A namespace now stays while its story is in progress, so the kept namespace of a send-back is the ordinary case. Its threads stay open until a deploy and checks run pass at the story's head, whether the agent or the review asked for them.
- **No deploy.** With no `deploy` in the block, flai serve deploys nothing, at a request or in review, and still tears down a namespace already deployed ([ADR-0259](0259-enabling-kluster-on-a-project-with-no-kluster-block-writes-one-naming-only-its.md)).
- Jobs still run one at a time per project, in the order their causes came, the agent's requests among them.

### 4. How many namespaces live at once

- **`kluster.agent_namespaces`**, an optional key of the `kluster:` block, is the most stories in progress that hold a namespace at once. Default 2. 0 turns deploys before review off, and the project works as ADR-0254 first had it. `flai manifest set kluster.agent_namespaces` writes it, and the dashboard's setup form shows it under the `settings` host action, as the block's other keys. `flai check` reports a value that is not a whole number of 0 or more as `manifest.kluster`.
- **Stories in review do not count.** Their deploys are the acceptance gate's evidence and are never refused for the limit; the review WIP limit bounds them, as ADR-0254 said.
- **The limit makes no namespace; it tears none down.** A story sent back from review holds its namespace whatever the count, and counts while it is in progress. A request for a story whose namespace is there is never over the limit.
- **Over the limit, a request is refused at once,** naming the limit and each story in progress that holds a namespace, with its last deploy's commit and time. It is not queued: a namespace frees only when another story is accepted, cancelled, or torn down, which could take days. The agent goes on with the work that needs no namespace, asks again later, or says in the story's Notes which criteria it could not verify, and the review's deploy checks them.

### 5. Waiting, and failures

- **The answer at once.** A request answers the job's id and the jobs ahead of it. The agent does the work that needs no namespace and holds `wait_for_events` for the event, again each time it returns.
- **The `kluster` event.** When a job the agent asked for ends, flai serve appends a `kluster` event to `.flai-cache/events.jsonl` with the story as its ID and fields `job`, `kind` (`deploy`, `checks`, or `teardown`), `commit`, `outcome` (`passed`, `failed`, `refused`, or `interrupted`), `failed_step`, `exit_code`, and `log`. The inbox reports it to the story's agent and to readers with no story of their own, as it does a `verified` event (ADR-0229 §4). `wait_for_events` does not answer `end` to a story's agent while a job of its story waits or runs, so the session is not ended under it.
- **How long.** A job starts when those ahead of it end, each bounded by `kluster.timeout`, and runs at most `kluster.timeout` itself. A deploy and its checks run are one job. So a request waits at most the jobs ahead, times the timeout, and its own. A job still waiting or running when flai serve stops ends `interrupted` when it starts again, with an event; the agent asks again.
- **A failure before review is the agent's.** The event names the step, its exit code, and its log; `flai kluster log` prints the log, and the namespace stays as the failure left it. flai serve opens no thread, moves nothing, and starts no agent for it. The agent fixes the change and asks again.
- **When a failure reaches the operator.** Only when the agent says on a thread on its story that the cause is outside the story, such as a bring-up that fails, a shared service that drifted, or a baseline that fails in the default namespace, naming the job, the step, and the log. Or at review: a deploy or checks run that fails there sends the story back as §5a has it, with its thread to the operator.

### 6. What is recorded, and who reads it

- **The namespace's record**, `.flai-cache/kluster/deployments/s-nnnn.json`, holds a deploy the agent asked for as it holds a review's: the commit, every step and its outcome and log, and `checks_commit` and `checks_outcome`. It gains `asked_by`, the agent that asked or `flai serve` for the review pass, and `checks_log`, the checks run's log. A checks run the agent asks for writes `checks_commit`, `checks_outcome`, and `checks_log` as the review's does. What the agent runs itself with `exec` is recorded nowhere by flai: its results go in the criteria it ticks and the narrative.
- **The move to review** (`flai move S-nnnn review` and `item_move`) reads the record and says what was deployed and checked at the head, or why not, as it warns of an unticked criterion. It refuses nothing: the review deploys what has not passed at the head.
- **The close-out** (`flai verify`) gains a step, `kluster`, that reads the record against the head and runs nothing. It passes when the action is off, the block has no `deploy`, or `kluster.agent_namespaces` is 0, and otherwise reports the record's commit, outcomes, and whether it is at the head. It fails only on a deploy or checks run that failed at the head: a story whose change the cluster could not run is not ready for review. Not deploying is no failure.
- **An acceptance** reads the record as §5a's `kluster` blocker does, with at the head as 3 says. A deploy and checks run the agent asked for that passed at the head lift the blocker without a review deploy. The review page shows the record (S-0662).

## Consequences

- S-0660 builds the requests, the queue entries, the `kluster` event, the limit and its key, the namespace's kubeconfig, `flai kluster env`, `exec`, and `log`, and the guard's rules of 2. S-0661 builds `kluster_checks` and the checks run before review. S-0662 builds what the move to review, the close-out, and the review page say, and the at-the-head rule in `kluster.AcceptGate`. S-0663 shows the namespaces of stories in progress on the Deployments page. `design/system/kluster.md` describes the whole and marks what is not built.
- More namespaces live at once: one for each story in progress that asked, up to the limit, beside those in review. On the 24-core, 62 GB host of S-0575, two authstar namespaces and the default ran with room to spare, so the default of 2 adds at most that much.
- An agent's deploy waits behind review deploys and promotions in one queue. A review deploy waits behind an agent's in turn. A queue that grows long is the operator's to see on the Deployments page, and to shorten with a lower limit.
- The review deploys less: a story checked at its head before review is not deployed again there. The `kluster` gate reads the same record, so its evidence is the same.
- A namespace now outlives a send-back and stays while its story is in progress, up to acceptance or cancellation. Signet keeps a torn-down namespace's bundles and policies, as ADR-0254 says, and they accumulate no faster: one namespace per story still.
- The namespace's kubeconfig holds a token that works against the local cluster only, which `safety.md` allows. It lasts 24 hours, and a deploy makes a new one. An agent in a container reaches the API server only when its network reaches the host's loopback, where k3d publishes it; S-0660 says how, as [ADR-0269](0269-flai-serve-refuses-a-container-start-whose-http-mcp-server-listens-on-loopback.md) says for the MCP server.
- The `edit` role reaches only the namespace's Kubernetes objects. A project's own commands that reach a shared service, such as Postgres in the profile's namespace, do so over the cluster's network from a pod in the story's namespace, or through the project's checks, which flai runs with the whole kubeconfig.

## Alternatives considered

- **flai serve deploys a story in progress at each new head on its own.** The agent would wait for nothing, but every commit of every story in progress would rebuild and redeploy, minutes each, in one queue. The agent knows when its change is worth deploying.
- **The agent runs `flai kluster deploy` itself.** A deploy takes minutes, past the harness's tool-call limit, and changes shared infrastructure on the operator's machine. A job in flai serve behind the host action, as the close-out is (ADR-0229), keeps both off the agent's turn.
- **Queue a request over the limit until a namespace frees.** A namespace frees only at another story's acceptance or cancellation, so the agent would wait days on a request it cannot see move. A refusal it can act on is better.
- **Evict the namespace deployed longest ago to make room.** It would tear down another agent's environment under it, mid-check.
- **Count review namespaces against the limit.** A story entering review could then not be deployed, and its acceptance would wait on stories in progress.
- **Let the guard alone hold the agent to its namespace.** The guard reads shell lines, not the scripts they run, and a `kubectl` call inside the project's test script would reach every namespace with flai's kubeconfig. Kubernetes's own RBAC, through a kubeconfig for the namespace, holds whatever the agent runs.
- **Run the agent's own commands as flai serve jobs too.** Every `kubectl get` would wait in the queue behind deploys. Reads and the project's tests against one namespace do not race the jobs the way builds and loads do.
- **Tear the namespace down when a story leaves review for in-progress, as ADR-0254 §5 had it.** The agent would lose the environment it is about to fix the change in, and deploy again from nothing.
- **Fail the close-out when the story was not deployed before review.** Not every story needs the running system to verify, and the review deploys it anyway. A failure at the head, though, says the change cannot run, so it stops the close-out.
- **Supersede ADR-0254.** Its cluster, namespaces, block, and image load (§1 to §3, §6) stand unchanged. Refining it, and naming the sections amended, keeps one source for each.
