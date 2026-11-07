---
id: T-2563
type: task
nature: improvement
title: The Klusters page leads an operator through first-time setup and then manages the cluster, with no shell or file edit
status: backlog
parent: S-0586
owner: alex
created: 2026-11-07T08:31:36Z
updated: 2026-11-07T08:31:36Z
transitions: []
stream: S-0586
tags: [kluster, flaiover]
touches: [flaiover/src/routes/ops/klusters/+page.svelte, flaiover/src/routes/ops/klusters/klusters.svelte.test.ts]
after: [T-2562]
---
# T-2563 The Klusters page leads an operator through first-time setup and then manages the cluster, with no shell or file edit

## Work

Waits for T-2562, whose fields and route this page reads.

Rewrite `flaiover/src/routes/ops/klusters/+page.svelte` so that the page holds everything setup and management need. Follow `design/system/design-system.md` for its colours, spacing, and components.

**First-time setup**, shown until the cluster is set up and collapsible afterwards, as steps that each say whether they are done:

1. **Tools.** `kluster`, `docker`, and `kubectl`, each found or missing, from `tools`.
2. **Turn it on.** A button for the `kluster` host action, through `settings.action`, which also writes the block and the config (T-2560). It says what it wrote and the commit.
3. **The cluster.** The config file, edited in place through `file.get` and `file.save`, with its name, profile, and provider read back from `kluster.status`.
4. **Deploying.** Fields for `namespace`, `build`, `deploy`, `teardown`, and `timeout`, each saved through `settings.manifest`, with flai's refusal beside the field it names. With no `deploy`, the page says that the cluster can run but no story is deployed yet. A **Write these with an agent** button makes a backlog story to write the project's build, deploy, and teardown scripts and name them in the block, through the dashboard's existing item write, and links it. The scripts are the project's code, so an agent writes them on a story and the operator reviews them.
5. **Bring it up.** The **Up** button.

**Managing**, once set up:

- What the page has today: the cluster's state, **Up**, **Down**, the running job's log, and the last job.
- A way to turn the action off.
- `findings` from `kluster.status`.
- The nodes and pods from `kluster.workloads`, read when the page opens and after each job.
- A link to Deployments.

Remove every sentence that tells the operator to run a command or edit `system-flow.yaml`: the `kluster-off` notice's command and the `kluster-no-block` notice. With `settings` off, the setup steps that write say that one host action, `settings`, is off, and what it allows. The plan thread settles what that notice says. flai's refusals and errors stay word for word.

If the page grows past what one file carries well, take a part out into a component under `flaiover/src/lib/components/` and give it its `*.examples.json` under `design/components/`, as that folder's `README.md` says.

## Done when

- `flaiover/src/routes/ops/klusters/klusters.svelte.test.ts` walks a project with no block from step 1 to **Up**, and covers a tool missing, a refused field naming its field, `settings` off, no `deploy`, a failed job, and the managing view with nodes and pods.
- No text on the page, and no test's expected text, names a `flai` command or tells the operator to edit a file.
- The story's agent drives the page in a browser against a scratch project, with `kluster` faked or real, and says in the narrative's `## Decisions` which it was and what it saw.
- `flai test` passes on the paths this task changes.

## Notes
