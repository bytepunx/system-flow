// A story's task plan (S-0176): which of its tasks can start, which wait for which, and the layers
// of tasks that can run at once, as flai works it out from each task's `after`. flai's item.get
// gives the plan beside the item and its children; board.get gives a story's card the counts.

export type TaskState = 'ready' | 'waiting' | 'in-progress' | 'done' | 'cancelled';

/** One task of the plan: `after` and `waiting_for` are absent when empty. */
export type PlanTask = {
	id: string;
	state: TaskState;
	/** the tasks of the same story it starts after */
	after?: string[];
	/** those of `after` not yet done */
	waiting_for?: string[];
};

/** A story's plan: its tasks, and the layers of tasks that can run at once, first to last. */
export type TaskPlan = { tasks: PlanTask[]; layers: string[][] };

/** A story's tasks counted for its card on the board. */
export type TaskSummary = {
	ready: number;
	waiting: number;
	in_progress: number;
	done: number;
	layers: number;
};

/** The plan as flai marshals it, with nulls where a list was empty. */
export type FlaiPlan = {
	tasks: (Omit<PlanTask, 'after' | 'waiting_for'> & {
		after?: string[] | null;
		waiting_for?: string[] | null;
	})[];
	layers: (string[] | null)[] | null;
} | null;

/** The plan in the shape the page reads, or undefined when flai sent none. */
export function planFrom(p: FlaiPlan | undefined): TaskPlan | undefined {
	if (!p) return undefined;
	return {
		tasks: (p.tasks ?? []).map((t) => ({
			id: t.id,
			state: t.state,
			...(t.after?.length ? { after: t.after } : {}),
			...(t.waiting_for?.length ? { waiting_for: t.waiting_for } : {})
		})),
		layers: (p.layers ?? []).map((l) => l ?? [])
	};
}

const STATE_LABEL: Record<TaskState, string> = {
	ready: 'ready to start',
	waiting: 'waiting',
	'in-progress': 'in progress',
	done: 'done',
	cancelled: 'cancelled'
};

/** What a task's state says on the page; a state this dashboard does not know is shown as it is. */
export const stateLabel = (s: string) => STATE_LABEL[s as TaskState] ?? s;

/** What a waiting task waits for: the tasks not yet done, or its whole `after` if flai sent no subset. */
export const waitsFor = (t: PlanTask) => t.waiting_for ?? t.after ?? [];

/**
 * Open tasks in no layer: flai leaves a task whose `after` runs in a cycle out of every layer. A
 * done or cancelled task is not counted, whether or not flai lays it out.
 */
export function unplaced(plan: TaskPlan): string[] {
	const placed = new Set(plan.layers.flat());
	return plan.tasks
		.filter((t) => t.state !== 'done' && t.state !== 'cancelled' && !placed.has(t.id))
		.map((t) => t.id);
}

const plural = (n: number, one: string, many = `${one}s`) => `${n} ${n === 1 ? one : many}`;

/**
 * The card's line for a story's tasks: the open ones by state, those at zero left out, and the
 * layers; all done says so.
 */
export function tasksLine(s: TaskSummary): string {
	const open = [
		[s.in_progress, 'in progress'],
		[s.ready, 'ready'],
		[s.waiting, 'waiting']
	] as const;
	const parts = open.filter(([n]) => n > 0).map(([n, what]) => `${n} ${what}`);
	if (!parts.length) parts.push(`${s.done} done`);
	return `tasks ${parts.join(' · ')} · ${plural(s.layers, 'layer')}`;
}

/** The same line in full, for the card's tooltip and accessible name. */
export function tasksTitle(s: TaskSummary): string {
	return (
		`tasks: ${s.in_progress} in progress, ${s.ready} ready to start, ${s.waiting} waiting, ` +
		`${s.done} done; the plan has ${plural(s.layers, 'layer')}`
	);
}
