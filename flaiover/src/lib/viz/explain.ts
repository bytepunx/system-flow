// The explainer of each chart (S-0485): what the chart shows, how each figure it draws is
// measured, and how to read it, written from design/system/metrics.md, which defines every figure,
// with the anchor of the heading there that does. Plain text: a span in backticks names a field,
// a command, or a file, and the component draws it as code.
import type { Kind } from './catalog';
import { plural } from './charts';

/** What a chart shows, how its figures are measured, how to read it, and where metrics.md defines them. */
export type Explainer = { what: string; how: string; read: string; metrics: string };
/** What an explainer may name: the report's item type (story, task, epic) and the bucket the chart is laid out in (hour, day, week). */
export type ExplainContext = { type: string; bucket: string };

/** The words an explainer is written with: the type, its plural and article, the bucket, and its plural. */
type Words = {
	type: string;
	types: string;
	Types: string;
	a: string;
	A: string;
	bucket: string;
	buckets: string;
};

const capital = (s: string) => s.charAt(0).toUpperCase() + s.slice(1);

function words({ type, bucket }: ExplainContext): Words {
	const a = /^[aeiou]/.test(type) ? `an ${type}` : `a ${type}`;
	return {
		type,
		types: plural(type),
		Types: capital(plural(type)),
		a,
		A: capital(a),
		bucket,
		buckets: plural(bucket)
	};
}

// The anchors of the headings of metrics.md, as the Docs page ids them (markdown.ts `slug`).
const PER_ITEM = 'per-item-durations';
const AGGREGATES = 'aggregates';
const PRECISION = 'precision-rules';
const USAGE = 'usage-(s-0143%2C-s-0163)';
const SPEND = 'spend-over-time-(s-0163)';
const FORECASTS = 'forecasts-and-estimates';
const COST_OF_DELAY = 'cost-of-delay';
const PULL_ORDER = 'what-the-pull-order-costs-(s-0213)';
const WAITING = 'waiting';
const CLAIMS = 'claims-and-touches';
const STRATEGIC_DAYS = 'strategic-use-per-day';
const COORDINATION = 'coordination';
const QUEUE = 'merge-queue-(s-0404)';

/** The note every chart of spend over time carried under it. */
const spendNote = (w: Words) =>
	`Each ${w.bucket} holds the items that entered \`done\` in it, in UTC, with everything their agents spent on them.`;
/** The note under Tokens / Day and $ / Day. */
const meanNote = (w: Words) =>
	`The dashed line is the mean per ${w.bucket} so far, from the first with spend.`;
const TOKENS =
	"An item's tokens are the sum over its models of `input + output + cache_read + cache_write`.";
const COST = "An item's cost is the sum over its models of `cost`, in US dollars.";
/** What the star beside an agents' cost marks. */
const STAR =
	"A star marks a cost estimated in part: a task's share of its story's session, or a run that ended without its totals.";
const BOTH_MODELS =
	'An item two models worked on counts whole for each, with all its agent time, since flai does not split agent time by model.';
const ALL_MODELS = 'With two models or more, a dashed line is all of them together.';
const NEAREST_RANK =
	'nearest-rank percentiles: the p-th is the value at position `ceil(p/100 × n)` of the sorted values';
const CYCLE =
	'from its first move to `in-progress`, `started`, to its last transition, the move to `done` or `cancelled`, `completed`';

const EXPLAINERS: Record<Kind, (w: Words) => Explainer> = {
	'cycle-time': (w) => ({
		what: `One point per ${w.type} completed in the window, at the moment it was completed and as high as its cycle time in days, coloured by its nature. Two dashed lines mark the p50 and the p85 of the cycle time.`,
		how: `${w.A}'s cycle time runs ${CYCLE}: \`items[].cycle_time_seconds\` in \`flai stats --json\`. ${w.A} never started has none. The lines are \`summary.cycle_time.p50_seconds\` and \`p85_seconds\`, over the ${w.types} completed in the window, cancelled ones left out, ${NEAREST_RANK}.`,
		read: `Half the ${w.types} took no longer than the p50 line, and 85% no longer than the p85 line, the figure forecasts use. Points far above the p85 line are the ones to look into.`,
		metrics: PER_ITEM
	}),
	'burn-up': (w) => ({
		what: `Two lines over each day of the window: scope, the ${w.types} created so far, and done, those done so far. Choose an epic to see its own.`,
		how: `From \`burnup\` in \`flai stats --json\`, per UTC day. ${w.A} counts in scope from the day of its \`created\`; cancelled ${w.types} are left out. The lines run from the day that holds the window's start, or the first ${w.type}'s \`created\` day if that is later, to today.`,
		read: `The gap between the lines is the work left. Done rising as fast as scope keeps the gap steady; scope rising faster means work arrives faster than it is finished.`,
		metrics: PRECISION
	}),
	cfd: (w) => ({
		what: `Stacked bands, one per state, backlog, ready, in progress, review, and done: how many ${w.types} were in each at the end of each day of the window.`,
		how: `From \`cfd\` in \`flai stats --json\`: each ${w.type}'s state at the end of the day, 23:59:59 UTC, reconstructed from its transitions. It runs from the day that holds the window's start, or the first ${w.type}'s \`created\` day if that is later, to today.`,
		read: `Widening bands show where work piles up. A band that keeps its width means work leaves that state as fast as it enters.`,
		metrics: PRECISION
	}),
	'time-in-state': (w) => ({
		what: `One stacked bar per day of the window on which ${w.types} were completed: the mean hours those ${w.types} spent in each state. Hover a bar for the ${w.types} it averages. Under it, one bar of each state's share of all their time.`,
		how: `${w.A}'s time in a state is the sum of its intervals in that state, each from a transition's \`at\` to the next one's: \`items[].time_in_state_seconds\` in \`flai stats --json\`. A bar averages it over the ${w.types} completed that UTC day. The share is \`summary.time_in_state_share\`: the total seconds per state over the total lead time, from \`created\` to \`completed\`, of the ${w.types} completed in the window.`,
		read: `The largest part is where the process spends its time, the part to optimise first. Time in ready or review is waiting rather than work.`,
		metrics: PER_ITEM
	}),
	throughput: (w) => ({
		what: `One bar per ISO week of the window: the ${w.types} completed in it, stacked by nature. A week with none is drawn at zero.`,
		how: `From \`throughput\` in \`flai stats --json\`: the ${w.types} whose \`completed\`, the last transition, to \`done\` or \`cancelled\`, falls in the week. Weeks are ISO weeks starting Monday, UTC, from the one that holds the window's start to the one that holds now. \`summary.throughput_per_week\` is the ${w.types} completed in the window over its days divided by seven.`,
		read: `Bars of even height are a steady pace to forecast from. Bars that swing widely make a forecast from throughput less certain.`,
		metrics: AGGREGATES
	}),
	'agent-waiting': (w) => ({
		what: `One stacked bar per ISO week of the window: the hours the agents of the ${w.types} completed that week waited on their threads and in review. Under pull-request acceptance, review is split into its time on the required checks and on the reviewers. The dashed line is the mean wait per ${w.type} done that week. The table lists the ten longest waits.`,
		how: `A thread's wait runs from its first entry to the first later entry by another author that is not a recommendation; only its part in the ${w.type}'s \`in-progress\` intervals counts, \`items[].wait_threads_seconds\`. A review wait is the time in \`review\`, \`items[].wait_review_seconds\`. A checks wait runs from each push to the read at which no required check is pending, and a reviewers wait from then to the first approval or requested change. A bar is \`waiting.weeks[]\`, summed over the ${w.types} completed in the week, so a wait counts in the week its ${w.type} was done.`,
		read: `Tall thread bars mean agents wait on answers; tall review bars mean finished work waits to be accepted. A week with nothing done has no point on the dashed line.`,
		metrics: WAITING
	}),
	queue: () => ({
		what: `The lines are the mean and the longest time the stories the merge queue landed or sent back each week waited in it, in hours. On the right are the mean number of stories in the batches it built, and as bars the red batches it bisected.`,
		how: `From \`queue.weeks[]\` in \`flai stats --json\`, read from the queue's event log, \`.flai-cache/merge-queue/events.jsonl\`. A story's wait runs from its enqueueing, its latest \`enqueued\` event, to its \`landed\` or \`sent_back\` event: \`wait_mean_seconds\` and \`wait_max_seconds\` are the mean and the longest of the week's. \`batch_size_mean\` is the mean number of stories on the week's \`batch_built\` events, and \`bisections\` counts its \`bisected\` events, each a failed batch's bisection. Each event counts in the ISO week that holds it.`,
		read: `A week with none landed or sent back has no wait, and one with no batch no size: a gap, not zero. Weeks before the release that writes the log show nothing.`,
		metrics: QUEUE
	}),
	'token-rate': (w) => ({
		what: `One line per model: the tokens per minute of agent work over the ${w.types} done in each ${w.bucket} that the model worked on. ${ALL_MODELS}`,
		how: `${spendNote(w)} ${TOKENS} The rate is tokens over agent minutes, the usage \`seconds\` over 60: \`tokens_per_minute\` of each model in \`usage.spend.${w.type}.buckets[]\` of \`flai stats --json\`. A model's rate is its own tokens over the agent minutes of the items it worked on. ${BOTH_MODELS}`,
		read: `A higher line is a model that moves more tokens per minute of agent work. A ${w.bucket} without agent time has no point.`,
		metrics: SPEND
	}),
	'tokens-spent': (w) => ({
		what: `One bar per ${w.bucket}: the tokens spent on the ${w.types} done in it, stacked by model, with the mean per ${w.bucket} so far as a dashed line.`,
		how: `${spendNote(w)} ${meanNote(w)} ${TOKENS} A bar is the \`tokens\` of each model in \`usage.spend.${w.type}.buckets[]\` of \`flai stats --json\`. The dashed line is \`mean_tokens\`: the sum from the first ${w.bucket} of the series to this one over the number of those ${w.buckets}, empty ones counted.`,
		read: `Bars above the dashed line are ${w.buckets} that took more tokens than usual. Spend counts when its item is done, not when the tokens were used, so a ${w.bucket} with nothing done is empty.`,
		metrics: SPEND
	}),
	'tokens-per-item': (w) => ({
		what: `One line per item type, epics, stories, and tasks: the tokens an item took on average, over the items of that type done in each ${w.bucket}.`,
		how: `${spendNote(w)} ${TOKENS} Tokens per item is tokens over items: \`tokens_per_item\` in \`usage.spend.epic.buckets[]\`, \`usage.spend.story.buckets[]\`, and \`usage.spend.task.buckets[]\` of \`flai stats --json\`. A story's usage holds its tasks', so the types are never added together.`,
		read: `A rising line means items of that type take more tokens each. Click a name in the legend to hide a line, such as epics, which dwarf the others.`,
		metrics: SPEND
	}),
	'tokens-per-dollar': (w) => ({
		what: `One line per model: how many tokens a dollar bought over the ${w.types} done in each ${w.bucket} that the model worked on. ${ALL_MODELS}`,
		how: `${spendNote(w)} ${TOKENS} ${COST} The figure is tokens over cost: \`tokens_per_dollar\` of each model in \`usage.spend.${w.type}.buckets[]\` of \`flai stats --json\`.`,
		read: `A higher line is a model whose tokens cost less. A ${w.bucket} without cost has no point.`,
		metrics: SPEND
	}),
	'cost-spent': (w) => ({
		what: `One bar per ${w.bucket}: the dollars spent on the ${w.types} done in it, stacked by model, with what the strategic agents spent on them stacked on top as a series of its own, strategic. The dashed line is the mean per ${w.bucket} so far.`,
		how: `${spendNote(w)} ${meanNote(w)} ${COST} A bar is the \`cost\` of each model in \`usage.spend.${w.type}.buckets[]\` of \`flai stats --json\`, and strategic is \`buckets[].strategic.cost\`, what the planner, the orchestrator, and the analyzer spent on those ${w.types}, never added to the agents' figures. The dashed line is \`mean_cost\`, the agents' alone: the sum from the first ${w.bucket} of the series to this one over the number of those ${w.buckets}, empty ones counted. ${STAR} Strategic cost is always estimated.`,
		read: `Bars above the dashed line are ${w.buckets} that cost more than usual. Spend counts when its item is done, not when it was spent, so a ${w.bucket} with nothing done is empty.`,
		metrics: SPEND
	}),
	'cost-per-item': (w) => ({
		what: `One line per item type, epics, stories, and tasks: what an item cost on average, over the items of that type done in each ${w.bucket}.`,
		how: `${spendNote(w)} ${COST} Cost per item is cost over items: \`cost_per_item\` in \`usage.spend.epic.buckets[]\`, \`usage.spend.story.buckets[]\`, and \`usage.spend.task.buckets[]\` of \`flai stats --json\`. A story's usage holds its tasks', so the types are never added together. ${STAR}`,
		read: `A rising line means items of that type cost more each. Click a name in the legend to hide a line, such as epics, which dwarf the others.`,
		metrics: SPEND
	}),
	cost: (w) => ({
		what: `One bar per ${w.type} completed in the window with usage, in order of completion, stacked by model. What the strategic agents spent on ${w.a} is stacked on top as a series of its own, strategic; ${w.a} only they spent on shows that alone.`,
		how: `${COST} It comes from the ${w.type}'s \`usage\` front matter: \`items[].usage.models[].cost\` in \`flai stats --json\`. Strategic is its \`items[].usage.strategic\` costs summed, what the planner, the orchestrator, and the analyzer spent on it, apart from its total, which is the agents'. A star after an ID marks a cost estimated in part: a task's share of its story's session, a run that ended without its totals, or strategic cost, which is always estimated.`,
		read: `The tallest bars are the ${w.types} that cost most. Hover a bar for its cost per model and what the strategic agents spent on it.`,
		metrics: USAGE
	}),
	'time-per-model': (w) => ({
		what: `One line per model: the minutes of agent work ${w.a} took on average, over the ${w.types} done in each ${w.bucket} that the model worked on. ${ALL_MODELS}`,
		how: `${spendNote(w)} Minutes per item is the items' agent \`seconds\` over their number, in minutes: \`minutes_per_item\` of each model in \`usage.spend.${w.type}.buckets[]\` of \`flai stats --json\`. ${BOTH_MODELS}`,
		read: `A lower line is a model that finishes ${w.a} in less agent time. A ${w.bucket} in which none of its ${w.types} took agent time has no point.`,
		metrics: SPEND
	}),
	'cost-per-model': (w) => ({
		what: `One line per model: what ${w.a} cost on average, over the ${w.types} done in each ${w.bucket} that the model worked on. ${ALL_MODELS}`,
		how: `${spendNote(w)} ${COST} A model's cost per item is its own cost over the items it worked on: \`cost_per_item\` of each model in \`usage.spend.${w.type}.buckets[]\` of \`flai stats --json\`. ${STAR}`,
		read: `A lower line is a model that delivers ${w.a} for less. Read it beside Avg. Time / Model.`,
		metrics: SPEND
	}),
	'forecast-accuracy': (w) => ({
		what: `One point per ${w.type} done in the window with a forecast, at its completion: how much longer, above zero, or shorter, below, it took than its forecast duration. Triangles show the same against the estimate, where it has one. Dashed lines sit at plus and minus the p50 and the p85 of the forecast error.`,
		how: `The forecast error is cycle time, ${CYCLE}, minus \`forecast.duration\`: \`items[].forecast_error_seconds\` in \`flai stats --json\`. The estimate error is cycle time minus \`estimate\`, \`estimate_error_seconds\`. The lines are \`forecasts.forecast.p50_seconds\` and \`p85_seconds\`, ${NEAREST_RANK}, of the absolute errors. Under a nature or a model they are its \`by_nature\` or \`by_model\` entry; under both, the percentiles of the points shown. Cancelled ${w.types} are left out.`,
		read: `Half the errors fall between the two p50 lines and 85% between the two p85 lines: the closer together the lines, the better the forecasts. Points mostly above zero mean forecasts are too short.`,
		metrics: FORECASTS
	}),
	'delivery-accuracy': (w) => ({
		what: `One point per ${w.type} done in the window with a forecast delivery date, at its completion: how many days after that date it was done, early below zero, with dashed lines at plus and minus the p50 and the p85. A line on the right-hand axis shows each week's share delivered on time.`,
		how: `The delivery error is \`completed\` minus \`forecast.delivery\`: \`items[].delivery_error_seconds\` in \`flai stats --json\`, drawn in days. The lines are \`forecasts.delivery.p50_seconds\` and \`p85_seconds\`, ${NEAREST_RANK}, of the absolute errors, narrowed by nature and model as on Forecast Accuracy. A week's share is its ${w.types} with an error of 0 or less over those with a delivery forecast; weeks are ISO weeks starting Monday, UTC. Cancelled ${w.types} are left out.`,
		read: `A week with no such ${w.type} has no share: the line breaks there rather than drop to 0. The closer the p50 and p85 lines sit to zero, the better the dates.`,
		metrics: FORECASTS
	}),
	'forecast-by-model': (w) => ({
		what: `One line per model: the p50 of the absolute forecast errors of its ${w.types} done in each ${w.bucket}.`,
		how: `${w.A}'s model is its \`agent.model\`, \`(none)\` without one: \`items[].model\` in \`flai stats --json\`. Its forecast error is cycle time, ${CYCLE}, minus \`forecast.duration\`, \`items[].forecast_error_seconds\`. A point is the nearest-rank p50 of the absolute errors of the model's ${w.types} completed in the ${w.bucket}, in UTC, cancelled ones left out.`,
		read: `A ${w.bucket} with no such ${w.type} is a gap, not 0. A lower line is a model whose ${w.types} are forecast better.`,
		metrics: FORECASTS
	}),
	parallelism: (w) => ({
		what: `For each day of the window: the ${w.types} in progress at its end, and the stories in ready held then, beside a dashed line at the board's in-progress limit.`,
		how: `From \`claims.days[]\` in \`flai stats --json\`: \`in_progress\` counts the ${w.types} in \`in-progress\` at 23:59:59 UTC, and \`held\` the stories in \`ready\` the hold rules held then, \`overlap\`, \`no-touches\`, or \`after\`, replayed from the states of the time and today's touches and \`after\`, since flai records neither a hold nor older touches. The limit is \`claims.limit\`, the board's \`in-progress\` limit today; without one there is no line.`,
		read: `Held stories on a day when in progress is below the limit could have started but for their claims. Held stories with in progress at the limit would have waited anyway.`,
		metrics: CLAIMS
	}),
	'hold-time': () => ({
		what: `One bar per ISO week of the window: the hours stories in ready were held, stacked by reason. Overlap is a claim that overlaps a story in progress; after, a story it waits for is not done; empty claim, its touches or the in-progress story's are empty.`,
		how: `From \`claims.weeks[].held_seconds\` in \`flai stats --json\`, per reason, replayed at every creation and transition from the states of the time and today's touches and \`after\`. A hold counts under the one reason it is named by: \`after\` when a story in its \`after\` is not done, otherwise the reason of the first story in progress, by ID, whose claim holds it. The window's first week is whole, and this week ends at now.`,
		read: `A bar's height is the time held, since the reasons add up to it. A tall overlap or empty-claim bar points at touches drawn too wide or left out.`,
		metrics: CLAIMS
	}),
	coordination: () => ({
		what: `For each ISO week of the window, bars of the conflicts found early against those met late, and a line each for the conversations between stories opened, agreed, escalated, and closed with their story. Beside it, the stories started on a share and the hold time they saved.`,
		how: `The bars are the conflicts a trial merge found while both stories were open, early, \`coordination.weeks[].conflicts.found\`, against those a rebase met at sync or at acceptance, late, \`met_sync\` plus \`met_accept\`, from the conflict log \`.flai-cache/conflicts.jsonl\`. The lines are the conversations between stories under \`wip/messages\` opened each week, and those that agreed, by a share or a clean merge at sync, were escalated, or closed with their story. The hold time saved, \`shares.held_seconds\` summed, is what the stories started on a share would have waited for their holder: from their first move to \`in-progress\` to the holder's first move to \`review\`, \`done\`, or \`cancelled\`.`,
		read: `More conflicts found early than met late, and more conversations agreed than escalated, mean stories find their conflicts while both are open and settle them between themselves. Agreed is a floor: agreement in a reply, or a narrowed claim, is not recorded. Weeks before the release that writes the conflict log show no conflicts.`,
		metrics: COORDINATION
	}),
	'touches-drift': () => ({
		what: `One bar per story done in the window, at its completion: the files its commits changed outside its touches, stacked on its touches no commit changed. A line on the right-hand axis shows each week's share of stories whose touches were exact.`,
		how: `From \`claims.drift[]\` in \`flai stats --json\`. A story's commits are those on the main branch and the \`story/\` branches whose subject names it, \`[S-nnnn]\`, merges and files under \`wip\` left out; its touches include those of its tasks not cancelled. \`outside_count\` counts the files changed under none of its touches, and \`unchanged_count\` the touches no committed file is under. The share is \`claims.weeks[].exact_share\`: the stories completed in the week with both counts 0 over those completed. Cancelled stories are left out.`,
		read: `A file outside the touches is one the planner or the story's agent should have declared; a touch never changed is one they should have dropped. A week with no story done has no share, a gap, not 0.`,
		metrics: CLAIMS
	}),
	'strategic-cost': (w) => ({
		what: `The bars are what planning, orchestration, and analysis cost each day, the planner, the orchestrator, and the analyzer stacked. The line is what the agents spent per ${w.type} completed that day, the strategic agents left out; the dashed line is what the strategic agents add to each ${w.type} over the window.`,
		how: `From \`strategic_days[]\` in \`flai stats --json\`, one point per UTC day. A bar is the \`cost\` of an agent under \`agents\`, the cost of its activities that ended that day, from its activity document under \`wip/agents\`. The line is \`cost_per_item\`: the agents' usage cost of the ${w.types} completed that day on which they spent, over their number, strategic usage left out. The dashed line is the sum of \`cost\` over the window over the sum of \`completed\`, labelled with its share of what the agents spent per ${w.type}, \`usage.cost\` over \`usage.items\`. A star marks a cost estimated in part: an activity of that agent that day carries an estimated cost.`,
		read: `At 10%, the strategic agents add a tenth to each ${w.type}'s cost. A day with nothing completed has no point on the line, and with nothing completed in the window, or nothing spent by the agents, there is no dashed line.`,
		metrics: STRATEGIC_DAYS
	}),
	'strategic-use': (w) => ({
		what: `The bars are the hours the planner, the orchestrator, and the analyzer worked each day, stacked. The lines are the mean cycle time and the mean waiting of the ${w.types} completed that day, in hours.`,
		how: `From \`strategic_days[]\` in \`flai stats --json\`, one point per UTC day. A bar is the \`seconds\` of an agent under \`agents\`, the wall-clock seconds of its activities that ended that day. The mean cycle time is \`cycle_time_seconds\`, over the ${w.types} completed that day with one, each ${CYCLE}. The mean waiting is each such ${w.type}'s \`wait_threads_seconds\`, its thread waits while in progress, plus its \`wait_review_seconds\`, its time in review, averaged.`,
		read: `Waiting or cycle time falling while the strategic agents' hours rise is the return on their time. A day with nothing completed has no point on the lines.`,
		metrics: STRATEGIC_DAYS
	}),
	'cod-outstanding': (w) => ({
		what: `Stacked areas over each day of the window: the cost of delay per week of the ${w.types} in each column at the end of the day, backlog, ready, in progress, and review, in the colours of Cumulative Flow.`,
		how: `\`cost_of_delay.days[].outstanding\` in \`flai stats --json\`: per column, the sum of the \`cost_of_delay.value\` of the ${w.types} in it at 23:59:59 UTC, archived ones included. A value is what a week of waiting costs, in the project's currency, and today's value stands for all of ${w.a}'s history. ${w.Types} without one add nothing; \`days[].without_value\` counts them per column.`,
		read: `A tall backlog or ready area is value waiting to be pulled. A column full of ${w.types} without a value looks cheaper than it is.`,
		metrics: COST_OF_DELAY
	}),
	'cod-incurred': (w) => ({
		what: `One bar per ISO week of the window: what waiting cost in it, with the mean per week so far as a dashed line.`,
		how: `\`cost_of_delay.weeks[].incurred\` in \`flai stats --json\`: over every ${w.type}, its \`cost_of_delay.value\` times the seconds of the week it spent in \`backlog\` or \`ready\`, over 604800, a week. The window's first week is whole. The dashed line is the sum from the window's first week to this one over the number of those weeks.`,
		read: `Bars above the dashed line are weeks in which waiting cost more than usual. ${w.Types} without a value add nothing, so the bars leave them out.`,
		metrics: COST_OF_DELAY
	}),
	'cod-order': () => ({
		what: `One step line per order, the board's pull order, by cost of delay, and by WSJF: what the stories in ready now will cost by the time each is pulled, rising at each story's projected pull. The time axis runs from now to the last pull, not over the window.`,
		how: `\`cost_of_delay.order.series[].points\` in \`flai stats --json\`. Each order is played out on the board's \`in-progress\` limit as \`flai forecast\` plays out the pull order, and a story's cost is its \`cost_of_delay.value\` times the seconds from now to its pull, over 604800. By cost of delay is \`flai order --by cod\`, value per week, highest first; by WSJF is \`flai order --by wsjf\`, value over \`forecast.duration\` in hours, highest first. Only ready stories with both a value and a \`forecast.duration\` are placed; the rest are in \`left_out\`.`,
		read: `The lowest line at the right is the cheapest order. A story's cost ends when it is pulled, not when it is delivered. With no in-progress limit every story is pulled at once and every order costs nothing.`,
		metrics: PULL_ORDER
	})
};

/** The explainer of a chart kind. */
export function explain(kind: Kind, ctx: ExplainContext): Explainer {
	return EXPLAINERS[kind](words(ctx));
}
