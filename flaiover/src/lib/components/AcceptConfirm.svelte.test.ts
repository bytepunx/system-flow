import { afterEach, describe, expect, it, vi } from 'vitest';
import { flushSync, mount, unmount } from 'svelte';
import AcceptConfirm from './AcceptConfirm.svelte';

const api = vi.fn();
vi.mock('$lib/api', () => ({ api: (...args: unknown[]) => api(...args) }));

const settle = async () => {
	for (let i = 0; i < 5; i++) await new Promise((r) => setTimeout(r, 0));
	flushSync();
};
const json = (body: unknown, ok = true) => ({
	ok,
	statusText: ok ? 'OK' : 'Conflict',
	json: async () => body
});

describe('AcceptConfirm', () => {
	afterEach(() => {
		api.mockReset();
		document.body.innerHTML = '';
	});

	it('shows the release plan and the branch, and only acts on confirm', async () => {
		api.mockResolvedValue(
			json({
				id: 'S-0046',
				branch: 'story/S-0046',
				plan: {
					level: 'minor',
					commits: ['a', 'b'],
					steps: [
						{
							component: { name: 'flai' },
							delivered: true,
							level: 'minor',
							from: { Major: 1, Minor: 1, Patch: 6 },
							to: { Major: 1, Minor: 2, Patch: 0 },
							files: ['x', 'y']
						},
						{
							component: { name: 'flaiover' },
							delivered: false,
							level: 'patch',
							from: '0.10.2',
							to: '0.10.3',
							files: ['z']
						}
					]
				}
			})
		);
		const onconfirm = vi.fn();
		const oncancel = vi.fn();
		const c = mount(AcceptConfirm, {
			target: document.body,
			props: { id: 'S-0046', onconfirm, oncancel }
		});
		await settle();
		expect(api).toHaveBeenCalledWith('/api/items/S-0046/acceptance');
		const text = document.body.textContent ?? '';
		expect(text).toContain('story/S-0046');
		expect(text).toContain('minor release from 2');
		expect(text).toMatch(/flai\s+1\.1\.6 → 1\.2\.0/);
		expect(text).toMatch(/flaiover\s+0\.10\.2 → 0\.10\.3/);
		const [cancel, accept] = [...document.querySelectorAll('button')];
		cancel.click();
		expect(oncancel).toHaveBeenCalledOnce();
		expect(onconfirm).not.toHaveBeenCalled();
		accept.click();
		await settle();
		expect(onconfirm).toHaveBeenCalledOnce();
		unmount(c);
	});

	it('accepts a research story, saying it releases nothing and what lands unreleased (ADR-0025)', async () => {
		api.mockResolvedValue(
			json({
				id: 'S-0052',
				branch: 'story/S-0052',
				blockers: [],
				plan: {
					level: 'none',
					commits: ['a'],
					steps: [],
					skipped:
						'S-0052 is research: its findings land on main and are pushed, and research cuts no release whatever it touched (ADR-0025)',
					unreleased: [{ component: 'flai', files: ['flai/cmd/x.go', 'flai/cmd/x_test.go'] }]
				}
			})
		);
		const onconfirm = vi.fn();
		const c = mount(AcceptConfirm, {
			target: document.body,
			props: { id: 'S-0052', onconfirm, oncancel: vi.fn() }
		});
		await settle();
		const text = (document.body.textContent ?? '').replace(/\s+/g, ' ');
		expect(text).toContain('No release: S-0052 is research');
		expect(text).toContain('flai lands on main without a release (2 files)');
		expect(text).not.toContain('Tag and push');
		const accept = [...document.querySelectorAll('button')].at(-1)!;
		expect(accept.disabled).toBe(false);
		accept.click();
		await settle();
		expect(onconfirm).toHaveBeenCalledOnce();
		unmount(c);
	});

	it('lists nothing as unreleased for research that touched no component', async () => {
		api.mockResolvedValue(
			json({
				id: 'S-0052',
				blockers: [],
				plan: { level: 'none', commits: ['a'], steps: [], skipped: 'S-0052 is research' }
			})
		);
		const c = mount(AcceptConfirm, {
			target: document.body,
			props: { id: 'S-0052', onconfirm: vi.fn(), oncancel: vi.fn() }
		});
		await settle();
		expect(document.querySelector('[data-testid="unreleased"]')).toBeNull();
		expect(document.body.textContent).toContain('No release: S-0052 is research');
		unmount(c);
	});

	// ADR-0066: an experiment is accepted like research once its results document is on its
	// branch, and flai's blocker names the document while it is missing.
	it('accepts an experiment with its results document, releasing nothing', async () => {
		api.mockResolvedValue(
			json({
				id: 'S-0176',
				branch: 'story/S-0176',
				blockers: [],
				plan: {
					level: 'none',
					commits: ['a'],
					steps: [],
					skipped:
						'S-0176 is an experiment: its results land on main, and an experiment cuts no release whatever it touched (ADR-0066)',
					unreleased: [{ component: 'flai', files: ['flai/cmd/x.go'] }]
				}
			})
		);
		const onconfirm = vi.fn();
		const c = mount(AcceptConfirm, {
			target: document.body,
			props: { id: 'S-0176', onconfirm, oncancel: vi.fn() }
		});
		await settle();
		const text = (document.body.textContent ?? '').replace(/\s+/g, ' ');
		expect(text).toContain('No release: S-0176 is an experiment');
		expect(text).toContain('flai lands on main without a release (1 file)');
		expect(document.querySelector('[role=alert]')).toBeNull();
		const accept = [...document.querySelectorAll('button')].at(-1)!;
		expect(accept.disabled).toBe(false);
		accept.click();
		await settle();
		expect(onconfirm).toHaveBeenCalledOnce();
		unmount(c);
	});

	it('refuses an experiment without its results document, naming the document', async () => {
		api.mockResolvedValue(
			json({
				id: 'S-0176',
				branch: 'story/S-0176',
				blockers: [
					'S-0176 is an experiment and has no results document on its branch story/S-0176: write design/experiments/S-0176-trial.md with its hypothesis, success measure, what was done, results, and recommendation (adopt, adapt, or drop), commit it, and accept again (ADR-0066)'
				],
				plan: null
			})
		);
		const onconfirm = vi.fn();
		const c = mount(AcceptConfirm, {
			target: document.body,
			props: { id: 'S-0176', onconfirm, oncancel: vi.fn() }
		});
		await settle();
		expect(document.querySelector('[role=alert]')?.textContent).toContain(
			'design/experiments/S-0176-trial.md'
		);
		const accept = [...document.querySelectorAll('button')][1] as HTMLButtonElement;
		expect(accept.disabled).toBe(true);
		accept.click();
		expect(onconfirm).not.toHaveBeenCalled();
		unmount(c);
	});

	it('lists blockers from the preview and keeps accept disabled', async () => {
		api.mockResolvedValue(
			json({ id: 'S-0046', blockers: ['git has no committer identity here'], plan: null })
		);
		const onconfirm = vi.fn();
		const c = mount(AcceptConfirm, {
			target: document.body,
			props: { id: 'S-0046', onconfirm, oncancel: vi.fn() }
		});
		await settle();
		expect(document.querySelector('[role=alert]')?.textContent).toContain('committer identity');
		expect(([...document.querySelectorAll('button')][1] as HTMLButtonElement).disabled).toBe(true);
		unmount(c);
	});

	it('shows flai’s refusal verbatim and keeps accept disabled', async () => {
		api.mockResolvedValue(
			json({ error: 'S-0046 is in-progress; a story is accepted from review' }, false)
		);
		const onconfirm = vi.fn();
		const c = mount(AcceptConfirm, {
			target: document.body,
			props: { id: 'S-0046', onconfirm, oncancel: vi.fn() }
		});
		await settle();
		expect(document.querySelector('[role=alert]')?.textContent).toContain(
			'a story is accepted from review'
		);
		const accept = [...document.querySelectorAll('button')][1] as HTMLButtonElement;
		expect(accept.disabled).toBe(true);
		accept.click();
		expect(onconfirm).not.toHaveBeenCalled();
		unmount(c);
	});

	it('passes false to onconfirm when there is nothing uncommitted', async () => {
		api.mockResolvedValue(json({ id: 'S-0051', plan: null }));
		const onconfirm = vi.fn();
		const c = mount(AcceptConfirm, {
			target: document.body,
			props: { id: 'S-0051', onconfirm, oncancel: vi.fn() }
		});
		await settle();
		expect(document.querySelector('input[type=checkbox]')).toBeNull();
		([...document.querySelectorAll('button')][1] as HTMLButtonElement).click();
		await settle();
		expect(onconfirm).toHaveBeenCalledWith(false);
		unmount(c);
	});

	it('lists uncommitted paths and waits for the choice before accepting', async () => {
		api.mockResolvedValue(
			json({ id: 'S-0051', plan: null, uncommitted: ['.claude/', 'docs/stray.md'] })
		);
		const onconfirm = vi.fn();
		const c = mount(AcceptConfirm, {
			target: document.body,
			props: { id: 'S-0051', onconfirm, oncancel: vi.fn() }
		});
		await settle();
		const text = document.body.textContent ?? '';
		expect(text).toContain('.claude/');
		expect(text).toContain('docs/stray.md');
		expect(text).toContain('holds only acceptance');
		const box = document.querySelector('input[type=checkbox]') as HTMLInputElement;
		const accept = [...document.querySelectorAll('button')][1] as HTMLButtonElement;
		expect(box.checked).toBe(false); // off by default
		expect(accept.disabled).toBe(true);
		accept.click();
		expect(onconfirm).not.toHaveBeenCalled();

		box.click();
		await settle();
		expect(accept.disabled).toBe(false);
		accept.click();
		await settle();
		expect(onconfirm).toHaveBeenCalledWith(true);
		unmount(c);
	});

	it('keeps accept disabled for a blocker even when the choice is made', async () => {
		api.mockResolvedValue(
			json({ id: 'S-0051', plan: null, blockers: ['no identity'], uncommitted: ['x.md'] })
		);
		const c = mount(AcceptConfirm, {
			target: document.body,
			props: { id: 'S-0051', onconfirm: vi.fn(), oncancel: vi.fn() }
		});
		await settle();
		(document.querySelector('input[type=checkbox]') as HTMLInputElement).click();
		await settle();
		expect(([...document.querySelectorAll('button')][1] as HTMLButtonElement).disabled).toBe(true);
		unmount(c);
	});

	// S-0140: work left uncommitted in the story's worktree blocks acceptance, and the confirmation
	// offers to have an agent commit it.
	it('offers to have an agent commit what the worktree holds, and keeps Accept off', async () => {
		api.mockImplementation(async (path: string, init?: RequestInit) => {
			if (path === '/api/items/S-0140/acceptance')
				return json({
					id: 'S-0140',
					branch: 'story/S-0140',
					worktree_uncommitted: ['flai/cmd/move.go', 'docs/left.md'],
					blockers: [
						'the worktree .flai-cache/worktrees/S-0140 has uncommitted changes (flai/cmd/move.go, docs/left.md); commit them on story/S-0140, or discard them, before accepting',
						'git has no committer identity here'
					],
					plan: null
				});
			if (path === '/api/items/S-0140/agent') {
				expect(JSON.parse(String(init?.body))).toEqual({ action: 'commit' });
				return json({ story: 'S-0140', agent: 'agent-S-0140', pid: 44 });
			}
			throw new Error(`unexpected ${path}`);
		});
		const onconfirm = vi.fn();
		const c = mount(AcceptConfirm, {
			target: document.body,
			props: { id: 'S-0140', onconfirm, oncancel: vi.fn() }
		});
		await settle();
		const box = document.querySelector('[data-testid="worktree-uncommitted"]');
		expect(box?.textContent).toContain('flai/cmd/move.go');
		expect(box?.textContent).toContain('docs/left.md');
		const alert = document.querySelector('[role="alert"]')?.textContent ?? '';
		expect(alert).toContain('git has no committer identity');
		expect(alert).not.toContain('has uncommitted changes');
		const accept = [...document.querySelectorAll('button')].find((b) => b.textContent === 'Accept');
		expect(accept?.disabled).toBe(true);
		const commit = [...document.querySelectorAll('button')].find(
			(b) => b.textContent === 'Have an agent commit them'
		);
		commit!.click();
		await settle();
		expect(api).toHaveBeenCalledWith('/api/items/S-0140/agent', expect.anything());
		expect(box?.textContent).toContain('agent-S-0140 is committing them (pid 44)');
		expect(onconfirm).not.toHaveBeenCalled();
		unmount(c);
	});

	it('says why an agent could not be started to commit', async () => {
		api.mockImplementation(async (path: string) =>
			path.endsWith('/acceptance')
				? json({ id: 'S-0140', worktree_uncommitted: ['x.go'], blockers: [], plan: null })
				: json({ error: 'the host action "agent" is not enabled: flai serve enable agent' }, false)
		);
		const c = mount(AcceptConfirm, {
			target: document.body,
			props: { id: 'S-0140', onconfirm: vi.fn(), oncancel: vi.fn() }
		});
		await settle();
		[...document.querySelectorAll('button')]
			.find((b) => b.textContent === 'Have an agent commit them')!
			.click();
		await settle();
		expect(document.body.textContent).toContain('flai serve enable agent');
		expect(
			[...document.querySelectorAll('button')].some(
				(b) => b.textContent === 'Have an agent commit them'
			)
		).toBe(true);
		unmount(c);
	});
});
