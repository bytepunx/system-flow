import { describe, expect, it } from 'vitest';
import { EventEmitter } from 'node:events';
import { sseStream } from './sse';
import type { Repo } from './repo';

function fakeRepo(): Repo {
	return new EventEmitter() as unknown as Repo;
}

async function readChunk(reader: ReadableStreamDefaultReader<Uint8Array>): Promise<string> {
	const { value } = await reader.read();
	return new TextDecoder().decode(value);
}

describe('sseStream', () => {
	it('sends ready then change events', async () => {
		const r = fakeRepo();
		const ac = new AbortController();
		const reader = sseStream(r, ac.signal, 60000).getReader();
		expect(await readChunk(reader)).toContain('event: ready');
		r.emit('change', 'design/x.md');
		expect(await readChunk(reader)).toContain('data: {"path":"design/x.md"}');
		r.emit('change', 'wip/agents/S-0161.md', 'narrative');
		expect(await readChunk(reader)).toContain(
			'data: {"path":"wip/agents/S-0161.md","kind":"narrative"}'
		);
		ac.abort();
		expect((await reader.read()).done).toBe(true);
		expect(r.listenerCount('change')).toBe(0);
	});
	it("sends a story's agent starting or ending as an agent event (S-0154)", async () => {
		const r = fakeRepo();
		const ac = new AbortController();
		const reader = sseStream(r, ac.signal, 60000).getReader();
		await readChunk(reader);
		r.emit('agent', 'S-0154');
		const got = await readChunk(reader);
		expect(got).toContain('event: agent');
		expect(got).toContain('data: {"story":"S-0154"}');
		ac.abort();
		expect(r.listenerCount('agent')).toBe(0);
	});
	it('survives cancel followed by abort and abort followed by cancel', async () => {
		const r = fakeRepo();
		const ac = new AbortController();
		const stream = sseStream(r, ac.signal, 60000);
		await stream.cancel();
		expect(() => ac.abort()).not.toThrow();
		const r2 = fakeRepo();
		const ac2 = new AbortController();
		const stream2 = sseStream(r2, ac2.signal, 60000);
		ac2.abort();
		await expect(stream2.cancel()).resolves.toBeUndefined();
		expect(r2.listenerCount('change')).toBe(0);
		r2.emit('change', 'ignored'); // no controller access after close
	});
});
