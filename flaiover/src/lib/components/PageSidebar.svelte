<script lang="ts">
	// A page's sidebar of its sections (S-0460): a left-aligned list of entries beside the page's
	// content, each a link to its hash. A pane entry shows its pane, through the pane snippet, and
	// hides the others; a scroll entry scrolls to the element with its id, and is current while that
	// is in view. The entry chosen is kept in the URL's hash (pagesidebar.svelte), so back, forward,
	// and a shared link open the same pane: a pane's link is followed as the browser follows any
	// hash link, adding a history entry, and a scroll entry's replaces the hash as it scrolls, at
	// once when the reader asks for reduced motion. From md the list is sticky at the left of the
	// content; below md it is one row above it that scrolls within itself. A pane's id should not
	// be an element id on the page, or following its link scrolls to that element.
	import { tick, untrack, type Snippet } from 'svelte';
	import { PageSidebarState, scrollBehavior, type SidebarEntry } from '$lib/pagesidebar.svelte';

	let {
		label,
		entries,
		fallback,
		match,
		current = $bindable(),
		pane,
		children
	}: {
		/** The sidebar's accessible name, such as "Story sections". */
		label: string;
		/** The entries, in order; an entry's children are listed under it. */
		entries: SidebarEntry[];
		/** The entry shown when the hash names none, as on arrival; the first entry by default. */
		fallback?: string;
		/** The entry that holds a hash that is not an entry's id, such as a message's anchor. */
		match?: (hash: string) => string | undefined;
		/** The entry current, bound: setting it chooses that entry without changing the hash. */
		current?: string;
		/** Renders the pane shown, given its id; left out on a page with no panes. */
		pane?: Snippet<[string]>;
		/** Content after the pane, on every pane, such as the sections scroll entries name. */
		children?: Snippet;
	} = $props();

	const sidebar = new PageSidebarState(() => ({ entries, fallback, match }));
	let region = $state<HTMLElement>();

	// The bound current both ways: a value the page sets is chosen, and the entry current is given
	// back. synced is the value last passed either way, so neither direction undoes the other.
	let synced: string | undefined;
	$effect.pre(() => {
		const want = current;
		untrack(() => {
			if (want === undefined || want === synced) return;
			synced = want;
			sidebar.choose(want);
		});
	});
	$effect.pre(() => {
		const now = sidebar.current;
		untrack(() => {
			if (now === current) return;
			synced = now;
			current = now;
		});
	});

	$effect(() => sidebar.listen());
	$effect(() => {
		void sidebar.watched;
		return untrack(() => sidebar.watch());
	});

	// Arriving with a scroll entry's hash scrolls to its section once the content is drawn.
	$effect(() => {
		const arrived = untrack(() => sidebar.entry(sidebar.chosen));
		if (arrived?.mode === 'scroll')
			void tick().then(() => document.getElementById(arrived.id)?.scrollIntoView?.());
	});

	/** A click that opens the link elsewhere, which the browser handles. */
	const elsewhere = (e: MouseEvent) =>
		e.button !== 0 || e.metaKey || e.ctrlKey || e.shiftKey || e.altKey;

	// A pane's link is left to the browser, which adds the history entry and fires hashchange; it is
	// chosen here at once, and the content brought back into view if the reader had scrolled past it.
	async function choosePane(e: MouseEvent, id: string) {
		if (elsewhere(e)) return;
		sidebar.choose(id);
		await tick();
		if (region && region.getBoundingClientRect().top < 0)
			region.scrollIntoView?.({ block: 'start', behavior: scrollBehavior() });
	}

	// A scroll entry's link is handled here, in the capture phase so that the router never sees it:
	// its section is scrolled to, smoothly unless the reader asks for reduced motion, and the hash
	// replaced, as the browser's own jump would land at once.
	async function chooseSection(e: MouseEvent, id: string) {
		if (elsewhere(e)) return;
		e.preventDefault();
		sidebar.choose(id);
		history.replaceState(history.state, '', `#${encodeURIComponent(id)}`);
		await tick();
		document.getElementById(id)?.scrollIntoView?.({ block: 'start', behavior: scrollBehavior() });
	}
</script>

{#snippet list(items: SidebarEntry[], top: boolean)}
	<ul
		class={top
			? 'flex gap-1 overflow-x-auto md:flex-col md:overflow-visible'
			: 'flex gap-1 md:mt-0.5 md:ml-3 md:flex-col md:border-l md:border-line md:pl-2'}
		data-sidebar-list={top ? 'top' : 'nested'}
	>
		{#each items as e (e.id)}
			{@const isCurrent = sidebar.current === e.id}
			{@const shown = isCurrent || sidebar.pane === e.id}
			<li class="flex shrink-0 gap-1 md:block">
				<a
					class="flex items-center gap-2 rounded px-2 py-1 whitespace-nowrap hover:bg-raised {shown
						? 'font-semibold text-accent'
						: 'text-ink-soft hover:text-accent'}"
					href="#{encodeURIComponent(e.id)}"
					aria-current={isCurrent ? 'true' : undefined}
					data-sidebar-entry={e.id}
					data-mode={e.mode}
					onclick={e.mode === 'pane' ? (ev) => choosePane(ev, e.id) : undefined}
					onclickcapture={e.mode === 'scroll' ? (ev) => chooseSection(ev, e.id) : undefined}
				>
					<span class="min-w-0 truncate">{e.label}</span>
					{#if e.count !== undefined}
						<span
							class="ml-auto rounded-full bg-raised px-1.5 text-xs font-normal text-muted"
							data-sidebar-count>{e.count}</span
						>
					{/if}
				</a>
				{#if e.children?.length}{@render list(e.children, false)}{/if}
			</li>
		{/each}
	</ul>
{/snippet}

<div class="flex min-w-0 flex-col gap-4 md:flex-row md:items-start md:gap-6" data-page-sidebar>
	<nav
		class="min-w-0 text-sm md:sticky md:top-4 md:max-h-screen md:w-48 md:shrink-0 md:overflow-y-auto"
		aria-label={label}
		data-sidebar-nav
	>
		{@render list(entries, true)}
	</nav>
	<div bind:this={region} class="min-w-0 flex-1" data-sidebar-content>
		{#if pane && sidebar.pane}
			<div data-sidebar-pane={sidebar.pane}>{@render pane(sidebar.pane)}</div>
		{/if}
		{@render children?.()}
	</div>
</div>
