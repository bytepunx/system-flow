<script lang="ts">
	// A log path on the host as a link to the log page, /ops/log?path=<path>, with the project this
	// tab views (S-0785), so the operator reads the log in the dashboard rather than in a shell. Its
	// text is the path, in the code style paths are shown in, unless given a label.
	import { logHref } from '$lib/logs';
	import { projectState } from '$lib/project.svelte';

	let {
		path,
		label = ''
	}: {
		/** The log's path as flai names it: absolute, ~/-prefixed, or from the project's root. */
		path: string;
		/** What the link says in place of the path. */
		label?: string;
	} = $props();

	const href = $derived(logHref(path, projectState.current));
</script>

<!-- eslint-disable svelte/no-navigation-without-resolve -- logHref resolve()s the page; only ?path= and ?project= are added -->
<a class="underline" {href} title={`Show ${path}`} data-testid="log-link" data-path={path}
	>{#if label}{label}{:else}<code class="font-mono break-all">{path}</code>{/if}</a
>
<!-- eslint-enable svelte/no-navigation-without-resolve -->
