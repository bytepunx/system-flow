---
id: D-0481
title: "Accept ADR-0254: how flai runs a project's kluster for a story in review: 1 - yes, bring up on demand is correct 2 - a failed deploy is necessary when using a cluster since we can't verify the…"
kind: recommendation
sources: [TH-0832, S-0576]
topics: [kluster, testing, deployment, manifest, all]
paths: [design/adrs, design/adrs/README.md, design/system/kluster.md, design/system/README.md, design/conventions/tooling.md]
date: 2026-11-06
status: absorbed
absorbed_by: "design/adrs/0254-flai-runs-one-kluster-per-project-from-the-config-its-manifest-s-kluster-block.md#decision"
asked_again: []
cited: [S-0576, S-0694, TH-0832, TH-1036]
---

# D-0481 Accept ADR-0254: how flai runs a project's kluster for a story in review: 1 - yes, bring up on demand is correct 2 - a failed deploy is necessary when using a cluster since we can't verify the…

The operator, alex, answered TH-0832 at 2026-11-06T23:57:20Z.

What was asked:

```text
**Recommendation:** accept ADR-0254 as written. Reply "accept", or name the points by number to change first.
```

The recommendation it overruled:

```text
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
```

The operator's answer:

```text
1 - yes, bring up on demand is correct 2 - a failed deploy is necessary when using a cluster since we can't verify the changes without it. it may punish the story but it also alerts the operator to the fact that something is amiss and needs to be corrected before shipping broken code (why ship if it could be bad) 3 - a single cluster per project makes sense. systems like authstar provide the notion of tenants so that each story could be automated to use a tenant and then tare it down once the story is clear; this keeps story "environments" separate within one cluster sharing infrastructure outside the project itself. 4 - manifests that are going to work with flai and kluster in this context need to allow for customizable namespaces that have a default but can be overriden via parameters.
```

The operator, alex, answered TH-0832 again at 2026-11-07T00:09:46Z.

The answer:

```text
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
```

Source: TH-0832.
