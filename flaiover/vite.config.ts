import tailwindcss from '@tailwindcss/vite';
import { defineConfig } from 'vitest/config';
import adapter from '@sveltejs/adapter-node';
import { sveltekit } from '@sveltejs/kit/vite';
import { execSync } from 'node:child_process';
import { readFileSync } from 'node:fs';

// Build-time version for flaiover_build_info: the nearest flaiover release tag
// when building inside the repo, else the package version.
const pkg = JSON.parse(readFileSync(new URL('./package.json', import.meta.url), 'utf8')) as {
	version: string;
};
const version = (() => {
	try {
		return execSync("git describe --tags --match 'flaiover/v*' --abbrev=0", {
			stdio: ['ignore', 'pipe', 'ignore']
		})
			.toString()
			.trim()
			.replace(/^flaiover\/v/, '');
	} catch {
		return pkg.version;
	}
})();

export default defineConfig({
	define: { __FLAIOVER_VERSION__: JSON.stringify(version) },
	plugins: [
		tailwindcss(),
		{
			// The development server's hook for /agent, as server.js is the image's (ADR-0029).
			// Vite's own WebSocket upgrades are left alone.
			name: 'flaiover-agent-upgrade',
			configureServer(server) {
				server.httpServer?.on('upgrade', (req, socket, head) => {
					if ((req.url ?? '').split('?')[0] === '/agent')
						(
							globalThis as {
								__flaioverAgentUpgrade?: (r: unknown, s: unknown, h: unknown) => void;
							}
						).__flaioverAgentUpgrade?.(req, socket, head);
				});
			}
		},
		sveltekit({
			compilerOptions: {
				// Force runes mode for the project, except for libraries. Can be removed in svelte 6.
				runes: ({ filename }) =>
					filename.split(/[/\\]/).includes('node_modules') ? undefined : true
			},
			adapter: adapter()
		})
	],
	test: {
		expect: { requireAssertions: true },
		projects: [
			{
				extends: './vite.config.ts',
				test: {
					name: 'server',
					environment: 'node',
					env: { LOG_LEVEL: 'error' },
					include: ['src/**/*.{test,spec}.{js,ts}'],
					exclude: ['src/**/*.svelte.{test,spec}.{js,ts}']
				}
			},
			{
				// Component tests: Svelte's client build in jsdom.
				extends: './vite.config.ts',
				resolve: { conditions: ['browser'] },
				test: {
					name: 'client',
					environment: 'jsdom',
					include: ['src/**/*.svelte.{test,spec}.{js,ts}']
				}
			}
		]
	}
});
