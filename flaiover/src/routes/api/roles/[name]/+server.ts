// A role page's read and its Run, Stop, and Start (S-0420): the role as agent.status reports it under
// roles, by its name, and its activity document, each asked of flai on the host, and the write each
// control is, as roleWrite says.

import { repo, RepoError } from '$lib/server/repo';
import { respond } from '$lib/server/respond';
import type { StrategicStatus } from '$lib/server/agent';
import {
	ROLE_ACTIONS,
	roleState,
	roleWrite,
	type RawRoleState,
	type RoleActionName,
	type RoleState,
	type RoleView
} from '$lib/roles';
import { activityDocument, ROLE_NAME, type RawActivityDocument } from '$lib/strategic';
import type { RequestHandler } from './$types';

/** The path's name, or 400 when it is not a role's name, asking flai nothing. */
function named(name: string): string {
	if (!ROLE_NAME.test(name))
		throw new RepoError(
			400,
			`${JSON.stringify(name)} is not a role's name: one in kebab case, such as reviewer`
		);
	return name;
}

/** The role agent.status reports by the name, or 404 naming the roles it does report. */
async function reported(name: string): Promise<RoleState> {
	// with finished, since this page lists the role's ended runs (S-0611)
	const status = await repo().ask<StrategicStatus>('agent.status', { finished: true });
	const roles = (status.state?.roles ?? {}) as Record<string, RawRoleState>;
	const raw = roles[name];
	if (!raw) {
		const known = Object.keys(roles).sort();
		throw new RepoError(
			404,
			known.length
				? `the project runs no role named ${name}; it runs ${known.join(', ')}`
				: `flai on the host reports no roles: it may be older than this dashboard`
		);
	}
	return roleState(raw);
}

/**
 * GET: the role agent.status reports by the path's name, with its host action's state, its runs, and
 * whether it is held, and its activity document, wip/agents/<name>.md. Either failing is the
 * answer's error; a role flai does not report is 404.
 */
export const GET: RequestHandler = ({ params }) =>
	respond(async (): Promise<RoleView> => {
		const name = named(params.name);
		const [role, activity] = await Promise.all([
			reported(name),
			repo().ask<RawActivityDocument>('activity.document', { kind: name }).then(activityDocument)
		]);
		return { role, activity };
	});

/**
 * POST { action, item? }: run, stop, or start the role, as the write roleWrite names for it, and
 * answer what flai answered. It is 403 with what enables it while the operator has not turned the
 * role's host action on, 400 for an action that is none of them, an item that is not a string, or
 * one roleWrite refuses (Stop or Start of a role with no hold), and 404 for a role flai does not
 * report.
 */
export const POST: RequestHandler = ({ params, request }) =>
	respond(async () => {
		const name = named(params.name);
		const body = (await request.json().catch(() => null)) as {
			action?: unknown;
			item?: unknown;
		} | null;
		const action = body?.action;
		if (!ROLE_ACTIONS.includes(action as RoleActionName))
			throw new RepoError(
				400,
				`action must be one of ${ROLE_ACTIONS.join(', ')}, not ${JSON.stringify(action ?? null)}`
			);
		const item = body?.item ?? '';
		if (typeof item !== 'string')
			throw new RepoError(
				400,
				`item must be an epic's or a story's ID, not ${JSON.stringify(item)}`
			);
		const write = roleWrite(await reported(name), action as RoleActionName, item.trim());
		if ('refused' in write) throw new RepoError(400, write.refused);
		const { data } = await repo().write(write.method, write.params);
		return (data ?? {}) as Record<string, unknown>;
	});
