---
id: T-1235
type: task
nature: feature
title: flaiover reads and verifies its release stamp from a file, labels flaiover_build_info signed, and warns at start without a valid one
status: backlog
parent: S-0235
owner: alex
created: 2026-10-07T22:48:44Z
updated: 2026-10-07T22:48:44Z
transitions: []
stream: S-0235
tags: [dashboard]
touches: [flaiover/src/lib/server/release.ts, flaiover/src/lib/server/release.test.ts, flaiover/src/lib/server/metrics.ts, flaiover/src/lib/server/metrics.test.ts, flaiover/src/hooks.server.ts]
---
# T-1235 flaiover reads and verifies its release stamp from a file, labels flaiover_build_info signed, and warns at start without a valid one

## Work

The flaiover side of criteria 2, 3, and 4. It waits for no task: the public key is S-0232's constant in `flaiover/src/lib/server/release.ts`, on main before this story starts, and it shares no path with T-1234 beside it.

- In `flaiover/src/lib/server/release.ts`, read the stamp from a file in the image, `/app/release-stamp` by default, two lines: the statement `flaiover <version> <commit>` and its base64 signature. An environment variable, such as `FLAIOVER_RELEASE_STAMP`, names another file, for tests and for anyone who mounts one; the task that builds the image uses the same path and format.
- Verify it with `node:crypto` (`crypto.verify('sha256', statement, publicKey, signature)`, the ASN.1 ECDSA P-256 signature `cosign sign-blob` makes), and check that it names `flaiover`, the image's `FLAIOVER_VERSION`, and its `FLAIOVER_COMMIT`. Take the key as a parameter so that tests use a key pair they generate. Return one state: `signed`, `missing`, or `invalid` with the reason.
- A development build is one whose version is unset, `0.0.0`, or starts with `0.0.0-`; it is not warned about.
- `flaiover_build_info` in `flaiover/src/lib/server/metrics.ts` gains a `signed` label, `true` or `false`.
- In `init` in `flaiover/src/hooks.server.ts`, check the stamp once and log one `warn` (`component: 'server'`) when a versioned build has no valid stamp, naming the reason; add the state to the `server started` line.
- Tests in `flaiover/src/lib/server/release.test.ts`: a valid stamp signed by a test key, a missing file, one signed by another key, and a wrong component; `flaiover/src/lib/server/metrics.test.ts` sees the `signed` label.

## Done when

- The vitest tests for a valid stamp, a missing one, and one signed by another key pass, with `flai test` on the changed paths.
- `/metrics` shows `flaiover_build_info` with `version`, `commit`, and `signed` labels.
- The server logs exactly one `warn` at start for a versioned build without a valid stamp, and none for a development build.

## Notes

Drafted by the planner. The file path, its format, and the variable name are the story's agent's to choose; T-1237, which bakes the file into the image, follows what this task picks.
