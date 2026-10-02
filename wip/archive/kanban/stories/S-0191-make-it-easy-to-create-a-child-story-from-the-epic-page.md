---
id: S-0191
type: story
nature: feature
title: Make it easy to create a child story from the epic page
status: done
owner: alex
created: 2026-10-01T11:10:50Z
updated: 2026-10-02T16:27:23Z
transitions:
  - to: ready
    at: 2026-10-01T11:10:50Z
    by: alex
  - to: in-progress
    at: 2026-10-02T16:07:23Z
    by: agent-S-0191
  - to: review
    at: 2026-10-02T16:15:19Z
    by: agent-S-0191
  - to: done
    at: 2026-10-02T16:27:23Z
    by: alex
tags: [dashboard]
topics: [client-side-epic-page]
touches: [flaiover/src, design/system/flaiover-dashboard.md, docs/users/flaiover.md, design/issues/I-0056-flai-s-wip-markdown-lint-does-not-flag-a-bare-email-address-so-a-thread-entry-with-one-reached-main-and-failed-a-story-s-close-out.md, design/issues/I-0057-flai-check-strict-stops-a-story-s-close-out-on-wip-findings-outside-the-story.md, design/issues/summary.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 501
  models:
    - model: claude-opus-5-5
      input: 146
      output: 32595
      cache_read: 6372890
      cache_write: 118485
      cost: 2.8749
    - model: claude-sonnet-5-5
      input: 32
      output: 7860
      cache_read: 390215
      cache_write: 106351
      cost: 0.4226
---
# S-0191 Make it easy to create a child story from the epic page

## Goal

Add a `Create story` button on the epic page after the `New epic` button.

## Acceptance criteria
- [x] Navigates to the new story page with the parent epic pre-selected

## Tasks
- T-0700 An open epic's page links Create story to the new story form under it
- T-0701 Main's markdown lint passes again and I-0056 records why flai let TH-0067's bare email through
- T-0702 I-0057 counts close-outs stopped by flai check findings outside the story

## Notes

- The criterion is verified by `item.svelte.test.ts`: an open epic's **Create story** links `/new?type=story&parent=<epic>`, and `/new` hands `?parent=` to the form as its chosen epic (S-0171). A done, cancelled, or archived epic has no such link, because the form offers only open epics as parents.
- A close-out stopped once at `flai check --strict` on a warning outside the story, `threads.archived` on TH-0067 (answered on archived S-0231), recorded in I-0057 as TH-0056 asks. TH-0067 was resolved afterwards, and the last close-out ended clean: markdown lint, `flai check --strict` with no findings, and flaiover's lint, types, and 741 tests.
