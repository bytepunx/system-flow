import { afterEach, describe, expect, it, vi } from 'vitest';
import { flushSync, mount, unmount } from 'svelte';
import Review from './Review.svelte';

const api = vi.fn();
vi.mock('$lib/api', () => ({ api: (...args: unknown[]) => api(...args) }));
vi.mock('$app/paths', () => ({
	resolve: (route: string, params: Record<string, string>) =>
		route.replace('[...path]', params.path ?? '').replace('[id]', params.id ?? '')
}));

const settle = async () => {
	for (let i = 0; i < 8; i++) await new Promise((r) => setTimeout(r, 0));
	flushSync();
};
const json = (body: unknown, status = 200) => ({
	ok: status < 400,
	status,
	statusText: String(status),
	body: null,
	json: async () => body,
	text: async () => JSON.stringify(body)
});
const ndjson = (lines: unknown[]) => ({
	ok: true,
	status: 200,
	statusText: 'OK',
	body: null,
	json: async () => ({}),
	text: async () => lines.map((l) => JSON.stringify(l)).join('\n')
});

const STORY = {
	id: 'S-0041',
	type: 'story',
	title: 'Review and acceptance',
	status: 'review',
	path: 'wip/kanban/stories/S-0041-x.md',
	body: '# S-0041\n\n## Acceptance criteria\n- [x] A review page\n- [ ] Accept as the designer\n\n## Notes\n'
};
const NARRATIVE =
	'# S-0041\n\n## Current state\nAll tasks done.\n\n## Next steps\n1. Accept.\n\n## Decisions\n';
const DIFF = {
	branch: 'story/S-0041',
	base: 'abc1234',
	commits: 1,
	additions: 2,
	deletions: 1,
	truncated: false,
	files: [
		{
			path: 'docs/a.md',
			status: 'modified',
			additions: 2,
			deletions: 1,
			binary: false,
			truncated: false,
			patch: '@@ -1,2 +1,3 @@\n one\n-two\n+2\n+three'
		}
	]
};

function backend(over: {
	item?: Record<string, unknown>;
	preview?: unknown;
	accept?: () => unknown;
	move?: () => unknown;
	diff?: () => unknown;
	checksStatuses?: unknown[];
	checksRun?: () => unknown;
	checksCancel?: () => unknown;
	checksTail?: () => unknown;
	narrative?: string;
}) {
	let statusIdx = 0;
	api.mockImplementation(async (url: string, init?: { method?: string; body?: string }) => {
		if (url.includes('/checks/tail'))
			return over.checksTail
				? over.checksTail()
				: ndjson([{ event: 'done', offset: 0, running: false }]);
		if (url.endsWith('/checks') && init?.method === 'POST') {
			const body = JSON.parse(init.body ?? '{}') as { action?: string };
			if (body.action === 'run')
				return over.checksRun ? over.checksRun() : json({ story: 'S-0041', running: false });
			return over.checksCancel ? over.checksCancel() : json({ story: 'S-0041', running: false });
		}
		if (url.endsWith('/checks')) {
			const list = over.checksStatuses ?? [
				{ story: 'S-0041', running: false, checks_enabled: false }
			];
			const v = list[Math.min(statusIdx, list.length - 1)];
			statusIdx++;
			return json(v);
		}
		if (url.endsWith('/accept'))
			return over.accept ? over.accept() : ndjson([{ event: 'done', result: {} }]);
		if (url.endsWith('/move') && init?.method === 'POST') return over.move ? over.move() : json({});
		if (url.endsWith('/diff')) return over.diff ? over.diff() : json(DIFF);
		if (url.endsWith('/acceptance'))
			return json(
				over.preview ?? {
					branch: 'story/S-0041',
					plan: {
						level: 'minor',
						commits: ['a'],
						steps: [
							{
								component: { name: 'flaiover' },
								delivered: true,
								level: 'minor',
								from: '0.12.1',
								to: '0.13.0',
								files: ['x']
							}
						]
					}
				}
			);
		if (url.startsWith('/api/board')) return json({ writable: true });
		if (url.startsWith('/api/docs/file')) return json({ body: over.narrative ?? NARRATIVE });
		if (url.startsWith('/api/threads')) return json([]);
		if (url.startsWith('/api/items/')) return json({ item: over.item ?? STORY, children: [] });
		return json({});
	});
}
const button = (label: string) =>
	[...document.querySelectorAll('button')].find((b) =>
		(b.textContent ?? '').trim().startsWith(label)
	)!;
const calls = (suffix: string) =>
	api.mock.calls.filter(
		(c) => String(c[0]).endsWith(suffix) && (c[1] as { method?: string })?.method === 'POST'
	);

describe('Review', () => {
	let c: ReturnType<typeof mount> | undefined;
	afterEach(() => {
		if (c) unmount(c);
		c = undefined;
		api.mockReset();
		document.body.innerHTML = '';
	});

	it('puts criteria, narrative, diff, and the release plan on one page', async () => {
		backend({});
		c = mount(Review, { target: document.body, props: { id: 'S-0041' } });
		await settle();
		const text = document.body.textContent ?? '';
		expect(text).toContain('1 of 2 checked');
		expect(text).toContain('not checked: Accept as the designer');
		expect(text).toContain('All tasks done.');
		expect(text).toContain('Accept.');
		expect(text).toContain('story/S-0041');
		expect(text).toContain('docs/a.md');
		expect(text).toContain('0.12.1 → 0.13.0');
		// hunks on demand, with + and - kept
		expect(text).not.toContain('+three');
		[...document.querySelectorAll('button')]
			.find((b) => b.textContent?.includes('docs/a.md'))!
			.click();
		flushSync();
		expect(document.body.textContent).toContain('+three');
		expect(document.body.textContent).toContain('-two');
	});

	it('accepts as a stream: each step, then the tags, and says it was not pushed', async () => {
		backend({
			accept: () =>
				ndjson([
					{ event: 'progress', step: 'merged', msg: 'story/S-0041 rebased and fast-forwarded' },
					{ event: 'progress', step: 'committed', msg: 'chore: [S-0041] accept and archive' },
					{ event: 'warning', msg: 'run: git push origin HEAD flaiover/v0.13.0' },
					{ event: 'done', result: { status: 'done', tags: ['flaiover/v0.13.0'], pushed: false } }
				])
		});
		c = mount(Review, { target: document.body, props: { id: 'S-0041' } });
		await settle();
		button('Accept').click();
		await settle();
		const text = document.body.textContent ?? '';
		expect(calls('/accept')).toHaveLength(1);
		expect(text).toContain('story/S-0041 rebased and fast-forwarded');
		expect(text).toContain('chore: [S-0041] accept and archive');
		expect(text).toContain('S-0041 is accepted');
		expect(text).toContain('flaiover/v0.13.0');
		expect(text).toContain('Not pushed');
		expect(button('Accept')).toBeUndefined();
	});

	it('shows flai’s failure verbatim and that the story is still in review', async () => {
		const message = 'rebase of story/S-0041 onto main stopped with conflicts in docs/users/flai.md';
		backend({
			accept: () =>
				ndjson([
					{ event: 'progress', step: 'merged', msg: 'x' },
					{ event: 'error', status: 500, error: message }
				])
		});
		c = mount(Review, { target: document.body, props: { id: 'S-0041' } });
		await settle();
		button('Accept').click();
		await settle();
		const alert = [...document.querySelectorAll('[role=alert]')]
			.map((a) => a.textContent)
			.join(' ');
		expect(alert).toContain(message);
		expect(alert).toContain('S-0041 is review');
		expect(document.body.textContent).not.toContain('is accepted');
	});

	it('keeps accept disabled for blockers and until uncommitted files are chosen', async () => {
		backend({ preview: { branch: 'story/S-0041', plan: null, uncommitted: ['docs/stray.md'] } });
		c = mount(Review, { target: document.body, props: { id: 'S-0041' } });
		await settle();
		expect(button('Accept').disabled).toBe(true);
		(document.querySelector('input[type=checkbox]') as HTMLInputElement).click();
		flushSync();
		expect(button('Accept').disabled).toBe(false);
		button('Accept').click();
		await settle();
		expect(JSON.parse(calls('/accept')[0][1].body)).toEqual({ include_uncommitted: true });
	});

	it('offers a research story for acceptance, with no release and what lands unreleased (ADR-0025)', async () => {
		backend({
			preview: {
				branch: 'story/S-0041',
				blockers: [],
				plan: {
					level: 'none',
					commits: ['a'],
					steps: [],
					skipped: 'S-0041 is research: its findings land on main and are pushed',
					unreleased: [{ component: 'flai', files: ['flai/x.go'] }]
				}
			}
		});
		c = mount(Review, { target: document.body, props: { id: 'S-0041' } });
		await settle();
		const text = (document.body.textContent ?? '').replace(/\s+/g, ' ');
		expect(text).toContain('No release: S-0041 is research');
		expect(text).toContain('flai lands on main without a release (1 file)');
		expect(button('Accept').disabled).toBe(false);
	});

	it('sends a story back with a reason typed in the page', async () => {
		backend({});
		c = mount(Review, { target: document.body, props: { id: 'S-0041' } });
		await settle();
		button('Send back').click();
		flushSync();
		expect(button('Send back to in-progress').disabled).toBe(true);
		const box = document.querySelector('textarea#send-back-reason') as HTMLTextAreaElement;
		box.value = 'the diff misses the docs';
		box.dispatchEvent(new Event('input', { bubbles: true }));
		flushSync();
		button('Send back to in-progress').click();
		await settle();
		expect(JSON.parse(calls('/move')[0][1].body)).toEqual({
			to: 'in-progress',
			reason: 'the diff misses the docs'
		});
	});

	it('says so when the story is not in review, and offers nothing to accept', async () => {
		backend({ item: { ...STORY, status: 'in-progress' } });
		c = mount(Review, { target: document.body, props: { id: 'S-0041' } });
		await settle();
		expect(document.body.textContent).toContain('is in-progress, not in review');
		expect(button('Accept')).toBeUndefined();
	});

	it('shows why there is no diff when the story has no branch', async () => {
		backend({ diff: () => json({ error: 'S-0041 has no branch story/S-0041' }, 500) });
		c = mount(Review, { target: document.body, props: { id: 'S-0041' } });
		await settle();
		expect(document.body.textContent).toContain('has no branch story/S-0041');
	});

	describe('checks (S-0082)', () => {
		it('says checks are off and what enables them, and offers no button', async () => {
			backend({ checksStatuses: [{ story: 'S-0041', running: false, checks_enabled: false }] });
			c = mount(Review, { target: document.body, props: { id: 'S-0041' } });
			await settle();
			const section = document.querySelector('[data-testid="checks-section"]')!;
			expect(section.textContent).toContain('flai serve enable checks');
			expect(section.querySelector('[data-testid="checks-run"]')).toBeNull();
		});

		it('shows a run that finished before this page was opened, not nothing', async () => {
			backend({
				checksStatuses: [
					{
						story: 'S-0041',
						running: false,
						checks_enabled: true,
						outcome: 'passed',
						started: '2026-09-22T00:00:00Z',
						ended: '2026-09-22T00:05:00Z',
						steps: [{ name: 'flai', command: 'scripts/flai-test.sh', exit: 0 }]
					}
				]
			});
			c = mount(Review, { target: document.body, props: { id: 'S-0041' } });
			await settle();
			const text = document.querySelector('[data-testid="checks-section"]')!.textContent!;
			expect(text).toContain('passed');
			expect(text).toContain('flai: ok');
		});

		// A live checksRun state, held here rather than in a fixed backend() array: watchChecks
		// polls checks.status in a loop of its own once a run is found active, at a pace this test
		// does not control, so the mock must answer from real state (flipped by the actions under
		// test) rather than a sequence indexed by call count. Everything else still goes through
		// the ordinary backend().
		function liveChecksBackend(initial: Record<string, unknown> = {}) {
			backend({});
			const base = api.getMockImplementation()!;
			let run: Record<string, unknown> = {
				story: 'S-0041',
				running: false,
				checks_enabled: true,
				...initial
			};
			let tailSent = false;
			api.mockImplementation(async (url: string, init?: { method?: string; body?: string }) => {
				if (url.includes('/checks/tail')) {
					// A real flai checks.tail blocks server-side, for real time; a mock that answers
					// on the spot lets watchChecks's loop spin on nothing but resolved microtasks,
					// which starves the event loop before it ever runs settle()'s own timers. One
					// real tick keeps this test from hanging (found live, running this very test).
					await new Promise((r) => setTimeout(r, 0));
					if (!tailSent && run.running) {
						tailSent = true;
						return ndjson([
							{ event: 'line', text: '=== flai: scripts/flai-test.sh ===' },
							{ event: 'done', offset: 40, running: run.running }
						]);
					}
					return ndjson([{ event: 'done', offset: 40, running: run.running }]);
				}
				if (url.endsWith('/checks') && init?.method === 'POST') {
					const body = JSON.parse(init.body ?? '{}') as { action?: string };
					if (body.action === 'run') {
						run = {
							story: 'S-0041',
							running: true,
							checks_enabled: true,
							current: 'flai',
							started: '2026-09-22T00:00:00Z'
						};
						tailSent = false;
					} else {
						run = { ...run, running: false, outcome: 'cancelled' };
					}
					return json(run);
				}
				if (url.endsWith('/checks')) return json(run);
				return base(url, init);
			});
			return {
				end: (outcome: string, steps: unknown[]) => {
					run = { ...run, running: false, outcome, ended: '2026-09-22T00:01:00Z' };
					if (steps.length) run.steps = steps;
				}
			};
		}

		it('runs checks: live output while running, the outcome once the run ends', async () => {
			const live = liveChecksBackend();
			c = mount(Review, { target: document.body, props: { id: 'S-0041' } });
			await settle();
			document.querySelector<HTMLButtonElement>('[data-testid="checks-run"]')!.click();
			await settle();
			expect(document.querySelector('[data-testid="checks-section"]')!.textContent).toContain(
				'scripts/flai-test.sh'
			);
			live.end('passed', [{ name: 'flai', command: 'scripts/flai-test.sh', exit: 0 }]);
			await settle();
			await settle();
			expect(calls('/checks')[0][1]).toMatchObject({
				method: 'POST',
				body: JSON.stringify({ action: 'run' })
			});
			const text = document.querySelector('[data-testid="checks-section"]')!.textContent!;
			expect(text).toContain('scripts/flai-test.sh');
			expect(text).toContain('passed');
		});

		it('cancels a running run', async () => {
			liveChecksBackend({ running: true, current: 'flai', started: '2026-09-22T00:00:00Z' });
			c = mount(Review, { target: document.body, props: { id: 'S-0041' } });
			await settle();
			expect(document.querySelector('[data-testid="checks-run"]')).toBeNull();
			document.querySelector<HTMLButtonElement>('[data-testid="checks-cancel"]')!.click();
			await settle();
			expect(calls('/checks')[0][1]).toMatchObject({
				method: 'POST',
				body: JSON.stringify({ action: 'cancel' })
			});
			expect(document.querySelector('[data-testid="checks-section"]')!.textContent).toContain(
				'cancelled'
			);
		});

		it('shows what flai said when a run could not be started', async () => {
			backend({
				checksStatuses: [{ story: 'S-0041', running: false, checks_enabled: true }],
				checksRun: () => json({ error: 'a checks run for S-0041 is already active (pid 123)' }, 409)
			});
			c = mount(Review, { target: document.body, props: { id: 'S-0041' } });
			await settle();
			document.querySelector<HTMLButtonElement>('[data-testid="checks-run"]')!.click();
			await settle();
			expect(document.querySelector('[data-testid="checks-error"]')!.textContent).toContain(
				'already active'
			);
		});
	});

	// Found live (S-0088): the narrative pane showed markdown as literal text ("1. Accept.", "**x**")
	// instead of rendering it, unlike every other place in the dashboard that shows a document body.
	it('renders markdown in the narrative pane rather than showing it as literal text', async () => {
		backend({
			narrative:
				'# S-0041\n\n## Current state\n**All** tasks done.\n\n## Next steps\n1. Accept.\n2. Then archive.\n\n## Decisions\n'
		});
		c = mount(Review, { target: document.body, props: { id: 'S-0041' } });
		await settle();
		const strong = document.querySelector('strong');
		expect(strong?.textContent).toBe('All');
		const items = [...document.querySelectorAll('ol li')].map((li) => li.textContent);
		expect(items).toEqual(['Accept.', 'Then archive.']);
		expect(document.body.textContent).not.toContain('**All**');
	});
});
