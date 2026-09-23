import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { inboxState } from './inbox.svelte';
import { projectState, resetForTests } from './project.svelte';

/** Stands in for the browser's EventSource: records where each one listens and whether it closed. */
class FakeEventSource {
	static opened: FakeEventSource[] = [];
	closed = false;
	constructor(public url: string) {
		FakeEventSource.opened.push(this);
	}
	addEventListener() {}
	close() {
		this.closed = true;
	}
}

describe('inboxState across a project switch (S-0095)', () => {
	beforeEach(() => {
		resetForTests();
		FakeEventSource.opened = [];
		vi.stubGlobal('EventSource', FakeEventSource);
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
		expect(FakeEventSource.opened.map((e) => e.url)).toEqual(['/api/events?project=alpha']);

		projectState.pick('beta', false);
		inboxState.restart();
		expect(FakeEventSource.opened[0].closed).toBe(true);
		expect(FakeEventSource.opened.map((e) => e.url)).toEqual([
			'/api/events?project=alpha',
			'/api/events?project=beta'
		]);
		// the new project's inbox is asked, not the old one's kept on show
		await vi.waitFor(() =>
			expect(inboxState.data?.entries[0].key).toBe('review:/api/inbox?project=beta')
		);
	});

	it('restart does nothing before the inbox has started (the login page)', () => {
		inboxState.restart();
		expect(FakeEventSource.opened).toEqual([]);
	});
});
