---
id: S-0442
type: story
nature: feature
title: The Components page renders a story's worktree, so the operator sees a component an agent is changing live while the story is in progress or in review
status: in-progress
parent: E-0026
owner: alex
created: 2026-10-16T21:37:05Z
updated: 2026-10-17T02:35:31Z
transitions:
  - to: ready
    at: 2026-10-17T01:02:20Z
    by: alex
  - to: in-progress
    at: 2026-10-17T02:35:31Z
    by: agent-S-0442
tags: [dashboard, cli]
topics: [components, dashboard, cli]
touches: [flai/internal/components/worktree.go, flai/internal/components/worktree_test.go, flai/internal/components/devserver.go, flai/internal/components/devserver_test.go, flai/internal/serve/components.go, flai/internal/serve/components_test.go, flai/internal/hostapi/components.go, flai/internal/hostapi/components_test.go, flai/internal/hostapi/writes.go, flai/internal/hostapi/writes_test.go, flai/cmd/components.go, flai/cmd/components_test.go, flai/cmd/serve_actions.go, docs/users/flai-reference.md, flaiover/src/lib/components/explorer/SourcePicker.svelte, flaiover/src/lib/components/explorer/SourcePicker.svelte.test.ts, flaiover/src/routes/design/components/+page.svelte, flaiover/src/routes/design/components/components.svelte.test.ts, flaiover/src/routes/api/components/+server.ts, flaiover/src/routes/api/components/start/+server.ts, flaiover/src/routes/api/components/stop/+server.ts, flaiover/src/routes/api/components/components.test.ts, flaiover/src/lib/server/repo.ts, "flaiover/src/routes/review/[id]/+page.svelte", flaiover/src/lib/explorer/changed.ts, flaiover/src/lib/explorer/changed.test.ts, docs/operators/index.md, docs/operators/settings.md, docs/users/flaiover.md, design/system/component-library.md, design/system/flaiover-dashboard.md]
after: [S-0439]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: sum
  seconds: 0
  models: []
  strategic:
    - kind: orchestrator
      seconds: 655
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 66
          output: 1198
          cache_read: 11407911
          cache_write: 26076
          cost: 2.8665
cost_of_delay:
  value: 79.82
  by: planner-E-0026
  at: 2026-10-16T21:50:06Z
forecast:
  duration: 58m
  delivery: 2026-10-17T03:38:00Z
  basis: "Its own forecast of 58m; 1st in the pull order with an in-progress limit of 5, behind S-0441."
  by: flai
  at: 2026-10-17T02:34:36Z
finalized:
  by: orchestrator
  at: 2026-10-16T21:52:04Z
---
# S-0442 The Components page renders a story's worktree, so the operator sees a component an agent is changing live while the story is in progress or in review

## Goal

Changes to the underlying source reach the operator live (outcome 3) where they are made: in a story's worktree, by its agent, before anything is merged. The Components page offers a source: the main checkout, or a story in progress or in review. For a story it runs a second dev server in the story's worktree, so the operator watches a component change as the agent saves it. They then judge it in review with the same props and examples as on main. The review page links to the Components page for a story whose branch changes a component, and lists the changed components first.

## Acceptance criteria

- [ ] `components.start`, `components.status`, `components.stop`, and `flai components serve` take `story`, and run the dev server in that story's worktree with its own port. A story with no worktree is refused as ADR-0145 refuses. A worktree without the sub-project's dependencies installed gets them first, by the sub-project's lock file (`pnpm install --frozen-lockfile`, `npm ci`, or `yarn install --frozen-lockfile`), and a failed install is the start's error with its last lines.
- [ ] Each story's dev server stops when its story leaves review, as well as on the idle stop. At most `components.max_servers` dev servers run per project (3 by default), and a start beyond it is refused with the ones running.
- [ ] The Components page has a source picker, Main or a story in progress or in review, kept in the URL. With a story picked, the catalogue is the worktree's, and the components its branch changes are marked and listed first, from `stream.diff`'s paths.
- [ ] The review page shows "Open in Components" for a story whose branch changes a path the catalogue's globs reach, and opens the page with the story picked.
- [ ] Tests cover the worktree start, the refusal, the install step with a fake package manager, the stop on leaving review, the limit, the picker, and the review link. `docs/operators/index.md`, `docs/operators/settings.md`, `docs/users/flaiover.md`, `design/system/flaiover-dashboard.md`, and `design/system/component-library.md` describe them.

## Tasks
- T-1879 flai/internal/components starts a dev server in a story's worktree, installing the sub-project's dependencies by its lock file when they are missing
- T-1880 The components methods and flai components serve take a story, flai serve stops a story's dev server when it leaves review, and components.max_servers caps them
- T-1881 The Components page's source picker chooses main or a story in progress or in review, and marks and lists first the components the story's branch changes
- T-1882 The review page offers Open in Components for a story whose branch changes a component, opening the page with the story picked
- T-1883 The operator and user docs and the component library's design describe a story's dev server, the install step, the cap, and the review link

## Notes

- Drafted by the planner for E-0026, after S-0439, whose page it extends.
- A story's worktree under `.flai-cache/worktrees/` has no `node_modules` of its own: this repository's `scripts/flaiover-lint.sh` installs them with `scripts/flaiover-install.sh --if-needed` in a worktree (I-0080). flai has to do the same for any project, so the dev server installs by the lock file.

### Planning

- **Touches.** The story declared none; every touch comes from its five tasks, file by file.
  - **From design:** ADR-0145 for resolving a story's worktree; ADR-0019 for where worktrees live.
  - **From layout:** `worktree.go` and S-0438's `devserver.go` in `flai/internal/components`; S-0438's supervision, methods, and commands; the picker beside S-0439's components; the routes and `repo.ts`; the review page at `flaiover/src/routes/review/[id]/+page.svelte`; `changed.ts` beside S-0439's client.
  - **Co-change, added:** `design/system/flaiover-dashboard.md` (33%), for the picker, the routes' `story`, and the review link.
  - **Co-change, not added:** `flai-cli.md` and `flai.md` (37%), since `serve --story` is documented by the reference, which is regenerated; add them if the agent judges the flag needs prose. `flaiover/src/lib/server/agent.ts` (10%), since no method is added.
  - No folder touch.
- **Forecast.** Kept at flai's 58m: size 36, at 96 s a unit over 61 done large feature stories. Delivery is flai's play-out, after S-0439.
- **Cost of delay.** 79.82 USD a week, flai's share of E-0026's 600 USD penalty by forecast duration, kept.
- The forecast carries `by: system-flow`: the CLI wrote it under the host's `FLAI_AGENT`. The planner wrote it.
