---
id: T-3117
type: task
nature: remediation
title: component-library.md says where the dev server's port comes from and how flai knows the answer is its vite's, and I-0342 is closed
status: backlog
parent: S-0714
owner: alex
created: 2026-11-11T09:05:48Z
updated: 2026-11-11T09:05:48Z
transitions: []
stream: S-0714
tags: [docs, components, issues]
touches: [design/system/component-library.md, design/issues/I-0342-testeverylistedflaiovercomponentcompilesintheharness-fails-a-close-out-with-the.md, design/issues/summary.md]
after: [T-3116]
---
# T-3117 component-library.md says where the dev server's port comes from and how flai knows the answer is its vite's, and I-0342 is closed

## Work

`design/system/component-library.md` says flai finds the free port by listening on port 0 and closing it, and that the gap cannot be closed (S-0447, I-0162). Rewrite that bullet: the port now comes from below the host's ephemeral range, so no port 0 listener or outgoing connection on the host is given it, and the harness marks each answer with a token that the start's wait and the integration test check, naming S-0714 and I-0342. Bump the document's `updated`.

Then close the issue with `flai issue close I-0342 --reason "<what fixed it>"`, naming T-3114, T-3115, and T-3116, which updates `design/issues/summary.md`.

Waits for T-3116: the issue closes only once the fix and its test are in, and the document describes what they built.

## Done when

- `design/system/component-library.md` describes the port's range and the token, with no claim that the gap cannot be closed.
- I-0342 is closed with a reason naming the fix, and `design/issues/summary.md` no longer lists it open.
- `flai check --strict` and the markdown lint pass on the three files.

## Notes
