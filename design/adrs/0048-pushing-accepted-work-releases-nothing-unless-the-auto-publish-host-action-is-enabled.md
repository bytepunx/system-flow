---
id: ADR-0048
title: Pushing accepted work releases nothing unless the auto-publish host action is enabled
status: accepted
date: 2026-09-29
supersedes: []
superseded_by: []
refines: [ADR-0032]
---

# ADR-0048 Pushing accepted work releases nothing unless the auto-publish host action is enabled

## Context

ADR-0032 made publishing a deliberate step over everything accumulated: acceptance merges, and a release is cut when the operator chooses, so several stories batch into one release. S-0094 then had `flai push --pending` compute, apply, and tag that release before every push, so that a release was never a step someone forgot. Pushing is not deliberate in the same way. Agents run `flai push --pending` whenever `inbox` reports an unpushed acceptance, and the board's Push now runs it too. Every acceptance was therefore released on its own soon after it was made, and the operator lost the batching ADR-0032 gave them (S-0144).

## Decision

A push releases nothing unless the operator enables the `auto-publish` host action for the project (`flai serve enable auto-publish`, or its toggle in the dashboard's settings), which is off by default like every host action. Off, `flai push --pending` pushes the merged commits and any tags already made, and tags nothing new; releasing waits for `flai release --pending` or the board's Publish, which stay gated by `push` as ADR-0032 says. On, a push first tags everything accumulated, as S-0094 made it do. Enabling `push` alone never implies a release.

## Consequences

- Acceptances batch again: they reach the remote as they are accepted, and are released together when the operator publishes.
- An operator who wants every push to publish says so once, per project or for every project, and can turn it off again at once.
- `flai push --pending` reads the host configuration it runs with, so an agent's shell and `flai serve` decide by the configuration each uses; both default to off.
- Nothing pushed without a release misses its version change: what a release covers is worked out from local history since each component's last tag, whether or not those commits were pushed, so the next publish includes them (TH-0031). Its tags are then on commits the remote already has, which `pending.Detect` cannot see, so `flai release --pending` and `flai push --pending` push tags they have just made on their own; before this, such a tag was made and left local.
- Once published, a batch's bump is still the highest delivery type among everything accepted since the last tag (ADR-0032). Only when it is cut changes.

## Alternatives considered

- Revert S-0094 outright, so a push never releases. Lost because the operator asked for a setting, and an operator who does want releases at every push would lose that with no way back.
- A key in `system-flow.yaml` rather than a host action. Lost because pushing and publishing act with the operator's credentials on the host, which is what host actions govern (ADR-0029). The manifest is shared by every clone.
- Stop agents pushing after an acceptance as well. Declined by the operator in TH-0031: pushing merged commits cuts no release, so batching does not depend on it.
