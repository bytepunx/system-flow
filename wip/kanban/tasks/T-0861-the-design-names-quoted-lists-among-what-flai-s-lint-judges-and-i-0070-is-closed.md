---
id: T-0861
type: task
nature: remediation
title: The design names quoted lists among what flai's lint judges, and I-0070 is closed
status: backlog
parent: S-0258
owner: alex
created: 2026-10-05T03:13:35Z
updated: 2026-10-05T03:13:35Z
transitions: []
stream: S-0258
tags: [flai, mdlint]
touches: [design/system/flai-cli.md, design/issues]
after: [T-0860]
---
# T-0861 The design names quoted lists among what flai's lint judges, and I-0070 is closed

## Work

1. In `design/system/flai-cli.md`, add S-0258 to the layout line for `mdlint/`, which names the stories that changed the lint (S-0179, MD007 S-0240), and say that it parses blockquote content. Change any sentence in that document that says flai's lint skips what is inside a quote.
2. Close I-0070 with `flai issue close I-0070 --reason`. Run it in the story's worktree, so that it writes under the story's checkout. The reason says what T-0860 changed and names the fixture that reproduces the issue.

This task waits for T-0860, because the close reason and the design describe the fix that T-0860 makes.

## Done when

- `design/system/flai-cli.md` says that flai's lint judges lists inside blockquotes (S-0258).
- I-0070's status is closed, with a reason that names the parser change and the fixture.
- `design/issues/summary.md` is current.
- `flai check --strict` reports nothing on the story.

## Notes

Drafted by planner-S-0258.
