import { afterEach, describe, expect, it, vi } from 'vitest';
import { flushSync, mount, unmount } from 'svelte';
import ActivityView from './ActivityView.svelte';
import InboxView from './InboxView.svelte';
import InboxBadge from './InboxBadge.svelte';
import { inboxState, newEntries, type Inbox } from '$lib/inbox.svelte';

vi.mock('$lib/api', () => ({ api: vi.fn() }));
vi.mock('$app/paths', () => ({
	resolve: (route: string, params: Record<string, string>) =>
		route.replace('[...path]', params.path ?? '').replace('[id]', params.id ?? '')
}));

const BOX: Inbox = {
	total: 3,
	counts: { thread: 1, question: 0, review: 1, blocked: 1, overlap: 0 },
	notes: ['Overlapping touches are not listed: flai is not available to this dashboard.'],
	entries: [
		{
			key: 'thread:TH-0004',
			kind: 'thread',
			title: 'Which port?',
			detail: 'claude wrote last, on S-0042',
			href: '/items/S-0042',
			at: '2026-09-19T04:00:00Z'
		},
		{
			key: 'review:S-0041',
			kind: 'review',
			title: 'S-0041 Review and acceptance',
			detail: 'in review: accept it or send it back',
			href: '/review/S-0041'
		},
		{
			key: 'blocked:T-0003:x',
			kind: 'blocked',
			title: 'T-0003 Wire it',
			detail: 'blocked: waiting on keys',
			href: '/items/T-0003'
		}
	]
};

describe('inbox and activity views', () => {
	let c: ReturnType<typeof mount> | undefined;
	afterEach(() => {
		if (c) unmount(c);
		c = undefined;
		inboxState.data = null;
		document.body.innerHTML = '';
	});

	it('groups inbox entries by kind, reviews first, each linking to its page', () => {
		c = mount(InboxView, { target: document.body, props: { inbox: BOX } });
		flushSync();
		const headings = [...document.querySelectorAll('h2')].map((h) =>
			h.textContent?.replace(/\s+/g, ' ').trim()
		);
		expect(headings).toEqual(['Stories in review 1', 'Threads awaiting you 1', 'Blocked items 1']);
		const links = [...document.querySelectorAll('a')].map((a) => a.getAttribute('href'));
		expect(links).toEqual(['/review/S-0041', '/items/S-0042', '/items/T-0003']);
		expect(document.body.textContent).toContain('blocked: waiting on keys');
		expect(document.body.textContent).toContain('flai is not available');
	});

	it('says so when nothing needs the designer', () => {
		c = mount(InboxView, {
			target: document.body,
			props: { inbox: { total: 0, counts: BOX.counts, entries: [], notes: [] } }
		});
		flushSync();
		expect(document.body.textContent).toContain('Nothing needs you right now.');
	});

	it('shows the count beside Inbox, and nothing when the inbox is empty', () => {
		c = mount(InboxBadge, { target: document.body });
		flushSync();
		expect(document.querySelector('[data-testid="inbox-count"]')).toBeNull();
		inboxState.data = BOX;
		flushSync();
		const badge = document.querySelector('[data-testid="inbox-count"]')!;
		expect(badge.textContent).toBe('3');
		expect(badge.getAttribute('aria-label')).toBe('3 in the inbox');
	});

	it('shows each stream with agent, session, age, task, blocked flag, and last log entry', () => {
		c = mount(ActivityView, {
			target: document.body,
			props: {
				streams: [
					{
						stream: 'S-0042',
						title: 'Agent presence',
						agent: 'system-flow',
						session: '5f268c00',
						updated: '2026-09-19T04:48:45Z',
						age_seconds: 7200,
						status: 'in-progress',
						blocked: true,
						task: { id: 'T-0193', title: 'Activity and inbox pages' },
						last_log: { at: '2026-09-19T04:48:45Z', text: 'T-0192 done.' },
						path: 'wip/agents/S-0042.md'
					}
				]
			}
		});
		flushSync();
		const text = (document.body.textContent ?? '').replace(/\s+/g, ' ');
		for (const want of [
			'S-0042',
			'system-flow',
			'session 5f268c00',
			'2h ago',
			'T-0193',
			'BLOCKED',
			'T-0192 done.'
		])
			expect(text).toContain(want);
		expect([...document.querySelectorAll('a')].map((a) => a.getAttribute('href'))).toContain(
			'/docs/wip/agents/S-0042.md'
		);
	});

	it('says when there are no active streams', () => {
		c = mount(ActivityView, { target: document.body, props: { streams: [] } });
		flushSync();
		expect(document.body.textContent).toContain('No active streams');
	});
});

describe('new inbox entries', () => {
	it('is nothing on the first look, then only keys not seen before', () => {
		expect(newEntries(null, BOX.entries)).toEqual([]);
		const known = new Set(['thread:TH-0004', 'review:S-0041']);
		expect(newEntries(known, BOX.entries).map((e) => e.key)).toEqual(['blocked:T-0003:x']);
		expect(newEntries(new Set(BOX.entries.map((e) => e.key)), BOX.entries)).toEqual([]);
	});
});
