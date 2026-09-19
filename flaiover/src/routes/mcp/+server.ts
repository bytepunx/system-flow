import { json } from '@sveltejs/kit';
import { repo, RepoError } from '$lib/server/repo';
import { getSession, startSession } from '$lib/server/mcpbridge';
import type { RequestHandler } from './$types';

// MCP over Streamable HTTP (ADR-0024): the same server as `flai mcp`, one flai
// mcp process per session. Bearer token only; see decide() in auth.ts.

const SESSION_HEADER = 'mcp-session-id';
type Msg = Record<string, unknown>;

const rpcError = (status: number, message: string, id: unknown = null) =>
	json({ jsonrpc: '2.0', id, error: { code: -32000, message } }, { status });

const isRequest = (m: Msg) => typeof m.method === 'string' && m.id !== undefined && m.id !== null;

export const POST: RequestHandler = async ({ request }) => {
	let body: unknown;
	try {
		body = await request.json();
	} catch {
		return rpcError(400, 'the body is not JSON');
	}
	const messages = (Array.isArray(body) ? body : [body]) as Msg[];
	if (!messages.length || messages.some((m) => m === null || typeof m !== 'object'))
		return rpcError(400, 'the body is not a JSON-RPC message or batch');

	const sid = request.headers.get(SESSION_HEADER);
	let session = sid ? getSession(sid) : undefined;
	const headers: Record<string, string> = {};
	if (sid && !session)
		return rpcError(404, 'unknown or expired MCP session; initialize again', messages[0].id);
	if (!session) {
		const init = messages.find((m) => m.method === 'initialize');
		if (!init)
			return rpcError(
				400,
				'no Mcp-Session-Id header: start with an initialize request',
				messages[0].id
			);
		const client = (init.params as { clientInfo?: { name?: unknown } } | undefined)?.clientInfo
			?.name;
		const agent =
			request.headers.get('x-flai-agent') || (typeof client === 'string' ? client : '') || 'agent';
		try {
			session = await startSession(repo().root, agent);
		} catch (e) {
			if (e instanceof RepoError) return rpcError(e.status, e.message, init.id);
			throw e;
		}
		headers[SESSION_HEADER] = session.id;
	}

	try {
		const replies = (await Promise.all(messages.map((m) => session!.send(m)))).filter(
			(r): r is Msg => r !== null
		);
		// nothing to answer: notifications and responses are acknowledged
		if (!messages.some(isRequest)) return new Response(null, { status: 202, headers });
		return json(Array.isArray(body) ? replies : replies[0], { headers });
	} catch (e) {
		return rpcError(
			502,
			`flai mcp did not answer: ${e instanceof Error ? e.message : String(e)}`,
			messages[0].id
		);
	}
};

/** No server-initiated stream is offered; the transport allows 405 here. */
export const GET: RequestHandler = () =>
	new Response(null, { status: 405, headers: { allow: 'POST, DELETE' } });

export const DELETE: RequestHandler = ({ request }) => {
	const sid = request.headers.get(SESSION_HEADER);
	const session = sid ? getSession(sid) : undefined;
	if (!session) return rpcError(404, 'unknown or expired MCP session');
	session.close();
	return new Response(null, { status: 204 });
};
