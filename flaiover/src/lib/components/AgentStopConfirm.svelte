<script lang="ts">
	// Stopping a story's agent ends what it is doing and cannot be taken back (S-0170): say what it
	// does before it happens. Closing the dialog changes nothing.
	import type { StoryActivity } from '$lib/activity';

	let {
		story,
		activity,
		status = '',
		error = null,
		onconfirm,
		oncancel
	}: {
		story: string;
		activity: StoryActivity;
		/** The story's state, when it is known. */
		status?: string;
		/** Why flai would not stop it, when it would not. */
		error?: string | null;
		onconfirm: () => void | Promise<void>;
		oncancel: () => void;
	} = $props();

	let busy = $state(false);
	// An agent that ended waiting for an answer has no process: stopping it only keeps the answer
	// from starting it again.
	const asked = $derived(!!activity.run.ended && activity.run.outcome === 'asked');
	const thread = $derived(activity.thread ?? activity.run.thread);

	async function confirm() {
		busy = true;
		try {
			await onconfirm();
		} finally {
			busy = false;
		}
	}
</script>

<div
	class="fixed inset-0 z-50 flex items-start justify-center bg-ink/40 p-4 pt-24"
	role="presentation"
	onclick={(e) => {
		if (e.target === e.currentTarget && !busy) oncancel();
	}}
>
	<div
		class="w-full max-w-lg rounded-lg border border-line bg-surface p-5 text-sm shadow-lg"
		role="dialog"
		aria-modal="true"
		aria-labelledby="agent-stop-title"
	>
		<h2 id="agent-stop-title" class="text-base font-semibold">Stop {story}'s agent?</h2>
		<p class="mt-1 text-muted">{activity.run.agent}, started {activity.run.started}</p>
		<div class="mt-3 rounded border border-warn bg-warn-soft p-2" role="note" data-warning>
			{#if asked}
				<p>
					It ended waiting for an answer{thread ? ` to ${thread}` : ''}. Stopped, it is not started
					again when the answer comes.
				</p>
			{:else}
				<p>
					Its process, and every process it started, is ended now{activity.run.pid
						? ` (pid ${activity.run.pid})`
						: ''}: whatever it is in the middle of stops there, and it cannot be resumed.
				</p>
			{/if}
			<p class="mt-2">
				What it changed in the story's worktree stays as it left it, committed or not; nothing is
				undone.
			</p>
			<p class="mt-2">
				{story} stays {status ? `in ${status}` : 'where it is'}, with no agent. It gets another only
				when you press Retry on its page or move it back to ready.
			</p>
		</div>
		{#if error}
			<p class="mt-3 rounded border border-danger bg-danger-soft p-2 text-danger" role="alert">
				{error}
			</p>
		{/if}
		<div class="mt-4 flex justify-end gap-2">
			<button
				type="button"
				class="rounded border border-line-strong px-3 py-1"
				disabled={busy}
				onclick={oncancel}>{asked ? 'Keep it' : 'Keep it running'}</button
			>
			<button
				type="button"
				class="rounded bg-danger px-3 py-1 text-on-primary disabled:opacity-50"
				disabled={busy}
				onclick={confirm}
				data-testid="agent-stop-confirm">{busy ? 'Stopping…' : 'Stop the agent'}</button
			>
		</div>
	</div>
</div>
