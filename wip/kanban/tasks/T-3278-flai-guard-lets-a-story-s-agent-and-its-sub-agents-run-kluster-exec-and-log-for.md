---
id: T-3278
type: task
nature: feature
title: flai guard lets a story's agent and its sub-agents run kluster exec and log for their own story, and refuses other namespaces, the kluster program but status, and flai's cluster kubeconfig
status: backlog
parent: S-0763
owner: alex
created: 2026-11-13T09:03:11Z
updated: 2026-11-13T09:03:11Z
transitions: []
stream: S-0763
tags: [cli]
touches: [flai/internal/guard/guard_rules.go, flai/internal/guard/call.go, flai/internal/guard/kluster_test.go]
---
# T-3278 flai guard lets a story's agent and its sub-agents run kluster exec and log for their own story, and refuses other namespaces, the kluster program but status, and flai's cluster kubeconfig

## Work

Criterion 5, as ADR-0292 §2 lists it. Today `exec` and `log` fall through every kluster rule, and no rule names the `kluster` program or `.flai-cache/kluster/kubeconfig`.

In `guard_rules.go`, beside `klustering`, `klusterWrite`, and `ownStory`:

- `flai kluster exec S-nnnn -- <command>` and `flai kluster log S-nnnn --step <step>`: let a story's agent and its sub-agents run them for their own story, the story from `FLAI_STORY` or the name `agent-S-nnnn`, as `ownStory` reads it. Refuse them for another story, for no story, for a namespace that is not a story's (the default one), and for any session with no story of its own, the planner and the orchestrator among them, saying the command is a story's agent's for its own namespace. Read the flags as the request rules do (flai's own flags before `kluster`, `scripts/flai.sh`, a path, `cd <dir> &&`, `bash -c`), and stop reading `exec`'s own words at `--`: what follows is the command it runs, read by the line rules below and nothing else.
- Every agent and sub-agent: refuse a line that runs the `kluster` program but `kluster status`, in the forms a shell line takes (a path, after `&&` or `;`, in `bash -c`), saying the cluster is flai's and naming `flai kluster status` and `flai kluster exec`.
- Every agent and sub-agent: refuse a line that names `.flai-cache/kluster/kubeconfig`, whole or as the end of a path, saying it is flai's kubeconfig for the whole cluster and that `flai kluster exec` gives the namespace's own. A line naming `.flai-cache/kluster/agents/` is not refused by this rule.
- Hook the two line rules into the Bash call in `call.go` where `klusterWrite` is called.

The operator's own session (no `FLAI_STARTED_BY`, `FLAI_ROLE`, or `FLAI_STORY`) runs all of it, as it runs the writes. Tests in `kluster_test.go`: each form allowed for the own story, by the agent and a sub-agent; each refusal, with its reason; the operator's session passing.

Waits for nothing: it reads only the command forms the story's criteria fix, and changes no file another task changes.

## Done when

- A story's agent and its sub-agents may run `flai kluster exec` and `flai kluster log` for their own story, and are refused them for another story or the default namespace.
- Every agent and sub-agent is refused the `kluster` program but `status`, and a line naming `.flai-cache/kluster/kubeconfig`.
- `go test ./internal/guard/` passes.

## Notes
