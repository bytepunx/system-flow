---
id: T-1462
type: task
nature: improvement
title: docs/users/flai.md says which sub-agents flai guard refuses a shell write, what counts as one, and what passes
status: done
parent: S-0368
owner: alex
created: 2026-10-12T16:57:39Z
updated: 2026-10-12T17:35:10Z
transitions:
  - to: ready
    at: 2026-10-12T17:21:56Z
    by: agent-S-0368
  - to: in-progress
    at: 2026-10-12T17:21:57Z
    by: agent-S-0368
  - to: done
    at: 2026-10-12T17:35:10Z
    by: agent-S-0368
stream: S-0368
tags: [guard, docs]
touches: [docs/users/flai.md]
after: [T-1460]
usage:
  source: log
  seconds: 793
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 16
      output: 7907
      cache_read: 858376
      cache_write: 38291
      cost: 0.578
---
# T-1462 docs/users/flai.md says which sub-agents flai guard refuses a shell write, what counts as one, and what passes

## Work

Update `## Sub-agents` in `docs/users/flai.md` from T-1460's ADR. It waits for T-1460 because the ADR decides what the page describes. It runs beside T-1461, since the two change no path in common.

- Add a row or sentence to the section's table saying a read-only sub-agent, such as the explorer or the verifier, is refused a shell command that writes a file in the project.
- List the commands counted as writes and what passes: `/dev/null`, the temporary folder, and descriptor duplications.
- Keep the paragraph that says the guard is not a shell. Add that a script or interpreter that writes is not seen, and that `git status` shows its file before review.
- Link the ADR.

## Done when

- [ ] `docs/users/flai.md` `## Sub-agents` states the rule as the ADR decides it, with its limits.
- [ ] The markdown lint passes on the file.

## Notes

S-0233 and S-0350, both in progress, also touch `docs/users/flai.md`, in other sections. Expect `flai stream sync` to bring their changes in.
