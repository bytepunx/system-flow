---
id: S-0087
type: story
nature: feature
title: "Moving a story to done merges it; a publish button on the board tags and pushes everything accumulated since the last one"
status: done
parent: E-0003
owner: alex
created: 2026-09-21T00:13:58Z
updated: 2026-09-21T22:09:00Z
transitions:
  - to: ready
    at: 2026-09-21T00:14:18Z
    by: system-flow
  - to: in-progress
    at: 2026-09-21T03:29:49Z
    by: alex
  - to: review
    at: 2026-09-21T04:45:25Z
    by: system-flow
  - to: done
    at: 2026-09-21T22:09:00Z
    by: alex
tags: [cli, dashboard]
touches: [flai/cmd, flai/internal/release, flaiover/src, design/system, design/adrs, docs/operators]
---
# S-0087 Moving a story to done merges it; a publish button on the board tags and pushes everything accumulated since the last one

## Goal
Today, accepting a story computes its component's bump, tags, and (when the push host action is enabled) pushes, all as one inseparable step. A run of small stories against the same component each cuts its own release, so releases come out smaller than they need to be. Moving a story from review to done should merge its branch into main and nothing else; a new Publish action, reachable from the board's done column, computes and creates the tag(s) for everything merged and unreleased since the last publish, then pushes main and the tags together.

`flai accept`, `release.Compute`/`Apply`/`Tag` (`flai/internal/release`), and `flai release <id>` already separate these steps as distinct function calls glued together by `acceptItem`; `flai release` already computes and tags a single item without pushing. This story splits accept's merge from its tag-and-push, and extends the release machinery from one item to everything pending.

## Acceptance criteria
- [x] `flai accept` merges the story's branch into main, moves it to done, and archives it as before, but computes no bump, creates no tag, and pushes nothing
- [x] A component's bump for a publish is the highest delivery type among everything merged and unpublished for it since the last publish (an epic done outranks a feature story, which outranks a remediation/improvement/patch-level change), not one bump per story
- [x] A new `flai release --pending` (or equivalent) computes every component's pending bump across everything unpublished, creates the tag(s), commits any version-file bumps, and pushes main and the tags together, batched at most three tags per push (`pending.Batches`, per I-0026) so a large batch cannot silently produce no release workflow run
- [x] The done column shows, per card, whether it has been published or is merged and waiting, and a Publish action at the top of the column runs the batch above; disabled or absent when nothing is pending. Reachable only when the push host action is enabled, the same gate acceptance's automatic push used before this story
- [x] Research stories still land on main with no release of their own; an experiment is still refused at accept. Cancelling, blocking, and the rest of the workflow are unaffected
- [x] The changelog entry for a batched release names every story it covers, not just one
- [x] A publish that fails partway (one component's tag made, the next fails; or the push itself fails) leaves the state such that running Publish again finishes what is left, without re-tagging what already succeeded
- [x] An ADR refines ADR-0019 (which established tag-at-accept as the baseline) and cross-references ADR-0025 (research/experiment) and the push host action (ADR-0031, S-0078): what changed, why, and what a component's version means when it was last cut from a batch rather than one story
- [x] `docs/operators` and any user-facing documentation describing acceptance and release are updated

## Tasks
- T-0313 flai accept merges and moves to done; it computes no bump, creates no tag, pushes nothing
- T-0314 A component's pending bump is the highest delivery type among everything merged and unpublished for it since its last tag
- T-0315 A publish step tags and pushes everything pending, batched three tags per push, resumable on partial failure
- T-0316 The done column shows published versus merged-and-waiting, with a Publish action gated on the push host action
- T-0317 An ADR records the move from tag-at-accept to a batched publish, and the conventions and documentation are corrected
- T-0318 Tried end to end: a batch of several stories against one component, published once

## Notes
Raised by the operator on 2026-09-20: releases were coming out "too small" because each accepted card cuts its own tag immediately. Discussed two decisions before this story was written: a batch's bump is the highest delivery type among what it bundles, not a sum or one-per-story; and Publish always releases everything currently merged and unpublished, with no per-story selection — the only way to hold a merged story back from a release is not to move it to done yet. `flai release <id>` already exists (git.md's note that it "computes the bump, tags, and pushes once it exists" is stale: it computes and tags today, it just does not push). I-0026's three-tags-per-push limit, already handled by `pending.Batches`, becomes load-bearing under this design rather than an edge case, since a batched publish routinely produces more tags at once than a single acceptance did.
