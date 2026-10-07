---
id: T-0920
type: task
nature: improvement
title: flai guard reads a heredoc's text as input, not as commands, unless a shell runs it
status: done
parent: S-0246
owner: alex
created: 2026-10-05T05:44:52Z
updated: 2026-10-07T09:03:06Z
transitions:
  - to: ready
    at: 2026-10-07T09:00:25Z
    by: agent-S-0246
  - to: in-progress
    at: 2026-10-07T09:00:25Z
    by: agent-S-0246
  - to: done
    at: 2026-10-07T09:03:06Z
    by: agent-S-0246
stream: S-0246
tags: [flai]
touches: [flai/internal/guard/guard.go, flai/internal/guard/guard_test.go, flai/cmd/guard_test.go]
usage:
  source: log
  seconds: 161
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 10
      output: 2613
      cache_read: 556041
      cache_write: 11658
      cost: 0.2568
---
# T-0920 flai guard reads a heredoc's text as input, not as commands, unless a shell runs it

## Work

`commands` in `flai/internal/guard/guard.go` splits a command line at newlines and reads every line as a command, including the lines of a heredoc. A heredoc's text is the input of the program before it, not a command. Both instances on I-0058 were python heredocs that edited Go code. In S-0209's T-0791, an apostrophe in the text opened a quote that never closed, and the words `flai story new` were read as a flai write. The fix is in `commands`, which `Guard.Check` and `Guard.plan` share, so a sub-agent and the planner both get it.

- Recognise each heredoc redirection on a line: `<<WORD` and `<<-WORD`, with the delimiter bare, quoted (`'EOF'`, `"EOF"`), or escaped (`\EOF`), and more than one on a line. Leave `<<<`, a here-string, as it is. From the end of that line, take each body as text up to the line that is exactly its delimiter. With `<<-`, strip leading tabs first. A body with no closing line runs to the end of the input.
- Keep the guard where a heredoc does run commands:
  - A body fed to a shell (`sh`, `bash`, `zsh`, `dash`, `ksh` with no `-c` script, `-s` among its flags allowed) is a script. Check it as `script` checks a `-c` script, so `bash <<EOF` followed by `flai move S-0001 review` is still refused.
  - In a body whose delimiter is not quoted, the shell runs command substitutions, `$(...)` and backticks. Check those.
- Tests in `guard_test.go`:
  - Both I-0058 instances pass for a sub-agent and for the planner: a `python3 - <<'EOF'` body with an apostrophe and the words `flai story new`, and a heredoc whose text is a Go comment.
  - A heredoc fed to `bash` with a flai write is refused, and so is a `$(flai move …)` in an unquoted body.
  - A command after the closing delimiter is still checked.
  - Cover `<<-` with tabs, a quoted and an escaped delimiter, two heredocs on one line, and an unclosed body.
- Add one test in `flai/cmd/guard_test.go` that feeds S-0209's instance as hook input to `flai guard` and expects exit 0, so the multi-line command is shown to survive the JSON.

This task waits for none. The documentation task waits for it.

## Done when

- `commands` takes heredoc bodies out of the command line, and checks a body as a script when a shell reads it, along with command substitutions in an unquoted body.
- The tests above pass, and every existing test in `guard_test.go` and `flai/cmd/guard_test.go` passes unchanged.
- `scripts/flai-test.sh` passes.

## Notes
