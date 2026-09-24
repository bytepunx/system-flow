---
title: Settings index
updated: 2026-09-24
status: active
---

# Settings index

Every setting of flai and of the flaiover dashboard, by where it is set. Each row says what the setting is for in a phrase and links to where it is explained. Look a setting up by its name with your browser's find.

| Kind | Where it is kept | Changed with |
|------|------------------|--------------|
| [flai configuration](#flai-configuration) | `~/.flai/config.json`, one per user and host, or the file `--config` or `FLAI_CONFIG` names | `flai config set` for the first table; the `flai serve` commands for the second |
| [Project manifest](#project-manifest) | `system-flow.yaml` at the project's root, committed | By hand for `name`, `description`, `dashboard`; flai for the rest |
| [Other project files](#other-project-files) | The board's front matter, `.mcp.json` | The board, `flai order`, by hand |
| [Environment variables](#environment-variables) | The shell flai runs in, or what flai hands a process it starts | Your shell profile or the command line |
| [flaiover container](#flaiover-container) | The dashboard container's environment and mounts | `flai dashboard`, which starts it; by hand only for a container you run yourself |
| [Flags](#flags) | The command line | Each command's flags |

When a setting can be made in more than one place, the first of these wins: a flag, then the project's `system-flow.yaml`, then `~/.flai/config.json`, then the default. Secrets are not settings and are listed under [where their files are](index.md#authentication); none of them is ever an environment variable of the dashboard.

## flai configuration

`~/.flai/config.json`, created with the defaults on the first command that needs it. A key flai does not know is refused when the file is read. The whole file is explained in [the flai guide's Configuration](../users/flai.md#configuration). `flai config path` names the file in use.

Keys `flai config set` reaches:

| Key | Default | For |
|-----|---------|-----|
| `template.repo` | `https://github.com/bytepunx/system-flow-template` | The template `flai new`, `flai import`, and `flai upgrade` use when `--template` is not given |
| `template.ref` | `main` | The template's branch, tag, or commit when `--ref` is not given |
| `dashboard.image` | `ghcr.io/bytepunx/flaiover` | The image `flai dashboard` runs ([Access to the image](index.md#access-to-the-image)) |
| `dashboard.tag` | `latest` | The image's tag |
| `dashboard.port` | `4242` | The host port the dashboard is published on |
| `dashboard.bind` | `0.0.0.0` | The host address it is published on; `127.0.0.1` keeps it to this machine ([Security posture](index.md#security-posture)) |
| `dashboard.push_key` | empty | Retired and ignored; `flai dashboard` names it while it is set ([Upgrading from a release that mounted the repository](index.md#upgrading-from-a-release-that-mounted-the-repository)) |
| `dashboard.push_known_hosts` | empty | Retired and ignored, as `dashboard.push_key` |
| `cache_dir` | `~/.flai/cache`, or `FLAI_CACHE_DIR` when the file is created | Where git templates are cloned |
| `author` | your login name | The default owner of new items and `--by` of transitions and acceptances |
| `worktrees.relative_paths` | `false` | Link story worktrees with relative paths, git 2.48 or newer ([Relative worktree links](../users/flai.md#relative-worktree-links-opt-in)) |

Keys only the `flai serve` commands on the host change, or the dashboard's Settings page once the `settings` host action is on ([the settings host action](index.md#the-settings-host-action-changing-the-hosts-settings-from-the-dashboard-s-0105)):

| Key | Default | Changed with | For |
|-----|---------|--------------|-----|
| `host_actions.<name>` | absent: every action off | `flai serve enable`, `flai serve disable` | For each host action (`push`, `agent`, `checks`, `dashboard`, `host`, `settings`), the main checkouts it is on for, or `*` for every project ([The push host action](index.md#the-push-host-action)) |
| `agent.command` | none | `flai serve agent set -- ...`, `flai serve agent clear` | What starts a ready story's agent when the story names no harness ([Starting an agent](index.md#starting-an-agent-when-a-story-becomes-ready)) |
| `agent.name` | `agent` | `flai serve agent set --name` | The `FLAI_AGENT` prefix of the agents flai serve starts |
| `agent.attended_minutes` | `6` | `flai serve agent set --attended-minutes` | How recent a sign of someone attending must be, and how long it holds a ready story back |
| `agent.harnesses.<name>.program` | the harness's own program, `claude` for `claude-code` | `flai serve agent harness <name> --program` | The binary a harness runs on this host |
| `agent.harnesses.<name>.args` | the adapter's, `--permission-mode acceptEdits --allowedTools Bash,mcp__flai` for `claude-code` | `flai serve agent harness <name> -- ...`, `--reset` | What an agent that harness starts may do |
| `checks.commands[].name` | none: the manifest's `checks` apply | `flai serve checks set --name`, `flai serve checks clear` | The name of one check a story in review is run with ([The checks host action](index.md#the-checks-host-action-s-0082)) |
| `checks.commands[].command` | none | `flai serve checks set --name <name> -- ...` | That check's argument list, never run through a shell |
| `checks.timeout_minutes` | `15` | `flai serve checks timeout` | The bound on one run of every check together |
| `import_roots` | none | `flai serve import add`, `flai serve import remove` | Folders whose git repositories the board offers to import ([Importing repositories](index.md#importing-repositories-from-the-board-s-0098)) |

Beside the file, in the folders `serve` and `host`, flai keeps state, tokens, and logs, not settings; the [backup runbook](runbooks/backup.md) lists them.

## Project manifest

`system-flow.yaml` at the project's root. The schema and its rules are in [design/system/project-manifest.md](../../design/system/project-manifest.md). It is committed, so a setting here is the same for everyone who works on the project.

| Key | Default | For |
|-----|---------|-----|
| `version` | required | The manifest's schema version, `1` |
| `name` | required | The project's name, in titles and on the dashboard |
| `key` | none; `flai check` warns | The project's short key, in the dashboard's addresses and the `X-Flai-Project-Key` header ([Project identity](index.md#project-identity)) |
| `description` | empty | One line about the project |
| `owner` | empty | Who the dashboard's writes are recorded as; `designer` when empty ([Who acts when the dashboard writes](index.md#who-acts-when-the-dashboard-writes)) |
| `repo` | empty | The repository's URL |
| `template.repo` | set by `flai new` or `flai import` | The template the project came from, which `flai upgrade` fetches |
| `template.ref` | set by `flai new` or `flai import` | The template's branch, tag, or commit |
| `template.version` | set by flai | The template version applied, which `flai upgrade` compares against |
| `template.applied` | set by flai | When it was applied |
| `layout.<name>` | `design`, `docs`, `wip` | The folder names; `design`, `docs`, and `wip` are required |
| `projects[].name` | none | A releasable component's name, the prefix of its release tags |
| `projects[].path` | none | Its folder at the root |
| `projects[].kind` | none | `go`, `sveltekit`, `template`, and so on; `template` is released by its version file, not a tag |
| `projects[].tags` | none | Story tags that mean a story delivers to it |
| `dashboard.image` | the configuration's | The image, for this project, over `dashboard.image` in the configuration |
| `dashboard.tag` | the configuration's | The image's tag, likewise |
| `dashboard.port` | the configuration's | The host port, likewise |
| `dashboard.bind` | the configuration's | The host address, likewise |
| `dashboard.autocommit` | `true` | Commit documents saved from the dashboard, one path per commit |
| `dashboard.notify_url` | unset | A webhook the dashboard posts new inbox entries to |
| `checks[].name` | none | A check a story in review is run with, when the host names none ([The checks host action](index.md#the-checks-host-action-s-0082)) |
| `checks[].command` | none | Its argument list |
| `agent.harness` | none | The harness every new story gets, such as `claude-code` ([Starting an agent](index.md#starting-an-agent-when-a-story-becomes-ready)) |
| `agent.model` | none | The model every new story gets |
| `agent.config.<name>` | none | Options for the harness; `claude-code` takes `effort` (`low`, `medium`, `high`, `xhigh`, `max`), `max_budget_usd`, and `fallback_model` |

`dashboard.image`, `tag`, `port`, and `bind` matter only when they decide what the one shared container is started with, which is the first `flai dashboard` on the host. A story's own `agent:` front matter, copied from `agent` when the story is made, is the story's, not a setting; `flai edit` and the story's page change it.

## Other project files

| Setting | Where | Default | For |
|---------|-------|---------|-----|
| `wip_limits` | the front matter of `wip/kanban/board.md` | as the template sets them | The WIP limit of each column; flai serve starts agents only while `in-progress` has room ([design/system/workflow.md](../../design/system/workflow.md)) |
| `order` | the front matter of `wip/kanban/board.md` | empty | The pull order of `ready` and `backlog` stories; drag cards on the board or run `flai order` |
| `mcpServers.flai` | `.mcp.json` | `flai mcp` on stdio | How an agent on the host reaches flai's MCP server ([MCP over HTTP](index.md#mcp-over-http)) |

## Environment variables

### Read by flai

| Variable | Default | For |
|----------|---------|-----|
| `FLAI_CONFIG` | `~/.flai/config.json` | The configuration file when `--config` is not given; `flai host` and `flai serve` keep their state in `host` and `serve` beside it |
| `FLAI_CACHE_DIR` | `~/.flai/cache` | The `cache_dir` written when the configuration file is created |
| `FLAI_AGENT` | none | Who narrative entries, thread entries, and transitions are recorded as; the MCP server's agent unless `flai mcp --agent` names one |
| `FLAI_SESSION` | none | Which session a narrative entry belongs to |
| `LOG_LEVEL` | `info` | `debug`, `info`, `warn`, `error`, or `fatal` ([Logging](../users/flai.md#logging)) |
| `LOG_FORMAT` | text on a terminal, JSON otherwise | `text` or `json` |
| `FLAI_HOST_ADDR` | `127.0.0.1:4241` | Where `flai host` listens; set it for every flai you run when another program holds the port ([flai host](index.md#flai-host-the-process-that-runs-the-others-s-0106)) |
| `FLAI_RELEASES_API` | `https://api.github.com` | Another source of releases for `flai self-upgrade`, `flai host check`, and `flai host upgrade` |
| `FLAI_INSTALL_DIR` | `~/.flai/bin` | Where `flai self-upgrade`, run from a checkout's own build, installs a release |
| `GITHUB_TOKEN`, `GH_TOKEN` | `gh auth token` when neither is set | GitHub API and registry authentication for `flai self-upgrade`, `flai host upgrade`, and `flai dashboard`'s image pull |
| `USER` | none | The `author` written when the configuration file is created, when the system cannot name the login |
| `GITHUB_ACTOR` | a placeholder | The user name `flai dashboard` logs in to `ghcr.io` with; the registry checks only the token |

### Set by flai for the processes it starts

Do not set these yourself; a command you write for `flai serve agent set` may read them.

| Variable | Set by | Holds |
|----------|--------|-------|
| `FLAI_HOST_URL`, `FLAI_HOST_TOKEN` | `flai host`, for `flai serve` and the MCP servers | The host's address and the token its API takes |
| `FLAI_STORY` | `flai serve`, for an agent it starts | The story the agent is to work |
| `FLAI_STARTED_BY` | `flai serve` | `flai-serve` |
| `FLAI_AGENT`, `FLAI_SESSION` | `flai serve` | The agent's own name, `<agent.name>-<story>`, and a session made from the start time |
| `FLAI_MODEL`, `FLAI_HARNESS` | `flai serve`, for your `agent.command` | The story's model and harness |
| `FLAI_AGENT_CONFIG` | `flai serve`, for your `agent.command` | The story's `agent.config` as a JSON object |
| `FLAI_ANSWERED` | `flai serve`, for an agent started again because its question was answered | The thread's ID |

### Read by install.sh

| Variable | Default | For |
|----------|---------|-----|
| `FLAI_INSTALL_DIR` | `$HOME/.flai/bin` | Where it installs flai ([install runbook](runbooks/install.md)) |
| `FLAI_VERSION` | the newest release | The release to install, such as `1.16.1` |
| `FLAI_REPO` | `bytepunx/system-flow` | The repository whose releases it installs from |
| `FLAI_API` | `https://api.github.com` | The GitHub API it asks |
| `GITHUB_TOKEN`, `GH_TOKEN` | `gh auth token` | Authentication while the repository is private |

## flaiover container

`flai dashboard` starts the container with everything below already set; change it only for a container you run yourself ([the docker run command](index.md#running-the-dashboard)).

| Setting | Default | For |
|---------|---------|-----|
| `FLAIOVER_TOKEN_FILE` | `/run/secrets/flaiover_token`, mounted read-only | The login token's file ([Authentication](index.md#authentication)); never the token itself |
| `FLAIOVER_AGENT_KEY_FILE` | `/run/secrets/flaiover_agent_key`, mounted read-only | The credential flai on the host proves itself with ([The connection from flai on the host](index.md#the-connection-from-flai-on-the-host)) |
| `FLAIOVER_METRICS_PUBLIC` | unset | `true` lets `/metrics` be read without the token ([Telemetry](index.md#telemetry)) |
| `FLAIOVER_AUTH` | unset | `off` turns authentication off, outside production only (`NODE_ENV` not `production`); for development |
| `PORT` | `3000` | The port the server listens on inside the container |
| `HOST` | `0.0.0.0` | The address it listens on inside the container |
| `PROTOCOL_HEADER` | `x-forwarded-proto` | The header a TLS-terminating proxy names the original scheme in ([How the dashboard knows whether it is served over HTTPS](index.md#how-the-dashboard-knows-whether-it-is-served-over-https)) |
| `SHUTDOWN_TIMEOUT` | `5` | Seconds after SIGTERM before every connection is closed |
| `NODE_ENV` | `production` | Anything else allows `FLAIOVER_AUTH=off` and text logs |
| `LOG_LEVEL` | `info` | `debug`, `info`, `warn`, or `error` |
| `LOG_FORMAT` | `json` | `text` for key-value text, outside production only |
| `OTEL_EXPORTER_OTLP_ENDPOINT`, `OTEL_EXPORTER_OTLP_TRACES_ENDPOINT` | unset: no traces | Where traces are sent over OTLP/HTTP |
| `OTEL_SERVICE_NAME` | `flaiover` | The service name on its spans |
| `OTEL_TRACES_SAMPLER`, `OTEL_TRACES_SAMPLER_ARG` | the OpenTelemetry SDK's | Sampling |
| `FLAIOVER_VERSION`, `FLAIOVER_COMMIT` | baked into the image | The release and commit `flaiover_build_info` reports; not for changing |
| `PROJECT_DIR` | the working directory | The project the development server and tests ask flai about; names nothing in the container |
| `--publish` | `<dashboard.bind>:<dashboard.port>:3000` | Where the host reaches the container |
| `--user` | your user and group | So that the two secret files, readable only by you, can be read |

Build arguments, for an image built with `flai dashboard --build` or `make flaiover-image`: `FLAIOVER_VERSION` and `FLAI_VERSION` from the newest `flaiover/v*` and `flai/v*` tags, `FLAI_COMMIT` the commit, and `FLAI_DATE` the build time. `flai dashboard --build` sets them.

## Flags

Every flag of every flai command is in [the command reference](../users/flai-reference.md), generated from `flai --help`. The flags that change a setting for one run are the ones named in the tables above.
