// S-0398: the merge queue above the review lane, and each card in review's mark from it; S-0399:
// its hold switch, and the batch's step, result, and drop.
import { afterEach, describe, expect, it, vi } from 'vitest';
import { flushSync, mount, unmount } from 'svelte';
import MergeQueue, {
	queueMarks,
	refusalOf,
	verifying,
	type MergeQueueState
} from './MergeQueue.svelte';

const api = vi.fn();
vi.mock('$lib/api', () => ({ api: (...args: unknown[]) => api(...args) }));
const settle = async () => {
	for (let i = 0; i < 4; i++) await new Promise((r) => setTimeout(r, 0));
	flushSync();
};
const answer = (body: unknown, status = 200) => ({
	ok: status < 400,
	status,
	statusText: String(status),
	json: async () => body
});

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
	const render = (
		queue: MergeQueueState | undefined,
		more: { enabled?: boolean; writable?: boolean; onchanged?: () => void } = {}
	) => {
		component = mount(MergeQueue, { target: document.body, props: { queue, ...more } });
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

	it('shows a passed batch’s result with its summary', () => {
		render({
			...running,
			batch: {
				stories: ['S-0003', 'S-0001'],
				step: 'verified',
				result: { passed: true, summary: 'every tier passed', finished_at: '2026-10-18T11:00:00Z' }
			}
		});
		expect(text('merge-queue-batch')).toMatch(
			/^Batch S-0003, S-0001 at step verified Passed, 2026-10-1\d \d\d:\d\d .+: every tier passed$/
		);
		expect(el('merge-queue-result')!.className).toContain('text-good');
		expect(el('merge-queue-dropped')).toBeNull();
	});

	it('shows a failed batch’s result: the tier it stopped at and why', () => {
		render({
			...running,
			batch: {
				stories: ['S-0003'],
				step: 'verified',
				result: {
					passed: false,
					stopped_at: 'integration',
					summary: 'TestLanding failed',
					finished_at: '2026-10-18T11:00:00Z'
				}
			}
		});
		expect(text('merge-queue-result')).toMatch(
			/^Failed at integration, 2026-10-1\d \d\d:\d\d .+: TestLanding failed$/
		);
		expect(el('merge-queue-result')!.className).toContain('text-danger');
	});

	it('shows a dropped batch’s reason, who dropped it, and when', () => {
		render({
			...running,
			batch: {
				stories: ['S-0003'],
				step: 'verify',
				dropped: {
					reason: 'S-0003 was removed from the queue by alex',
					by: 'alex',
					at: '2026-10-18T11:00:00Z'
				}
			}
		});
		expect(text('merge-queue-dropped')).toMatch(
			/^Dropped by alex, 2026-10-1\d \d\d:\d\d .+: S-0003 was removed from the queue by alex$/
		);
		expect(el('merge-queue-result')).toBeNull();
	});

	it('says who held the queue and when', () => {
		render({ ...running, held: true, held_by: 'alex', held_at: '2026-10-18T10:30:00Z' });
		expect(text('merge-queue-held')).toMatch(
			/^Held by alex at 2026-10-1\d \d\d:\d\d .+: no batch is built until it is resumed\.$/
		);
	});

	describe('the hold switch', () => {
		afterEach(() => api.mockReset());
		const toggle = () => {
			const s = el('merge-queue-hold') as HTMLInputElement;
			s.click();
			flushSync();
		};

		it('holds the queue through flai, then has the page read it again', async () => {
			const onchanged = vi.fn();
			api.mockResolvedValue(answer({ action: 'hold', changed: true, by: 'alex' }));
			render(running, { enabled: true, writable: true, onchanged });
			const s = el('merge-queue-hold') as HTMLInputElement;
			expect(s.getAttribute('role')).toBe('switch');
			expect(s.checked).toBe(false);
			toggle();
			await settle();
			expect(api).toHaveBeenCalledWith('/api/queue', expect.objectContaining({ method: 'POST' }));
			expect(JSON.parse(api.mock.calls[0][1].body)).toEqual({ held: true });
			expect(onchanged).toHaveBeenCalledOnce();
			expect(el('merge-queue-hold-error')).toBeNull();
		});

		it('resumes a held queue', async () => {
			api.mockResolvedValue(answer({ action: 'resume', changed: true, by: 'alex' }));
			render({ ...running, held: true, held_by: 'alex' }, { enabled: true, writable: true });
			expect((el('merge-queue-hold') as HTMLInputElement).checked).toBe(true);
			toggle();
			await settle();
			expect(JSON.parse(api.mock.calls[0][1].body)).toEqual({ held: false });
		});

		it('shows what flai said when the hold is refused, and has nothing read again', async () => {
			const onchanged = vi.fn();
			api.mockResolvedValue(
				answer({ error: 'the host action "merge-queue" is not enabled for this project' }, 403)
			);
			render(running, { enabled: true, writable: true, onchanged });
			toggle();
			await settle();
			expect(text('merge-queue-hold-error')).toBe(
				'the host action "merge-queue" is not enabled for this project'
			);
			expect(onchanged).not.toHaveBeenCalled();
		});

		it('is offered only while the action is on and the dashboard can write, and says what enables it', () => {
			render(running, { enabled: false, writable: true });
			expect(el('merge-queue-hold')).toBeNull();
			expect(text('merge-queue-disabled')).toBe(
				'No batch runs: the operator turns the queue on in a shell on the host with flai serve enable merge-queue.'
			);
			unmount(component!);
			render(running, { enabled: true, writable: false });
			expect(el('merge-queue-hold')).toBeNull();
			expect(el('merge-queue-disabled')).toBeNull();
			unmount(component!);
			render(running);
			expect(el('merge-queue-hold')).toBeNull();
			expect(el('merge-queue-disabled')).toBeNull();
		});
	});
});

describe('verifying', () => {
	it('is true only for a batch with neither a result nor a drop, in a queue that runs', () => {
		const batch = { stories: ['S-0003'], step: 'verify' };
		expect(verifying({ ...running, batch })).toBe(true);
		expect(
			verifying({ ...running, batch: { ...batch, result: { passed: true, finished_at: 'x' } } })
		).toBe(false);
		expect(
			verifying({ ...running, batch: { ...batch, dropped: { reason: 'r', by: 'a', at: 'x' } } })
		).toBe(false);
		expect(verifying({ ...running, batch: undefined })).toBe(false);
		expect(verifying({ ...off, batch })).toBe(false);
		expect(verifying(undefined)).toBe(false);
		expect(verifying(null)).toBe(false);
	});
});

describe('refusalOf', () => {
	it('takes flai’s sentence and reasons, else its error, else the status', () => {
		const reasons = [{ kind: 'verify_missing', text: 'no result' }];
		expect(refusalOf({ text: 'not enqueued', reasons, error: 'x' }, '422')).toEqual({
			text: 'not enqueued',
			reasons
		});
		expect(refusalOf({ error: 'disabled' }, '403')).toEqual({ text: 'disabled', reasons: [] });
		expect(refusalOf(null, 'Bad Gateway')).toEqual({ text: 'Bad Gateway', reasons: [] });
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
