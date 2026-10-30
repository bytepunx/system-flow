---
id: D-0428
title: "release_publish found nothing pending while release_evaluate lists S-0544 as unreleased: yes, add the tag \"flaiover\""
kind: answer
sources: [TH-0759]
topics: [flaiover, dashboard, sveltekit, code, all]
paths: [flaiover/src, flaiover/src/lib/components/NavRail.svelte, flaiover/src/lib/components/NavRail.svelte.test.ts, flaiover/src/lib/hostflai.svelte.ts, flaiover/src/routes/api/agent/+server.ts, flaiover/src/routes/api/agent/agent.test.ts, flaiover/src/lib/components/HostPanel.svelte, flaiover/src/lib/components/HostPanel.svelte.test.ts, flaiover/src/routes/api/dashboard/+server.ts, flaiover/src/routes/api/dashboard/dashboard.test.ts, docs/users/flaiover.md, design/system/flaiover-dashboard.md]
date: 2026-10-29
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0428 release_publish found nothing pending while release_evaluate lists S-0544 as unreleased: yes, add the tag "flaiover"

The operator, alex, answered TH-0759 at 2026-10-29T18:20:39Z.

What was asked: `release_publish` published nothing for S-0544 while `release_evaluate` listed it as pending. The orchestrator then found the cause in the publish warning: S-0544 touched flai and flaiover, but no tag said which project it delivers to.

The answer:

```text
yes, add the tag "flaiover"
```

So a story whose change touches two projects and carries no tag naming one is given the tag of the project it delivers to, and is then published with that project. The orchestrator cannot set the tag: `flai guard` refused its `item_edit` (ADR-0060), so it asked the operator to run `flai edit S-0544 --tags flaiover`. D-0075 holds the same fact for acceptance (I-0024), with this case added.

Source: TH-0759.
