---
title: Settings index
updated: 2026-10-03
status: active
---

# Settings index

Every setting of flai and of the flaiover dashboard, by where it is set. Each row says what the setting is for in a phrase and links to where it is explained. Look a setting up by its name with your browser's find.

| Kind | Where it is kept | Changed with |
|------|------------------|--------------|
| [flai configuration](#flai-configuration) | `~/.flai/config.json`, one per user and host, or the file `--config` or `FLAI_CONFIG` names | `flai config set` for the first table; the `flai serve` commands for the second |
| [Project manifest](#project-manifest) | `system-flow.yaml` at the project's root, committed | By hand for `name`, `description`, `dashboard`; flai for the rest |
| [Other project files](#other-project-files) | The board's front matter, `.mcp.json`, `.claude/settings.json` | The board, `flai order`, by hand |
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
| `dashboard.no_restart` | `false` | `true` stops flai host restarting a dashboard that is gone or not answering ([The dashboard's watch](index.md#the-dashboards-watch)) |
| `dashboard.push_key` | empty | Retired and ignored; `flai dashboard` names it while it is set ([Upgrading from a release that mounted the repository](index.md#upgrading-from-a-release-that-mounted-the-repository)) |
| `dashboard.push_known_hosts` | empty | Retired and ignored, as `dashboard.push_key` |
| `cache_dir` | `~/.flai/cache`, or `FLAI_CACHE_DIR` when the file is created | Where git templates are cloned |
| `author` | your login name | The default owner of new items and `--by` of transitions and acceptances |
| `worktrees.relative_paths` | `false` | Link story worktrees with relative paths, git 2.48 or newer ([Relative worktree links](../users/flai.md#relative-worktree-links-opt-in)) |

Keys only the `flai serve` commands on the host change, or the dashboard's Settings page once the `settings` host action is on ([the settings host action](index.md#the-settings-host-action-changing-the-hosts-settings-from-the-dashboard-s-0105)):

| Key | Default | Changed with | For |
|-----|---------|--------------|-----|
| `host_actions.<name>` | absent: every action off | `flai serve enable`, `flai serve disable` | For each host action (`push`, `auto-publish`, `agent`, `checks`, `dashboard`, `host`, `settings`), the main checkouts it is on for, or `*` for every project ([The push host action](index.md#the-push-host-action), which lets the board publish). `auto-publish` has `flai push --pending` release first; it is shell only, outside the workflow, and the dashboard's Settings page neither lists nor changes it ([Pushing outside the workflow](index.md#accepting-releases-nothing-publishing-reaches-the-remote-s-0087-s-0195)) |
| `agent.command` | none | `flai serve agent set -- ...`, `flai serve agent clear` | What starts a ready story's agent when the story names no harness ([Starting an agent](index.md#starting-an-agent-when-a-story-becomes-ready)) |
| `agent.name` | `agent` | `flai serve agent set --name` | The `FLAI_AGENT` prefix of the agents flai serve starts |
| `agent.attended_minutes` | `6` | `flai serve agent set --attended-minutes` | How recent a sign of someone attending must be, and how long it holds a ready story back |
| `agent.harnesses.<name>.program` | the harness's own program, `claude` for `claude-code` | `flai serve agent harness <name> --program` | The binary a harness runs on this host |
| `agent.harnesses.<name>.args` | the adapter's, `--permission-mode acceptEdits --allowedTools Bash,mcp__flai` for `claude-code` | `flai serve agent harness <name> -- ...`, `--reset` | What an agent that harness starts may do |
| `checks.commands[].name` | none: the manifest's `checks` apply | `flai serve checks set --name`, `flai serve checks clear` | The name of one check a story in review is run with ([The checks host action](index.md#the-checks-host-action-s-0082)) |
| `checks.commands[].command` | none | `flai serve checks set --name <name> -- ...` | That check's argument list, never run through a shell |
| `checks.timeout_minutes` | `15` | `flai serve checks timeout` | The bound on one run of every check together |
| `import_roots` | none | `flai serve import add`, `flai serve import remove` | Folders whose git repositories the board offers to import, and whose repositories with a `system-flow.yaml` `flai serve` serves ([Importing repositories](index.md#importing-repositories-from-the-board-s-0098)) |

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
| `prime.budget` | `80KB` | The size a story's context pack fits, for `flai prime --story` and the MCP `prime` tool; bytes, or a number with `KB` or `MB`; `--budget` overrides it for one run ([Prime a session](../users/flai.md#prime-a-session)) |
| `agent.config.<name>` | none | Options for the harness; `claude-code` takes `effort` (`low`, `medium`, `high`, `xhigh`, `max`), `max_budget_usd`, and `fallback_model` |
| `issues.story_after` | `168h` | How long an issue may stay open with no open story linking it before `flai check` warns (`issues.no-story`); a Go duration such as `24h`, or `0` to turn the warning off ([Record recurring friction](../users/flai.md#record-recurring-friction)) |
| `flai.minimum` | unset | The oldest flai release that may read the project, `X.Y.Z`; an older one stops before reading any item and names the version needed. Publishing a flai release that changes the front-matter fields flai reads raises it ([Keeping the host's flai current](index.md#keeping-the-hosts-flai-current)) |
| `agent.roles.<name>.harness` | none | The harness of a sub-agent role (`explore`, `verify`); for `claude-code`, only `claude-code`, since a sub-agent runs in the story's session |
| `agent.roles.<name>.model` | none | The model a sub-agent role runs, over the one its definition in `.claude/agents/` names, such as `haiku` or `claude-sonnet-5-5` |
| `agent.roles.<name>.config.<name>` | none | Options for a sub-agent role; `claude-code` takes none |

`dashboard.image`, `tag`, `port`, and `bind` matter only when they decide what the one shared container is started with, which is the first `flai dashboard` on the host. A story's own `agent:` front matter, copied from `agent` when the story is made, is the story's, not a setting; `flai edit` and the story's page change it.

## Other project files

| Setting | Where | Default | For |
|---------|-------|---------|-----|
| `wip_limits` | the front matter of `wip/kanban/board.md` | as the template sets them | The WIP limit of each column; flai serve starts agents only while `in-progress` has room ([design/system/workflow.md](../../design/system/workflow.md)) |
| `order` | the front matter of `wip/kanban/board.md` | empty | The pull order of `ready` and `backlog` stories; drag cards on the board or run `flai order` |
| `mcpServers.flai` | `.mcp.json` | `flai mcp` on stdio | How an agent on the host reaches flai's MCP server ([MCP over HTTP](index.md#mcp-over-http)) |
| `hooks.PreToolUse` | `.claude/settings.json` | `flai guard` before `Bash` and `mcp__flai__.*` | Refuses a Claude Code sub-agent's writes to work items, threads, and history ([Sub-agents](../users/flai.md#sub-agents)) |

## Environment variables

### Read by flai

| Variable | Default | For |
|----------|---------|-----|
| `FLAI_CONFIG` | `~/.flai/config.json` | The configuration file when `--config` is not given; `flai host` and `flai serve` keep their state in `host` and `serve` beside it. Not passed on to an agent `flai serve` starts |
| `FLAI_CACHE_DIR` | `~/.flai/cache` | The `cache_dir` written when the configuration file is created |
| `FLAI_AGENT` | none | Who narrative entries, thread entries, and transitions are recorded as; the MCP server's agent unless `flai mcp --agent` names one |
| `FLAI_SESSION` | none | Which session a narrative entry belongs to |
| `LOG_LEVEL` | `info` | `debug`, `info`, `warn`, `error`, or `fatal` ([Logging](../users/flai.md#logging)) |
| `LOG_FORMAT` | text on a terminal, JSON otherwise | `text` or `json` |
| `FLAI_SLOW_REQUEST` | `500ms` | How long a request to `flai serve` or `flai mcp` may take before its `request answered` event is logged at `info` rather than `debug` ([Request timing](../users/flai.md#request-timing-s-0152)) |
| `FLAI_PPROF_ADDR` | none | A loopback address, such as `127.0.0.1:6060`, where `flai serve` offers Go's CPU, heap, and goroutine profiles at `/debug/pprof/`; an address beyond this machine is refused ([Request timing](../users/flai.md#request-timing-s-0152)) |
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
| `FLAI_HOST_URL`, `FLAI_HOST_TOKEN` | `flai host`, for `flai serve` and the MCP servers; never passed on to an agent `flai serve` starts | The host's address and the token its API takes; `flai serve` uses them only when the host's token is the one beside its own config |
| `FLAI_STORY` | `flai serve`, for an agent it starts | The story the agent is to work |
| `FLAI_STARTED_BY` | `flai serve` | `flai-serve` |
| `FLAI_AGENT`, `FLAI_SESSION` | `flai serve` | The agent's own name, `<agent.name>-<story>`, and a session made from the start time |
| `FLAI_MODEL`, `FLAI_HARNESS` | `flai serve`, for your `agent.command` | The story's model and harness |
| `FLAI_AGENT_CONFIG` | `flai serve`, for your `agent.command` | The story's `agent.config` as a JSON object |
| `FLAI_AGENT_ROLES` | `flai serve`, for your `agent.command`, when the story's agent has roles | The story's `agent.roles` as a JSON object, each role with its `harness`, `model`, and `config` |
| `FLAI_ANSWERED` | `flai serve`, for an agent started again because its question was answered | The thread's ID |
| `FLAI_COMMIT` | `flai serve`, for your `agent.command` started to commit what a story's worktree holds (`flai serve agent commit`) | The worktree's path |

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
| `FLAIOVER_LICENSE` | unset: `LICENSE.md` beside the server (`/app/LICENSE.md` in the image), else one directory up | The license file the Host › License page and `GET /api/license` show (S-0231); for a build that keeps it elsewhere |
| `--publish` | `<dashboard.bind>:<dashboard.port>:3000` | Where the host reaches the container |
| `--user` | your user and group | So that the two secret files, readable only by you, can be read |

Build arguments, for an image built with `flai dashboard --build` or `make flaiover-image`: `FLAIOVER_VERSION` and `FLAI_VERSION` from the newest `flaiover/v*` and `flai/v*` tags, `FLAI_COMMIT` the commit, and `FLAI_DATE` the build time. `flai dashboard --build` sets them.

## Flags

Every flag of every flai command, with the commands that take it. What a flag means and its default are in the command's section of [the command reference](../users/flai-reference.md), which, like this table, is generated from `flai --help`. The flags that change a setting for one run are named in the tables above.

<!-- flags:begin: generated by flai reference from the command help; run make flai-reference -->
| Flag | Taken by |
|------|----------|
| `--addr` | [flai mcp http](../users/flai-reference.md#flai-mcp-http), [flai mcp start](../users/flai-reference.md#flai-mcp-start) |
| `--after` | [flai edit](../users/flai-reference.md#flai-edit), [flai order](../users/flai-reference.md#flai-order), [flai story new](../users/flai-reference.md#flai-story-new), [flai task new](../users/flai-reference.md#flai-task-new) |
| `--agent` | [flai mcp](../users/flai-reference.md#flai-mcp) |
| `--agent-config` | [flai edit](../users/flai-reference.md#flai-edit), [flai story new](../users/flai-reference.md#flai-story-new) |
| `--all` | [flai board](../users/flai-reference.md#flai-board), [flai issue list](../users/flai-reference.md#flai-issue-list), [flai serve agent usage](../users/flai-reference.md#flai-serve-agent-usage), [flai thread list](../users/flai-reference.md#flai-thread-list) |
| `--all-projects` | [flai serve disable](../users/flai-reference.md#flai-serve-disable), [flai serve enable](../users/flai-reference.md#flai-serve-enable) |
| `--apply` | [flai release](../users/flai-reference.md#flai-release) |
| `--attach` | [flai dashboard](../users/flai-reference.md#flai-dashboard) |
| `--autocommit` | [flai adr accept](../users/flai-reference.md#flai-adr-accept), [flai adr new](../users/flai-reference.md#flai-adr-new), [flai adr topics](../users/flai-reference.md#flai-adr-topics), [flai agent clear](../users/flai-reference.md#flai-agent-clear), [flai agent set](../users/flai-reference.md#flai-agent-set), [flai edit](../users/flai-reference.md#flai-edit), [flai epic new](../users/flai-reference.md#flai-epic-new), [flai issue story](../users/flai-reference.md#flai-issue-story), [flai story new](../users/flai-reference.md#flai-story-new), [flai task new](../users/flai-reference.md#flai-task-new) |
| `--before` | [flai order](../users/flai-reference.md#flai-order) |
| `--bind` | [flai dashboard](../users/flai-reference.md#flai-dashboard), [flai dashboard restart](../users/flai-reference.md#flai-dashboard-restart), [flai dashboard upgrade](../users/flai-reference.md#flai-dashboard-upgrade) |
| `--body-stdin` | [flai adr new](../users/flai-reference.md#flai-adr-new), [flai edit](../users/flai-reference.md#flai-edit), [flai epic new](../users/flai-reference.md#flai-epic-new), [flai story new](../users/flai-reference.md#flai-story-new), [flai task new](../users/flai-reference.md#flai-task-new) |
| `--bottom` | [flai order](../users/flai-reference.md#flai-order) |
| `--bucket` | [flai stats](../users/flai-reference.md#flai-stats) |
| `--budget` | [flai prime](../users/flai-reference.md#flai-prime) |
| `--build` | [flai dashboard](../users/flai-reference.md#flai-dashboard) |
| `--by` | [flai accept](../users/flai-reference.md#flai-accept), [flai edit](../users/flai-reference.md#flai-edit), [flai move](../users/flai-reference.md#flai-move), [flai stats](../users/flai-reference.md#flai-stats), [flai stream answer](../users/flai-reference.md#flai-stream-answer), [flai thread new](../users/flai-reference.md#flai-thread-new), [flai thread reply](../users/flai-reference.md#flai-thread-reply), [flai thread resolve](../users/flai-reference.md#flai-thread-resolve) |
| `--cat` | [flai prime](../users/flai-reference.md#flai-prime) |
| `--check` | [flai self-upgrade](../users/flai-reference.md#flai-self-upgrade) |
| `--class` | [flai issue new](../users/flai-reference.md#flai-issue-new) |
| `--clear` | [flai touches](../users/flai-reference.md#flai-touches) |
| `--clear-after` | [flai edit](../users/flai-reference.md#flai-edit) |
| `--clear-agent` | [flai edit](../users/flai-reference.md#flai-edit) |
| `--clear-tags` | [flai edit](../users/flai-reference.md#flai-edit) |
| `--clear-topics` | [flai edit](../users/flai-reference.md#flai-edit) |
| `--clear-touches` | [flai edit](../users/flai-reference.md#flai-edit) |
| `--commit` | [flai import](../users/flai-reference.md#flai-import) |
| `--config` | every command ([global flags](../users/flai-reference.md#flai)), [flai agent set](../users/flai-reference.md#flai-agent-set) |
| `--cost` | [flai issue bump](../users/flai-reference.md#flai-issue-bump), [flai issue new](../users/flai-reference.md#flai-issue-new) |
| `--defaults` | [flai new](../users/flai-reference.md#flai-new) |
| `--deliver` | [flai release](../users/flai-reference.md#flai-release) |
| `--dir` | [flai self-upgrade](../users/flai-reference.md#flai-self-upgrade) |
| `--dry-run` | [flai accept](../users/flai-reference.md#flai-accept), [flai archive](../users/flai-reference.md#flai-archive), [flai import](../users/flai-reference.md#flai-import), [flai migrate ids](../users/flai-reference.md#flai-migrate-ids), [flai move](../users/flai-reference.md#flai-move), [flai push](../users/flai-reference.md#flai-push), [flai release](../users/flai-reference.md#flai-release), [flai template push](../users/flai-reference.md#flai-template-push), [flai upgrade](../users/flai-reference.md#flai-upgrade) |
| `--epic` | [flai issue story](../users/flai-reference.md#flai-issue-story), [flai story new](../users/flai-reference.md#flai-story-new) |
| `-f`, `--follow` | [flai dashboard logs](../users/flai-reference.md#flai-dashboard-logs), [flai serve agent stream](../users/flai-reference.md#flai-serve-agent-stream) |
| `--force` | [flai import](../users/flai-reference.md#flai-import), [flai new](../users/flai-reference.md#flai-new), [flai template push](../users/flai-reference.md#flai-template-push), [flai upgrade](../users/flai-reference.md#flai-upgrade) |
| `--from` | [flai checks tail](../users/flai-reference.md#flai-checks-tail), [flai serve agent stream](../users/flai-reference.md#flai-serve-agent-stream) |
| `--grace-seconds` | [flai checks cancel](../users/flai-reference.md#flai-checks-cancel) |
| `--harness` | [flai agent set](../users/flai-reference.md#flai-agent-set), [flai edit](../users/flai-reference.md#flai-edit), [flai story new](../users/flai-reference.md#flai-story-new) |
| `--hash` | [flai doc save](../users/flai-reference.md#flai-doc-save), [flai edit](../users/flai-reference.md#flai-edit) |
| `--heading` | [flai doc show](../users/flai-reference.md#flai-doc-show), [flai thread new](../users/flai-reference.md#flai-thread-new) |
| `--idle` | [flai mcp http](../users/flai-reference.md#flai-mcp-http), [flai mcp start](../users/flai-reference.md#flai-mcp-start) |
| `--image` | [flai dashboard](../users/flai-reference.md#flai-dashboard), [flai dashboard check](../users/flai-reference.md#flai-dashboard-check), [flai dashboard restart](../users/flai-reference.md#flai-dashboard-restart), [flai dashboard upgrade](../users/flai-reference.md#flai-dashboard-upgrade) |
| `--json` | every command ([global flags](../users/flai-reference.md#flai)) |
| `--keep-all` | [flai upgrade](../users/flai-reference.md#flai-upgrade) |
| `-n`, `--last` | [flai serve journal](../users/flai-reference.md#flai-serve-journal) |
| `--layout` | [flai import](../users/flai-reference.md#flai-import), [flai new](../users/flai-reference.md#flai-new) |
| `--limit` | [flai doc search](../users/flai-reference.md#flai-doc-search) |
| `--max-sessions` | [flai mcp http](../users/flai-reference.md#flai-mcp-http), [flai mcp start](../users/flai-reference.md#flai-mcp-start) |
| `--message` | [flai doc save](../users/flai-reference.md#flai-doc-save), [flai edit](../users/flai-reference.md#flai-edit) |
| `--model` | [flai agent set](../users/flai-reference.md#flai-agent-set), [flai edit](../users/flai-reference.md#flai-edit), [flai story new](../users/flai-reference.md#flai-story-new) |
| `--name` | [flai serve agent set](../users/flai-reference.md#flai-serve-agent-set), [flai serve checks set](../users/flai-reference.md#flai-serve-checks-set) |
| `--nature` | [flai edit](../users/flai-reference.md#flai-edit), [flai epic new](../users/flai-reference.md#flai-epic-new), [flai story new](../users/flai-reference.md#flai-story-new), [flai task new](../users/flai-reference.md#flai-task-new) |
| `--no-branch` | [flai stream open](../users/flai-reference.md#flai-stream-open) |
| `--no-commit` | [flai doc save](../users/flai-reference.md#flai-doc-save) |
| `--no-descriptions` | [flai completion bash](../users/flai-reference.md#flai-completion-bash), [flai completion fish](../users/flai-reference.md#flai-completion-fish), [flai completion powershell](../users/flai-reference.md#flai-completion-powershell), [flai completion zsh](../users/flai-reference.md#flai-completion-zsh) |
| `--no-git` | [flai new](../users/flai-reference.md#flai-new) |
| `--no-restart` | [flai dashboard token](../users/flai-reference.md#flai-dashboard-token) |
| `--no-serve` | [flai dashboard](../users/flai-reference.md#flai-dashboard) |
| `--note` | [flai issue bump](../users/flai-reference.md#flai-issue-bump), [flai issue new](../users/flai-reference.md#flai-issue-new) |
| `--on` | [flai thread list](../users/flai-reference.md#flai-thread-list), [flai thread new](../users/flai-reference.md#flai-thread-new) |
| `--open` | [flai dashboard](../users/flai-reference.md#flai-dashboard) |
| `--owner` | [flai epic new](../users/flai-reference.md#flai-epic-new), [flai issue story](../users/flai-reference.md#flai-issue-story), [flai story new](../users/flai-reference.md#flai-story-new), [flai task new](../users/flai-reference.md#flai-task-new) |
| `--parent` | [flai edit](../users/flai-reference.md#flai-edit) |
| `--pending` | [flai push](../users/flai-reference.md#flai-push), [flai release](../users/flai-reference.md#flai-release) |
| `--port` | [flai dashboard](../users/flai-reference.md#flai-dashboard), [flai dashboard restart](../users/flai-reference.md#flai-dashboard-restart), [flai dashboard upgrade](../users/flai-reference.md#flai-dashboard-upgrade) |
| `--print-body` | [flai adr new](../users/flai-reference.md#flai-adr-new), [flai epic new](../users/flai-reference.md#flai-epic-new), [flai story new](../users/flai-reference.md#flai-story-new), [flai task new](../users/flai-reference.md#flai-task-new) |
| `--program` | [flai serve agent harness](../users/flai-reference.md#flai-serve-agent-harness) |
| `--publish` | [flai push](../users/flai-reference.md#flai-push) |
| `--pull` | [flai dashboard](../users/flai-reference.md#flai-dashboard) |
| `--reason` | [flai block](../users/flai-reference.md#flai-block), [flai issue close](../users/flai-reference.md#flai-issue-close), [flai move](../users/flai-reference.md#flai-move), [flai thread resolve](../users/flai-reference.md#flai-thread-resolve) |
| `--ref` | [flai import](../users/flai-reference.md#flai-import), [flai new](../users/flai-reference.md#flai-new), [flai template push](../users/flai-reference.md#flai-template-push), [flai template show](../users/flai-reference.md#flai-template-show), [flai template use](../users/flai-reference.md#flai-template-use), [flai upgrade](../users/flai-reference.md#flai-upgrade) |
| `--refines` | [flai adr new](../users/flai-reference.md#flai-adr-new) |
| `--relock` | [flai upgrade](../users/flai-reference.md#flai-upgrade) |
| `--remote` | [flai template push](../users/flai-reference.md#flai-template-push) |
| `--replace` | [flai agent set](../users/flai-reference.md#flai-agent-set) |
| `--replace-all` | [flai upgrade](../users/flai-reference.md#flai-upgrade) |
| `--repo` | [flai self-upgrade](../users/flai-reference.md#flai-self-upgrade) |
| `--reset` | [flai serve agent harness](../users/flai-reference.md#flai-serve-agent-harness) |
| `--role` | [flai prime](../users/flai-reference.md#flai-prime) |
| `--role-config` | [flai agent set](../users/flai-reference.md#flai-agent-set), [flai edit](../users/flai-reference.md#flai-edit), [flai story new](../users/flai-reference.md#flai-story-new) |
| `--role-harness` | [flai agent set](../users/flai-reference.md#flai-agent-set), [flai edit](../users/flai-reference.md#flai-edit), [flai story new](../users/flai-reference.md#flai-story-new) |
| `--role-model` | [flai agent set](../users/flai-reference.md#flai-agent-set), [flai edit](../users/flai-reference.md#flai-edit), [flai story new](../users/flai-reference.md#flai-story-new) |
| `--rotate` | [flai dashboard token](../users/flai-reference.md#flai-dashboard-token), [flai mcp token](../users/flai-reference.md#flai-mcp-token) |
| `--show` | [flai edit](../users/flai-reference.md#flai-edit) |
| `--since` | [flai stats](../users/flai-reference.md#flai-stats) |
| `--status` | [flai adr new](../users/flai-reference.md#flai-adr-new) |
| `--story` | [flai issue bump](../users/flai-reference.md#flai-issue-bump), [flai issue list](../users/flai-reference.md#flai-issue-list), [flai issue new](../users/flai-reference.md#flai-issue-new), [flai issue story](../users/flai-reference.md#flai-issue-story), [flai prime](../users/flai-reference.md#flai-prime), [flai task new](../users/flai-reference.md#flai-task-new) |
| `--strict` | [flai check](../users/flai-reference.md#flai-check) |
| `--supersedes` | [flai adr new](../users/flai-reference.md#flai-adr-new) |
| `--tag` | [flai dashboard](../users/flai-reference.md#flai-dashboard), [flai dashboard check](../users/flai-reference.md#flai-dashboard-check), [flai dashboard restart](../users/flai-reference.md#flai-dashboard-restart), [flai dashboard upgrade](../users/flai-reference.md#flai-dashboard-upgrade), [flai edit](../users/flai-reference.md#flai-edit), [flai epic new](../users/flai-reference.md#flai-epic-new), [flai story new](../users/flai-reference.md#flai-story-new), [flai task new](../users/flai-reference.md#flai-task-new), [flai template push](../users/flai-reference.md#flai-template-push) |
| `--template` | [flai import](../users/flai-reference.md#flai-import), [flai new](../users/flai-reference.md#flai-new), [flai template show](../users/flai-reference.md#flai-template-show), [flai upgrade](../users/flai-reference.md#flai-upgrade) |
| `--timing` | [flai hostapi](../users/flai-reference.md#flai-hostapi) |
| `--title` | [flai edit](../users/flai-reference.md#flai-edit) |
| `--top` | [flai order](../users/flai-reference.md#flai-order) |
| `--topics` | [flai edit](../users/flai-reference.md#flai-edit), [flai epic new](../users/flai-reference.md#flai-epic-new), [flai story new](../users/flai-reference.md#flai-story-new) |
| `--touches` | [flai edit](../users/flai-reference.md#flai-edit), [flai epic new](../users/flai-reference.md#flai-epic-new), [flai story new](../users/flai-reference.md#flai-story-new), [flai task new](../users/flai-reference.md#flai-task-new) |
| `--trailer` | [flai accept](../users/flai-reference.md#flai-accept), [flai adr accept](../users/flai-reference.md#flai-adr-accept), [flai adr new](../users/flai-reference.md#flai-adr-new), [flai adr topics](../users/flai-reference.md#flai-adr-topics), [flai agent clear](../users/flai-reference.md#flai-agent-clear), [flai agent set](../users/flai-reference.md#flai-agent-set), [flai doc save](../users/flai-reference.md#flai-doc-save), [flai edit](../users/flai-reference.md#flai-edit), [flai epic new](../users/flai-reference.md#flai-epic-new), [flai import](../users/flai-reference.md#flai-import), [flai issue story](../users/flai-reference.md#flai-issue-story), [flai move](../users/flai-reference.md#flai-move), [flai story new](../users/flai-reference.md#flai-story-new), [flai task new](../users/flai-reference.md#flai-task-new) |
| `--type` | [flai stats](../users/flai-reference.md#flai-stats) |
| `--unset` | [flai agent set](../users/flai-reference.md#flai-agent-set) |
| `--unset-role` | [flai agent set](../users/flai-reference.md#flai-agent-set), [flai edit](../users/flai-reference.md#flai-edit) |
| `--var` | [flai import](../users/flai-reference.md#flai-import), [flai new](../users/flai-reference.md#flai-new), [flai upgrade](../users/flai-reference.md#flai-upgrade) |
| `-v`, `--verbose` | every command ([global flags](../users/flai-reference.md#flai)) |
| `--version` | [flai self-upgrade](../users/flai-reference.md#flai-self-upgrade) |
| `--wait` | [flai checks tail](../users/flai-reference.md#flai-checks-tail) |
| `--write` | [flai serve agent usage](../users/flai-reference.md#flai-serve-agent-usage) |
| `-y`, `--yes` | every command ([global flags](../users/flai-reference.md#flai)) |
<!-- flags:end -->
