import { afterEach, describe, expect, it, vi } from 'vitest';
import { flushSync, mount, unmount } from 'svelte';
import ActivityView from './ActivityView.svelte';
import type { StoryActivity } from '$lib/activity';

const api = vi.fn();
vi.mock('$lib/api', () => ({ api: (...args: unknown[]) => api(...args) }));
vi.mock('$app/paths', () => ({
	resolve: (route: string, params: Record<string, string>) =>
		route.replace('[id]', params.id ?? '').replace('[...path]', params.path ?? '')
}));

const settle = async () => {
	for (let i = 0; i < 20; i++) await Promise.resolve();
	flushSync();
};
const stream = (id: string) => ({
	stream: id,
	title: `Story ${id}`,
	agent: `agent-${id}`,
	session: 's',
	updated: '2026-09-29T05:00:00Z',
	age_seconds: 60,
	status: 'in-progress',
	blocked: false,
	path: `wip/agents/${id}.md`
});
const run = (id: string, started = '2026-09-29T05:00:00Z') => ({
	story: id,
	command: 'claude',
	agent: `agent-${id}`,
	started
});
const card = (id: string) =>
	[...document.querySelectorAll('li.rounded')].find((li) => li.textContent!.includes(id)) as
		HTMLElement | undefined;

describe('ActivityView agent streams (S-0142)', () => {
	let c: ReturnType<typeof mount> | undefined;
	afterEach(() => {
		if (c) unmount(c);
		c = undefined;
		api.mockReset();
		document.body.innerHTML = '';
	});

	it("shows each started agent's stream on its story's card, open while it runs, and a card for an agent with no narrative yet", async () => {
		api.mockImplementation(async (path: string) => ({
			ok: true,
			status: 200,
			json: async () => ({
				story: path.split('/').at(-1),
				agent: 'a',
				started: '2026-09-29T05:00:00Z',
				running: false,
				from: 0,
				next: 0,
				size: 0,
				entries: [{ kind: 'text', text: `said in ${path.split('/').at(-1)}` }]
			})
		}));
		const agents: Record<string, StoryActivity> = {
			'S-0001': { state: 'working', run: run('S-0001') },
			'S-0002': { state: 'worked', run: { ...run('S-0002'), outcome: 'worked' } },
			// a held story that never had an agent: flai's stand-in run has not started
			'S-0003': {
				state: 'waiting',
				run: { story: 'S-0003', command: '', agent: '', started: '' },
				hold: { code: 'overlap', reason: 'held' }
			},
			'S-0009': { state: 'working', run: run('S-0009') }
		};
		c = mount(ActivityView, {
			target: document.body,
			props: { streams: [stream('S-0001'), stream('S-0002'), stream('S-0003')], agents }
		});
		await settle();

		const box = (id: string) =>
			card(id)?.querySelector<HTMLDetailsElement>('[data-testid="agent-stream"]') ?? null;
		expect(box('S-0001')!.open).toBe(true);
		expect(box('S-0001')!.textContent).toContain('said in S-0001');
		expect(card('S-0001')!.textContent).toContain('agent working (claude)');
		expect(box('S-0002')!.open).toBe(false);
		expect(box('S-0002')!.textContent).toContain('said in S-0002');
		expect(box('S-0003')).toBeNull();
		const alone = document.querySelector<HTMLElement>('[data-testid="unnarrated"]')!;
		expect(alone.textContent).toContain('S-0009');
		expect(alone.textContent).toContain('no narrative yet');
		expect(alone.querySelector('[data-testid="agent-stream"]')!.textContent).toContain(
			'said in S-0009'
		);
		expect(api.mock.calls.map(([p]) => p).sort()).toEqual([
			'/api/agent-stream/S-0001',
			'/api/agent-stream/S-0002',
			'/api/agent-stream/S-0009'
		]);
	});

	it('shows when the narrative was written and its last log entry in the local zone (S-0329)', async () => {
		// just after midnight in UTC, the evening before in New York, where the tests run
		c = mount(ActivityView, {
			target: document.body,
			props: {
				streams: [
					{
						...stream('S-0001'),
						updated: '2026-01-15T00:20:00Z',
						last_log: { at: '2026-01-15T00:10:00Z', text: 'T-0001 done.' }
					}
				]
			}
		});
		await settle();
		const text = card('S-0001')!.textContent!.replace(/\s+/g, ' ');
		expect(text).toContain('2026-01-14 19:10 EST T-0001 done.');
		expect(text).not.toContain('2026-01-15');
		expect(card('S-0001')!.querySelector('dd[title]')!.getAttribute('title')).toBe(
			'2026-01-14 19:20 EST'
		);
	});

	it('shows the narratives alone when flai knows of no agent', async () => {
		c = mount(ActivityView, { target: document.body, props: { streams: [stream('S-0001')] } });
		await settle();
		expect(card('S-0001')).toBeDefined();
		expect(document.querySelector('[data-testid="agent-stream"]')).toBeNull();
		expect(api).not.toHaveBeenCalled();
	});
});

describe('ActivityView stopping an agent (S-0170)', () => {
	let c: ReturnType<typeof mount> | undefined;
	afterEach(() => {
		if (c) unmount(c);
		c = undefined;
		api.mockReset();
		document.body.innerHTML = '';
	});
	const stopButton = (id: string) =>
		card(id)?.querySelector<HTMLButtonElement>('[data-testid="agent-stop"]') ?? null;
	const agents = (): Record<string, StoryActivity> => ({
		'S-0001': { state: 'working', run: { ...run('S-0001'), pid: 4242 } },
		'S-0002': { state: 'worked', run: { ...run('S-0002'), ended: 'e', outcome: 'worked' } },
		'S-0003': {
			state: 'waiting',
			thread: 'TH-0007',
			why: 'waiting for an answer to TH-0007',
			run: { ...run('S-0003'), ended: 'e', outcome: 'asked', thread: 'TH-0007' }
		}
	});
	// the streams answer with nothing; a stop answers as the route does
	const answer = (stop: { ok: boolean; body: unknown }) =>
		api.mockImplementation(async (path: string) =>
			path.endsWith('/agent')
				? { ok: stop.ok, status: stop.ok ? 200 : 400, json: async () => stop.body }
				: {
						ok: true,
						status: 200,
						json: async () => ({ running: true, from: 0, next: 0, size: 0, entries: [] })
					}
		);
	const posts = () => api.mock.calls.filter(([p]) => String(p).endsWith('/agent'));

	it('offers Stop only for an agent that runs or waits for an answer, and only while it may', async () => {
		answer({ ok: true, body: {} });
		const streams = [stream('S-0001'), stream('S-0002'), stream('S-0003')];
		c = mount(ActivityView, { target: document.body, props: { streams, agents: agents() } });
		await settle();
		expect(document.querySelector('[data-testid="agent-stop"]')).toBeNull();
		unmount(c);
		c = mount(ActivityView, {
			target: document.body,
			props: { streams, agents: agents(), canStop: true }
		});
		await settle();
		expect(stopButton('S-0001')).not.toBeNull();
		expect(stopButton('S-0002')).toBeNull();
		expect(stopButton('S-0003')).not.toBeNull();
	});

	it('warns what stopping does, and stops the agent only once the operator confirms', async () => {
		answer({ ok: true, body: { story: 'S-0001' } });
		const onstopped = vi.fn();
		c = mount(ActivityView, {
			target: document.body,
			props: { streams: [stream('S-0001')], agents: agents(), canStop: true, onstopped }
		});
		await settle();
		stopButton('S-0001')!.click();
		flushSync();
		const dialog = document.querySelector<HTMLElement>('[role="dialog"]')!;
		expect(dialog.textContent).toContain("Stop S-0001's agent?");
		// when it started, in the local zone (S-0329): the tests run in New York
		expect(dialog.textContent).toContain('agent-S-0001, started 2026-09-29 01:00 EDT');
		expect(dialog.textContent).not.toContain('2026-09-29T05:00:00Z');
		const warning = dialog.querySelector('[data-warning]')!.textContent!.replace(/\s+/g, ' ');
		expect(warning).toContain('every process it started, is ended now (pid 4242)');
		expect(warning).toContain('stays as it left it, committed or not');
		expect(warning).toContain('S-0001 stays in in-progress, with no agent');
		expect(warning).toContain('Retry');
		expect(posts()).toHaveLength(0);

		// keeping it changes nothing
		[...dialog.querySelectorAll('button')]
			.find((b) => b.textContent === 'Keep it running')!
			.click();
		flushSync();
		expect(document.querySelector('[role="dialog"]')).toBeNull();
		expect(posts()).toHaveLength(0);

		stopButton('S-0001')!.click();
		flushSync();
		document.querySelector<HTMLButtonElement>('[data-testid="agent-stop-confirm"]')!.click();
		await settle();
		expect(posts()).toEqual([
			[
				'/api/items/S-0001/agent',
				expect.objectContaining({ method: 'POST', body: JSON.stringify({ action: 'stop' }) })
			]
		]);
		expect(onstopped).toHaveBeenCalledWith('S-0001');
		expect(document.querySelector('[role="dialog"]')).toBeNull();
	});

	it("says an agent waiting for an answer is not started again, and flai's refusal", async () => {
		answer({ ok: false, body: { error: "S-0003's agent is not running: it ended" } });
		const onstopped = vi.fn();
		c = mount(ActivityView, {
			target: document.body,
			props: { streams: [stream('S-0003')], agents: agents(), canStop: true, onstopped }
		});
		await settle();
		stopButton('S-0003')!.click();
		flushSync();
		const warning = document.querySelector('[data-warning]')!.textContent!.replace(/\s+/g, ' ');
		expect(warning).toContain('It ended waiting for an answer to TH-0007');
		expect(warning).toContain('not started again when the answer comes');
		document.querySelector<HTMLButtonElement>('[data-testid="agent-stop-confirm"]')!.click();
		await settle();
		expect(document.querySelector('[role="alert"]')!.textContent).toContain('is not running');
		expect(document.querySelector('[role="dialog"]')).not.toBeNull();
		expect(onstopped).not.toHaveBeenCalled();
	});
});
