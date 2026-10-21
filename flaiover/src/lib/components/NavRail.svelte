<script lang="ts">
	// The site's nav rail (S-0458, replacing S-0172's two-tier header menu): the five groups in a
	// column, each a button that shows or hides its pages, each page an icon and its label. An
	// expanded group's icon is in the accent colour and a collapsed one's muted; a collapsed group
	// still carries its pages' indicators, so Flow shows the Inbox count. The group of the page shown
	// expands on each arrival; the groups the operator expands or collapses are kept per browser
	// (navrail.svelte). License sits at the foot; Overview and Updates are in the top bar (TH-0553).
	// The rail fills the width its container gives it.
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import ChevronRightIcon from '@lucide/svelte/icons/chevron-right';
	import { projectState } from '$lib/project.svelte';
	import { navRail, type NavRailState } from '$lib/navrail.svelte';
	import { SITE_MENU, groupBadges, locate, withProject } from '$lib/sitemenu';
	import InboxBadge from './InboxBadge.svelte';

	let { onnavigate, rail = navRail }: { onnavigate?: () => void; rail?: NavRailState } = $props();

	const BADGE = { inbox: InboxBadge };
	const LICENSE = resolve('/license');
	const here = $derived(locate(page.url.pathname, page.url.hash));

	// Arriving on a page in a group expands it, before the rail is drawn.
	$effect.pre(() => {
		const g = here.group;
		if (g) rail.expand(g);
	});
</script>

<nav class="flex h-full min-w-0 flex-col gap-4 p-2 text-sm" aria-label="Site" data-nav-rail>
	<ul class="flex flex-col gap-1">
		{#each SITE_MENU as g (g.key)}
			{@const open = rail.isExpanded(g.key)}
			<li>
				<button
					type="button"
					class="flex w-full cursor-pointer items-center gap-2 rounded px-2 py-1.5 text-left text-ink hover:bg-raised"
					aria-expanded={open}
					aria-controls="nav-rail-{g.key}"
					data-rail-group={g.key}
					onclick={() => rail.toggle(g.key)}
				>
					<g.icon
						class="size-5 shrink-0 {open ? 'text-accent' : 'text-muted'}"
						aria-hidden="true"
						data-rail-icon
					/>
					<span class="min-w-0 truncate font-medium">{g.label}</span>
					{#if !open}
						{#each groupBadges(g) as b (b)}
							{@const Badge = BADGE[b]}
							<Badge />
						{/each}
					{/if}
					<ChevronRightIcon
						class="ml-auto size-4 shrink-0 text-muted {open ? 'rotate-90' : ''}"
						aria-hidden="true"
					/>
				</button>
				<ul
					id="nav-rail-{g.key}"
					class="mt-1 ml-4 flex flex-col gap-0.5 border-l border-line pl-2"
					hidden={!open}
					data-rail-pages={g.key}
				>
					{#each g.pages as p (p.key)}
						<li>
							{#if p.soon}
								<span
									class="flex items-center gap-2 rounded px-2 py-1 text-muted"
									data-rail-page={p.key}
									data-soon
								>
									<p.icon class="size-4 shrink-0" aria-hidden="true" />
									<span class="min-w-0 truncate">{p.label}</span>
									<span class="ml-auto rounded border border-line px-1 text-xs" aria-hidden="true">
										soon
									</span>
									<span class="sr-only">(coming soon)</span>
								</span>
							{:else}
								{@const current = here.page === p.key}
								<!-- eslint-disable svelte/no-navigation-without-resolve -- p.href is resolve()d; only ?project= is added -->
								<a
									class="flex items-center gap-2 rounded px-2 py-1 hover:bg-raised {current
										? 'font-semibold text-accent'
										: 'text-ink-soft hover:text-accent'}"
									href={withProject(p.href, projectState.current)}
									aria-current={current ? 'page' : undefined}
									data-rail-page={p.key}
									onclick={() => onnavigate?.()}
								>
									<p.icon
										class="size-4 shrink-0 {current ? 'text-accent' : 'text-muted'}"
										aria-hidden="true"
									/>
									<span class="min-w-0 truncate">{p.label}</span>
									{#if p.badge}
										{@const Badge = BADGE[p.badge]}
										<Badge />
									{/if}
								</a>
								<!-- eslint-enable svelte/no-navigation-without-resolve -->
							{/if}
						</li>
					{/each}
				</ul>
			</li>
		{/each}
	</ul>
	<div class="mt-auto border-t border-line px-2 pt-2">
		<!-- eslint-disable svelte/no-navigation-without-resolve -- LICENSE is resolve()d; only ?project= is added -->
		<a
			class="text-xs {page.url.pathname === LICENSE
				? 'font-semibold text-accent'
				: 'text-muted hover:text-accent'}"
			href={withProject(LICENSE, projectState.current)}
			aria-current={page.url.pathname === LICENSE ? 'page' : undefined}
			data-rail-license
			onclick={() => onnavigate?.()}>License</a
		>
		<!-- eslint-enable svelte/no-navigation-without-resolve -->
	</div>
</nav>
