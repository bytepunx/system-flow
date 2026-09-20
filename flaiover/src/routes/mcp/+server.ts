import { json } from '@sveltejs/kit';
import type { RequestHandler } from './$types';

// MCP was served here until flaiover 0.22 (ADR-0024). It is flai's now, on the host, where the
// agents are (ADR-0030); this answer is what is left, so a client configured for the old address
// is told where to go instead of getting a 404. It needs no token: it says nothing about the project.

const MOVED =
	'MCP is no longer served by the dashboard. On the host, in the project: an agent there starts `flai mcp` itself (stdio, as .mcp.json does); any other reaches it over HTTP after `flai mcp start`, and `flai mcp status` prints the address and the configuration (ADR-0030).';

const gone: RequestHandler = () =>
	json({ jsonrpc: '2.0', id: null, error: { code: -32000, message: MOVED } }, { status: 410 });

export const GET = gone;
export const POST = gone;
export const DELETE = gone;
