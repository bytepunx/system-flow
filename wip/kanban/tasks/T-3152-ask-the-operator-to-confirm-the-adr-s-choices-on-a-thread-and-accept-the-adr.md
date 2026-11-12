---
id: T-3152
type: task
nature: feature
title: Ask the operator to confirm the ADR's choices on a thread, and accept the ADR with their answer
status: backlog
parent: S-0733
owner: alex
created: 2026-11-12T15:55:13Z
updated: 2026-11-12T15:55:21Z
transitions: []
stream: S-0733
tags: [cli]
touches: [design/adrs/README.md, design/system/flai-cli.md]
after: [T-3149, T-3151]
---
# T-3152 Ask the operator to confirm the ADR's choices on a thread, and accept the ADR with their answer

## Work

Open one thread on S-0733 naming the ADR's choices an operator could reasonably make otherwise, each with its recommendation first:

- pull, build, or both;
- the registry and image names;
- consent for the first automatic build;
- the step limit and the retry count;
- whether strategic agents wait for the layer;
- where the operator's own build command is named.

Publishing images uses the release key and a public registry, and an automatic build uses the host's network and minutes, so these are the operator's to confirm.

On the answer, change the ADR and `design/system/flai-cli.md` where they differ. Then accept the ADR with `flai adr accept`, which dates it and writes `design/adrs/README.md` again. It waits for T-3149 and T-3151, so the operator confirms the ADR and the living design together; it shares `design/adrs/README.md` with T-3149 and `design/system/flai-cli.md` with T-3151.

## Done when

- The operator's answer is on the thread, and the ADR says what it decided.
- The ADR's status is accepted, and `design/adrs/README.md` lists it so.

## Notes
