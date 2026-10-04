// What a story's agent is doing (S-0104), as flai on the host reports it in agent.status: working
// (green), waiting for the designer or, queued by a retry, for room (yellow), failed (red), or
// worked, which shows no dot. A story in ready that another's claim holds is waiting too, with the
// hold (S-0129), and so is a story in progress that this host has had no agent for, with where it was
// begun (S-0177).

export type AgentRun = {
	story: string;
	harness?: string;
	model?: string;
	command: string;
	agent: string;
	pid?: number;
	started: string;
	ended?: string;
	exit?: number;
	error?: string;
	log?: string;
	outcome?: 'worked' | 'failed' | 'asked' | 'stopped';
	why?: string;
	/** When the operator queued another agent for its story, waiting for room (S-0118). */
	queued?: string;
	/** When the operator stopped it (S-0170). */
	stopped?: string;
	/** The question it ended waiting on, when it did. */
	thread?: string;
};

export type ActivityState = 'working' | 'waiting' | 'failed' | 'worked';

/**
 * One thing a story's agent said or did, as flai reads it from the agent's log (S-0142): `session`
 * (it started), `text`, `thinking`, `tool` (a call, with the tool's name), `result` (what a tool
 * answered), `task` (a background task), `end` (how the session ended), or `output` (a line that is
 * not the harness's stream). `error` marks a result or an end that failed.
 */
export type AgentStreamEntry = {
	kind: 'session' | 'text' | 'thinking' | 'tool' | 'result' | 'task' | 'end' | 'output';
	tool?: string;
	text: string;
	error?: boolean;
};

/**
 * One read of a story's agent's stream (flai's agent.stream): the entries from byte `from` to
 * `next`, where the next read starts; `more` when the log goes on past it. `started` names the run:
 * another run writes another log.
 */
export type AgentStreamRead = {
	story: string;
	agent: string;
	started: string;
	ended?: string;
	running: boolean;
	outcome?: string;
	from: number;
	next: number;
	size: number;
	more?: boolean;
	skipped?: number;
	entries: AgentStreamEntry[];
};

/**
 * Why a story in ready waits for another's claim (S-0128, ADR-0046): `overlap` or `no-touches`, or
 * for a story it names in after: (`after`, S-0130), and flai's reason, which names every story that
 * holds it and says what clears it.
 */
export type Hold = { code: string; reason: string };

/**
 * Where a story in progress that this host has had no agent for was begun (S-0177, ADR-0064): who
 * last moved it to in-progress and when, the agent and host its narrative names, and `here` when that
 * host is this one, outside flai serve.
 */
export type Elsewhere = { by: string; at: string; agent?: string; host?: string; here?: boolean };

export type StoryActivity = {
	state: ActivityState;
	why?: string;
	run: AgentRun;
	thread?: string;
	/** Set for a held story in ready, whether or not it has had an agent (S-0129). */
	hold?: Hold;
	/** Set for a story in progress begun with no agent of this host's (S-0177). */
	elsewhere?: Elsewhere;
};

export type HostAgent = {
	enabled: boolean;
	/** Whether the plan host action is on for the project (S-0263). */
	plan_enabled?: boolean;
	state?: {
		command: string;
		running?: AgentRun | null;
		last?: AgentRun | null;
		waiting?: string;
		stories?: Record<string, StoryActivity>;
		/** The newest planner run for each epic or story it planned, by the item's ID (S-0208). */
		plans?: Record<string, PlanRun>;
	};
};

/** A planner run (S-0208): an agent's run for the epic or story it plans, which is no story's. */
export type PlanRun = Omit<AgentRun, 'story'> & { item: string };

/** The dot's colour, or none for an agent that finished its story. */
export function dotClass(state: ActivityState): string | null {
	switch (state) {
		case 'working':
			return 'bg-dot-working';
		case 'waiting':
			return 'bg-dot-waiting';
		case 'failed':
			return 'bg-dot-failed';
		default:
			return null;
	}
}

/**
 * What each story's agent is doing, for the dots: every story while the `agent` action is on, and
 * otherwise only the held ones and those begun elsewhere, since flai reports both whether or not it
 * starts agents (S-0129, S-0177).
 */
export function storyActivity(h: HostAgent | null): Record<string, StoryActivity> {
	const stories = h?.state?.stories ?? {};
	if (h?.enabled) return stories;
	return Object.fromEntries(Object.entries(stories).filter(([, a]) => a.hold || a.elsewhere));
}

/** The card's short line for a story begun elsewhere: by whom and where, not when. */
export function elsewhereLine(e: Elsewhere): string {
	const who = e.agent || e.by || 'someone';
	const where = e.here
		? 'on this host, outside flai serve'
		: e.host
			? `on ${e.host}`
			: 'on another host';
	return `begun by ${who} ${where}; no agent here`;
}

const storyID = /\bS-\d+\b/g;

/**
 * The stories a hold waits for: those its reason names in each `starts when` clause, up to the next
 * semicolon, in order, once each. A hold that is `after` and `overlap` too has two such clauses, and
 * a note after one may name the story itself, which it does not wait for (S-0130).
 */
export function holdWaitsFor(h: Hold): string[] {
	const clauses = [...h.reason.matchAll(/starts when ([^;]*)/g)].map((m) => m[1]);
	return [...new Set(clauses.flatMap((c) => c.match(storyID) ?? []))];
}

/** The card's short line for a hold: its code and the stories it waits for. */
export function holdLine(h: Hold): string {
	const ids = holdWaitsFor(h);
	return `held (${h.code})${ids.length ? `: ${ids.join(', ')}` : ''}`;
}

/** A hold's reason cut into text and the story IDs in it, so each ID can be a link. */
export function reasonParts(reason: string): { text: string; id?: string }[] {
	return reason
		.split(/\b(S-\d+)\b/)
		.map((text, i) => (i % 2 ? { text, id: text } : { text }))
		.filter((p) => p.text);
}

/** One line for a tooltip or a screen reader: what the agent is doing and why. */
export function activityLine(a: StoryActivity): string {
	if (a.hold) return a.hold.reason;
	if (a.elsewhere) return a.why ?? elsewhereLine(a.elsewhere);
	const who = [a.run.harness, a.run.model].filter(Boolean).join(', ') || a.run.command;
	switch (a.state) {
		case 'working':
			return `agent working (${who})`;
		case 'waiting':
			return `agent waiting (${who})${a.why ? `: ${a.why}` : ''}`;
		case 'failed':
			if (a.run.outcome === 'stopped') return `agent stopped by the operator (${who})`;
			return `agent failed (${who})${a.why ? `: ${a.why}` : ''}`;
		default:
			return `agent finished (${who})`;
	}
}

/**
 * Whether the operator can stop a story's agent (S-0170): one flai started that runs, whether at
 * work or waiting for the designer, or that ended waiting for an answer and would be started again
 * by it. Not a held story's stand-in, nor a retry queued for room: neither has a process.
 */
export function stoppable(a: StoryActivity | undefined): a is StoryActivity {
	if (!a?.run.started || a.hold || a.run.stopped) return false;
	if (a.state === 'working') return true;
	return a.state === 'waiting' && (!a.run.ended || a.run.outcome === 'asked');
}

// What a writer may have flai do for a story's agent, with the agent action on (S-0115, S-0116,
// S-0177): one rule for the story page and the card's menu (S-0202).

/** Whether a writer may retry a story in ready or in progress whose agent failed (S-0116). */
export function retryable(
	a: StoryActivity | undefined,
	enabled: boolean,
	status: string,
	writable: boolean
): a is StoryActivity {
	return (
		writable &&
		enabled &&
		(a?.state === 'failed' || (!!a?.hold && a.run.outcome === 'failed')) &&
		(status === 'ready' || status === 'in-progress')
	);
}

/** Whether a writer may start the agent of a story in ready that has had none (S-0115, S-0129). */
export function startable(
	a: StoryActivity | undefined,
	enabled: boolean,
	status: string,
	writable: boolean
): boolean {
	return writable && enabled && status === 'ready' && (!a || (!!a.hold && !a.run.started));
}

/** Whether a writer may start an agent here for a story in progress begun elsewhere (S-0177). */
export function startableHere(
	a: StoryActivity | undefined,
	enabled: boolean,
	status: string,
	writable: boolean
): boolean {
	return writable && enabled && status === 'in-progress' && !!a?.elsewhere;
}

/** The one thing a writer may have flai do for a story's agent, Retry first, or null. */
export function agentAction(
	activity: StoryActivity | undefined,
	enabled: boolean,
	status: string,
	writable: boolean
): { action: 'start' | 'restart'; label: 'Start agent' | 'Retry' } | null {
	if (retryable(activity, enabled, status, writable)) return { action: 'restart', label: 'Retry' };
	if (startable(activity, enabled, status, writable))
		return { action: 'start', label: 'Start agent' };
	if (startableHere(activity, enabled, status, writable))
		return { action: 'restart', label: 'Start agent' };
	return null;
}

/**
 * Whether any agent is still running, so its state is worth asking for again. A story begun elsewhere
 * has none here, and changes only when its files do.
 */
export function anyRunning(h: HostAgent | null): boolean {
	return Object.values(h?.state?.stories ?? {}).some(
		(a) => (a.state === 'working' || a.state === 'waiting') && !a.elsewhere
	);
}
