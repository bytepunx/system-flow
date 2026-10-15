---
id: T-1540
type: task
nature: feature
title: An isolated preview frame styles a gallery of sample components from a set of tokens, in light and dark, apart from the dashboard's theme
status: done
parent: S-0382
owner: alex
created: 2026-10-13T17:47:07Z
updated: 2026-10-15T01:46:19Z
transitions:
  - to: ready
    at: 2026-10-15T01:34:25Z
    by: agent-S-0382
  - to: in-progress
    at: 2026-10-15T01:34:27Z
    by: agent-S-0382
  - to: done
    at: 2026-10-15T01:46:19Z
    by: agent-S-0382
stream: S-0382
tags: [dashboard]
touches: [flaiover/src/lib/design/preview.ts, flaiover/src/lib/design/preview.test.ts, flaiover/src/lib/design/gallery.css, flaiover/src/lib/components/design/PreviewFrame.svelte, flaiover/src/lib/components/design/PreviewFrame.svelte.test.ts, flaiover/src/lib/components/design/Gallery.svelte]
usage:
  source: log
  seconds: 712
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 133
      output: 50206
      cache_read: 7611588
      cache_write: 219989
      cost: 3.8134
---
# T-1540 An isolated preview frame styles a gallery of sample components from a set of tokens, in light and dark, apart from the dashboard's theme

## Work

`preview.ts` turns tokens into the custom properties the gallery reads, and works out each colour's contrast against the ground, with tests. `Gallery.svelte` holds the sample components: buttons, inputs, a form, a card, an alert, a table, navigation, and headings with body text, styled only through those properties. `PreviewFrame.svelte` renders the gallery isolated, in the way S-0380 decided (an iframe or a shadow root), in light and dark side by side, and restyles it when the tokens it is given change, with no build. The component test checks that a changed token reaches the frame and that the dashboard's tokens do not. It waits for nothing.

## Done when

- The frame shows the gallery in both themes from given tokens, follows a change at once, and its tests pass.

## Notes

- Layer 1.
