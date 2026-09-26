---
id: E-0010
type: epic
nature: feature
title: Agents are primed with relevant context documents
status: backlog
owner: alex
created: 2026-09-26T07:20:38Z
updated: 2026-09-26T17:46:25Z
transitions: []
tags: [cli]
touches: [flai/cmd]
---
# E-0010 Agents are primed with relevant context documents

## Outcome

Instead of feeding agents dispatched to work stories with every ADR and convention document, flak should be able to determine relevant documentation for a given story based on the component being worked on, the technologies involved, and the architecture affected.

## Stories
- S-0125 Determine how to prime context with relevant documentation only
- S-0134 Conventions, design, tech files, and ADRs carry topics on the file and on headings, and flai check keeps them honest
- S-0135 Stories and epics carry topics, and flai works out a story's topics from them, its epic, and the projects its tags and claim reach
- S-0136 flai prime --story prints the conventions a story's topics select, section by section, and lists what it left out
- S-0137 flai prime --story adds the design, tech, and ADRs a story's topics, links, and ranking select, with a catalog of the rest
- S-0138 Agents flai starts, and any agent with a story, prime with flai prime --story

## Notes
