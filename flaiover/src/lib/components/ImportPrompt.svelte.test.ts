// S-0098: picking a repository the host offers for import asks whether to import it; a decline
// says to pick another and is not remembered; an import shows its tests and, once committed,
// switches to the project the repository became.
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { flushSync, mount, unmount } from 'svelte';
import ImportPrompt from './ImportPrompt.svelte';
import { projectState, resetForTests, type Project } from '$lib/project.svelte';

const candidate: Project = {
	key: 'import-widget',
	name: 'widget',
	connected: true,
	candidate: true
};
const preview = {
	analysis: { root: '/home/op/git/widget' },
	tests: [{ name: 'go test ./...', dir: '.', command: ['go', 'test', './...'] }],
	tests_from: 'repository'
};

function ndjson(lines: unknown[]): Response {
	return new Response(lines.map((l) => JSON.stringify(l)).join('\n') + '\n', {
		headers: { 'content-type': 'application/x-ndjson' }
	});
}

describe('ImportPrompt (S-0098)', () => {
	let c: ReturnType<typeof mount> | undefined;
	let posts: string[];
	let importAnswer: unknown[];
	let projects: unknown[];
	const settle = async () => {
		for (let i = 0; i < 30; i++) await new Promise((r) => setTimeout(r, 0));
		flushSync();
	};
	const text = () => document.body.textContent ?? '';
	const button = (label: string) =>
		[...document.querySelectorAll('button')].find((b) => b.textContent?.trim() === label)!;

	beforeEach(() => {
		resetForTests();
		projectState.list = [candidate, { key: 'harbour', name: 'Harbour', connected: true }];
		projectState.pick(candidate.key, false);
		posts = [];
		projects = [candidate];
		vi.stubGlobal(
			'fetch',
			vi.fn(async (url: string, init?: RequestInit) => {
				if (url.startsWith('/api/projects'))
					return new Response(JSON.stringify({ projects }), { status: 200 });
				if (url.startsWith('/api/import') && init?.method === 'POST') {
					posts.push(url);
					return ndjson(importAnswer);
				}
				if (url.startsWith('/api/import'))
					return new Response(JSON.stringify(preview), { status: 200 });
				return new Response('{}', { status: 404 });
			})
		);
	});
	afterEach(() => {
		if (c) unmount(c);
		c = undefined;
		vi.unstubAllGlobals();
		document.body.innerHTML = '';
		resetForTests();
	});

	it('asks, saying what an import does and which tests it runs', async () => {
		c = mount(ImportPrompt, { target: document.body, props: { project: candidate } });
		await settle();
		expect(text()).toContain('not a system-flow project yet (/home/op/git/widget)');
		expect(text()).toContain('go test ./...');
		expect(text()).toContain('leave it uncommitted if one fails');
	});

	it('a decline says to pick a different project, and imports nothing', async () => {
		c = mount(ImportPrompt, { target: document.body, props: { project: candidate } });
		await settle();
		button('Not now').click();
		flushSync();
		expect(document.querySelector('[data-testid="import-declined"]')?.textContent).toContain(
			'select a different project'
		);
		expect(posts).toEqual([]);
		// declined is this prompt's own state: a fresh one (the repository picked again) asks again
		unmount(c);
		c = mount(ImportPrompt, { target: document.body, props: { project: candidate } });
		await settle();
		expect(button('Import')).toBeTruthy();
		expect(document.querySelector('[data-testid="import-declined"]')).toBeNull();
	});

	it('a failed test shows what failed, and says nothing was committed', async () => {
		importAnswer = [
			{ event: 'progress', msg: 'importing' },
			{ event: 'progress', msg: 'running go test ./...' },
			{ event: 'progress', msg: 'go test ./... failed' },
			{
				event: 'done',
				committed: false,
				result: {
					key: 'widget',
					commit: {
						committed: false,
						reason: 'tests failed: go test ./...',
						tests: [
							{
								name: 'go test ./...',
								ok: false,
								exit_code: 1,
								seconds: 1,
								output: '--- FAIL: TestAdd'
							}
						]
					}
				}
			}
		];
		c = mount(ImportPrompt, { target: document.body, props: { project: candidate } });
		await settle();
		button('Import').click();
		await settle();
		expect(posts).toEqual(['/api/import?project=import-widget']);
		expect(document.querySelector('[data-testid="import-not-committed"]')?.textContent).toContain(
			'not committed'
		);
		expect(text()).toContain('--- FAIL: TestAdd');
		// flai serve no longer offers it: it stays on screen until the operator moves on
		projects = [];
		await projectState.refresh();
		expect(projectState.project?.key).toBe('import-widget');
	});

	it('a committed import switches to the project it became, once it is served', async () => {
		importAnswer = [
			{ event: 'progress', msg: 'importing' },
			{
				event: 'done',
				committed: true,
				result: { key: 'widget', commit: { committed: true, commit: 'abc1234', tests: [] } }
			}
		];
		projects = [{ key: 'widget', name: 'widget', connected: true }];
		c = mount(ImportPrompt, { target: document.body, props: { project: candidate } });
		await settle();
		button('Import').click();
		await settle();
		expect(text()).toContain('Imported and committed as abc1234');
		expect(projectState.current).toBe('widget');
	});
});
