---
id: S-0837
type: story
nature: improvement
title: A publish the orchestrator starts is a release_publish call and makes no story, and the baseline publish formula is retired
status: backlog
owner: alex
created: 2026-11-17T08:21:17Z
updated: 2026-11-17T08:21:17Z
transitions: []
tags: [flai, template]
topics: [release, orchestration]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
---
# S-0837 A publish the orchestrator starts is a release_publish call and makes no story, and the baseline publish formula is retired

## Goal

[ADR-0316](../../../design/adrs/0316-a-publish-the-orchestrator-starts-is-a-release-publish-call-and-makes-no-story.md): a publish the orchestrator starts makes no story. Today each one runs the baseline `publish` formula (ADR-0240), which made 51 stories in this repository between 2026-11-07 and 2026-11-17, 131 USD and 2.5 hours of agent time, each taking an in-progress slot, a close-out, a review slot, and an acceptance though its branch changes nothing. The operator asked on 2026-11-17 that this stop. Accept the ADR, or record what the operator changes in it, as this story's first task.

## Acceptance criteria

- [ ] The ADR is accepted, ADR-0240 and ADR-0273 name it as what supersedes them, and no accepted ADR is otherwise edited.
- [ ] With `publish` on and the policy met, the orchestrator's publish is one `release_publish` call: no story, task, branch, or worktree is made, shown by a test of the orchestrator's publish that reads the board before and after.
- [ ] `release_publish` fetches the remote's branch and tags, evaluates the policy again, publishes only when it is met and something is pending, and says which release covers each component when another publish got there first, changing nothing.
- [ ] `release_publish` no longer refuses for a formula named `publish`, flai refuses the orchestrator `flai formula run` for every formula, and the `push` host action still gates the pending publish for an agent flai serve started.
- [ ] The template and this repository no longer hold `design/workflows/formulas/publish.md`; `flai upgrade` removes a project's unchanged copy and reports, without removing, one with project additions; `scripts/fetch.sh` goes with it where nothing else uses it.
- [ ] The orchestrator's role file, in the baseline, the template, and this repository, tells it to publish with `release_publish` only, and the harness definition is generated again from it.
- [ ] The conventions (`work-management.md`, `strategic-agents.md`), the design (`strategic-agents.md`, `flai-cli.md`, `project-manifest.md`), the user and operator guides, and the formulas' README no longer describe a publish as a story, in the template's copies too.
- [ ] S-0811's criterion that names the publish formula is rewritten to speak of `release_publish`, by a thread to the operator or the orchestrator if this story's agent may not edit it.

## Tasks

## Notes

- The stopgap until this lands, if the operator wants it sooner: remove `design/workflows/formulas/publish.md` from this repository's main checkout. The orchestrator's role file already says to call `release_publish` where `flai formula list` shows no `publish`, and `release_publish` already works with no formula. Check first that `flai check` and `flai upgrade` do not put the file back or report it missing.
- The archived publish stories stay as they are.
