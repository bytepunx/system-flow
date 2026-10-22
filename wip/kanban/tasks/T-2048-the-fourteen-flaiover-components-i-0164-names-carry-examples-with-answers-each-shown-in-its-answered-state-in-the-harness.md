---
id: T-2048
type: task
nature: improvement
title: The fourteen flaiover components I-0164 names carry examples with answers, each shown in its answered state in the harness
status: backlog
parent: S-0476
owner: alex
created: 2026-10-22T21:00:37Z
updated: 2026-10-22T21:00:37Z
transitions: []
stream: S-0476
tags: [flaiover]
touches: [design/components/flaiover/src/lib/components/AcceptConfirm.examples.json, design/components/flaiover/src/lib/components/ActivityView.examples.json, design/components/flaiover/src/lib/components/AnalyzerPanel.examples.json, design/components/flaiover/src/lib/components/CancelConfirm.examples.json, design/components/flaiover/src/lib/components/ImportPrompt.examples.json, design/components/flaiover/src/lib/components/InboxView.examples.json, design/components/flaiover/src/lib/components/Messages.examples.json, design/components/flaiover/src/lib/components/NewAdrForm.examples.json, design/components/flaiover/src/lib/components/NewItemForm.examples.json, design/components/flaiover/src/lib/components/OrchestratorPanel.examples.json, design/components/flaiover/src/lib/components/PlannerPanel.examples.json, design/components/flaiover/src/lib/components/PublishBanner.examples.json, design/components/flaiover/src/lib/components/Threads.examples.json, design/components/flaiover/src/lib/components/WorktreeUncommitted.examples.json]
after: [T-2038, T-2045]
---
# T-2048 The fourteen flaiover components I-0164 names carry examples with answers, each shown in its answered state in the harness

## Work

Give each of the fourteen components I-0164 names an example whose `answers` let the harness show its answered state: AcceptConfirm, ActivityView, AnalyzerPanel, CancelConfirm, ImportPrompt, InboxView, Messages, NewAdrForm, NewItemForm, OrchestratorPanel, PlannerPanel, PublishBanner, Threads, and WorktreeUncommitted, all under `flaiover/src/lib/components/`. Twelve have an examples file under `design/components/flaiover/src/lib/components/` already; Messages and Threads have none, so make theirs.

For each component, read what it asks of the app's server (`api()` calls: the method, the path, and what it reads from the answer) and write one answer per request, with a body shaped as the real route answers, small and plausible: a short list for a list view, one item for a confirmation's preview, an empty list where the empty state is the one worth showing too. Where a component opens an `EventSource`, give its path `events` only when its view changes on them; otherwise leave it to fail as before. Keep every example that exists, and add the answers to the example whose props match the answered state, or add a new example named for it, such as "With two threads".

Verify each in the harness: run `flai components serve` (or the dashboard with the `components` host action on) and choose each example on the Components page, or drive the harness headless as S-0449's pass did, and check that the component renders its answered state and that no `fetch` event in the log carries a 501 for a request the example meant to answer. `flai check --strict` passes on the folder.

This task waits for T-2038, so that `flai check` accepts the key, and for T-2045, so that the page sends the answers to the harness; through T-2045 it waits for T-2040 and T-2041 too.

## Done when

- [ ] Each of the fourteen components has an example with `answers` that renders its answered state in the harness, with no 501 for a request the example answers
- [ ] Messages and Threads have an examples file
- [ ] `flai check --strict` reports no `components.examples` finding on the folder
- [ ] Which components were verified, and how, is noted in the story's narrative

## Notes

Planned by planner-S-0476. The bodies are fixtures, not a contract: when a route's answer shape changes, the example changes with it in the same story, as the check keeps props in step with the component.
