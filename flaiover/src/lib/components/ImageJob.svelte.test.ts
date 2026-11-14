// S-0757: ImageJob shows one image job, its images and the way it left each, its state, outcome,
// step, and times, and its output in a LogFeed read through readLog: open at first for a running
// job, which it follows from where the last read stopped about every second and once more when it
// ends, and for a failed one with output; closed, reading nothing, for a queued job and one that
// passed. A failed job shows its cause and failed step with Retry, which calls the page's onretry
// and shows its refusal word for word.
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { flushSync, mount, unmount } from 'svelte';
import type { ImageJob as ImageJobRecord, ImageLog } from '$lib/images';
import ImageJob, { imageJobAnchor } from './ImageJob.svelte';

const settle = async () => {
	for (let i = 0; i < 10; i++) await Promise.resolve();
	flushSync();
};
const job = (over: Partial<ImageJobRecord>): ImageJobRecord => ({
	id: 'ij-0001',
	images: ['flai-agent:1.94.0', 'acme-agent:1.94.0'],
	state: 'running',
	ways: [],
	requested: '2026-11-13T10:00:00Z',
	started: '2026-11-13T10:00:02Z',
	log_size: 0,
	...over
});
const log = (over: Partial<ImageLog>): ImageLog => ({
	id: 'ij-0001',
	state: 'running',
	exists: true,
	offset: 0,
	next: 0,
	size: 0,
	text: '',
	...over
});

let target: HTMLElement;
let component: ReturnType<typeof mount> | undefined;
let props: Record<string, unknown>;
const readLog = vi.fn<(id: string, offset: number) => Promise<ImageLog>>();

function show(over: Record<string, unknown> = {}) {
	const shown = $state({ job: job({}), readLog, ...over });
	props = shown;
	component = mount(ImageJob, { target, props: shown as never });
	flushSync();
}
const q = (id: string) => target.querySelector(`[data-testid="${id}"]`);
const text = () => q('log-feed-text')!.textContent;
const feed = () => q('log-feed') as HTMLDetailsElement;
const offsets = () => readLog.mock.calls.map(([, o]) => o);

beforeEach(() => {
	vi.useFakeTimers({ toFake: ['setTimeout', 'clearTimeout'] });
	target = document.createElement('div');
	document.body.appendChild(target);
	readLog.mockReset();
});
afterEach(() => {
	if (component) unmount(component);
	component = undefined;
	target.remove();
	vi.useRealTimers();
});

describe('ImageJob', () => {
	it('follows a running job from where it stopped, reads once more when it ends, and stops', async () => {
		readLog
			.mockResolvedValueOnce(log({ next: 4, size: 4, text: 'one\n' }))
			.mockResolvedValueOnce(log({ offset: 4, next: 8, size: 8, text: 'two\n' }))
			.mockResolvedValueOnce(
				log({ offset: 8, next: 14, size: 14, text: 'three\n', state: 'ended', outcome: 'passed' })
			);
		show({ job: job({ step: 'agent image: building acme-agent:1.94.0' }) });
		await settle();
		expect(q('image-job-status')!.textContent).toBe('running');
		expect(q('image-job-step')!.textContent).toBe('agent image: building acme-agent:1.94.0');
		expect(target.querySelector('article')!.id).toBe('image-job-ij-0001');
		expect(imageJobAnchor('ij-0001')).toBe('image-job-ij-0001');
		expect(feed().open).toBe(true);
		expect(text()).toBe('one\n');
		vi.advanceTimersByTime(1000);
		await settle();
		expect(text()).toBe('one\ntwo\n');
		props.job = job({
			state: 'ended',
			outcome: 'passed',
			ended: '2026-11-13T10:02:00Z',
			log_size: 14
		});
		flushSync();
		await settle();
		// the output stays open once the job has passed
		expect(feed().open).toBe(true);
		expect(text()).toBe('one\ntwo\nthree\n');
		expect(offsets()).toEqual([0, 4, 8]);
		expect(readLog.mock.calls.every(([id]) => id === 'ij-0001')).toBe(true);
		vi.advanceTimersByTime(10000);
		await settle();
		expect(readLog).toHaveBeenCalledTimes(3);
	});

	it('starts the text over on a reset, and passes over a read from another offset', async () => {
		readLog
			.mockResolvedValueOnce(log({ next: 4, size: 4, text: 'old\n' }))
			.mockResolvedValueOnce(log({ offset: 0, next: 4, size: 4, text: 'old\n' }))
			.mockResolvedValueOnce(log({ offset: 0, next: 4, size: 4, text: 'new\n', reset: true }));
		show();
		await settle();
		vi.advanceTimersByTime(1000);
		await settle();
		expect(text()).toBe('old\n');
		vi.advanceTimersByTime(1000);
		await settle();
		expect(text()).toBe('new\n');
	});

	it('shows an ended job that passed, the way it left each image, and reads its log once opened', async () => {
		readLog.mockResolvedValueOnce(
			log({ next: 5, size: 5, text: 'done\n', state: 'ended', outcome: 'passed' })
		);
		show({
			job: job({
				state: 'ended',
				outcome: 'passed',
				ended: '2026-11-13T10:02:05Z',
				duration_ms: 123000,
				ways: [
					{ image: 'flai-agent:1.94.0', way: 'pulled' },
					{ image: 'acme-agent:1.94.0', way: 'layered' }
				],
				log_size: 5
			}),
			onretry: vi.fn()
		});
		await settle();
		expect(feed().open).toBe(false);
		expect(readLog).not.toHaveBeenCalled();
		feed().open = true;
		feed().dispatchEvent(new Event('toggle'));
		flushSync();
		await settle();
		expect(q('image-job-status')!.textContent).toBe('passed');
		expect(q('image-job-status')!.className).toContain('text-good-strong');
		expect(q('image-job-duration')!.textContent).toBe('2m 03s');
		expect(q('image-job-ended')).not.toBeNull();
		const rows = [...target.querySelectorAll('[data-testid="image-job-image"]')].map((li) =>
			li.textContent!.replace(/\s+/g, ' ').trim()
		);
		expect(rows).toEqual(['flai-agent:1.94.0 pulled', 'acme-agent:1.94.0 layered']);
		expect(text()).toBe('done\n');
		expect(q('image-job-failure')).toBeNull();
		expect(q('image-job-retry')).toBeNull();
		vi.advanceTimersByTime(10000);
		await settle();
		expect(readLog).toHaveBeenCalledTimes(1);
	});

	it('keeps closed, reading nothing, a stopped job that wrote no output', async () => {
		show({ job: job({ state: 'ended', outcome: 'stopped', cause: 'flai serve stopped it' }) });
		await settle();
		expect(readLog).not.toHaveBeenCalled();
		expect(feed().open).toBe(false);
		expect(q('image-job-cause')!.textContent).toBe('flai serve stopped it');
	});

	it('reads nothing of a queued job', async () => {
		show({ job: job({ state: 'queued', started: undefined }) });
		await settle();
		vi.advanceTimersByTime(5000);
		await settle();
		expect(readLog).not.toHaveBeenCalled();
		expect(q('image-job-status')!.textContent).toBe('queued');
		expect(feed().open).toBe(false);
		expect(feed().dataset.outcome).toBe('pending');
	});

	it('shows a failed job its cause, its failed step, and its output, and Retry calls onretry', async () => {
		readLog.mockResolvedValueOnce(
			log({
				next: 31,
				size: 31,
				text: '#4 ERROR: apt-get install: 100\n',
				state: 'ended',
				outcome: 'failed'
			})
		);
		const onretry = vi.fn(async () => {});
		show({
			job: job({
				state: 'ended',
				outcome: 'failed',
				step: 'RUN apt-get install -y git',
				cause: 'docker build exited with status 1',
				ended: '2026-11-13T10:00:40Z',
				duration_ms: 38000,
				log_size: 31
			}),
			onretry
		});
		await settle();
		expect(q('image-job-status')!.className).toContain('text-danger');
		expect(q('image-job-cause')!.textContent).toBe('docker build exited with status 1');
		expect(q('image-job-failed-step')!.textContent).toBe('RUN apt-get install -y git');
		expect(q('image-job-step')).toBeNull();
		expect(text()).toBe('#4 ERROR: apt-get install: 100\n');
		(q('image-job-retry') as HTMLButtonElement).click();
		await settle();
		expect(onretry).toHaveBeenCalledTimes(1);
		expect(q('image-job-retry-error')).toBeNull();
	});

	it("shows a retry's refusal word for word, and disables Retry while the host action is off", async () => {
		readLog.mockResolvedValue(log({ state: 'ended', outcome: 'failed' }));
		const onretry = vi.fn(async () => {
			throw new Error('the image job ij-0002, which makes flai-agent:1.94.0, runs since 10:05');
		});
		show({ job: job({ state: 'ended', outcome: 'failed', cause: 'x', log_size: 1 }), onretry });
		await settle();
		(q('image-job-retry') as HTMLButtonElement).click();
		await settle();
		expect(q('image-job-retry-error')!.textContent!.trim()).toBe(
			'the image job ij-0002, which makes flai-agent:1.94.0, runs since 10:05'
		);
		props.disabled = true;
		flushSync();
		expect((q('image-job-retry') as HTMLButtonElement).disabled).toBe(true);
	});

	it('links the log a cause or a refusal names to the log page (S-0785)', async () => {
		readLog.mockResolvedValue(log({ state: 'ended', outcome: 'failed' }));
		const ownLog = '/home/op/.flai/serve/image-jobs/ij-0001.log';
		const busyLog = '/home/op/.flai/serve/image-jobs/ij-0002.log';
		const cause = `docker build exited with status 1; read its log ${ownLog}, mend what it names`;
		const onretry = vi.fn(async () => {
			throw new Error(`the image job ij-0002 runs; its log is ${busyLog}`);
		});
		show({ job: job({ state: 'ended', outcome: 'failed', cause, log_size: 1 }), onretry });
		await settle();
		const links = (id: string) =>
			[...q(id)!.querySelectorAll<HTMLAnchorElement>('[data-testid="log-link"]')].map((a) => ({
				path: a.dataset.path,
				page: new URL(a.href).pathname
			}));
		expect(q('image-job-cause')!.textContent).toBe(cause);
		expect(links('image-job-cause')).toEqual([{ path: ownLog, page: '/ops/log' }]);
		(q('image-job-retry') as HTMLButtonElement).click();
		await settle();
		expect(links('image-job-retry-error')).toEqual([{ path: busyLog, page: '/ops/log' }]);
	});

	it("shows a failed job's log read refused word for word", async () => {
		readLog.mockRejectedValueOnce(new Error('images.log: no image job "ij-0001"'));
		show({ job: job({ state: 'ended', outcome: 'failed', log_size: 10 }) });
		await settle();
		expect(q('log-feed-error')!.textContent).toBe('images.log: no image job "ij-0001"');
	});
});
