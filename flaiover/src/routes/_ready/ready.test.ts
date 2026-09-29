// S-0162: readiness asks flai for the manifest alone, not the whole archive with its bodies.
import { afterEach, describe, expect, it, vi } from 'vitest';
import type { Ask } from '$lib/server/repo';

const project = { name: 'Harbour', key: 'harbour', layout: { design: 'd', docs: 'o', wip: 'w' } };

describe('GET /_ready', () => {
	afterEach(() => {
		vi.resetModules();
		vi.doUnmock('$lib/server/agent');
	});

	async function ready(ask: Ask) {
		vi.doMock('$lib/server/agent', async (importOriginal) => ({
			...(await importOriginal<typeof import('$lib/server/agent')>()),
			agent: () => ({ status: () => ({ connected: true, flai: '1.2.3' }) })
		}));
		const { useRepo: use, Repo: R } = await import('$lib/server/repo');
		use(new R('/nowhere', ask));
		const { GET } = await import('./+server');
		const res = await (GET as unknown as () => Promise<Response>)();
		return { status: res.status, body: await res.json() };
	}

	it('asks for the manifest and nothing else', async () => {
		const asked: string[] = [];
		const { status, body } = await ready((async (method: string) => {
			asked.push(method);
			return project;
		}) as unknown as Ask);
		expect(asked).toEqual(['project.info']);
		expect(status).toBe(200);
		expect(Object.keys(body.checks).sort()).toEqual(['host_flai', 'project']);
		expect(body.checks.project).toEqual({ ok: true, detail: 'Harbour' });
	});

	it('is not ready when flai cannot give the manifest', async () => {
		const { status, body } = await ready((async () => {
			throw new Error('no answer');
		}) as unknown as Ask);
		expect(status).toBe(503);
		expect(body.checks.project.ok).toBe(false);
	});
});
