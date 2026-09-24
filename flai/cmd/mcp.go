package cmd

import (
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"

	"github.com/bytepunx/system-flow/flai/internal/buildinfo"
	"github.com/bytepunx/system-flow/flai/internal/mcpserver"
)

func newMCPCmd(a *app) *cobra.Command {
	var agentFlag string
	c := &cobra.Command{
		Use:   "mcp",
		Short: "Serve this repository to agents over the Model Context Protocol: on stdio, or over HTTP with flai mcp start",
		Long: `An MCP server for agent sessions (ADR-0020): the inbox of threads awaiting
the agent, replies, work item and document reads, transitions with the
workflow rules, who is touching a path, and wait_for_events, which blocks
until the designer or another agent changes a thread, an item, or a
narrative. Every write goes through the same code as the CLI. The agent is
FLAI_AGENT. Register it in .mcp.json:

  { "mcpServers": { "flai": { "command": "flai", "args": ["mcp"] } } }

Standard output is the protocol channel; log events go to standard error.

Started in a folder that is not a project, ~/git say, it serves every
system-flow project in the folder and below it (S-0101): inbox,
wait_for_work, and wait_for_events cover them all, and the other tools take
the project's key.

An agent that cannot start a process here reaches the same server over
Streamable HTTP (ADR-0030): flai mcp start runs it in the background for
this project, flai mcp status says where it listens and what an agent's
configuration looks like, and flai mcp token prints its bearer token.`,
		Example: `  flai mcp              # stdio, as .mcp.json starts it
  flai mcp start        # over HTTP, in the background, on this machine only
  flai mcp status
  flai mcp token
  flai mcp stop
  flai mcp http         # over HTTP in the foreground; Ctrl-C stops it`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			repo, err := a.projectOrNone()
			if err != nil {
				return err
			}
			agent := mcpAgent(agentFlag)
			opt := mcpserver.Options{Repo: repo, Agent: agent, Version: buildinfo.Version, Now: a.now, Runner: a.runner}
			if repo == nil {
				// Not in a project (S-0101): every project in this folder and below it.
				if opt.Folder, err = a.workingDir(); err != nil {
					return err
				}
				a.logger().Info("mcp server started", "component", "mcp", "agent", agent, "folder", opt.Folder)
			} else {
				a.logger().Info("mcp server started", "component", "mcp", "agent", agent, "root", repo.Root)
			}
			return mcpserver.New(opt).Run(cmd.Context(), &mcp.StdioTransport{})
		},
	}
	c.Flags().StringVar(&agentFlag, "agent", "", "the agent this server serves, over FLAI_AGENT; flai serve names the agent it started this way")
	c.AddCommand(newMCPHTTPCmds(a)...)
	return c
}

// mcpAgent is who the server serves: --agent when given, else FLAI_AGENT.
// flai serve passes the name as an argument because a harness's own settings
// can replace the environment of the servers it starts, which made every
// agent it started one name (S-0114, I-0037).
func mcpAgent(flag string) string {
	if flag != "" {
		return flag
	}
	agent, _ := agentIdentity()
	return agent
}
