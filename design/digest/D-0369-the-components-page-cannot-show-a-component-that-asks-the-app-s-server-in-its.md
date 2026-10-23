---
id: D-0369
title: The Components page cannot show a component that asks the app's server in its answered state, since the harness answers every such request with 501
kind: cause
sources: [I-0164]
topics: [components, dashboard, cli, flai, go, flaiover, sveltekit, template, conventions, code, all]
paths: [design/adrs/README.md, design/system/component-library.md, design/components/README.md, template/root/design/components/README.md, flai/internal/components/examples.go, flai/internal/components/examples_test.go, flai/internal/check/components_test.go, flai/cmd/components.go, flai/cmd/components_test.go, flai/internal/hostapi/components_test.go, flai/internal/components/harness/harness.mjs, flai/internal/components/integration_test.go, flai/internal/components/testdata/harness/src/lib/Fetcher.svelte, flaiover/src/routes/api/components/examples/+server.ts, flaiover/src/routes/api/components/examples/examples.test.ts, flaiover/src/lib/server/repo.ts, flaiover/src/lib/explorer/protocol.ts, flaiover/src/lib/explorer/protocol.test.ts, flaiover/src/lib/components/explorer/ComponentFrame.svelte, flaiover/src/lib/components/explorer/ComponentFrame.svelte.test.ts, flaiover/src/lib/components/explorer/ExampleList.svelte, flaiover/src/lib/components/explorer/ExampleList.svelte.test.ts, flaiover/src/routes/design/components/+page.svelte, flaiover/src/routes/design/components/components.svelte.test.ts, design/components/flaiover/src/lib/components/AcceptConfirm.examples.json, design/components/flaiover/src/lib/components/ActivityView.examples.json, design/components/flaiover/src/lib/components/AnalyzerPanel.examples.json, design/components/flaiover/src/lib/components/CancelConfirm.examples.json, design/components/flaiover/src/lib/components/ImportPrompt.examples.json, design/components/flaiover/src/lib/components/InboxView.examples.json, design/components/flaiover/src/lib/components/Messages.examples.json, design/components/flaiover/src/lib/components/NewAdrForm.examples.json, design/components/flaiover/src/lib/components/NewItemForm.examples.json, design/components/flaiover/src/lib/components/OrchestratorPanel.examples.json, design/components/flaiover/src/lib/components/PlannerPanel.examples.json, design/components/flaiover/src/lib/components/PublishBanner.examples.json, design/components/flaiover/src/lib/components/Threads.examples.json, design/components/flaiover/src/lib/components/WorktreeUncommitted.examples.json, docs/users/flai.md, docs/users/flaiover.md, design/issues/I-0164-the-components-page-cannot-show-a-component-that-asks-the-app-s-server-in-its-answered-state-since-the-harness-answers-every-such-request-with-501.md, design/issues/summary.md, design/adrs/0187-an-example-may-carry-the-answers-the-component-harness-gives-a-framed-component.md, template/CHANGELOG.md, flai/internal/components/testdata/harness/src/lib/Answered.svelte, docs/users/flai-reference.md, design/issues/I-0194-the-mcp-task-done-checks-a-story-with-the-installed-flai-so-a-story-that-widens-a-check-s-format-fails-its-own-task-close.md]
date: 2026-10-22
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0369 The Components page cannot show a component that asks the app's server in its answered state, since the harness answers every such request with 501

I-0164 was closed on 2026-10-22.

The cause:

```text
The Components page cannot show a component that asks the app's server in its answered state, since the harness answers every such request with 501
```

The fix:

```text
S-0476: an example may carry answers (ADR-0187), which the component harness gives a framed component's fetch and EventSource to the app's own server in place of the 501, and the Components page sends, edits, and saves them. The fourteen components carry examples with answers, and a headless pass rendered their 17 answered examples with no error and no 501. integration_test.go's harness test reproduces it: an answered fetch and EventSource, and an unmatched request still 501.
```

Source: I-0164.
