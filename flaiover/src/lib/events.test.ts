import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { debounced, listen } from './events';
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
		const stopB = listen({ change: (p) => paths.push('b:' + p), agent: (s) => stories.push(s) });
		expect(FakeEventSource.opened.map((e) => e.url)).toEqual(['/api/events?project=alpha']);
		const source = FakeEventSource.opened[0];

		source.emit('change', { path: 'wip/threads/TH-0001-x.md' });
		source.emit('agent', { story: 'S-0154' });
		expect(paths).toEqual(['a:wip/threads/TH-0001-x.md', 'b:wip/threads/TH-0001-x.md']);
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

describe('debounced', () => {
	it('calls once when the calls settle, and not at all once stopped', () => {
		vi.useFakeTimers();
		const f = vi.fn();
		const later = debounced(f, 300);
		later();
		later();
		vi.advanceTimersByTime(299);
		expect(f).not.toHaveBeenCalled();
		vi.advanceTimersByTime(1);
		expect(f).toHaveBeenCalledTimes(1);
		later();
		later.stop();
		vi.advanceTimersByTime(300);
		expect(f).toHaveBeenCalledTimes(1);
		vi.useRealTimers();
	});
});
