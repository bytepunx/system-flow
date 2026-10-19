<script lang="ts">
	// Props typed by imports through $lib, one that does not exist, and a package that is not installed.
	import type { Deep, Status, Story, Tree } from '$lib/types';
	import type { Missing } from '$lib/missing';
	import type { Ghost } from 'not-installed';

	let {
		label,
		story,
		stories,
		tree,
		status,
		deep,
		missing,
		ghost,
		extra,
		local
	}: {
		label: string;
		story: Story;
		stories?: Story[] | null;
		tree: Tree;
		status: Status;
		deep: Deep;
		missing: Missing;
		ghost: Ghost[];
		extra: Record<string, number>;
		local: { at: Date; ok?: boolean | undefined };
	} = $props();
</script>

<p>{label} {story} {stories} {tree} {status} {deep} {missing} {ghost} {extra} {local}</p>
