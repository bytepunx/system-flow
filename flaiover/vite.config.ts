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
			}
		]
	}
});
