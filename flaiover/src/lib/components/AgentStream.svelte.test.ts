import { afterEach, describe, expect, it, vi } from 'vitest';
import { flushSync, mount, unmount } from 'svelte';
import AgentStream from './AgentStream.svelte';
import type { AgentStreamRead } from '$lib/activity';

const api = vi.fn();
vi.mock('$lib/api', () => ({ api: (...args: unknown[]) => api(...args) }));

// Microtask turns only: the tests that fake setTimeout step the stream's own timers themselves.
const settle = async () => {
	for (let i = 0; i < 20; i++) await Promise.resolve();
	flushSync();
};
const answer = (body: unknown) => ({ ok: true, status: 200, json: async () => body });
const refusal = (status: number, error: string) => ({
	ok: false,
	status,
	statusText: 'x',
	json: async () => ({ error })
});
const read = (over: Partial<AgentStreamRead>): AgentStreamRead => ({
	story: 'S-0142',
	agent: 'agent-S-0142',
	started: '2026-09-29T05:42:53Z',
	running: true,
	from: 0,
	next: 100,
	size: 100,
	entries: [],
	...over
});
const box = () => document.querySelector<HTMLDetailsElement>('[data-testid="agent-stream"]')!;
const lines = () =>
	[...box().querySelectorAll('li')].map((li) => li.textContent!.replace(/\s+/g, ' ').trim());

describe('AgentStream (S-0142)', () => {
	let c: ReturnType<typeof mount> | undefined;
	afterEach(() => {
		if (c) unmount(c);
		c = undefined;
		api.mockReset();
		vi.useRealTimers();
		document.body.innerHTML = '';
	});
	const show = async (props: { started?: string; open?: boolean } = {}) => {
		c = mount(AgentStream, {
			target: document.body,
			props: { story: 'S-0142', started: '2026-09-29T05:42:53Z', ...props }
		});
		await settle();
	};

	it('shows what the agent said and did, and reads on from where it stopped while it runs', async () => {
		vi.useFakeTimers({ toFake: ['setTimeout', 'clearTimeout'] });
		api.mockResolvedValueOnce(
			answer(
				read({
					from: 4000,
					entries: [
						{ kind: 'text', text: 'Priming the session.' },
						{ kind: 'tool', tool: 'Bash', text: 'List files' },
						{ kind: 'result', text: 'no such file', error: true }
					]
				})
			)
		);
		await show();
		expect(api).toHaveBeenCalledWith('/api/agent-stream/S-0142');
		expect(lines()).toEqual([
			"… earlier entries are in the agent's log on the host",
			'Priming the session.',
			'▸ Bash List files',
			'↳ no such file'
		]);
		expect(box().querySelector('[data-kind="result"]')!.className).toContain('text-danger');
		expect(box().dataset.running).toBe('true');
		expect(box().textContent).toContain('live');

		// two seconds later it asks from the offset the last read gave, and appends
		api.mockResolvedValueOnce(
			answer(
				read({
					from: 100,
					next: 180,
					running: false,
					entries: [{ kind: 'end', text: 'session ended (success)' }]
				})
			)
		);
		await vi.advanceTimersByTimeAsync(2000);
		await settle();
		expect(api).toHaveBeenLastCalledWith('/api/agent-stream/S-0142?after=100');
		expect(lines().at(-1)).toBe('end session ended (success)');
		expect(lines()).toHaveLength(5);
		// the agent has ended: nothing more is asked
		await vi.advanceTimersByTimeAsync(10000);
		expect(api).toHaveBeenCalledTimes(2);
		expect(box().dataset.running).toBe('false');
	});

	it('reads on at once while the log has more, and starts over for another run', async () => {
		vi.useFakeTimers({ toFake: ['setTimeout', 'clearTimeout'] });
		api
			.mockResolvedValueOnce(answer(read({ more: true, entries: [{ kind: 'text', text: 'one' }] })))
			.mockResolvedValueOnce(
				answer(read({ from: 100, next: 200, entries: [{ kind: 'text', text: 'two' }] }))
			);
		await show();
		await vi.advanceTimersByTimeAsync(1);
		await settle();
		expect(api).toHaveBeenNthCalledWith(2, '/api/agent-stream/S-0142?after=100');
		expect(lines()).toEqual(['one', 'two']);

		// the story's agent was started again: the offset means nothing in its log
		api
			.mockResolvedValueOnce(
				answer(read({ started: '2026-09-29T06:00:00Z', from: 200, next: 200, entries: [] }))
			)
			.mockResolvedValueOnce(
				answer(
					read({
						started: '2026-09-29T06:00:00Z',
						entries: [{ kind: 'session', text: 'session started' }]
					})
				)
			);
		await vi.advanceTimersByTimeAsync(2000);
		await settle();
		await vi.advanceTimersByTimeAsync(1);
		await settle();
		expect(api).toHaveBeenNthCalledWith(4, '/api/agent-stream/S-0142');
		expect(lines()).toEqual(['session session started']);
	});

	it('shows an agent that has ended once, closed when told', async () => {
		api.mockResolvedValueOnce(
			answer(read({ running: false, entries: [{ kind: 'output', text: 'Error: no API key' }] }))
		);
		await show({ open: false });
		expect(box().open).toBe(false);
		expect(lines()).toEqual(['output Error: no API key']);
		expect(box().textContent).not.toContain('live');
		await settle();
		expect(api).toHaveBeenCalledTimes(1);
	});

	it('says what flai answered when it cannot read, and tries again unless there is no agent', async () => {
		vi.useFakeTimers({ toFake: ['setTimeout', 'clearTimeout'] });
		api.mockResolvedValueOnce(refusal(500, 'method not found: agent.stream'));
		await show();
		expect(box().querySelector('[role="alert"]')!.textContent).toBe(
			'method not found: agent.stream'
		);
		expect(lines()).toEqual([]);
		api.mockResolvedValueOnce(answer(read({ entries: [{ kind: 'text', text: 'back' }] })));
		await vi.advanceTimersByTimeAsync(10000);
		await settle();
		expect(box().querySelector('[role="alert"]')).toBeNull();
		expect(lines()).toEqual(['back']);
		unmount(c!);
		c = undefined;
		document.body.innerHTML = '';

		api.mockReset();
		api.mockResolvedValueOnce(refusal(404, 'flai serve has started no agent for S-0142'));
		await show();
		expect(box().querySelector('[role="alert"]')!.textContent).toContain('no agent');
		await vi.advanceTimersByTimeAsync(30000);
		expect(api).toHaveBeenCalledTimes(1);
	});
});
