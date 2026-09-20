---
id: ADR-0026
title: The dashboard container may push acceptances with a key the operator gives it
status: accepted
date: 2026-09-19
supersedes: []
superseded_by: [ADR-0031]
refines: [ADR-0018]
---

# ADR-0026 The dashboard container may push acceptances with a key the operator gives it

## Context

[ADR-0018](0018-dashboard-token.md) guards the dashboard with one project token, and `flai dashboard` has so far given the container the host's git identity and global excludes and never a credential. An acceptance from the board is therefore committed and tagged in the clone and pushed by someone on the host. The operator works from the board, often remotely and without a shell, was surprised twice that an accepted release had not reached the remote, and wants to drive the project from the dashboard.

S-0052 researched the ways the push could happen, inside the container and outside it, and tried them against a scratch remote: [Pushing an acceptance made from the board](../system/pushing-from-the-board.md). The operator's answers on 2026-09-19: the dashboard will be reached by them alone over a public tunnel, so the token is its only lock; a push that nobody on the host approved is acceptable, for every acceptance; the board should show a standing "accepted, not pushed" state. For the credential the finding recommended a key dedicated to this repository. The operator, told what each choice exposes, chose their own SSH key, mounted.

## Decision

`flai dashboard` may mount an SSH private key that the operator names into the container, read-only, and the bundled `flai accept` then pushes the acceptance commit and its release tags with it. It is an explicit opt-in on the host, off by default; a project that does not opt in behaves as before, and the container holds nothing.

The key is whichever file the operator names. The operator of this repository names their own key. The mechanism does not distinguish, so a key dedicated to one repository (a deploy key with write access) is a matter of naming a different file, and the operators' documentation presents that as the safer choice, first.

With the opt-in on:

- The image carries an SSH client. The container's user gets a passwd entry when the image has none for the host's user ID, because OpenSSH refuses to run without one.
- The remote's host key is pinned from the host's own `known_hosts`; the container never accepts a host key on first use.
- A key protected by a passphrase cannot push unattended. `flai dashboard` refuses it at start and says why, rather than failing at the first acceptance. The agent socket is not forwarded as a way round this: it would lend every key in the agent and fail whenever the operator is logged out.
- `flai dashboard` says at start what the container holds (the key's fingerprint and comment) and what that means: whoever holds the dashboard token can publish a release by accepting a story.
- A push that fails leaves the acceptance standing, as now.

Whether or not a key is given, the board and the item page show a standing "accepted, not pushed" state while `main` is ahead of its remote-tracking branch by an acceptance, and agents see the same fact in `inbox`.

## Consequences

- The security posture of ADR-0018 changes for a project that opts in: the dashboard token becomes the power to publish a release of any story in review, and a compromise of the container yields the key. With the operator's own key that is every repository of their account and every server that trusts the key; with a dedicated key it is one repository. The operators' documentation says so in those words.
- A stolen key can publish a release from a commit of the thief's making, whichever key it is. Repository rules that block force pushes and the deletion of branches and tags narrow what else it can do, and are recommended with the opt-in.
- The opt-in is a host setting, not part of the manifest: a path to a key is a fact about one machine and does not belong in the repository.
- The research also found that a container with the clone mounted read-write can already leave a git hook that runs on the host as the operator. That is independent of this decision and is remediated separately (S-0064); until then it is the more likely route to the operator's credentials for a project that has not opted in.
- An HTTPS fine-grained token, the best scoped credential the research found, is not built. It needs the push redirected for a repository whose remote is SSH. It is the natural next opt-in if a project with an HTTPS remote asks.
