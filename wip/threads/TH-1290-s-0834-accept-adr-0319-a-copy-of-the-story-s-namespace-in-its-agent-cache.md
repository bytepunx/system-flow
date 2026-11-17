---
id: TH-1290
title: "S-0834: accept ADR-0319, a copy of the story's namespace in its agent cache folder, before any code is written?"
anchor:
  path: wip/kanban/stories/S-0834-flai-kluster-exec-and-log-read-flai-cache-kluster-which-a-story-s-agent-in-a.md
  item: S-0834
status: open
participants: [agent-S-0834, orchestrator]
created: 2026-11-17T12:07:30Z
updated: 2026-11-17T12:08:37Z
---

# TH-1290 S-0834: accept ADR-0319, a copy of the story's namespace in its agent cache folder, before any code is written?

On wip/kanban/stories/S-0834-flai-kluster-exec-and-log-read-flai-cache-kluster-which-a-story-s-agent-in-a.md.

## Entries

### 2026-11-17T12:07:30Z agent-S-0834
**Recommended: accept ADR-0319 as written.** T-3546 needs your acceptance before T-3547 to T-3550 build it, as I-0361 asks.

The ADR is `design/adrs/0319-flai-copies-a-story-namespace-s-record-step-logs-state-and-kubeconfig-into-the.md` on `story/S-0834` (commit 9d7e135), status proposed.

**What it decides:**

1. **A copy.** Each time flai writes a story namespace's record, it copies that namespace's record, step logs, state folder, and kubeconfig into `.flai-cache/agents/S-nnnn/kluster/`. The story's container already mounts that folder. The copy holds nothing of another namespace and no cluster kubeconfig. A passing teardown removes it.
2. **The readers.** `flai kluster log`, `env`, and `exec` read the main checkout's `.flai-cache/kluster` when it is there. Otherwise they read the copy. The kluster config is read from the worktree when the main checkout's root lacks it.
3. **The API server.** On any network but `host`, the container gets `FLAI_KLUSTER_API_HOST=host.docker.internal`. Containers already get `--add-host host.docker.internal:host-gateway`. `env` and `exec` then give a kubeconfig whose server names that host, with `tls-server-name` set to the original host, so the certificate is still verified.
4. **No new mount.** ADR-0253's list stays.

**What I found that changes the plan a little:**

- The worktree is mounted at its host path. So `exec` needs no change of directory, only the config read and the copy.
- `--job` requests from a container's shell still fail: their request folder is not mounted. The MCP tools work. The ADR names this as outside its scope. In T-3549 I will check it and record an issue if it holds.

**Risk.** Whether `tls-server-name` verifies on `bridge` depends on k3d's certificate naming the host the kubeconfig gives. This repository runs no kluster, so the tests cover the rewrite and the layout, not a live API server.

**Other answers:**

- **Change it:** say what. I rewrite the ADR, then build.
- **Reject it:** say which alternative instead. The ADR lists five, for example mounting `.flai-cache/kluster` read-only.

### 2026-11-17T12:08:37Z orchestrator (recommendation)
**Recommended: accept ADR-0319 as written** (`story/S-0834`, eaa744e87), so that T-3547 to T-3550 can build it.

**Why:**

- ADR-0292 §2 already lets a story's agent and its sub-agents run `flai kluster env`, `exec`, and `log` for their own namespace. In a container those commands find nothing (I-0361). ADR-0319 makes that grant work there, and it grants nothing new.
- It keeps ADR-0292's boundary. The copy holds only the namespace's own kubeconfig, which the API server limits to that namespace. It never holds flai's cluster kubeconfig, and the guard's refusal of `.flai-cache/kluster/kubeconfig` is unchanged.
- It keeps ADR-0253's mount list. The copy goes in `.flai-cache/agents/S-nnnn`, which each story's container mounts already. The ADR's Alternatives explain why a bind mount of the story's own files would go stale, since they are written by rename.
- `tls-server-name` keeps the API server's certificate verified when the host name is rewritten.

**Why I am asking you:** this is a new ADR, so it is a decision no one has recorded yet. The agent names one risk, and its tests cannot cover it here: `bridge` verification depends on k3d's certificate. That risk sits in the ADR's Consequences.

**Alternatives:** change the ADR (say what), or take one of its five alternatives, such as mounting `.flai-cache/kluster` read-only. I would not take that one, because it exposes the cluster's kubeconfig and every namespace's token.

Source: design/adrs/0292-a-story-s-agent-asks-flai-serve-to-deploy-check-and-tear-down-its-own-story-s.md § Decision
