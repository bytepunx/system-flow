import { afterAll, afterEach, beforeAll, describe, expect, it } from 'vitest';
import { cp, mkdtemp, rm, writeFile } from 'node:fs/promises';
import { tmpdir } from 'node:os';
import { join, resolve } from 'node:path';
import { AgentError } from '$lib/server/agent';
import { Repo, useRepo, type Ask, type DesignSystem } from '$lib/server/repo';
import { flaiAsk, haveFlai, shell } from '$lib/server/testing';

// GET /api/design answers flai's design.get on the host (S-0382), of a story's worktree when asked
// (S-0384), through the flai built from this tree, against copies of the fixture: one without a
// design system, and one flai design init wrote.
const fixture = resolve('../flai/internal/metrics/testdata/good');
const TOKENS = 'design/system/design-system.tokens.json';

type Handler = (event: never) => Promise<Response>;

describe.skipIf(!haveFlai)('design system endpoint', () => {
	let bare: string;
	let designed: string;
	let GET: Handler;
	const get = (query = '') => GET({ url: new URL(`http://x/api/design${query}`) } as never);

	beforeAll(async () => {
		bare = await mkdtemp(join(tmpdir(), 'flaiover-design-bare-'));
		designed = await mkdtemp(join(tmpdir(), 'flaiover-design-'));
		await cp(fixture, bare, { recursive: true });
		await cp(fixture, designed, { recursive: true });
		await shell(designed, ['design', 'init']);
		GET = (await import('./+server')).GET as unknown as Handler;
	});
	afterEach(() => useRepo(null));
	afterAll(async () => {
		await rm(bare, { recursive: true, force: true });
		await rm(designed, { recursive: true, force: true });
	});

	it('answers the tokens per theme, the guidance, the renderings, and the contrast', async () => {
		useRepo(new Repo(designed, flaiAsk(designed)));
		const r = await get();
		expect(r.status).toBe(200);
		const body = (await r.json()) as DesignSystem;
		expect(body.path).toBe(TOKENS);
		expect(body.exists).toBe(true);
		expect(body.no_token_file).toBeUndefined();
		expect(body.guidance).toBe('design/system/design-system.md');
		expect(body.headings).toContain('Colour');
		expect(body.themes.map((t) => t.name)).toEqual(['light', 'dark']);
		const primary = body.themes[0].tokens.find((t) => t.path === 'color.primary');
		expect(primary?.type).toBe('color');
		expect(primary?.variables[0].name).toBe('--color-primary');
		expect(body.problems).toEqual([]);
		expect(body.css).toContain('--color-primary');
		expect(body.tailwind).toContain('@theme');
		expect(body.contrast.map((c) => c.theme)).toEqual(['light', 'dark']);
		expect(body.contrast[0].pairs.every((p) => p.passes)).toBe(true);
	});

	it('answers a token file problem by its line', async () => {
		const broken = await mkdtemp(join(tmpdir(), 'flaiover-design-broken-'));
		try {
			await cp(designed, broken, { recursive: true });
			await writeFile(join(broken, TOKENS), '{\n  "color": {\n');
			useRepo(new Repo(broken, flaiAsk(broken)));
			const r = await get();
			expect(r.status).toBe(200);
			const body = (await r.json()) as DesignSystem;
			expect(body.exists).toBe(true);
			expect(body.problems.length).toBeGreaterThan(0);
			expect(body.problems[0].line).toBeGreaterThan(0);
			expect(body.themes).toEqual([]);
		} finally {
			await rm(broken, { recursive: true, force: true });
		}
	});

	it('answers that there is no token file, naming the command that writes one', async () => {
		useRepo(new Repo(bare, flaiAsk(bare)));
		const r = await get();
		expect(r.status).toBe(200);
		const body = (await r.json()) as DesignSystem;
		expect(body.path).toBe(TOKENS);
		expect(body.exists).toBe(false);
		expect(body.no_token_file).toContain('flai design init');
		expect(body.themes).toEqual([]);
		expect(body.problems).toEqual([]);
		expect(body.contrast).toEqual([]);
	});

	it('keeps the answer until a document changes', async () => {
		let asked = 0;
		const flai = flaiAsk(designed);
		const counting = (async (method: string, params?: Record<string, unknown>) => {
			if (method === 'design.get') asked++;
			return flai(method, params);
		}) as unknown as Ask;
		const r = new Repo(designed, counting);
		useRepo(r);
		await r.manifest();
		await get();
		await get();
		expect(asked).toBe(1);
		r.changed('wip/kanban/stories/S-0001-x.md');
		await get();
		expect(asked).toBe(1);
		r.changed(TOKENS);
		await get();
		expect(asked).toBe(2);
	});

	it("passes the story to design.get and asks afresh each time, keeping the main checkout's", async () => {
		const asked: (Record<string, unknown> | undefined)[] = [];
		const flai = flaiAsk(designed);
		const counting = (async (method: string, params?: Record<string, unknown>) => {
			if (method !== 'design.get') return flai(method, params);
			asked.push(params);
			// the designed copy stands in for the story's worktree
			return flai(method, {});
		}) as unknown as Ask;
		useRepo(new Repo(designed, counting));
		await get();
		const r = await get('?story=S-004');
		expect(r.status).toBe(200);
		expect(((await r.json()) as DesignSystem).path).toBe(TOKENS);
		await get('?story=S-004');
		await get();
		expect(asked).toEqual([{}, { story: 'S-004' }, { story: 'S-004' }]);
	});

	it("answers flai's refusal of a story with its status and message", async () => {
		useRepo(new Repo(bare, flaiAsk(bare)));
		const missing = await get('?story=S-999');
		expect(missing.status).toBe(404);
		expect((await missing.json()).error).toContain('S-999 not found');
		const malformed = await get('?story=nope');
		expect(malformed.status).toBe(400);
		expect((await malformed.json()).error).toContain('not an item ID');
		const epic = await get('?story=E-001');
		expect(epic.status).toBe(400);
		expect((await epic.json()).error).toContain('name a story');
		const unopened = await get('?story=S-004');
		expect(unopened.status).toBe(400);
		expect((await unopened.json()).error).toContain('no worktree');
	});

	it("answers flai's refusal with its status and message", async () => {
		const flai = flaiAsk(bare);
		const refusing = (async (method: string, params?: Record<string, unknown>) => {
			if (method === 'design.get')
				throw new AgentError(
					502,
					`read the design token file ${TOKENS}: permission denied`,
					-32603
				);
			return flai(method, params);
		}) as unknown as Ask;
		useRepo(new Repo(bare, refusing));
		const r = await get();
		expect(r.status).toBe(500);
		expect((await r.json()).error).toContain('permission denied');
	});

	it('answers 501 when the host flai is older than design.get', async () => {
		const flai = flaiAsk(bare);
		const older = (async (method: string, params?: Record<string, unknown>) => {
			if (method === 'design.get')
				throw new AgentError(502, 'flai offers no method design.get', -32601);
			return flai(method, params);
		}) as unknown as Ask;
		useRepo(new Repo(bare, older));
		const r = await get();
		expect(r.status).toBe(501);
		expect((await r.json()).error).toContain('design.get');
	});
});
