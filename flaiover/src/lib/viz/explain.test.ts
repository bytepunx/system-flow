// S-0485: each chart's explainer and its anchor in design/system/metrics.md.
import { readFileSync } from 'node:fs';
import { describe, expect, it } from 'vitest';
import { slug } from '../markdown';
import { KINDS, type Kind } from './catalog';
import { explain, type ExplainContext } from './explain';

const STORY_DAY: ExplainContext = { type: 'story', bucket: 'day' };
const CONTEXTS: ExplainContext[] = [
	STORY_DAY,
	{ type: 'epic', bucket: 'week' },
	{ type: 'task', bucket: 'hour' }
];
const METRICS_URL = new URL('../../../../design/system/metrics.md', import.meta.url);
const METRICS = readFileSync(METRICS_URL, 'utf8');
/** The anchors of the headings of metrics.md, as the Docs page ids them; fenced lines left out. */
const anchors = new Set(
	METRICS.split('```')
		.filter((_, i) => i % 2 === 0)
		.flatMap((part) => part.split('\n'))
		.flatMap((line) => {
			const m = /^#{1,6}\s+(.+)$/.exec(line);
			return m ? [slug(m[1])] : [];
		})
);
const balanced = (text: string) => (text.match(/`/g) ?? []).length % 2 === 0;
/** The three parts of a kind's explainer as one text. */
const all = (kind: Kind, ctx: ExplainContext = STORY_DAY) => {
	const e = explain(kind, ctx);
	return `${e.what} ${e.how} ${e.read}`;
};

const SPEND_NOTE =
	'Each week holds the items that entered `done` in it, in UTC, with everything their agents spent on them.';
const AGENTS_STAR =
	/star .*estimated in part: a task's share of its story's session, .*a run that ended without its totals/;
const STRATEGIC_STAR =
	'A star marks a cost estimated in part: an activity of that agent that day carries an estimated cost.';
const STRATEGIC_LINES =
	'The line is what the agents spent per epic completed that day, the strategic agents left out; ' +
	'the dashed line is what the strategic agents add to each epic over the window.';
const STRATEGIC_RETURN =
	"Waiting or cycle time falling while the strategic agents' hours rise is the return on their time.";
const QUEUE_LINES =
	'The lines are the mean and the longest time the stories the merge queue landed or sent back ' +
	'each week waited in it';
const QUEUE_RIGHT =
	'On the right are the mean number of stories in the batches it built, and as bars the red ' +
	'batches it bisected.';
const DIGEST_BARS =
	'One stacked bar per ISO week of the window: the threads opened to the operator per story ' +
	'completed that week.';

describe('explain', () => {
	it.each(KINDS)('gives %s every part, in plain text with balanced code spans', (kind) => {
		for (const ctx of CONTEXTS) {
			const e = explain(kind, ctx);
			expect(e.what.trim()).not.toBe('');
			expect(e.how.trim()).not.toBe('');
			expect(e.read.trim()).not.toBe('');
			for (const part of [e.what, e.how, e.read]) {
				expect(balanced(part)).toBe(true);
				expect(part).not.toMatch(/<[a-z]/i);
				expect(part).not.toMatch(/\$\{|undefined/);
			}
		}
	});

	it('reads the headings of metrics.md from the repository', () => {
		expect(anchors).toContain('spend-over-time-(s-0163)');
		expect(anchors.size).toBeGreaterThan(20);
	});

	it.each(KINDS)('anchors %s at a heading of metrics.md', (kind) => {
		expect(anchors).toContain(explain(kind, STORY_DAY).metrics);
	});

	it('names cycle time by its transitions, and the p50 and p85 as nearest-rank', () => {
		const { how } = explain('cycle-time', STORY_DAY);
		expect(how).toContain('first move to `in-progress`');
		expect(how).toContain('`done`');
		expect(how).toContain('`summary.cycle_time.p50_seconds`');
		expect(how).toContain('nearest-rank');
	});

	it('carries the note of spend over time, with the bucket chosen', () => {
		const week: ExplainContext = { type: 'story', bucket: 'week' };
		const spend: Kind[] = [
			'token-rate',
			'tokens-spent',
			'tokens-per-item',
			'tokens-per-dollar',
			'cost-spent',
			'cost-per-item',
			'time-per-model',
			'cost-per-model'
		];
		for (const kind of spend) expect(all(kind, week)).toContain(SPEND_NOTE);
		expect(all('tokens-spent', week)).toContain(
			'The dashed line is the mean per week so far, from the first with spend.'
		);
		expect(all('cost-spent', { type: 'task', bucket: 'hour' })).toContain(
			'The dashed line is the mean per hour so far, from the first with spend.'
		);
		expect(all('cost-per-item')).not.toContain('dashed line is the mean per');
		expect(all('cost')).not.toContain('Each day holds');
	});

	it('says what a star marks on the cost charts', () => {
		const costs: Kind[] = ['cost', 'cost-spent', 'cost-per-item', 'cost-per-model'];
		for (const kind of costs) expect(all(kind)).toMatch(AGENTS_STAR);
		expect(all('strategic-cost')).toContain(STRATEGIC_STAR);
		expect(all('tokens-spent')).not.toContain('star');
	});

	it('explains the strategic charts with the type', () => {
		const cost = all('strategic-cost', { type: 'epic', bucket: 'day' });
		expect(cost).toContain('The bars are what planning, orchestration, and analysis cost each day');
		expect(cost).toContain(STRATEGIC_LINES);
		const use = all('strategic-use');
		expect(use).toContain(
			'The bars are the hours the planner, the orchestrator, and the analyzer worked each day'
		);
		expect(use).toContain('the mean cycle time and the mean waiting of the stories completed');
		expect(use).toContain(STRATEGIC_RETURN);
	});

	it('explains coordination and the merge queue', () => {
		const coordination = all('coordination');
		expect(coordination).toContain('a trial merge found while both stories were open, early');
		expect(coordination).toContain('those a rebase met at sync or at acceptance, late');
		expect(coordination).toContain('agreed');
		expect(coordination).toContain('escalated');
		expect(coordination).toContain('what the stories started on a share would have waited for');
		const queue = all('queue');
		expect(queue).toContain(QUEUE_LINES);
		expect(queue).toContain('from its enqueueing');
		expect(queue).toContain('A week with none landed or sent back has no wait');
		expect(queue).toContain(QUEUE_RIGHT);
	});

	it('explains the digest by its threads to the operator, its repeats, and its gaps (S-0434)', () => {
		const digest = all('digest');
		expect(digest).toContain(DIGEST_BARS);
		expect(digest).toContain("so a bar's height is the week's threads per story");
		expect(digest).toContain("neither the operator, the manifest's `owner`, nor flai itself");
		expect(digest).toContain('`per_story` is `threads_to_operator` over it');
		expect(digest).toContain('a digest entry of an earlier day lists in its `asked_again`');
		expect(digest).toContain('A week with no story completed has no rate: a gap, not zero.');
		expect(digest).toContain('`digest.entries`');
		// the curator's runs, and the weeks it has not reached, which have no repeat count (S-0504)
		expect(explain('digest', STORY_DAY).what).toContain(
			"On the right, the curator's runs each week."
		);
		expect(digest).toContain("`curator_runs` counts the curator's runs that ended in the week");
		expect(digest).toContain(
			'A week the curator has not reached, after its newest run or with none at all, has no repeat count'
		);
		expect(digest).toContain('its asked-again part is a gap, not zero');
		expect(explain('digest', STORY_DAY).metrics).toBe('digest-(s-0434)');
		// the digest counts stories whatever the type chosen
		expect(all('digest', { type: 'epic', bucket: 'week' })).toBe(digest);
	});

	it('names the type and the bucket from the context', () => {
		const epic: ExplainContext = { type: 'epic', bucket: 'day' };
		expect(explain('cycle-time', epic).what).toContain('One point per epic completed');
		expect(explain('cycle-time', epic).how).toContain('An epic never started has none.');
		expect(explain('cod-incurred', epic).read).toContain('Epics without a value add nothing');
		const task: ExplainContext = { type: 'task', bucket: 'day' };
		expect(explain('throughput', task).what).toContain('the tasks completed in it');
		const hour: ExplainContext = { type: 'story', bucket: 'hour' };
		expect(explain('tokens-spent', hour).what).toContain('One bar per hour: the tokens spent');
		expect(explain('tokens-spent', hour).how).toContain('`usage.spend.story.buckets[]`');
		const week: ExplainContext = { type: 'story', bucket: 'week' };
		expect(explain('forecast-by-model', week).what).toContain('its stories done in each week');
	});
});
