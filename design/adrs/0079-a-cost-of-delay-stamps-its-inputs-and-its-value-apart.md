---
id: ADR-0079
title: A cost of delay stamps its inputs and its value apart
status: accepted
date: 2026-10-03
supersedes: []
superseded_by: []
refines: [ADR-0074]
---

# ADR-0079 A cost of delay stamps its inputs and its value apart

## Context

ADR-0074 gives `cost_of_delay` one `by` and `at`, the block's last editor, and says the keys tell whose a figure is: the inputs are the operator's, the value the planner's. S-0204 shows the planner's value on the item page with who set it and when, and notes when the value is stale because the inputs changed after it. With one stamp, the operator's saving of an input overwrites the planner's stamp, so neither the value's author nor whether the inputs are newer can be told. The question is TH-0090.

## Decision

A cost of delay stamps its inputs and its value apart.

- `cost_of_delay.inputs` carry their own `by` and `at`, set when an input is added, changed, or removed and the inputs remain. A change that leaves no input removes the inputs and their stamp.
- The block's `by` and `at` are the value's, set when `value` is added or changed; when the value is removed, they go with it.
- `flai edit`, MCP `item_edit`, and hostapi `item.edit` stamp each with the editor, as before; `flai story new` and `flai epic new` stamp the inputs with the owner, and a story made from an issue stamps them `flai`.
- The value is stale when there is a value and inputs, and the inputs' `at` is later than the block's.
- Validation requires `inputs.by` and `inputs.at` on inputs, and the block's `by` and `at` on a value; neither stamp without what it stamps.
- A block written with one stamp (ADR-0074) is read as: the stamp is the inputs' when there is no value, and when there is a value the stamp is the value's and is given to the inputs too, so it reads as not stale. flai writes it back in the new shape on its next edit.

## Consequences

- The item page can show who set the value and when, and say when the planner should recompute it.
- Clearing every input leaves the value with its stamp and no stale note: nothing newer is left to compare.
- `inputs.by` and `inputs.at` are new keys in `front-matter-fields.txt`, so the flai release that carries them raises `flai.minimum`.

## Alternatives considered

- Keep one stamp and call the value stale whenever the block's `by` is not the planner: wrong whenever anyone else sets a value, and the planner's stamp is still lost.
- Per-key stamps on every input: more than any reader needs; the inputs are set together by one person.
