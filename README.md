# system-flow

An agentic lean project management system that lives inside your monorepo. Design, documentation, and work in process are markdown files with front matter. A CLI (`flai`) creates and manages conforming repositories. A dashboard (`flaiover`) turns the files into a kanban board, flow charts, and a searchable documentation explorer.

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
| `flai/` | The Go CLI (not yet started, see E-002) |
| `flaiover/` | The SvelteKit dashboard (not yet started, see E-003) |

This repository follows its own standard. `CLAUDE.md` is how agents work here.
