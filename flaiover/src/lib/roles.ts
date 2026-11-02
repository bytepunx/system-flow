// The role page (S-0420): each role the project runs as agent.status's state carries it under
// roles, by its name, the built-in ones included; the /api/roles/<name> answer; and what the page
// reads of a role: its runs, what its controls are, where its run's stream is read, which agent
// events concern it, and which of flai's writes each control is.

import type { ActivityDocument, RoleRun } from './strategic';

/** A role's host action by name, and whether the operator has it on for the project. */
export type RoleAction = { name: string; enabled: boolean };

/**
 * One role as agent.status's state carries it under roles (flai's hostapi RoleState): its name,
 * its verb, and its file; the item kinds its run takes and its triggers' kinds; its host action;
 * its newest run under way, every run under way, and its runs that ended, newest first, at most
 * twenty; and whether the operator holds it stopped, which only the orchestrator can be.
 */
export type RoleState<R extends RoleRun = RoleRun> = {
	name: string;
	/** Its verb, such as plan, orchestrate, analyze, curate, or review. */
	role: string;
	/** Its file, relative to the project's root. */
	path: string;
	/** The item kinds its run takes, epic or story; none for a role that runs on no item. */
	item: string[];
	/** Its triggers' kinds, asked, event, or schedule, each once. */
	trigger: string[];
	action: RoleAction;
	run: R | null;
	running: R[];
	runs: R[];
	held: boolean;
};

/** A role entry as flai may answer it, a list null or absent. */
export type RawRoleState = Omit<RoleState, 'item' | 'trigger' | 'running' | 'runs' | 'run'> & {
	item?: string[] | null;
	trigger?: string[] | null;
	run?: RoleRun | null;
	running?: RoleRun[] | null;
	runs?: RoleRun[] | null;
};

/** The entry with its missing lists read as empty and a missing run as none. */
export function roleState(raw: RawRoleState): RoleState {
	return {
		...raw,
		item: raw.item ?? [],
		trigger: raw.trigger ?? [],
		run: raw.run ?? null,
		running: raw.running ?? [],
		runs: raw.runs ?? []
	};
}

/** The /api/roles/<name> answer: the role as agent.status reports it, and its activity document. */
export type RoleView = { role: RoleState; activity: ActivityDocument };

/** Every run of the role flai knows of, those under way first, newest first. */
export const roleRuns = <R extends RoleRun>(role: RoleState<R>): R[] => [
	...role.running,
	...role.runs
];

/** The verb of the one role that has a hold, the orchestrator's (S-0228). */
export const HOLDING_ROLE = 'orchestrate';

/**
 * What the role page offers for a role: Run for a role its file lets the operator ask for, with an
 * item when its run takes one, but for the orchestrator, which is started rather than run; Stop
 * while the orchestrator is not held, and Start while it is, since only the orchestrator has a hold.
 * Until flai has role.stop and role.start (TH-0797), no other role has either.
 */
export type RoleControls = { run: boolean; item: boolean; stop: boolean; start: boolean };

export function roleControls(
	role: Pick<RoleState, 'role' | 'item' | 'trigger' | 'held'>
): RoleControls {
	const holds = role.role === HOLDING_ROLE;
	const asked = role.trigger.includes('asked') && !holds;
	return {
		run: asked,
		item: asked && role.item.length > 0,
		stop: holds && !role.held,
		start: holds && role.held
	};
}

/** What the operator may have flai do to a role from its page. */
export const ROLE_ACTIONS = ['run', 'stop', 'start'] as const;
export type RoleActionName = (typeof ROLE_ACTIONS)[number];

/** One of flai's writes, by method and params, or why the page's action is none. */
export type RoleWrite = { method: string; params: Record<string, unknown> } | { refused: string };

/**
 * Which of flai's writes a role page's action is. Run is role.run with the role's name and its item,
 * save for a built-in role, which flai starts by a method of its own: the planner's plan.run on its
 * item, the analyzer's analyze.run, and the curator's curate.run; the orchestrator is not run on
 * asking. Stop and Start are the orchestrator's orchestrate.stop and orchestrate.start: no other role
 * has a hold until flai has role.stop and role.start (TH-0797).
 */
export function roleWrite(
	role: Pick<RoleState, 'name' | 'role'>,
	action: RoleActionName,
	item = ''
): RoleWrite {
	if (action !== 'run') {
		if (role.role !== HOLDING_ROLE)
			return {
				refused: `only the orchestrator has a hold to ${action}; the ${role.name} has none`
			};
		return { method: `orchestrate.${action}`, params: {} };
	}
	switch (role.role) {
		case 'plan':
			return item
				? { method: 'plan.run', params: { id: item } }
				: { refused: `the ${role.name} plans an epic or a story: give its item` };
		case 'analyze':
			return { method: 'analyze.run', params: {} };
		case 'curate':
			return { method: 'curate.run', params: {} };
		case HOLDING_ROLE:
			return {
				refused: `flai serve runs the ${role.name} while its host action is on: Start it rather than Run it`
			};
		default:
			return { method: 'role.run', params: { name: role.name, ...(item ? { item } : {}) } };
	}
}

/**
 * Where a role's run's stream is read (AgentStream's story, and its plan or role): a planner run's
 * by its item, the orchestrator's and the analyzer's by their verbs, and any other role's by its
 * name, which reads its newest run.
 */
export function roleStream(
	role: Pick<RoleState, 'name' | 'role'>,
	run: Pick<RoleRun, 'item'>
): { story: string; plan?: boolean; role?: string } {
	if (role.role === 'plan' && run.item) return { story: run.item, plan: true };
	if (role.role === 'orchestrate' || role.role === 'analyze')
		return { story: role.name, role: role.role };
	return { story: role.name, role: role.name };
}

const ITEM = /^[EST]-\d+$/;

/**
 * Whether flai serve's agent event for id concerns the role: one that names the role, as the hub
 * names the orchestrator's, the analyzer's, the curator's, and a defined role's run on no item, or,
 * for a role whose run takes an item, one that names an item, since its run on one is named by it.
 */
export function concernsRole(role: Pick<RoleState, 'name' | 'item'>, id: string): boolean {
	return id === role.name || (role.item.length > 0 && ITEM.test(id));
}

const TRIGGER_WORDS: Record<string, string> = {
	asked: 'when you Run it',
	event: 'on the events its file names',
	schedule: 'when its schedule comes round'
};

/** What the role's host action does while it is on, from its triggers, for the page's "On:" line. */
export function triggerLine(role: Pick<RoleState, 'name' | 'role' | 'trigger'>): string {
	if (role.role === HOLDING_ROLE)
		return `flai serve runs the ${role.name} for as long as the action stays on, and starts it again when a run ends.`;
	const when = role.trigger.map((t) => TRIGGER_WORDS[t] ?? `on ${t}`);
	if (!when.length) return `flai serve starts the ${role.name} as its file says.`;
	const said = when.length < 2 ? when[0] : `${when.slice(0, -1).join(', ')} and ${when.at(-1)}`;
	return `flai serve starts the ${role.name} ${said}.`;
}
