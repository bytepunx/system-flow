# Architecture Decision Records

One decision per file, numbered in order of acceptance. An ADR is never edited after acceptance except to record `superseded_by`, or its `topics` with `flai adr topics`. To change a decision, write a new ADR that supersedes the old one in whole or in part, then update the living documents in `design/system` and `design/tech`.

Record a new one with `flai adr new "<title>"`, or from the dashboard's ADRs page: the number, the file, its front matter, and the row below are made for you, from the sections in `0000-template.md`. `flai check` warns when a file has no row here or a row has no file.

| ADR | Title | Status |
|-----|-------|--------|
| [0001](0001-markdown-and-adrs.md) | Markdown is the only documentation format and decisions are recorded as ADRs | accepted |
| [0002](0002-repository-layout.md) | Three top-level documentation folders: design, docs, wip | accepted |
| [0003](0003-work-item-hierarchy.md) | Epics, stories, tasks, each with a nature | accepted |
| [0004](0004-workflow-states-and-transitions.md) | One state machine, transitions recorded in front matter, blocked is a flag | accepted |
| [0005](0005-template-as-separate-repository.md) | The template is a separate git repository rendered by flai | accepted |
| [0006](0006-go-for-the-cli.md) | Go for the flai CLI | accepted |
| [0007](0007-sveltekit-spa-with-node-adapter-in-docker.md) | SvelteKit SPA with node adapter, shipped as a Docker image | accepted, superseded in part by 0016 |
| [0008](0008-cli-config-location.md) | CLI configuration lives in ~/.flai/config.json | accepted |
| [0009](0009-agent-narrative-per-story.md) | One agent narrative per story in wip/agents | accepted |
| [0010](0010-shell-out-to-git-and-docker.md) | flai shells out to git and docker | accepted |
| [0011](0011-project-manifest.md) | system-flow.yaml marks a project and owns folder names | accepted |
| [0012](0012-flat-kanban-folders-and-archive.md) | Kanban items are flat per type, hierarchy by parent key, archive on completion | accepted |
| [0013](0013-conventions-folder.md) | design/conventions holds agent norms, one file per topic | accepted |
| [0014](0014-design-issues.md) | design/issues records recurring friction with counts and cost | accepted |
| [0015](0015-template-lock-file.md) | system-flow.lock.yaml records what the template rendered so upgrade can tell edits from baseline | accepted |
| [0016](0016-dashboard-delegates-to-flai.md) | flaiover delegates writes and metrics to the flai binary | accepted, supersedes 0007 in part |
| [0017](0017-four-digit-ids.md) | Work item IDs are zero-padded to four digits | accepted, refines 0003 |
| [0018](0018-dashboard-token.md) | flaiover authenticates every request with a per-project token | accepted, refines 0016 |
| [0019](0019-story-branches-and-touches.md) | Agents work each story on a branch in a worktree while wip stays on main | accepted, refined by 0022 |
| [0020](0020-files-plus-mcp.md) | The designer and agents communicate through files, served to agents by flai mcp | accepted, refines 0009 and 0016 |
| [0021](0021-story-ready-without-tasks.md) | A story is ready without tasks; the pulling agent writes them and review requires them | accepted |
| [0022](0022-repository-mounted-at-its-host-path.md) | The dashboard mounts the repository at its host path; relative-path worktrees are an opt-in | accepted, refines 0019, superseded by 0031 |
| [0023](0023-documents-are-saved-through-flai.md) | Documents edited in the dashboard are saved through flai | accepted, refines 0016 |
| [0024](0024-mcp-over-http-and-project-identity.md) | flaiover serves MCP over HTTP, and every response names its project | accepted, refines 0018 and 0020, superseded by 0030 |
| [0025](0025-research-is-accepted-without-a-release.md) | A research story is accepted and pushed without a release; an experiment stays on its branch | accepted |
| [0026](0026-the-dashboard-may-push-with-a-key-the-operator-gives-it.md) | The dashboard container may push acceptances with a key the operator gives it | accepted, refines 0018, superseded by 0031 |
| [0027](0027-git-hooks-config-and-info-are-read-only-in-the-dashboard-container.md) | Git hooks, config, and info are read-only in the dashboard container | accepted, refines 0018, 0022 and 0026, superseded by 0031 |
| [0028](0028-cancelling-an-item-cancels-everything-open-under-it.md) | Cancelling an item cancels everything open under it | accepted, refines 0004 |
| [0029](0029-the-dashboard-reaches-a-project-only-through-flai-on-the-host.md) | The dashboard reaches a project only through flai on the host | accepted, refines 0007, 0016, 0018 and 0024 |
| [0030](0030-mcp-is-served-by-flai-on-the-host-over-stdio-and-http-and-the-dashboard-s-api.md) | MCP is served by flai on the host, over stdio and HTTP, and the dashboard's API still names its project | accepted, supersedes 0024, refines 0020 and 0029 |
| [0031](0031-the-dashboard-s-container-holds-nothing-of-the-project-a-port-and-two-secrets.md) | The dashboard's container holds nothing of the project: a port and two secrets | accepted, supersedes 0022, 0026 and 0027, refines 0018 and 0029 |
| [0032](0032-accepting-a-story-merges-it-publishing-is-a-deliberate-batched-step-over.md) | Accepting a story merges it; publishing is a deliberate, batched step over everything accumulated | accepted, refines 0019, 0025 and 0031 |
| [0033](0033-one-login-token-and-one-agent-credential-per-user-serve-every-project.md) | One login token and one agent credential per user serve every project | accepted, refines 0018, 0029 and 0031 |
| [0034](0034-flai-serve-keeps-each-served-project-s-http-mcp-server-running.md) | flai serve keeps each served project's HTTP MCP server running | accepted, refines 0030, superseded by 0040 |
| [0035](0035-repositories-under-folders-the-operator-names-can-be-imported-from-the-board.md) | Repositories under folders the operator names can be imported from the board | accepted, refines 0029 |
| [0036](0036-a-folder-that-is-not-a-project-is-served-whole-by-flai-mcp-and-flai-dashboard.md) | A folder that is not a project is served whole by flai mcp and flai dashboard | accepted, refines 0030 |
| [0037](0037-a-story-carries-its-agent-copied-from-the-project-s-default-when-it-is-made.md) | A story carries its agent, copied from the project's default when it is made | accepted |
| [0038](0038-flai-serve-starts-a-story-s-own-agent-through-an-adapter-with-what-the-operator.md) | flai serve starts a story's own agent through an adapter, with what the operator allows | accepted, refines 0029 and 0037, superseded in part by 0041 |
| [0039](0039-a-settings-host-action-turned-on-only-in-a-shell-lets-the-dashboard-change-the.md) | A settings host action, turned on only in a shell, lets the dashboard change the host's settings | accepted, refines 0029 |
| [0040](0040-one-flai-host-per-machine-runs-flai-serve-and-each-project-s-mcp-server-as-its.md) | One flai host per machine runs flai serve and each project's MCP server as its children | accepted, supersedes 0034, refines 0029 |
| [0041](0041-a-story-in-ready-gets-its-agent-whether-it-entered-ready-before-or-after-flai.md) | A story in ready gets its agent whether it entered ready before or after flai serve started | accepted, supersedes 0038 in part |
| [0042](0042-someone-attending-holds-a-ready-story-back-for-the-attended-window-then-flai.md) | Someone attending holds a ready story back for the attended window, then flai serve starts it | accepted, refines 0038 and 0041, superseded by 0043 |
| [0043](0043-flai-serve-starts-a-ready-story-s-agent-whenever-the-in-progress-limit-has-room.md) | flai serve starts a ready story's agent whenever the in-progress limit has room, and the operator can restart one that dropped or failed | accepted, supersedes 0042, refines 0038 and 0041 |
| [0044](0044-a-retry-of-a-failed-agent-on-a-ready-story-with-the-in-progress-limit-full-is.md) | A retry of a failed agent on a ready story with the in-progress limit full is queued until there is room | accepted, refines 0043 |
| [0045](0045-a-system-flow-repository-under-a-folder-named-for-import-is-served-whether-or.md) | A system-flow repository under a folder named for import is served whether or not it is registered | accepted, refines 0035 |
| [0046](0046-a-ready-story-whose-claim-overlaps-an-open-story-s-is-held-yellow-and-with-its.md) | A ready story whose claim overlaps an open story's is held, yellow and with its reason, until that story is accepted | accepted, refines 0019 |
| [0047](0047-an-agent-is-primed-with-what-its-story-s-topics-claim-and-links-select.md) | An agent is primed with what its story's topics, claim, and links select: conventions and design by topic, the documents it names, ranked sections, and a catalog of the rest | accepted, refines 0001 and 0013 |
| [0048](0048-pushing-accepted-work-releases-nothing-unless-the-auto-publish-host-action-is-enabled.md) | Pushing accepted work releases nothing unless the auto-publish host action is enabled | accepted, refines 0032 |
| [0049](0049-a-story-s-context-pack-fits-a-size-budget-what-the-story-names-loads-whole-what.md) | A story's context pack fits a size budget: what the story names loads whole, what its topics and links select loads as briefs, ranking fills the rest, and the agent fetches sections on demand | accepted, refines 0047 |
| [0050](0050-a-document-a-story-names-only-by-its-path-written-out-is-briefed-when-it-is.md) | A document a story names only by its path written out is briefed when it is larger than an eighth of the pack's budget | accepted, refines 0049 |
| [0051](0051-work-items-record-the-tokens-and-cost-their-agents-spent-measured-from-the.md) | Work items record the tokens and cost their agents spent, measured from the agents' logs and summed up the hierarchy when an item is done | accepted |
| [0052](0052-flai-asks-its-questions-a-line-at-a-time-with-its-own-prompt-package-not-with.md) | flai asks its questions a line at a time with its own prompt package, not with huh | accepted, refines 0006 |
| [0053](0053-usage-is-charted-as-spend-over-time-in-buckets-for-every-item-type-at-once-and.md) | Usage is charted as spend over time in buckets, for every item type at once, and agent time is counted in minutes | accepted, refines 0051 |
| [0054](0054-every-chart-spans-the-window-chosen-its-time-axis-runs-from-the-window-s-start.md) | Every chart spans the window chosen: its time axis runs from the window's start to now, and it plots only what the window holds | accepted, refines 0053 |
| [0055](0055-a-story-moves-back-one-column-from-ready-in-progress-review-or-cancelled-and.md) | A story moves back one column from ready, in progress, review, or cancelled, and done is final; an item reopened from cancelled is not completed | accepted, refines 0004 |
| [0056](0056-time-in-state-is-one-stacked-bar-per-day-of-the-window-the-mean-hours-per-state.md) | Time in state is one stacked bar per day of the window, the mean hours per state of the items completed that day | accepted, refines 0054 |
| [0057](0057-the-dashboard-charts-agent-time-and-cost-per-item-per-model-over-time-instead.md) | The dashboard charts agent time and cost per item per model over time, instead of completion, aging, and estimates, and charts per item draw item types only | accepted, refines 0053 |
| [0058](0058-the-agent-host-action-lets-the-dashboard-stop-a-story-s-agent-and-a-stop.md) | The agent host action lets the dashboard stop a story's agent, and a stop signals only the process flai started | accepted, refines 0043 |
| [0059](0059-a-story-s-agent-hands-search-test-runs-and-verification-to-an-explorer-and-a.md) | A story's agent hands search, test runs, and verification to an explorer and a verifier sub-agent that can only read, primed by role, and asks the designer itself | accepted |
| [0060](0060-a-claude-code-pretooluse-hook-flai-guard-refuses-any-sub-agent-s-call-that.md) | A Claude Code PreToolUse hook, flai guard, refuses any sub-agent's call that would change a work item, a thread, or the repository's history | accepted, refines 0059 |
| [0061](0061-flai-lints-the-markdown-it-writes-in-wip-with-its-own-implementation-of-the.md) | flai lints the markdown it writes in wip with its own implementation of the markdownlint rules it can break, configured by the project's markdownlint file | accepted |
| [0062](0062-flai-host-restarts-the-dashboard-container-flai-dashboard-recorded-when-it-is.md) | flai host restarts the dashboard container flai dashboard recorded when it is gone or not answering, unless dashboard.no_restart is set | accepted, refines 0040 |
| [0063](0063-the-dashboard-refuses-a-second-flai-for-a-project-with-close-code-4409-while.md) | The dashboard refuses a second flai for a project with close code 4409 while the one it has answers, and flai waits a minute before it dials again | accepted, refines 0029 |
| [0064](0064-a-story-in-ready-or-in-progress-with-no-agent-run-on-this-host-is-started-here.md) | A story in ready or in progress with no agent run on this host is started here on the operator's word, and its agent is told who began it, where, and when | accepted, refines 0043 |
| [0065](0065-a-story-s-agent-carries-a-model-per-sub-agent-role-and-claude-code-runs-each.md) | A story's agent carries a model per sub-agent role, and claude-code runs each role's sub-agent on it | accepted, refines 0037 and 0059 |
| [0066](0066-an-experiment-story-is-accepted-like-any-other-and-records-its-results-in-a.md) | An experiment story is accepted like any other and records its results in a document under design/experiments | accepted, refines 0025 and 0032 |
| [0067](0067-accepted-work-reaches-the-remote-only-when-it-is-published-and-agents-publish.md) | Accepted work reaches the remote only when it is published, and agents publish only when the operator asks | proposed, refines 0032 and 0048 |
| [0068](0068-a-convention-s-roles-list-every-agent-that-reads-it-the-story-s-agent-included.md) | A convention's roles list every agent that reads it, the story's agent included, and a convention without roles is read by all | accepted, refines 0059 |
| [0069](0069-agents-commit-each-task-on-the-story-branch-after-flai-stream-sync-and-never.md) | Agents commit each task on the story branch after flai stream sync, and never rebase or merge it by hand | accepted, refines 0019 |
| [0070](0070-releases-are-signed-with-a-cosign-key-pair-verified-before-they-are-installed.md) | Releases are signed with a cosign key pair, verified before they are installed or run, and each component refuses a peer without a signed release stamp | accepted, refines 0029 and 0031 |
| [0071](0071-a-task-s-usage-is-the-calls-of-the-sub-agents-started-for-it-and-an-even-share.md) | A task's usage is the calls of the sub-agents started for it and an even share of the story's agent's calls while it was in progress | accepted, refines 0051 |
| [0072](0072-a-story-s-agent-plans-its-tasks-with-after-and-hands-each-to-a-task-sub-agent.md) | A story's agent plans its tasks with after and hands each to a task sub-agent, and runs a layer's tasks at once only when they are long | accepted, refines 0059 |
| [0073](0073-a-full-review-holds-the-pull-and-flai-check-strict-passes-over-review-over-its.md) | A full review holds the pull, and flai check --strict passes over review over its limit | accepted, refines 0043 |
| [0074](0074-work-items-carry-planning-data-a-story-s-draft-flag-an-epic-s-or-story-s-cost.md) | Work items carry planning data: a story's draft flag, an epic's or story's cost of delay, and a story's forecast, each saying who set it | accepted |
| [0075](0075-the-planner-the-orchestrator-and-the-analyzer-prime-by-role-plan-orchestrate-or.md) | The planner, the orchestrator, and the analyzer prime by role, plan, orchestrate, or analyze, with a full pack of their own and no story | accepted, refines 0059 and 0068 |
| [0076](0076-an-epic-follows-its-stories-to-ready-and-in-progress-with-the-first-to-review.md) | An epic follows its stories: to ready and in-progress with the first, to review with the last open one, to done when that one is accepted | accepted, refines 0004 |
| [0077](0077-a-story-that-was-a-draft-records-who-finalized-it-and-when-and-the-dashboard.md) | A story that was a draft records who finalized it and when, and the dashboard finalizes through item.finalize | proposed, refines 0074 |
| [0078](0078-an-accepted-epic-contributes-no-release-bump-of-its-own-its-stories-carry-theirs.md) | An accepted epic contributes no release bump of its own: its stories carry theirs | accepted, refines 0076 |
| [0079](0079-the-planner-the-orchestrator-and-the-analyzer-each-log-their-activities-in-one.md) | The planner, the orchestrator, and the analyzer each log their activities in one document per project under wip/agents, written by flai with its totals in front matter that flai stats reports | accepted |
| [0080](0080-a-cost-of-delay-stamps-its-inputs-and-its-value-apart.md) | A cost of delay stamps its inputs and its value apart | accepted, refines 0074 |
| [0081](0081-flai-stats-reports-forecast-error-cost-of-delay-waiting-holds-touches-drift-and.md) | flai stats reports forecast error, cost of delay, waiting, holds, touches drift, and strategic use per day from the files and git | accepted, refines 0074 |
