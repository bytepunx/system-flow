---
id: D-0476
title: A field added to manifest.Agent is a work-item front-matter field too, since a story's agent shares the type, so it changes front-matter-fields.txt and raises flai.minimum even when no story carries it
kind: decision
sources: [S-0510, S-0745]
topics: [release, manifest, flai, code, go]
paths: [flai/internal/manifest/agent.go, flai/internal/workitem/front-matter-fields.txt, flai/internal/workitem/store_test.go, docs/operators/settings.md]
date: 2026-11-03
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0476 A field added to manifest.Agent is a work-item front-matter field too, since a story's agent shares the type, so it changes front-matter-fields.txt and raises flai.minimum even when no story carries it

S-0510 added `agent.image.layer` to the project's agent, as a field of `manifest.Agent` in `flai/internal/manifest/agent.go`. A story's front matter carries an `agent` of the same type, so the field showed up as `item.agent` gaining `image` in `flai/internal/workitem/front-matter-fields.txt`. From S-0510's Decisions (`wip/archive/agents/S-0510.md`):

```text
`Agent.Image` makes `item.agent` gain `image` in `flai/internal/workitem/front-matter-fields.txt`, so the release raises `flai.minimum` though no story carries it. I accepted this rather than give the project's agent a type of its own, which every caller of `Manifest.Agent` would have had to change.
```

A later story that adds a key to the manifest's `agent`, even one `Agent.With` never copies into a story, should expect the same raise at the next publish (D-0401). If the raise matters, the other choice is a type of the project's own for its agent. S-0510 turned that down for its cost.

Since the field is a field of `Item` too, `TestCloneSharesNothing` in `flai/internal/workitem/store_test.go` must set it: the test fails on any nil pointer, slice, or map of `Item` with `fill it in this test`. S-0745 met this when it added `agent.image.args` to `manifest.AgentImage`. From S-0745's Decisions (`wip/archive/agents/S-0745.md`):

```text
T-3173 also changed `flai/internal/workitem/store_test.go`, whose clone test must fill every field of Item, and added the `agent.image.args.<name>` row to `docs/operators/settings.md`, which TestSettingsIndex needs.
```

So a task that adds a field to `manifest.Agent` or a type under it touches `store_test.go` and, for a manifest key, `docs/operators/settings.md`, beside `agent.go`.

Sources: S-0510, S-0745.
