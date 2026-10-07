---
id: T-1267
type: task
nature: research
title: Show a headless session flai serve starts offering subagent_type fork, and record its Claude Code version
status: backlog
parent: S-0241
owner: alex
created: 2026-10-07T23:12:43Z
updated: 2026-10-07T23:12:43Z
transitions: []
stream: S-0241
tags: [template]
touches: [design/system/agent-context.md]
---
# T-1267 Show a headless session flai serve starts offering subagent_type fork, and record its Claude Code version

## Work

Start a `claude -p` session with the command `flai serve` builds for a story's agent (harness `claude-code`, model `claude-opus-5-5`, effort `high`, the same environment and permission prompt), as S-0230 did on 2026-10-02. In it, call the Agent tool with `subagent_type: fork` on a trivial prompt, and list the agent types it offers. Record the Claude Code version (`claude --version`), the types offered, and whether the fork started, in `design/system/agent-context.md` § Tasks in parallel › Forks are not offered headless; rename that heading if forks are now offered.

If forks are still refused, record the version and the refusal, block the story with `flai block --reason`, and stop: the other tasks cannot be worked.

This task waits for no other.

## Done when

- `design/system/agent-context.md` names the Claude Code version whose headless session offered `subagent_type: fork`, and the agent types it listed.
- A fork was started in that session and returned, and the section says so.

## Notes

Claude Code 2.1.290, the planner's headless session on 2026-10-07, still offered no fork.
