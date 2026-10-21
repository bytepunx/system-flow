<script lang="ts">
	// /settings goes on to the first Configuration page (S-0459), and an old link to one of the
	// settings page's anchors (S-0458) to the page that holds that section, keeping the query and so
	// the project. The hash never reaches a load, so the browser goes on from here, in place of this
	// page in the history.
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import { onMount } from 'svelte';
	import { pageOfAnchor } from '$lib/settingspages';

	onMount(() => {
		const to = resolve('/settings/[page]', { page: pageOfAnchor(page.url.hash) });
		// eslint-disable-next-line svelte/no-navigation-without-resolve -- the path is resolve()d; the rule does not follow the query added to it
		void goto(to + page.url.search, { replaceState: true });
	});
</script>
