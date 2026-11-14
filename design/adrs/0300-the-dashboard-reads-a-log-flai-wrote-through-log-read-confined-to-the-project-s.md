---
id: ADR-0300
title: "The dashboard reads a log flai wrote through log.read, confined to the project's .flai-cache and to its own logs in flai serve's state folder"
status: proposed
date: 2026-11-14
supersedes: []
superseded_by: []
refines: [ADR-0029]
topics: [flai, cli, dashboard]
---

# ADR-0300 The dashboard reads a log flai wrote through log.read, confined to the project's .flai-cache and to its own logs in flai serve's state folder

## Context

Statuses and errors on the dashboard name log files on the host: an agent run's under flai serve's state folder, an image job's, a checks run's, a kluster step's, a close-out's, a merge queue tier's. Reading one needs a shell on the host. The dashboard reaches the project only through flai (ADR-0029, ADR-0031), so a log's content must come through a host API method. flai's state folder also holds config.json with the operator's tokens, and flai serve's holds every served project's agent runs, so a method that reads any path a message names would hand them to whoever holds the dashboard token.

## Decision

The host API method `log.read`, a read under no host action, takes `{path, offset}` and answers `{path, exists, offset, next, size, text, reset}`, at most 64 KB from the offset ending on a whole character, as `images.log` and `kluster.log` answer. The path is absolute, `~/` for the home of the user flai runs as, or relative to the project's root. It reads a regular file whose name ends in `.log`, by its real path with every link followed, and only when that path lies:

1. under the project's main checkout's `.flai-cache`, but not in `.flai-cache/worktrees`, nor in a `serve` or `host` folder there, which is the state folder of a flai whose config file is in `.flai-cache`; or
2. in flai serve's state folder, `Host.ServeDir`, which flai serve sets to its own: `serve.log`; an agent run under `agents` named `<key>-<subject>-<stamp>.log` with the project's key and, as subject, a story's ID, a built-in strategic agent's kind, or a role the project defines; an image job's log under `image-jobs` that `Host.Images` gives as the project's; a checks run under `checks`, `<story>-<stamp>.log`, whose story's checks record names the project's root; or `host.log` in the `host` folder beside it.

Every other path is refused as bad params, naming the path and saying what is readable. A file that changes to a link or another file between the check and the read is refused. A log not there within these folders is answered as none, so a viewer can follow one not written yet.

## Consequences

A message's log path can be a link to a dashboard page that reads it. flai serve's folder layout is now read by hostapi too, by its names (`agents`, `image-jobs`, `checks`, `serve.log`, `host/host.log`), so a change to them changes `log.read`. A path that names a folder, such as a kluster job's step logs folder, is not a log `log.read` reads; its files are. `serve.log` and `host.log` hold every served project's lines and are readable from each project's dashboard. A story's checks runs older than its newest are readable only while its newest is the project's. `log.read` is not a guard read, as `images.log` is not, so a story's agent does not call it.

## Alternatives considered

- Any path a message names: hands the token holder flai's config and other projects' logs.
- A method per kind of log, as images.log and kluster.log: every new producer needs a method and a page of its own.
- flai serve answering whether a log in its folder is a project's, through a Host function: keeps its layout private, but adds a function per kind for a check of names hostapi can make itself.
