package cmd

import (
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"

	"github.com/bytepunx/system-flow/flai/internal/buildinfo"
	"github.com/bytepunx/system-flow/flai/internal/mcpserver"
)

func newMCPCmd(a *app) *cobra.Command {
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
			repo, err := a.project()
			if err != nil {
				return err
			}
			agent, _ := agentIdentity()
			srv := mcpserver.New(mcpserver.Options{Repo: repo, Agent: agent, Version: buildinfo.Version, Now: a.now, Runner: a.runner})
			a.logger().Info("mcp server started", "component", "mcp", "agent", agent, "root", repo.Root)
			return srv.Run(cmd.Context(), &mcp.StdioTransport{})
		},
	}
	c.AddCommand(newMCPHTTPCmds(a)...)
	return c
}
