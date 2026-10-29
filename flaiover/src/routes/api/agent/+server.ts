import { agent, heldBudget, type Refusal } from '$lib/server/agent';
import { respond } from '$lib/server/respond';
import type { RequestHandler } from './$types';

/** What the operator is told of a flai this dashboard refused as unsigned, and what to do (S-0237). */
function refusedError(r: Refusal): string {
	const who = r.flai ? `flai ${r.flai}` : 'a flai';
	return `this dashboard refused ${who} on the host: unsigned (${r.reason}); install a release of flai on the host with flai host upgrade, or with flai self-upgrade and then flai serve stop and flai dashboard; a development build of flai carries no stamp and is refused`;
}

/**
 * GET: whether a flai on the host has this dashboard connected (ADR-0029), and when it has, what it
 * says about the project: a round trip over the channel, so "connected" means it answers. Beside it,
 * `budget`, the providers flai serve holds every agent start on (S-0534), present only while it holds one.
 */
export const GET: RequestHandler = () =>
	respond(async () => {
		const status = agent().status();
		// a flai refused as unsigned is never connected; the refusal is kept so the page can say why.
		// While a flai is connected the project is served, and a refused second flai is not news here.
		if (!status.connected)
			return status.refused ? { ...status, error: refusedError(status.refused) } : status;
		// a repository offered for import has no project.info to ask for; answering at all is enough
		if (status.candidate) return status;
		if (status.missing?.length)
			return {
				...status,
				error: `flai ${status.flai} on the host is older than this dashboard and lacks ${status.missing.length} of the things it asks for (${status.missing.slice(0, 3).join(', ')}${status.missing.length > 3 ? ', …' : ''}); upgrade flai on the host, then run flai serve stop and flai dashboard`
			};
		try {
			const [info, budget] = await Promise.all([
				agent().ask('project.info', {}, 2000),
				heldBudget(agent(), 2000)
			]);
			return budget ? { ...status, info, budget } : { ...status, info };
		} catch (err) {
			return { ...status, error: err instanceof Error ? err.message : String(err) };
		}
	});
