---
id: D-0651
title: "A criterion on what a flai serve job does, that names no live run, is verified by the story's Go and vitest tests, since flai serve here runs the installed flai, which predates the story; the live check waits for a release"
kind: decision
sources: [S-0804]
topics: [verify, acceptance, flai, flaiover, dashboard]
paths: [flai/internal/serve/resolvejob.go, flai/internal/serve/verifyjob.go, flai/internal/closeout, flaiover/src/lib/components/Review.svelte]
date: 2026-11-17
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0651 A criterion on what a flai serve job does, that names no live run, is verified by the story's Go and vitest tests, since flai serve here runs the installed flai, which predates the story; the live check waits for a release

S-0804 made Rebase and resolve a flai serve job and gave the close-out job progress lines that the review page follows. Its criteria 6 and 7 say what the job's record holds while it runs and what the review page shows. Neither names a live run. Its narrative, `wip/archive/agents/S-0804.md`, records under Decisions:

> Criteria 6 and 7 are verified by the Go and vitest tests, not on a live flai serve. flai serve here runs the installed flai, which predates S-0804, so this story's own close-out log holds no progress lines; a live check waits for a release.

The story ticked both criteria on those tests. Its Next steps left the live check for after a release: once flai serve runs a flai that carries S-0804, a real Rebase and resolve's log should hold the `verify:` and `resolve:` lines. The orchestrator accepted S-0804 on that evidence.

This differs from the answers in D-0461, D-0465, D-0466, and D-0469. There a criterion itself asked for a live run that could not happen before acceptance or a release, and the operator had it moved to a follow-up story. A criterion that names no live run stays in its story and is verified by tests.
