---
id: TH-0589
title: "S-0483 not enqueued: component-library.md still says the main checkout is never installed into, which TH-0578 asked to change"
anchor:
  path: wip/kanban/stories/S-0483-accepting-a-story-that-adds-a-flaiover-dependency-leaves-the-main-checkout-s-flaiover-node-modules-without-it-so-every-close-out-s-integration-tier-fails-testeverylistedflaiovercomponentcompilesintheharness-until-someone-installs.md
  item: S-0483
status: resolved
participants: [orchestrator, alex]
created: 2026-10-21T20:52:39Z
updated: 2026-10-22T20:48:17Z
---

# TH-0589 S-0483 not enqueued: component-library.md still says the main checkout is never installed into, which TH-0578 asked to change

On wip/kanban/stories/S-0483-accepting-a-story-that-adds-a-flaiover-dependency-leaves-the-main-checkout-s-flaiover-node-modules-without-it-so-every-close-out-s-integration-tier-fails-testeverylistedflaiovercomponentcompilesintheharness-until-someone-installs.md.

## Entries

### 2026-10-21T20:52:39Z orchestrator
**Recommended:** send S-0483 back to in-progress for one small change, then it goes to the queue.

- Rewrite the last sentence of "The install step" in `design/system/component-library.md` (line 232). It should say that only `flai accept` installs into the main checkout, by the manifest's `install` command when the accepted change touched its `install.paths`. Give the reason, I-0179, and link ADR-0185.
- Add the file to the story's touches.

## What I checked

`flai verify S-0483` passed at the branch head, 42b03609. Integration and smoke are left to the merge queue. The verifier matched both criteria to changed files within the touches:

- **Criterion 1:** `flai accept` runs `projects[].install.command` after the commit (ADR-0185, `flai/cmd/accept_install.go`), with four accept tests and the manifest check.
- **Criterion 2:** I-0179 is closed with its reason.

## The blocker

On TH-0578 you said: "Go ahead with 483 as planned and update the rule it is breaking to include the exception and reason for this edge case." The diff does not touch `design/system/component-library.md`, and ADR-0185 does not mention the rule. As it stands, the living design contradicts the shipped behaviour.

## A minor point, not blocking

The planner said ADR-0185 would settle the install's timeout. It states none, and `installAccepted` runs the command under `context.Background()`, so a hung install would hang the acceptance. The fix: one line in the ADR saying there is no timeout, or a timeout in the code.

### 2026-10-22T20:48:17Z alex
Resolved: S-0483 was accepted
