---
id: D-0345
title: flai move to done refuses unfixed lint as text only, so the dashboard picks each file's checkout from the refusal's words
kind: cause
sources: [I-0145]
topics: [cli, dashboard, flai, go, flaiover, sveltekit, code, all]
paths: [flai/cmd/move.go, flai/cmd/move_refused_test.go, flai/internal/hostapi/writes.go, flai/internal/hostapi/writes_test.go, docs/users/flai.md, design/system/flai-cli.md, flaiover/src/lib/components/AcceptFindings.svelte, flaiover/src/lib/components/AcceptFindings.svelte.test.ts, design/components/flaiover/src/lib/components/AcceptFindings.examples.json, flaiover/src/lib/components/Review.svelte, flaiover/src/lib/components/Review.svelte.test.ts, flaiover/src/routes/board/+page.svelte, flaiover/src/routes/board/board.svelte.test.ts, "flaiover/src/routes/items/[id]/+page.svelte", "flaiover/src/routes/items/[id]/item.svelte.test.ts", "flaiover/src/routes/api/items/[id]/move/+server.ts", flaiover/src/lib/filelinks.ts, flaiover/src/lib/filelinks.test.ts, design/system/flaiover-dashboard.md, docs/users/flaiover.md, design/issues/I-0145-flai-move-to-done-refuses-unfixed-lint-as-text-only-so-the-dashboard-picks-each-file-s-checkout-from-the-refusal-s-words.md, design/issues/summary.md, docs/users/flai-reference.md]
date: 2026-10-21
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0345 flai move to done refuses unfixed lint as text only, so the dashboard picks each file's checkout from the refusal's words

I-0145 was closed on 2026-10-21.

The cause:

```text
flai move to done refuses unfixed lint as text only, so the dashboard picks each file's checkout from the refusal's words
```

The fix:

```text
S-0465: flai move <story> done --json now refuses as flai accept --json does (exit 4, the refused object with each finding's checkout), item.move answers it as Refused with that data, and the board and the item page list the findings through AcceptFindings, so linkSegments no longer reads flai's wording to pick a file's checkout.
```

Source: I-0145.
