import { afterEach, describe, expect, it } from 'vitest';
import { flushSync, mount, unmount } from 'svelte';
import ForgeStatus, { type PullRequestStatus } from './ForgeStatus.svelte';

const STATUS: PullRequestStatus = {
	story: 'S-0394',
	pull_request: 7,
	url: 'https://github.com/bytepunx/system-flow/pull/7',
	head_sha: '0123456789abcdef0123456789abcdef01234567',
	read_at: '2026-10-17T10:00:00Z',
	checks: [
		{
			name: 'flai',
			state: 'failed',
			url: 'https://github.com/bytepunx/system-flow/actions/runs/1',
			summary: 'go test failed in internal/forge',
			since: '2026-10-17T09:58:00Z'
		},
		{ name: 'flaiover', state: 'pending', since: '2026-10-17T09:57:00Z' },
		{
			name: 'template',
			state: 'passed',
			url: 'https://github.com/bytepunx/system-flow/actions/runs/2',
			since: '2026-10-17T09:56:00Z'
		}
	],
	reviews: [
		{
			author: 'alex',
			state: 'changes-requested',
			url: 'https://github.com/bytepunx/system-flow/pull/7#pullrequestreview-1',
			body: 'Name the check.',
			submitted_at: '2026-10-17T09:59:00Z'
		},
		{ author: 'sam', state: 'approved', submitted_at: '2026-10-17T09:55:00Z' }
	]
};

describe('ForgeStatus', () => {
	let component: ReturnType<typeof mount> | undefined;
	const render = (status: PullRequestStatus | null) => {
		component = mount(ForgeStatus, { target: document.body, props: { status } });
		flushSync();
	};
	afterEach(() => {
		if (component) unmount(component);
		component = undefined;
		document.body.innerHTML = '';
	});
	const text = (el: Element) => (el.textContent ?? '').replace(/\s+/g, ' ').trim();

	it('renders nothing when the story has no pull request', () => {
		render(null);
		expect(document.querySelector('[data-testid="forge-status"]')).toBeNull();
		expect(document.body.textContent).toBe('');
	});

	it('links to the pull request on the forge, with its head commit and when it was read', () => {
		render(STATUS);
		const link = document.querySelector<HTMLAnchorElement>('h2 a')!;
		expect(text(link)).toBe('Pull request #7');
		expect(link.getAttribute('href')).toBe(STATUS.url);
		expect(link.getAttribute('target')).toBe('_blank');
		const line = text(document.querySelector('[data-testid="forge-status"] p')!);
		expect(line).toContain('0123456789ab');
		expect(line).not.toContain('0123456789abc');
		// in the local zone (S-0329): the tests run in New York
		expect(line).toContain('2026-10-17 06:00 EDT');
	});

	it("lists each required check's name and state, its log link when it has one, and its summary", () => {
		render(STATUS);
		const rows = [...document.querySelectorAll('[data-check]')];
		expect(rows.map(text)).toEqual([
			'failed flai log go test failed in internal/forge',
			'pending flaiover',
			'passed template log'
		]);
		const failed = rows[0].querySelector('[data-state]')!;
		expect(failed.getAttribute('data-state')).toBe('failed');
		expect(failed.className).toContain('text-danger');
		expect(rows[1].querySelector('[data-state]')!.className).toContain('text-warn');
		expect(rows[2].querySelector('[data-state]')!.className).toContain('text-good');
		const log = rows[0].querySelector<HTMLAnchorElement>('a')!;
		expect(log.getAttribute('href')).toBe(STATUS.checks[0].url);
		expect(log.getAttribute('aria-label')).toBe('flai log');
		expect(rows[1].querySelector('a')).toBeNull();
		expect(rows[0].closest('ul')).not.toBeNull();
	});

	it("lists each review's author and state, and its link when it has one", () => {
		render(STATUS);
		const rows = [...document.querySelectorAll('[data-review]')];
		expect(rows.map(text)).toEqual(['changes requested alex review', 'approved sam']);
		expect(rows[0].querySelector('[data-state]')!.className).toContain('text-danger');
		expect(rows[1].querySelector('[data-state]')!.className).toContain('text-good');
		expect(rows[0].querySelector('a')!.getAttribute('href')).toBe(STATUS.reviews[0].url);
		expect(rows[1].querySelector('a')).toBeNull();
	});

	it('says so when there are no required checks and no reviews', () => {
		render({ ...STATUS, checks: [], reviews: [] });
		const body = text(document.querySelector('[data-testid="forge-status"]')!);
		expect(body).toContain('Required checks None required.');
		expect(body).toContain('Reviews None yet.');
	});
});
