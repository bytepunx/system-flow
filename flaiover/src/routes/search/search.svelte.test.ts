// S-0201: a search hit for a draft story says so in words, beside its status.
import { afterEach, describe, expect, it, vi } from 'vitest';
import { flushSync, mount, unmount } from 'svelte';

vi.mock('$app/paths', () => ({
	resolve: (route: string, params: Record<string, string>) =>
		route.replace('[...path]', params.path)
}));

const answer = (body: unknown) => ({ ok: true, json: async () => body });
const hit = { kind: 'item', scope: 'wip', status: 'backlog', type: 'story', nature: 'feature' };

describe('the search page', () => {
	let c: ReturnType<typeof mount> | undefined;
	afterEach(() => {
		if (c) unmount(c);
		c = undefined;
		vi.unstubAllGlobals();
		vi.useRealTimers();
		document.body.innerHTML = '';
	});

	it('marks a draft story’s hit Draft, and no other hit', async () => {
		vi.useFakeTimers();
		vi.stubGlobal(
			'fetch',
			vi.fn(async () =>
				answer({
					query: 'cranes',
					indexed: 3,
					hits: [
						{
							...hit,
							path: 'wip/kanban/stories/S-0002.md',
							itemId: 'S-0002',
							title: 'Cranes',
							draft: true
						},
						{ ...hit, path: 'wip/kanban/stories/S-0001.md', itemId: 'S-0001', title: 'Berths' },
						{
							...hit,
							path: 'wip/kanban/tasks/T-0001.md',
							itemId: 'T-0001',
							title: 'Crane',
							type: 'task',
							draft: true
						}
					].map((h) => ({ ...h, snippet: '', route: `/items/${h.itemId}` }))
				})
			)
		);
		const { default: Search } = await import('./+page.svelte');
		c = mount(Search, { target: document.body });
		const input = document.querySelector('input')!;
		input.value = 'cranes';
		input.dispatchEvent(new Event('input', { bubbles: true }));
		await vi.runAllTimersAsync();
		flushSync();
		const rows = [...document.querySelectorAll('li[data-kind="item"]')];
		expect(rows).toHaveLength(3);
		const marks = rows.map(
			(r) => r.querySelector('[data-testid="hit-draft"]')?.textContent ?? null
		);
		expect(marks).toEqual(['Draft', null, null]);
		const mark = rows[0].querySelector<HTMLElement>('[data-testid="hit-draft"]')!;
		expect(mark.className).toContain('text-warn');
		expect(mark.parentElement!.textContent!.replace(/\s+/g, ' ')).toContain('backlog·Draft');
	});
});
