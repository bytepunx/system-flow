<script>
	// The fixture of the harness's in-place update (S-0440): open is its own
	// state, which a change of label or children must leave as it is, and
	// onToggle is called with a value, an event, and a function, the last two
	// of which do not serialise.
	let { label = 'Menu', onToggle, children } = $props();
	let open = $state(false);
</script>

<div class="stateful">
	<button
		type="button"
		onclick={(event) => {
			open = !open;
			onToggle?.(open, event, () => {});
		}}>{label}</button
	>
	{#if open}<p class="open">open</p>{/if}
	{#if children}<div class="slot">{@render children()}</div>{/if}
</div>
