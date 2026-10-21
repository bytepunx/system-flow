// End-to-end tests of the built dashboard in Chromium (scripts/flaiover-e2e.sh, make flaiover-e2e).
// The server is the image's own entry, node server.js over the adapter-node build that the script
// makes first, not vite preview: only server.js answers the WebSocket upgrades of /agent and
// /api/events (ADR-0029, ADR-0182). It listens on 127.0.0.1 alone, with a login token and a flai
// credential made for the run in a scratch directory, which the tests read to sign in and to connect
// their fake flai; nothing of this repository's or the operator's is used.
import { defineConfig } from '@playwright/test';
import { randomBytes } from 'node:crypto';
import { mkdtempSync, rmSync, writeFileSync } from 'node:fs';
import { tmpdir } from 'node:os';
import { join } from 'node:path';

const port = Number(process.env.FLAIOVER_E2E_PORT ?? 4173);

// The config is loaded again in every worker, which inherits the runner's environment: the runner
// makes the scratch directory once, and removes it when it exits.
let dir = process.env.FLAIOVER_E2E_DIR;
if (!dir) {
	const made = mkdtempSync(join(tmpdir(), 'flaiover-e2e-'));
	dir = process.env.FLAIOVER_E2E_DIR = made;
	writeFileSync(join(made, 'token'), randomBytes(24).toString('hex') + '\n', { mode: 0o600 });
	writeFileSync(join(made, 'agent-key'), randomBytes(24).toString('hex') + '\n', { mode: 0o600 });
	process.on('exit', () => rmSync(made, { recursive: true, force: true }));
}

export default defineConfig({
	testDir: 'e2e',
	testMatch: '**/*.e2e.{ts,js}',
	use: { baseURL: `http://127.0.0.1:${port}`, browserName: 'chromium' },
	webServer: {
		command: 'node server.js',
		url: `http://127.0.0.1:${port}/_health`,
		reuseExistingServer: false,
		env: {
			HOST: '127.0.0.1',
			PORT: String(port),
			FLAIOVER_TOKEN_FILE: join(dir, 'token'),
			FLAIOVER_AGENT_KEY_FILE: join(dir, 'agent-key'),
			// the tests' fake flai carries no release stamp
			FLAIOVER_ALLOW_UNSIGNED: '1'
		}
	}
});
