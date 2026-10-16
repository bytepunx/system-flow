---
id: T-1716
type: task
nature: experiment
title: The target story and the non-Claude model are chosen with the operator, and the results document records the setup
status: backlog
parent: S-0412
owner: alex
created: 2026-10-16T05:20:26Z
updated: 2026-10-16T05:20:26Z
transitions: []
stream: S-0412
tags: [cli]
touches: [design/experiments/S-0412-a-non-claude-model-works-one-story-of-this-repository-end-to-end-on-opencode-through-openrouter-with-the-guard-and-the-permission-prompt-on-and-the-results-record-what-broke-and-what-it-cost.md, design/experiments/README.md]
---
# T-1716 The target story and the non-Claude model are chosen with the operator, and the results document records the setup

## Work

- Open a thread on S-0412, recommendation first: a small open story of this repository, a remediation or improvement with few touches and no `.claude/` path, and an OpenRouter model that is not Claude and supports tool calls. Ask the operator to set the target's `agent` to `harness: opencode`, the host's OpenRouter `provider`, and that model, since a story's agent is the operator's to set.
- Start the results document from `design/experiments/template.md`, named for S-0412, and list it in `design/experiments/README.md`. Record under its setup the model, the `opencode` and flai versions, the host's harness and provider settings by name, never a key, and the target story.
- Waits for nothing in this story.

## Done when

The operator has answered the thread, the target's agent is set, and the results document records the setup.

## Notes
