---
id: T-1341
type: task
nature: remediation
title: flai-cli.md and flai.md say install.sh and self-upgrade retry a call the network drops
status: backlog
parent: S-0291
owner: alex
created: 2026-10-08T08:00:12Z
updated: 2026-10-08T08:00:12Z
transitions: []
stream: S-0291
tags: [docs, install]
touches: [design/system/flai-cli.md, docs/users/flai.md]
after: [T-1338, T-1339]
---
# T-1341 flai-cli.md and flai.md say install.sh and self-upgrade retry a call the network drops

## Work

- In `design/system/flai-cli.md`, in the `flai self-upgrade` row of the commands table, say that `install.sh` and `self-upgrade` retry a GitHub API call that fails on the network, up to the attempts T-1338 and T-1339 settled, and never an HTTP error. Say that `install.sh` resolves the latest release from a small page. Name I-0086.
- In `docs/users/flai.md`, under the `self-upgrade` paragraph, say in one sentence that both retry a dropped connection to GitHub and that a 401 or 404 still fails at once.
- Bump each file's `updated`.

Waits for T-1338 and T-1339: it describes the attempts and the page size they settle.

## Done when

- Both documents say what T-1338 and T-1339 built, with the same number of attempts.
- `flai test design/system/flai-cli.md docs/users/flai.md` passes.

## Notes
