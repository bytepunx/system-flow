// S-0215: the table under the agent-waiting chart lists the longest waits of the window with their
// item, the thread waited on or review, when they started and ended, and who was awaited.
import { afterEach, describe, expect, it, vi } from 'vitest';
import { flushSync, mount, unmount } from 'svelte';
import WaitTable from './WaitTable.svelte';
import type { LongestWait } from '$lib/viz/charts';

vi.mock('$app/paths', () => ({
	resolve: (route: string, params: Record<string, string>) => route.replace('[id]', params.id)
}));

const rows: LongestWait[] = [
	{
		item: 'S-0004',
		kind: 'thread',
		thread: 'TH-0012',
		started: '2026-09-28T09:15:00Z',
		ended: '2026-09-28T13:45:00Z',
		seconds: 16200,
		awaited: 'alex'
	},
	{
		item: 'S-0007',
		kind: 'review',
		started: '2026-09-29T10:00:00Z',
		ended: '2026-09-29T12:00:00Z',
		seconds: 7200,
		awaited: 'orchestrator'
	},
	{ item: 'S-0009', kind: 'review', started: '2026-09-29T20:30:00Z', seconds: 1800 }
];

describe('WaitTable', () => {
	let component: ReturnType<typeof mount> | undefined;
	afterEach(() => {
		if (component) unmount(component);
		component = undefined;
		document.body.innerHTML = '';
	});
	const show = (props: { rows: LongestWait[]; type: string }) => {
		component = mount(WaitTable, { target: document.body, props });
		flushSync();
	};
	const cells = () =>
		[...document.querySelectorAll('[data-testid="wait-table"] tbody tr')].map((tr) =>
			[...tr.querySelectorAll('td')].map((td) => td.textContent?.trim())
		);
	const links = (row: number) =>
		[
			...document.querySelectorAll<HTMLAnchorElement>('[data-testid="wait-table"] tbody tr')[row]
				.querySelectorAll('a')
		].map((a) => a.getAttribute('href'));

	it("lists a thread's wait with its item and its thread linked, and who answered it", () => {
		show({ rows, type: 'story' });
		const heads = [...document.querySelectorAll('th')].map((th) => th.textContent?.trim());
		expect(heads).toEqual(['story', 'kind', 'thread', 'started', 'ended', 'waited', 'awaited']);
		expect(cells()[0]).toEqual([
			'S-0004',
			'thread',
			'TH-0012',
			'2026-09-28 09:15',
			'2026-09-28 13:45',
			'4.5h',
			'alex'
		]);
		expect(links(0)).toEqual(['/items/S-0004', '/items/S-0004?thread=TH-0012']);
	});

	it('lists a spell in review as review, awaited by who moved the item on', () => {
		show({ rows, type: 'story' });
		expect(cells()[1]).toEqual([
			'S-0007',
			'review',
			'review',
			'2026-09-29 10:00',
			'2026-09-29 12:00',
			'2h',
			'orchestrator'
		]);
		expect(links(1)).toEqual(['/items/S-0007']);
	});

	it('says an open wait is open, with no one awaited yet', () => {
		show({ rows, type: 'task' });
		expect(document.querySelector('th')?.textContent).toBe('task');
		expect(cells()[2]).toEqual([
			'S-0009',
			'review',
			'review',
			'2026-09-29 20:30',
			'open',
			'30m',
			'—'
		]);
	});

	it('says so when no wait has a part in the window', () => {
		show({ rows: [], type: 'story' });
		expect(document.querySelector('[data-testid="wait-table"]')).toBeNull();
		expect(document.querySelector('[data-testid="wait-none"]')?.textContent).toBe(
			'No wait has a part in the window.'
		);
	});
});
