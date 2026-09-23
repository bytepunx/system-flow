<script lang="ts">
	// /settings: the host's settings for this project (S-0105), changed through flai once the
	// operator has turned on the settings host action in a shell.
	import SettingsPanel from '$lib/components/SettingsPanel.svelte';
</script>

<svelte:head><title>flaiover — settings</title></svelte:head>

<h1 class="mb-4 text-lg font-semibold text-ink">Settings</h1>
<SettingsPanel />
