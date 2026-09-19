---
id: T-0221
type: task
nature: feature
title: "The opt-in and its checks: dashboard.push_key and --push-key, the key refused when it cannot push unattended, host keys pinned from the host, a passwd entry, and what flai dashboard says at start and in status"
status: done
parent: S-0062
owner: alex
created: 2026-09-19T08:21:01Z
updated: 2026-09-19T08:29:06Z
transitions:
  - to: ready
    at: 2026-09-19T08:21:02Z
    by: system-flow
  - to: in-progress
    at: 2026-09-19T08:21:02Z
    by: system-flow
  - to: done
    at: 2026-09-19T08:29:06Z
    by: system-flow
stream: S-0062
tags: []
---

# T-0221 The opt-in and its checks: dashboard.push_key and --push-key, the key refused when it cannot push unattended, host keys pinned from the host, a passwd entry, and what flai dashboard says at start and in status

## Work
Host config key `dashboard.push_key` (never the manifest) and `flai dashboard --push-key <path>`, off by default. With it unset nothing new is mounted or set, pinned by a test. With it set, before the container starts: the file exists, is a regular file, is no looser than 0600, is a private key, and has no passphrase (`ssh-keygen -y -P ""`), each refusal with its reason and what to do; the remote is an SSH URL; the remote's host keys are copied from the host's `known_hosts` (`ssh-keygen -F`, which handles hashed entries) into `.flai-cache/dashboard.known_hosts`, and a missing entry refuses with how to add one; a passwd file with an entry for the host's user ID is generated. The key, the host keys, and the passwd file are mounted read-only, and `GIT_SSH_COMMAND` names only that key and that file, with strict host key checking, batch mode, and no user configuration. The agent socket is never forwarded. `flai dashboard` prints what the container holds (fingerprint and comment) and what it means; `--json` and `flai dashboard status` carry the same. Windows hosts are refused for now, with the reason. Tests with the fake runner for every branch, and one with the real `ssh-keygen` for the passphrase and not-a-key cases.

## Done when
- Every refusal and the happy path are tested, and the default is pinned
- `make flai-test` passes

## Notes
