---
id: D-0447
title: storygit cannot import issues, since issues imports storygit, so storygit.NewTrialJudge takes a storygit.Files the caller builds with issues.FilesIn, and workitem takes the judge as a plain func type
kind: decision
sources: [S-0551]
topics: [flai, go, code, git]
paths: [flai/internal/storygit/trial.go, flai/internal/workitem/hold.go, flai/internal/issues/generated.go, flai/internal/serve/agents.go, flai/internal/mcpserver/work.go, flai/cmd/project.go]
date: 2026-11-01
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0447 storygit cannot import issues, since issues imports storygit, so storygit.NewTrialJudge takes a storygit.Files the caller builds with issues.FilesIn, and workitem takes the judge as a plain func type

S-0551 built the trial-merge judge that ADR-0230 gives the holds. Its narrative records the seam: "`storygit.NewTrialJudge(r execx.Runner, repo *workitem.Repo, files storygit.Files)` builds the judge, in `trial.go`. It takes `files` because the code that builds a `Files` (`issues.FilesIn`) imports storygit. Callers pass `issues.FilesIn(repo, repo.MainRoot, r, time.Now)`."

The judge reaches `workitem` as `workitem.TrialJudge func(held, holder string) ([]string, error)`, set with `Holds.WithTrial(j)` or the `Repo.Trial` field, since `storygit` imports `workitem`. So the three callers that wire it, `flai serve`'s look (`flai/internal/serve/agents.go`), the MCP server (`flai/internal/mcpserver/work.go`), and the CLI's project open (`flai/cmd/project.go`), each build the `Files` themselves.

A review fix in the same story made `NewTrialJudge` return no judge without a runner or outside a git checkout, so a hold there gives the verdict it gave before the trial merge, as ADR-0230 says of no judge.

A later change that has `storygit` need what `issues` computes should take it as an argument, as this does, rather than import `issues`.

Source: S-0551 (wip/archive/agents/S-0551.md, Decisions).
