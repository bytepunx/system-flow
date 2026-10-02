// S-0178: the activity page holds still while the operator watches an agent's stream: a reload that
// fails, or an answer from the host-agent route with no flai connected, leaves what it shows in place.
import { afterEach, describe, expect, it, vi } from 'vitest';
import { flushSync, mount, unmount } from 'svelte';

const api = vi.fn();
vi.mock('$lib/api', () => ({ api: (...args: unknown[]) => api(...args) }));
vi.mock('$app/paths', () => ({
	resolve: (route: string, params: Record<string, string>) =>
		route.replace('[id]', params.id ?? '').replace('[...path]', params.path ?? '')
}));
// each followed kind's reload, to call as a change event would
const followed: { kinds: string[]; f: () => void }[] = [];
vi.mock('$lib/events', () => ({
	follow: (kinds: string[], f: () => void) => {
		followed.push({ kinds, f });
		return () => {};
	},
	listen: () => () => {},
	debounced: (f: () => void) => Object.assign(() => f(), { stop: () => {} })
}));

import ActivityPage from './+page.svelte';

const answer = (body: unknown) => ({ ok: true, status: 200, json: async () => body });
const refused = (error: string) => ({
	ok: false,
	status: 500,
	statusText: 'Internal Server Error',
	json: async () => ({ error })
});
const settle = async () => {
	for (let i = 0; i < 20; i++) await new Promise((r) => setTimeout(r, 0));
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
const read = (story: string) => ({
	story,
	agent: 'a',
	started: '2026-09-29T05:00:00Z',
	running: true,
	from: 0,
	next: 0,
	size: 0,
	entries: []
});
const withAgent = {
	enabled: true,
	state: {
		command: 'claude',
		stories: {
			'S-0001': {
				state: 'working',
				run: {
					story: 'S-0001',
					command: 'claude',
					agent: 'agent-S-0001',
					pid: 4242,
					started: '2026-09-29T05:00:00Z'
				}
			}
		}
	}
};
const reload = (kind: string) => followed.find((x) => x.kinds.includes(kind))!.f();
const card = () =>
	[...document.querySelectorAll('li.rounded')].find((li) => li.textContent!.includes('S-0001'));

describe('the activity page holds still (S-0178)', () => {
	let c: ReturnType<typeof mount> | undefined;
	afterEach(() => {
		if (c) unmount(c);
		c = undefined;
		api.mockReset();
		followed.length = 0;
		document.body.innerHTML = '';
	});

	it('keeps the streams it showed when a reload fails, with the error above them', async () => {
		let failing = false;
		api.mockImplementation(async (url: string) => {
			if (url === '/api/activity')
				return failing ? refused('flai is away') : answer({ streams: [stream('S-0001')] });
			if (url === '/api/host-agent') return answer({ enabled: false });
			return answer({});
		});
		c = mount(ActivityPage, { target: document.body });
		await settle();
		expect(card()).toBeDefined();
		expect(document.querySelector('[role="alert"]')).toBeNull();

		failing = true;
		reload('narrative');
		await settle();
		expect(document.querySelector('[role="alert"]')!.textContent).toContain('flai is away');
		expect(card()).toBeDefined();

		failing = false;
		reload('narrative');
		await settle();
		expect(document.querySelector('[role="alert"]')).toBeNull();
		expect(card()).toBeDefined();
	});

	it('says only the error when nothing was ever loaded', async () => {
		api.mockImplementation(async (url: string) =>
			url === '/api/activity' ? refused('flai is away') : answer({ enabled: false })
		);
		c = mount(ActivityPage, { target: document.body });
		await settle();
		expect(document.querySelector('[role="alert"]')!.textContent).toContain('flai is away');
		expect(document.body.textContent).not.toContain('Loading…');
		expect(card()).toBeUndefined();
	});

	it("keeps each agent's stream window when flai goes away, and no longer offers to stop it", async () => {
		let host: unknown = withAgent;
		api.mockImplementation(async (url: string) => {
			if (url === '/api/activity') return answer({ streams: [stream('S-0001')] });
			if (url === '/api/host-agent') return answer(host);
			if (url.startsWith('/api/agent-stream/')) return answer(read(url.split('/').at(-1)!));
			return answer({});
		});
		c = mount(ActivityPage, { target: document.body });
		await settle();
		const box = () => card()?.querySelector('[data-testid="agent-stream"]') ?? null;
		const stop = () => card()?.querySelector('[data-testid="agent-stop"]') ?? null;
		const before = box();
		expect(before).not.toBeNull();
		expect(stop()).not.toBeNull();

		// what the host-agent route answers with no flai connected
		host = { enabled: false };
		reload('thread');
		await settle();
		expect(box()).toBe(before);
		expect(stop()).toBeNull();
	});
});
