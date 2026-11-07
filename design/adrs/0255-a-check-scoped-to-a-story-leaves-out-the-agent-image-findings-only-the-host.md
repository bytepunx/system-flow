---
id: ADR-0255
title: "A check scoped to a story leaves out the agent.image findings only the host clears, docker not on PATH, docker unable to say, and the layer's image not built, unless the story changes system-flow.yaml or the layer's Dockerfile"
status: accepted
date: 2026-11-07
supersedes: []
superseded_by: []
refines: [ADR-0085, ADR-0171]
---

# ADR-0255 A check scoped to a story leaves out the agent.image findings only the host clears, docker not on PATH, docker unable to say, and the layer's image not built, unless the story changes system-flow.yaml or the layer's Dockerfile

## Context

S-0510 gave `flai verify` the rule `agent.image` (`flai/internal/check/image.go`, `AgentImage`). When `system-flow.yaml` names `agent.image.layer`, it gives one of four warnings, each on the manifest's `layer:` line:

| Finding | Advisory | What clears it |
|---------|----------|----------------|
| the layer is not a file | no | the story: fix the path or add the Dockerfile |
| docker is not on `PATH` | yes | the host: install Docker |
| docker could not say whether the image is built | yes | the host: start Docker |
| the layer's image is not built | yes | the host: `make agent-images` |

[I-0264](../issues/I-0264-flai-check-finds-agent-image-outside-the-story-at-close-out.md) has one instance, from S-0520's close-out run inside the agent container on 2026-11-06. There `docker` is never on `PATH`, so the check gave the note on `system-flow.yaml`, which S-0520 did not change. It was outside the story, so `--record-issues` recorded it ([ADR-0085](0085-a-close-out-s-flai-check-reports-findings-outside-the-story-as-notes-and.md)). The close-out's commit step then stopped on the uncommitted issue file, since commits inside a container are the host's ([ADR-0253](0253-a-container-that-runs-a-story-s-agent-mounts-the-main-checkout-s-wip-and-system.md)). So every close-out inside an agent container stops this way. On a host without the layer's image built, every close-out records the not-built warning the same way.

None of the three advisory findings is the closing story's to clear: the manifest and the Dockerfile are right, and only the host's Docker and its built image are missing. [ADR-0171](0171-a-sub-project-whose-svelte-is-not-installed-gets-one-components-examples-note.md) settled the same case for `components.examples`: a note that only installing a tool clears is left out of a scoped check unless the story changes what the note is about.

## Decision

A check scoped to a story leaves out the three `agent.image` findings only the host clears, docker not on `PATH`, docker unable to say, and the layer's image not built, unless the story changes `system-flow.yaml` or the layer's Dockerfile.

1. **Which findings.** The three advisory `agent.image` findings are marked, as ADR-0171's notes are, with the paths whose change makes them the story's: the manifest and the layer's Dockerfile. The warning for a layer that is not a file is not marked: a story can fix it, and it stays a finding on the manifest, inside or outside the story, as before.
2. **Left out unless the story changes one of those paths.** `flai check --story S-nnnn`, and so `flai verify` and the close-out, drop a marked finding when no path the story's branch changes, or that is uncommitted in its worktree, is the manifest or the Dockerfile, or a folder holding either, and take back the counts it added. When one is, the finding is inside the story and stays, advisory, so `--strict` passes over it: the story changed the image, and its agent should know the image went unchecked or unbuilt.
3. **`--record-issues` records none of them outside the story.** What is left out of the result is not recorded, so a close-out inside an agent container no longer stops on an issue file it cannot commit.
4. **Without `--story` they are reported.** The main checkout's `flai check`, and `flai verify` on the host, still say that the image is unchecked or unbuilt.

This refines ADR-0085, whose rule that every finding outside the story is recorded in an issue now leaves these out, and ADR-0171, whose mark for a note only the host clears now names the paths whose change brings it inside, beside the examples folder it names.

## Consequences

- A close-out inside an agent container, where docker is never on `PATH`, no longer records `agent.image` or stops on its issue file. I-0264 can be closed.
- A close-out on a host whose layer image is not built no longer records the not-built warning.
- A story that changes `system-flow.yaml` or the layer's Dockerfile still sees the finding, advisory, inside the story.
- A layer path that is not a file is still a warning wherever it is found.

## Alternatives considered

- **Leave the finding out whenever the check runs inside a container.** It needs flai to know it is in one, and it would still record the not-built warning on a host without the image.
- **Leave out every `agent.image` finding outside the story.** It would hide a broken layer path, which a story can fix.
- **Record the note but let the close-out commit inside the container.** Commits are the host's (ADR-0253), and the issue would only count how often a close-out ran inside a container.
