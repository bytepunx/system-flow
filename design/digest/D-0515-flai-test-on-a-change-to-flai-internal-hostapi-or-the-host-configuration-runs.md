---
id: D-0515
title: flai test on a change to flai/internal/hostapi or the host configuration runs neither flaiover's server tests that drive the real flai nor TestSettingsIndexIsComplete, so such a story runs them itself before its close-out (I-0304)
kind: cause
sources: [S-0587, I-0304, S-0592]
topics: [flai, flaiover, testing, docs]
paths: [flai/internal/hostapi/writes.go, flaiover/src/lib/server/writes.test.ts, flai/cmd/settings_doc_test.go, docs/operators/settings.md, flai/internal/config, system-flow.yaml, flai/internal/verify/select_test.go]
date: 2026-11-07
status: absorbed
absorbed_by: "design/system/devex.md#used-by-this-monorepo"
asked_again: []
cited: [S-0634, TH-0947]
---

# D-0515 flai test on a change to flai/internal/hostapi or the host configuration runs neither flaiover's server tests that drive the real flai nor TestSettingsIndexIsComplete, so such a story runs them itself before its close-out (I-0304)

S-0587's first two close-outs each failed on a test that `flai test` had not run for the story's paths. Both are recorded in I-0304.

- First close-out: `flaiover/src/lib/server/writes.test.ts` failed. It runs the built flai and still expected a 403 from the `settings.agent` gate that S-0587 made per project in `flai/internal/hostapi/writes.go`. vitest's `related` run picks tests by the TypeScript import graph, so no change under `flai/internal/hostapi` selects it.
- Second close-out: `TestSettingsIndexIsComplete` (`flai/cmd/settings_doc_test.go`) failed, because `docs/operators/settings.md` named `agent.projects.<root>` instead of each key below it. That test skipped under `-short`, so `flai test` on `docs/operators/settings.md` and `flai/internal/config` ran it nowhere. A task sub-agent had reported that it passed.

Fixed by S-0592, which closed I-0304. Its fix, from I-0304: "a change under flai/internal/hostapi/ or to flai/cmd/hostapi.go now selects the flaiover-flai tier, which runs vitest related on src/lib/server/testing.ts, so writes.test.ts and the other tests that drive the built flai run in flai test. TestSettingsIndexHasEveryKey checks the configuration and manifest key rows under -short, so go-test catches a schema change, and the settings-index tier runs it when docs/operators/settings.md changes." Why the new tier's paths stay narrow is D-0521.

Sources: S-0587 (wip/archive/agents/S-0587.md, Decisions); I-0304, Instances and Fix; S-0592.
