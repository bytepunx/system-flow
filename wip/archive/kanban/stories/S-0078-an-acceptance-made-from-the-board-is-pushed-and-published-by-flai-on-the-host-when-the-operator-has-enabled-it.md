---
id: S-0078
type: story
nature: feature
title: An acceptance made from the board is pushed and published by flai on the host, when the operator has enabled it
status: done
parent: E-0003
owner: alex
created: 2026-09-20T07:26:51Z
updated: 2026-09-20T13:41:50Z
transitions:
  - to: ready
    at: 2026-09-20T07:52:55Z
    by: alex
  - to: in-progress
    at: 2026-09-20T13:05:03Z
    by: system-flow
  - to: review
    at: 2026-09-20T13:25:50Z
    by: system-flow
  - to: done
    at: 2026-09-20T13:41:50Z
    by: alex
tags: [cli, dashboard]
touches: [flai/cmd, flai/internal, flaiover/src, design/system, docs/operators]
---
# S-0078 An acceptance made from the board is pushed and published by flai on the host, when the operator has enabled it

## Goal
The first host action. After an acceptance from the board, flai on the host pushes `main` and the release tags and publishes the template when its version moved, with the operator's own git credentials, because the operator enabled that action. The board's "accepted, not pushed" state becomes "pushing", then pushed or the reason it was not.

## Acceptance criteria
- [x] Host actions are off by default and enabled by name in the host's flai configuration; a disabled action answers with what to enable; nothing about an action can be configured from the dashboard
- [x] The push is `flai push --pending` and the template publish is what `flai template push` does: never forced, refused when the remote has moved, with the refusal shown on the board
- [x] Every host action is journalled on the host: what, which project, asked by whom, when, the outcome; `flai serve` has a command to read the journal
- [x] The board and the story's page show pushing, pushed with the tags, or not pushed with the reason and the command to run by hand; the standing notice from S-0063 stays for acceptances made while the action was disabled
- [x] The operators' documentation says what enabling it means: a holder of the dashboard token can then publish any story an agent has put in review

## Tasks
- T-0286 Host actions are off until enabled by name on the host, and an acceptance from the board pushes nothing unless push is (I-0028)
- T-0287 push.run pushes pending acceptances and publishes the template, and every host action is journalled
- T-0288 The board and the story's page show pushing, pushed with the tags, or not pushed with the reason
- T-0289 The design and the operators' documentation say what enabling push means, and correct what S-0075 and S-0077 said
- T-0290 End to end with scratch remotes: disabled, enabled, refused, and the journal

## Notes
From S-0071's finding, `design/system/dashboard-host-channel.md`, and ADR-0029. Depends on the writes story. Replaces the push key (ADR-0026), which the mount story retires; the two can land in either order.
