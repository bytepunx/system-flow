import { json } from '@sveltejs/kit';
import { repo } from '$lib/server/repo';
import { flaiBinary } from '$lib/server/flai';

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
