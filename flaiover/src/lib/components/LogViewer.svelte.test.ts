// S-0785: LogViewer shows a log's path, reads its contents through /api/log from the start, then
// from where it stopped about every second while the log grows, says so when the log is not there
// and shows it once written, shows flai's refusal word for word with no badge of an outcome, and
// starts over for another path.
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { flushSync, mount, unmount } from 'svelte';
import type { LogRead } from '$lib/logs';

const api = vi.fn();
vi.mock('$lib/api', () => ({ api: (...args: unknown[]) => api(...args) }));

import LogViewer from './LogViewer.svelte';

const PATH = '.flai-cache/kluster/logs/s-0595/build-api.log';

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
const read = (over: Partial<LogRead>): LogRead => ({
	path: PATH,
	exists: true,
	offset: 0,
	next: 0,
	size: 0,
	text: '',
	...over
});
const asked = (): Record<string, string>[] =>
	api.mock.calls.map(([p]) => {
		const url = new URL(String(p), 'http://x');
		return { pathname: url.pathname, ...Object.fromEntries(url.searchParams) };
	});

let target: HTMLElement;
let component: ReturnType<typeof mount> | undefined;
let props: { path: string };

function show(path = PATH) {
	const shown = $state({ path });
	props = shown;
	component = mount(LogViewer, { target, props: shown });
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

describe('LogViewer', () => {
	it('shows the path and reads the log from the start through /api/log', async () => {
		api.mockResolvedValueOnce(answer(read({ next: 6, size: 6, text: 'built\n' })));
		show();
		await settle();
		expect(q('log-viewer-path')!.textContent).toBe(PATH);
		expect(asked()).toEqual([{ pathname: '/api/log', path: PATH, offset: '0' }]);
		expect(text()).toBe('built\n');
		expect(q('log-feed-status')).toBeNull();
		expect(target.textContent).toContain('live');
	});

	it('follows a growing log from where it stopped', async () => {
		api
			.mockResolvedValueOnce(answer(read({ next: 4, size: 4, text: 'one\n' })))
			.mockResolvedValueOnce(answer(read({ offset: 4, next: 4, size: 4, text: '' })))
			.mockResolvedValueOnce(answer(read({ offset: 4, next: 8, size: 8, text: 'two\n' })));
		show();
		await settle();
		expect(text()).toBe('one\n');
		vi.advanceTimersByTime(1000);
		await settle();
		expect(text()).toBe('one\n');
		vi.advanceTimersByTime(1000);
		await settle();
		expect(text()).toBe('one\ntwo\n');
		expect(asked().map((a) => a.offset)).toEqual(['0', '4', '4']);
	});

	it('says when the log is not there, and shows it once it is written', async () => {
		api
			.mockResolvedValueOnce(answer(read({ exists: false })))
			.mockResolvedValueOnce(answer(read({ next: 6, size: 6, text: 'begun\n' })));
		show();
		await settle();
		expect(q('log-feed-empty')!.textContent).toContain('There is no log at this path');
		vi.advanceTimersByTime(1000);
		await settle();
		expect(text()).toBe('begun\n');
		expect(q('log-feed-empty')).toBeNull();
	});

	it("shows flai's refusal of a path word for word", async () => {
		const refusal =
			'log.read reads a log under the project root, ~/.flai, or the story worktrees, not /etc/passwd';
		api.mockResolvedValueOnce(answer({ error: refusal }, 400));
		show('/etc/passwd');
		await settle();
		expect(q('log-feed-error')!.textContent).toBe(refusal);
		expect(q('log-feed-error')!.getAttribute('role')).toBe('alert');
	});

	it('starts over for another path', async () => {
		api
			.mockResolvedValueOnce(answer(read({ next: 4, size: 4, text: 'one\n' })))
			.mockResolvedValueOnce(
				answer(read({ path: '/tmp/b.log', next: 6, size: 6, text: 'other\n' }))
			);
		show();
		await settle();
		props.path = '/tmp/b.log';
		flushSync();
		await settle();
		expect(q('log-viewer-path')!.textContent).toBe('/tmp/b.log');
		expect(text()).toBe('other\n');
		expect(asked()[1]).toEqual({ pathname: '/api/log', path: '/tmp/b.log', offset: '0' });
	});
});
