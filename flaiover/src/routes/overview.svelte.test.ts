// S-0080: the root page shows a project list, with each project's glance, while more than one is
// known and none is chosen; otherwise the single-project overview it always has.
import { afterEach, describe, expect, it, vi } from 'vitest';
import { flushSync, mount, unmount } from 'svelte';
import { resetForTests } from '$lib/project.svelte';
import { baselineFlow, testingFlow } from '$lib/lanes.fixture';

const answer = (body: unknown) => ({ ok: true, json: async () => body });

describe('the root page (S-0080)', () => {
	let c: ReturnType<typeof mount> | undefined;
	afterEach(() => {
		if (c) unmount(c);
		c = undefined;
		vi.unstubAllGlobals();
		document.body.innerHTML = '';
		resetForTests();
	});
	const settle = async () => {
		for (let i = 0; i < 20; i++) await new Promise((r) => setTimeout(r, 0));
		flushSync();
	};

	it('shows the single-project overview while at most one project is known', async () => {
		const asked: string[] = [];
		vi.stubGlobal(
			'fetch',
			vi.fn(async (url: string) => {
				asked.push(url);
				if (url.startsWith('/api/projects'))
					return answer({ projects: [{ key: 'harbour', name: 'Harbour', connected: true }] });
				if (url.startsWith('/api/manifest'))
					return answer({ name: 'Harbour', description: 'a project' });
				if (url.startsWith('/api/items/count')) return answer({ active: 1, archived: 812 });
				if (url.startsWith('/api/items'))
					return answer([{ id: 'S-0001', type: 'story', status: 'ready', title: 'Berths' }]);
				return answer({});
			})
		);
		const { default: Overview } = await import('./+page.svelte');
		c = mount(Overview, { target: document.body });
		await settle();
		expect(document.querySelector('[data-testid="project-list"]')).toBeNull();
		expect(document.body.textContent).toContain('Harbour');
		// the active items and the archive's size, not the archive (S-0162)
		expect(
			asked
				.filter((u) => u.startsWith('/api/items'))
				.map((u) => u.replace(/[?&]project=[^&]*$/, ''))
				.sort()
		).toEqual(['/api/items/count', '/api/items?archived=false']);
		expect(document.body.textContent!.replace(/\s+/g, ' ')).toContain(
			'1 active items · 812 archived'
		);
	});

	// S-0753: a column for each of the board's lanes but the cancelled kind, as its workflow defines them
	it("counts the items in each of the board's lanes, all but the cancelled", async () => {
		const table = async (flow: typeof baselineFlow) => {
			vi.stubGlobal(
				'fetch',
				vi.fn(async (url: string) => {
					if (url.startsWith('/api/projects'))
						return answer({ projects: [{ key: 'harbour', name: 'Harbour', connected: true }] });
					if (url.startsWith('/api/manifest')) return answer({ name: 'Harbour' });
					if (url.startsWith('/api/board')) return answer(flow);
					if (url.startsWith('/api/items/count')) return answer({ active: 3, archived: 0 });
					if (url.startsWith('/api/items'))
						return answer([
							{ id: 'S-0001', type: 'story', status: 'ready', title: 'Berths' },
							{ id: 'S-0002', type: 'story', status: 'testing', title: 'Cranes' },
							{ id: 'T-0001', type: 'task', status: 'testing', title: 'Bollards' }
						]);
					return answer({});
				})
			);
			const { default: Overview } = await import('./+page.svelte');
			c = mount(Overview, { target: document.body });
			await settle();
			const rows = [...document.querySelectorAll('tr')].map((tr) =>
				[...tr.children].map((td) => td.textContent!.trim())
			);
			unmount(c);
			c = undefined;
			document.body.innerHTML = '';
			resetForTests();
			return rows;
		};
		expect((await table(baselineFlow))[0]).toEqual([
			'',
			'backlog',
			'ready',
			'in-progress',
			'review',
			'done'
		]);
		expect(await table(testingFlow)).toEqual([
			['', 'backlog', 'ready', 'doing', 'testing', 'review', 'done'],
			['epics', '0', '0', '0', '0', '0', '0'],
			['storys', '0', '1', '0', '1', '0', '0'],
			['tasks', '0', '0', '0', '1', '0', '0']
		]);
	});

	it('shows a list of projects, each with its glance, while more than one is known and none chosen', async () => {
		vi.stubGlobal(
			'fetch',
			vi.fn(async () =>
				answer({
					projects: [
						{
							key: 'harbour',
							name: 'Harbour',
							connected: true,
							review: 2,
							threadsAwaiting: 1,
							agentAttending: true
						},
						{ key: 'quay', name: 'Quay', connected: false }
					]
				})
			)
		);
		const { default: Overview } = await import('./+page.svelte');
		c = mount(Overview, { target: document.body });
		await settle();
		const list = document.querySelector('[data-testid="project-list"]')!;
		expect(list).not.toBeNull();
		const text = list.textContent!.replace(/\s+/g, ' ');
		expect(text).toContain('Harbour');
		expect(text).toContain('2 in review');
		expect(text).toContain('1 thread');
		expect(text).toContain('agent attending');
		expect(text).toContain('Quay');
		expect(text).toContain('not connected');
	});

	it('one load asks for the projects once: the layout and the page share the request (S-0186)', async () => {
		const asked: string[] = [];
		vi.stubGlobal(
			'fetch',
			vi.fn(async (url: string) => {
				asked.push(url);
				await new Promise((r) => setTimeout(r, 0)); // in flight while the page mounts
				return answer({
					projects: [
						{ key: 'harbour', name: 'Harbour', connected: true, review: 3 },
						{ key: 'quay', name: 'Quay', connected: true }
					]
				});
			})
		);
		const { default: Overview } = await import('./+page.svelte');
		const { projectState } = await import('$lib/project.svelte');
		void projectState.refresh(); // what +layout.svelte asks once signed in, in the same flush
		c = mount(Overview, { target: document.body });
		await settle();
		expect(asked.filter((u) => u.startsWith('/api/projects'))).toEqual(['/api/projects']);
		// the glances come from that one answer
		expect(document.querySelector('[data-testid="project-list"]')!.textContent).toContain(
			'3 in review'
		);
	});

	it('filters the project list by name or key', async () => {
		vi.stubGlobal(
			'fetch',
			vi.fn(async () =>
				answer({
					projects: [
						{ key: 'harbour', name: 'Harbour', connected: true },
						{ key: 'quay', name: 'Quay', connected: true }
					]
				})
			)
		);
		const { default: Overview } = await import('./+page.svelte');
		c = mount(Overview, { target: document.body });
		await settle();
		const filterInput = document.querySelector<HTMLInputElement>('[data-testid="project-filter"]')!;
		expect(filterInput).not.toBeNull();
		filterInput.value = 'har';
		filterInput.dispatchEvent(new Event('input', { bubbles: true }));
		flushSync();
		const list = document.querySelector('[data-testid="project-list"]')!;
		expect(list.textContent).toContain('Harbour');
		expect(list.textContent).not.toContain('Quay');
		filterInput.value = 'nothing matches this';
		filterInput.dispatchEvent(new Event('input', { bubbles: true }));
		flushSync();
		expect(list.textContent).toContain('No project matches');
	});

	it('picking a project from the list chooses it', async () => {
		vi.stubGlobal(
			'fetch',
			vi.fn(async () =>
				answer({
					projects: [
						{ key: 'harbour', name: 'Harbour', connected: true },
						{ key: 'quay', name: 'Quay', connected: true }
					]
				})
			)
		);
		const reload = vi.fn();
		vi.stubGlobal('location', { ...window.location, reload });
		const { default: Overview } = await import('./+page.svelte');
		const { projectState } = await import('$lib/project.svelte');
		c = mount(Overview, { target: document.body });
		await settle();
		document.querySelector<HTMLButtonElement>('button')!.click();
		expect(projectState.current).toBe('harbour');
		// in place (S-0095): the layout remounts the page for the project picked, no reload
		expect(reload).not.toHaveBeenCalled();
		await settle(); // let anything still in flight land before unmount, not after
	});
});
