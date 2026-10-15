import { afterAll, afterEach, beforeAll, describe, expect, it } from 'vitest';
import { cp, mkdtemp, rm } from 'node:fs/promises';
import { tmpdir } from 'node:os';
import { join, resolve } from 'node:path';
import { AgentError } from '$lib/server/agent';
import { Repo, useRepo, type Ask, type DesignProposal } from '$lib/server/repo';
import { flaiAsk, haveFlai } from '$lib/server/testing';

// GET /api/design/proposals answers flai's design.proposals on the host (S-0384), through the flai
// built from this tree, against a copy of the fixture, whose open stories have no branch.
const fixture = resolve('../flai/internal/metrics/testdata/good');

type Handler = (event: never) => Promise<Response>;

describe.skipIf(!haveFlai)('design proposals endpoint', () => {
	let dir: string;
	let GET: Handler;
	const get = () => GET({ url: new URL('http://x/api/design/proposals') } as never);

	beforeAll(async () => {
		dir = await mkdtemp(join(tmpdir(), 'flaiover-design-proposals-'));
		await cp(fixture, dir, { recursive: true });
		GET = (await import('./+server')).GET as unknown as Handler;
	});
	afterEach(() => useRepo(null));
	afterAll(async () => {
		await rm(dir, { recursive: true, force: true });
	});

	it('answers each proposing story with its ID, title, state, and paths, asked afresh each time', async () => {
		const proposals: DesignProposal[] = [
			{
				id: 'S-004',
				title: 'Four',
				status: 'in-progress',
				paths: ['design/system/design-system.tokens.json', 'design/system/design-system.md']
			}
		];
		let asked = 0;
		const flai = flaiAsk(dir);
		const proposing = (async (method: string, params?: Record<string, unknown>) => {
			if (method !== 'design.proposals') return flai(method, params);
			asked++;
			return proposals;
		}) as unknown as Ask;
		useRepo(new Repo(dir, proposing));
		const r = await get();
		expect(r.status).toBe(200);
		expect(await r.json()).toEqual(proposals);
		await get();
		expect(asked).toBe(2);
	});

	it('answers [] when no open story changes the design system', async () => {
		useRepo(new Repo(dir, flaiAsk(dir)));
		const r = await get();
		expect(r.status).toBe(200);
		expect(await r.json()).toEqual([]);
	});

	it('answers 501 when the host flai is older than design.proposals', async () => {
		const flai = flaiAsk(dir);
		const older = (async (method: string, params?: Record<string, unknown>) => {
			if (method === 'design.proposals')
				throw new AgentError(502, 'flai offers no method design.proposals', -32601);
			return flai(method, params);
		}) as unknown as Ask;
		useRepo(new Repo(dir, older));
		const r = await get();
		expect(r.status).toBe(501);
		expect((await r.json()).error).toContain('design.proposals');
	});
});
