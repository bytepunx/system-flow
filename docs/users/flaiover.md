---
title: flaiover dashboard
updated: 2026-10-23
status: active
---

# flaiover

flaiover is the dashboard for a system-flow project: the board, the work items, the documentation, the decisions, and the flow charts, in a browser. It reads and changes nothing by itself. Everything it shows and everything it changes it asks of `flai serve` on your machine, so a move made on the board is the same move `flai move` makes, committed with your own git identity. The design is in [design/system/flaiover-dashboard.md](../../design/system/flaiover-dashboard.md).

Running it, securing it, and the settings behind it are in the [operators guide](../operators/index.md#running-the-dashboard). Read its warning before you make the dashboard reachable from anywhere but your own machine: whoever holds its token can change your project as you.

## Starting it and logging in

In a project, run `flai dashboard`. It starts the dashboard, and `flai serve` behind it, if they are not running, registers the project with them, and prints a login link. Open the link: <http://localhost:4242/login#token=...>. The token in it never leaves the browser; the page exchanges it for a session cookie and removes it from the address and the history. `flai dashboard token` prints the link again.

**Login** is the page you land on without a session. Paste the token, or open the login link. When your session ends, because the token was rotated or the cookie was lost, any page sends you back here and returns you to where you were once you log in. If the token is right but the browser did not keep the cookie, the page says so and why, rather than asking for the token again; the operators guide explains the case of a tunnel that reports HTTPS.

**Live updates.** A page shows what changes in the project without a reload. Each tab follows its project over one WebSocket, which does not count against the six connections a browser keeps to the dashboard. So any number of tabs can stay open. Before S-0454 each tab held one of those six, and with six tabs open every page of the dashboard stopped loading. When the connection drops, the tab connects again by itself and reads again what it shows.

## The header

Every page but the login page has a top bar across the top and a nav rail at the left.

The top bar holds, from the left:

- **Menu button**, only on a window narrower than 768 pixels: opens the nav rail over the page. Choose a page, press Escape, or click outside the rail to close it.
- **flaiover**: the name, a link to [Overview](#overview).
- **Project switcher**: the project you are looking at, and the others this dashboard serves. See [More than one project](#more-than-one-project).
- **Theme button**: see [Theme](#theme).
- **Host flai badge**, at the right: whether a flai on the host has the dashboard connected, and a link to the [Host](#host) page. See [Host flai](#host-flai).

The nav rail lists five groups. Each page is a section below.

| Group | Pages |
|-------|-------|
| Flow | [Board](#board), [Inbox](#inbox) (with the number of things that need you), [Charts](#charts) (which opens to list the charts), and Workflow, still to come and marked soon |
| Agents | Orchestration (the [Orchestrator](#orchestrator)), Planning (the [Planner](#planner)), Analysis (the [Analyzer](#analyzer)), [Activity](#activity), [Threads](#threads), [Messages](#messages) |
| Design | System (the [Design system](#design-system)), [Components](#components) |
| Documentation | System, Conventions, and Issues, which open those folders in [Docs](#docs); [ADRs](#adrs); [Search](#search) |
| Configuration | [Features](#features), [Agents](#agents), [File System](#file-system), [Projects](#projects), and [Harnesses](#harnesses), the five pages of the [Settings](#settings); then [Host](#host) |

Click a group to show or hide its pages. An open group's icon is coloured and a closed group's is grey. A closed group still shows its pages' numbers, so the Inbox count stands beside Flow. The page you are on is bold and coloured, and its group opens when you arrive. The browser remembers which groups you opened and closed. [License](#license) is at the foot of the rail.

Charts, in Flow, opens and closes as a group does; it is not a link itself. Open, it lists the charts under four grey labels: flow, usage, planning, and strategic. Click a chart to open it. The chart shown is bold and coloured, and Charts is bold while any chart is shown. Arriving on a chart opens Flow and Charts. The browser remembers whether you left Charts open or closed, as it does the groups. On a narrow window, choosing a chart closes the rail.

On a window 768 pixels wide or more, the rail stays at the left. On a narrower one, the menu button opens it.

Every time the dashboard shows is in your browser's time zone, such as `2026-10-07 15:35 EDT`, and every date is the date there. flai records times in UTC; the dashboard converts them as it shows them, including a time flai names in a message it answers with. The text of a document or a narrative, and its front matter, is shown as it is written.

## Host flai

At the right of the header the dashboard says whether a flai on the host has it connected: "host flai: connected" (hover for the version and since when), or "host flai: not connected" on a yellow ground. `flai dashboard` starts that process, `flai serve`, for you; `flai serve status` on the host says why it is not connected. Every page is read through it: the dashboard asks flai for the board, the items, the threads, the documents, the ADRs, the inbox, activity, and search, and flai tells it when a file changes, so a move made in a terminal shows on an open board in under a second. A page waits half a second for changes that arrive together and then asks once, and only for a change to what it shows: an agent's narrative or a thread reply does not make the board ask for itself again. Without it pages have nothing to show, and a yellow banner under the header says so on every page, with the commands that bring it back; the page recovers by itself when flai returns.

Writes go the same way: a move, a new story, a saved document, an acceptance are asked of flai on the host, which does them with your own git identity and, for commits made for you, the dashboard's trailer. The dashboard's image holds no flai and no git. Throughout this guide, "when the dashboard can write" means while flai on the host has it connected. If the flai on the host is older than the dashboard, the banner says what it lacks and asks you to upgrade flai. If the connection is lost in the middle of an acceptance, the page says the acceptance may have completed on the host: reload the board before trying again.

The dashboard and the flai on the host each check the other's release stamp when flai connects, and refuse a build that is not a release (S-0237). A development build of either carries no stamp and is refused, unless the operator allows unsigned builds on the host (see [Unsigned builds](#unsigned-builds)). When this dashboard refuses the flai on the host:

- The badge reads "host flai: refused (unsigned)", with what to do as its tooltip.
- The banner on every page of that project reads "This dashboard refused the flai on the host as unsigned, so the project cannot be shown.", with the reason and what to do under it:

  ```text
  this dashboard refused flai 1.43.0 on the host: unsigned (missing stamp); install a release of flai on the host with flai host upgrade, or with flai self-upgrade and then flai serve stop and flai dashboard; a development build of flai carries no stamp and is refused
  ```

- The reason is `missing stamp`, `invalid signature`, or `wrong component <name>`. It says "a flai" in place of "flai 1.43.0" when the flai named no version.
- The refusal stays until a signed flai for that project connects. While a flai is connected, no refusal is shown.

When it is flai that refuses this dashboard, the dashboard only sees no flai connect, and says "host flai: not connected". `flai serve status` on the host says "flai refused the dashboard: unsigned", with the reason. Either way, run a released dashboard and a released flai from the same release; the [operators guide](../operators/index.md#the-stamp-on-the-connection-s-0237) says how.

### Unsigned builds

The operator can let development builds connect on purpose, with `dashboard.allow_unsigned` on the host (S-0239, [the operators guide](../operators/index.md#allowing-unsigned-builds-s-0239)). While a project's connection was let through that way, a banner under the header names the side with no release stamp, on every page of that project but the login page. It has no X: it stays while the connection does. It reads, for the flai on the host:

```text
unsigned allowed: flai
The flai on the host is a development build, with no release stamp.
The host lets it through on purpose, with dashboard.allow_unsigned; a release of flai with the published dashboard image needs no allowance and shows no banner.
```

For this dashboard's own build, the heading is `unsigned allowed: dashboard` and the second line "This dashboard is a development build, with no release stamp." With both, the heading is `unsigned allowed: flai and dashboard` and both lines show. A released flai and a released dashboard show no banner, whatever the host allows.

## Notices

What an action did is said in a banner on the page: a move on the board, a ticked criterion or a save on an item's page, an accepted ADR, a saved document, a publish. The banner stays until your next action there, or until you dismiss it with the X at its right. Banners that say how things stand, such as flai not being connected, an acceptance not yet pushed, or what is ready to publish, have no X: they go when what they say stops being true.

## More than one project

One dashboard, and one `flai serve` on the host, can serve several projects. The header always has a project switcher showing the project you are looking at, even when there is only one. It lists the projects, then, under "Not imported yet", the repositories on the host that can be imported (see below). Pick another project and every screen changes to it where you are, without reloading. The board you are on shows the other project's board and follows its changes from then on, and the inbox count and the host flai badge follow too. Your choice is remembered in this browser and shown in the page's address and in the header's links, so a bookmark or a copied link opens on the same project. Nobody sees your inbox, board, or documents mixed with another project's; picking one shows that one and nothing else.

To add a project, run `flai dashboard` in it on the host: when the dashboard is already running, that only registers the project with the running `flai serve`, which connects it within a second or two, and it joins the switcher by itself (the list is looked at again every 30 seconds, and whenever you open the switcher). The project needs a `key` in its `system-flow.yaml`; `flai check` says how to add one. If a project you expect is missing, its flai has not connected here yet: `flai serve status` on the host says why. A project `flai serve` serves that is not connected stays in the switcher as "(served, not connected)", with the reason when you hover over it, and a link beside the switcher ("1 not connected: why?") opens the [Projects](#projects) settings page, where each project's state says why. Projects can also be served and removed there, without a shell.

The switcher also says, per project, whether its flai is signed, that is, whether its [release stamp](#host-flai) verified (S-0237). The label after the project's name:

| Label | What it means |
|-------|---------------|
| "(signed)" | Connected, and its flai's release stamp verified |
| "(unsigned)" | Connected, and its flai's release stamp did not verify |
| "(unsigned allowed: flai)", "(unsigned allowed: dashboard)", "(unsigned allowed: flai and dashboard)" | Connected because the host allows unsigned builds, with the sides that have no release stamp ([Unsigned builds](#unsigned-builds)); in place of "(signed)" or "(unsigned)" |
| "(refused: unsigned flai)" | Not connected: this dashboard refused its flai as unsigned |
| "(served, not connected)" | Not connected, for another reason |

For a refused project, the option's tooltip and its line in the "N not connected: why?" link's tooltip give the reason and the version that flai named, such as `Quay: refused: unsigned flai: missing stamp (flai 1.43.0)`. Its state on the [Projects](#projects) page reads `not connected: refused: unsigned flai:` and the reason.

A repository that is not a system-flow project yet can be imported from the board, once the operator has named the folder it is in on the host (`flai serve import add ~/git`, in the operators guide). It appears in the switcher and the project list as "(not imported)". Picking it asks whether to import it and says which tests will run. **Not now** tells you to select a different project; nothing is remembered, so picking the repository again asks again. **Import** adds the standard's folders and files to the repository (nothing that is there is overwritten), runs its tests, and commits the import if they pass. The page then shows the project it has become. If a test fails, nothing is committed: the page shows what failed, the imported files wait in the repository to be fixed and committed, and **Open the project** shows it, since it is a project either way. A repository with uncommitted changes is refused until they are committed or stashed.

## Overview

The front page, which the flaiover name in the top bar opens. For one project it shows the project's name and description from `system-flow.yaml`, the template version it follows, how many items are active and how many archived, and a table counting epics, stories, and tasks in each state from backlog to done. It is the quickest answer to "how much is there, and where is it".

While no project is chosen and more than one is known, the front page shows a list of every project instead, with a filter box, each with whether it is connected, how many stories are in review, how many threads are awaiting you, and whether an agent is attending, so you can tell at a glance which one needs you. Pick one to open it.

## Board

Board shows one column per state with the WIP count against the limit from `wip/kanban/board.md`. The legend above the columns chooses which types of card are shown: click epic, story, or task there to hide the cards of that type, and again to show them. A highlighted type is shown. Every type is shown until you change it, and your choice is remembered in this browser, so it stays the same across reloads and when you come back to the board. The WIP counts count stories whichever types are shown. Below each column's name a line counts the epics, stories, and tasks it holds, such as `1 | 2 | 4`; hover over it, or use a screen reader, for the long form, `1 epic, 2 stories, 4 tasks`. It too counts every card in the column whichever types are shown. Each card shows the ID, title, nature, how long it has sat in its column, and a flag when it is blocked. A draft story, one an agent wrote that you have not finalized, says `DRAFT` after `BLOCKED` (see [The item page](#the-item-page)). A story in ready, in progress, or review whose branch undoes main's change to some paths says `undoes main on 2 paths` after `DRAFT`; hover over it for the paths. Click it, or focus it and press Enter, to open the story's review page instead of its item page. In review, that page shows each path as a card to resolve ([Reviewing a story](#reviewing-a-story)); before review it shows the branch's changes and no cards. The nature, the type, the blocked flag, and the parent sit in one row under a thin line below the title, all in one size. In that row's right corner a story card shows the ID of its epic, and a task card the ID of its story; hover over it, or use a screen reader, for the parent's title. The whole card is one link to the item, so the parent ID is not a link of its own.

**Wide and narrow screens.** On a screen at least 640 pixels wide the six columns sit side by side in one row across the whole page beside the navigation rail, and never wrap. When they do not fit, scroll the row sideways, with a trackpad, Shift and the mouse wheel, or the scroll bar under the columns; the page itself does not scroll sideways. Each column scrolls its own cards, and its name stays at the top while you do. On a phone the columns stack, and each is folded to its name, its WIP count, and its counts line: tap it, or focus it and press Enter or Space, to open or close the column. Ready, In progress, and Review are open the first time, and the rest are closed. The columns you open or close are remembered in this browser, and used only on a narrow screen: on a wide one every column shows its cards. A closed column still takes a card you move into it, its counts change at once, and its menu still opens from its name.

**Colours.** Cards are colour coded: the background is a pastel for the item's nature (feature, improvement, remediation, research, experiment) and the stripe on the left edge is a colour for its type (epic, story, task); the legend above the columns names each colour. The legend's types and nature tags are also filters: click one to hide its cards, and again to show them. A highlighted type, tinted in its own colour with an edge of that colour, means its cards are shown; a plain one, with only its swatch, means they are hidden. A lit tag, brighter than the card pastel and with a darker edge, means its nature is shown; a dim one, on the plain pastel, means it is hidden. Every nature is shown until you change it, and your choice is remembered in this browser, as the types are, so it lasts across reloads and when you come back to the board. The WIP counts and the line under each column's name count every card whichever types and natures are shown. The colours only repeat what the card says in words, so nothing depends on telling them apart. The same colours follow an item elsewhere: the header of its page, its hits in Search (documents stay plain), and the table under the cycle time charts, where each nature's dots and bars are a stronger shade of the same hue. A blocked card also has a red outline, and the card you are dragging turns faded with a dashed edge.

**Agents.** When flai on the host starts agents for ready stories (the operator turns that on with `flai serve enable agent`), a notice above the columns says what it last did: the agent it started, for which story, by which command, or why it could not, or why a ready story is waiting. While review is at its WIP limit flai starts no agent for a ready story, and the notice says so (`review is full (5 of 5): accept or send back a story`): accepting or sending back a story in review makes room. The caret beside its bold title collapses it to that title and opens it again; the browser remembers it collapsed until what the notice says changes, and then it opens by itself. A story's card has a dot beside its ID while its agent is at it: green while it works, yellow while it waits for you, red when it failed. Hover over the dot, or use a screen reader, for what it is doing and why. An agent that finished its story shows no dot. A yellow dot means the agent asked you something on the story, or the story is blocked: answer it on the story's page or from the inbox, and the agent goes on, started again if it had ended. It also means the agent waits on the agent of another story to reply to a message (S-0335): the story's page says `It waits on the agent of S-0004 in MS-0002: it goes on when one replies or writes to this story.`, with the story linked, and flai starts it again when that agent replies or a new message comes to this story, with nothing for you to do. The story's page says the same on its Agent section, with the harness, the model, when it started, and for a red dot where the agent left the story and where its output is on the host. Moving a story with a red dot back to ready starts another agent, and so does changing its agent while it is in ready. For a story in ready or in progress, **Retry**, at the top right of the agent section, starts a new one once the operator has turned on the `agent` host action. Retry disappears when you press it. For a story in ready while the in-progress limit or review is full, the new agent is queued: the dot turns yellow, the page says it waits for room, and flai starts it as soon as there is some. If flai refuses, the page says why and Retry comes back. A story in ready that has no agent yet says why flai has not started one, and **Start agent** starts it at once, under the same conditions. A story in progress that was begun on another host, or outside flai serve on this one, has a yellow dot and the line `begun by agent-S-0173 on another host; no agent here` on its card. Its page says when, too, and **Start agent** starts an agent for it here. That agent is told who began the story, where, and when, takes up its branch from the remote if it was pushed, and goes on from what was committed. Before you press it, check that the agent on the other host has stopped: nothing stops two hosts working the same story.

**Held stories.** A story in ready is *held* while a story in progress touches the same paths (see [Touches](flai.md#touches) in flai.md): flai starts no agent for it and offers it to no agent until that story moves to review, is cancelled, or is sent back. A story in review holds nothing, and neither does an overlap inside the project's shared paths, which the [File System](#file-system) page lists. It is held too while a story it waits for, named in its **Waits for** field, is not done (`held (after): S-0129`). Its card has a yellow dot and, after `BLOCKED` if it is blocked too, a line such as `held (overlap): S-0128`, the reason and the stories it waits for; hover over the dot for flai's whole reason. Once flai has asked the agent of the story that holds it about the hold ([Asked about a hold](flai.md#asked-about-a-hold) in flai.md), a line at the foot of the card says so: `asked S-0128 in MS-0009`. When that agent shares the paths instead, the story can start, and its card says `shares docs/users/flai.md with S-0128 in MS-0009` for as long as the share holds, in ready, in progress, and in review. Click the conversation's ID to open it on the [Messages](#messages) page. Hover over the line for when the agent was asked, or for who shared, when, and the split. The story's page says `HELD` beside its state and gives the whole reason on its Agent section, with each story it names as a link. **Start agent** still starts it, and **Retry** still starts another for one whose agent failed: your word overrides the hold. The card and the page stop showing the hold as soon as it clears, without a reload, and show it whether or not the operator has turned on the `agent` host action.

**Pull order.** The backlog and ready columns are in pull order, the order agents take work in: drag a story up or down within either column to change it, and a line shows where it will land. Stories you have never placed come last, by ID. Without a mouse, focus a card and press Alt+↑ or Alt+↓, or use the ↑ and ↓ buttons that appear on the card when it is hovered or focused. The other columns have no order, so dragging within them does nothing.

**Moving.** Drag a card to another column to move it: the same rules `flai move` enforces apply, and a refused move shows the rule. A story's move takes its epic along, as `flai move` does ([flai.md](flai.md#moving-work)). Dropping a card on cancelled, or choosing cancelled on an item's page, opens a confirmation before anything happens: cancelling an epic also cancels every open story under it and their tasks, cancelling a story its open tasks, and the confirmation lists them, points out any story in review (its work stays on its branch, unmerged, so accept it first if you want it), names the branches and worktrees that will be left for you, and asks for the reason, which is written into every item cancelled. A cancelled story can be moved back to backlog later from its lane's menu; what its cancellation cancelled under it stays cancelled. "Keep it" closes the confirmation and changes nothing.

**The card menu.** Right-click a card for what you can do with it without opening it. On a phone or tablet, press and hold the card; without a mouse, focus the card and press the context-menu key or Shift+F10, or Tab to the **⋯** button that appears at its top right when it is hovered or focused, and press Enter. The menu opens on the card and lists, in order:

| Entry | Offered on | What it does |
|-------|------------|--------------|
| **Open** | every card | opens the item's page, or a story in review's review page, as clicking the card does |
| **Finalize** | a draft story that is not archived | takes the story as yours, as **Finalize** on its page does: `DRAFT` leaves the card |
| **Start agent** or **Retry** | a story, when its page offers it | the same button as on the story's page (see **Agents** above) |
| **Plan** | an open epic or story, while the `plan` host action is on | starts the planner for it on the host, as **Plan** on its page does (see [The item page](#the-item-page)) |
| **Block…** or **Unblock** | an item not done or cancelled, or one that is blocked | asks why and blocks it, or unblocks it, as its page does |
| **Cancel…** | an item not done or cancelled | opens the same confirmation as dropping it on cancelled |
| **Archive** | a cancelled item that is not archived | archives it, as `flai archive` does: a story takes its tasks with it (see **Archive All** below) |

Under the column's name the menu then lists the lane menu's entries for that column (below), with the card ticked first in a move. A notice above the columns says what each choice did, or why flai refused it, and the card shows the result at once. Arrow keys move through the menu; Escape closes it and puts you back on the card; clicking elsewhere, tabbing away, or choosing an entry closes it too. Hold Shift while you right-click to get the browser's own menu instead, for example to open a card in a new tab. A read-only board offers no menu and no **⋯** button.

**The lane menu.** Right-click a column's empty space for what that column allows. Hold Shift while you right-click to get the browser's own menu instead. A read-only board offers no menu.

| Column | Create item here | Move stories forward | Move stories back | Change WIP limit | Archive All |
|--------|------------------|----------------------|-------------------|------------------|-------------|
| backlog | yes, starts in backlog | to ready | | | |
| ready | yes, starts in ready | | to backlog | yes | |
| in-progress | yes, starts in in-progress | | to ready | yes | |
| review | yes, starts in backlog | | to in-progress, with a reason | yes | |
| done | yes, starts in backlog | | | | |
| cancelled | yes, starts in backlog | | to backlog | | yes, while it shows a card |

**Create item here** opens the new item page with that column chosen as where it starts. **Move stories** opens a list of the column's stories to tick, with **All** to tick every one; chosen from a card's menu, that card starts ticked. Each story is moved on its own with the rules of `flai move`, so a story without an acceptance criterion is refused on its way to ready while the others move, and the notice names what moved and what was refused and why. A story moved back to ready goes last in the ready order and, like any story entering ready, may have its agent started by flai on the host. Done is final: nothing moves out of it. **Change WIP limit** sets the column's limit, 0 for none, in `wip/kanban/board.md`, the file flai, flai serve, and `flai check` read it from; a move past a limit warns and is made, and flai serve starts agents only while in-progress has room.

A cancelled item stays in the cancelled column, and `flai check` fails on it, until it is archived. **Archive All** archives every card the cancelled column shows (the type and nature toggles of the legend decide which), and a cancelled card's **Archive** archives that one item. Each item is archived on its own with `flai archive`, stories first, then tasks, then epics: a story takes its tasks with it, so they are not asked for again, and an epic is refused while a story of it is still on the board, so it goes last. A refusal stops only its item, and the notice names what was archived and what was refused and why. Nothing is committed: the items and narratives move to `wip/archive/` in the project's checkout on the host, as `flai archive` moves them, for you to commit.

**Merge queue.** The merge queue verifies the enqueued stories in review as one batch, on the host, and is to land the green ones; landing is not built yet, so you still accept each story yourself ([The merge queue](flai.md#the-merge-queue) in flai.md). Each story in review says `queue 1`, `queue 2`, and so on, its place in the queue, or `not eligible`. Hover over `not eligible`, or use a screen reader, for every reason, such as an unticked criterion, an open thread, or a `flai verify` that did not pass at the branch's head. At the top of the review column a panel shows:

- the batch being verified and its step, such as `rebasing S-0042` or `verifying go-test`, or `No batch.`;
- once the batch is verified, its result: `Passed`, or `Failed at <tier>`, with when and flai's summary;
- a dropped batch's reason, who dropped it, and when, such as a story taken out of the queue or accepted by hand;
- the eligible stories in the queue's order, each marked enqueued when it is, or that none is eligible, and how many are not eligible.

With the `merge-queue` host action on and a dashboard that can write, the panel has a **Hold** switch. Switched on, it holds the queue: no new batch starts, and a batch under verification finishes and waits. The panel then says `Held by <who> at <time>`. Switched off, the queue builds batches again. While the action is off, the panel says no batch runs, and names the command the operator runs on the host, `flai serve enable merge-queue`. While a batch is under verification, the board asks for the queue every 10 seconds, so its step and result show as they change. Under pull-request acceptance the panel says the queue does not run there, and why, and no card is marked. When the flai on the host sends no queue, being older than the queue or unable to read it, the panel says so.

**Accepting.** Dropping a story on done accepts it, so the dashboard first shows what acceptance will do — the branch to be merged, anything that blocks it — and only proceeds when you confirm. When the story is its epic's last open one, the epic is accepted and archived with it, as with `flai accept` ([flai.md](flai.md#accept-and-release)). When the story's agent left changes uncommitted in its worktree, the confirmation lists them and offers **Have an agent commit them**. Accepting merges, archives, and commits; it cuts no release by itself (S-0087). An experiment story is accepted once its results document, `design/experiments/<S-nnnn>-<slug>.md`, is committed on its branch; until then the confirmation names the document to write and its accept button stays disabled, and cancelling leaves the card in review. Like research, an experiment releases nothing, and the confirmation lists the components it touched as landing on main without a release. If files outside `wip/` are uncommitted, the confirmation lists them: acceptance refuses them by default so that its commit holds only acceptance, so either cancel and commit or stash them first, or tick the box to include them, which is the dashboard's form of `flai accept --yes`. The accept button waits for that choice. When acceptance cannot run from the dashboard, the confirmation lists why and what to do, and its accept button stays disabled; accept from a shell with `flai accept S-nnnn` instead. An acceptance is made on the host by flai, committed there. The review page (below) is the fuller way to accept.

**Publishing.** Accepted work stays in the project's clone on the host until it is published, and publishing is the one way it reaches the remote ([ADR-0067](../../design/adrs/0067-accepted-work-reaches-the-remote-only-when-it-is-published-and-agents-publish.md)). Publishing is the operator's choice; agents publish only when the operator asks. The board has no push notice and no push button. Whether a story has been **published** (released, with a version bump and a tag, and pushed) is shown in the done column itself: a card says **published** or **waiting to publish**, and a banner at the top of the column lists what publishing now would release, by component, with the story IDs it bundles (several small stories against the same component release together, at the highest bump among them, not one release each). The banner's **Publish** button, when the operator has enabled it on the host (`flai serve enable push`), runs `flai release --pending` there: it tags and pushes the branch and the tags, never forcing. A publish that fails partway is safe to run again; it does not redo what already succeeded. When this clone is behind its remote, the banner offers nothing to publish and says why, with the command to run in the project's checkout on the host: release tags the remote has and this clone lacks, as after publishing from another clone and pulling without tags (`git fetch --tags origin`), or commits on the remote branch this clone lacks (`git fetch origin && git merge origin/main`, or rebase onto it instead of merging). While tags are missing, the done column leaves out the stories it would otherwise list as waiting, since they may already be published. When the remote cannot be reached, the plan is shown with a warning and no **Publish** button. An accepted story no release could cover, such as one touching two components with no tag saying which it delivers to, is listed in the banner with the reason.

### Creating work

When the dashboard can write, the board has a "+ new" action, which opens the new item page, and so do **New story**, **New sibling**, **New epic**, and **Create story** at the top of an item page (see [The item page](#the-item-page)). Choose epic or story, for a story an epic to belong to or "No epic" (not every story fits an active one, S-0092), and the nature (each is shown with what it means and what bump it contributes once published). Give it a title, optionally tags (a component's name decides which component it delivers to), topics (what it is about beyond its components, such as `logging`; see [A story's topics](flai.md#a-storys-topics)), and the paths it touches. For a story you may also name its agent: the harness, the model, the provider its model's calls go to (one of the host's [providers](#agents), such as `litellm`; [Where its model's calls go](flai.md#where-its-models-calls-go-a-provider)), and options one `key=value` a line. The project's default agent (`flai agent set`) shows in grey in those fields, `none` for a provider it does not name, and is what the story gets for anything you leave empty. For an epic or a story, the **Cost of delay** panel, closed until you open it, takes what waiting for it costs: revenue per week and penalty per week, in the project's currency (`planning.currency` in `system-flow.yaml`, default USD), and the time lost per cycle, a duration such as `6h`. Each is optional, and while the panel is closed it says in one line what you gave. They are recorded as set by you, and the planner works out the value from them ([Drafts, cost of delay, and forecasts](flai.md#drafts-cost-of-delay-and-forecasts)). Then write what it is for in markdown: the form starts from your project's template sections and shows the preview beside the text; acceptance criteria are `- [ ]` lines. You write no front matter. The ID, the file name, the front matter, the link in the epic's list, and the commit (your git identity, with a flaiover trailer) are made by `flai`, so the file is exactly what `flai story new` would have made. If `flai check` would report anything with the new item in place, nothing is created, the findings are shown as `flai check` prints them, with the file and the line, such as `wip/kanban/stories/S-0231-author-license-md-file.md:23: warning: markdown.MD007: MD007/ul-indent Unordered list indentation [Expected: 0; Actual: 1]` for a list indented one space, and your text stays in the form. The line counts from the top of the file flai would have written, front matter included. Choose the lane it starts in: backlog, ready, or in-progress. Opened from a lane's menu on the board, the form starts on that lane, or on backlog for review, done, and cancelled. flai makes every item in backlog, and the form then moves it on, one column at a time, with the rules of each move: a story needs an acceptance criterion to be ready. If a move is refused the item still exists where it got to; the form says which move and why, links the item, and will not make it again. After creating you land on the item's page, and the board shows it without a reload. Tasks are not created here: the agent that pulls a story writes them, or the planner drafts them. Nothing is pushed.

### The item page

Click a card for the item page. Its top has the ID and title, `[Draft]` beside the title of a draft story, the type and nature, the status, a BLOCKED flag, `HELD` on a held story, `undoes main on N paths` on a story whose branch undoes main's change (as on its card, a link to its review page, with the paths as its tooltip), the parent, the page's notices, and the buttons for what you can do with the item. These stay at the top whatever you read below them.

Below them, a list of the page's sections sits at the left, and the page shows one section at a time: click an entry to read it. The list is a column beside the section on a wide window, and stays in view as you scroll; on a narrow one it is a row above the section, which you can scroll sideways. Every entry is a link, so Tab reaches it and Enter opens it, and the section you read is in the page's address: back and forward go between the sections you read, and a link you copy opens the same section. The sections, in this order, each left out when the item has nothing for it:

| Section | What it holds |
|---------|---------------|
| Story (Epic, Task) | The rendered body, or its editor, and on a story in review a **Review this story** link. Shown when you arrive |
| Tasks (an epic's Stories) | The children, with their states, and a story's task plan |
| Questions and threads (Threads off a story) | On a story, its narrative's open questions, each answerable there when the dashboard can write; then the threads on the item. Its entry counts those still open. A link from the inbox to a thread or a question opens this section, with that one outlined and scrolled to |
| Messages | A story's conversations with other stories' agents ([Messages](#messages)). A board card's link to a conversation opens this section |
| Agent | What the story's agent is doing, and its pull request when `flai serve` has read one ([below](#a-storys-pull-request)) |
| History | Every transition with when and by whom, and the intervals the item was blocked with their reasons |
| Files | The item's file, which opens in the explorer, and a story's narrative, with a box to add a log entry to it |
| Planning and usage | The usage, the cost of delay, the forecast, the expected cost, and the owner, estimate, tags, topics, touches, and agent from the front matter |

Pressing **edit…** from any section opens the editor on the Story section.

Threads are shown one at a time, wherever they are shown: with more than one, a pager reading "← 2 of 5 →" beside the Threads heading, and again under the thread, says which you are reading, and its arrows go to the previous and next; with a pager's arrow focused, the Left and Right keys page too. A thread you open becomes the one shown, and a reply or a resolution keeps your place. A thread you reach from the inbox is the one shown, and the page scrolls to it. A thread of more than two entries shows its last two, with **show N earlier entries** above them to read the rest and **hide earlier entries** to fold them again; each thread keeps its own choice while you page. A thread's entries render as markdown, as the body does, wherever threads are shown. Your entries sit on the right on a blue ground and an agent's on the left on a neutral one, so who said what reads at a glance; you are the project's `owner` in `system-flow.yaml`. When the dashboard can write, the buttons at the top make the moves the item's state allows, block it with a reason or unblock it, and open **edit…** for a story or an epic (see [Editing a story or an epic](#editing-a-story-or-an-epic)). On a draft story, **Finalize** among them takes the story as yours: it may then go to ready, flai records you as who finalized it and when, as `flai edit --no-draft` does ([Drafts, cost of delay, and forecasts](flai.md#drafts-cost-of-delay-and-forecasts)), and `[Draft]` goes without a reload. A story someone finalized meanwhile is refused, with a word to reload it. On an open epic or story, **Plan** among them starts the planner for it on the host, as `flai plan` does: for an epic it drafts the stories and their tasks, or revisits those it has and drafts the tasks of each draft among them that has none; for a story it adds touches, a forecast, and a cost of delay value, and drafts its tasks or revisits those it has ([Running the planner](flai.md#running-the-planner)). It shows only while the operator has turned on the `plan` host action (`flai serve enable plan`, or the [Features](#features) page). The page says that the planner started, with its log on the host, or why flai refused, and says `planner running since` while it runs. The stories it drafts appear under the epic as drafts, each with its tasks, for you to read and finalize. On a story the acceptance criteria are live checkboxes: tick one where you read it. At the top right, beside the title, **New story** on a story's page opens the new item page for a story on its own, and **New epic** on an epic's page opens it for another epic. On a story under an epic, **New sibling** follows **New story** and opens the new item page for a story already under the same epic; it shows only while that epic is open, so a story under a done, cancelled, or archived epic has none. After **New epic**, an open epic's page has **Create story**, which opens the new item page for a story already under that epic; a done, cancelled, or archived epic has none, because a story cannot be added to it. You can still change any of these choices there. A task's page has no such link, because tasks are written by the agent that pulls a story, or drafted by the planner.

The page keeps itself current: when the story moves, gains a task, is edited, or its agent starts or ends, the page shows it without a reload, and a reply in its threads appears as it is written. After you answer the agent in a thread, three moving dots under your reply say the agent is working on it, until its next entry arrives or it stops.

On Tasks a draft story says `Draft` after its state. On Planning and usage the agent is on one line with its provider as `via <provider>`. An epic or a story with a cost of delay shows the planner's value per week with who set it and when, or that there is no value yet, then the inputs with who set them and when. When the inputs changed after the value was set, a note says the planner's value is out of date: it is worked out again from the new inputs.

An item's usage shows what its agents spent, then what the planner spent planning it on a line of its own, marked estimated, apart from the agents' figures. That line shows even before any agent works the item. An item with a forecast duration, or else an estimate, shows its expected cost after the forecast: that duration priced at the project's cost per agent hour, marked estimated, with what it was priced from. It shows once a story has been measured from its logs ([ADR-0083](../../design/adrs/0083-a-planner-activity-s-usage-is-charged-to-the-item-it-planned-and-the-items.md)).

### Reviewing a story

A story in review has a review page: open it from the card in the review column, or from **Review this story** on the item page. It puts what you need to decide in one place.

Its top has the story's ID, title, and status, a notice once it is accepted or when it is not in review, and the **Acceptance** box: what blocks acceptance, **Accept**, **Send back**, and **Enqueue**, each described below, and what each step of an acceptance did. These stay at the top whatever you read below them.

**What blocks acceptance.** The box lists the paths the branch undoes of main's change, as cards (see **A branch that undoes main's change** below), then anything else that blocks acceptance from here, and any uncommitted files outside `wip/`, which you choose to include or deal with first. When the branch changes a path Claude Code protects, such as a file in `.claude/` or `.mcp.json`, a box says that only the operator accepts the story, names who that is, and lists the files: read them in the diff before you accept, since they change what every agent runs with once merged ([ADR-0106](../../design/adrs/0106-a-story-whose-branch-changes-a-path-claude-code-protects-is-accepted-by-the.md)). Changes left uncommitted in the story's worktree block it: **Have an agent commit them** starts the story's agent to commit them on its branch, once the operator has turned on the `agent` host action, and you accept when it has finished. The board accepts as the project's `owner`, so your **Accept** goes through when `owner` in `system-flow.yaml` names you ([Accept and release](flai.md#accept-and-release)).

Below them, a list of the page's sections sits at the left. Click an entry to scroll to its section; the entry of the section you are reading is highlighted. The list is a column beside the sections on a wide window, and stays in view as you scroll; on a narrow one it is a row above them, which you can scroll sideways. The sections, in this order:

| Section | What it holds |
|---------|---------------|
| Criteria | The acceptance criteria, with which are ticked and how many |
| State | The narrative's current state and next steps, as the agent left them, with a link to the whole narrative |
| Threads | The threads on the story, where you can ask before deciding |
| Changes | What the story's branch changes against main, file by file. Under its entry, one entry per changed file, by its file name: click one to open that file's changes and scroll to them. A link copied from that entry opens the page at the same file |
| Release plan | What accepting will do: the branch that is merged, the move to done and the archive, and the release it would cut once published |
| Issues | The open issues that no open story is already working. Not shown when there are none |
| Checks | The host's checks for the story: the last run, its steps and output, and **Run checks** or **Cancel** when the dashboard can write. Off until the operator runs `flai serve enable checks` on the host |
| Verification | The story's last `flai verify`. Not shown when the story has none |

**Changes.** Each file is a line with its path and the lines added and removed. Click the line, which has an arrow at its start, to open the file's changes below it, and again to close them. A removed line has `-` in the margin and a dim red background, an added line `+` and a dim green one, and lines removed or added one after another share one brighter outline, so a block that was replaced reads as two blocks, the old above the new. Large patches are cut and marked; the rest is a `git diff` away. A story that was not worked on a branch says so instead. When the branch changes a component, the section names the components first, with **Open in Components**, which shows them from the story's worktree ([A story's components](#a-storys-components)).

**Release plan.** Accepting cuts no release by itself (S-0087); publishing what has accumulated is a step of its own, from the done column.

**Issues.** The issues this story recorded or bumped come first, checked. Every other open issue follows, unchecked. Each issue's ID opens its document.

**The pull request.** The story's pull request shows after Checks, as on the item page ([A story's pull request](#a-storys-pull-request)).

**Verification.** The story's last `flai verify`, which its agent ran through the close-out before moving it to review ([Verify a story before review](flai.md#verify-a-story-before-review)). It says whether it passed every step or the step it stopped at, the commit it verified and the main branch it was verified against, when, and how long it took. Each step follows, in the order it ran, with its state (passed, failed, not reached, or reused) and its duration; a test or lint step is marked `tier`. A reused tier passed in an earlier run at the same commit and main branch, and was not run again ([Verify a story before review](flai.md#verify-a-story-before-review) says when). In place of its duration it shows `from` and when that run began, in the browser's time zone, as `from 2026-10-08 05:12 EDT`. Its state is a lighter green than a pass. Under the step that failed are its findings, each with its file, line, and message. The findings of `flai check` outside the story are notes that do not fail it, collapsed under a line that counts them. When the commit is not the branch's head, the result is older than the branch. The page only shows the result; it has no button to run it. When the story has no stored result, the section is not shown.

**Enqueue** puts a story in review in the merge queue, as the project's owner ([Enqueueing a story](flai.md#enqueueing-a-story)). Once it is enqueued the page says who enqueued it and when, and the button reads **Dequeue**, which takes it out. The button shows when the dashboard can write, the operator has the `merge-queue` host action on, and the queue runs, which it does not under pull-request acceptance. When flai refuses, for a story that is not eligible say, the page shows flai's sentence and each reason under it, and the queue is unchanged. Accepting the story yourself takes it out of the queue.

**Accept** runs the acceptance as you, the project's owner, and shows each step as it completes: the branch merged, the story moved to done, the archive, the commit. Nothing is tagged or pushed by acceptance: the page says the accepted work stays local until it is published, from **Publish** on the board or with `flai release --pending` on the host. If flai refuses or fails, for example a rebase that stops on a conflict or an unticked criterion, the page shows flai's message word for word, and the story is still in review. Each file the message names opens in the editor.

**Fixing markdown that blocks an acceptance.** When the acceptance is refused for markdown its lint fix could not fix ([Accept and release](flai.md#accept-and-release)), the page lists each finding: the file and line, the rule, and the message. Click a file to open it in the editor at that line, on the story's branch when the story's branch holds it, otherwise on `main`. Correct it and save, and the finding is marked saved. Once every listed file is saved, **Accept** is enabled again; press it to accept, with no reload. A move to done refused this way on the board or on the item page lists the findings as this page does, and each file opens on the same checkout. Any other refusal shown there links each file it names to the editor too. Once accepted, the page says who accepted it: `S-0031 is accepted by alex.`

**A branch that undoes main's change.** When the story's branch, rebased onto main, removes lines main added since the branch left it, accepting would take main's change away ([ADR-0155](../../design/adrs/0155-a-rebased-story-branch-that-undoes-main-s-change-to-a-path-is-recorded-at-sync.md), [Accept and release](flai.md#accept-and-release)). The page shows each such path as a card, above anything else that blocks: the path, main's lines the branch removes, and beside them the branch's change to that file. Each card offers three choices.

| Choice | What it does |
|--------|--------------|
| **Restore main's version** | flai commits main's version of the path on the story's branch. The card goes, and the branch's changes show the file as it now is. If flai refuses, for example because the worktree has uncommitted changes, the card says why. |
| **Keep** | For a removal the story means. Type why, in one line, and press **Keep the removal**. The card turns yellow and shows your reason; **Undo** takes it back. The reason goes into the story's Notes, under `Removals kept`, when you accept. |
| **Send back, naming the undone paths** | Opens a thread on the story that names the paths for its agent, then moves the story to in-progress with a reason that names them. Its agent restores the paths or redoes the change. |

**Accept** stays disabled until every card is restored or kept, so that main's change is never taken away without your decision. If the acceptance's own rebase finds another such path, the refusal adds its card, and you resolve it the same way. Only you keep a removal. The orchestrator never does, so such a story waits for you even with `accept_reviews` on.

**Who accepted.** A done story's page says who accepted it and on which day, under its type and status: `Accepted by alex on 2026-10-06`. With `accept_reviews` on, the orchestrator accepts stories in review itself ([When the orchestrator accepts](flai.md#when-the-orchestrator-accepts), [ADR-0093](../../design/adrs/0093-with-accept-reviews-on-the-orchestrator-accepts-a-story-in-review-through-flai.md)). A story it accepted says `Accepted by orchestrator`, and the line links to its evidence: the section `Accepted by the orchestrator` in the story's Notes, with the commit its verifier passed, when, the verdict, and the files that meet each criterion. The history lists each move with who made it, as before.

With any issue checked, the button reads **Accept and Create Stories**. It accepts first. Only if the acceptance succeeds does it make one backlog story for each checked issue, as `flai issue story` does. Each story is committed on its own. The page names each story it made. If some could not be made, it names each issue that failed and what flai said. If the acceptance fails, no story is made. Uncheck every box to accept without making stories.

**Send back** asks why, in the page, and moves the story to in-progress with your reason recorded in its notes, where the agent reads it.

### A story's pull request

Under pull-request acceptance, a story's page and its review page show its pull request as `flai serve` last read it from the forge ([Checks, reviews, and the send back](flai.md#checks-reviews-and-the-send-back)). On the story's page it is on the Agent section, after what the story's agent is doing. On the review page it comes just before Verification.

- **Pull request #12** links to the pull request on the forge. Below it are the head commit and when `flai serve` read it, in your time zone.
- **Required checks** lists each check the base branch requires, with its state, its name, a **log** link, and the forge's summary. The state is green for passed, red for failed, and yellow for pending. `None required.` means the branch requires none.
- **Reviews** lists each review, with its state, the reviewer, and a **review** link. The state is green for approved, red for changes requested, yellow for pending, and plain for a comment or a dismissed review. `None yet.` means nobody has reviewed it.

The page reads what flai keeps, never the forge, so it is at most a minute behind while the `forge` host action is on. The story's page asks again whenever the story changes, so a story flai sends back shows the failure as it moves. Nothing shows for a story `flai serve` has read no pull request for, or when the read fails.

When a required check fails or a reviewer requests changes, flai opens a thread on the story and moves it back to in-progress, and `flai serve` starts its agent. The thread shows on the story's page with its other threads.

## Inbox

Inbox is the list of things that need you, and the number beside its link in the nav rail, and beside Flow while Flow is closed, is how many there are. It holds five kinds of entry, each a link to where you deal with it, never to a document you could not answer it from: stories in review, which open their review page; threads where someone other than you wrote last, which open the story or item they are on (or whose file or narrative they are on) at that thread, and otherwise the [Threads](#threads) page at that thread; open questions agents left in their narratives, which open the story's page at that question; blocked items, with the reason, which open the item; and overlapping touches, where two stories in progress say they change the same files, which open the first story named. An entry leaves the list when its cause does: you answer the thread, accept the story, the item is unblocked. A hand-written open question (one with no thread of its own) can be answered right there, in the inbox, or on its story's page, when the dashboard can write: type the answer and submit it, and it moves to the narrative's Decisions and drops off the list. A thread where an agent, such as the orchestrator, recommends an answer waits for you until you confirm it: its entry says so, shows the recommendation, marked as one, with the document and heading it cites as a link, and, when the dashboard can write, **Confirm**. Confirm makes the recommendation the thread's answer, and the entry drops off the list; to answer otherwise, open the thread and reply. The list and the count refresh by themselves when files change.

You can ask the browser to tell you when something new arrives: tick "Desktop notification" on the inbox page and allow it when the browser asks. It is off until you turn it on, it is remembered per browser, and it only fires for entries that appear while a dashboard tab is open, never for what was already there.

## Threads

Threads shows every open thread in the project, whatever it is on, one at a time with the same pager, **show resolved**, reply, resolve, and **Confirm** as a story's page. An entry that recommends an answer is marked **recommendation**, and the document and heading it cites are a link under it. A recommendation is not the answer until you confirm it: under a thread whose recommendation awaits you, a line says so and, when the dashboard can write, **Confirm** makes it the answer, and the thread shows answered at once. Each thread links to the item or document it is on. A thread in the inbox that is on no item, such as one on a design issue, opens here at that thread. New threads are started on the item or document they are about, not here.

## Messages

Messages shows the conversations between the agents of two stories in progress or in review, about the paths both touch: one asking the other who changes a file first, or flai telling both of a conflict it found at a sync. Each shows its two stories, linked, the paths it is about, which story it awaits, how long ago it was last written to, and its entries. Open conversations come first, newest first; tick **show closed** for those closed, by an agent or because one of the stories was accepted or cancelled. A conversation the two agents could not settle is marked **escalated**, with a link to the thread it opened for you: that thread is where you answer. When the agent of a story in progress shared paths with a ready story it held, the conversation shows the share under the paths it is about: which story shared which paths with which, who shared them and how long ago, and the split, who changes what. A link from a board card opens the page at its conversation.

Messages are read-only here. They never reach your inbox or its count: the agents settle them between themselves. A story's page shows its own conversations under its threads, each named by the other story.

## Activity

Activity shows who is working on what: one card per story with an open narrative, with the agent and session that last wrote it, how long ago, the story's state, the task in progress, whether anything in it is blocked, and the last line of its log. Nothing is recorded to produce this. An agent that stops writing does not disappear; its card simply grows older.

When `flai serve` started an agent for the story, the card also shows the agent's dot and what it is doing, and a **stream** window: what the agent said, each tool it called (`▸ Bash List files`) and what the tool answered (`↳ …`, in red when it failed), its background tasks, and how its session ended, read by flai from the log it gave the agent. The window opens when the agent starts and follows it every two seconds, staying at the newest entry while you are at the end; scroll up and it stays where you put it until you scroll back to the end. It closes only when you close it, so one you were reading stays open when its agent ends; an agent that had ended before you opened the page shows its window closed, to open when you want to read what it did. The page holds still under you while it reads the streams, the agents, and the projects again: a reload that fails says so above the cards rather than in their place, and the windows stay while flai on the host is away. Long answers are cut to their start, and older entries are in the log on the host (`flai serve agent stream <story>` prints it there). An agent at work on a story that has no narrative yet gets a card of its own. Your flai on the host needs to be as new as the dashboard for this: an older one is named in the banner, and the window says it cannot read.

**Stop** is beside an agent that runs, or that ended waiting for an answer or another story's agent's reply, once the operator has turned on the `agent` host action. It asks first, and says what stopping does: the agent's process and everything it started end at once, and cannot be resumed; an agent waiting for an answer is not started again when the answer comes; the story's worktree keeps what the agent changed, committed or not; and the story stays where it is, with no agent, until you press **Retry** on its page or move it back to ready. Only **Stop the agent** in that dialog stops it. The agent then reads as stopped by the operator, with a red dot. If flai refuses, for an agent that has already ended, the dialog says why.

**The orchestrator.** When the operator has turned on the `orchestrate` host action (`flai serve enable orchestrate`), `flai serve` runs the project's orchestrator, and its newest run has a card above the stories: what it is doing (working, or stopped by the operator, ended, or failed, and why), its agent, harness, and model, its session, and when it started and ended, with its stream window, open while it runs. A run that could not be started says why and has no window. The page reads it again when the orchestrator starts or ends. Its decisions, its runs, and **Stop** and **Start** are on the [Orchestrator](#orchestrator) page.

## Planner

Planner shows the agent that plans epics and stories: what it is doing, what it has done, and what that cost, and below that its settings. Find it in the rail's Agents group as Planning, after Orchestration.

- **The plan host action** says whether `flai serve` starts the planner here. While it is off, the page gives the command that turns it on, to run on the host in the project: `flai serve enable plan`.
- **Running now** shows the planner run under way, if there is one: the epic or story it plans, the agent, when it started, and why (`asked`, or what changed to start it again). Its stream window shows what it is saying and doing, as the Activity page does for a story's agent.
- **Activity** is the planner's activity document, `wip/agents/planner.md`. It shows the totals (what the planner has cost, how long it has worked, how many activities it has logged, and when it last ran), then each activity, newest first, with what it did, why it started, the items it changed, how long it took, and what it cost.
- **Runs** lists the newest planner run for each epic or story it has planned, newest first. Each shows when it started and ended, how it went, and what it cost: the activities logged while it ran. A run that failed says why.
- **Plan an epic or a story** takes an ID, such as `E-0016` or `S-0259`, and starts the planner for it, as **Plan** on the item's page and the card's menu do. The page says what flai answered: the run it started, or why it refused, such as a planner already running for the item. While the plan host action is off, the form is disabled and says so.
- **The planner's settings** are the `planning` keys of `system-flow.yaml` but `currency`, edited here as [Editing a strategic agent's settings](#editing-a-strategic-agents-settings) says:

  | Setting | What it does |
  |---------|--------------|
  | `agent` | The planner's agent, over the project's: what it sets wins, and what it leaves out is the project's agent's. Its harness, model, and options; empty, the planner runs on the project's agent |
  | `replan` | What flai serve does when a story is accepted or cancelled or the pull order changes: `never` nothing, `deterministic`, the default, plays the board out again and moves forecast deliveries, and `agent` does that and queues the planner for each story whose delivery moved |
  | `schedule` | When flai serve runs the planner over the ready column, a five-field cron expression in UTC, such as `0 6 * * 1-5`, or `daily`; unset, there is no schedule |
  | `hour_rate` | What an hour of work costs, in the project's currency, which prices a cost of delay's time lost; unset means unknown, not free |
  | `cycle` | The period a cost of delay's time lost is counted over; `168h`, a week, by default |
  | `default_duration` | The work a story is forecast to take when there is too little history to forecast it from; `1h` by default |

  `agent` and `schedule` start planner sessions you pay for without asking each time, and `replan` set to `agent` does too. The currency is set by hand in `system-flow.yaml`, since changing it would re-denominate every amount on the items.

The page follows the planner as it works: it reads itself again when the planner's document or a work item changes, and when a planner run starts or ends. It reads its settings again after a save and when `system-flow.yaml` changes, never on a timer, so a form you are typing in is not refilled under you.

## Orchestrator

Orchestrator shows the agent that keeps the project's work moving: whether it runs, what it is doing, what it has decided, and what that cost, and below that its settings. Find it first in the rail's Agents group, as Orchestration. The orchestrator runs while the operator has the `orchestrate` host action on ([Running the orchestrator](flai.md#running-the-orchestrator)).

- **The orchestrate host action** says whether `flai serve` runs the orchestrator: on; off, with the command that turns it on, `flai serve enable orchestrate`; or held, because you stopped it.
- **Stop** ends the orchestrator's run and holds it stopped. While the action is on, `flai serve` starts the orchestrator again whenever a run ends; while it is held, it does not. **Start** lifts the hold and starts it at once. Stop shows while the orchestrator is not held, and Start while it is. Both need the `orchestrate` host action on, and are disabled, saying so, while it is off. The page says what flai did, or why it refused, such as an orchestrator already held. The run you stopped logs its activity as `stopped: stopped from the dashboard`. Turning the action off and on again lifts the hold too. On the host, `flai serve orchestrate stop` and `flai serve orchestrate start` do the same.
- **Running now** shows the run under way, with its agent, when it started, and its stream window, as the Planner page does.
- **Last decisions** are the five newest entries of its activity: what it did and why, with the policy figure behind it.
- **Activity** is `wip/agents/orchestrator.md`: its totals, then every decision, newest first, as the Planner page shows the planner's.
- **Runs** lists the orchestrator's runs, newest first, the newest and up to 20 before it: when each started and ended, how it went, what it cost, and why it failed.

The page reads itself again when the orchestrator's document changes, and when a run starts or ends or is held or let go.

Its settings are the `orchestration` keys of `system-flow.yaml` but `agent`, edited as [Editing a strategic agent's settings](#editing-a-strategic-agents-settings) says. **Permissions** are what it may do without you. Each is off until you turn it on, and a change holds from the orchestrator's next step: it is not restarted, and it reads its permissions again before each decision. Beside each the page says what it allows and what can go wrong:

| Permission | Lets the orchestrator | Risk |
|------------|-----------------------|------|
| `plan_backlog_epics` | Ask the planner to draft stories for an epic in the backlog | The planner runs, and spends, on an epic you may not mean to start yet, and its drafts fill the backlog |
| `plan_backlog_stories` | Ask the planner to plan each story in the backlog that lacks touches, a forecast, a cost of delay value, or tasks, answer and resolve the threads that planner opens, and give such a story its cost of delay inputs when neither it nor its epic has any | The planner runs, and spends, on every unplanned backlog story, and the orchestrator answers its questions and sets cost of delay inputs, money decisions, without you |
| `finalize_drafts` | Finalize a draft story, as `flai edit --no-draft` does | A story the planner drafted becomes one that can be promoted without your having read it |
| `promote_to_ready` | Move a story to ready | Agents may pull, work, and spend on a story you have not chosen to start |
| `order_ready` | Write the order of the ready column by its policy | It replaces the order you set on the board, so the next story pulled is the policy's choice, not yours |
| `answer_threads` | Answer threads: `off` leaves them to you, `recommend` replies with a recommendation for you, and `autonomous` answers them itself | With `autonomous`, agents act on answers you did not give; with `recommend`, the decision stays yours |
| `accept_reviews` | Accept a story in review, as `flai accept` does | Work is merged to the main branch without your review, and a story accepted wrongly has to be undone by hand |
| `publish` | Release and push accepted work, as `flai release --pending` and `flai push` do | A release reaches everyone who pulls or installs the project, and a published release cannot be taken back |

Below them are its policies:

| Setting | What it does |
|---------|--------------|
| `policy` | How the ready column is ordered: `cod` by cost of delay, `wsjf` by cost of delay over forecast duration, `throughput` shortest first, and `fifo`, the default, leaves your order alone |
| `release` `policy` | When accepted stories not yet released are due a release: `judgement`, the default, leaves it to you, `threshold` at a value or a count of unreleased work, and `theme` when every story of an epic or a tag is accepted |
| `release` `value` | Under `threshold`, the unreleased cost of delay per week, in the project's currency, at which a release is due |
| `release` `count` | Under `threshold`, the number of accepted stories not yet released at which a release is due |
| `release` `epic` | Under `theme`, the epic, such as `E-0001`, whose stories, every one accepted, make a release due |
| `release` `tag` | Under `theme`, the tag whose stories, every one accepted, make a release due |
| `release` `whole_epics` | Under every policy, holds a release back while a story accepted and not yet released belongs to an epic in neither review nor done |

The orchestrator's own agent is set by hand in `system-flow.yaml`, under `orchestration.agent`.

## Analyzer

Analyzer shows the agent that looks at the whole project and writes a report under `design/analysis/` ([Running the analyzer](flai.md#running-the-analyzer)): whether it may run, what it is doing, what it has found, and what that cost, and below that its settings. Find it in the rail's Agents group as Analysis, after Planning. Its reports are on the [Docs](#docs) page.

- **The analyze host action** says whether `flai serve` may start the analyzer here. While it is off, the page gives the command that turns it on: `flai serve enable analyze`.
- **Run** starts the analyzer now, looking for the **Focus** you choose: `bottlenecks`, `intent`, `risk`, or all three, the default. The page says what flai answered: the run it started, or why it refused. One analyzer runs at a time, so Run is disabled while one runs, and while the host action is off. The page says which.
- **Running now** shows the run under way: its focus, its agent, when it started, why it started, and its stream window.
- **Activity** is `wip/agents/analyzer.md`: its totals, then each run's entry, newest first, with a link to the report it wrote, or `no report`.
- **Runs** lists its runs, newest first, the newest and up to 20 before it. Each shows its focus, its report, when it started and ended, how it went, what it cost, and why it failed.

The page reads itself again when the analyzer's document changes, and when a run starts or ends.

Its settings are the `analysis` keys of `system-flow.yaml`, edited as the next section says:

| Setting | What it does |
|---------|--------------|
| `agent` | The analyzer's agent, over the project's: what it sets wins, and what it leaves out is the project's agent's |
| `schedule` | When flai serve runs the analyzer, a five-field cron expression in UTC, such as `0 6 * * 1`, or `daily`; unset, it runs only when you ask. Each scheduled run is a session you pay for |

### Editing a strategic agent's settings

The Planner, Orchestrator, and Analyzer pages edit their agents' settings the same way. Each setting shows its value, or while it is unset its default, and what it does.

- **Read-only until the operator allows it.** The settings change `system-flow.yaml`, so they can be changed here only while the operator has turned on the `settings` host action for the project. Until then every field is disabled, and the page says why and gives the command to run on the host, in the project: `flai serve enable settings`.
- **Save** writes only the settings you changed, through flai on the host, and commits `system-flow.yaml` with the dashboard's trailer; the page then says the commit. A field you empty is removed from the file, so its default applies. For an agent, give its harness, its model, and its options, one `key=value` a line; empty all three to remove it. Its sub-agents' roles are kept as they are.
- **A refusal.** flai checks the whole `system-flow.yaml` as it would be after your change, as `flai check` does. When it finds something wrong, nothing is saved, what you typed stays, and the reason is said under each field it is about, such as a release value below zero or a schedule that never comes round. A problem elsewhere in the file is said at the top of the panel: it blocks every save until it is fixed by hand.

The same settings are changed in a shell with `flai manifest set` ([Changing the strategic agents' settings](flai.md#changing-the-strategic-agents-settings)).

## Charts

Charts plots the flow metrics `flai stats` computes, so the numbers are the same in both places. You choose a chart from Charts in the nav rail, which lists them in four groups: flow, usage, planning, and strategic ([The header](#the-header)). The page has no list of its own and takes its whole width. It reads from the top down:

1. The chart's title.
2. Any notice that replaces or qualifies the chart: the flai on the host is older than the dashboard, no item carries usage yet, no story has a forecast, or flai cannot read the conflict log or git.
3. The chart, with what belongs to it: the share bar of Time in State, the Longest waits table of Agent Waiting, and the figures a chart states, such as Strategic Cost's ratio, Coordination's shares, CoD by Order's saving, and the stories without a value.
4. **About this chart**, the chart's explainer.
5. The controls: the window and the controls the chart offers, with the summary strip at their end.
6. The table view, which you expand to read every number the chart draws.

Pick a window, and the controls the chart offers: an item type, an epic, on the usage charts over time and Forecast Error / Model what a bar or a point covers (an hour, a day, or a week), and on the forecast charts of planning a nature and a model. The strategic charts offer only the item type. The controls show while a chart loads, and when it fails to load, so you can change the window. Every chart follows the light or dark theme. The window sets what each chart shows: its time axis runs from the start of the window to now, and the charts of items and their tables list only the items completed in the window. The window you pick last is kept in this browser, so every chart opens at it, after a reload or a visit to another page too.

About this chart explains the chart in three short paragraphs: what it draws, how each figure is measured, naming the transitions, timestamps, and `flai stats --json` fields behind it, and what a shape of the chart means. It names the item type and the bucket you chose. A link at its end opens the heading of [metrics.md](../../design/system/metrics.md) that defines the chart's figures; metrics.md is the source of those definitions, and the explainer follows it. The explainer is open on your first visit. Close or open it on one chart, and every chart shows it that way in this browser, after a reload too.

| Chart | Shows |
|-------|-------|
| Cycle Time | One point per item completed in the window, with the p50 and p85 lines |
| Burn-up | Scope against done, each day of the window |
| Cumulative Flow | How many items sit in each state, each day of the window |
| Time in State | A bar for each day of the window with items completed: how long, on average, those items spent in each state. Hover a bar for the items it averages; the table lists each item. Under it, the share of all their time per state |
| Throughput | Completions per week of the window, by nature, weeks with none included |
| Agent Waiting | How long agents waited on someone else, mostly you, each week of the window: on their threads and in review, stacked, with review split into its time on checks and on reviewers when stories are accepted through pull requests, and the mean wait per item as a dashed line. The table under it lists the ten longest waits ([Agent waiting](#agent-waiting)) |
| Merge Queue | How long stories waited in the merge queue each week of the window, from being enqueued to being landed or sent back: the mean and the longest, in hours. On the right-hand axis, the mean number of stories in the batches it built, and as bars the failed batches it bisected. Choose the window only. Says so when the flai on the host sends no merge queue, because the queue there has logged nothing yet or the flai is older ([Merge queue](../../design/system/metrics.md#merge-queue-s-0404)) |
| Tokens / Min | Tokens per minute of agent work over time, one line per model |
| Tokens / Day | The tokens spent on the items done each day, stacked by model, with the mean per day so far as a dashed line |
| Tokens per item | The tokens an item took on average, over time: a line for epics, for stories, and for tasks |
| Tokens / $ | How many tokens a dollar bought, over time, one line per model |
| $ / Day | The dollars spent on the items done each day, stacked by model, with the mean per day so far as a dashed line. What the strategic agents spent on them, the planner, the orchestrator, and the analyzer, is stacked on top as its own series, strategic, estimated; the mean leaves it out |
| $ / Work Type | What an item cost on average, over time: a line for epics, for stories, and for tasks |
| $ / Item | What each item completed in the window cost, stacked by model; an asterisk marks an item whose cost is estimated in part. What the strategic agents spent on an item, the planner, the orchestrator, and the analyzer, is stacked on top as its own series, strategic, estimated, and apart from the item's total. An item only the strategic agents spent on shows that alone |
| Avg. Time / Model | The minutes of agent work an item of the chosen type took on average, over time, one line per model |
| Avg. Cost / Model | What an item of the chosen type cost on average, over time, one line per model |
| Forecast Accuracy | One point per story done in the window with a forecast: how much longer or shorter it took than forecast. Triangles show the same against your estimate, where a story has one. The p50 and p85 lines sit either side of zero |
| Delivery Accuracy | One point per story done in the window with a forecast delivery date: how many days after that date it was done, early below zero, with the p50 and p85 lines. A line on the right-hand axis shows the share delivered on time each week |
| Forecast Error / Model | How far forecasts were off, over time, one line per model: the p50 of the forecast errors of the stories each model's agents did in each hour, day, or week |
| Parallelism | Each day of the window: how many items were in progress at its end, and how many stories in ready were held then, beside a dashed line at the board's in-progress limit. Stories held while in progress sits below the limit are work the claims kept from starting, not the limit |
| Hold Time | The hours stories in ready were held in each week of the window, stacked by reason: overlap, a claim that overlaps a story in progress; after, a story it waits for is not done; empty claim, its touches or the in-progress story's are empty. Each hold counts under one reason, so a bar's height is the time held |
| Coordination | Whether the agents of open stories settle their overlaps early: each week of the window, a bar of the conflicts found while both stories were open beside a bar of those met late, at sync and at acceptance, and the conversations between stories opened, agreed, escalated, and closed with their story. Beside it, the stories started on a share and the hold time they saved ([Coordination](#coordination)) |
| Touches Drift | One bar per story done in the window: the files it changed outside its touches, which the planner or its agent should have declared, stacked on the touches it never changed, which they should have dropped. Hover a bar for the paths. A line on the right-hand axis shows each week's share of stories whose touches were exact |
| Strategic Cost | A bar for each day of the window: what the planner, the orchestrator, and the analyzer spent that day, stacked. A line for what the agents spent per item completed that day, and a dashed line for what the strategic agents spent per item over the window, with its share of the agents' cost per item |
| Strategic Use | A bar for each day of the window: the hours the planner, the orchestrator, and the analyzer worked that day, stacked. Two lines for the items completed that day: their mean cycle time and their mean waiting time, in hours |

The usage charts read what agents spent on each item: its tokens and cost, which flai records from the logs of the agents `flai serve` starts (see [Tokens and cost](flai.md#tokens-and-cost)). Each model keeps its colour and its mark on every chart. Until items carry usage, they say so and how to fill in stories worked before. An item's page shows its usage, per model, beside its other fields.

### Spend over time

Every usage chart but $ / Item lays what was spent out over time.

- **What a point covers.** Choose **per** hour, day, or week. Tokens / Day and $ / Day are renamed for it: $ / Hour, Tokens / Week. An hour is offered over a window of 31 days or less. Days and weeks are UTC days and weeks, as flai counts them, each drawn on its own date and named by it; an hour is labelled in your time zone.
- **When an item counts.** Everything an item's agents spent counts in the hour, day, or week in which the item was accepted as done, not when the tokens were used. An item that is not done, or was cancelled, is not counted.
- **Which items.** The charts cover the item type you choose, except Tokens per item and $ / Work Type, which put epics, stories, and tasks side by side. A story's spend includes its tasks', and an epic's its stories', so the types are never added together.
- **Per model.** A model's line averages the items of the chosen type it worked on. An item two models worked on counts whole for each, with all its agent time, since flai does not split agent time by model. Click a name in the legend to hide a line, such as epics, which dwarf the others.
- **The dashed lines.** On Tokens / Day and $ / Day, the mean per day so far, from the first day with spend, days with nothing done counted. On the other charts per model with more than one model, all models together.
- **The summary strip**, at the end of the controls, says what the items done in the window spent, and what that comes to per item, per minute of agent work, and per dollar.

The table view lists each hour, day, or week in which items were done, newest first, with every number the charts draw.

These charts need a flai on the host as new as the dashboard. With an older one they say so: upgrade it with `flai self-upgrade`.

### Planning

Forecast Accuracy, Delivery Accuracy, and Forecast Error / Model compare the planner's forecasts, and your estimates, with what happened. They show stories that are done; a cancelled story is left out.

- **A positive error is late.** On Forecast Accuracy a point above zero took longer than its forecast duration, and one below took less. On Delivery Accuracy a point above zero was done after its forecast date, and one below was early.
- **The p50 and p85 lines.** Half the errors fall between the two p50 lines, and 85% between the two p85 lines. The closer together the lines, the better the forecasts. With no filter or one, they are the figures `flai stats --json` gives; with both, they are worked out from the points shown.
- **The on-time share.** On Delivery Accuracy, each week's point on the right-hand axis is the share of that week's stories with a forecast date that were done on or before it. A week with no such story has no point, and the line breaks there; it does not drop to 0. Weeks start on Monday, UTC, and are drawn on their Monday's date.
- **Filters.** Narrow a chart to one nature or one model. A story's model is its agent's model; a story without one is listed as `(none)`. Forecast Error / Model shows every model, so it offers only the nature.
- **Which stories appear.** Only stories with a forecast. The planner sets one when it plans a story, or set it yourself with `flai edit <story> --forecast-duration 6h --forecast-delivery 2026-10-09T17:00:00Z` ([Drafts, cost of delay, and forecasts](flai.md#drafts-cost-of-delay-and-forecasts)). Until a story done in the window has one, the chart says so.

The summary strip, at the end of the controls, gives how many stories had a forecast, and the p50 and p85 of the forecast and delivery errors. The table view lists each story, newest first: when it was completed, its nature and model, its forecast, how long it took, and its forecast, delivery, and estimate errors.

Parallelism, Hold Time, and Touches Drift show what claims cost. A story in ready is held while its claim overlaps a story in progress, while a story it waits for is not done, or while either claim is empty (see [Held stories](#board) on the board, and [Touches](flai.md#touches)).

- **Held beside the limit.** On Parallelism, held stories on a day when in progress is below the limit could have started but for their claims. Held stories with in progress at the limit would have waited anyway.
- **Hold hours by reason.** On Hold Time, a hold with more than one reason counts under the one the board shows, `held (after)` before `held (overlap)`, so a tall overlap or empty-claim bar points at touches drawn too wide or left out.
- **Drift.** On Touches Drift, a file changed outside the touches is one the planner or the story's agent should have declared; a touch never changed is one they should have dropped. A story with neither had exact touches. A week with no story done has no share, and the line breaks there.
- **Controls.** These charts offer only the window. They need a flai on the host as new as the dashboard, and say so when it is older. Touches Drift is drawn from git: when flai cannot read it, the chart says so.

The table views list Parallelism by day, with in progress, held, and the limit; Hold Time by week, with each reason's time and the total; and Touches Drift by story, with both counts, then by week, with the stories done, how many had exact touches, and the share.

#### Coordination

Coordination, listed between Hold Time and Touches Drift, shows whether the agents of stories open at once settle their overlaps while both are open, rather than when one is accepted, which is what epic E-0018 set out to achieve ([Messages between stories](flai.md#messages-between-stories)). It reads `coordination` in `flai stats --json`; [metrics.md](../../design/system/metrics.md#coordination) defines each value.

- **Found early against met late.** Each week has a bar of the conflicts a trial merge found while both stories were open, and a bar of those a rebase met late: at sync, after another story was accepted, plus at acceptance.
- **How the conversations went.** Each week also shows the conversations between stories that were opened, agreed, escalated to you, and closed with their story. Agreed counts a share or a clean merge at sync; agreement in a reply, or a narrowed claim, is not recorded, so it is a floor.
- **Shares.** Beside the chart, the stories started on a share over the window, and the hold time they saved: how long each would have waited for the story holding it ([Asked about a hold](flai.md#asked-about-a-hold)).
- **How to read it.** More conflicts found early than met late, and more conversations agreed than escalated, means the coordination of E-0018 works: stories find their conflicts while both are open and settle them between themselves.
- **The conflict log.** The conflicts come from `.flai-cache/conflicts.jsonl` on the host, which flai writes from the release that brought it, so earlier weeks show none found or met. When flai cannot read the log, the chart says so and draws no conflicts; `flai stats` says why ([Flow metrics](flai.md#flow-metrics)).
- **Controls.** It offers only the window, as the other charts of claims do. With a flai on the host older than the dashboard, it says to upgrade it.

### Cost of delay

The planning group also has three charts of cost of delay: what a week of waiting for a story costs, the value you set on it ([Drafts, cost of delay, and forecasts](flai.md#drafts-cost-of-delay-and-forecasts)). CoD Outstanding and CoD Incurred cover the item type you choose, stories or epics; CoD by Order always projects the ready stories, and offers only the window.

| Chart | Shows |
|-------|-------|
| CoD Outstanding | For each day of the window, the cost of delay per week of the stories in each column at the end of the day, stacked: backlog, ready, in progress, and review, in the colours of Cumulative Flow |
| CoD Incurred | What waiting cost in each week of the window, as bars, with the mean per week so far as a dashed line |
| CoD by Order | What the stories in ready now will cost by the time each is pulled, under three orders: the board's pull order, by cost of delay, and by WSJF. One step line per order, rising at each story's projected pull. The time axis runs from now to the last pull, not over the window |

- **Stories without a value.** Under each chart, how many stories in each column have no cost of delay now. They add nothing to the sums, so a column full of them looks cheaper than it is.
- **The saving.** Under CoD by Order, each order's total, and what ordering the ready column by cost of delay or by WSJF, whichever is cheaper, would save against the pull order. Each story costs its value for every week from now until it is projected to be pulled, with no more stories in progress at once than the in-progress limit allows. When the pull order is already the cheapest, it says so. To take the saving, reorder the column with `flai order --by cod` or `--by wsjf` ([Ordering by a policy](flai.md#ordering-by-a-policy)), or drag the cards on the board.
- **Stories left out.** CoD by Order prices only the ready stories with both a cost of delay and a forecast duration, so the three totals cover the same stories. It names the others.
- **No in-progress limit.** Without one, every story is pulled at once, and every order costs nothing.

The table view lists each day, each week, or each projected pull, with the numbers the chart draws. These charts need a flai on the host as new as the dashboard; with an older one the page says so.

### Agent waiting

Agent Waiting, under Flow, shows the time agents spent waiting on someone else, mostly on you: for an answer on a thread, and for a story in review to be accepted or sent back. The exact definitions are in [metrics.md](../../design/system/metrics.md#waiting).

- **A bar** is a week of the window, Monday to Sunday, UTC, drawn on its Monday's date. It stacks the hours the agents of the items done that week waited on threads and in review. A wait counts in the week its item was done, not in the week it happened.
- **A thread wait** runs from a thread's first entry to the first reply by someone else. A recommendation does not end it; your confirmation does. Only the part while the story was in progress counts. **A review wait** runs from the move to review to the move out of it.
- **On checks and on reviewers.** When the project accepts stories through pull requests, two more series show the parts of review spent waiting on the required checks and then on the reviewers, from the forge state flai serve keeps ([ADR-0168](../../design/adrs/0168-under-pull-request-acceptance-flai-stats-reports-a-checks-wait-and-a-reviewers.md)). A checks wait runs from each push until no required check is pending. A reviewers wait runs from then until the first approval or requested change. They stack on review, which then shows the rest of it, so a bar is still as tall as the time waited. The tooltip's review is the whole review. In `local` mode, or with no kept forge state, the two series do not appear.
- **The dashed line** is the mean wait per item done that week, threads and review together. A week with nothing done has no point, and the line breaks there.
- **Filters.** Choose the item type. There is no epic filter, since flai adds up the waits over every epic.
- **The table** under the chart lists the ten longest waits with a part in the window, longest first: the story, whether it waited on a thread, in review, on checks, or on reviewers, the thread (`—` for checks and reviewers), when the wait started and when it ended, or `open`, the hours of it inside the window, and who it waited for: the one who answered the thread or moved the story out of review, the required check that finished last, or the reviewer whose review ended the wait. It lists stories in any state but cancelled, so a story waiting now is in the table before any bar holds it. When nothing waited in the window, it says so.

### Strategic

The strategic charts set what the planner, the orchestrator, and the analyzer cost against the delivery they serve, so you can tell whether they pay for themselves. Each is one bar per UTC day of the window, the three agents stacked, each in its own colour. They count the item type you choose. The figures are those of `strategic_days` in `flai stats --json`, defined in [Strategic use per day](../../design/system/metrics.md#strategic-use-per-day).

- **Strategic Cost, in dollars.** The bars are what each strategic agent spent that day; hover one to see whether its cost is estimated in part. The solid line is what the agents spent on average per item completed that day, without the strategic agents; a day with nothing completed has no point.
- **The ratio.** The dashed line is what the strategic agents spent over the window, divided by the items completed in it. Its label gives that as a share of what the agents spent per item over the window, the `usage $… per item` that `flai stats` prints under "strategic agents in the window", beside their total and the items completed. At 10%, the strategic agents add a tenth to each item's cost. With nothing completed, or nothing spent by the agents, there is no ratio.
- **Strategic Use, in hours.** The bars are the hours each strategic agent worked that day. The two lines cover the items completed that day: their mean cycle time, and their mean waiting time, the time each waited on its threads while in progress plus its time in review ([Waiting](../../design/system/metrics.md#waiting)). A day with nothing completed has no point.
- **Reading Strategic Use.** The return on the strategic agents' time is the waiting and cycle-time lines falling while their bars rise.

Under Strategic Cost, a line states the ratio, or why there is none. About this chart says how to read each of the two. The summary strip, at the end of the controls, gives what the strategic agents spent and worked over the window, and the items completed in it. The table view lists each day: each agent's cost or hours, the total, the items completed, and the figures of the lines. When the flai on the host is too old to send these figures, the page says so; upgrade it with `flai self-upgrade`.

## Docs

Docs shows every markdown file under `design/`, `docs/`, and `wip/` in a collapsible tree. The rail's Documentation group opens it at three places: System at `design/system/README.md`, Conventions at `design/conventions/README.md`, and Issues at `design/issues/summary.md`. The analyzer's reports are there, in `design/analysis/`, one per run beside the folder's README, as experiments' results are in `design/experiments/`; a report's front matter names its focus and the window its figures cover, `from` and `to`. A document renders with its Mermaid diagrams, highlighted code, task-list checkboxes, and heading anchors; links between documents open in the explorer, and a link to a heading lands on it. A wide table scrolls sideways inside the document. The front matter is shown in a panel above the text, and the threads on the document below it. When a story in progress or in review says it touches the document, the page names the story: an edit here may collide with that story's branch. When the dashboard can write, an **Edit** link opens the editor.

The tree sits at the left of the document, in a box that scrolls on its own. Under it, **On this page** lists the open document's headings, each nested under the heading it belongs to: click one to scroll to it. The heading you are reading is highlighted. On a wide window the tree and the headings are a column beside the document, and stay in view as you scroll; on a narrow one they are above it, the headings in a row you can scroll sideways.

### Editing documents

Every document page and every work item page has an Edit link when the dashboard can write. The editor shows the markdown on the left and, on the right, a preview drawn by the same renderer as the explorer, diagrams and code highlighting included.

What you can change depends on the file, and flai decides it, not the dashboard. Design and docs files are yours entirely: body and front matter, with the front matter checked when you save. For work items, narratives, and the board, flai owns the front matter, because it is the item's state; it is shown read-only and you edit the body. Move and block items from the board or with `flai`; change what a story or an epic says about itself from its own page, as described next. Generated files (`wip/agents/index.md`, `design/issues/summary.md`), threads, issues, and anything in the archive are not editable here, and the page says why.

Saving does three things. flai checks the repository with your change in place, and if the change introduces any finding, the save is refused, the findings are listed, the file is left as it was, and your text stays in the editor. If someone else changed the document after you opened it, an agent or a colleague, you get a conflict instead: the page shows what is there now against what you are saving, and you choose to load the current version, which discards your edits, or to save yours over it. Otherwise the file is saved and committed on its own, with you as the author, the line you typed under "What changed" as the subject, and a trailer naming the dashboard. For design and docs files the `updated` date is set to today unless you set it yourself. A project can turn the commit off with `dashboard.autocommit: false` in `system-flow.yaml`; the edit is then saved and left for you to commit.

When a story in progress or in review says it touches the document, the editor warns you, as the document page does, and asks you to tick a box before the first save: your edit lands on `main`, and that story's branch will meet it at its next sync. Click anywhere in the body and "Open a thread on" names the heading you are under, so a question can be attached to the section it is about.

**From a link.** A pencil beside a link to a markdown file, or beside a file name in code such as `design/system/overview.md:42`, opens the editor over the page you are on, with the cursor on that line. Following the link itself still shows the rendered document. Pages, threads, the analyzer's reports, and refusals on the board, an item page, and the review page have them. Files the editor cannot edit have none: the archive, threads, issues, and `wip/agents/index.md`. Close the editor to go back to the page. A link of the form `/edit/<path>?line=42` opens the editor in a tab of its own.

**On a story's branch.** On the page of a story in progress or in review, a design or docs file named in its body or its threads opens on the story's branch: the editor says `saves on branch story/S-nnnn`, and the save is committed there, in the story's worktree, not on `main` ([ADR-0145](../../design/adrs/0145-a-dashboard-save-may-name-a-story-and-is-then-committed-on-that-story-s-branch.md)). A file under `wip/` always opens on `main`, where it lives.

**Lint as you type.** Under the text the editor lists the markdown lint findings, each with its line, its rule, and the message, when it opens and after each pause in typing. Click one to put the cursor on its line. They are flai's subset of the markdownlint rules with the project's configuration; the project's markdown tier still decides when a story is accepted. With a flai on the host too old to lint, the list is not shown.

The editor asks before you leave, close it, or open another document with unsaved changes. Creating, renaming, and deleting documents is not something it does. Commits made here are not pushed; they reach the remote with the next publish, as a story accepted from the board does.

### Editing a story or an epic

A story's or an epic's page has an **edit…** button while the item is open (not done, cancelled, or archived). It turns the page into a form: the title, the nature, the tags, the topics beside them (one word each, comma separated; emptied, they are removed), and for a story what it touches, the stories it waits for (**Waits for**: it stays held in ready until each is done), its parent epic, and its agent (the harness, the model, the provider, and options; what you save replaces the story's agent, and emptying all four removes it), and for a story or an epic the **Cost of delay** panel as on the new item page (clearing all three inputs removes them and leaves the planner's value, to be worked out again), with the body below as Markdown and a preview. The line at the top names what is not yours to change there, because it is the item's state and flai's: the ID, the type, the status (move the card instead), the owner, and the dates. On a draft story a warning says it cannot be moved to ready until you finalize it, with a **Finalize** button that finalizes it there and keeps what you have typed.

Save sends only what you changed, and flai on the host does the rest in one step. A new title is kept in step everywhere it appears: the file's name, the heading, the line in the parent's list, the story's narrative, and links to the old file name in other documents. A new parent must be an open epic; the story leaves the old epic's list and joins the new one. The repository is checked with your change in place, exactly as for a document: if the check finds anything the change introduces, say a story in progress left without acceptance criteria, nothing is changed, the findings are listed, and your text stays in the form. If an agent or someone else changed the item after you opened it, you get the same choice as for a document: load the current version, or save yours over it, and only the fields you changed are written over. Everything is committed in one commit, as you, with the dashboard's trailer.

An agent working on the story is told. Its `inbox` reports the edit with what changed (title, nature, tags, topics, touches, after, parent, cost of delay, goal, criteria, notes, or body), and one that is waiting hears within a second, so it reads the story again before it goes on.

Tasks are not edited here: they are the agent's to write. A task's page says so, and its body can still be opened as a document.

## ADRs

ADRs lists the architecture decisions with status, date, and which decisions supersede which; each opens in the explorer. When the dashboard can write, "+ new ADR" records a decision from there.

### Recording a decision

"+ new ADR" opens a form: a title (the decision, as a sentence), whether it is proposed or accepted, the existing decisions it supersedes or refines, chosen from the list, and the markdown, which starts from the sections of your project's `design/adrs/0000-template.md` with the preview beside it. You write no front matter. The number (one past the highest ADR in your checkout, on the main branch, in any story's worktree, or on any story branch, so two stories never take the same number), the file name, the front matter, the date, the row in `design/adrs/README.md`, `superseded_by` on a decision it replaces, and the commit (your git identity, with a flaiover trailer) are made by `flai adr new`. If `flai check` would report anything with the new record in place, nothing is created and your text stays in the form. Afterwards you land on the new ADR in the explorer, and the list shows it without a reload, with what it supersedes and refines.

A proposed ADR is a draft: it can still be edited, and its row on the ADRs page has an "accept" button, which sets it to accepted with today's date. An accepted ADR is immutable, here as everywhere: the editor will not open its body, and the way to change a decision is to record a new one that supersedes it. Nothing is pushed.

## Search

Search covers `design/` and `wip/` by default and `docs/` when you tick the box. Type an item ID, a title, or words from the body; each result shows where it is, its status, `Draft` after the status for a draft story, and a snippet around the match. It runs in flai on the host and matches item IDs, titles, tags, headings, and bodies, by prefix and with small typos forgiven; an item ID such as `S-0052` finds its item first. Results open in the explorer.

## Host

Host is the Updates page. Host in the rail's Configuration group links it, and so does the host flai badge in the top bar. It has three parts.

**Dashboard** shows what image and version the dashboard container is running, and, once the operator has enabled it (`flai serve enable dashboard`), Restart, Upgrade, and Stop. **Check for updates** always works and changes nothing. Restart and a successful upgrade stop the very container answering the page, so the page expects the connection to drop and shows "Reconnecting…" rather than an error, then says what came back once it does. An upgrade never touches the running container until the new image has proven itself healthy: if it does not, the page says so and nothing changed.

**Addresses** lists the host names the dashboard accepts besides its own address: a tunnel's or a domain's name that reaches it, such as `flow.example.com, board.example.net:8443`. The address a request comes to is always accepted. Through any other, without its name here, live updates fail and the page keeps showing "Reconnecting…". Once the operator has run `flai serve enable settings --all-projects` on the host, the area has a field and **Save**, which sets `dashboard.allowed_hosts` on the host. flai refuses an entry that is not a host name, such as one with a scheme, a path, or a space: the reason shows under the field and nothing is written. The running dashboard keeps the list it was started with, so after a save the area says to apply it with Restart in the Dashboard area, which needs `flai serve enable dashboard`. While the running dashboard's list differs from the saved one, the area names both. Without the settings action, it shows the list, read-only, and the command that allows changing it. A name here does not let anyone in: every request still needs the dashboard's token. Setting it up from a shell is in the [operators guide](../operators/index.md#the-tunnel-expectation).

**flai host** shows the flai host process (S-0106): its version and pid, and a row each for `serve` and the MCP servers. Each row gives the state, the version, and the restart count, and the MCP row lists the projects served. **Check for upgrade** always works and changes nothing. Once the operator has enabled it with `flai serve enable host`, each row has Start, Stop, and Restart, and one **Upgrade** installs the newest flai and has the host start serve and the MCP servers again from it. The dashboard reaches the host through serve, so it asks before stopping serve: only `flai host start serve` in a shell on the host brings serve back. A serve restart and an upgrade drop the page's connection, so the page shows "Reconnecting…" and then says what version came back. With no host running, the area says to start one with `flai host start`.

Each area has **Versions** beside its check, to go back to an earlier release or to any published one (S-0298). It always works and changes nothing: it lists the published releases, newest first. The Dashboard area marks the release running, the newest, and the configured tag. It marks **not deployable**, with flai's reason, a release with no digest list the host's flai can verify, such as one from before releases were signed: flai refuses to deploy it, so the page offers no way to. The flai host area gives each release's date and marks the installed one, the newest, and any below a served project's `flai.minimum`. It marks **not installable**, with flai's reason, a release whose signature the host's flai cannot verify, such as one from before releases were signed, which publishes no signature: flai refuses to install it, so the page offers no way to. Once the area's host action is on, each release other than the one running or installed, and not marked not deployable or not installable, has **Deploy**. It asks first, then deploys that release as Upgrade deploys the newest: the same "Reconnecting…" and the same report of what came back. For flai, the confirmation warns when the release is below a project's `flai.minimum`: `flai serve` leaves that project unserved, and flai's commands refuse it, until flai is at least the minimum again. A dashboard release deployed here applies to the container, not the configuration. It keeps running through restarts until the next Upgrade, which goes to the configured tag, or a start after a stop. The page has no control to pin a release; on the host, `flai config set dashboard.tag` does. The same lists and deploys from a shell are in the [update runbook](../operators/runbooks/update.md).

## Settings

The settings are five pages of the rail's Configuration group, which lists [Host](#host) after them. Together they show how flai on the host is set up for this project and, once the operator allows it, change it. Each page has its own address:

| Page | Address | Sections |
|------|---------|----------|
| [Features](#features) | `/settings/features` | Host actions, Checks |
| [Agents](#agents) | `/settings/agents` | Default agent, Providers, Planning again |
| [File System](#file-system) | `/settings/file-system` | Shared paths, Import folders |
| [Projects](#projects) | `/settings/projects` | Projects, Tokens |
| [Harnesses](#harnesses) | `/settings/harnesses` | Starting agents, and a block for each harness |

`/settings` opens Features. A link from when the settings were one page, such as `/settings#projects`, opens the page its anchor names, or Features for any other anchor, and keeps the project it names. Any other address below `/settings` is not found, and the error names the five pages.

Nothing changes until the operator runs `flai serve enable settings` on the host. The top of every page says whether settings can be changed here, for this project or for every project, or are read-only, and gives the command that allows them. It also says that whoever holds the dashboard token can then change what flai runs on the host, as the operator, and that only a shell turns it off: `flai serve disable settings`. Each section that cannot be changed says so and names the command. Once `settings` is on, this project's own settings can be changed: its host actions, its default agent, its shared paths, and its MCP token. The settings the host keeps for every project, its checks, providers, import folders, dashboard token, and the command and harnesses that start agents, also need `flai serve enable settings --all-projects`. The same switch lets the [Planner](#planner), [Orchestrator](#orchestrator), and [Analyzer](#analyzer) pages change their agents' settings. When flai on the host keeps no settings for the project, each page says so.

Commands are written one argument a line, and are run exactly as written, with no shell in between. Each section says whether a change was saved, or why flai refused it, and the page then reads the settings again.

### Features

Host actions lists every host action flai has, `plan` and `orchestrate` among them, each with a checkbox and what it does. Tick one to turn it on for this project, or clear it to turn it off. `settings` itself is shown and changed only in a shell. An action on for every project is shown as such and not changed here.

Checks are the commands run on a story in review on this host, each with its name. **Save** changes one, **Remove** deletes it, and **Add check** takes a name and a command. With none set on the host, the section says so and names the manifest's checks, which are used instead. The time limit, in minutes, bounds one run of every check together. Checks are kept for every project on the host.

### Agents

Default agent is the agent this project's new stories get: its harness, its model, its provider, and its options, one `key=value` a line. **Save default agent** writes it to `system-flow.yaml` and commits it. A field left empty leaves the default agent without it.

Providers lists the gateways on this host, such as a LiteLLM proxy or OpenRouter, that an agent's model calls can go to when it names one ([Where its model's calls go](flai.md#where-its-models-calls-go-a-provider)). Each shows the API it speaks (`anthropic-messages`, `openai-chat`, or `openai-responses`), its base URL, `key_env`, the name of the variable on the host that holds its key, and its models, one `alias=name` a line. Below them a line says `key set`, or that the variable is not set in `flai serve`'s environment and an agent that names the provider is refused at start. The key itself is never shown or asked for: set the variable where `flai serve` runs. **Save** writes the fields you typed; models left empty keep the aliases the provider has. **Remove** deletes the provider. To add one, give its name, API, base URL, and `key_env`, and models if you need them. The default agent's fields have a provider too. Providers are kept for every project on the host, so changing them needs `flai serve enable settings --all-projects`; until then the section is read-only and names that command.

Planning again, for this project, is read-only here: `replan` and `schedule`, like the other planning settings, are edited on the [Planner](#planner) page, which the section links to. It shows whether edits to a planned story start the planner, which is whether the `plan` host action is on, since replanning and the schedule act only while it is; what happens when work ahead completes or the pull order changes, `never`, `deterministic`, or `agent`, marked when it is the default, with what that means; and the schedule, a cron expression read in UTC, with its next run in your time zone, or none. A value flai cannot read is shown with its error. [Running the planner](flai.md#running-the-planner) says what each does.

### File System

Shared paths lists the patterns in `claims.shared` in the project's `system-flow.yaml`: an overlap of two stories' claims inside one holds no ready story ([Shared paths](flai.md#shared-paths)). With the `settings` host action on, you can add a pattern or remove one. In a pattern, `*` matches within a folder, `**` across folders, and a plain path covers everything under it. flai commits each change. When flai refuses a pattern, because it is not valid or is listed already, the reason is shown beside the field. Check, which works even while settings are read-only, takes a path or a story's ID. For a path, or for each entry of the story's claim, it says which shared pattern the entry lies inside, or that none does. The section is missing while the host's flai is older than the one that brought shared paths.

Import folders lists the folders on the host whose git repositories the board offers to import, and whose repositories with a `system-flow.yaml` `flai serve` serves. **Add folder** adds one and **Remove** takes one away; removing one stops serving every project below it. Import folders are kept for every project on the host.

### Projects

Projects lists every project `flai serve` serves, with its key, name, and folder, and whether it is connected and since when, the last error, or why it cannot be served. A project below an import folder is served as soon as `flai serve` next looks, and joins the switcher without a reload. One it does not serve is listed under "Below the import folders or the folder flai serve was started in, not served" with why: you removed it, its `system-flow.yaml` has no key, or its key is served already for another folder. **Remove** asks first, then stops serving the project, as `flai serve project remove` does, and the switcher drops it without a reload; none of its files is touched. A registered project is unregistered, and the page says how to serve it again on the host. A project served because it is below an import folder, or below the folder `flai serve` was started in, says which folder serves it; removing it puts it on `flai serve`'s list of removed projects, and it moves to the not-served list as removed (S-0123). **Serve** on a removed project takes it off that list, as `flai serve project add` does, and the switcher gains it without a reload. On any other project not served, Serve asks `flai serve` to register it and shows the answer; for no key or a key served already that is a refusal saying what to fix. Removing a folder under Import folders, on the [File System](#file-system) page, stops every project below it at once. Serve and Remove need `flai serve enable settings` for the project served or removed, and each says so, with the command, where it is off.

Tokens shows the address of the MCP server for agents over HTTP, or that it is not running. **Rotate the MCP token** makes a new one and shows it once: every agent connected over HTTP must be given it. **Rotate the dashboard token** keeps you logged in and shows the new login link once. Everyone else, and every agent using the token, is logged out. The dashboard token is kept for every project on the host.

### Harnesses

Starting agents is how `flai serve` starts a story's agent on this host. The name is the prefix its agents work as, such as `agent-S-0001`. The command starts the agent of a story that names no harness, one argument a line, with `{story}`, `{root}`, `{model}`, and `{harness}` replaced; **Remove the command** clears it. Below it is a block for each harness with its program; its arguments, what its agent may do, one a line; its guard, `flai guard` or `none`; and `deny_protected`. Under the last two the block says what choosing them gives up, and that they change nothing on a harness that runs `flai guard` or asks before a protected write. The `command` harness, the command above, takes only the guard and `deny_protected`. A harness you have not set shows its defaults, and **Back to the defaults** drops what you set. The keys behind each field are in the [settings index](../operators/settings.md#flai-configuration).

## License

License, at the foot of the nav rail, shows the Bytepunx Shield License: the terms the dashboard, flai, and the system-flow repository are distributed under. The text is the `LICENSE.md` the dashboard's image carries beside its server, so it is the license of the very dashboard you are looking at, and the same text `flai license` prints and the repository keeps at its root. It is the one page that asks flai on the host for nothing: the file is the dashboard's own, not the project's. An image built without the file says so instead.

## Design system

Design system, System in the rail's Design group, shows the project's design system as it stands: the design tokens in the file the manifest's `design.tokens` names, `design/system/design-system.tokens.json` by default, and the guidance beside them, `design/system/design-system.md`. You see the effect of the tokens in the browser without building or deploying the project.

- **Problems**: what flai found wrong in the token file, by line, as `flai check` reports it.
- **Preview**: sample components styled by the tokens, light and dark side by side: buttons, inputs, a form, a card, an alert per state, a table, navigation, and headings over body text. They are generic components, not the project's own, so they show its tokens, not its markup. The dashboard's own theme never reaches them, and switching it changes nothing in the preview.
- **Tokens**: each colour as a swatch with its contrast against the ground, `color.background`, in both themes; the colour pairs flai checks, with the ratio each must reach and whether it does; the type scale as sample text; the font stacks as your browser resolves them; and spacing, radii, and shadows drawn to size.
- **Guidance**: the guidance document, rendered.

A project with no token file is told so: run `flai design init` to write the starter tokens and guidance. The page follows a change to the token file or the guidance made elsewhere, as the other pages follow theirs.

### Changing tokens

**Edit tokens**, beside the page's introduction, opens the tokens for editing next to the preview. Nothing is written until you save.

- **Colours**: a picker and a text field for each theme. Type `#rrggbb` or `oklch(L C H)`, such as `oklch(52% 0.18 255)`.
- **Lengths**, such as spacing and radii: a number and its unit, `px` or `rem`.
- **Type scale**: each step's font size, as a number and a unit, and its line height.
- **Font stacks**: the font names, separated by commas.
- **Shadows and aliases** are read-only here. An alias follows the token it names, so change that one.

The preview, the token sheet, and the contrast follow each change as you type. When a colour leaves a pair below its threshold, 4.5:1 for text or 3:1 for controls, the colour says so beside its field: which pair, in which theme, at what ratio. A value flai would not accept stays in its field, marked with the reason, and changes nothing until you correct it. A changed token is marked, with **Reset** to take it back. A change to a light value reaches dark too, unless dark has a value of its own. In a token file with no dark context, a colour has one value for both themes: change its light value, because a dark change cannot be saved there.

**Save** writes your changes into the token file, leaves the rest of it as it was, and commits it on main, as saving a document does ([Editing documents](#editing-documents)). flai on the host saves it with `flai file save` ([flai.md § Edit a file that is not markdown through flai](flai.md#edit-a-file-that-is-not-markdown-through-flai)). When the manifest names a design stylesheet, flai writes the stylesheet from the new tokens and commits it with the token file in the same commit ([flai.md § The generated stylesheet](flai.md#the-generated-stylesheet)). The line under "What changed" becomes the commit's subject. flai checks the file first. If it finds a problem, the save is refused, the findings are listed, and your changes stay. If someone changed the token file after you started, you choose, as in the document editor: load the current version, which keeps your changes where its tokens still take them, or save yours over it.

**Discard changes** goes back to the saved tokens. **Done** closes the editor once nothing is unsaved. The page asks before you leave it with unsaved changes. A change to the token file made elsewhere while you edit updates the saved tokens and keeps your changes.

**Unsaved changes are kept in this browser**, per project, as you make them. They are not on the host: another browser does not see them. Reload the page, or come back to it later, and they are restored over the saved tokens, with "Your unsaved changes, kept in this browser, are restored: save or discard them." A save or **Discard changes** drops them. A Design page open in another tab of this browser takes up your changes too, and the Components page shows them as tried tokens ([Tried and saved tokens](#tried-and-saved-tokens)).

**Saving needs a flai that reads and saves the token file.** With a flai on the host older than the dashboard, the banner says what it lacks, and the page says "This flai cannot save the token file yet; upgrade flai on the host." You can still change tokens and see the preview, but Save is off. Your changes stay in this browser until you discard them.

### Proposals

An agent proposes a change to the design system on its story's branch. **Proposals**, below the guidance, lists each open story whose branch changes the token file, the guidance, or both: its ID, linked to the story, its title, its state, and what it changes. `flai design proposals` lists the same ([flai.md § Proposals](flai.md#proposals)).

**Compare** shows the proposal beside the main branch's design system:

- **Preview**: the sample components styled by main's tokens and by the proposal's, side by side, each in light and dark.
- **Tokens that differ**: per theme, each token the proposal changes, adds, or removes, with main's value and the proposal's. A token differs when it renders differently, so an alias replaced by the same value is no difference.

Problems in the proposal's token file are listed above the preview. **Close comparison**, or **Compare** again, closes it. Compare is off while you edit tokens. If flai cannot read the story's design system, such as when its worktree is gone, the comparison says why.

To accept a proposal, accept its story on the board. Its branch is merged, the page shows its tokens as the project's, and it leaves the list. A cancelled story leaves the list too.

### Requesting a change

**Request a change** asks an agent for a change to the design system. Choose the epic, the nature, a title, and write what to change and why in your own words. **Request the change** makes a story under that epic and moves it to ready for a story's agent to pull ([ADR-0151](../../design/adrs/0151-the-operator-and-agents-change-a-design-system-by-edits-on-the-design-page.md)). The story has topic `design`, your words as its goal, and three criteria: the tokens changed as asked, the generated stylesheet in step with them, and the UI's tests passing. It claims the token file, and the stylesheet when the manifest names one, so a second request is held on the board until the first is accepted. Once its agent changes the tokens on its branch, the story shows under **Proposals**.

If flai refuses the story, your words stay and its findings are listed. If the story is made but cannot move to ready, the page says so and the story stays in backlog: move it from the board.

## Components

Components, the second page of the Design group, lists the project's Svelte components and renders the one you pick, live. The component renders from its sub-project's own Vite dev server, with the sub-project's own config, plugins, and stylesheet, so it looks as it does in the app. `flai serve` on the host runs that dev server for the page ([Render a component live](flai.md#render-a-component-live) in flai.md).

**The list.** At the left is each sub-project of the manifest, with its kind, and its components grouped by folder. The filter box matches a component's name or path. A sub-project that is not SvelteKit or Svelte is listed with its kind and no components. A component whose props flai could not read is still listed, marked `props unreadable`: hover over it for why. It can still be framed.

**The frame.** Pick a component to frame it. The address keeps your pick as `?p=<sub-project>&c=<path>`, so a reload or a copied link frames it again. The component renders with the props the props panel gives it. The frame is as tall as what the component draws, with a small margin around it, so the page scrolls to show all of it. A component wider than the frame scrolls sideways inside the frame. Only the component you picked last shows in the frame, and only its errors: a late answer about the one you picked before is ignored. Above it:

| Control | Choices |
|---------|---------|
| Theme | **Light**, **Dark**, or **Both**, side by side; it starts at the dashboard's theme |
| Width | **375**, **768**, or **1280** pixels, or **Fill**, the panel's width |

**Props.** The props panel lists the picked component's props, each with its name, its type and its default as the source writes them, and `required` when it must have a value. On a wide screen the panel sits right of the frames; on a narrower one, below them. Each change renders in the frame as you make it. The component is not mounted again, so it keeps its own state, such as an open menu or typed text. Each prop has a control by its kind:

| Kind | Control |
|------|---------|
| `string` | A text box. Emptied, it gives the empty string |
| `number` | A number box |
| `boolean` | A checkbox. It shows `not set` until you click it |
| A union, such as `'sm' \| 'md' \| 'lg'` | A list of its values, with **(not set)** first |
| An array, an object, or a type flai could not resolve | A JSON box. Any valid JSON is taken |
| `Snippet` | A text box. The component renders your text as the snippet's content, as text, not HTML |
| A function | No control. It shows **Logged**: its calls are listed under **Events** |

A control starts at the prop's default when the source writes the default as a plain value, such as `'md'`, `3`, `false`, or `[]`. Otherwise a JSON box starts with the empty shape of the prop's type, such as `[{"id": "", "title": "", "tasks": 0}]` for a list of stories. Fill in its values; the shape is sent as it stands, so the component gets what the box shows. Where flai cannot work out the type, the box starts at `[]`, `{}`, or `null`. Any other control starts empty, and the component's own default applies. A prop without a value is not sent, so it takes its default. **Reset** puts every prop back to its default or its empty shape.

Text that does not parse, such as JSON with a missing quote, stays as you typed it, with its error and `not sent` under it. While any prop's text does not parse, nothing is sent: the frame keeps the props it had, and the address is not changed.

A required prop without a value says `Required: give it a value`. Until each required prop has one, the frame's place says `<path> needs a value for <props>.` Svelte does not stop a component that lacks a required prop, so this is the page's own notice, not an error from the component.

**Events.** Under the props, **Events** lists the calls of the component's function props, newest first, up to the last 50. Each shows the time to the millisecond, the prop's name, and its arguments as JSON. An event argument shows as its type, and its detail for a custom event; an element as its tag, such as `<button>`. **Clear** empties the list. Picking another component empties it too.

The frame does not send the component's requests to the app's own server, such as a fetch of `/api/...`, which the dev server cannot answer. It answers each itself and lists it under **Events**. A fetch shows as `fetch`, with its method, its path, and the status it got. One that the **Answers** panel answers gets that answer (**Answers**, under [Examples](#examples)); any other is answered 501 Not Implemented, with an error that names it. A component that shows a request's error shows that one. An `EventSource`, a live stream, shows as `EventSource`, with its path and 200 when an answer sent its messages, or 501 when it failed without connecting. A `WebSocket` shows as `WebSocket`, with its path, and always fails without connecting, closing as a refused one does. The dev server's own socket, which brings a saved change to the frame, is not affected. Loading the component's own files and assets is not affected.

**Links that carry props.** The address keeps the props' values too, as `&props=`, so a reload, or the link opened in another tab, renders the component with the same props. The address has no `props` when no prop has a value. Picking another component drops the address's props, your changes, and the events: it starts at its defaults.

**The dev server.** Below the controls the page says whether its sub-project's dev server runs. **Start** starts it. A start can take up to a minute while Vite gets ready. While it runs, the page gives its address and **Stop**. A change you save to a component reaches the frame without a reload.

The page asks the host again every minute while it is open, and when a document of the project changes. Each ask keeps the dev server running: the host stops one that nobody has looked at for its idle time, 30 minutes by default. Each ask also reads the list again, so a component you add or remove shows within a minute.

**The frame needs a browser on the host.** The dev server listens on 127.0.0.1 only, so the frame loads only in a browser on the host machine, such as one at `127.0.0.1:4242`. From another machine the list shows, but the frame cannot load.

**Before it works.** The operator turns on the `components` host action for the project. Starting a dev server runs the sub-project's code as the operator; the [operators guide](../operators/index.md#the-components-host-action-s-0438) says what that means. Listing the components needs no action.

What the page can say:

| The page says | What it means |
|---------------|---------------|
| The components host action is off, with `flai serve enable components` | The action is off for this project. Run the command on the host, in the project. The list still shows; the dev server cannot start |
| This flai cannot show the components yet; upgrade flai on the host | The flai on the host is older than the dashboard. The line under it says what it lacks |
| An error in red under **Start** | The dev server did not start. flai's message is shown whole, with Vite's last lines of output. A sub-project without `vite` installed is refused with the fix: `npm install` or `pnpm install` in its folder |
| Its last dev server failed, with its output | The last dev server failed to start or exited on its own. Its last lines stay until one starts again |
| `<path>` did not render, in the frame's place | The component failed to import or render. The harness's message follows, and `at <file>:<line>:<column>` when it knows where. The frame comes back when the component renders again, such as after you save a fix, or when you pick another |
| Its props could not be read | flai could not read the component's props. It still renders, with no props |
| `<path>` needs a value for `<props>`, in the frame's place | Those required props have no value. Give each one in the props panel and the frame comes back |
| An error and `not sent` under a prop | Its text does not parse. Nothing is sent until it does |
| `Invalid JSON: …; not sent.` or `Not an array; not sent.` under **Answers** | The answers' text is not a JSON array. The frame keeps the last answers that were, until it is |
| No sub-project is Svelte, so there is nothing to frame | The manifest has no sub-project of kind `sveltekit` or `svelte` |
| Could not read the components | The host did not answer. What was shown stays until the next ask |

### Examples

An example is a named set of a component's props, kept in the repository, so a component that needs data renders without typing its props each time. Each component's examples are one file under `design/components/` ([Component examples](flai.md#component-examples) in flai.md).

**The list.** Above the props panel, **Examples** lists the picked component's examples by name, in its file's order. A component without any says "No examples yet." Click an example to load it: its props replace the panel's, and the frame shows the theme and the width it was saved with, when it was saved with one. The address keeps the props, as for any change. The example you loaded is marked.

**Answers.** Under the examples, **Answers** holds what the frame answers the component's requests to the app's own server with, as a JSON array, so a component that asks the server for its data, such as a list of threads, shows answered. Each answer gives the request's `path`, and its `method`, `status`, `headers`, and `body`, or, for an `EventSource`, its `events`; [Component examples](flai.md#component-examples) in flai.md gives the format. The first answer whose path and method match a request answers it; a request none matches gets 501. Loading an example loads its answers, or empties the box when it has none, and picking another component empties it. A change reaches the frame at once, for the requests the component makes after it: pick the component again, or click what makes the request, to see a request it made before. Text that is not a JSON array shows its error and `not sent`, and the frame keeps the last answers that were. The address does not keep the answers.

**Save as example** keeps the panel's props as an example, with the answers when **Answers** has any. Give it a name and click **Save**. It starts from the name of the example you loaded; a name the component already has says "Saving replaces the example of this name.", and the save replaces that example where it is, while a new name goes last. The example keeps the theme when the page shows **Light** or **Dark**, not **Both**, and the width when you picked one, not **Fill**. A function prop is never kept. The form says where the example goes: on main, or on the story's branch when **Source** names a story ([A story's components](#a-storys-components)). The save is committed there. A component's first example creates its file.

| The page says | What it means |
|---------------|---------------|
| Saved the example "`<name>`" | It is saved and committed, and listed |
| An error with findings | flai refused the save, such as for an example whose props do not fit the component. Your name and the panel's props stay; fix them and save again |
| The file changed while it was saved | Someone saved the file meanwhile. Save again to add yours to it as it is now |

With a flai on the host older than examples, the page lists none and offers no **Save as example**: upgrade flai on the host with `flai self-upgrade`.

### Tried and saved tokens

The frames show the design tokens you change on the Design page ([Changing tokens](#changing-tokens)): first as you try them, then once you save them.

**Tried.** While the Design page has unsaved changes, the frames render with them, in each theme. Above the frames the page says "Showing tried tokens: the Design page's unsaved edits, kept in this browser." Keep the Design page open in another tab and the frames follow each change as you make it, without a reload or a save. The tried tokens are the saved tokens with your changes laid over them, as the generated stylesheet would carry them. They show only in the browser you edited in.

**Drop tried tokens** drops those unsaved changes: the frames go back to the saved tokens, and the Design page has none to restore. A Design page open in another tab drops them too.

**Saved.** A save on the Design page writes the token file and, when the manifest names one, the design stylesheet, in one commit. The sub-project's dev server sees the stylesheet change and updates it in the frame: the component shows the saved value without a reload, and keeps its props and its own state. The save also drops the tried tokens, so what you see is what was saved. The stylesheet reaches the frame only when the sub-project's code imports it, such as through the setup module the components key names ([Render a component live](flai.md#render-a-component-live) in flai.md).

If the page cannot read the saved tokens to lay your changes over, it says "Could not read the saved tokens to lay the tried ones over", with why, and keeps what it last read.

### A story's components

The page can show a story's worktree instead of the main checkout, so you watch a component change as its agent saves it, before anything is merged.

**The source.** Above the list, **Source** offers **Main** and each story in progress or in review. Pick a story and the list, the dev server, and the frame are its worktree's. A component the branch adds is listed too. The address keeps the story as `?story=S-nnnn`, so a reload or a copied link opens the same source. Pick **Main** to go back.

**What the branch changes.** Under the source, **Changed on its branch** lists the components the story's branch changes, each marked `changed`; click one to frame it. In the list they are marked `changed` too, first in their folder. The branch's changes are its commits: a file the agent has saved but not committed is not marked, though the frame shows it.

**The story's dev server.** **Start** starts a dev server in the story's worktree, beside main's, on a port of its own. The first start in a worktree can take a few minutes: a worktree has no installed dependencies, so the host installs them first from the sub-project's lock file. The host stops a story's dev server when the story leaves review: when it is accepted, cancelled, or sent back. It also stops one nobody has looked at, as it does main's. The host runs at most three dev servers per project at once, main's and the stories' together, unless the operator sets another limit.

| The page says | What it means |
|---------------|---------------|
| An error below **Source**, in place of the list | The host refused the story: it is no longer in progress or in review, or it has no worktree. Pick another source |
| Could not read what its branch changes | The branch's diff could not be read. The list still shows, without marks |
| An error under **Start** naming `npm ci`, `pnpm install`, or `yarn install` | Installing the worktree's dependencies failed. The package manager's last lines follow |
| An error under **Start** naming `components.max_servers` | As many dev servers as the host allows already run, each named. Stop one, or ask the operator to raise the limit |

**From review.** On a story's review page, under **Changes on the branch**, a story whose branch changes a component names them, with **Open in Components**. It opens this page with the story picked.

## Theme

flaiover uses the brand palette in a light and a dark theme. It follows your system preference by default; the theme button in the top bar cycles system, light, and dark, and the choice is remembered per browser. Both themes are checked for readable contrast, and the charts use palettes validated for colour-vision deficiency in each theme.

## Agents on other machines

Agents do not work through the dashboard. They use flai's MCP server on the host: over stdio as `.mcp.json` sets up, or over HTTP after `flai mcp start` for an agent that cannot start a process there; setting that up is in [the flai guide](flai.md). What an agent does shows up here like anyone's work. The dashboard's own `/mcp` address, which served MCP until flaiover 0.22, now answers 410 and says so.
