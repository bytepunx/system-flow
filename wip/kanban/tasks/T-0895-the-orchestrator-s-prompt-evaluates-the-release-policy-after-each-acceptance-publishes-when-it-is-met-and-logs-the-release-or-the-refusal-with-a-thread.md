---
id: T-0895
type: task
nature: feature
title: The orchestrator's prompt evaluates the release policy after each acceptance, publishes when it is met, and logs the release or the refusal with a thread
status: backlog
parent: S-0222
owner: alex
created: 2026-10-05T04:46:53Z
updated: 2026-10-05T04:47:05Z
transitions: []
stream: S-0222
tags: [flai]
touches: [flai/internal/harness/harness.go, flai/internal/harness/harness_test.go, ".claude/agents/orchestrator.md", template/root/.claude/agents/orchestrator.md]
after: [T-0886]
---
# T-0895 The orchestrator's prompt evaluates the release policy after each acceptance, publishes when it is met, and logs the release or the refusal with a thread

## Work

Add publishing to the orchestrator's prompt, which S-0218 writes in `flai/internal/harness/harness.go` beside `planPrompt`. Add the same rules to its definition in `.claude/agents/orchestrator.md` and in the template's copy, `template/root/.claude/agents/orchestrator.md`. With `publish` on, whenever `wait_for_events` reports a story moved to `done`:

1. Call `release_evaluate`.
2. Under `threshold` or `theme`, when the evaluation is met, call `release_publish` with the figure that was met as its reason.
3. Under `judgement`, publish only when it judges the unreleased work coherent and complete, and give that reasoning as the reason.
4. Never publish a batch that the evaluation holds back under `whole_epics`, under any policy.
5. Log each release with `activity_log`: the policy, its figures, the versions and tags, and the items bundled, all from what `release_publish` returned.
6. Log a decision not to publish too, with the figures.
7. When `release_publish` refuses, log the refusal with `activity_log`. Open a thread to the operator with `thread_open`, on the most recently accepted story of the batch, with the refusal's words and what fixes it. Do not try again until that thread is answered or the next acceptance.

With `publish` off, it neither evaluates nor publishes. Publishing goes only through `release_publish`, since the guard refuses its other routes.

This task waits for T-0886, whose rules the prompt states. It runs with T-0891, whose paths it does not share.

## Done when

- a harness test pins that the orchestrator's prompt names `release_evaluate`, `release_publish`, `whole_epics`, the log of the release and of the refusal, and the thread on a refusal
- `.claude/agents/orchestrator.md` and `template/root/.claude/agents/orchestrator.md` are the same, as the planner's two copies are, and `diff` between them prints nothing
- `go test ./internal/harness/` passes

## Notes
