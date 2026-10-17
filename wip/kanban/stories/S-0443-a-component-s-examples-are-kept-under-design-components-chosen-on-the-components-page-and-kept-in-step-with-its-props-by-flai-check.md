---
id: S-0443
type: story
nature: feature
title: A component's examples are kept under design/components, chosen on the Components page, and kept in step with its props by flai check
status: ready
parent: E-0026
owner: alex
created: 2026-10-16T21:38:09Z
updated: 2026-10-17T03:34:53Z
transitions:
  - to: ready
    at: 2026-10-17T01:02:23Z
    by: alex
tags: [cli, dashboard, template]
topics: [components, dashboard, cli, conventions]
touches: [design/adrs, design/adrs/README.md, design/system/component-library.md, flai/internal/components/examples.go, flai/internal/components/examples_test.go, flai/internal/check/check.go, flai/internal/check/components_test.go, flai/internal/docedit/file.go, flai/internal/docedit/file_test.go, flai/internal/hostapi/components.go, flai/internal/hostapi/components_test.go, flai/cmd/components.go, flai/cmd/components_test.go, flaiover/src/lib/components/explorer/ExampleList.svelte, flaiover/src/lib/components/explorer/ExampleList.svelte.test.ts, flaiover/src/routes/api/components/examples/+server.ts, flaiover/src/routes/api/components/examples/examples.test.ts, flaiover/src/routes/design/components/+page.svelte, design/system/flai-cli.md, docs/users/flai.md, docs/users/flaiover.md, design/system/flaiover-dashboard.md, design/conventions/documentation.md, template/root/design/conventions/documentation.md, design/system/documentation-standard.md, design/components/README.md, template/root/design/components/README.md, design/README.md, template/root/design/README.md.tmpl, CLAUDE.md, template/root/CLAUDE.md.tmpl, design/system/repository-layout.md, template/CHANGELOG.md, template/template.yaml]
after: [S-0436, S-0440, S-0442]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: sum
  seconds: 0
  models: []
  strategic:
    - kind: orchestrator
      seconds: 380
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 43
          output: 796
          cache_read: 5625626
          cache_write: 10275
          cost: 1.413
cost_of_delay:
  value: 88.07
  by: planner-E-0026
  at: 2026-10-16T21:50:07Z
forecast:
  duration: 1h4m
  delivery: 2026-10-17T04:48:00Z
  basis: "Its own forecast of 1h4m; 2nd in the pull order with an in-progress limit of 5, behind S-0442 and S-0444."
  by: flai
  at: 2026-10-17T03:34:53Z
finalized:
  by: orchestrator
  at: 2026-10-16T21:52:05Z
---
# S-0443 A component's examples are kept under design/components, chosen on the Components page, and kept in step with its props by flai check

## Goal

Every component can be seen and worked with (outcome 1), including one that renders only with real data, such as a card that needs an item. A component's examples are named sets of props. On TH-0483 the operator put supporting files under `design/components/`, so each component's examples are one JSON file there, at a path that mirrors the component's, versioned with the design and in reach of `file.get` and `file.save` as the design folder already is. The Components page lists them under the component. Choosing one loads its props into the panel, and the operator saves what they built as a new example, on main or on the story the page shows. `flai check` reports an example that no longer fits its component, so a story that renames a prop cannot pass its close-out with a stale example.

## Acceptance criteria

- [ ] An ADR records the examples on TH-0483's answer 4. They live under `design/components/`, one file per component at `design/components/<sub-project>/<the component's path in the sub-project, .svelte replaced by .examples.json>`, unless the ADR finds a better layout and says why. Each file holds `examples`, each with `name`, `props`, and optional `theme` and `width`. `design/components/` is in `file.get`'s and `file.save`'s reach as part of the design folder, so ADR-0147's reach is unchanged. A save with no hash creates a file that is absent.
- [ ] The documentation convention allows the examples files as a second exception to "Everything is markdown", in `design/conventions/documentation.md`, `template/root/design/conventions/documentation.md`, and `design/system/documentation-standard.md`. `design/components/README.md`, here and in the template, says what the folder holds. `design/README.md` and its template, `CLAUDE.md` and its template, and `design/system/repository-layout.md` list it.
- [ ] `flai/internal/components` reads each component's examples, and `components.list` and `flai components list` carry them.
- [ ] `flai check` reports `components.examples` for an examples file that does not parse, an example that names a prop the component lacks, misses a required prop, or gives a value that does not fit its prop's kind, and an examples file whose component does not exist. Each finding names the example and the line. A project with no Svelte sub-project, or without `node`, is passed over with one note.
- [ ] The Components page lists the picked component's examples. Choosing one loads its props, theme, and width. "Save as example" asks for a name and saves through `file.save` on main, or on the story the page's source names. A conflict or a refusal is shown with its findings, and the panel's values are kept.
- [ ] Tests cover the reading, each check finding, the create, and the page's list and save. `design/system/component-library.md`, `design/system/flai-cli.md`, `docs/users/flai.md`, `docs/users/flaiover.md`, `design/system/flaiover-dashboard.md`, and `template/CHANGELOG.md` describe examples.

## Tasks
- T-1884 An ADR records the examples files under design/components, their path and JSON, and that a save with no hash creates one
- T-1885 flai/internal/components reads each component's examples under design/components, and flai check reports components.examples for one that does not fit its component
- T-1886 file.save with no hash creates a file that is absent, so the first example of a component can be saved
- T-1887 components.list and flai components list carry each component's examples and the path of its examples file
- T-1888 The Components page lists a component's examples, loads one into the panel, and saves the panel's values as a new example through file.save
- T-1889 The CLI design, the user docs, and the component library's design describe examples, their check, and saving one from the page
- T-1894 The documentation convention allows the examples files under design/components, the folder has a README here and in the template, and the layout lists it

## Notes

- Drafted by the planner for E-0026, after S-0436, whose `file.save` it saves through, S-0440, whose panel it fills, and S-0442, whose source picker names the story a save goes to.
- TH-0483 answer 4, 2026-10-16: "lets put supporting files under design/components". The draft had them beside each component; this version follows the answer. The folder is named for supporting files, so a later story may keep more there.
- `flai check` reads no code today (`design/system/design-system-tool.md` § How agents build from it). This check reads only the examples files and S-0437's catalogue, which runs `node`, so it is passed over where `node` is missing rather than failing every project. `flai check` passes over files under `design/` that are not markdown, so the JSON adds no finding of its own; only the convention's words need the exception.

### Planning

- **Touches.** The story declared none; every touch comes from its seven tasks, file by file.
  - **From design:** ADR-0147, whose reach already covers the design folder, so only the create needs `flai/internal/docedit/file.go` from S-0436; the documentation convention's token-file exception, which S-0381 carried to both copies and the standard, the pattern for this one; `design-system-tool.md` § How agents build from it, for a check that reads no code.
  - **From layout:** `examples.go` in `flai/internal/components`; a `checker.components` in `flai/internal/check/check.go` with a test file of its own; S-0437's and S-0438's list in the host API and the CLI; `ExampleList.svelte` and a route under `/api/components/examples/`; the folder's README here and in `template/root/design/`, with the layout's indexes as S-0428 adds `design/digest/`.
  - **Co-change, added:** `design/system/flaiover-dashboard.md` (28%), for the new route; `template/CHANGELOG.md` and `template/template.yaml`, for the template's new folder and convention text.
  - **Co-change, not added:** `docs/operators/index.md` and `settings.md` (19%, 13%), since nothing on the host changes; `flai/internal/check/check_test.go` (4%), since a new test file is used.
  - **Folder touch kept:** `design/adrs`, for the new ADR, in `claims.shared`.
- **Order.** `after` S-0442 as well as S-0436 and S-0440: the save on a story's branch needs S-0442's picker. Both stories change `flaiover/src/routes/design/components/+page.svelte`, so the hold would run them one after the other anyway.
- **Forecast.** Kept at flai's 1h4m for this version: size 40, at 96 s a unit over 61 done large feature stories. The convention exception, the folder README, and the layout indexes added twelve touches after TH-0483's answer 4. Delivery is flai's play-out, after S-0442.
- **Cost of delay.** 88.07 USD a week, flai's share of E-0026's 600 USD penalty by forecast duration, kept.
