// What a story's agent is doing (S-0104), as flai on the host reports it in agent.status: working
// (green), waiting for the designer or, queued by a retry, for room (yellow), failed (red), or
// worked, which shows no dot. A story in ready that another's claim holds is waiting too, with the
// hold (S-0129).

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
	outcome?: 'worked' | 'failed';
	why?: string;
	/** When the operator queued another agent for its story, waiting for room (S-0118). */
	queued?: string;
};

export type ActivityState = 'working' | 'waiting' | 'failed' | 'worked';

/**
 * Why a story in ready waits for another's claim (S-0128, ADR-0046): `overlap` or `no-touches`, or
 * for a story it names in after: (`after`, S-0130), and flai's reason, which names every story that
 * holds it and says what clears it.
 */
export type Hold = { code: string; reason: string };

export type StoryActivity = {
	state: ActivityState;
	why?: string;
	run: AgentRun;
	thread?: string;
	/** Set for a held story in ready, whether or not it has had an agent (S-0129). */
	hold?: Hold;
};

export type HostAgent = {
	enabled: boolean;
	state?: {
		command: string;
		running?: AgentRun | null;
		last?: AgentRun | null;
		waiting?: string;
		stories?: Record<string, StoryActivity>;
	};
};

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
 * otherwise only the held ones, since flai reports a hold whether or not it starts agents (S-0129).
 */
export function storyActivity(h: HostAgent | null): Record<string, StoryActivity> {
	const stories = h?.state?.stories ?? {};
	if (h?.enabled) return stories;
	return Object.fromEntries(Object.entries(stories).filter(([, a]) => a.hold));
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
	const who = [a.run.harness, a.run.model].filter(Boolean).join(', ') || a.run.command;
	switch (a.state) {
		case 'working':
			return `agent working (${who})`;
		case 'waiting':
			return `agent waiting (${who})${a.why ? `: ${a.why}` : ''}`;
		case 'failed':
			return `agent failed (${who})${a.why ? `: ${a.why}` : ''}`;
		default:
			return `agent finished (${who})`;
	}
}

/** Whether any agent is still running, so its state is worth asking for again. */
export function anyRunning(h: HostAgent | null): boolean {
	return Object.values(h?.state?.stories ?? {}).some(
		(a) => a.state === 'working' || a.state === 'waiting'
	);
}
