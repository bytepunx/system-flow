import { json } from '@sveltejs/kit';
import { repo } from '$lib/server/repo';
import { flaiBinary } from '$lib/server/flai';
import { agent } from '$lib/server/agent';

const TIMEOUT_MS = 2000;

function withTimeout<T>(p: Promise<T>): Promise<T> {
	return Promise.race([
		p,
		new Promise<T>((_, reject) => setTimeout(() => reject(new Error('timed out')), TIMEOUT_MS))
	]);
}

/** Readiness: each dependency checked with a short timeout; names the one that failed. */
export const GET = async () => {
	const checks: Record<string, { ok: boolean; detail?: string }> = {};
	// The project and its items are asked of flai on the host (ADR-0029): without it nothing is ready.
	const host = agent().status();
	checks.host_flai = {
		ok: host.connected,
		detail: host.connected
			? `flai ${host.flai}`
			: 'not connected; run flai dashboard, or flai serve start'
	};
	try {
		const m = await withTimeout(repo().manifest());
		checks.project = { ok: true, detail: m.name };
	} catch (e) {
		checks.project = { ok: false, detail: e instanceof Error ? e.message : String(e) };
	}
	try {
		await withTimeout(repo().items());
		checks.items = { ok: true };
	} catch (e) {
		checks.items = { ok: false, detail: e instanceof Error ? e.message : String(e) };
	}
	checks.flai = { ok: true, detail: (await flaiBinary()) ? 'writable' : 'read-only' };
	const ok = Object.values(checks).every((c) => c.ok);
	return json({ status: ok ? 'ready' : 'not ready', checks }, { status: ok ? 200 : 503 });
};
