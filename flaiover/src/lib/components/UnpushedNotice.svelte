<script lang="ts">
	// A standing notice while the clone holds an acceptance its remote has not got (S-0063).
	// It asks on mount, whenever `refresh` changes (the board bumps it on every change event),
	// when the window regains focus, and every minute while it is showing, so it clears by
	// itself after a push from this clone. There is no retry button: pushing is done on the
	// host, or by the container at acceptance when the operator gave it a key (ADR-0026).
	import { api } from '$lib/api';
	import { onMount } from 'svelte';

	type Unpushed = {
		upstream: string;
		commits: number;
		acceptances: string[];
		tags: string[];
		command: string;
	};

	let { item, refresh = 0 }: { item?: string; refresh?: number } = $props();
	let unpushed = $state<Unpushed | null>(null);
	let problem = $state<string | null>(null);

	async function ask() {
		try {
			const r = await api('/api/unpushed');
			if (!r.ok) return;
			const body = await r.json();
			unpushed = body.unpushed ?? null;
			problem = body.problem ?? null;
		} catch {
			// keep what we had: a failed question is not news
		}
	}
	onMount(() => {
		const timer = setInterval(() => (unpushed || problem) && ask(), 60_000);
		const focus = () => ask();
		window.addEventListener('focus', focus);
		return () => {
			clearInterval(timer);
			window.removeEventListener('focus', focus);
		};
	});
	$effect(() => {
		void refresh;
		void ask();
	});

	const mine = $derived(
		unpushed && (!item || unpushed.acceptances.includes(item)) ? unpushed : null
	);
</script>

{#if mine}
	<div
		class="mb-3 rounded border border-warn bg-warn-soft p-2 text-sm text-warn"
		role="status"
		data-testid="unpushed"
	>
		<p>
			<span class="font-semibold">Accepted, not pushed:</span>
			{mine.acceptances.join(', ')}
			({mine.commits} commit{mine.commits === 1 ? '' : 's'} ahead of {mine.upstream}{mine.tags
				.length
				? `, tags ${mine.tags.join(', ')}`
				: ''}). On the host, run
			<code class="rounded bg-surface px-1 text-ink">{mine.command}</code>.
		</p>
		<p class="mt-1 text-xs">
			This is what this clone knows. A push made from another clone is not seen here until someone
			fetches.
		</p>
	</div>
{:else if problem && !item}
	<p
		class="mb-3 rounded border border-warn bg-warn-soft p-2 text-sm text-warn"
		role="status"
		data-testid="unpushed-problem"
	>
		{problem}
	</p>
{/if}
