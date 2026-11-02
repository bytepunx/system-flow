---
id: D-0464
title: "flai issue close refuses an issue already closed, so a story that fixes more of a closed issue records its fix as body text in the issue's ## Remediation, after the close line"
kind: decision
sources: [S-0540, T-2329, I-0227]
topics: [flai, cli, issues]
paths: [design/issues]
date: 2026-11-02
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0464 flai issue close refuses an issue already closed, so a story that fixes more of a closed issue records its fix as body text in the issue's ## Remediation, after the close line

S-0540 was made from I-0227 to remediate it, but S-0548 closed I-0227 first, on 2026-10-31, with part of the fix. S-0540's planner had written T-2329 to close the issue with `flai issue close`. The story's agent found that command cannot close an issue that is already closed, and rewrote the task: "the fix goes into I-0227's `## Remediation`, as body text, after S-0548's close line" (`wip/archive/agents/S-0540.md`, § Decisions).

I-0227's `## Remediation` now ends with a dated line, "2026-11-02: S-0540 fixed the second cause, …", naming the task and the file it changed, below the `Closed …` line `flai issue close` wrote for S-0548. So a later story whose issue was closed before it finished records what it fixed the same way, in the issue's body, rather than reopening or closing the issue again.

Sources: S-0540's narrative, § Decisions; T-2329; I-0227, § Remediation.
