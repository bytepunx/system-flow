// Test support, never imported by the application: a Repo's questions answered by the flai built
// from this tree (`flai hostapi`), the same method table flai serve offers a dashboard, read from a
// fixture folder. So the tests exercise the answers the dashboard will really get, and flaiover
// keeps no parser of its own to test against (S-0073).
import { execFile } from 'node:child_process';
import { existsSync } from 'node:fs';
import { resolve } from 'node:path';
import { AgentError } from './agent';
import type { Ask } from './repo';

export const flaiBin = process.env.FLAI_BIN ?? resolve('../bin/flai');
export const haveFlai = existsSync(flaiBin);

export function flaiAsk(root: string): Ask {
	return <T>(method: string, params: Record<string, unknown> = {}) =>
		new Promise<T>((done, fail) => {
			execFile(
				flaiBin,
				['hostapi', method, JSON.stringify(params)],
				{ cwd: root, maxBuffer: 64 * 1024 * 1024, env: { ...process.env, FLAI_AGENT: '' } },
				(err, stdout) => {
					let out: unknown;
					try {
						out = JSON.parse(stdout);
					} catch {
						return fail(err ?? new Error(`flai hostapi ${method}: no JSON`));
					}
					const e = (out as { error?: { code: number; message: string } } | null)?.error;
					if (err && e) return fail(new AgentError(502, e.message, e.code));
					if (err) return fail(err);
					done(out as T);
				}
			);
		});
}
