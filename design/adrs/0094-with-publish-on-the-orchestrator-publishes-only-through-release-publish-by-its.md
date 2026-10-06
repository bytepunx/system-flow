---
id: ADR-0094
title: "With publish on, the orchestrator publishes only through release_publish, by its release policy, never a batch whole_epics holds back, and raises each refusal on a thread"
status: accepted
date: 2026-10-06
supersedes: []
superseded_by: []
refines: [ADR-0067, ADR-0087]
topics: [cli, orchestration, release]
---

# ADR-0094 With publish on, the orchestrator publishes only through release_publish, by its release policy, never a batch whole_epics holds back, and raises each refusal on a thread

## Context

[ADR-0067](0067-accepted-work-reaches-the-remote-only-when-it-is-published-and-agents-publish.md) made publishing, `git fetch` then `flai release --pending`, the one way accepted work reaches the remote, and said agents publish only when the operator asks. It left room for a later agent, enabled by configuration, to publish. [ADR-0087](0087-flai-serve-runs-one-orchestrator-per-project-behind-the-orchestrate-host-action.md) gave the orchestrator the permission `orchestration.permissions.publish`, off by default, and its guard let the orchestrator run `flai release --pending` and `flai push` while it is on. Nothing said when the orchestrator should publish, what it must check first, or what it does when publishing is refused.

S-0217 gave the orchestrator the figures: `flai release --evaluate` says whether `orchestration.release` is met, under `threshold`, `theme`, or `judgement`, the default, which is never met by itself. Publishing acts with the operator's credentials, and the dashboard's Publish (`publish.run`) runs it only while the operator has the `push` host action on. A release can also bundle a story whose epic is still being built, which an operator may not want released before the epic is whole. And `flai release --pending` refuses when the remote has release tags this clone lacks (S-0174) or its branch has moved; an agent that retries or works around that refusal could push over the operator's remote.

## Decision

With `orchestration.permissions.publish` on, the orchestrator publishes, and only through the MCP tool `release_publish`, by the project's release policy. The permission is the operator's asking that ADR-0067 requires.

1. **When it evaluates.** After each acceptance, its own or one `wait_for_events` reports as a story moved to done, the orchestrator calls `release_evaluate`. While `publish` is off it neither evaluates nor publishes.
2. **What each policy publishes.** Under `threshold` or `theme`, it publishes when the evaluation is met, with the figure that was met as its reason. Under `judgement`, it publishes when it judges the unreleased work coherent and complete, and gives that reasoning as its reason. Each release publishes everything accepted and not yet released, as `flai release --pending` does.
3. **whole_epics.** With `orchestration.release.whole_epics` set, a batch with a story whose epic is in neither review nor done is held back under every policy. `flai release --evaluate` names those stories, with their epics, in `held_by_epic`, and the policy is not met while one is there.
4. **What release_publish checks.** It weighs the batch as `release_evaluate` does, and refuses, changing nothing, when nothing is accepted and not yet released; when `whole_epics` holds the batch back, naming each story held and its epic; under `threshold` or `theme`, when the policy is not met, with its figures; under `judgement`, without a reason; when the `push` host action is off; when the caller is not the orchestrator; and when `publish` is off. Otherwise it runs what `publish.run` runs, `flai release --pending`, under the same `push` host action, never forced.
5. **What is recorded.** The host's journal records the run as `publish.run`'s are, with method `mcp.release_publish`, the orchestrator as who asked, and `<policy>: <reason>` in its detail, so a release the orchestrator cut is told apart from the operator's. The tool returns the policy, the reason, the evaluation, the versions and tags released, and the items bundled. The orchestrator logs each release with `activity_log`, with those figures, and logs each decision not to publish with the evaluation's.
6. **A refusal goes to the operator.** When `release_publish` refuses, the orchestrator logs the refusal and opens a thread to the operator on the most recently accepted story of the batch, with the refusal's words and what fixes it. A refusal by `flai release --pending` with exit 3, for newer remote tags or a moved remote, is returned as flai's message unchanged, beginning `conflict:`. The orchestrator does not try again until that thread is answered or the next acceptance, and never works around the refusal.
7. **The guard is the fence.** Under `FLAI_ROLE=orchestrate`, `flai guard` lets `release_publish` through only while `publish` is on, and refuses the orchestrator `flai release` other than `--evaluate`, `flai push`, `git push`, and `git tag`, whatever its permissions. It refuses `release_publish` to every other session, and the tool refuses any caller that is not the orchestrator as well.

## Consequences

- The operator can let accepted work reach the remote without publishing by hand, and turns that on with one permission and the `push` host action.
- Every release the orchestrator cuts passes the same checks the evaluation reports. Under `threshold` and `theme` the orchestrator cannot publish early; under `judgement` it cannot publish without saying why, and the reason is in the journal and its log.
- `publish` alone publishes nothing: the `push` host action must be on too, so the dashboard token and the orchestrator share one switch for the operator's credentials.
- A conflict with the remote stops publishing until the operator fetches and merges or rebases. Publishing waits on the operator's answer rather than on a retry.
- An operator who sets `whole_epics` gets no release from the orchestrator with a part-built epic in it. The operator's own publishing, `flai release --pending` or the board's Publish, is not held by it.
- ADR-0087's `publish` row changes: the permission no longer allows `flai release --pending` or `flai push` on the command line.

## Alternatives considered

- **Let the orchestrator run `flai release --pending` and `flai push`, as ADR-0087 allowed.** The guard reads a command line, not the release policy, so nothing would stop a release the policy does not allow, and the run would be journalled as the operator's.
- **Publish under `judgement` only when the operator confirms on a thread.** That is the operator publishing. The reason, logged and journalled, keeps the call reviewable.
- **Retry a conflict after fetching.** Fetching and merging change the operator's checkout, which the orchestrator never writes. The thread hands it to the operator.
- **Hold back only the stories of an epic not yet whole, and publish the rest.** `flai release --pending` publishes everything pending, and a release that leaves out stories already merged on the branch cannot be tagged on it.
