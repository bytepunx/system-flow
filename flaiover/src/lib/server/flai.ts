// Where a flai binary is, if the image has one. Since S-0075 every read and write of the project is
// asked of flai on the host over the channel, and the image carries no flai; the one thing left that
// would run one in the container is the MCP bridge, which S-0076 removes. Until then it answers 503
// when there is none.
import { access } from 'node:fs/promises';
import { constants } from 'node:fs';
import { delimiter, join } from 'node:path';

let resolved: string | null | undefined;

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
