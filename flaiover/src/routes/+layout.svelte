<script lang="ts">
	import './layout.css';
	import favicon from '$lib/assets/favicon.svg';
	import { resolve } from '$app/paths';
	import { onMount } from 'svelte';
	import { themeState } from '$lib/theme.svelte';
	import { inboxState } from '$lib/inbox.svelte';
	import InboxBadge from '$lib/components/InboxBadge.svelte';
	import HostFlai from '$lib/components/HostFlai.svelte';
	import HostFlaiBanner from '$lib/components/HostFlaiBanner.svelte';
	import { hostFlai } from '$lib/hostflai.svelte';
	import { projectState } from '$lib/project.svelte';
	import ProjectSwitcher from '$lib/components/ProjectSwitcher.svelte';
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
		{ href: resolve('/search'), label: 'Search' },
		{ href: resolve('/host'), label: 'Host' }
	];
	const modeLabel = {
		system: 'theme: system',
		light: 'theme: light',
		dark: 'theme: dark'
	} as const;
	onMount(() => {
		themeState.start();
	});

	// The inbox, the host flai badge, and the project switcher all need a session; the login page has
	// none yet. Reactive, not onMount-once (I-0032): the login page's own redirect is a client-side
	// navigation, so onMount alone would never see signedIn become true once the app has booted on
	// /login itself, which following a login link always does.
	let started = false;
	$effect(() => {
		const signedIn = page.url.pathname !== '/login';
		if (signedIn && !started) {
			started = true;
			inboxState.start();
			hostFlai.start();
			void projectState.refresh();
		} else if (!signedIn && started) {
			started = false;
			inboxState.stop();
			hostFlai.stop();
		}
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
			{#if page.url.pathname !== '/login'}<ProjectSwitcher />{/if}
			<span class="ml-auto"
				>{#if page.url.pathname !== '/login'}<HostFlai />{/if}</span
			>
			<button
				type="button"
				class="rounded border border-line px-2 py-1 text-xs text-muted hover:text-ink"
				title="Cycle system, light, dark"
				aria-label={modeLabel[themeState.mode]}
				onclick={() => themeState.cycle()}
			>
				{modeLabel[themeState.mode]}
			</button>
		</nav>
	</header>
	{#if page.url.pathname !== '/login'}<HostFlaiBanner />{/if}
	<main class="mx-auto max-w-6xl px-4 py-6">
		{@render children()}
	</main>
</div>
