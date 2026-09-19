---
title: flaiover dashboard
updated: 2026-09-19
status: draft
---

# flaiover

Run it with `flai dashboard` and open <http://localhost:4242>. The design is in [design/system/flaiover-dashboard.md](../../design/system/flaiover-dashboard.md); the full guide arrives with story S-0019.

## Theme

flaiover uses the brand palette in a light and a dark theme. It follows your system preference by default; the button at the right of the navigation cycles system, light, and dark, and the choice is remembered per browser. Both themes are checked for readable contrast, and the charts use palettes validated for colour-vision deficiency in each theme.

## Documentation explorer

Docs shows every markdown file under `design/`, `docs/`, and `wip/` in a collapsible tree. A document renders with its Mermaid diagrams, highlighted code, task-list checkboxes, and heading anchors; links between documents open in the explorer. The front matter is shown in a panel above the text. ADRs lists the architecture decisions with status, date, and which decisions supersede which.

## Agents on other machines

The dashboard also serves the project's MCP tools over HTTP at `/mcp`, so an agent that is not on this machine can read the inbox, move items, and reply to threads with the project token. What it does shows up here like any other agent's work. Setting it up is in [the flai guide](flai.md) and, for what to put in front of it, the operator guide.

## Inbox and activity

Inbox is the list of things that need you, and the number beside its link in the navigation is how many there are. It holds five kinds of entry, each a link to where you deal with it: stories in review, which open their review page; threads where someone other than you wrote last; open questions agents left in their narratives; blocked items, with the reason; and overlapping touches, where two stories in progress say they change the same files. An entry leaves the list when its cause does: you answer the thread, accept the story, the item is unblocked. The list and the count refresh by themselves when files change.

You can ask the browser to tell you when something new arrives: tick "Desktop notification" on the inbox page and allow it when the browser asks. It is off until you turn it on, it is remembered per browser, and it only fires for entries that appear while a dashboard tab is open, never for what was already there.

Activity shows who is working on what: one card per story with an open narrative, with the agent and session that last wrote it, how long ago, the story's state, the task in progress, whether anything in it is blocked, and the last line of its log. Nothing is recorded to produce this. An agent that stops writing does not disappear; its card simply grows older.

## Reviewing a story

A story in review has a review page: open it from the card in the review column, or from "Review this story" on the item page. It puts what you need to decide in one place.

- The acceptance criteria, with which are ticked and how many.
- The narrative's current state and next steps, as the agent left them, with a link to the whole narrative.
- What the story's branch changes against main: the files, with lines added and removed, and the hunks when you open a file. Large patches are cut and marked; the rest is a `git diff` away. A story that was not worked on a branch says so instead.
- What accepting will do: the branch that is merged, the release that is cut, anything that blocks acceptance from here, and any uncommitted files outside `wip/`, which you choose to include or deal with first.
- The threads on the story, where you can ask before deciding.

Accept runs the acceptance as you, the project's owner, and shows each step as it completes: the branch merged, the story moved to done, the archive, the commit, the tags. Nothing is pushed from the dashboard; the page says so, and you push from a shell. If flai refuses or fails, for example a rebase that stops on a conflict or an unticked criterion, the page shows flai's message word for word, and the story is still in review.

Send back asks why, in the page, and moves the story to in-progress with your reason recorded in its notes, where the agent reads it.

## Editing documents

Every document page and every work item page has an Edit link when the dashboard can write. The editor shows the markdown on the left and, on the right, a preview drawn by the same renderer as the explorer, diagrams and code highlighting included.

What you can change depends on the file, and flai decides it, not the dashboard. Design and docs files are yours entirely: body and front matter, with the front matter checked when you save. For work items, narratives, and the board, flai owns the front matter, because it is the item's state; it is shown read-only and you edit the body. Move, block, and retitle items from the board or with `flai`. Generated files (`wip/agents/index.md`, `design/issues/summary.md`), threads, issues, and anything in the archive are not editable here, and the page says why.

Saving does three things. flai checks the repository with your change in place, and if the change introduces any finding, the save is refused, the findings are listed, the file is left as it was, and your text stays in the editor. If someone else changed the document after you opened it, an agent or a colleague, you get a conflict instead: the page shows what is there now against what you are saving, and you choose to load the current version, which discards your edits, or to save yours over it. Otherwise the file is saved and committed on its own, with you as the author, the line you typed under "What changed" as the subject, and a trailer naming the dashboard. For design and docs files the `updated` date is set to today unless you set it yourself. A project can turn the commit off with `dashboard.autocommit: false` in `system-flow.yaml`; the edit is then saved and left for you to commit.

Commits made here are not pushed; push from a shell, as with a story accepted from the board.

When a story in progress or in review says it touches the document, the editor warns you, as the document page does, and asks you to tick a box before the first save: your edit lands on `main`, and that story's branch will meet it at its next sync. Click anywhere in the body and "Open a thread on" names the heading you are under, so a question can be attached to the section it is about.

The editor asks before you leave with unsaved changes. Creating, renaming, and deleting documents is not something it does.

## Board

Board shows one column per state with the WIP count against the limit from `wip/kanban/board.md`. Cards are stories by default; a toggle adds epics and tasks. Each card shows the ID, title, nature, how long it has sat in its column, and a flag when it is blocked. The nature, the type, the blocked flag, and the parent sit in one row under a thin line below the title, all in one size. In that row's right corner a story card shows the ID of its epic, and a task card the ID of its story; hover over it, or use a screen reader, for the parent's title. Cards are colour coded: the background is a pastel for the item's nature (feature, improvement, remediation, research, experiment) and the stripe on the left edge is a colour for its type (epic, story, task); the legend above the columns names each colour. The colours only repeat what the card says in words, so nothing depends on telling them apart. The same colours follow an item elsewhere: the header of its page, its hits in Search (documents stay plain), and the table under the cycle time charts, where each nature's dots and bars are a stronger shade of the same hue. A blocked card also has a red outline, and the card you are dragging turns faded with a dashed edge. The whole card is one link to the item, so the parent ID is not a link of its own. Drag a card to another column to move it: the same rules `flai move` enforces apply, and a refused move shows the rule. Dropping a story on done accepts it, so the dashboard first shows what acceptance will do, the branch to be merged and the release it will cut, and only proceeds when you confirm; cancelling leaves the card in review. If files outside `wip/` are uncommitted, the confirmation lists them: acceptance refuses them by default so that its commit holds only acceptance, so either cancel and commit or stash them first, or tick the box to include them, which is the dashboard's form of `flai accept --yes`. The accept button waits for that choice. When acceptance cannot run from the dashboard, the confirmation lists why and what to do, and its accept button stays disabled: the usual case is a dashboard that sees the repository at a different path than your machine does, so git cannot open the story's worktree; accept from a shell with `flai accept S-nnnn`, or stop the dashboard and start it with `flai dashboard`, which mounts the repository at its own path. If the dashboard cannot push, it says the story was accepted locally and asks you to push from a shell. Every write is made by the bundled `flai`, so the files change exactly as they would from the terminal and appear in `git status` for you to commit.

Click a card for the item page: front matter, the rendered body, children, the transition history, blocked intervals, a link to the narrative, and buttons for the allowed moves, block, unblock, and adding a narrative log entry.

## Charts

Charts plots the flow metrics `flai stats` computes, so the numbers are the same in both places. Pick a window, an item type, and where it applies an epic. Cycle time shows one point per completed item with the p50 and p85 lines; burn-up shows scope against done; cumulative flow shows how many items sit in each state each day; time in state shows where each completed item spent its time and the share across all of them; throughput counts completions per week by nature; aging lists in-progress work against the p85 line; estimates compares estimated with actual hours. Every chart has a table view and follows the light or dark theme.

## Search

Search covers `design/` and `wip/` by default and `docs/` when you tick the box. Type an item ID, a title, or words from the body; each result shows where it is, its status, and a snippet around the match. Results open in the explorer.
