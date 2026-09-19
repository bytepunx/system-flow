---
id: S-0059
type: story
nature: feature
title: Operators create epics and stories from the board
status: backlog
parent: E-0006
owner: alex
created: 2026-09-19T05:34:09Z
updated: 2026-09-19T05:34:09Z
transitions: []
tags: [dashboard, cli]
touches: [flaiover, flai/cmd]
---

# S-0059 Operators create epics and stories from the board

## Goal
An operator creates a new epic or story from the board by writing what it is for, in markdown, and nothing else. The ID, file name, front matter, timestamps, the link in the parent's list, and the place on the board are supplied behind the scenes, so the file that results is exactly what `flai epic new` or `flai story new` would have made, without the operator creating a file by hand or running a command.

## Acceptance criteria
- [ ] The board has a "new" action that opens a form for an epic or a story. For a story the operator picks the parent epic from the open epics and the nature from the five natures, each shown with what it means; for an epic, the nature. Title is required. Tags and touches are optional and offered, not demanded
- [ ] The operator writes the content in markdown in the sections the item template already has: outcome for an epic; goal, acceptance criteria, and notes for a story. The form starts from the template's headings, shows the explorer's preview beside the text, and checkboxes typed as `- [ ]` become the acceptance criteria. No front matter is shown for editing, because none of it is the operator's to write
- [ ] Creation goes through flai, not through file writes in the dashboard (ADR-0016, ADR-0023): one call creates the item from the project's template with the operator's content as its body, gives it the next free ID, links it into its parent, and records the manifest's owner as the owner. If `flai story new` and `flai epic new` cannot take a body today, they gain a way to (for example content on standard input), so the same is possible from a script
- [ ] The result is validated before it is kept: `flai check` runs with the new item in place, and anything it would report is shown in the form with the text kept, and nothing is left behind, the way a refused document save behaves (S-0040)
- [ ] A created item is committed on its own unless the project sets `dashboard.autocommit: false`, with the operator's git identity as author and the dashboard's trailer, like a document saved from the editor. Nothing is pushed
- [ ] After creating, the operator lands on the new item with its ID shown, the board shows it in backlog without a reload, and a story can be moved to ready from there when its goal and criteria are written (ADR-0021: no tasks needed)
- [ ] A read-only dashboard (no flai) shows no "new" action. Tasks are not created here: the agent that pulls a story writes them
- [ ] Tests: the flai side with real git (body, parent link, ID, check refusal leaving nothing behind, commit), the endpoint, and the form as a component; `design/system/flaiover-dashboard.md`, `flai-cli.md`, and `docs/users` updated

## Tasks

## Notes
Raised by the operator on 2026-09-19: "On the board, operators should be able to create new epics or stories ... provide the boilerplate values behind the scenes and allow the operator to focus on providing new guidance in markdown that will create the underlying documents in the correct format and with the correct front-matter without needing to hand-create files or use the flai CLI commands since the goal is to eventually drive all project activity via the dashboard."

What exists. `flai epic new "<title>"` and `flai story new "<title>" --epic --nature --owner --tag --touches` create the file from the template's item bodies with the next free ID and link it into the parent's list; the body they write has empty `## Goal`, a blank checkbox, `## Tasks`, and `## Notes`. The dashboard has no create endpoint: `flaiover/src/routes/api/items/+server.ts` only lists. Since S-0040 the dashboard can save a document's body through `flai doc save`, with a check before it is kept and a path-limited commit; creation can reuse that machinery rather than invent another.

Two ways to get the operator's content into the new item, to choose when refining: `flai story new` accepts the body on standard input and does create, check, and commit as one step; or the dashboard calls `flai story new` and then `flai doc save` on the new file. The first is one operation that either happens or does not, which is what the fourth criterion needs; the second leaves an empty item behind when the save is refused.

The template decides the body's sections (`template/root` item bodies, which a project may have changed), so the form should take its headings from the template flai renders, not hard-code them.

Tasks are out of scope on purpose (ADR-0021). Editing an existing item's body is already possible (S-0040); retitling, changing nature, tags, or parent from the dashboard is not, and is not part of this story.

Sits with S-0060 (ADRs from the dashboard): both are "the operator writes markdown, flai supplies the rest", and whichever is built second should reuse the first's form and endpoint shape.
