<script>
	// The fixture of a dialog (S-0449), drawn as flaiover's are: a fixed
	// backdrop over the whole frame, centring a dialog taller than the page's
	// unsized frame, 160px, so that only its own size frames it whole.
</script>

<div class="backdrop">
	<div class="dialog" role="dialog">A dialog</div>
</div>

<style>
	.backdrop {
		position: fixed;
		inset: 0;
		display: flex;
		align-items: center;
		justify-content: center;
		background: rgb(0 0 0 / 0.4);
	}
	.dialog {
		box-sizing: border-box;
		width: 240px;
		height: 300px;
		padding: 16px;
		background: white;
	}
</style>
