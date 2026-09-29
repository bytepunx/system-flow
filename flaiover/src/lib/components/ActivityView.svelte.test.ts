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

	it('shows the narratives alone when flai knows of no agent', async () => {
		c = mount(ActivityView, { target: document.body, props: { streams: [stream('S-0001')] } });
		await settle();
		expect(card('S-0001')).toBeDefined();
		expect(document.querySelector('[data-testid="agent-stream"]')).toBeNull();
		expect(api).not.toHaveBeenCalled();
	});
});
