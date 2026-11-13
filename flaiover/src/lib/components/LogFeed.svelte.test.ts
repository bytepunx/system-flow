// S-0595 (ADR-0263): LogFeed shows a kluster step's name and a badge of its outcome, reads nothing
// while closed, reads the step's log from the start once opened, then from where it stopped about
// every second while the step runs, once more when it ends, and no more after; a reset starts the
// text over, a read from another offset is passed over, and a refusal is shown word for word; given
// read, it reads through read in place of kluster.log (S-0757).
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { flushSync, mount, unmount } from 'svelte';
import type { KlusterLog } from '$lib/server/repo';

const api = vi.fn();
vi.mock('$lib/api', () => ({ api: (...args: unknown[]) => api(...args) }));

import LogFeed, { type LogPart } from './LogFeed.svelte';

const answer = (body: unknown, status = 200) => ({
	ok: status < 400,
	status,
	statusText: '',
	json: async () => body
});
const settle = async () => {
	for (let i = 0; i < 10; i++) await Promise.resolve();
	flushSync();
};
const log = (over: Partial<KlusterLog>): KlusterLog => ({
	namespace: 's-0595',
	step: 'build-api',
	log: '.flai-cache/kluster/logs/s-0595/build-api.log',
	exists: true,
	offset: 0,
	next: 0,
	size: 0,
	text: '',
	...over
});
const offsets = () =>
	api.mock.calls.map(([p]) => Number(new URL(String(p), 'http://x').searchParams.get('offset')));

let target: HTMLElement;
let component: ReturnType<typeof mount> | undefined;
let props: Record<string, unknown>;

function show(over: Record<string, unknown> = {}) {
	const shown = $state({ namespace: 's-0595', step: 'build-api', label: 'api', ...over });
	props = shown;
	component = mount(LogFeed, { target, props: shown });
	flushSync();
}
const q = (id: string) => target.querySelector(`[data-testid="${id}"]`);
const text = () => q('log-feed-text')!.textContent;

beforeEach(() => {
	vi.useFakeTimers({ toFake: ['setTimeout', 'clearTimeout'] });
	target = document.createElement('div');
	document.body.appendChild(target);
	api.mockReset();
});
afterEach(() => {
	if (component) unmount(component);
	component = undefined;
	target.remove();
	vi.useRealTimers();
});

describe('LogFeed', () => {
	it('shows the step and a badge per outcome', () => {
		show({ outcome: 'pending' });
		expect(target.textContent).toContain('api');
		const cases: [string, string][] = [
			['pending', 'text-muted'],
			['running', 'text-info'],
			['passed', 'text-good-strong'],
			['failed', 'text-danger'],
			['skipped', 'text-muted']
		];
		for (const [o, cls] of cases) {
			props.outcome = o;
			flushSync();
			expect(q('log-feed-status')!.textContent).toBe(o);
			expect(q('log-feed-status')!.className).toContain(cls);
		}
	});

	it('shows how long an ended step took', () => {
		show({ outcome: 'passed', startedAt: '2026-11-07T01:00:00Z', endedAt: '2026-11-07T01:02:05Z' });
		expect(q('log-feed-duration')!.textContent).toBe('2m 05s');
	});

	it('reads nothing while closed', async () => {
		show({ outcome: 'running' });
		await settle();
		vi.advanceTimersByTime(5000);
		await settle();
		expect(api).not.toHaveBeenCalled();
	});

	it('reads from the start once opened, and asks the step it shows', async () => {
		api.mockResolvedValueOnce(
			answer(log({ next: 6, size: 6, text: 'built\n', outcome: 'passed' }))
		);
		show({ outcome: 'passed', open: true });
		await settle();
		expect(api).toHaveBeenCalledTimes(1);
		expect(api.mock.calls[0][0]).toBe('/api/kluster/log?namespace=s-0595&step=build-api&offset=0');
		expect(text()).toBe('built\n');
		// an ended step is not read again
		vi.advanceTimersByTime(10000);
		await settle();
		expect(api).toHaveBeenCalledTimes(1);
	});

	it('follows a running step from where it stopped, then reads once more when it ends, and stops', async () => {
		api
			.mockResolvedValueOnce(answer(log({ next: 4, size: 4, text: 'one\n', outcome: 'running' })))
			.mockResolvedValueOnce(
				answer(log({ offset: 4, next: 8, size: 8, text: 'two\n', outcome: 'running' }))
			)
			.mockResolvedValueOnce(
				answer(log({ offset: 8, next: 14, size: 14, text: 'three\n', outcome: 'passed' }))
			);
		show({ outcome: 'running', open: true });
		await settle();
		expect(text()).toBe('one\n');
		vi.advanceTimersByTime(1000);
		await settle();
		expect(text()).toBe('one\ntwo\n');
		props.outcome = 'passed';
		flushSync();
		await settle();
		expect(text()).toBe('one\ntwo\nthree\n');
		expect(offsets()).toEqual([0, 4, 8]);
		vi.advanceTimersByTime(10000);
		await settle();
		expect(api).toHaveBeenCalledTimes(3);
	});

	it('reads at once again while more of the log is left than one read answers', async () => {
		api
			.mockResolvedValueOnce(answer(log({ next: 2, size: 4, text: 'ab', outcome: 'passed' })))
			.mockResolvedValueOnce(
				answer(log({ offset: 2, next: 4, size: 4, text: 'cd', outcome: 'passed' }))
			);
		show({ outcome: 'passed', open: true });
		await settle();
		vi.advanceTimersByTime(0);
		await settle();
		expect(text()).toBe('abcd');
		expect(offsets()).toEqual([0, 2]);
	});

	it('starts the text over on a reset, and passes over a read from another offset', async () => {
		api
			.mockResolvedValueOnce(answer(log({ next: 4, size: 4, text: 'old\n', outcome: 'running' })))
			.mockResolvedValueOnce(answer(log({ offset: 0, next: 4, size: 4, text: 'old\n' })))
			.mockResolvedValueOnce(
				answer(log({ offset: 0, next: 4, size: 4, text: 'new\n', reset: true, outcome: 'running' }))
			);
		show({ outcome: 'running', open: true });
		await settle();
		vi.advanceTimersByTime(1000);
		await settle();
		// answered from 0 when 4 was asked, with no reset: nothing appended
		expect(text()).toBe('old\n');
		vi.advanceTimersByTime(1000);
		await settle();
		expect(text()).toBe('new\n');
	});

	it('shows a refusal word for word', async () => {
		api.mockResolvedValueOnce(answer({ error: 'kluster.log: the step "x" is not a step' }, 400));
		show({ outcome: 'failed', open: true });
		await settle();
		expect(q('log-feed-error')!.textContent).toBe('kluster.log: the step "x" is not a step');
	});

	it('says a pending step has not started', async () => {
		api.mockResolvedValueOnce(answer(log({ exists: false })));
		show({ outcome: 'pending', open: true });
		await settle();
		expect(q('log-feed-empty')!.textContent).toContain('has not started');
	});

	it('reads through read when given, in place of kluster.log, and names the log it shows', async () => {
		const read = vi
			.fn<(offset: number) => Promise<LogPart>>()
			.mockResolvedValueOnce(log({ next: 4, size: 4, text: 'one\n', outcome: 'running' }))
			.mockResolvedValueOnce(
				log({ offset: 4, next: 8, size: 8, text: 'two\n', outcome: 'running' })
			);
		show({
			namespace: undefined,
			step: 'the image job 20261113T093000Z-001',
			label: 'Output',
			outcome: 'running',
			open: true,
			read
		});
		await settle();
		vi.advanceTimersByTime(1000);
		await settle();
		expect(text()).toBe('one\ntwo\n');
		expect(read.mock.calls.map(([o]) => o)).toEqual([0, 4]);
		expect(api).not.toHaveBeenCalled();
		expect(q('log-feed-text')!.getAttribute('aria-label')).toBe(
			'the log of the image job 20261113T093000Z-001'
		);
	});

	it('stops reading once closed', async () => {
		api.mockResolvedValue(answer(log({ next: 1, size: 1, text: 'x', outcome: 'running' })));
		show({ outcome: 'running', open: true });
		await settle();
		props.open = false;
		flushSync();
		const asked = api.mock.calls.length;
		vi.advanceTimersByTime(5000);
		await settle();
		expect(api.mock.calls.length).toBe(asked);
	});
});
