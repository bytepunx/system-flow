import { connectedProjects, registry, type ConnectedProject } from '$lib/server/agent';
import { respond } from '$lib/server/respond';
import type { ProjectsView } from '$lib/settings';
import type { RequestHandler } from './$types';

/**
 * GET: every project flai on the host has ever named here, for a switcher (S-0080), each with a
 * glance at its state: stories in review, threads awaiting the designer, and whether an agent is
 * attending. No `?project=` scopes this one, on purpose: it is how a page learns what projects
 * exist at all, before it can name one.
 *
 * Each project is asked directly (not through repo(), which is scoped to one project at a time) and
 * on its own short timeout, all at once: a project whose flai is slow or has stopped answering
 * shows without its glance rather than holding up the rest of the list.
 *
 * The projects flai serve serves (S-0122) come from one connected project's settings.get: each is
 * marked served, one this dashboard has no connection to is added, and one that is not connected
 * carries flai serve's last error, so the switcher can tell it apart and point at the settings page.
 */
export type ProjectGlance = ConnectedProject & {
	review?: number;
	threadsAwaiting?: number;
	agentAttending?: boolean;
	/** flai serve serves it (S-0122). */
	served?: boolean;
	/** Why flai serve has it not connected, or not served, when it says. */
	lastError?: string;
};

const GLANCE_TIMEOUT_MS = 3000;

async function glanceOf(p: ConnectedProject): Promise<ProjectGlance> {
	// a repository offered for import has no board or inbox to glance at (S-0098)
	if (!p.connected || p.candidate) return p;
	const hub = registry().peek(p.key);
	if (!hub) return p;
	const [board, inbox, agent] = await Promise.allSettled([
		hub.ask<{ columns?: Record<string, { type: string }[]> }>('board.get', {}, GLANCE_TIMEOUT_MS),
		hub.ask<{ counts?: Record<string, number> }>('inbox.designer', {}, GLANCE_TIMEOUT_MS),
		hub.ask<{ state?: { running?: unknown } }>('agent.status', {}, GLANCE_TIMEOUT_MS)
	]);
	return {
		...p,
		review:
			board.status === 'fulfilled'
				? (board.value.columns?.review?.filter((c) => c.type === 'story').length ?? 0)
				: undefined,
		threadsAwaiting: inbox.status === 'fulfilled' ? (inbox.value.counts?.thread ?? 0) : undefined,
		agentAttending: agent.status === 'fulfilled' ? !!agent.value.state?.running : undefined
	};
}

/** What flai serve serves, asked of the first connected project, or undefined when none answers. */
async function servedOf(projects: ConnectedProject[]): Promise<ProjectsView | undefined> {
	const via = projects.find((p) => p.connected && !p.candidate);
	const hub = via && registry().peek(via.key);
	if (!hub) return undefined;
	try {
		const view = await hub.ask<{ host?: { projects?: ProjectsView } }>(
			'settings.get',
			{},
			GLANCE_TIMEOUT_MS
		);
		return view.host?.projects;
	} catch {
		// an older flai, or one that does not answer: the list is what the registry knows
		return undefined;
	}
}

function withServed(list: ProjectGlance[], view: ProjectsView | undefined): ProjectGlance[] {
	if (!view?.served) return list;
	const byKey = new Map(view.served.map((s) => [s.key, s]));
	const out = list.map((p): ProjectGlance => {
		const s = byKey.get(p.key);
		if (!s || p.candidate) return p;
		const why = s.last_error ?? s.reason;
		return { ...p, served: true, ...(!p.connected && why ? { lastError: why } : {}) };
	});
	for (const s of view.served) {
		if (list.some((p) => p.key === s.key)) continue;
		const why = s.last_error ?? s.reason;
		out.push({
			key: s.key,
			name: s.name,
			connected: false,
			served: true,
			...(why ? { lastError: why } : {})
		});
	}
	return out.sort((a, b) => a.key.localeCompare(b.key));
}

export const GET: RequestHandler = () =>
	respond(async () => {
		const known = connectedProjects();
		const [glances, served] = await Promise.all([
			Promise.all(known.map(glanceOf)),
			servedOf(known)
		]);
		return { projects: withServed(glances, served) };
	});
