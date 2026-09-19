// Every mutation and every metric goes through the flai binary (ADR-0016).
// The server only spawns it; rules, front matter writing, and numbers stay
// in one implementation.
import { execFile } from 'node:child_process';
import { access } from 'node:fs/promises';
import { constants } from 'node:fs';
import { join, delimiter } from 'node:path';
import { RepoError } from './repo';

let resolved: string | null | undefined;

/** Locate flai: FLAI_BIN, then PATH. Cached per process; null when absent. */
export async function flaiBinary(): Promise<string | null> {
	if (resolved !== undefined) return resolved;
	const candidates: string[] = [];
	if (process.env.FLAI_BIN) candidates.push(process.env.FLAI_BIN);
	for (const dir of (process.env.PATH ?? '').split(delimiter)) {
		if (dir) candidates.push(join(dir, 'flai'), join(dir, 'flai.exe'));
	}
	for (const c of candidates) {
		try {
			await access(c, constants.X_OK);
			resolved = c;
			return c;
		} catch {
			// next
		}
	}
	resolved = null;
	return null;
}

/** For tests: forget the cached location. */
export function resetFlaiBinary() {
	resolved = undefined;
}

export type FlaiResult<T> = { data: T; warnings: string[] };

export type FlaiOptions = {
	/** Written to flai's standard input. */
	input?: string;
	/** flai exit codes that are an HTTP status of their own, with the JSON on stdout as the response data. */
	exitStatus?: Record<number, number>;
};

/**
 * Run flai with --json in the project directory and return its JSON.
 * A failure becomes a RepoError carrying flai's rule text (from the fatal
 * log event on stderr) so the browser can show why a write was refused.
 */
export async function flai<T = unknown>(
	projectDir: string,
	args: string[],
	opt: FlaiOptions = {}
): Promise<FlaiResult<T>> {
	const bin = await flaiBinary();
	if (!bin)
		throw new RepoError(
			503,
			'flai is not available to this dashboard; writes and metrics are disabled (set FLAI_BIN or put flai on PATH)'
		);
	const env = {
		...process.env,
		FLAI_CONFIG: process.env.FLAI_CONFIG ?? join(projectDir, '.flai-cache', 'config.json'),
		FLAI_CACHE_DIR: process.env.FLAI_CACHE_DIR ?? join(projectDir, '.flai-cache', 'cache'),
		FLAI_AGENT: process.env.FLAI_AGENT ?? 'flaiover',
		LOG_FORMAT: 'json'
	};
	return new Promise((resolvePromise, reject) => {
		const child = execFile(
			bin,
			[...args, '--json'],
			{ cwd: projectDir, env, maxBuffer: 16 * 1024 * 1024 },
			(err, stdout, stderr) => {
				const events = parseEvents(stderr);
				const warnings = events
					.filter((e) => e.level === 'WARN')
					.map((e) => String(e.detail ?? e.msg));
				if (err) {
					const fatal = events.find((e) => e.level === 'FATAL');
					const message = String(fatal?.err ?? stderr.trim() ?? err.message);
					// An exit code the caller knows (flai doc save: conflict, refused) carries its
					// payload as JSON on stdout.
					const status = typeof err.code === 'number' ? opt.exitStatus?.[err.code] : undefined;
					if (status) {
						let data: Record<string, unknown> | undefined;
						try {
							data = JSON.parse(stdout);
						} catch {
							data = undefined;
						}
						reject(new RepoError(status, message.replace(/^(conflict|refused):\s*/, ''), data));
						return;
					}
					reject(
						new RepoError(message.startsWith('rule:') ? 400 : 500, message.replace(/^rule:\s*/, ''))
					);
					return;
				}
				try {
					resolvePromise({
						data: stdout.trim() ? (JSON.parse(stdout) as T) : (null as T),
						warnings
					});
				} catch {
					reject(new RepoError(500, `flai returned non-JSON output: ${stdout.slice(0, 200)}`));
				}
			}
		);
		// Content goes on standard input, never in arguments or through a shell.
		if (opt.input !== undefined) child.stdin?.end(opt.input);
	});
}

function parseEvents(stderr: string): Record<string, unknown>[] {
	const out: Record<string, unknown>[] = [];
	for (const line of stderr.split('\n')) {
		if (!line.startsWith('{')) continue;
		try {
			out.push(JSON.parse(line));
		} catch {
			// not an event line
		}
	}
	return out;
}
