import { afterEach, describe, expect, it, vi } from 'vitest';
import { flushSync, mount, unmount } from 'svelte';
import Threads from './Threads.svelte';

const api = vi.fn();
vi.mock('$lib/api', () => ({ api: (...args: unknown[]) => api(...args) }));
vi.mock('$app/paths', () => ({
	resolve: (route: string, params: Record<string, string>) => route.replace('[id]', params.id ?? '')
}));

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

describe('Threads', () => {
	let c: ReturnType<typeof mount> | undefined;
	afterEach(() => {
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

	it('shows no pager for a single thread (S-0133)', async () => {
		api.mockResolvedValue({ ok: true, json: async () => [numbered(1)] });
		c = mount(Threads, { target: document.body, props: { on: 'S-0001' } });
		await settle();

		expect(document.querySelectorAll('article')).toHaveLength(1);
		expect(pagers()).toHaveLength(0);
	});
});
