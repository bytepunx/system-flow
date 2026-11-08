---
id: T-2661
type: task
nature: feature
title: Mark this repository's independent tiers, time a close-out serial and at once on one commit, and record it in server-performance.md
status: backlog
parent: S-0610
owner: alex
created: 2026-11-07T23:49:41Z
updated: 2026-11-07T23:50:02Z
transitions: []
stream: S-0610
tags: [cli]
touches: [system-flow.yaml, design/system/server-performance.md]
after: [T-2653]
---
# T-2661 Mark this repository's independent tiers, time a close-out serial and at once on one commit, and record it in server-performance.md

## Work

For criterion 3. It waits for T-2653, the last code task, since it measures the run that T-2644 and T-2653 build and marks the tiers with T-2634's key.

- In `system-flow.yaml`'s `tests`, mark `integration` (Go) and `flaiover` (node) as able to run beside each other, and `smoke` as waiting for the tiers whose verdicts it reads, with T-2634's key. Leave the cheap per-path tiers as they are unless the ADR says otherwise. Check first that the flai `flai serve` runs on this host reads the manifest with the key; if it would refuse it, say so on a thread on the story before you commit it, since the close-outs of every story read this file.
- S-0610's own branch changes nothing under `flaiover/`, so its close-out does not select the `flaiover` tier. Time a run that selects `integration`, `flaiover`, and `smoke` on one commit: a close-out's tiers over a range that changes `flai/` and `flaiover/` files, such as `verify.Tiers` with a base before an accepted story that changed both (S-0595), or `flai test --all` if the run's selection is the same. Run it once with the key taken out (serial, as before) and once with it, both with the `go` and `vitest` caches warm alike. Run each twice if the host's load moves the figures by more than a third, and give the load average beside them.
- Record in `design/system/server-performance.md`, under cause 20 of § A second pass at 596 stories (S-0596), the before and after: what was run and on which commit, each tier's duration, the whole run's, the steps smoke skipped, and the host's load. Bump the document's `updated`.

## Done when

- `system-flow.yaml` marks the tiers, and `flai check --strict` passes.
- The two runs on one commit show the run with the key at least 150 s faster than the serial one, or the shortfall is stated with its cause and the criterion left unticked.
- `design/system/server-performance.md` records both runs.

## Notes
