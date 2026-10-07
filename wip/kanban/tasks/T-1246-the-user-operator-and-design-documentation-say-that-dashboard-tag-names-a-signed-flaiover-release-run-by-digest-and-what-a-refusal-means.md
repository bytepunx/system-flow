---
id: T-1246
type: task
nature: feature
title: The user, operator, and design documentation say that dashboard.tag names a signed flaiover release run by digest, and what a refusal means
status: backlog
parent: S-0236
owner: alex
created: 2026-10-07T22:54:36Z
updated: 2026-10-07T22:54:36Z
transitions: []
stream: S-0236
tags: [docs]
touches: [docs/operators/settings.md, docs/operators/runbooks/update.md, docs/users/flai.md, docs/users/flai-reference.md, design/system/flai-cli.md, design/system/flaiover-dashboard.md, design/tech/docker.md]
after: [T-1244, T-1245]
---
# T-1246 The user, operator, and design documentation say that dashboard.tag names a signed flaiover release run by digest, and what a refusal means

## Work

Describe what T-1241 to T-1245 built, in the files that say how the dashboard's image is chosen and run:

- `docs/operators/settings.md`: the two `dashboard.tag` rows, the configuration's and the manifest's. `latest` is the newest flaiover release, `X.Y.Z` is that release, and either is run by the digest its signed list names.
- `docs/operators/runbooks/update.md`, its flaiover section: `latest` no longer follows the main branch; a release with no verifiable list is refused; `--image` and `--build` run unsigned.
- `docs/users/flai.md`: the `dashboard.tag` row, and `flai dashboard`, `status`, `check`, `upgrade`, and `versions`: resolution, the refusal and its reasons, the digest and signed state in status, and the deployable mark.
- `docs/users/flai-reference.md`: regenerate with `make flai-reference` when a help text changed.
- `design/system/flai-cli.md`: the dashboard commands' rows and "Listing and installing a published release", including where the verified lists are cached.
- `design/system/flaiover-dashboard.md` § Deploying a chosen release: a release that is not deployable.
- `design/tech/docker.md`: its Run row says the image is run as `ghcr.io/bytepunx/flaiover@sha256:…`.

Waits for T-1244 and T-1245, the last of the behaviour it describes; T-1244 waits for T-1243 and T-1245 for T-1242.

## Done when

- Each file above says what `dashboard.tag` now means, how the image is run, and what a refusal means, and none still says `latest` follows the main branch or that the image is run by tag.
- `docs/users/flai-reference.md` matches the help text.
- `flai test` passes on the changed files, the markdown lint and the settings index test among them.

## Notes
