<script lang="ts">
	// The header's two-tier site menu (S-0172): the groups in the top row, and the open group's pages
	// in a row below it. A group opens on hover, click, or touch; with none opened that way, the group
	// of the page shown is open. The open group and the page shown are bold in the accent colour, and a
	// group carries its pages' indicators. Each group's row follows its button in the DOM, so the
	// keyboard reaches a group's pages straight after it; order-last draws the row below the header.
	import type { Snippet } from 'svelte';
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import { projectState } from '$lib/project.svelte';
	import { SITE_MENU, groupBadges, locate, type GroupKey, type PageKey } from '$lib/sitemenu';
	import InboxBadge from './InboxBadge.svelte';

	let { start, end }: { start?: Snippet; end?: Snippet } = $props();

	const HREF: Record<PageKey, string> = {
		overview: resolve('/'),
		board: resolve('/board'),
		inbox: resolve('/inbox'),
		threads: resolve('/threads'),
		activity: resolve('/activity'),
		charts: resolve('/charts/[kind]', { kind: 'cycle-time' }),
		adrs: resolve('/adrs'),
		docs: resolve('/docs/[...path]', { path: '' }),
		search: resolve('/search'),
		updates: resolve('/host'),
		settings: resolve('/settings')
	};
	const BADGE = { inbox: InboxBadge };

	/** The group opened by hover, click, or touch; null leaves the page's own group open. */
	let opened = $state<GroupKey | null>(null);
	const here = $derived(locate(page.url.pathname));
	const open = $derived(opened ?? here.group);

	// Arriving on a page shows its group again, whichever was opened on the way.
	$effect(() => {
		void page.url.pathname;
		opened = null;
	});

	/** A menu link keeps the project it was followed from, so a copied link shows the same one. */
	function withProject(href: string): string {
		const key = projectState.current;
		return key ? `${href}?project=${encodeURIComponent(key)}` : href;
	}

	// A touch ends with a pointerleave, so only a mouse or a pen opens and closes a group by hovering.
	const hovering = (e: PointerEvent) => e.pointerType !== 'touch';
</script>

<!-- Leaving the menu or pressing Escape in it only undoes an opened group: every group opens from its
button, so the nav's listeners add no interaction that is not also reachable by the keyboard. -->
<!-- svelte-ignore a11y_no_noninteractive_element_interactions -->
<nav
	class="mx-auto flex max-w-6xl flex-wrap items-center gap-x-6 gap-y-2 px-4 py-3"
	aria-label="Site"
	onpointerleave={(e) => hovering(e) && (opened = null)}
	onkeydown={(e) => e.key === 'Escape' && (opened = null)}
>
	{@render start?.()}
	{#each SITE_MENU as g (g.key)}
		{@const active = open === g.key}
		<button
			type="button"
			class="text-sm {active
				? 'font-semibold text-accent'
				: 'text-ink-soft hover:text-accent'} cursor-pointer"
			aria-expanded={open === g.key}
			aria-controls="site-menu-{g.key}"
			data-menu-group={g.key}
			data-active={active}
			onpointerenter={(e) => hovering(e) && (opened = g.key)}
			onclick={() => (opened = g.key)}
			>{g.label}{#each groupBadges(g) as b (b)}{@const Badge = BADGE[b]}<Badge />{/each}</button
		>
		{#if open === g.key}
			<div
				id="site-menu-{g.key}"
				class="order-last flex basis-full flex-wrap gap-x-5 gap-y-1 border-t border-line pt-2"
				data-menu-pages={g.key}
			>
				{#each g.pages as p (p.key)}
					{@const current = here.page === p.key}
					<!-- eslint-disable svelte/no-navigation-without-resolve -- HREF is resolve()d; only ?project= is added -->
					<a
						class="text-sm {current
							? 'font-semibold text-accent'
							: 'text-ink-soft hover:text-accent'}"
						href={withProject(HREF[p.key])}
						aria-current={current ? 'page' : undefined}
						data-menu-page={p.key}
						>{p.label}{#if p.badge}{@const Badge = BADGE[p.badge]}<Badge />{/if}</a
					>
					<!-- eslint-enable svelte/no-navigation-without-resolve -->
				{/each}
			</div>
		{/if}
	{/each}
	{@render end?.()}
</nav>
