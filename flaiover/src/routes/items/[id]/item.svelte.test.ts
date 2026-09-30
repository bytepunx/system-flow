// S-0154: a story page follows the project: its status, history, children, and body change on the
// page as they change in the project, without a reload.
import { afterEach, describe, expect, it, vi } from 'vitest';
import { flushSync, mount, unmount } from 'svelte';

const api = vi.fn();
vi.mock('$lib/api', () => ({ api: (...args: unknown[]) => api(...args) }));
vi.mock('$app/state', () => ({
	page: { params: { id: 'S-0154' }, url: new URL('http://localhost/items/S-0154') }
}));
vi.mock('$app/paths', () => ({
	resolve: (route: string, params: Record<string, string>) =>
		route.replace(/\[(\.\.\.)?(\w+)\]/, (_, __, k) => params[k] ?? '')
}));
const enhance = vi.hoisted(() => vi.fn(async () => {}));
vi.mock('$lib/markdown', async (actual) => ({
	...(await actual<typeof import('$lib/markdown')>()),
	enhance
}));
const events = vi.hoisted(() => ({ heard: [] as { change?: (path: string) => void }[] }));
vi.mock('$lib/events', () => ({
	listen: (l: { change?: (path: string) => void }) => {
		events.heard.push(l);
		return () => events.heard.splice(events.heard.indexOf(l), 1);
	},
	follow: (_kinds: string[], f: () => void) => {
		const l = { change: () => f() };
		events.heard.push(l);
		return () => events.heard.splice(events.heard.indexOf(l), 1);
	}
}));
const changed = (path: string) => events.heard.forEach((l) => l.change?.(path));

import ItemPage from './+page.svelte';

const answer = (body: unknown) => ({ ok: true, json: async () => body });
const settle = async () => {
	for (let i = 0; i < 20; i++) await new Promise((r) => setTimeout(r, 0));
	flushSync();
};
const story = {
	id: 'S-0154',
	type: 'story',
	nature: 'feature',
	title: 'Story pages receive live updates',
	status: 'in-progress',
	created: '2026-09-29T05:55:11Z',
	updated: '2026-09-29T07:06:54Z',
	transitions: [{ to: 'in-progress', at: '2026-09-29T07:06:54Z', by: 'agent-S-0154' }],
	path: 'wip/kanban/stories/S-0154-story-pages-receive-live-updates.md',
	archived: false,
	body: '# S-0154\n\n## Goal\n\nPages follow the project.\n'
};

describe('the item page (S-0154)', () => {
	let c: ReturnType<typeof mount> | undefined;
	afterEach(() => {
		if (c) unmount(c);
		c = undefined;
		api.mockReset();
		enhance.mockClear();
		document.body.innerHTML = '';
	});

	const serve = (item: typeof story, children: unknown[] = []) =>
		api.mockImplementation(async (url: string) => {
			if (url === '/api/items/S-0154') return answer({ item, children });
			if (url === '/api/board') return answer({ writable: false });
			if (url.startsWith('/api/threads')) return answer([]);
			return answer({ enabled: false });
		});
	const asked = (url: string) => api.mock.calls.filter(([u]) => u === url).length;
	const history = () => document.querySelector('aside ol')!.textContent!.replace(/\s+/g, ' ');

	it('shows a move, a new task, and a changed body when the project changes', async () => {
		serve(story);
		c = mount(ItemPage, { target: document.body });
		await settle();
		expect(document.querySelector('h1')!.textContent).toContain('Story pages receive live updates');
		expect(enhance).toHaveBeenCalledTimes(1);

		// the agent writes a task and moves the story to review
		const review = {
			...story,
			status: 'review',
			transitions: [
				...story.transitions,
				{ to: 'review', at: '2026-09-29T08:00:00Z', by: 'agent-S-0154' }
			]
		};
		serve(review, [{ ...story, id: 'T-0553', type: 'task', title: 'A task', status: 'done' }]);
		changed('wip/kanban/stories/S-0154-story-pages-receive-live-updates.md');
		await settle();
		expect(document.querySelector('[data-testid="item-line"]')!.textContent).toContain('review');
		expect(history()).toContain('2026-09-29T08:00:00Z review by agent-S-0154');
		expect(document.body.textContent).toContain('T-0553');
		expect(document.body.textContent).toContain('Review this story');
		// the body did not change: it is not drawn again, and writable is not asked again
		expect(enhance).toHaveBeenCalledTimes(1);
		expect(asked('/api/board')).toBe(1);

		serve({ ...review, body: '# S-0154\n\n## Goal\n\nPages follow the project, live.\n' });
		changed('wip/kanban/stories/S-0154-story-pages-receive-live-updates.md');
		await settle();
		expect(document.querySelector('article')!.textContent).toContain(
			'Pages follow the project, live.'
		);
		expect(enhance).toHaveBeenCalledTimes(2);
	});

	it("says in the story's thread that its agent is working on the operator's reply", async () => {
		const reply = {
			id: 'TH-0040',
			title: 'Which animation?',
			anchor: { path: story.path, item: 'S-0154' },
			status: 'answered',
			participants: ['agent-S-0154', 'alex'],
			updated: '2026-09-29T08:05:00Z',
			entries: [
				{ at: '2026-09-29T08:00:00Z', author: 'agent-S-0154', text: 'Dots or a bar?' },
				{ at: '2026-09-29T08:05:00Z', author: 'alex', text: 'Dots.', operator: true }
			]
		};
		const run = {
			story: 'S-0154',
			command: 'claude',
			agent: 'agent-S-0154',
			started: '2026-09-29T07:06:00Z'
		};
		api.mockImplementation(async (url: string) => {
			if (url === '/api/items/S-0154') return answer({ item: story, children: [] });
			if (url === '/api/board') return answer({ writable: false });
			if (url.startsWith('/api/threads')) return answer([reply]);
			if (url === '/api/host-agent')
				return answer({
					enabled: true,
					state: { command: 'claude', stories: { 'S-0154': { state: 'working', run } } }
				});
			return answer({});
		});
		c = mount(ItemPage, { target: document.body });
		await settle();
		expect(document.querySelector('[data-testid="agent-working"]')!.textContent).toContain(
			'agent-S-0154 is working on your reply'
		);
	});

	it('dismisses the notice a refused edit leaves by its X (S-0151)', async () => {
		const ticked = {
			...story,
			body: '# S-0154\n\n## Goal\n\nPages follow the project.\n\n## Acceptance criteria\n\n- [ ] Live\n'
		};
		api.mockImplementation(async (url: string) => {
			if (url === '/api/items/S-0154') return answer({ item: ticked, children: [] });
			if (url === '/api/board') return answer({ writable: true });
			if (url === '/api/items/S-0154/edit')
				return { ok: false, json: async () => ({ error: 'the story is being edited' }) };
			if (url.startsWith('/api/threads')) return answer([]);
			return answer({ enabled: false });
		});
		c = mount(ItemPage, { target: document.body });
		await settle();
		const box = document.querySelector<HTMLInputElement>('input[aria-label="criterion 1"]')!;
		box.click();
		await settle();
		const notice = () => document.querySelector('[data-testid="item-notice"]');
		expect(notice()!.textContent).toBe('refused: the story is being edited');
		document.querySelector<HTMLButtonElement>('[data-testid="dismiss"]')!.click();
		flushSync();
		expect(notice()).toBeNull();
	});

	describe('the New link at the top (S-0171)', () => {
		const show = async (item: Record<string, unknown>, writable = true) => {
			api.mockImplementation(async (url: string) => {
				if (url === '/api/items/S-0154') return answer({ item, children: [] });
				if (url === '/api/board') return answer({ writable });
				if (url.startsWith('/api/threads')) return answer([]);
				return answer({ enabled: false });
			});
			c = mount(ItemPage, { target: document.body });
			await settle();
			return document.querySelector<HTMLAnchorElement>('[data-testid="new-same-type"]');
		};

		it("names a story and opens the form for a story under this one's epic", async () => {
			const link = await show({ ...story, parent: 'E-0013' });
			expect(link!.textContent).toBe('New story');
			expect(link!.getAttribute('href')).toBe('/new?type=story&parent=E-0013');
		});

		it('names a story with no epic and opens the form for a story alone', async () => {
			const link = await show(story);
			expect(link!.getAttribute('href')).toBe('/new?type=story');
		});

		it('names an epic and opens the form for an epic', async () => {
			const link = await show({ ...story, id: 'E-0013', type: 'epic', parent: undefined });
			expect(link!.textContent).toBe('New epic');
			expect(link!.getAttribute('href')).toBe('/new?type=epic');
		});

		it("is not offered on a task, which is the agent's to write, or when the dashboard cannot write", async () => {
			expect(await show({ ...story, id: 'T-0606', type: 'task', parent: 'S-0171' })).toBeNull();
			unmount(c!);
			c = undefined;
			document.body.innerHTML = '';
			expect(await show(story, false)).toBeNull();
		});
	});

	it('stops following the project once it is left', async () => {
		serve(story);
		c = mount(ItemPage, { target: document.body });
		await settle();
		expect(events.heard.length).toBeGreaterThan(0);
		unmount(c);
		c = undefined;
		expect(events.heard).toHaveLength(0);
	});
});
