---
id: D-0321
title: The Components page frames a component that draws a fixed overlay, such as a dialog, at 160px, as its content's height is 0, cutting the dialog off
kind: cause
sources: [I-0160]
topics: [flai, cli, go, flaiover, dashboard, sveltekit, code, all]
paths: [design/system/component-library.md, docs/users/flai.md, flai/internal/components/harness/harness.html, flai/internal/components/harness/harness.mjs, flai/internal/components/integration_test.go, flai/internal/components/testdata/harness/src/lib/Dialog.svelte, flai/internal/components/testdata/harness/src/lib/Leaky.svelte, design/system/flaiover-dashboard.md, docs/users/flaiover.md, flaiover/src/lib/components/explorer/ComponentFrame.svelte, flaiover/src/lib/components/explorer/ComponentFrame.svelte.test.ts, design/issues/I-0160-the-components-page-frames-a-component-that-draws-a-fixed-overlay-such-as-a-dialog-at-160px-as-its-content-s-height-is-0-cutting-the-dialog-off.md, design/issues/summary.md, design/adrs/0172-the-component-harness-answers-a-framed-component-s-requests-to-the-app-s-own.md, design/adrs/README.md, flai/internal/components/testdata/harness/src/lib/Fetcher.svelte, design/issues/I-0164-the-components-page-cannot-show-a-component-that-asks-the-app-s-server-in-its-answered-state-since-the-harness-answers-every-such-request-with-501.md, design/issues/I-0165-flai-guard-refuses-a-story-s-sub-agent-flai-components-list-and-flai-components-serve-so-a-headless-pass-of-the-harness-falls-to-the-story-s-agent.md, design/system/flai-cli.md, docs/users/flai-reference.md, flai/cmd/components.go, flai/internal/components/props.go, flai/internal/components/props.mjs, flai/internal/components/props_test.go, flai/internal/components/testdata/typed/package.json, flai/internal/components/testdata/typed/src/lib/Board.svelte, flai/internal/components/testdata/typed/src/lib/Imported.svelte, flai/internal/components/testdata/typed/src/lib/types.ts, flai/internal/components/testdata/typed/tsconfig.json, flaiover/src/lib/components/explorer/PropsPanel.svelte.test.ts, flaiover/src/lib/explorer/controls.test.ts, flaiover/src/lib/explorer/controls.ts, flaiover/src/lib/server/repo.ts, flaiover/src/routes/design/components/components.svelte.test.ts, design/issues/I-0166-flai-components-list-and-the-components-examples-check-take-about-2-3s-for-flaiover-since-props-mjs-resolves-json-props-types-with-typescript-up-from-0-4s.md]
date: 2026-10-19
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0321 The Components page frames a component that draws a fixed overlay, such as a dialog, at 160px, as its content's height is 0, cutting the dialog off

I-0160 was closed on 2026-10-19.

The cause:

```text
The Components page frames a component that draws a fixed overlay, such as a dialog, at 160px, as its content's height is 0, cutting the dialog off
```

The fix:

```text
S-0449 (T-1904): the harness measures every element a component draws, fixed dialogs among them, leaving out a backdrop pinned to the viewport, so the frame holds a dialog whole; a browser check in flai/internal/components/integration_test.go frames a 300px dialog whole.
```

Source: I-0160.
