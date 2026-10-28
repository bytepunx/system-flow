---
id: D-0386
title: "Which story does Claude Code work inside the container for T-2166: Recommended for the operator: answer \"as recommended\""
kind: answer
sources: [TH-0676, S-0505]
topics: [harness, docker, experiments, flai, cli, go, code, all]
paths: [flai/images/bedrock/Dockerfile, flai/images/claude-code/Dockerfile, scripts/agent-image.sh, scripts/agent-container.sh, Makefile, scripts/README.md, ".dockerignore", design/experiments/S-0505-claude-code-works-one-story-of-this-repository-inside-a-hand-built-container-from-a-bedrock-image-and-a-claude-code-layer-and-the-results-settle-what-adr-0203-leaves-open.md, design/experiments/README.md, design/system/agent-adapters.md, design/tech/docker.md, design/issues/summary.md, design/issues]
date: 2026-10-28
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0386 Which story does Claude Code work inside the container for T-2166: Recommended for the operator: answer "as recommended"

The operator, alex, answered TH-0676 at 2026-10-28T12:11:11Z.

What was asked:

```text
S-0505's criterion 3 needs one story of this repository, chosen with you, worked inside the container from its first task to review. The story's notes ask for a documentation-only story, since the images carry no Go toolchain. No ready story is documentation-only: the backlog's are all Go or flaiover work.
```

The answer:

```text
**Recommended for the operator:** answer "as recommended". Let agent-S-0505 create the documentation-only runbook story under E-0029, finalized, and move it to ready.

## Why I recommend it

- S-0505's notes ask for a documentation-only story unless you name another, because the bedrock and claude-code images carry no Go toolchain. No ready or backlog story is documentation-only today.
- A runbook for building the agent images and running an agent in a container by hand is work E-0029 needs anyway. A throwaway story would spend the same subscription session and leave nothing behind.
- Its two touches, `docs/operators/runbooks/agent-container.md` (new) and `docs/operators/runbooks/README.md`, overlap no open story, so it can start as soon as T-2166 needs it.

## Why it is yours

- Criterion 3 says the story is "chosen with the operator".
- The run spends one session of your Claude subscription, about 30 to 60 minutes, and mounts your login read-only.
- A new story is a change of scope. `work-management.md` § Rules keeps an agent from inventing scope.

**The alternatives:** name an existing story, or answer "throwaway" and the agent cancels the story after the run.

S-0505 goes on building the images (T-2149) and the container script (T-2159) meanwhile. Only T-2166 waits on your answer.

Source: design/conventions/work-management.md § Rules
```

Source: TH-0676.

The outcome: the operator confirmed the recommendation, and S-0513, the runbook story under E-0029, was the story worked inside the container. S-0505 started it on the host with `flai story start`, because flai serve starts a ready story's agent within a second; see S-0505's results in design/experiments.

Source: S-0505.
