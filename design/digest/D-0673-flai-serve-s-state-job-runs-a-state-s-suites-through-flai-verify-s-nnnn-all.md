---
id: D-0673
title: flai serve's state job runs a state's suites through flai verify S-nnnn --all --suite, not flai test, because flai test's tier selection skips the all_only tiers without --all
kind: decision
sources: [S-0810]
topics: [verify, workflow, flai, go, code]
paths: [flai/internal/serve/statejob.go, flai/internal/verify/select.go, flai/cmd/verify.go]
date: 2026-11-21
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0673 flai serve's state job runs a state's suites through flai verify S-nnnn --all --suite, not flai test, because flai test's tier selection skips the all_only tiers without --all

S-0810 made `flai serve` run a state's `runs:` suites as a job when a story enters the state (ADR-0312 §4). It had two commands to run a suite with. Its narrative's Decisions (`wip/archive/agents/S-0810.md`) record the choice:

> The state job runs a state's suites through flai verify S-nnnn limited to those suites (a repeatable --suite), with the tiers a story's paths select as the close-out selects them (selectStoryTiers, all_only tiers whose paths match included), not through flai test, whose Select skips all_only tiers without --all.

So a suite step selects tiers as the close-out does, and an `all_only` tier, such as integration or smoke here, runs when the story's paths match its `when`. Through `flai test` it would not run at all without `--all`, and with `--all` it would run whatever the paths. A later step that runs a list of suites for a story, such as the publish's or CI's (`verify_at`), should weigh the same difference before it reaches for `flai test`.

How the job runs each suite is in `design/system/flai-cli.md` (the state job) and `design/system/workflow.md` § What a state's `runs` starts (S-0810). A suite-limited run stores no report and no record (`docs/users/flai.md`, `--suite`).

Source: S-0810's narrative, § Decisions.
