---
id: TH-0017
title: "main's system-flow check is red since S-0122's acceptance: may an agent fix three archived files' markdown lint?"
anchor:
  path: design/issues/I-0027-work-items-written-in-the-main-checkout-reach-ci-without-a-markdown-lint.md
status: resolved
participants: [agent-S-0123, alex, agent-S-0134]
created: 2026-09-26T07:23:33Z
updated: 2026-10-06T09:55:29Z
---

# TH-0017 main's system-flow check is red since S-0122's acceptance: may an agent fix three archived files' markdown lint?

On design/issues/I-0027-work-items-written-in-the-main-checkout-reach-ci-without-a-markdown-lint.md.

## Entries

### 2026-09-26T07:23:33Z agent-S-0123
`system-flow check` on main (run 36225691919, on the S-0122 push) fails at Lint markdown: a double blank line (MD012) in `wip/archive/agents/S-0122.md` line 37 and in the archived S-0122 story line 68, and emphasis used as a heading (MD036) in `wip/threads/TH-0012-…md` line 79. The archive is not an agent's to edit, so I left them.

Recommended: A. allow an agent a one-commit `chore` fix of the three on main (collapse the blank lines, make "Done" a plain line), and add a story so that `flai accept` lints what it archives and commits (I-0027, now count 2). B. You fix them by hand. Either way S-0123 does not depend on it.

### 2026-09-26T08:01:16Z alex
A

### 2026-09-27T04:08:35Z agent-S-0134
Done under your answer A, from S-0134, whose smoke tier failed on it: 88e345f on main fixes all six lint errors main carried. They were the three here, plus three in the archived S-0124 narrative (lines left over from an edit, removed). Main's markdown lint is clean now. I-0027 is bumped to count 4. The story that makes `flai accept` lint what it archives is still not written.

### 2026-09-29T05:23:48Z alex
yes

### 2026-10-06T09:55:29Z alex
Resolved.
