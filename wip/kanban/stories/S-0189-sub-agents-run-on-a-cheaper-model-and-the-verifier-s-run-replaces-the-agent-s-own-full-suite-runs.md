---
id: S-0189
type: story
nature: improvement
title: Sub-agents run on a cheaper model, and the verifier's run replaces the agent's own full-suite runs
status: backlog
owner: arobson
created: 2026-10-01T09:19:43Z
updated: 2026-10-01T10:11:43Z
transitions: []
tags: [template, cli]
topics: [conventions]
touches: [template/root/.claude/agents/, ".claude/agents/", flai/internal/harness, design/conventions/delegation.md, template/root/design/conventions/delegation.md, design/system/agent-context.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
---
# S-0189 Sub-agents run on a cheaper model, and the verifier's run replaces the agent's own full-suite runs

## Goal

S-0188 measured two stories run with the released delegation prompt (S-0184, S-0185) against comparable runs that did not delegate (`design/system/agent-context.md § Sub-agents › Measured`). Sub-agents took 10 to 12% of each run's tokens, and no delegating run cost less, finished sooner, or kept a smaller context than its comparables. The story's agent still ran the whole test suite three times itself, before and after its verifiers ran it too. What the verifier bought was defects found before review. Keep that, and stop paying twice for it: run the explorer and the verifier on a cheaper model, and have the verifier's run replace the agent's own full-suite runs.

Add the backing data structure fields to support agent/model configuration based on its role: "story" (what exists now), "explore", and "verify" (with future additions possible).

## Acceptance criteria
- [ ] The template's `explorer` and `verifier` definitions, and this repository's copies, name a model cheaper than the story's agent's, and `agent-context.md § Sub-agents` says which and why
- [ ] `delegation.md` and the `claude-code` prompt (`harness.delegation`) tell the story's agent to run only the tests for what it changed, and to leave the whole suite, lint, and `flai check` to one verifier before review, plus one more after fixing what that verifier found
- [ ] When the agent and/or model configured for the verifier step differs from the one configured for the story, do not allow the verifier agent to perform any necessary corrections found by the verifier.
- [ ] Two stories run with a flai that includes this change are measured as S-0188 measured S-0184 and S-0185, against the same comparables, and the table in `agent-context.md` shows whether the cost of delegation fell below the cost of not delegating

## Tasks

## Notes

- Follow-up of S-0188. The tables and the reading are in `design/system/agent-context.md § Sub-agents › Measured`.
- The defects S-0184's and S-0185's verifiers found are the benefit to keep. Check that a cheaper verifier still finds them, for instance by replaying its prompt on S-0185's branch before the fixes.
