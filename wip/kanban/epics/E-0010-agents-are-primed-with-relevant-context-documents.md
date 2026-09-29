---
id: E-0010
type: epic
nature: feature
title: Agents are primed with relevant context documents
status: in-progress
owner: alex
created: 2026-09-26T07:20:38Z
updated: 2026-09-29T04:57:00Z
transitions:
  - to: ready
    at: 2026-09-26T17:48:00Z
    by: alex
  - to: in-progress
    at: 2026-09-28T23:04:29Z
    by: alex
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
- S-0138 The MCP prime tool returns a story's context pack, built by the same code as flai prime --story
- S-0145 Look for relevant approaches to prime context without overloading agent
- S-0146 A story's context pack fits a size budget, with briefs of what its topics and links select
- S-0147 Agents fetch a design section on demand: doc_search, and a heading on doc_get and flai doc show
- S-0148 Agents flai starts, and any agent with a story, prime with the budgeted pack and fetch what it briefs
- S-0149 Provide briefs in cases where large files are part of context

## Notes
