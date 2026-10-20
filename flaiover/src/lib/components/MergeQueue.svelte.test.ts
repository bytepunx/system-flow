// S-0398: the merge queue above the review lane, and each card in review's mark from it; S-0399:
// its hold switch, and the batch's step, result, and drop; S-0400: a red batch's bisection, and
// what each batched story comes to in the batch; S-0402: a story the main checkout's uncommitted
// files hold back, with its files and the command that frees them.
import { afterEach, describe, expect, it, vi } from 'vitest';
import { flushSync, mount, unmount } from 'svelte';
import MergeQueue, {
	batchMark,
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

	describe('a red batch’s bisection and landing (S-0400)', () => {
		const failed = {
			passed: false,
			stopped_at: 'unit',
			summary: 'TestRebase failed',
			finished_at: '2026-10-19T09:52:00Z'
		};
		const red = {
			stories: ['S-0001', 'S-0002', 'S-0003', 'S-0004'],
			started_at: '2026-10-19T09:30:00Z',
			result: failed
		};
		const marks = () =>
			[...el('merge-queue-batch')!.querySelectorAll('[data-story]')].map((s) => {
				const next = s.nextElementSibling;
				const mark = next?.getAttribute('data-testid') === 'merge-queue-mark' ? next : null;
				return `${s.textContent} ${mark?.textContent ?? ''}`.trim().replace(/,$/, '');
			});
		const steps = () =>
			[...el('merge-queue-bisect')!.querySelectorAll('li')].map((li) =>
				li.textContent!.replace(/\s+/g, ' ').trim()
			);
		const passedStep = {
			stories: ['S-0001', 'S-0002'],
			commit: '0123456789abcdef',
			started_at: '2026-10-19T09:53:00Z',
			result: { passed: true, finished_at: '2026-10-19T10:00:00Z' }
		};

		it('shows each step of a bisection under way: one passed, one running', () => {
			render({
				...running,
				batch: {
					...red,
					step: 'bisecting 2 of 4',
					bisect: {
						steps: [
							passedStep,
							{
								stories: ['S-0001', 'S-0002', 'S-0003'],
								commit: '89abcdef01234567',
								started_at: '2026-10-19T10:01:00Z'
							}
						],
						good: 2,
						bad: 4
					}
				}
			});
			expect(steps()).toEqual([
				'Bisect 1: main with S-0001, S-0002 at 0123456, passed',
				'Bisect 2: main with S-0001, S-0002, S-0003 at 89abcde, running'
			]);
			expect(el('merge-queue-bisect')!.getAttribute('aria-label')).toBe('Bisection');
			expect(el('merge-queue-failing')).toBeNull();
			expect(el('merge-queue-mark')).toBeNull();
			expect(text('merge-queue-batch')).toMatch(
				/^Batch S-0001, S-0002, S-0003, S-0004 at step bisecting 2 of 4 /
			);
		});

		it('names the story a bisection found as sent back, and the stories after it as still enqueued', () => {
			render({
				...running,
				batch: {
					...red,
					step: 'landing S-0001',
					bisect: {
						steps: [
							passedStep,
							{
								stories: ['S-0001', 'S-0002', 'S-0003'],
								commit: '89abcdef01234567',
								started_at: '2026-10-19T10:01:00Z',
								result: failed
							}
						],
						good: 2,
						bad: 3,
						failing: 'S-0003',
						commit: '89abcdef01234567',
						result: failed
					},
					landing: { stories: ['S-0001', 'S-0002'], landed: [], head: 'fedcba9876543210' }
				}
			});
			expect(steps()[1]).toBe(
				'Bisect 2: main with S-0001, S-0002, S-0003 at 89abcde, failed at unit: TestRebase failed'
			);
			expect(text('merge-queue-failing')).toBe(
				'The batch fails with S-0003 at 89abcde: sent back to in-progress.'
			);
			expect(marks()).toEqual([
				'S-0001 (to land)',
				'S-0002 (to land)',
				'S-0003 (sent back)',
				'S-0004 (still enqueued)'
			]);
			const sentBack = document.querySelectorAll('[data-testid="merge-queue-mark"]')[2];
			expect(sentBack.className).toContain('text-danger');
		});

		it('marks the stories a landing has landed and those it has still to land', () => {
			render({
				...running,
				batch: {
					stories: ['S-0001', 'S-0002'],
					started_at: '2026-10-19T09:30:00Z',
					step: 'landing S-0002',
					result: { passed: true, finished_at: '2026-10-19T09:48:00Z' },
					landing: { stories: ['S-0001', 'S-0002'], landed: ['S-0001'], head: 'fedcba98' }
				},
				outcomes: [{ story: 'S-0001', outcome: 'landed', at: '2026-10-19T09:50:00Z' }]
			});
			expect(marks()).toEqual(['S-0001 (landed)', 'S-0002 (to land)']);
			expect(text('merge-queue-batch')).toMatch(
				/^Batch S-0001 \(landed\), S-0002 \(to land\) at step landing S-0002 Passed/
			);
			expect(el('merge-queue-bisect')).toBeNull();
		});

		it('takes a failing outcome since the batch started as sent back, and none from before', () => {
			const batch = { ...red, step: 'verified' };
			const at = (t: string) => [{ story: 'S-0002', outcome: 'failing', at: t }];
			expect(batchMark(batch, 'S-0002', at('2026-10-19T10:30:00Z'))).toBe('sent back');
			expect(batchMark(batch, 'S-0003', at('2026-10-19T10:30:00Z'))).toBe('still enqueued');
			expect(batchMark(batch, 'S-0001', at('2026-10-19T10:30:00Z'))).toBeUndefined();
			expect(batchMark(batch, 'S-0002', at('2026-10-18T10:30:00Z'))).toBeUndefined();
			expect(batchMark(batch, 'S-0002', null)).toBeUndefined();
		});
	});

	describe('a story the main checkout’s uncommitted files hold back (S-0402)', () => {
		const command = "git stash push -- a.txt 'b c.md'";
		const dirty: MergeQueueState = {
			...running,
			ineligible: [
				...running.ineligible!,
				{
					...entry('S-0004', 'Overwrites'),
					reasons: [
						{
							kind: 'dirty_main',
							text: 'S-0004’s change would overwrite a.txt, b c.md, which the main checkout has not committed',
							paths: ['a.txt', 'b c.md'],
							command
						}
					]
				}
			]
		};
		const clipboard = (writeText: unknown) =>
			Object.defineProperty(navigator, 'clipboard', { value: { writeText }, configurable: true });
		afterEach(() => Reflect.deleteProperty(navigator, 'clipboard'));

		it('names the story, its files, and the command that frees them, and says committing frees it too', () => {
			render(dirty);
			const rows = document.querySelectorAll('[data-testid="merge-queue-held-back-row"]');
			expect(rows).toHaveLength(1);
			expect(el('merge-queue-held-back-row')!.dataset.story).toBe('S-0004');
			expect(text('merge-queue-held-back-row')).toBe(`S-0004: a.txt, b c.md ${command} Copy`);
			expect(text('merge-queue-held-back-command')).toBe(command);
			expect(text('merge-queue-held-back')).toMatch(
				/^Held back: .+ Commit them, or set them aside/
			);
			expect(text('merge-queue-ineligible')).toBe(
				'2 not eligible: the reasons are on their cards.'
			);
		});

		it('copies the command to the clipboard', async () => {
			const writeText = vi.fn().mockResolvedValue(undefined);
			clipboard(writeText);
			render(dirty);
			el('merge-queue-held-back-copy')!.click();
			await settle();
			expect(writeText).toHaveBeenCalledWith(command);
			expect(text('merge-queue-held-back-copy')).toBe('Copied');
			expect(el('merge-queue-held-back-copy-error')).toBeNull();
		});

		it('says when the command could not be copied', async () => {
			const writeText = vi.fn().mockRejectedValue(new Error('denied'));
			clipboard(writeText);
			render(dirty);
			el('merge-queue-held-back-copy')!.click();
			await settle();
			expect(text('merge-queue-held-back-copy-error')).toBe(
				'Not copied (denied): select the command and copy it.'
			);
			expect(text('merge-queue-held-back-copy')).toBe('Copy');
		});

		it('lists none when no story is held back', () => {
			render(running);
			expect(el('merge-queue-held-back')).toBeNull();
			expect(el('merge-queue-held-back-row')).toBeNull();
		});

		it('marks a story the landing held back as held back, not to land', () => {
			const batch = {
				stories: ['S-0001', 'S-0002'],
				step: 'landing S-0002',
				landing: {
					stories: ['S-0001', 'S-0002'],
					landed: ['S-0001'],
					head: 'fedcba98',
					held_back: [{ story: 'S-0002', paths: ['a.txt'] }]
				}
			};
			expect(batchMark(batch, 'S-0002')).toBe('held back');
			expect(batchMark(batch, 'S-0001')).toBe('landed');
			render({ ...running, batch });
			const mark = document.querySelectorAll('[data-testid="merge-queue-mark"]')[1];
			expect(mark.textContent).toBe('(held back)');
			expect(mark.className).toContain('text-warn');
		});
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
	it('is true for a batch not dropped, in a queue that runs, being bisected or landed too (S-0400)', () => {
		const batch = { stories: ['S-0003'], step: 'verify' };
		expect(verifying({ ...running, batch })).toBe(true);
		expect(
			verifying({ ...running, batch: { ...batch, result: { passed: true, finished_at: 'x' } } })
		).toBe(true);
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
