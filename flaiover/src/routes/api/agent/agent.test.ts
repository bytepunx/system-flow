// S-0098: /api/agent asks a connected project's flai for project.info, to say it answers. A
// repository offered for import has none to give, by design, and is not reported as failing.
import { afterEach, describe, expect, it, vi } from 'vitest';

describe('GET /api/agent', () => {
	afterEach(() => {
		vi.resetModules();
		vi.doUnmock('$lib/server/agent');
	});

	async function answerFor(status: Record<string, unknown>, ask = vi.fn()) {
		vi.doMock('$lib/server/agent', async (importOriginal) => ({
			...(await importOriginal<typeof import('$lib/server/agent')>()),
			agent: () => ({ status: () => status, ask })
		}));
		const { GET } = await import('./+server');
		const res = await (GET as unknown as () => Promise<Response>)();
		return { body: await res.json(), ask };
	}

	it('does not ask a candidate for project.info', async () => {
		const { body, ask } = await answerFor({
			configured: true,
			connected: true,
			candidate: true,
			missing: []
		});
		expect(ask).not.toHaveBeenCalled();
		expect(body.error).toBeUndefined();
		expect(body.candidate).toBe(true);
	});

	it('still asks a project', async () => {
		const ask = vi.fn(async () => ({ name: 'Harbour' }));
		const { body } = await answerFor({ configured: true, connected: true, missing: [] }, ask);
		expect(ask).toHaveBeenCalledWith('project.info', {}, 2000);
		expect(body.info).toEqual({ name: 'Harbour' });
	});
});
