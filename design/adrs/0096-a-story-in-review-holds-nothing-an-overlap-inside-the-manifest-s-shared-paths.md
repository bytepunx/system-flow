---
id: ADR-0096
title: "A story in review holds nothing, an overlap inside the manifest's shared paths never holds, and a story's tasks narrow its folder touches in its claim"
status: accepted
date: 2026-10-06
supersedes: []
superseded_by: []
refines: [ADR-0046]
---

# ADR-0096 A story in review holds nothing, an overlap inside the manifest's shared paths never holds, and a story's tasks narrow its folder touches in its claim

## Context

[ADR-0046](0046-a-ready-story-whose-claim-overlaps-an-open-story-s-is-held-yellow-and-with-its.md) holds a ready story while its claim overlaps an open story's. A story is open while it is in progress or in review. Its claim is its `touches` plus those of its open tasks, and an entry claims everything below it.

[I-0087](../issues/I-0087-one-in-progress-story-holds-every-ready-story-by-overlap-through-folder-wide-touches-and-docs-users-flai-md-so-the-board-runs-one-story-at-a-time-under-a-limit-of-three.md) recorded three occurrences on 2026-10-06. Each time, one open story held every ready story, and the board ran one story at a time under an in-progress limit of three:

- S-0220 held all ten ready stories from 05:58Z to 09:57Z, through `flai/internal/harness`, `flai/internal/serve`, `flai/internal/mcpserver`, and `flaiover/src/routes`. For 2h45m of those four hours it waited in review and no agent ran at all.
- S-0283 held all nine ready stories. Seven were held through the folder `flai/internal/mcpserver`, though S-0283 changed one file there, `permission.go`. Two were held through `docs/users/flai.md`.
- S-0285 held all eight ready stories, through `flai/internal/harness`, `flai/internal/mcpserver`, and `docs/users`. Its touches were a first guess written with the story, before any task named a file.

On TH-0180 the planner also found S-0278 and S-0284 held through `design/adrs`, where every story that records a decision writes a new file.

Three things made the overlap total. A story in review holds until the operator accepts it, though its branch is finished and synced. A folder touch holds every file in the folder, even once the story's tasks name the few files it changes. And a few paths are changed by almost every story, each in its own section or a new file: `docs/users/flai.md`, `design/system/flai-cli.md`, `design/adrs`, and `design/issues`.

On TH-0180 (2026-10-06) the planner offered the three remedies as one ADR refining ADR-0046. It also offered rule 1 as a story of its own, and a split of `docs/users/flai.md` by command in place of a shared list. The operator took all three as one. They asked that operators of any project can edit the shared list through the CLI, HTTP, MCP, and the dashboard's project settings, with globs.

## Decision

**Only a story in progress holds a ready story: a story in review holds nothing. An overlap that lies wholly inside a pattern of the manifest's new `claims.shared` list never holds. In a story's claim, a folder touch that holds touches of the story's tasks is replaced by those touches.** This refines ADR-0046. The rest of ADR-0046 stands: the reasons, `after:`, skip-ahead, the operator's word, the trial merge at sync, and the notice at acceptance.

### 1. A story in review holds nothing

- Only a story in progress holds a ready story, by `overlap` or by `no-touches`. A story in review holds none.
- A story in review keeps its claim for everything else ADR-0046 gives a claim: `flai check`'s `wip.overlap`, the trial merge at `flai stream sync`, and the overlap notice at acceptance.
- Its branch is finished and synced. When it is accepted, the notice tells each open story whose claim covers a changed path which paths those are, and that story's next `flai stream sync` rebases onto them.
- flai serve's launcher is unchanged. To it, a story in ready whose agent it has started counts as in progress, because the agent moves it there next. So one look never starts two stories that overlap.

### 2. Shared paths never hold

**The key.** `claims.shared` in `system-flow.yaml` is a list of glob patterns. Each names paths that many stories change in separate sections or new files.

```yaml
claims:
  shared:
    - docs/users/flai.md
    - design/adrs
```

**The dialect.** Paths are relative to the repository root and separated by `/`.

| In a pattern | Matches |
|--------------|---------|
| no glob character | the path itself and everything below it, so `design/adrs` and `design/adrs/**` mean the same |
| `*` | any characters within one segment |
| `**` | zero or more whole segments |
| `?` | one character within a segment |

**The rule.** A touches entry is first read as a path, a component name becoming its sub-project's path as ADR-0046 says. When two entries overlap, the narrower of the two is the deeper one, or either when they are equal. The pair does not hold when the narrower lies wholly inside a shared pattern. A file lies inside a pattern that matches it. A folder lies inside only when the whole folder does: `docs/users` is inside `docs/users/**` and `docs` is not, and `docs/users` is not inside `docs/users/*.md`, because the folder may hold other files. A ready story is held by overlap only through a pair that does not lie inside.

**Where it applies.** The same exception applies to `flai check`'s `wip.overlap` and to the notice a write gives when it grows a claim onto another open story's. It does not apply to the trial merge at `flai stream sync` or to the overlap notice at acceptance. Both still report on shared paths, so a real textual conflict there is still caught.

**Validation.** The manifest's checks reject an empty pattern, an absolute path, a `..` segment, and a malformed glob. `flai check` reports each as an error on its line, and `flai shared add` refuses it. A pattern that is not valid frees nothing.

**Defaults.** This project's list is `docs/users/flai.md`, `docs/users/flai-reference.md`, `design/system/flai-cli.md`, `design/adrs`, `design/issues`, and `template/CHANGELOG.md`. A new project from the template starts with `<design>/adrs` and `<design>/issues` under its design layout folder, by default `design/adrs` and `design/issues`. The planner proposed this project's list for every project, but only those two come from the template's layout. The other four are this repository's own files.

**Editing the list.** Every surface reads, checks, and changes the same list. Add and remove rewrite only `claims.shared`, keep the file's other keys and comments, refuse a duplicate or a pattern that is not valid, and say what changed.

| Surface | Read and check | Change |
|---------|----------------|--------|
| CLI | `flai shared list`; `flai shared check` with paths, touches entries, or a story ID, reporting each entry, a story's claim entry by entry, and which pattern matched | `flai shared add <pattern>`, `flai shared remove <pattern>` |
| Host API | `settings.get` carries the list; `settings.shared_check` checks paths or a story, gated like `settings.get` | `settings.shared` adds or removes one pattern, gated by the `settings` host action, through `flai shared add` or `remove` ([ADR-0016](0016-dashboard-delegates-to-flai.md)) |
| Dashboard | the project settings panel lists the patterns and checks a path | the same panel adds and removes a pattern |
| MCP | `shared_paths` lists and checks, for every agent | `shared_paths_edit` adds or removes |

**Who may change it.** The list decides what holds, so an agent that could add to it could free its own story. `flai guard` refuses `shared_paths_edit`, and `flai shared add` and `remove`, to every session flai serve starts: a story's agent and its sub-agents, the planner, the orchestrator, and the analyzer. The refusal says to ask the operator on a thread. The operator's own session may change the list. The orchestrator is granted nothing here. Listing and checking are reads, open to every agent.

### 3. Tasks narrow a folder claim

A story's claim is worked out in three steps, each entry read as a path first:

1. A story touch that is a folder holding at least one touch of the story's tasks is replaced by those task touches. Done tasks count, because the branch changed their files. Cancelled tasks do not.
2. A story touch that no task names inside stays whole.
3. A touch of an open task that lies outside every story touch is added, as ADR-0046 says.

This changes ADR-0046's claim, which counted only open tasks: a done task's touches now stay claimed where they narrow a folder. `flai stream sync`'s out-of-claim report reads the narrowed claim. It reports a file a task changed outside it, which tells the agent to widen the task's touches. The planner and the story's agent are told, in `strategic-agents.md` and `work-management.md`, to declare touches file by file. They keep a folder only where files that no task can name yet may be added.

## Consequences

- Rule 1 alone frees the 2h45m in which S-0220 waited in review and nothing ran.
- Rule 3 frees the seven stories S-0283 held through `flai/internal/mcpserver`, once a task named `permission.go`. Rule 2 frees the two it held through `docs/users/flai.md`.
- Rule 2 frees S-0278 and S-0284 from `design/adrs`. S-0285 is freed by rules 2 and 3 together: its `docs/users` touch narrows to the files its tasks name, and those are shared.
- A shared path can carry a real textual conflict. It is caught only by the trial merge at sync and the notice at acceptance, after both stories have started, and an agent resolves it at its next sync.
- Narrowing makes the claim depend on the tasks' touches being right. A task that changes a file it did not name lets an overlapping story start. The out-of-claim report at sync is the net, and it reports the file to the agent that changed it.
- A story's claim can now shrink as its tasks are written. It changes only through the story's own items.
- An older flai ignores `claims.shared`, because the manifest is decoded leniently, and holds as ADR-0046 says, review included. The key is not a front-matter field, so it raises no `flai.minimum`. A host whose flai serve is older holds as before until it is upgraded.
- The implementation is split into tasks under S-0295: review holds nothing, the manifest key with its matcher and edit, the conventions, the narrowing, the CLI, the MCP tools and the guard, the exception in the hold and `flai check`, the host API, the dashboard, and the guides.

## Alternatives considered

- **Split `docs/users/flai.md` by command.** Declined on TH-0180. The shared list covers it, and `design/adrs` and `design/issues` too, which no split helps.
- **Rule 1 as a story of its own.** Offered on TH-0180 as option (b), and declined.
- **An empty default list, set per project.** Declined. The template's own layout gives every project `design/adrs` and `design/issues`, which every story writes new files in.
- **Let agents edit the list.** Declined: an agent could free its own story by adding the paths it holds others through.
- **A story in review holds as before.** This is what I-0087 measured: four hours with one story running, 2h45m of them with none.
