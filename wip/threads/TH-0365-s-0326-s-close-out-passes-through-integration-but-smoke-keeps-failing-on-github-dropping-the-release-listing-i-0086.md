---
id: TH-0365
title: S-0326's close-out passes through integration but smoke keeps failing on GitHub dropping the release listing (I-0086)
anchor:
  path: wip/kanban/stories/S-0326-flai-verify-record-issues-opens-an-issue-on-the-story-branch-that-another-branch-opened-under-the-same-title-meanwhile.md
  item: S-0326
status: resolved
participants: [agent-S-0326, alex]
created: 2026-10-08T07:39:02Z
updated: 2026-10-08T07:58:04Z
---

# TH-0365 S-0326's close-out passes through integration but smoke keeps failing on GitHub dropping the release listing (I-0086)

On wip/kanban/stories/S-0326-flai-verify-record-issues-opens-an-issue-on-the-story-branch-that-another-branch-opened-under-the-same-title-meanwhile.md.

## Entries

### 2026-10-08T07:39:02Z agent-S-0326
Recommendation: let me move S-0326 to review once one more close-out has passed every tier but smoke, and run `make smoke` yourself before acceptance, or tell me to keep retrying.

What happened:

- Four close-outs passed rebase, sync, narrative, check, gofmt, vet, golangci-lint, go-test, markdown, and integration. Integration runs the git-backed tests that reproduce both of I-0112's instances.
- Each time, smoke then failed in `scripts/install-test.sh`. `install.sh` fetches `releases?per_page=50` (about 1 MB) and the connection drops: `curl: (92) HTTP/2 stream 1 was not closed cleanly` or `curl: (56) ... unexpected eof while reading`.
- Measured on this host: 2 of 5 fetches of that listing dropped; `per_page=5` came back 3 of 3. That is a network fault, not this story's code. It is I-0086, which S-0291 remediates. I added this evidence to my instance there.

The conventions let me move to review only on a clean close-out, so I am asking. Alternatives:

- Keep retrying the close-out until smoke passes. Each run takes about 5 minutes.
- Fix `install.sh` here: fetch fewer releases and retry once on a network error. That is S-0291's scope, not S-0326's.

Every criterion is ticked, and the branch is synced and committed.

### 2026-10-08T07:58:04Z alex
Resolved.
