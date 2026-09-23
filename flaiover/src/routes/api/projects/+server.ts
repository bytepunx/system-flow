import { connectedProjects, registry, type ConnectedProject } from '$lib/server/agent';
import { respond } from '$lib/server/respond';
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
 */
export type ProjectGlance = ConnectedProject & {
	review?: number;
	threadsAwaiting?: number;
	agentAttending?: boolean;
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

export const GET: RequestHandler = () =>
	respond(async () => ({ projects: await Promise.all(connectedProjects().map(glanceOf)) }));
