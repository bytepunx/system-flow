---
id: S-0084
type: story
nature: remediation
title: "One flai serve per repository: a second one steps back instead of taking the dashboard's connection"
status: backlog
parent: E-0003
owner: alex
created: 2026-09-20T13:56:29Z
updated: 2026-09-20T13:56:29Z
transitions: []
tags: [cli, dashboard]
touches: [flai/cmd, flai/internal, flaiover/src, design/system, docs/operators]
---
# S-0084 One flai serve per repository: a second one steps back instead of taking the dashboard's connection

## Goal
A repository is served by one `flai serve` at a time. A second one, under another configuration or from another flai, does not take the dashboard's connection away from the first; it is told who serves the project and steps back.

## Acceptance criteria
- [ ] One `flai serve` per repository: before serving a project, `flai serve` claims it with a lock kept with the repository (under the main checkout's `.flai-cache`), naming its PID, flai version, configuration, and start; a claim whose process is gone is taken over, a live one is respected
- [ ] A `flai serve` that finds a project claimed by another does not dial that project's dashboard; `flai serve status`, `flai dashboard`, and `flai dashboard status` say which process serves it, from which flai and configuration, and what to run to hand it over
- [ ] The dashboard no longer lets a second flai with the same credential silently replace a healthy one: what it does instead (refuse the newcomer while the holder answers, replace only a holder that has stopped answering) is decided, tested from both sides, and recorded; a flai that was replaced or refused backs off and says so instead of redialling every few seconds
- [ ] Host actions and the journal follow the claim: the serve that holds the repository is the one whose configuration decides what is enabled, and `flai serve actions` says so when the project is served by another
- [ ] Tried with two flai builds and two configurations against a scratch project and dashboard: the second start is refused with the message, no connection flaps, a killed holder is taken over, and a write in flight during the attempt is neither lost nor done twice
- [ ] What "eventually one per machine" would need (one registry for every configuration, where it lives, how an installed flai and a tree's flai share it) is written down as a follow-up, not built here

## Tasks

## Notes
Asked for by the operator on 2026-09-20 after I-0029: the installed flai 1.6.1 under `~/.flai` and the tree's flai under `.flai-cache` both served this repository with the same agent credential, and replaced each other on the dashboard every few seconds (457 connections in 45 minutes). Whichever held the connection answered with its own version, methods, and host actions, so the dashboard said push was disabled and `push.run` missing although the tree's serve had both. "One flai serve per user" (ADR-0029) is only kept per configuration folder today. The operator's direction: one per repository now, one per machine eventually.
