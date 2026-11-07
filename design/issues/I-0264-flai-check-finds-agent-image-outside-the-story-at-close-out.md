---
id: I-0264
title: "flai check finds `agent.image` outside the story at close-out"
class: efficiency
status: closed
count: 1
first_reported: 2026-11-06T22:37:17Z
last_reported: 2026-11-06T22:37:17Z
updated: 2026-11-07T00:43:37Z
---

# I-0264 flai check finds `agent.image` outside the story at close-out

## Description
flai check finds `agent.image` outside the story at close-out

## Instances

### 2026-11-06T22:37:17Z
Story: S-0520.
flai check found outside the story:
`system-flow.yaml`: the image flai-agent-sf-claude-code:local of agent.image.layer flai/images/system-flow/Dockerfile is not checked: docker is not on PATH; install Docker to check it, and build it with 'make agent-images'

It came from `scripts/close-out.sh S-0520` run by hand inside `flai-agent-claude-code:local` (`scripts/agent-container.sh S-0520 -- scripts/close-out.sh S-0520`), where `docker` is never on `PATH`. So every close-out inside an agent container records this note as an issue through `--record-issues`, and its commit step then stops on the uncommitted issue file, since commits are the host's (ADR-0253). A check that knows it runs inside a container, or one scoped to a story that changes neither the manifest's `agent.image.layer` nor its Dockerfile, could leave the note out as ADR-0171 does for `components.examples`.

## Remediation

Story S-0582 remediates this issue, created from it at 2026-11-06T23:25:47Z.
Closed 2026-11-07T00:43:37Z: S-0582 (ADR-0255): check.AgentImage marks the three agent.image findings only the host clears, docker not on PATH, docker unable to say, and the layer's image not built, and a check scoped to a story leaves them out unless the story changes system-flow.yaml or the layer's Dockerfile. So a close-out inside an agent container, where docker is never on PATH, no longer records this issue or stops on its uncommitted file. TestScopedCheckLeavesOutAnAgentImageFindingOnlyTheHostClears reproduces S-0520's case for each of the three findings.
