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
| [0007](0007-sveltekit-spa-with-node-adapter-in-docker.md) | SvelteKit SPA with node adapter, shipped as a Docker image | accepted |
| [0008](0008-cli-config-location.md) | CLI configuration lives in ~/.flai/config.json | accepted |
| [0009](0009-agent-narrative-per-story.md) | One agent narrative per story in wip/agents | accepted |
| [0010](0010-shell-out-to-git-and-docker.md) | flai shells out to git and docker | accepted |
| [0011](0011-project-manifest.md) | system-flow.yaml marks a project and owns folder names | accepted |
| [0012](0012-flat-kanban-folders-and-archive.md) | Kanban items are flat per type, hierarchy by parent key, archive on completion | accepted |
| [0013](0013-conventions-folder.md) | design/conventions holds agent norms, one file per topic | accepted |
