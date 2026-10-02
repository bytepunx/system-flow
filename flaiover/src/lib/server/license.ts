// The license the dashboard is distributed under (S-0231): LICENSE.md, which the image carries
// beside the server at /app and a checkout keeps at the repository root. It is the dashboard's
// own file, not the project's, so it is read here rather than asked of flai on the host.
import { readFileSync } from 'node:fs';
import { resolve } from 'node:path';

export const LICENSE_FILE = 'LICENSE.md';

/** Where the license may be, first match wins: FLAIOVER_LICENSE, then beside the server, then one up (a checkout). */
export function licensePaths(
	env: Record<string, string | undefined> = process.env,
	cwd: string = process.cwd()
): string[] {
	const paths = [resolve(cwd, LICENSE_FILE), resolve(cwd, '..', LICENSE_FILE)];
	if (env.FLAIOVER_LICENSE) paths.unshift(resolve(cwd, env.FLAIOVER_LICENSE));
	return paths;
}

export type License = { name: string; text: string; path: string };

/** The license's title is its first heading. */
export function licenseName(text: string): string {
	const first = text.split('\n', 1)[0] ?? '';
	return first.replace(/^#+/, '').trim();
}

/** The first license file that reads, or null when the image carries none. */
export function readLicense(paths: string[] = licensePaths()): License | null {
	for (const path of paths) {
		try {
			const text = readFileSync(path, 'utf8');
			return { name: licenseName(text), text, path };
		} catch {
			// not here; try the next place
		}
	}
	return null;
}
