<script lang="ts">
	import './layout.css';
	import favicon from '$lib/assets/favicon.svg';
	import { resolve } from '$app/paths';

	let { children } = $props();
	const nav = [
		{ href: resolve('/'), label: 'Overview' },
		{ href: resolve('/board'), label: 'Board' },
		{ href: resolve('/charts/[kind]', { kind: 'cycle-time' }), label: 'Charts' },
		{ href: resolve('/docs/[...path]', { path: '' }), label: 'Docs' },
		{ href: resolve('/adrs'), label: 'ADRs' },
		{ href: resolve('/search'), label: 'Search' }
	];
</script>

<svelte:head><link rel="icon" href={favicon} /><title>flaiover</title></svelte:head>

<div class="min-h-screen bg-zinc-50 text-zinc-900 dark:bg-zinc-950 dark:text-zinc-100">
	<header class="border-b border-zinc-200 bg-white dark:border-zinc-800 dark:bg-zinc-900">
		<nav class="mx-auto flex max-w-6xl items-center gap-6 px-4 py-3">
			<span class="font-semibold tracking-tight">flaiover</span>
			{#each nav as { href, label } (href)}
				<a
					class="text-sm text-zinc-600 hover:text-zinc-900 dark:text-zinc-400 dark:hover:text-zinc-100"
					{href}>{label}</a
				>
			{/each}
		</nav>
	</header>
	<main class="mx-auto max-w-6xl px-4 py-6">
		{@render children()}
	</main>
</div>
