---
id: T-1040
type: task
nature: improvement
title: flai self-upgrade --list and flai dashboard versions list the published releases, marking the running and the newest
status: backlog
parent: S-0298
owner: alex
created: 2026-10-06T21:44:53Z
updated: 2026-10-06T21:45:36Z
transitions: []
stream: S-0298
tags: [cli]
touches: [flai/cmd/selfupgrade.go, flai/cmd/selfupgrade_test.go, flai/cmd/dashboard.go, flai/cmd/dashboard_upgrade.go, flai/cmd/dashboard_test.go]
after: [T-1038, T-1039]
---
# T-1040 flai self-upgrade --list and flai dashboard versions list the published releases, marking the running and the newest

## Work

Add the CLI listings the ADR names, both built on T-1039's `List`:

- `flai self-upgrade --list`: the published flai releases, newest first, marking the installed one and the newest, and any below a served project's `flai.minimum`; `--json` prints them as a list with `version`, `tag`, `published`, `installed`, and `latest`.
- `flai dashboard versions`: the published dashboard releases (`flaiover/v`), marking the one the running container's image carries and the configured tag; `--json` the same shape.

`flai self-upgrade --version` and `flai dashboard upgrade --tag` already install a chosen release. Make the first refuse a version that is not published, with the list, through `Resolve`. Give `flai dashboard upgrade` a `--published` flag that refuses a `--tag` that is not a published dashboard release, with the list; the host API passes it on every tag the dashboard names, while a tag the operator gives in a shell, such as a mirror's, is left as it is today. Register `versions` in `flai/cmd/dashboard.go` beside `check` and `upgrade`. Waits for T-1038 for the names and flags, and for T-1039 for `List`.

## Done when

- Both listings print in text and JSON, tested against a stand-in releases API.
- `flai self-upgrade --version` with an unpublished version, and `flai dashboard upgrade --published --tag` with an unpublished tag, fail naming the published ones before anything is pulled or replaced.
- `scripts/flai-test.sh` passes for `flai/cmd`.

## Notes
