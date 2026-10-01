<script lang="ts">
	// The answer to a hand-written open question in a story's narrative (S-0090), on the inbox and
	// on the story's page (S-0173). Answering moves it out of Open questions into the narrative's
	// Decisions; the inbox is refreshed, and the question drops out of both lists once it does.
	import type { InboxEntry } from '$lib/inbox.svelte';
	import { inboxState } from '$lib/inbox.svelte';
	import { api } from '$lib/api';

	let { entry }: { entry: InboxEntry } = $props();

	let draft = $state('');
	let answering = $state(false);
	let failed = $state('');

	async function answer() {
		const text = draft.trim();
		if (!text || !entry.item) return;
		answering = true;
		failed = '';
		try {
			const r = await api(`/api/streams/${entry.item}/answer`, {
				method: 'POST',
				headers: { 'content-type': 'application/json' },
				body: JSON.stringify({ question: entry.title, answer: text })
			});
			const data = await r.json().catch(() => ({}));
			if (!r.ok) {
				failed = data.error ?? r.statusText;
				return;
			}
			draft = '';
			await inboxState.refresh();
		} finally {
			answering = false;
		}
	}
</script>

<form
	class="mt-2 flex gap-2"
	onsubmit={(ev) => {
		ev.preventDefault();
		void answer();
	}}
>
	<input
		class="min-w-0 flex-1 rounded border border-line-strong px-2 py-1 text-xs"
		placeholder="answer"
		data-testid="question-answer-input"
		bind:value={draft}
	/>
	<button
		type="submit"
		class="rounded border border-line-strong px-2 py-1 text-xs disabled:opacity-50"
		disabled={!draft.trim() || answering}
		data-testid="question-answer-submit">{answering ? 'answering…' : 'answer'}</button
	>
</form>
{#if failed}
	<p class="mt-1 text-xs text-danger" role="alert">{failed}</p>
{/if}
