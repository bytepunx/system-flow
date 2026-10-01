---
id: E-0015
type: epic
nature: feature
title: Secure Dashboard and CLI Releases
status: backlog
owner: alex
created: 2026-10-01T11:06:14Z
updated: 2026-10-01T11:12:56Z
transitions: []
tags: [dashboard, cli]
topics: [releases]
---
# E-0015 Secure Dashboard and CLI Releases

## Outcome

The dashboard and CLI are both cryptographically signed in CI with a private key and the CLI verifies new releases during self-upgrade with the public key.

Determine how each component can use this same mechanism to ensure they are only talking to a signed version of the other such that flai terminates the connection to an unsigned flaiover and flaiover terminates the connection from an unsigned flai.

## Stories
- S-0192 Add a way to create a sibling story

## Notes
