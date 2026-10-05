---
id: T-0978
type: task
nature: improvement
title: The docs and the design say flai touches replaces unless given --add or --remove
status: backlog
parent: S-0254
owner: alex
created: 2026-10-05T05:49:41Z
updated: 2026-10-05T05:49:41Z
transitions: []
stream: S-0254
tags: [docs]
touches: [docs/users/flai.md, docs/users/flai-reference.md, docs/operators/settings.md, design/system/flai-cli.md]
after: [T-0976]
---
# T-0978 The docs and the design say flai touches replaces unless given --add or --remove

## Work

Say what T-0976 built wherever `flai touches` is described:

- `docs/users/flai.md` § Touches: an example with `--add`, and one sentence each saying that paths given alone replace the list, that `--add` adds, and that `--remove` removes. § Outside the touches: the example sync output shows the new `flai touches <id> --add <paths>` hint.
- `design/system/flai-cli.md`: the `flai touches` row shows `[--add | --remove]` and says that paths given alone replace the list (S-0254, I-0067).
- `docs/users/flai-reference.md` and the flag index in `docs/operators/settings.md`: regenerate both with `scripts/flai-reference.sh`. Never edit them by hand.

Waits for T-0976, because the reference is generated from its help and the docs describe its flags and its sync hint as built.

## Done when

- `docs/users/flai.md` and `design/system/flai-cli.md` say that paths given alone replace the list, and show `--add` and `--remove`
- `scripts/flai-reference.sh` leaves no diff after it is run again, and `settings.md` lists `--add` and `--remove` under `flai touches`
- The markdown lint passes

## Notes
