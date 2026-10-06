---
id: ADR-0097
title: "permission_prompt takes an answer from the story's owner or the project's owner"
status: accepted
date: 2026-10-06
supersedes: []
superseded_by: []
refines: [ADR-0086]
topics: [cli]
---

# ADR-0097 permission_prompt takes an answer from the story's owner or the project's owner

## Context

[ADR-0086](0086-flai-serve-gives-a-claude-code-agent-flai-s-permission-prompt-as-its-permission.md) has flai's `permission_prompt` ask before a write under `.claude/` in an in-progress story's worktree. Its decision 2 takes an answer only from the story's owner: "Entries by anyone else, other agents included, are not answers."

[I-0081](../issues/I-0081-flai-s-permission-prompt-waits-for-an-answer-from-the-story-s-owner-and-the-operator-s-thread-replies-carry-another-name-so-their-allow-is-never-seen-and-the-write-times-out.md) counts two instances, from S-0218 and S-0222. In both the story's owner was `arobson`, the name the story was made under. The operator replied `allow` as `alex`, the `owner` in `system-flow.yaml` and the name the operator's flai writes threads as (`author` in its config). The prompt never saw an answer. The write stayed held until Claude Code's MCP idle timeout ended it after thirty minutes, and the agent stood blocked all that time. The operator then had to copy the file in by hand.

A story's owner is the project's owner when flai makes the story, so the two names differ only when a story was made under an older owner or by hand. Each name is one the operator chose. flai has no list of agent names, so it cannot tell a person's name from an agent's.

On TH-0164 (2026-10-06) the planner proposed letting the project's owner answer as well, and the operator accepted the plan.

## Decision

**An answer to a permission thread may come from the story's owner or from the project's owner, the manifest's `owner`.** This refines decision 2 of ADR-0086 and changes nothing else in it.

- The first entry after the request by either person answers it. allow, yes, approve, approved, or ok as its first word lets the write through, and anything else refuses it with their words as the reason.
- An entry by the asking agent is never an answer, even if the agent's name is one of the owners. Entries by anyone else, other agents included, are still not answers.
- When the story has no owner and the manifest names none, any entry but the agent's answers, as before.
- The request names whoever may answer: the story's owner, and the project's owner too when the two differ.

## Consequences

- The operator's reply under the project's owner name lets the write through, whatever name the story was made under. No agent sits blocked for thirty minutes on an answer it cannot see.
- One more person can answer: the project's owner. That is the operator who chose the manifest, so it gives an agent no permission the operator did not give.
- A story made by hand under a third name still needs that name, or the project's owner, to answer.
- `auto-approve` and the rest of ADR-0086 are unchanged.

## Alternatives considered

- **Any entry not by an agent answers.** flai keeps no list of agent names, so it cannot tell which entries are a person's, and another agent could then allow the write.
- **Any entry but the asking agent's answers.** Another story's agent, or a sub-agent replying under a name of its own, could then allow a write to the project's hooks.
- **Rewrite story owners to the manifest's owner.** It mends the stories already made but not one made by hand later, and it rewrites history that the owner field records.
