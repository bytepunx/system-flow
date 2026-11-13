import { describe, expect, it } from 'vitest';
import {
	CHECKS_LOG_HREF,
	criteriaOf,
	klusterView,
	patchLines,
	patchRuns,
	readNdjson,
	sectionOf,
	toggleCriterion,
	verifyView,
	type KlusterEvidence,
	type VerifyReport
} from './review';

const STORY = `# S-0041 Review

## Goal
Review in one place.

## Acceptance criteria
- [x] A review page per story
- [ ] Accept runs as the designer
* [X] Failures are shown verbatim

## Tasks
- T-0185 Living design
- [ ] not a criterion, it is under Tasks

## Notes
n
`;

describe('review helpers', () => {
	it('reads the acceptance criteria and their ticks, and nothing from other sections', () => {
		expect(criteriaOf(STORY)).toEqual([
			{ checked: true, text: 'A review page per story' },
			{ checked: false, text: 'Accept runs as the designer' },
			{ checked: true, text: 'Failures are shown verbatim' }
		]);
		expect(criteriaOf('# no criteria\n')).toEqual([]);
	});
	it('returns a section up to the next heading', () => {
		expect(sectionOf(STORY, 'Goal')).toBe('Review in one place.');
		expect(sectionOf(STORY, 'goal')).toBe('Review in one place.');
		expect(sectionOf(STORY, 'Notes')).toBe('n');
		expect(sectionOf(STORY, 'Missing')).toBe('');
	});
	it('classifies patch lines', () => {
		const lines = patchLines('@@ -1,2 +1,2 @@\n one\n-two\n+2\n\\ No newline at end of file');
		expect(lines.map((l) => l.kind)).toEqual(['hunk', 'context', 'del', 'add', 'note']);
		expect(patchLines('')).toEqual([]);
	});
	it('takes the sign of an added or removed line out of its text (S-0164)', () => {
		const lines = patchLines(
			'@@ -1,3 +1,3 @@\n one\n-\ttwo\n+-2\n \n\\ No newline at end of file\n'
		);
		expect(lines).toEqual([
			{ kind: 'hunk', sign: '', text: '@@ -1,3 +1,3 @@' },
			{ kind: 'context', sign: '', text: 'one' },
			{ kind: 'del', sign: '-', text: '\ttwo' },
			{ kind: 'add', sign: '+', text: '-2' },
			{ kind: 'context', sign: '', text: '' },
			{ kind: 'note', sign: '', text: '\\ No newline at end of file' }
		]);
	});
	it('gathers consecutive added lines and removed lines into runs (S-0164)', () => {
		const runs = patchRuns(
			'@@ -1,4 +1,5 @@\n-a\n-b\n+c\n+d\n+e\n same\n+f\n@@ -9,2 +10,1 @@\n-g\n last\n'
		);
		expect(runs.map((r) => [r.kind, r.lines.map((l) => l.text)])).toEqual([
			['hunk', ['@@ -1,4 +1,5 @@']],
			['del', ['a', 'b']],
			['add', ['c', 'd', 'e']],
			['context', ['same']],
			['add', ['f']],
			['hunk', ['@@ -9,2 +10,1 @@']],
			['del', ['g']],
			['context', ['last']]
		]);
		expect(patchRuns('')).toEqual([]);
	});
	it('delivers NDJSON lines as they arrive, across chunk boundaries', async () => {
		const enc = new TextEncoder();
		const chunks = [
			'{"event":"progress","st',
			'ep":"merged"}\n{"event":"pro',
			'gress","step":"done"}\nnot json\n{"event":"done"}'
		];
		const body = new ReadableStream<Uint8Array>({
			start(c) {
				for (const ch of chunks) c.enqueue(enc.encode(ch));
				c.close();
			}
		});
		const got: Record<string, unknown>[] = [];
		await readNdjson(new Response(body), (l) => got.push(l));
		expect(got).toEqual([
			{ event: 'progress', step: 'merged' },
			{ event: 'progress', step: 'done' },
			{ event: 'done' }
		]);
	});
});

describe('toggleCriterion (S-0085)', () => {
	const body = [
		'## Goal',
		'- [ ] not a criterion',
		'',
		'## Acceptance criteria',
		'- [ ] first',
		'text between',
		'* [x] second',
		'',
		'## Notes',
		'- [ ] nor this'
	].join('\n');
	it('ticks and unticks the nth criterion and nothing else', () => {
		expect(toggleCriterion(body, 0)).toBe(body.replace('- [ ] first', '- [x] first'));
		expect(toggleCriterion(body, 1)).toBe(body.replace('* [x] second', '* [ ] second'));
	});
	it('knows when there is no such criterion', () => {
		expect(toggleCriterion(body, 2)).toBeNull();
		expect(toggleCriterion(body, -1)).toBeNull();
		expect(toggleCriterion('## Goal\n- [ ] x\n', 0)).toBeNull();
	});
});

describe('verifyView (S-0270)', () => {
	const failing: VerifyReport = {
		story: 'S-0270',
		commit: '0123456789abcdef0123',
		base: 'main',
		ran_at: '2026-10-06T22:00:00Z',
		duration_ms: 42100,
		duration: '42.1s',
		passed: false,
		stopped_at: 'flaiover',
		steps: [
			{ name: 'rebase', state: 'passed', duration_ms: 3, duration: '3ms' },
			{ name: 'check', state: 'passed', duration_ms: 1500 },
			{
				name: 'flaiover',
				tier: true,
				state: 'failed',
				duration_ms: 40000,
				duration: '40s',
				findings: [
					{
						name: 'Review > shows it',
						path: 'flaiover/src/lib/x.test.ts',
						line: 12,
						message: 'expected 1\nto be 2'
					},
					{ message: 'exit status 1' }
				],
				omitted: 3
			},
			{ name: 'flai', tier: true, state: 'not-reached', duration_ms: 0 }
		],
		notes: [
			{
				rule: 'stale-narrative',
				level: 'warning',
				path: 'wip/agents/S-0001.md',
				line: 4,
				message: 'not touched in a week'
			}
		]
	};

	it('makes a row per step in order, with its state, duration, and the failed step’s findings', () => {
		const v = verifyView(failing);
		expect(v.rows.map((r) => [r.name, r.tier, r.label, r.duration])).toEqual([
			['rebase', false, 'passed', '3ms'],
			['check', false, 'passed', '1.5s'],
			['flaiover', true, 'failed', '40s'],
			['flai', true, 'not reached', '']
		]);
		expect(v.rows[2].findings).toEqual([
			{ where: 'flaiover/src/lib/x.test.ts:12 Review > shows it', message: 'expected 1\nto be 2' },
			{ where: '', message: 'exit status 1' }
		]);
		expect(v.rows[2].omitted).toBe(3);
		expect(v.rows[0].findings).toEqual([]);
	});

	it('says where it stopped, the commit short, when it ran, and the notes apart', () => {
		const v = verifyView(failing);
		expect(v).toMatchObject({
			passed: false,
			outcome: 'Stopped at flaiover',
			commit: '0123456789ab',
			base: 'main',
			ranAt: '2026-10-06T22:00:00Z',
			duration: '42.1s'
		});
		expect(v.notes).toEqual([
			{ where: 'wip/agents/S-0001.md:4 stale-narrative', message: 'warning: not touched in a week' }
		]);
	});

	it('shows a covered tier as skipped, covered by the tier, with no duration (S-0342)', () => {
		const v = verifyView({
			...failing,
			passed: true,
			stopped_at: undefined,
			steps: [
				{ name: 'integration', tier: true, state: 'passed', duration_ms: 61000, duration: '1m1s' },
				{
					name: 'go-short',
					tier: true,
					state: 'skipped',
					covered_by: 'integration',
					duration_ms: 0
				},
				{ name: 'old', tier: true, state: 'skipped', duration_ms: 0 }
			]
		});
		expect(v.rows.map((r) => [r.name, r.state, r.label, r.duration])).toEqual([
			['integration', 'passed', 'passed', '1m1s'],
			['go-short', 'skipped', 'skipped: covered by integration', ''],
			['old', 'skipped', 'skipped', '']
		]);
	});

	it('shows a tier left to the merge queue as queued, with no duration (ADR-0174)', () => {
		const v = verifyView({
			...failing,
			passed: true,
			stopped_at: undefined,
			steps: [
				{ name: 'go-test', tier: true, state: 'passed', duration_ms: 12000, duration: '12s' },
				{ name: 'smoke', tier: true, state: 'queued', duration_ms: 0 }
			]
		});
		expect(v.rows.map((r) => [r.name, r.state, r.label, r.duration])).toEqual([
			['go-test', 'passed', 'passed', '12s'],
			['smoke', 'queued', 'left to the merge queue', '']
		]);
	});

	it('says a passing report passed every step, with no findings and no notes', () => {
		const v = verifyView({
			...failing,
			passed: true,
			stopped_at: undefined,
			duration: undefined,
			duration_ms: 340,
			steps: [{ name: 'rebase', state: 'passed', duration_ms: 3, findings: [{ message: 'x' }] }],
			notes: undefined
		});
		expect(v.outcome).toBe('Passed every step');
		expect(v.duration).toBe('340ms');
		expect(v.rows[0].findings).toEqual([]);
		expect(v.notes).toEqual([]);
	});

	it('shows a reused tier with the time of the run it comes from in place of its duration, and no findings (S-0341)', () => {
		const v = verifyView({
			...failing,
			passed: false,
			steps: [
				{ name: 'rebase', state: 'passed', duration_ms: 3, duration: '3ms' },
				{
					name: 'flai',
					tier: true,
					state: 'reused',
					duration_ms: 61000,
					duration: '1m1s',
					reused_from: '2026-10-08T11:41:07Z',
					findings: [{ message: 'from the earlier run' }],
					omitted: 2,
					command: ['go', 'test', './...'],
					dir: 'flai',
					exit_code: 0
				},
				{ name: 'flaiover', tier: true, state: 'failed', duration_ms: 40000, duration: '40s' }
			]
		});
		expect(v.rows.map((r) => [r.name, r.tier, r.state, r.label, r.duration])).toEqual([
			['rebase', false, 'passed', 'passed', '3ms'],
			['flai', true, 'reused', 'reused', 'from 2026-10-08 07:41 EDT'],
			['flaiover', true, 'failed', 'failed', '40s']
		]);
		expect(v.rows[1].findings).toEqual([]);
		expect(v.rows[1].omitted).toBe(0);
	});

	it('falls back to a reused tier’s duration when it does not say where it comes from', () => {
		const v = verifyView({
			...failing,
			steps: [{ name: 'flai', tier: true, state: 'reused', duration_ms: 1500 }]
		});
		expect(v.rows[0]).toMatchObject({ label: 'reused', duration: '1.5s', findings: [] });
	});
});

describe('klusterView (S-0662)', () => {
	const HEAD = 'c0ffee000000000000000000000000000000beef';
	const OLDER = '0dd000000000000000000000000000000000c0de';
	const atHead = { commit: HEAD, head: HEAD, at_head: true, exact: true };
	const AT_HEAD: KlusterEvidence = {
		applies: true,
		story: 'S-0662',
		namespace: 's-0662',
		head: HEAD,
		checks_on: true,
		deploy: {
			commit: HEAD,
			outcome: 'passed',
			job: 'deploy',
			said: 'passed',
			asked_by: 'agent-S-0662',
			started_at: '2026-11-13T09:29:00Z',
			ended_at: '2026-11-13T09:30:00Z',
			steps: [
				{ name: 'build', outcome: 'passed', log: '.flai-cache/kluster/logs/s-0662/build.log' },
				{ name: 'deploy', outcome: 'passed', log: '.flai-cache/kluster/logs/s-0662/deploy.log' }
			],
			head: atHead,
			passed_at_head: true
		},
		checks: {
			commit: HEAD,
			outcome: 'passed',
			log: '.flai-cache/checks/S-0662.log',
			head: atHead,
			passed_at_head: true
		},
		passed: true
	};

	it('shows a deploy and checks run that passed at the head, each with its log', () => {
		const v = klusterView(AT_HEAD, '/ops/deployments')!;
		expect(v).toMatchObject({
			namespace: 's-0662',
			head: 'c0ffee000000',
			passed: true,
			headline: 'Its namespace is deployed and checked at the head c0ffee000000'
		});
		expect(v.deploy).toEqual({
			kind: 'deploy',
			label: 'Deploy',
			commit: 'c0ffee000000',
			outcome: 'passed',
			tone: 'passed',
			asked: 'asked by agent-S-0662',
			endedAt: '2026-11-13T09:30:00Z',
			atHead: true,
			head: 'at the head',
			logHref: '/ops/deployments#logs',
			logLabel: "its steps' logs"
		});
		expect(v.checks).toEqual({
			kind: 'checks',
			label: 'Checks',
			commit: 'c0ffee000000',
			outcome: 'passed',
			tone: 'passed',
			asked: 'who asked is not recorded',
			endedAt: '',
			atHead: true,
			head: 'at the head',
			logHref: CHECKS_LOG_HREF,
			logLabel: 'its log'
		});
		expect(CHECKS_LOG_HREF).toBe('#review-checks');
	});

	it('counts a commit followed only by design, docs, and wip changes as at the head', () => {
		const since = { commit: OLDER, head: HEAD, at_head: true };
		const v = klusterView(
			{
				...AT_HEAD,
				deploy: { ...AT_HEAD.deploy!, commit: OLDER, head: since },
				checks: { ...AT_HEAD.checks!, commit: OLDER, head: since }
			},
			'/ops/deployments'
		)!;
		expect(v.deploy).toMatchObject({ commit: '0dd000000000', atHead: true });
		expect(v.deploy!.head).toBe(
			'at the head: from 0dd000000000, followed only by changes under design, docs, and wip'
		);
	});

	it('says why a deploy and checks run at an older commit are not at the head', () => {
		const why =
			'1 path outside design, docs, and wip changed from 0dd000000000 to the head c0ffee000000, such as flai/cmd/verify.go';
		const older = {
			commit: OLDER,
			head: HEAD,
			at_head: false,
			outside: ['flai/cmd/verify.go'],
			why
		};
		const v = klusterView(
			{
				...AT_HEAD,
				passed: false,
				deploy: { ...AT_HEAD.deploy!, commit: OLDER, head: older, passed_at_head: false },
				checks: { ...AT_HEAD.checks!, commit: OLDER, head: older, passed_at_head: false }
			},
			'/ops/deployments'
		)!;
		expect(v.passed).toBe(false);
		expect(v.headline).toBe('Its namespace is not deployed and checked at the head c0ffee000000');
		expect(v.deploy).toMatchObject({ commit: '0dd000000000', tone: 'passed', atHead: false });
		expect(v.deploy!.head).toBe(`not at the head: ${why}`);
		expect(v.checks!.head).toBe(`not at the head: ${why}`);
	});

	it('shows a failed deploy with the failed step’s log, and a failed or timed-out checks run', () => {
		const v = klusterView(
			{
				...AT_HEAD,
				passed: false,
				deploy: {
					...AT_HEAD.deploy!,
					outcome: 'failed',
					said: 'failed at the deploy step',
					failed_step: 'deploy',
					exit_code: 1,
					asked_by: undefined,
					passed_at_head: false
				},
				checks: undefined
			},
			'/base/ops/deployments'
		)!;
		expect(v.deploy).toMatchObject({
			outcome: 'failed at the deploy step',
			tone: 'failed',
			asked: 'who asked is not recorded',
			logHref: '/base/ops/deployments#logs',
			logLabel: "the deploy step's log"
		});
		expect(v.checks).toBeNull();
		expect(v.noChecks).toBe('No checks run against it has ended since its deploy.');
		const timedOut = klusterView(
			{
				...AT_HEAD,
				passed: false,
				checks: { ...AT_HEAD.checks!, outcome: 'timed-out', passed_at_head: false }
			},
			'/ops/deployments'
		)!;
		expect(timedOut.checks).toMatchObject({ outcome: 'timed out', tone: 'failed' });
		const failed = klusterView(
			{
				...AT_HEAD,
				passed: false,
				checks: { ...AT_HEAD.checks!, outcome: 'failed', log: undefined, passed_at_head: false }
			},
			'/ops/deployments'
		)!;
		expect(failed.checks).toMatchObject({ outcome: 'failed', tone: 'failed', logHref: '' });
	});

	it('says when nothing is deployed, the checks action is off, or the story has no branch', () => {
		const none = klusterView(
			{ ...AT_HEAD, passed: false, deploy: undefined, checks: undefined },
			'/ops/deployments'
		)!;
		expect(none.deploy).toBeNull();
		expect(none.noDeploy).toBe('No deploy of the namespace is recorded.');
		const off = klusterView(
			{ ...AT_HEAD, checks_on: false, checks: undefined },
			'/ops/deployments'
		)!;
		expect(off.headline).toBe('Its namespace is deployed at the head c0ffee000000');
		expect(off.noChecks).toBe('The checks host action is off, so no checks run is needed.');
		const branchless = klusterView(
			{ ...AT_HEAD, head: undefined, passed: false, deploy: undefined, checks: undefined },
			'/ops/deployments'
		)!;
		expect(branchless.headline).toBe(
			'S-0662 has no branch story/S-0662, so nothing is at its head'
		);
	});

	it('is absent when no evidence applies or the report has none', () => {
		expect(klusterView(undefined, '/ops/deployments')).toBeNull();
		expect(klusterView(null, '/ops/deployments')).toBeNull();
		expect(
			klusterView(
				{
					applies: false,
					why: 'the kluster host action is off for this project',
					story: 'S-0662',
					checks_on: false,
					passed: false
				},
				'/ops/deployments'
			)
		).toBeNull();
	});
});
