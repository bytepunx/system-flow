---
id: D-0322
title: "How deep should the pre-filled JSON for a json prop go (criterion 3): Recommendation: go with (A) and a one-element array"
kind: answer
sources: [TH-0515]
topics: [flai, cli, go, flaiover, dashboard, sveltekit, code, all]
paths: [design/system/component-library.md, docs/users/flai.md, flai/internal/components/harness/harness.html, flai/internal/components/harness/harness.mjs, flai/internal/components/integration_test.go, flai/internal/components/testdata/harness/src/lib/Dialog.svelte, flai/internal/components/testdata/harness/src/lib/Leaky.svelte, design/system/flaiover-dashboard.md, docs/users/flaiover.md, flaiover/src/lib/components/explorer/ComponentFrame.svelte, flaiover/src/lib/components/explorer/ComponentFrame.svelte.test.ts, design/issues/I-0160-the-components-page-frames-a-component-that-draws-a-fixed-overlay-such-as-a-dialog-at-160px-as-its-content-s-height-is-0-cutting-the-dialog-off.md, design/issues/summary.md, design/adrs/0172-the-component-harness-answers-a-framed-component-s-requests-to-the-app-s-own.md, design/adrs/README.md, flai/internal/components/testdata/harness/src/lib/Fetcher.svelte, design/issues/I-0164-the-components-page-cannot-show-a-component-that-asks-the-app-s-server-in-its-answered-state-since-the-harness-answers-every-such-request-with-501.md, design/issues/I-0165-flai-guard-refuses-a-story-s-sub-agent-flai-components-list-and-flai-components-serve-so-a-headless-pass-of-the-harness-falls-to-the-story-s-agent.md, design/system/flai-cli.md, docs/users/flai-reference.md, flai/cmd/components.go, flai/internal/components/props.go, flai/internal/components/props.mjs, flai/internal/components/props_test.go, flai/internal/components/testdata/typed/package.json, flai/internal/components/testdata/typed/src/lib/Board.svelte, flai/internal/components/testdata/typed/src/lib/Imported.svelte, flai/internal/components/testdata/typed/src/lib/types.ts, flai/internal/components/testdata/typed/tsconfig.json, flaiover/src/lib/components/explorer/PropsPanel.svelte.test.ts, flaiover/src/lib/explorer/controls.test.ts, flaiover/src/lib/explorer/controls.ts, flaiover/src/lib/server/repo.ts, flaiover/src/routes/design/components/components.svelte.test.ts, design/issues/I-0166-flai-components-list-and-the-components-examples-check-take-about-2-3s-for-flaiover-since-props-mjs-resolves-json-props-types-with-typescript-up-from-0-4s.md]
date: 2026-10-19
status: absorbed
absorbed_by: "design/system/component-library.md#skeletons-s-0449"
asked_again: []
cited: []
---

# D-0322 How deep should the pre-filled JSON for a json prop go (criterion 3): Recommendation: go with (A) and a one-element array

The operator, alex, answered TH-0515 at 2026-10-19T13:57:46Z.

What was asked:

```text
Criterion 3 asks for the JSON of a json prop to be pre-filled with its schema, even when every value is empty. Most flaiover props are typed with imported names, such as `stories: Story[]` or `from: Version`. Today flai reads a prop's type only as far as the component's own script declares it, so for those props it knows only "array" or "unknown".
```

The answer:

```text
**Recommendation:** go with (A) and a one-element array. Resolve each json prop's type with the sub-project's own `typescript`, including imports and `$lib` aliases. Fill the field with the empty skeleton and send it. Fall back to `[]`, `{}` or `null` where a type cannot be resolved. In the same story, update `design/system/component-library.md` to record both rule changes.

**Why this goes to the operator.** (A) changes two rules the design records today:

- **§ The catalogue › Props** says "An imported type is not [resolved]: its kind is `unknown`, with its text kept." (A) resolves imported types.
- **§ Controls for props › What is sent** says "A prop without a value is left out and takes its default." Under (A), a filled skeleton is a value, so it is sent and overrides the component's own default.

These are changes of scope that no ADR settles. The operator wrote the criterion, so the operator decides what it means.

**Why (A) over the others:**

- **(B)** keeps both rules, but it leaves most flaiover json fields at `[]` or `{}`, because their props use imported types (`Story[]`, `Version`). The goal names exactly that problem: "The JSON required doesn't provide any kind of schema." (B) would not meet the criterion's "always".
- **(C)** keeps the "what is sent" rule but not the catalogue rule. The fields would show a schema while required props stay unset, so the required-prop notice would still stand in the frame's place. That is the friction the goal describes.
- **One element** makes an array's item shape visible where the operator edits it.

**Consequences to record:**

- The skeleton replaces the component's default only for json props whose source default is not a literal. A literal default still wins under § Defaults.
- Reset returns a field to its skeleton.

If you agree, confirm this. If you would rather keep the design's rules, answer (C), or (B) for the smallest change.

Source: design/system/component-library.md § Controls for props (S-0440)
```

Source: TH-0515.
