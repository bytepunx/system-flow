<script lang="ts">
	// A standing notice while the clone holds an acceptance its remote has not got (S-0063).
	// It asks on mount, whenever `refresh` changes (the board bumps it on every change event),
	// when the window regains focus, and every minute while it is showing, so it clears by
	// itself after a push from this clone. Pushing is done on the host, by flai, as the operator.
	// When they have enabled the push action there (flai serve enable push, S-0078) this offers
	// it; when they have not, it says what to run by hand. Nothing here can enable it.
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
	let pushEnabled = $state(false);
	let pushing = $state(false);
	let pushed = $state<string | null>(null);
	let refused = $state<string | null>(null);

	async function ask() {
		try {
			const r = await api('/api/unpushed');
			if (!r.ok) return;
			const body = await r.json();
			unpushed = body.unpushed ?? null;
			problem = body.problem ?? null;
			pushEnabled = body.push_enabled === true;
		} catch {
			// keep what we had: a failed question is not news
		}
	}
	async function push() {
		if (pushing) return;
		pushing = true;
		pushed = refused = null;
		try {
			const r = await api('/api/unpushed', { method: 'POST' });
			const body = await r.json().catch(() => ({}));
			if (!r.ok) {
				refused = body.error ?? r.statusText;
				return;
			}
			const tags: string[] = body.unpushed?.tags ?? [];
			const published: string[] = body.published ?? [];
			pushed = body.pushed
				? `Pushed${tags.length ? ` with tags ${tags.join(', ')}` : ''}${published.length ? `; published ${published.join(', ')}` : ''}.`
				: `Nothing was pushed: ${body.reason ?? 'nothing was pending'}.`;
		} catch (e) {
			refused = e instanceof Error ? e.message : String(e);
		} finally {
			pushing = false;
			await ask();
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

{#if pushed && !mine}
	<p
		class="mb-3 rounded border border-good bg-good-soft p-2 text-sm text-good"
		role="status"
		data-testid="pushed"
	>
		{pushed}
	</p>
{/if}
{#if mine}
	<div
		class="mb-3 rounded border border-warn bg-warn-soft p-2 text-sm text-warn"
		role="status"
		data-testid="unpushed"
	>
		<p>
			<span class="font-semibold">{pushing ? 'Pushing:' : 'Accepted, not pushed:'}</span>
			{mine.acceptances.join(', ')}
			({mine.commits} commit{mine.commits === 1 ? '' : 's'} ahead of {mine.upstream}{mine.tags
				.length
				? `, tags ${mine.tags.join(', ')}`
				: ''}).
			{#if !pushEnabled}On the host, run
				<code class="rounded bg-surface px-1 text-ink">{mine.command}</code>.{/if}
		</p>
		{#if pushEnabled}
			<p class="mt-1">
				<button
					type="button"
					class="rounded border border-warn px-2 py-0.5 font-medium disabled:opacity-60"
					onclick={push}
					disabled={pushing}
					data-testid="push-now">{pushing ? 'Pushing…' : 'Push now'}</button
				>
				<span class="ml-1 text-xs"
					>flai pushes from the host with the operator's credentials, never forced, and publishes
					the template if its version moved.</span
				>
			</p>
		{/if}
		{#if refused}
			<p class="mt-1" data-testid="push-refused">
				<span class="font-semibold">Not pushed:</span>
				{refused} By hand, on the host:
				<code class="rounded bg-surface px-1 text-ink">{mine.command}</code>.
			</p>
		{/if}
		<p class="mt-1 text-xs">
			This is what this clone knows. A push made from another clone is not seen here until someone
			fetches.{#if !pushEnabled}
				Pushing from the board is off; the operator turns it on in a shell on the host with
				<code class="rounded bg-surface px-1 text-ink">flai serve enable push</code>.{/if}
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
