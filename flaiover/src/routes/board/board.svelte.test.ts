// S-0161: the board asks for its board, the Publish banner, and the unpushed notice again only when
// a work item changes, once for changes that arrive together; the agents also when a thread does.
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { flushSync, mount, unmount } from 'svelte';
import { GATHER_MS } from '$lib/events';
import { resetForTests } from '$lib/project.svelte';

const api = vi.fn();
vi.mock('$lib/api', () => ({ api: (...args: unknown[]) => api(...args) }));
vi.mock('$app/paths', () => ({ resolve: (route: string) => route }));

import BoardPage from './+page.svelte';

/** Stands in for the browser's EventSource, as /api/events speaks through it. */
class FakeEventSource {
	static opened: FakeEventSource[] = [];
	listeners: Record<string, (e: { data: string }) => void> = {};
	constructor(public url: string) {
		FakeEventSource.opened.push(this);
	}
	addEventListener(kind: string, f: (e: { data: string }) => void) {
		this.listeners[kind] = f;
	}
	close() {}
	emit(kind: string, data: unknown) {
		this.listeners[kind]?.({ data: JSON.stringify(data) });
	}
}

const answer = (body: unknown) => ({ ok: true, json: async () => body });
const board = { wip_limits: {}, order: [], writable: false, columns: {} };
const settle = async (ms = 0) => {
	await new Promise((r) => setTimeout(r, ms));
	for (let i = 0; i < 10; i++) await new Promise((r) => setTimeout(r, 0));
	flushSync();
};

describe('the board follows the work items (S-0161)', () => {
	let c: ReturnType<typeof mount> | undefined;
	beforeEach(() => {
		resetForTests();
		FakeEventSource.opened = [];
		vi.stubGlobal('EventSource', FakeEventSource);
		api.mockImplementation(async (url: string) => {
			if (url === '/api/board') return answer(board);
			if (url === '/api/publish') return answer({ plans: [], push_enabled: false });
			if (url === '/api/unpushed') return answer({ unpushed: null });
			return answer({ enabled: false });
		});
	});
	afterEach(() => {
		if (c) unmount(c);
		c = undefined;
		api.mockReset();
		vi.unstubAllGlobals();
		document.body.innerHTML = '';
	});

	const counts = () => {
		const n = (url: string) => api.mock.calls.filter(([u]) => u === url).length;
		return {
			board: n('/api/board'),
			publish: n('/api/publish'),
			unpushed: n('/api/unpushed'),
			agents: n('/api/host-agent')
		};
	};
	const open = async () => {
		c = mount(BoardPage, { target: document.body });
		await settle();
		// the unpushed notice asks when it is drawn and again once the first board is in
		expect(counts()).toEqual({ board: 1, publish: 1, unpushed: 2, agents: 1 });
		api.mockClear();
		return FakeEventSource.opened[0];
	};

	it('asks nothing again when a narrative or a document changes', async () => {
		const events = await open();
		events.emit('change', { path: 'wip/agents/S-0161.md', kind: 'narrative' });
		events.emit('change', { path: 'design/system/overview.md', kind: 'document' });
		events.emit('change', { path: 'design/adrs/0052-x.md', kind: 'adr' });
		await settle(GATHER_MS + 50);
		expect(counts()).toEqual({ board: 0, publish: 0, unpushed: 0, agents: 0 });
	});

	it('asks only for the agents when a thread changes', async () => {
		const events = await open();
		events.emit('change', { path: 'wip/threads/TH-0001-x.md', kind: 'thread' });
		await settle(GATHER_MS + 50);
		expect(counts()).toEqual({ board: 0, publish: 0, unpushed: 0, agents: 1 });
	});

	it('asks once for everything when work items change together', async () => {
		const events = await open();
		events.emit('change', { path: 'wip/kanban/stories/S-0001-a.md', kind: 'item' });
		events.emit('change', { path: 'wip/kanban/tasks/T-0001-a.md', kind: 'item' });
		await settle(GATHER_MS / 2);
		events.emit('change', { path: 'wip/kanban/board.md', kind: 'item' });
		events.emit('change', { path: 'wip/agents/S-0001.md', kind: 'narrative' });
		await settle(GATHER_MS / 2);
		expect(counts()).toEqual({ board: 0, publish: 0, unpushed: 0, agents: 0 });
		await settle(GATHER_MS);
		expect(counts()).toEqual({ board: 1, publish: 1, unpushed: 1, agents: 1 });
	});

	it('asks for everything when the manifest changes', async () => {
		const events = await open();
		events.emit('change', { path: 'system-flow.yaml', kind: 'project' });
		await settle(GATHER_MS + 50);
		expect(counts()).toEqual({ board: 1, publish: 1, unpushed: 1, agents: 1 });
	});

	it("asks for the agents when flai serve says a story's agent started or ended", async () => {
		const events = await open();
		events.emit('agent', { story: 'S-0161' });
		await settle(GATHER_MS + 50);
		expect(counts()).toEqual({ board: 0, publish: 0, unpushed: 0, agents: 1 });
	});
});

// S-0174: while the clone lags its remote's release tags, the done lane leaves out the archived
// stories it shows only because they look unpublished, and the banner says why.
describe('the done lane of a clone missing published tags (S-0174)', () => {
	let c: ReturnType<typeof mount> | undefined;
	const card = (id: string, archived: boolean) => ({
		id,
		type: 'story',
		title: id,
		nature: 'feature',
		status: 'done',
		blocked: false,
		age_seconds: 0,
		archived
	});
	const withDone = {
		...board,
		columns: { done: [card('S-0001', true), card('S-0002', false)] }
	};
	const remote = {
		remote: 'origin',
		behind: [{ component: 'cli', local: 'cli/v1.0.0', remote: 'cli/v1.4.0' }],
		fix: 'git fetch --tags origin',
		message: 'missing'
	};
	const open = async (publish: unknown) => {
		api.mockImplementation(async (url: string) => {
			if (url === '/api/board') return answer(withDone);
			if (url === '/api/publish') return answer(publish);
			if (url === '/api/unpushed') return answer({ unpushed: null });
			return answer({ enabled: false });
		});
		c = mount(BoardPage, { target: document.body });
		await settle();
	};
	const shown = () =>
		[...document.querySelectorAll('[data-card]')].map((e) => e.getAttribute('data-card'));
	beforeEach(() => {
		resetForTests();
		vi.stubGlobal('EventSource', FakeEventSource);
	});
	afterEach(() => {
		if (c) unmount(c);
		c = undefined;
		api.mockReset();
		vi.unstubAllGlobals();
		document.body.innerHTML = '';
	});

	it('leaves out the archived cards and says how to fetch the tags', async () => {
		await open({ plans: [], remote, unplanned: [], push_enabled: true });
		expect(shown()).toEqual(['S-0002']);
		const banner = document.querySelector('[data-testid="publish-missing-tags"]')!.textContent!;
		expect(banner).toContain('cli/v1.4.0');
		expect(banner).toContain('git fetch --tags origin');
	});

	it('shows them while the clone is in step', async () => {
		await open({
			plans: [
				{
					component: { name: 'cli' },
					level: 'minor',
					from: '1.0.0',
					to: '1.1.0',
					items: [{ id: 'S-0001', title: 'S-0001', level: 'minor' }]
				}
			],
			remote: null,
			unplanned: [],
			push_enabled: true
		});
		expect(shown()).toEqual(['S-0001', 'S-0002']);
		expect(document.querySelector('[data-testid="publish-missing-tags"]')).toBeNull();
	});
});
