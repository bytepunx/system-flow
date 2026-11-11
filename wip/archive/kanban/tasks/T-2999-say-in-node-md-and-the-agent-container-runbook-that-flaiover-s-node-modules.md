---
id: T-2999
type: task
nature: remediation
title: Say in node.md and the agent-container runbook that flaiover's node_modules holds both libc variants, and close I-0259
status: done
parent: S-0706
owner: alex
created: 2026-11-10T10:31:01Z
updated: 2026-11-11T09:24:32Z
transitions:
  - to: ready
    at: 2026-11-11T09:23:56Z
    by: agent-S-0706
  - to: in-progress
    at: 2026-11-11T09:23:56Z
    by: agent-S-0706
  - to: done
    at: 2026-11-11T09:24:32Z
    by: agent-S-0706
stream: S-0706
tags: [docs, issues]
touches: [design/tech/node.md, docs/operators/runbooks/agent-container.md, design/issues/I-0259-a-story-worktree-s-flaiover-node-modules-installed-inside-the-alpine-agent.md, design/issues/summary.md]
after: [T-2997]
usage:
  source: log
  seconds: 36
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 17
      output: 4555
      cache_read: 1837384
      cache_write: 41940
      cost: 0.7942
---
# T-2999 Say in node.md and the agent-container runbook that flaiover's node_modules holds both libc variants, and close I-0259

## Work

It waits for T-2997, whose test the issue's closing reason names, and through it for T-2993's setting and T-2994's rule, which the documents describe.

Add a short note to `design/tech/node.md` beside the store and frozen-lockfile notes. It says that flaiover's install puts in both the glibc and the musl native bindings (`supportedArchitectures`), so that one `node_modules` serves the host and the Alpine agent container. It also says that `--if-needed` installs again when the install configuration changes. In `docs/operators/runbooks/agent-container.md`, near the mounts, say that the worktree's `flaiover/node_modules` is shared with the host. Say that it now holds both variants. Say that an install made before this story is replaced on the next `--if-needed` install, or by `make flaiover-install`.

Then run `flai issue close I-0259 --reason "<reason>"` in the story's worktree. The reason names T-2993's setting, T-2994's rule, and the test T-2997 added. `flai issue close` writes the issue file and regenerates `design/issues/summary.md`.

## Done when

- `design/tech/node.md` and `docs/operators/runbooks/agent-container.md` say what is above, with `updated` bumped.
- I-0259 is closed, with a reason that names the fix and its test.
- `design/issues/summary.md` no longer lists I-0259 as open.
- `flai test` on the changed documents passes, and `flai check --strict` reports nothing on them.

## Notes
