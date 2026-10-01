---
id: S-0189
type: story
nature: improvement
title: Sub-agents run on a cheaper model, and the verifier's run replaces the agent's own full-suite runs
status: done
owner: arobson
created: 2026-10-01T09:19:43Z
updated: 2026-10-01T11:36:28Z
transitions:
  - to: ready
    at: 2026-10-01T10:11:55Z
    by: alex
  - to: in-progress
    at: 2026-10-01T10:42:36Z
    by: agent-S-0189
  - to: review
    at: 2026-10-01T11:36:06Z
    by: agent-S-0189
  - to: done
    at: 2026-10-01T11:36:28Z
    by: alex
tags: [template, cli]
topics: [conventions]
touches: [template/root/.claude/agents, ".claude/agents", flai/internal/harness, flai/internal/guard, flai/internal/manifest, flai/internal/workitem, flai/cmd, flai/internal/mcpserver, flai/internal/hostapi, flai/internal/serve, flaiover/src, design/adrs, design/conventions/delegation.md, template/root/design/conventions/delegation.md, design/system/agent-context.md, design/system/conventions.md, design/system/flai-cli.md, design/system/work-hierarchy.md, design/system/project-manifest.md, template/CHANGELOG.md, template/template.yaml, docs/users/flai.md, docs/users/flai-reference.md, docs/users/conventions.md, docs/operators/settings.md, design/issues]
after: [S-0181]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 3195
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 564
      output: 36737
      cache_read: 50047111
      cache_write: 482819
      cost: 23.4656
    - model: claude-sonnet-5-5
      input: 142
      output: 8186
      cache_read: 3654226
      cache_write: 322682
      cost: 2.3571
---
# S-0189 Sub-agents run on a cheaper model, and the verifier's run replaces the agent's own full-suite runs

## Goal

S-0188 measured two stories run with the released delegation prompt (S-0184, S-0185) against comparable runs that did not delegate (`design/system/agent-context.md § Sub-agents › Measured`). Sub-agents took 10 to 12% of each run's tokens, and no delegating run cost less, finished sooner, or kept a smaller context than its comparables. The story's agent still ran the whole test suite three times itself, before and after its verifiers ran it too. What the verifier bought was defects found before review. Keep that, and stop paying twice for it: run the explorer and the verifier on a cheaper model, and have the verifier's run replace the agent's own full-suite runs.

Add the backing data structure fields to support agent/model configuration based on its role: "story" (what exists now), "explore", and "verify" (with future additions possible).

## Acceptance criteria
- [x] The template's `explorer` and `verifier` definitions, and this repository's copies, name a model cheaper than the story's agent's, and `agent-context.md § Sub-agents` says which and why
- [x] `delegation.md` and the `claude-code` prompt (`harness.delegation`) tell the story's agent to run only the tests for what it changed, and to leave the whole suite, lint, and `flai check` to one verifier before review, plus one more after fixing what that verifier found
- [x] When the agent and/or model configured for the verifier step differs from the one configured for the story, do not allow the verifier agent to perform any necessary corrections found by the verifier.
- [x] A story's and the project's `agent` carry an optional `roles` map (`explore`, `verify`, open to more), each with its own harness, model, and config, and `claude-code` starts the story's session with each role's model over its definition's (TH-0057)

## Tasks
- T-0671 The explorer and verifier run on a cheaper model
- T-0672 Delegation and the prompt leave the whole suite to the verifier
- T-0673 A story's and the project's agent carry per-role config
- T-0674 claude-code runs each role's model, and the command harness is told the roles
- T-0675 A verifier never makes the corrections it finds
- T-0676 Replay S-0185's verifier on sonnet against its branch before the fixes

## Notes

- Follow-up of S-0188. The tables and the reading are in `design/system/agent-context.md § Sub-agents › Measured`.
- The defects S-0184's and S-0185's verifiers found are the benefit to keep. Check that a cheaper verifier still finds them, for instance by replaying its prompt on S-0185's branch before the fixes.
- The criterion to measure two stories run with a flai that includes this change moved to S-0190, which is `after` S-0189. It can be measured only after S-0189 is accepted and released (TH-0057, the designer's word). The roles criterion states the goal's second paragraph in the shape TH-0057 settled.
