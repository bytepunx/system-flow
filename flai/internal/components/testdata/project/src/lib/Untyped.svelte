<script>
	// No type: each kind comes from its default, or is unknown.
	let { name = 'x', step = 1, label } = $props();
</script>

<p>{name} {step} {label}</p>
