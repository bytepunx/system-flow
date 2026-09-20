---
id: S-0075
type: story
nature: feature
title: Every write the dashboard makes goes through the channel, and the image no longer carries flai, git, or ssh
status: ready
parent: E-0003
owner: alex
created: 2026-09-20T07:26:50Z
updated: 2026-09-20T07:27:58Z
transitions:
  - to: ready
    at: 2026-09-20T07:27:58Z
    by: alex
tags: [cli, dashboard]
touches: [flai/cmd, flai/internal, flaiover/src, flaiover/Dockerfile, design/system]
---
# S-0075 Every write the dashboard makes goes through the channel, and the image no longer carries flai, git, or ssh

## Goal
The twenty flai commands the dashboard runs inside the container become named methods served by flai on the host. Acceptance streams its progress over the channel. With nothing left to run in the container, the image drops `flai`, `git`, and `ssh`.

## Acceptance criteria
- [ ] One method per command the dashboard uses today (move with its cancellation preview, order, block, unblock, epic and story new and their template bodies, accept and its dry run, stream diff, stream log, thread new, reply, and resolve, doc show and save with the hash, adr new, template, and accept, stats, check for overlaps, the pending push dry run), each with typed arguments validated in flai; the argument builders move from TypeScript to Go with their tests
- [ ] No method takes a command line, a flag list, or a path outside the manifest's folders; a test enumerates the methods and fails when one is added without a schema
- [ ] Acceptance streams each step as a notification and the board shows them as they happen, as it does now; a connection lost during an acceptance ends the browser's request with an error that says the acceptance may have completed and how to find out
- [ ] Every write carries a request ID; flai answers a repeat of a completed ID with the recorded result, so a retry cannot move, accept, or save twice; tested by repeating each kind
- [ ] Exit codes that carry meaning today (3 conflict, 4 refused, rule violations as 400) arrive as typed errors and the pages behave as before
- [ ] `flai.ts` no longer spawns anything; the Dockerfile drops the Go build stage, `git`, and `openssh`; `/_ready` no longer looks for a binary; the image size before and after is recorded
- [ ] The design documents and the users' and operators' documentation are updated

## Tasks

## Notes
From S-0071's finding, `design/system/dashboard-host-channel.md`, and ADR-0029. Depends on both read stories. The clone is still mounted after this story, unused; the next story but one removes it.
