import { describe, expect, it } from 'vitest';
import {
	criteriaOf,
	patchLines,
	patchRuns,
	readNdjson,
	sectionOf,
	toggleCriterion,
	verifyView,
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
});
