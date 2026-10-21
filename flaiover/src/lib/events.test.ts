import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { debounced, follow, GATHER_MS, listen, RETRY_MAX_MS, RETRY_MS } from './events';
import { projectState, resetForTests } from './project.svelte';

/** Stands in for the browser's WebSocket: where it connects, what the server says on it, and
 * whether it closed. */
class FakeWebSocket {
	static opened: FakeWebSocket[] = [];
	onmessage: ((e: { data: unknown }) => void) | null = null;
	onclose: (() => void) | null = null;
	closed = false;
	constructor(public url: string) {
		FakeWebSocket.opened.push(this);
	}
	close() {
		this.closed = true;
		this.onclose?.();
	}
	/** The server sends a message. */
	emit(message: unknown) {
		this.onmessage?.({ data: JSON.stringify(message) });
	}
	/** The connection drops, as when the server restarts. */
	drop() {
		this.onclose?.();
	}
}

const PAGE = { protocol: 'http:', host: 'dash.test:4242' };
const latest = () => FakeWebSocket.opened[FakeWebSocket.opened.length - 1];

function stubSocket() {
	resetForTests();
	FakeWebSocket.opened = [];
	vi.stubGlobal('WebSocket', FakeWebSocket);
	vi.stubGlobal('location', PAGE);
}

describe('listen (S-0154, ADR-0182)', () => {
	beforeEach(() => {
		stubSocket();
		vi.useFakeTimers();
	});
	afterEach(() => {
		vi.useRealTimers();
		vi.unstubAllGlobals();
	});

	it('serves every listener on the page from one socket per project, closed with the last', () => {
		projectState.pick('alpha', false);
		const paths: string[] = [];
		const stories: string[] = [];
		const stopA = listen({ change: (p) => paths.push('a:' + p) });
		const stopB = listen({
			change: (p, kind) => paths.push(`b:${p}:${kind}`),
			agent: (s) => stories.push(s)
		});
		expect(FakeWebSocket.opened.map((e) => e.url)).toEqual([
			'ws://dash.test:4242/api/events?project=alpha'
		]);
		const socket = FakeWebSocket.opened[0];

		socket.emit({ type: 'ready' });
		socket.emit({ type: 'change', path: 'wip/threads/TH-0001-x.md', kind: 'thread' });
		socket.emit({ type: 'agent', story: 'S-0154' });
		expect(paths).toEqual(['a:wip/threads/TH-0001-x.md', 'b:wip/threads/TH-0001-x.md:thread']);
		expect(stories).toEqual(['S-0154']);

		stopA();
		stopA(); // twice is once
		expect(socket.closed).toBe(false);
		stopB();
		expect(socket.closed).toBe(true);
		// closing it on purpose opens nothing again
		vi.advanceTimersByTime(RETRY_MAX_MS);
		expect(FakeWebSocket.opened).toHaveLength(1);
		// the next listener opens a socket of its own
		listen({})();
		expect(FakeWebSocket.opened).toHaveLength(2);
	});

	it("connects to the picked project's events, wss: on a page served over https:", () => {
		vi.stubGlobal('location', { protocol: 'https:', host: 'dash.test' });
		projectState.pick('beta', false);
		const stop = listen({});
		expect(FakeWebSocket.opened.map((e) => e.url)).toEqual([
			'wss://dash.test/api/events?project=beta'
		]);
		stop();
	});

	it('reads a change without a kind as other, and ignores what it cannot read', () => {
		const seen: string[] = [];
		const stop = listen({
			change: (p, kind) => seen.push(`${p}:${kind}`),
			agent: (s) => seen.push(`agent:${s}`)
		});
		const socket = latest();
		socket.emit({ type: 'change', path: 'y' });
		socket.onmessage?.({ data: 'not json' });
		socket.emit({ type: 'ping' });
		socket.emit(null);
		expect(seen).toEqual(['y:other']);
		stop();
	});

	it('opens a dropped socket again after a backoff that doubles to a cap, reset by ready', () => {
		const stop = listen({});
		const waits: number[] = [];
		for (let i = 0; i < 7; i++) {
			const before = FakeWebSocket.opened.length;
			latest().drop();
			let waited = 0;
			while (FakeWebSocket.opened.length === before) {
				vi.advanceTimersByTime(100);
				waited += 100;
			}
			waits.push(waited);
		}
		expect(waits).toEqual([1000, 2000, 4000, 8000, 16000, RETRY_MAX_MS, RETRY_MAX_MS]);

		latest().emit({ type: 'ready' });
		latest().drop();
		vi.advanceTimersByTime(RETRY_MS - 1);
		const before = FakeWebSocket.opened.length;
		vi.advanceTimersByTime(1);
		expect(FakeWebSocket.opened).toHaveLength(before + 1);
		stop();
	});

	it('tells every listener to ask again once a dropped socket is ready, not on the first ready', () => {
		const a = { change: vi.fn(), agent: vi.fn() };
		const b = { change: vi.fn() };
		const stopA = listen(a);
		const stopB = listen(b);
		latest().emit({ type: 'ready' });
		expect(a.change).not.toHaveBeenCalled();
		expect(b.change).not.toHaveBeenCalled();

		latest().drop();
		vi.advanceTimersByTime(RETRY_MS);
		expect(FakeWebSocket.opened).toHaveLength(2);
		expect(a.change).not.toHaveBeenCalled();
		latest().emit({ type: 'ready' });
		// the manifest affects every kind a page follows, and an agent with no story is any agent
		expect(a.change).toHaveBeenCalledExactlyOnceWith('system-flow.yaml', 'project');
		expect(a.agent).toHaveBeenCalledExactlyOnceWith('');
		expect(b.change).toHaveBeenCalledExactlyOnceWith('system-flow.yaml', 'project');

		// a second ready on the same socket is not another reconnection
		latest().emit({ type: 'ready' });
		expect(a.change).toHaveBeenCalledTimes(1);
		stopA();
		stopB();
	});

	it('stops opening it again once the last listener leaves while it waits', () => {
		const stop = listen({});
		latest().drop();
		stop();
		vi.advanceTimersByTime(RETRY_MAX_MS);
		expect(FakeWebSocket.opened).toHaveLength(1);
	});

	it('does nothing where there is no WebSocket', () => {
		vi.stubGlobal('WebSocket', undefined);
		expect(() => listen({ change: () => {} })()).not.toThrow();
		expect(FakeWebSocket.opened).toEqual([]);
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
		stubSocket();
		vi.useFakeTimers();
	});
	afterEach(() => {
		vi.useRealTimers();
		vi.unstubAllGlobals();
	});

	it('calls once for the changes of its kinds that arrive together, and ignores the rest', () => {
		const f = vi.fn();
		const stop = follow(['item'], f);
		const socket = FakeWebSocket.opened[0];
		socket.emit({ type: 'change', path: 'wip/agents/S-0161.md', kind: 'narrative' });
		socket.emit({ type: 'change', path: 'wip/threads/TH-0001-x.md', kind: 'thread' });
		socket.emit({ type: 'change', path: 'design/system/overview.md', kind: 'document' });
		vi.advanceTimersByTime(GATHER_MS);
		expect(f).not.toHaveBeenCalled();

		socket.emit({ type: 'change', path: 'wip/kanban/stories/S-0001-a.md', kind: 'item' });
		vi.advanceTimersByTime(100);
		socket.emit({ type: 'change', path: 'wip/kanban/board.md', kind: 'item' });
		vi.advanceTimersByTime(100);
		socket.emit({ type: 'change', path: 'wip/agents/S-0161.md', kind: 'narrative' });
		vi.advanceTimersByTime(GATHER_MS);
		expect(f).toHaveBeenCalledTimes(1);
		expect(f).toHaveBeenLastCalledWith([
			{ path: 'wip/kanban/stories/S-0001-a.md', kind: 'item' },
			{ path: 'wip/kanban/board.md', kind: 'item' }
		]);
		stop();
		expect(socket.closed).toBe(true);
	});

	it('counts the manifest, a file the server could not tell apart, and a change without a kind', () => {
		const f = vi.fn();
		const stop = follow(['adr'], f);
		const socket = FakeWebSocket.opened[0];
		for (const data of [
			{ path: 'system-flow.yaml', kind: 'project' },
			{ path: 'x', kind: 'other' },
			{ path: 'y' }
		]) {
			socket.emit({ type: 'change', ...data });
			vi.advanceTimersByTime(GATHER_MS);
		}
		expect(f).toHaveBeenCalledTimes(3);
		stop();
	});

	it('asks nothing once stopped, even for changes it had gathered', () => {
		const f = vi.fn();
		const stop = follow(['item'], f);
		FakeWebSocket.opened[0].emit({ type: 'change', path: 'wip/kanban/board.md', kind: 'item' });
		stop();
		vi.advanceTimersByTime(GATHER_MS);
		expect(f).not.toHaveBeenCalled();
	});

	it('asks again, whatever its kinds, once a dropped socket is ready again', () => {
		const f = vi.fn();
		const stop = follow(['adr'], f);
		latest().drop();
		vi.advanceTimersByTime(RETRY_MS);
		latest().emit({ type: 'ready' });
		vi.advanceTimersByTime(GATHER_MS);
		expect(f).toHaveBeenCalledExactlyOnceWith([{ path: 'system-flow.yaml', kind: 'project' }]);
		stop();
	});
});
