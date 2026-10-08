---
id: T-1350
type: task
nature: remediation
title: A local release server answers the three GitHub release endpoints from a GoReleaser dist folder
status: backlog
parent: S-0340
owner: alex
created: 2026-10-08T08:05:44Z
updated: 2026-10-08T08:05:44Z
transitions: []
stream: S-0340
tags: [cli]
touches: [flai/internal/releaseserver/releaseserver.go, flai/internal/releaseserver/releaseserver_test.go, flai/internal/releaseserver/serve/main.go, scripts/release-server.sh]
---
# T-1350 A local release server answers the three GitHub release endpoints from a GoReleaser dist folder

## Work

Add `flai/internal/releaseserver`: an `http.Handler` that serves the releases in a GoReleaser `dist` folder, for a repository name it is given, in the shape GitHub's API answers.

- `GET /repos/<repo>/releases`, with any `per_page`: a JSON list of releases, each with `tag_name`, `draft`, `prerelease`, `published_at`, and `assets` with `name` and `url`, and no `Link` header.
- `GET /repos/<repo>/releases/tags/<tag>`, with the tag's `/` escaped as `%2F` as `install.sh` asks: the one release.
- `GET /repos/<repo>/releases/assets/<n>`: the asset's bytes. Each asset's `url` ends in `/releases/assets/<n>`, which `install.sh`'s `asset_url` matches.
- Anything else answers 404.

Add `flai/internal/releaseserver/serve/main.go`, which listens on `127.0.0.1` on a free port and prints its base URL. Add `scripts/release-server.sh`, which starts it in the background for a calling script, waits until it answers, prints the base URL, and stops it when the caller exits. Use only the standard library's `net/http`: the story allows no new dependency.

Test it with `httptest`: `selfupgrade.Resolve` and `selfupgrade.Download` against it find and verify the release, and an unknown tag or asset answers 404.

Waits for no task: the tests build their own `dist` fixture rather than T-1348's build.

## Done when

- `go test ./internal/releaseserver/...` passes, with `Resolve` and `Download` against the server among its tests.
- `scripts/release-server.sh`, given a folder holding a `flai_<X.Y.Z>_<os>_<arch>.tar.gz` and its `checksums.txt`, prints a base URL that `install.sh` installs from with `FLAI_API` at it and `FLAI_REPO` at the stand-in name.

## Notes

Drafted by the planner from S-0340's criterion 1 and its note that the server serves the shape the real API answers, so `install.sh` and `selfupgrade` run unchanged. The JSON fields are those `selfupgrade`'s `ghRelease` and `ghAsset` decode. A Go server under `flai/` is preferred to a shell one: it needs no interpreter on the host, and a change under `flai/` makes the story's own close-out run the smoke tier.
