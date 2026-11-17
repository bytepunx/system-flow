---
id: D-0658
title: The section merge's word merge stops on an insertion where the other side's change begins or ends, not only on two insertions at one point, which is how ADR-0314's 'meet at one word or one point' reads
kind: decision
sources: [S-0817]
topics: [merge, sync, acceptance, flai, go]
paths: [flai/internal/storygit/mdmerge.go, flai/internal/storygit/mdmerge_test.go, design/adrs/0314-a-section-merge-merges-a-line-both-sides-change-by-its-words-keeping-both.md]
date: 2026-11-17
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0658 The section merge's word merge stops on an insertion where the other side's change begins or ends, not only on two insertions at one point, which is how ADR-0314's 'meet at one word or one point' reads

ADR-0314 merges by their words the lines both sides of a section merge change, and leaves the stop when both change the same word or both insert at the same point; its consequences say "Edits that meet at one word or one point still stop".

S-0817 read that sentence wider than two insertions at one point: "The word merge's strict points also stop on an insertion where the other side's change begins or ends, not only on two insertions at one point. This reads ADR-0314's 'edits that meet at one word or one point still stop'. Without it, one side replacing the word mat with 'old rug' while the other inserts 'old' before mat gave 'old old rug'" (wip/archive/agents/S-0817.md, Decisions). The code is `merge3` with `strictPoints` and `meetAtPoint` in `flai/internal/storygit/mdmerge.go`.

The same story made `lineHunks` place each insertion as late as the base allows, since without it "a word both sides inserted alike landed twice"; line-level results were unchanged. A later change to the word merge keeps both, or a test of `mdmerge_test.go` that pins them fails.

Sources: S-0817.
