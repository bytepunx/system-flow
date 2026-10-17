<script>
	let { name = 'world' } = $props();
</script>

<p class="hello">Hello, {name}!</p>

<style>
	.hello {
		color: var(--fixture-ink);
	}
</style>
