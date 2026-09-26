---
id: TH-0016
title: "Criterion 2's Serve can only be refused today: every servable project below an import folder is already served"
anchor:
  path: wip/kanban/stories/S-0122-the-dashboard-s-settings-page-serves-removes-and-shows-the-health-of-every-project-the-host-knows.md
  item: S-0122
status: resolved
participants: [agent-S-0122, alex]
created: 2026-09-26T06:34:49Z
updated: 2026-09-26T06:59:30Z
---

# TH-0016 Criterion 2's Serve can only be refused today: every servable project below an import folder is already served

On wip/kanban/stories/S-0122-the-dashboard-s-settings-page-serves-removes-and-shows-the-health-of-every-project-the-host-knows.md.

## Entries

### 2026-09-26T06:34:49Z agent-S-0122
Criterion 2's positive path cannot happen as flai stands. It says Serve registers an unserved project below an import folder, and the switcher gains it.

**Why.** Since S-0120, flai serve looks below the import folders every scan and serves every project there that has a key. So a project listed as "below the import folders, not served" always has one of three reasons (`serve.Place`):

| Reason | What Serve does |
|--------|-----------------|
| its `system-flow.yaml` has no key | `flai serve project add` refuses it, for the same reason |
| its key is served already, for another folder | refused by S-0121's key check |
| no dashboard is known to serve it for | nothing can reach the settings page then |

So the page shows Serve, and flai answers with the refusal and its reason. The page shows that answer. A project that is served shows up in the switcher without a reload after its next scan, but that is S-0120's work, not Serve's.

Serve would do something real once a project could be taken off below an import folder. That is the "list of roots the operator removed" that TH-0014 left for later. Remove would add the root to that list, and Serve would take it off again.

**A. Recommended.** Keep Serve as built. In the browser trial, show its refusal with flai's reason, and show that a project added below an import folder joins the Projects list and the switcher without a reload. Reword criterion 2 to say that. Put the removed-roots list, and a Serve that undoes Remove, in a new backlog story under E-0003.

B. Build the removed-roots list in this story. Remove on an import-served project adds its root to the list; Serve takes it off, and the switcher gains the project. It adds work in flai serve (the registry, Place, and the CLI) and about one more session.

C. Drop Serve from the page, and reword criterion 2 to show the unserved projects and why only.

I am running the rest of the trial now (health, Remove, refusal with settings off, the switcher), and it does not depend on this answer.

### 2026-09-26T06:54:55Z agent-S-0122
The browser trial is done, and everything but criterion 2 is ticked. The results are in S-0122's Notes.

- Serve on gamma (no key) returned flai's refusal: "not served: its system-flow.yaml has no key".
- A project committed below the import folder (eta) joined the switcher without a reload after flai serve's next scan.

So A needs only the criterion reworded and a backlog story. B needs new work in flai serve. C removes Serve from the page and the docs. S-0122 goes to review once you choose.

### 2026-09-26T06:58:43Z alex
A

### 2026-09-26T06:59:30Z agent-S-0122
Resolved: A: criterion 2 reworded to what flai does today and checked, the users guide and dashboard design say Serve shows flai's refusal, and backlog story S-0123 holds the list of removed roots and a Serve that undoes Remove.
