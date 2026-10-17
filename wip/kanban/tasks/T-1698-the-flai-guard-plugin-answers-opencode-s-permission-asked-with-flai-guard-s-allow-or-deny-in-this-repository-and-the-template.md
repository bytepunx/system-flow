---
id: T-1698
type: task
nature: feature
title: The flai-guard plugin answers OpenCode's permission.asked with flai guard's allow or deny, in this repository and the template
status: backlog
parent: S-0408
owner: alex
created: 2026-10-16T05:16:38Z
updated: 2026-10-17T02:39:22Z
transitions: []
stream: S-0408
tags: [cli, template]
touches: [".opencode/plugins/flai-guard.js", template/root/.opencode/plugins/flai-guard.js, template/CHANGELOG.md]
after: [T-1696]
---
# T-1698 The flai-guard plugin answers OpenCode's permission.asked with flai guard's allow or deny, in this repository and the template

## Work

- Under hold-and-ask: in both copies of `flai-guard.js`, on `permission.asked` run `flai guard --harness opencode` with the request, wait for its answer within the ADR's bound, and reply to OpenCode with allow or deny. Under `--auto` with deny rules there is no `ask` to answer: leave the plugin's handling of it out and say so in a comment.
- Keep the two copies the same, and name the change in `template/CHANGELOG.md`.
- Waits for T-1696: the plugin reads the answer that task prints.

## Done when

Both copies are the same, and the plugin's reply to a recorded `permission.asked` input is the allow or deny `flai guard` prints for it. The run against a real `opencode` is T-1699's.

## Notes

- Done when reworded by planner-S-0408, as the orchestrator approved on TH-0490: a task should not depend on the integration run a later layer builds.
