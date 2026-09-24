import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { flushSync, mount, unmount } from 'svelte';
import StoryAgent from './StoryAgent.svelte';

const api = vi.fn();
vi.mock('$lib/api', () => ({ api: (...args: unknown[]) => api(...args) }));

class FakeEventSource {
	static opened: FakeEventSource[] = [];
	listeners: Record<string, () => void> = {};
	constructor(public url: string) {
		FakeEventSource.opened.push(this);
	}
	addEventListener(kind: string, f: () => void) {
		this.listeners[kind] = f;
	}
	close() {}
}

const settle = async () => {
	for (let i = 0; i < 5; i++) await new Promise((r) => setTimeout(r, 0));
	flushSync();
};
const answer = (body: unknown) => ({ ok: true, json: async () => body });
const section = () => document.querySelector<HTMLElement>('[data-testid="story-agent"]');
const text = () => section()!.textContent!.replace(/\s+/g, ' ');
const run = {
	story: 'S-0104',
	harness: 'claude-code',
	model: 'claude-haiku-4-5',
	command: 'claude',
	agent: 'agent-S-0104',
	started: '2026-09-23T18:00:00Z'
};

describe('StoryAgent (S-0104)', () => {
	let c: ReturnType<typeof mount> | undefined;
	beforeEach(() => {
		FakeEventSource.opened = [];
		vi.stubGlobal('EventSource', FakeEventSource);
	});
	afterEach(() => {
		if (c) unmount(c);
		c = undefined;
		api.mockReset();
		vi.unstubAllGlobals();
		document.body.innerHTML = '';
	});
	const show = async (
		stories: unknown,
		enabled = true,
		props: { status?: string; writable?: boolean } = {}
	) => {
		api.mockResolvedValue(answer({ enabled, state: { command: '', stories } }));
		c = mount(StoryAgent, { target: document.body, props: { story: 'S-0104', ...props } });
		await settle();
	};

	it('shows nothing for a story no agent was started for, or while the action is off', async () => {
		await show({ 'S-0001': { state: 'working', run } });
		expect(section()).toBeNull();
		unmount(c!);
		c = undefined;
		await show({ 'S-0104': { state: 'working', run } }, false);
		expect(section()).toBeNull();
	});

	it('shows the agent at work, and asks again when the project changes', async () => {
		await show({ 'S-0104': { state: 'working', run } });
		expect(text()).toContain('agent working (claude-code, claude-haiku-4-5)');
		expect(text()).toContain('agent-S-0104, started 2026-09-23 18:00 UTC');
		expect(document.querySelector<HTMLElement>('[data-testid="agent-dot"]')!.dataset.state).toBe(
			'working'
		);
		// the agent asks the designer: the page hears of the thread's file changing
		api.mockResolvedValue(
			answer({
				enabled: true,
				state: {
					command: '',
					stories: {
						'S-0104': {
							state: 'waiting',
							why: 'waiting for an answer to TH-0009: Which port?',
							thread: 'TH-0009',
							run
						}
					}
				}
			})
		);
		FakeEventSource.opened[0].listeners.change();
		await settle();
		expect(text()).toContain(
			'agent waiting (claude-code, claude-haiku-4-5): waiting for an answer to TH-0009: Which port?'
		);
		expect(text()).toContain('It asked in TH-0009: answer it below and it goes on.');
	});

	it('says why it failed and where its output is', async () => {
		await show({
			'S-0104': {
				state: 'failed',
				why: 'ended (exit 1) with S-0104 in in-progress',
				run: {
					...run,
					ended: '2026-09-23T18:30:00Z',
					exit: 1,
					log: '/home/me/.flai/serve/agents/flow-S-0104.log'
				}
			}
		});
		expect(text()).toContain(
			'agent failed (claude-code, claude-haiku-4-5): ended (exit 1) with S-0104 in in-progress'
		);
		expect(text()).toContain('ended 2026-09-23 18:30 UTC (exit 1)');
		const failed = document
			.querySelector('[data-testid="story-agent-failed"]')!
			.textContent!.replace(/\s+/g, ' ');
		expect(failed).toContain('flow-S-0104.log');
		// S-0116: either way to have flai start another
		expect(failed).toContain(
			'Moving the story back to ready starts another, and so does changing its agent while it is in ready.'
		);
	});

	// S-0116: a story in ready or in progress whose agent dropped or failed gets a new one
	describe('Restart agent', () => {
		const failed = {
			'S-0104': {
				state: 'failed',
				why: 'ended (exit 1) with S-0104 in in-progress',
				run: { ...run, ended: '2026-09-23T18:30:00Z', exit: 1 }
			}
		};
		const button = () =>
			document.querySelector<HTMLButtonElement>('[data-testid="story-agent-restart"]');

		it('is offered only for a failed agent of a story in ready or in progress, to a writer', async () => {
			for (const [stories, props, offered] of [
				[failed, { status: 'in-progress', writable: true }, true],
				[failed, { status: 'ready', writable: true }, true],
				[failed, { status: 'review', writable: true }, false],
				[failed, { status: 'in-progress', writable: false }, false],
				[{ 'S-0104': { state: 'working', run } }, { status: 'in-progress', writable: true }, false]
			] as const) {
				await show(stories, true, props);
				expect(button() !== null, JSON.stringify(props)).toBe(offered);
				unmount(c!);
				c = undefined;
			}
		});

		it('asks flai to restart it, and says why when flai refuses', async () => {
			await show(failed, true, { status: 'in-progress', writable: true });
			api.mockImplementation(async (path: string, init?: RequestInit) =>
				init?.method === 'POST'
					? { ok: false, json: async () => ({ error: 'the in-progress limit leaves no room' }) }
					: answer({ enabled: true, state: { command: '', stories: failed } })
			);
			button()!.click();
			await settle();
			const post = api.mock.calls.find((call) => (call[1] as RequestInit)?.method === 'POST')!;
			expect(post[0]).toBe('/api/items/S-0104/agent');
			expect(JSON.parse((post[1] as RequestInit).body as string)).toEqual({ action: 'restart' });
			expect(
				document.querySelector('[data-testid="story-agent-restart-error"]')!.textContent
			).toContain('the in-progress limit leaves no room');
			// a restart that is taken: the page asks again and shows the new agent at work
			api.mockImplementation(async (path: string, init?: RequestInit) =>
				init?.method === 'POST'
					? { ok: true, json: async () => ({ story: 'S-0104' }) }
					: answer({
							enabled: true,
							state: { command: '', stories: { 'S-0104': { state: 'working', run } } }
						})
			);
			button()!.click();
			await settle();
			expect(text()).toContain('agent working');
			expect(button()).toBeNull();
		});
	});

	// S-0115: a ready story no agent was started for can have one started now
	describe('Start agent', () => {
		const button = () =>
			document.querySelector<HTMLButtonElement>('[data-testid="story-agent-start"]');
		const failed = {
			'S-0104': {
				state: 'failed',
				why: 'ended (exit 1)',
				run: { ...run, ended: '2026-09-23T18:30:00Z' }
			}
		};

		it('is offered only for a ready story with no agent, to a writer, with the action on', async () => {
			for (const [stories, enabled, props, offered] of [
				[{}, true, { status: 'ready', writable: true }, true],
				[{}, true, { status: 'backlog', writable: true }, false],
				[{}, true, { status: 'in-progress', writable: true }, false],
				[{}, true, { status: 'ready', writable: false }, false],
				[{}, false, { status: 'ready', writable: true }, false],
				[{ 'S-0104': { state: 'working', run } }, true, { status: 'ready', writable: true }, false],
				// Restart agent is the way for one whose agent failed
				[failed, true, { status: 'ready', writable: true }, false]
			] as const) {
				await show(stories, enabled, props);
				expect(button() !== null, JSON.stringify([enabled, props])).toBe(offered);
				unmount(c!);
				c = undefined;
			}
		});

		it('says why flai serve has not started it', async () => {
			api.mockResolvedValue(
				answer({
					enabled: true,
					state: {
						command: '',
						waiting:
							'S-0001 names no harness; the in-progress limit leaves no room for S-0104, S-0105',
						stories: {}
					}
				})
			);
			c = mount(StoryAgent, {
				target: document.body,
				props: { story: 'S-0104', status: 'ready', writable: true }
			});
			await settle();
			expect(text()).toContain(
				'flai serve has started no agent for this story: the in-progress limit leaves no room for S-0104, S-0105'
			);
		});

		it('asks flai to start it, and says why when flai refuses', async () => {
			await show({}, true, { status: 'ready', writable: true });
			api.mockImplementation(async (path: string, init?: RequestInit) =>
				init?.method === 'POST'
					? {
							ok: false,
							json: async () => ({ error: 'S-0104 names no harness, and no command is set' })
						}
					: answer({ enabled: true, state: { command: '', stories: {} } })
			);
			button()!.click();
			await settle();
			const post = api.mock.calls.find((call) => (call[1] as RequestInit)?.method === 'POST')!;
			expect(post[0]).toBe('/api/items/S-0104/agent');
			expect(JSON.parse((post[1] as RequestInit).body as string)).toEqual({ action: 'start' });
			expect(
				document.querySelector('[data-testid="story-agent-start-error"]')!.textContent
			).toContain('names no harness, and no command is set');
			// a start that is taken: the page asks again and shows the agent at work
			api.mockImplementation(async (path: string, init?: RequestInit) =>
				init?.method === 'POST'
					? { ok: true, json: async () => ({ story: 'S-0104' }) }
					: answer({
							enabled: true,
							state: { command: '', stories: { 'S-0104': { state: 'working', run } } }
						})
			);
			button()!.click();
			await settle();
			expect(text()).toContain('agent working');
			expect(button()).toBeNull();
		});
	});
});
