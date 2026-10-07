import { afterEach, describe, expect, it, vi } from 'vitest';
import { flushSync, mount, unmount } from 'svelte';
import ActivityView from './ActivityView.svelte';
import InboxView from './InboxView.svelte';
import InboxBadge from './InboxBadge.svelte';
import { inboxState, newEntries, type Inbox } from '$lib/inbox.svelte';

const api = vi.fn();
vi.mock('$lib/api', () => ({ api: (...args: unknown[]) => api(...args) }));
vi.mock('$app/paths', () => ({
	resolve: (route: string, params: Record<string, string>) =>
		route.replace('[...path]', params.path ?? '').replace('[id]', params.id ?? '')
}));
const json = (body: unknown, status = 200) => ({
	ok: status < 400,
	status,
	statusText: String(status),
	json: async () => body
});

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
		api.mockReset();
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
		// when the thread changed, in the local zone (S-0329): the tests run in New York
		expect(document.body.textContent).toContain('2026-09-19 00:00 EDT');
		expect(document.body.textContent).not.toContain('2026-09-19T04:00:00Z');
	});

	const QUESTION_BOX: Inbox = {
		total: 1,
		counts: { thread: 0, question: 1, review: 0, blocked: 0, overlap: 0 },
		notes: [],
		entries: [
			{
				key: 'question:S-0001:abc',
				kind: 'question',
				title: 'Which port should it use?',
				detail: 'asked in the narrative of S-0001',
				href: '/items/S-0001?question=question%3AS-0001%3Aabc',
				item: 'S-0001'
			}
		]
	};

	it('offers no way to answer a question when the dashboard is not writable', () => {
		c = mount(InboxView, {
			target: document.body,
			props: { inbox: QUESTION_BOX, writable: false }
		});
		flushSync();
		expect(document.querySelector('[data-testid="question-answer-input"]')).toBeNull();
	});

	it('answers an open question in place, posting to its story and refreshing the inbox', async () => {
		api.mockResolvedValueOnce(json({}));
		api.mockResolvedValueOnce(json({ total: 0, counts: BOX.counts, entries: [], notes: [] }));
		c = mount(InboxView, { target: document.body, props: { inbox: QUESTION_BOX, writable: true } });
		flushSync();
		const input = document.querySelector<HTMLInputElement>(
			'[data-testid="question-answer-input"]'
		)!;
		input.value = 'Nine.';
		input.dispatchEvent(new Event('input', { bubbles: true }));
		flushSync();
		document.querySelector<HTMLButtonElement>('[data-testid="question-answer-submit"]')!.click();
		for (let i = 0; i < 5; i++) await Promise.resolve();
		expect(api).toHaveBeenCalledWith(
			'/api/streams/S-0001/answer',
			expect.objectContaining({
				method: 'POST',
				body: JSON.stringify({ question: 'Which port should it use?', answer: 'Nine.' })
			})
		);
		expect(api).toHaveBeenLastCalledWith('/api/inbox');
	});

	it('shows what flai said when answering is refused', async () => {
		api.mockResolvedValueOnce(json({ error: 'no open question matching that' }, 400));
		c = mount(InboxView, { target: document.body, props: { inbox: QUESTION_BOX, writable: true } });
		flushSync();
		const input = document.querySelector<HTMLInputElement>(
			'[data-testid="question-answer-input"]'
		)!;
		input.value = 'Nine.';
		input.dispatchEvent(new Event('input', { bubbles: true }));
		flushSync();
		document.querySelector<HTMLButtonElement>('[data-testid="question-answer-submit"]')!.click();
		for (let i = 0; i < 5; i++) await Promise.resolve();
		flushSync();
		expect(document.body.textContent).toContain('no open question matching that');
	});

	// ADR-0090: a thread whose recommendation awaits the designer is confirmed in place.
	const RECOMMENDATION_BOX: Inbox = {
		total: 1,
		counts: { thread: 1, question: 0, review: 0, blocked: 0, overlap: 0 },
		notes: [],
		entries: [
			{
				key: 'thread:TH-0009',
				kind: 'thread',
				title: 'Which port?',
				detail: 'orchestrator recommends an answer to confirm, on S-0042',
				href: '/items/S-0042?thread=TH-0009',
				item: 'S-0042',
				recommendation: {
					author: 'orchestrator',
					at: '2026-10-01T09:00:00Z',
					text: 'Use 8080.',
					source: { path: 'design/system/overview.md', heading: 'Delivery sequence' }
				}
			}
		]
	};
	const confirmButton = () => document.querySelector<HTMLButtonElement>('button[data-confirm]');

	it('shows a pending recommendation with its text and source, linking to the thread', () => {
		c = mount(InboxView, {
			target: document.body,
			props: { inbox: RECOMMENDATION_BOX, writable: true }
		});
		flushSync();
		const rec = document.querySelector('[data-recommendation="thread:TH-0009"]')!;
		expect(rec.textContent).toContain('recommendation');
		expect(rec.textContent).toContain('orchestrator');
		expect(rec.textContent).toContain('Use 8080.');
		const source = rec.querySelector('[data-source] a')!;
		expect(source.getAttribute('href')).toBe('/docs/design/system/overview.md#delivery-sequence');
		expect(source.textContent).toBe('design/system/overview.md § Delivery sequence');
		expect(document.querySelector('a')!.getAttribute('href')).toBe('/items/S-0042?thread=TH-0009');
		expect(confirmButton()!.textContent).toBe('Confirm');
	});

	it('offers no Confirm on a read-only dashboard', () => {
		c = mount(InboxView, {
			target: document.body,
			props: { inbox: RECOMMENDATION_BOX, writable: false }
		});
		flushSync();
		expect(document.querySelector('[data-recommendation]')).not.toBeNull();
		expect(confirmButton()).toBeNull();
	});

	it('confirms in place, posting to the thread and refreshing the inbox', async () => {
		api.mockResolvedValueOnce(json({ id: 'TH-0009', status: 'answered' }));
		api.mockResolvedValueOnce(json({ total: 0, counts: BOX.counts, entries: [], notes: [] }));
		c = mount(InboxView, {
			target: document.body,
			props: { inbox: RECOMMENDATION_BOX, writable: true }
		});
		flushSync();
		confirmButton()!.click();
		for (let i = 0; i < 5; i++) await Promise.resolve();
		expect(api).toHaveBeenCalledWith(
			'/api/threads/TH-0009/confirm',
			expect.objectContaining({ method: 'POST' })
		);
		expect(api).toHaveBeenLastCalledWith('/api/inbox');
	});

	it('shows what flai said when confirming is refused', async () => {
		api.mockResolvedValueOnce(json({ error: 'TH-0009 has no recommendation awaiting' }, 400));
		c = mount(InboxView, {
			target: document.body,
			props: { inbox: RECOMMENDATION_BOX, writable: true }
		});
		flushSync();
		confirmButton()!.click();
		for (let i = 0; i < 5; i++) await Promise.resolve();
		flushSync();
		expect(document.querySelector('[role="alert"]')!.textContent).toBe(
			'TH-0009 has no recommendation awaiting'
		);
		expect(confirmButton()!.disabled).toBe(false);
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
