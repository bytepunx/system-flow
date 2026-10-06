import { afterEach, describe, expect, it, vi } from 'vitest';
import { flushSync, mount, unmount } from 'svelte';
import Threads from './Threads.svelte';

const api = vi.fn();
vi.mock('$lib/api', () => ({ api: (...args: unknown[]) => api(...args) }));
vi.mock('$app/paths', () => ({
	resolve: (route: string, params: Record<string, string>) =>
		route.replace(/\[(\.\.\.)?(\w+)\]/, (_, __, k) => params[k] ?? '')
}));

// The page's live events, as listen() would deliver them, with no wait for changes to settle.
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

const settle = async () => {
	for (let i = 0; i < 6; i++) await new Promise((r) => setTimeout(r, 0));
	flushSync();
};

const thread = (text: string) => ({
	id: 'TH-0001',
	title: 'A question',
	anchor: { path: 'wip/kanban/stories/S-0001-x.md', item: 'S-0001' },
	status: 'open',
	participants: ['agent'],
	updated: '2026-09-26T07:00:00Z',
	entries: [{ at: '2026-09-26T07:00:00Z', author: 'agent', text }]
});

const scrollIntoView = Element.prototype.scrollIntoView;

describe('Threads', () => {
	let c: ReturnType<typeof mount> | undefined;
	afterEach(() => {
		Element.prototype.scrollIntoView = scrollIntoView;
		if (c) unmount(c);
		c = undefined;
		api.mockReset();
		document.body.innerHTML = '';
	});

	it('renders an entry as markdown, not as raw text (S-0126)', async () => {
		const text =
			'Pick one:\n\n1. **Hold** at pull, see `touches`.\n2. Read [the survey](design/system/agent-coordination.md).\n';
		api.mockResolvedValue({ ok: true, json: async () => [thread(text)] });
		c = mount(Threads, { target: document.body, props: { on: 'S-0001' } });
		await settle();

		const entry = document.querySelector('article li')!;
		expect(entry.querySelectorAll('ol > li')).toHaveLength(2);
		expect(entry.querySelector('strong')?.textContent).toBe('Hold');
		expect(entry.querySelector('code')?.textContent).toBe('touches');
		expect(entry.querySelector('a')?.getAttribute('href')).toBe(
			'/docs/design/system/agent-coordination.md'
		);
		expect(entry.textContent).not.toContain('**');
	});

	it("sets the operator's entries apart from agents' by side and colour (S-0127)", async () => {
		const t = {
			...thread('A question for you.'),
			entries: [
				{
					at: '2026-09-26T07:00:00Z',
					author: 'agent-S-0001',
					text: 'A question for you.',
					operator: false
				},
				{ at: '2026-09-26T07:05:00Z', author: 'alex', text: 'Yes, do that.', operator: true }
			]
		};
		api.mockResolvedValue({ ok: true, json: async () => [t] });
		c = mount(Threads, { target: document.body, props: { on: 'S-0001' } });
		await settle();

		const [agent, operator] = [...document.querySelectorAll('article > ol > li')] as HTMLElement[];
		expect(agent.dataset.from).toBe('agent');
		expect(agent.className).toContain('items-start');
		expect(agent.querySelector('.prose')!.className).toContain('bg-raised');
		expect(operator.dataset.from).toBe('operator');
		expect(operator.className).toContain('items-end');
		expect(operator.querySelector('.prose')!.className).toContain('bg-primary-soft');
		expect(operator.textContent).toContain('alex');
	});

	const numbered = (n: number) => ({
		...thread(`Question ${n}.`),
		id: `TH-000${n}`,
		title: `Q${n}`
	});
	const pagers = () => [...document.querySelectorAll('nav[data-pager]')] as HTMLElement[];
	const arrow = (pager: HTMLElement, which: 'previous' | 'next') =>
		pager.querySelector(`button[aria-label="${which} thread"]`) as HTMLButtonElement;

	it('shows a reply as soon as the project says its thread changed, and stops listening once gone (S-0154)', async () => {
		api.mockResolvedValue({ ok: true, json: async () => [thread('Which port?')] });
		c = mount(Threads, { target: document.body, props: { on: 'S-0001' } });
		await settle();
		expect(document.querySelectorAll('article > ol > li')).toHaveLength(1);

		const answered = {
			...thread('Which port?'),
			status: 'answered',
			entries: [
				...thread('Which port?').entries,
				{ at: '2026-09-26T07:05:00Z', author: 'alex', text: 'Use port 8080.', operator: true }
			]
		};
		api.mockResolvedValue({ ok: true, json: async () => [answered] });
		changed('wip/threads/TH-0001-a-question.md');
		await settle();
		const entries = document.querySelectorAll('article > ol > li');
		expect(entries).toHaveLength(2);
		expect(entries[1].textContent).toContain('Use port 8080.');
		expect(document.querySelector('article header')!.textContent).toContain('answered');

		unmount(c);
		c = undefined;
		expect(events.heard).toHaveLength(0);
	});

	describe("the agent working on the operator's reply (S-0154)", () => {
		const replied = (status = 'answered') => ({
			...thread('Which port?'),
			status,
			entries: [
				{ at: '2026-09-26T07:00:00Z', author: 'agent-S-0001', text: 'Which port?' },
				{ at: '2026-09-26T07:05:00Z', author: 'alex', text: 'Use port 8080.', operator: true }
			]
		});
		const run = {
			story: 'S-0001',
			command: 'claude',
			agent: 'agent-S-0001',
			started: '2026-09-26T06:00:00Z'
		};
		const shown = async (t: unknown, agent?: unknown) => {
			api.mockResolvedValue({ ok: true, json: async () => [t] });
			c = mount(Threads, { target: document.body, props: { on: 'S-0001', agent } as never });
			await settle();
			return document.querySelector<HTMLElement>('[data-testid="agent-working"]');
		};

		it("shows a moving line on the agent's side while the agent works after a reply", async () => {
			const line = await shown(replied(), { state: 'working', run });
			expect(line).not.toBeNull();
			expect(line!.getAttribute('role')).toBe('status');
			expect(line!.textContent).toContain('agent-S-0001 is working on your reply');
			const dots = line!.querySelectorAll('[aria-hidden="true"] > span');
			expect(dots).toHaveLength(3);
			expect(dots[0].className).toContain('motion-safe:animate-bounce');
		});

		it('shows it while flai starts again the agent that ended asking this thread', async () => {
			const ended = { ...run, ended: '2026-09-26T07:01:00Z' };
			expect(
				await shown(replied(), { state: 'waiting', thread: 'TH-0001', run: ended })
			).not.toBeNull();
		});

		it("shows nothing after an agent's entry, on a resolved thread, or with no agent", async () => {
			expect(await shown(thread('Which port?'), { state: 'working', run })).toBeNull();
			unmount(c!);
			expect(await shown(replied('resolved'), { state: 'working', run })).toBeNull();
			unmount(c!);
			expect(await shown(replied())).toBeNull();
			unmount(c!);
			// waiting on another question, or finished: not at work on this reply
			expect(await shown(replied(), { state: 'waiting', thread: 'TH-0002', run })).toBeNull();
			unmount(c!);
			expect(await shown(replied(), { state: 'worked', run })).toBeNull();
		});
	});

	it('shows one thread at a time with its place above and below (S-0133)', async () => {
		api.mockResolvedValue({ ok: true, json: async () => [1, 2, 3].map(numbered) });
		c = mount(Threads, { target: document.body, props: { on: 'S-0001' } });
		await settle();

		expect(document.querySelectorAll('article')).toHaveLength(1);
		expect(document.querySelector('article')!.dataset.thread).toBe('TH-0001');
		const [above, below] = pagers();
		expect(above.dataset.pager).toBe('above');
		expect(below.dataset.pager).toBe('below');
		expect(above.compareDocumentPosition(document.querySelector('article')!)).toBe(
			Node.DOCUMENT_POSITION_FOLLOWING
		);
		for (const p of [above, below]) {
			expect(p.querySelector('span')!.textContent).toBe('1 of 3');
			expect(arrow(p, 'previous').disabled).toBe(true);
			expect(arrow(p, 'next').disabled).toBe(false);
		}
	});

	it('reads previous, the count, then next, with the upper pager beside the heading (S-0153)', async () => {
		api.mockResolvedValue({ ok: true, json: async () => [1, 2, 3].map(numbered) });
		c = mount(Threads, { target: document.body, props: { on: 'S-0001' } });
		await settle();

		const [above, below] = pagers();
		for (const p of [above, below])
			expect([...p.children].map((e) => e.getAttribute('aria-label') ?? e.textContent)).toEqual([
				'previous thread',
				'1 of 3',
				'next thread'
			]);
		expect(above.parentElement!.querySelector('h2')!.textContent).toBe('Threads');
	});

	it("pages with the left and right arrow keys on a pager's arrows, and stops at the ends (S-0153)", async () => {
		api.mockResolvedValue({ ok: true, json: async () => [1, 2, 3].map(numbered) });
		c = mount(Threads, { target: document.body, props: { on: 'S-0001' } });
		await settle();

		const press = (p: HTMLElement, key: string) => {
			const e = new KeyboardEvent('keydown', { key, bubbles: true, cancelable: true });
			(p.querySelector('button:not(:disabled)') as HTMLButtonElement).dispatchEvent(e);
			flushSync();
			return e;
		};
		const shownId = () => document.querySelector('article')!.dataset.thread;
		expect(press(pagers()[0], 'ArrowRight').defaultPrevented).toBe(true);
		expect(shownId()).toBe('TH-0002');
		press(pagers()[1], 'ArrowRight');
		press(pagers()[1], 'ArrowRight');
		expect(shownId()).toBe('TH-0003');
		press(pagers()[0], 'ArrowLeft');
		expect(shownId()).toBe('TH-0002');
		expect(press(pagers()[0], 'Enter').defaultPrevented).toBe(false);
		expect(shownId()).toBe('TH-0002');
	});

	it('moves between threads with the arrows, above or below, and stops at the ends (S-0133)', async () => {
		api.mockResolvedValue({ ok: true, json: async () => [1, 2, 3].map(numbered) });
		c = mount(Threads, { target: document.body, props: { on: 'S-0001' } });
		await settle();

		arrow(pagers()[0], 'next').click();
		flushSync();
		expect(document.querySelector('article')!.dataset.thread).toBe('TH-0002');
		expect(document.querySelector('article')!.textContent).toContain('Question 2.');
		arrow(pagers()[1], 'next').click();
		flushSync();
		expect(document.querySelector('article')!.dataset.thread).toBe('TH-0003');
		for (const p of pagers()) {
			expect(p.querySelector('span')!.textContent).toBe('3 of 3');
			expect(arrow(p, 'next').disabled).toBe(true);
		}
		arrow(pagers()[1], 'previous').click();
		flushSync();
		expect(pagers()[0].querySelector('span')!.textContent).toBe('2 of 3');
	});

	it('keeps the thread being read when the list reloads, and its place when it leaves (S-0133)', async () => {
		let list = [1, 2, 3].map(numbered);
		api.mockImplementation(async (path: string) => ({
			ok: true,
			json: async () => (path.startsWith('/api/threads?') ? list : {})
		}));
		c = mount(Threads, { target: document.body, props: { on: 'S-0001' } });
		await settle();
		arrow(pagers()[0], 'next').click();
		flushSync();

		// A reply reloads the list; a new thread ahead of this one must not move the reader.
		list = [numbered(4), ...list];
		const input = document.querySelector('article input') as HTMLInputElement;
		input.value = 'ok';
		input.dispatchEvent(new Event('input'));
		(document.querySelector('article form') as HTMLFormElement).requestSubmit();
		await settle();
		expect(document.querySelector('article')!.dataset.thread).toBe('TH-0002');
		expect(pagers()[0].querySelector('span')!.textContent).toBe('3 of 4');

		// Resolving it drops it from the list: the thread now in its place is shown.
		list = list.filter((t) => t.id !== 'TH-0002');
		(
			[...document.querySelectorAll('article button')].find(
				(b) => b.textContent === 'resolve'
			) as HTMLButtonElement
		).click();
		await settle();
		expect(document.querySelector('article')!.dataset.thread).toBe('TH-0003');
		expect(pagers()[0].querySelector('span')!.textContent).toBe('3 of 3');
	});

	it('shows a thread just opened (S-0133)', async () => {
		let list = [1, 2].map(numbered);
		api.mockImplementation(async (path: string, init?: { method?: string }) => {
			if (init?.method === 'POST') {
				list = [...list, numbered(3)];
				return { ok: true, json: async () => ({ id: 'TH-0003' }) };
			}
			return { ok: true, json: async () => list };
		});
		c = mount(Threads, { target: document.body, props: { on: 'S-0001' } });
		await settle();

		(
			[...document.querySelectorAll('button')].find(
				(b) => b.textContent === 'new thread'
			) as HTMLButtonElement
		).click();
		flushSync();
		const form = document.querySelector('section > form') as HTMLFormElement;
		for (const [sel, value] of [
			['input', 'Q3'],
			['textarea', 'Question 3.']
		]) {
			const el = form.querySelector(sel) as HTMLInputElement;
			el.value = value;
			el.dispatchEvent(new Event('input'));
		}
		form.requestSubmit();
		await settle();
		expect(document.querySelector('article')!.dataset.thread).toBe('TH-0003');
		expect(pagers()[0].querySelector('span')!.textContent).toBe('3 of 3');
	});

	const long = (n: number, entries: number) => ({
		...numbered(n),
		entries: Array.from({ length: entries }, (_, i) => ({
			at: `2026-09-26T07:0${i}:00Z`,
			author: i % 2 ? 'alex' : 'agent',
			text: `Entry ${i + 1}.`
		}))
	});
	const entryTexts = () =>
		[...document.querySelectorAll('article > ol > li .prose')].map((e) => e.textContent!.trim());
	const toggle = () => document.querySelector('article button[data-earlier]') as HTMLButtonElement;

	it('shows a long thread by its last two entries, the earlier ones behind a toggle (S-0153)', async () => {
		api.mockResolvedValue({ ok: true, json: async () => [long(1, 5), long(2, 3)] });
		c = mount(Threads, { target: document.body, props: { on: 'S-0001' } });
		await settle();

		expect(entryTexts()).toEqual(['Entry 4.', 'Entry 5.']);
		expect(toggle().textContent).toBe('show 3 earlier entries');
		expect(toggle().getAttribute('aria-expanded')).toBe('false');
		toggle().click();
		flushSync();
		expect(entryTexts()).toEqual(['Entry 1.', 'Entry 2.', 'Entry 3.', 'Entry 4.', 'Entry 5.']);
		expect(toggle().textContent).toBe('hide earlier entries');
		expect(toggle().getAttribute('aria-expanded')).toBe('true');

		// Each thread keeps its own choice while paging.
		arrow(pagers()[0], 'next').click();
		flushSync();
		expect(entryTexts()).toEqual(['Entry 2.', 'Entry 3.']);
		expect(toggle().textContent).toBe('show 1 earlier entry');
		arrow(pagers()[0], 'previous').click();
		flushSync();
		expect(entryTexts()).toHaveLength(5);
		toggle().click();
		flushSync();
		expect(entryTexts()).toEqual(['Entry 4.', 'Entry 5.']);
	});

	it('shows every entry of a thread of two, with no toggle (S-0153)', async () => {
		api.mockResolvedValue({ ok: true, json: async () => [long(1, 2)] });
		c = mount(Threads, { target: document.body, props: { on: 'S-0001' } });
		await settle();

		expect(entryTexts()).toEqual(['Entry 1.', 'Entry 2.']);
		expect(toggle()).toBeNull();
	});

	it('opens on the thread a link names and brings it into view (S-0155)', async () => {
		const scroll = vi.fn();
		Element.prototype.scrollIntoView = scroll;
		let list = [1, 2, 3].map(numbered);
		api.mockImplementation(async (path: string) => ({
			ok: true,
			json: async () => (path.startsWith('/api/threads?') ? list : {})
		}));
		c = mount(Threads, { target: document.body, props: { on: 'S-0001', select: 'TH-0002' } });
		await settle();

		expect(document.querySelector('article')!.dataset.thread).toBe('TH-0002');
		for (const p of pagers()) expect(p.querySelector('span')!.textContent).toBe('2 of 3');
		expect(scroll).toHaveBeenCalledTimes(1);
		expect(scroll.mock.instances[0]).toBe(document.querySelector('section[data-threads]'));

		// Paging away and a reload keep the reader where they went, not where the link pointed.
		arrow(pagers()[0], 'next').click();
		flushSync();
		list = [...list, numbered(4)];
		const input = document.querySelector('article input') as HTMLInputElement;
		input.value = 'ok';
		input.dispatchEvent(new Event('input'));
		(document.querySelector('article form') as HTMLFormElement).requestSubmit();
		await settle();
		expect(document.querySelector('article')!.dataset.thread).toBe('TH-0003');
		expect(scroll).toHaveBeenCalledTimes(1);
	});

	it('opens on the first thread when the one a link names is not there (S-0155)', async () => {
		api.mockResolvedValue({ ok: true, json: async () => [1, 2].map(numbered) });
		c = mount(Threads, { target: document.body, props: { on: 'S-0001', select: 'TH-0009' } });
		await settle();

		expect(document.querySelector('article')!.dataset.thread).toBe('TH-0001');
	});

	it('shows no pager for a single thread (S-0133)', async () => {
		api.mockResolvedValue({ ok: true, json: async () => [numbered(1)] });
		c = mount(Threads, { target: document.body, props: { on: 'S-0001' } });
		await settle();

		expect(document.querySelectorAll('article')).toHaveLength(1);
		expect(pagers()).toHaveLength(0);
	});

	// S-0173: the threads page shows every thread, each linking to what it is anchored on.
	it('without an anchor, lists every thread, links each to its anchor, and starts none', async () => {
		api.mockResolvedValue({
			ok: true,
			json: async () => [
				thread('on a story'),
				{
					...numbered(2),
					anchor: { path: 'design/issues/I-0027-a.md', heading: 'Instances' }
				}
			]
		});
		c = mount(Threads, { target: document.body, props: { writable: true } });
		await settle();

		expect(api).toHaveBeenCalledWith('/api/threads');
		expect(document.querySelector('section')!.dataset.threads).toBe('all');
		expect([...document.querySelectorAll('button')].map((b) => b.textContent)).not.toContain(
			'new thread'
		);
		expect(document.querySelector('article header a')!.getAttribute('href')).toBe('/items/S-0001');
		arrow(pagers()[0], 'next').click();
		flushSync();
		expect(document.querySelector('article header a')!.getAttribute('href')).toBe(
			'/docs/design/issues/I-0027-a.md'
		);
		// still answerable here
		expect(document.querySelector('article form input')).not.toBeNull();
	});
	// ADR-0090: a recommendation is marked, names its source, and is the answer once confirmed.
	describe('a pending recommendation', () => {
		const rec = {
			at: '2026-09-26T07:05:00Z',
			author: 'orchestrator',
			text: 'Use port 8080.\n\nSource: design/system/overview.md § Delivery sequence',
			recommendation: true,
			source: { path: 'design/system/overview.md', heading: 'Delivery sequence' }
		};
		const pending = {
			...thread('Which port?'),
			entries: [...thread('Which port?').entries, rec],
			pending_recommendation: rec
		};
		const confirmButton = () => document.querySelector<HTMLButtonElement>('button[data-confirm]');

		it('marks the entry, links its source to the document and heading, and offers Confirm', async () => {
			api.mockResolvedValue({ ok: true, json: async () => [pending] });
			c = mount(Threads, { target: document.body, props: { on: 'S-0001' } });
			await settle();

			const [question, recommended] = [
				...document.querySelectorAll('article > ol > li')
			] as HTMLElement[];
			expect(question.querySelector('[data-recommendation]')).toBeNull();
			expect(recommended.querySelector('[data-recommendation]')!.textContent).toBe(
				'recommendation'
			);
			expect(recommended.querySelector('.prose')!.textContent!.trim()).toBe('Use port 8080.');
			const source = recommended.querySelector('[data-source] a')!;
			expect(source.getAttribute('href')).toBe('/docs/design/system/overview.md#delivery-sequence');
			expect(source.textContent).toBe('design/system/overview.md § Delivery sequence');
			expect(document.querySelector('[data-pending]')!.textContent).toContain(
				"orchestrator's recommendation awaits your confirmation"
			);
			expect(confirmButton()!.textContent).toBe('Confirm');
		});

		it('confirms with one action, and the thread shows answered without a reload', async () => {
			let list: unknown[] = [pending];
			api.mockImplementation(async (path: string, init?: { method?: string }) => {
				if (init?.method === 'POST') {
					list = [
						{
							...pending,
							status: 'answered',
							pending_recommendation: null,
							entries: [
								...pending.entries,
								{
									at: '2026-09-26T07:10:00Z',
									author: 'alex',
									text: 'Confirmed the recommendation of 2026-09-26T07:05:00Z orchestrator.',
									operator: true
								}
							]
						}
					];
					return { ok: true, json: async () => ({ id: 'TH-0001', status: 'answered' }) };
				}
				return { ok: true, json: async () => list };
			});
			c = mount(Threads, { target: document.body, props: { on: 'S-0001' } });
			await settle();
			confirmButton()!.click();
			await settle();

			expect(api).toHaveBeenCalledWith(
				'/api/threads/TH-0001/confirm',
				expect.objectContaining({ method: 'POST' })
			);
			expect(document.querySelector('article header')!.textContent).toContain('answered');
			expect(document.querySelector('[data-pending]')).toBeNull();
			expect(confirmButton()).toBeNull();
		});

		it('offers no Confirm on a read-only dashboard, and still marks the recommendation', async () => {
			api.mockResolvedValue({ ok: true, json: async () => [pending] });
			c = mount(Threads, { target: document.body, props: { on: 'S-0001', writable: false } });
			await settle();

			expect(confirmButton()).toBeNull();
			expect(document.querySelector('[data-recommendation]')).not.toBeNull();
			expect(document.querySelector('[data-source] a')).not.toBeNull();
		});

		it('says why flai refused a confirmation', async () => {
			api.mockImplementation(async (_path: string, init?: { method?: string }) =>
				init?.method === 'POST'
					? {
							ok: false,
							statusText: 'Bad Request',
							json: async () => ({ error: 'TH-0001 has no recommendation awaiting confirmation' })
						}
					: { ok: true, json: async () => [pending] }
			);
			c = mount(Threads, { target: document.body, props: { on: 'S-0001' } });
			await settle();
			confirmButton()!.click();
			await settle();

			expect(document.querySelector('section')!.textContent).toContain(
				'refused: TH-0001 has no recommendation awaiting confirmation'
			);
		});
	});
});
