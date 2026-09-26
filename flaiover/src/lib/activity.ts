// What a story's agent is doing (S-0104), as flai on the host reports it in agent.status: working
// (green), waiting for the designer or, queued by a retry, for room (yellow), failed (red), or
// worked, which shows no dot.

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

export type StoryActivity = {
	state: ActivityState;
	why?: string;
	run: AgentRun;
	thread?: string;
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

/** One line for a tooltip or a screen reader: what the agent is doing and why. */
export function activityLine(a: StoryActivity): string {
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
