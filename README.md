# system-flow

An agentic lean project management system that lives inside your monorepo. Design, documentation, and work in process are markdown files with front matter. A CLI (`flai`) creates and manages conforming repositories. A dashboard (`flaiover`) turns the files into a kanban board, flow charts, and a searchable documentation explorer.

## Install flai

```bash
curl -fsSL https://raw.githubusercontent.com/bytepunx/system-flow/main/install.sh | sh
```

The script detects your platform, downloads the latest release, verifies its checksum, and installs `flai` to `$HOME/.flai/bin` (set `FLAI_INSTALL_DIR` to change that), a directory you already own, so it never needs `sudo`. While this repository is private, log in with `gh auth login` or set `GITHUB_TOKEN` first. Later, `flai self-upgrade` does the same from the installed binary. Details and alternatives: [docs/users/flai.md](docs/users/flai.md#install).

| I want to | Go to |
|-----------|-------|
| Understand the conventions | [docs/users/index.md](docs/users/index.md) |
| Use the CLI | [docs/users/flai.md](docs/users/flai.md) |
| Run the dashboard | [docs/operators/index.md](docs/operators/index.md) |
| Contribute | [docs/contributors/index.md](docs/contributors/index.md) |
| Read the design | [design/system/README.md](design/system/README.md) |
| See what is being worked on | [wip/kanban/board.md](wip/kanban/board.md) |

## Repository

| Path | What |
|------|------|
| `design/` | ADRs, living system design, technology choices |
| `docs/` | Documentation for users, operators, contributors |
| `wip/` | Board, work items, agent narratives |
| `template/` | Prototype of the template repository |
| `flai/` | The Go CLI |
| `flaiover/` | The SvelteKit dashboard |

This repository follows its own standard. `CLAUDE.md` is how agents work here.

## License

system-flow, flai, and flaiover are distributed under the [Bytepunx Shield License 1.0](LICENSE.md), an adaptation of the PolyForm Shield License. Use is free for 501(c)(3) nonprofits, for educators and students at state-funded educational institutions, for security research on the software, for start-ups under the revenue and funding limits the license states, and for individuals whose work with it grosses under one million US dollars a year. Any other commercial use needs a paid license from the licensor (<alex@robsonandmilligan.com>). Nobody, paid or free, may offer system-flow or a fork of it as a hosted service or license it to others. A copy of the license must travel with every copy of the software: `flai license` prints the one built into the binary, and the dashboard shows the one its image carries under Host › License. The license text is what binds; this paragraph only points at it.
