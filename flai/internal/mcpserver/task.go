package mcpserver

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/bytepunx/system-flow/flai/internal/taskdone"
)

const taskDoneDescription = "Close a task in one call from its story's worktree (ADR-0107), with the answer flai task done --json prints. In order: commit everything in the worktree with message (nothing to commit is no failure), message each other story in progress or in review whose claim covers a path the commit changed, the shared paths included, naming the task, the commit, its subject, and those paths, and answer them as told (a message that cannot be sent is logged and stops nothing), sync the story's branch as flai stream sync does, move the task to done with what moves with it, log message's subject line, or log, to the story's narrative, add the paths the commit changed to the task's and the story's touches where they do not cover them, run flai check --strict scoped to the story, and read your inbox, as inbox does. It stops at the first step that fails, and answers stopped naming it, error saying why and what to do, and what the steps before it did; a stop is an answer, not a tool error. After a stopped sync, resolve each conflicting path in the worktree, git add it, run git rebase --continue there, and call again; after a failed check, fix its findings and call again, which commits the fix. A task already done is not moved again. Running the tests and ticking criteria stay yours."

// taskDoneInstructions sends the agent to task_done at every task
// transition (ADR-0107).
const taskDoneInstructions = "At every task transition, close the task with task_done, in place of committing, syncing, moving the task, logging, widening touches, checking, and calling inbox one by one: it does all seven in that order, stops at the first that fails with what to do, and is called again once that is done."

// TaskDoneIn closes a task (ADR-0107).
type TaskDoneIn struct {
	Project string `json:"project,omitempty" jsonschema:"the project, by key or folder: needed only when the server serves more than one"`
	Task    string `json:"task" jsonschema:"the task to close, T-nnnn, in any zero padding"`
	Message string `json:"message" jsonschema:"the commit message; its subject line is the narrative's log entry unless log is given"`
	Log     string `json:"log,omitempty" jsonschema:"the narrative's log entry, when it is not to be the message's subject line"`
}

func (in TaskDoneIn) project() string { return in.Project }

// TaskDoneOut is taskdone's result itself, so that task_done answers field
// for field what flai task done --json prints.
type TaskDoneOut = taskdone.Result

// taskDone closes the task as the server's agent. A step that fails is the
// answer's stopped and error; only a close that cannot start is the tool's
// error.
func (s *server) taskDone(ctx context.Context, _ *mcp.CallToolRequest, in TaskDoneIn) (*mcp.CallToolResult, TaskDoneOut, error) {
	res, err := taskdone.Run(ctx, taskdone.Options{
		Repo: s.repo, Runner: s.runner, Task: in.Task, Message: in.Message, LogEntry: in.Log,
		Agent: s.agent, Session: s.session, Now: s.now, Version: s.version, Logger: s.logger,
	})
	if err != nil {
		return nil, TaskDoneOut{}, err
	}
	return nil, res, nil
}
