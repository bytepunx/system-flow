import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { inboxState } from './inbox.svelte';
import { projectState, resetForTests } from './project.svelte';

/** Stands in for the browser's WebSocket: records where each one connects and whether it closed. */
class FakeWebSocket {
	static opened: FakeWebSocket[] = [];
	closed = false;
	constructor(public url: string) {
		FakeWebSocket.opened.push(this);
	}
	close() {
		this.closed = true;
	}
}

/** The path and query each socket was opened on. */
const paths = () =>
	FakeWebSocket.opened.map((e) => {
		const u = new URL(e.url);
		return u.pathname + u.search;
	});

describe('inboxState across a project switch (S-0095)', () => {
	beforeEach(() => {
		resetForTests();
		FakeWebSocket.opened = [];
		vi.stubGlobal('WebSocket', FakeWebSocket);
		vi.stubGlobal(
			'fetch',
			vi.fn(
				async (url: string) =>
					new Response(
						JSON.stringify({
							total: 1,
							counts: { thread: 0, question: 0, review: 1, blocked: 0, overlap: 0 },
							entries: [{ key: `review:${url}`, kind: 'review', title: 'x', href: '/' }],
							notes: []
						})
					)
			)
		);
	});
	afterEach(() => {
		inboxState.stop();
		vi.unstubAllGlobals();
	});

	it("restart listens to the picked project's events and drops the old stream", async () => {
		projectState.pick('alpha', false);
		inboxState.start();
		await vi.waitFor(() => expect(inboxState.data).not.toBeNull());
		expect(paths()).toEqual(['/api/events?project=alpha']);

		projectState.pick('beta', false);
		inboxState.restart();
		expect(FakeWebSocket.opened[0].closed).toBe(true);
		expect(paths()).toEqual(['/api/events?project=alpha', '/api/events?project=beta']);
		// the new project's inbox is asked, not the old one's kept on show
		await vi.waitFor(() =>
			expect(inboxState.data?.entries[0].key).toBe('review:/api/inbox?project=beta')
		);
	});

	it('restart does nothing before the inbox has started (the login page)', () => {
		inboxState.restart();
		expect(FakeWebSocket.opened).toEqual([]);
	});
});
