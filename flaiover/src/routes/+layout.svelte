<script lang="ts">
	import './layout.css';
	import favicon from '$lib/assets/favicon.png';
	import { onMount } from 'svelte';
	import { themeState } from '$lib/theme.svelte';
	import { inboxState } from '$lib/inbox.svelte';
	import SiteMenu from '$lib/components/SiteMenu.svelte';
	import HostFlai from '$lib/components/HostFlai.svelte';
	import HostFlaiBanner from '$lib/components/HostFlaiBanner.svelte';
	import { hostFlai } from '$lib/hostflai.svelte';
	import { projectState } from '$lib/project.svelte';
	import ProjectSwitcher from '$lib/components/ProjectSwitcher.svelte';
	import ImportPrompt from '$lib/components/ImportPrompt.svelte';
	import { page } from '$app/state';

	let { children } = $props();
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
	let listTimer: ReturnType<typeof setInterval> | null = null;
	$effect(() => {
		const signedIn = page.url.pathname !== '/login';
		if (signedIn && !started) {
			started = true;
			inboxState.start();
			hostFlai.start();
			void projectState.refresh();
			// A project whose flai connects later joins the switcher without a reload (S-0095).
			listTimer = setInterval(() => void projectState.refresh(), 30000);
		} else if (!signedIn && started) {
			started = false;
			inboxState.stop();
			hostFlai.stop();
			if (listTimer) clearInterval(listTimer);
			listTimer = null;
		}
	});

	// Switching project (S-0095): the badges in the header outlive the page, so they are pointed at
	// the new project here; the page itself is keyed on the project below and remounts.
	let shownProject: string | null | undefined;
	$effect(() => {
		const current = projectState.current;
		if (shownProject !== undefined && current !== shownProject && started) {
			inboxState.restart();
			void hostFlai.refresh();
		}
		shownProject = current;
	});
</script>

<svelte:head><link rel="icon" href={favicon} /><title>flaiover</title></svelte:head>

<div class="min-h-screen bg-ground text-ink">
	<header class="border-b border-line bg-surface">
		<SiteMenu>
			{#snippet start()}<span class="font-semibold tracking-tight text-accent">flaiover</span
				>{/snippet}
			{#snippet end()}
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
			{/snippet}
		</SiteMenu>
	</header>
	{#if page.url.pathname !== '/login'}<HostFlaiBanner />{/if}
	<main class="mx-auto max-w-6xl px-4 py-6">
		{#key projectState.current}
			{#if projectState.project?.candidate}
				<ImportPrompt project={projectState.project} />
			{:else}
				{@render children()}
			{/if}
		{/key}
	</main>
</div>
