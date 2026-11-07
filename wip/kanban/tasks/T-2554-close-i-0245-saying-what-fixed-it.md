---
id: T-2554
type: task
nature: improvement
title: Close I-0245 saying what fixed it
status: in-progress
parent: S-0568
owner: alex
created: 2026-11-07T08:19:20Z
updated: 2026-11-07T08:37:00Z
transitions:
  - to: ready
    at: 2026-11-07T08:36:59Z
    by: system-flow
  - to: in-progress
    at: 2026-11-07T08:37:00Z
    by: system-flow
stream: S-0568
tags: [flai, issues]
touches: [design/issues/I-0245-flai-guard-refuses-a-sub-agent-every-flai-role-subcommand-the-reads-list-show.md, design/issues/summary.md]
after: [T-2551, T-2553]
---
# T-2554 Close I-0245 saying what fixed it

## Work

Layer 3: this task waits for T-2551 and T-2553, because the reason it writes names the change and its test, and the issue closes only once both are done.

The story's agent runs this itself, since `flai guard` refuses a sub-agent `flai issue close`:

```bash
flai issue close I-0245 --reason "<what fixed it>"
```

The reason names S-0568, the guard's `flagReads` entry for `role` or what T-2551 chose instead, the test that reproduces the refusal, and the commit it passes at.

## Done when

- `flai issue list` no longer lists I-0245 as open, and `design/issues/summary.md` agrees.
- The issue's reason names S-0568, the change, and the test.
- `flai check --strict` is clean for the story.

## Notes
