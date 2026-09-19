# Architecture Decision Records

One decision per file, numbered in order of acceptance. An ADR is never edited after acceptance except to record `superseded_by`. To change a decision, write a new ADR that supersedes the old one in whole or in part, then update the living documents in `design/system` and `design/tech`.

Copy `0000-template.md` to start a new one.

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
| [0022](0022-repository-mounted-at-its-host-path.md) | The dashboard mounts the repository at its host path; relative-path worktrees are an opt-in | accepted, refines 0019 |
| [0023](0023-documents-are-saved-through-flai.md) | Documents edited in the dashboard are saved through flai | accepted, refines 0016 |
