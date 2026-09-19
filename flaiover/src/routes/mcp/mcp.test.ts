import { afterAll, afterEach, beforeAll, describe, expect, it } from 'vitest';
import { cp, mkdtemp, readdir, rm } from 'node:fs/promises';
import { existsSync } from 'node:fs';
import { tmpdir } from 'node:os';
import { join, resolve } from 'node:path';
import { resetFlaiBinary } from '$lib/server/flai';
import { closeAllSessions, sessionCount } from '$lib/server/mcpbridge';
import { decide } from '$lib/server/auth';

const fixture = resolve('../flai/internal/metrics/testdata/good');
const bin = process.env.FLAI_BIN ?? resolve('../bin/flai');
type Handler = (event: never) => Promise<Response> | Response;

describe('who may call /mcp', () => {
	const h = (o: Record<string, string> = {}) => new Headers({ host: 'dash.example:4242', ...o });
	it('is the bearer token only, never the session cookie', () => {
		expect(decide('/mcp', 'POST', 'bearer', h()).kind).toBe('allow');
		expect(decide('/mcp', 'POST', 'cookie', h({ 'x-requested-with': 'flaiover' })).kind).toBe(
			'unauthorized'
		);
		expect(decide('/mcp', 'POST', null, h()).kind).toBe('unauthorized');
		expect(decide('/mcp', 'GET', null, h()).kind).toBe('unauthorized'); // no redirect to the login page
		expect(decide('/mcp', 'POST', 'off', h()).kind).toBe('allow');
	});
	it('refuses a request sent from another site’s page', () => {
		expect(decide('/mcp', 'POST', 'bearer', h({ origin: 'https://evil.example' })).kind).toBe(
			'forbidden'
		);
		expect(decide('/mcp', 'POST', 'bearer', h({ origin: 'http://dash.example:4242' })).kind).toBe(
			'allow'
		);
		expect(decide('/mcp', 'POST', 'bearer', h({ origin: 'null' })).kind).toBe('forbidden');
	});
});

describe.skipIf(!existsSync(bin))('MCP over HTTP, bridged to flai mcp', () => {
	let dir: string;
	let POST: Handler, GET: Handler, DELETE: Handler;
	const call = async (body: unknown, headers: Record<string, string> = {}) =>
		POST({
			request: new Request('http://x/mcp', {
				method: 'POST',
				headers: { 'content-type': 'application/json', ...headers },
				body: JSON.stringify(body)
			})
		} as never);
	const init = (name: string, id = 1) => ({
		jsonrpc: '2.0',
		id,
		method: 'initialize',
		params: { protocolVersion: '2025-06-18', capabilities: {}, clientInfo: { name, version: '0' } }
	});
	const open = async (name: string, headers: Record<string, string> = {}) => {
		const r = await call(init(name), headers);
		const sid = r.headers.get('mcp-session-id')!;
		await call({ jsonrpc: '2.0', method: 'notifications/initialized' }, { 'mcp-session-id': sid });
		return { sid, result: (await r.json()).result };
	};

	beforeAll(async () => {
		dir = await mkdtemp(join(tmpdir(), 'flaiover-mcp-'));
		await cp(fixture, dir, { recursive: true });
		process.env.PROJECT_DIR = dir;
		process.env.FLAI_BIN = bin;
		resetFlaiBinary();
		({ POST, GET, DELETE } = (await import('./+server')) as unknown as Record<string, Handler>);
	});
	afterEach(() => {
		closeAllSessions();
		delete process.env.FLAIOVER_MCP_MAX_SESSIONS;
	});
	afterAll(async () => {
		await rm(dir, { recursive: true, force: true });
	});

	it('initializes a session, lists the tools, and calls one', async () => {
		const r = await call(init('remote-agent'));
		expect(r.status).toBe(200);
		expect(r.headers.get('content-type')).toContain('application/json');
		const sid = r.headers.get('mcp-session-id')!;
		expect(sid).toMatch(/^[0-9a-f-]{36}$/);
		const first = await r.json();
		expect(first).toMatchObject({
			jsonrpc: '2.0',
			id: 1,
			result: { serverInfo: { name: 'flai' } }
		});

		const ack = await call(
			{ jsonrpc: '2.0', method: 'notifications/initialized' },
			{ 'mcp-session-id': sid }
		);
		expect(ack.status).toBe(202);

		const tools = await (
			await call({ jsonrpc: '2.0', id: 2, method: 'tools/list' }, { 'mcp-session-id': sid })
		).json();
		const names = tools.result.tools.map((t: { name: string }) => t.name);
		expect(names).toEqual(
			expect.arrayContaining(['inbox', 'board', 'item_get', 'wait_for_events'])
		);

		const got = await (
			await call(
				{
					jsonrpc: '2.0',
					id: 'a',
					method: 'tools/call',
					params: { name: 'item_get', arguments: { id: 'S-004' } }
				},
				{ 'mcp-session-id': sid }
			)
		).json();
		expect(got.id).toBe('a');
		expect(got.result.structuredContent).toMatchObject({ id: 'S-004', status: 'in-progress' });
	});

	it('gives each session its own flai mcp process, named for its agent', async () => {
		const a = await open('first-client');
		const b = await open('ignored-name', { 'x-flai-agent': 'claude@laptop' });
		expect(a.sid).not.toBe(b.sid);
		expect(sessionCount()).toBe(2);
		const who = async (sid: string) =>
			(
				await (
					await call(
						{
							jsonrpc: '2.0',
							id: 9,
							method: 'tools/call',
							params: { name: 'inbox', arguments: {} }
						},
						{ 'mcp-session-id': sid }
					)
				).json()
			).result.structuredContent.agent;
		expect(await who(a.sid)).toBe('first-client');
		expect(await who(b.sid)).toBe('claude-laptop'); // made safe for a file name
		// each agent has its own cursor, written by its own process
		expect((await readdir(join(dir, '.flai-cache', 'mcp'))).sort()).toEqual([
			'claude-laptop.json',
			'first-client.json'
		]);
	});

	it('answers a batch in order, and only the requests in it', async () => {
		const { sid } = await open('batcher');
		const r = await call(
			[
				{ jsonrpc: '2.0', id: 1, method: 'tools/list' },
				{ jsonrpc: '2.0', method: 'notifications/roots/list_changed' },
				{ jsonrpc: '2.0', id: 2, method: 'ping' }
			],
			{ 'mcp-session-id': sid }
		);
		const replies = await r.json();
		expect(replies.map((x: { id: number }) => x.id)).toEqual([1, 2]);
	});

	it('needs a session, refuses an unknown one, and forgets one that is deleted', async () => {
		expect((await call({ jsonrpc: '2.0', id: 1, method: 'tools/list' })).status).toBe(400);
		expect(
			(await call({ jsonrpc: '2.0', id: 1, method: 'tools/list' }, { 'mcp-session-id': 'nope' }))
				.status
		).toBe(404);
		const { sid } = await open('short-lived');
		const del = await DELETE({
			request: new Request('http://x/mcp', { method: 'DELETE', headers: { 'mcp-session-id': sid } })
		} as never);
		expect(del.status).toBe(204);
		expect(sessionCount()).toBe(0);
		expect(
			(await call({ jsonrpc: '2.0', id: 3, method: 'tools/list' }, { 'mcp-session-id': sid }))
				.status
		).toBe(404);
		expect((await GET({} as never)).status).toBe(405);
		expect((await call('not an object')).status).toBe(400);
	});

	it('caps the number of sessions and says how to free one', async () => {
		process.env.FLAIOVER_MCP_MAX_SESSIONS = '1';
		await open('one');
		const r = await call(init('two'));
		expect(r.status).toBe(503);
		expect((await r.json()).error.message).toContain('FLAIOVER_MCP_MAX_SESSIONS');
		expect(sessionCount()).toBe(1);
	});
});
