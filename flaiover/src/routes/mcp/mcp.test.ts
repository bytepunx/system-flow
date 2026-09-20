import { describe, expect, it } from 'vitest';
import { decide } from '$lib/server/auth';
import { DELETE, GET, POST } from './+server';

// S-0076: MCP is flai's on the host (ADR-0030). The dashboard keeps only an answer that says so.
describe('/mcp after the move', () => {
	it('says where MCP lives now, to every method, as a JSON-RPC error', async () => {
		for (const handler of [GET, POST, DELETE]) {
			const res = await (handler as (e: never) => Response | Promise<Response>)({} as never);
			expect(res.status).toBe(410);
			const body = await res.json();
			expect(body.jsonrpc).toBe('2.0');
			expect(body.error.message).toContain('flai mcp start');
			expect(body.error.message).toContain('ADR-0030');
		}
	});
	it('needs no token, because it says nothing about the project', () => {
		const h = new Headers({ host: 'dash.example:4242' });
		expect(decide('/mcp', 'POST', null, h).kind).toBe('allow');
		expect(decide('/mcp', 'GET', null, h).kind).toBe('allow');
		expect(decide('/mcp/other', 'GET', null, h).kind).toBe('login');
	});
});
