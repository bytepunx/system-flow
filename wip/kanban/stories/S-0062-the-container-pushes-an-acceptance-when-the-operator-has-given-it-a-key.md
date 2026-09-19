---
id: S-0062
type: story
nature: feature
title: The container pushes an acceptance when the operator has given it a key
status: backlog
parent: E-0006
owner: alex
created: 2026-09-19T07:54:42Z
updated: 2026-09-19T07:55:41Z
transitions: []
tags: []
touches: [flai/cmd, flaiover, design/adrs, docs/operators]
---

# S-0062 The container pushes an acceptance when the operator has given it a key

## Goal
An acceptance made from the board reaches the remote at once, with its release tags, when the operator has opted in by naming an SSH key for `flai dashboard` to mount. A project that has not opted in behaves exactly as today, and the container holds nothing.

## Acceptance criteria
- [ ] The opt-in is a host setting (`flai config set dashboard.push_key <path>` and a `--push-key` flag), never the manifest, and is off by default; with it unset `flai dashboard` mounts no key and sets nothing new, which a test pins
- [ ] With it set, `flai dashboard` mounts the key read-only, sets git's SSH command to use only that key, and pins the remote's host key from the host's own `known_hosts`; when the host has no entry for the remote it refuses with the command to add one, and the container never accepts a host key on first use
- [ ] The image carries an SSH client, and the container's user has a passwd entry whatever the host's user ID is; an acceptance from the board in a container run as a user ID the image does not know pushes successfully (OpenSSH otherwise refuses with "No user exists for uid")
- [ ] A key protected by a passphrase is refused at start with the reason (it cannot push unattended) and what to do; so are a missing file, a file that is not a private key, and permissions looser than 0600. The agent socket is not forwarded
- [ ] `flai dashboard` says at start what the container holds (the key's fingerprint and comment, never the key) and what it means: whoever holds the dashboard token can publish a release by accepting a story. `flai dashboard status` shows the same
- [ ] An acceptance from the board then ends with "pushed", tags included, and the template component is published when the plan has one; a push that fails leaves the acceptance standing with the existing not-pushed result
- [ ] Tried end to end against a scratch remote over SSH with a throwaway key, never this repository's remote or the operator's key: the default, the opt-in, the unknown user ID, the passphrase refusal, the missing host key
- [ ] `docs/operators/index.md` documents the opt-in with the safer choice first: a key dedicated to this repository (a deploy key with write access, with the commands to make one and where to add it), then the operator's own key with what that exposes in plain words, the repository rules to set (no force push, no deletion of branches and tags), and how to revoke. `flaiover-dashboard.md`, `flai-cli.md`, and the security notes that say the container holds no credential are updated

## Tasks

## Notes
From S-0052's finding, `design/system/pushing-from-the-board.md`, and ADR-0026. The operator decided on 2026-09-19 that every board acceptance should push, and chose their own SSH key, mounted, over the dedicated deploy key the finding recommended, having been told what it exposes; the dashboard will be reached over a public tunnel with the token as its only lock. The mechanism is the same for any key file, so this story builds "a key the operator names" and the documentation says which is safer.

What the research tried that this story can lean on: the published image has git and no SSH client; with `openssh-client-default` added (about 0.7 MB) a key mounted read-only pushed as the host user with `GIT_SSH_COMMAND="ssh -i <key> -o IdentitiesOnly=yes -o UserKnownHostsFile=<file> -o StrictHostKeyChecking=yes"`; OpenSSH refused a user ID with no passwd entry; it worked on the research host only because UID 1000 is `node` in the image. Mounting a generated passwd file read-only is the simplest fix found, not tried.

The acceptance flow already pushes whenever it can (`flai/cmd/accept.go`), so nothing changes there beyond the environment git runs in.

Do S-0064 first or alongside if it is ready: a container that holds a key is a better target, and the hook route it closes is the other way to the operator's credentials.
