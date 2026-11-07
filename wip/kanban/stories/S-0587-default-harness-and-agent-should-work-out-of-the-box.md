---
id: S-0587
type: story
nature: remediation
title: Default harness and agent should work out of the box
status: backlog
owner: alex
created: 2026-11-07T08:40:55Z
updated: 2026-11-07T08:40:55Z
transitions: []
tags: []
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
---
# S-0587 Default harness and agent should work out of the box

## Goal

In new projects, once an agent has been enabled, new stories that an operator tries to start the agent for report the following error message:

```text
S-#### names no harness, and no command is set on the host (flai serve agent set -- <program> [args...])
```

This indicates a failure of flai to provide the default settings (the system already has a default harness and agent selected).

Finally, if it were missing, the error message to direct the user to the configuration pages or sub-pages necessary to complete configuration.

## Acceptance criteria
- [ ] A default agent should always be provided to stories
- [ ] A default harness should always be provided to stories
- [ ] Error messages must provide a link to the page(s) where the user can enter the missing information if no default exists
- [ ] Harness and agent settings should be independent of other projects (setting these for one should never change another project's harness or agent)

## Tasks

## Notes
