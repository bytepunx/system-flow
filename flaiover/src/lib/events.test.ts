import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { debounced, follow, GATHER_MS, listen } from './events';
import { projectState, resetForTests } from './project.svelte';

/** Stands in for the browser's EventSource: where it listens, what for, and whether it closed. */
class FakeEventSource {
	static opened: FakeEventSource[] = [];
	listeners: Record<string, (e: { data: string }) => void> = {};
	closed = false;
	constructor(public url: string) {
		FakeEventSource.opened.push(this);
	}
	addEventListener(kind: string, f: (e: { data: string }) => void) {
		this.listeners[kind] = f;
	}
	close() {
		this.closed = true;
	}
	emit(kind: string, data: unknown) {
		this.listeners[kind]?.({ data: JSON.stringify(data) });
	}
}

describe('listen (S-0154)', () => {
	beforeEach(() => {
		resetForTests();
		FakeEventSource.opened = [];
		vi.stubGlobal('EventSource', FakeEventSource);
	});
	afterEach(() => vi.unstubAllGlobals());

	it('serves every listener on the page from one stream per project, closed with the last', () => {
		projectState.pick('alpha', false);
		const paths: string[] = [];
		const stories: string[] = [];
		const stopA = listen({ change: (p) => paths.push('a:' + p) });
		const stopB = listen({
			change: (p, kind) => paths.push(`b:${p}:${kind}`),
			agent: (s) => stories.push(s)
		});
		expect(FakeEventSource.opened.map((e) => e.url)).toEqual(['/api/events?project=alpha']);
		const source = FakeEventSource.opened[0];

		source.emit('change', { path: 'wip/threads/TH-0001-x.md', kind: 'thread' });
		source.emit('agent', { story: 'S-0154' });
		expect(paths).toEqual(['a:wip/threads/TH-0001-x.md', 'b:wip/threads/TH-0001-x.md:thread']);
		expect(stories).toEqual(['S-0154']);

		stopA();
		stopA(); // twice is once
		expect(source.closed).toBe(false);
		stopB();
		expect(source.closed).toBe(true);
		// the next listener opens a stream of its own
		listen({})();
		expect(FakeEventSource.opened).toHaveLength(2);
	});

	it("listens to the picked project's stream", () => {
		projectState.pick('beta', false);
		const stop = listen({});
		expect(FakeEventSource.opened.map((e) => e.url)).toEqual(['/api/events?project=beta']);
		stop();
	});

	it('does nothing where there is no EventSource', () => {
		vi.stubGlobal('EventSource', undefined);
		expect(() => listen({ change: () => {} })()).not.toThrow();
	});
});

describe('debounced (S-0161)', () => {
	afterEach(() => vi.useRealTimers());

	it('calls once, with every value, when calls within ms of each other settle', () => {
		vi.useFakeTimers();
		const f = vi.fn();
		const later = debounced<string>(f);
		later('a');
		vi.advanceTimersByTime(GATHER_MS - 1);
		later('b');
		vi.advanceTimersByTime(GATHER_MS - 1);
		later();
		vi.advanceTimersByTime(GATHER_MS - 1);
		expect(f).not.toHaveBeenCalled();
		vi.advanceTimersByTime(1);
		expect(f).toHaveBeenCalledTimes(1);
		expect(f).toHaveBeenLastCalledWith(['a', 'b']);
		later('c');
		vi.advanceTimersByTime(GATHER_MS);
		expect(f).toHaveBeenCalledTimes(2);
		expect(f).toHaveBeenLastCalledWith(['c']);
	});

	it('does not wait past maxMs for a burst that never settles', () => {
		vi.useFakeTimers();
		const f = vi.fn();
		const later = debounced(f, 300, 1000);
		for (let t = 0; t < 1000; t += 200) {
			later(t);
			vi.advanceTimersByTime(200);
		}
		expect(f).toHaveBeenCalledTimes(1);
		expect(f).toHaveBeenLastCalledWith([0, 200, 400, 600, 800]);
	});

	it('calls not at all once stopped', () => {
		vi.useFakeTimers();
		const f = vi.fn();
		const later = debounced(f, 300);
		later();
		later.stop();
		vi.advanceTimersByTime(300);
		expect(f).not.toHaveBeenCalled();
	});
});

describe('follow (S-0161)', () => {
	beforeEach(() => {
		resetForTests();
		FakeEventSource.opened = [];
		vi.stubGlobal('EventSource', FakeEventSource);
		vi.useFakeTimers();
	});
	afterEach(() => {
		vi.useRealTimers();
		vi.unstubAllGlobals();
	});

	it('calls once for the changes of its kinds that arrive together, and ignores the rest', () => {
		const f = vi.fn();
		const stop = follow(['item'], f);
		const source = FakeEventSource.opened[0];
		source.emit('change', { path: 'wip/agents/S-0161.md', kind: 'narrative' });
		source.emit('change', { path: 'wip/threads/TH-0001-x.md', kind: 'thread' });
		source.emit('change', { path: 'design/system/overview.md', kind: 'document' });
		vi.advanceTimersByTime(GATHER_MS);
		expect(f).not.toHaveBeenCalled();

		source.emit('change', { path: 'wip/kanban/stories/S-0001-a.md', kind: 'item' });
		vi.advanceTimersByTime(100);
		source.emit('change', { path: 'wip/kanban/board.md', kind: 'item' });
		vi.advanceTimersByTime(100);
		source.emit('change', { path: 'wip/agents/S-0161.md', kind: 'narrative' });
		vi.advanceTimersByTime(GATHER_MS);
		expect(f).toHaveBeenCalledTimes(1);
		expect(f).toHaveBeenLastCalledWith([
			{ path: 'wip/kanban/stories/S-0001-a.md', kind: 'item' },
			{ path: 'wip/kanban/board.md', kind: 'item' }
		]);
		stop();
		expect(source.closed).toBe(true);
	});

	it('counts the manifest, a file the server could not tell apart, and a change without a kind', () => {
		const f = vi.fn();
		const stop = follow(['adr'], f);
		const source = FakeEventSource.opened[0];
		for (const data of [
			{ path: 'system-flow.yaml', kind: 'project' },
			{ path: 'x', kind: 'other' },
			{ path: 'y' }
		]) {
			source.emit('change', data);
			vi.advanceTimersByTime(GATHER_MS);
		}
		expect(f).toHaveBeenCalledTimes(3);
		stop();
	});

	it('asks nothing once stopped, even for changes it had gathered', () => {
		const f = vi.fn();
		const stop = follow(['item'], f);
		FakeEventSource.opened[0].emit('change', { path: 'wip/kanban/board.md', kind: 'item' });
		stop();
		vi.advanceTimersByTime(GATHER_MS);
		expect(f).not.toHaveBeenCalled();
	});
});
