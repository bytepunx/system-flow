---
id: S-0091
type: story
nature: improvement
title: The flai CLI should install itself to a user's home directory or a configurable path
status: done
parent: E-0003
owner: alex
created: 2026-09-21T04:09:25Z
updated: 2026-09-22T21:29:57Z
transitions:
  - to: ready
    at: 2026-09-21T04:13:50Z
    by: alex
  - to: in-progress
    at: 2026-09-22T21:01:55Z
    by: system-flow
  - to: review
    at: 2026-09-22T21:06:51Z
    by: system-flow
  - to: done
    at: 2026-09-22T21:29:57Z
    by: alex
tags: [cli]
touches: [flai/cmd]
---
# S-0091 The flai CLI should install itself to a user's home directory or a configurable path

## Goal

As of now, installing flai requires sudo because the path the installer chooses is almost always blocked for writes by any user (`/usr/local/bin`). This means updating flai via `self-update` requires an operator with direct host access to `sudo flai self-update`, keeping agents that would be otherwise capable from doing it as well as keeping flai from detecting an out-of-date situation and running it itself.

## Acceptance criteria
- [x] the installer should create/choose a path that the user has access to, such as ~/.flai to install the CLI
- [x] the installer no longer requires sudo
- [x] the installer reports where the installation occurred and gives a PATH line to prefix $PATH with the path to the flai CLI binary

## Tasks
- T-0321 The installer and self-upgrade default to a directory the user already owns, with no sudo path

## Notes
