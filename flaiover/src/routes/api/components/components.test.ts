import { afterAll, beforeAll, beforeEach, describe, expect, it } from 'vitest';
import { AgentError } from '$lib/server/agent';
import {
	Repo,
	useRepo,
	type Ask,
	type ComponentsCatalogue,
	type ComponentsStatus
} from '$lib/server/repo';

// /api/components, /api/components/start, and /api/components/stop ask flai's components.list,
// components.status, components.start, and components.stop (S-0439) of a stubbed host: running a
// real dev server needs a Svelte sub-project's node_modules, which the flai side's tests cover.
type Handler = (event: never) => Promise<Response>;

const CATALOGUE: ComponentsCatalogue = {
	projects: [
		{
			name: 'web',
			path: 'web',
			kind: 'sveltekit',
			supported: true,
			components: [
				{
					project: 'web',
					path: 'src/lib/Button.svelte',
					name: 'Button',
					props: [
						{
							name: 'variant',
							kind: 'union',
							type: "'primary' | 'ghost'",
							required: false,
							default: "'primary'",
							literals: ['"primary"', '"ghost"']
						}
					]
				}
			]
		},
		{ name: 'api', path: 'api', kind: 'go', supported: false, components: [] }
	]
};

const status = (running: boolean): ComponentsStatus => ({
	enabled: true,
	enable: 'flai serve enable components',
	idle_seconds: 1800,
	servers: [
		running
			? {
					project: 'web',
					running: true,
					url: 'http://127.0.0.1:41000/__flai/components/',
					port: 41000,
					pid: 4242,
					since: '2026-10-16T09:00:00Z'
				}
			: { project: 'web', running: false }
	]
});

let asked: { method: string; params: Record<string, unknown> }[] = [];
let script: Record<string, { data?: unknown; error?: unknown }> = {};
// the reads answer raw, as flai's do; the two writes as commands do
const WRITES = new Set(['components.start', 'components.stop']);
const host = (async (method: string, params: Record<string, unknown> = {}) => {
	asked.push({ method, params });
	const s = script[method];
	if (!s) throw new AgentError(502, `flai offers no method ${method}`, -32601);
	if (s.error) throw s.error;
	return WRITES.has(method) ? { data: s.data, warnings: [] } : s.data;
}) as Ask;

const disabled = (method: string) =>
	new AgentError(
		502,
		`the host action "components" is not enabled for p. On the host, in the project, run: flai serve enable components (${method})`,
		-32012,
		{ action: 'components', enable: 'flai serve enable components' }
	);

describe('/api/components (S-0439)', () => {
	let GET: Handler;
	let START: Handler;
	let STOP: Handler;
	const get = () => GET({ url: new URL('http://x/api/components') } as never);
	const post = (fn: Handler, body?: unknown) =>
		fn({
			request: new Request('http://x/api/components/x', {
				method: 'POST',
				...(body === undefined ? {} : { body: JSON.stringify(body) })
			})
		} as never);
	const asks = (method: string) => asked.filter((a) => a.method === method);

	beforeAll(async () => {
		GET = (await import('./+server')).GET as unknown as Handler;
		START = (await import('./start/+server')).POST as unknown as Handler;
		STOP = (await import('./stop/+server')).POST as unknown as Handler;
	});
	beforeEach(() => {
		asked = [];
		script = {
			'project.info': { data: { name: 'p', key: 'p' } },
			'components.list': { data: CATALOGUE },
			'components.status': { data: status(false) }
		};
		useRepo(new Repo('/nowhere', host));
	});
	afterAll(() => useRepo(null));

	it('answers the catalogue and the status together, asking flai afresh each time', async () => {
		const r = await get();
		expect(r.status).toBe(200);
		expect(await r.json()).toMatchObject({ catalogue: CATALOGUE, status: status(false) });
		const [web, api] = CATALOGUE.projects;
		const card = { project: 'web', path: 'src/lib/Card.svelte', name: 'Card', props: [] };
		script['components.list'] = {
			data: { projects: [{ ...web, components: [...web.components, card] }, api] }
		};
		const again = (await (await get()).json()) as { catalogue: ComponentsCatalogue };
		expect(again.catalogue.projects[0].components.map((c) => c.name)).toEqual(['Button', 'Card']);
		expect(asks('components.list')).toHaveLength(2);
		expect(asks('components.status')).toHaveLength(2);
	});

	it('starts the named dev server with a request ID, then answers every server status', async () => {
		script['components.start'] = { data: status(true).servers[0] };
		script['components.status'] = { data: status(true) };
		const r = await post(START, { project: 'web' });
		expect(r.status).toBe(200);
		expect(await r.json()).toMatchObject(status(true));
		const params = asks('components.start')[0].params;
		expect(params.project).toBe('web');
		expect(typeof params.request_id).toBe('string');
		expect(asks('components.status')).toEqual([{ method: 'components.status', params: {} }]);
	});

	it('leaves the sub-project to flai when the body names none', async () => {
		script['components.start'] = { data: status(true).servers[0] };
		expect((await post(START)).status).toBe(200);
		expect((await post(START, { project: ' ' })).status).toBe(200);
		for (const a of asks('components.start')) expect(a.params.project).toBeUndefined();
	});

	it('refuses a project that is not a string with 400, before asking flai', async () => {
		const r = await post(START, { project: 7 });
		expect(r.status).toBe(400);
		expect((await r.json()).error).toContain('sub-project');
		expect(asks('components.start')).toEqual([]);
	});

	it("answers a failed start with flai's message, Vite's output with it, as 500", async () => {
		const message =
			"sub-project web's component dev server exited before it answered: vite: Error: Cannot find module 'svelte'";
		script['components.start'] = { error: new AgentError(502, message, -32603) };
		const r = await post(START, { project: 'web' });
		expect(r.status).toBe(500);
		expect((await r.json()).error).toBe(message);
	});

	it('answers a sub-project without vite, refused by rule, as 400', async () => {
		const message = 'sub-project web has no vite in its node_modules: run npm install in web';
		script['components.start'] = { error: new AgentError(502, message, -32011) };
		const r = await post(START, { project: 'web' });
		expect(r.status).toBe(400);
		expect((await r.json()).error).toBe(message);
	});

	it('answers 403 with the command that enables the action while it is off', async () => {
		script['components.start'] = { error: disabled('start') };
		script['components.stop'] = { error: disabled('stop') };
		for (const fn of [START, STOP]) {
			const r = await post(fn, { project: 'web' });
			expect(r.status).toBe(403);
			const body = await r.json();
			expect(body).toMatchObject({ action: 'components', enable: 'flai serve enable components' });
			expect(body.error).toContain('flai serve enable components');
		}
		expect(asks('components.status')).toEqual([]);
	});

	it('stops the named dev server with a request ID, then answers every server status', async () => {
		script['components.stop'] = { data: { project: 'web', stopped: true, pid: 4242 } };
		const r = await post(STOP, { project: 'web' });
		expect(r.status).toBe(200);
		expect(await r.json()).toMatchObject(status(false));
		const params = asks('components.stop')[0].params;
		expect(params.project).toBe('web');
		expect(typeof params.request_id).toBe('string');
	});

	it('answers 501 telling to upgrade a host flai without the four methods', async () => {
		script = { 'project.info': { data: { name: 'p', key: 'p' } } };
		const responses = [await get(), await post(START), await post(STOP)];
		for (const r of responses) {
			expect(r.status).toBe(501);
			const error = (await r.json()).error as string;
			expect(error).toContain('the host flai is too old');
			expect(error).toContain('flai self-upgrade');
		}
	});
});
