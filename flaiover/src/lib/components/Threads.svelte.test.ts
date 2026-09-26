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
});
