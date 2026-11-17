---
id: T-3582
type: task
nature: remediation
title: Close I-0364 with the reason that flai manifest set now writes agent.image.layer and agent.image.args
status: backlog
parent: S-0839
owner: alex
created: 2026-11-17T12:30:08Z
updated: 2026-11-17T12:30:08Z
transitions: []
stream: S-0839
tags: [flai]
touches: [design/issues/I-0364-no-flai-command-writes-agent-image-args-so-a-project-s-layer-build-arguments.md, design/issues/summary.md]
after: [T-3577, T-3581]
---
# T-3582 Close I-0364 with the reason that flai manifest set now writes agent.image.layer and agent.image.args

## Work

Close I-0364 in the story's worktree: `flai issue close I-0364 --reason "<reason>"`. The reason names what fixed it. `flai manifest set agent.image.args.<NAME>=<value>` and `agent.image.layer=<path>` write the keys, and `--unset` removes them, checked by the manifest's image validation (S-0839). The test that reproduces the issue is in `flai/cmd/manifest_test.go`. The command regenerates `design/issues/summary.md`.

This task waits for T-3577 and T-3581: the reason names the test T-3577 writes, and the issue closes only once the documents T-3581 changes no longer say the keys are written by hand.

## Done when

- I-0364's status is closed with that reason, and `design/issues/summary.md` no longer lists it.
- `flai check --strict` is clean for the story.

## Notes
