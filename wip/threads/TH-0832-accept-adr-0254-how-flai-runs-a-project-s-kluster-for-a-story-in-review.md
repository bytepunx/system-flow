---
id: TH-0832
title: "Accept ADR-0254: how flai runs a project's kluster for a story in review"
anchor:
  path: wip/kanban/stories/S-0576-an-adr-settles-how-flai-runs-a-project-s-kluster-for-a-story-in-review-one.md
  item: S-0576
status: resolved
participants: [agent-S-0576, orchestrator, alex]
created: 2026-11-06T23:29:42Z
updated: 2026-11-07T00:10:35Z
---

# TH-0832 Accept ADR-0254: how flai runs a project's kluster for a story in review

On wip/kanban/stories/S-0576-an-adr-settles-how-flai-runs-a-project-s-kluster-for-a-story-in-review-one.md.

## Entries

### 2026-11-06T23:29:42Z agent-S-0576
**Recommendation:** accept ADR-0254 as written. Reply "accept", or name the points by number to change first.

ADR-0254 is proposed on `story/S-0576`, at `design/adrs/0254-flai-runs-one-kluster-per-project-from-the-config-its-manifest-s-kluster-block.md`. It is decided from S-0575's results. S-0577 to S-0581 build on it, and T-2470 (`design/system/kluster.md`) and T-2471 (the tooling convention) wait for your answer.

## What it decides

| § | Decision |
|---|----------|
| 1 | One cluster per project, from `kluster.config` as it stands on main; provider `k3d` only. flai keeps its own kubeconfig under `.flai-cache/kluster/`, never in `~/.kube/config`. |
| 2 | A story deploys into its lower-case ID, `s-0576`; main into `kluster.namespace`, default the manifest's name. Prefix for shared-service names: `s0576_`, empty for main. |
| 3 | `kluster:` keys: `config` and `deploy` required; `namespace`, `build`, `teardown`, `timeout` (30m) optional. Argument lists, no shell. |
| 4 | Commands get `FLAI_STORY`, `FLAI_KLUSTER_NAME`, `_NAMESPACE`, `_PREFIX`, `_COMMIT`, `_TAG` (`s-0576-<12 hex>`), `_IMAGES`, `_STATE`, `_CONTEXT`, and `KUBECONFIG`. |
| 5 | Jobs in `flai serve`, one at a time per project: deploy on entering review or a new head; teardown on leaving review; on acceptance, teardown then promote main's head to the default namespace. A failed job is recorded only: no thread, no move. |
| 6 | flai loads images itself, `docker save` into `ctr import` on each k3d node, until kluster gains an import. |
| 7 | Host action `kluster`, off by default; turning it off tears nothing down. `flai kluster` writes work from your shell regardless; `flai guard` refuses them to every agent. |
| 8 | With `checks` on too, flai serve starts the checks run after a passing deploy, and the run gets §4's variables. Kluster checks run one at a time per project. |

Consequences list seven proposals for kluster's repository: an image import, `status --json`, namespace add/remove for the shared services, a pinned signet, a Dex for authstar's tenant, JSON progress from `up`, and a profile deploy into a named namespace.

## Points you may want to change

1. **§1, bring-up on demand.** With the action on, flai serve brings the cluster up itself (about 11.5 minutes) when a deploy finds it down. Alternative: deploys fail until you run `flai kluster up`.
2. **§5, a failed deploy gates nothing.** Alternative: a failed deploy sends the story back to in-progress, as a failing required check does under pull-request acceptance.
3. **§5, one job at a time per project.** Safe for a shared build cache and authstar's fixed ports, but two stories entering review together wait on each other.
4. **§2, the default namespace defaults to the manifest's name.** Alternative: make `kluster.namespace` required.

Source: design/experiments/S-0575-two-builds-of-authstar-run-side-by-side-in-one-kluster-each-in-a-namespace.md § Recommendation

### 2026-11-06T23:30:54Z orchestrator (recommendation)
**Recommended:** accept ADR-0254 as written, and keep all four points as they stand.

I read the ADR on `story/S-0576`. Each decision follows S-0575's recommendation to adapt: one cluster per project, a namespace per story named for its serial, and a deploy that belongs to the project, not to kluster.

| Point | Keep it? | Why |
|-------|----------|-----|
| 1. Bring-up on demand | Yes | Turning on the `kluster` host action is your consent (§7). It is off by default, and flai serve never takes the cluster down. |
| 2. A failed deploy gates nothing | Yes | S-0575's baseline failed because of drift between authstar, signet, and kluster, not because of any story. Sending a story back for that would punish the wrong change. |
| 3. One job at a time per project | Yes | S-0575 found that authstar's checks use fixed local ports and leak port-forwards. Builds also share one Docker cache. |
| 4. The default namespace defaults to the manifest's name | Yes | For authstar this gives `authstar`, the namespace its profile deploys into. `flai check` reports a default that is not a DNS label or that looks like a story namespace. |

**Why this is marked for you alone.** The ADR lets flai serve bring up and change shared infrastructure on your machine (§1, §5, §6). `safety.md` asks you to confirm that kind of change. S-0576's first criterion also asks for acceptance by you on a thread. So flai serve should not confirm this after `confirm_after`.

**Next.** Reply "accept", or name the points to change. Then S-0576's agent goes on with T-2470 and T-2471.

For the operator alone: flai serve does not confirm this after orchestration.permissions.confirm_after.

Source: design/experiments/S-0575-two-builds-of-authstar-run-side-by-side-in-one-kluster-each-in-a-namespace.md § Recommendation

### 2026-11-06T23:57:20Z alex
1 - yes, bring up on demand is correct 2 - a failed deploy is necessary when using a cluster since we can't verify the changes without it. it may punish the story but it also alerts the operator to the fact that something is amiss and needs to be corrected before shipping broken code (why ship if it could be bad) 3 - a single cluster per project makes sense. systems like authstar provide the notion of tenants so that each story could be automated to use a tenant and then tare it down once the story is clear; this keeps story "environments" separate within one cluster sharing infrastructure outside the project itself. 4 - manifests that are going to work with flai and kluster in this context need to allow for customizable namespaces that have a default but can be overriden via parameters.

### 2026-11-06T23:59:36Z agent-S-0576
ADR-0254 now follows your four answers. It is still proposed. Reply "accept" to accept it as it now stands, or say what differs from what you meant. It is on `story/S-0576`, at `design/adrs/0254-flai-runs-one-kluster-per-project-from-the-config-its-manifest-s-kluster-block.md`.

| Point | Your answer | What ADR-0254 says now |
|-------|-------------|------------------------|
| 1 | Bring up on demand | Unchanged (§1). |
| 2 | A failed deploy must hold the story | New §5a. A failed deploy, or a failed checks run against the story's namespace, opens one thread to you on the story and sends it back to in-progress, starting its agent, as a failing required check does (ADR-0167). The namespace stays for diagnosis until the story is redeployed, cancelled, or accepted. While `kluster` is on, every acceptance (yours, the orchestrator's, the merge queue's) is refused with the blocker `kluster` until the deploy and its checks passed at the story's head. `flai accept --without-kluster` lands it anyway, on your word only, recorded in the story's Notes, for a cause outside the story. A failed promotion opens a thread to you on `system-flow.yaml`. |
| 3 | One cluster per project; a tenant per story | Unchanged: one cluster, one job at a time per project. §2 adds the story's tenant, where the project has tenants, among the settings its deploy makes and its teardown removes. |
| 4 | Namespaces customisable, with a default, overridable by parameter | §2 and §3: a project's manifests and scripts take the namespace as a parameter, with the default namespace as the default, and never name it as a constant. The default namespace defaults to the manifest's name, and `kluster.namespace` overrides it. |

The two details I added under point 2 are the `--without-kluster` flag and keeping the namespace of a story that was sent back. Say if you want either one otherwise.

Meanwhile I am writing `design/system/kluster.md` (T-2470) and the tooling convention (T-2471) to match this text.

Source: design/experiments/S-0575-two-builds-of-authstar-run-side-by-side-in-one-kluster-each-in-a-namespace.md § Recommendation

### 2026-11-07T00:00:38Z orchestrator (recommendation)
**Recommended:** reply "accept". ADR-0254, as revised, matches each of your four answers (D-0481). Keep both details the agent added.

I read the revision in S-0576's worktree. It is not committed yet: `story/S-0576` at 21cd0d060 still has the first text. The agent's close-out commits it.

| Point | Your answer | Where the revision says it |
|-------|-------------|----------------------------|
| 1 | Bring up on demand | §1, unchanged. |
| 2 | A failed deploy holds the story | §5a. A failed deploy or checks run opens a thread to you and sends the story back. Acceptance is refused with the blocker `kluster`. The rejected alternative records your reason. |
| 3 | One cluster, a tenant per story | §2. The project's deploy makes the story's tenant and its teardown removes it. |
| 4 | Namespace as a parameter, with a default | §2 and §3. Manifests and scripts take the namespace as a parameter. `kluster.namespace` overrides the default. |

**The two added details.**
- `flai accept --without-kluster`: keep it. It is yours alone, as `--without-queue-tiers` is (ADR-0221). It is your only way past a cause outside the story, such as baseline drift, short of turning the action off.
- Keeping a sent-back story's namespace: keep it. Without it, neither you nor the agent can see what failed.

**Know this before you accept.** Under §5a, until a project's baseline passes in its default namespace, every story that enters review is sent back. The Consequences say so. Turn the `kluster` action on for authstar only once its baseline passes there.

**Why this is yours alone.** S-0576's first criterion asks for your acceptance on a thread. The ADR also lets flai serve change shared infrastructure on your host.

For the operator alone: flai serve does not confirm this after orchestration.permissions.confirm_after.

Source: design/digest/D-0481-accept-adr-0254-how-flai-runs-a-project-s-kluster-for-a-story-in-review-1-yes.md

### 2026-11-07T00:09:46Z alex
Confirmed the recommendation of 2026-11-07T00:00:38Z orchestrator.

Source: design/digest/D-0481-accept-adr-0254-how-flai-runs-a-project-s-kluster-for-a-story-in-review-1-yes.md

### 2026-11-07T00:10:35Z agent-S-0576
Resolved: The operator accepted ADR-0254 as revised to their four answers, keeping --without-kluster and the namespace a sent-back story keeps (D-0481). ADR-0254 is accepted on story/S-0576 at 4a7cf8f56.
