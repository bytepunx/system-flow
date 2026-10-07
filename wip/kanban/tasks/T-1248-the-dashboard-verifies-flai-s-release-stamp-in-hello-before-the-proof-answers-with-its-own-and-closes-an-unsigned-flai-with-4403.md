---
id: T-1248
type: task
nature: feature
title: The dashboard verifies flai's release stamp in hello before the proof, answers with its own, and closes an unsigned flai with 4403
status: backlog
parent: S-0237
owner: alex
created: 2026-10-07T23:00:26Z
updated: 2026-10-07T23:00:26Z
transitions: []
stream: S-0237
tags: [dashboard]
touches: [flaiover/src/lib/server/agent.ts, flaiover/src/lib/server/agent.test.ts, flaiover/src/lib/server/release.ts, flaiover/src/lib/server/release.test.ts]
---
# T-1248 The dashboard verifies flai's release stamp in hello before the proof, answers with its own, and closes an unsigned flai with 4403

## Work

The dashboard side of criteria 1 and 4. It waits for no task of this story, and shares no path with the flai side: the two meet only in the wire format the story's criterion 1 and `release-signing.md § A release stamp in hello` fix, and the stamp file and its check are S-0235's (`flaiover/src/lib/server/release.ts`).

- In `flaiover/src/lib/server/release.ts`, add a check of a peer's stamp beside the one S-0235 wrote for the dashboard's own: the signature with the embedded public key, and the component the statement names (`flai`). Tests in `release.test.ts`: a valid peer stamp, another key, a wrong component.
- In `flaiover/src/lib/server/agent.ts`, bump `PROTOCOL` to 2 and accept a protocol 1 `hello` far enough to judge its stamp, so that an older flai is closed as unsigned, not as malformed.
- In `AgentRegistry.handshake()`, verify `hello.release` before the credential proof. A missing, unverifiable, or wrong-component stamp closes the socket with a new `CLOSE_UNSIGNED` (4403, reason `unsigned peer`). Otherwise answer with the dashboard's own `release` beside its nonce and proof.
- Keep per project what the last refusal was: when, the reason, and the version that flai named. Keep it after the socket closes, since the refused flai is never adopted. Mark each adopted connection signed. `status()` and the registry's list expose both, for the task that shows them.
- Tests in `agent.test.ts`: a signed flai is adopted and marked signed; an unsigned one, a wrong component, and a protocol 1 `hello` without the field are each closed with 4403 before any proof is checked, and the refusal is kept against the project.

## Done when

- The hub closes with 4403 every `hello` whose stamp is missing, invalid, or names another component, before the proof, and answers a signed one with its own stamp.
- The refusal and the signed mark are kept per project and readable from the registry.
- `agent.test.ts` and `release.test.ts` pass with `flai test` on the changed paths.

## Notes

Drafted by the planner. The names of the kept fields are the story's agent's to choose; the UI task reads them.
