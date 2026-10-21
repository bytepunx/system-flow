<script lang="ts">
	import './layout.css';
	import favicon from '$lib/assets/favicon.png';
	import { flushSync, onMount } from 'svelte';
	import { themeState } from '$lib/theme.svelte';
	import { inboxState } from '$lib/inbox.svelte';
	import MenuIcon from '@lucide/svelte/icons/menu';
	import NavRail from '$lib/components/NavRail.svelte';
	import HostFlai from '$lib/components/HostFlai.svelte';
	import HostFlaiBanner from '$lib/components/HostFlaiBanner.svelte';
	import UnsignedBanner from '$lib/components/UnsignedBanner.svelte';
	import { hostFlai } from '$lib/hostflai.svelte';
	import { projectState } from '$lib/project.svelte';
	import ProjectSwitcher from '$lib/components/ProjectSwitcher.svelte';
	import ImportPrompt from '$lib/components/ImportPrompt.svelte';
	import EditorDialog from '$lib/components/EditorDialog.svelte';
	import { withProject } from '$lib/sitemenu';
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import { afterNavigate } from '$app/navigation';

	let { children } = $props();
	const OVERVIEW = resolve('/');
	// The login page has no session: no rail, no menu button, no project picker or status pills.
	const signedIn = $derived(page.url.pathname !== '/login');
	// The board takes the whole main area beside the rail, its lanes one row (S-0462); every other
	// page keeps the reading width.
	const BOARD = resolve('/board');
	const wide = $derived(page.url.pathname === BOARD);
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

	// The nav rail (S-0458): docked at the left from Tailwind's md width, and below it a drawer the
	// menu button opens. The drawer is a native <dialog> opened modal: the browser closes it on Escape,
	// keeps focus inside it, makes the page behind it inert, and draws it in the top layer, so closed
	// or open it never widens the page. Closing it, however it closes, returns focus to the button.
	// Only one rail is in the DOM at a time, so NavRail's ids stay unique: the docked one gives way
	// while the drawer is open, which it only is below md, where the docked one is not shown.
	let drawer = $state<HTMLDialogElement>();
	let toggle = $state<HTMLButtonElement>();
	let drawerOpen = $state(false);
	function openDrawer() {
		drawerOpen = true;
		// The rail is drawn before the dialog opens, so the browser focuses its first group.
		flushSync();
		drawer?.showModal();
	}
	function closeDrawer() {
		drawer?.close();
	}
	function closed() {
		drawerOpen = false;
		toggle?.focus();
	}
	// A page chosen in the drawer closes it before the navigation ends, and SvelteKit then resets
	// focus to the body; afterNavigate runs after that reset, so focus goes back to the button there.
	let chose = false;
	function chosen() {
		chose = true;
		closeDrawer();
	}
	afterNavigate(() => {
		if (chose) toggle?.focus();
		chose = false;
	});
	// Widening the window to md docks the rail, so an open drawer closes.
	$effect(() => {
		if (typeof window.matchMedia !== 'function') return;
		const wide = window.matchMedia('(min-width: 48rem)');
		const onchange = (e: MediaQueryListEvent) => e.matches && closeDrawer();
		wide.addEventListener('change', onchange);
		return () => wide.removeEventListener('change', onchange);
	});
</script>

<svelte:head><link rel="icon" href={favicon} /><title>flaiover</title></svelte:head>

<!-- The rail runs the height of the window at the left, from the top, and the top bar and the page
sit beside it; the column is min-w-0 so a wide page cannot push the window sideways. -->
<div class="flex min-h-dvh bg-ground text-ink">
	{#if signedIn && !drawerOpen}
		<aside class="hidden w-56 shrink-0 border-r border-line bg-surface md:block" data-nav-docked>
			<div class="sticky top-0 h-dvh overflow-y-auto"><NavRail /></div>
		</aside>
	{/if}
	<div class="min-w-0 flex-1">
		<header class="border-b border-line bg-surface">
			<div
				class="mx-auto flex max-w-6xl flex-wrap items-center gap-x-4 gap-y-2 px-4 py-3"
				data-top-bar
			>
				{#if signedIn}
					<button
						bind:this={toggle}
						type="button"
						class="cursor-pointer rounded p-1 text-muted hover:text-ink md:hidden"
						aria-label="Open navigation"
						aria-expanded={drawerOpen}
						aria-controls="nav-drawer"
						data-nav-toggle
						onclick={openDrawer}><MenuIcon class="size-5" aria-hidden="true" /></button
					>
				{/if}
				<!-- eslint-disable svelte/no-navigation-without-resolve -- OVERVIEW is resolve()d; only ?project= is added -->
				<a
					class="font-semibold tracking-tight text-accent no-underline"
					href={withProject(OVERVIEW, projectState.current)}
					data-top-name>flaiover</a
				>
				<!-- eslint-enable svelte/no-navigation-without-resolve -->
				{#if signedIn}<ProjectSwitcher />{/if}
				<button
					type="button"
					class="rounded border border-line px-2 py-1 text-xs text-muted hover:text-ink"
					title="Cycle system, light, dark"
					aria-label={modeLabel[themeState.mode]}
					data-top-theme
					onclick={() => themeState.cycle()}
				>
					{modeLabel[themeState.mode]}
				</button>
				{#if signedIn}<span class="ml-auto" data-top-pills><HostFlai /></span>{/if}
			</div>
		</header>
		{#if signedIn}<HostFlaiBanner /><UnsignedBanner />{/if}
		<main class="mx-auto px-4 py-6 {wide ? '' : 'max-w-6xl'}">
			{#key projectState.current}
				{#if projectState.project?.candidate}
					<ImportPrompt project={projectState.project} />
				{:else}
					{@render children()}
				{/if}
			{/key}
		</main>
	</div>
	{#if signedIn}
		<!-- The drawer below md. Escape closes it natively; a click on the backdrop, which lands on the
		dialog itself, closes it too, a pointer convenience over what the keyboard already has. -->
		<dialog
			bind:this={drawer}
			id="nav-drawer"
			class="fixed inset-y-0 left-0 h-full max-h-none w-56 border-r border-line bg-surface text-ink backdrop:bg-ink/40"
			aria-label="Navigation"
			data-nav-drawer
			onclose={closed}
			onclick={(e) => e.target === e.currentTarget && closeDrawer()}
		>
			{#if drawerOpen}<NavRail onnavigate={chosen} />{/if}
		</dialog>
	{/if}
	<!-- The document editor over any page: links open it with openEditor (S-0378). -->
	{#if signedIn}<EditorDialog />{/if}
</div>
