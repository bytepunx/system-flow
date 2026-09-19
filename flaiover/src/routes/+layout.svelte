<script lang="ts">
	import './layout.css';
	import favicon from '$lib/assets/favicon.svg';
	import { resolve } from '$app/paths';
	import { onMount } from 'svelte';
	import { themeState } from '$lib/theme.svelte';
	import { inboxState } from '$lib/inbox.svelte';
	import InboxBadge from '$lib/components/InboxBadge.svelte';
	import { page } from '$app/state';

	let { children } = $props();
	const nav = [
		{ href: resolve('/'), label: 'Overview' },
		{ href: resolve('/board'), label: 'Board' },
		{ href: resolve('/inbox'), label: 'Inbox' },
		{ href: resolve('/activity'), label: 'Activity' },
		{ href: resolve('/charts/[kind]', { kind: 'cycle-time' }), label: 'Charts' },
		{ href: resolve('/docs/[...path]', { path: '' }), label: 'Docs' },
		{ href: resolve('/adrs'), label: 'ADRs' },
		{ href: resolve('/search'), label: 'Search' }
	];
	const modeLabel = {
		system: 'theme: system',
		light: 'theme: light',
		dark: 'theme: dark'
	} as const;
	onMount(() => {
		themeState.start();
		// The inbox needs a session; the login page has none yet.
		if (page.url.pathname !== '/login') inboxState.start();
		return () => inboxState.stop();
	});
</script>

<svelte:head><link rel="icon" href={favicon} /><title>flaiover</title></svelte:head>

<div class="min-h-screen bg-ground text-ink">
	<header class="border-b border-line bg-surface">
		<nav class="mx-auto flex max-w-6xl flex-wrap items-center gap-x-6 gap-y-2 px-4 py-3">
			<span class="font-semibold tracking-tight text-accent">flaiover</span>
			{#each nav as { href, label } (href)}
				<a class="text-sm text-ink-soft hover:text-accent" {href}
					>{label}{#if label === 'Inbox'}<InboxBadge />{/if}</a
				>
			{/each}
			<button
				type="button"
				class="ml-auto rounded border border-line px-2 py-1 text-xs text-muted hover:text-ink"
				title="Cycle system, light, dark"
				aria-label={modeLabel[themeState.mode]}
				onclick={() => themeState.cycle()}
			>
				{modeLabel[themeState.mode]}
			</button>
		</nav>
	</header>
	<main class="mx-auto max-w-6xl px-4 py-6">
		{@render children()}
	</main>
</div>
