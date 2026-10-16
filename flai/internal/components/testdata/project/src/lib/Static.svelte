<script lang="ts">
	// No props at all; $props.id() is not $props().
	const id = $props.id();
</script>

<p {id}>static</p>
