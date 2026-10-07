---
id: T-1238
type: task
nature: feature
title: The operator and design documentation say what the release stamp is, where each build carries it, and what it does not prove
status: backlog
parent: S-0235
owner: alex
created: 2026-10-07T22:49:20Z
updated: 2026-10-07T22:49:20Z
transitions: []
stream: S-0235
tags: [cli, dashboard, docs]
touches: [docs/operators/index.md, docs/operators/settings.md, docs/users/flai.md, design/system/flai-cli.md, design/system/release-signing.md, design/tech/docker.md]
after: [T-1236, T-1237]
---
# T-1238 The operator and design documentation say what the release stamp is, where each build carries it, and what it does not prove

## Work

Criterion 5, and the documentation the code tasks change. It waits for T-1236 and T-1237, so that it describes the statement, the variables, the file, and the build arguments as they were built.

- `docs/operators/index.md`, under `## Security posture` or a section of its own beside it: what the stamp is (CI signs `<component> <version> <commit>` with the release key before the build; flai carries it in `buildinfo`, the image as a file), that each component checks its own at start and warns when a versioned build has none, and what it does not prove, in the words of `design/system/release-signing.md § Verifying the peer`: it names the peer's release and refuses builds that are not one; it does not prove that the peer is that release against anyone who holds the agent credential and a release binary, and the credential stays the boundary of the channel. Update the `## Telemetry` row for `flaiover_build_info` with its `signed` label.
- `docs/operators/settings.md`: a row for the variable that names the stamp file (T-1235), and the `FLAIOVER_VERSION`, `FLAIOVER_COMMIT` row naming `signed`.
- `docs/users/flai.md` and `design/system/flai-cli.md`: what `flai version` and `flai version --json` show of the stamp.
- `design/tech/docker.md`: the image row names the stamp file and the two build arguments, and that `latest` from `main` and `flai dashboard --build` carry none.
- `design/system/release-signing.md`: under the stamp's section or the decision table, how it was built (S-0235).

## Done when

- The operator documentation says what the stamp is and what it does not prove, in the words of `release-signing.md § Verifying the peer`.
- Every changed setting and output is documented, and the markdown lint and `flai check --strict` pass on the changed files.

## Notes

Drafted by the planner.
