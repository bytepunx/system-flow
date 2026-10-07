<script lang="ts">
	// A banner an action leaves on the page until the next one (S-0151): what happened, and a
	// right-aligned X that dismisses it. The caller holds the text and clears it on ondismiss, so
	// the next action's notice shows again. Notices that show state, not an action's outcome, are
	// not dismissible and do not use this. A time flai names in its answer shows in the local zone.
	import { localTimes } from '$lib/localtime';

	let {
		text,
		ondismiss,
		class: classes = '',
		role = 'status',
		testid
	}: {
		text: string;
		ondismiss: () => void;
		/** The banner's border, colours, spacing, and size. */
		class?: string;
		role?: 'status' | 'alert';
		/** Set on the text, so that its textContent is the notice without the X. */
		testid?: string;
	} = $props();
</script>

<div class="flex items-start gap-2 {classes}" {role}>
	<p class="min-w-0 flex-1" data-testid={testid}>{localTimes(text)}</p>
	<button
		type="button"
		class="shrink-0 cursor-pointer px-1 leading-none opacity-70 hover:opacity-100"
		aria-label="Dismiss"
		title="Dismiss"
		data-testid="dismiss"
		onclick={ondismiss}><span aria-hidden="true">✕</span></button
	>
</div>
