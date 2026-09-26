---
title: flaiover dashboard
updated: 2026-09-26
status: active
---

# flaiover

flaiover is the dashboard for a system-flow project: the board, the work items, the documentation, the decisions, and the flow charts, in a browser. It reads and changes nothing by itself. Everything it shows and everything it changes it asks of `flai serve` on your machine, so a move made on the board is the same move `flai move` makes, committed with your own git identity. The design is in [design/system/flaiover-dashboard.md](../../design/system/flaiover-dashboard.md).

Running it, securing it, and the settings behind it are in the [operators guide](../operators/index.md#running-the-dashboard). Read its warning before you make the dashboard reachable from anywhere but your own machine: whoever holds its token can change your project as you.

## Starting it and logging in

In a project, run `flai dashboard`. It starts the dashboard, and `flai serve` behind it, if they are not running, registers the project with them, and prints a login link. Open the link: <http://localhost:4242/login#token=...>. The token in it never leaves the browser; the page exchanges it for a session cookie and removes it from the address and the history. `flai dashboard token` prints the link again.

**Login** is the page you land on without a session. Paste the token, or open the login link. When your session ends, because the token was rotated or the cookie was lost, any page sends you back here and returns you to where you were once you log in. If the token is right but the browser did not keep the cookie, the page says so and why, rather than asking for the token again; the operators guide explains the case of a tunnel that reports HTTPS.

## The header

The header is the same on every page.

- **Project switcher**, at the left: the project you are looking at, and the others this dashboard serves. See [More than one project](#more-than-one-project).
- **Navigation**: Overview, Board, Inbox (with the number of things that need you), Activity, Charts, Docs, ADRs, Search, Host, and Settings. Each is a section below.
- **Host flai badge**, at the right: whether a flai on the host has the dashboard connected. See [Host flai](#host-flai).
- **Theme button**: see [Theme](#theme).

## Host flai

At the right of the header the dashboard says whether a flai on the host has it connected: "host flai: connected" (hover for the version and since when), or "host flai: not connected" on a yellow ground. `flai dashboard` starts that process, `flai serve`, for you; `flai serve status` on the host says why it is not connected. Every page is read through it: the dashboard asks flai for the board, the items, the threads, the documents, the ADRs, the inbox, activity, and search, and flai tells it when a file changes, so a move made in a terminal shows on an open board in about half a second. Without it pages have nothing to show, and a yellow banner under the header says so on every page, with the commands that bring it back; the page recovers by itself when flai returns.

Writes go the same way: a move, a new story, a saved document, an acceptance are asked of flai on the host, which does them with your own git identity and, for commits made for you, the dashboard's trailer. The dashboard's image holds no flai and no git. Throughout this guide, "when the dashboard can write" means while flai on the host has it connected. If the flai on the host is older than the dashboard, the banner says what it lacks and asks you to upgrade flai. If the connection is lost in the middle of an acceptance, the page says the acceptance may have completed on the host: reload the board before trying again.

## More than one project

One dashboard, and one `flai serve` on the host, can serve several projects. The header always has a project switcher showing the project you are looking at, even when there is only one. It lists the projects, then, under "Not imported yet", the repositories on the host that can be imported (see below). Pick another project and every screen changes to it where you are, without reloading. The board you are on shows the other project's board and follows its changes from then on, and the inbox count and the host flai badge follow too. Your choice is remembered in this browser and shown in the page's address and in the header's links, so a bookmark or a copied link opens on the same project. Nobody sees your inbox, board, or documents mixed with another project's; picking one shows that one and nothing else.

To add a project, run `flai dashboard` in it on the host: when the dashboard is already running, that only registers the project with the running `flai serve`, which connects it within a second or two, and it joins the switcher by itself (the list is looked at again every 30 seconds, and whenever you open the switcher). The project needs a `key` in its `system-flow.yaml`; `flai check` says how to add one. If a project you expect is missing, its flai has not connected here yet: `flai serve status` on the host says why. A project `flai serve` serves that is not connected stays in the switcher as "(served, not connected)", with the reason when you hover over it, and a link beside the switcher ("1 not connected: why?") opens [Settings](#settings), where each project's state says why. Projects can also be served and removed there, without a shell.

A repository that is not a system-flow project yet can be imported from the board, once the operator has named the folder it is in on the host (`flai serve import add ~/git`, in the operators guide). It appears in the switcher and the project list as "(not imported)". Picking it asks whether to import it and says which tests will run. **Not now** tells you to select a different project; nothing is remembered, so picking the repository again asks again. **Import** adds the standard's folders and files to the repository (nothing that is there is overwritten), runs its tests, and commits the import if they pass. The page then shows the project it has become. If a test fails, nothing is committed: the page shows what failed, the imported files wait in the repository to be fixed and committed, and **Open the project** shows it, since it is a project either way. A repository with uncommitted changes is refused until they are committed or stashed.

## Overview

The front page. For one project it shows the project's name and description from `system-flow.yaml`, the template version it follows, how many items are active and how many archived, and a table counting epics, stories, and tasks in each state from backlog to done. It is the quickest answer to "how much is there, and where is it".

While no project is chosen and more than one is known, the front page shows a list of every project instead, with a filter box, each with whether it is connected, how many stories are in review, how many threads are awaiting you, and whether an agent is attending, so you can tell at a glance which one needs you. Pick one to open it.

## Board

Board shows one column per state with the WIP count against the limit from `wip/kanban/board.md`. Cards are stories by default; a toggle adds epics and tasks. Each card shows the ID, title, nature, how long it has sat in its column, and a flag when it is blocked. The nature, the type, the blocked flag, and the parent sit in one row under a thin line below the title, all in one size. In that row's right corner a story card shows the ID of its epic, and a task card the ID of its story; hover over it, or use a screen reader, for the parent's title. The whole card is one link to the item, so the parent ID is not a link of its own.

**Colours.** Cards are colour coded: the background is a pastel for the item's nature (feature, improvement, remediation, research, experiment) and the stripe on the left edge is a colour for its type (epic, story, task); the legend above the columns names each colour. The colours only repeat what the card says in words, so nothing depends on telling them apart. The same colours follow an item elsewhere: the header of its page, its hits in Search (documents stay plain), and the table under the cycle time charts, where each nature's dots and bars are a stronger shade of the same hue. A blocked card also has a red outline, and the card you are dragging turns faded with a dashed edge.

**Agents.** When flai on the host starts agents for ready stories (the operator turns that on with `flai serve enable agent`), a notice above the columns says what it last did: the agent it started, for which story, by which command, or why it could not, or why a ready story is waiting. A story's card has a dot beside its ID while its agent is at it: green while it works, yellow while it waits for you, red when it failed. Hover over the dot, or use a screen reader, for what it is doing and why. An agent that finished its story shows no dot. A yellow dot means the agent asked you something on the story, or the story is blocked: answer it on the story's page or from the inbox, and the agent goes on, started again if it had ended. The story's page says the same at the top of the side column, with the harness, the model, when it started, and for a red dot where the agent left the story and where its output is on the host. Moving a story with a red dot back to ready starts another agent, and so does changing its agent while it is in ready. For a story in ready or in progress, **Retry**, at the top right of the agent section, starts a new one once the operator has turned on the `agent` host action. Retry disappears when you press it. For a story in ready while the in-progress limit is full, the new agent is queued: the dot turns yellow, the page says it waits for room, and flai starts it as soon as there is some. If flai refuses, the page says why and Retry comes back. A story in ready that has no agent yet says why flai has not started one, and **Start agent** starts it at once, under the same conditions.

**Pull order.** The backlog and ready columns are in pull order, the order agents take work in: drag a story up or down within either column to change it, and a line shows where it will land. Stories you have never placed come last, by ID. Without a mouse, focus a card and press Alt+↑ or Alt+↓, or use the ↑ and ↓ buttons that appear on the card when it is hovered or focused. The other columns have no order, so dragging within them does nothing.

**Moving.** Drag a card to another column to move it: the same rules `flai move` enforces apply, and a refused move shows the rule. Dropping a card on cancelled, or choosing cancelled on an item's page, opens a confirmation before anything happens: cancelling an epic also cancels every open story under it and their tasks, cancelling a story its open tasks, and the confirmation lists them, points out any story in review (its work stays on its branch, unmerged, so accept it first if you want it), names the branches and worktrees that will be left for you, and asks for the reason, which is written into every item cancelled. Cancelled is final; "Keep it" closes the confirmation and changes nothing.

**Accepting.** Dropping a story on done accepts it, so the dashboard first shows what acceptance will do — the branch to be merged, anything that blocks it — and only proceeds when you confirm. Accepting merges, archives, and commits; it cuts no release by itself (S-0087). An experiment story cannot be accepted: the confirmation gives the reason and its accept button stays disabled; cancelling leaves the card in review. If files outside `wip/` are uncommitted, the confirmation lists them: acceptance refuses them by default so that its commit holds only acceptance, so either cancel and commit or stash them first, or tick the box to include them, which is the dashboard's form of `flai accept --yes`. The accept button waits for that choice. When acceptance cannot run from the dashboard, the confirmation lists why and what to do, and its accept button stays disabled; accept from a shell with `flai accept S-nnnn` instead. An acceptance is made on the host by flai, committed there. The review page (below) is the fuller way to accept.

**Pushing and publishing.** Whether an acceptance has reached the remote is the operator's choice, made on the host: by default it has not, and the board keeps a notice that says what was accepted and not pushed and the command to run on the host (`flai push --pending`), which a running agent session does by itself. Pushing tags whatever has accumulated and is unreleased first, so a release is never a separate step someone has to remember: the notice's own **Push now** button, when the operator has enabled pushing from the board (`flai serve enable push`), reports what it pushed with the tags it just created, and clears once the push has happened from this clone (a push made from another clone is not seen until someone fetches here); a push that is refused, because the remote has moved, shows its reason and the command to run by hand. Whether a story has been **published** — released, with a version bump and a tag — is shown in the done column itself, since a push and an accepted story do not always land together: a card says **published** or **waiting to publish**, and a banner at the top of the column lists what publishing now would release, by component, with the story IDs it bundles (several small stories against the same component release together, at the highest bump among them, not one release each). The banner's **Publish** button (same host action, same enable command) does that same computing and tagging explicitly, ahead of or instead of a push; a publish that fails partway is safe to run again, it does not redo what already succeeded.

### Creating work

When the dashboard can write, the board has a "+ new" action, which opens the new item page. Choose epic or story, for a story an epic to belong to or "No epic" (not every story fits an active one, S-0092), and the nature (each is shown with what it means and what bump it contributes once published). Give it a title, optionally tags (a component's name decides which component it delivers to) and the paths it touches. For a story you may also name its agent: the harness, the model, and options one `key=value` a line. The project's default agent (`flai agent set`) shows in grey in those fields and is what the story gets for anything you leave empty. Then write what it is for in markdown: the form starts from your project's template sections and shows the preview beside the text; acceptance criteria are `- [ ]` lines. You write no front matter. The ID, the file name, the front matter, the link in the epic's list, and the commit (your git identity, with a flaiover trailer) are made by `flai`, so the file is exactly what `flai story new` would have made. If `flai check` would report anything with the new item in place, nothing is created, the findings are shown, and your text stays in the form. After creating you land on the item's page; the board shows it in backlog without a reload, and a story with its goal and criteria written can be moved to ready from there. Tasks are not created here: the agent that pulls a story writes them. Nothing is pushed.

### The item page

Click a card for the item page. The main column has the ID and title, the type and nature, the status, a BLOCKED flag, the parent, and the rendered body, with the threads on the item below it. A thread's entries render as markdown, as the body does, wherever threads are shown. Your entries sit on the right on a blue ground and an agent's on the left on a neutral one, so who said what reads at a glance; you are the project's `owner` in `system-flow.yaml`. When the dashboard can write, buttons above the body make the moves the item's state allows, block it with a reason or unblock it, and open **edit…** for a story or an epic (see [Editing a story or an epic](#editing-a-story-or-an-epic)). On a story the acceptance criteria are live checkboxes: tick one where you read it. A story in review has a **Review this story** link; a done story that was accepted and not pushed says so.

The side column has the story's agent, when it has one; the history, every transition with when and by whom, and the intervals it was blocked with their reasons; the children, with their states; the file, which opens in the explorer; the narrative, with a box to add a log entry to it; and the owner, estimate, tags, touches, and agent from the front matter.

### Reviewing a story

A story in review has a review page: open it from the card in the review column, or from **Review this story** on the item page. It puts what you need to decide in one place.

- The acceptance criteria, with which are ticked and how many.
- The narrative's current state and next steps, as the agent left them, with a link to the whole narrative.
- What the story's branch changes against main: the files, with lines added and removed, and the hunks when you open a file. Large patches are cut and marked; the rest is a `git diff` away. A story that was not worked on a branch says so instead.
- What accepting will do: the branch that is merged, anything that blocks acceptance from here, and any uncommitted files outside `wip/`, which you choose to include or deal with first. Accepting cuts no release by itself (S-0087); publishing what has accumulated is a step of its own, from the done column.
- The threads on the story, where you can ask before deciding.

**Accept** runs the acceptance as you, the project's owner, and shows each step as it completes: the branch merged, the story moved to done, the archive, the commit. Nothing is tagged or pushed by acceptance; the page says so. If flai refuses or fails, for example a rebase that stops on a conflict or an unticked criterion, the page shows flai's message word for word, and the story is still in review.

**Send back** asks why, in the page, and moves the story to in-progress with your reason recorded in its notes, where the agent reads it.

## Inbox

Inbox is the list of things that need you, and the number beside its link in the navigation is how many there are. It holds five kinds of entry, each a link to where you deal with it: stories in review, which open their review page; threads where someone other than you wrote last; open questions agents left in their narratives; blocked items, with the reason; and overlapping touches, where two stories in progress say they change the same files. An entry leaves the list when its cause does: you answer the thread, accept the story, the item is unblocked. A hand-written open question (one with no thread of its own) can be answered right there, in the inbox, when the dashboard can write: type the answer and submit it, and it moves to the narrative's Decisions and drops off the list. The list and the count refresh by themselves when files change.

You can ask the browser to tell you when something new arrives: tick "Desktop notification" on the inbox page and allow it when the browser asks. It is off until you turn it on, it is remembered per browser, and it only fires for entries that appear while a dashboard tab is open, never for what was already there.

## Activity

Activity shows who is working on what: one card per story with an open narrative, with the agent and session that last wrote it, how long ago, the story's state, the task in progress, whether anything in it is blocked, and the last line of its log. Nothing is recorded to produce this. An agent that stops writing does not disappear; its card simply grows older.

## Charts

Charts plots the flow metrics `flai stats` computes, so the numbers are the same in both places. Pick a window, an item type, and where it applies an epic. Every chart has a table view under it and follows the light or dark theme.

| Chart | Shows |
|-------|-------|
| Cycle time | One point per completed item, with the p50 and p85 lines |
| Burn-up | Scope against done over time |
| Cumulative flow | How many items sit in each state each day |
| Time in state | Where each completed item spent its time, and the share across all of them |
| Throughput | Completions per week, by nature |
| Aging work in progress | In-progress work against the p85 cycle time line |
| Estimate versus actual | Estimated against actual hours |

## Docs

Docs shows every markdown file under `design/`, `docs/`, and `wip/` in a collapsible tree. A document renders with its Mermaid diagrams, highlighted code, task-list checkboxes, and heading anchors; links between documents open in the explorer. The front matter is shown in a panel above the text, and the threads on the document below it. When a story in progress or in review says it touches the document, the page names the story: an edit here may collide with that story's branch. When the dashboard can write, an **Edit** link opens the editor.

### Editing documents

Every document page and every work item page has an Edit link when the dashboard can write. The editor shows the markdown on the left and, on the right, a preview drawn by the same renderer as the explorer, diagrams and code highlighting included.

What you can change depends on the file, and flai decides it, not the dashboard. Design and docs files are yours entirely: body and front matter, with the front matter checked when you save. For work items, narratives, and the board, flai owns the front matter, because it is the item's state; it is shown read-only and you edit the body. Move and block items from the board or with `flai`; change what a story or an epic says about itself from its own page, as described next. Generated files (`wip/agents/index.md`, `design/issues/summary.md`), threads, issues, and anything in the archive are not editable here, and the page says why.

Saving does three things. flai checks the repository with your change in place, and if the change introduces any finding, the save is refused, the findings are listed, the file is left as it was, and your text stays in the editor. If someone else changed the document after you opened it, an agent or a colleague, you get a conflict instead: the page shows what is there now against what you are saving, and you choose to load the current version, which discards your edits, or to save yours over it. Otherwise the file is saved and committed on its own, with you as the author, the line you typed under "What changed" as the subject, and a trailer naming the dashboard. For design and docs files the `updated` date is set to today unless you set it yourself. A project can turn the commit off with `dashboard.autocommit: false` in `system-flow.yaml`; the edit is then saved and left for you to commit.

When a story in progress or in review says it touches the document, the editor warns you, as the document page does, and asks you to tick a box before the first save: your edit lands on `main`, and that story's branch will meet it at its next sync. Click anywhere in the body and "Open a thread on" names the heading you are under, so a question can be attached to the section it is about.

The editor asks before you leave with unsaved changes. Creating, renaming, and deleting documents is not something it does. Commits made here are not pushed; push from a shell, as with a story accepted from the board.

### Editing a story or an epic

A story's or an epic's page has an **edit…** button while the item is open (not done, cancelled, or archived). It turns the page into a form: the title, the nature, the tags, and for a story what it touches, its parent epic, and its agent (the harness, the model, and options; what you save replaces the story's agent, and emptying all three removes it), with the body below as Markdown and a preview. The line at the top names what is not yours to change there, because it is the item's state and flai's: the ID, the type, the status (move the card instead), the owner, and the dates.

Save sends only what you changed, and flai on the host does the rest in one step. A new title is kept in step everywhere it appears: the file's name, the heading, the line in the parent's list, the story's narrative, and links to the old file name in other documents. A new parent must be an open epic; the story leaves the old epic's list and joins the new one. The repository is checked with your change in place, exactly as for a document: if the check finds anything the change introduces, say a story in progress left without acceptance criteria, nothing is changed, the findings are listed, and your text stays in the form. If an agent or someone else changed the item after you opened it, you get the same choice as for a document: load the current version, or save yours over it, and only the fields you changed are written over. Everything is committed in one commit, as you, with the dashboard's trailer.

An agent working on the story is told. Its `inbox` reports the edit with what changed (title, nature, tags, touches, parent, goal, criteria, notes, or body), and one that is waiting hears within a second, so it reads the story again before it goes on.

Tasks are not edited here: they are the agent's to write. A task's page says so, and its body can still be opened as a document.

## ADRs

ADRs lists the architecture decisions with status, date, and which decisions supersede which; each opens in the explorer. When the dashboard can write, "+ new ADR" records a decision from there.

### Recording a decision

"+ new ADR" opens a form: a title (the decision, as a sentence), whether it is proposed or accepted, the existing decisions it supersedes or refines, chosen from the list, and the markdown, which starts from the sections of your project's `design/adrs/0000-template.md` with the preview beside it. You write no front matter. The number (one more than the highest ADR file present), the file name, the front matter, the date, the row in `design/adrs/README.md`, `superseded_by` on a decision it replaces, and the commit (your git identity, with a flaiover trailer) are made by `flai adr new`. If `flai check` would report anything with the new record in place, nothing is created and your text stays in the form. Afterwards you land on the new ADR in the explorer, and the list shows it without a reload, with what it supersedes and refines.

A proposed ADR is a draft: it can still be edited, and its row on the ADRs page has an "accept" button, which sets it to accepted with today's date. An accepted ADR is immutable, here as everywhere: the editor will not open its body, and the way to change a decision is to record a new one that supersedes it. Nothing is pushed.

## Search

Search covers `design/` and `wip/` by default and `docs/` when you tick the box. Type an item ID, a title, or words from the body; each result shows where it is, its status, and a snippet around the match. It runs in flai on the host and matches item IDs, titles, tags, headings, and bodies, by prefix and with small typos forgiven; an item ID such as `S-0052` finds its item first. Results open in the explorer.

## Host

The host flai badge in the header is a link to Host (also in the nav). It has two parts.

**Dashboard** shows what image and version the dashboard container is running, and, once the operator has enabled it (`flai serve enable dashboard`), Restart, Upgrade, and Stop. **Check for updates** always works and changes nothing. Restart and a successful upgrade stop the very container answering the page, so the page expects the connection to drop and shows "Reconnecting…" rather than an error, then says what came back once it does. An upgrade never touches the running container until the new image has proven itself healthy: if it does not, the page says so and nothing changed.

**flai host** shows the flai host process (S-0106): its version and pid, and a row each for `serve` and the MCP servers. Each row gives the state, the version, and the restart count, and the MCP row lists the projects served. **Check for upgrade** always works and changes nothing. Once the operator has enabled it with `flai serve enable host`, each row has Start, Stop, and Restart, and one **Upgrade** installs the newest flai and has the host start serve and the MCP servers again from it. The dashboard reaches the host through serve, so it asks before stopping serve: only `flai host start serve` in a shell on the host brings serve back. A serve restart and an upgrade drop the page's connection, so the page shows "Reconnecting…" and then says what version came back. With no host running, the area says to start one with `flai host start`.

## Settings

Settings shows how flai on the host is set up for this project:

- which host actions are on;
- the default agent new stories get;
- the command and harnesses that start agents;
- the checks run on a story in review;
- the folders the board offers repositories from;
- the projects `flai serve` serves, and those below the import folders it does not;
- the MCP server.

It changes nothing until the operator runs `flai serve enable settings` on the host. Each section says so and names the command. After that, this project's own settings can be changed here: its host actions, its default agent, and its MCP token. The settings the host keeps for every project also need `flai serve enable settings --all-projects`.

Projects lists every project `flai serve` serves, with its key, name, and folder, and whether it is connected and since when, the last error, or why it cannot be served. A project below an import folder is served as soon as `flai serve` next looks, and joins the switcher without a reload. One it does not serve is listed under "Below the import folders or the folder flai serve was started in, not served" with why: you removed it, its `system-flow.yaml` has no key, or its key is served already for another folder. **Remove** asks first, then stops serving the project, as `flai serve project remove` does, and the switcher drops it without a reload; none of its files is touched. A registered project is unregistered, and the page says how to serve it again on the host. A project served because it is below an import folder, or below the folder `flai serve` was started in, says which folder serves it; removing it puts it on `flai serve`'s list of removed projects, and it moves to the not-served list as removed (S-0123). **Serve** on a removed project takes it off that list, as `flai serve project add` does, and the switcher gains it without a reload. On any other project not served, Serve asks `flai serve` to register it and shows the answer; for no key or a key served already that is a refusal saying what to fix. Removing the folder under Import folders stops every project below it at once. Serve and Remove need `flai serve enable settings` for the project served or removed, and each says so, with the command, where it is off.

Commands are written one argument a line, and are run exactly as written, with no shell in between. Each section says whether the change was saved, or why flai refused it. Rotating the dashboard token keeps you logged in and shows the new login link once. Everyone else is logged out.

## Theme

flaiover uses the brand palette in a light and a dark theme. It follows your system preference by default; the button at the right of the navigation cycles system, light, and dark, and the choice is remembered per browser. Both themes are checked for readable contrast, and the charts use palettes validated for colour-vision deficiency in each theme.

## Agents on other machines

Agents do not work through the dashboard. They use flai's MCP server on the host: over stdio as `.mcp.json` sets up, or over HTTP after `flai mcp start` for an agent that cannot start a process there; setting that up is in [the flai guide](flai.md). What an agent does shows up here like anyone's work. The dashboard's own `/mcp` address, which served MCP until flaiover 0.22, now answers 410 and says so.
