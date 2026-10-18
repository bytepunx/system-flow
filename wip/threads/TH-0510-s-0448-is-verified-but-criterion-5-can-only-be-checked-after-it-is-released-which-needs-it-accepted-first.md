---
id: TH-0510
title: S-0448 is verified but criterion 5 can only be checked after it is released, which needs it accepted first
anchor:
  path: wip/kanban/stories/S-0448-the-components-page-lists-starts-and-stops-through-the-dashboard-the-components-methods-take-their-sub-project-as-sub-project-and-no-call-s-params-can-overwrite-the-channel-s-project-key.md
  item: S-0448
status: open
participants: [orchestrator]
created: 2026-10-18T02:54:34Z
updated: 2026-10-18T02:54:34Z
---

# TH-0510 S-0448 is verified but criterion 5 can only be checked after it is released, which needs it accepted first

On wip/kanban/stories/S-0448-the-components-page-lists-starts-and-stops-through-the-dashboard-the-components-methods-take-their-sub-project-as-sub-project-and-no-call-s-params-can-overwrite-the-channel-s-project-key.md.

## Entries

### 2026-10-18T02:54:34Z orchestrator
**Recommendation:** accept S-0448 yourself, or tick criterion 5 and tell me here so that I accept it. Then do criterion 5's hand check after the release, once the host runs the new flai and flaiover.

**Why the order has to change.** Criterion 5 asks for the live Components page to work on this host with the fixed flai and flaiover. That needs:
- this story released, and `flai release --pending` publishes only accepted stories;
- the host's `flai serve` and dashboard upgraded to that release, which only you can do.

So the check cannot come before the acceptance.

## What is verified

- **flai verify passed** at f47f65a64b86, the head of story/S-0448, against main 299af23. Every tier passed, integration and smoke included.
- **Criteria 1 to 4 and 6 are met.**
  - **Sub-project parameter:** the methods take `sub_project`, and no params type reads `project`.
  - **Repo:** it sends `sub_project`, and the routes keep `project`.
  - **The channel key:** the agent's channel key comes after the params, so params can't overwrite it.
  - **Tests:**
    - A test sends all four methods through a real channel with the project key.
    - A source-scanning test fails on any params field named `project`.
  - **Docs:** they name the parameter and say why it is named that way.
- **Agent's partial check of criterion 5 (notes, 2026-10-18).**
  - The installed flai 1.57.0 refuses `components.status '{"project":"sf"}'`.
  - A flai built from the branch answers the same call, and lists flaiover's 54 components.

## What is missing

- **Criterion 5:** the live page on this host, covering list, start, frame and stop, with the versions used. It is unticked.

Ticking or deferring a criterion is yours, not mine, so S-0448 stays in review. S-0446 and S-0447, which share `component-library.md` with it, are already in ready.
