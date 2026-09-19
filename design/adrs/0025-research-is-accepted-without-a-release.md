---
id: ADR-0025
title: A research story is accepted and pushed without a release; an experiment stays on its branch
status: accepted
date: 2026-09-19
supersedes: []
superseded_by: []
---

# ADR-0025 A research story is accepted and pushed without a release; an experiment stays on its branch

## Context

The git convention said that `research` and `experiment` stories "stay on a branch and get their additional testing there; they do not release from main", and the tooling enforced it bluntly: `release.LevelFor` returned an error for both natures, `release.Compute` called it first, and so `flai accept`, `flai move <story> done`, the dry run, and the dashboard's acceptance preview all refused such a story. The only way through was `flai accept --no-release` from a shell, which the board cannot do.

That rule was written with a prototype in mind: code that is tried on a branch and may never be kept. It does not fit what a research story produces here. S-0052, the story that raised this, is research whose findings are a decision, an ADR or an edit to the living design, and the stories that follow from it. Those belong on main, where the next agent reads them, and the operator said so on 2026-09-18: "research stories should be able to get moved to done. their findings can produce ADRs, other stories, etc. and those should result in a push even if there is no release to cut." On 2026-09-19 the operator blocked S-0052 with the reason "research stories cannot presently be moved to done".

The two natures were asked about separately. For `experiment` the operator chose on 2026-09-19 to leave the rule as it is.

## Decision

A `research` story is accepted like any other story: its branch is merged, it is archived, the acceptance is committed, and the commit is pushed. It cuts no release, whatever files its commits touched. A finding is a deliverable of the repository, not of a component, so there is no version for it to bump: no tag is created, and no component's version file or changelog changes.

When a research story's commits did touch a component's files, the release plan names each such component as landing on main without a release. It does not refuse. The operator accepting the story sees that code is reaching main unreleased, and the next story that delivers to that component releases it, as it would release any earlier unreleased commit.

An `experiment` story stays on its branch. Acceptance refuses it before anything is merged, with a message that names the nature and the rule, and `flai accept --no-release` remains the deliberate way to land one.

The push follows the acceptance whether or not a release was cut. That was already what `flai accept` did; it is now part of the rule and pinned by a test. From the dashboard's container, which has no git credentials, the existing behaviour stands: accepted locally, with the push command shown. Whether the dashboard should ever push is S-0052's question.

The pre-release pattern the convention mentions (`1.3.0-rc.1`) stays undefined.

## Consequences

- `release.LevelFor` returns a no-release level for research and an error for experiment; `release.Compute` returns a plan whose `skipped` reason says research releases nothing and whose `unreleased` list names the components touched. Everything downstream of a skipped plan (no tags, no version bump, no template publish, push as usual) already existed for stories that touch no component.
- The baseline sentence in `design/conventions/git.md` changes, in the template and here, so projects made from the template get the rule with their next upgrade.
- Code can now reach main through a research story without a release. The plan says so at acceptance; nothing prevents it. If that turns out to be abused, the fix is a check that a research story's component changes are tests or scaffolding, not a return to refusing the nature.
- S-0052 can be accepted from the board once this is released and the dashboard runs the new flai.
