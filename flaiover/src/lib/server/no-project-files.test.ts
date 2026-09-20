import { describe, expect, it } from 'vitest';
import { readdirSync, readFileSync, statSync } from 'node:fs';
import { join, relative } from 'node:path';

// S-0074, ADR-0029: the dashboard's server reads no file of the project; it asks flai on the host.
// The files below touch the filesystem for something else, each for the reason given. Anything new
// that imports node:fs fails here, and the answer is a method in flai, not a line in this list.
const ALLOWED: Record<string, string> = {
	'src/lib/server/auth.ts': 'the login token, a secret handed to the container',
	'src/lib/server/agent.ts': 'the agent credential, a secret handed to the container',
	'src/lib/server/flai.ts': 'looks for a flai binary for the MCP bridge, which S-0076 removes',
	'src/lib/server/testing.ts': 'test support, never imported by the application'
};

function sources(dir: string, out: string[] = []): string[] {
	for (const name of readdirSync(dir)) {
		const path = join(dir, name);
		if (statSync(path).isDirectory()) sources(path, out);
		else if (/\.(ts|js)$/.test(name) && !/\.test\.ts$/.test(name)) out.push(path);
	}
	return out;
}

describe('the server reads no project file', () => {
	it('imports node:fs only where it is accounted for', () => {
		const server = [
			...sources('src/lib/server'),
			...sources('src/routes').filter((p) => /\+server\.ts$|\.server\.ts$/.test(p)),
			'src/hooks.server.ts'
		];
		expect(server.length).toBeGreaterThan(30);
		const offenders = server
			.map((p) => relative('.', p).split('\\').join('/'))
			.filter((p) => /from ['"](node:)?fs(\/promises)?['"]/.test(readFileSync(p, 'utf8')))
			.filter((p) => !(p in ALLOWED));
		expect(offenders).toEqual([]);
	});

	it('no longer depends on a YAML parser, a search library, or a file watcher', () => {
		const pkg = JSON.parse(readFileSync('package.json', 'utf8')) as {
			dependencies: Record<string, string>;
		};
		for (const gone of ['yaml', 'minisearch', 'chokidar'])
			expect(Object.keys(pkg.dependencies)).not.toContain(gone);
	});
});
