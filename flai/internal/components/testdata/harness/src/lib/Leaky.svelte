<script>
	// The fixture of a component that leaves things behind (S-0449): as it
	// mounts, it adds an element where it is mounted, ahead of its own markup,
	// and one to the body, as a portal does, and removes neither, so that
	// unmounting it leaves both.
	let el;
	$effect(() => {
		const stray = document.createElement('p');
		stray.className = 'leaked';
		el.parentNode.prepend(stray);
		const portal = document.createElement('div');
		portal.className = 'leaked';
		document.body.append(portal);
	});
</script>

<p class="leaky" bind:this={el}>Leaky</p>
