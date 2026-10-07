---
id: T-1234
type: task
nature: feature
title: flai carries its release stamp in buildinfo, shows it in flai version --json, and flai serve warns at start without a valid one
status: backlog
parent: S-0235
owner: alex
created: 2026-10-07T22:48:34Z
updated: 2026-10-07T22:48:34Z
transitions: []
stream: S-0235
tags: [cli]
touches: [flai/internal/buildinfo/buildinfo.go, flai/internal/buildinfo/stamp.go, flai/internal/buildinfo/stamp_test.go, flai/cmd/version.go, flai/cmd/cmd_test.go, flai/cmd/serve.go, flai/cmd/serve_test.go]
---
# T-1234 flai carries its release stamp in buildinfo, shows it in flai version --json, and flai serve warns at start without a valid one

## Work

The flai side of criteria 1, 3, and 4. It waits for no task: the public key it verifies with is S-0232's `flai/internal/buildinfo/releasekey.go`, which is on main before this story starts.

- Add two linker variables beside `Version`, `Commit`, and `Date` in `flai/internal/buildinfo/buildinfo.go`: the release statement (`flai <version> <commit>`, one line, no newline) and its signature (base64, as `cosign sign-blob` prints it). Both empty by default.
- In a new `flai/internal/buildinfo/stamp.go`, verify a statement and signature against a public key with the standard library: ECDSA P-256, ASN.1 signature over the SHA-256 of the statement's bytes (`ecdsa.VerifyASN1`), and check that the statement names the expected component, flai's own `Version`, and its `Commit`. Take the key as a parameter so that tests can use a key pair they generate; the build checks with the embedded release key. Return one state: `signed`, `missing`, or `invalid` with the reason (bad signature, another key, wrong component, version or commit that does not match).
- A development build is one whose `Version` is `dev`, empty, or starts with `0.0.0` (`scripts/flai.sh`, `goreleaser --snapshot`); it carries no stamp and is not warned about.
- `flai version --json` gains `release`: the statement, the signature, and the state; a build with none shows the state `missing` and stays valid JSON for existing readers. The plain output gains one line saying whether the build is a signed release.
- `flai serve` (`flai/cmd/serve.go`, before `serve.Run`) checks its own stamp once and logs one `warn` with `component: serve` when a versioned build's stamp is missing or invalid, naming the reason.
- Tests in `flai/internal/buildinfo/stamp_test.go`: a valid stamp signed by a test key, a missing one, one signed by another key, and a wrong component; `flai/cmd/cmd_test.go` extends `TestVersionPlainAndJSON` for the new field; `flai/cmd/serve_test.go` shows the warn for a versioned unstamped build and none for a development build.

## Done when

- `flai version --json` from a test build with a stamp shows the statement, the signature, and `signed`, and from `scripts/flai.sh` shows no stamp.
- The stamp tests pass for a valid stamp, a missing one, and one signed by another key, with `flai test` on the changed paths.
- `flai serve` logs exactly one `warn` at start for a versioned build without a valid stamp, and none for a development build.

## Notes

Drafted by the planner. The variable names, the JSON shape, and the state names are the story's agent's to choose; the next task's `ldflags` use the variable names this one picks.
