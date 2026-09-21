---
id: S-0091
type: story
nature: improvement
title: The flai CLI should install itself to a user's home directory or a configurable path
status: backlog
parent: E-0003
owner: alex
created: 2026-09-21T04:09:25Z
updated: 2026-09-21T04:09:25Z
transitions: []
tags: [cli]
touches: [flai/cmd]
---
# S-0091 The flai CLI should install itself to a user's home directory or a configurable path

## Goal

As of now, installing flai requires sudo because the path the installer chooses is almost always blocked for writes by any user (`/usr/local/bin`). This means updating flai via `self-update` requires an operator with direct host access to `sudo flai self-update`, keeping agents that would be otherwise capable from doing it as well as keeping flai from detecting an out-of-date situation and running it itself.

## Acceptance criteria
- [ ] the installer should create/choose a path that the user has access to, such as ~/.flai to install the CLI
- [ ] the installer no longer requires sudo
- [ ] the installer reports where the installation occurred and gives a PATH line to prefix $PATH with the path to the flai CLI binary

## Tasks

## Notes
