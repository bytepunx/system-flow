// The Vite config of the fixture sub-project the components integration test
// serves (S-0438): a plain Svelte sub-project, whose own config flai's
// harness config loads and merges its plugin into.
import { svelte } from '@sveltejs/vite-plugin-svelte';

export default {
	plugins: [svelte()]
};
