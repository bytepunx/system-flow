---
id: D-0395
title: Svelte's lint reads an assignment to a bindable prop as a useless assignment, so a child component reports a value to its parent through a callback prop
kind: decision
sources: [S-0129]
topics: [sveltekit, flaiover, code]
paths: [flaiover/src/lib/components/StoryAgent.svelte]
date: 2026-10-28
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0395 Svelte's lint reads an assignment to a bindable prop as a useless assignment, so a child component reports a value to its parent through a callback prop

S-0129: "StoryAgent tells the page the hold through an `onhold` callback rather than a bindable prop, which the lint reads as a useless assignment." A new child component that reports a value back to its page takes a callback prop the same way.

Source: S-0129.
