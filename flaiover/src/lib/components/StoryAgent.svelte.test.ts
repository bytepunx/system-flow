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
	const show = async (stories: unknown, enabled = true) => {
		api.mockResolvedValue(answer({ enabled, state: { command: '', stories } }));
		c = mount(StoryAgent, { target: document.body, props: { story: 'S-0104' } });
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
		expect(document.querySelector('[data-testid="story-agent-failed"]')!.textContent).toContain(
			'flow-S-0104.log'
		);
	});
});
