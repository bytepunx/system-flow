// S-0154: a story page follows the project: its status, history, children, and body change on the
// page as they change in the project, without a reload. S-0460: its sections are panes chosen
// from a sidebar, so each case opens the pane that holds what it checks.
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
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
const openEditor = vi.fn();
vi.mock('$lib/editor.svelte', () => ({
	openEditor: (...args: unknown[]) => openEditor(...args)
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
import { inboxState } from '$lib/inbox.svelte';
import { page } from '$app/state';

const home = 'http://localhost/items/S-0154';
const scrollIntoView = Element.prototype.scrollIntoView;

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
const thread = (id: string, title: string, status = 'open') => ({
	id,
	title,
	anchor: { path: story.path, item: 'S-0154' },
	status,
	participants: ['agent-S-0154'],
	updated: '2026-09-29T08:00:00Z',
	entries: [{ at: '2026-09-29T08:00:00Z', author: 'agent-S-0154', text: `${title}?` }]
});
const conversation = {
	id: 'MS-0003',
	title: 'Who owns the port?',
	from: 'S-0154',
	to: 'S-0155',
	about: ['flaiover/src/lib/server/port.ts'],
	status: 'open',
	closed: false,
	closed_reason: '',
	awaiting: 'S-0155',
	participants: ['agent-S-0154'],
	created: '2026-09-29T07:10:00Z',
	updated: '2026-09-29T07:10:00Z',
	path: 'wip/messages/MS-0003-who-owns-the-port.md',
	entries: [
		{
			at: '2026-09-29T07:10:00Z',
			author: 'agent-S-0154',
			story: 'S-0154',
			text: 'Do you change the port?'
		}
	]
};

// Arrive on a pane by its hash, as a shared link does.
const at = (hash: string) => history.replaceState(null, '', `/items/S-0154#${hash}`);
// Arrive by a link with a query, as the inbox makes.
const linked = (query: string) => Object.assign(page, { url: new URL(home + query) });
// The panes the sidebar lists, by id, and the one shown.
const panes = () =>
	[...document.querySelectorAll('[data-sidebar-list="top"] > li > a')].map((a) =>
		a.getAttribute('data-sidebar-entry')
	);
const showing = () =>
	document.querySelector('[data-sidebar-pane]')?.getAttribute('data-sidebar-pane');
// Record what is scrolled into view.
const scrolls = () => {
	const into: Element[] = [];
	Element.prototype.scrollIntoView = function (this: Element) {
		into.push(this);
	};
	return into;
};

describe('the item page (S-0154)', () => {
	let c: ReturnType<typeof mount> | undefined;
	beforeEach(() => history.replaceState(null, '', '/items/S-0154'));
	// jsdom follows a link's hash a task later: let it land before the next case sets its own.
	afterEach(async () => {
		await settle();
		if (c) unmount(c);
		c = undefined;
		api.mockReset();
		enhance.mockClear();
		document.body.innerHTML = '';
		Element.prototype.scrollIntoView = scrollIntoView;
		history.replaceState(null, '', '/');
		linked('');
	});
	// Choose a pane from the sidebar, as the reader does.
	const open = async (id: string) => {
		document.querySelector<HTMLAnchorElement>(`[data-sidebar-entry="${id}"]`)!.click();
		await settle();
	};

	const serve = (item: typeof story, children: unknown[] = []) =>
		api.mockImplementation(async (url: string) => {
			if (url === '/api/items/S-0154') return answer({ item, children });
			if (url === '/api/board') return answer({ writable: false });
			if (url.startsWith('/api/threads')) return answer([]);
			return answer({ enabled: false });
		});
	const asked = (url: string) => api.mock.calls.filter(([u]) => u === url).length;
	const moves = () =>
		document.querySelector('[data-sidebar-pane="history"] ol')!.textContent!.replace(/\s+/g, ' ');

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
		expect(document.querySelector('[data-sidebar-pane="story"]')!.textContent).toContain(
			'Review this story'
		);
		// the body did not change: it is not drawn again, and writable is not asked again
		expect(enhance).toHaveBeenCalledTimes(1);
		expect(asked('/api/board')).toBe(1);
		await open('history');
		// every time in the local zone (S-0329): the tests run in New York
		expect(moves()).toContain('2026-09-29 01:55 EDT created');
		expect(moves()).toContain('2026-09-29 04:00 EDT review by agent-S-0154');
		expect(moves()).not.toContain('Z ');
		await open('children');
		expect(document.querySelector('[data-sidebar-pane="children"]')!.textContent).toContain(
			'T-0553'
		);
		// the body's pane draws a new article when it is shown again
		await open('story');
		expect(enhance).toHaveBeenCalledTimes(2);

		serve({ ...review, body: '# S-0154\n\n## Goal\n\nPages follow the project, live.\n' });
		changed('wip/kanban/stories/S-0154-story-pages-receive-live-updates.md');
		await settle();
		expect(document.querySelector('article')!.textContent).toContain(
			'Pages follow the project, live.'
		);
		expect(enhance).toHaveBeenCalledTimes(3);
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
		at('threads');
		c = mount(ItemPage, { target: document.body });
		await settle();
		// the agent's pane is not shown, but what its agent is doing reaches the threads
		expect(showing()).toBe('threads');
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

		it('names a story and opens the form for a story on its own, even under an epic (S-0192)', async () => {
			const link = await show({ ...story, parent: 'E-0013' });
			expect(link!.textContent).toBe('New story');
			expect(link!.getAttribute('href')).toBe('/new?type=story');
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

	describe('the Create story link on an epic (S-0191)', () => {
		const epic = { ...story, id: 'E-0013', type: 'epic', parent: undefined, status: 'in-progress' };
		const show = async (item: Record<string, unknown>, writable = true) => {
			api.mockImplementation(async (url: string) => {
				if (url === '/api/items/S-0154') return answer({ item, children: [] });
				if (url === '/api/board') return answer({ writable });
				if (url.startsWith('/api/threads')) return answer([]);
				return answer({ enabled: false });
			});
			c = mount(ItemPage, { target: document.body });
			await settle();
			return document.querySelector<HTMLAnchorElement>('[data-testid="new-child-story"]');
		};
		const again = () => {
			unmount(c!);
			c = undefined;
			document.body.innerHTML = '';
		};

		it('comes after New epic and opens the form for a story under the epic', async () => {
			const link = await show(epic);
			expect(link!.textContent).toBe('Create story');
			expect(link!.getAttribute('href')).toBe('/new?type=story&parent=E-0013');
			const links = [...link!.parentElement!.querySelectorAll('a')].map((a) => a.textContent);
			expect(links).toEqual(['New epic', 'Create story']);
		});

		it('is not offered on a closed or archived epic, whose stories the form would not take', async () => {
			expect(await show({ ...epic, status: 'done' })).toBeNull();
			again();
			expect(await show({ ...epic, status: 'cancelled' })).toBeNull();
			again();
			expect(await show({ ...epic, archived: true })).toBeNull();
		});

		it('is not offered on a story or a task, or when the dashboard cannot write', async () => {
			expect(await show({ ...story, parent: 'E-0013' })).toBeNull();
			again();
			expect(await show({ ...story, id: 'T-0606', type: 'task', parent: 'S-0171' })).toBeNull();
			again();
			expect(await show(epic, false)).toBeNull();
		});
	});

	describe('the New sibling link on a story (S-0192)', () => {
		const sibling = { ...story, parent: 'E-0013' };
		const epic = { ...story, id: 'E-0013', type: 'epic', parent: undefined, status: 'in-progress' };
		const show = async (
			item: Record<string, unknown>,
			{ writable = true, parent = (): unknown => answer({ item: epic, children: [] }) } = {}
		) => {
			api.mockImplementation(async (url: string) => {
				if (url === '/api/items/S-0154') return answer({ item, children: [] });
				if (url === '/api/items/E-0013') return parent();
				if (url === '/api/board') return answer({ writable });
				if (url.startsWith('/api/threads')) return answer([]);
				return answer({ enabled: false });
			});
			c = mount(ItemPage, { target: document.body });
			await settle();
			return document.querySelector<HTMLAnchorElement>('[data-testid="new-sibling-story"]');
		};
		const again = () => {
			unmount(c!);
			c = undefined;
			document.body.innerHTML = '';
		};
		const withEpic = (e: Record<string, unknown>) => ({
			parent: () => answer({ item: { ...epic, ...e }, children: [] })
		});

		it('comes after New story and opens the form for a story under the same epic', async () => {
			const link = await show(sibling);
			expect(link!.textContent).toBe('New sibling');
			expect(link!.getAttribute('href')).toBe('/new?type=story&parent=E-0013');
			const links = [...link!.parentElement!.querySelectorAll('a')].map((a) => a.textContent);
			expect(links).toEqual(['New story', 'New sibling']);
		});

		it('asks after the epic once, not at every change in the project', async () => {
			await show(sibling);
			changed(story.path);
			await settle();
			expect(api.mock.calls.filter(([u]) => u === '/api/items/E-0013').length).toBe(1);
			expect(document.querySelector('[data-testid="new-sibling-story"]')).not.toBeNull();
		});

		it('is not offered on a story with no epic', async () => {
			expect(await show(story)).toBeNull();
			expect(api.mock.calls.some(([u]) => u === '/api/items/E-0013')).toBe(false);
		});

		it('is not offered under a closed or archived epic, which the form would not take', async () => {
			expect(await show(sibling, withEpic({ status: 'done' }))).toBeNull();
			again();
			expect(await show(sibling, withEpic({ status: 'cancelled' }))).toBeNull();
			again();
			expect(await show(sibling, withEpic({ archived: true }))).toBeNull();
		});

		it('is not offered when the epic cannot be read', async () => {
			const failed = () => ({ ok: false, json: async () => ({ error: 'no such item' }) });
			expect(await show(sibling, { parent: failed })).toBeNull();
			again();
			const offline = () => {
				throw new Error('offline');
			};
			expect(await show(sibling, { parent: offline })).toBeNull();
		});

		it('is not offered when the dashboard cannot write, or on an epic', async () => {
			expect(await show(sibling, { writable: false })).toBeNull();
			again();
			expect(await show(epic)).toBeNull();
		});
	});

	// S-0173: the inbox leads an open question here, so the story's page shows its open questions.
	it("shows the story's open questions from the inbox, and a task's page none", async () => {
		inboxState.data = {
			total: 1,
			counts: { thread: 0, question: 1, review: 0, blocked: 0, overlap: 0 },
			notes: [],
			entries: [
				{
					key: 'question:S-0154:abc',
					kind: 'question',
					title: 'Which port should it use?',
					href: '/items/S-0154?question=question%3AS-0154%3Aabc',
					item: 'S-0154'
				}
			]
		};
		serve(story);
		at('threads');
		c = mount(ItemPage, { target: document.body });
		await settle();
		expect(document.querySelector('[data-testid="open-questions"]')?.textContent).toContain(
			'Which port should it use?'
		);
		unmount(c);
		serve({ ...story, type: 'task' });
		c = mount(ItemPage, { target: document.body });
		await settle();
		expect(document.querySelector('[data-testid="open-questions"]')).toBeNull();
		inboxState.data = null;
	});

	// S-0199: a story's page shows its cost of delay and forecast, read only, each only when set;
	// S-0204: the inputs and the planner's value each with who set them and when (ADR-0080)
	describe('its cost of delay and forecast', () => {
		beforeEach(() => at('usage'));
		const text = (testid: string) =>
			[...(document.querySelector(`[data-testid="${testid}"]`)?.children ?? [])].map((d) =>
				d.textContent!.trim()
			);
		const stale = () => document.querySelector('[data-testid="item-cost-of-delay-stale"]');
		const costOfDelay = (inputsAt: string) => ({
			inputs: {
				revenue_per_week: 1200,
				time_lost_per_cycle: '6h',
				by: 'alex',
				at: inputsAt
			},
			value: 1650,
			by: 'planner',
			at: '2026-10-03T09:20:00Z'
		});

		it("shows a story's cost of delay and forecast when set, and neither when not", async () => {
			serve({
				...story,
				cost_of_delay: costOfDelay('2026-10-03T09:00:00Z'),
				forecast: { delivery: '2026-10-07T17:00:00Z', by: 'planner', at: '2026-10-03T09:30:00Z' }
			} as typeof story);
			c = mount(ItemPage, { target: document.body });
			await settle();
			expect(text('item-cost-of-delay')).toEqual([
				'cost of delay: 1,650 per week, set by planner at 2026-10-03 05:20 EDT',
				'inputs: revenue 1,200 per week · 6h lost per cycle, set by alex at 2026-10-03 05:00 EDT'
			]);
			expect(stale()).toBeNull();
			expect(text('item-forecast')).toEqual([
				'forecast: delivery 2026-10-07 13:00 EDT',
				'set by planner at 2026-10-03 05:30 EDT'
			]);
			unmount(c);
			serve(story);
			c = mount(ItemPage, { target: document.body });
			await settle();
			expect(document.querySelector('[data-testid="item-cost-of-delay"]')).toBeNull();
			expect(document.querySelector('[data-testid="item-forecast"]')).toBeNull();
		});

		it('notes that the value is out of date when the inputs were set after it', async () => {
			serve({ ...story, cost_of_delay: costOfDelay('2026-10-03T10:00:00Z') } as typeof story);
			c = mount(ItemPage, { target: document.body });
			await settle();
			expect(stale()?.textContent?.trim()).toBe(
				"The inputs changed after the value was set, so the planner's value is out of date."
			);
			expect(text('item-cost-of-delay')[1]).toBe(
				'inputs: revenue 1,200 per week · 6h lost per cycle, set by alex at 2026-10-03 06:00 EDT'
			);
		});

		it('says an epic has no value yet when only its inputs are set, with no note', async () => {
			const { inputs } = costOfDelay('2026-10-03T10:00:00Z');
			serve({ ...story, type: 'epic', cost_of_delay: { inputs } } as typeof story);
			c = mount(ItemPage, { target: document.body });
			await settle();
			expect(text('item-cost-of-delay')).toEqual([
				'cost of delay: no value yet',
				'inputs: revenue 1,200 per week · 6h lost per cycle, set by alex at 2026-10-03 06:00 EDT'
			]);
			expect(stale()).toBeNull();
		});

		// S-0225, ADR-0083: what the planner expects it to cost, and what planning it cost, apart
		it("shows a planned story's expected cost and the planner's usage before any agent works it", async () => {
			const planned = {
				...story,
				status: 'backlog',
				forecast: { duration: '3h', by: 'planner', at: '2026-10-03T09:30:00Z' },
				usage: {
					source: 'sum',
					seconds: 0,
					models: [],
					strategic: [
						{
							kind: 'planner',
							seconds: 412,
							estimated: true,
							models: [
								{
									model: 'claude-opus-5-5',
									input: 12000,
									output: 0,
									cache_read: 800000,
									cache_write: 0,
									cost: 0.81
								}
							]
						}
					]
				}
			};
			api.mockImplementation(async (url: string) => {
				if (url === '/api/items/S-0154')
					return answer({
						item: planned,
						children: [],
						expected_cost: { cost: 12, from: 'forecast', estimated: true }
					});
				if (url === '/api/board') return answer({ writable: false });
				if (url.startsWith('/api/threads')) return answer([]);
				return answer({ enabled: false });
			});
			c = mount(ItemPage, { target: document.body });
			await settle();
			expect(
				document.querySelector('[data-testid="item-expected-cost"]')?.textContent?.trim()
			).toBe('expected cost: $12.00 (estimated, from the forecast)');
			const usage = document.querySelector('[data-testid="item-usage"]')!.textContent!;
			expect(usage).toContain('planner, strategic: 812.0K tokens · $0.810 (estimated) · 6m');
			// no agent has worked it, so there is no agents' line of nothing
			expect(usage).not.toContain('of agent work');
			unmount(c);
			serve(story);
			c = mount(ItemPage, { target: document.body });
			await settle();
			expect(document.querySelector('[data-testid="item-expected-cost"]')).toBeNull();
			expect(document.querySelector('[data-testid="item-usage"]')).toBeNull();
		});
	});

	// S-0201: a draft story says so beside its title and is finalized from its page without a reload
	describe('a draft story (S-0201)', () => {
		const draftStory = { ...story, status: 'backlog', draft: true };
		const epic = { ...story, id: 'E-0016', type: 'epic', draft: undefined };
		const indicator = () => document.querySelector('[data-testid="item-draft"]');
		const button = () => document.querySelector<HTMLButtonElement>('[data-testid="item-finalize"]');
		const show = async (
			item: Record<string, unknown>,
			children: unknown[] = [],
			writable = true
		) => {
			api.mockImplementation(async (url: string) => {
				if (url === '/api/items/S-0154') return answer({ item, children });
				if (url === '/api/board') return answer({ writable });
				if (url.startsWith('/api/threads')) return answer([]);
				return answer({ enabled: false });
			});
			c = mount(ItemPage, { target: document.body });
			await settle();
		};
		const again = () => {
			unmount(c!);
			c = undefined;
			document.body.innerHTML = '';
		};

		it('shows [Draft] beside the title and a Finalize button', async () => {
			await show(draftStory);
			expect(indicator()!.textContent).toBe('[Draft]');
			expect(indicator()!.closest('h1')).not.toBeNull();
			expect(button()!.textContent).toBe('Finalize');
		});

		it('shows neither on a story that is not a draft, nor on an epic', async () => {
			await show(story);
			expect(indicator()).toBeNull();
			expect(button()).toBeNull();
			again();
			await show({ ...epic, draft: true });
			expect(indicator()).toBeNull();
			expect(button()).toBeNull();
		});

		it('offers Finalize only where the dashboard may write, and not on an archived story', async () => {
			await show(draftStory, [], false);
			expect(indicator()).not.toBeNull();
			expect(button()).toBeNull();
			again();
			await show({ ...draftStory, archived: true });
			expect(button()).toBeNull();
		});

		it('finalizes through the API and, once reloaded, shows neither', async () => {
			let current: Record<string, unknown> = draftStory;
			api.mockImplementation(async (url: string, init?: RequestInit) => {
				if (url === '/api/items/S-0154') return answer({ item: current, children: [] });
				if (url === '/api/items/S-0154/finalize' && init?.method === 'POST') {
					current = { ...draftStory, draft: false };
					return answer({ id: 'S-0154', draft: false, warnings: [] });
				}
				if (url === '/api/board') return answer({ writable: true });
				if (url.startsWith('/api/threads')) return answer([]);
				return answer({ enabled: false });
			});
			c = mount(ItemPage, { target: document.body });
			await settle();
			button()!.click();
			await settle();
			expect(asked('/api/items/S-0154/finalize')).toBe(1);
			expect(indicator()).toBeNull();
			expect(button()).toBeNull();
			expect(document.querySelector('[data-testid="item-notice"]')!.textContent).toBe('done');
		});

		it('says why when flai refuses, and keeps the draft', async () => {
			api.mockImplementation(async (url: string) => {
				if (url === '/api/items/S-0154') return answer({ item: draftStory, children: [] });
				if (url === '/api/items/S-0154/finalize')
					return { ok: false, json: async () => ({ error: 'S-0154 is not a draft' }) };
				if (url === '/api/board') return answer({ writable: true });
				if (url.startsWith('/api/threads')) return answer([]);
				return answer({ enabled: false });
			});
			c = mount(ItemPage, { target: document.body });
			await settle();
			button()!.click();
			await settle();
			expect(document.querySelector('[data-testid="item-notice"]')!.textContent).toBe(
				'refused: S-0154 is not a draft'
			);
			expect(indicator()).not.toBeNull();
		});

		it("marks an epic's draft story among its children", async () => {
			at('children');
			await show(epic, [
				{ ...story, id: 'S-0201', title: 'Drafted', draft: true },
				{ ...story, id: 'S-0200', title: 'Finalized' }
			]);
			const marks = [...document.querySelectorAll('[data-testid="child-draft"]')];
			expect(marks).toHaveLength(1);
			expect(marks[0].textContent).toBe('Draft');
			expect(marks[0].closest('li')!.textContent).toContain('S-0201');
		});
	});

	// S-0208: an open epic's or story's page offers Plan while the plan host action is on
	describe('the Plan action (S-0208)', () => {
		const epic = { ...story, id: 'E-0016', type: 'epic' };
		const button = () => document.querySelector<HTMLButtonElement>('[data-testid="item-plan"]');
		const show = async (item: Record<string, unknown>, writable = true, enabled = true) => {
			api.mockImplementation(async (url: string, init?: RequestInit) => {
				if (url === '/api/items/S-0154') return answer({ item, children: [] });
				if (url === '/api/board') return answer({ writable });
				if (url.startsWith('/api/threads')) return answer([]);
				if (url === `/api/items/${item.id}/plan` && init?.method === 'POST')
					return {
						ok: false,
						json: async () => ({
							error: `the planner is already running for ${item.id} (pid 42, started 2026-10-03T10:00:00Z); one item has one planner at a time`
						})
					};
				if (url === `/api/items/${item.id}/plan`)
					return answer({ plan_enabled: enabled, run: null });
				return answer({ enabled: false });
			});
			c = mount(ItemPage, { target: document.body });
			await settle();
		};
		const again = () => {
			unmount(c!);
			c = undefined;
			document.body.innerHTML = '';
		};

		it("is among an epic's and a story's actions while the host action is on", async () => {
			await show(epic);
			expect(button()!.textContent).toBe('Plan');
			again();
			await show(story);
			expect(button()!.textContent).toBe('Plan');
			again();
			await show(story, true, false);
			expect(button()).toBeNull();
		});

		it('is not offered on a task, a closed or archived item, or when the dashboard cannot write', async () => {
			for (const item of [
				{ ...story, id: 'T-0788', type: 'task' },
				{ ...story, status: 'done' },
				{ ...epic, status: 'cancelled' },
				{ ...story, archived: true }
			]) {
				await show(item);
				expect(button()).toBeNull();
				again();
			}
			await show(story, false);
			expect(button()).toBeNull();
		});

		// S-0300: an epic's page plans the epic, and only while the epic is open
		it("is on an open epic's page and asks the planner for the epic, but not on a done or archived epic's", async () => {
			await show({ ...epic, status: 'in-progress' });
			expect(button()!.textContent).toBe('Plan');
			button()!.click();
			await settle();
			expect(
				api.mock.calls.some(
					([u, init]) => u === '/api/items/E-0016/plan' && init?.method === 'POST'
				)
			).toBe(true);
			again();
			for (const item of [
				{ ...epic, status: 'done' },
				{ ...epic, status: 'in-progress', archived: true }
			]) {
				await show(item);
				expect(button()).toBeNull();
				again();
			}
		});

		it("says in the page's notice why flai refused", async () => {
			await show(epic);
			button()!.click();
			await settle();
			expect(document.querySelector('[data-testid="item-notice"]')!.textContent).toContain(
				'refused: the planner is already running for E-0016'
			);
		});
	});

	// S-0176: a story's page shows its task plan, and follows it as tasks move
	it("shows a story's task plan when flai sends one, and follows it as the tasks move", async () => {
		const plan = {
			tasks: [
				{ id: 'T-0001', state: 'waiting', after: ['T-0002'], waiting_for: ['T-0002'] },
				{ id: 'T-0002', state: 'in-progress' },
				{ id: 'T-0003', state: 'ready' }
			],
			layers: [['T-0002', 'T-0003'], ['T-0001']]
		};
		const children = [{ ...story, id: 'T-0002', type: 'task', title: 'The flai side' }];
		const show = (p?: unknown) =>
			api.mockImplementation(async (url: string) => {
				if (url === '/api/items/S-0154') return answer({ item: story, children, plan: p });
				if (url === '/api/board') return answer({ writable: false });
				if (url.startsWith('/api/threads')) return answer([]);
				return answer({ enabled: false });
			});
		const shown = () => document.querySelector('[data-testid="task-plan"]');
		const tasks = () =>
			[...document.querySelectorAll('[data-testid="plan-task"]')].map((t) =>
				t.textContent!.replace(/\s+/g, ' ').trim()
			);
		show(plan);
		at('children');
		c = mount(ItemPage, { target: document.body });
		await settle();
		expect(shown()!.closest('[data-sidebar-pane="children"]')).not.toBeNull();
		expect(tasks()).toEqual([
			'T-0001 waiting for T-0002',
			'T-0002 in progress',
			'T-0003 ready to start'
		]);
		expect(document.querySelectorAll('[data-testid="plan-layer"]')).toHaveLength(2);
		expect(shown()!.querySelector('a[href="/items/T-0002"]')!.getAttribute('title')).toBe(
			'The flai side'
		);

		// T-0002 is done: T-0001 can start
		show({
			...plan,
			tasks: [
				{ id: 'T-0001', state: 'ready', after: ['T-0002'] },
				{ id: 'T-0002', state: 'done' },
				{ id: 'T-0003', state: 'ready' }
			]
		});
		changed('wip/kanban/tasks/T-0002-the-flai-side.md');
		await settle();
		expect(tasks()).toEqual(['T-0001 ready to start', 'T-0002 done', 'T-0003 ready to start']);

		// a story without tasks has no plan, and shows none
		show(undefined);
		changed('wip/kanban/tasks/T-0002-the-flai-side.md');
		await settle();
		expect(shown()).toBeNull();
	});

	it('shows when the item was blocked and unblocked in the local zone (S-0329)', async () => {
		serve({
			...story,
			blocked: [
				{ from: '2026-09-29T08:10:00Z', until: '2026-09-29T09:20:00Z', reason: 'waiting on X' },
				{ from: '2026-09-30T03:00:00Z', reason: 'waiting on Y' }
			]
		} as typeof story);
		at('history');
		c = mount(ItemPage, { target: document.body });
		await settle();
		const headings = [...document.querySelectorAll('[data-sidebar-pane="history"] h3')];
		const heading = headings.find((h) => h.textContent === 'Blocked')!;
		const blocked = [...heading.nextElementSibling!.querySelectorAll('li')].map((li) =>
			li.textContent!.replace(/\s+/g, ' ').trim()
		);
		expect(blocked).toEqual([
			'2026-09-29 04:10 EDT → 2026-09-29 05:20 EDT: waiting on X',
			'2026-09-29 23:00 EDT → open: waiting on Y'
		]);
	});

	// ADR-0093: a done story says who accepted it and on which day, from its last transition to
	// done; the orchestrator's acceptance links to the evidence it wrote in the story's Notes.
	describe('who accepted a done story (ADR-0093)', () => {
		const done = (by: string, notes = '') => ({
			...story,
			status: 'done',
			transitions: [
				...story.transitions,
				{ to: 'review', at: '2026-10-05T22:00:00Z', by: 'agent-S-0154' },
				// just after midnight in UTC, which is the day before in New York (S-0329)
				{ to: 'done', at: '2026-10-06T02:30:00Z', by }
			],
			body: story.body + notes
		});
		const line = () =>
			document
				.querySelector('[data-testid="item-accepted"]')
				?.textContent!.replace(/\s+/g, ' ')
				.trim();
		const link = () =>
			document.querySelector<HTMLAnchorElement>('[data-testid="item-accepted-evidence"]');
		const notes =
			'\n## Notes\n\n### Accepted by the orchestrator\n\n- Verified: 0123456789abcdef0123456789abcdef01234567\n- At: 2026-10-06T09:30:00Z\n\nVerdict: pass, tests and lint clean\n- 1: `flai/cmd/accept.go`\n';

		it('names the orchestrator and links to its evidence in the Notes', async () => {
			serve(done('orchestrator', notes));
			c = mount(ItemPage, { target: document.body });
			await settle();
			expect(line()).toBe(
				'Accepted by orchestrator on 2026-10-05, with its evidence in the Notes.'
			);
			expect(link()!.getAttribute('href')).toBe('#accepted-by-the-orchestrator');
			// the link lands on the heading the body renders
			expect(document.getElementById('accepted-by-the-orchestrator')!.textContent).toContain(
				'Accepted by the orchestrator'
			);
		});

		it("shows the body's pane when its evidence link is followed from another (S-0460)", async () => {
			serve(done('orchestrator', notes));
			at('history');
			c = mount(ItemPage, { target: document.body });
			await settle();
			expect(document.getElementById('accepted-by-the-orchestrator')).toBeNull();
			link()!.click();
			await settle();
			expect(showing()).toBe('story');
			expect(location.hash).toBe('#accepted-by-the-orchestrator');
			expect(document.getElementById('accepted-by-the-orchestrator')).not.toBeNull();
		});

		it('names the operator with no evidence link, and says nothing of a story that is not done', async () => {
			serve(done('alex'));
			c = mount(ItemPage, { target: document.body });
			await settle();
			expect(line()).toBe('Accepted by alex on 2026-10-05.');
			expect(link()).toBeNull();
			unmount(c);
			document.body.innerHTML = '';
			serve(story);
			c = mount(ItemPage, { target: document.body });
			await settle();
			expect(line()).toBeUndefined();
		});
	});

	// S-0388: an open story whose branch undoes main's change says so beside HELD, with the paths'
	// count, and leads to its review page, where each is kept or restored (ADR-0155).
	describe('a story whose branch undoes main (S-0388)', () => {
		const show = async (undoes_main?: string[]) => {
			api.mockImplementation(async (url: string) => {
				if (url === '/api/items/S-0154')
					return answer({ item: story, children: [], ...(undoes_main ? { undoes_main } : {}) });
				if (url === '/api/board') return answer({ writable: false });
				if (url.startsWith('/api/threads')) return answer([]);
				return answer({ enabled: false });
			});
			c = mount(ItemPage, { target: document.body });
			await settle();
		};
		const badge = () => document.querySelector<HTMLAnchorElement>('[data-testid="undoes-main"]');

		it('says one path and links to the review page', async () => {
			await show(['docs/a.md']);
			expect(badge()!.textContent!.trim()).toBe('undoes main on 1 path');
			expect(badge()!.getAttribute('href')).toBe('/review/S-0154');
			expect(badge()!.closest('[data-testid="item-line"]')).not.toBeNull();
		});

		it('counts several paths', async () => {
			await show(['docs/a.md', 'flai/cmd/accept.go', 'design/system/workflow.md']);
			expect(badge()!.textContent!.trim()).toBe('undoes main on 3 paths');
			expect(badge()!.getAttribute('href')).toBe('/review/S-0154');
		});

		it('is not shown when flai sends none', async () => {
			await show();
			expect(badge()).toBeNull();
		});
	});

	// S-0394: a story's page shows its pull request as flai serve last read it from the forge
	describe("a story's pull request (S-0394)", () => {
		const pullRequest = {
			story: 'S-0154',
			pull_request: 7,
			url: 'https://github.com/bytepunx/system-flow/pull/7',
			head_sha: '0123456789abcdef0123456789abcdef01234567',
			read_at: '2026-09-29T08:00:00Z',
			checks: [
				{
					name: 'flai',
					state: 'failed',
					url: 'https://github.com/bytepunx/system-flow/actions/runs/1',
					since: '2026-09-29T07:58:00Z'
				}
			],
			reviews: [{ author: 'alex', state: 'approved', submitted_at: '2026-09-29T07:59:00Z' }]
		};
		const show = async (item: Record<string, unknown>, forge: () => unknown) => {
			api.mockImplementation(async (url: string) => {
				if (url === '/api/items/S-0154') return answer({ item, children: [] });
				if (url.endsWith('/forge')) return forge();
				if (url === '/api/board') return answer({ writable: false });
				if (url.startsWith('/api/threads')) return answer([]);
				return answer({ enabled: false });
			});
			c = mount(ItemPage, { target: document.body });
			await settle();
		};
		const section = () => document.querySelector('[data-testid="forge-status"]');

		it('shows the link to it, each required check, and each review', async () => {
			at('agent');
			await show(story, () => answer({ story: 'S-0154', pull_request: pullRequest }));
			expect(asked('/api/items/S-0154/forge')).toBeGreaterThan(0);
			expect(showing()).toBe('agent');
			const pane = section()!.closest<HTMLElement>('[data-testid="item-agent-pane"]')!;
			expect(pane.hidden).toBe(false);
			const link = section()!.querySelector<HTMLAnchorElement>('h2 a')!;
			expect(link.textContent!.trim()).toBe('Pull request #7');
			expect(link.getAttribute('href')).toBe(pullRequest.url);
			const check = section()!.querySelector('[data-check="flai"]')!;
			expect(check.textContent).toContain('failed');
			expect(check.querySelector('a')!.getAttribute('href')).toBe(pullRequest.checks[0].url);
			expect(section()!.querySelector('[data-review="alex"]')!.textContent).toContain('approved');
		});

		it('shows nothing for a story with none, and with no agent either, no Agent pane', async () => {
			await show(story, () => answer({ story: 'S-0154', pull_request: null }));
			expect(section()).toBeNull();
			expect(panes()).not.toContain('agent');
		});

		it('shows nothing, and the page still, when it cannot be read', async () => {
			await show(story, () => ({ ok: false, json: async () => ({ error: 'no flai connected' }) }));
			expect(section()).toBeNull();
			expect(document.querySelector('h1')!.textContent).toContain(story.title);
			unmount(c!);
			document.body.innerHTML = '';
			await show(story, () => Promise.reject(new Error('network down')));
			expect(section()).toBeNull();
			expect(document.querySelector('h1')!.textContent).toContain(story.title);
		});

		it("is not asked for on an epic's or a task's page", async () => {
			const forge = () => answer({ story: 'S-0154', pull_request: pullRequest });
			await show({ ...story, id: 'E-0013', type: 'epic' }, forge);
			expect(section()).toBeNull();
			unmount(c!);
			document.body.innerHTML = '';
			await show({ ...story, id: 'T-0606', type: 'task', parent: 'S-0154' }, forge);
			expect(section()).toBeNull();
			expect(api.mock.calls.some(([u]) => String(u).endsWith('/forge'))).toBe(false);
		});
	});

	// S-0336: a story's page lists its conversations with other stories' agents, on their own pane
	// since S-0460, which a story with none leaves out
	describe("a story's conversations (S-0336)", () => {
		beforeEach(() => at('messages'));
		// all is what the sidebar reads, closed ones too; conversations what the pane lists at first
		const show = async (
			item: Record<string, unknown>,
			conversations: unknown[] = [],
			all = conversations
		) => {
			api.mockImplementation(async (url: string) => {
				if (url === '/api/items/S-0154') return answer({ item, children: [] });
				if (url === '/api/board') return answer({ writable: false });
				if (url.startsWith('/api/threads')) return answer([]);
				if (url.startsWith('/api/messages'))
					return answer(url.includes('all=1') ? all : conversations);
				return answer({ enabled: false });
			});
			c = mount(ItemPage, { target: document.body });
			await settle();
		};
		const section = () => document.querySelector('section[data-messages]');

		it('lists them on their own pane, each with the other story linked and its entries', async () => {
			await show(story, [conversation]);
			expect(asked('/api/messages?story=S-0154')).toBe(1);
			expect(showing()).toBe('messages');
			expect(section()!.getAttribute('data-messages')).toBe('S-0154');
			expect(document.querySelector('section[data-threads]')).toBeNull();
			const li = section()!.querySelector('[data-conversation="MS-0003"]')!;
			expect(li.textContent).toContain('Who owns the port?');
			expect(li.querySelector('[data-with]')!.textContent!.replace(/\s+/g, ' ')).toBe('to S-0155');
			expect(li.querySelector('a[data-story="S-0155"]')!.getAttribute('href')).toBe(
				'/items/S-0155'
			);
			expect(li.querySelector('[data-entry="S-0154"]')!.textContent).toContain(
				'Do you change the port?'
			);
		});

		it('says so when the story has only closed ones', async () => {
			await show(story, [], [{ ...conversation, status: 'closed', closed: true }]);
			expect(section()!.querySelector('[data-empty]')!.textContent).toContain(
				'No open conversations here.'
			);
			expect(section()!.querySelector('[data-conversation]')).toBeNull();
		});

		it('leaves the pane out for a story with none, and shows its body instead', async () => {
			await show(story);
			expect(panes()).not.toContain('messages');
			expect(showing()).toBe('story');
			expect(section()).toBeNull();
		});

		it("is not on an epic's or a task's page", async () => {
			await show({ ...story, id: 'E-0013', type: 'epic' }, [conversation]);
			expect(section()).toBeNull();
			unmount(c!);
			document.body.innerHTML = '';
			await show({ ...story, id: 'T-0606', type: 'task', parent: 'S-0154' }, [conversation]);
			expect(section()).toBeNull();
			expect(api.mock.calls.some(([u]) => String(u).startsWith('/api/messages'))).toBe(false);
		});
	});

	// S-0465, I-0145: an acceptance from the item page refused for lint the fix left lists the
	// refusal's findings, each file opened on the checkout flai names, never on one read from its words.
	describe('a move to done refused for lint the fix left (S-0465)', () => {
		const finding = (over: Record<string, unknown>) => ({
			tier: 'markdown',
			checkout: 'story',
			story: 'S-0154',
			path: 'docs/bad.md',
			line: 2,
			rule: 'MD009',
			message: 'Trailing spaces',
			...over
		});
		const refuse = (body: Record<string, unknown>) =>
			api.mockImplementation(async (url: string) => {
				if (url === '/api/items/S-0154')
					return answer({ item: { ...story, status: 'review' }, children: [] });
				if (url === '/api/board') return answer({ writable: true });
				if (url.startsWith('/api/threads')) return answer([]);
				if (url === '/api/items/S-0154/acceptance')
					return answer({ id: 'S-0154', blockers: [], plan: null });
				if (url === '/api/items/S-0154/move') return { ok: false, json: async () => body };
				return answer({ enabled: false });
			});
		const click = (text: string) =>
			[...document.querySelectorAll('button')].find((b) => b.textContent === text)!.click();
		const accept = async () => {
			c = mount(ItemPage, { target: document.body });
			await settle();
			click('→ done');
			await settle();
			click('Accept');
			await settle();
		};
		const link = (label: string) =>
			document.querySelector<HTMLButtonElement>(`button[aria-label="${label}"]`)!;
		afterEach(() => openEditor.mockReset());

		it("lists the findings and opens each file on its own checkout, whatever the refusal's words say", async () => {
			refuse({
				error: "rule: S-0154 cannot be accepted: the main checkout's uncommitted files hold lint",
				id: 'S-0154',
				reason: 'S-0154 cannot be accepted',
				findings: [
					finding({}),
					finding({ checkout: 'main', story: undefined, path: 'design/notes.md', line: 7 })
				]
			});
			await accept();
			const rows = [...document.querySelectorAll('[data-testid="accept-finding"]')].map(
				(e) => e.textContent ?? ''
			);
			expect(rows).toHaveLength(2);
			expect(rows[0]).toContain('docs/bad.md:2');
			expect(rows[1]).toContain('design/notes.md:7');
			expect(document.querySelector('[data-testid="item-notice"]')).toBeNull();
			link('Edit docs/bad.md at line 2').click();
			expect(openEditor).toHaveBeenLastCalledWith(
				expect.objectContaining({ path: 'docs/bad.md', line: 2, story: 'S-0154' })
			);
			link('Edit design/notes.md at line 7').click();
			const main = openEditor.mock.calls.at(-1)![0];
			expect(main).toMatchObject({ path: 'design/notes.md', line: 7 });
			expect(main.story).toBeUndefined();
		});

		it('keeps the text notice for a refusal without findings', async () => {
			refuse({ error: 'S-0154 has an open blocker' });
			await accept();
			expect(document.querySelector('[data-testid="accept-findings"]')).toBeNull();
			expect(document.querySelector('[data-testid="item-notice"]')!.textContent).toBe(
				'refused: S-0154 has an open blocker'
			);
		});
	});

	// S-0460: the header and actions stay on top, and the sections are panes chosen from a sidebar,
	// each left out when the item has nothing for it; the body is shown on arrival, and the inbox's
	// thread and question links open Questions and threads.
	describe('its panes (S-0460)', () => {
		const task = { ...story, id: 'T-0553', type: 'task', title: 'A task', status: 'done' };
		const run = {
			story: 'S-0154',
			command: 'claude',
			agent: 'agent-S-0154',
			started: '2026-09-29T07:06:00Z'
		};
		// A story with tasks, threads, conversations, an agent, and an owner: every pane has something.
		const full = (url: string) => {
			if (url === '/api/items/S-0154')
				return answer({ item: { ...story, owner: 'alex' }, children: [task] });
			if (url === '/api/board') return answer({ writable: true });
			if (url.startsWith('/api/threads'))
				return answer([
					thread('TH-0040', 'Dots or a bar'),
					thread('TH-0041', 'Which port', 'answered'),
					thread('TH-0042', 'Which colour', 'resolved')
				]);
			if (url.startsWith('/api/messages')) return answer([conversation]);
			if (url === '/api/host-agent')
				return answer({
					enabled: true,
					state: { command: 'claude', stories: { 'S-0154': { state: 'working', run } } }
				});
			return answer({ enabled: false });
		};
		const show = async (
			item: Record<string, unknown>,
			{ children = [] as unknown[], writable = false, threads = [] as unknown[] } = {}
		) => {
			api.mockImplementation(async (url: string) => {
				if (url === '/api/items/S-0154') return answer({ item, children });
				if (url === '/api/board') return answer({ writable });
				if (url.startsWith('/api/threads')) return answer(threads);
				return answer({ enabled: false });
			});
			c = mount(ItemPage, { target: document.body });
			await settle();
		};
		const again = () => {
			unmount(c!);
			c = undefined;
			document.body.innerHTML = '';
			history.replaceState(null, '', '/items/S-0154');
		};
		const labels = () =>
			[...document.querySelectorAll('[data-sidebar-list="top"] > li > a')].map((a) =>
				a.textContent!.replace(/\s+/g, ' ').trim()
			);
		const agentPane = () => document.querySelector<HTMLElement>('[data-testid="item-agent-pane"]')!;
		afterEach(() => {
			inboxState.data = null;
		});

		it('lists every pane in order for a story with tasks, threads, and messages, each showing its section', async () => {
			api.mockImplementation(async (url: string) => full(url));
			c = mount(ItemPage, { target: document.body });
			await settle();
			expect(panes()).toEqual([
				'story',
				'children',
				'threads',
				'messages',
				'agent',
				'history',
				'files',
				'usage'
			]);
			// the threads not resolved are counted
			expect(labels()).toEqual([
				'Story',
				'Tasks',
				'Questions and threads 2',
				'Messages',
				'Agent',
				'History',
				'Files',
				'Planning and usage'
			]);
			const holds: Record<string, string> = {
				children: '[data-sidebar-pane="children"] a[href="/items/T-0553"]',
				threads: 'section[data-threads="S-0154"]',
				messages: '[data-conversation="MS-0003"]',
				agent: '[data-testid="story-agent"]',
				history: '[data-sidebar-pane="history"] ol',
				files: `[data-sidebar-pane="files"] a[href="/docs/${story.path}"]`,
				usage: '[data-sidebar-pane="usage"] section',
				story: '[data-sidebar-pane="story"] article'
			};
			for (const [id, selector] of Object.entries(holds)) {
				await open(id);
				expect(showing()).toBe(id);
				expect(location.hash).toBe(`#${id}`);
				expect(document.querySelector(selector), id).not.toBeNull();
				expect(agentPane().hidden, id).toBe(id !== 'agent');
			}
			// the header and the actions stay above every pane
			for (const testid of ['item-line', 'edit-item']) {
				const el = document.querySelector(`[data-testid="${testid}"]`)!;
				expect(el.closest('[data-page-sidebar]'), testid).toBeNull();
			}
		});

		it('leaves out the panes an item has nothing for, and names an epic’s and a task’s own', async () => {
			// not writable, with no tasks, threads, conversations, agent, or planning
			await show(story);
			expect(panes()).toEqual(['story', 'history', 'files']);
			expect(labels()[0]).toBe('Story');
			again();
			await show(
				{ ...story, id: 'E-0013', type: 'epic' },
				{ children: [{ ...story, id: 'S-0155' }] }
			);
			expect(panes()).toEqual(['story', 'children', 'history', 'files']);
			expect(labels().slice(0, 2)).toEqual(['Epic', 'Stories']);
			again();
			// a writer opens threads on the item, so its pane stays; a task has no Messages or Agent
			await show({ ...story, id: 'T-0606', type: 'task', parent: 'S-0154' }, { writable: true });
			expect(panes()).toEqual(['story', 'threads', 'history', 'files']);
			expect(labels().slice(0, 2)).toEqual(['Task', 'Threads']);
			await open('files');
			expect(document.querySelector('[data-sidebar-pane="files"]')!.textContent).toContain(
				"Tasks are the agent's to write."
			);
		});

		it('shows the body on arrival, and for a hash that names no pane the item has', async () => {
			api.mockImplementation(async (url: string) => full(url));
			c = mount(ItemPage, { target: document.body });
			await settle();
			expect(showing()).toBe('story');
			expect(
				document.querySelector('[data-sidebar-entry="story"]')!.getAttribute('aria-current')
			).toBe('true');
			expect(document.querySelector('[data-sidebar-pane="story"] article')).not.toBeNull();
			expect(agentPane().hidden).toBe(true);
			again();
			at('children');
			await show(story);
			expect(showing()).toBe('story');
		});

		it('opens the pane a link’s hash names, and Messages for a conversation’s anchor', async () => {
			api.mockImplementation(async (url: string) => full(url));
			at('history');
			c = mount(ItemPage, { target: document.body });
			await settle();
			expect(showing()).toBe('history');
			again();
			const into = scrolls();
			at('MS-0003');
			c = mount(ItemPage, { target: document.body });
			await settle();
			expect(showing()).toBe('messages');
			expect(into).toContain(document.getElementById('MS-0003'));
		});

		it("opens Questions and threads on the inbox's link to a thread, with that thread in view", async () => {
			const into = scrolls();
			linked('?thread=TH-0041');
			await show(story, {
				threads: [thread('TH-0040', 'Dots or a bar'), thread('TH-0041', 'Which port')]
			});
			expect(showing()).toBe('threads');
			const section = document.querySelector('section[data-threads="S-0154"]')!;
			expect(section.textContent).toContain('Which port');
			expect(section.textContent).not.toContain('Dots or a bar');
			expect(into).toContain(section);
		});

		it("opens Questions and threads on the inbox's link to a question, with that question in view", async () => {
			const key = 'question:S-0154:abc';
			inboxState.data = {
				total: 1,
				counts: { thread: 0, question: 1, review: 0, blocked: 0, overlap: 0 },
				notes: [],
				entries: [
					{
						key,
						kind: 'question',
						title: 'Which port should it use?',
						href: `/items/S-0154?question=${encodeURIComponent(key)}`,
						item: 'S-0154'
					}
				]
			};
			const into = scrolls();
			linked(`?question=${encodeURIComponent(key)}`);
			await show(story);
			expect(showing()).toBe('threads');
			expect(labels()[1]).toBe('Questions and threads 1');
			const question = document.querySelector(`[data-question="${key}"]`)!;
			expect(question.getAttribute('aria-current')).toBe('true');
			expect(into).toContain(question);
		});

		it("shows the body's pane when edit… is pressed on another", async () => {
			api.mockImplementation(async (url: string) => {
				if (url === '/api/items/S-0154') return answer({ item: story, children: [] });
				if (url === '/api/board') return answer({ writable: true });
				if (url === '/api/items/S-0154/edit')
					return { ok: false, json: async () => ({ error: 'the story is being edited' }) };
				if (url.startsWith('/api/threads')) return answer([]);
				return answer({ enabled: false });
			});
			at('history');
			c = mount(ItemPage, { target: document.body });
			await settle();
			document.querySelector<HTMLButtonElement>('[data-testid="edit-item"]')!.click();
			await settle();
			expect(showing()).toBe('story');
			expect(asked('/api/items/S-0154/edit')).toBeGreaterThan(0);
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
