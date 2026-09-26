---
id: E-0009
type: epic
nature: feature
title: Agent Safety and Coordination
status: ready
owner: alex
created: 2026-09-26T07:10:42Z
updated: 2026-09-26T17:47:56Z
transitions:
  - to: ready
    at: 2026-09-26T17:47:56Z
    by: alex
tags: [cli]
touches: [flai/cmd]
---
# E-0009 Agent Safety and Coordination

## Outcome

Right now, the operator can drag stories into ready and they will always be worked on if there is capacity. There aren’t any checks whether agents will make conflicting changes in the same files or whether their work is compatible.

flai needs mechanisms in place to avoid parallel changes to the same files.

## Stories
- S-0124 Survey state of the art coordination mechanisms
- S-0128 flai serve and wait_for_work hold a ready story whose claim overlaps an open story's, and say why
- S-0129 A held story's card is yellow on the board and its page says why it waits
- S-0130 A story that names another in after: waits until that story is done
- S-0131 flai stream sync reports conflicts with other open story branches and changes outside the story's touches
- S-0132 Accepting a story tells every open story that overlaps it which paths changed

## Notes
