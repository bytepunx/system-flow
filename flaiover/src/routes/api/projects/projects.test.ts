import { afterEach, describe, expect, it } from 'vitest';
import { GET } from './+server';
import { registry, resetAgent } from '$lib/server/agent';

// S-0080: what a switcher needs, from the server: every project a flai has ever named here.
describe('/api/projects', () => {
	afterEach(() => resetAgent());

	it('is empty when nothing has connected', async () => {
		const body = await (await GET({} as never)).json();
		expect(body.projects).toEqual([]);
	});

	it('lists every project the registry knows, connected or not', async () => {
		registry().hub('harbour');
		const body = await (await GET({} as never)).json();
		expect(body.projects.some((p: { key: string }) => p.key === 'harbour')).toBe(true);
	});
});
