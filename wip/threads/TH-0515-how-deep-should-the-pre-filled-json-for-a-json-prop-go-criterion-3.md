---
id: TH-0515
title: How deep should the pre-filled JSON for a json prop go (criterion 3)?
anchor:
  path: wip/kanban/stories/S-0449-improvements-to-design-components-screen-part-1.md
  item: S-0449
status: resolved
participants: [agent-S-0449, orchestrator, alex]
created: 2026-10-19T05:02:12Z
updated: 2026-10-19T17:30:21Z
---

# TH-0515 How deep should the pre-filled JSON for a json prop go (criterion 3)?

On wip/kanban/stories/S-0449-improvements-to-design-components-screen-part-1.md.

## Entries

### 2026-10-19T05:02:12Z agent-S-0449
Criterion 3 asks for the JSON of a json prop to be pre-filled with its schema, even when every value is empty. Most flaiover props are typed with imported names, such as `stories: Story[]` or `from: Version`. Today flai reads a prop's type only as far as the component's own script declares it, so for those props it knows only "array" or "unknown".

Recommendation: (A) flai resolves each json prop's type with the sub-project's own `typescript` package, imports and `$lib` aliases included, and the catalogue carries an empty skeleton built from it, such as `[{"id": "", "title": "", "status": "", "tasks": 0}]` for `Story[]`. The props panel fills the field with that skeleton and sends it to the component, so what the field shows is what the component gets. Where flai cannot resolve a type, the field falls back to `[]`, `{}`, or `null`. This is the larger change: it adds type resolution to `flai/internal/components/props.mjs`.

Alternatives:
- (B) A skeleton only from types the component declares itself, inline or in its own script. Imported types fall back to `[]` or `{}`. This is small, but most flaiover fields would show only `[]` or `{}`.
- (C) As (A), but the skeleton is only shown, as a placeholder or template, and is not sent until you edit it. Components keep getting `undefined` until then.

Also: should a one-element array be the skeleton for an array type, so its item's shape shows (`[{...}]`), or an empty `[]` with the shape elsewhere? I recommend one element.

Until you answer, I am working the other four criteria: the 10px margin, sizing the frame, rendering only the selected component, and the 500s.

### 2026-10-19T05:02:42Z orchestrator (recommendation)
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

### 2026-10-19T13:57:46Z alex
Confirmed the recommendation of 2026-10-19T05:02:42Z orchestrator.

Source: design/system/component-library.md § Controls for props (S-0440)

### 2026-10-19T17:30:21Z alex
Resolved: S-0449 was accepted
