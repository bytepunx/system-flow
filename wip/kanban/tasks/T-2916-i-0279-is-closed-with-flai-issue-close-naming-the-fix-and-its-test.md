---
id: T-2916
type: task
nature: remediation
title: I-0279 is closed with flai issue close, naming the fix and its test
status: backlog
parent: S-0670
owner: alex
created: 2026-11-10T08:48:03Z
updated: 2026-11-10T08:48:03Z
transitions: []
stream: S-0670
tags: [issues]
touches: [design/issues/I-0279-scripts-agent-container-sh-always-runs-flai-agent-claude-code-local-and-ignores.md, design/issues/summary.md]
after: [T-2911, T-2914]
---
# T-2916 I-0279 is closed with flai issue close, naming the fix and its test

## Work

Criterion 2. Run `flai issue close I-0279 --reason "<reason>"` in the story's worktree, the reason saying that `scripts/agent-container.sh` runs the manifest's layer image, `flai-agent-<key>-claude-code:local`, or the image `FLAI_AGENT_IMAGE` names, that it refuses one not built, and that `scripts/agent-container-test.sh`, run by the smoke tier, reproduces the case. flai regenerates `design/issues/summary.md`.

Waits for T-2911 and T-2914: the issue is closed once the fix is tested and documented.

## Done when

- [ ] `flai issue show I-0279` reads `status: closed` with the reason
- [ ] `design/issues/summary.md` no longer lists I-0279, and `flai check --strict` is clean

## Notes
