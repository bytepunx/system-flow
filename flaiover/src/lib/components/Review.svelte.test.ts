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
}) {
	api.mockImplementation(async (url: string, init?: { method?: string }) => {
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
		if (url.startsWith('/api/docs/file')) return json({ body: NARRATIVE });
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
		expect(text).toContain('1. Accept.');
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
});
