---
id: D-0632
title: "A workflow rule's check gives one refusal text and one warning text, the same words on any transition that lists it, and the transition's severity picks which; criteria-ticked alone keeps two of its own"
kind: decision
sources: [S-0738]
topics: [flai, go, code]
paths: [flai/internal/workitem/rules.go, flai/internal/workflow/rules.go]
date: 2026-11-15
status: active
absorbed_by: ""
asked_again: []
cited: [TH-1218, TH-1219, TH-1220]
---

# D-0632 A workflow rule's check gives one refusal text and one warning text, the same words on any transition that lists it, and the transition's severity picks which; criteria-ticked alone keeps two of its own

S-0738 made `Repo.Move` run the rules a transition lists, each refusing or warning by its severity (ADR-0301 §4). Its narrative's Decisions (`wip/archive/agents/S-0738.md`) record what a rule says:

> T-3202: a rule's check gives a refusal text and a warning text, alike for all but criteria-ticked, which keeps today's two (refuse at done, warn at review); a transition's severity picks one. A worktree git cannot read always warns.
  
> Accepted as it is: a rule an operator's workflow puts on another transition keeps its words, so has-task says 'before it goes to review' and criteria-ticked at warn says 'goes to review with unticked acceptance criteria' wherever they are listed. Today's words are what criterion 3 holds; texts per transition can come with a story that asks for them.

So a rule's message does not follow the transition it sits on: an operator's workflow that lists `has-task` or `criteria-ticked` elsewhere still names review. A test of a rule's words reads the one text per severity. Words per transition are a story of their own, when one asks for them.

Source: S-0738's narrative, `wip/archive/agents/S-0738.md`, § Decisions.
