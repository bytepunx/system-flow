---
id: T-1356
type: task
nature: remediation
title: flai.md, ci.md, devex.md, flai-cli.md, and scripts/README.md say what the smoke tier installs from and what the GitHub check covers
status: backlog
parent: S-0340
owner: alex
created: 2026-10-08T08:06:17Z
updated: 2026-10-08T08:06:17Z
transitions: []
stream: S-0340
tags: [cli]
touches: [docs/users/flai.md, design/tech/ci.md, design/system/devex.md, design/system/flai-cli.md, scripts/README.md]
after: [T-1354, T-1352]
---
# T-1356 flai.md, ci.md, devex.md, flai-cli.md, and scripts/README.md say what the smoke tier installs from and what the GitHub check covers

## Work

Say in each document what the smoke tier installs from, what the GitHub check covers, and how to run each by hand.

- The smoke tier installs and self-upgrades from a release built from the tree and served on `127.0.0.1` by `scripts/release-server.sh`. It needs no token and no network to GitHub. Run it by hand with `make install-test`.
- `scripts/install-published-test.sh` checks that the latest published release installs from GitHub. `install-published.yml` runs it on push to `main`, daily, and by hand, and a failure is recorded through `flai issue`, as T-1352 settles. Run it by hand with `make install-published-test`.

By document:

- `docs/users/flai.md`: beside the `self-upgrade` paragraph, for a user testing an installer against a mirror or a local build.
- `design/tech/ci.md`: the `system-flow-check.yml` row no longer installs the latest release; add a row for `install-published.yml`.
- `design/system/devex.md`: the Test tiers row, on what smoke installs from.
- `design/system/flai-cli.md`: the `flai self-upgrade` row says `scripts/install-test.sh` runs "against the real latest release in the smoke tier"; correct it.
- `scripts/README.md`: rows for `release-server.sh` and `install-published-test.sh`, and the new mode of `flai-snapshot.sh`.

Bump each design and docs file's `updated`.

Waits for T-1354 and T-1352: it describes what they built, under the names they chose.

## Done when

- Each of the five files says what the smoke tier installs from, what the GitHub check covers, and how to run each, where its subject reaches them.
- `flai test` on the five paths passes the markdown lint.

## Notes

Drafted by the planner from S-0340's criterion 5. `design/system/flai-cli.md` and `scripts/README.md` are added to the criterion's three: the first's `self-upgrade` row describes the smoke test, and the second indexes the scripts.
