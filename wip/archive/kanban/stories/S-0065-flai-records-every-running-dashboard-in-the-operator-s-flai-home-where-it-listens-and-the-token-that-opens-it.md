---
id: S-0065
type: story
nature: feature
title: "flai records every running dashboard in the operator's flai home: where it listens and the token that opens it"
status: cancelled
parent: E-0007
owner: alex
created: 2026-09-20T01:18:58Z
updated: 2026-09-20T06:46:19Z
transitions:
  - to: ready
    at: 2026-09-20T03:31:55Z
    by: alex
  - to: cancelled
    at: 2026-09-20T06:46:19Z
    by: alex
tags: [cli, dashboard]
touches: [flai/cmd, flai/internal, flaiover/src]
---
# S-0065 flai records every running dashboard in the operator's flai home: where it listens and the token that opens it

## Goal
Every flaiover dashboard that `flai dashboard` starts on this host is recorded in the operator's `~/.flai`, one YAML file per project, with what is needed to find it and to talk to it: where it listens and the token that opens it. The record is kept true as dashboards start, stop, restart, and rotate their tokens, so that something else on the host (flairport, S-0067) can list every project without being told about each one.

## Acceptance criteria
- [ ] `flai dashboard` writes `~/.flai/instances/<project key>.yaml` when it starts a dashboard: the project's name and key, the repository's path on the host, the container's name and image, the URL, bind address, and port, when it was started and by which flai version, and the token or the path of the token file (decided in this story with the reason, see the notes). The schema is versioned and documented in `design/system`
- [ ] The directory is 0700 and every file 0600, written atomically (temporary file and rename), and two `flai dashboard` commands running at once for different projects never corrupt or lose a record, which a test shows
- [ ] `flai dashboard stop` marks the instance stopped (or removes the record, decided here and documented), `flai dashboard token --rotate` updates it, and a restart rewrites it; a record whose container is gone (killed, host rebooted) is recognised as stale by whoever reads it and is repaired by the next `flai dashboard`, `stop`, or `status` for that project
- [ ] `flai dashboard list` (the name is open) prints every recorded instance with its state checked against docker, as a table and as `--json`, from any directory
- [ ] flaiover's part is stated and built: the container is not given `~/.flai`, so flai is the only writer; flaiover already answers `/_health` and names its project in every API response (ADR-0024), and anything more a reader of the registry needs from it (version, whether it can write, whether it holds a push key) is added to an endpoint that needs no more than the token
- [ ] The location follows the existing rule for flai's home (`FLAI_CONFIG`'s directory, or an explicit override such as `FLAI_HOME`), so tests and this repository's scripts never touch the operator's real `~/.flai`; the tests use a scratch home
- [ ] The token's exposure is written down: a second copy of each project's token now lives outside its repository, readable by the operator's user only; `docs/operators/index.md` says so and what reads it. An ADR records the registry, refining ADR-0018 and ADR-0024
- [ ] Tests with the fake runner and a scratch home: start, stop, rotate, restart, stale record, two projects, a project with no key in its manifest; `flai-cli.md` and `docs/users/flai.md` updated

## Tasks

## Notes
From E-0007, the operator's first bullet: "Extend flai and flaiover to write and manage a metadata file about each running instance into ~/.flai (home directory belongs to the current operator) in YAML format - this includes the port for each flaiover instance as well as the necessary token to access it."

What exists. `flai dashboard` runs one container per project, `flaiover-<key>`, and keeps the token at `.flai-cache/dashboard.token` inside the repository (0600, ignored by git, mounted read-only twice since S-0064). Nothing outside the repository knows a dashboard is running; `flai dashboard status` works only from inside a project. `~/.flai/config.json` is flai's host config for every project made from the template.

The token or a pointer to it. The operator asked for the token in the file. The alternative is the path of the project's token file: one copy of the secret instead of two, and rotation cannot leave the registry stale, at the cost of the reader needing read access to each repository (it has it: same user). Decide in the story and say why; either satisfies "the necessary token to access it".

A project's key comes from `system-flow.yaml` (`key`, a warning when missing since S-0043). Two clones of the same project on one host would collide on the file name; the record should be keyed so that does not silently overwrite (key plus a short hash of the path, or refuse with the reason).

In this repository flai's config is `.flai-cache/config.json` through `scripts/flai.sh`, not `~/.flai`, and agents here must not write under the operator's `~/.flai` while testing. The override in the sixth criterion exists for that.

S-0066 (port negotiation) builds on this record, and S-0067 (flairport) reads it.
- 2026-09-20T06:46:19Z: moved to cancelled: E-0007 cancelled: going a different direction
