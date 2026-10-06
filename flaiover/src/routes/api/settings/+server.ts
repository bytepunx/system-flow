import { repo, RepoError } from '$lib/server/repo';
import { respond } from '$lib/server/respond';
import { authDisabled, isSecure, sessionCookie, setToken } from '$lib/server/auth';
import {
	SETTINGS_KINDS,
	SETTINGS_READS,
	type SettingsKind,
	type SettingsRead,
	type SettingsView,
	type SharedCheck
} from '$lib/settings';
import type { RequestHandler } from './$types';

/**
 * GET: the host's settings as they apply to this project (S-0105), flai's settings.get: every host
 * action and where it is on, the default agent, the agent's command and harnesses, the checks,
 * the import folders, the MCP server, and whether the dashboard may change them, here and for
 * every project. No token is ever in it.
 */
export const GET: RequestHandler = () => respond(() => repo().ask<SettingsView>('settings.get'));

/**
 * POST {kind, ...params}: one change, flai's settings.<kind>, gated on the host by the settings
 * action (403 with the command that enables it). A rotated dashboard token is taken here at once:
 * the server checks the new one from now on, and the session that asked gets a cookie with it, so
 * it stays logged in while every other session and agent must log in again. The answer carries
 * the new login link, never the token alone.
 */
export const POST: RequestHandler = async ({ request, url }) => {
	const body = (await request.json().catch(() => ({}))) as Record<string, unknown>;
	const kind = body.kind as SettingsKind | SettingsRead;
	let cookie: string | null = null;
	const res = await respond(async () => {
		const params = { ...body };
		delete params.kind;
		if (SETTINGS_READS.includes(kind as SettingsRead)) {
			const { data, warnings } = await repo().run<SharedCheck[]>(`settings.${kind}`, params);
			return { entries: data ?? [], warnings };
		}
		if (!SETTINGS_KINDS.includes(kind as SettingsKind)) {
			throw new RepoError(
				400,
				`kind must be one of ${[...SETTINGS_KINDS, ...SETTINGS_READS].join(', ')}`
			);
		}
		const timeoutMs = kind === 'mcp_token' || kind === 'dashboard_token' ? 60000 : 30000;
		const { data, warnings } = await repo().write<Record<string, unknown>>(
			`settings.${kind}`,
			params,
			{ timeoutMs }
		);
		if (kind === 'dashboard_token') {
			const token = typeof data?.token === 'string' ? data.token : '';
			if (token && !authDisabled()) {
				setToken(token);
				cookie = sessionCookie(token, isSecure(url, request.headers));
			}
			return { login_url: data?.login_url, warnings };
		}
		return { ...(data ?? {}), warnings };
	});
	if (cookie) res.headers.append('set-cookie', cookie);
	return res;
};
