// S-0080: the root page shows a project list, with each project's glance, while more than one is
// known and none is chosen; otherwise the single-project overview it always has.
import { afterEach, describe, expect, it, vi } from 'vitest';
import { flushSync, mount, unmount } from 'svelte';
import { resetForTests } from '$lib/project.svelte';

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
		vi.stubGlobal(
			'fetch',
			vi.fn(async (url: string) => {
				if (url.startsWith('/api/projects'))
					return answer({ projects: [{ key: 'harbour', name: 'Harbour', connected: true }] });
				if (url.startsWith('/api/manifest'))
					return answer({ name: 'Harbour', description: 'a project' });
				if (url.startsWith('/api/items')) return answer([]);
				return answer({});
			})
		);
		const { default: Overview } = await import('./+page.svelte');
		c = mount(Overview, { target: document.body });
		await settle();
		expect(document.querySelector('[data-testid="project-list"]')).toBeNull();
		expect(document.body.textContent).toContain('Harbour');
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
