---
id: ADR-0263
title: "A kluster block may list components, each built as its own step with its own log, a deploy records every step it will run before the first starts, kluster.log reads a step's log from an offset, and the Deployments page holds the deploy setup, the builds, and the live logs in an in-page rail"
status: accepted
date: 2026-11-07
supersedes: []
superseded_by: []
refines: [ADR-0254, ADR-0259]
---

# ADR-0263 A kluster block may list components, each built as its own step with its own log, a deploy records every step it will run before the first starts, kluster.log reads a step's log from an offset, and the Deployments page holds the deploy setup, the builds, and the live logs in an in-page rail

## Context

S-0595: under Ops › Deployments a project with a `kluster:` block but no `deploy` showed one button, and pressing it answered `the kluster: block of system-flow.yaml names no deploy command`. The deploy setup lived on the Klusters page alone. A deploy's build was one command and one step, `build`, with one log, so a project that builds three images showed one build and could not say which image failed. A step's log reached the dashboard only as the last 4 KB of a finished step, and a running step was not in the record until it ended ([ADR-0254](0254-flai-runs-one-kluster-per-project-from-the-config-its-manifest-s-kluster-block.md), [ADR-0259](0259-enabling-kluster-on-a-project-with-no-kluster-block-writes-one-naming-only-its.md)).

The operator asked that the operator never edit `system-flow.yaml` to set up or control deployments, that the Deployments page group its controls, live logs, and status in the page's own vertical rail, and that a Build sub-page show one panel per component's build, with its status and an expandable live feed.

## Decision

A `kluster:` block may list components, each built as its own step with its own log; a deploy records every step it will run before the first starts; `kluster.log` reads a step's log from an offset; and the Deployments page holds the deploy setup, the builds, and the live logs in an in-page rail.

1. **Components.** The `kluster:` block gains an optional key `components`, a list of `{name, build}`: `name` a DNS label unique in the list, `build` an argument list as `build` is. A block names `build` or `components`, not both; `flai check` reports both, a missing or repeated name, a name that is not a DNS label, and a component with no program, as `manifest.kluster`. `flai manifest set kluster.components '<JSON list>'` writes it, `[]` unsets it, and the dashboard writes it through `settings.manifest` under the `settings` host action, as the other keys are.
2. **Steps.** With components, a deploy runs one build step per component, in the list's order, named `build-<name>`, with its log at `logs/<namespace>/build-<name>.log`, its own images file at `logs/<namespace>/images-<name>` as `FLAI_KLUSTER_IMAGES`, and `FLAI_KLUSTER_COMPONENT` set to the name; then one load of every image they listed, then `deploy`. A failed build stops the deploy, and the steps after it are recorded `skipped`. With `build`, the step stays `build`, as before.
3. **The plan in the record.** A deploy, a teardown, and a promotion write every step they will run into the record as `pending` before the first starts, mark each `running` with its start and its log when it starts, and give it its outcome when it ends. So the record says at every moment which steps passed, which runs, which failed, and which wait.
4. **`kluster.log`.** A host read, whatever the host action says, `{namespace, step, offset}`, answers at most 64 KB of that step's log from `offset`, with `next`, the offset to ask from next, the log's `size`, the step's `outcome`, and `reset` true when the log is shorter than `offset`, as when a new deploy cleared it, and the text then starts at 0. The dashboard asks it about once a second while the step runs, which is the live feed.
5. **The Deployments page.** It has an in-page vertical rail, the item page's `PageSidebar`, with the panes Overview (main's namespace and each story's, their status and their actions), Build (one panel per build step of a chosen namespace, each with its status and an expandable live feed), Logs (every step of the chosen namespace's last job, each with its live feed), Access (the port-forwards), and Setup (the block's `namespace`, `components`, `build`, `deploy`, `teardown`, and `timeout`, and **Write these with an agent**). The choice is kept in the URL's hash, as the item page keeps it, rather than in routes of their own. With no `deploy`, **Deploy** and **Promote** open the Setup pane rather than ask flai and show its refusal. The Klusters page's step 4 uses the same setup form.

## Consequences

- A project can say what it builds, one component a line, and the operator sees each component's build pass, run, fail, or wait, with its own log, live.
- The record's step `outcome` gains `pending`. A reader of an older flai's record finds no pending step.
- The settings catalog gains a kind, `components`, beside `text`, `command`, and `duration`.
- The setup is on both Ops pages, through one form, so they cannot drift.
- The live feed polls; it does not stream through the events socket ([ADR-0182](0182-the-dashboard-s-live-events-reach-the-browser-over-a-websocket.md)). A request a second while a step runs is within what the channel carries.

## Alternatives considered

- **One build command whose output marks each image.** Parsing a project's output is fragile and gives one log for all images.
- **Components built in parallel.** Faster on a large host, but their logs interleave in the host's resources and a failure stops nothing else cleanly; one at a time keeps the job lock's single slot simple. A later story can add a `parallel` key.
- **Routes of their own for each pane, such as `/ops/deployments/build`.** flaiover has no nested layout yet, and the item page's hash-kept panes are the existing pattern; links to a pane still work through the hash.
- **The live feed over the events socket.** It would push lines with no polling, but the socket carries work-item changes, not files, and would need flai serve to watch log files.
