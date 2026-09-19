---
id: T-0168
type: task
nature: research
title: "Research: what pushing from the container would need, and what each way exposes"
status: done
parent: S-0052
owner: alex
created: 2026-09-19T01:55:25Z
updated: 2026-09-19T07:36:30Z
transitions:
  - to: ready
    at: 2026-09-19T07:29:03Z
    by: system-flow
  - to: in-progress
    at: 2026-09-19T07:29:03Z
    by: system-flow
  - to: done
    at: 2026-09-19T07:36:30Z
    by: system-flow
stream: S-0052
tags: []
---

# T-0168 Research: what pushing from the container would need, and what each way exposes

## Work
For each way the container could push (forwarded SSH agent socket, a deploy key mounted read-only, an HTTPS token given to a git credential helper, and any other found on the way), work out: what `flai dashboard` would have to mount or set; whether it works with the container running as the host user on Linux and WSL2, and what differs on macOS and Windows with Docker Desktop; what someone holding the dashboard token could do to the remote with it, including force pushes and tags that start release workflows; how narrowly the credential can be scoped (one repository, contents only, no workflow scope); how it is rotated and revoked. Try each that looks viable in a scratch repository pushing to a scratch bare remote on this machine, and for the hosted-remote specifics read the provider's documentation and cite it. Never use this repository's remote or the operator's real credentials.

## Done when
- A section of the finding covers each mechanism under the same headings
- Each claim is marked tried (with what was run) or not tried

## Notes
