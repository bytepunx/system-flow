---
id: D-0515
title: flai test on a change to flai/internal/hostapi or the host configuration runs neither flaiover's server tests that drive the real flai nor TestSettingsIndexIsComplete, so such a story runs them itself before its close-out (I-0304)
kind: cause
sources: [S-0587, I-0304]
topics: [flai, flaiover, testing, docs]
paths: [flai/internal/hostapi/writes.go, flaiover/src/lib/server/writes.test.ts, flai/cmd/settings_doc_test.go, docs/operators/settings.md, flai/internal/config]
date: 2026-11-07
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0515 flai test on a change to flai/internal/hostapi or the host configuration runs neither flaiover's server tests that drive the real flai nor TestSettingsIndexIsComplete, so such a story runs them itself before its close-out (I-0304)

S-0587's first two close-outs each failed on a test that `flai test` had not run for the story's paths. Both are recorded in I-0304.

- First close-out: `flaiover/src/lib/server/writes.test.ts` failed. It runs the built flai and still expected a 403 from the `settings.agent` gate that S-0587 made per project in `flai/internal/hostapi/writes.go`. vitest's `related` run picks tests by the TypeScript import graph, so no change under `flai/internal/hostapi` selects it.
- Second close-out: `TestSettingsIndexIsComplete` (`flai/cmd/settings_doc_test.go`) failed, because `docs/operators/settings.md` named `agent.projects.<root>` instead of each key below it. That test skips under `-short`, so `flai test` on `docs/operators/settings.md` and `flai/internal/config` runs it nowhere. A task sub-agent had reported that it passed.

S-0587's narrative: "`flai test` runs neither for these paths, recorded as I-0304." Until S-0592, which remediates I-0304, changes this, a story that changes a host API gate or the host configuration's shape runs these tests itself before it asks for the close-out. Run vitest on `flaiover/src/lib/server/writes.test.ts` with a freshly built flai. Run `go test -run TestSettingsIndexIsComplete ./cmd` in `flai/` without `-short`.

Sources: S-0587 (wip/archive/agents/S-0587.md, Decisions); I-0304, Instances.
