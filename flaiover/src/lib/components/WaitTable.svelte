<script lang="ts">
	// The longest waits of the window under the agent-waiting chart (S-0215): each with its item,
	// the thread it waited on or review, when it started and ended, its time in the window, and who
	// was awaited.
	import { resolve } from '$app/paths';
	import { human, type LongestWait } from '$lib/viz/charts';
	import { localTime } from '$lib/localtime';

	let { rows, type }: { rows: LongestWait[]; type: string } = $props();
</script>

{#if rows.length === 0}
	<p class="text-muted" data-testid="wait-none">No wait has a part in the window.</p>
{:else}
	<!-- eslint-disable svelte/no-navigation-without-resolve -- the path is resolve()d; the rule does not follow the query added to it -->
	<table class="min-w-full" data-testid="wait-table">
		<thead
			><tr class="text-left text-muted"
				><th class="pr-4">{type}</th><th class="pr-4">kind</th><th class="pr-4">thread</th><th
					class="pr-4">started</th
				><th class="pr-4">ended</th><th class="pr-4">waited</th><th>awaited</th></tr
			></thead
		><tbody
			>{#each rows as row (row.item + (row.thread ?? '') + row.started)}<tr
					><td class="pr-4"
						><a class="font-mono underline" href={resolve('/items/[id]', { id: row.item })}
							>{row.item}</a
						></td
					><td class="pr-4">{row.kind}</td><td class="pr-4"
						>{#if row.thread}<a
								class="font-mono underline"
								href={resolve('/items/[id]', { id: row.item }) +
									`?thread=${encodeURIComponent(row.thread)}`}>{row.thread}</a
							>{:else}review{/if}</td
					><td class="pr-4 font-mono whitespace-nowrap">{localTime(row.started)}</td><td
						class="pr-4 font-mono whitespace-nowrap">{row.ended ? localTime(row.ended) : 'open'}</td
					><td class="pr-4">{human(row.seconds)}</td><td>{row.awaited ?? '—'}</td></tr
				>{/each}</tbody
		>
	</table>
	<!-- eslint-enable svelte/no-navigation-without-resolve -->
{/if}
