---
title: Settings index
updated: 2026-10-24
status: active
---

# Settings index

Every setting of flai and of the flaiover dashboard, by where it is set. Each row says what the setting is for in a phrase and links to where it is explained. Look a setting up by its name with your browser's find.

| Kind | Where it is kept | Changed with |
|------|------------------|--------------|
| [flai configuration](#flai-configuration) | `~/.flai/config.json`, one per user and host, or the file `--config` or `FLAI_CONFIG` names | `flai config set` for the first table; the `flai serve` commands for the second |
| [Project manifest](#project-manifest) | `system-flow.yaml` at the project's root, committed | By hand for `name`, `description`, `dashboard`; flai for the rest, the strategic agents' settings with `flai manifest set` or on the dashboard's Orchestrator, Planner, and Analyzer pages |
| [Other project files](#other-project-files) | The board's front matter, `.mcp.json`, `.claude/settings.json`, `.github/pull_request_template.md` | The board, `flai order`, by hand |
| [Environment variables](#environment-variables) | The shell flai runs in, or what flai hands a process it starts | Your shell profile or the command line |
| [flaiover container](#flaiover-container) | The dashboard container's environment and mounts | `flai dashboard`, which starts it; by hand only for a container you run yourself |
| [Release signing secrets](#release-signing-secrets) | GitHub Actions secrets of the environment `release` of `bytepunx/system-flow` | `gh secret set --env release`, by whoever holds the release key |
| [Flags](#flags) | The command line | Each command's flags |

When a setting can be made in more than one place, the first of these wins: a flag, then the project's `system-flow.yaml`, then `~/.flai/config.json`, then the default. Secrets are not settings and are listed under [where their files are](index.md#authentication); none of them is ever an environment variable of the dashboard. The release workflows' secrets are under [Release signing secrets](#release-signing-secrets).

## flai configuration

`~/.flai/config.json`, created with the defaults on the first command that needs it. A key flai does not know is refused when the file is read. The whole file is explained in [the flai guide's Configuration](../users/flai.md#configuration). `flai config path` names the file in use.

Keys `flai config set` reaches:

| Key | Default | For |
|-----|---------|-----|
| `template.repo` | `https://github.com/bytepunx/system-flow-template` | The template `flai new`, `flai import`, and `flai upgrade` use when `--template` is not given |
| `template.ref` | `main` | The template's branch, tag, or commit when `--ref` is not given; empty, the default branch, or a release tag follows releases: the newest version tag |
| `dashboard.image` | `ghcr.io/bytepunx/flaiover` | The image `flai dashboard` runs ([Access to the image](index.md#access-to-the-image)); another repository runs `dashboard.tag` as given, reported unsigned |
| `dashboard.tag` | `latest` | For `ghcr.io/bytepunx/flaiover`, the flaiover release to run: `latest` for the newest, `X.Y.Z` for that one. Either is run by the digest its signed digest list names; any other tag, and a release with no list flai can verify, is refused ([What the dashboard verifies](../users/flai.md#what-the-dashboard-verifies)). For another image, its tag |
| `dashboard.port` | `4242` | The host port the dashboard is published on |
| `dashboard.bind` | `0.0.0.0` | The host address it is published on; `127.0.0.1` keeps it to this machine ([Security posture](index.md#security-posture)) |
| `dashboard.no_restart` | `false` | `true` stops flai host restarting a dashboard that is gone or not answering ([The dashboard's watch](index.md#the-dashboards-watch)) |
| `dashboard.allow_unsigned` | `false` | `true` allows unsigned builds ([ADR-0070](../../design/adrs/0070-releases-are-signed-with-a-cosign-key-pair-verified-before-they-are-installed.md)): flai dials a dashboard whose image is on no signed list and accepts a dashboard that sends no release stamp, and `flai dashboard` starts the container with `FLAIOVER_ALLOW_UNSIGNED=1`, so the dashboard accepts a flai with no stamp. A stamp that does not verify is refused either way, and the credential is always checked. While it is in force, `flai serve` logs a `warn`, the statuses say `unsigned allowed`, and the dashboard shows a banner ([Allowing unsigned builds](index.md#allowing-unsigned-builds-s-0239)) |
| `dashboard.allowed_hosts` | empty | Host names, each with an optional port, separated by commas, such as `flow.example.com, board.example.net:8443`, that the dashboard accepts in a browser's `Origin` beside the address it is published on: the name of a tunnel or a domain that reaches it. No scheme, path, or spaces; an IPv6 address in brackets. `flai dashboard` starts the container with `FLAIOVER_ALLOWED_HOSTS`, each entry trimmed and lower-cased, and `restart`, `upgrade`, and flai host's restart give it the list as it is then, so a change applies on the next start. `flai dashboard status` says what the container runs with and when this names other hosts. One machine's fact, so never in the manifest |
| `dashboard.push_key` | empty | Retired and ignored; `flai dashboard` names it while it is set ([Upgrading from a release that mounted the repository](index.md#upgrading-from-a-release-that-mounted-the-repository)) |
| `dashboard.push_known_hosts` | empty | Retired and ignored, as `dashboard.push_key` |
| `cache_dir` | `~/.flai/cache`, or `FLAI_CACHE_DIR` when the file is created | Where git templates are cloned, and where `flai self-upgrade --list` and `--check` keep each release's signature verdict, `releases/verdicts.json` ([What an upgrade verifies](../users/flai.md#what-an-upgrade-verifies)) |
| `author` | your login name | The default owner of new items and `--by` of transitions and acceptances |
| `worktrees.relative_paths` | `false` | Link story worktrees with relative paths, git 2.48 or newer ([Relative worktree links](../users/flai.md#relative-worktree-links-opt-in)) |

Keys only the `flai serve` commands on the host change, or the dashboard's [Configuration pages](../users/flaiover.md#settings) once the `settings` host action is on: host actions and checks on Features, providers on Agents, import folders on File System, and the agent's name and command and the harnesses on Harnesses ([the settings host action](index.md#the-settings-host-action-changing-the-hosts-settings-from-the-dashboard-s-0105)):

| Key | Default | Changed with | For |
|-----|---------|--------------|-----|
| `host_actions.<name>` | absent: every action off | `flai serve enable`, `flai serve disable` | For each host action (`push`, `forge`, `auto-publish`, `agent`, `checks`, `components`, `dashboard`, `host`, `settings`, `plan`, `orchestrate`, `analyze`, `curate`, `merge-queue`, `auto-approve`, `auto-upgrade`), the main checkouts it is on for, or `*` for every project ([The push host action](index.md#the-push-host-action), which lets the board publish; [the forge host action](index.md#the-forge-host-action-pull-requests-on-the-move-to-review-s-0393), which has `flai serve` push each story in review under `acceptance.mode: pull-request` and open or update its pull request, [ADR-0166](../../design/adrs/0166-under-pull-request-acceptance-flai-serve-pushes-a-story-s-branch-and-opens-or.md); [the plan host action](index.md#the-plan-host-action-the-planner-s-0208), which lets you start the planner; `orchestrate`, which runs the orchestrator for as long as it is on ([Running the orchestrator](../users/flai.md#running-the-orchestrator)); [the analyze host action](index.md#the-analyze-host-action-the-analyzer-s-0223), which lets you start the analyzer and `analysis.schedule` start it; [the curate host action](index.md#the-curate-host-action-the-curator-s-0431), which lets you start the curator; [the components host action](index.md#the-components-host-action-s-0438), which lets the Components page start and stop a sub-project's own Vite dev server, [ADR-0163](../../design/adrs/0163-flai-runs-a-sub-project-s-own-vite-dev-server-on-127-0-0-1-behind-the.md); [the merge-queue host action](index.md#the-merge-queue-host-action-stories-in-review-verified-as-one-batch-s-0399), which has `flai serve` build the enqueued stories in review into one batch and verify it once, under local acceptance with `orchestration.permissions.accept_reviews` on, [ADR-0173](../../design/adrs/0173-a-merge-queue-in-flai-serve-with-no-model-runs-behind-the-merge-queue-host.md)). `auto-approve` is shell only, like `auto-publish`: on, it lets a story's agent write a path flai protects in its worktree, such as a `.claude/` file or `.mcp.json` but never `.git`, with no thread, and lets the story's sub-agents make those writes, which `flai guard` refuses them while it is off, so that the story's agent makes them itself; a story whose branch changes such a path is accepted by the operator only ([Writes to protected paths](../users/flai.md#writes-to-protected-paths), [ADR-0106](../../design/adrs/0106-a-story-whose-branch-changes-a-path-claude-code-protects-is-accepted-by-the.md)). `auto-upgrade` is shell only too: on, `flai serve` has `flai host` install the newest flai release tagged in the project's history, and restart on it, whenever the flai serving the project is older; off, it asks you on a thread ([Keeping the host's flai current](index.md#keeping-the-hosts-flai-current), [ADR-0158](../../design/adrs/0158-flai-serve-checks-each-served-project-for-a-newer-flai-release-every-minute.md)). `auto-publish` has `flai push --pending` release first; it is shell only, outside the workflow, and the dashboard's Features page neither lists nor changes it ([Pushing outside the workflow](index.md#accepting-releases-nothing-publishing-reaches-the-remote-s-0087-s-0195)) |
| `agent.command` | none | `flai serve agent set -- ...`, `flai serve agent clear` | What starts a ready story's agent when the story names no harness ([Starting an agent](index.md#starting-an-agent-when-a-story-becomes-ready)) |
| `agent.name` | `agent` | `flai serve agent set --name` | The `FLAI_AGENT` prefix of the agents flai serve starts |
| `agent.attended_minutes` | `6` | `flai serve agent set --attended-minutes` | How recent a sign of someone attending must be, and how long it holds a ready story back |
| `agent.harnesses.<name>.program` | the harness's own program, `claude` for `claude-code` | `flai serve agent harness <name> --program` | The binary a harness runs on this host |
| `agent.harnesses.<name>.args` | the adapter's, `--permission-mode acceptEdits --allowedTools Bash,mcp__flai` for `claude-code` | `flai serve agent harness <name> -- ...`, `--reset` | What an agent that harness starts may do |
| `agent.harnesses.<name>.guard` | empty: on a harness that cannot run `flai guard`, a story whose agent has roles is refused | `flai serve agent harness <name> --guard none`, `--guard ""`, `--reset`; the `command` harness takes it too | `none` is your taking the risk that a story's sub-agents run unguarded on that harness, where a sub-agent could move items, write threads, or commit; no other value is taken ([ADR-0131](../../design/adrs/0131-a-harness-without-a-guard-hook-runs-a-story-s-agent-but-none-of-its-roles.md)) |
| `agent.harnesses.<name>.deny_protected` | `false`: on a harness that cannot ask before a protected write, a start is refused unless the project's `auto-approve` is on | `flai serve agent harness <name> --deny-protected`, `--deny-protected=false`, `--reset`; the `command` harness takes it too | Your say that your arguments for that harness already deny the protected paths in its own vocabulary, which flai cannot read, so flai starts its agent ([ADR-0131](../../design/adrs/0131-a-harness-without-a-guard-hook-runs-a-story-s-agent-but-none-of-its-roles.md)) |
| `agent.providers.<name>.api` | none: no provider is set | `flai serve agent provider <name> --api`, `--reset` | What the provider a story's agent may name speaks: `anthropic-messages`, `openai-chat`, or `openai-responses`. A new entry needs `api`, `base_url`, and `key_env`; a project may override `api`, `base_url`, and `models` under its manifest's `providers.<name>` ([ADR-0129](../../design/adrs/0129-a-story-s-agent-names-a-provider-beside-its-harness-and-model-and-the-providers.md)) |
| `agent.providers.<name>.base_url` | none | `flai serve agent provider <name> --base-url` | Where the provider answers, an absolute `http` or `https` URL with no user or password in it |
| `agent.providers.<name>.key_env` | none | `flai serve agent provider <name> --key-env` | The name of the environment variable on this host that holds the provider's key, which `flai serve` reads when it starts the agent. The key itself is never written to the configuration, and no flag takes it |
| `agent.providers.<name>.spend_key_env` | none: the spend log is read with the agent's own key | `flai serve agent provider <name> --spend-key-env`, `--spend-key-env ""` to clear | The name of the environment variable on this host that holds a key that can read the gateway's spend log where the agent's own key cannot, such as LiteLLM's master key; optional. The key itself is never written to the configuration, and no flag takes it ([ADR-0132](../../design/adrs/0132-a-story-s-cost-comes-from-its-provider-s-spend-log-when-it-has-a-provider-from.md)) |
| `agent.providers.<name>.models.<name>` | none: the harness's own names | `flai serve agent provider <name> --model <alias>=<name>`, repeated; given at all, it replaces the entry's aliases | A model alias of the harness, such as `opus`, and the model's name on the provider, such as `anthropic/claude-opus-5-5` |
| `forges.<name>.token_env` | none: `GITHUB_TOKEN` for `github` | `flai serve forge <name> --token-env`, `--reset` | The name of the environment variable on this host that holds the forge's token, for a project whose `acceptance.mode` is `pull-request`; `github` is the only forge. flai reads it only when it calls the forge, as `flai forge check`, `flai forge push`, and `flai serve` under the `forge` host action do, and no agent `flai serve` starts gets the variable. The token itself is never written to the configuration, and no flag takes it ([A forge token for pull-request acceptance](index.md#a-forge-token-for-pull-request-acceptance-s-0392), [ADR-0165](../../design/adrs/0165-a-project-chooses-local-or-pull-request-acceptance-in-its-manifest-flai-reaches.md)) |
| `agent.auto_restarts` | `2` | `flai serve agent set --auto-restarts N` | How many times `flai serve` restarts a story's agent on its own after it ends with its story in progress, failed and not asking; `0` turns the automatic restart off, and past the limit `flai serve` opens one thread on the story for you ([ADR-0108](../../design/adrs/0108-flai-serve-restarts-a-story-s-agent-that-ended-with-its-story-in-progress-up-to.md)) |
| `checks.commands[].name` | none: the manifest's `checks` apply | `flai serve checks set --name`, `flai serve checks clear` | The name of one check a story in review is run with ([The checks host action](index.md#the-checks-host-action-s-0082)) |
| `checks.commands[].command` | none | `flai serve checks set --name <name> -- ...` | That check's argument list, never run through a shell |
| `checks.timeout_minutes` | `15` | `flai serve checks timeout` | The bound on one run of every check together |
| `components.idle` | `30m` | `flai serve components idle` | How long a component dev server goes without a `components.status` before `flai serve` stops it, a Go duration such as `45m`, for every project on this host and applied without a restart ([The components host action](index.md#the-components-host-action-s-0438), [ADR-0163](../../design/adrs/0163-flai-runs-a-sub-project-s-own-vite-dev-server-on-127-0-0-1-behind-the.md)) |
| `components.max_servers` | `3` | `flai serve components max-servers` | How many component dev servers `flai serve` runs or starts at once per project, the main checkout's and every story worktree's together; a start beyond it is refused, naming the ones running. For every project on this host, applied without a restart ([The components host action](index.md#the-components-host-action-s-0438), [ADR-0163](../../design/adrs/0163-flai-runs-a-sub-project-s-own-vite-dev-server-on-127-0-0-1-behind-the.md)) |
| `import_roots` | none | `flai serve import add`, `flai serve import remove` | Folders whose git repositories the board offers to import, and whose repositories with a `system-flow.yaml` `flai serve` serves ([Importing repositories](index.md#importing-repositories-from-the-board-s-0098)) |

Beside the file, in the folders `serve` and `host`, flai keeps state, tokens, and logs, not settings; the [backup runbook](runbooks/backup.md) lists them. One of them, `serve/claude-checks.json`, records each Claude Code version flai serve checked against `permission_prompt`, keyed by version, with `checked`, `program`, `outcome` (`passed` or `failed`), and `error`; flai never checks a version with an entry again, so remove the entry to have it checked again ([Writes to protected paths](../users/flai.md#writes-to-protected-paths)).

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
| `template.ref` | set by `flai new`, `flai import`, or `flai upgrade` | The template's branch, tag, or commit last applied; edited by hand, the next `flai upgrade` applies it, asking when it is not the newest release |
| `template.version` | set by flai | The template version last applied; edited by hand, the next `flai upgrade` applies its `v<version>` tag, asking when it is not the newest release |
| `template.applied` | set by flai | When it was applied |
| `layout.<name>` | `design`, `docs`, `wip` | The folder names; `design`, `docs`, and `wip` are required |
| `projects[].name` | none | A releasable component's name, the prefix of its release tags |
| `projects[].path` | none | Its folder at the root |
| `projects[].kind` | none | `go`, `sveltekit`, `svelte`, `template`, and so on; `template` is released by its version file, not a tag. flai reads the components of a `sveltekit` or `svelte` sub-project, a Svelte project that is not a SvelteKit app |
| `projects[].tags` | none | Story tags that mean a story delivers to it |
| `projects[].components.include` | `src/lib/**/*.svelte`, so routes are left out | For a `sveltekit` or `svelte` sub-project, the globs of its component files, relative to its folder, with the syntax of `claims.shared`. A glob that does not parse, or one that leaves the sub-project, is a `flai check` error, and so is `components` on a sub-project of any other kind. Changed with `flai components config <sub-project> [--include] [--exclude] [--setup] [--unset]`, which keeps the manifest's other keys, order, and comments and refuses what the check reports, or by hand ([Project manifest](../../design/system/project-manifest.md)) |
| `projects[].components.exclude` | none | The globs of files `include` finds that are not components to list, relative to the sub-project's folder, as `include`. Changed with `flai components config` or by hand |
| `projects[].components.setup` | none | A module of the sub-project, relative to its folder, that the harness imports before any component, such as the app's stylesheet `src/app.css`. A path that is not clean, leaves the sub-project, or names no file is a `flai check` error. Changed with `flai components config` or by hand |
| `projects[].install.command` | none: `flai accept` installs nothing for the sub-project | For a sub-project of any kind, the program and its arguments, as a list, that installs its dependencies, run with no shell in its folder. `flai accept` runs it in the main checkout after an acceptance that changed a path `paths` matches; a failed install is a warning, not a failed acceptance. Required when `install` is set, its first element not empty, or it is a `flai check` error. This repository sets flaiover's to `["../scripts/flaiover-install.sh", "--if-needed"]`. Changed by hand ([Project manifest](../../design/system/project-manifest.md), [ADR-0185](../../design/adrs/0185-flai-accept-runs-a-sub-project-s-install-command-in-the-main-checkout-when-the.md)) |
| `projects[].install.paths` | none | The globs of the files whose change needs `command` run, relative to the sub-project's folder, as `components.include`'s, such as `package.json` and `pnpm-lock.yaml`. Required when `install` is set, at least one; a glob that does not parse, or one that leaves the sub-project, is a `flai check` error. Changed by hand |
| `dashboard.image` | the configuration's | The image, for this project, over `dashboard.image` in the configuration |
| `dashboard.tag` | the configuration's | The flaiover release, `latest` or `X.Y.Z`, run by digest, or another image's tag, likewise |
| `dashboard.port` | the configuration's | The host port, likewise |
| `dashboard.bind` | the configuration's | The host address, likewise |
| `dashboard.autocommit` | `true` | Commit documents saved from the dashboard, one path per commit |
| `dashboard.allow_unsigned` | the configuration's | `true` or `false` allows unsigned builds for this project, or refuses them, over `dashboard.allow_unsigned` in the configuration; unset leaves the configuration's |
| `dashboard.notify_url` | unset | A webhook the dashboard posts new inbox entries to |
| `checks[].name` | none | A check a story in review is run with, when the host names none ([The checks host action](index.md#the-checks-host-action-s-0082)) |
| `checks[].command` | none | Its argument list |
| `tests[].name` | none; with no `tests` key, one tier `test` | A test tier's name, unique among them. `flai test` runs the tiers, cheapest first, for the paths each selects. With no `tests` key, one `plain` tier named `test` runs `scripts/test.sh` for every path when that file exists; `tests: []` means none. A project made from the template starts with `test`, `integration`, and `smoke`, its three scripts. Changed whole with `flai manifest set tests='<JSON list>'`, such as `[{"name":"test","command":["scripts/test.sh"],"paths":["**"]}]`, or by hand ([Project manifest](../../design/system/project-manifest.md)) |
| `tests[].command` | none | Its argument list, run without a shell; `{packages}`, `{importers}`, or `{files}`, each as an argument of its own, stands for what its paths selected: the Go packages, those packages and every package of the module that imports them, or the files ([Project manifest](../../design/system/project-manifest.md)) |
| `tests[].dir` | the root | The folder it runs in, relative to the root and inside it |
| `tests[].paths` | none; required unless `all_only` | Glob patterns from the root, as `claims.shared`'s, that select the tier for a changed or named path; one beginning with `!` takes paths out |
| `tests[].format` | `plain` | How its output becomes findings: `go-test-json`, `vitest-json`, `golangci-json`, `gofmt-list`, or `plain`, whose one finding is the exit status after the last 20 lines of output and, ahead of them, the `go test` failure lines those left out (`--- FAIL:` with the lines under it, `FAIL<tab><package>`, `panic:`) ([Project manifest](../../design/system/project-manifest.md)) |
| `tests[].all_only` | `false` | Run the tier only under `flai test --all`, as integration and smoke tests are |
| `tests[].all_command` | `command` | The argument list run instead of `command` under `--all`; it takes no placeholder |
| `tests[].fix_command` | none | The argument list, with the placeholders `command` takes, that rewrites in place the files the tier selects: before `flai task done` commits them, after each clean `flai stream sync` over what the branch changes, and before `flai accept` commits the main checkout; its exit status is ignored, but what it leaves for the tier's `command` refuses an acceptance ([ADR-0140](../../design/adrs/0140-flai-accept-refuses-an-acceptance-whose-files-a-tier-with-a-fix-command-still.md), [ADR-0139](../../design/adrs/0139-a-test-tier-s-fix-command-rewrites-the-files-a-story-changes-before-flai-task.md)) |
| `tests[].covers` | none | The tiers this one runs everything of, by name, which are skipped when it is selected with them |
| `tests[].queue` | `false` | Leave the tier to the merge queue's batch. While the queue can batch for the project, with the `merge-queue` host action on for it in the configuration the running flai reads, `orchestration.permissions.accept_reviews` on, and acceptance local, `flai verify` and the close-out list the tier as `queued` and do not run it; the batch runs it. `flai verify --all` runs it. This repository marks `integration` and `smoke`. An older flai ignores the key ([ADR-0174](../../design/adrs/0174-while-the-merge-queue-can-batch-flai-verify-and-the-close-out-leave-the-test.md), [Project manifest](../../design/system/project-manifest.md)) |
| `agent.harness` | none | The harness every new story gets, such as `claude-code` ([Starting an agent](index.md#starting-an-agent-when-a-story-becomes-ready)) |
| `agent.model` | none | The model every new story gets |
| `agent.provider` | none; the harness's own default | The provider every new story's model calls go to: the name of an entry the host sets with `flai serve agent provider <name>` (`agent.providers` in the [flai configuration](#flai-configuration)), with the project's override under `providers.<name>` laid over it. `flai serve` refuses to start an agent whose provider the host has no entry for ([ADR-0129](../../design/adrs/0129-a-story-s-agent-names-a-provider-beside-its-harness-and-model-and-the-providers.md)) |
| `prime.budget` | `80KB` | The size a story's context pack fits, for `flai prime --story` and the MCP `prime` tool; bytes, or a number with `KB` or `MB`; `--budget` overrides it for one run ([Prime a session](../users/flai.md#prime-a-session)). It bounds the pack's text, not one tool result: the MCP `prime` tool returns the pack in parts of at most 40,000 bytes whatever the budget ([ADR-0104](../../design/adrs/0104-the-mcp-prime-tool-returns-a-context-pack-in-parts-of-at-most-40-000-bytes-of.md)) |
| `agent.config.<name>` | none | Options for the harness; `claude-code` takes `effort` (`low`, `medium`, `high`, `xhigh`, `max`), `max_budget_usd`, and `fallback_model` |
| `issues.story_after` | `168h` | How long an issue may stay open with no open story linking it before `flai check` warns (`issues.no-story`); a Go duration such as `24h`, or `0` to turn the warning off ([Record recurring friction](../users/flai.md#record-recurring-friction)) |
| `planning.currency` | `USD` | The ISO 4217 code, three capital letters such as `EUR`, of every amount in an item's `cost_of_delay` and of `planning.hour_rate`. Changed by hand alone: `flai manifest set` refuses it, since changing it re-denominates every amount on the items, which flai does not convert |
| `planning.hour_rate` | unset | What an hour of work costs, in `planning.currency`; a number of zero or more. Unset means unknown, not free. Changed on the [Planner page](../users/flaiover.md#planner) or with `flai manifest set` |
| `planning.cycle` | `168h` | The period a cost of delay's `time_lost_per_cycle` is counted over; a Go duration longer than zero. Changed on the [Planner page](../users/flaiover.md#planner) or with `flai manifest set` |
| `planning.default_duration` | `1h` | The work `flai forecast` gives a story when fewer than three done stories with usage are on record to forecast it from; a Go duration longer than zero. Changed on the [Planner page](../users/flaiover.md#planner) or with `flai manifest set` |
| `planning.agent.harness` | `agent.harness` | The planner's harness, over the project's `agent` ([Running the planner](../users/flai.md#running-the-planner)). Changed as `planning.agent`, whole, on the [Planner page](../users/flaiover.md#planner) or with `flai manifest set` |
| `planning.agent.model` | `agent.model` | The planner's model. Changed as `planning.agent`, whole, on the [Planner page](../users/flaiover.md#planner) or with `flai manifest set` |
| `planning.agent.provider` | `agent.provider` | The planner's provider. Changed as `planning.agent`, whole, on the [Planner page](../users/flaiover.md#planner) or with `flai manifest set` |
| `planning.agent.config.<name>` | `agent.config` | The planner's options for its harness, merged key by key over `agent.config`. Changed as `planning.agent`, whole, on the [Planner page](../users/flaiover.md#planner) or with `flai manifest set` |
| `planning.agent.roles.<name>.harness` | `agent.roles` | A sub-agent role's harness in the planner's session, merged role by role over `agent.roles`. Changed with `flai manifest set planning.agent=<json>` or by hand; the Planner page keeps them as they are |
| `planning.agent.roles.<name>.model` | `agent.roles` | A sub-agent role's model in the planner's session. Changed with `flai manifest set planning.agent=<json>` or by hand; the Planner page keeps them as they are |
| `planning.agent.roles.<name>.config.<name>` | `agent.roles` | A sub-agent role's options in the planner's session; `claude-code` takes none. Changed with `flai manifest set planning.agent=<json>` or by hand; the Planner page keeps them as they are |
| `claims.shared` | empty; `design/adrs` and `design/issues`, under the design folder, in a project made from the template | Glob patterns of paths that many stories change in separate sections or new files: an overlap inside one holds no ready story and is no `wip.overlap` ([ADR-0096](../../design/adrs/0096-a-story-in-review-holds-nothing-an-overlap-inside-the-manifest-s-shared-paths.md)). Paths are from the repository's root; `*` matches within one segment, `**` zero or more whole segments, `?` one character, and a plain path covers itself and everything below it. Change it with `flai shared add` and `remove`, the dashboard's [File System page](../users/flaiover.md#file-system) or `settings.shared` in the host API while the `settings` host action is on, or the MCP tool `shared_paths_edit`. Only your own session may: `flai guard` refuses every sub-agent and every session `flai serve` starts. `flai shared list` and `check`, `settings.get` and `settings.shared_check`, and the MCP tool `shared_paths` read it ([Shared paths](../users/flai.md#shared-paths)) |
| `planning.replan` | `deterministic` | What `flai serve` does, while the `plan` host action is on, when a story is accepted or cancelled or the pull order changes: `never` does nothing; `deterministic` plays the board out again with no agent and moves each forecast's delivery; `agent` does that and starts the planner for each story whose delivery moved ([Running the planner](../users/flai.md#running-the-planner)). Changed on the [Planner page](../users/flaiover.md#planner) or with `flai manifest set` |
| `planning.schedule` | unset | When `flai serve`, while the `plan` host action is on, runs the planner over every ready story: a five-field cron expression in UTC such as `0 6 * * 1-5`, or `daily` for 00:00 UTC. Unset means no schedule. Changed on the [Planner page](../users/flaiover.md#planner) or with `flai manifest set` |
| `orchestration.policy` | `fifo` | How the ready column is ordered, by the names `flai order --by` takes: `cod` by cost of delay, the largest first; `wsjf` by cost of delay over forecast duration, the largest first; `throughput` by forecast duration, the shortest first; `fifo` leaves the operator's order alone. Changed on the [Orchestrator page](../users/flaiover.md#orchestrator) or with `flai manifest set` |
| `orchestration.release.policy` | `judgement` | When accepted stories not yet released are due a release: `judgement` leaves it to the operator and is never met by itself; `threshold` uses `value` and `count`; `theme` uses `epic` or `tag`. Changed on the [Orchestrator page](../users/flaiover.md#orchestrator) or with `flai manifest set` |
| `orchestration.release.value` | unset | For `threshold`: the unreleased cost of delay per week, in `planning.currency`, at which a release is due; zero or more. Give it, `count`, or both. Changed on the [Orchestrator page](../users/flaiover.md#orchestrator) or with `flai manifest set` |
| `orchestration.release.count` | unset | For `threshold`: the number of accepted stories not yet released at which a release is due; zero or more. Changed on the [Orchestrator page](../users/flaiover.md#orchestrator) or with `flai manifest set` |
| `orchestration.release.epic` | unset | For `theme`: an epic ID such as `E-0001`; a release is due when every story of it is accepted. Give it or `tag`, not both. Changed on the [Orchestrator page](../users/flaiover.md#orchestrator) or with `flai manifest set` |
| `orchestration.release.tag` | unset | For `theme`: a tag; a release is due when every story with it is accepted. Changed on the [Orchestrator page](../users/flaiover.md#orchestrator) or with `flai manifest set` |
| `orchestration.release.whole_epics` | `false` | For every policy: no release is due while a story accepted and not yet released belongs to an epic in neither `review` nor `done`. `flai release --evaluate` names those stories with their epics; under `judgement` they are what the orchestrator must not publish. A story with no epic is never held. Changed on the [Orchestrator page](../users/flaiover.md#orchestrator) or with `flai manifest set` |
| `orchestration.permissions.plan_backlog_epics` | `false` | The orchestrator may start the planner for an epic `flai plan --candidates` lists; the `plan` host action must be on too ([Running the orchestrator](../users/flai.md#running-the-orchestrator)). Changed on the [Orchestrator page](../users/flaiover.md#orchestrator) or with `flai manifest set` |
| `orchestration.permissions.plan_backlog_stories` | `false` | The orchestrator may start the planner, one at a time, for a story in the backlog that `flai plan --candidates` lists for lacking touches, a forecast duration, a cost of delay value, or a task; answer and resolve the threads that story's planner opens; and give a backlog story its cost of delay inputs, and change nothing else of it, when neither it nor its epic has any. Those inputs are money decisions it then makes without you. The `plan` host action must be on too ([Running the orchestrator](../users/flai.md#running-the-orchestrator)). Changed on the [Orchestrator page](../users/flaiover.md#orchestrator) or with `flai manifest set` |
| `orchestration.permissions.finalize_drafts` | `false` | The orchestrator may finalize a draft story that `flai promote --drafts` finds complete, and change nothing else of it. Changed on the [Orchestrator page](../users/flaiover.md#orchestrator) or with `flai manifest set` |
| `orchestration.permissions.promote_to_ready` | `false` | The orchestrator may move a story `flai promote --candidates` lists to `ready`, while ready is under its WIP limit. Changed on the [Orchestrator page](../users/flaiover.md#orchestrator) or with `flai manifest set` |
| `orchestration.permissions.order_ready` | `false` | The orchestrator may write the ready column's order, `flai order --by <policy> --apply`, which keeps a story placed by hand in the last day; it never places a story by hand. Changed on the [Orchestrator page](../users/flaiover.md#orchestrator) or with `flai manifest set` |
| `orchestration.permissions.answer_threads` | `off` | How the orchestrator replies on threads: `off` never; `recommend` with a recommendation for the operator; `autonomous` with an answer of its own. `flai guard` lets it reply with either of the last two; its prompt says which. Changed on the [Orchestrator page](../users/flaiover.md#orchestrator) or with `flai manifest set` |
| `orchestration.permissions.accept_reviews` | `false` | The orchestrator may accept a story in review, `flai accept --by orchestrator`, once its verifier passed and nothing blocks it; the `orchestrate` host action must be on too ([When the orchestrator accepts](../users/flai.md#when-the-orchestrator-accepts)). A ready story waiting on that acceptance, by `after:` or a full review column, then starts without you ([What waits on an acceptance](../../design/system/workflow.md#what-waits-on-an-acceptance-s-0296)). The merge queue runs only while it is on, whether or not `orchestrate` is ([the merge-queue host action](index.md#the-merge-queue-host-action-stories-in-review-verified-as-one-batch-s-0399)). Changed on the [Orchestrator page](../users/flaiover.md#orchestrator) or with `flai manifest set` |
| `orchestration.permissions.publish` | `false` | The orchestrator may publish what is accepted and not yet released, through the MCP tool `release_publish` alone, when `orchestration.release` allows it; the `push` host action must be on too ([When the orchestrator publishes](../users/flai.md#when-the-orchestrator-publishes)). Changed on the [Orchestrator page](../users/flaiover.md#orchestrator) or with `flai manifest set` |
| `orchestration.permissions.edit_subprojects` | `false` | The orchestrator may add, change, or remove a sub-project, an entry of `projects`, when a story or a thread asks for it, with `flai subproject add`, `set`, or `remove` and `--autocommit`, committed on the main branch without your review; adding one creates its folder ([Running the orchestrator](../users/flai.md#running-the-orchestrator), [ADR-0195](../../design/adrs/0195-with-edit-subprojects-on-the-orchestrator-adds-changes-and-removes-a-sub.md)). Changed on the [Orchestrator page](../users/flaiover.md#orchestrator) or with `flai manifest set` |
| `orchestration.agent.harness` | `agent.harness` | The orchestrator's harness, over the project's `agent`. Changed by hand: `flai manifest set` does not write `orchestration.agent` |
| `orchestration.agent.model` | `agent.model` | The orchestrator's model. Changed by hand: `flai manifest set` does not write `orchestration.agent` |
| `orchestration.agent.provider` | `agent.provider` | The orchestrator's provider. Changed by hand: `flai manifest set` does not write `orchestration.agent` |
| `orchestration.agent.config.<name>` | `agent.config` | The orchestrator's options for its harness, merged key by key over `agent.config`. Changed by hand: `flai manifest set` does not write `orchestration.agent` |
| `orchestration.agent.roles.<name>.harness` | `agent.roles` | A sub-agent role's harness in the orchestrator's session, merged role by role over `agent.roles`. Changed by hand: `flai manifest set` does not write `orchestration.agent` |
| `orchestration.agent.roles.<name>.model` | `agent.roles` | A sub-agent role's model in the orchestrator's session. Changed by hand: `flai manifest set` does not write `orchestration.agent` |
| `orchestration.agent.roles.<name>.config.<name>` | `agent.roles` | A sub-agent role's options in the orchestrator's session; `claude-code` takes none. Changed by hand: `flai manifest set` does not write `orchestration.agent` |
| `analysis.schedule` | unset | When `flai serve`, while the `analyze` host action is on, runs the analyzer with no focus: a five-field cron expression in UTC such as `0 6 * * 1`, or `daily` for 00:00 UTC. Unset means no schedule. Changed on the [Analyzer page](../users/flaiover.md#analyzer) or with `flai manifest set` |
| `analysis.agent.harness` | `agent.harness` | The analyzer's harness, over the project's `agent`. Changed as `analysis.agent`, whole, on the [Analyzer page](../users/flaiover.md#analyzer) or with `flai manifest set` |
| `analysis.agent.model` | `agent.model` | The analyzer's model. Changed as `analysis.agent`, whole, on the [Analyzer page](../users/flaiover.md#analyzer) or with `flai manifest set` |
| `analysis.agent.provider` | `agent.provider` | The analyzer's provider. Changed as `analysis.agent`, whole, on the [Analyzer page](../users/flaiover.md#analyzer) or with `flai manifest set` |
| `analysis.agent.config.<name>` | `agent.config` | The analyzer's options for its harness, merged key by key over `agent.config`. Changed as `analysis.agent`, whole, on the [Analyzer page](../users/flaiover.md#analyzer) or with `flai manifest set` |
| `analysis.agent.roles.<name>.harness` | `agent.roles` | A sub-agent role's harness in the analyzer's session, merged role by role over `agent.roles`. Changed with `flai manifest set analysis.agent=<json>` or by hand; the Analyzer page keeps them as they are |
| `analysis.agent.roles.<name>.model` | `agent.roles` | A sub-agent role's model in the analyzer's session. Changed with `flai manifest set analysis.agent=<json>` or by hand; the Analyzer page keeps them as they are |
| `analysis.agent.roles.<name>.config.<name>` | `agent.roles` | A sub-agent role's options in the analyzer's session; `claude-code` takes none. Changed with `flai manifest set analysis.agent=<json>` or by hand; the Analyzer page keeps them as they are |
| `digest.agent.harness` | `agent.harness` | The curator's harness, over the project's `agent`. Changed as `digest.agent`, whole, with `flai manifest set` |
| `digest.agent.model` | `agent.model` | The curator's model. Changed as `digest.agent`, whole, with `flai manifest set` |
| `digest.agent.provider` | `agent.provider` | The curator's provider. Changed as `digest.agent`, whole, with `flai manifest set` |
| `digest.agent.config.<name>` | `agent.config` | The curator's options for its harness, merged key by key over `agent.config`. Changed as `digest.agent`, whole, with `flai manifest set` |
| `digest.agent.roles.<name>.harness` | `agent.roles` | A sub-agent role's harness in the curator's session, merged role by role over `agent.roles`. Changed with `flai manifest set digest.agent=<json>` or by hand |
| `digest.agent.roles.<name>.model` | `agent.roles` | A sub-agent role's model in the curator's session. Changed with `flai manifest set digest.agent=<json>` or by hand |
| `digest.agent.roles.<name>.config.<name>` | `agent.roles` | A sub-agent role's options in the curator's session; `claude-code` takes none. Changed with `flai manifest set digest.agent=<json>` or by hand |
| `digest.schedule` | unset | When `flai serve`, while the `curate` host action is on, runs the curator besides after each acceptance: a five-field cron expression in UTC such as `0 6 * * 1`, or `daily` for 00:00 UTC. Unset means no schedule ([ADR-0199](../../design/adrs/0199-the-curator-runs-after-acceptances-coalesced-and-on-digest-schedule-marks.md)). Changed with `flai manifest set` |
| `digest.compact_after` | `720h` | How long an active digest entry may go with no work item, narrative, or thread citing it before `flai digest compact` turns it into a one-line entry: a Go duration such as `720h`; `0` turns compaction off. Changed with `flai manifest set` |
| `design.tokens` | `design/system/design-system.tokens.json`, under the `design` layout folder | The repository path of the project's design token file, a DTCG file, beside its guidance in `design/system/design-system.md` ([ADR-0146](../../design/adrs/0146-a-project-s-design-tokens-are-a-dtcg-file-the-manifest-names-beside-its.md)). Relative to the root and inside it: an absolute path, one that leaves the repository, or one that is not clean stops every command, naming the key. Changed by hand ([Project manifest](../../design/system/project-manifest.md)) |
| `design.stylesheet` | unset | The repository path of the Tailwind v4 stylesheet `flai design css --write` writes the design tokens to, a file under a sub-project, such as `web/src/theme.css` ([ADR-0159](../../design/adrs/0159-the-generated-design-stylesheet-puts-the-dark-values-under-data-theme-dark-and.md)). Unset means none, and nothing is written or checked. An absolute path, one that leaves the repository, or one that is not clean stops every command, naming the key; one outside every sub-project, and a stylesheet that has fallen behind the tokens (`design.stale`), `flai check` reports. Changed by hand ([Project manifest](../../design/system/project-manifest.md)) |
| `acceptance.mode` | `local` | How a story in review is accepted: `local` merges its branch on this machine, as `flai accept` does; `pull-request` accepts it through a pull request on the project's forge. Under `pull-request`, with the `forge` host action on, `flai serve` pushes each story in review and opens or updates its pull request ([the forge host action](index.md#the-forge-host-action-pull-requests-on-the-move-to-review-s-0393)); nothing merges one yet, and `flai accept` merges locally under either mode until later releases of E-0021. A value outside these stops every command that reads the manifest, naming the key, rather than falling back to `local`. Changed with `flai manifest set` or by hand ([Accepting through a pull request](../users/flai.md#accepting-through-a-pull-request), [ADR-0165](../../design/adrs/0165-a-project-chooses-local-or-pull-request-acceptance-in-its-manifest-flai-reaches.md)) |
| `acceptance.forge` | `github` | The forge pull-request acceptance goes through, at the repository the `origin` remote names; `github` is the only one. Allowed and unused under `local`. Its token is the host's, named by `forges.<name>.token_env` in the [flai configuration](#flai-configuration). Changed with `flai manifest set` or by hand |
| `acceptance.merge_method` | `merge` | How the forge merges an accepted pull request: `merge` with a merge commit, `squash` into one commit, or `rebase` the commits onto the base. Allowed and unused under `local`. Changed with `flai manifest set` or by hand |
| `merge_queue.batch_max` | `4` | The most enqueued stories the merge queue takes into one batch, verified together once, a whole number of one or more; the queue runs in `flai serve` behind the `merge-queue` host action, while `orchestration.permissions.accept_reviews` is on, under local acceptance only ([the merge-queue host action](index.md#the-merge-queue-host-action-stories-in-review-verified-as-one-batch-s-0399), [ADR-0173](../../design/adrs/0173-a-merge-queue-in-flai-serve-with-no-model-runs-behind-the-merge-queue-host.md)). Zero, a negative number, or a value that is not a whole number stops every command that reads the manifest, naming the key. Changed with `flai manifest set` or by hand |
| `flai.minimum` | unset | The oldest flai release that may read the project, `X.Y.Z`; an older one stops before reading any item and names the version needed. Publishing a flai release that changes the front-matter fields flai reads raises it ([Keeping the host's flai current](index.md#keeping-the-hosts-flai-current)) |
| `agent.roles.<name>.harness` | none | The harness of a sub-agent role (`explore`, `verify`); for `claude-code`, only `claude-code`, since a sub-agent runs in the story's session |
| `agent.roles.<name>.model` | none | The model a sub-agent role runs, over the one its definition in `.claude/agents/` names, such as `haiku` or `claude-sonnet-5-5` |
| `agent.roles.<name>.config.<name>` | none | Options for a sub-agent role; `claude-code` takes none |
| `providers.<name>.api` | the host's | The API the provider `<name>` speaks for this project, over the host's entry for it: `anthropic-messages`, `openai-chat`, or `openai-responses`. An override for a provider the host has no entry for is refused when an agent naming it starts ([ADR-0129](../../design/adrs/0129-a-story-s-agent-names-a-provider-beside-its-harness-and-model-and-the-providers.md)). Changed by hand |
| `providers.<name>.base_url` | the host's | The provider's base URL for this project, over the host's: an absolute `http` or `https` URL with no user or password in it. Changed by hand |
| `providers.<name>.key_env` | refused | Not a project setting: the name of the variable that holds the key is the host's, set with `flai serve agent provider <name> --key-env`, and `flai check` refuses it here |
| `providers.<name>.spend_key_env` | refused | Not a project setting: the name of the variable that holds the spend log's key is the host's, set with `flai serve agent provider <name> --spend-key-env`, and `flai check` refuses it here |
| `providers.<name>.models.<name>` | the host's | The harness's model aliases on the provider, such as `opus: anthropic/claude-opus-5-5`, merged alias by alias over the host's, the project's winning. Changed by hand |

The strategic agents' settings, the `orchestration` keys but `agent`, the `planning` and `analysis` keys but `planning.currency`, and the `digest` keys, are changed with `flai manifest set <key>=<value>` or `--unset <key>` in a shell, which checks the whole manifest as `flai check` does and refuses with the field and the reason, writing nothing ([flai manifest set](../users/flai.md#changing-the-strategic-agents-settings)). The dashboard's Orchestrator, Planner, and Analyzer pages change the same keys through it, the `digest` keys aside, while the `settings` host action is on for the project, and show them read-only while it is off ([ADR-0101](../../design/adrs/0101-with-the-settings-host-action-on-the-dashboard-edits-the-strategic-agents.md), [the settings host action](index.md#the-settings-host-action-changing-the-hosts-settings-from-the-dashboard-s-0105)). A change to a permission holds from the orchestrator's next call, with no restart. `flai manifest set` also writes `tests` and the `acceptance` keys, which no dashboard page changes.

`dashboard.image`, `tag`, `port`, and `bind` matter only when they decide what the one shared container is started with, which is the first `flai dashboard` on the host. A story's own `agent:` front matter, copied from `agent` when the story is made, is the story's, not a setting; `flai edit` and the story's page change it.

## Other project files

| Setting | Where | Default | For |
|---------|-------|---------|-----|
| `wip_limits` | the front matter of `wip/kanban/board.md` | as the template sets them | The WIP limit of each column; flai serve starts agents only while `in-progress` has room ([design/system/workflow.md](../../design/system/workflow.md)) |
| `order` | the front matter of `wip/kanban/board.md` | empty | The pull order of `ready` and `backlog` stories; drag cards on the board or run `flai order` |
| `<!-- flai:story -->`, `criteria`, `definition-of-done`, `verify` | `.github/pull_request_template.md` | all four, in a project made from the template | The places flai fills in a story's pull request under pull-request acceptance, each between `<!-- flai:<part> -->` and `<!-- /flai:<part> -->`: the story's ID, its acceptance criteria, the definition of done, kept as written, and its last `flai verify` result. A part not marked is added at the end under its own heading ([The pull request on the move to review](../users/flai.md#the-pull-request-on-the-move-to-review), [ADR-0166](../../design/adrs/0166-under-pull-request-acceptance-flai-serve-pushes-a-story-s-branch-and-opens-or.md)) |
| `mcpServers.flai` | `.mcp.json` | `flai mcp` on stdio | How an agent on the host reaches flai's MCP server ([MCP over HTTP](index.md#mcp-over-http)) |
| `hooks.PreToolUse` | `.claude/settings.json` | `flai guard` before `Bash` and `mcp__flai__.*`, and before `Edit\|MultiEdit\|Write\|NotebookEdit` when `FLAI_ROLE` is `plan`, `orchestrate`, `analyze`, or `curate`, or `FLAI_STORY` is set | Refuses a Claude Code sub-agent's writes to work items, threads, and history ([Sub-agents](../users/flai.md#sub-agents)), and, while `auto-approve` is off, its writes to a path flai protects, such as a file in a `.claude/` folder ([Writes to protected paths](../users/flai.md#writes-to-protected-paths)), holds the planner to planning ([Running the planner](../users/flai.md#running-the-planner)), holds the orchestrator to `orchestration.permissions` ([Running the orchestrator](../users/flai.md#running-the-orchestrator)), holds the analyzer to its report under `design/analysis/` ([Running the analyzer](../users/flai.md#running-the-analyzer)), holds the curator to `flai digest new` and `edit` ([Running the curator](../users/flai.md#running-the-curator)), and refuses a story's agent `wait_for_events` while a sub-agent of its session runs and no thread on its story is open ([Sub-agents](../users/flai.md#sub-agents)) |
| `hooks.SubagentStart` | `.claude/settings.json` | `flai guard` when `FLAI_STORY` is set, its output discarded and its exit ignored | Records a story's agent's sub-agent as running, in `.flai-cache/guard/<session_id>.json` in the main checkout, for the refusal of `wait_for_events` ([Sub-agents](../users/flai.md#sub-agents), S-0285) |
| `hooks.SubagentStop` | `.claude/settings.json` | `flai guard` when `FLAI_STORY` is set, its output discarded and its exit ignored | Records that the sub-agent no longer runs; the record is removed when none does |

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
| The variable a provider's `key_env` names, such as `LITELLM_API_KEY` | none | The provider's key, read by `flai serve` at each start of an agent that names the provider; a start is refused while it is unset or empty. Set it in the environment `flai host` or `flai serve` runs in ([Starting an agent](index.md#starting-an-agent-when-a-story-becomes-ready)) |
| The variable a provider's `spend_key_env` names, such as `LITELLM_MASTER_KEY` | none: the spend log is read with the key `key_env` names | A key that can read the gateway's spend log, read by `flai serve` when it prices a run through the provider; never passed to the agent. Unset or empty, the run is priced by the harness or estimated, and `flai serve` logs why ([ADR-0132](../../design/adrs/0132-a-story-s-cost-comes-from-its-provider-s-spend-log-when-it-has-a-provider-from.md)) |
| The variable a forge's `token_env` names, `GITHUB_TOKEN` for `github` by default | none | The forge's token, for a project that accepts through a pull request; read only when flai calls the forge, such as `flai forge check`, which names the variable and says whether the token can read and push, `flai forge push`, and `flai serve` under the `forge` host action, which reads its own environment. No agent `flai serve` starts gets it. Unlike the row for `GITHUB_TOKEN` above, `gh auth token` is not borrowed when it is unset. Set it in the environment `flai host` or `flai serve` runs in ([A forge token for pull-request acceptance](index.md#a-forge-token-for-pull-request-acceptance-s-0392)) |

### Set by flai for the processes it starts

Do not set these yourself; a command you write for `flai serve agent set` may read them.

| Variable | Set by | Holds |
|----------|--------|-------|
| `FLAI_HOST_URL`, `FLAI_HOST_TOKEN` | `flai host`, for `flai serve` and the MCP servers; never passed on to an agent `flai serve` starts | The host's address and the token its API takes; `flai serve` uses them only when the host's token is the one beside its own config |
| `FLAI_STORY` | `flai serve`, for an agent it starts for a story | The story the agent is to work; with no `FLAI_ROLE`, `flai guard` reads it to refuse the agent's own `wait_for_events` while a sub-agent of its session runs and no thread on the story is open, and the `SubagentStart` and `SubagentStop` hooks run only when it is set, as does the `PreToolUse` hook on file writes outside a strategic agent's session; when a command is not given the story, it is the default before `FLAI_AGENT` and the branch checked out: the story `flai message send` and `reply` write from ([Messages between stories](../users/flai.md#messages-between-stories)), the story `flai issue new` and `bump` record, and the story whose branch `--commit` commits on for them, `flai issue close`, and `flai adr new` |
| `FLAI_ROLE` | `flai serve`, for the planner it starts (`flai plan`), the orchestrator it runs (the `orchestrate` host action), the analyzer it starts (the `analyze` host action), and the curator it starts (the `curate` host action) | `plan`, `orchestrate`, `analyze`, or `curate`; `flai guard` reads it to hold the session to planning, to the orchestrator's permissions, to the analyzer's report, or to the curator's digest commands |
| `FLAI_ITEM` | `flai serve`, for the planner it starts | The epic or story the planner plans |
| `FLAI_FOCUS` | `flai serve`, for an analyzer it starts with a focus under the `command` harness | `bottlenecks`, `intent`, or `risk`; unset for a run that looks for all three. The `claude-code` harness says it in the prompt instead |
| `FLAI_SINCE` | `flai serve`, for a curator it starts with a since under the `command` harness | A `YYYY-MM-DD` day in UTC or an RFC 3339 time: the since it was given, else the time of its last run; unset on its first run, which takes in everything. The `claude-code` harness says it in the prompt instead |
| `FLAI_STARTED_BY` | `flai serve` | `flai-serve` |
| `FLAI_AGENT`, `FLAI_SESSION` | `flai serve` | The agent's own name, `<agent.name>-<story>`, `planner-<item>` for the planner, or `orchestrator` for the orchestrator, and a session made from the start time |
| `FLAI_MODEL`, `FLAI_HARNESS` | `flai serve`, for your `agent.command` | The story's model and harness |
| `FLAI_AGENT_CONFIG` | `flai serve`, for your `agent.command` | The story's `agent.config` as a JSON object |
| `FLAI_AGENT_ROLES` | `flai serve`, for your `agent.command`, when the story's agent has roles | The story's `agent.roles` as a JSON object, each role with its `harness`, `model`, and `config` |
| `FLAI_ANSWERED` | `flai serve`, for an agent started again because its question was answered | The thread's ID |
| `FLAI_COMMIT` | `flai serve`, for your `agent.command` started to commit what a story's worktree holds (`flai serve agent commit`) | The worktree's path |
| `FLAI_PROVIDER`, `FLAI_PROVIDER_API`, `FLAI_PROVIDER_BASE_URL`, `FLAI_PROVIDER_KEY_ENV` | `flai serve`, for your `agent.command`, when the agent names a provider | The provider's name, its `api`, its `base_url`, and the name of the variable that holds its key, which the command reads itself ([Where its model's calls go](../users/flai.md#where-its-models-calls-go-a-provider)) |
| `ANTHROPIC_BASE_URL`, `ANTHROPIC_AUTH_TOKEN`, `ANTHROPIC_API_KEY`, `ANTHROPIC_DEFAULT_OPUS_MODEL`, `ANTHROPIC_DEFAULT_SONNET_MODEL`, `ANTHROPIC_DEFAULT_HAIKU_MODEL` | `flai serve`, for a `claude-code` agent that names a provider | The provider's `base_url`; its key, copied from the variable its `key_env` names; empty; and its `models`' `opus`, `sonnet`, and `haiku`, each only when named |

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
| `FLAIOVER_ALLOW_UNSIGNED` | unset | `1`, and only `1`, serves a flai that sends no release stamp, and shows the banner; `flai dashboard` sets it when `dashboard.allow_unsigned` is in force or with `--build` |
| `FLAIOVER_ALLOWED_HOSTS` | unset | Host names, each with an optional port, separated by commas, lower-case, with no spaces and no scheme, that the dashboard accepts in a browser's `Origin` beside its own address; `flai dashboard` sets it from `dashboard.allowed_hosts` when that names any |
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
| `FLAIOVER_VERSION`, `FLAIOVER_COMMIT` | baked into the image | The release and commit `flaiover_build_info` reports, beside its `signed` label, and the release stamp must name; not for changing |
| `FLAIOVER_RELEASE_STAMP` | unset: `/app/release-stamp` | The file holding the image's release stamp, which the server checks at start ([The release stamp](index.md#the-release-stamp)); for a build that keeps it elsewhere |
| `PROJECT_DIR` | the working directory | The project the development server and tests ask flai about; names nothing in the container |
| `FLAIOVER_LICENSE` | unset: `LICENSE.md` beside the server (`/app/LICENSE.md` in the image), else one directory up | The license file the Host › License page and `GET /api/license` show (S-0231); for a build that keeps it elsewhere |
| `--publish` | `<dashboard.bind>:<dashboard.port>:3000` | Where the host reaches the container |
| `--user` | your user and group | So that the two secret files, readable only by you, can be read |

Build arguments, for an image built with `flai dashboard --build` or `make flaiover-image`: `FLAIOVER_VERSION` and `FLAI_VERSION` from the newest `flaiover/v*` and `flai/v*` tags, `FLAI_COMMIT` the commit, and `FLAI_DATE` the build time. `flai dashboard --build` sets them. `FLAIOVER_STAMP_STATEMENT` and `FLAIOVER_STAMP_SIGNATURE` are the release stamp, which only the release workflow gives, both or neither; a local build has none.

## Release signing secrets

Stored secrets of the repository's CI, not of flai or the dashboard: no flai or flaiover reads them. They are in the GitHub Actions environment `release` of `bytepunx/system-flow`, which only `flai/v*` and `flaiover/v*` tags deploy to. How to make, store, rotate, and replace them is in the [release key runbook](runbooks/release-key.md).

| Secret | Read by | For |
|--------|---------|-----|
| `COSIGN_PRIVATE_KEY` | `release-flai.yml` and `release-flaiover.yml`, on a release tag | The release key's private half, the encrypted PEM `cosign generate-key-pair` wrote, which signs flai's `checksums.txt`, flaiover's digest list, and the flaiover image ([The release signing key](runbooks/release-key.md#store-the-secrets)) |
| `COSIGN_PASSWORD` | the same | The password that decrypts it |

## Flags

Every flag of every flai command, with the commands that take it. What a flag means and its default are in the command's section of [the command reference](../users/flai-reference.md), which, like this table, is generated from `flai --help`. The flags that change a setting for one run are named in the tables above.

<!-- flags:begin: generated by flai reference from the command help; run make flai-reference -->
| Flag | Taken by |
|------|----------|
| `--about` | [flai message send](../users/flai-reference.md#flai-message-send) |
| `--add` | [flai touches](../users/flai-reference.md#flai-touches) |
| `--addr` | [flai mcp http](../users/flai-reference.md#flai-mcp-http), [flai mcp start](../users/flai-reference.md#flai-mcp-start) |
| `--after` | [flai edit](../users/flai-reference.md#flai-edit), [flai order](../users/flai-reference.md#flai-order), [flai story new](../users/flai-reference.md#flai-story-new), [flai task new](../users/flai-reference.md#flai-task-new) |
| `--agent` | [flai mcp](../users/flai-reference.md#flai-mcp) |
| `--agent-config` | [flai edit](../users/flai-reference.md#flai-edit), [flai story new](../users/flai-reference.md#flai-story-new) |
| `--all` | [flai board](../users/flai-reference.md#flai-board), [flai issue list](../users/flai-reference.md#flai-issue-list), [flai message list](../users/flai-reference.md#flai-message-list), [flai serve agent usage](../users/flai-reference.md#flai-serve-agent-usage), [flai test](../users/flai-reference.md#flai-test), [flai thread list](../users/flai-reference.md#flai-thread-list), [flai verify](../users/flai-reference.md#flai-verify) |
| `--all-projects` | [flai serve disable](../users/flai-reference.md#flai-serve-disable), [flai serve enable](../users/flai-reference.md#flai-serve-enable) |
| `--allow-unsigned` | [flai dashboard](../users/flai-reference.md#flai-dashboard) |
| `--api` | [flai serve agent provider](../users/flai-reference.md#flai-serve-agent-provider) |
| `--apply` | [flai order](../users/flai-reference.md#flai-order), [flai release](../users/flai-reference.md#flai-release) |
| `--asked-again` | [flai digest edit](../users/flai-reference.md#flai-digest-edit) |
| `--attach` | [flai dashboard](../users/flai-reference.md#flai-dashboard) |
| `--auto-restarts` | [flai serve agent set](../users/flai-reference.md#flai-serve-agent-set) |
| `--autocommit` | [flai adr accept](../users/flai-reference.md#flai-adr-accept), [flai adr new](../users/flai-reference.md#flai-adr-new), [flai adr topics](../users/flai-reference.md#flai-adr-topics), [flai agent clear](../users/flai-reference.md#flai-agent-clear), [flai agent set](../users/flai-reference.md#flai-agent-set), [flai criteria tick](../users/flai-reference.md#flai-criteria-tick), [flai criteria untick](../users/flai-reference.md#flai-criteria-untick), [flai edit](../users/flai-reference.md#flai-edit), [flai epic new](../users/flai-reference.md#flai-epic-new), [flai issue bump](../users/flai-reference.md#flai-issue-bump), [flai issue close](../users/flai-reference.md#flai-issue-close), [flai issue new](../users/flai-reference.md#flai-issue-new), [flai issue story](../users/flai-reference.md#flai-issue-story), [flai manifest set](../users/flai-reference.md#flai-manifest-set), [flai shared add](../users/flai-reference.md#flai-shared-add), [flai shared remove](../users/flai-reference.md#flai-shared-remove), [flai story new](../users/flai-reference.md#flai-story-new), [flai subproject add](../users/flai-reference.md#flai-subproject-add), [flai subproject remove](../users/flai-reference.md#flai-subproject-remove), [flai subproject set](../users/flai-reference.md#flai-subproject-set), [flai task new](../users/flai-reference.md#flai-task-new) |
| `--base-url` | [flai serve agent provider](../users/flai-reference.md#flai-serve-agent-provider) |
| `--batch` | [flai accept](../users/flai-reference.md#flai-accept) |
| `--before` | [flai order](../users/flai-reference.md#flai-order) |
| `--bind` | [flai dashboard](../users/flai-reference.md#flai-dashboard), [flai dashboard restart](../users/flai-reference.md#flai-dashboard-restart), [flai dashboard upgrade](../users/flai-reference.md#flai-dashboard-upgrade) |
| `--body` | [flai digest edit](../users/flai-reference.md#flai-digest-edit), [flai digest new](../users/flai-reference.md#flai-digest-new) |
| `--body-file` | [flai digest edit](../users/flai-reference.md#flai-digest-edit), [flai digest new](../users/flai-reference.md#flai-digest-new) |
| `--body-stdin` | [flai adr new](../users/flai-reference.md#flai-adr-new), [flai edit](../users/flai-reference.md#flai-edit), [flai epic new](../users/flai-reference.md#flai-epic-new), [flai story new](../users/flai-reference.md#flai-story-new), [flai task new](../users/flai-reference.md#flai-task-new) |
| `--bottom` | [flai order](../users/flai-reference.md#flai-order) |
| `--bucket` | [flai stats](../users/flai-reference.md#flai-stats) |
| `--budget` | [flai prime](../users/flai-reference.md#flai-prime), [flai story start](../users/flai-reference.md#flai-story-start) |
| `--build` | [flai dashboard](../users/flai-reference.md#flai-dashboard) |
| `--by` | [flai accept](../users/flai-reference.md#flai-accept), [flai criteria tick](../users/flai-reference.md#flai-criteria-tick), [flai criteria untick](../users/flai-reference.md#flai-criteria-untick), [flai edit](../users/flai-reference.md#flai-edit), [flai message escalate](../users/flai-reference.md#flai-message-escalate), [flai message reply](../users/flai-reference.md#flai-message-reply), [flai message send](../users/flai-reference.md#flai-message-send), [flai message share](../users/flai-reference.md#flai-message-share), [flai move](../users/flai-reference.md#flai-move), [flai order](../users/flai-reference.md#flai-order), [flai queue add](../users/flai-reference.md#flai-queue-add), [flai stats](../users/flai-reference.md#flai-stats), [flai stream answer](../users/flai-reference.md#flai-stream-answer), [flai thread confirm](../users/flai-reference.md#flai-thread-confirm), [flai thread new](../users/flai-reference.md#flai-thread-new), [flai thread reply](../users/flai-reference.md#flai-thread-reply), [flai thread resolve](../users/flai-reference.md#flai-thread-resolve) |
| `--candidates` | [flai plan](../users/flai-reference.md#flai-plan), [flai promote](../users/flai-reference.md#flai-promote) |
| `--cat` | [flai prime](../users/flai-reference.md#flai-prime) |
| `--check` | [flai self-upgrade](../users/flai-reference.md#flai-self-upgrade) |
| `--class` | [flai issue new](../users/flai-reference.md#flai-issue-new) |
| `--clear` | [flai touches](../users/flai-reference.md#flai-touches) |
| `--clear-after` | [flai edit](../users/flai-reference.md#flai-edit) |
| `--clear-agent` | [flai edit](../users/flai-reference.md#flai-edit) |
| `--clear-cost-of-delay` | [flai edit](../users/flai-reference.md#flai-edit) |
| `--clear-forecast` | [flai edit](../users/flai-reference.md#flai-edit) |
| `--clear-tags` | [flai edit](../users/flai-reference.md#flai-edit) |
| `--clear-topics` | [flai edit](../users/flai-reference.md#flai-edit) |
| `--clear-touches` | [flai edit](../users/flai-reference.md#flai-edit) |
| `--commit` | [flai adr new](../users/flai-reference.md#flai-adr-new), [flai import](../users/flai-reference.md#flai-import), [flai issue bump](../users/flai-reference.md#flai-issue-bump), [flai issue close](../users/flai-reference.md#flai-issue-close), [flai issue new](../users/flai-reference.md#flai-issue-new) |
| `--config` | every command ([global flags](../users/flai-reference.md#flai)), [flai agent set](../users/flai-reference.md#flai-agent-set) |
| `--cost` | [flai issue bump](../users/flai-reference.md#flai-issue-bump), [flai issue new](../users/flai-reference.md#flai-issue-new) |
| `--cost-of-delay-value` | [flai edit](../users/flai-reference.md#flai-edit) |
| `--current` | [flai stream state](../users/flai-reference.md#flai-stream-state) |
| `--defaults` | [flai new](../users/flai-reference.md#flai-new) |
| `--deliver` | [flai release](../users/flai-reference.md#flai-release) |
| `--deny-protected` | [flai serve agent harness](../users/flai-reference.md#flai-serve-agent-harness) |
| `--dir` | [flai self-upgrade](../users/flai-reference.md#flai-self-upgrade) |
| `--draft` | [flai edit](../users/flai-reference.md#flai-edit), [flai story new](../users/flai-reference.md#flai-story-new) |
| `--drafts` | [flai promote](../users/flai-reference.md#flai-promote) |
| `--dry-run` | [flai accept](../users/flai-reference.md#flai-accept), [flai archive](../users/flai-reference.md#flai-archive), [flai digest seed](../users/flai-reference.md#flai-digest-seed), [flai import](../users/flai-reference.md#flai-import), [flai migrate ids](../users/flai-reference.md#flai-migrate-ids), [flai move](../users/flai-reference.md#flai-move), [flai push](../users/flai-reference.md#flai-push), [flai release](../users/flai-reference.md#flai-release), [flai template push](../users/flai-reference.md#flai-template-push), [flai upgrade](../users/flai-reference.md#flai-upgrade) |
| `--epic` | [flai issue story](../users/flai-reference.md#flai-issue-story), [flai prime](../users/flai-reference.md#flai-prime), [flai story new](../users/flai-reference.md#flai-story-new) |
| `--evaluate` | [flai release](../users/flai-reference.md#flai-release) |
| `--evidence` | [flai accept](../users/flai-reference.md#flai-accept), [flai issue bump](../users/flai-reference.md#flai-issue-bump), [flai issue new](../users/flai-reference.md#flai-issue-new), [flai move](../users/flai-reference.md#flai-move), [flai queue add](../users/flai-reference.md#flai-queue-add) |
| `--exclude` | [flai components config](../users/flai-reference.md#flai-components-config) |
| `--focus` | [flai analyze](../users/flai-reference.md#flai-analyze) |
| `-f`, `--follow` | [flai dashboard logs](../users/flai-reference.md#flai-dashboard-logs), [flai serve agent stream](../users/flai-reference.md#flai-serve-agent-stream) |
| `--force` | [flai import](../users/flai-reference.md#flai-import), [flai new](../users/flai-reference.md#flai-new), [flai template push](../users/flai-reference.md#flai-template-push), [flai upgrade](../users/flai-reference.md#flai-upgrade) |
| `--forecast-basis` | [flai edit](../users/flai-reference.md#flai-edit) |
| `--forecast-delivery` | [flai edit](../users/flai-reference.md#flai-edit) |
| `--forecast-duration` | [flai edit](../users/flai-reference.md#flai-edit) |
| `--fresh` | [flai verify](../users/flai-reference.md#flai-verify) |
| `--from` | [flai checks tail](../users/flai-reference.md#flai-checks-tail), [flai message escalate](../users/flai-reference.md#flai-message-escalate), [flai message reply](../users/flai-reference.md#flai-message-reply), [flai message send](../users/flai-reference.md#flai-message-send), [flai message share](../users/flai-reference.md#flai-message-share), [flai serve agent stream](../users/flai-reference.md#flai-serve-agent-stream) |
| `--grace-seconds` | [flai checks cancel](../users/flai-reference.md#flai-checks-cancel) |
| `--guard` | [flai serve agent harness](../users/flai-reference.md#flai-serve-agent-harness) |
| `--harness` | [flai agent set](../users/flai-reference.md#flai-agent-set), [flai edit](../users/flai-reference.md#flai-edit), [flai guard](../users/flai-reference.md#flai-guard), [flai story new](../users/flai-reference.md#flai-story-new) |
| `--hash` | [flai criteria tick](../users/flai-reference.md#flai-criteria-tick), [flai criteria untick](../users/flai-reference.md#flai-criteria-untick), [flai doc save](../users/flai-reference.md#flai-doc-save), [flai edit](../users/flai-reference.md#flai-edit), [flai file save](../users/flai-reference.md#flai-file-save) |
| `--heading` | [flai doc show](../users/flai-reference.md#flai-doc-show), [flai thread new](../users/flai-reference.md#flai-thread-new) |
| `--idle` | [flai mcp http](../users/flai-reference.md#flai-mcp-http), [flai mcp start](../users/flai-reference.md#flai-mcp-start) |
| `--image` | [flai dashboard](../users/flai-reference.md#flai-dashboard), [flai dashboard check](../users/flai-reference.md#flai-dashboard-check), [flai dashboard restart](../users/flai-reference.md#flai-dashboard-restart), [flai dashboard upgrade](../users/flai-reference.md#flai-dashboard-upgrade) |
| `--include` | [flai components config](../users/flai-reference.md#flai-components-config) |
| `--json` | every command ([global flags](../users/flai-reference.md#flai)) |
| `--keep-all` | [flai upgrade](../users/flai-reference.md#flai-upgrade) |
| `--keep-placed` | [flai order](../users/flai-reference.md#flai-order) |
| `--keep-removal` | [flai accept](../users/flai-reference.md#flai-accept), [flai move](../users/flai-reference.md#flai-move) |
| `--key-env` | [flai serve agent provider](../users/flai-reference.md#flai-serve-agent-provider) |
| `--kind` | [flai digest harvest](../users/flai-reference.md#flai-digest-harvest), [flai digest new](../users/flai-reference.md#flai-digest-new), [flai subproject add](../users/flai-reference.md#flai-subproject-add), [flai subproject set](../users/flai-reference.md#flai-subproject-set) |
| `--last` | [flai verify](../users/flai-reference.md#flai-verify) |
| `-n`, `--last` | [flai serve journal](../users/flai-reference.md#flai-serve-journal) |
| `--layout` | [flai import](../users/flai-reference.md#flai-import), [flai new](../users/flai-reference.md#flai-new) |
| `--limit` | [flai doc search](../users/flai-reference.md#flai-doc-search), [flai promote](../users/flai-reference.md#flai-promote), [flai touches suggest](../users/flai-reference.md#flai-touches-suggest) |
| `--list` | [flai self-upgrade](../users/flai-reference.md#flai-self-upgrade) |
| `--log` | [flai task done](../users/flai-reference.md#flai-task-done) |
| `--max` | [flai test](../users/flai-reference.md#flai-test), [flai verify](../users/flai-reference.md#flai-verify) |
| `--max-sessions` | [flai mcp http](../users/flai-reference.md#flai-mcp-http), [flai mcp start](../users/flai-reference.md#flai-mcp-start) |
| `--message` | [flai doc save](../users/flai-reference.md#flai-doc-save), [flai edit](../users/flai-reference.md#flai-edit), [flai file save](../users/flai-reference.md#flai-file-save) |
| `-m`, `--message` | [flai task done](../users/flai-reference.md#flai-task-done) |
| `--min` | [flai touches suggest](../users/flai-reference.md#flai-touches-suggest) |
| `--model` | [flai agent set](../users/flai-reference.md#flai-agent-set), [flai edit](../users/flai-reference.md#flai-edit), [flai serve agent provider](../users/flai-reference.md#flai-serve-agent-provider), [flai story new](../users/flai-reference.md#flai-story-new) |
| `--name` | [flai serve agent set](../users/flai-reference.md#flai-serve-agent-set), [flai serve checks set](../users/flai-reference.md#flai-serve-checks-set), [flai subproject set](../users/flai-reference.md#flai-subproject-set) |
| `--nature` | [flai edit](../users/flai-reference.md#flai-edit), [flai epic new](../users/flai-reference.md#flai-epic-new), [flai story new](../users/flai-reference.md#flai-story-new), [flai task new](../users/flai-reference.md#flai-task-new) |
| `--next` | [flai stream state](../users/flai-reference.md#flai-stream-state) |
| `--no-branch` | [flai stream open](../users/flai-reference.md#flai-stream-open) |
| `--no-commit` | [flai doc save](../users/flai-reference.md#flai-doc-save), [flai file save](../users/flai-reference.md#flai-file-save) |
| `--no-descriptions` | [flai completion bash](../users/flai-reference.md#flai-completion-bash), [flai completion fish](../users/flai-reference.md#flai-completion-fish), [flai completion powershell](../users/flai-reference.md#flai-completion-powershell), [flai completion zsh](../users/flai-reference.md#flai-completion-zsh) |
| `--no-draft` | [flai edit](../users/flai-reference.md#flai-edit) |
| `--no-git` | [flai new](../users/flai-reference.md#flai-new) |
| `--no-restart` | [flai dashboard token](../users/flai-reference.md#flai-dashboard-token) |
| `--no-serve` | [flai dashboard](../users/flai-reference.md#flai-dashboard) |
| `--note` | [flai issue bump](../users/flai-reference.md#flai-issue-bump), [flai issue new](../users/flai-reference.md#flai-issue-new) |
| `--on` | [flai thread list](../users/flai-reference.md#flai-thread-list), [flai thread new](../users/flai-reference.md#flai-thread-new) |
| `--open` | [flai dashboard](../users/flai-reference.md#flai-dashboard) |
| `--owner` | [flai epic new](../users/flai-reference.md#flai-epic-new), [flai issue story](../users/flai-reference.md#flai-issue-story), [flai story new](../users/flai-reference.md#flai-story-new), [flai task new](../users/flai-reference.md#flai-task-new) |
| `--parent` | [flai edit](../users/flai-reference.md#flai-edit) |
| `--part` | [flai prime](../users/flai-reference.md#flai-prime) |
| `--path` | [flai digest edit](../users/flai-reference.md#flai-digest-edit), [flai digest new](../users/flai-reference.md#flai-digest-new), [flai stream restore](../users/flai-reference.md#flai-stream-restore), [flai subproject add](../users/flai-reference.md#flai-subproject-add), [flai subproject set](../users/flai-reference.md#flai-subproject-set) |
| `--paths` | [flai message share](../users/flai-reference.md#flai-message-share) |
| `--penalty-per-week` | [flai edit](../users/flai-reference.md#flai-edit), [flai epic new](../users/flai-reference.md#flai-epic-new), [flai issue bump](../users/flai-reference.md#flai-issue-bump), [flai issue new](../users/flai-reference.md#flai-issue-new), [flai story new](../users/flai-reference.md#flai-story-new) |
| `--pending` | [flai push](../users/flai-reference.md#flai-push), [flai release](../users/flai-reference.md#flai-release) |
| `--placed-by` | [flai order](../users/flai-reference.md#flai-order) |
| `--port` | [flai dashboard](../users/flai-reference.md#flai-dashboard), [flai dashboard restart](../users/flai-reference.md#flai-dashboard-restart), [flai dashboard upgrade](../users/flai-reference.md#flai-dashboard-upgrade) |
| `--print-body` | [flai adr new](../users/flai-reference.md#flai-adr-new), [flai epic new](../users/flai-reference.md#flai-epic-new), [flai story new](../users/flai-reference.md#flai-story-new), [flai task new](../users/flai-reference.md#flai-task-new) |
| `--program` | [flai serve agent harness](../users/flai-reference.md#flai-serve-agent-harness) |
| `--project` | [flai components list](../users/flai-reference.md#flai-components-list), [flai components serve](../users/flai-reference.md#flai-components-serve), [flai design init](../users/flai-reference.md#flai-design-init) |
| `--provider` | [flai agent set](../users/flai-reference.md#flai-agent-set), [flai edit](../users/flai-reference.md#flai-edit), [flai story new](../users/flai-reference.md#flai-story-new) |
| `--publish` | [flai push](../users/flai-reference.md#flai-push) |
| `--published` | [flai dashboard upgrade](../users/flai-reference.md#flai-dashboard-upgrade) |
| `--pull` | [flai dashboard](../users/flai-reference.md#flai-dashboard) |
| `--reason` | [flai accept](../users/flai-reference.md#flai-accept), [flai block](../users/flai-reference.md#flai-block), [flai issue close](../users/flai-reference.md#flai-issue-close), [flai move](../users/flai-reference.md#flai-move), [flai thread resolve](../users/flai-reference.md#flai-thread-resolve) |
| `--recommend` | [flai thread reply](../users/flai-reference.md#flai-thread-reply) |
| `--record-issues` | [flai check](../users/flai-reference.md#flai-check), [flai verify](../users/flai-reference.md#flai-verify) |
| `--ref` | [flai import](../users/flai-reference.md#flai-import), [flai new](../users/flai-reference.md#flai-new), [flai template push](../users/flai-reference.md#flai-template-push), [flai template show](../users/flai-reference.md#flai-template-show), [flai template use](../users/flai-reference.md#flai-template-use), [flai upgrade](../users/flai-reference.md#flai-upgrade) |
| `--refines` | [flai adr new](../users/flai-reference.md#flai-adr-new) |
| `--relock` | [flai upgrade](../users/flai-reference.md#flai-upgrade) |
| `--remote` | [flai template push](../users/flai-reference.md#flai-template-push) |
| `--remove` | [flai touches](../users/flai-reference.md#flai-touches) |
| `--replace` | [flai agent set](../users/flai-reference.md#flai-agent-set) |
| `--replace-all` | [flai upgrade](../users/flai-reference.md#flai-upgrade) |
| `--repo` | [flai dashboard upgrade](../users/flai-reference.md#flai-dashboard-upgrade), [flai dashboard versions](../users/flai-reference.md#flai-dashboard-versions), [flai self-upgrade](../users/flai-reference.md#flai-self-upgrade) |
| `--report` | [flai issue bump](../users/flai-reference.md#flai-issue-bump), [flai issue new](../users/flai-reference.md#flai-issue-new) |
| `--reset` | [flai serve agent harness](../users/flai-reference.md#flai-serve-agent-harness), [flai serve agent provider](../users/flai-reference.md#flai-serve-agent-provider), [flai serve forge](../users/flai-reference.md#flai-serve-forge) |
| `--revenue-per-week` | [flai edit](../users/flai-reference.md#flai-edit), [flai epic new](../users/flai-reference.md#flai-epic-new), [flai issue bump](../users/flai-reference.md#flai-issue-bump), [flai issue new](../users/flai-reference.md#flai-issue-new), [flai story new](../users/flai-reference.md#flai-story-new) |
| `--role` | [flai prime](../users/flai-reference.md#flai-prime) |
| `--role-config` | [flai agent set](../users/flai-reference.md#flai-agent-set), [flai edit](../users/flai-reference.md#flai-edit), [flai story new](../users/flai-reference.md#flai-story-new) |
| `--role-harness` | [flai agent set](../users/flai-reference.md#flai-agent-set), [flai edit](../users/flai-reference.md#flai-edit), [flai story new](../users/flai-reference.md#flai-story-new) |
| `--role-model` | [flai agent set](../users/flai-reference.md#flai-agent-set), [flai edit](../users/flai-reference.md#flai-edit), [flai story new](../users/flai-reference.md#flai-story-new) |
| `--rotate` | [flai dashboard token](../users/flai-reference.md#flai-dashboard-token), [flai mcp token](../users/flai-reference.md#flai-mcp-token) |
| `--settled` | [flai message reply](../users/flai-reference.md#flai-message-reply) |
| `--setup` | [flai components config](../users/flai-reference.md#flai-components-config) |
| `--show` | [flai edit](../users/flai-reference.md#flai-edit) |
| `--since` | [flai curate](../users/flai-reference.md#flai-curate), [flai digest harvest](../users/flai-reference.md#flai-digest-harvest), [flai digest seed](../users/flai-reference.md#flai-digest-seed), [flai stats](../users/flai-reference.md#flai-stats) |
| `--source` | [flai digest edit](../users/flai-reference.md#flai-digest-edit), [flai digest new](../users/flai-reference.md#flai-digest-new), [flai thread reply](../users/flai-reference.md#flai-thread-reply) |
| `--spend-key-env` | [flai serve agent provider](../users/flai-reference.md#flai-serve-agent-provider) |
| `--status` | [flai adr new](../users/flai-reference.md#flai-adr-new), [flai digest list](../users/flai-reference.md#flai-digest-list) |
| `--stdin` | [flai doc lint](../users/flai-reference.md#flai-doc-lint) |
| `--story` | [flai check](../users/flai-reference.md#flai-check), [flai components serve](../users/flai-reference.md#flai-components-serve), [flai design show](../users/flai-reference.md#flai-design-show), [flai doc save](../users/flai-reference.md#flai-doc-save), [flai doc show](../users/flai-reference.md#flai-doc-show), [flai file save](../users/flai-reference.md#flai-file-save), [flai file show](../users/flai-reference.md#flai-file-show), [flai issue bump](../users/flai-reference.md#flai-issue-bump), [flai issue list](../users/flai-reference.md#flai-issue-list), [flai issue new](../users/flai-reference.md#flai-issue-new), [flai issue story](../users/flai-reference.md#flai-issue-story), [flai message list](../users/flai-reference.md#flai-message-list), [flai prime](../users/flai-reference.md#flai-prime), [flai task new](../users/flai-reference.md#flai-task-new) |
| `--strict` | [flai check](../users/flai-reference.md#flai-check) |
| `--supersedes` | [flai adr new](../users/flai-reference.md#flai-adr-new) |
| `--sync-only` | [flai verify](../users/flai-reference.md#flai-verify) |
| `--tag` | [flai dashboard](../users/flai-reference.md#flai-dashboard), [flai dashboard check](../users/flai-reference.md#flai-dashboard-check), [flai dashboard restart](../users/flai-reference.md#flai-dashboard-restart), [flai dashboard upgrade](../users/flai-reference.md#flai-dashboard-upgrade), [flai edit](../users/flai-reference.md#flai-edit), [flai epic new](../users/flai-reference.md#flai-epic-new), [flai story new](../users/flai-reference.md#flai-story-new), [flai task new](../users/flai-reference.md#flai-task-new), [flai template push](../users/flai-reference.md#flai-template-push) |
| `--tags` | [flai subproject add](../users/flai-reference.md#flai-subproject-add), [flai subproject set](../users/flai-reference.md#flai-subproject-set) |
| `--tailwind` | [flai design css](../users/flai-reference.md#flai-design-css) |
| `--template` | [flai import](../users/flai-reference.md#flai-import), [flai new](../users/flai-reference.md#flai-new), [flai template show](../users/flai-reference.md#flai-template-show), [flai upgrade](../users/flai-reference.md#flai-upgrade) |
| `--time-lost-per-cycle` | [flai edit](../users/flai-reference.md#flai-edit), [flai epic new](../users/flai-reference.md#flai-epic-new), [flai issue bump](../users/flai-reference.md#flai-issue-bump), [flai issue new](../users/flai-reference.md#flai-issue-new), [flai story new](../users/flai-reference.md#flai-story-new) |
| `--timing` | [flai hostapi](../users/flai-reference.md#flai-hostapi) |
| `--title` | [flai edit](../users/flai-reference.md#flai-edit) |
| `--token-env` | [flai serve forge](../users/flai-reference.md#flai-serve-forge) |
| `--top` | [flai order](../users/flai-reference.md#flai-order) |
| `--topic` | [flai digest edit](../users/flai-reference.md#flai-digest-edit), [flai digest list](../users/flai-reference.md#flai-digest-list), [flai digest new](../users/flai-reference.md#flai-digest-new) |
| `--topics` | [flai design init](../users/flai-reference.md#flai-design-init), [flai edit](../users/flai-reference.md#flai-edit), [flai epic new](../users/flai-reference.md#flai-epic-new), [flai story new](../users/flai-reference.md#flai-story-new) |
| `--touches` | [flai edit](../users/flai-reference.md#flai-edit), [flai epic new](../users/flai-reference.md#flai-epic-new), [flai story new](../users/flai-reference.md#flai-story-new), [flai task new](../users/flai-reference.md#flai-task-new) |
| `--trailer` | [flai accept](../users/flai-reference.md#flai-accept), [flai adr accept](../users/flai-reference.md#flai-adr-accept), [flai adr new](../users/flai-reference.md#flai-adr-new), [flai adr topics](../users/flai-reference.md#flai-adr-topics), [flai agent clear](../users/flai-reference.md#flai-agent-clear), [flai agent set](../users/flai-reference.md#flai-agent-set), [flai criteria tick](../users/flai-reference.md#flai-criteria-tick), [flai criteria untick](../users/flai-reference.md#flai-criteria-untick), [flai doc save](../users/flai-reference.md#flai-doc-save), [flai edit](../users/flai-reference.md#flai-edit), [flai epic new](../users/flai-reference.md#flai-epic-new), [flai file save](../users/flai-reference.md#flai-file-save), [flai import](../users/flai-reference.md#flai-import), [flai issue bump](../users/flai-reference.md#flai-issue-bump), [flai issue close](../users/flai-reference.md#flai-issue-close), [flai issue new](../users/flai-reference.md#flai-issue-new), [flai issue story](../users/flai-reference.md#flai-issue-story), [flai manifest set](../users/flai-reference.md#flai-manifest-set), [flai move](../users/flai-reference.md#flai-move), [flai shared add](../users/flai-reference.md#flai-shared-add), [flai shared remove](../users/flai-reference.md#flai-shared-remove), [flai story new](../users/flai-reference.md#flai-story-new), [flai stream restore](../users/flai-reference.md#flai-stream-restore), [flai subproject add](../users/flai-reference.md#flai-subproject-add), [flai subproject remove](../users/flai-reference.md#flai-subproject-remove), [flai subproject set](../users/flai-reference.md#flai-subproject-set), [flai task new](../users/flai-reference.md#flai-task-new) |
| `--type` | [flai stats](../users/flai-reference.md#flai-stats) |
| `--unset` | [flai agent set](../users/flai-reference.md#flai-agent-set), [flai components config](../users/flai-reference.md#flai-components-config), [flai manifest set](../users/flai-reference.md#flai-manifest-set) |
| `--unset-role` | [flai agent set](../users/flai-reference.md#flai-agent-set), [flai edit](../users/flai-reference.md#flai-edit) |
| `--var` | [flai import](../users/flai-reference.md#flai-import), [flai new](../users/flai-reference.md#flai-new), [flai upgrade](../users/flai-reference.md#flai-upgrade) |
| `-v`, `--verbose` | every command ([global flags](../users/flai-reference.md#flai)) |
| `--verified` | [flai accept](../users/flai-reference.md#flai-accept), [flai move](../users/flai-reference.md#flai-move), [flai queue add](../users/flai-reference.md#flai-queue-add) |
| `--version` | [flai host upgrade](../users/flai-reference.md#flai-host-upgrade), [flai self-upgrade](../users/flai-reference.md#flai-self-upgrade) |
| `--wait` | [flai checks tail](../users/flai-reference.md#flai-checks-tail) |
| `--write` | [flai design css](../users/flai-reference.md#flai-design-css), [flai serve agent usage](../users/flai-reference.md#flai-serve-agent-usage) |
| `-y`, `--yes` | every command ([global flags](../users/flai-reference.md#flai)) |
<!-- flags:end -->
