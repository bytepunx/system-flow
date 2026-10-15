<!--
	A proposal for the design system beside the main branch's (S-0384): the gallery previewed with
	each one's tokens, side by side, and the tokens they differ in, per theme, each marked changed,
	added, or removed, with main's value and the proposal's. A proposal whose token file has
	problems shows them; one with no tokens to compare says so rather than list every token removed.
-->
<script lang="ts">
	import type { DesignSystem, DesignTheme, DesignToken } from '$lib/server/repo';
	import { diffDesign, valueText, type TokenChange } from '$lib/design/diff';
	import { unsafeVariable } from '$lib/design/preview';
	import PreviewFrame from './PreviewFrame.svelte';

	let {
		main,
		proposal,
		story,
		title
	}: {
		/** The main branch's design system, as design.get answers it. */
		main: DesignSystem;
		/** The proposal's design system, as design.get answers it for the story. */
		proposal: DesignSystem;
		/** The proposal's story ID, such as S-0384. */
		story: string;
		/** The proposal's story title. */
		title: string;
	} = $props();

	const LABEL: Record<DesignTheme['name'], string> = { light: 'Light', dark: 'Dark' };
	const TONE: Record<TokenChange, string> = {
		changed: 'border-warn text-warn',
		added: 'border-good text-good',
		removed: 'border-danger text-danger'
	};

	const differences = $derived(diffDesign(main, proposal));

	const tokensOf = (ds: DesignSystem, name: DesignTheme['name']) =>
		ds.themes.find((t) => t.name === name)?.tokens ?? [];

	/** The value, when it may be written into a style; undefined when it is empty or refused. */
	const safe = (value: string | undefined): string | undefined =>
		value && unsafeVariable('--v', value) === null ? value : undefined;

	const h4 = 'mb-3 text-lg font-semibold';
	const th = 'py-2 pr-4 text-left font-medium text-muted';
</script>

{#snippet value(token: DesignToken | undefined, absent: string)}
	{#if token}
		<div class="flex items-start gap-2">
			{#if token.type === 'color'}
				<span
					class="mt-0.5 inline-block size-6 shrink-0 rounded border border-line-strong"
					style:background-color={safe(token.css)}
					aria-hidden="true"
					data-testid="swatch"
				></span>
			{/if}
			<span class="font-mono text-xs break-all">{valueText(token)}</span>
		</div>
	{:else}
		<span class="text-xs text-muted">{absent}</span>
	{/if}
{/snippet}

<section aria-labelledby="proposal-compare-heading" data-testid="proposal-compare">
	<h3 id="proposal-compare-heading" class="text-lg font-semibold">
		Proposal {story}: {title}
	</h3>
	<p class="mt-1 text-sm text-muted">
		The sample components with the main branch's tokens and with {story}'s, and the tokens they
		differ in.
	</p>
	{#if proposal.problems.length}
		<div
			class="mt-4 rounded border border-danger bg-danger-soft p-3 text-sm text-danger"
			data-testid="proposal-problems"
		>
			<p class="font-medium">{story}'s token file has problems:</p>
			<ul class="mt-1 list-disc pl-5">
				{#each proposal.problems as p, i (i)}
					<li>
						line {p.line}{#if p.path}&nbsp;· <code>{p.path}</code>{/if}: {p.message}
					</li>
				{/each}
			</ul>
		</div>
	{/if}

	<div class="mt-4 grid grid-cols-1 gap-6 xl:grid-cols-2">
		<section class="min-w-0" aria-labelledby="proposal-compare-main" data-testid="compare-main">
			<h4 id="proposal-compare-main" class={h4}>Main branch</h4>
			<PreviewFrame
				light={tokensOf(main, 'light')}
				dark={tokensOf(main, 'dark')}
				title="The main branch's design system preview"
			/>
		</section>
		<section
			class="min-w-0"
			aria-labelledby="proposal-compare-proposal"
			data-testid="compare-proposal"
		>
			<h4 id="proposal-compare-proposal" class={h4}>Proposal {story}</h4>
			<PreviewFrame
				light={tokensOf(proposal, 'light')}
				dark={tokensOf(proposal, 'dark')}
				title="{story}'s proposed design system preview"
			/>
		</section>
	</div>

	<section
		class="mt-6"
		aria-labelledby="proposal-compare-differences"
		data-testid="proposal-differences"
	>
		<h4 id="proposal-compare-differences" class="text-lg font-semibold">Tokens that differ</h4>
		{#if !proposal.themes.length}
			<p class="mt-2 text-sm text-warn" data-testid="proposal-no-tokens">
				{story}'s token file {proposal.exists ? 'does not parse' : 'is missing'}, so it has no
				tokens to compare.
			</p>
		{:else if !differences.length}
			<p class="mt-2 text-sm text-muted" data-testid="proposal-no-difference">
				{story} changes no token; its guidance may still differ.
			</p>
		{:else}
			{#each differences as group (group.theme)}
				<h5 id="proposal-differences-{group.theme}" class="mt-4 font-semibold">
					{LABEL[group.theme]}
				</h5>
				<table class="mt-2 w-full text-sm" aria-labelledby="proposal-differences-{group.theme}">
					<thead>
						<tr>
							<th class={th}>Token</th>
							<th class={th}>Change</th>
							<th class={th}>Main branch</th>
							<th class={th}>Proposal</th>
						</tr>
					</thead>
					<tbody>
						{#each group.differences as d (d.path)}
							<tr
								class="border-t border-line align-top"
								data-testid="difference-{group.theme}-{d.path}"
								data-change={d.change}
							>
								<td class="py-2 pr-4 font-mono text-xs">{d.path}</td>
								<td class="py-2 pr-4">
									<span
										class={['rounded border px-1.5 text-xs font-medium', TONE[d.change]]}
										data-testid="change">{d.change}</span
									>
								</td>
								<td class="py-2 pr-4" data-side="main">
									{@render value(d.main, 'not on the main branch')}
								</td>
								<td class="py-2 pr-4" data-side="proposal">
									{@render value(d.proposal, 'not in the proposal')}
								</td>
							</tr>
						{/each}
					</tbody>
				</table>
			{/each}
		{/if}
	</section>
</section>
