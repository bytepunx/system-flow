// S-0173: a story's page shows its narrative's open questions, answerable there, and marks and
// brings into view the one an inbox link names.
import { afterEach, describe, expect, it, vi } from 'vitest';
import { flushSync, mount, unmount } from 'svelte';

const api = vi.fn();
vi.mock('$lib/api', () => ({ api: (...args: unknown[]) => api(...args) }));

import OpenQuestions from './OpenQuestions.svelte';
import { inboxState, type Inbox } from '$lib/inbox.svelte';

const json = (body: unknown, status = 200) => ({
	ok: status < 400,
	status,
	statusText: String(status),
	json: async () => body
});
const settle = async () => {
	for (let i = 0; i < 5; i++) await new Promise((r) => setTimeout(r, 0));
	flushSync();
};

const question = (story: string, key: string, title: string) => ({
	key: `question:${story}:${key}`,
	kind: 'question' as const,
	title,
	detail: `asked in the narrative of ${story}`,
	href: `/items/${story}?question=question%3A${story}%3A${key}`,
	item: story
});
const BOX: Inbox = {
	total: 3,
	counts: { thread: 0, question: 3, review: 0, blocked: 0, overlap: 0 },
	notes: [],
	entries: [
		question('S-0001', 'a', 'Which port should it use?'),
		question('S-0001', 'b', 'Which host?'),
		question('S-0002', 'c', 'Another story’s question')
	]
};

describe('a story’s open questions (S-0173)', () => {
	let c: ReturnType<typeof mount> | undefined;
	afterEach(() => {
		if (c) unmount(c);
		c = undefined;
		api.mockReset();
		inboxState.data = null;
		document.body.innerHTML = '';
	});

	it('lists only this story’s questions, and none at all when it has none', () => {
		inboxState.data = BOX;
		c = mount(OpenQuestions, { target: document.body, props: { story: 'S-0001' } });
		flushSync();
		const items = [...document.querySelectorAll('[data-question]')].map((li) => li.textContent);
		expect(items).toHaveLength(2);
		expect(items.join()).not.toContain('Another story');
		// not writable: nothing to answer with
		expect(document.querySelector('[data-testid="question-answer-input"]')).toBeNull();
		unmount(c);
		c = mount(OpenQuestions, { target: document.body, props: { story: 'S-0009' } });
		flushSync();
		expect(document.querySelector('[data-testid="open-questions"]')).toBeNull();
	});

	it('marks the question a link names and brings it into view', async () => {
		const seen: string[] = [];
		Element.prototype.scrollIntoView = function (this: Element) {
			seen.push(this.getAttribute('data-question') ?? '');
		};
		inboxState.data = BOX;
		c = mount(OpenQuestions, {
			target: document.body,
			props: { story: 'S-0001', select: 'question:S-0001:b' }
		});
		await settle();
		const marked = document.querySelector('[aria-current="true"]')!;
		expect(marked.getAttribute('data-question')).toBe('question:S-0001:b');
		expect(seen).toEqual(['question:S-0001:b']);
		expect(document.querySelector('[data-testid="question-gone"]')).toBeNull();
	});

	it('says so when the question a link names is answered already', () => {
		inboxState.data = { ...BOX, entries: [] };
		c = mount(OpenQuestions, {
			target: document.body,
			props: { story: 'S-0001', select: 'question:S-0001:b' }
		});
		flushSync();
		expect(document.querySelector('[data-testid="question-gone"]')).not.toBeNull();
	});

	it('answers a question on the page, posting to its story and refreshing the inbox', async () => {
		inboxState.data = BOX;
		api.mockResolvedValueOnce(json({}));
		api.mockResolvedValueOnce(json({ ...BOX, entries: BOX.entries.slice(1) }));
		c = mount(OpenQuestions, {
			target: document.body,
			props: { story: 'S-0001', writable: true }
		});
		flushSync();
		const input = document.querySelector<HTMLInputElement>(
			'[data-testid="question-answer-input"]'
		)!;
		input.value = 'Nine.';
		input.dispatchEvent(new Event('input', { bubbles: true }));
		flushSync();
		document.querySelector<HTMLButtonElement>('[data-testid="question-answer-submit"]')!.click();
		await settle();
		expect(api).toHaveBeenCalledWith(
			'/api/streams/S-0001/answer',
			expect.objectContaining({
				body: JSON.stringify({ question: 'Which port should it use?', answer: 'Nine.' })
			})
		);
		expect(api).toHaveBeenLastCalledWith('/api/inbox');
		expect(document.querySelectorAll('[data-question]')).toHaveLength(1);
	});
});
