<script lang="ts">
	// At the top of the done column: everything accepted and unreleased since each component's
	// last tag (S-0087), and a button to publish it — apply, tag, and push together, as the
	// operator, when they have enabled the push action (flai serve enable push, S-0078; the same
	// gate the old automatic push used). The plan is shown before it runs, same as an acceptance's
	// preview; what happened after, success or flai's error verbatim, same as the review page.
	import { api } from '$lib/api';

	type PendingItem = { id: string; title: string; level: string };
	type PendingPlan = {
		component: { name: string };
		level: string;
		from: string;
		to: string;
		items: PendingItem[];
	};

	let {
		plans,
		enabled,
		onpublished
	}: {
		plans: PendingPlan[];
		enabled: boolean;
		/** Told after a publish attempt, success or not, so the board and the plan can be asked
		 * again. */
		onpublished?: () => void;
	} = $props();

	let publishing = $state(false);
	let result = $state<string | null>(null);
	let failed = $state<string | null>(null);

	async function publish() {
		if (publishing) return;
		publishing = true;
		result = failed = null;
		try {
			const r = await api('/api/publish', { method: 'POST' });
			const body = await r.json().catch(() => ({}));
			if (!r.ok) {
				failed = body.error ?? r.statusText;
				return;
			}
			const tags: string[] = body.tags ?? [];
			const published: string[] = body.published ?? [];
			result = body.pushed
				? `Published${tags.length ? ` ${tags.join(', ')}` : ''}, pushed.${published.length ? ` Published ${published.join(', ')}.` : ''}`
				: `Applied and tagged locally; not pushed.`;
		} catch (e) {
			failed = e instanceof Error ? e.message : String(e);
		} finally {
			publishing = false;
			onpublished?.();
		}
	}
</script>

{#if result}
	<p
		class="mb-2 rounded border border-good bg-good-soft p-2 text-xs text-good"
		role="status"
		data-testid="published"
	>
		{result}
	</p>
{/if}
{#if plans.length > 0}
	<div
		class="mb-2 rounded border border-warn bg-warn-soft p-2 text-xs text-warn"
		role="status"
		data-testid="publish-pending"
	>
		<p class="font-semibold">Ready to publish:</p>
		<ul class="mt-1 space-y-1">
			{#each plans as p (p.component.name)}
				<li>
					<span class="font-mono">{p.component.name}</span>
					{p.from} → {p.to} ({p.level}): {p.items.map((it) => it.id).join(', ')}
				</li>
			{/each}
		</ul>
		{#if enabled}
			<p class="mt-1">
				<button
					type="button"
					class="rounded border border-warn px-2 py-0.5 font-medium disabled:opacity-60"
					onclick={publish}
					disabled={publishing}
					data-testid="publish-now">{publishing ? 'Publishing…' : 'Publish'}</button
				>
			</p>
		{:else}
			<p class="mt-1">
				Publishing from the board is off; the operator turns it on in a shell on the host with
				<code class="rounded bg-surface px-1 text-ink">flai serve enable push</code>.
			</p>
		{/if}
		{#if failed}
			<p class="mt-1" data-testid="publish-refused">
				<span class="font-semibold">Not published:</span>
				{failed}
			</p>
		{/if}
	</div>
{/if}
