// S-0098: /api/agent asks a connected project's flai for project.info, to say it answers. A
// repository offered for import has none to give, by design, and is not reported as failing.
// S-0534: it carries agent.status's budget beside the connection. S-0544: every answer carries the
// dashboard's own version, beside the connected flai's.
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

	// S-0534: beside the connection, the providers flai serve holds every agent start on, from
	// agent.status's state.budget; none while it holds none, and a failed ask costs the page nothing.
	const budget = [
		{
			provider: 'claude-code',
			limit: 'usage_limit_reached',
			status: 429,
			since: '2026-10-28T05:30:00Z',
			until: '2026-10-28T09:00:00Z'
		},
		{
			provider: 'codex',
			since: '2026-10-28T07:55:00Z',
			until: '2026-10-28T08:00:00Z',
			probing: true
		}
	];
	const connected = { configured: true, connected: true, missing: [] };
	const answering = (state: Record<string, unknown> | Error) =>
		vi.fn(async (method: string) => {
			if (method === 'project.info') return { name: 'Harbour' };
			if (state instanceof Error) throw state;
			return { enabled: true, state };
		});

	it('returns the budget agent.status holds', async () => {
		const { body, ask } = await answerFor(connected, answering({ command: 'claude', budget }));
		expect(ask).toHaveBeenCalledWith('agent.status', {}, 2000);
		expect(body.info).toEqual({ name: 'Harbour' });
		expect(body.budget).toEqual(budget);
		expect(body.error).toBeUndefined();
	});

	it.each([
		['no budget', { command: 'claude' }],
		['an empty budget', { command: 'claude', budget: [] }]
	])('returns no budget while agent.status holds %s', async (_, state) => {
		const { body } = await answerFor(connected, answering(state));
		expect(body.budget).toBeUndefined();
		expect(body.info).toEqual({ name: 'Harbour' });
	});

	it('returns no budget, and no error, when agent.status is not answered', async () => {
		const { body } = await answerFor(connected, answering(new Error('did not answer in time')));
		expect(body.budget).toBeUndefined();
		expect(body.error).toBeUndefined();
		expect(body.info).toEqual({ name: 'Harbour' });
	});

	// S-0237: a flai refused as unsigned never connects, so the status alone would say "not
	// connected"; the kept refusal becomes an error naming the flai, the reason, and what to do.
	const refused = { at: '2026-10-12T10:00:00Z', reason: 'missing stamp', flai: '1.44.0' };

	it('names a refused flai, why, and what to do, with no connection up', async () => {
		const { body, ask } = await answerFor({ configured: true, connected: false, refused });
		expect(ask).not.toHaveBeenCalled();
		expect(body.connected).toBe(false);
		expect(body.refused).toEqual(refused);
		expect(body.error).toContain(
			'this dashboard refused flai 1.44.0 on the host: unsigned (missing stamp)'
		);
		expect(body.error).toContain('flai host upgrade');
		expect(body.error).toContain('flai self-upgrade');
		expect(body.error).toContain('a development build of flai carries no stamp');
	});

	it('names a refused flai that gave no version without an empty name', async () => {
		const { body } = await answerFor({
			configured: true,
			connected: false,
			refused: { ...refused, reason: 'invalid signature', flai: '' }
		});
		expect(body.error).toContain(
			'this dashboard refused a flai on the host: unsigned (invalid signature)'
		);
	});

	it('says nothing of a refusal while a flai is connected and serves the project', async () => {
		const ask = vi.fn(async () => ({ name: 'Harbour' }));
		const { body } = await answerFor(
			{ configured: true, connected: true, signed: true, missing: [], refused },
			ask
		);
		expect(body.error).toBeUndefined();
		expect(body.info).toEqual({ name: 'Harbour' });
	});

	it('says nothing more of a flai that was never refused and is not connected', async () => {
		const { body } = await answerFor({ configured: true, connected: false });
		expect(body.error).toBeUndefined();
	});

	it("carries the dashboard's own version beside the connected flai's", async () => {
		const { version } = await import('$lib/server/metrics');
		const { body } = await answerFor(
			{ ...connected, flai: '1.77.2' },
			vi.fn(async () => ({ name: 'Harbour' }))
		);
		expect(body.flai).toBe('1.77.2');
		expect(body.dashboard).toBe(version);
		expect(version).not.toBe('');
	});

	it("carries the dashboard's own version with no flai connected", async () => {
		const { version } = await import('$lib/server/metrics');
		const { body } = await answerFor({ configured: true, connected: false });
		expect(body.dashboard).toBe(version);
	});
});
