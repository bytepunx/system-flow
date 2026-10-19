// S-0398: the merge queue above the review lane, and each card in review's mark from it.
import { afterEach, describe, expect, it } from 'vitest';
import { flushSync, mount, unmount } from 'svelte';
import MergeQueue, { queueMarks, type MergeQueueState } from './MergeQueue.svelte';

const entry = (story: string, title: string) => ({
	story,
	title,
	entered_review: '2026-10-18T09:00:00Z'
});
const running: MergeQueueState = {
	mode: 'local',
	eligible: [
		{ ...entry('S-0003', 'Oldest in review'), head: 'abc', enqueued: true },
		{ ...entry('S-0001', 'Next in review'), head: 'def' }
	],
	ineligible: [
		{
			...entry('S-0002', 'Unticked'),
			reasons: [{ kind: 'criteria_unticked', text: 'acceptance 2 unticked', criteria: [2] }]
		}
	],
	enqueued: [{ story: 'S-0003', by: 'alex', at: '2026-10-18T10:00:00Z' }],
	batch: { stories: ['S-0003'], step: 'verify' },
	outcomes: []
};
const off: MergeQueueState = {
	mode: 'pull-request',
	off: true,
	off_reason: 'the merge queue runs under local acceptance only',
	eligible: null,
	ineligible: null,
	enqueued: null,
	outcomes: null
};

describe('MergeQueue', () => {
	let component: ReturnType<typeof mount> | undefined;
	afterEach(() => {
		if (component) unmount(component);
		component = undefined;
		document.body.innerHTML = '';
	});
	const render = (queue: MergeQueueState | undefined) => {
		component = mount(MergeQueue, { target: document.body, props: { queue } });
		flushSync();
	};
	const el = (id: string) => document.querySelector<HTMLElement>(`[data-testid="${id}"]`);
	const text = (id: string) => el(id)?.textContent?.replace(/\s+/g, ' ').trim();

	it('lists the eligible stories in queue order, and the batch with its step', () => {
		render(running);
		const items = [...el('merge-queue-order')!.querySelectorAll('li')];
		expect(items.map((li) => li.dataset.story)).toEqual(['S-0003', 'S-0001']);
		expect(items[0].textContent!.replace(/\s+/g, ' ').trim()).toBe(
			'1. S-0003 Oldest in review, enqueued'
		);
		expect(items[1].textContent!.replace(/\s+/g, ' ').trim()).toBe('2. S-0001 Next in review');
		expect(text('merge-queue-batch')).toBe('Batch S-0003 at step verify');
		expect(text('merge-queue-ineligible')).toBe('1 not eligible: the reasons are on its card.');
		expect(el('merge-queue-empty')).toBeNull();
		expect(el('merge-queue')!.getAttribute('aria-label')).toBe('Merge queue');
	});

	it('says the queue is empty and there is no batch', () => {
		render({ ...running, eligible: [], ineligible: [], batch: undefined });
		expect(text('merge-queue-empty')).toBe('Empty: no story in review is eligible.');
		expect(text('merge-queue-batch')).toBe('No batch.');
		expect(el('merge-queue-order')).toBeNull();
		expect(el('merge-queue-ineligible')).toBeNull();
	});

	it('takes null lists from flai as empty', () => {
		render({ mode: 'local', eligible: null, ineligible: null, enqueued: null, outcomes: null });
		expect(text('merge-queue-empty')).toBe('Empty: no story in review is eligible.');
	});

	it('says the queue does not run under pull-request acceptance, with flai’s reason, in place of it', () => {
		render(off);
		expect(text('merge-queue-off')).toBe(
			'The merge queue does not run under pull-request acceptance. Why: the merge queue runs under local acceptance only'
		);
		expect(el('merge-queue-order')).toBeNull();
		expect(el('merge-queue-batch')).toBeNull();
		expect(el('merge-queue-empty')).toBeNull();
	});

	it('says the flai on the host sends no merge queue when the board has none', () => {
		render(undefined);
		expect(text('merge-queue-absent')).toBe(
			'No merge queue: the flai on the host sends none, being older than the queue or unable to read it.'
		);
		expect(el('merge-queue-batch')).toBeNull();
	});
});

describe('queueMarks', () => {
	it('gives each eligible story its place and each other story its reasons', () => {
		expect(queueMarks(running)).toEqual({
			'S-0003': { place: 1 },
			'S-0001': { place: 2 },
			'S-0002': {
				reasons: [{ kind: 'criteria_unticked', text: 'acceptance 2 unticked', criteria: [2] }]
			}
		});
	});

	it('marks no story while the queue is off or absent', () => {
		expect(queueMarks(off)).toEqual({});
		expect(queueMarks(undefined)).toEqual({});
	});
});
