import { afterEach, describe, expect, it, vi } from 'vitest';
import { flushSync, mount, unmount } from 'svelte';
import Review from './Review.svelte';

const api = vi.fn();
vi.mock('$lib/api', () => ({ api: (...args: unknown[]) => api(...args) }));
const openEditor = vi.fn();
vi.mock('$lib/editor.svelte', () => ({
	openEditor: (...args: unknown[]) => openEditor(...args)
}));
vi.mock('$app/paths', () => ({
	resolve: (route: string, params: Record<string, string> = {}) =>
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
	/** The acceptance preview, or a function giving it afresh on each read. */
	preview?: unknown;
	accept?: () => unknown;
	restore?: () => unknown;
	thread?: () => unknown;
	move?: () => unknown;
	diff?: () => unknown;
	checksStatuses?: unknown[];
	checksRun?: () => unknown;
	checksCancel?: () => unknown;
	checksTail?: () => unknown;
	narrative?: string;
	agent?: () => unknown;
	issues?: unknown[];
	ownIssues?: unknown[];
	issueStory?: (id: string) => unknown;
	verify?: () => unknown;
	/** The story's worktree catalogue, GET /api/components?story=. */
	components?: () => unknown;
	/** The story's pull request as flai serve kept it, GET /api/items/[id]/forge. */
	forge?: () => unknown;
}) {
	let statusIdx = 0;
	api.mockImplementation(async (url: string, init?: { method?: string; body?: string }) => {
		const made = /^\/api\/issues\/([^/]+)\/story$/.exec(url);
		if (made && init?.method === 'POST')
			return over.issueStory
				? over.issueStory(made[1])
				: json({ id: 'S-0099', title: 'x', issue: made[1], committed: true });
		if (url.startsWith('/api/issues?story=')) return json(over.ownIssues ?? []);
		if (url === '/api/issues') return json(over.issues ?? []);
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
		if (url.startsWith('/api/components?story='))
			return over.components ? over.components() : json({ error: 'no components' }, 502);
		if (url.endsWith('/verify')) return over.verify ? over.verify() : json({ report: null });
		if (url.endsWith('/forge'))
			return over.forge ? over.forge() : json({ story: 'S-0041', pull_request: null });
		if (url.endsWith('/accept'))
			return over.accept ? over.accept() : ndjson([{ event: 'done', result: {} }]);
		if (url.endsWith('/move') && init?.method === 'POST') return over.move ? over.move() : json({});
		if (url.endsWith('/restore') && init?.method === 'POST')
			return over.restore ? over.restore() : json({ restored: true });
		if (url === '/api/threads' && init?.method === 'POST')
			return over.thread ? over.thread() : json({ id: 'TH-0001' });
		if (url.endsWith('/diff')) return over.diff ? over.diff() : json(DIFF);
		if (url.endsWith('/agent') && init?.method === 'POST')
			return over.agent ? over.agent() : json({ error: 'no agent' }, 400);
		if (url.endsWith('/acceptance'))
			return json(
				(typeof over.preview === 'function' ? over.preview() : over.preview) ?? {
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
		openEditor.mockReset();
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
		// ADR-0067: accepting pushes nothing, and the page names no push command
		const local = document.querySelector('[data-testid="accept-stays-local"]')!.textContent!;
		expect(local.replace(/\s+/g, ' ')).toContain('stays local until it is published');
		expect(text).not.toContain('flai push');
		// hunks on demand, with + and - kept
		expect(text).not.toContain('+three');
		[...document.querySelectorAll('button')]
			.find((b) => b.textContent?.includes('docs/a.md'))!
			.click();
		flushSync();
		expect(document.body.textContent).toContain('+three');
		expect(document.body.textContent).toContain('-two');
	});

	// S-0140: work left uncommitted in the story's worktree blocks acceptance, and the page offers
	// to have an agent commit it.
	it('offers to have an agent commit what the worktree holds, and keeps Accept off', async () => {
		backend({
			preview: {
				branch: 'story/S-0041',
				worktree_uncommitted: ['docs/left.md'],
				blockers: [
					'the worktree .flai-cache/worktrees/S-0041 has uncommitted changes (docs/left.md); commit them on story/S-0041, or discard them, before accepting'
				],
				plan: null
			},
			agent: () => json({ story: 'S-0041', agent: 'agent-S-0041', pid: 7 })
		});
		c = mount(Review, { target: document.body, props: { id: 'S-0041' } });
		await settle();
		const box = document.querySelector('[data-testid="worktree-uncommitted"]');
		expect(box?.textContent).toContain('docs/left.md');
		expect(document.body.textContent).not.toContain('This cannot be accepted from here yet');
		expect(button('Accept').disabled).toBe(true);
		button('Have an agent commit them').click();
		await settle();
		expect(calls('/agent')).toHaveLength(1);
		expect(JSON.parse(String((calls('/agent')[0][1] as { body: string }).body))).toEqual({
			action: 'commit'
		});
		expect(box?.textContent).toContain('agent-S-0041 is committing them (pid 7)');
	});

	// ADR-0067: acceptance pushes nothing; what is accepted stays local until it is published.
	it('accepts as a stream: each step, then says it stays local until published', async () => {
		backend({
			accept: () =>
				ndjson([
					{ event: 'progress', step: 'merged', msg: 'story/S-0041 rebased and fast-forwarded' },
					{ event: 'progress', step: 'committed', msg: 'chore: [S-0041] accept and archive' },
					{ event: 'done', result: { id: 'S-0041', status: 'done', merged: true, archived: 3 } }
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
		expect(text).toContain('Nothing was released');
		const local = document.querySelector('[data-testid="accept-unpublished"]')!.textContent!;
		expect(local.replace(/\s+/g, ' ')).toContain(
			'stays local until it is published: Publish on the board, or flai release --pending'
		);
		expect(text).not.toContain('flai push');
		expect(api.mock.calls.map(([u]) => u)).not.toContain('/api/unpushed');
		expect(button('Accept')).toBeUndefined();
	});

	// ADR-0093: the message names who accepted, from the acceptance's by, and the commit the
	// orchestrator's verifier passed; an older flai names neither, and the message stays as it was.
	it('names who accepted, and the verified commit when there is one', async () => {
		const accepted = async (result: Record<string, unknown>) => {
			backend({ accept: () => ndjson([{ event: 'done', result }]) });
			c = mount(Review, { target: document.body, props: { id: 'S-0041' } });
			await settle();
			button('Accept').click();
			await settle();
			const text = document.querySelector('[data-testid="accepted"]')!.textContent!;
			unmount(c);
			c = undefined;
			document.body.innerHTML = '';
			return text.replace(/\s+/g, ' ').trim();
		};
		expect(
			await accepted({
				id: 'S-0041',
				status: 'done',
				by: 'orchestrator',
				verified: '0123456789abcdef0123456789abcdef01234567'
			})
		).toBe('S-0041 is accepted by orchestrator at 0123456789ab.');
		expect(await accepted({ id: 'S-0041', status: 'done', by: 'alex' })).toBe(
			'S-0041 is accepted by alex.'
		);
		expect(await accepted({ id: 'S-0041', status: 'done' })).toBe('S-0041 is accepted.');
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

	// ADR-0140, S-0378: a refusal for lint the fix left lists each finding as a link that opens the
	// editor at its line, on the story's branch for the story's files; Accept is on again once each
	// listed file has been saved.
	describe('a refusal for lint the fix left', () => {
		const finding = (over: Record<string, unknown>) => ({
			tier: 'markdown',
			checkout: 'story',
			story: 'S-0041',
			path: 'docs/bad.md',
			line: 2,
			column: 3,
			level: 'error',
			rule: 'MD009',
			alias: 'no-trailing-spaces',
			message: 'Trailing spaces [Expected: 0; Actual: 1]',
			...over
		});
		// The first acceptance is refused with the findings; the next is accepted.
		const refuse = (findings: unknown[]) => {
			const refusal = ndjson([
				{
					event: 'error',
					status: 422,
					error: 'rule: S-0041 cannot be accepted: docs/bad.md:2: Trailing spaces',
					data: { id: 'S-0041', reason: 'S-0041 cannot be accepted', findings }
				}
			]);
			const accepted = ndjson([{ event: 'done', result: { id: 'S-0041', status: 'done' } }]);
			let answers = 0;
			backend({ accept: () => (++answers > 1 ? accepted : refusal) });
		};
		const link = (label: string) =>
			document.querySelector<HTMLButtonElement>(`button[aria-label="${label}"]`)!;
		const words = (e: Element) => (e.textContent ?? '').replace(/\s+/g, ' ').trim();
		const saveLast = () => {
			const target = openEditor.mock.calls.at(-1)![0] as {
				path: string;
				onsaved: (s: unknown) => void;
			};
			target.onsaved({ path: target.path, commit: 'abc1234', committed: true });
			flushSync();
		};

		it('lists each finding with its rule and message, each file a link to the editor at its line', async () => {
			refuse([
				finding({}),
				finding({
					checkout: 'main',
					story: undefined,
					path: 'design/notes.md',
					line: 7,
					rule: 'MD041',
					alias: 'first-line-heading',
					message: 'First line in a file should be a top-level heading'
				})
			]);
			c = mount(Review, { target: document.body, props: { id: 'S-0041' } });
			await settle();
			button('Accept').click();
			await settle();
			const list = document.querySelector('[data-testid="accept-findings"]')!;
			const rows = [...list.querySelectorAll('[data-testid="accept-finding"]')].map(words);
			expect(rows[0]).toContain('docs/bad.md:2');
			expect(rows[0]).toContain('MD009/no-trailing-spaces');
			expect(rows[0]).toContain('Trailing spaces [Expected: 0; Actual: 1]');
			expect(rows[1]).toContain('design/notes.md:7');
			expect(rows[1]).toContain('MD041/first-line-heading');
			expect(link('Edit docs/bad.md at line 2')).toBeTruthy();
			expect(link('Edit design/notes.md at line 7')).toBeTruthy();
			// the refusal's text is not repeated beside its findings
			expect(list.textContent).not.toContain('rule: S-0041');
		});

		it('opens a story’s file on its branch, and a main-checkout file without a story', async () => {
			refuse([
				finding({}),
				finding({ checkout: 'main', story: undefined, path: 'design/notes.md', line: 7 })
			]);
			c = mount(Review, { target: document.body, props: { id: 'S-0041' } });
			await settle();
			button('Accept').click();
			await settle();
			link('Edit docs/bad.md at line 2').click();
			expect(openEditor).toHaveBeenLastCalledWith(
				expect.objectContaining({ path: 'docs/bad.md', line: 2, story: 'S-0041' })
			);
			link('Edit design/notes.md at line 7').click();
			const main = openEditor.mock.calls.at(-1)![0];
			expect(main).toMatchObject({ path: 'design/notes.md', line: 7 });
			expect(main.story).toBeUndefined();
		});

		it('turns Accept on again once each listed file is saved, and accepts from the same page', async () => {
			refuse([
				finding({}),
				finding({
					line: 9,
					rule: 'MD012',
					alias: 'no-multiple-blanks',
					message: 'Multiple blanks'
				}),
				finding({ checkout: 'main', story: undefined, path: 'design/notes.md', line: 7 })
			]);
			c = mount(Review, { target: document.body, props: { id: 'S-0041' } });
			await settle();
			button('Accept').click();
			await settle();
			expect(button('Accept').disabled).toBe(true);
			link('Edit docs/bad.md at line 2').click();
			saveLast();
			expect(document.querySelectorAll('[data-testid="finding-saved"]')).toHaveLength(2);
			expect(button('Accept').disabled).toBe(true);
			link('Edit design/notes.md at line 7').click();
			saveLast();
			expect(button('Accept').disabled).toBe(false);
			button('Accept').click();
			await settle();
			expect(calls('/accept')).toHaveLength(2);
			expect(document.body.textContent).toContain('S-0041 is accepted');
		});

		it('leaves Accept on for a file the editor cannot open, shown without a link', async () => {
			refuse([finding({ path: 'flaiover/README.md', line: 1 })]);
			c = mount(Review, { target: document.body, props: { id: 'S-0041' } });
			await settle();
			button('Accept').click();
			await settle();
			expect(document.body.textContent).toContain('flaiover/README.md:1');
			expect(document.querySelector('button[aria-label^="Edit flaiover"]')).toBeNull();
			expect(button('Accept').disabled).toBe(false);
		});
	});

	// ADR-0155, S-0388: a path the branch undoes of main's change is a card with main's lines beside
	// the branch's change; it is restored, kept with a reason, or the story goes back naming it.
	describe('paths the branch undoes of main’s change', () => {
		const undone = (path: string, lines = ['kept on main']) => ({
			path,
			lines,
			old_base: '0123456789abcdef'
		});
		// as flai's preview.UndoneBlocker words it
		const blocker = (path: string) =>
			`story/S-0041 undoes the main branch's change to ${path}, removing 1 line main added since the branch left it; restore main's version with flai stream restore S-0041 --path ${path}, keep the removal with flai accept S-0041 --keep-removal ${path} --reason "<why>", or send S-0041 back to in-progress (ADR-0155)`;
		const previewOf = (...paths: string[]) => ({
			branch: 'story/S-0041',
			plan: null,
			blockers: paths.map(blocker),
			undoes_main: paths.map((p) => undone(p))
		});
		const cards = () => [...document.querySelectorAll<HTMLElement>('[data-testid="undoes-main"]')];
		const cardOf = (path: string) => cards().find((c) => c.textContent?.includes(path))!;
		const within = (card: HTMLElement, label: string) =>
			[...card.querySelectorAll('button')].find((b) =>
				(b.textContent ?? '').trim().startsWith(label)
			)!;
		const words = (e: Element) => (e.textContent ?? '').replace(/\s+/g, ' ').trim();
		const type = (input: HTMLInputElement, value: string) => {
			input.value = value;
			input.dispatchEvent(new Event('input', { bubbles: true }));
			flushSync();
		};
		const body = (call: unknown[]) => JSON.parse(String((call[1] as { body: string }).body));
		const reads = () => api.mock.calls.filter((c) => String(c[0]).endsWith('/acceptance')).length;

		it('shows a card per path: main’s lines beside the branch’s change, and not the blocker twice', async () => {
			backend({
				preview: {
					...previewOf('docs/a.md', 'scripts/gone.sh'),
					undoes_main: [undone('docs/a.md', ['two', 'kept on main']), undone('scripts/gone.sh')]
				}
			});
			c = mount(Review, { target: document.body, props: { id: 'S-0041' } });
			await settle();
			expect(cards()).toHaveLength(2);
			const a = cardOf('docs/a.md');
			expect(words(a)).toContain("story/S-0041 undoes main's change to docs/a.md");
			expect(words(a)).toContain('removes 2 lines main added');
			expect(a.querySelector('[data-testid="undoes-main-lines"]')!.textContent).toBe(
				'two\nkept on main'
			);
			// the branch's version, from the story's diff, shown without opening it
			const branch = a.querySelector('[data-testid="undoes-main-branch"]')!;
			expect([...branch.querySelectorAll('[data-line]')].map(words)).toEqual([
				'@@ -1,2 +1,3 @@',
				'one',
				'-two',
				'+2',
				'+three'
			]);
			expect(words(cardOf('scripts/gone.sh'))).toContain(
				"The branch's version of this path is not in its diff."
			);
			expect(document.body.textContent).not.toContain('This cannot be accepted from here yet');
			expect(button('Accept').disabled).toBe(true);
		});

		it('restores a path, reads the preview again, and the card clears', async () => {
			let restored = false;
			backend({
				preview: () => (restored ? previewOf('docs/b.md') : previewOf('docs/a.md', 'docs/b.md')),
				restore: () => {
					restored = true;
					return json({ id: 'S-0041', restored: ['docs/a.md'] });
				}
			});
			c = mount(Review, { target: document.body, props: { id: 'S-0041' } });
			await settle();
			expect(reads()).toBe(1);
			within(cardOf('docs/a.md'), "Restore main's version").click();
			await settle();
			expect(calls('/api/items/S-0041/restore')).toHaveLength(1);
			expect(body(calls('/api/items/S-0041/restore')[0])).toEqual({ path: 'docs/a.md' });
			expect(reads()).toBe(2);
			expect(cards().map((card) => card.textContent?.includes('docs/b.md'))).toEqual([true]);
			expect(button('Accept').disabled).toBe(true);
		});

		it('shows on the card why a restore failed, and keeps the card', async () => {
			const why = 'the worktree .flai-cache/worktrees/S-0041 has uncommitted changes';
			backend({ preview: previewOf('docs/a.md'), restore: () => json({ error: why }, 422) });
			c = mount(Review, { target: document.body, props: { id: 'S-0041' } });
			await settle();
			within(cardOf('docs/a.md'), "Restore main's version").click();
			await settle();
			const alert = cardOf('docs/a.md').querySelector('[role="alert"]')!;
			expect(words(alert)).toContain(`docs/a.md was not restored. flai said: ${why}`);
			expect(reads()).toBe(1);
		});

		it('keeps a removal with a reason, which Accept sends as keep_removals', async () => {
			backend({ preview: previewOf('docs/a.md', 'docs/b.md') });
			c = mount(Review, { target: document.body, props: { id: 'S-0041' } });
			await settle();
			for (const path of ['docs/a.md', 'docs/b.md']) {
				within(cardOf(path), 'Keep').click();
				flushSync();
				expect(within(cardOf(path), 'Keep the removal').disabled).toBe(true);
				type(cardOf(path).querySelector('input')!, `the story rewrites ${path}`);
				within(cardOf(path), 'Keep the removal').click();
				flushSync();
				expect(words(cardOf(path).querySelector('[data-testid="undoes-main-kept"]')!)).toContain(
					`Kept: the story rewrites ${path}`
				);
				// one kept, one open: Accept waits for the second
				expect(button('Accept').disabled).toBe(path === 'docs/a.md');
			}
			// taken back before Accept, then kept again
			within(cardOf('docs/b.md'), 'Undo').click();
			flushSync();
			expect(button('Accept').disabled).toBe(true);
			within(cardOf('docs/b.md'), 'Keep').click();
			flushSync();
			type(cardOf('docs/b.md').querySelector('input')!, 'gone with the old build');
			within(cardOf('docs/b.md'), 'Keep the removal').click();
			flushSync();
			button('Accept').click();
			await settle();
			expect(body(calls('/accept')[0])).toEqual({
				keep_removals: [
					{ path: 'docs/a.md', reason: 'the story rewrites docs/a.md' },
					{ path: 'docs/b.md', reason: 'gone with the old build' }
				]
			});
			expect(document.body.textContent).toContain('S-0041 is accepted');
		});

		it('sends the story back: a thread naming the paths, then the move with that reason', async () => {
			backend({ preview: previewOf('docs/a.md', 'docs/b.md') });
			c = mount(Review, { target: document.body, props: { id: 'S-0041' } });
			await settle();
			button('Send back, naming the undone paths').click();
			await settle();
			const order = api.mock.calls
				.filter((call) => (call[1] as { method?: string })?.method === 'POST')
				.map((call) => String(call[0]));
			expect(order).toEqual(['/api/threads', '/api/items/S-0041/move']);
			const thread = body(calls('/api/threads')[0]);
			expect(thread.on).toBe('S-0041');
			expect(thread.title).toContain('docs/a.md, docs/b.md');
			expect(thread.text).toContain('docs/a.md, docs/b.md');
			const move = body(calls('/move')[0]);
			expect(move.to).toBe('in-progress');
			expect(move.reason).toContain("undoes main's change to docs/a.md, docs/b.md");
		});

		it('moves nothing when the thread could not be opened, and says so', async () => {
			backend({
				preview: previewOf('docs/a.md'),
				thread: () => json({ error: 'flai is not available' }, 503)
			});
			c = mount(Review, { target: document.body, props: { id: 'S-0041' } });
			await settle();
			button('Send back, naming the undone paths').click();
			await settle();
			expect(calls('/move')).toHaveLength(0);
			expect(document.body.textContent).toContain(
				'No thread was opened, and S-0041 stays in review: flai is not available'
			);
		});

		it('shows the cards a refused acceptance names, and holds Accept until each is kept', async () => {
			const refusal = ndjson([
				{
					event: 'error',
					status: 422,
					error: "rule: S-0041 cannot be accepted: story/S-0041 undoes main's change to docs/a.md",
					data: {
						id: 'S-0041',
						reason: 'S-0041 cannot be accepted',
						undoes_main: [undone('docs/a.md')]
					}
				}
			]);
			const accepted = ndjson([{ event: 'done', result: { id: 'S-0041', status: 'done' } }]);
			let answers = 0;
			backend({ accept: () => (++answers > 1 ? accepted : refusal) });
			c = mount(Review, { target: document.body, props: { id: 'S-0041' } });
			await settle();
			expect(cards()).toHaveLength(0);
			button('Accept').click();
			await settle();
			expect(cards()).toHaveLength(1);
			expect(
				cardOf('docs/a.md').querySelector('[data-testid="undoes-main-lines"]')!.textContent
			).toBe('kept on main');
			expect(button('Accept').disabled).toBe(true);
			within(cardOf('docs/a.md'), 'Keep').click();
			flushSync();
			type(cardOf('docs/a.md').querySelector('input')!, 'meant');
			within(cardOf('docs/a.md'), 'Keep the removal').click();
			flushSync();
			button('Accept').click();
			await settle();
			expect(body(calls('/accept')[1])).toEqual({
				keep_removals: [{ path: 'docs/a.md', reason: 'meant' }]
			});
		});
	});

	it('links the files another refusal names to the editor on the story’s branch', async () => {
		const message =
			'rebase of story/S-0041 onto main stopped with conflicts in docs/users/flai.md:12';
		backend({ accept: () => ndjson([{ event: 'error', status: 500, error: message }]) });
		c = mount(Review, { target: document.body, props: { id: 'S-0041' } });
		await settle();
		button('Accept').click();
		await settle();
		const label = 'Edit docs/users/flai.md at line 12';
		const open = document.querySelector<HTMLButtonElement>(`button[aria-label="${label}"]`)!;
		expect(open.closest('pre')!.textContent).toContain(message);
		open.click();
		expect(openEditor).toHaveBeenCalledWith(
			expect.objectContaining({ path: 'docs/users/flai.md', line: 12, story: 'S-0041' })
		);
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

	// ADR-0106: a branch that changes paths Claude Code protects is accepted by the operator only;
	// the page names the files and says so, and keeps the operator's Accept on.
	it('names the protected files and that only the operator accepts, and keeps Accept on', async () => {
		const sentence =
			'only the operator (alex) accepts S-0041: its branch changes paths Claude Code protects';
		backend({
			preview: {
				branch: 'story/S-0041',
				plan: null,
				protected: ['.claude/settings.json', '.mcp.json'],
				operator_only: sentence
			}
		});
		c = mount(Review, { target: document.body, props: { id: 'S-0041' } });
		await settle();
		const box = document.querySelector('[data-testid="accept-protected"]')!;
		expect(box.textContent).toContain(sentence);
		expect([...box.querySelectorAll('li')].map((li) => li.textContent)).toEqual([
			'.claude/settings.json',
			'.mcp.json'
		]);
		expect(button('Accept').disabled).toBe(false);
	});

	it('says nothing of protected paths when the branch changes none', async () => {
		backend({});
		c = mount(Review, { target: document.body, props: { id: 'S-0041' } });
		await settle();
		expect(document.querySelector('[data-testid="accept-protected"]')).toBeNull();
		expect(document.body.textContent).not.toContain('only the operator');
		expect(button('Accept').disabled).toBe(false);
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

	// S-0442: a branch that changes one of the worktree's components offers the Components page
	// with the story picked; any other branch, or a catalogue flai cannot read, offers nothing.
	describe('Open in Components (S-0442)', () => {
		const component = (path: string, name: string) => ({
			project: 'flaiover',
			path,
			name,
			props: []
		});
		const catalogue = {
			catalogue: {
				projects: [
					{
						name: 'flaiover',
						path: 'flaiover',
						supported: true,
						components: [
							component('src/lib/components/Badge.svelte', 'Badge'),
							component('src/lib/components/Card.svelte', 'Card'),
							component('src/lib/components/Added.svelte', 'Added')
						]
					}
				]
			},
			status: { servers: [] }
		};
		const changing =
			(...paths: string[]) =>
			() =>
				json({ ...DIFF, files: paths.map((path) => ({ ...DIFF.files[0], path })) });
		const offer = () => document.querySelector('[data-testid="open-in-components"]');

		it('links to the Components page with the story picked and names the changed components', async () => {
			backend({
				diff: changing(
					'docs/a.md',
					'flaiover/src/lib/components/Card.svelte',
					'flaiover/src/lib/components/Added.svelte'
				),
				components: () => json(catalogue)
			});
			c = mount(Review, { target: document.body, props: { id: 'S-0041' } });
			await settle();
			expect(api.mock.calls.map((call) => call[0])).toContain('/api/components?story=S-0041');
			const link = [...(offer()?.querySelectorAll('a') ?? [])].find(
				(a) => a.textContent === 'Open in Components'
			);
			expect(link?.getAttribute('href')).toBe('/design/components?story=S-0041');
			expect(offer()?.textContent?.replace(/\s+/g, ' ')).toContain(
				'It changes the components Card, Added.'
			);
			expect(offer()?.textContent).not.toContain('Badge');
		});

		it('offers nothing when the branch changes no component', async () => {
			backend({ components: () => json(catalogue) });
			c = mount(Review, { target: document.body, props: { id: 'S-0041' } });
			await settle();
			expect(api.mock.calls.map((call) => call[0])).toContain('/api/components?story=S-0041');
			expect(offer()).toBeNull();
			expect(document.body.textContent).not.toContain('Open in Components');
		});

		it('offers nothing, and shows no error, when the catalogue cannot be read', async () => {
			backend({
				diff: changing('flaiover/src/lib/components/Card.svelte'),
				components: () => json({ error: 'flai 1.55.1 cannot list the components' }, 502)
			});
			c = mount(Review, { target: document.body, props: { id: 'S-0041' } });
			await settle();
			expect(offer()).toBeNull();
			expect(document.body.textContent).not.toContain('cannot list the components');
			expect(document.querySelector('[role="alert"]')).toBeNull();
		});
	});

	describe('issues (S-0198)', () => {
		const issue = (id: string, over: Record<string, unknown> = {}) => ({
			id,
			title: `Title of ${id}`,
			class: 'defect',
			status: 'open',
			count: 1,
			path: `/repo/design/issues/${id}-slug.md`,
			stories: [],
			story: '',
			...over
		});
		const MAIN = [issue('I-0003'), issue('I-0004', { story: 'S-0050' })];
		const OWN = [issue('I-0012', { stories: ['S-0041'] })];
		const boxes = () =>
			[...document.querySelectorAll<HTMLInputElement>('[data-testid="issues-section"] input')].map(
				(b) => [b.closest('li')!.querySelector('a')!.textContent, b.checked]
			);
		const posted = () =>
			api.mock.calls
				.filter((c) => /\/api\/issues\/.*\/story$/.test(String(c[0])))
				.map((c) => String(c[0]));

		it('offers the story’s own issues checked, then the other open ones, and names the button for it', async () => {
			backend({ issues: MAIN, ownIssues: OWN });
			c = mount(Review, { target: document.body, props: { id: 'S-0041' } });
			await settle();
			expect(boxes()).toEqual([
				['I-0012', true],
				['I-0003', false]
			]);
			const link = document.querySelector('[data-testid="issues-section"] a')!;
			expect(link.getAttribute('href')).toBe('/docs/design/issues/I-0012-slug.md');
			expect(button('Accept').textContent?.trim()).toBe('Accept and Create Stories');
			expect(document.querySelector('[data-testid="accept-makes-stories"]')?.textContent).toContain(
				'I-0012'
			);
			document.querySelector<HTMLInputElement>('[data-testid="issues-section"] input')!.click();
			flushSync();
			expect(button('Accept').textContent?.trim()).toBe('Accept');
		});

		it('shows no issues section when there is none to offer', async () => {
			backend({ issues: [issue('I-0001', { status: 'closed' })] });
			c = mount(Review, { target: document.body, props: { id: 'S-0041' } });
			await settle();
			expect(document.querySelector('[data-testid="issues-section"]')).toBeNull();
			expect(button('Accept').textContent?.trim()).toBe('Accept');
		});

		it('accepts first, then makes a story for each checked issue, naming each that failed and why', async () => {
			backend({
				issues: MAIN,
				ownIssues: OWN,
				issueStory: (id) =>
					id === 'I-0012'
						? json({ id: 'S-0060', title: 'Title of I-0012', issue: id, committed: true })
						: json({ error: 'I-0003 is already linked by open story S-0061' }, 400)
			});
			c = mount(Review, { target: document.body, props: { id: 'S-0041' } });
			await settle();
			const second = document.querySelectorAll<HTMLInputElement>(
				'[data-testid="issues-section"] input'
			)[1];
			second.click();
			flushSync();
			button('Accept').click();
			await settle();
			const order = api.mock.calls.map((c) => String(c[0]));
			expect(order.indexOf('/api/items/S-0041/accept')).toBeLessThan(
				order.indexOf('/api/issues/I-0012/story')
			);
			expect(posted()).toEqual(['/api/issues/I-0012/story', '/api/issues/I-0003/story']);
			expect(JSON.parse(String(calls('/api/issues/I-0012/story')[0][1].body))).toEqual({});
			expect(document.querySelector('[data-testid="issue-story-made"]')?.textContent).toContain(
				'S-0060'
			);
			const failed = document.querySelector('[data-testid="issue-story-failed"]')!.textContent!;
			expect(failed).toContain('I-0003');
			expect(failed).toContain('already linked by open story S-0061');
			expect(failed).not.toContain('I-0012');
		});

		it('makes no story when the acceptance fails', async () => {
			backend({
				issues: MAIN,
				ownIssues: OWN,
				accept: () => ndjson([{ event: 'error', status: 500, error: 'rebase stopped' }])
			});
			c = mount(Review, { target: document.body, props: { id: 'S-0041' } });
			await settle();
			button('Accept').click();
			await settle();
			expect(calls('/accept')).toHaveLength(1);
			expect(posted()).toEqual([]);
		});
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
			// the run's times in the local zone (S-0329): the tests run in New York, where midnight
			// in UTC is the evening before
			const summary = document.querySelector('[data-testid="checks-summary"]')!.textContent!;
			expect(summary.replace(/\s+/g, ' ')).toContain(
				'passed, started 2026-09-21 20:00 EDT, ended 2026-09-21 20:05 EDT.'
			);
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
			expect(document.querySelector('[data-testid="checks-summary"]')!.textContent).toContain(
				'Running flai, started 2026-09-21 20:00 EDT.'
			);
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

	describe('verification (S-0270)', () => {
		const REPORT = {
			story: 'S-0041',
			commit: '0123456789abcdef0123',
			base: 'main',
			ran_at: '2026-10-06T22:00:00Z',
			duration_ms: 42100,
			duration: '42.1s',
			passed: true,
			steps: [
				{ name: 'rebase', state: 'passed', duration_ms: 3, duration: '3ms' },
				{ name: 'check', state: 'passed', duration_ms: 1500, duration: '1.5s' },
				{ name: 'flaiover', tier: true, state: 'passed', duration_ms: 40000, duration: '40s' }
			]
		};
		const section = () => document.querySelector('[data-testid="verify-section"]');
		const steps = () =>
			[...document.querySelectorAll('[data-testid="verify-step"]')].map((li) =>
				(li.textContent ?? '').replace(/\s+/g, ' ').trim()
			);

		it('shows a stored passing result: the outcome, the commit and when, each step’s state and duration', async () => {
			backend({ verify: () => json({ report: REPORT }) });
			c = mount(Review, { target: document.body, props: { id: 'S-0041' } });
			await settle();
			const summary = document.querySelector('[data-testid="verify-summary"]')!.textContent!;
			expect(summary).toContain('Passed every step');
			expect(summary).toContain('0123456789ab');
			expect(summary).not.toContain('0123456789abcdef');
			// when it ran, in the local zone (S-0329): the tests run in New York
			expect(summary).toContain('2026-10-06 18:00 EDT');
			expect(summary).not.toContain('2026-10-06T22:00:00Z');
			expect(steps()).toEqual([
				'passed rebase 3ms',
				'passed check 1.5s',
				'passed flaiover tier 40s'
			]);
			expect(section()!.querySelector('button')).toBeNull();
			expect(section()!.querySelector('[data-testid="verify-notes"]')).toBeNull();
		});

		it('shows where a failing result stopped, the failed step’s findings, and the notes apart', async () => {
			backend({
				verify: () =>
					json({
						report: {
							...REPORT,
							passed: false,
							stopped_at: 'flaiover',
							steps: [
								REPORT.steps[0],
								REPORT.steps[1],
								{
									name: 'flaiover',
									tier: true,
									state: 'failed',
									duration_ms: 40000,
									duration: '40s',
									findings: [
										{
											name: 'Review > shows it',
											path: 'flaiover/src/lib/x.test.ts',
											line: 12,
											message: 'expected 1 to be 2'
										}
									]
								},
								{ name: 'flai', tier: true, state: 'not-reached', duration_ms: 0 }
							],
							notes: [
								{
									rule: 'stale-narrative',
									level: 'warning',
									path: 'wip/agents/S-0001.md',
									message: 'not touched in a week'
								}
							]
						}
					})
			});
			c = mount(Review, { target: document.body, props: { id: 'S-0041' } });
			await settle();
			expect(document.querySelector('[data-testid="verify-summary"]')!.textContent).toContain(
				'Stopped at flaiover'
			);
			const rows = steps();
			expect(rows[2]).toContain('failed flaiover tier 40s');
			expect(rows[2]).toContain(
				'flaiover/src/lib/x.test.ts:12 Review > shows it: expected 1 to be 2'
			);
			expect(rows[3]).toBe('not reached flai tier');
			expect(rows[0]).toBe('passed rebase 3ms');
			const notes = section()!.querySelector('[data-testid="verify-notes"]')!;
			expect(notes.tagName).toBe('DETAILS');
			const noteText = (notes.textContent ?? '').replace(/\s+/g, ' ');
			expect(noteText).toContain('1 note from flai check outside the story');
			expect(noteText).toContain('wip/agents/S-0001.md stale-narrative: warning: not touched');
			expect(rows.join(' ')).not.toContain('stale-narrative');
		});

		it('shows nothing and no error when the story has no result', async () => {
			backend({});
			c = mount(Review, { target: document.body, props: { id: 'S-0041' } });
			await settle();
			expect(api.mock.calls.map((call) => call[0])).toContain('/api/items/S-0041/verify');
			expect(section()).toBeNull();
			expect(document.body.textContent).not.toContain('Verification');
			expect(document.querySelector('[role="alert"]')).toBeNull();
		});
	});

	describe('the pull request (S-0394)', () => {
		const PULL_REQUEST = {
			story: 'S-0041',
			pull_request: 7,
			url: 'https://github.com/bytepunx/system-flow/pull/7',
			head_sha: '0123456789abcdef0123456789abcdef01234567',
			read_at: '2026-10-17T10:00:00Z',
			checks: [
				{
					name: 'flai',
					state: 'failed',
					url: 'https://github.com/bytepunx/system-flow/actions/runs/1',
					summary: 'go test failed',
					since: '2026-10-17T09:58:00Z'
				}
			],
			reviews: [
				{
					author: 'alex',
					state: 'changes-requested',
					url: 'https://github.com/bytepunx/system-flow/pull/7#pullrequestreview-1'
				}
			]
		};
		const section = () => document.querySelector('[data-testid="forge-status"]');

		it('shows the link to it, each required check, and each review, after the checks', async () => {
			backend({ forge: () => json({ story: 'S-0041', pull_request: PULL_REQUEST }) });
			c = mount(Review, { target: document.body, props: { id: 'S-0041' } });
			await settle();
			expect(api.mock.calls.map((call) => call[0])).toContain('/api/items/S-0041/forge');
			const checks = document.querySelector('[data-testid="checks-section"]')!;
			expect(checks.nextElementSibling).toBe(section());
			const link = section()!.querySelector<HTMLAnchorElement>('h2 a')!;
			expect(link.textContent!.trim()).toBe('Pull request #7');
			expect(link.getAttribute('href')).toBe(PULL_REQUEST.url);
			const check = section()!.querySelector('[data-check="flai"]')!;
			expect(check.textContent).toContain('failed');
			expect(check.textContent).toContain('go test failed');
			expect(check.querySelector('a')!.getAttribute('href')).toBe(PULL_REQUEST.checks[0].url);
			const review = section()!.querySelector('[data-review="alex"]')!;
			expect(review.textContent).toContain('changes requested');
			expect(review.querySelector('a')!.getAttribute('href')).toBe(PULL_REQUEST.reviews[0].url);
		});

		it('shows nothing for a story with none, nor when it cannot be read', async () => {
			backend({});
			c = mount(Review, { target: document.body, props: { id: 'S-0041' } });
			await settle();
			expect(section()).toBeNull();
			unmount(c!);
			document.body.innerHTML = '';
			backend({ forge: () => json({ error: 'no flai connected' }, 503) });
			c = mount(Review, { target: document.body, props: { id: 'S-0041' } });
			await settle();
			expect(section()).toBeNull();
			expect(document.querySelector('[data-testid="checks-section"]')).not.toBeNull();
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
