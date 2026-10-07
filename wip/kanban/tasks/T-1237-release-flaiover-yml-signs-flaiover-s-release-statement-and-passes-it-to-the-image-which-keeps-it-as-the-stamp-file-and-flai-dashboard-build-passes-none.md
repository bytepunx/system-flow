---
id: T-1237
type: task
nature: feature
title: release-flaiover.yml signs flaiover's release statement and passes it to the image, which keeps it as the stamp file, and flai dashboard --build passes none
status: backlog
parent: S-0235
owner: alex
created: 2026-10-07T22:49:09Z
updated: 2026-10-07T22:49:09Z
transitions: []
stream: S-0235
tags: [dashboard, ci]
touches: [".github/workflows/release-flaiover.yml", flaiover/Dockerfile, flai/cmd/dashboard_test.go]
after: [T-1235]
---
# T-1237 release-flaiover.yml signs flaiover's release statement and passes it to the image, which keeps it as the stamp file, and flai dashboard --build passes none

## Work

The build half of criterion 2. It waits for T-1235, whose stamp file path and format the image must write.

- In `.github/workflows/release-flaiover.yml`, on a `flaiover/v*` tag only, sign the line `flaiover <version> <commit>` (no trailing newline) with `cosign sign-blob --key env://COSIGN_PRIVATE_KEY --tlog-upload=false`, the version taken from the tag and the commit the one passed as `FLAI_COMMIT`, and pass the statement and the signature as two build arguments. The key's secrets live in the environment `release`, which only release tags may deploy to (S-0232's `docs/operators/runbooks/release-key.md`): give the job that environment on a tag alone, install cosign as `release-flai.yml` does, and fail a tag build whose key is missing. A push to `main` builds `latest` with no stamp, as ADR-0070 says.
- In `flaiover/Dockerfile`, take the two build arguments in the runtime stage and write them as the stamp file T-1235 reads, only when they are given; with none, write no file. Check that the version `FLAIOVER_VERSION` gets on a tag is the version the statement names.
- `flai dashboard --build` (`buildDashboardImage` in `flai/cmd/dashboard.go`) and `scripts/flaiover-image.sh` pass no stamp arguments; a test in `flai/cmd/dashboard_test.go` asserts that the build's arguments carry none.
- Build the image locally with `scripts/flaiover-image.sh`, and once more with a statement and signature made with a throwaway key, and confirm the first has no stamp file and the second has it where T-1235 reads it.

## Done when

- `release-flaiover.yml` signs the statement on a release tag and passes it and its signature as build arguments; a `main` build carries none.
- The image keeps them as the file the server reads, and an image built without them, as `flai dashboard --build` builds it, has no stamp.
- The `flai/cmd/dashboard_test.go` test passes with `flai test`.

## Notes

Drafted by the planner. `flai dashboard --build` keeps the version `git describe` gives it, so a dashboard built that way warns once at start that it has no stamp; S-0239's `dashboard.allow_unsigned` is where that is made quiet and shown, not this story.
