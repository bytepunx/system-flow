---
id: TH-0074
title: "S-0198: which issues the review page offers, and the check's default threshold"
anchor:
  path: wip/kanban/stories/S-0198-give-the-operator-the-option-to-have-all-issues-turned-into-stories.md
  item: S-0198
status: resolved
participants: [agent-S-0198, alex]
created: 2026-10-03T01:26:54Z
updated: 2026-10-03T02:11:09Z
---

# TH-0074 S-0198: which issues the review page offers, and the check's default threshold

On wip/kanban/stories/S-0198-give-the-operator-the-option-to-have-all-issues-turned-into-stories.md.

## Entries

### 2026-10-03T01:26:54Z agent-S-0198
Two choices change the work. I am going on with everything else while they are open.

1. Which open issues does a story's review page offer, each with a checkbox?

   Recommendation: every open issue no open story links, the ones this story recorded or bumped listed first and checked by default, the rest unchecked. A "Create stories" button beside the list makes one backlog story per checked issue, apart from Accept. The alternative is only the issues the story recorded or bumped.

2. What does `flai check` do when the project sets no threshold?

   Recommendation: `issues.story_after` in `system-flow.yaml`, a Go duration, defaulting to `168h` (7 days) when unset, with `0` turning the check off. The rule is `issues.no-story`, a warning. The cost: until each of today's open issues older than that has an open story, every close-out here fails `flai check --strict`. Today that is only I-0007, from 2026-09-15. Making one story for it in the backlog clears it. The alternative is no default, so the check is off until a project sets the key.

Two choices I made without asking, which you can change:

- A story made from an issue copies the issue's `## Remediation` text as its recommended solution, and links the issue in its goal.
- Nature: `defect` and `blocker` issues get `remediation`; `efficiency` and `impression` issues get `improvement`.

### 2026-10-03T01:46:03Z alex
1 - every open issue that this story touched should show first in the list with a checked checkbox; all issues that don't have a story for their remediation or improvement should appear in the list, unchecked. If any issues are checked, the `Accept` button text changes to `Accept and Create Stories` in which the action should perform both tasks (create stories only takes place if the accept succeeds, otherwise, do not proceed to the create stories step). 2 - proceed with the recommendation AND create a story for issue 7 so that flai check does not fail.

### 2026-10-03T02:11:09Z agent-S-0198
Resolved: Built as answered: T-0734 (review page, Accept and Create Stories), T-0732 (issues.story_after 168h, issues.no-story), and S-0243 made for I-0007. Recorded in the S-0198 narrative's Decisions.
